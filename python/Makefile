VERSION := 3.7.1

all: /opt/python/bin/python
/opt/python/bin/python: Python-${VERSION}/pybuilddir.txt
	make -C Python-${VERSION} install
Python-${VERSION}/pybuilddir.txt: Python-${VERSION}/Makefile
	make -C Python-${VERSION} -j 5
Python-${VERSION}/Makefile: Python-${VERSION}
	cd Python-${VERSION}/ && ./configure --prefix=/opt/python
Python-${VERSION}: Python-${VERSION}.tar.xz
	tar xf Python-${VERSION}.tar.xz
	touch Python-${VERSION} # update mtime so make behaves nicely
Python-${VERSION}.tar.xz:
	curl -O https://www.python.org/ftp/python/${VERSION}/Python-${VERSION}.tar.xz

clean:
	rm -rf python Python-* Python-*.tar.xz
dockerized:
	docker run --rm -it -v "$(CURDIR)":/opt lambci/lambda:build-python3.7 make -C /opt/
.PHONY: clean dockerized
