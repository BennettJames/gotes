package gotes

import (
	"math"
	"time"
)

// Looper will swap through the provided set of wave functions, playing each one
// for "dur" amount of time. There is no transition; so unless the wave
// functions themselves blend well this can easily lead to popping.
func Looper(
	dur time.Duration,
	fns ...WaveFn,
) WaveFn {
	if len(fns) == 0 {
		return ZeroWave()
	}
	durT := float64(dur) / float64(time.Second)
	return func(t float64) float64 {
		nextI := (int(t/durT) + len(fns)) % len(fns)
		t = math.Mod(t, durT) // note (bs): inefficient
		return fns[nextI](t)
	}
}

// LinearFadeLooper will swap through the provided set of wave functions,
// playing each one for "dur" amount of time and transitioning in time "fade".
// It uses an linear transition function between waves.
func LinearFadeLooper(
	dur, fade time.Duration,
	fns ...WaveFn,
) WaveFn {
	return genericFadeLooper(dur, fade, LinearMix, fns...)
}

// ExpFadeLooper will swap through the provided set of wave functions,
// playing each one for "dur" amount of time and transitioning in time "fade".
// It uses an exponential transition function between waves.
func ExpFadeLooper(
	dur, fade time.Duration,
	fns ...WaveFn,
) WaveFn {
	return genericFadeLooper(dur, fade, ExpMix, fns...)
}

func SinFadeLooper(
	dur, fade time.Duration,
	fns ...WaveFn,
) WaveFn {
	return genericFadeLooper(dur, fade, SinMix, fns...)
}

func SigmoidFadeLooper(
	dur, fade time.Duration,
	fns ...WaveFn,
) WaveFn {
	return genericFadeLooper(dur, fade, SigmoidMix, fns...)
}

func genericFadeLooper(
	dur, fade time.Duration,
	mix func(t float64, v1, v2 float64) float64,
	fns ...WaveFn,
) WaveFn {
	if len(fns) == 0 {
		fns = []WaveFn{
			ZeroWave(),
		}
	}
	durT := float64(dur) / float64(time.Second)
	fadeT := float64(fade) / float64(time.Second)
	return func(t float64) float64 {
		nextI := (int(t/durT) + len(fns)) % len(fns)
		lastI := (nextI - 1 + len(fns)) % len(fns)
		lastFn, nextFn := fns[lastI], fns[nextI]
		if nextI == 0 && t < durT {
			lastFn = ZeroWave()
		}
		// fixme (bs): math.mod is very inefficient
		return mix(math.Mod(t, durT)/fadeT, lastFn(t), nextFn(t))
	}
}

// LinearMix blends v1 and v2 based on the time value. It's v1 for t<0, v2 for t>1,
// and a linear transition between the two from 0->1.
func LinearMix(t float64, v1, v2 float64) float64 {
	per := math.Min(1, math.Max(0, t))
	fadeOut := 1 - per
	fadeIn := per
	return v1*fadeOut + v2*fadeIn
}

// ExpMix blends v1 and v2 based on the time value. It's v1 for t<0, v2 for t>1,
// and an exponential transition in 0->1.
func ExpMix(t float64, v1, v2 float64) float64 {
	c := 0.1 // note (bs) - technically, this could be made configurable.
	fadeOut := (1 / (1 - c)) * (math.Pow(c, t) - c)
	fadeIn := 1 - fadeOut
	return fadeOut*v1 + fadeIn*v2
}

func SinMix(t float64, v1, v2 float64) float64 {
	t = math.Min(1, math.Max(0, t))
	fadeOut := (math.Cos(t*math.Pi) + 1) / 2
	fadeIn := 1 - fadeOut
	return v1*fadeOut + v2*fadeIn
}

func SigmoidMix(t float64, v1, v2 float64) float64 {
	t = math.Min(1, math.Max(0, t))
	fadeIn := 1 / (1 + math.Pow(math.E, 6-12*t))
	fadeOut := 1 - fadeIn
	return v1*fadeOut + v2*fadeIn
}
