package main

import "github.com/gdamore/tcell"

// note (bs): so, let's fool around with some approaches to components.
//
// - First up, I think it'd be good to represent the entire screen in some way.
// It doesn't have to be direct; I just mean a way to create a component with

type (
	Drawable interface {
		Draw(dc DrawCell)
	}

	DrawCell func(x, y int, c ScreenCell)

	AreaComponent struct {
		// note (bs): I suspect I should make some composite structures for
		// positions, mostly to make composing element types easier.

		X, Y int
		W, H int

		Overflow bool

		BgColor     tcell.Color
		Border      bool
		BorderColor tcell.Color

		// note (bs): I think I need a more generic render interface here. That
		// probably means I'll need to bind render methods to elements directly. A
		// little annoying, but not the end of the world.
		Children []Drawable
	}

	TextComponent struct {
		X, Y                int
		MaxWidth, MaxHeight int

		Text string

		TextColor tcell.Color
		BgColor   tcell.Color
	}
)

func (ac *AreaComponent) Draw(dc DrawCell) {
	// todo (bs): render first
	//
	// so - let's just do bg to start; and maybe try out a few different render
	// patterns. Borders can wait.
	//
	// ques - should I do a short diversion to make instrumentation more modular?
	// I'm gonna say no - let's do a copy-pase

	for x := ac.X; x < ac.X+ac.W; x++ {
		for y := ac.Y; y < ac.Y+ac.H; y++ {
			c := ' '
			if ac.Border {
				if x == ac.X || x == ac.X+ac.W-1 {
					c = '█'
				} else if y == ac.Y {
					c = '▀'
				} else if y == ac.Y+ac.H-1 {
					c = '▄'
				}
			}
			dc(x, y, ScreenCell{
				Char: c,
				FG:   ac.BorderColor,
				BG:   ac.BgColor,
			})
		}
	}

	var dcWrap DrawCell = offsetDrawCell(dc, ac.X, ac.Y)
	if !ac.Overflow {
		if ac.Border {
			dcWrap = boundDrawCell(dcWrap, 1, 1, ac.W-1, ac.H-1)
		} else {
			dcWrap = boundDrawCell(dcWrap, 0, 0, ac.W, ac.H)
		}
	}
	for _, c := range ac.Children {
		c.Draw(dcWrap)
	}
}

func (tc *TextComponent) Draw(dc DrawCell) {
	// note (bs): this makes it really easy to override the container element with
	// text. That's not necessarily a problem, but it would be nice if I could
	// limit the context of a child via the parent. I *can* technically do that -
	// I can wrap the draw with someting that
	//
	// I suppose that's not a bad take - a parent can decide to limit the render
	// by triggering a cutoff. Let's take that for a spin.
	x, y := tc.X, tc.Y
	for _, r := range tc.Text {
		if tc.MaxWidth > 0 && x >= tc.X+tc.MaxWidth {
			x = tc.X
			y++
		}
		if tc.MaxHeight > 0 && y >= tc.Y+tc.MaxHeight {
			break
		}
		dc(x, y, ScreenCell{
			Char: r,
			FG:   tc.TextColor,
			BG:   tc.BgColor,
		})
		x++
	}
}

func boundDrawCell(
	dc DrawCell,
	// todo (bs): this is one of my bugbears in that it's a lot of args in a row
	// with identical types and no real ordering. I've already said this, but I'd
	// really like ot figure out a good way to structure the basic x/y/h/w
	// structure, which I think could be applied here.
	minX, minY, maxX, maxY int,
) DrawCell {
	return func(x, y int, c ScreenCell) {
		if x < minX || x >= maxX || y < minY || y >= maxY {
			return
		}
		dc(x, y, c)
	}
}

func offsetDrawCell(dc DrawCell, xOff, yOff int) DrawCell {
	return func(x, y int, c ScreenCell) {
		dc(x+xOff, y+yOff, c)
	}
}
