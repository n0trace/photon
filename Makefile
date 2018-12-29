#!/usr/bin/env bash

all: dev run

fmt:
	goimports -l -w -local "gitlab.com/n0trace/photon" ./


test:
	echo "" > coverage.txt
	for d in $(shell go list ./... | grep -v vendor); do \
		go test -race -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out; \
	done
