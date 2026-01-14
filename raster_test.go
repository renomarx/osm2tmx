package main

import (
	"testing"

	"github.com/paulmach/osm"
	"github.com/stretchr/testify/assert"
)

func TestRaster(t *testing.T) {

	osmfilename := "test.osm.pbf"

	mapper := NewMapper()

	raster := NewRaster(mapper)

	result, err := raster.Parse(osmfilename)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Map.Layers))
	assert.Equal(t, 352, result.Map.Layers[0].SizeY())
	assert.Equal(t, 410, result.Map.Layers[0].SizeX())
	assert.Equal(t, 352, result.Map.Layers[1].SizeY())
	assert.Equal(t, 410, result.Map.Layers[1].SizeX())
	assert.Equal(t, RasterResultMeta{
		Bounds: osm.Bounds{
			MinLat: 46.159768999,
			MaxLat: 46.161954,
			MinLon: 6.670234999000001,
			MaxLon: 6.673915,
		},
		MapSizeX:    410,
		MapSizeY:    352,
		MaxEasting:  742936.82,
		MaxNorthing: 5806340.56,
		MinEasting:  742527.16,
		MinNorthing: 5805989.38,
	}, result.Meta)
	assert.Equal(t, 1694, len(result.Nodes))
	assert.Equal(t, 306, len(result.Ways))
	assert.Equal(t, 26, len(result.Relations))
	assert.Equal(t, 3322, len(result.NodesOutOfBounds))
}
