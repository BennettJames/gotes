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
		X, Y  int
		Color tcell.Color
	}

	DotCharComponent struct {
		X, Y  int
		Char  byte
		Color tcell.Color
	}

	DotCharColumnComponent struct {
		X, Y int
		W, H int

		ActiveNotes []DotCharComponent
		Char        byte

		// note (bs): this should probably be done on a per-character basis; or
		// perhaps take a set that will assign different colors at different
		// offsets.
		CharColor tcell.Color

		BG tcell.Color
	}

	// DotRowComponent draws a line at the given coordinates.
	DotRowComponent struct {
		X, Y  int
		W     int
		Color tcell.Color
	}

	DotColumnComponent struct {
		X, Y  int
		H     int
		Color tcell.Color
	}
)

func (da *DotAreaComponent) Draw(dc DrawCell) {
	dc = offsetDrawCell(dc, da.X, da.Y)
	dc = boundDrawCell(dc, 0, 0, da.W, da.H)
	w, h := da.W, da.H*2

	dotBuffer := make([][]tcell.Color, w)
	for i := 0; i < len(dotBuffer); i++ {
		dotBuffer[i] = make([]tcell.Color, h)
		for j := 0; j < len(dotBuffer[i]); j++ {
			dotBuffer[i][j] = da.BG
		}
	}

	var dd DrawDot = func(x, y int, c tcell.Color) {
		if x >= 0 && x < w && y >= 0 && y < h {
			dotBuffer[x][y] = c
		}
	}

	for _, c := range da.Children {
		c.DrawDots(dd)
	}

	for i := 0; i < len(dotBuffer); i++ {
		for j := 0; j+1 < len(dotBuffer[i]); j += 2 {
			c1, c2 := dotBuffer[i][j], dotBuffer[i][j+1]
			dc(i, j/2, ScreenCell{
				Char: '▄',
				FG:   c2,
				BG:   c1,
			})
		}
	}
}

func (da *DotAreaComponent) DrawDots(dd DrawDot) {
	// note (bs): I'm not entirely sure if this makes sense. Main issue: there's
	// an x/y scale transformation between this and the plain draw. Can they
	// really be used interchangably, given that those values will have
	// substantially different effects?
	for _, c := range da.Children {
		c.DrawDots(dd)
	}
}

func (sd *SingleDotComponent) DrawDots(dd DrawDot) {
	dd(sd.X, sd.Y, sd.Color)
}

func (dc *DotCharComponent) DrawDots(dd DrawDot) {
	dd = offsetDrawDots(dd, dc.X, dc.Y)
	charData := GetCellChar(dc.Char)
	for _, cell := range charData.Cells {
		dd(cell.X, cell.Y, dc.Color)
	}
}

func (dcc *DotCharColumnComponent) DrawDots(dd DrawDot) {
	dd = offsetDrawDots(dd, dcc.X, dcc.Y)
	for i := 0; i < dcc.W; i++ {
		for j := 0; j < dcc.H; j++ {
			dd(i, j, dcc.BG)
		}
	}
	for _, cc := range dcc.ActiveNotes {
		// note (bs): I think this is o.k. for now; but is pretty clumsy.
		cc.DrawDots(dd)
	}
}

func (dr *DotRowComponent) DrawDots(dd DrawDot) {
	dd = offsetDrawDots(dd, dr.X, dr.Y)
	for i := 0; i < dr.W; i++ {
		dd(i, 0, dr.Color)
	}
}

func (dc *DotColumnComponent) DrawDots(dd DrawDot) {
	dd = offsetDrawDots(dd, dc.X, dc.Y)
	for i := 0; i < dc.H; i++ {
		dd(0, i, dc.Color)
	}
}

func offsetDrawDots(dd DrawDot, xOff, yOff int) DrawDot {
	return func(x, y int, c tcell.Color) {
		dd(xOff+x, yOff+y, c)
	}
}
