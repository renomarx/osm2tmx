package model

import (
	"fmt"
	"strings"
)

type Map struct {
	Layers []Layer
}

type Layer struct {
	M [][]*Cell // M[y][x]
}

func (l *Layer) SetTile(x, y int, tile Tile) {
	if l.M[y][x] == nil {
		l.M[y][x] = &Cell{
			X: x,
			Y: y,
		}
	}
	l.M[y][x].Tile = tile
}

type Cell struct {
	Tile Tile
	X    int
	Y    int
}

type Tile int

func (m *Map) Print() {
	for z, l := range m.Layers {
		fmt.Printf("Layer %d -----------------------\n", z)
		l.Print()
	}
}

func (l *Layer) Print() {
	fmt.Print(l.String())
}

func (l *Layer) String() string {
	var line strings.Builder
	line.WriteString("\n")
	for y := range l.M {
		for _, c := range l.M[y] {
			var tile Tile = 0
			if c != nil {
				tile = c.Tile
			}
			line.WriteString(fmt.Sprintf("%d,", tile))
		}
		line.WriteString("\n")
	}
	return line.String()
}
