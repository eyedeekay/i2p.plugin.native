
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

railroad-example: all
	./scripts/bin/i2p.plugin.native -name=railroad \
		-signer=hankhill19580@gmail.com \
		-version 0.0.031 \
		-author=hankhill19580@gmail.com \
		-autostart=true \
		-clientname=railroad \
		-command=railroad \
		-consolename="Railroad Blog" \
		-delaystart="5" \
		-desc="$(cat desc)" \
		-exename=railroad \
		-license=mit \
		-res=config
#		-icondata=

fmt:
	find . -name '*.go' -exec gofmt -w -s {} \;