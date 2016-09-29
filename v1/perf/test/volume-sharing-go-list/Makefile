.PHONY: all

all: Dockerfile go-list.sh
	rm -rf machine
	git clone "https://github.com/docker/machine.git"
	docker build -t volume-sharing-go-list .
	docker run -it --rm \
		-v `pwd`/machine:/machine \
		volume-sharing-go-list \
		time ./go-list.sh
