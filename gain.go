package gotes

import (
	"math"
	"sync"
)

const (
	// gainTransitionT is the amount of time in seconds that it effectively takes
	// to transition between two gains.
	//
	// Technically this value should really be dynamic based on the size of the
	// gap, but that would be harder to do.
	gainTransitionT = 0.2
)

// Gain creates a constant amplier function by the given amount; but adjusted
// expontentially. This results in a more "natural" wave adjustment than a static
// amount.
func Gain(amt float64) AmpFn {
	g := math.Pow(2, math.Max(0, amt)) - 1
	return func(t float64) float64 {
		return g
	}
}

func ScaleGain(amt1, amt2 float64) AmpFn {
	amt1, amt2 = math.Max(0, amt1), math.Max(0, amt2)
	return func(t float64) float64 {
		// I don't think I like the function "getTransitionAmt". It's a more general
		// question of mapping a value to [0,1] w/ different formula. This is a
		// linear interpolation function that I need to give a period for (and
		// arguably a starting point, if you feel so inclined)
		amt := getTransitionAmt(t)
		return math.Pow(2, amt1*(1-amt)+amt2*amt) - 1
	}
}

func ScaleToLinear(startT, periodT float64) func(t float64) float64 {
	// so - here's an experiment to kind of generalize the type of function that
	// scales two values in a time period. Not sure I 100% agree with this API - I
	// sort of think that even if I did it, I'd be better as a more direct
	// interpolate fn like
	//
	//  type InterpolateFn func(t float64, v1, v2 float64) float64
	//
	// or maybe
	//
	//  type Interpolate func(sFn ScaleFn, w1, w2 WaveFn) WaveFn
	//
	//

	return func(t float64) float64 {
		if t <= startT {
			return 0
		} else if t >= startT+periodT {
			return 1
		}
		return (t - startT) / periodT
	}
}

// OffsetWave will play the wave function "starting" at the offset. That is:
// before the offset; it will be zero; afterwards; it will return the wave
// function offset by that amount.
func OffsetWave(tOff float64, fn WaveFn) WaveFn {
	return func(t float64) float64 {
		if t < tOff {
			return 0
		}
		return fn(t - tOff)
	}
}

func MoveGain(amt1, amt2 float64, fn WaveFn) WaveFn {
	amt1, amt2 = math.Max(0, amt1), math.Max(0, amt2)
	return func(t float64) float64 {
		tAmt := getTransitionAmt(t)
		return fn(t) * (amt1*(1-tAmt) + amt2*tAmt)
	}
}

type GainStreamer struct {
	l sync.Mutex

	sr SampleRate
	fn WaveFn

	totalSamples int

	gainV     float64
	lastGainV float64
	lastGainT float64
}

func NewGainStreamer(sr SampleRate, fn WaveFn, initGain float64) *GainStreamer {
	return &GainStreamer{
		sr:        sr,
		fn:        fn,
		gainV:     math.Max(0, initGain),
		lastGainV: math.Max(0, initGain),
	}
}

func (s *GainStreamer) Stream(samples [][2]float64) {
	s.l.Lock()
	defer s.l.Unlock()

	g1, g2 := s.lastGainV, s.gainV

	for i := range samples {
		t := float64(s.totalSamples) / float64(s.sr)
		tAmt := getTransitionAmt(t - s.lastGainT)
		v := s.fn(t) * (g1*(1-tAmt) + g2*tAmt)
		samples[i][0], samples[i][1] = v, v
		s.totalSamples++
	}
}

// SetGain will set the node to the new gain. Note that this is not
// instantaneous; and it will scale to the new gain value over several
// milliseconds.
func (s *GainStreamer) SetGain(gainV float64) {
	s.l.Lock()
	defer s.l.Unlock()

	t := float64(s.totalSamples) / float64(s.sr)
	tAmt := getTransitionAmt(t)
	s.lastGainV = s.gainV*tAmt + s.lastGainV*(1-tAmt)
	s.lastGainT = t
	s.gainV = math.Max(0, gainV)
}

// getTransitionAmt normalizes t to a value in [0, 1] based on how far within
// the transition range the time is.
func getTransitionAmt(t float64) float64 {
	t = math.Min(gainTransitionT, math.Max(0, t))
	return 1 - (gainTransitionT-t)/gainTransitionT
}
