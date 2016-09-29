
TAG=volume-sharing-dd-write

.PHONY: test build

test: build
	mkdir -p logs
	./run.sh | tee logs/stdout
	./results.sh > results.dat

build: Dockerfile
	rm -f volumes/*
	docker build -f Dockerfile -t ${TAG} .

clean:
	docker images -q ${TAG} | xargs docker rmi -f
	rm -rf logs
	rm -f results.dat
	rm -f volumes

depends:
	which ocaml || brew install ocaml
	docker pull ubuntu
