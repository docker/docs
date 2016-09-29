
.PHONY: test build depends

test: depends build
	mkdir -p logs
	./run.sh > logs/stdout
	./results.sh > results.dat

build: Dockerfile
	docker build -f Dockerfile -t example:hello .

clean:
	docker images -q example:hello | xargs docker rmi -f
	rm -rf logs
	rm -f results.dat

depends:
	which ocaml || brew install ocaml
