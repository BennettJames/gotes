package gotes

import (
	"math"
	"time"
)

// PianoNote creates a piano-like note at the given frequency over the time
// period, while using some caching for efficiency.
func PianoNote(dur time.Duration, freq float64) WaveFn {
	durT := float64(dur) / float64(time.Second)
	dampen := math.Pow(0.5*math.Log(freq*0.3), 2)
	return AmplifyWave(
		AttackAndDecay(durT, dampen),
		PianoWave(freq),
	)
}

// PianoWave creates a piano-wave structure at the given frequency.
func PianoWave(freq float64) WaveFn {
	return IntegrateWave(
		MultiplyTime(freq),
		lookupPiano,
	)
}

// BasicPianoFn creates a waveform with a degree of internal resonance that
// can be shaped to sound somewhat piano-like.
func BasicPianoFn(t float64) float64 {
	fn := func(o float64) float64 {
		return math.Sin(2*math.Pi*t + o)
	}
	return fn(math.Pow(fn(0), 2) + 0.75*fn(0.25) + 0.1*fn(0.5))
}

// AttackAndDecay is an amplitude function that has an initial rapid
// gain phase (the "attack"), and a fadeout over the course of durT. dampen
// controls the rate of the fadeout, with higher values increasing the severity
// of it.
func AttackAndDecay(
	durT float64,
	dampen float64,
) AmpFn {
	attackT := 0.002
	durT = math.Max(attackT*2, durT)
	// The exponential decay function is highly dynamic based on parameters, which
	// makes it hard to have a shared cache. However, it's also fairly simple and
	// easy to approximate; so a small cache is perpared each time. With 16 samples,
	// this generally gets w/in 0.1% accuracy.
	cache := MakeCache(
		func(t float64) float64 {
			return math.Pow(1-t, dampen)
		},
		16,
	)
	cache[len(cache)-1] = 0 // clumsy fix for cache wraparound
	invAttack := 1 / attackT
	invDur := 1 / (durT - attackT)
	return func(t float64) float64 {
		if t < 0 {
			return 0
		} else if t < attackT {
			return t * invAttack
		} else if t < durT {
			return CacheInterpolateLookup(cache, (t-attackT)*invDur)
		}
		return 0
	}
}

// WeirdPianoWave produces a sound akin to piano note, but with a different
// waveform.
func WeirdPianoWave(dur time.Duration, freq float64) WaveFn {
	durT := float64(dur) / float64(time.Second)
	dampen := math.Pow(0.5*math.Log(freq*0.3), 2)
	return AmplifyWave(
		AttackAndDecay(durT, dampen),
		IntegrateWave(
			MultiplyTime(freq),
			WeirdWave1,
		),
	)
}

// WeirdWave1 is a random attempt at an alternate waveform.
func WeirdWave1(t float64) float64 {
	return math.Sin(2*math.Pi*t + math.Sin(3.4*math.Pi*t))
}

// PeriodicSinWave cycles between the two given frequencies in each cycle.
func PeriodicSinWave(cycle time.Duration, f1, f2 float64) WaveFn {
	if f1 > f2 {
		f2, f1 = f1, f2
	}
	return IntegrateWave(
		SinTime(f1, f2, float64(time.Second)/float64(cycle)),
		BasicSinFn,
	)
}

// OscillateTime creates a time function that varies how "fast" it moves.
// "peakAccel" is the peak rate it reached (e.g. peakAccel of 1.0 -> 100% as
// fast as normal; 0.0 -> constant rate); period is the time between
func OscillateTime(peakAccel, period float64) TimeFn {
	c1 := 1 + (peakAccel / 2)
	c2 := peakAccel / (4 * math.Pi * period)
	return func(t float64) float64 {
		return c1*t - c2*CacheInterpolateLookup(sinWaveCache, period*t)
	}
}

// SinTime is a time function that varies the rate of change according to a sin
// wave. The rate varies between low and high in each period.
func SinTime(low, high, period float64) TimeFn {
	// note (bs): so, this is a little better. Still a little curious if it could
	// be better. Perhaps not.
	d := 2 * math.Pi * period
	return func(t float64) float64 {
		return (low+high)/2*t + (high-low)/(2*d)*math.Sin(d*t)
	}
}
