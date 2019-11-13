package main

import (
	"context"
	"log"
	"time"

	"github.com/bennettjames/gotes"
)

func main() {
	const (
		sampleRate   = gotes.SampleRate(48000)
		speakerCache = 100 * time.Millisecond
		noteDuration = 2000 * time.Millisecond
	)

	wave := gotes.Looper(
		time.Second,
		gotes.PianoNote(noteDuration, gotes.NoteA3),
		gotes.PianoNote(noteDuration, gotes.NoteB3),
		gotes.PianoNote(noteDuration, gotes.NoteG4),
		gotes.PianoNote(noteDuration, gotes.NoteF4),
	)
	streamer := gotes.StreamerFromWave(sampleRate, wave)
	speaker := gotes.NewSpeaker(sampleRate, streamer, sampleRate.N(speakerCache))

	log.Fatal(speaker.Run(context.Background()))
}
