package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

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
	s = s.Background(atBg)
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

func randPastel() tcell.Color {
	return tcell.NewRGBColor(
		127+rand.Int31n(128),
		127+rand.Int31n(128),
		127+rand.Int31n(128),
	)
}
