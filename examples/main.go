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

	wave = gotes.WaveFnSequence(
		250*time.Millisecond,
		gotes.SinWave(gotes.NoteA1),
		gotes.SinWave(gotes.NoteA2),
		gotes.SinWave(gotes.NoteA3),
		gotes.SinWave(gotes.NoteA4),
		gotes.SinWave(gotes.NoteA5),
	)

	wave = gotes.ExpFadeLooperWave(
		1000*time.Millisecond, 1000*time.Millisecond,
		gotes.SinWave(gotes.NoteA2),
		gotes.SinWave(gotes.NoteA5),
	)

	wave = gotes.SawWave(gotes.NoteA3)

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

	wave = gotes.ExpFadeLooperWave(
		250*time.Millisecond,
		200*time.Millisecond,
		gotes.SinWave(gotes.NoteF4),
		gotes.SinWave(gotes.NoteA4),
		gotes.SinWave(gotes.NoteC4),
		gotes.SinWave(gotes.NoteC4),
		gotes.SinWave(gotes.NoteD4),
		gotes.SinWave(gotes.NoteD4),
		gotes.SinWave(gotes.NoteC4),
		gotes.SinWave(gotes.NoteC4),
	)

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

		info.average += (math.Abs(s0) + math.Abs(s1)) / 2
	}
	info.average /= float64(len(samples))
	return info
}
