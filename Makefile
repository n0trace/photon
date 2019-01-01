#!/usr/bin/env bash

all: dev run

fmt:
	goimports -l -w -local "gitlab.com/n0trace/photon" ./


test:
	go test -race  -covermode=atomic ./ ./common ./middleware -count 1
