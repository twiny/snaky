package snaky

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRenderer struct {
	mock.Mock
}

// Listen
func (mr *MockRenderer) Listen(events chan<- Event) error {
	return nil
}

// Render
func (mr *MockRenderer) Render(b *Board) error {
	return nil
}

func TestMoveSnake(t *testing.T) {
	assert := assert.New(t)

	mockRenderer := &MockRenderer{}
	game, err := NewGame(100, 100, Fast, mockRenderer)
	if err != nil {
		assert.Error(err)
		return
	}

	if err := game.moveSnake(); err != nil {
		assert.Error(err)
		return
	}
}
