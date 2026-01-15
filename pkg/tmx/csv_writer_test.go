package tmx

import (
	"testing"

	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestPrintCSVWithLastComma(t *testing.T) {
	sizeY := 6
	sizeX := 12
	layer := model.Layer{}
	layer.Init(sizeX, sizeY, 0)
	for y := range sizeY {
		for x := range sizeX {
			layer.SetTile(x, y, model.Tile(x*y%3))
		}
	}

	csv := PrintCSVWithLastComma(&layer)
	expectedCSV := `0,0,0,0,0,0,0,0,0,0,0,0,
0,1,2,0,1,2,0,1,2,0,1,2,
0,2,1,0,2,1,0,2,1,0,2,1,
0,0,0,0,0,0,0,0,0,0,0,0,
0,1,2,0,1,2,0,1,2,0,1,2,
0,2,1,0,2,1,0,2,1,0,2,1
`
	assert.Equal(t, expectedCSV, csv)
}
