package gotes

import (
	"math"
)

func SinWave(freq float64) WaveFn {
	return func(t float64) float64 {
		return math.Sin(2 * math.Pi * freq * t)
	}
}

func SquareWave(freq float64) WaveFn {
	w := SinWave(freq)
	return func(t float64) float64 {
		v := w(t)
		if v > 0 {
			return 1
		}
		return -1
	}
}

func TriangleWave(freq float64) WaveFn {
	return func(t float64) float64 {
		x := freq * t
		return math.Abs(4*(x-math.Floor(x)-0.5)) - 1
	}
}

func SawWave(freq float64) WaveFn {
	return func(t float64) float64 {
		x := freq * t
		return 2 * (x - math.Floor(x) - 0.5)
	}
}

func ZeroWave() WaveFn {
	return func(t float64) float64 {
		return 0
	}
}
