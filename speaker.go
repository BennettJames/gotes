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

	samples := make([][2]float64, s.bufSize)
	buf := make([]byte, s.bufSize*4)

	update := func() error {

		// ques (bs): is this getting samples at buffer size rather than 512? If so,
		// should I deliberately slow it down? Actually, this doesn't even really
		// look like a proper buffer.
		//
		// I think streamFrom does that - let's consider using it. I'd like to do a
		// more thorough rethink of buffering before committing to anything in
		// particular though.
		//
		// streamFrom(s.stream, samples)

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

func streamFrom(stream BiStreamer, samples [][2]float64) {

	var tmp [512][2]float64

	for len(samples) > 0 {
		toStream := len(tmp)
		if toStream > len(samples) {
			toStream = len(samples)
		}

		// clear the samples
		for i := range samples[:toStream] {
			samples[i] = [2]float64{}
		}

		// mix the stream
		stream.Stream(tmp[:toStream])
		for i := range tmp[:toStream] {
			samples[i][0] += tmp[i][0]
			samples[i][1] += tmp[i][1]
		}

		samples = samples[toStream:]
	}

	/*
		var tmp [512][2]float64

		for len(samples) > 0 {
			toStream := len(tmp)
			if toStream > len(samples) {
				toStream = len(samples)
			}

			// clear the samples
			for i := range samples[:toStream] {
				samples[i] = [2]float64{}
			}

			stream.Stream(tmp[:toStream])
			for i := range tmp[:toStream] {
				samples[i][0] += tmp[i][0]
				samples[i][1] += tmp[i][1]
			}

			samples = samples[toStream:]
		} */
}
