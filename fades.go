package gotes

import (
	"math"
	"time"
)

func LinearFadeWave(cycle time.Duration, w1, w2 WaveFn) WaveFn {
	cycleT := float64(cycle) / float64(time.Second)
	return func(t float64) float64 {
		if t > cycleT {
			return w2(t)
		}
		if t < 0 {
			return w1(t)
		}
		return (w1(t) * (1 - t/cycleT)) + (w2(t) * (t / cycleT))
	}
}

func LinearFadeLooperWave(
	dur, fade time.Duration,
	fns ...WaveFn,
) WaveFn {
	if len(fns) == 0 {
		fns = []WaveFn{
			ZeroWave(),
		}
	}
	durT := float64(dur) / float64(time.Second)
	fadeT := float64(fade) / float64(time.Second)

	return func(t float64) float64 {
		nextI := (int(t/durT) + len(fns)) % len(fns)
		lastI := (nextI - 1 + len(fns)) % len(fns)
		lastFn, nextFn := fns[lastI], fns[nextI]
		if nextI == 0 && t < durT {
			// this is a little inelegant, but w/e
			lastFn = ZeroWave()
		}

		completion := math.Min(1, math.Mod(t, durT)/fadeT)
		return (lastFn(t) * (1 - completion)) + (nextFn(t) * completion)
	}
}

func ExpFadeLooperWave(
	dur, fade time.Duration,
	fns ...WaveFn,
) WaveFn {
	if len(fns) == 0 {
		fns = []WaveFn{
			ZeroWave(),
		}
	}
	durT := float64(dur) / float64(time.Second)
	fadeT := float64(fade) / float64(time.Second)

	return func(t float64) float64 {
		nextI := (int(t/durT) + len(fns)) % len(fns)
		lastI := (nextI - 1 + len(fns)) % len(fns)
		lastFn, nextFn := fns[lastI], fns[nextI]
		if nextI == 0 && t < durT {
			lastFn = ZeroWave()
		}

		// note - these are slightly off initially: the 0.01 value causes a
		// discontinuity. I think it's kinda useless as-is; will need more proper
		// math to get the right value. I'm just gonna leave it with this warning
		// for now; but yeah seems like it
		completion := math.Min(1, math.Mod(t, durT)/fadeT)
		down := math.Pow(0.1, 2*completion) - 0.01
		up := 1.01 - math.Pow(0.1, 2*completion)
		return (lastFn(t) * down) + (nextFn(t) * up)
	}
}
