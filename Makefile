# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=tcpsb
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build build-linux

build:
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f ./build/$(BINARY_NAME)
	rm -f ./build/$(BINARY_UNIX)
run:
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v ./...
	./build/$(BINARY_NAME)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_UNIX) -v
