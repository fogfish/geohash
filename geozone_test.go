package geohash_test

import (
	"testing"

	"github.com/fogfish/geohash"
	"github.com/fogfish/it/v2"
)

func TestZone(t *testing.T) {
	for _, hash := range []geohash.GeoHash{
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
			it.Equal(geohash.ToZone(hash), geohash.GeoZone("g:"+hash)),
			it.Equal(geohash.FromZone(geohash.ToZone(hash)), hash),
		)
	}

}
