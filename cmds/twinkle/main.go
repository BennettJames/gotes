package main

import (
	"log"
	"time"

	"github.com/bennettjames/gotes"
	"github.com/bennettjames/gotes/internal/iutil"
)

func main() {
	ctx, cancel := iutil.RootContext()
	defer cancel()

	sr := gotes.SampleRate(48000)

	kb := gotes.NewKeyboard(sr, 2000*time.Millisecond)
	go func() {
		twinkleNotes := []float64{
			gotes.NoteC4, gotes.NoteC4, gotes.NoteG4, gotes.NoteG4,
			gotes.NoteA4, gotes.NoteA4, gotes.NoteG4, 0,

			gotes.NoteF4, gotes.NoteF4, gotes.NoteE4, gotes.NoteE4,
			gotes.NoteD4, gotes.NoteD4, gotes.NoteC4, 0,

			gotes.NoteG4, gotes.NoteG4, gotes.NoteF4, gotes.NoteF4,
			gotes.NoteE4, gotes.NoteE4, gotes.NoteD4, 0,

			gotes.NoteG4, gotes.NoteG4, gotes.NoteF4, gotes.NoteF4,
			gotes.NoteE4, gotes.NoteE4, gotes.NoteD4, 0,

			gotes.NoteC4, gotes.NoteC4, gotes.NoteG4, gotes.NoteG4,
			gotes.NoteA4, gotes.NoteA4, gotes.NoteG4, 0,

			gotes.NoteF4, gotes.NoteF4, gotes.NoteE4, gotes.NoteE4,
			gotes.NoteD4, gotes.NoteD4, gotes.NoteC4, 0,
		}

		for {
			for _, n := range twinkleNotes {
				time.Sleep(400 * time.Millisecond)
				kb.Add(n)
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	speaker := gotes.NewSpeaker(sr, kb, sr.N(200*time.Millisecond))
	log.Fatal(speaker.Run(ctx))
}
