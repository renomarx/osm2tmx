package tmx

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTMXWriter(t *testing.T) {

	writer := NewWriter("tileset/basechip_pipo.tsx", 16, 16)

	m := model.Map{
		Layers: []model.Layer{
			*generateLayerTest(t),
			*generateLayerTest(t),
		},
	}
	rasterResult := model.RasterMap{
		Map: &m,
		Meta: model.RasterMapMeta{
			MapSizeX: 12,
			MapSizeY: 6,
		},
	}

	f, err := os.CreateTemp("", fmt.Sprintf("osm2tmx_writer_test_%s.osm.tmx", uuid.New()))
	require.NoError(t, err)
	defer os.Remove(f.Name())

	writer.Write(rasterResult, f.Name())

	bytes, err := os.ReadFile(f.Name())
	require.NoError(t, err)

	assert.Equal(t, expectedXML, string(bytes))
}

func generateLayerTest(t *testing.T) *model.Layer {
	sizeY := 6
	sizeX := 12
	layer := model.Layer{}
	layer.Init(sizeX, sizeY, 0)
	for y := range sizeY {
		for x := range sizeX {
			tile := model.Tile(0)
			if (x == 1 || x == sizeX-2) && y != 0 && y != sizeY-1 {
				tile = 2
			}
			if (y == 1 || y == sizeY-2) && x != 0 && x != sizeX-1 {
				tile = 2
			}
			layer.SetTile(x, y, tile)
		}
	}
	layerVue := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	require.Equal(t, layerVue, layer.String())
	return &layer
}

const expectedXML = `<map version="1.4" tiledversion="1.4.3" orientation="orthogonal" renderorder="right-down" width="12" height="6" tilewidth="16" tileheight="16">
  <tileset firstgid="1" source="tileset/basechip_pipo.tsx"></tileset>
  <layer id="1" name="Calque 1" width="12" height="6">
    <data encoding="csv">0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0
</data>
  </layer>
  <layer id="2" name="Calque 2" width="12" height="6">
    <data encoding="csv">0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0
</data>
  </layer>
</map>`
