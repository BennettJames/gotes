package main

import (
	"context"
	"time"

	"github.com/bennettjames/gotes"
	"github.com/gdamore/tcell"
)

func runComponentRenderer(
	ctx context.Context,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sr, srErr := NewScreenRunner()
	if srErr != nil {
		return srErr
	}

	sampleRate := gotes.SampleRate(48000)
	kb := gotes.NewKeyboard(sampleRate, 2000*time.Millisecond)
	speaker := gotes.NewSpeaker(sampleRate, kb, sampleRate.N(100*time.Millisecond))
	go speaker.Run(ctx) // todo (bs): handle error; possibly via brun

	sm := NewGameStateManager()

	// todo (bs): let's reset this on a timer after you deplete it
	sm.SetNotes(time.Now(), getTwinkleNotes())

	go func() {
		for {
			stepTime := time.Now()
			sm.Tick(stepTime)
			gs := sm.State()

			renderGameState(sr, gs)

			select {
			case <-time.After(1000 / 60 * time.Millisecond):
				continue
			case <-ctx.Done():
				return
			}
		}
	}()

	return sr.Run(ctx, func(ev tcell.Event) {
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC, tcell.KeyEscape:
				cancel()

			// note (bs): I kinda feel like these probably aren't long for this world;
			// but oh well. I think I'd likely still want to forward arrow keys, but
			// it'll be in a more generic format that the state manager will
			// conditionally handle based on it's own state.
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
		}
	})
}

func renderGameState(sr *ScreenRunner, gs GameState) {
	// so - this is coming along nicely. Let's spend some time sketching out
	// what this should look like for handling the actual game.
	//
	// first up, I'd say that I should try to make a one page with a few
	// different practical elements on it. How do I want to lay out the
	// state->graphical matching?
	//
	// At the core, I should have a state manager that contains all the game
	// state. For now, this will be singular and total - anything controlling
	// things like different views

	grid := sr.NewGrid()
	w, h := grid.Width(), grid.Height()

	var draw DrawCell = func(x, y int, c ScreenCell) {
		if c.BG == tcell.ColorDefault {
			c.BG = grid.Get(x, y).BG
		}
		grid.Set(x, y, c)
	}

	// so - time to get individual "dot numbers" to render here as well.

	mainC := &AreaComponent{
		X: 0, Y: 0,
		W: w, H: h,
		BgColor: tcell.ColorMidnightBlue,
		Children: []Drawable{
			&AreaComponent{
				X: 10, Y: 5,
				W: 12, H: 5,
				BgColor: tcell.ColorLime,
			},
			&AreaComponent{
				X: 30, Y: 5,
				W: 12, H: 5,
				BgColor:     tcell.ColorPurple,
				Border:      true,
				BorderColor: tcell.ColorSilver,
				Children: []Drawable{
					&TextComponent{
						X: 2, Y: 1,
						MaxWidth:  7,
						MaxHeight: 3,
						Text:      "hello world!!!",
						TextColor: tcell.ColorWhite,
						BgColor:   tcell.ColorDefault,
					},
				},
			},
			&AreaComponent{
				X: 50, Y: 5,
				W: 12, H: 5,
				BgColor:     tcell.ColorDefault,
				Border:      true,
				BorderColor: tcell.ColorSilver,
			},
			&AreaComponent{
				X: 10, Y: 14,
				W: 52, H: 20,
				BgColor:     tcell.ColorDefault,
				Border:      true,
				BorderColor: tcell.ColorSilver,
				// so - let's see about getting a basic "dot drawer" in here. Rough
				// idea: a "dot area" can be a drawable; but it also can be a
				// "DotDrawable".
				Children: []Drawable{
					&DotAreaComponent{
						X: 1, Y: 1,
						W: 50, H: 18,
						BG: tcell.ColorGreen,
						// Children: []DotDrawable{
						// 	&SingleDotComponent{
						// 		X: 5, Y: 4,
						// 		C: tcell.ColorBlack,
						// 	},
						// 	&SingleDotComponent{
						// 		X: 6, Y: 4,
						// 		C: tcell.ColorBlack,
						// 	},
						// 	&SingleDotComponent{
						// 		X: 7, Y: 4,
						// 		C: tcell.ColorBlack,
						// 	},
						// 	&SingleDotComponent{
						// 		X: 8, Y: 4,
						// 		C: tcell.ColorBlack,
						// 	},
						// },
					},
				},
			},
		},
	}

	mainC.Draw(draw)

	sr.Render(grid)
}
