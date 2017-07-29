REPO=ecrsidecar

.PHONY: build docker

build:
	docker run --rm -v $$(pwd):/go/src/github.com/andrepinto/$(REPO) \
		-w /go/src/github.com/andrepinto/$(REPO) \
		golang:1.8 go build -v -a -tags netgo -installsuffix netgo -ldflags '-w'

docker:
	docker build -t andrepinto/$(REPO) .
