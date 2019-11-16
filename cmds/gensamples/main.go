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
	if err := genWavs(*dir); err != nil {
		log.Fatal(err)
	}
}

// genWavs writes out a set of wav samples of various waves to the specified
// directory.
func genWavs(dir string) error {
	waveData := []struct {
		fName string
		wave  gotes.WaveFn
		dur   time.Duration
	}{
		{
			fName: "sin-a3.wav",
			wave:  gotes.SinWave(gotes.NoteA3),
		},
		{
			fName: "square-a3.wav",
			wave: gotes.AmplifyWave(
				gotes.FixedAmplify(0.3),
				gotes.SquareWave(gotes.NoteA3),
			),
		},
		{
			fName: "tri-a3.wav",
			wave:  gotes.TriangleWave(gotes.NoteA3),
		},
		{
			fName: "saw-a3.wav",
			wave: gotes.AmplifyWave(
				gotes.FixedAmplify(0.4),
				gotes.SawWave(gotes.NoteA3),
			),
		},
		{
			fName: "good-osc-1.wav",
			wave: gotes.AmplifyWave(
				gotes.Gain(0.5),
				gotes.IntegrateWave(
					gotes.OscillateTime(2.0, 0.2),
					gotes.SinWave(gotes.NoteA3),
				),
			),
			dur: 10 * time.Second,
		},
		{
			fName: "bad-osc-1.wav",
			wave: gotes.AmplifyWave(
				gotes.Gain(0.5),
				gotes.IntegrateWave(
					gotes.BadOscillateTime(2.0, 0.2),
					gotes.SinWave(gotes.NoteA3),
				),
			),
			dur: 10 * time.Second,
		},
		{
			fName: "bad-osc-2.wav",
			wave: gotes.AmplifyWave(
				gotes.Gain(0.5),
				gotes.IntegrateWave(
					gotes.BadOscillateTime2(1.0, 0.2),
					gotes.SinWave(gotes.NoteA3),
				),
			),
			dur: 10 * time.Second,
		},
	}

	for _, w := range waveData {
		dur := w.dur
		if dur == 0 {
			dur = defaultDur
		}
		streamer := gotes.StreamerFromWave(w.wave, sampleRate)
		buf := gotes.WriteWav(
			gotes.SampleStreamer(streamer, sampleRate, dur),
			gotes.WavConfig{
				SampleRate: sampleRate,
				Channels:   channels,
				BitDepth:   bitDepth,
			},
		)
		if err := writeBuf(buf, dir, w.fName); err != nil {
			return fmt.Errorf("Error writing '%s': %w", w.fName, err)
		}
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
