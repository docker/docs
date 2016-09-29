SERVER_CONTAINER := vsock-ping

all: results.png

client: client.c
	cc -Wall -O3 -o $@ $<

server: server.c Dockerfile
	docker build -t $(SERVER_CONTAINER) .
	touch server

.PHONY: results.dat # Always regenerate
results.dat: server client
	./ping.sh $(SERVER_CONTAINER) > results.dat

results.png: results.dat vping.gp
	gnuplot \
	-e "fifty=`../../lib/percentiles.py results.dat 50`" \
	-e "ninety=`../../lib/percentiles.py results.dat 90`" \
	vping.gp

clean:
	rm -f client server
	rm -f results.png results.dat
	docker rm $(SERVER_CONTAINER) || true
