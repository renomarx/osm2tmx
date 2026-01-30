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

var (
	ErrTifNotFound  error = fmt.Errorf("tif not found")
	ErrNotSupported error = fmt.Errorf("not supported lat,lon (not in SRTM range)")
)

type TifParser struct {
	tifs       map[string]string
	Topography *model.Topography
}

func NewTifParser(topography *model.Topography) *TifParser {
	if topography.Altitudes == nil {
		topography.Altitudes = make(map[model.GeoPoint]model.Altitude)
	}
	return &TifParser{
		tifs:       make(map[string]string),
		Topography: topography,
	}
}

// AddDirectory: for each .tif in dirpath (recursively),
// add tif filepath to TifParser (for future loading)
func (tp *TifParser) AddDirectory(dirpath string, recursive bool) error {
	// TODO
	return nil
}

func (tp *TifParser) AddTif(filepath string) error {
	basename := path.Base(filepath)
	basename = strings.TrimSuffix(basename, ".tif")
	// TODO: add some controls on basename format
	tp.tifs[basename] = filepath
	return nil
}

// Preload parse all tifs corresponding to the range between minlat and maxlat, minlon and maxlon
// return a custom error if missing tifs, but preload the others anyway
func (tp *TifParser) Preload(minlat, maxlat, minlon, maxlon float64) error {
	// TODO
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
	// Altitude not found, search for corresponding tif and parse it
	// then return the corresponding altitude from topography
	basename := ""
	switch {
	case lat > 60:
		return 0, fmt.Errorf("%w: lat %f", ErrNotSupported, lat)
	case lat < -57:
		return 0, fmt.Errorf("%w: lon %f", ErrNotSupported, lat)
	case lat >= 0:
		north := int(math.Floor(lat))
		basename = fmt.Sprintf("N%02d", north)
	case lat < 0:
		south := int(math.Floor(lat))
		basename = fmt.Sprintf("S%02d", south)
	}
	switch {
	case lon >= 0:
		east := int(math.Floor(lon))
		basename += fmt.Sprintf("E%03d", east)
	case lon < 0:
		west := int(math.Floor(-1 * lon))
		basename += fmt.Sprintf("W%03d", west)
	}

	filepath, exists := tp.tifs[basename]
	if !exists {
		return 0, fmt.Errorf("%w: %s", ErrTifNotFound, basename)
	}

	tp.parseTif(filepath, precision)

	alt, exists = tp.Topography.Altitudes[geopoint]
	if exists {
		return alt, nil
	}

	// If still not found, we can definitively consider that we do not have the
	// altitude for this geopoint, and set it to 0
	tp.Topography.Altitudes[geopoint] = 0

	return model.Altitude(0), nil
}

func (tp *TifParser) parseTif(filepath string, precision int) error {
	input, err := os.Open(filepath)
	if err != nil {
		return err
	}

	// Expected format for filename: N26W080.tif
	north, east, err := parseTifFilepath(filepath)
	if err != nil {
		return err
	}
	latOffset := float64(north)
	lngOffset := float64(east)

	img, err := tiff.Decode(bufio.NewReader(input))
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Normalize X and Y to a 0 to 1 space based on the size of the image.
			// Add the offsets to get coordinates.
			latAbs := math.Abs(latOffset) + (float64(y) / float64(bounds.Max.Y))
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
		}
	}

	return nil
}

func parseTifFilepath(filepath string) (int, int, error) {
	errMsg := "bad tif filename: expected matching '[N|S][0-6][0-9][E|W][0-1][0-8][0-9]', got %s"
	basename := path.Base(filepath)
	basename = strings.TrimSuffix(basename, ".tif")

	if len(basename) != 7 {
		return 0, 0, fmt.Errorf(errMsg, basename)
	}

	northFactor := 1
	switch basename[0] {
	case 'N':
		northFactor = 1
	case 'S':
		northFactor = -1
	default:
		return 0, 0, fmt.Errorf(errMsg, basename)
	}

	eastFactor := 1
	switch basename[3] {
	case 'E':
		eastFactor = 1
	case 'W':
		eastFactor = -1
	default:
		return 0, 0, fmt.Errorf(errMsg, basename)
	}

	northStr := basename[1:3]
	eastStr := basename[4:]
	north, err := strconv.Atoi(northStr)
	if err != nil {
		return 0, 0, fmt.Errorf(errMsg+": %w", basename, err)
	}
	east, err := strconv.Atoi(eastStr)
	if err != nil {
		return 0, 0, fmt.Errorf(errMsg+": %w", basename, err)
	}

	return north * northFactor, east * eastFactor, nil
}
