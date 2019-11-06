package gotes

import (
	"math"
)

// FixedAmplify applies a constant
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
