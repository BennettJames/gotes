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

func Gain(amt float64, fn WaveFn) WaveFn {
	g := math.Pow(2, math.Max(0, amt)) - 1
	return func(t float64) float64 {
		return fn(t) * g
	}
}

func OffsetWave(tOff float64, fn WaveFn) WaveFn {
	return func(t float64) float64 {
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
		sr:    sr,
		fn:    fn,
		gainV: math.Max(0, initGain),
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
