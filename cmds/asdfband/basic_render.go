package main

import (
	"context"
	"time"

	"github.com/bennettjames/gotes"
	"github.com/gdamore/tcell"
)

func runRenderer(
	ctx context.Context,
	screen tcell.Screen,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sm := NewGameStateManager()

	sr := gotes.SampleRate(48000)
	kb := gotes.NewKeyboard(sr, 2000*time.Millisecond)
	speaker := gotes.NewSpeaker(sr, kb, sr.N(100*time.Millisecond))
	go speaker.Run(ctx)

	go func() {
		defer cancel()

		for {
			// todo (bs): I think this can likely block indefinitely; I'd rather it be
			// instrumented
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlC:
					// note (bs): this needs to force-propagate a cancel all the way back
					// to the root. It implicitly works now, but this really shouldn't be
					// intercepting baseline signaling key chords like this.
					cancel()
					return
				case tcell.KeyEscape, tcell.KeyEnter:
					return
				case tcell.KeyCtrlL:
					screen.Sync()

				case tcell.KeyDown:
					sm.ModifyOffset(0, 1)
				case tcell.KeyUp:
					sm.ModifyOffset(0, -1)
				case tcell.KeyLeft:
					sm.ModifyOffset(-1, 0)
				case tcell.KeyRight:
					sm.ModifyOffset(1, 0)

				case tcell.KeyRune:
					switch ev.Rune() {
					case 'A', 'a':
						kb.Add(gotes.NoteC4)
					case 'S', 's':
						kb.Add(gotes.NoteD4)
					case 'D', 'd':
						kb.Add(gotes.NoteE4)
					case 'F', 'f':
						kb.Add(gotes.NoteF4)
					case 'J', 'j':
						kb.Add(gotes.NoteG4)
					case 'K', 'k':
						kb.Add(gotes.NoteA4)
					case 'L', 'l':
						kb.Add(gotes.NoteB4)
					}

				}

			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()

	// note (bs): I suspect the render step should be decoupled from the state
	// updates for the most part, and they should instead be coupled via

	for i, c := range AvailableChars() {
		// let's see if I can randomize the colors here
		sm.AddChar(BoardChar{
			Pos: Pos{
				X: (i % 13) * 5,
				Y: (i / 13) * 7,
			},
			Color: randPastel(),
			// Color: tcell.ColorLime,
			Char: c,
		})
	}

	for {
		screen.Clear()

		bgColor := tcell.NewRGBColor(0, 0, 0)
		w, h := screen.Size()
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				screen.SetContent(
					x, y,
					' ', nil,
					tcell.StyleDefault.Background(bgColor),
				)
			}
		}

		curState := sm.State()
		drawBoard(screen, curState.Board, curState.Offset)

		screen.Show()

		select {
		case <-time.After(50 * time.Millisecond):
			// note (bs): in the long term, I this this should be eliminated. All
			// changes should be keyed to a higher-level notion of state updates, and
			// that should be read in from the channel.
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
