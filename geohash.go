package geohash

import (
	"math"

	hasher "github.com/cabify/geohash-golang"
)

// Precision of GeoHash
type Precision int

// GeoHash representation
type GeoHash string

// Create new instance of GeoHash from (lat,lng) pair.
// Optionally, control precision
func New(lat, lng float64, precision ...Precision) GeoHash {
	if len(precision) == 0 {
		return GeoHash(hasher.Encode(lat, lng))
	}

	p := precision[0]
	return GeoHash(hasher.EncodeWithPrecision(lat, lng, int(p)))
}

// Round GeoHash to precision
func Round(hash GeoHash, precision Precision) GeoHash {
	if len(hash) < int(precision) {
		return hash
	}

	return hash[:precision]
}

// GeoHash to north of given one
func NorthOf(hash GeoHash) GeoHash {
	return GeoHash(hasher.CalculateAdjacent(string(hash), "top"))
}

func ToNorth(hash GeoHash, n int) (north GeoHash) {
	north = hash
	for range n {
		north = NorthOf(north)
	}
	return
}

// GeoHash to south of given one
func SouthOf(hash GeoHash) GeoHash {
	return GeoHash(hasher.CalculateAdjacent(string(hash), "bottom"))
}

func ToSouth(hash GeoHash, n int) (south GeoHash) {
	south = hash
	for range n {
		south = SouthOf(south)
	}
	return
}

// GeoHash to east of given one
func EastOf(hash GeoHash) GeoHash {
	return GeoHash(hasher.CalculateAdjacent(string(hash), "right"))
}

func ToEast(hash GeoHash, n int) (east GeoHash) {
	east = hash
	for range n {
		east = EastOf(east)
	}
	return
}

// GeoHash to west of given one
func WestOf(hash GeoHash) GeoHash {
	return GeoHash(hasher.CalculateAdjacent(string(hash), "left"))
}

func ToWest(hash GeoHash, n int) (west GeoHash) {
	west = hash
	for range n {
		west = WestOf(west)
	}
	return
}

// Converts GeoHash to (lat,lng) pair
func ToLatLng(hash GeoHash) (lat float64, lng float64) {
	coords := hasher.Decode(string(hash)).Center()
	lat, lng = coords.Lat(), coords.Lng()
	return
}

// Converts GeoHash to bounding box
func ToBBox(hash GeoHash) (n float64, e float64, s float64, w float64) {
	bbox := hasher.Decode(string(hash))
	ne := bbox.NorthEast()
	sw := bbox.SouthWest()
	n, e, s, w = ne.Lat(), ne.Lng(), sw.Lat(), sw.Lng()
	return
}

func Expand(hash GeoHash, radius float64) (hashNorthEast, hashSouthWest GeoHash) {
	// Get center coordinates of the input geohash
	centerLat, centerLng := ToLatLng(hash)

	// Earth's radius in kilometers
	earthRadius := 6378.0

	// Calculate the angular distance in radians
	angularDistance := radius / earthRadius

	// Convert to degrees
	degToRad := math.Pi / 180.0
	radToDeg := 180.0 / math.Pi

	centerLatRad := centerLat * degToRad

	// Calculate latitude offset
	latOffset := angularDistance * radToDeg

	// Calculate longitude offset (accounting for latitude)
	lngOffset := angularDistance * radToDeg / math.Cos(centerLatRad)

	// Calculate corner coordinates
	northLat := centerLat + latOffset
	southLat := centerLat - latOffset
	eastLng := centerLng + lngOffset
	westLng := centerLng - lngOffset

	// Clamp latitude to valid range [-90, 90]
	if northLat > 90 {
		northLat = 90
	}
	if southLat < -90 {
		southLat = -90
	}

	// Handle longitude wrapping around [-180, 180]
	if eastLng > 180 {
		eastLng = eastLng - 360
	}
	if westLng < -180 {
		westLng = westLng + 360
	}

	// Create geohashes with the same precision as input
	precision := Precision(len(hash))
	hashNorthEast = New(northLat, eastLng, precision)
	hashSouthWest = New(southLat, westLng, precision)

	return
}

// Unordered Set of GeoHashes
type GeoHashes []GeoHash

// Calculate all GeoHashes that overlaps with area defined by two points
func GeoHashesIn(hashNorthEast, hashSouthWest GeoHash) GeoHashes {
	// TODO: simple trunk do not work (re-coding is needed)
	p := len(hashNorthEast)
	if len(hashSouthWest) < p {
		p = len(hashSouthWest)
	}

	n, e := ToLatLng(hashNorthEast)
	s, w := ToLatLng(hashSouthWest)

	hashNorthWest := New(n, w, Precision(p))
	hashSouthEast := New(s, e, Precision(p))

	// first and last columns of the box
	head := hashesNorthToSouth(hashNorthWest, hashSouthWest)
	tail := hashesNorthToSouth(hashNorthEast, hashSouthEast)

	// fill the box
	area := []GeoHash{}
	for i := 0; i < len(head); i++ {
		area = append(area, hashesWestToEast(head[i], tail[i])...)
	}

	return area
}

