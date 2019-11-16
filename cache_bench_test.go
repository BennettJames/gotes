package gotes

import (
	"fmt"
	"math"
	"testing"
	"time"
)

var _cacheBenchInt = int64(0)
var _cacheBenchFloat = float64(0)

func Benchmark_waveFns(b *testing.B) {
	const sr = float64(48_000)
	const stepFloat = float64(time.Second) / sr

	b.Run("attackAndDecay", func(b *testing.B) {
		fn := AttackAndDecay(2.0, 6.0)
		for i := 0; i < b.N; i++ {
			t := float64(i) * stepFloat
			t = t - float64(int(t))
			_cacheBenchFloat = fn(t)
		}
	})

	b.Run("pseudoCacheAttackAndDecay", func(b *testing.B) {
		fn := AttackAndDecay(2.0, 6.0)
		fn = AmpFn(InterpolateCache(WaveFn(fn), 1024))
		for i := 0; i < b.N; i++ {
			t := float64(i) * stepFloat
			t = t - float64(int(t))
			_cacheBenchFloat = fn(t)
		}
	})
}

func Benchmark_cacheLookups(b *testing.B) {
	const freq = NoteA3
	const cacheSize = 512
	const sr = float64(48_000)
	const stepFloat = float64(time.Second) / sr

	fn := BasicSinFn

	// This benchmark test several different ways of performing lookups on
	// interpolated caches. The best-performing is the one that's actually used.

	b.Run("interpolate", func(b *testing.B) {
		// This tests whether injecting the size versus looking it up incurs
		// different costs.
		//
		// Conclusion: injection doesn't add anything.
		b.Run("injectedSize", func(b *testing.B) {
			b.Run("explicitSize", func(b *testing.B) {
				cache := benchBuildCache(fn, cacheSize)
				for i := 0; i < b.N; i++ {
					t := float64(i) * stepFloat
					_cacheBenchFloat = interCacheLookupOpt1(cache, cacheSize, t)
				}
			})

			b.Run("calculatedSize", func(b *testing.B) {
				cache := benchBuildCache(fn, cacheSize)
				for i := 0; i < b.N; i++ {
					t := float64(i) * stepFloat
					_cacheBenchFloat = interCacheLookupOpt2(cache, t)
				}
			})
		})

		// This tests an option with and without safe mapping for t values less than
		// zero.
		//
		// Conclusion: negative handling adds ~0.4ns of overhead; which is worth it.
		b.Run("negativeGuard", func(b *testing.B) {
			b.Run("noNegativeHandling", func(b *testing.B) {
				cache := benchBuildCache(fn, cacheSize)
				for i := 0; i < b.N; i++ {
					t := float64(i) * stepFloat
					_cacheBenchFloat = interCacheLookupOpt2(cache, t)
				}
			})

			b.Run("negativeHandling", func(b *testing.B) {
				cache := benchBuildCache(fn, cacheSize)
				for i := 0; i < b.N; i++ {
					t := float64(i) * stepFloat
					_cacheBenchFloat = interCacheLookupOpt3(cache, t)
				}
			})
		})

		// This tests a few different approaches for mapping t to the cache domain.
		// The main one in use, "normT", uses some float->int tricks to map t to a
		// value in [0, 1). As all wave caching operates in that time window, that's
		// a safe way to quickly convert the value.
		//
		// Conclusion: normT is the best. float.Mod is absolutely terrible. int
		// modding isn't atrocious; but unless it's by a constant value it's not
		// great either.
		b.Run("floatModulus", func(b *testing.B) {
			b.Run("normalizedT", func(b *testing.B) {
				cache := benchBuildCache(fn, cacheSize)
				for i := 0; i < b.N; i++ {
					t := float64(i) * stepFloat
					_cacheBenchFloat = interCacheLookupOpt1(cache, cacheSize, t)
				}
			})

			b.Run("intModulus", func(b *testing.B) {
				cache := benchBuildCache(fn, cacheSize)
				for i := 0; i < b.N; i++ {
					t := float64(i) * stepFloat
					_cacheBenchFloat = interCacheLookupOpt4(cache, t)
				}
			})

			b.Run("floatModulus", func(b *testing.B) {
				cache := benchBuildCache(fn, cacheSize)
				for i := 0; i < b.N; i++ {
					t := float64(i) * stepFloat
					_cacheBenchFloat = interCacheLookupOpt5(cache, t)
				}
			})
		})

		// This tests the usage of a float64-based cache and an int16 based cache. In
		// practice, the precision of the sound is limited to 16 bits per sample, so
		// most of the space occupied by a float64 represents unusable precision.
		//
		// Conclusion: int16's require a fair bit more horsepower; 4.5ns vs 7.2ns.
		// Might be worth it though; as it's only 1/4 of the space.
		//
		// Question: would it be possible to "cheat" here some in terms of casting and
		// storage as to make this more efficient?
		b.Run("floatVsInt16", func(b *testing.B) {
			b.Run("float", func(b *testing.B) {
				cache := benchBuildCache(fn, cacheSize)
				for i := 0; i < b.N; i++ {
					t := float64(i) * stepFloat
					_cacheBenchFloat = interCacheLookupOpt3(cache, t)
				}
			})

			b.Run("int16", func(b *testing.B) {
				cache := benchBuildInt16Cache(fn, cacheSize)
				for i := 0; i < b.N; i++ {
					t := float64(i) * stepFloat
					_cacheBenchFloat = interCacheLookupOpt6(cache, t)
				}
			})
		})
	})

	// This performs a few tests on fixed cache. For the most part, this behaves
	// the behavior as interpolation, but a little simpler. This tries a few
	// options before settling on a similar option.
	b.Run("fixed", func(b *testing.B) {

		// opt1 takes a size argument, and uses this to map t to the cache.
		b.Run("opt1", func(b *testing.B) {
			cache := benchBuildCache(fn, cacheSize)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = fixedCacheLookupOpt1(cache, cacheSize, t)
			}
		})

		// opt2 calculates size internally, and adjusts t to fit inside of [0, 1).
		// This performs a little worse than opt1, at the benefit of working for
		// negative t and not requiring an argument.
		b.Run("opt2", func(b *testing.B) {
			cache := benchBuildCache(fn, cacheSize)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = fixedCacheLookupOpt2(cache, t)
			}
		})

		// opt3 calculates it's own size argument, and uses it to mod the index.
		// That drastically reduces performance (12ns/op); most likely due to broken
		// inlining.
		b.Run("opt3", func(b *testing.B) {
			cache := benchBuildCache(fn, cacheSize)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = fixedCacheLookupOpt3(cache, t)
			}
		})
	})

	// This compares the overall performance/safety winners of the fixed cache to
	// the interpolated cache. Note that this is just performance; on average the
	// interpolated cache is more accurate for the same number of samples.
	//
	// Conclusion: fixed inline cache is best. It costs 3ns at base. Adding a
	// wrapper adds 2.5ns and interpolation adds 1.5ns.
	b.Run("fixedVsInterpolate", func(b *testing.B) {

		b.Run("bestFixedInline", func(b *testing.B) {
			cache := benchBuildCache(fn, cacheSize)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = fixedCacheLookupOpt2(cache, t)
			}
		})

		b.Run("bestFixedWrapper", func(b *testing.B) {
			cacheFn := fixedCacheWrapper(fn, cacheSize)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = cacheFn(t)
			}
		})

		b.Run("bestInterpolatedInline", func(b *testing.B) {
			cache := benchBuildCache(fn, cacheSize)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = interCacheLookupOpt3(cache, t)
			}
		})

		b.Run("bestInterpolatedWrapper", func(b *testing.B) {
			cacheFn := interCacheWrapper(fn, cacheSize)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = cacheFn(t)
			}
		})
	})
}

