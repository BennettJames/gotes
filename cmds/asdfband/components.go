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

	DrawCell func(x, y int, char rune, fg, bg tcell.Color)

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
			dc(x, y, c, ac.BorderColor, ac.BgColor)
		}
	}

	var dcWrap DrawCell = func(x, y int, char rune, fg, bg tcell.Color) {
		if !ac.Overflow {
			// todo (bs): forbid border overflow as well
			if ac.Border {
				if x < 1 || x >= ac.W-1 || y < 1 || y > ac.H-1 {
					return
				}
			} else {
				if x < 0 || x >= ac.W || y < 0 || y > ac.H {
					return
				}
			}
		}
		dc(x+ac.X, y+ac.Y, char, fg, bg)
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
		dc(x, y, r, tc.TextColor, tc.BgColor)
		x++
	}
}

func AreaWithPos(a AreaComponent, x, y int) AreaComponent {
	newArea := a
	newArea.X = x
	newArea.Y = y
	return newArea
}

func AreaWithSize(a AreaComponent, w, h int) AreaComponent {
	newArea := a
	newArea.W = w
	newArea.H = h
	return newArea
}

func AreaWithBgColor(a AreaComponent, bg tcell.Color) AreaComponent {
	newArea := a
	newArea.BgColor = bg
	return newArea
}

func AreaWithBorder(a AreaComponent, border bool) AreaComponent {
	newArea := a
	newArea.Border = border
	return newArea
}

func AreaWithBorderColor(a AreaComponent, bc tcell.Color) AreaComponent {
	newArea := a
	newArea.BorderColor = bc
	return newArea
}

func TextComponentWithText(t TextComponent, text string) TextComponent {
	newText := t
	newText.Text = text
	return newText
}
