package gotes

type (
	WaveFn   func(t float64) float64
	BiWaveFn func(t float64) [2]float64

	// note (bs): for the next two values, I'm not convinced I shouldn't just
	// create an annotated function rather than the two

	streamerWave struct {
		sr           SampleRate
		fn           WaveFn
		totalSamples int
	}

	biStreamerWave struct {
		sr           SampleRate
		fn           BiWaveFn
		totalSamples int
	}
)

func (s *streamerWave) Stream(samples []float64) {
	for i := range samples {
		t := float64(s.totalSamples) / float64(s.sr)
		samples[i] = s.fn(t)
		s.totalSamples++
	}
}

func (s *biStreamerWave) Stream(samples [][2]float64) {
	// so - need to create an effective t value based on the total number of samples and the
	for i := range samples {
		t := float64(s.totalSamples) / float64(s.sr)
		samples[i] = s.fn(t)
		s.totalSamples++
	}
}

func BiWaveFromWave(fn WaveFn) BiWaveFn {
	return func(t float64) [2]float64 {
		v := fn(t)
		return [2]float64{v, v}
	}
}

func StreamerFromWave(sr SampleRate, fn WaveFn) Streamer {
	return &streamerWave{
		sr: sr,
		fn: fn,
	}
}

func BiStreamerFromWave(sr SampleRate, fn WaveFn) BiStreamer {
	return &biStreamerWave{
		sr: sr,
		fn: BiWaveFromWave(fn),
	}
}

func BiStreamerFromBiWave(sr SampleRate, fn BiWaveFn) BiStreamer {
	return &biStreamerWave{
		sr: sr,
		fn: fn,
	}
}
