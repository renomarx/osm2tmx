package model

import (
	"fmt"
	"strings"
)

type Map struct {
	Layers []Layer
}

type Layer struct {
	m [][]Cell // M[y][x]
}

type Cell struct {
	Tile Tile
	X    int
	Y    int
}

type Tile int

func (m *Map) Init(layers, mapSizeX, mapSizeY int, getTile func(x, y int) Tile) {
	m.Layers = make([]Layer, layers)
	if layers == 0 {
		return
	}
	// Important: we only set default tile for the layer 0, keeping empty other layers,
	// to avoid overloading all future other layer[0] tiles
	m.Layers[0].Init(mapSizeX, mapSizeY, getTile)
	for z := 1; z < layers; z++ {
		m.Layers[z] = Layer{}
		m.Layers[z].Init(mapSizeX, mapSizeY, func(x, y int) Tile { return 0 })
	}
}

func (l *Layer) Init(mapSizeX, mapSizeY int, getTile func(x, y int) Tile) {
	l.m = make([][]Cell, mapSizeY)
	for y := range l.m {
		l.m[y] = make([]Cell, mapSizeX)
		for x := range l.m[y] {
			l.m[y][x] = Cell{Tile: getTile(x, y), X: x, Y: y}
		}
	}
}

func (m *Map) SizeY() int {
	if len(m.Layers) == 0 {
		return 0
	}
	return m.Layers[0].SizeY()
}

func (m *Map) SizeX() int {
	if len(m.Layers) == 0 {
		return 0
	}
	return m.Layers[0].SizeX()
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

func (l *Layer) Row(y int) []Cell {
	return l.m[y]
}

func (l *Layer) GetCell(x, y int) Cell {
	if y < 0 || y >= len(l.m) {
		return Cell{X: x, Y: y}
	}
	if x < 0 || x >= len(l.m[y]) {
		return Cell{X: x, Y: y}
	}
	return l.m[y][x]
}

func (l *Layer) SetTile(x, y int, tile Tile) {
	if y < 0 || y >= len(l.m) {
		return
	}
	if x < 0 || x >= len(l.m[y]) {
		return
	}
	l.m[y][x] = Cell{
		X:    x,
		Y:    y,
		Tile: tile,
	}
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
			line.WriteString(fmt.Sprintf("%d,", c.Tile))
		}
		line.WriteString("\n")
	}
	return line.String()
}
