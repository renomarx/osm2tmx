package draw

import (
	"fmt"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

type UI struct {
	background  *ebiten.Image
	pressedKeys []ebiten.Key
	tx, ty      float64
	zoom        float64
	speed       float64
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
	ui.pressedKeys = inpututil.AppendPressedKeys(ui.pressedKeys[:0])
	return nil
}

type action string

const (
	zoomUp   action = "zoomUp"
	zoomDown action = "zoomDown"
	goLeft   action = "goLeft"
	goDown   action = "goDown"
	goUp     action = "goUp"
	goRight  action = "goRight"
)

func (ui *UI) getAction() action {
	if len(ui.pressedKeys) == 0 {
		return ""
	}
	isCtrlPressed := false
	isKeyUpPressed := false
	isKeyLeftPressed := false
	isKeyDownPressed := false
	isKeyRightPressed := false
	for _, key := range ui.pressedKeys {
		switch key {
		case ebiten.KeyControl:
			isCtrlPressed = true
		case ebiten.KeyArrowUp:
			isKeyUpPressed = true
		case ebiten.KeyArrowLeft:
			isKeyLeftPressed = true
		case ebiten.KeyArrowDown:
			isKeyDownPressed = true
		case ebiten.KeyArrowRight:
			isKeyRightPressed = true
		}
	}
	switch {
	case isCtrlPressed && isKeyUpPressed:
		return zoomUp
	case isCtrlPressed && isKeyDownPressed:
		return zoomDown
	case isKeyUpPressed:
		return goUp
	case isKeyLeftPressed:
		return goLeft
	case isKeyDownPressed:
		return goDown
	case isKeyRightPressed:
		return goRight
	}

	return ""
}

func (ui *UI) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	switch ui.getAction() {
	case zoomDown:
		ui.zoom = ui.zoom / 1.1
	case zoomUp:
		ui.zoom = ui.zoom * 1.1
	case goLeft:
		ui.tx += ui.speed / ui.zoom
	case goRight:
		ui.tx -= ui.speed / ui.zoom
	case goUp:
		ui.ty += ui.speed / ui.zoom
	case goDown:
		ui.ty -= ui.speed / ui.zoom
	}
	op.GeoM.Translate(ui.tx, ui.ty)
	op.GeoM.Scale(ui.zoom, ui.zoom)
	screen.DrawImage(ui.background, op)
}

func (ui *UI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func Draw(tmxFilename string) error {
	ebiten.SetWindowSize(800, 600)
	basename := path.Base(tmxFilename)
	ebiten.SetWindowTitle(basename)

	ui := UI{
		speed: 2,
		zoom:  1,
		tx:    -200,
		ty:    -200,
	}
	if err := ui.loadTmx(tmxFilename); err != nil {
		return err
	}

	if err := ebiten.RunGame(&ui); err != nil {
		return err
	}

	return nil
}
