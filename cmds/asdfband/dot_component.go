package main

import "github.com/gdamore/tcell"

type (
	DotDrawable interface {
		DrawDots(dd DrawDot)
	}

	DrawDot func(x, y int, c tcell.Color)

	DotAreaComponent struct {
		X, Y int
		W, H int

		BG tcell.Color

		Children []DotDrawable
	}

	SingleDotComponent struct {
		X, Y int
		C    tcell.Color
	}
)

func (dac *DotAreaComponent) Draw(dc DrawCell) {

	// so - what's the basic idea of what this does? I think it just needs to map
	// dots to dc, and *possibly* apply some bg rules. Let's ignore bg rules to
	// start,
	//
	// also - note that this can have kinda weird behavior with borders on account
	// of the offset width at the top. I'll hypothesize that's an issue I can
	// ignore, but I may well be proven wrong.

	dc = offsetDrawCell(dc, dac.X, dac.Y)
	dc = boundDrawCell(dc, 0, 0, dac.W, dac.H)

	// so - this is moderately complicated. I can do the naive thing of halve the
	// value and write out, but that'll partially fail. I need an intermediary
	// buffer here, that upon completion writes out completely.

	dotBuffer := make([][]tcell.Color, dac.W)
	for i := 0; i < len(dotBuffer); i++ {
		dotBuffer[i] = make([]tcell.Color, dac.H*2)
		for j := 0; j < len(dotBuffer[i]); j++ {
			dotBuffer[i][j] = dac.BG
		}
	}

	var dd DrawDot = func(x, y int, c tcell.Color) {
		if x >= 0 && x < dac.W && y >= 0 && y < dac.H {
			dotBuffer[x][y] = c
		}
	}

	for _, c := range dac.Children {
		c.DrawDots(dd)
	}

	if len(dotBuffer) > 10 && len(dotBuffer[0]) > 10 {
		dd(1, 1, tcell.ColorPurple)
		dd(1, 2, tcell.ColorPurple)
		dd(1, 3, tcell.ColorPurple)
		dd(1, 4, tcell.ColorPurple)
		dd(1, 5, tcell.ColorPurple)
		dd(6, 5, tcell.ColorPurple)
		dd(7, 6, tcell.ColorPurple)
		dd(8, 7, tcell.ColorPurple)
		dd(9, 8, tcell.ColorPurple)
		dd(10, 9, tcell.ColorPurple)
	}

	// fixme (bs): wrong coordinate basis
	for i := 0; i < len(dotBuffer); i++ {
		for j := 0; j < len(dotBuffer[i]); j++ {
			c := dotBuffer[i][j]
			dc(i, j/2, ScreenCell{
				Char: '█',
				FG:   c,
				BG:   c,
			})
		}
	}
}

func (dac *DotAreaComponent) DrawDots(dd DrawDot) {
	// note (bs): I'm not entirely sure if this makes sense. Main issue: there's
	// an x/y scale transformation between this and the plain draw. Can they
	// really be used interchangably, given that those values will have
	// substantially different effects?
	for _, c := range dac.Children {
		c.DrawDots(dd)
	}
}

func (sdc *SingleDotComponent) DrawDots(dd DrawDot) {
	dd(sdc.X, sdc.Y, sdc.C)
}
