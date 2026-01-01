package model

import (
	"fmt"
	"strings"
)

type Map struct {
	Layers []Layer
}

type Layer struct {
	M [][]Case // M[y][x]
}

type Case struct {
	Tile Tile
}

type Tile int

func (m *Map) Print() {
	for z, l := range m.Layers {
		fmt.Printf("Layer %d -----------------------\n", z)
		for y := range l.M {
			var line strings.Builder
			for _, c := range l.M[y] {
				line.WriteString(fmt.Sprintf("%d ", c.Tile))
			}
			fmt.Println(line.String())
		}
	}
}
