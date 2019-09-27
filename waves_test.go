package gotes

import (
	"fmt"
	"testing"
	"time"
)

func Test_PianoWave(t *testing.T) {
	wave := PianoWave(2000*time.Millisecond, NoteA4)

	fmt.Println(wave(3))
}
