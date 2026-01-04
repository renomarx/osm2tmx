package tmx

import (
	"encoding/xml"
)

// ----------------------------------------------------------------
// Structures qui mappent la structure TMX

type Map struct {
	XMLName     xml.Name  `xml:"map"`
	Version     string    `xml:"version,attr"`
	TiledVer    string    `xml:"tiledversion,attr,omitempty"`
	Orientation string    `xml:"orientation,attr"`
	RenderOrder string    `xml:"renderorder,attr,omitempty"`
	Width       int       `xml:"width,attr"`
	Height      int       `xml:"height,attr"`
	TileWidth   int       `xml:"tilewidth,attr"`
	TileHeight  int       `xml:"tileheight,attr"`
	Tilesets    []Tileset `xml:"tileset"`
	Layers      []Layer   `xml:"layer"`
}

type Tileset struct {
	FirstGID int    `xml:"firstgid,attr"`
	Source   string `xml:"source,attr,omitempty"` // TSX externe
}

type Layer struct {
	ID     int    `xml:"id,attr"`
	Name   string `xml:"name,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Data   Data   `xml:"data"`
}

type Data struct {
	Encoding    string `xml:"encoding,attr,omitempty"`
	Compression string `xml:"compression,attr,omitempty"`
	CSV         string `xml:",innerxml"`
}
