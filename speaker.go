package gotes

import (
	"context"

	"github.com/hajimehoshi/oto"
)

type (
	Streamer interface {
		Stream(samples []float64)
	}

	BiStreamer interface {
		Stream(samples [][2]float64)
	}
)

type Speaker struct {
	// todo (bs): see if you can get rid of sample rate as a foreign dependency.
	// Even if I just copy/paste it that's fine.
	sr      SampleRate
	stream  BiStreamer
	bufSize int
}

func NewSpeaker(
	sr SampleRate,
	stream BiStreamer,
	bufSize int,
) *Speaker {
	return &Speaker{
		sr:      sr,
		stream:  stream,
		bufSize: bufSize,
	}
}

func (s *Speaker) Run(ctx context.Context) error {

	otoCtx, otoCtxErr := oto.NewContext(int(s.sr), 2, 2, s.bufSize)
	if otoCtxErr != nil {
		return otoCtxErr
	}
	player := otoCtx.NewPlayer()
	defer player.Close()

	sampleSize := 512
	samples := make([][2]float64, sampleSize)
	buf := make([]byte, sampleSize*4)

	update := func() error {

		s.stream.Stream(samples)

		for i := range samples {
			for c := range samples[i] {
				val := samples[i][c]
				if val < -1 {
					val = -1
				}
				if val > 1 {
					val = 1
				}
				packedVal := int16(val * (1<<15 - 1))
				buf[i*4+c*2+0] = byte(packedVal)
				buf[i*4+c*2+1] = byte(packedVal >> 8)
			}
		}
		_, writeErr := player.Write(buf)
		return writeErr
	}

	for {
		select {
		case <-ctx.Done():
			// todo (bs): it would be good to change this to do a very quick fade-out
			// rather than an immediate cutout - e.g. fade out the volume in 50ms,
			// then return.
			return ctx.Err()
		default:
			if err := update(); err != nil {
				return err
			}
		}
	}
}
