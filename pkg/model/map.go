package model

import (
	"fmt"
	"strings"
)

type Map struct {
	Layers []Layer
}

type Layer struct {
	m [][]*Cell // M[y][x]
}

type Cell struct {
	Tile Tile
	X    int
	Y    int
}

type Tile int

func (m *Map) Init(layers, mapSizeX, mapSizeY int, defaultTile Tile) {
	m.Layers = make([]Layer, layers)
	if layers == 0 {
		return
	}
	// Important: we only set default tile for the layer 0, keeping empty other layers,
	// to avoid overloading all future other layer[0] tiles
	m.Layers[0].Init(mapSizeX, mapSizeY, defaultTile)
	for z := 1; z < layers; z++ {
		m.Layers[z] = Layer{}
		m.Layers[z].Init(mapSizeX, mapSizeY, 0)
	}
}

func (l *Layer) Init(mapSizeX, mapSizeY int, tile Tile) {
	l.m = make([][]*Cell, mapSizeY)
	for y := range l.m {
		l.m[y] = make([]*Cell, mapSizeX)
		for x := range l.m[y] {
			l.m[y][x] = &Cell{Tile: tile, X: x, Y: y}
		}
	}
}

func (l *Layer) SizeY() int {
	return len(l.m)
}

func (l *Layer) SizeX() int {
	if len(l.m) == 0 {
		return 0
	}
	return len(l.m[0])
}

func (l *Layer) Row(y int) []*Cell {
	return l.m[y]
}

func (l *Layer) GetCell(x, y int) *Cell {
	return l.m[y][x]
}

func (l *Layer) SetTile(x, y int, tile Tile) {
	if l.m[y][x] == nil {
		l.m[y][x] = &Cell{
			X: x,
			Y: y,
		}
	}
	l.m[y][x].Tile = tile
}

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
	for y := range l.m {
		for _, c := range l.m[y] {
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
