package tmx

import (
	"testing"

	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestPrintCSVWithLastComma(t *testing.T) {
	sizeY := 6
	sizeX := 12
	layer := model.Layer{
		M: make([][]*model.Cell, sizeY),
	}
	for y := range sizeY {
		layer.M[y] = make([]*model.Cell, sizeX)
		for x := range sizeX {
			layer.M[y][x] = &model.Cell{
				Tile: model.Tile(x * y % 3),
				X:    x,
				Y:    y,
			}
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
