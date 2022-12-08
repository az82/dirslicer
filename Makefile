# Name of the command to build
# code is at /cmd/$(CMD)
# binaries will named $(CMD)
CMD := dirslice

# Platforms for which to build distributions
# Can be overriden with an environment variable
PLATFORMS ?= linux/amd64 windows/amd64 darwin/amd64 darwin/arm64

# Version name
# Will be used to name distributions
# Default is Git tag || branch || commit ID
# Can be overriden with an environment variable
VERSION ?= $(shell \
	git describe --tags --exact-match 2> /dev/null \
	|| git symbolic-ref -q --short HEAD \
	|| git rev-parse --short HEAD)

PLATFORM = $(subst /, ,$@)
OS = $(word 1, $(PLATFORM))
ARCH = $(word 2, $(PLATFORM))
DIST = $(CMD)-$(OS)-$(ARCH)-$(VERSION)

.PHONY: build
build: test
	mkdir -p bin
	go build -o 'bin/$(CMD)' ./cmd/$(CMD)

.PHONY: dist
dist: test $(PLATFORMS)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p dist/$(DIST)
	GOOS=$(OS) GOARCH=$(ARCH) go build -o 'dist/$(DIST)/$(CMD)' ./cmd/$(CMD)
	cp LICENSE README.md dist/$(DIST)
	cd dist; tar cjf $(DIST).tar.bz2 $(DIST)
	rm -r dist/$(DIST)

.PHONY: test
test:
	go test ./cmd/$(CMD)

.PHONY: clean
clean:
	rm -rf bin
	rm -rf DIST