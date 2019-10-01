package main

import (
	"context"
	"time"

	"github.com/gdamore/tcell"
)

func runComponentRenderer(
	ctx context.Context,
	screen tcell.Screen,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer cancel()

		// todo (bs): let's extract this. Let's also try to see if I can add a
		// little (just a little) more distance between top-level
		for {
			// todo (bs): I think this can likely block indefinitely; I'd rather it be
			// instrumented. Not sure I can do that though.
			//
			// I could kinda of indirectly instrument this with a run-block style
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
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()

	for {
		screen.Clear()

		var draw DrawCell = func(x, y int, char rune, fg, bg tcell.Color) {
			s := tcell.StyleDefault
			s = s.Foreground(fg)
			if bg == tcell.ColorDefault {
				// note (bs): this preserves the parent bg style, in the case that
				_, _, atStyle, _ := screen.GetContent(x, y)
				_, atBg, _ := atStyle.Decompose()
				s = s.Background(atBg)
			} else {
				s = s.Background(bg)
			}
			screen.SetContent(x, y, char, []rune{}, s)
		}

		w, h := screen.Size()

		// so - this is coming along nicely. Let's spend some time sketching out
		// what this should look like for handling the actual game.
		//
		// first up, I'd say that I should try to make a one page with
		//
		// ques - do I want to try to brun-itize the screen? Possibly. Rough idea: I
		// only to to initialize, then run looking for events. You could expose
		// whatever methods from the main object you'd like still; but truthfully I
		// think I can reduce it down to "run" and "render frame".
		//
		// also, do I want to make use of someting CellBuffer-like directly? I have
		// wondered before if I'd be better off having a mutable data state to
		// reflect the screen state, then pass it to the screen to indicate a new
		// render.

		mainC := &AreaComponent{
			X: 0, Y: 0,
			W: w, H: h,
			BgColor: tcell.ColorBlack,
			Children: []Drawable{
				&AreaComponent{
					X: 10, Y: 10,
					W: 10, H: 5,
					BgColor:     tcell.ColorLime,
					Border:      true,
					BorderColor: tcell.ColorSilver,
				},
				&AreaComponent{
					X: 30, Y: 10,
					W: 10, H: 5,
					BgColor:     tcell.ColorPurple,
					Border:      true,
					BorderColor: tcell.ColorSilver,
					Children: []Drawable{
						&TextComponent{
							X: 2, Y: 1,
							MaxWidth:  6,
							MaxHeight: 3,
							Text:      "hello world!!!",
							TextColor: tcell.ColorWhite,
							BgColor:   tcell.ColorDefault,
						},
					},
				},
				&AreaComponent{
					X: 50, Y: 10,
					W: 10, H: 5,
					BgColor:     tcell.ColorDefault,
					Border:      true,
					BorderColor: tcell.ColorSilver,
				},
			},
		}

		mainC.Draw(draw)

		screen.Show()

		select {
		case <-time.After(50 * time.Millisecond):
			// note (bs): in the long term, I this this should be eliminated. All
			// changes should be keyed to a higher-level notion of state updates, and
			// that should be read in from the channel.
			//
			// So you could

			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
