package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_abc(t *testing.T) {
	start := time.Now()
	noteB := NoteBoard{
		StartTime: start,
		LastTime:  start,

		ScheduledNotes: []ScheduledNote{
			ScheduledNote{At: 0 * time.Millisecond, DispChar: 'A'},
			ScheduledNote{At: 100 * time.Millisecond, DispChar: 'B'},
			ScheduledNote{At: 200 * time.Millisecond, DispChar: 'C'},
			ScheduledNote{At: 300 * time.Millisecond, DispChar: 'D'},
			ScheduledNote{At: 400 * time.Millisecond, DispChar: 'E'},
		},
	}

	for i := 0; i < 10; i++ {
		now := start.Add(time.Duration(150*i) * time.Millisecond)
		noteB = noteBoardUpdate(noteB, now)
	}

	// todo (bs): consider turning this into a proper test

	fmt.Println("@@@ active notes", noteB.ActiveNotes)
}
