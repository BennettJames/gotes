package gotes

import (
	"math"
	"time"
)

// In the course of creating gotes, many mistakes were made. Some were interesting
// enough to preserve; those are kept here for posterity.

// MistakenPeriodicSinWave is an interesting mistake that was made while
// creating PeriodicSinWave - has an interesting swoopiness to it.
func MistakenPeriodicSinWave(
	cycle time.Duration,
	f1, f2 float64,
) WaveFn {
	if f1 > f2 {
		f2, f1 = f1, f2
	}
	return func(t float64) float64 {
		d := float64(time.Second) / float64(cycle)
		base := (f1 + f2) / 2 * t
		variation := (f2 - f1) / (2 * d) * math.Cos(t*d*2*math.Pi)
		return math.Sin((base + variation) * 2 * math.Pi)
	}
}

// BadOscillateTime an incorrect time oscillator that sounds really interesting.
func BadOscillateTime(peakAccel, period float64) TimeFn {
	// todo (bs): let's see if I can inline cacheV here. "mistakes" like this
	// which are too dependent on sub-behavior like that can accidentally
	// disappear.
	return func(t float64) float64 {
		cacheV := CacheDirectLookup(sinWaveCache, period*t)
		return (1+(peakAccel/2))*t - peakAccel/(4*math.Pi*period)*cacheV
	}
}

// BadOscillateTime2 is another incorrect time oscillator that sounds really
// interesting.
func BadOscillateTime2(peakAccel, period float64) TimeFn {
	return func(t float64) float64 {
		cacheV := CacheDirectLookup(sinWaveCache, 2*math.Pi*period*t)
		return (1+(peakAccel/2))*t - peakAccel/(4*math.Pi*period)*cacheV
	}
}
