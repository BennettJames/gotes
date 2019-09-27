package main

import (
	"sync"

	"github.com/gdamore/tcell"
)

type (
	GameStateManager struct {
		l sync.RWMutex

		state GameState
	}

	GameState struct {
		Offset Pos
		Board  Board
	}

	Board struct {
		// note (bs): I doubt this array is really good enough for further
		// "mutability"; but it'll be ok for now.
		Chars []BoardChar
	}

	BoardChar struct {
		Pos   Pos
		Color tcell.Color
		Char  byte
	}

	Pos struct {
		X, Y int
	}
)

// note (bs): let's try to keep types explicit and functions disjoint. For now,
// let's make all all render and update functions proper functions (not
// methods). I'll be ok with query methods on types for now.

func NewGameStateManager() *GameStateManager {
	return &GameStateManager{}
}

func (m *GameStateManager) AddChar(c BoardChar) {
	m.l.Lock()
	defer m.l.Unlock()

	// note (bs): this starts running into an issue I've had before, which is how
	// to properly handle nested objects in a single state manager. Let's sit on
	// that question initially - just having one top level manager for events
	// should be fine to start - but that is an important question I'd like some
	// clarity on.
	m.state = gameWithBoard(
		m.state,
		boardWithChar(m.state.Board, c),
	)
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

func gameWithBoard(game GameState, b Board) GameState {
	newGame := game
	newGame.Board = b
	return newGame
}

func gameWithOffset(game GameState, offset Pos) GameState {
	newGame := game
	newGame.Offset = offset
	return newGame
}

func boardWithChar(b Board, c BoardChar) Board {
	newBoard := b
	newBoard.Chars = append(b.Chars[:], c)
	return newBoard
}

func updatePos(p Pos, x, y int) Pos {
	newPos := p
	newPos.X += x
	newPos.Y += y
	return newPos
}
