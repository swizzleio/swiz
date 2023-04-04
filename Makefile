.DEFAULT_GOAL := help

COMMIT_HASH="$(shell git rev-parse --short HEAD)"
PROJECTNAME="$(shell basename "$(PWD)")"
VERSION="0.1"
#VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"

# Project params
BINARY_DIR=./out
PROJECT_LOC=.
MAIN_LOC=./cmd/
BINARY_NAME=$(BINARY_DIR)/$(PROJECTNAME)
CMD_PACKAGE="github.com/swizzleio/swiz/internal/cmd"


# Use linker flags to provide version/build settings
LDFLAGS="-X $(CMD_PACKAGE).CommitHash=$(COMMIT_HASH) -X $(CMD_PACKAGE).Version=$(VERSION)"

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