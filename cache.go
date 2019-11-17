package gotes

import "math"

const (
	// DefaultCacheSize is a good cache size to use when no tuning has been
	// performed.
	DefaultCacheSize = 512
)

// Cache will generate a cache for the given function in [0,1). It will
// generate a set of samples from the given function, store them, and return
// a function that calculates values from the cache rather than the underlying
// function.
//
// This uses an interpolated cache with the default cache size of 512; which
// is a good option for most use cases.
func Cache(fn WaveFn) WaveFn {
	return InterpolateCache(fn, DefaultCacheSize)
}

// DirectCache creates a function that will cache via direct lookups. For each
// value in the returned wave function, it finds the closest value and returns
// it.
func DirectCache(fn WaveFn, size int) WaveFn {
	cache := MakeCache(fn, size)
	return func(t float64) float64 {
		size := len(cache) - 1
		t = normT(t)
		i := int(t * float64(size))
		return cache[i]
	}
}

// InterpolateCache creates a function that will perform interpolated cache
// lookups. For each value passed to the returned wave function, it finds the
// to closest values and averages between them.
func InterpolateCache(fn WaveFn, size int) WaveFn {
	cache := MakeCache(fn, size)
	return func(t float64) float64 {
		size := len(cache) - 1
		t = normT(t)
		floatI := t * float64(size)
		i := int(floatI)
		amt := floatI - float64(i)
		return cache[i]*(1-amt) + cache[i+1]*(amt)
	}
}

// MakeCache creates a cached set of values of the given size. The final
// value will always wrap around to the first, which makes for easier lookups
// in certain situations.
func MakeCache(fn func(float64) float64, size int) []float64 {
	cache := make([]float64, size)
	size-- // the last index is reserved for easy wraparound calculations
	for i := 0; i < size; i++ {
		t := float64(i) / float64(size)
		cache[i] = fn(t)
	}
	cache[size] = cache[0]
	return cache
}

// CacheInterpolateLookup performs cache lookups by finding the closes value
// for t and returning it.
func CacheInterpolateLookup(cache []float64, t float64) float64 {
	size := len(cache) - 1
	t = normT(t)
	floatI := t * float64(size)
	i := int(floatI)
	amt := floatI - float64(i)
	return cache[i]*(1-amt) + cache[i+1]*(amt)
}

// CacheDirectLookup performs cache lookups by finding the two closest values
// to t and returning an average of them.
func CacheDirectLookup(cache []float64, t float64) float64 {
	size := len(cache) - 1
	t = normT(t)
	i := int(t * float64(size))
	return cache[i]
}

// CalcWaveRMSE returns the root mean square error of the difference between the
// two errors when at the given sample rate. That can be useful when tuning cache
// size.
func CalcWaveRMSE(
	actualFn, approxFn func(float64) float64,
	numSamples int,
) float64 {
	errSum := 0.0
	for i := 0; i < numSamples*2; i++ {
		t := float64(i) / float64(numSamples)
		errSum += math.Pow(actualFn(t)-approxFn(t), 2)
	}
	return math.Sqrt(errSum / float64(numSamples))
}

// normT maps all t values to range [0, 1); e.g. 0.5->0.5, 35.2->0.2, -0.1->0.9.
func normT(t float64) float64 {
	t = t - float64(int(t))
	if t < 0 {
		return 1 + t
	}
	return t
}
