package main

import (
	"sync"
	"time"

	"github.com/gdamore/tcell"
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
		Steps               int

		PlaybackRate float64
		NoteLimit    int

		// note (bs): not necessary quite yet, but there is also a two-handed
		// version of this. Could of course try to

		// how do I want to represent notes? Input will essentially consist of a
		// "schedule" of notes - a list of notes, and when they should be played
		// relative to the start.
		//
		// That can be done via times, but need not - the playback can be adjusted;
		// though it's not a bad idea to have a built-in notion of times
		ScheduledNotes []ScheduledNote

		ActiveNotes []BoardChar
	}

	// so, thinking more about the underlying note structure here: I think there's
	// really two things I need from them:
	//
	// - an underlying set of note information that is fixed; i.e. the actual keys
	//
	// - A dynamic set of keys. Velocity/direction is fixed; but they should start
	// at (x, 0), and go down until y exceed height (let's populate a "range")
	//
	// I think the actual key presses in the end will just be based on the second;
	// as they're the "interactive" portion of the data.

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

func (m *GameStateManager) SetNotes(now time.Time, notes []ScheduledNote) {
	m.l.Lock()
	defer m.l.Unlock()

	m.state = gameWithBoard(
		m.state,
		NoteBoard{
			StartTime:      now,
			LastTime:       now,
			NoteLimit:      32, // todo (bs): this really shouldn't be hard coded
			ScheduledNotes: notes,
		},
	)
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

func noteBoardUpdate(noteB NoteBoard, now time.Time) NoteBoard {
	stepSize := 100 * time.Millisecond

	newBoard := noteB

	rate := newBoard.PlaybackRate
	if newBoard.PlaybackRate <= 0 {
		rate = 1
	}

	for now.Sub(newBoard.LastTime) >= stepSize {

		elapsed := newBoard.LastTime.Sub(newBoard.StartTime)
		newBoard.LastTime = newBoard.LastTime.Add(stepSize)

		for _, sn := range newBoard.ScheduledNotes {
			// todo (bs): let's add a "playback rate" that lets you adjust the speed

			diff := elapsed - time.Duration(float64(sn.At)/rate)
			if diff > 0 || diff <= -stepSize {
				continue
			}

			x := 0
			off := 7 // todo (bs): a little hacky, but real tempted to make this a global
			switch sn.DispChar {
			case 'A':
				x = 0 * off
			case 'S':
				x = 1 * off
			case 'D':
				x = 2 * off
			case 'F':
				x = 3 * off
			case 'J':
				x = 4 * off
			case 'K':
				x = 5 * off
			case 'L':
				x = 6 * off
			}

			newBoard.ActiveNotes = append(newBoard.ActiveNotes, BoardChar{
				Pos: Pos{
					X: x,
					Y: 0,
				},
				Color: randPastel(),
				Char:  sn.DispChar,
			})
		}

		// this may also need some notion of "cutoff", at which point a note is
		// removed. Let's just hardcode it for now.
		newActiveNotes := []BoardChar{}
		for _, an := range newBoard.ActiveNotes {
			updatedNote := charWithPos(an, updatePos(an.Pos, 0, 1))
			if newBoard.NoteLimit > 0 && updatedNote.Pos.Y > newBoard.NoteLimit {
				continue
			}
			newActiveNotes = append(newActiveNotes, updatedNote)
		}
		newBoard.ActiveNotes = newActiveNotes
	}

	return newBoard
}
