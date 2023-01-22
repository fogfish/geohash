package geohash

import "github.com/fogfish/curie"

// GeoZone globally unique identity of area/box, which is derived from GeoHash
// GeoZone is GeoHash of reduced precision
//   - g:u
//   - ...
//   - g:ud8q81
type GeoZone curie.IRI

func ToZone(hash GeoHash) GeoZone {
	return GeoZone(curie.New("g:%s", hash))
}

func FromZone(zone GeoZone) GeoHash {
	return GeoHash(curie.Reference(curie.IRI(zone)))
}
