package gotes

type (
	// Streamer
	Streamer interface {
		Stream(samples []float64)
	}

	// BiStreamer
	BiStreamer interface {
		Stream(samples [][2]float64)
	}

	streamerWave struct {
		sr           SampleRate
		fn           WaveFn
		totalSamples int
	}

	biStreamerWave struct {
		sr           SampleRate
		fn           WaveFn
		totalSamples int
	}
)

// StreamerFromWave creates a streamer for the wave at the given sample rate.
func StreamerFromWave(fn WaveFn, sr SampleRate) Streamer {
	return &streamerWave{
		sr: sr,
		fn: fn,
	}
}

func (s *streamerWave) Stream(samples []float64) {
	for i := range samples {
		t := float64(s.totalSamples) / float64(s.sr)
		v := s.fn(t)
		samples[i] = v
		s.totalSamples++
	}
}

// BiStreamerFromWave creates a bi-streamer for the wave at the given sample rate.
func BiStreamerFromWave(fn WaveFn, sr SampleRate) BiStreamer {
	return &biStreamerWave{
		sr: sr,
		fn: fn,
	}
}

func (s *biStreamerWave) Stream(samples [][2]float64) {
	for i := range samples {
		t := float64(s.totalSamples) / float64(s.sr)
		v := s.fn(t)
		samples[i][0], samples[i][1] = v, v
		s.totalSamples++
	}
}
