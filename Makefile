
CGO_ENABLED=0
GOPATH=$(HOME)/go/
VERSION="0.0.001"

all: fmt build readme

build:
	go build -o ./scripts/bin/i2p.plugin.native ./scripts/src

readme:
	echo "I2P native plugin generation tool" | tee README.md
	echo "=================================" | tee -a README.md
	echo "" | tee -a README.md
	echo "I wrote this way faster than I documented it. Shocking, right?" | tee -a README.md
	echo "" | tee -a README.md
	echo "This is a handy little tool for assembling I2P plugins when those" | tee -a README.md
	echo "plugins don't have a clean way to interface with the JVM, or just don't" | tee -a README.md
	echo "need one. Think of it a little like \`checkinstall\` but for I2P Plugins." | tee -a README.md
	echo "Right now it mostly works, and it's pretty cleanly put together." | tee -a README.md
	echo "" | tee -a README.md
	echo "There are some examples in the Makefile for now." | tee -a README.md
	echo "" | tee -a README.md
	echo "Here's a copy of the usage while I work on a better README.md:" | tee -a README.md
	echo "" | tee -a README.md
	echo "\`\`\`markdown" | tee -a README.md
	./scripts/bin/i2p.plugin.native -h 2>&1 | tee -a README.md
	echo "\`\`\`" | tee -a README.md

examples: clean railroad-example brb-example railroad-example-win brb-example-win

railroad-lin:
	cp -v $(GOPATH)src/i2pgit.org/idk/railroad/railroad .

railroad-example: all clean railroad-lin
	./scripts/bin/i2p.plugin.native -name=railroad \
		-signer=hankhill19580@gmail.com \
		-version 0.0.031 \
		-author=hankhill19580@gmail.com \
		-autostart=true \
		-clientname=railroad \
		-consolename="Railroad Blog" \
		-delaystart="1" \
		-desc="$(cat desc)" \
		-exename=railroad \
		-exeperm=0755 \
		-license=MIT \
		-res=config
	cp -v *.su3 ../railroad-linux.su3
	unzip -o railroad.zip -d railroad-zip

brb-lin:
	cp -v $(GOPATH)src/github.com/eyedeekay/brb/brb .

brb-example: all clean brb-lin
	./scripts/bin/i2p.plugin.native -name=brb \
		-signer=hankhill19580@gmail.com \
		-version 0.0.09 \
		-author=hankhill19580@gmail.com \
		-autostart=true \
		-clientname=brb \
		-command="\$$PLUGIN/lib/brb -dir=\$$PLUGIN/lib -eris=true -i2psite=true" \
		-consolename="BRB IRC" \
		-delaystart="1" \
		-desc="$(cat ircdesc)" \
		-exename=brb \
		-exeperm=0755 \
		-license=MIT
	cp -v *.su3 ../brb-linux.su3
	unzip -o brb.zip -d brb-zip

railroad-win:
	cp -v $(GOPATH)src/i2pgit.org/idk/railroad/railroad.exe .

railroad-example-win: all clean railroad-win
	./scripts/bin/i2p.plugin.native -name=railroad \
		-signer=hankhill19580@gmail.com \
		-version 0.0.031 \
		-author=hankhill19580@gmail.com \
		-autostart=true \
		-clientname=railroad.exe \
		-consolename="Railroad Blog" \
		-delaystart="1" \
		-desc="$(cat desc)" \
		-exename=railroad.exe \
		-license=MIT \
		-targetos="windows" \
		-res=config
	cp -v *.su3 ../railroad-windows.su3
	unzip -o railroad.zip -d railroad-zip-win

brb-win:
	cp -v $(GOPATH)src/github.com/eyedeekay/brb/brb.exe .

brb-example-win: all clean brb-win
	./scripts/bin/i2p.plugin.native -name=brb \
		-signer=hankhill19580@gmail.com \
		-version 0.0.09 \
		-author=hankhill19580@gmail.com \
		-autostart=true \
		-clientname=brb.exe \
		-command="\$$PLUGIN/lib/brb.exe -dir=\$$PLUGIN/lib -eris=true -i2psite=true" \
		-consolename="BRB IRC" \
		-delaystart="1" \
		-desc="$(cat ircdesc)" \
		-exename=brb.exe \
		-license=MIT \
		-targetos="windows" \
		-res=windll
	cp -v *.su3 ../brb-windows.su3
	unzip -o brb.zip -d brb-zip-win

clean:
	rm -rf plugin *.zip *.su3

fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;

export sumrrlinux=`sha256sum "../railroad-linux.su3"`
export sumrrwindows=`sha256sum "../railroad-windows.su3"`
export sumbblinux=`sha256sum "../brb-linux.su3"`
export sumbbwindows=`sha256sum "../brb-windows.su3"`

upload:
	gothub upload -R -u eyedeekay -r "railroad" -t 0.0.031 -l "$(sumrrlinux)" -n "railroad-linux.su3" -f "../railroad-linux.su3"
	gothub upload -R -u eyedeekay -r "railroad" -t 0.0.031 -l "$(sumrrwindows)" -n "railroad-windows.su3" -f "../railroad-windows.su3"
	gothub upload -R -u eyedeekay -r "brb" -t v0.0.09 -l "$(sumbblinux)" -n "brb-linux.su3" -f "../brb-linux.su3"
	gothub upload -R -u eyedeekay -r "brb" -t v0.0.09 -l "$(sumbbwindows)" -n "brb-windows.su3" -f "../brb-windows.su3"

karens: fmt
	GOOS=windows go build -o karen.exe -ldflags "-extldflags -static" -tags netgo ./scripts/src/karen
	GOOS=linux go build -o karen -ldflags "-extldflags -static" -tags netgo ./scripts/src/karen
	GOOS=darwin go build -o karen-darwin -ldflags "-extldflags -static" -tags netgo ./scripts/src/karen
	file karen*

export sumklinux=`sha256sum "karen"`
export sumkwindows=`sha256sum "karen.exe"`
export sumkdarwin=`sha256sum "karen-darwin"`

upload-karens: karens
	gothub release -u eyedeekay -r "i2p.plugin.native" -t v$(VERSION) -d "I2P Plugin Generator and Supervisor"
	gothub upload -R -u eyedeekay -r "i2p.plugin.native" -t v$(VERSION) -l "$(sumklinux)" -n "karen" -f "karen"
	gothub upload -R -u eyedeekay -r "i2p.plugin.native" -t v$(VERSION) -l "$(sumkwindows)" -n "karen.exe" -f "karen.exe"
	gothub upload -R -u eyedeekay -r "i2p.plugin.native" -t v$(VERSION) -l "$(sumkdarwin)" -n "karen-darwin" -f "karen-darwin"
