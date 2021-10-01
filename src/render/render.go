package render

import (
	"context"
	"snaky/src/snaky"
	"strconv"

	"github.com/gdamore/tcell/v2/termbox"
)

// Render
type Render struct {
	ctx  context.Context
	quit context.CancelFunc
}

// NewRender
func NewRender() (*Render, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}

	ctx, quit := context.WithCancel(context.Background())

	return &Render{
		ctx:  ctx,
		quit: quit,
	}, nil
}

// Listen
// to keyboard events and send it to a receive only channel
func (r *Render) Listen(events chan<- snaky.Event) error {
	for {
		select {
		case <-r.ctx.Done():
			return nil
		default:
			if r.ctx.Err() != nil {
				return r.ctx.Err()
			}

			// get terminal events
			event := termbox.PollEvent()

			switch event.Type {
			// check error
			case termbox.EventError:
				return event.Err
			//
			case termbox.EventInterrupt:
				events <- snaky.EventQuit // close signal
				return nil

			// keyboard events
			case termbox.EventKey:
				switch {
				// exit
				case event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC:
					events <- snaky.EventQuit // close signal
					return nil

				// pause
				case event.Ch == 'p':
					events <- snaky.EventPause

				// start
				case event.Ch == 's':
					events <- snaky.EventStart
					// fmt.Println("Key S")

				// right arrow
				case event.Key == termbox.KeyArrowRight:
					events <- snaky.EventMoveRight

				// left arrow
				case event.Key == termbox.KeyArrowLeft:
					events <- snaky.EventMoveLeft

				// up arrow
				case event.Key == termbox.KeyArrowUp:
					events <- snaky.EventMoveUp

				// down arrow
				case event.Key == termbox.KeyArrowDown:
					events <- snaky.EventMoveDown

				// R key
				case event.Ch == 'r':
					events <- snaky.EventRestart
				}
			}
		}
	}
}

// Render
func (r *Render) Render(b *snaky.Board) error {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	r.render(b)

	return termbox.Flush()
}

// render the game
// the origin of the actual game board grid is shifted by (1,3)
func (r *Render) render(b *snaky.Board) {
	var (
		titleHeight int = 2
		//
		scoreWidth  int = 20 // default side board width
		scoreHeight int = b.Y
		gameWidth   int = 1 + b.X + scoreWidth + 1
		gameHeight  int = 2 + b.Y + 1
		//
		o1x, o1y int = 1 + b.X + 1, titleHeight + 1
	)

	lineX := termbox.Cell{
		Ch: '─',
		Fg: termbox.ColorDefault,
		Bg: termbox.ColorDefault,
	}

	lineY := termbox.Cell{
		Ch: '|',
		Fg: termbox.ColorDefault,
		Bg: termbox.ColorDefault,
	}

	// border
	fill(0, 0, gameWidth, 1, lineX)           // top
	fill(0, titleHeight, gameWidth, 1, lineX) // horizontal line in middle
	//
	fill(o1x, o1y, 1, scoreHeight, lineY)    // vertical line in middle
	fill(gameWidth, 0, 1, gameHeight, lineY) // right
	fill(0, 0, 1, gameHeight, lineY)         // left
	//
	fill(0, gameHeight, gameWidth, 1, lineX) // line in middle

	// corners
	termbox.SetCell(0, 0, '┌', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(0, gameHeight, '└', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(gameWidth, 0, '┐', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(gameWidth, gameHeight, '┘', termbox.ColorDefault, termbox.ColorDefault)

	// print title
	mid := gameWidth / 2
	writeXln(mid-5, 1, termbox.ColorDefault, termbox.ColorDefault, "Snaky Game")

	// head coord
	hi, hj := b.Head.Coord()

	// print score
	writeXln(o1x+2, o1y+1, termbox.ColorDefault, termbox.ColorDefault, "Size: "+strconv.Itoa(b.X)+" x "+strconv.Itoa(b.Y))
	writeXln(o1x+2, o1y+2, termbox.ColorDefault, termbox.ColorDefault, "Speed: "+string(b.Speed))
	writeXln(o1x+2, o1y+3, termbox.ColorDefault, termbox.ColorDefault, "Head: "+strconv.Itoa(hi)+" x "+strconv.Itoa(hj))
	//
	writeXln(o1x+2, o1y+5, termbox.ColorDefault, termbox.ColorDefault, "Round: "+strconv.Itoa(b.Round))
	writeXln(o1x+2, o1y+6, termbox.ColorDefault, termbox.ColorDefault, "Score: "+strconv.Itoa(b.Score))
	writeXln(o1x+2, o1y+7, termbox.ColorDefault, termbox.ColorDefault, "Length: "+strconv.Itoa(b.Length))
	//
	writeXln(o1x+2, o1y+9, termbox.ColorDefault, termbox.ColorDefault, "Arrow: move")
	writeXln(o1x+2, o1y+10, termbox.ColorDefault, termbox.ColorDefault, "R: restart")
	writeXln(o1x+2, o1y+11, termbox.ColorDefault, termbox.ColorDefault, "ESC: quit")

	// message
	// using variables to center the messages
	// on the gird
	msg1 := `press "R" to restart`
	msg2 := `or "ESC" to quit`
	msg3 := `game paused`
	msg4 := `press "S" to start`

	// if error
	if b.Errors != "" {
		writeXln((b.X-len(b.Errors))/2, gameHeight/2, termbox.ColorDefault, termbox.ColorDefault, b.Errors)
		writeXln((b.X-len(msg1))/2, (gameHeight/2)+2, termbox.ColorDefault, termbox.ColorDefault, msg1)
		writeXln((b.X-len(msg2))/2, (gameHeight/2)+3, termbox.ColorDefault, termbox.ColorDefault, msg2)

		return
	}

	// if paused
	if b.IsPaused {
		writeXln((b.X-len(msg3))/2, (gameHeight / 2), termbox.ColorDefault, termbox.ColorDefault, msg3)
		writeXln((b.X-len(msg4))/2, (gameHeight/2)+2, termbox.ColorDefault, termbox.ColorDefault, msg4)
		writeXln((b.X-len(msg2))/2, (gameHeight/2)+3, termbox.ColorDefault, termbox.ColorDefault, msg2)

		return
	}

	// grid
	head := termbox.Cell{
		Ch: 'o',
		Fg: termbox.AttrBold,
		Bg: termbox.ColorDefault,
	}
	body := termbox.Cell{
		Ch: 'x',
		Fg: termbox.ColorDefault,
		Bg: termbox.ColorDefault,
	}
	food := termbox.Cell{
		Ch: '*',
		Fg: termbox.AttrBold,
		Bg: termbox.ColorDefault,
	}

	// print snake & food only if game is running
	// and there are no erros
	for i := range b.Grid {
		col := b.Grid[i]
		for j := range col {
			icon := b.Grid[i][j]
			switch icon {
			case snaky.IconSnakeHead:
				fill(i+1, j+3, 1, 1, head)
			case snaky.IconSnakeBody:
				fill(i+1, j+3, 1, 1, body)
			case snaky.IconFood:
				fill(i+1, j+3, 1, 1, food)
			}
		}
	}
}

// Close
func (r *Render) Close() {
	r.quit()
	termbox.Close()
}
