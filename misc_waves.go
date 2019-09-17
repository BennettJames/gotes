package gotes

import (
	"math"
	"time"
)

func WaveFnSequence(
	step time.Duration,
	waves ...WaveFn,
) WaveFn {
	return func(t float64) float64 {
		// todo (bs): need to port linear fade to this so it won't pop. I actually
		// think it doesn't right now due to stability, but that's not a very stable
		// system.
		waveI := int(t / (float64(step) / float64(time.Second)))
		return waves[waveI%len(waves)](t)
	}
}

func WeirdWave1(freq float64) WaveFn {
	return func(t float64) float64 {
		return math.Sin(1.2*math.Pi*freq*t + math.Sin(2*math.Pi*freq*t))
	}
}

func WeirdWave2(freq float64) WaveFn {
	// note (bs): fundamentally, I think this is very smilar to a tonewheel organ
	// note. Let's see if I can figure out the internals for that, and perhaps
	// generalize this.
	fn := func(x, o float64) float64 {
		return math.Sin(2*math.Pi*x*freq + o)
	}

	return func(t float64) float64 {
		return fn(t, math.Pow(fn(t, 0), 2)+0.75*fn(t, 0.25)+0.1*fn(t, 0.5))
	}
}

func PianoWave(dur time.Duration, freq float64) WaveFn {
	fn := func(x, o float64) float64 {
		return math.Sin(2*math.Pi*x*freq + o)
	}

	// note (bs): rather than enforce duration > attack, I think it'd be better to
	// set attack to a shorter amount if duration is short; e.g. if it's <500ms
	// then shorten attack proportionally.
	attackT := 0.002
	durT := math.Max(attackT*2, float64(dur)/float64(time.Second))
	dampen := math.Pow(0.5*math.Log(freq*0.3), 2)

	return func(t float64) float64 {
		if t > durT {
			return 0
		}

		// ques (bs): would it make sense to separate this out some? The underlying
		// wave can be handled separately, as it is in "WeirdWave2". This could
		// completely ignore the underlying wave; and instead just apply an envelope
		// to it. That is sort of pleasant in that it highlights the different
		// elements - there is a pseudo-resonant wave that a note consists of via a
		// shaped amplitude-envelope.
		v := fn(t, math.Pow(fn(t, 0), 2)+0.75*fn(t, 0.25)+0.1*fn(t, 0.5))
		if t < attackT {
			return v * t / attackT
		}
		return v * math.Pow(1-(t-attackT)/(durT-attackT), dampen)
	}
}

func PeriodicSinWave(cycle time.Duration, f1, f2 float64) WaveFn {
	if f1 > f2 {
		f2, f1 = f1, f2
	}
	return func(t float64) float64 {
		// ques (bs): is it possible to abstract this transition some to make it
		// more of a pure "time variation"? Let's play around some.
		d := float64(time.Second) / float64(cycle) * 2 * math.Pi
		base := (f1 + f2) / 2 * t
		variation := (f2 - f1) / (2 * d) * math.Cos(t*d)
		return math.Sin((base + variation) * 2 * math.Pi)
	}
}

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
