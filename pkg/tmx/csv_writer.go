package tmx

import (
	"fmt"
	"strings"

	"github.com/renomarx/osm2tmx/pkg/model"
)

// PrintCSVWithLastComma writes non-standard csv,
// adding a comma at the end of each line except the last line
// because it seems that it's the expected format for Tiled software
func PrintCSVWithLastComma(l *model.Layer) string {
	var csvStr strings.Builder

	for y := range l.M {
		for _, c := range l.M[y] {
			var tile model.Tile = 0
			if c != nil {
				tile = c.Tile
			}
			csvStr.WriteString(fmt.Sprintf("%d,", tile))
		}
		csvStr.WriteString("\n")
	}
	result := csvStr.String()

	// removing last comma
	return result[:len(result)-2] + "\n"
}
