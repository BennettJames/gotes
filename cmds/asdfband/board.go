package main

import (
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

const (
	// hitRange is the amount of time in either direction of the target time that
	// a note can be played.
	hitRange = 200 * time.Millisecond

	// stepSize is how much time passes between each calculation of game state.
	stepSize = 100 * time.Millisecond
)

type (
	GameStateManager struct {
		l sync.RWMutex

		lastTime time.Time

		state GameState
	}

	GameState struct {
		Offset Pos
		Board  NoteBoard
	}

	BoardChar struct {
		Pos   Pos
		Color tcell.Color
		Char  byte
	}

	Pos struct {
		X, Y int
	}

	NoteBoard struct {
		StartTime, LastTime time.Time
		Delay               time.Duration
		Steps               int

		PlaybackRate float64
		Score        int

		ScheduledNotes []ScheduledNote
		Columns        map[byte]NoteColumn
	}

	NoteColumn struct {
		Char        byte
		ActiveNotes []ColumnChar
		Event       ColumnEvent
	}

	ColumnChar struct {
		Offset time.Duration
		Color  tcell.Color
	}

	ColumnEvent struct {
		Type     string // one of "miss", "hit", or "noop"
		GameTime time.Duration
	}

	ScheduledNote struct {
		At time.Duration

		// note (bs): for now, I'm going to represent notes with bytes A-G with
		// corresponding keys and range in the 4th octave. This should be better.
		DispChar byte

		Note float64
	}
)

// note (bs): let's try to keep types explicit and functions disjoint. For now,
// let's make all all render and update functions proper functions (not
// methods). I'll be ok with query methods on types for now.

func NewGameStateManager() *GameStateManager {
	return &GameStateManager{}
}

func (m *GameStateManager) ModifyOffset(x, y int) {
	m.l.Lock()
	defer m.l.Unlock()

	m.state = gameWithOffset(
		m.state,
		updatePos(m.state.Offset, x, y),
	)
}

func (m *GameStateManager) State() GameState {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.state
}

func (m *GameStateManager) Tick(now time.Time) {
	m.l.Lock()
	defer m.l.Unlock()

	// note (bs): I don't necessarily agree with the current design of this - it's
	// kind of ad hoc grown between a few different styles. Still, this should
	// work o.k. for now.

	m.state = gameWithBoard(
		m.state,
		noteBoardUpdate(m.state.Board, now),
	)
}

func (m *GameStateManager) SetNotes(
	now time.Time,
	delay time.Duration,
	notes []ScheduledNote,
) {
	m.l.Lock()
	defer m.l.Unlock()
	m.state = gameWithBoard(
		m.state,
		initializeBoard(now, delay, notes),
	)
}

func (m *GameStateManager) KeyPress(now time.Time, char byte) bool {
	m.l.Lock()
	defer m.l.Unlock()

	column, hasColumn := m.state.Board.Columns[char]
	if !hasColumn {
		return false
	}

	gameTime := (time.Duration(m.state.Board.Steps) * stepSize) +
		now.Sub(m.state.Board.LastTime)
	newColumn := columnWithHit(
		column,
		gameTime,
	)

	// todo (bs): analyze time versus offsets to determine if event type should by
	// "hit" or "noop". Arguably, that logic should be pushed down into an
	// updater.
	newState := gameWithBoard(
		m.state,
		boardWithColumn(m.state.Board, newColumn),
	)
	hit := len(column.ActiveNotes) > 0 &&
		len(column.ActiveNotes) != len(newColumn.ActiveNotes)
	if hit {
		// todo (bs): eventually, this should be replaced with more precise grading
		// of the hit.
		newState.Board.Score++
	}
	m.state = newState

	return hit
}

func gameWithBoard(game GameState, nb NoteBoard) GameState {
	newGame := game
	newGame.Board = nb
	return newGame
}

func gameWithOffset(game GameState, offset Pos) GameState {
	newGame := game
	newGame.Offset = offset
	return newGame
}

func charWithPos(c BoardChar, p Pos) BoardChar {
	newC := c
	newC.Pos = p
	return newC
}

func updatePos(p Pos, x, y int) Pos {
	newPos := p
	newPos.X += x
	newPos.Y += y
	return newPos
}

func boardWithColumn(nb NoteBoard, nc NoteColumn) NoteBoard {
	newBoard := nb
	newColumns := map[byte]NoteColumn{}
	for b, c := range nb.Columns {
		if b != nc.Char {
			newColumns[b] = c
		}
	}
	newColumns[nc.Char] = nc
	newBoard.Columns = newColumns
	return newBoard
}

func columnCutoff(nc NoteColumn, gameTime time.Duration) NoteColumn {
	newColumn := nc
	newNotes := []ColumnChar{}
	for _, n := range nc.ActiveNotes {
		if n.Offset+hitRange > gameTime {
			newNotes = append(newNotes, n)
		}
	}
	newColumn.ActiveNotes = newNotes
	if len(newNotes) != len(nc.ActiveNotes) {
		newColumn.Event = ColumnEvent{
			Type:     "miss",
			GameTime: gameTime,
		}
	}
	return newColumn
}

func columnWithHit(nc NoteColumn, gameTime time.Duration) NoteColumn {
	newColumn := nc
	// todo (bs): need to amend this to have better time tracking, either in this
	// level or one up, to decide if the time has really been hit.
	numNotes := len(nc.ActiveNotes)
	if numNotes > 0 {
		nextNote := nc.ActiveNotes[0]
		if nextNote.Offset-hitRange <= gameTime &&
			nextNote.Offset+hitRange >= gameTime {
			newColumn.ActiveNotes = nc.ActiveNotes[1:]
			newColumn.Event = ColumnEvent{
				Type:     "hit",
				GameTime: gameTime,
			}
			return newColumn
		}
	}
	newColumn.Event = ColumnEvent{
		Type:     "noop",
		GameTime: gameTime,
	}
	return newColumn
}

func noteBoardUpdate(noteB NoteBoard, now time.Time) NoteBoard {
	newBoard := noteB

	for now.Sub(newBoard.LastTime) >= stepSize {
		gameTime := (time.Duration(newBoard.Steps) * stepSize)

		newCols := map[byte]NoteColumn{}
		for b, c := range newBoard.Columns {
			newCols[b] = columnCutoff(c, gameTime)
		}
		newBoard.Columns = newCols

		// note (bs): I'd consider keeping "time indexing" like this and the actual
		// state on different structural levels of the object.
		newBoard.LastTime = newBoard.LastTime.Add(stepSize)
		newBoard.Steps++
	}

	return newBoard
}

func initializeBoard(
	now time.Time,
	delay time.Duration,
	notes []ScheduledNote,
) NoteBoard {

	cols := map[byte]NoteColumn{}
	// todo (bs): I think this array should be abstracted in some fashion. Nothing
	// fancy; maybe just make a function that returns it somewhere and make that
	// spot in the code responsible for the general rules surrounding
	// keys<->notes.
	for _, b := range []byte{'A', 'S', 'D', 'F', 'J', 'K', 'L'} {
		// todo (bs): color selection should be more orderly, rather than purely
		// random
		chars := []ColumnChar{}
		for _, sn := range notes {
			if sn.DispChar != b {
				continue
			}
			chars = append(chars, ColumnChar{
				Offset: sn.At + delay,
				Color:  randPastel(),
			})
		}

		// note (bs): consider having some sort of "initialized event"; perhaps just
		// have everything fade in.
		newColumn := NoteColumn{
			Char:        b,
			ActiveNotes: chars,
		}
		cols[b] = newColumn
	}

	return NoteBoard{
		StartTime:      now,
		LastTime:       now,
		Delay:          delay,
		Score:          0,
		ScheduledNotes: notes,
		Columns:        cols,
	}
}
