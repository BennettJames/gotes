package gotes

import (
	"context"

	"github.com/hajimehoshi/oto"
)

// Speaker is used to play back a streamer over physical speakers.
type Speaker struct {
	sr      SampleRate
	stream  Streamer
	bufSize int
}

// NewSpeaker initializes a speaker around the given stream. The bufSize controls
// how many samples are kept in a buffer, and the sample rate is how many
// samples per second to take from the stream.
func NewSpeaker(
	sr SampleRate,
	stream Streamer,
	bufSize int,
) *Speaker {
	return &Speaker{
		sr:      sr,
		stream:  stream,
		bufSize: bufSize,
	}
}

// Run will play back the underlying stream until it's cancelled or an error
// occurs.
func (s *Speaker) Run(ctx context.Context) error {

	otoCtx, otoCtxErr := oto.NewContext(int(s.sr), 2, 2, s.bufSize)
	if otoCtxErr != nil {
		return otoCtxErr
	}
	player := otoCtx.NewPlayer()
	defer player.Close()

	sampleSize := 512
	samples := make([]float64, sampleSize)
	buf := make([]byte, sampleSize*4)

	update := func() error {

		s.stream.Stream(samples)

		for i := range samples {
			// note (bs): this float array is a not-very-good way to make this
			// a 2-channel sample. Ideally, there'd be better management of
			// mono-vs-stereo via the injected interfaces.
			for c, val := range []float64{samples[i], samples[i]} {
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
