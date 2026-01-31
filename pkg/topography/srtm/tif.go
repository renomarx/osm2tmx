package srtm

import (
	"bufio"
	"fmt"
	"image/color"
	"io/fs"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chai2010/tiff"
	"github.com/renomarx/osm2tmx/pkg/model"
)

var (
	ErrTifNotFound  error = fmt.Errorf("tif not found")
	ErrNotSupported error = fmt.Errorf("not supported lat,lon (not in SRTM range)")
)

type TifParser struct {
	tifs       map[tifID]string
	loaded     map[tifID]bool
	Topography *model.Topography
}

type tifID struct {
	Lat, Lon int
}

func NewTifParser(topography *model.Topography) *TifParser {
	if topography.Altitudes == nil {
		topography.Altitudes = make(map[model.GeoPoint]model.Altitude)
	}
	return &TifParser{
		tifs:       make(map[tifID]string),
		loaded:     make(map[tifID]bool),
		Topography: topography,
	}
}

// HasTifFiles usefull to know if there is any need to use the parser
func (tp *TifParser) HasTifFiles() bool {
	return len(tp.tifs) > 0
}

// AddDirectory: for each .tif in dirpath (recursively),
// add tif filepath to TifParser (for future loading)
func (tp *TifParser) AddDirectory(dirpath string) error {
	return filepath.WalkDir(dirpath, tp.walkDirectory)
}

func (tp *TifParser) walkDirectory(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
		return nil
	}

	// Voluntarly ignore errors because we want to skip unsupported files
	_ = tp.AddTif(s)

	return nil
}

func (tp *TifParser) AddTif(filepath string) error {
	tifID, err := parseTifFilepath(filepath)
	if err != nil {
		return err
	}
	tp.tifs[tifID] = filepath
	return nil
}

// Preload parse all tifs corresponding to the range between minlat and maxlat, minlon and maxlon
// return a custom error if missing tifs, but preload the others anyway
func (tp *TifParser) Preload(minlat, maxlat, minlon, maxlon float64, precision int) error {
	if len(tp.tifs) > int(math.Ceil((maxlat-minlat)*(maxlon-minlon))) {
		// Preloading by lat,lon range
		for lat := minlat; lat <= maxlat; lat += 1 {
			for lon := minlon; lon <= maxlon; lon += 1 {
				id := tifID{Lat: int(lat), Lon: int(lon)}
				filepath, exists := tp.tifs[id]
				if !exists {
					continue
				}
				tp.parseTif(id, filepath, precision)
			}
		}
		return nil
	}

	// Preloading by tifs added
	for tifID, filepath := range tp.tifs {
		// Parse correspoding tif if included withmin min,max range
		if tifID.Lat >= int(minlat) && tifID.Lat <= int(maxlat) &&
			tifID.Lon >= int(minlon) && tifID.Lon <= int(maxlon) {
			tp.parseTif(tifID, filepath, precision)
		}
	}

	return nil
}

func (tp *TifParser) GetAltitude(lat, lon float64, precision int) (model.Altitude, error) {
	// Get altitude from topography if already loaded
	mult := math.Pow(10, float64(precision))
	roundedLat := math.Round(lat*mult) / mult
	roundedLon := math.Round(lon*mult) / mult
	geopoint := model.GeoPoint{Lat: roundedLat, Lon: roundedLon}
	alt, exists := tp.Topography.Altitudes[geopoint]
	if exists {
		return alt, nil
	}

	// search for corresponding tif and parse it
	// then return the corresponding altitude from topography
	id := tifID{Lat: int(lat), Lon: int(lon)}
	filepath, exists := tp.tifs[id]
	if !exists {
		return 0, fmt.Errorf("%w: %+v", ErrTifNotFound, id)
	}

	tp.parseTif(id, filepath, precision)

	alt, exists = tp.Topography.Altitudes[geopoint]
	if exists {
		return alt, nil
	}

	// If still not found, we can definitively consider that we do not have the
	// altitude for this geopoint, and set it to 0
	tp.Topography.Altitudes[geopoint] = 0

	return model.Altitude(0), nil
}

func (tp *TifParser) parseTif(id tifID, filepath string, precision int) error {
	_, exists := tp.loaded[id]
	if exists {
		// Tif already loaded, no need to parse it again
		return nil
	}
	log.Printf("loading SRTM tif file %s with precision %d ...", filepath, precision)
	input, err := os.Open(filepath)
	if err != nil {
		return err
	}

	latOffset := float64(id.Lat)
	lngOffset := float64(id.Lon)

	img, err := tiff.Decode(bufio.NewReader(input))
	if err != nil {
		return err
	}

	var maxAlt uint16
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Normalize X and Y to a 0 to 1 space based on the size of the image.
			// Add the offsets to get coordinates.
			latAbs := math.Abs(latOffset) + (float64(bounds.Max.Y-y) / float64(bounds.Max.Y))
			lngAbs := math.Abs(lngOffset) + (float64(x) / float64(bounds.Max.X))

			height := img.At(x, y).(color.Gray16)
			if height.Y == 0 {
				// optimization: we ignore 0-value height to avoid filling memory with unnecessary data
				continue
			}
			mult := math.Pow(10, float64(precision))
			latFactor := 1
			if latOffset < 0 {
				latFactor = -1
			}
			lngFactor := 1
			if lngOffset < 0 {
				lngFactor = -1
			}
			roundedLat := math.Round(latAbs*float64(latFactor)*mult) / mult
			roundedLon := math.Round(lngAbs*float64(lngFactor)*mult) / mult

			tp.Topography.Altitudes[model.GeoPoint{Lat: roundedLat, Lon: roundedLon}] = model.Altitude(height.Y)
			if height.Y > maxAlt {
				maxAlt = height.Y
			}
		}
	}
	tp.loaded[id] = true
	log.Printf("finished loading SRTM file %s, max alt: %d", filepath, maxAlt)

	return nil
}

func parseTifFilepath(filepath string) (tifID, error) {
	errMsg := "bad tif filename: expected matching '[N|S]dd[E|W]ddd', got %s"
	basename := path.Base(filepath)
	basename = strings.TrimSuffix(basename, ".tif")

	if len(basename) != 7 {
		return tifID{}, fmt.Errorf(errMsg, basename)
	}

	latFactor := 1
	switch basename[0] {
	case 'N':
		latFactor = 1
	case 'S':
		latFactor = -1
	default:
		return tifID{}, fmt.Errorf(errMsg, basename)
	}

	lonFactor := 1
	switch basename[3] {
	case 'E':
		lonFactor = 1
	case 'W':
		lonFactor = -1
	default:
		return tifID{}, fmt.Errorf(errMsg, basename)
	}

	latStr := basename[1:3]
	lonStr := basename[4:]
	latAbs, err := strconv.Atoi(latStr)
	if err != nil {
		return tifID{}, fmt.Errorf(errMsg+": %w", basename, err)
	}
	lonAbs, err := strconv.Atoi(lonStr)
	if err != nil {
		return tifID{}, fmt.Errorf(errMsg+": %w", basename, err)
	}

	return tifID{Lat: latAbs * latFactor, Lon: lonAbs * lonFactor}, nil
}
