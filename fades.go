package gotes

import (
	"math"
	"time"
)

// LinearFadeLooperWave will swap through the provided set of wave functions,
// playing each one for "dur" amount of time and transitioning in time "fade".
// It uses an linear transition function between waves.
func LinearFadeLooperWave(
	dur, fade time.Duration,
	fns ...WaveFn,
) WaveFn {
	return genericFadeLooperWave(dur, fade, linMix, fns...)
}

// ExpFadeLooperWave will swap through the provided set of wave functions,
// playing each one for "dur" amount of time and transitioning in time "fade".
// It uses an exponential transition function between waves.
func ExpFadeLooperWave(
	dur, fade time.Duration,
	fns ...WaveFn,
) WaveFn {
	return genericFadeLooperWave(dur, fade, expMix, fns...)
}

func genericFadeLooperWave(
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
		return mix(math.Mod(t, durT)/fadeT, lastFn(t), nextFn(t))
	}
}

// linMix blends v1 and v2 based on the time value. It's v1 for t<0, v2 for t>1,
// and a linear transition between the two from 0->1.
func linMix(t float64, v1, v2 float64) float64 {
	per := math.Min(1, math.Max(0, t))
	fadeOut := 1 - per
	fadeIn := per
	return v1*fadeOut + v2*fadeIn
}

// expMix blends v1 and v2 based on the time value. It's v1 for t<0, v2 for t>1,
// and an exponential transition in 0->1.
func expMix(t float64, v1, v2 float64) float64 {
	c := 0.1 // note (bs) - technically, this could be made configurable.
	fadeOut := (1 / (1 - c)) * (math.Pow(c, t) - c)
	fadeIn := 1 - fadeOut
	return fadeOut*v1 + fadeIn*v2
}
