IMAGE_NAME=test_volumes_socket_chmod

.PHONY: test depends clean

test: depends
	rm -rf /tmp/sock
	docker run --rm -v /tmp:/tmp -t $(IMAGE_NAME) /socket_chmod.sh

depends: Dockerfile socket_chmod.sh
	docker build -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
