# Gotes

Gotes ("go notes") is a just-for-fun, simple sound synthesis library for Go.
Design-wise, there is a strong emphasis on immutable function composition to
achieve synthesis. It's not terribly practical; but

It is based on [oto](oto), which is a cross-platform Go library that tries to
link into the platform-specific audio libraries. It is not well tested across
platforms: I've tried it on a few different computers, and it sometimes works;
sometimes can run into buffering issues depending on the platform.


## Basics of Sound

(note [bs]: let's revisit these graphs; there fine to get started though)

| Wave Type | Wave Form                                           | Sample |
| :---:     |   :---:                                             | :---:  |
| Sine      | ![sin](doc/sin-sample.png)                          |        |
| Sawtooth  | ![sawtooth](doc/sawtooth-sample.png)                |        |
| Triangle  | ![triangle](doc/triangle-sample.png)                |        |
| Square    | ![square](doc/square-sample.png)                    |        |






## Composing Gotes









(todo [bs]: let's at least *try* to create a simple wasm example. Don't spend
too long on it - honestly I'd say create a plan, try to execute it within ten
minutes, then if you don't have a result decide whether it's worth continuing.)

(todo [bs]: let's also try to staple wav generation to this. There's )




[oto]:https://github.com/hajimehoshi/oto
