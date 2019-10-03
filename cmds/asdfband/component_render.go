package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bennettjames/gotes"
	"github.com/gdamore/tcell"
)

const (
	// initialNoteDelay is how long there is between the notes first loading and
	// the first time they become "playable".
	initialNoteDelay = 3000 * time.Millisecond

	// targetPlayHeight is the y-offset of where the notes "proper play time" is
	// supposed to occur.
	targetPlayHeight = 34

	// dotCharWidth is the width of a dot character in a dot grid.
	dotCharWidth = 4

	// dotCharHeight is the height of a dot character in a dot grid.
	dotCharHeight = 6
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

	sm.SetNotes(time.Now(), initialNoteDelay, getTwinkleNotes())

	go func() {
		// Give time for screen to boot. Note that this does (clumsily) avoid a race
		// condition; wherein this tries to render before the screen has been
		// initialized. I think it should be handled better though; any bootup wait
		// time caused by this should be better handled.
		time.Sleep(100 * time.Millisecond)
		for {
			stepTime := time.Now()
			sm.Tick(stepTime)
			gs := sm.State()

			drawGameState(sr, gs)

			select {
			case <-time.After(1000 / 30 * time.Millisecond):
				// todo (bs): consider using a more sophisticated waiting system to make
				// the framerate more consistent.
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

			// note (bs): may wish to make this more data-structure oriented
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'A', 'a':
					if sm.KeyPress(time.Now(), 'A') {
						kb.Add(gotes.NoteC4)
					}
				case 'S', 's':
					if sm.KeyPress(time.Now(), 'S') {
						kb.Add(gotes.NoteD4)
					}
				case 'D', 'd':
					if sm.KeyPress(time.Now(), 'D') {
						kb.Add(gotes.NoteE4)
					}
				case 'F', 'f':
					if sm.KeyPress(time.Now(), 'F') {
						kb.Add(gotes.NoteF4)
					}
				case 'J', 'j':
					if sm.KeyPress(time.Now(), 'J') {
						kb.Add(gotes.NoteG4)
					}
				case 'K', 'k':
					if sm.KeyPress(time.Now(), 'K') {
						kb.Add(gotes.NoteA4)
					}
				case 'L', 'l':
					if sm.KeyPress(time.Now(), 'L') {
						kb.Add(gotes.NoteB4)
					}
				}
			}
		}
	})
}

func drawGameState(sr *ScreenRunner, gs GameState) {
	grid := sr.NewGrid()
	w, h := grid.Width(), grid.Height()

	var draw DrawCell = func(x, y int, c ScreenCell) {
		if c.BG == tcell.ColorDefault {
			c.BG = grid.Get(x, y).BG
		}
		grid.Set(x, y, c)
	}

	boardC := renderState(gs, w, h)
	boardC.Draw(draw)
	sr.Render(grid)
}

func renderState(gs GameState, w, h int) Drawable {
	gameTime := time.Duration(gs.Board.Steps) * stepSize
	return &AreaComponent{
		X: 0, Y: 0,
		W: w, H: h,
		BgColor: tcell.NewRGBColor(20, 20, 40),
		Children: []Drawable{
			renderBoard(gs.Board, tcell.NewRGBColor(50, 50, 50), gameTime),
			&AreaComponent{
				X: 84, Y: 20,
				W: 18, H: 3,
				BgColor:     tcell.NewRGBColor(80, 40, 80),
				Border:      true,
				BorderColor: tcell.ColorSilver,
				Children: []Drawable{
					&TextComponent{
						X: 2, Y: 1,
						MaxWidth:  14,
						MaxHeight: 1,
						Text:      fmt.Sprintf("Score: %06d", gs.Board.Score),
						TextColor: tcell.ColorWhite,
						BgColor:   tcell.ColorDefault,
					},
				},
			},
		},
	}
}

func renderBoard(nb NoteBoard, bg tcell.Color, gameTime time.Duration) Drawable {
	const innerWidth = 7 * 8 // todo (bs): hoist
	const innerHeight = 18

	dotChars := []DotDrawable{}
	for i, b := range []byte{'A', 'S', 'D', 'F', 'J', 'K', 'L'} {
		c, hasC := nb.Columns[b]
		if !hasC {
			continue
		}
		dotChars = append(
			dotChars,
			renderCharColumn(c, i*8+1, gameTime),
		)
	}

	borderColor := tcell.ColorDarkGray
	dotChars = append(dotChars, &DotRowComponent{
		X: 0, Y: 0,
		W:     innerWidth + 2,
		Color: borderColor,
	})
	dotChars = append(dotChars, &DotColumnComponent{
		X: innerWidth + 1, Y: 0,
		H:     innerHeight*2 + 3,
		Color: borderColor,
	})
	dotChars = append(dotChars, &DotRowComponent{
		X: 0, Y: innerHeight*2 + 3,
		W:     innerWidth + 2,
		Color: borderColor,
	})
	dotChars = append(dotChars, &DotColumnComponent{
		X: 0, Y: 0,
		H:     innerHeight*2 + 3,
		Color: borderColor,
	})

	dotChars = append(dotChars, &DotRowComponent{
		X: 1, Y: targetPlayHeight,
		W:     innerWidth,
		Color: tcell.ColorWhite,
	})

	return &DotAreaComponent{
		// note (bs): these dimensions are a little "over massaged" - probably would
		// be better if you could make this a little more injectable/dynamic based
		// on what you know about the
		X: 10, Y: 12,
		W: innerWidth + 2, H: innerHeight + 2,
		BG:       bg,
		Children: dotChars,
	}
}

func renderCharColumn(nc NoteColumn, xOff int, gameTime time.Duration) DotDrawable {
	const evDur = 300 * time.Millisecond
	appliedBG := tcell.NewRGBColor(50, 50, 50)

	per := float64(gameTime-nc.Event.GameTime) / float64(evDur)
	secondBG := appliedBG
	switch nc.Event.Type {
	case "miss":
		secondBG = tcell.NewRGBColor(100, 50, 50)
	case "hit":
		secondBG = tcell.NewRGBColor(50, 100, 50)
	case "noop":
		secondBG = tcell.NewRGBColor(70, 70, 80)
	}
	appliedBG = blendColors(secondBG, appliedBG, per)

	chars := []DotCharComponent{}
	for _, an := range nc.ActiveNotes {
		dist := int(an.Offset-gameTime) / int(100*time.Millisecond)
		yOff := targetPlayHeight - dist - dotCharHeight
		chars = append(chars, DotCharComponent{
			X:     2,
			Y:     yOff,
			Char:  nc.Char,
			Color: an.Color,
		})
	}

	return &DotCharColumnComponent{
		X: xOff, Y: 1,
		W: 8, H: 38,
		BG:          appliedBG,
		Char:        nc.Char,
		ActiveNotes: chars,
	}
}

func blendColors(c1, c2 tcell.Color, amt float64) tcell.Color {
	if amt <= 0 {
		return c1
	}
	if amt >= 1 {
		return c2
	}
	r1, g1, b1 := c1.RGB()
	r2, g2, b2 := c2.RGB()
	return tcell.NewRGBColor(
		int32(float64(r1)*(1-amt)+float64(r2)*amt),
		int32(float64(g1)*(1-amt)+float64(g2)*amt),
		int32(float64(b1)*(1-amt)+float64(b2)*amt),
	)
}
