# The "v" prefix exists for legacy reasons.
version:=$(shell git describe --tags)
commit:=$(shell git rev-parse HEAD)
timestamp:=$(shell date '+%Y-%m-%d %H:%M:%S')


pre:="World"
test:=$(date '+%Y-%m-%d %H:%M:%S')

versionLDFlags:=-X "github.com/emyrk/grow/internal/version.Version=${version}" \
               -X "github.com/emyrk/grow/internal/version.CommitSHA1=${commit}" \
               -X "github.com/emyrk/grow/internal/version.CompiledDate=${timestamp}"

.PHONY: build/client
build/client:
	go build \
    	-ldflags='$(versionLDFlags)' \
		-o ./bin/client \
		./cmd/client

build/server: bin/server
	go build \
    	-ldflags='$(versionLDFlags)' \
		-o ./bin/server \
		./


