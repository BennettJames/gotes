package gotes

// IntegrateWave will call the given wave function with the given mapped time
// function for any given time.
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
