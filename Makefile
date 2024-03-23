PKG = github.com/litencatt/rlp
COMMIT = $$(git describe --tags --always)
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

export GO111MODULE=on

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)
RLP_BINARY ?= ./rlp

deps:
	go install github.com/spf13/cobra-cli@latest

build:
	go build -ldflags="$(BUILD_LDFLAGS)" -o $(RLP_BINARY) cmd/rlp/main.go

prerelease:
	go mod tidy
	ghch -w -A --format=markdown -N $(NEXT_VER)
	gocredits -skip-missing -w .

release:
	git tag ${NEXT_VER}
	git push origin main --tag
	goreleaser --rm-dist