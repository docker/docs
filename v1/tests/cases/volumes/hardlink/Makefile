.PHONY: test clean

test: depends
	rm -rf /tmp/foo /tmp/bar
	docker run --rm -v /tmp:/tmp -t hardlink /test.sh
	./verify.sh

depends: Dockerfile test.sh
	docker build -t hardlink .

clean:
	docker rmi hardlink
