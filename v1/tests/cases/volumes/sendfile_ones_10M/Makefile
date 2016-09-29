IMAGE_NAME=test_volumes_sendfile_ones_tenmeg

.PHONY: test depends clean

test: depends
	docker run --rm -v /tmp:/tmp -t $(IMAGE_NAME) /sendfile_ones_10M.py

depends: Dockerfile sendfile_ones_10M.py
	docker build -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
