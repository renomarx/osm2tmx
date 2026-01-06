package model

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

type Map struct {
	Layers []Layer
}

type Layer struct {
	M [][]*Case // M[y][x]
}

type Case struct {
	Tile Tile
	X    int
	Y    int
}

type Tile int

func (m *Map) Print() {
	for z, l := range m.Layers {
		fmt.Printf("Layer %d -----------------------\n", z)
		for y := range l.M {
			var line strings.Builder
			for _, c := range l.M[y] {
				var tile Tile = 0
				if c != nil {
					tile = c.Tile
				}
				line.WriteString(fmt.Sprintf("%d ", tile))
			}
			fmt.Println(line.String())
		}
	}
}

func (l *Layer) PrintCSV() string {
	var csvStr strings.Builder
	writer := csv.NewWriter(&csvStr)

	for y := range l.M {
		var records []string = make([]string, len(l.M[y]))
		for x, c := range l.M[y] {
			var tile Tile = 0
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

func (l *Layer) PrintCSV2() string {
	var csvStr strings.Builder

	for y := range l.M {
		for _, c := range l.M[y] {
			var tile Tile = 0
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
