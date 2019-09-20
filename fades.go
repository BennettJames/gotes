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

func SigmoidFadeLooperWave(
	dur, fade time.Duration,
	fns ...WaveFn,
) WaveFn {
	return genericFadeLooperWave(dur, fade, sigmoidMix, fns...)
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

func sinMix(t float64, v1, v2 float64) float64 {
	t = math.Min(1, math.Max(0, t))
	fadeOut := (math.Cos(t*math.Pi) + 1) / 2
	fadeIn := 1 - fadeOut
	return v1*fadeOut + v2*fadeIn
}

func sigmoidMix(t float64, v1, v2 float64) float64 {
	t = math.Min(1, math.Max(0, t))
	fadeIn := 1 / (1 + math.Pow(math.E, 6-12*t))
	fadeOut := 1 - fadeIn
	return v1*fadeOut + v2*fadeIn

	// note - other than the popping, the inverse here was kinda interesting.
	// Could I do somehting similar w/out popping? Basically, that would mean
	// rapid inversion, then coming back up.
	//
	// I think rather than do anything now, I'll just make note that being able
	// to more minutely vary and mix notes would be nice.
}

// so - given the now-simplicity of swapping out mixing effects for transitions,
// would I want to try out any other fades? Other options are logarithmic, sine,
// sigmoid.
