package snaky

import (
	"errors"
)

var (
	ErrCellOutOfRange = errors.New("cell out of range")
)

// `Icon` is used to draw on grid
// and determine with what
// to fill the `Cell`.
type Icon int

const (
	IconSnakeHead Icon = iota + 1
	IconSnakeBody
	IconFood
)

// Cell represent the coordinate
// of a point on the grid
type Cell struct {
	i, j int
}

// NewCell
func NewCell(i, j int) Cell {
	return Cell{
		i: i,
		j: j,
	}
}

// Coord return cell coordinate
// i, j respectively
func (c *Cell) Coord() (int, int) {
	return c.i, c.j
}

// Board holds game information
// that will be rended every time.
type Board struct {
	X, Y         int // board size
	Round, Score int
	Speed        Speed
	Head         Cell // snake head coord
	Length       int  // snake length
	Food         Cell // food coord
	IsPaused     bool
	Errors       string
	Grid         [][]Icon // grid[i][j] where: // i = col, j = row
}

// NewBoard
func NewBoard(width, height int, snake *Snake, food *Food, speed Speed) *Board {
	grid := make([][]Icon, width)
	for i := range grid {
		grid[i] = make([]Icon, height)
	}

	return &Board{
		X:        width,
		Y:        height,
		Score:    0,
		Round:    0,
		Speed:    speed,
		Head:     snake.Head(),
		Length:   snake.Length(),
		IsPaused: true,
		Errors:   "",
		Grid:     grid,
	}
}

// Paint board
func (b *Board) paint(snake *Snake, food *Food) error {
	// update & clear grid
	b.update(snake, food)

	head := snake.Head()

	body := snake.body[:len(snake.body)-1]

	// mark food
	if err := b.mark(food.cell, IconFood); err != nil {
		return err
	}

	// mark snake head
	if err := b.mark(head, IconSnakeHead); err != nil {
		return err
	}

	// mark snake body
	for _, c := range body {
		if err := b.mark(c, IconSnakeBody); err != nil {
			return err
		}
	}

	return nil
}

// update
func (b *Board) update(snake *Snake, food *Food) {
	// clear grid
	grid := make([][]Icon, b.X)
	for i := range grid {
		grid[i] = make([]Icon, b.Y)
	}
	//
	b.Grid = grid

	b.Head = snake.Head()
	b.Length = snake.Length()
	b.Food = food.cell
}

// mark a Cell on the board with a Icon
func (b *Board) mark(c Cell, i Icon) error {
	if c.j > b.X || c.j > b.Y {
		return ErrCellOutOfRange
	}

	b.Grid[c.i][c.j] = i

	return nil
}
