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
