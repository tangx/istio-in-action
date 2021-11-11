PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = v$(shell cat .version)
COMMIT_SHA ?= $(shell git describe --always)

VERSION_MAJAR=$(shell echo $(VERSION) | cut -d '.' -f 1)

GOBUILD=CGO_ENABLED=0 go build -a -ldflags "-X ${PKG}/version.Version=$(VERSION)"

APPNAME ?= prod

build:
	$(GOBUILD) -o out/$(APPNAME) cmd/$(APPNAME)/main.go

docker: build
	docker build -t uyinn28/istio-in-action-$(APPNAME):$(VERSION) -f dockerfiles/$(APPNAME).Dockerfile .

apply:
	version=$(VERSION_MAJAR) envsubst < scripts/deployment/$(APPNAME).yml.tmpl | kubectl apply -f -

docker.apply: docker.apply

clean:
	docker rmi `docker images -q -f dangling=true` || echo
	rm -rf out
