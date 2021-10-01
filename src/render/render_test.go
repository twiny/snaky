package render

import (
	"snaky/src/snaky"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	assert := assert.New(t)

	render, err := NewRender()
	if err != nil {
		assert.Error(err)
		return
	}

	mockBoard := &snaky.Board{
		X:      29,
		Y:      39,
		Score:  100,
		Round:  100,
		Speed:  snaky.Slow,
		Length: 10,
		Food:   snaky.NewCell(5, 4),
		Head:   snaky.NewCell(10, 24),
	}

	if err := render.Render(mockBoard); err != nil {
		assert.Error(err)
		return
	}
}
