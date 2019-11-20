package gotes

// CombineWaves returns a simple aggregated waveform of all the given waves.
func CombineWaves(fns ...WaveFn) WaveFn {
	if len(fns) == 0 {
		return ZeroWave()
	}
	return func(t float64) float64 {
		total := 0.0
		for _, fn := range fns {
			total += fn(t)
		}
		return total
	}
}

// IntegrateWave creates a wave function that will take the initial time, pass
// it to the given time function, and pass that result to the wave function.
func IntegrateWave(tFn TimeFn, wFn WaveFn) WaveFn {
	return func(t float64) float64 {
		return wFn(tFn(t))
	}
}

// AmplifyWave will amplify the given wave function at each time point by the
// value yielded by the amplify function.
func AmplifyWave(aFn AmpFn, wFn WaveFn) WaveFn {
	return func(t float64) float64 {
		return wFn(t) * aFn(t)
	}
}
