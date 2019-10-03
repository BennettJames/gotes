package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_abc(t *testing.T) {
	start := time.Now()
	delay := 1000 * time.Millisecond
	board := []ScheduledNote{
		ScheduledNote{At: 0 * time.Millisecond, DispChar: 'A'},
		ScheduledNote{At: 100 * time.Millisecond, DispChar: 'S'},
		ScheduledNote{At: 200 * time.Millisecond, DispChar: 'D'},
		ScheduledNote{At: 300 * time.Millisecond, DispChar: 'F'},
		ScheduledNote{At: 400 * time.Millisecond, DispChar: 'J'},
		ScheduledNote{At: 400 * time.Millisecond, DispChar: 'K'},
		ScheduledNote{At: 500 * time.Millisecond, DispChar: 'A'},
	}
	noteB := initializeBoard(start, delay, board)

	aCol := noteB.Columns['A']

	// fixme (bs): formalize these tests

	fmt.Printf("@@@ note board: %+v", aCol)
	fmt.Printf("@@@ note board 1: %+v", columnCutoff(aCol, 0*time.Millisecond))
	fmt.Printf("@@@ note board 2: %+v", columnCutoff(aCol, 1400*time.Millisecond))

}
