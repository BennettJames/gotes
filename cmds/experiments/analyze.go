package main

import (
	"fmt"
	"math"
	"time"

	"github.com/bennettjames/gotes"
)

// analyzeStream logs some simplistic stream information once per second.
// Occasionally useful/interesting when debugging.
func analyzeStream(
	sr gotes.SampleRate,
	srcStream gotes.Streamer,
) func(samples []float64) {
	startTime := time.Now()
	lastLog := startTime
	totalSamples := 0
	totalCalls := 0

	return func(samples []float64) {

		totalSamples += len(samples)
		totalCalls++
		now := time.Now()
		if now.Sub(lastLog) >= time.Second {
			lastLog = now
			fmt.Println("max/min/average", getInfo(samples))
		}
	}
}

type sampleInfo struct {
	min     float64
	max     float64
	average float64
}

func getInfo(samples []float64) sampleInfo {
	// some notes on sampling rate:
	//
	// - average chunk size seems to be 512, but can vary. Not sure who
	//   determines that; might just be a set value somewhere.
	//
	// - values are in range of [-1, 1]. Average of that per sample is
	//   zero, but can vary a bit.

	info := sampleInfo{
		min: math.MaxFloat64,
		max: -math.MaxFloat64,
	}
	for i := 0; i < len(samples); i++ {
		s0 := samples[i]
		info.min = math.Min(s0, info.min)
		info.max = math.Max(s0, info.max)
		info.average += s0
	}
	info.average /= float64(len(samples))
	return info
}
