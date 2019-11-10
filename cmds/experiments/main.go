package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/bennettjames/gotes"
	"github.com/bennettjames/gotes/internal/iutil"
)

func main() {
	ctx, cancel := iutil.RootContext()
	defer cancel()

	fmt.Println("running")

	sr := gotes.SampleRate(48000)
	var wave gotes.WaveFn
	var _ = wave
	var streamer gotes.BiStreamer
	var _ = streamer

	wave = gotes.LinearFadeLooper(
		250*time.Millisecond,
		2*time.Millisecond,
		gotes.SinWave(gotes.NoteA1),
		gotes.SinWave(gotes.NoteA2),
		gotes.SinWave(gotes.NoteA3),
		gotes.SinWave(gotes.NoteA4),
		gotes.SinWave(gotes.NoteA5),
	)

	wave = gotes.ExpFadeLooper(
		1000*time.Millisecond, 1000*time.Millisecond,
		gotes.SinWave(gotes.NoteA2),
		gotes.SinWave(gotes.NoteA5),
	)

	wave = gotes.SawWave(gotes.NoteA3)

	// wave = gotes.WeirdWave1(gotes.NoteA3)

	wave = gotes.PianoNote(2000*time.Millisecond, gotes.NoteA4)

	wave = gotes.PeriodicSinWave(10*time.Second, gotes.NoteA1, gotes.NoteA4)
	wave = gotes.MistakenPeriodicSinWave(10*time.Second, gotes.NoteA1, gotes.NoteA4)

	wave = gotes.LinearFadeLooper(
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

	wave = gotes.ExpFadeLooper(
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

	wave = gotes.SigmoidFadeLooper(
		250*time.Millisecond,
		100*time.Millisecond,
		gotes.SinWave(gotes.NoteF4),
		gotes.SinWave(gotes.NoteA4),
		gotes.SinWave(gotes.NoteC4),
		gotes.SinWave(gotes.NoteC4),
		gotes.SinWave(gotes.NoteD4),
		gotes.SinWave(gotes.NoteD4),
		gotes.SinWave(gotes.NoteC4),
		gotes.SinWave(gotes.NoteC4),
	)

	wave = gotes.PianoNote(2000*time.Millisecond, gotes.NoteA3)

	wave = gotes.WeirdPianoWave(2000*time.Millisecond, gotes.NoteA3)

	wave = gotes.PeriodicSinWave(3*time.Second, gotes.NoteA1, gotes.NoteA4)

	// this one's vaguely guitar-y
	wave = gotes.AmplifyWave(
		gotes.AttackAndDecay(2.0, 6.0),
		gotes.IntegrateWave(
			gotes.MultiplyTime(gotes.NoteA4),
			func(t float64) float64 {
				return gotes.BasicSinFn(t) + gotes.BasicSinFn(t/2) + gotes.BasicSinFn(t/8)
			},
		),
	)

	// Honestly, this isn't a bad stand-in for the piano note and it's much simpler.
	wave = gotes.AmplifyWave(
		gotes.AttackAndDecay(2.0, 6.0),
		gotes.IntegrateWave(
			gotes.MultiplyTime(gotes.NoteA3),
			func(t float64) float64 {
				return math.Cos(2*math.Pi*t) + math.Pow(math.Sin(2*math.Pi*t), 2)/2
			},
		),
	)

	wave = gotes.AmplifyWave(
		// gotes.AttackAndDecay(2.0, 6.0),
		gotes.FixedAmplify(1.0),
		gotes.IntegrateWave(
			gotes.MultiplyTime(gotes.NoteA3),
			func(t float64) float64 {
				return (0.5+0.3*math.Cos(2*math.Pi*t))*math.Sin(8*math.Pi*t) +
					(0.5-0.3*math.Cos(2*math.Pi*t))*math.Sin(8*math.Pi*t)*(1.01+0.99*math.Cos(2*math.Pi*t))/2
			},
		),
	)

	// note (bs): overall, I do like the direction the function composition here
	// is moving in. But: it still seems like there are some bits, like when to
	// parallelize vs wrap, that isn't well settled.

	wave = gotes.AmplifyWave(
		gotes.FixedAmplify(1.0),
		gotes.PianoWave(gotes.NoteA3),
	)

	wave = gotes.AmplifyWave(
		gotes.FixedAmplify(0.5),
		// let's see if this construction could be decoupled from frequency; e.g.
		// take two fn's rather than two notes.
		//
		// Well, it might not quite be that simple.
		gotes.PeriodicSinWave(
			1*time.Second,
			gotes.NoteA3,
			gotes.NoteA4,
		),
	)

	wave = gotes.AmplifyWave(
		gotes.FixedAmplify(0.5),
		gotes.IntegrateWave(
			func(t float64) float64 {
				// so - I *suspect* that if I'm clever with calculating ratios here; I
				// can achieve what the original was doing much more directly.
				//
				// alright; so this gets the period right. Now I want to make the ratio
				// adjustable. My preference is to always make the "floor" 1, and so the
				// given amount is strictly >= 1 and determines the max rate.
				return 1.5*t + math.Sin(t*2*math.Pi)/(4*math.Pi)
			},
			gotes.IntegrateWave(
				gotes.MultiplyTime(gotes.NoteA3),
				gotes.BasicSinFn,
			),
		),
	)

	// math.Sin(t * 2 * math.PI) / 2 + 0.5

	// just as another reminder; there was one other thing I wanted to experiment
	// with in terms of alternation - shifting directly between two wave forms of
	// identical frequency. Zooming into the mix functions - what I'd guess I'd
	// want to to create a "cyclic mixer" rather than a "plain mixer". That can be
	// done by

	wave = gotes.AmplifyWave(
		gotes.FixedAmplify(0.3),
		func(t float64) float64 {

			v1 := gotes.SinWave(gotes.NoteA3)(t)
			v2 := gotes.SquareWave(gotes.NoteA3)(t)
			return gotes.LinearMix(
				math.Sin(t*2*math.Pi*16)/2+0.5,
				v1, v2)
		},
	)

	streamer = gotes.StreamerFromWave(sr, wave)

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
	streamer = kb

	speaker := gotes.NewSpeaker(sr, streamer, sr.N(100*time.Millisecond))
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

func combineWave(fns ...gotes.WaveFn) gotes.WaveFn {
	return func(t float64) float64 {
		var v float64
		for _, f := range fns {
			v += f(t)
		}
		return v
	}
}
