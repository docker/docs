IMAGE_NAME=test_volumes_inotify

.PHONY: test depends clean

test: depends run_inotify.sh events.sh
	rm -rf /tmp/inotify*
	mkdir /tmp/inotify
	./run_inotify.sh $(IMAGE_NAME)

depends: Dockerfile inotify.sh
	docker build -t $(IMAGE_NAME) .

clean:
	docker rmi $(IMAGE_NAME)
