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
	var streamer gotes.Streamer
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

	wave = gotes.AmplifyWave(
		gotes.Gain(0.25),
		gotes.IntegrateWave(
			gotes.BadOscillateTime(2.0, 0.2),
			gotes.SinWave(gotes.NoteA3),
		),
	)

	wave = gotes.AmplifyWave(
		gotes.Gain(0.25),
		gotes.IntegrateWave(
			gotes.BadOscillateTime2(1.0, 0.2),
			gotes.SinWave(gotes.NoteA3),
		),
	)

	wave = gotes.AmplifyWave(
		gotes.Gain(0.25),
		gotes.IntegrateWave(
			gotes.OscillateTime(1.0, 0.2),
			gotes.SinWave(gotes.NoteA3),
		),
	)

	wave = gotes.AmplifyWave(
		gotes.Gain(0.5),
		linMixer(
			func(t float64) float64 {
				return math.Sin(t)/2 + 0.5
			},
			gotes.SinWave(gotes.NoteA3),
			gotes.SinWave(gotes.NoteA4),
		),
	)

	wave = gotes.AmplifyWave(
		gotes.Gain(0.5),
		sigMixer(
			func(t float64) float64 {
				return math.Sin(t)/2 + 0.5
			},
			gotes.SinWave(gotes.NoteA3),
			gotes.SinWave(gotes.NoteA4),
		),
	)

	streamer = gotes.StreamerFromWave(wave, sr)

	speaker := gotes.NewSpeaker(sr, streamer, sr.N(200*time.Millisecond))
	log.Fatal(speaker.Run(ctx))

}

func linMixer(
	mixFn func(t float64) float64,
	w1, w2 gotes.WaveFn,
) gotes.WaveFn {

	return func(t float64) float64 {
		mixAmt := mixFn(t)
		if mixAmt <= 0 {
			return w1(t)
		}
		if mixAmt >= 1 {
			return w2(t)
		}
		return w1(t)*(1-mixAmt) + w2(t)*(mixAmt)

	}
}

func sigMixer(
	mixFn func(t float64) float64,
	w1, w2 gotes.WaveFn,
) gotes.WaveFn {

	return func(t float64) float64 {
		mixAmt := mixFn(t)
		if mixAmt <= 0 {
			return w1(t)
		}
		if mixAmt >= 1 {
			return w2(t)
		}

		c1 := math.Sqrt(1 - mixAmt)
		c2 := math.Sqrt(mixAmt)

		return w1(t)*c1 + w2(t)*c2

	}
}
