
all: fmt build readme

build:
	go build -o ./scripts/bin/i2p.plugin.native ./scripts/src

readme:
	echo "I2P native plugin generation tool" | tee README.md
	echo "=================================" | tee -a README.md
	echo "" | tee -a README.md
	echo "\`\`\`bash" | tee -a README.md
	./scripts/bin/i2p.plugin.native -h 2>&1 | tee -a README.md
	echo "\`\`\`" | tee -a README.md

railroad-example:
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

brb-example: 
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

#		-icondata=

fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;