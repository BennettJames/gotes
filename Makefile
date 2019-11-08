
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
# Generates a set of samples that are used in the readme. Output is written to
# "doc", where they can be included in documentation.
#
.PHONY: gen-samples
gen-samples: bin/gensamples
	bin/gensamples -dir doc
