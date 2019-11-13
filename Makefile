
#
# runs tests for all subpackages.
#
.PHONY: test
test:
	go test ./...

#
# Generates a binary that when called will generate a set of .wav samples.
#
# Marked as phony as go has a fairly smart self-managed cache for building
# binaries.
#
.PHONY: bin/gensamples
bin/gensamples:
	go build -o bin/gensamples ./cmds/gensamples

#
# Builds experiments, a random set of extended audio experiments.
#
.PHONY: bin/experiments
bin/experiments:
	go build -o bin/experiments ./cmds/experiments

#
# Generates a set of samples that are used in the readme. Output is written to
# "doc", where they can be included in documentation.
#
.PHONY: run-gensamples
run-gensamples: bin/gensamples
	bin/gensamples -dir doc

#
# Executes the experiments.
#
.PHONY: run-experiments
run-experiments: bin/experiments
	bin/experiments

#
# Builds experiments, a random set of extended audio experiments.
#
.PHONY: bin/basicexample
bin/basicexample:
       go build -o bin/basicexample ./cmds/basicexample

#
# Executes the basicexample.
#
.PHONY:run-basicexample
run-basicexample: bin/basicexample
       bin/basicexample
