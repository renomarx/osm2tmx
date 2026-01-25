package raster

import (
	"testing"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestRaster(t *testing.T) {
	osmfilename := "test.osm.pbf"

	t.Run("original_size", func(t *testing.T) {
		mapper := mapper.New()

		raster := New(mapper, 1, Bounds{})

		result, err := raster.Parse(osmfilename)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(result.Map.Layers))
		assert.Equal(t, 352, result.Map.Layers[0].SizeY())
		assert.Equal(t, 410, result.Map.Layers[0].SizeX())
		assert.Equal(t, 352, result.Map.Layers[1].SizeY())
		assert.Equal(t, 410, result.Map.Layers[1].SizeX())
		assert.Equal(t, model.RasterMapMeta{
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
	})

	t.Run("downscale_4", func(t *testing.T) {
		mapper := mapper.New()

		raster := New(mapper, 4, Bounds{})

		result, err := raster.Parse(osmfilename)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(result.Map.Layers))
		assert.Equal(t, 88, result.Map.Layers[0].SizeY())
		assert.Equal(t, 102, result.Map.Layers[0].SizeX())
		assert.Equal(t, 88, result.Map.Layers[1].SizeY())
		assert.Equal(t, 102, result.Map.Layers[1].SizeX())
		assert.Equal(t, model.RasterMapMeta{
			Bounds: osm.Bounds{
				MinLat: 46.159768999,
				MaxLat: 46.161954,
				MinLon: 6.670234999000001,
				MaxLon: 6.673915,
			},
			MapSizeX:    102,
			MapSizeY:    88,
			MaxEasting:  742936.82,
			MaxNorthing: 5806340.56,
			MinEasting:  742527.16,
			MinNorthing: 5805989.38,
		}, result.Meta)
		assert.Equal(t, 1698, len(result.Nodes))
		assert.Equal(t, 306, len(result.Ways))
		assert.Equal(t, 26, len(result.Relations))
		assert.Equal(t, 3318, len(result.NodesOutOfBounds))
	})

	t.Run("downscale_2_with_bounds", func(t *testing.T) {
		mapper := mapper.New()

		raster := New(mapper, 2, Bounds{
			OffsetX: 50,
			OffsetY: 40,
			LimitX:  100,
			LimitY:  95,
		})

		result, err := raster.Parse(osmfilename)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(result.Map.Layers))
		assert.Equal(t, 95, result.Map.Layers[0].SizeY())
		assert.Equal(t, 100, result.Map.Layers[0].SizeX())
		assert.Equal(t, 95, result.Map.Layers[1].SizeY())
		assert.Equal(t, 100, result.Map.Layers[1].SizeX())
		assert.Equal(t, model.RasterMapMeta{
			Bounds: osm.Bounds{
				MinLat: 46.159768999,
				MaxLat: 46.161954,
				MinLon: 6.670234999000001,
				MaxLon: 6.673915,
			},
			MapSizeX:    100,
			MapSizeY:    95,
			MaxEasting:  742936.82,
			MaxNorthing: 5806340.56,
			MinEasting:  742527.16,
			MinNorthing: 5805989.38,
		}, result.Meta)
		assert.Equal(t, 697, len(result.Nodes))
		assert.Equal(t, 306, len(result.Ways))
		assert.Equal(t, 26, len(result.Relations))
		assert.Equal(t, 4319, len(result.NodesOutOfBounds))
	})
}