func Test_newRMSE(t *testing.T) {
	const freq = NoteA3
	const cacheSize = 128
	const sr = float64(48_000)
	const stepFloat = float64(time.Second) / sr

	fn := BasicSinFn

	cache := benchBuildCache(fn, cacheSize)

	rmse1 := CalcWaveRMSE(
		fn,
		func(t float64) float64 {
			return interCacheLookupOpt3(cache, t)
		},
		int(sr),
	)
	rmse2 := CalcWaveRMSE(
		fn,
		func(t float64) float64 {
			return fixedCacheLookupOpt2(cache, t)
		},
		int(sr),
	)
	fmt.Println("@@@ error rate is: ", rmse1, rmse2)

}

// benchBuildCache creates a cache of the specific size from the wave function. There
// is one extra value on the end that's equal to the first; that makes wrapping
// interpolation easier and less error prone.
func benchBuildCache(
	fn WaveFn,
	size int,
) []float64 {
	cache := make([]float64, size+1)
	for i := 0; i < size; i++ {
		t := float64(i) / float64(size)
		cache[i] = fn(t)
	}
	cache[size] = cache[0] // wraparound for easier lookups
	return cache
}

// benchBuildInt16Cache creates a cache that uses 16-bit integers rather than
// 64-bit floats. This has the same technical amount of clarity with one fourth
// of the space dedicated to caching.
func benchBuildInt16Cache(
	fn WaveFn,
	size int,
) []int16 {
	cache := make([]int16, size+1)
	for i := 0; i < size; i++ {
		t := float64(i) / float64(size)
		cache[i] = int16(fn(t) * math.MaxInt16)
	}
	cache[size] = cache[0] // wraparound for easier lookups
	return cache
}

