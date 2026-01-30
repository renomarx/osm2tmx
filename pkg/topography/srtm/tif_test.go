package srtm

import (
	"testing"

	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestParseTif(t *testing.T) {

	t.Run("single_tif", func(t *testing.T) {
		topo := model.Topography{}
		tp := NewTifParser(&topo)
		tp.AddTif("test/N26W080.tif")
		alt, err := tp.GetAltitude(26.769445, -80.276432, 4)
		assert.NoError(t, err)
		assert.Equal(t, model.Altitude(8), alt)
		assert.Equal(t, 17, len(topo.Altitudes))
	})

	t.Run("directory_without_preload", func(t *testing.T) {
		topo := model.Topography{}
		tp := NewTifParser(&topo)
		tp.AddDirectory("test")

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

	// t.Run("directory_with_preload", func(t *testing.T) {
	// 	topo := model.Topography{}
	// 	tp := NewTifParser(&topo)
	// 	tp.AddDirectory("test")
	// 	tp.Preload()

	// 	alt, err := tp.GetAltitude(26.769445, -80.276432, 4)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, model.Altitude(8), alt)
	// 	assert.Equal(t, 17, len(topo.Altitudes))

	// 	alt, err = tp.GetAltitude(-56.115745, 158.687832, 4)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, model.Altitude(5), alt)
	// 	t.Log(topo.Altitudes)
	// 	assert.Equal(t, 54, len(topo.Altitudes))
	// })
}
