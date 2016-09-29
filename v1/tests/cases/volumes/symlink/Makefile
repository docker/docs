.PHONY: test clean

test: depends
	rm -rf /tmp/foo /tmp/bar
	docker run --rm -v /tmp:/tmp -t symlink /test.sh
	./verify.sh

depends: Dockerfile test.sh
	docker build -t symlink .

clean:
	docker rmi symlink
