
#
# runs tests for all subpackages.
#
.PHONY: test
test:
	go test ./...

#
# Targets to build and run gensamples, which will create a series of wav sample
# files and output them to doc/.
#
.PHONY: bin/gensamples
bin/gensamples:
	go build -o bin/gensamples ./cmds/gensamples
.PHONY: run-gensamples
run-gensamples: bin/gensamples
	bin/gensamples -dir doc

#
# Targets to build and run experiments, a hodgepode set of one-off examples that
# are useful when developing waves.
#
.PHONY: bin/experiments
bin/experiments:
	go build -o bin/experiments ./cmds/experiments
.PHONY: run-experiments
run-experiments: bin/experiments
	bin/experiments

#
# Targets to build and run basicexample, which will play a very simple note
# progression on repeat.
#
.PHONY: bin/basicexample
bin/basicexample:
	go build -o bin/basicexample ./cmds/basicexample
.PHONY: run-basicexample
run-basicexample: bin/basicexample
	bin/basicexample

#
# Targets to build and run the "twinkle" command; which will play
# twinkle-twinkle little star on repeat via the piano synthesizer.
#
.PHONY: bin/twinkle run-twinkle
bin/twinkle:
	go build -o bin/twinkle ./cmds/twinkle
run-twinkle: bin/twinkle
	bin/twinkle
