package gotes

import (
	"math"
)

// FixedAmplify creates a fixed amplifier function that will multiply a wave fn
// by a constant.
//
// Note that this is a direct constant; more often than not `Gain` should
// be used; which adjusts for decibels.
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

// SinWave produces a sine wave function at the given frequency.
func SinWave(freq float64) WaveFn {
	return IntegrateWave(MultiplyTime(freq), lookupSin)
}

// SquareWave produces a square wave function at the given frequency.
func SquareWave(freq float64) WaveFn {
	return IntegrateWave(MultiplyTime(freq), lookupSquare)
}

// TriangleWave produces a triangle wave function at the given frequency.
func TriangleWave(freq float64) WaveFn {
	return IntegrateWave(MultiplyTime(freq), lookupTriangle)
}

// SawWave produces a sawtooth wave function at the given frequency.
func SawWave(freq float64) WaveFn {
	return IntegrateWave(MultiplyTime(freq), lookupSaw)
}

// ZeroWave is a wave that is always zero (produces no sound).
func ZeroWave() WaveFn {
	return func(t float64) float64 {
		return 0
	}
}

// BasicSinFn is a simple WaveFn that produces a 1hz sin wave.
func BasicSinFn(t float64) float64 {
	return math.Sin(2 * math.Pi * t)
}

// BasicSquareFn is a simple WaveFn that produces a 1hz square wave.
func BasicSquareFn(t float64) float64 {
	if BasicSinFn(t) > 0 {
		return 1
	}
	return -1
}

// BasicTriangleFn is a simple WaveFn that produces a 1hz triangle wave.
func BasicTriangleFn(t float64) float64 {
	return math.Abs(4*(t-math.Floor(t)-0.5)) - 1
}

// BasicSawFn is a simple WaveFn that produces a 1hz saw wave.
func BasicSawFn(t float64) float64 {
	return 2 * (t - math.Floor(t) - 0.5)
}
