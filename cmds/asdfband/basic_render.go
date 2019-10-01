package main

import (
	"context"
	"time"

	"github.com/bennettjames/gotes"
	"github.com/gdamore/tcell"
)

func runBasicRenderer(
	ctx context.Context,
	screen tcell.Screen,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sm := NewGameStateManager()

	noteHeight := 32

	noteB := NoteBoard{
		StartTime: time.Now(),
		LastTime:  time.Now(),

		PlaybackRate: 1,
		NoteLimit:    noteHeight,

		ScheduledNotes: getTwinkleNotes(),
	}

	sr := gotes.SampleRate(48000)
	kb := gotes.NewKeyboard(sr, 2000*time.Millisecond)
	speaker := gotes.NewSpeaker(sr, kb, sr.N(100*time.Millisecond))
	go speaker.Run(ctx)

	go func() {
		defer cancel()

		// todo (bs): let's extract this
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

		noteB = noteBoardUpdate(noteB, time.Now())

		// todo (bs): this should be extracted
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

		// curState := sm.State()
		// drawBoard(screen, curState.Board, curState.Offset)

		for _, c := range noteB.ActiveNotes {
			drawCellChar(screen, c.Pos.X, c.Pos.Y, GetCellChar(c.Char), c.Color)
		}

		for i := 0; i < 7*7; i++ {
			drawCell(screen, i, noteHeight+3, tcell.ColorWhite)
		}

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

func getTwinkleNotes() []ScheduledNote {
	baseT := 400 * time.Millisecond
	return []ScheduledNote{
		ScheduledNote{At: 1 * time.Millisecond, DispChar: 'A', Note: gotes.NoteC4},
		ScheduledNote{At: 1 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 2 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 3 * baseT, DispChar: 'K', Note: gotes.NoteA4},
		ScheduledNote{At: 4 * baseT, DispChar: 'K', Note: gotes.NoteA4},
		ScheduledNote{At: 5 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 7 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 8 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 9 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 10 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 11 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 12 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 13 * baseT, DispChar: 'A', Note: gotes.NoteC4},
		ScheduledNote{At: 15 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 16 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 17 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 18 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 19 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 20 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 21 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 23 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 24 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 25 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 26 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 27 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 28 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 29 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 31 * baseT, DispChar: 'A', Note: gotes.NoteC4},
		ScheduledNote{At: 32 * baseT, DispChar: 'A', Note: gotes.NoteC4},
		ScheduledNote{At: 33 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 34 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 35 * baseT, DispChar: 'K', Note: gotes.NoteA4},
		ScheduledNote{At: 36 * baseT, DispChar: 'K', Note: gotes.NoteA4},
		ScheduledNote{At: 37 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 39 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 40 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 41 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 42 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 43 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 44 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 45 * baseT, DispChar: 'A', Note: gotes.NoteC4},
	}
}
