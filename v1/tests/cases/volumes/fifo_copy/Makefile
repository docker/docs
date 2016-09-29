IMAGE_NAME=test_volumes_fifo_copy

.PHONY: test depends clean

test: depends
	rm -rf /tmp/fifo
	echo "twelve bytes" > /tmp/fifo_copy_input
	docker run --rm -v /tmp:/tmp -t $(IMAGE_NAME) /fifo_copy.sh
	diff /tmp/fifo_copy_input /tmp/fifo_copy_output

depends: Dockerfile fifo_copy.sh
	docker build -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