// interCacheLookupOpt1 performs the lookup using an explicit size argument.
func interCacheLookupOpt1(
	cache []float64,
	size int,
	t float64,
) float64 {
	t = t - float64(int(t))
	floatI := t * float64(size)
	i := int(floatI)
	amt := floatI - float64(i)
	return cache[i]*(1-amt) + cache[i+1]*(amt)
}

// interCacheLookupOpt2 is as #1, but with negative protection for t.
func interCacheLookupOpt2(
	cache []float64,
	t float64,
) float64 {
	size := len(cache) - 1
	t = t - float64(int(t))
	floatI := t * float64(size)
	i := int(floatI)
	amt := floatI - float64(i)
	return cache[i]*(1-amt) + cache[i+1]*(amt)
}

// interCacheLookupOpt3 is as #2, but with inline size calculation.
func interCacheLookupOpt3(
	cache []float64,
	t float64,
) float64 {
	size := len(cache) - 1
	t = benchNormT(t)
	floatI := t * float64(size)
	i := int(floatI)
	amt := floatI - float64(i)
	return cache[i]*(1-amt) + cache[i+1]*(amt)
}

// interCacheLookupOpt4 perfoms it's own size calculation and determines values
// using integer modulo operations.
func interCacheLookupOpt4(
	cache []float64,
	t float64,
) float64 {
	size := len(cache) - 1
	fs := float64(size)
	i := int(t*fs) % size
	amt := float64(i) / fs
	return cache[i]*(1-amt) + cache[i+1]*(amt)
}

// interCacheLookupOpt5 performs it's own size calculation and determines values
// using float modulo operations.
func interCacheLookupOpt5(
	cache []float64,
	t float64,
) float64 {
	size := len(cache) - 1
	t = math.Mod(t, 1)
	floatI := t * float64(size)
	i := int(floatI)
	amt := floatI - float64(i)
	return cache[i]*(1-amt) + cache[i+1]*(amt)
}

// interCacheLookupOpt6 makes use of 16-bit cache size.
func interCacheLookupOpt6(
	cache []int16,
	t float64,
) float64 {
	const intInv float64 = 1.0 / math.MaxInt16
	size := len(cache) - 1
	t = benchNormT(t)
	floatI := t * float64(size)
	i := int(floatI)
	amt := floatI - float64(i)
	return float64(cache[i])*(1-amt)*intInv + float64(cache[i+1])*(amt)*intInv
}

// fixedCacheLookupOpt1 does a basic resolution of a value without interpolation
// and with explicit size.
func fixedCacheLookupOpt1(
	cache []float64,
	size int,
	t float64,
) float64 {
	i := int(t*float64(size)) % size
	return cache[i]
}

// fixedCacheLookupOpt2 performs fixed lookup in the cache using time normalization.
func fixedCacheLookupOpt2(
	cache []float64,
	t float64,
) float64 {
	size := len(cache) - 1
	t = benchNormT(t)
	i := int(t * float64(size))
	return cache[i]
}

// fixedCacheLookupOpt3 performs fixed lookup in the cache using int modulus.
func fixedCacheLookupOpt3(
	cache []float64,
	t float64,
) float64 {
	size := len(cache) - 1
	i := int(t*float64(size)) % size
	return cache[i]
}

// fixedCacheWrapper builds a fixed cache of the given size for the wave fn,
// and returns a new wave fn that uses the cache.
func fixedCacheWrapper(
	fn WaveFn,
	size int,
) WaveFn {
	cache := benchBuildCache(fn, size)
	return func(t float64) float64 {
		t = benchNormT(t)
		i := int(t * float64(size))
		return cache[i]
	}
}

// interCacheWrapper builds an interpolation cache of the given size for
// the wave fn, and returns a new wave fn that uses the cache.
func interCacheWrapper(
	fn WaveFn,
	size int,
) WaveFn {
	cache := make([]float64, size)
	size--
	for i := 0; i < size; i++ {
		t := float64(i) / float64(size)
		cache[i] = fn(t)
	}
	cache[size] = cache[0] // wraparound for easier lookups
	return func(t float64) float64 {
		t = benchNormT(t)
		floatI := t * float64(size)
		i := int(floatI)
		amt := floatI - float64(i)
		return cache[i]*(1-amt) + cache[i+1]*(amt)
	}
}

// benchNormT maps all t values to range [0, 1); e.g. 0.5->0.5, 35.2->0.2,
// -0.1->0.9.
func benchNormT(t float64) float64 {
	t = t - float64(int(t))
	if t < 0 {
		return 1 + t
	}
	return t
}

