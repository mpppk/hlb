SHELL = /bin/bash

.PHONY: deps
deps:
	dep ensure -v

.PHONY: deps-ci
deps-ci:
	dep ensure -v -vendor-only=true

.PHONY: setup
setup:
	go get github.com/golang/dep/cmd/dep

.PHONY: lint
lint: deps-ci
	gometalinter

.PHONY: test
test: deps-ci
	go test ./...

.PHONY: coverage
coverage: deps-ci
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: codecov
codecov: deps-ci coverage
	bash <(curl -s https://codecov.io/bash)

.PHONY: build
build: deps-ci
	go build

.PHONY: cross-build-snapshot
cross-build: deps-ci
	goreleaser --rm-dist --snapshot

.PHONY: install
install: deps-ci
	go install

.PHONY: circleci
circleci:
	circleci build -e GITHUB_TOKEN=$GITHUB_TOKEN