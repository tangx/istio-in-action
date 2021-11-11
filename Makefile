PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = $(shell cat .version)
COMMIT_SHA ?= $(shell git describe --always)

GOBUILD=CGO_ENABLED=0 go build -a -ldflags "-X ${PKG}/version.Version=$(VERSION)"

build:
	$(GOBUILD) -o out/prod cmd/prod/main.go
