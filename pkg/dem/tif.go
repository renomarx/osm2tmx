package dem

import (
	"bufio"
	"image/color"
	"os"
	"strconv"

	"github.com/chai2010/tiff"
)

func main() {
	input, err := os.Open("N050E028_DSM.tif")
	if err != nil {
		panic(err)
	}

	// TODO: parse the offsets from the file name
	latOffset := float64(50)
	lngOffset := float64(28)

	output, err := os.Create("N050E028_DSM.xyz")
	if err != nil {
		panic(err)
	}
	outputWriter := bufio.NewWriter(output)

	img, err := tiff.Decode(bufio.NewReader(input))
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Normalize X and Y to a 0 to 1 space based on the size of the image.
			// Add the offsets to get coordinates.
			lng := lngOffset + (float64(x) / float64(bounds.Max.X))
			lat := latOffset + (float64(y) / float64(bounds.Max.Y))

			height := img.At(x, y).(color.Gray16)

			outputWriter.WriteString(strconv.FormatFloat(lng, 'f', 16, 64))
			outputWriter.WriteString(" ")
			outputWriter.WriteString(strconv.FormatFloat(lat, 'f', 16, 64))
			outputWriter.WriteString(" ")
			outputWriter.WriteString(strconv.Itoa(int(height.Y)))
			outputWriter.WriteString("\n")
		}
	}
}
