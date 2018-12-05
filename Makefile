all:
	make -C bootstrap
	make -C sidecartester
	make -C nodejs
	make -C python dockerized

clean:
	make -C bootstrap clean
	make -C sidecartester clean
	make -C nodejs clean
	make -C python clean

.PHONY: all clean
