IMAGE_NAME=test_volumes_open_flags

.PHONY: test depends clean

test: depends
	docker run --rm -v /tmp:/tmp -t $(IMAGE_NAME) /open_flags.py

depends: Dockerfile open_flags.py
	docker build -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
