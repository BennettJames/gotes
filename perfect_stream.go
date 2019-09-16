package gotes

type PerfectStream func(samples [][2]float64)

func (f PerfectStream) Stream(samples [][2]float64) (n int, ok bool) {
	f(samples)
	return len(samples), true
}

func (f PerfectStream) Err() error {
	return nil
}
