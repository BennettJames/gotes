# Gotes

_**Author's note** - unfortunately, github does not allow embedding of sound
files in a readme. A few examples are linked that can be downloaded and played
back. At a latter date, I will try to port this to a github page, which should
allow direct audio embedding._

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

Gotes is based on [oto](oto), a cross-platform Go library that dynamically links
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
that will oscillate the frequency between 220 and 440 -

```go
IntegrateWave(
  func(t float64) float64 {
		return 1.5*t + math.Sin(t*2*math.Pi)/(4*math.Pi)
  },
  SinWave(220),
)
```

Anyway: that's the basics of gotes. There are many other waves and modifier
functions; head over to the [docs](docs) to see them all.


[oto]:https://github.com/hajimehoshi/oto
[docs]:https://godoc.org/github.com/BennettJames/gotes