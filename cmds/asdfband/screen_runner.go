package main

import (
	"context"

	"github.com/gdamore/tcell"
)

type (
	// ScreenRunner is a process that when executed, will consume the current tty
	// with a rendered screen.
	ScreenRunner struct {
		screen tcell.Screen
	}

	// ScreenCell represents a single character's values in a terminal.
	ScreenCell struct {
		Char   rune
		FG, BG tcell.Color
	}

	// CellGrid is a two-dimensional grid of cells that represents the full
	// terminal.
	CellGrid [][]ScreenCell
)

// NewScreenRunner creates a new screen runner.
func NewScreenRunner() (*ScreenRunner, error) {
	screen, screenErr := tcell.NewScreen()
	if screenErr != nil {
		return nil, screenErr
	}
	return &ScreenRunner{
		screen: screen,
	}, nil
}

// Run will consume the current tty with the display process. All events will be
// serially passed back through onEv. Note that is a blocking call, and being
// slow to consume it can lead to upstream renderings becoming blocked.
func (sr *ScreenRunner) Run(
	ctx context.Context,
	onEv func(ev tcell.Event),
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if initErr := sr.screen.Init(); initErr != nil {
		return initErr
	}
	go func() {
		<-ctx.Done()
		sr.screen.Fini()
	}()

	for {
		// note (bs): it's possible certain events, like ctrl-c and resizing and
		// exits, should be handled here.
		ev := sr.screen.PollEvent()
		if _, isResize := ev.(*tcell.EventResize); isResize {
			sr.screen.Sync()
		}
		if ev != nil {
			// ques (bs): if this happens, will it strictly be indicative that the
			// screen has halted? If so, should return here.
			onEv(ev)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Render will clear the screen and display the given grid.
func (sr *ScreenRunner) Render(cg CellGrid) {
	sr.screen.Clear()
	for x := 0; x < len(cg); x++ {
		row := cg[x]
		for y := 0; y < len(row); y++ {
			cell := row[y]
			s := tcell.StyleDefault.Background(cell.BG).Foreground(cell.FG)
			sr.screen.SetContent(x, y, cell.Char, nil, s)
		}
	}
	sr.screen.Show()
}

// NewGrid returns a new CellGrid whose dimensions match the screen.
func (sr *ScreenRunner) NewGrid() CellGrid {
	x, y := sr.screen.Size()
	return NewCellGrid(x, y)
}

// NewCellGrid creates a new CellGrid with the given dimensions.
func NewCellGrid(x, y int) CellGrid {
	cg := make(CellGrid, x)
	for i := 0; i < x; i++ {
		cg[i] = make([]ScreenCell, y)
	}
	return cg
}

// Get returns the cell at the given location. Returns an empty cell if the
// location is not in the grid.
func (cg CellGrid) Get(x, y int) ScreenCell {
	if x < 0 || x >= len(cg) {
		return ScreenCell{}
	}
	row := cg[x]
	if y < 0 || y >= len(row) {
		return ScreenCell{}
	}
	return row[y]
}

// Set set's the cell at the given location to the provided value. Does nothing
// if the coordinates are out of range.
func (cg CellGrid) Set(x, y int, sc ScreenCell) {
	if x < 0 || x >= len(cg) {
		return
	}
	row := cg[x]
	if y < 0 || y >= len(row) {
		return
	}
	row[y] = sc
}

// Width returns the width of the grid.
func (cg CellGrid) Width() int {
	return len(cg)
}

// Height returns the height of the grid.
func (cg CellGrid) Height() int {
	if len(cg) == 0 {
		return 0
	}
	return len(cg[0])
}
