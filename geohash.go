package geohash

import (
	"math"

	"github.com/cabify/geohash-golang"
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

// GeoHash to south of given one
func SouthOf(hash GeoHash) GeoHash {
	return GeoHash(hasher.CalculateAdjacent(string(hash), "bottom"))
}

// GeoHash to east of given one
func EastOf(hash GeoHash) GeoHash {
	return GeoHash(hasher.CalculateAdjacent(string(hash), "right"))
}

// GeoHash to west of given one
func WestOf(hash GeoHash) GeoHash {
	return GeoHash(hasher.CalculateAdjacent(string(hash), "left"))
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

// Unordered Set of GeoHashes
type GeoHashes []GeoHash

// Calculate all GeoHashes that overlaps with area defined by two points
func GeoHashesIn(hashNorthEast, hashSouthWest GeoHash) GeoHashes {
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

func Weights(hashes GeoHashes, hashNorthEast, hashSouthWest GeoHash) []float64 {
	seq := make([]float64, len(hashes))
	boxN, boxE := ToLatLng(hashNorthEast)
	boxS, boxW := ToLatLng(hashSouthWest)
	areaViewport := areaOf(boxS, boxW, boxN, boxE)

	for i, x := range hashes {
		bb := geohash.Decode(string(x))
		s := math.Max(bb.SouthWest().Lat(), boxS)
		n := math.Min(bb.NorthEast().Lat(), boxN)
		w := math.Max(bb.SouthWest().Lng(), boxW)
		e := math.Min(bb.NorthEast().Lng(), boxE)

		areaVisible := areaOf(s, w, n, e)
		seq[i] = areaVisible / areaViewport
	}

	return seq
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
