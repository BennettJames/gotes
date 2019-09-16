package gotes

import (
	"math"
	"sync/atomic"
)

type Volume struct {
	// note (bs): I suspect this should be converted to a Gain and a BiGain node.
	stream BiStreamer
	vol    uint64
}

func NewVolume(stream BiStreamer, vol float64) *Volume {
	v := &Volume{
		stream: stream,
	}
	v.SetVolume(vol)
	return v
}

func (v *Volume) Stream(samples [][2]float64) {
	v.stream.Stream(samples)
	gain := math.Pow(2, v.GetVolume())
	for i := range samples {
		samples[i][0] *= gain
		samples[i][1] *= gain
	}
}

func (v *Volume) SetVolume(vol float64) {
	volBits := math.Float64bits(vol)
	atomic.StoreUint64(&v.vol, volBits)
}

func (v *Volume) GetVolume() float64 {
	volBits := atomic.LoadUint64(&v.vol)
	return math.Float64frombits(volBits)
}
