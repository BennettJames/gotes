package gotes

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func Test_WriteWav(t *testing.T) {
	// This does a very basic test of the wav file assembly utility. It's hard to
	// truly verify the behavior is correct without manual inspection; but this
	// performs at least a basic check against the sample->buffer pipeline.

	sampleRate := 48000
	dur := 2000 * time.Millisecond
	wave := PianoNote(dur, NoteA3)

	samples := SampleWave(wave, sampleRate, dur)
	expectedLen := sampleRate * int(float64(dur)/float64(time.Second)) * 2
	if len(samples) != expectedLen {
		t.Fatalf("Expected %d samples, got %d", expectedLen, len(samples))
	}

	dir, dirErr := ioutil.TempDir("", "sounds")
	if dirErr != nil {
		t.Fatal(dirErr)
	}
	defer os.RemoveAll(dir)

	outBuf := WriteWav(samples, WavConfig{
		SampleRate: uint32(sampleRate),
		Channels:   1,
		BitDepth:   16,
	})
	fileName := path.Join(dir, "out1.wav")
	outFile, outFileErr := os.Create(fileName)
	if outFileErr != nil {
		t.Fatal("could not write file: ", outFileErr)
	}
	_, copyErr := io.Copy(outFile, outBuf)
	if copyErr != nil {
		t.Fatal("could not write file: ", outFileErr)
	}
	stat, statErr := outFile.Stat()
	if statErr != nil {
		t.Fatal(statErr)
	}
	if stat.Size() < 50 {
		t.Fatalf("Outfile appears near empty: %d", stat.Size())
	}
}
