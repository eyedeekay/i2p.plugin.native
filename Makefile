
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

examples: clean railroad-example brb-example

railroad-example: all clean
	./scripts/bin/i2p.plugin.native -name=railroad \
		-signer=hankhill19580@gmail.com \
		-version 0.0.031 \
		-author=hankhill19580@gmail.com \
		-autostart=true \
		-clientname=railroad \
		-consolename="Railroad Blog" \
		-delaystart="5" \
		-desc="$(cat desc)" \
		-exename=railroad \
		-license=MIT \
		-res=config
	cp -v *.su3 ../
	unzip railroad.zip -d railroad-zip

brb-example: all clean
	./scripts/bin/i2p.plugin.native -name=brb \
		-signer=hankhill19580@gmail.com \
		-version 0.0.09 \
		-author=hankhill19580@gmail.com \
		-autostart=true \
		-clientname=brb \
		-command="\$$PLUGIN/lib/brb -dir=\$$PLUGIN/lib -eris=true -i2psite=true 2>&1 \$$PLUGIN/lib/brb.log" \
		-consolename="BRB IRC" \
		-delaystart="5" \
		-desc="$(cat ircdesc)" \
		-exename=brb \
		-license=MIT
	cp -v *.su3 ../
	unzip brb.zip -d brb-zip

clean:
	rm -rf plugin *.zip *.su3

fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;