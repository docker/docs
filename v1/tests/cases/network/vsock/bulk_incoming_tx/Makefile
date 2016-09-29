.PHONY: test clean

CONTAINER := bulk_incoming_tx

DIR := "$(HOME)/Library/Containers/com.docker.docker/Data"
GUEST := 3
PORT := 1 #$(shell echo $$RANDOM)

BS_KB := 128
COUNT := 16

test: depends
	docker run --rm -v /usr/bin/nc-vsock:/usr/bin/nc-vsock -t --net=host $(CONTAINER) /test.sh $(GUEST) $(PORT) $$(($(BS_KB)*1024)) $(COUNT) &
	sleep 1 # give the server a chance
	./client.sh $(DIR) $(GUEST) $(PORT) $$(($(BS_KB)*1024)) $(COUNT)

depends: Dockerfile test.sh
	which socat || brew install socat
	docker build -t $(CONTAINER) .

clean:
	rm -f received.dat
	docker rmi $(CONTAINER)
