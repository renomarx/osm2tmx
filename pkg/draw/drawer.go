package draw

import (
	"fmt"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

type UI struct {
	background *ebiten.Image
}

func (ui *UI) loadTmx(tmxFile string) error {
	// Parse .tmx file.
	gameMap, err := tiled.LoadFile(tmxFile)
	if err != nil {
		return fmt.Errorf("error parsing map: %w", err)
	}

	// You can also render the map to an in-memory image for direct
	// use with the default Renderer, or by making your own.
	renderer, err := render.NewRenderer(gameMap)
	if err != nil {
		return fmt.Errorf("map unsupported for rendering: %w", err)
	}

	// Render all layers to the Renderer.
	for i := range len(gameMap.Layers) {
		err = renderer.RenderLayer(i)
		if err != nil {
			return fmt.Errorf("layer %d unsupported for rendering: %w", i, err)
		}
	}

	// Get a reference to the Renderer's output, an image.NRGBA struct.
	img := renderer.Result
	ui.background = ebiten.NewImageFromImage(img)

	return nil
}

func (ui *UI) Update() error {
	return nil
}

func (ui *UI) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-200, -200)
	//op.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(ui.background, op)
}

func (ui *UI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func Draw(tmxFilename string) error {
	ebiten.SetWindowSize(800, 600)
	basename := path.Base(tmxFilename)
	ebiten.SetWindowTitle(basename)

	ui := UI{}
	if err := ui.loadTmx(tmxFilename); err != nil {
		return err
	}

	if err := ebiten.RunGame(&ui); err != nil {
		return err
	}

	return nil
}
