package tmx

import (
	"encoding/xml"
	"fmt"
	"os"
)

// LoadTMX ouvre et parse un fichier TMX
func LoadTMX(filename string) (*Map, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var m Map
	if err := decoder.Decode(&m); err != nil {
		return nil, fmt.Errorf("parse erreur: %v", err)
	}
	return &m, nil
}
