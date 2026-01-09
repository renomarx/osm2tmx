package tmx

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/renomarx/osm2tmx/pkg/model"
)

func PrintCSV(l *model.Layer) string {
	var csvStr strings.Builder
	writer := csv.NewWriter(&csvStr)

	for y := range l.M {
		var records []string = make([]string, len(l.M[y]))
		for x, c := range l.M[y] {
			var tile model.Tile = 0
			if c != nil {
				tile = c.Tile
			}
			records[x] = strconv.Itoa(int(tile))
		}
		err := writer.Write(records)
		if err != nil {
			panic(err)
		}
	}

	return csvStr.String()
}

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
