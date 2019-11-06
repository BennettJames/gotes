package gotes

type (
	// WaveFn is a type that represents a basic sound wave. Given time values
	// (where 1.0 == 1 second), this maps to wave amplitude in [-1, 1].
	WaveFn func(t float64) (v float64)

	// FiniteWaveFn is a "finite" WaveFn. It will fade out and stop producing any
	// nonzero values past a certain t, at which point it will not produce any
	// values for any greater t values.
	FiniteWaveFn func(t float64) (v float64, done bool)

	// TimeFn maps the given time to another time. This is useful for certain wave
	// transformations - for instance, the frequency of a given wave can be modified
	// by applying a constant multiplier to time.
	//
	// Note that while this has the same signature of WaveFn, it's implied contract
	// is different. Rather than being a periodic function that maps to [-1,1], this
	// will instead map t to another t in a monotonically increasing fashion.
	TimeFn func(t float64) float64

	// AmpFn yields a multiple for the base wave at the given time.
	AmpFn func(t float64) float64
)
