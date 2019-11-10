package gotes

import (
	"sync"
	"time"
)

type Keyboard struct {
	sr  SampleRate
	dur time.Duration

	l sync.Mutex

	waves        []FiniteWaveFn
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
	newWaves := []FiniteWaveFn{}
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
		AmplifyWave(
			AttackAndDecay(2.0, 6.0),
			IntegrateWave(
				MultiplyTime(freq),
				cachePiano,
			),
		))

	g.waves = append(g.waves, func(t float64) (float64, bool) {
		if t > baseT+durT {
			return 0, true
		}
		return w(t - baseT), false
	})
}

var cachePiano = cacheWave(BasicPianoFn)

func cacheWave(fn WaveFn) WaveFn {
	cacheSize := 2048 // note (bs): may wish to make this configurable
	cache := make([]float64, cacheSize)
	for i := 0; i < cacheSize; i++ {
		t := float64(i) / float64(cacheSize)
		cache[i] = fn(t)
	}
	return func(t float64) float64 {
		return cache[int(t*float64(cacheSize))%cacheSize]
	}
}
