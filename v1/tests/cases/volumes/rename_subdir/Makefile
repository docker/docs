.PHONY: test clean

test: depends
	rm -rf /tmp/foo
	docker run --rm -v /tmp:/tmp -t rename_subdir /test.py

depends: Dockerfile test.py
	docker build -t rename_subdir .

clean:
	docker rmi rename_subdir
