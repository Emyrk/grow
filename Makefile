# The "v" prefix exists for legacy reasons.
version=$(git describe --tags)
commit=$(git rev-parse HEAD)
timestamp=$(date '+%Y-%m-%d %H:%M:%S')


versionLDFlags=-X github.com/emyrk/grow/internal/version.Version=$(version) \
               -X github.com/emyrk/grow/internal/version.CommitSHA1=${commit} \
               -X github.com/emyrk/grow/internal/version.CompiledDate=$(timestamp)

build/client:
	go build \
    	-ldflags="$(versionLDFlags)" \
		-o ./bin/client \
		./cmd/client

build/server:
	echo $(version)
	go build \
    	-ldflags="$(versionLDFlags)" \
		-o ./bin/server \
		./

