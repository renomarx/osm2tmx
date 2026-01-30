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
		t.Log(topo.Altitudes)
		assert.Equal(t, 17, len(topo.Altitudes))
	})
}
