package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/bennettjames/gotes"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("running")

	sr := gotes.SampleRate(48000)
	var wave gotes.WaveFn
	// streamer := newNoteSequence(
	// 	250 * time.Millisecond,
	// 	sr,
	// 	[]float64{55, 110},
	// 	[]float64{220},
	// 	[]float64{440, 261},
	// 	[]float64{55, 110},
	// )

	// streamer := multiSin(sr, NoteA1)
	// streamer := multiSquare(sr, NoteA1, NoteA2)
	// streamer := newSquareStream(440, sr)
	// streamer := newSinStream(440, sr)

	// streamer := newNoteSequence(
	// 	500 * time.Millisecond, sr,
	// 	multiSin(sr, NoteA1),
	// 	multiSin(sr, NoteA1, NoteA2),
	// 	multiSin(sr, NoteA2, NoteA3),
	// 	multiSin(sr, NoteA3, NoteA4),
	// 	multiSin(sr, NoteA4, NoteA5),
	// 	multiSin(sr, NoteA5, NoteA6),
	// 	// multiSin(sr, NoteA6, NoteA7),
	// )

	// streamer := newNoteSequence(
	// 	1000 * time.Millisecond, sr,
	// 	multiSin(sr, NoteA4),
	// 	multiSquare(sr, NoteA4),
	// )

	// streamer = newNoteSequence(
	// 	250*time.Millisecond, sr,
	// 	multiSin(sr, NoteF4),
	// 	multiSin(sr, NoteA4),
	// 	multiSin(sr, NoteC4),
	// 	multiSin(sr, NoteC4),
	// 	multiSin(sr, NoteD4),
	// 	multiSin(sr, NoteD4),
	// 	multiSin(sr, NoteC4),
	// 	multiSin(sr, NoteC4),
	// )

	// streamer := newTriangleStream(sr, NoteA4)
	// streamer = multiSin(sr, NoteA2)

	// streamer = newLinearFadeLooper(
	// 	sr,
	// 	250*time.Millisecond, 50*time.Millisecond,
	// 	noteFn(sr, NoteF4),
	// 	noteFn(sr, NoteA4),
	// 	noteFn(sr, NoteC4),
	// 	noteFn(sr, NoteC4),
	// 	noteFn(sr, NoteD4),
	// 	noteFn(sr, NoteD4),
	// 	noteFn(sr, NoteC4),
	// 	noteFn(sr, NoteC4),
	// )

	// streamer = newLinearFadeLooper(
	// 	sr,
	// 	1000*time.Millisecond, 1000*time.Millisecond,
	// 	noteFn(sr, NoteA2),
	// 	noteFn(sr, NoteA5),
	// )

	// so - I think this works. That said, the effect is subtle, and due to my
	// mediocre ear I can't really say it's all that different. It seems like
	// the tone is a little "stronger" here, whereas the linear one spends more
	// time in transition.
	// streamer = newExpFadeLooper(
	// 	sr,
	// 	1000*time.Millisecond, 1000*time.Millisecond,
	// 	noteFn(sr, NoteA2),
	// 	noteFn(sr, NoteA5),
	// )

	// so interesting - square streams at full volume do not change w/
	// multiple. But sin does, and . That makes me happy - I think I
	// understand the mechanics of beep's mixing. It's not really a very
	// sophisticated technique.

	// question - from my understanding, mixing was just chopping values
	// greater than zero.

	// This is an interesting one. Definitely can tell it's quieter. But
	// also interesting is when I place two speakers together and don't
	// negate is that intensity varies by distance - I'm guessing there's
	// some weird interference patterns there.
	//
	// streamer = newNegateStream(streamer)

	// streamer = weirdNote1(sr, NoteA3)
	// streamer = newSinStream(sr, (NoteA2))
	// streamer = gotes.PeriodicFreqNote(sr, gotes.NoteA2, gotes.NoteA3, 10*time.Second)
	// streamer = mistakenPeriodicFreqNote(sr, NoteA2, NoteA3, 20 * time.Second)
	// streamer = basicSinNote(sr, NoteA3)

	// let's see if I can swap out the note functions here with an envelope for
	// streamer = newExpFadeLooper(
	// 	sr,
	// 	500*time.Millisecond, 50*time.Millisecond,
	// 	sinKeyFn(sr, NoteA4, 500*time.Millisecond),
	// )

	// streamer = gotes.WeirdNote2(sr, gotes.NoteA4)
	// streamer = newSinStream(sr, NoteA4)

	// streamer = newExpFadeLooper(
	// 	sr,
	// 	1500*time.Millisecond, 50*time.Millisecond,
	// 	func() beep.Streamer {
	// 		return pianoNote(sr, NoteA4, 1500*time.Millisecond)
	// 	},
	// )

	// streamer = gotes.PeriodicFreqNote(sr, gotes.NoteA2, gotes.NoteA3, 10*time.Second)

	wave = gotes.WaveFnSequence(
		250*time.Millisecond,
		gotes.SinWave(gotes.NoteA1),
		gotes.SinWave(gotes.NoteA2),
		gotes.SinWave(gotes.NoteA3),
		gotes.SinWave(gotes.NoteA4),
		gotes.SinWave(gotes.NoteA5),
	)

	// streamer = gotes.SawWave(gotes.NoteA3)

	wave = gotes.WeirdWave1(gotes.NoteA3)
	wave = gotes.WeirdWave2(gotes.NoteA3)

	wave = gotes.PianoWave(2000*time.Millisecond, gotes.NoteA4)

	wave = gotes.PeriodicSinWave(10*time.Second, gotes.NoteA1, gotes.NoteA4)
	wave = gotes.MistakenPeriodicSinWave(10*time.Second, gotes.NoteA1, gotes.NoteA4)

	wave = gotes.LinearFadeLooperWave(
		250*time.Millisecond,
		150*time.Millisecond,
		gotes.SinWave(gotes.NoteF4),
		gotes.SinWave(gotes.NoteA4),
		gotes.SinWave(gotes.NoteC4),
		gotes.SinWave(gotes.NoteC4),
		gotes.SinWave(gotes.NoteD4),
		gotes.SinWave(gotes.NoteD4),
		gotes.SinWave(gotes.NoteC4),
		gotes.SinWave(gotes.NoteC4),
	)

	// wave = gotes.ExpFadeLooperWave(
	// 	1000*time.Millisecond, 1000*time.Millisecond,
	// 	gotes.SinWave(gotes.NoteA2),
	// 	gotes.SinWave(gotes.NoteA5),
	// )

	streamer := gotes.BiStreamerFromWave(sr, wave)

	vol := gotes.NewVolume(streamer, -2)

	speaker := gotes.NewSpeaker(
		gotes.SampleRate(sr),
		vol,
		sr.N(1000*time.Millisecond))

	var _ = wave
	log.Fatal(speaker.Run(ctx))
}

