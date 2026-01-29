package srtm

import (
	"bufio"
	"fmt"
	"image/color"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/chai2010/tiff"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type TifParser struct {
	tifs map[string]string
}

func NewTifParser() *TifParser {
	return &TifParser{
		tifs: make(map[string]string),
	}
}

func (tp *TifParser) LoadDirectory(dirpath string, recursive bool) error {
	// TODO: for each .tif in dirpath (recursively),
	// add tif to tp.tifs[basename] = pathname
	return nil
}

func (tp *TifParser) Preload(minlat, maxlat, minlon, maxlon float64) error {
	// TODO: fParse all tifs corresponding to the range between minlat and maxlat, minlon and maxlon
	// return a custom error if missing tifs, but preload the others anyway
	return nil
}

func (tp *TifParser) GetAltitude(lat, lon float64) (model.Altitude, error) {
	// TODO: get altitude from topography already loaded
	// if not found, search for corresponding tif and parse it
	// then return the corresponding altitude from topography
	return model.Altitude(0), nil
}

func ParseTif(filename string, precision int, topo *model.Topography) error {
	input, err := os.Open(filename)
	if err != nil {
		return err
	}

	// Expected format for filename: N26W080.tif
	north, west, err := parseTifFilename(filename)
	if err != nil {
		return err
	}
	latOffset := float64(north)
	lngOffset := float64(180 - west) // West to East

	img, err := tiff.Decode(bufio.NewReader(input))
	if err != nil {
		return err
	}

	if topo.Altitudes == nil {
		topo.Altitudes = make(map[model.GeoPoint]model.Altitude)
	}

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Normalize X and Y to a 0 to 1 space based on the size of the image.
			// Add the offsets to get coordinates.
			lng := lngOffset + (float64(x) / float64(bounds.Max.X))
			lat := latOffset + (float64(y) / float64(bounds.Max.Y))

			height := img.At(x, y).(color.Gray16)
			if height.Y == 0 {
				// optimization: we ignore 0-value height to avoid filling memory with unnecessary data
				continue
			}
			mult := math.Pow(10, float64(precision))
			roundedLat := math.Round(lat*mult) / mult
			roundedLon := math.Round(lng*mult) / mult

			topo.Altitudes[model.GeoPoint{Lat: roundedLat, Lon: roundedLon}] = model.Altitude(height.Y)
		}
	}

	return nil
}

func parseTifFilename(filename string) (int, int, error) {
	errMsg := "bad tif filename: expected NxxWyy, got %s"
	basename := path.Base(filename)
	basename = strings.TrimSuffix(basename, ".tif")
	if len(basename) == 0 {
		return 0, 0, fmt.Errorf(errMsg, basename)
	}
	posN := strings.Index(basename, "N")
	if posN != 0 {
		return 0, 0, fmt.Errorf(errMsg, basename)
	}
	posW := strings.Index(basename, "W")
	if posW == -1 || len(basename) <= posW+1 || posN+1 >= posW {
		return 0, 0, fmt.Errorf(errMsg, basename)
	}
	northStr := basename[posN+1 : posW]
	westStr := basename[posW+1:]
	north, err := strconv.Atoi(northStr)
	if err != nil {
		return 0, 0, fmt.Errorf(errMsg+": %w", basename, err)
	}
	west, err := strconv.Atoi(westStr)
	if err != nil {
		return 0, 0, fmt.Errorf(errMsg+": %w", basename, err)
	}

	return north, west, nil
}
