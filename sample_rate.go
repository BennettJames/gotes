package gotes

import "time"

// SampleRate is the number of samples per second.
//
// This interface is essentially a copy of SampleRate from here -
// https://github.com/faiface/beep/blob/master/interface.go
type SampleRate int

// D returns the duration of n samples.
func (sr SampleRate) D(n int) time.Duration {
	return time.Second * time.Duration(n) / time.Duration(sr)
}

// N returns the number of samples that last for d duration.
func (sr SampleRate) N(d time.Duration) int {
	return int(d * time.Duration(sr) / time.Second)
}
