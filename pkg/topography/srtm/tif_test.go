package srtm

import (
	"testing"

	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestParseTif(t *testing.T) {
	topo := model.Topography{}
	err := ParseTif("test/N26W080.tif", 4, &topo)
	assert.NoError(t, err)
	assert.Equal(t, 17, len(topo.Altitudes))
}
