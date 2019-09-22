package gotes

import (
	"math"
	"time"
)

// TimeFn maps the given time to another time. This is useful for certain wave
// transformations - for instance, the frequency of a given wave can be modified
// by applying a constant multiplier to time.
//
// Note that while this has the same signature of WaveFn, it's implied contract
// is different. Rather than being a periodic function that maps to [-1,1], this
// will instead map t to another t in a monotonically increasing fashion.
type TimeFn func(t float64) float64

// AmpFn yields a multiple for the base wave at the given time.
type AmpFn func(t float64) float64

// IntegrateWave will call the given wave function with the given mapped time
// function for any given time.
func IntegrateWave(tFn TimeFn, wFn WaveFn) WaveFn {
	return func(t float64) float64 {
		return wFn(tFn(t))
	}
}

// AmplifyWave will amplify the given wave function at each time point by the
// value yielded by the amplify function.
func AmplifyWave(aFn AmpFn, wFn WaveFn) WaveFn {
	return func(t float64) float64 {
		return wFn(t) * aFn(t)
	}
}

func FixedAmplify(amp float64) AmpFn {
	return func(t float64) float64 {
		return amp
	}
}

// MultiplyTime applies a constant multiple to the given time.
func MultiplyTime(mult float64) TimeFn {
	return func(t float64) float64 {
		return mult * t
	}
}

func TimeLooper(dur time.Duration) TimeFn {
	durT := float64(dur) / float64(time.Second)
	return func(t float64) float64 {
		return math.Mod(t, durT)
	}
}

// SinWave produces a sine wave function at the given frequency.
func SinWave(freq float64) WaveFn {
	return IntegrateWave(MultiplyTime(freq), BasicSinFn)
}

// SquareWave produces a square wave function at the given frequency.
func SquareWave(freq float64) WaveFn {
	return IntegrateWave(MultiplyTime(freq), BasicSquareFn)
}

// TriangleWave produces a triangle wave function at the given frequency.
func TriangleWave(freq float64) WaveFn {
	return IntegrateWave(MultiplyTime(freq), BasicTriangleFn)
}

// SawWave produces a sawtooth wave function at the given frequency.
func SawWave(freq float64) WaveFn {
	return IntegrateWave(MultiplyTime(freq), BasicSawFn)
}

// ZeroWave is a wave that is always zero (produces no sound).
func ZeroWave() WaveFn {
	return func(t float64) float64 {
		return 0
	}
}

func BasicSinFn(t float64) float64 {
	return math.Sin(2 * math.Pi * t)
}

func BasicSquareFn(t float64) float64 {
	if BasicSinFn(t) > 0 {
		return 1
	}
	return -1
}

func BasicTriangleFn(t float64) float64 {
	return math.Abs(4*(t-math.Floor(t)-0.5)) - 1
}

func BasicSawFn(t float64) float64 {
	return 2 * (t - math.Floor(t) - 0.5)
}
