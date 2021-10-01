package snaky

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"
)

// Errors
var (
	ErrSnakeBite    = errors.New("snake bite")
	ErrSnakeHitWall = errors.New("wall hit")
)

// Event represent user interaction with the game
// it capture keyboard inputs.
type Event string

const (
	EventQuit    Event = "quit"
	EventRestart Event = "restart"
	EventPause   Event = "pause"
	EventStart   Event = "start"
	//
	EventMoveUp    Event = "up"
	EventMoveDown  Event = "down"
	EventMoveRight Event = "right"
	EventMoveLeft  Event = "left"
)

// Game Speed
type Speed string

const (
	Slow   Speed = "slow"
	Medium Speed = "medium"
	Fast   Speed = "fast"
)

// Renderer
// able to `Listen` to keyboard event and send it to an `Event` channel
// and render game `Board`
type Renderer interface {
	Listen(chan<- Event) error
	Render(b *Board) error
}

// Game
type Game struct {
	mu     *sync.Mutex // used to lock snake movement
	ui     Renderer
	board  *Board
	food   *Food
	snake  *Snake
	speed  Speed
	events chan Event
	paused bool
	hold   chan struct{}
	ctx    context.Context
	quit   context.CancelFunc
}

// NewGame
// height and width
func NewGame(width, height int, speed Speed, ui Renderer) (*Game, error) {
	// random seed
	// used when randomly generating food.
	rand.Seed(time.Now().UnixNano())

	// game context
	ctx, quit := context.WithCancel(context.Background())

	// new inital snake
	snake := NewSnake()

	// new inital food
	food := NewFood()

	paused := true

	return &Game{
		mu:     &sync.Mutex{},
		ui:     ui,
		board:  NewBoard(width, height, snake, food, speed),
		food:   food,
		snake:  snake,
		speed:  speed,
		events: make(chan Event, 1),
		paused: paused,
		hold:   make(chan struct{}, 1),
		ctx:    ctx,
		quit:   quit,
	}, nil
}

// Run
func (g *Game) Run() error {
	// listen to keyboard inputs
	var gerr error = nil
	go func() {
		if err := g.ui.Listen(g.events); err != nil {
			gerr = err
		}
	}()

	// handle events
	go g.eventHandler()

	// paint board
	if err := g.board.paint(g.snake, g.food); err != nil {
		return err
	}

	// render ui
	if err := g.ui.Render(g.board); err != nil {
		return err
	}

	for {
		select {
		case <-g.ctx.Done():
			close(g.events)
			close(g.hold)
			return nil
		default:
			// check if context is finished
			if g.ctx.Err() != nil {
				return g.ctx.Err()
			}

			// check Renderer Listen errors
			if gerr != nil {
				return gerr
			}

			// check if paused
			if g.paused {
				// update board state
				g.board.IsPaused = true

				// paint
				if err := g.board.paint(g.snake, g.food); err != nil {
					g.events <- EventPause
					continue
				}

				// render ui
				if err := g.ui.Render(g.board); err != nil {
					g.events <- EventPause
					continue
				}

				<-g.hold
			}

			if err := g.moveSnake(); err != nil {
				// add error
				g.board.Errors = err.Error()

				// paint
				if err := g.board.paint(g.snake, g.food); err != nil {
					g.events <- EventPause
					continue
				}

				// render ui
				if err := g.ui.Render(g.board); err != nil {
					g.events <- EventPause
					continue
				}

				g.events <- EventPause
				continue
			}

			// update board state
			g.board.IsPaused = false

			if err := g.board.paint(g.snake, g.food); err != nil {
				return err
			}

			// render ui
			if err := g.ui.Render(g.board); err != nil {
				return err
			}

			// sleep
			switch g.speed {
			case Slow:
				time.Sleep(200 * time.Millisecond)
			case Medium:
				time.Sleep(150 * time.Millisecond)
			case Fast:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// restart the game reset snake and food
// to the initial coords
func (g *Game) restart() {
	g.snake = NewSnake()
	g.food = NewFood()
	g.board = NewBoard(g.board.X, g.board.Y, g.snake, g.food, g.speed)
}

// eventHandler process incoming
func (g *Game) eventHandler() {
	for {
		select {
		case <-g.ctx.Done():
			return
		case event := <-g.events:
			g.mu.Lock()
			//
			switch event {
			case EventQuit:
				g.paused = false
				g.hold <- struct{}{} // unblock chan

				g.quit()

			case EventRestart:
				g.paused = false
				g.hold <- struct{}{} // unblock chan

				g.restart()

			case EventPause:
				g.paused = true
			case EventStart:
				g.paused = false
				g.hold <- struct{}{} // unblock chan
			case EventMoveRight:
				if g.snake.move == EventMoveUp || g.snake.move == EventMoveDown {
					g.snake.move = EventMoveRight
					g.board.Round++
				}
			case EventMoveLeft:
				if g.snake.move == EventMoveUp || g.snake.move == EventMoveDown {
					g.snake.move = EventMoveLeft
					g.board.Round++
				}
			case EventMoveUp:
				if g.snake.move == EventMoveRight || g.snake.move == EventMoveLeft {
					g.snake.move = EventMoveUp
					g.board.Round++
				}
			case EventMoveDown:
				if g.snake.move == EventMoveRight || g.snake.move == EventMoveLeft {
					g.snake.move = EventMoveDown
					g.board.Round++
				}
			}
			//
			g.mu.Unlock()
		}
	}
}

// moveSnake
// to move the snake from one `Cell` to the next
// we add a new `head` on the next `Cell`
// and we remove `tail`. if snake ate the `food`
// it grows by keeping the tail
func (g *Game) moveSnake() error {
	// snake head
	head := g.snake.Head()

	// next move
	switch g.snake.move {
	case EventMoveRight:
		head.i++
	case EventMoveLeft:
		head.i--
	case EventMoveUp:
		head.j--
	case EventMoveDown:
		head.j++
	}

	// check if hit the wall
	if (head.i < 0 || head.j < 0) || (head.i >= g.board.X || head.j >= g.board.Y) {
		return ErrSnakeHitWall
	}

	// check if hit itself
	// check body except the head
	body := g.snake.body
	for _, cell := range body[:len(body)-1] {
		if head == cell {
			return ErrSnakeBite
		}
	}

	// check if snake eaten food
	var ate bool = false
	if head == g.food.cell {
		ate = !ate

		// place new food
		var try bool = true

		var f Cell = Cell{}

		for try {
			f.i = rand.Intn(g.board.X) // rand.Intn(max - min) + min
			f.j = rand.Intn(g.board.Y)

			if f != g.food.cell && !g.snake.IsOnBody(f) {
				try = !try
			}
		}

		g.food.cell = f

		// increment score by 10
		g.board.Score = g.board.Score + 10
	}

	// if did not eat food
	// remove tail
	if !ate {
		g.snake.body = g.snake.body[1:]
	}

	// move snake had
	g.snake.body = append(g.snake.body, head)

	return nil
}