// analyzeStream logs some simplistic stream information once per second.
// Occasionally useful/interesting when debugging.
func analyzeStream(
	sr gotes.SampleRate,
	srcStream gotes.PerfectStream,
) gotes.PerfectStream {
	startTime := time.Now()
	lastLog := startTime
	totalSamples := 0
	totalCalls := 0

	return func(samples [][2]float64) {

		totalSamples += len(samples)
		totalCalls++
		now := time.Now()
		if now.Sub(lastLog) >= time.Second {
			lastLog = now
			fmt.Println("max/min/average", getInfo(samples))
		}
	}
}

type sampleInfo struct {
	min     float64
	max     float64
	average float64
}

func getInfo(samples [][2]float64) sampleInfo {
	// some notes on sampling rate:
	//
	// - average chunk size seems to be 512, but can vary. Not sure who
	//   determines that; might just be a set value somewhere.
	//
	// - values are in range of [-1, 1]. Average of that per sample is
	//   zero, but can vary a bit.

	info := sampleInfo{
		min: math.MaxFloat64,
		max: -math.MaxFloat64,
	}
	for i := 0; i < len(samples); i++ {
		s0, s1 := samples[i][0], samples[i][1]
		info.min = math.Min(s0, info.min)
		info.min = math.Min(s1, info.min)
		info.max = math.Max(s0, info.max)
		info.max = math.Max(s1, info.max)

		// note (bs): let's turn this into a "raw amplitute" with abs
		// values. This was useful to see the basic pattern, but it's
		// pretty boring now.
		info.average += (math.Abs(s0) + math.Abs(s1)) / 2
	}
	info.average /= float64(len(samples))
	return info
}
