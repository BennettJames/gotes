package gotes

import (
	"math"
	"testing"
)

func Test_caching(t *testing.T) {
	// A note on these tests: they mostly just ensure that the caching
	// code more-or-less seems to be fitting together in a sane way. Caching
	// by it's nature is a little fuzzy, which resists precise testing.

	const freq = NoteA3
	const samples = 48_000
	noteFn := func(fn WaveFn) WaveFn {
		return IntegrateWave(MultiplyTime(freq), fn)
	}

	// sanityCheck does a basic check of the given fn to ensure the RMSE of the
	// cache approximation function is basically sane.
	sanityCheck := func(t *testing.T, approxFn WaveFn) {
		t.Run("rmse", func(t *testing.T) {
			rmse := CalcWaveRMSE(noteFn(BasicPianoFn), noteFn(approxFn), samples)
			if rmse >= 0.01 {
				t.Fatalf("Unexpectedly large RMSE: %f", rmse)
			}
			if rmse <= 0.0 {
				t.Fatalf("unexpectedly small RMSE: %f", rmse)
			}
		})
		t.Run("negativeValue", func(t *testing.T) {
			neg, pos := approxFn(-0.9), approxFn(0.1)
			if math.Abs(neg-pos) > 0.001 {
				t.Fatalf(
					"Cache negative values not adjusted correctly [neg=%f pos=%f]",
					neg, pos)
			}
		})
	}

	t.Run("Cache", func(t *testing.T) {
		c := Cache(BasicPianoFn)
		sanityCheck(t, c)
	})

	t.Run("DirectCache", func(t *testing.T) {
		c := DirectCache(BasicPianoFn, 1024)
		sanityCheck(t, c)
	})

	t.Run("InterpolateCache", func(t *testing.T) {
		c := InterpolateCache(BasicPianoFn, 1024)
		sanityCheck(t, c)
	})

	t.Run("CacheInterpolateLookup", func(t *testing.T) {
		cache := MakeCache(BasicPianoFn, DefaultCacheSize)
		sanityCheck(t, func(t float64) float64 {
			return CacheInterpolateLookup(cache, t)
		})
	})

	t.Run("CacheDirectLookup", func(t *testing.T) {
		cache := MakeCache(BasicPianoFn, DefaultCacheSize)
		sanityCheck(t, func(t float64) float64 {
			return CacheDirectLookup(cache, t)
		})
	})
}
