
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

examples: clean railroad-example brb-example railroad-example-win brb-example-win

clean:
	rm -rf plugin *.zip *.su3

fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;

export sumrrlinux=`sha256sum "../railroad-linux.su3"`
export sumrrwindows=`sha256sum "../railroad-windows.su3"`
export sumbblinux=`sha256sum "../brb-linux.su3"`
export sumbbwindows=`sha256sum "../brb-windows.su3"`

karens: fmt
	GOOS=windows go build -o karen.exe -ldflags "-extldflags -static" -tags netgo ./cmd/i2p.plugin.native/karen
	GOOS=linux go build -o karen -ldflags "-extldflags -static" -tags netgo ./cmd/i2p.plugin.native/karen
	GOOS=darwin go build -o karen-darwin -ldflags "-extldflags -static" -tags netgo ./cmd/i2p.plugin.native/karen
	file karen*

export sumklinux=`sha256sum "karen"`
export sumkwindows=`sha256sum "karen.exe"`
export sumkdarwin=`sha256sum "karen-darwin"`

upload-karens: karens
	gothub release -u eyedeekay -r "i2p.plugin.native" -t v$(VERSION) -d "I2P Plugin Generator and Supervisor"
	gothub upload -R -u eyedeekay -r "i2p.plugin.native" -t v$(VERSION) -l "$(sumklinux)" -n "karen" -f "karen"
	gothub upload -R -u eyedeekay -r "i2p.plugin.native" -t v$(VERSION) -l "$(sumkwindows)" -n "karen.exe" -f "karen.exe"
	gothub upload -R -u eyedeekay -r "i2p.plugin.native" -t v$(VERSION) -l "$(sumkdarwin)" -n "karen-darwin" -fss "karen-darwin"

install: all
	install -m755 ./i2p.plugin.native ~/go/bin/i2p.plugin.native