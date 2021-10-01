package snaky

// Snake
type Snake struct {
	move Event
	body []Cell
}

// NewSnake
func NewSnake() *Snake {
	var body = []Cell{}

	// initially make a snake
	// of 04 cells
	for i := 0; i < 4; i++ {
		body = append(body, Cell{
			i: 1 + i, // initial x coord
			j: 3,     // initial y coord
		})
	}

	return &Snake{
		move: EventMoveRight,
		body: body,
	}
}

// Length returs current snake
// length
func (s *Snake) Length() int {
	return len(s.body)
}

// Head return snake head coord
func (s *Snake) Head() Cell {
	return s.body[len(s.body)-1]
}

// IsOnBody
// check whether a Cell in on
// top of snake body (including head)
func (s *Snake) IsOnBody(c Cell) bool {
	for _, b := range s.body {
		if c == b {
			return true
		}
	}
	return false
}
