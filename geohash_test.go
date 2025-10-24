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
	sw, ne := geohash.GeoHash("ud9wn"), geohash.GeoHash("ud9y8")
	seq := geohash.GeoHashesIn(ne, sw)

	it.Then(t).Should(
		it.Seq(seq).Equal("ud9ww", "ud9wx", "ud9y8", "ud9wq", "ud9wr", "ud9y2", "ud9wn", "ud9wp", "ud9y0"),
	)
}

func TestShortestPath(t *testing.T) {
	bag := geohash.GeoHashes{
		"ud9wpwv2kccp",
		"ud9wpr1trmrv",
		"ud9wpxjrf7wx",
		"ud9wpx61yuyw",
		"ud9wprkgjd38",
		"ud9wpxm6jh2t",
		"ud9wprs33e86",
		"ud9wprt60ez9",
		"ud9wnyee164z",
		"ud9wpxs7wbw9",
		"ud9wpxdv8zku",
		"ud9wpxthtvc6",
		"ud9wpxstky14",
		"ud9wnyg1cyny",
		"ud9wpxfgpys7",
		"ud9wpxcun8dn",
		"ud9wprvhhs3q",
	}

	seq, _ := geohash.ShortestPath(bag)
	it.Then(t).Should(
		it.Seq(seq).Equal(
			"ud9wpwv2kccp", "ud9wpxjrf7wx", "ud9wpxm6jh2t", "ud9wpxs7wbw9",
			"ud9wpxstky14", "ud9wpxthtvc6", "ud9wpxdv8zku", "ud9wpxfgpys7",
			"ud9wpxcun8dn", "ud9wpx61yuyw", "ud9wprkgjd38", "ud9wprs33e86",
			"ud9wprt60ez9", "ud9wprvhhs3q", "ud9wpr1trmrv", "ud9wnyee164z",
			"ud9wnyg1cyny",
		),
	)
}

func TestExpand(t *testing.T) {
	hash := geohash.GeoHash("ud9wr3xe")
	radius := 1.0 // 1 km

	ne, sw := geohash.Expand(hash, radius)

	// Test that we get valid geohashes back
	it.Then(t).Should(
		it.True(len(ne) > 0),
		it.True(len(sw) > 0),
		it.Equal(len(ne), len(hash)), // Same precision as input
		it.Equal(len(sw), len(hash)), // Same precision as input
	)

	// Test that the corners are roughly the expected distance away
	distanceNE := geohash.Distance(hash, ne)
	distanceSW := geohash.Distance(hash, sw)

	// Distance should be close to radius (within reasonable tolerance due to approximation)
	it.Then(t).Should(
		it.True(distanceNE > 0.5 && distanceNE < 1.5), // Within 50% tolerance
		it.True(distanceSW > 0.5 && distanceSW < 1.5), // Within 50% tolerance
	)

	// Test with different radius
	ne5, sw5 := geohash.Expand(hash, 5.0)
	distance5NE := geohash.Distance(hash, ne5)
	distance5SW := geohash.Distance(hash, sw5)

	it.Then(t).Should(
		it.True(distance5NE > distanceNE), // Larger radius should give larger distance
		it.True(distance5SW > distanceSW), // Larger radius should give larger distance
	)
}
