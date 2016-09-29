
TAG=volume-sharing-directory-traversal

.PHONY: test build

test: build
	mkdir -p logs
	./run.sh > logs/stdout
	./results.sh

build: Dockerfile create_hierarchy
	rm -rf volumes/*
	docker build -f Dockerfile -t ${TAG} .

clean:
	docker images -q ${TAG} | xargs docker rmi -f
	rm -rf logs
	rm -f results.dat
	rm -rf volumes
	rm -f create_hierarchy.cmx create_hierarchy.o create_hierarchy.cmi create_hierarchy

depends:
	which ocaml || brew install ocaml
	docker pull ubuntu

create_hierarchy: create_hierarchy.ml
	ocamlfind opt -o $@ -package sexplib -linkpkg $^
