.PHONY: test clean

test: depends
	rm -rf /tmp/foo /tmp/bar
	docker run --rm -v /tmp:/tmp -t rmdir_nonempty /test.sh

depends: Dockerfile test.sh
	docker build -t rmdir_nonempty .

clean:
	docker rmi rmdir_nonempty
