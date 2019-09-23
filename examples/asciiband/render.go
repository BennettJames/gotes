package main

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell"
)

func runRenderer(
	ctx context.Context,
	screen tcell.Screen,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// note (bs): while this is of course a trivial amount of data, I sorta think
	// this should be moved to a board state sooner rather than later. I'm not
	// going to care *quite* yet, as I'd also say the characters themselves are
	// just a hacked experiment that should be moved out w/ the board as well, but
	// it should happen soon all the same.
	var xOffset, yOffset int64

	update := make(chan struct{}, 1)

	// so - this is likely intercepting ctrl+c. Can I easily tell if the user is
	// hitting that themselves? Or will I need a fancy stack machine? Or is it
	// just impossible given the keyboard affordances that are granted here?
	go func() {
		defer cancel()

		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					return
				case tcell.KeyCtrlL:
					screen.Sync()

				case tcell.KeyDown:
					atomic.AddInt64(&yOffset, 1)
				case tcell.KeyUp:
					atomic.AddInt64(&yOffset, -1)
				case tcell.KeyLeft:
					atomic.AddInt64(&xOffset, -1)
				case tcell.KeyRight:
					atomic.AddInt64(&xOffset, 1)
				}

				select {
				case update <- struct{}{}:
				default:
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()

	for {
		screen.Clear()

		x, y := int(atomic.LoadInt64(&xOffset)), int(atomic.LoadInt64(&yOffset))

		for c := byte('A'); c <= 'Z'; c++ {
			i := int(c - 'A')
			drawCellChar(screen, x+(i%13)*5, y+(i/13)*7, GetCellChar(c), tcell.ColorLime)
		}

		for c := byte('0'); c <= '9'; c++ {
			i := int(c - '0')
			drawCellChar(screen, x+(i%13)*5, y+20+(i/13)*7, GetCellChar(c), tcell.ColorLime)
		}

		screen.Show()

		select {
		case <-update:
			continue
		case <-time.After(500 * time.Millisecond):
			// note (bs): in the long term, I this this should be eliminated. All
			// changes should be keyed to a higher-level notion of state updates, and
			// that should be read in from the channel.
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func drawCellChar(screen tcell.Screen, x, y int, c CellChar, color tcell.Color) {
	for _, cell := range c.Cells {
		// drawCellAt(screen, x+cell.X*2, y+cell.Y, color)
		drawCell(screen, x+cell.X, y+cell.Y, color)
	}
}

func drawCell(screen tcell.Screen, x, y int, color tcell.Color) {
	const lowC = '▄'
	const highC = '▀'

	s := tcell.StyleDefault
	s = s.Foreground(color)
	c := lowC
	if y%2 == 0 {
		c = highC
	}

	// preserves the existing "pair" to the cell if it was set already.
	atLoc, _, atStyle, _ := screen.GetContent(x, y/2)
	atFg, atBg, _ := atStyle.Decompose()
	if (atLoc == highC && c == lowC) ||
		(atLoc == lowC && c == highC) {
		s = s.Background(atFg)
	} else if (atLoc == highC && c == highC) ||
		(atLoc == lowC && c == lowC) {
		var _ = atBg
		s = s.Background(atBg)
	}

	screen.SetContent(x, y/2, c, []rune{}, s)
}
