package raster

import (
	"math"
	"testing"
	"time"

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

		begin := time.Now()
		result, err := raster.Parse(osmfilename)
		totalDuration := time.Since(begin)
		t.Logf("duration: %d ms", totalDuration/time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result.Map.Layers))
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
			MapSizeX:         410,
			MapSizeY:         352,
			MaxEasting:       742936.82,
			MaxNorthing:      5806340.56,
			MinEasting:       742527.16,
			MinNorthing:      5805989.38,
			Nodes:            1694,
			Ways:             306,
			Relations:        26,
			NodesOutOfBounds: 3322,
		}, result.Meta)

		lat, lon := raster.toLatLon(92, 42)
		assert.Equal(t, 46.1617, math.Round(lat*10000)/10000) // TODO: fix
		assert.Equal(t, 6.6711, math.Round(lon*10000)/10000)  // TODO: fix
	})

	t.Run("downscale_4", func(t *testing.T) {
		mapper := mapper.New()

		raster := New(mapper, 4, Bounds{})

		result, err := raster.Parse(osmfilename)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result.Map.Layers))
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
			MapSizeX:         102,
			MapSizeY:         88,
			MaxEasting:       742936.82,
			MaxNorthing:      5806340.56,
			MinEasting:       742527.16,
			MinNorthing:      5805989.38,
			Nodes:            1698,
			Ways:             306,
			Relations:        26,
			NodesOutOfBounds: 3318,
		}, result.Meta)

		lat, lon := raster.toLatLon(92, 42)
		assert.Equal(t, 46.1609, math.Round(lat*10000)/10000) // TODO: fix
		assert.Equal(t, 6.6735, math.Round(lon*10000)/10000)  // TODO: fix
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
		assert.Equal(t, 2, len(result.Map.Layers))
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
			MapSizeX:         100,
			MapSizeY:         95,
			MaxEasting:       742936.82,
			MaxNorthing:      5806340.56,
			MinEasting:       742527.16,
			MinNorthing:      5805989.38,
			Nodes:            697,
			Ways:             306,
			Relations:        26,
			NodesOutOfBounds: 4319,
		}, result.Meta)
	})
}
