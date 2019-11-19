# Gotes

_**Author's note** - unfortunately, github does not allow embedding of sound
files in a readme. A few examples are linked that can be downloaded and played
back. A few samples can be [found here][wav-page]._

[![GoDoc](https://godoc.org/github.com/BennettJames/gotes?status.svg)](https://godoc.org/github.com/BennettJames/gotes)

Gotes ("go notes") is a just-for-fun, simple sound synthesis library for Go.
Design-wise, there is a strong emphasis on immutable function composition to
achieve synthesis. It's not terribly practical, but leads to very easy and
obvious composition.

Detailed usage and explanation will follow; but here's a quick, complete example
of gotes being used to play a simple set of piano-like notes:

```go
import (
	"context"
	"log"
	"time"

	"github.com/bennettjames/gotes"
)

func main() {
	const (
		sampleRate   = gotes.SampleRate(48000)
		speakerCache = 100 * time.Millisecond
		noteDuration = 2000 * time.Millisecond
	)

	wave := gotes.Looper(
		time.Second,
		gotes.PianoNote(noteDuration, gotes.NoteA3),
		gotes.PianoNote(noteDuration, gotes.NoteB3),
		gotes.PianoNote(noteDuration, gotes.NoteG4),
		gotes.PianoNote(noteDuration, gotes.NoteF4),
  )

	streamer := gotes.StreamerFromWave(sampleRate, wave)
	speaker := gotes.NewSpeaker(sampleRate, streamer, sampleRate.N(speakerCache))
	log.Fatal(speaker.Run(context.Background()))
}
```


## Dependencies; Compatibility; Stability

Gotes is based on [oto][oto], a cross-platform Go library that dynamically links
platform-specific audio libraries. There can be platform specific requirements
for oto; I recommend viewing it's project page for help setting up any
dependencies.

This was primarily developed on an old mac mini; it should be portable to other
platforms but has not been thoroughly tested.

This is very much an experimental library; so the API is not terribly stable.

## Digital Sound Refresher

A quick, simple review of how sound works, for those who might not know or just
haven't thought much about it recently -

Sound is essentially a wave of rapid, tiny fluctuations in air pressure. A
computer speaker can create sound by rapidly vibrating it's diaphragm back and
forth, creating particular waves of high and low pressure that our ears map to
the sounds we all know and love.

To control the speaker, the computer sends tens of thousands of discrete values
to the speaker a second. These are all single values that describe the
underlying wave. They are not continuous; they are just *samples* at various
points that come close enough to representing the underlying sound for the
speaker to convincingly play.

## Basics of Wave Composition in Gotes

In gotes, the fundamental unit of sound is a _wave function_. This is a simple
function that maps a time argument to a sample. Time proceeds from 0, and goes
up by 1 for every second that passes. Samples are all in the range of -1 to 1.

To start, let's define the simplest wave possible: a basic sine wave. It starts
at zero; goes up to 1; down to -1; and back to zero once every second -

```go
func BasicSinFn(t float64) float64 {
	return math.Sin(2 * math.Pi * t)
}
```

This is quite slow - only 1 hz, which can't be heard by the human ear. We'd need
to increase the speed of the wave to actually hear it. In gotes, that's done by
applying a _time function_. A wave function maps a time value to a sample value;
a time function maps one time value to another.

Gotes composes time functions and wave functions together using `IntegrateWave`.
This takes two arguments - a time function and a wave function. It returns a new
wave function, where the given time argument is first passed to the time
function, then that value is given to the wave function.

If we want to say boost this wave up to a audible frequency, we can apply a
constant multiplier to the time wave -


```go
func SinWave(frequency float64) WaveFn {
  return IntegrateWave(
    MultiplyTime(frequency),
    BasicSinFn,
  )
}
```

Note that this is still just returning a wave function. If we'd like, we could
apply `IntegrateWave` all over again. For instance; here's a usage of `SinWave`
that will oscillate the frequency between 220 and 440 every five seconds -

```go
IntegrateWave(
  gotes.OscillateTime(1.0, 0.2),
  SinWave(220),
)
```

That's the basics of gotes waves. There are many other waves and
modifier functions; head over to the [docs][docs] to see them all.


## Mutability and Managing Output

Waves themselves are immutable and inert. They describe a wave over time; but we
still need to be processed, managed, and played.

An intermediary interface, `Streamer`, is used to create sets of samples from
waves that. Here's it's definition -

```go
type Streamer interface {
	Stream(samples []float64)
}
```

And here's a wave being converted to a streamer -

```go
sr := SampleRate(48_000)
streamer := StreamerFromWave(SinWave(NoteA3), sr)
```

This creates a streamer that will take 48,000 samples from the provided sin wave
per second. The streamer will be called repeatedly with a float array; and the
streamer is responsible for fully populating it with values. If say a sample
array of size 1,000 is being used, the streamer will be called 48 times every
second, and each time it will fill the array with the next 1,000 samples.

Streamer implementations at a minimum need to be aware of sample rate and the
passage of time. They can be extended to handle other time-sensitive and mutable
state. For example, the `Keyboard` class is a streamer that handles realtime
playback of piano notes. Notes can be dynamically triggered on keyboard, which
will then manage the playback and eventual fadeout/removal of the note.

Streamers can be used to output the sound in two ways: as .wav files, or as
direct playback on speakers. Here's an example of a second-long sample being
written to file -

```go
sampleLen := 1 * time.Second
buf := WriteWav(
	SampleStreamer(streamer, sr, 1*time.Second),
	WavConfig{
		SampleRate: sr,
	},
)
if err := ioutil.WriteFile("out.wav", buf, 0644); err != nil {
	return fmt.Errorf("Error writing 'out.wav': %w", err)
}
```

Here's the speaker being set up for playback -

```go
speakerBuffer := 100 * time.Millisecond
speaker := gotes.NewSpeaker(sr, streamer, sr.N(speakerBuffer))
log.Fatal(speaker.Run(context.Background()))
```

They have the same basic pattern: a sample rate is set; a streamer is set up;
and the playback system is initialized with both.


[oto]:https://github.com/hajimehoshi/oto
[docs]:https://godoc.org/github.com/BennettJames/gotes
[wav-page]:https://bennettjames.github.io/gotes/index.html
