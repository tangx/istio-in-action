PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = $(shell cat .version)
COMMIT_SHA ?= $(shell git describe --always)

GOBUILD=CGO_ENABLED=0 go build -a -ldflags "-X ${PKG}/version.Version=$(VERSION)"

APPNAME ?= prod

build:
	$(GOBUILD) -o out/$(APPNAME) cmd/$(APPNAME)/main.go

docker: build
	docker build -t uyinn28/istio-in-action-$(APPNAME):$(VERSION) -f dockerfiles/$(APPNAME).Dockerfile .

apply: 
	version=$(VERSION) envsubst < scripts/deployment/$(APPNAME).yml.tmpl | kubectl -f -

clean:
	docker rmi `docker images -q -f dangling=true` || echo
	rm -rf out
