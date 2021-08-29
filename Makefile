.DEFAULT_GOAL := help

#VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Project params
BINARY_DIR=./out
PROJECT_LOC=.
MAIN_LOC=./cmd/
BINARY_NAME=$(BINARY_DIR)/$(PROJECTNAME)

# Use linker flags to provide version/build settings
LDFLAGS="-X=main.Build=$(BUILD) -X=main.Version=v0.1"

all: clean unittest build-cli

## build-cli: Build the cli
.PHONY: build-cli
build-cli:
	go build -ldflags=$(LDFLAGS)  -o $(BINARY_NAME) $(MAIN_LOC)

## unittest: Run unit tests
.PHONY: unittest
unittest:
	go test -v ./...

## depupdate: Update dependencies
.PHONY: depupdate
depupdate:
	go get -u=patch $(MAIN_LOC)

## deptree: List dependency tree
.PHONY: deptree
deptree:
	goda tree $(PROJECT_LOC)/...:all

## clean: clean all dependencies
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_DIR)/*

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
