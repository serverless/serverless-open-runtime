VERSION := 10.13.0

all:
	curl -O https://nodejs.org/dist/v${VERSION}/node-v${VERSION}-linux-x64.tar.xz
	tar -xf node-v${VERSION}-linux-x64.tar.xz
	mv node-v${VERSION}-linux-x64 nodejs
	npm i

.PHONY: all
