PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = v$(shell cat .version)  # v1.1.1
COMMIT_SHA ?= $(shell git describe --always)  # abcd123

VERSION_MAJOR=$(shell echo $(VERSION) | cut -d '.' -f 1) # v1

GOBUILD=CGO_ENABLED=0 go build -a -ldflags "-X ${PKG}/version.Version=$(VERSION)"

APPNAME ?= prod

build:
	$(GOBUILD) -o out/$(APPNAME) cmd/$(APPNAME)/main.go

docker: build
	docker build -t uyinn28/istio-in-action-$(APPNAME):$(VERSION) -f dockerfiles/$(APPNAME).Dockerfile .

apply:
	version=$(VERSION) version_major=$(VERSION_MAJOR) envsubst < scripts/deployment/$(APPNAME).tmpl.yml | kubectl apply -f -

apply.docker: docker apply

apply.dryrun:
	version=$(VERSION) version_major=$(VERSION_MAJOR) envsubst < scripts/deployment/$(APPNAME).tmpl.yml

clean:
	docker rmi `docker images -q -f dangling=true` || echo
	rm -rf out


prod.apply: apply.docker
review.apply:
	APPNAME=review make apply.docker

dryrun.review:
	APPNAME=review make apply.dryrun


book:
	bash hugobook-builder.sh
