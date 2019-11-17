package gotes

import (
	"math"
	"testing"
	"time"
)

func Test_pianoCacheAccuracy(t *testing.T) {
	const sr = float64(48_000)
	const freq = NoteA3
	const dur = 1500 * time.Millisecond

	rmse1 := CalcWaveRMSE(
		uncachedPianoWave(dur, NoteA3),
		PianoNote(dur, NoteA3),
		48_000,
	)
	if rmse1 > 0.002 {
		t.Fatal("RMSE for cache piano should be below 0.2%", rmse1)
	}
}

func uncachedPianoWave(dur time.Duration, freq float64) WaveFn {
	durT := float64(dur) / float64(time.Second)
	dampen := math.Pow(0.5*math.Log(freq*0.3), 2)
	return AmplifyWave(
		AttackAndDecay(durT, dampen),
		// PianoWave(freq),
		IntegrateWave(
			MultiplyTime(freq),
			BasicPianoFn,
		),
	)
}
