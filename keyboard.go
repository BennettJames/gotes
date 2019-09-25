package gotes

import (
	"sync"
	"time"
)

type Keyboard struct {
	sr  SampleRate
	dur time.Duration

	l sync.Mutex

	waves        []WaveFFn
	totalSamples int
}

func NewKeyboard(sr SampleRate, dur time.Duration) *Keyboard {
	return &Keyboard{
		sr:  sr,
		dur: dur,
	}
}

func (g *Keyboard) Stream(samples [][2]float64) {
	g.l.Lock()
	defer g.l.Unlock()

	for i := range samples {
		t := float64(g.totalSamples) / float64(g.sr)
		v := float64(0)

		for _, w := range g.waves {
			waveV, _ := w(t)
			v += waveV
		}
		samples[i][0], samples[i][1] = v, v
		g.totalSamples++
	}

	t := float64(g.totalSamples) / float64(g.sr)
	newWaves := []WaveFFn{}
	for _, w := range g.waves {
		if _, done := w(t); !done {
			newWaves = append(newWaves, w)
		}
	}
	g.waves = newWaves
}

func (g *Keyboard) Add(freq float64) {
	g.l.Lock()
	defer g.l.Unlock()

	baseT := float64(g.totalSamples) / float64(g.sr)
	durT := float64(g.dur) / float64(time.Second)

	// todo (bs): this fixed gain is pretty clumsy. It acts as something of a
	// safeguard to ensure that
	w := AmplifyWave(Gain(0.4), PianoWave(g.dur, freq))

	g.waves = append(g.waves, func(t float64) (float64, bool) {
		if t > baseT+durT {
			return 0, true
		}
		return w(t - baseT), false
	})
}
