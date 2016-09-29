IMAGE_NAME=test_volumes_sendfile_ones_sixtyfourk

.PHONY: test depends clean

test: depends
	docker run --rm -v /tmp:/tmp -t $(IMAGE_NAME) /sendfile_ones_64K.py

depends: Dockerfile sendfile_ones_64K.py
	docker build -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
