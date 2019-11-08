package gotes

import (
	"bytes"
	"io"
	"math"
	"time"
)

type (
	// WavConfig specifies properties that
	WavConfig struct {
		SampleRate uint32
		Channels   uint32
		BitDepth   uint32
	}
)

// SampleRange will generate a range of samples from the wave and return a
// uint16-encoded byte array with the samples, up to the given duration.
func SampleRange(fn WaveFn, rate int, dur time.Duration) []byte {
	samples := []byte{}
	sampleSize := (float64(dur) / float64(time.Second)) * float64(rate)
	for i := 0; i < int(sampleSize); i++ {
		sample := fn(float64(i) / float64(rate))
		sample = math.Min(1, math.Max(-1, sample))
		packedSample := int16(sample * (1<<15 - 1))
		samples = append(samples, byte(packedSample), byte(packedSample>>8))
	}
	return samples
}

// WriteWav will format the given samples to a wav byte stream.
func WriteWav(samples []byte, config WavConfig) io.Reader {
	// Reference for wav file structure:
	// https://blogs.msdn.microsoft.com/dawate/2009/06/23/intro-to-audio-programming-part-2-demystifying-the-wav-format/
	sampleLen := uint32(len(samples))
	outBytes := flattenByteBuffers(
		/* header sGroupID */ []byte("RIFF"),
		/* header dwFileLength */ uint32Array(36+sampleLen),
		/* header sRiffType */ []byte("WAVE"),
		/* chunk sGroupID */ []byte("fmt "),
		/* chunk dwChunkSize */ uint32Array(config.BitDepth),
		/* chunk wFormatTag */ uint16Array(1),
		/* chunk wChannels */ uint16Array(uint16(config.Channels)),
		/* chunk dwSamplesPerSec */ uint32Array(config.SampleRate),
		/* chunk dwAvgBytesPerSec */ uint32Array(config.SampleRate*4),
		/* chunk wBlockAlign */ uint16Array(uint16(config.Channels*2)),
		/* chunk dwBitsPerSample */ uint16Array(16),
		/* data sGroupID */ []byte("data"),
		/* data dwChunkSize */ uint32Array(sampleLen),
		/* data sampleData */ samples,
	)
	return bytes.NewBuffer(outBytes)
}

func uint16Array(v uint16) []byte {
	return []byte{
		byte(v & 0xff),
		byte((v >> 8) & 0xff),
	}
}

func uint32Array(v uint32) []byte {
	return []byte{
		byte(v & 0xff),
		byte((v >> 8) & 0xff),
		byte((v >> 16) & 0xff),
		byte((v >> 24) & 0xff),
	}
}

func flattenByteBuffers(buffers ...[]byte) []byte {
	all := []byte{}
	for _, buf := range buffers {
		for _, b := range buf {
			all = append(all, b)
		}
	}
	return all
}
