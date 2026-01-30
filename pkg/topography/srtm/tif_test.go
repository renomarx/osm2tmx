package srtm

import (
	"testing"

	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestTifParser(t *testing.T) {

	t.Run("single_tif", func(t *testing.T) {
		topo := model.Topography{}
		tp := NewTifParser(&topo)
		err := tp.AddTif("test/N26W080.tif")
		assert.NoError(t, err)
		alt, err := tp.GetAltitude(26.769445, -80.276432, 4)
		assert.NoError(t, err)
		assert.Equal(t, model.Altitude(8), alt)
		assert.Equal(t, 17, len(topo.Altitudes))
	})

	t.Run("err_add_bad_tif", func(t *testing.T) {
		topo := model.Topography{}
		tp := NewTifParser(&topo)
		err := tp.AddTif("test/false_file.tif")
		assert.Error(t, err)
	})

	t.Run("directory_without_preload", func(t *testing.T) {
		topo := model.Topography{}
		tp := NewTifParser(&topo)
		err := tp.AddDirectory("test")
		assert.NoError(t, err)

		alt, err := tp.GetAltitude(26.769445, -80.276432, 4)
		assert.NoError(t, err)
		assert.Equal(t, model.Altitude(8), alt)
		assert.Equal(t, 17, len(topo.Altitudes))

		alt, err = tp.GetAltitude(-56.115745, 158.687832, 4)
		assert.NoError(t, err)
		assert.Equal(t, model.Altitude(5), alt)
		t.Log(topo.Altitudes)
		assert.Equal(t, 54, len(topo.Altitudes))
	})

	t.Run("directory_with_preload_ranging_over_tifs", func(t *testing.T) {
		topo := model.Topography{}
		tp := NewTifParser(&topo)
		err := tp.AddDirectory("test")
		assert.NoError(t, err)

		err = tp.Preload(-57, 27, -80, 158, 4)
		assert.NoError(t, err)
		assert.Equal(t, 54, len(topo.Altitudes))

		alt, err := tp.GetAltitude(26.769445, -80.276432, 4)
		assert.NoError(t, err)
		assert.Equal(t, model.Altitude(8), alt)

		alt, err = tp.GetAltitude(-56.115745, 158.687832, 4)
		assert.NoError(t, err)
		assert.Equal(t, model.Altitude(5), alt)
	})

	t.Run("directory_with_preload_ranging_over_lat_lon", func(t *testing.T) {
		topo := model.Topography{}
		tp := NewTifParser(&topo)
		err := tp.AddDirectory("test")
		assert.NoError(t, err)

		err = tp.Preload(-56, -56, 158, 158, 4)
		assert.NoError(t, err)
		assert.Equal(t, 37, len(topo.Altitudes))

		alt, err := tp.GetAltitude(-56.115745, 158.687832, 4)
		assert.NoError(t, err)
		assert.Equal(t, model.Altitude(5), alt)
	})
}