func fixedCacheToPrecision(
	fn WaveFn,
	sr SampleRate,
	errRate float64,
) WaveFn {
	startPow := int(math.Ceil(math.Log(float64(sr)) / math.Log(2)))
	measureSR := SampleRate(int(math.Pow(2, float64(startPow))))
	lastCacheFn := fixedCacheWrapper(fn, int(measureSR))

	for i := startPow; i > 0; i-- {
		cacheFn := interCacheWrapper(fn, int(math.Pow(2, float64(i))))
		cacheRMSE := benchGetRMSE(fn, cacheFn, int(measureSR))
		if cacheRMSE > errRate {
			return lastCacheFn
		}
		lastCacheFn = cacheFn
	}

	// probably shouldn't be reached; let's think about when it would be
	return lastCacheFn
}

func benchGetRMSE(
	actualFn, approxFn WaveFn,
	numSamples int,
) float64 {
	errSum := 0.0
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		errSum += math.Pow(actualFn(t)-approxFn(t), 2)
	}
	return math.Sqrt(errSum / float64(numSamples))
}

func Benchmark_Oscillate(b *testing.B) {

	const sr = float64(48_000)
	const stepFloat = float64(time.Second) / sr
	accel, period := 2.5, 0.3

	fmt.Println("@@@ rmse is", CalcWaveRMSE(
		oscillateTimeOpt1(accel, period),
		OscillateTime(accel, period),
		48_000,
	))

	oscWithTime := func(peakAccel, period float64) TimeFn {
		return func(t float64) float64 {
			return (1+peakAccel/2)*t - peakAccel/(4*math.Pi*period)*math.Sin(2*math.Pi*period*t)
		}
	}

	wave1 := func(freq float64) WaveFn {
		return AmplifyWave(
			FixedAmplify(0.5),
			IntegrateWave(
				oscWithTime(2.0, 0.2),
				IntegrateWave(
					MultiplyTime(freq),
					BasicSinFn,
				),
			),
		)
	}

	wave2 := func(freq float64) WaveFn {
		return AmplifyWave(
			FixedAmplify(0.5),
			IntegrateWave(
				OscillateTime(2.0, 0.2),
				IntegrateWave(
					MultiplyTime(freq),
					BasicSinFn,
				),
			),
		)
	}

	// alright; so the caching doesn't scale with frequency. Is that a fault of
	// the equations; or a relic of the

	fmt.Println("@@@ rmse @440hz is", CalcWaveRMSE(
		wave1(440),
		wave2(440),
		48_000,
	))

	fmt.Println("@@@ rmse @1hz is", CalcWaveRMSE(
		wave1(1),
		wave2(1),
		48_000,
	))

	b.Skip()

	// This tests the speed of directly calculating the oscillation each time vs
	// using a cached lookup of time.
	//
	// Conclusion: cached is, unsurprisingly, much faster - 8ns vs uncached ns.
	b.Run("cacheVsManual", func(b *testing.B) {
		b.Run("uncached", func(b *testing.B) {
			fn := oscillateTimeOpt1(accel, period)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = fn(t)
			}
		})

		b.Run("cached", func(b *testing.B) {
			fn := oscillateTimeOpt2(accel, period)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = fn(t)
			}
		})
	})

	// This hoists certain unchanging calculations up to be calculated once.
	//
	// Conclusion: hoisting the fixed multipliers cuts another 1ns off; down to
	// 7.5ns.
	b.Run("hoistedConsts", func(b *testing.B) {
		b.Run("inlineCalc", func(b *testing.B) {
			fn := oscillateTimeOpt2(accel, period)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = fn(t)
			}
		})

		b.Run("hoistedConsts", func(b *testing.B) {
			fn := oscillateTimeOpt3(accel, period)
			for i := 0; i < b.N; i++ {
				t := float64(i) * stepFloat
				_cacheBenchFloat = fn(t)
			}
		})
	})
}

func oscillateTimeOpt1(peakAccel, period float64) TimeFn {
	return func(t float64) float64 {
		return (1+peakAccel/2)*t - peakAccel/(4*math.Pi*period)*math.Sin(2*math.Pi*period*t)
	}
}

func oscillateTimeOpt2(peakAccel, period float64) TimeFn {
	return func(t float64) float64 {
		return (1+(peakAccel/2))*t - peakAccel/(4*math.Pi*period)*lookupSin(period*t)
	}
}

func oscillateTimeOpt3(peakAccel, period float64) TimeFn {
	c1 := 1 + (peakAccel / 2)
	c2 := peakAccel / (4 * math.Pi * period)
	return func(t float64) float64 {
		return c1*t - c2*lookupSin(period*t)
	}
}
