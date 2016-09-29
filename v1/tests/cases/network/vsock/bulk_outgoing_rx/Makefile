.PHONY: test clean

CONTAINER := bulk_outgoing_rx

DIR := "$(HOME)/Library/Containers/com.docker.docker/Data"
HOST := 2
PORT := 1

BS_KB := 128
COUNT := 16

test: depends
	./server.sh $(DIR) $(HOST) $(PORT) $$(($(BS_KB)*1024)) $(COUNT) &
	docker run --rm -v /usr/bin/nc-vsock:/usr/bin/nc-vsock -t --net=host $(CONTAINER) /test.sh $(HOST) $(PORT) $$(($(BS_KB)*1024)) $(COUNT)

depends: Dockerfile test.sh
	which socat || brew install socat
	docker build -t $(CONTAINER) .

clean:
	docker rmi $(CONTAINER)
