package tmx

import (
	"encoding/xml"
	"os"
)

// SaveTMX écrit une instance Map au format TMX XML
func SaveTMX(filename string, m *Map) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(m); err != nil {
		return err
	}
	return encoder.Flush()
}
