package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/bennettjames/gotes"
)

const (
	sampleRate = 48000
	channels   = 1
	bitDepth   = 16
	defaultDur = 1000 * time.Millisecond
)

func main() {
	var (
		flagSet = flag.NewFlagSet("flags", flag.PanicOnError)
		dir     = flagSet.String("dir", "", "directory to output samples")
	)
	flagSet.Parse(os.Args[1:])
	if err := genBasic(*dir); err != nil {
		log.Fatal(err)
	}
}

func genBasic(dir string) error {
	sinBuf := gotes.WriteWav(
		gotes.SampleRange(
			gotes.SinWave(gotes.NoteA3),
			sampleRate,
			defaultDur,
		),
		gotes.WavConfig{
			SampleRate: sampleRate,
			Channels:   channels,
			BitDepth:   bitDepth,
		},
	)
	if err := writeBuf(sinBuf, dir, "sin-a3.wav"); err != nil {
		return err
	}

	squareBuf := gotes.WriteWav(
		gotes.SampleRange(
			gotes.AmplifyWave(
				gotes.FixedAmplify(0.3),
				gotes.SquareWave(gotes.NoteA3),
			),
			sampleRate,
			defaultDur,
		),
		gotes.WavConfig{
			SampleRate: sampleRate,
			Channels:   channels,
			BitDepth:   bitDepth,
		},
	)
	if err := writeBuf(squareBuf, dir, "square-a3.wav"); err != nil {
		return err
	}

	triBuf := gotes.WriteWav(
		gotes.SampleRange(
			gotes.TriangleWave(gotes.NoteA3),
			sampleRate,
			defaultDur,
		),
		gotes.WavConfig{
			SampleRate: sampleRate,
			Channels:   channels,
			BitDepth:   bitDepth,
		},
	)
	if err := writeBuf(triBuf, dir, "tri-a3.wav"); err != nil {
		return err
	}

	sawBuf := gotes.WriteWav(
		gotes.SampleRange(
			gotes.AmplifyWave(
				gotes.FixedAmplify(0.4),
				gotes.SawWave(gotes.NoteA3),
			),
			sampleRate,
			defaultDur,
		),
		gotes.WavConfig{
			SampleRate: sampleRate,
			Channels:   channels,
			BitDepth:   bitDepth,
		},
	)
	if err := writeBuf(sawBuf, dir, "saw-a3.wav"); err != nil {
		return err
	}
	return nil
}

func writeBuf(data io.Reader, dir, file string) error {
	w, wErr := os.Create(path.Join(dir, file))
	if wErr != nil {
		return fmt.Errorf("could not write file: %w", wErr)
	}
	_, copyErr := io.Copy(w, data)
	if copyErr != nil {
		return fmt.Errorf("could not write file: %w", wErr)
	}
	return nil
}