func hashesNorthToSouth(n, s GeoHash) GeoHashes {
	seq := GeoHashes{n}
	for seq[len(seq)-1] != s {
		seq = append(seq, SouthOf(seq[len(seq)-1]))
	}

	return seq
}

func hashesWestToEast(w, e GeoHash) GeoHashes {
	seq := []GeoHash{w}
	for seq[len(seq)-1] != e {
		seq = append(seq, EastOf(seq[len(seq)-1]))
	}

	return seq
}

// Extends existing GeoHash by col and row offsets
func Neighbor(hash GeoHash, x, y int) GeoHashes {
	north := ToNorth(hash, y)
	south := ToSouth(hash, y)
	northEast := ToEast(north, x)
	northWest := ToWest(north, x)
	southEast := ToEast(south, x)
	southWest := ToWest(south, x)

	// first and last columns of the box
	head := hashesNorthToSouth(northWest, southWest)
	tail := hashesNorthToSouth(northEast, southEast)

	// fill the box
	area := []GeoHash{}
	for i := 0; i < len(head); i++ {
		area = append(area, hashesWestToEast(head[i], tail[i])...)
	}

	return area
}

// Calculates all predecessors (higher rank) GeoHashes from the set.
func Predecessors(seq GeoHashes) []GeoHashes {
	n := len(seq[0]) - 1
	layers := make([]GeoHashes, n)

	for rank := 0; rank < n; rank++ {
		index := map[GeoHash]struct{}{}
		layer := GeoHashes{}

		for i := 0; i < len(seq); i++ {
			hash := seq[i][:rank+1]
			if _, has := index[hash]; !has {
				index[hash] = struct{}{}
				layer = append(layer, hash)
			}
		}

		layers[rank] = layer
	}

	return layers
}

// Calculate "weights" of GeoHashes fitted into the area
func Weights(hashes GeoHashes, hashNorthEast, hashSouthWest GeoHash) []float64 {
	boxN, boxE := ToLatLng(hashNorthEast)
	boxS, boxW := ToLatLng(hashSouthWest)
	if hashNorthEast == hashSouthWest {
		boxN, boxE, boxS, boxW = ToBBox(hashNorthEast)
	}
	areaViewport := areaOf(boxS, boxW, boxN, boxE)

	seq := make([]float64, len(hashes))
	for i, x := range hashes {
		bb := hasher.Decode(string(x))
		s := math.Max(bb.SouthWest().Lat(), boxS)
		n := math.Min(bb.NorthEast().Lat(), boxN)
		w := math.Max(bb.SouthWest().Lng(), boxW)
		e := math.Min(bb.NorthEast().Lng(), boxE)

		areaVisible := areaOf(s, w, n, e)
		seq[i] = areaVisible / areaViewport
	}

	return seq
}

// Approximates distance between GeoHashes
func Distance(a, b GeoHash) float64 {
	lat1, lng1 := ToLatLng(a)
	lat2, lng2 := ToLatLng(b)
	return distance(lat1, lng1, lat2, lng2)
}

func distance(lat1, lng1, lat2, lng2 float64) float64 {
	r := 6378.0
	degToRad := math.Pi / 180.0

	phi1 := (90.0 - lat1) * degToRad
	phi2 := (90.0 - lat2) * degToRad

	theta1 := lng1 * degToRad
	theta2 := lng2 * degToRad

	a := math.Sin(phi1) * math.Sin(phi2) * math.Cos(theta1-theta2)
	b := math.Cos(phi1) * math.Cos(phi2)
	arc := math.Acos(a + b)

	return arc * r
}

func areaOf(latSouth, lngWest, latNorth, lngEast float64) float64 {
	a := distance(latSouth, lngWest, latNorth, lngWest)
	b := distance(latSouth, lngWest, latSouth, lngEast)
	return a * b
}

// Build shortest path from sequence of hashes
func ShortestPath(hashes GeoHashes) (GeoHashes, []float64) {
	hashSet := map[GeoHash][]float64{}
	for _, hash := range hashes {
		lat, lng := ToLatLng(hash)
		hashSet[hash] = []float64{lng, lat}
	}

	seq := make(GeoHashes, len(hashes))
	hop := make([]float64, len(hashes))
	seq[0] = hashes[0]

	src := hashSet[hashes[0]]
	delete(hashSet, hashes[0])

	for i := 1; i < len(hashes); i++ {
		hop[i-1] = 637800.0
		for hash, dst := range hashSet {
			d := distance(src[1], src[0], dst[1], dst[0]) * 1000
			if d < hop[i-1] {
				hop[i-1] = d
				seq[i] = hash
			}
		}

		src = hashSet[seq[i]]
		delete(hashSet, seq[i])
	}

	return seq, hop
}
