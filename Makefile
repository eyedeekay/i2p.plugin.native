
CGO_ENABLED=0
GOPATH=$(HOME)/go/
VERSION="0.0.001"

all: fmt build readme

build:
	go build -o ./i2p.plugin.native ./cmd/i2p.plugin.native

readme:
	cat desc | tee README.md
	echo "" | tee -a README.md
	echo "There are some examples in the Makefile for now." | tee -a README.md
	echo "" | tee -a README.md
	echo "Here's a copy of the usage while I work on a better README.md:" | tee -a README.md
	echo "" | tee -a README.md
	echo "\`\`\`markdown" | tee -a README.md
	i2p.plugin.native -h 2>&1 | tee -a README.md
	echo "\`\`\`" | tee -a README.md

clean:
	rm -rf plugin *.zip *.su3

fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;

install: all
	install -m755 ./i2p.plugin.native ~/go/bin/i2p.plugin.native