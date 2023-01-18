package geohash_test

import (
	"testing"

	"github.com/fogfish/geohash"
	"github.com/fogfish/it/v2"
)

func TestNew(t *testing.T) {
	for expect, coords := range map[geohash.GeoHash][]float64{
		"ud9wr3xe47zx": {24.938385, 60.169909},
		"ud9wr3xe47z":  {24.938385, 60.169909},
		"ud9wr3xe47":   {24.938385, 60.169909},
		"ud9wr3xe4":    {24.938385, 60.169909},
		"ud9wr3xe":     {24.93845, 60.16997},
		"ud9wr3x":      {24.9383, 60.1701},
		"ud9wr3":       {24.933, 60.1694},
		"ud9wr":        {24.939, 60.183},
		"ud9w":         {24.79, 60.21},
		"ud9":          {24.6, 59.8},
		"ud":           {28.0, 59.1},
		"u":            {23.0, 68.0},
	} {
		it.Then(t).Should(
			it.Equal(expect, geohash.New(coords[1], coords[0], geohash.Precision(len(expect)))),
			it.String(geohash.New(coords[1], coords[0])).HavePrefix(string(expect)),
		)
	}
}

func TestRound(t *testing.T) {
	for _, expect := range []geohash.GeoHash{
		"ud9wr3xe47zx",
		"ud9wr3xe47z",
		"ud9wr3xe47",
		"ud9wr3xe4",
		"ud9wr3xe",
		"ud9wr3x",
		"ud9wr3",
		"ud9wr",
		"ud9w",
		"ud9",
		"ud",
		"u",
	} {
		it.Then(t).Should(
			it.Equal(expect, geohash.Round("ud9wr3xe47zx", geohash.Precision(len(expect)))),
		)

	}
}

func TestNeighbor(t *testing.T) {
	for input, expect := range map[geohash.GeoHash]geohash.GeoHashes{
		"ud9wr3xe": {"ud9wr3xs", "ud9wr3xd", "ud9wr3x7", "ud9wr3xg"},
		"ud9wr3x":  {"ud9wr3z", "ud9wr3r", "ud9wr3w", "ud9wr98"},
		// TODO: with https://www.movable-type.co.uk/scripts/geohash.html
	} {
		it.Then(t).Should(
			it.Equal(expect[0], geohash.NorthOf(input)),
			it.Equal(expect[1], geohash.SouthOf(input)),
			it.Equal(expect[2], geohash.WestOf(input)),
			it.Equal(expect[3], geohash.EastOf(input)),
		)
	}
}

func TestToLatLng(t *testing.T) {
	// for input, expect := range map[geohash.GeoHash][]float64{
	// 	"ud9wr3xe4": {24.938385, 60.169909},
	// 	"ud9wr3xe":  {24.93845, 60.16997},
	// 	"ud9wr3x":   {24.9383, 60.1701},
	// 	"ud9wr3":    {24.933, 60.1694},
	// 	"ud9wr":     {24.939, 60.183},
	// 	"ud9w":      {24.79, 60.21},
	// 	"ud9":       {24.6, 59.8},
	// 	"ud":        {28.0, 59.1},
	// 	"u":         {23.0, 68.0},
	// } {
	// 	lat, lng := geohash.ToLatLng(input)

	// 	it.Then(t).Should(
	// 		it.Equal(lat, expect[1]),
	// 		it.Equal(lng, expect[0]),
	// 	)
	// }
}

func TestGeoHashesIn(t *testing.T) {
	// d := 0.01
	// ne := geohash.New(60.1699+d, 24.9384+d, 6)
	// sw := geohash.New(60.1699-d, 24.9384-d, 6)

	// seq := geohash.GeoHashesIn(ne, sw)
	// fmt.Println(seq)
	// fmt.Println(geohash.Predecessors(seq))
}
