PACKAGE = .

export GO15VENDOREXPERIMENT=1
export GO111MODULE=on

BUILD_VERSION=1.0
BUILD_NUMBER=$(strip $(if $(TRAVIS_BUILD_NUMBER), $(TRAVIS_BUILD_NUMBER), 0))
BUILD_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

SRC = $(shell glide nv $(PACKAGE))

linux: 
	GOARCH=amd64 GOOS=linux go build -v -ldflags "-X main.BuildVersion=$(BUILD_VERSION).$(BUILD_NUMBER) -X main.BuildCommit=$(BUILD_COMMIT) -X main.BuildDate=$(BUILD_DATE)" 

build: 
	go build -v -ldflags "-X main.BuildVersion=$(BUILD_VERSION).$(BUILD_NUMBER) -X main.BuildCommit=$(BUILD_COMMIT) -X main.BuildDate=$(BUILD_DATE)" 

run: build
	./vsub serve
