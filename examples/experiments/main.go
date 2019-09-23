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
	var streamer gotes.BiStreamer

	wave = gotes.LinearFadeLooperWave(
		250*time.Millisecond,
		2*time.Millisecond,
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

	// wave = gotes.WeirdWave1(gotes.NoteA3)

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

	wave = gotes.SigmoidFadeLooperWave(
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

	wave = gotes.PianoWave(2000*time.Millisecond, gotes.NoteA3)

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
		// gotes.AttackAndDecay(2.0, 6.0),
		gotes.FixedAmplify(1.0),
		gotes.IntegrateWave(
			gotes.MultiplyTime(gotes.NoteA3),
			func(t float64) float64 {
				return gotes.BasicPianoWave(t)
				// return 0.2*math.Sin(8*math.Pi*t) +
				// 	0.8*math.Sin(8*math.Pi*t)*(1.01+0.99*math.Cos(2*math.Pi*t))/2
				// return math.Cos(2*math.Pi*t) + math.Pow(math.Sin(2*math.Pi*t), 2)/2
			},
		),
	)

	// note (bs): so, this works, but is dreadfully inefficient. I think I'd like
	// to create a notion of "finite amplitude functions", wherein a function is
	// able to "know" that it is only active so long, then explicitly can be
	// removed from calculations.
	//
	// Also, I'd like to take a somewhat deeper look at "sustaining notes".
	// Particularly, I think the current decay pattern works best with simple, one
	// stroke full notes. I'd guess it works ok for short hits (just decrease the
	// time), but I'd *guess* there should be a more explicit sustain action in
	// cases where you hit a longer note. Also possible that just doing it short/long
	// should be fine.
	space := 0.4
	wave = combineWave(
		gotes.OffsetWave(space*0, gotes.PianoWave(2000*time.Millisecond, gotes.NoteC4)),
		gotes.OffsetWave(space*1, gotes.PianoWave(2000*time.Millisecond, gotes.NoteC4)),
		gotes.OffsetWave(space*2, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),
		gotes.OffsetWave(space*3, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),
		gotes.OffsetWave(space*4, gotes.PianoWave(2000*time.Millisecond, gotes.NoteA4)),
		gotes.OffsetWave(space*5, gotes.PianoWave(2000*time.Millisecond, gotes.NoteA4)),
		gotes.OffsetWave(space*6, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),

		gotes.OffsetWave(space*8, gotes.PianoWave(2000*time.Millisecond, gotes.NoteF4)),
		gotes.OffsetWave(space*9, gotes.PianoWave(2000*time.Millisecond, gotes.NoteF4)),
		gotes.OffsetWave(space*10, gotes.PianoWave(2000*time.Millisecond, gotes.NoteE4)),
		gotes.OffsetWave(space*11, gotes.PianoWave(2000*time.Millisecond, gotes.NoteE4)),
		gotes.OffsetWave(space*12, gotes.PianoWave(2000*time.Millisecond, gotes.NoteD4)),
		gotes.OffsetWave(space*13, gotes.PianoWave(2000*time.Millisecond, gotes.NoteD4)),
		gotes.OffsetWave(space*14, gotes.PianoWave(2000*time.Millisecond, gotes.NoteC4)),

		gotes.OffsetWave(space*16, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),
		gotes.OffsetWave(space*17, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),
		gotes.OffsetWave(space*18, gotes.PianoWave(2000*time.Millisecond, gotes.NoteF4)),
		gotes.OffsetWave(space*19, gotes.PianoWave(2000*time.Millisecond, gotes.NoteF4)),
		gotes.OffsetWave(space*20, gotes.PianoWave(2000*time.Millisecond, gotes.NoteE4)),
		gotes.OffsetWave(space*21, gotes.PianoWave(2000*time.Millisecond, gotes.NoteE4)),
		gotes.OffsetWave(space*22, gotes.PianoWave(2000*time.Millisecond, gotes.NoteD4)),

		gotes.OffsetWave(space*24, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),
		gotes.OffsetWave(space*25, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),
		gotes.OffsetWave(space*26, gotes.PianoWave(2000*time.Millisecond, gotes.NoteF4)),
		gotes.OffsetWave(space*27, gotes.PianoWave(2000*time.Millisecond, gotes.NoteF4)),
		gotes.OffsetWave(space*28, gotes.PianoWave(2000*time.Millisecond, gotes.NoteE4)),
		gotes.OffsetWave(space*29, gotes.PianoWave(2000*time.Millisecond, gotes.NoteE4)),
		gotes.OffsetWave(space*30, gotes.PianoWave(2000*time.Millisecond, gotes.NoteD4)),

		gotes.OffsetWave(space*32, gotes.PianoWave(2000*time.Millisecond, gotes.NoteC4)),
		gotes.OffsetWave(space*33, gotes.PianoWave(2000*time.Millisecond, gotes.NoteC4)),
		gotes.OffsetWave(space*34, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),
		gotes.OffsetWave(space*35, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),
		gotes.OffsetWave(space*36, gotes.PianoWave(2000*time.Millisecond, gotes.NoteA4)),
		gotes.OffsetWave(space*37, gotes.PianoWave(2000*time.Millisecond, gotes.NoteA4)),
		gotes.OffsetWave(space*38, gotes.PianoWave(2000*time.Millisecond, gotes.NoteG4)),

		gotes.OffsetWave(space*40, gotes.PianoWave(2000*time.Millisecond, gotes.NoteF4)),
		gotes.OffsetWave(space*41, gotes.PianoWave(2000*time.Millisecond, gotes.NoteF4)),
		gotes.OffsetWave(space*42, gotes.PianoWave(2000*time.Millisecond, gotes.NoteE4)),
		gotes.OffsetWave(space*43, gotes.PianoWave(2000*time.Millisecond, gotes.NoteE4)),
		gotes.OffsetWave(space*44, gotes.PianoWave(2000*time.Millisecond, gotes.NoteD4)),
		gotes.OffsetWave(space*45, gotes.PianoWave(2000*time.Millisecond, gotes.NoteD4)),
		gotes.OffsetWave(space*46, gotes.PianoWave(2000*time.Millisecond, gotes.NoteC4)),

		// gotes.PianoWave(2000*time.Millisecond, gotes.NoteA4),
		// gotes.OffsetWave(0.25, gotes.PianoWave(2000*time.Millisecond, gotes.NoteA4)),
		// gotes.PianoWave(2000*time.Millisecond, gotes.NoteA2),
	)

	oldWave := wave
	wave = func(t float64) float64 {
		return oldWave(math.Mod(t, space*55+2))
	}

	// streamer = gotes.NewGainStreamer(sr, wave, 0.4)
	streamer = gotes.NewGainStreamer(sr, wave, 0.4)

	speaker := gotes.NewSpeaker(
		gotes.SampleRate(sr),
		streamer,
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

func combineWave(fns ...gotes.WaveFn) gotes.WaveFn {
	return func(t float64) float64 {
		var v float64
		for _, f := range fns {
			v += f(t)
		}
		return v
	}
}
