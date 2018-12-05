VERSION := 10.13.0

all: nodejs node_modules
node_modules:
	npm i
nodejs: node-v${VERSION}-linux-x64
	ln -s node-v${VERSION}-linux-x64 nodejs
node-v${VERSION}-linux-x64: node-v${VERSION}-linux-x64.tar.xz
	tar -xf node-v${VERSION}-linux-x64.tar.xz
	touch node-v${VERSION}-linux-x64 # update mtime so make behaves nicely
node-v${VERSION}-linux-x64.tar.xz:
	curl -O https://nodejs.org/dist/v${VERSION}/node-v${VERSION}-linux-x64.tar.xz

clean:
	rm -rf nodejs node-v${VERSION}-linux-x64.tar.xz node-v${VERSION}-linux-x64 node_modules
.PHONY: clean
