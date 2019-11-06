package gotes

// note (bs): for the next two values, I'm not convinced I shouldn't just
// create an annotated function rather than the two types.

type streamerWave struct {
	sr           SampleRate
	fn           WaveFn
	totalSamples int
}

func (s *streamerWave) Stream(samples []float64) {
	for i := range samples {
		t := float64(s.totalSamples) / float64(s.sr)
		samples[i] = s.fn(t)
		s.totalSamples++
	}
}

func StreamerFromWave(sr SampleRate, fn WaveFn) Streamer {
	return &streamerWave{
		sr: sr,
		fn: fn,
	}
}
