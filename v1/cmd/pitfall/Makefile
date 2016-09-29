SOURCES=$(wildcard *.go)
.PHONY=all

all: pitfall

vendor:
	mkdir -p ./vendor
	cp -r ../../vendor ./vendor/

pitfall: Dockerfile $(SOURCES) vendor
	docker build -t pitfall:build . 
	docker run --name pitfall-builder pitfall:build sh
	docker cp pitfall-builder:/go/src/github.com/docker/pinata/v1/cmd/pitfall/pitfall-linux pitfall-linux
	docker cp pitfall-builder:/go/src/github.com/docker/pinata/v1/cmd/pitfall/pitfall-darwin pitfall-darwin
	docker rm pitfall-builder

clean:
	rm -f pitfall-linux pitfall pitfall-darwin .container_built
	rm -rf ./vendor
	docker stop pitfall-builder || true
	docker rm pitfall-builder || true
	docker images -q pitfall:build | xargs docker rmi || true
