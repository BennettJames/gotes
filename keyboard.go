package gotes

import (
	"sync"
	"time"
)

// Kyeboard is an object that can be used to dynamically play piano-style notes
// on over time.
type Keyboard struct {
	sr  SampleRate
	dur time.Duration

	l sync.Mutex

	waves        []FiniteWaveFn
	totalSamples int
}

// NewKeyboard initializes a new keyboard with the given sample rate and note
// duration.
func NewKeyboard(sr SampleRate, dur time.Duration) *Keyboard {
	return &Keyboard{
		sr:  sr,
		dur: dur,
	}
}

// Stream implements the Streamer interface by providing the next set of samples.
func (g *Keyboard) Stream(samples []float64) {
	g.l.Lock()
	defer g.l.Unlock()

	for i := range samples {
		t := float64(g.totalSamples) / float64(g.sr)
		v := float64(0)

		for _, w := range g.waves {
			waveV, _ := w(t)
			v += waveV
		}
		samples[i] = v
		g.totalSamples++
	}

	t := float64(g.totalSamples) / float64(g.sr)
	newWaves := []FiniteWaveFn{}
	for _, w := range g.waves {
		if _, done := w(t); !done {
			newWaves = append(newWaves, w)
		}
	}
	g.waves = newWaves
}

// Add will create a note of the given frequency and add it to the keyboard.
func (g *Keyboard) Add(freq float64) {
	g.l.Lock()
	defer g.l.Unlock()

	baseT := float64(g.totalSamples) / float64(g.sr)
	durT := float64(g.dur) / float64(time.Second)

	// todo (bs): this fixed gain is pretty clumsy. It acts as something of a
	// safeguard to ensure that multiple notes can be played at the same time
	// without overwhelming the volume. I'd kinda guess this should be more
	// adaptive based on the number of concurrent notes - e.g. make a mapping
	// like: 1 note -> 0.4 gain; 2 notes -> 0.35 gain each; 3 notes -> 0.28 gain
	// each; 4 notes -> 0.25 gain each
	//
	// and have further notes have a fixed fraction of 1. This would need some
	// good internal smarts about how to downscale past gains for existing notes;
	// I'd say it'd require some better struct-based functions to make variability
	// easier to manage.
	// w := AmplifyWave(Gain(0.4), PianoNote(g.dur, freq))

	w := AmplifyWave(
		Gain(0.4),
		PianoNote(g.dur, freq),
	)

	g.waves = append(g.waves, func(t float64) (float64, bool) {
		if t > baseT+durT {
			return 0, true
		}
		return w(t - baseT), false
	})
}
