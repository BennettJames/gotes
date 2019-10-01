package main

import (
	"time"

	"github.com/bennettjames/gotes"
)

func getTwinkleNotes() []ScheduledNote {
	baseT := 400 * time.Millisecond
	return []ScheduledNote{
		ScheduledNote{At: 1 * time.Millisecond, DispChar: 'A', Note: gotes.NoteC4},
		ScheduledNote{At: 1 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 2 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 3 * baseT, DispChar: 'K', Note: gotes.NoteA4},
		ScheduledNote{At: 4 * baseT, DispChar: 'K', Note: gotes.NoteA4},
		ScheduledNote{At: 5 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 7 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 8 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 9 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 10 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 11 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 12 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 13 * baseT, DispChar: 'A', Note: gotes.NoteC4},
		ScheduledNote{At: 15 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 16 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 17 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 18 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 19 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 20 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 21 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 23 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 24 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 25 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 26 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 27 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 28 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 29 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 31 * baseT, DispChar: 'A', Note: gotes.NoteC4},
		ScheduledNote{At: 32 * baseT, DispChar: 'A', Note: gotes.NoteC4},
		ScheduledNote{At: 33 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 34 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 35 * baseT, DispChar: 'K', Note: gotes.NoteA4},
		ScheduledNote{At: 36 * baseT, DispChar: 'K', Note: gotes.NoteA4},
		ScheduledNote{At: 37 * baseT, DispChar: 'J', Note: gotes.NoteG4},
		ScheduledNote{At: 39 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 40 * baseT, DispChar: 'F', Note: gotes.NoteF4},
		ScheduledNote{At: 41 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 42 * baseT, DispChar: 'D', Note: gotes.NoteE4},
		ScheduledNote{At: 43 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 44 * baseT, DispChar: 'S', Note: gotes.NoteD4},
		ScheduledNote{At: 45 * baseT, DispChar: 'A', Note: gotes.NoteC4},
	}
}
