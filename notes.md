
# Misc Notes

alright, let's see if I can just hack around a little on this. Let's try
some simple things, like:

- speedup/slowdown

- [x] try to get a middleman function for stream. I think there are some
  built-ins for this; effects if I'm not mistaken are something of a wrapper.

- [-] inverting channels, maybe. I *think* this is trivial if you get an
  interception function working, but we'll see.

- [x] See about doing "pure processing", independent of playback.

- [x] try to understand a little better what, exactly, is in a sample. How many
  are given at once, how often it's called under normal playback, etc.

- [x] let's try to create a slight motion effect between left and right; just
  oscillate on a loop.

- [x] create a looping effect. Perhaps provide a time range to "cut" and then
  loop from.

- this seems sorta unlikely, but is there any decent way to perform fft on this?
  Either an out-of-the-box package, or a fft library that I can staple to this?
  At a quick glance I'd *guess* go-dsp fft.FFTReal essentially does what I want;
  only question is do I actually know what I'm doing well enough to actually
  make practical sense of it (probably not).

- [x] Is there a way to generate notes? I mean, on some level of course there is
  an easy way to say "I want to generate an 'a' note for a period of time" using
  some proscribed technique. I don't think there is, which would lead to the
  next question: how hard would it be to do that?

- [x] as sort of a primer for the last one: let's see if I can just instrument a
  sine wave. Let's try a few different notes with it and see what happens. I
  think there is some theory about how to do things like that should be worth
  looking up, but only afterwards - let's try a hacked version that does 440 ->
  880 in a few steps.

- [-] note: I'm not worrying about this too much now as it's not super necessary for
  hacking, but I will add that I don't think my volume modifications are quite
  kosher - I'm given to belief that scaling up/down would require something more
  complicated than straight scalars.

- [x] So, I have what's frankly a pretty bad way to generate notes. Nonetheless,
  I'd be curious to see about making a simple system that can play multiple
  notes in sequence. To start, let's make a fixed frequency. I think it would
  also be interesting to add a fadeout for different notes (done as a wrapper)

- [x] even though they are also pretty ugly, let's see about implementing
  saw/square/triangle

- [x] A nice explanation of how to improve notes with an alternative note
  structure: https://keithwhor.com/music/. Let's see if I can clone the basic
  intensity structure of that.

- [x] add a linear fader and a way to loop notes with it.

- [x] Do I want to try some other fading techniques? Wouldn't be bad. Maybe
  exponential, logarithmic, equal power, sine, sigmoid. Or some subset of those
  - I don't need to punish myself here, and given the subtlety of the effect I
    may not have a terrible amount of patience.

- Also - do I want to add some slight white noise fuzzers? Sure. I'd like to
  spend a bit of time thinking about how to do this relatively efficiently
  though - I think a kludgy pseudo rng would probably work well.

- there are pops on cancellation. Can I eliminate those? Would require better
  instrumentation around the speaker itself. I'm inclined against trying this
  right now. I believe instrumenting my own speaker & volume effects should come
  first. At that point, I'd consider purging the mp3's, copying the useful
  parts, and working against oto directly. I'd add that I don't really need
  error handling for this either - I could simplify composition provided it's

- [x] let's see if I can do variadic frequency. Basically: have a note that say
  goes from A3 to A4 (and back again) every second. To be honest I actually
  don't really remember how to do this - clearly there'd be some way to
  structure the nature of rotation speed with time, but you'd have to be careful
  to make it continuous (or would I?) and properly periodic. I'd guess there are
  some fairly straightforward equations for that; and that


## Misc Links

- A fairly simple but useful look at a few different synths and the basics of
  their sound -
  http://www.acoustics.salford.ac.uk/acoustics_info/sound_synthesis/

- Here's an interesting article on the psychology of sound -
  https://en.wikipedia.org/wiki/Psychoacoustics. Includes a nice graph on
  relative loudness. If I were to spin a task off of this, I'd say play with a
  system that tries to normalize apparent loudness of a wave.

- A large set of synthesis tutorials form "sound on sound" -
  https://sonicbloom.net/en/63-in-depth-synthesis-tutorials-by-sound-on-sound/