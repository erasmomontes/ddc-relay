GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

all: clean fmt test ddc-relayy

.PHONY: ddc-relayy
ddc-relayy:
	$(GOBUILD) -o ./bin/ddc-relayy ./cmd/ddc-relayy/main.go

.PHONY: test
test:
	$(GOTEST) ./...

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f ./bin/*

.PHONY: fmt
fmt:
	$(GOFMT) ./...

.PHONY: bump
bump:
	bumpversion --current-version $(shell cat VERSION) minor --verbose --allow-dirty ./VERSION

.PHONY: aws-ecr
aws-ecr:
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 327921414325.dkr.ecr.us-west-2.amazonaws.com

.PHONY: build
build:
	VERSION=$(shell cat VERSION) docker-compose build --no-cache


.PHONY: push
push:
	VERSION=$(shell cat VERSION) docker-compose push
