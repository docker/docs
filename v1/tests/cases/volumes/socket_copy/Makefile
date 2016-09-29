IMAGE_NAME=test_volumes_socket_copy

.PHONY: test depends clean

test: depends
	rm -rf /tmp/sock
	echo "twelve bytes" > /tmp/socket_copy_input
	docker run --rm -v /tmp:/tmp -t $(IMAGE_NAME) /socket_copy.sh
	diff /tmp/socket_copy_input /tmp/socket_copy_output

depends: Dockerfile socket_copy.sh
	docker build -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
