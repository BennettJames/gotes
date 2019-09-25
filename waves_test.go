package gotes

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func Test_PianoWave(t *testing.T) {
	wave := PianoWave(2000*time.Millisecond, NoteA4)

	fmt.Println(wave(3))
}

func Test_ibmProblem(t *testing.T) {
	type noteSet struct {
		Fifty  bool
		Twenty bool
		Ten    bool
		Five   bool
		One    bool
	}

	// note (bs): this is an overly-simplistic approach to fitting that is based
	// on the observation you can be maximally greedy when the denominators/sizes
	// are structured this way.
	notesSet := func(amt int) noteSet {
		var s noteSet
		reduce := func(note int) {
			for amt-note >= 0 {
				amt -= note
			}
		}
		if amt >= 50 {
			s.Fifty = true
			reduce(50)
		}
		if amt >= 20 {
			s.Twenty = true
			reduce(20)
		}
		if amt >= 10 {
			s.Ten = true
			reduce(10)
		}
		if amt >= 5 {
			s.Five = true
			reduce(5)
		}
		if amt >= 1 {
			s.One = true
			reduce(1)
		}
		return s
	}

	metaSet := map[noteSet]int{}
	for i := 0; i < 100; i++ {
		metaSet[notesSet(i)]++
	}

	// so - how to calculate probability based on this? Basically, for any selection
	sum := float64(0)
	for _, c := range metaSet {
		sum += float64(c) * math.Pow(float64(c)/100, 2)
	}

	fmt.Println("@@@ set size", len(metaSet), sum)
}
