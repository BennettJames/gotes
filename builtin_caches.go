package gotes

// This file contains some internal cache values that are used in wave
// functions. Generally, they are all caches of "pure wave functions"
// that are used with frequency to optimize performance.
var (
	pianoWaveCache    = MakeCache(BasicPianoFn, DefaultCacheSize)
	sinWaveCache      = MakeCache(BasicSinFn, DefaultCacheSize)
	squareWaveCache   = MakeCache(BasicSquareFn, DefaultCacheSize)
	triangleWaveCache = MakeCache(BasicTriangleFn, DefaultCacheSize)
	sawWaveCache      = MakeCache(BasicSawFn, DefaultCacheSize)

	// so - can this be cached; really? Maybe not. Obviously for fixed decay it
	// can; but I think it'd be hard to avoid the need
	//
	// It seems *conceivable* that you can do some sort of power-hack to make it a
	// little better; In practice the values are all in [0, 4]. let's say you did
	// this:
	//
	// create a cache of
	//
	// actually scratch that. I was about to suggest this: create a cace of
	// different square values. But that's not constant around different powers or
	// values. Hrm.
	//
	// I'd say this: the nature of the decay doesn't need to be hyper-precise.
	//
	// Another option that I don't love would be to just precalculate a few values
	// on a by-case basis. Like; let's say I generate 64 values an interpolate
	// them. How well would that work?
	//
	// Hmm, so a
	//
	//
	decayCache = MakeCache(
		decayFn(1.0, 1.0),
		DefaultCacheSize,
	)
)

func lookupPiano(t float64) float64 {
	return CacheInterpolateLookup(pianoWaveCache, t)
}

func lookupSin(t float64) float64 {
	return CacheDirectLookup(sinWaveCache, t)
}

func lookupSquare(t float64) float64 {
	return CacheDirectLookup(squareWaveCache, t)
}

func lookupTriangle(t float64) float64 {
	return CacheDirectLookup(triangleWaveCache, t)
}

func lookupSaw(t float64) float64 {
	return CacheDirectLookup(sawWaveCache, t)
}
