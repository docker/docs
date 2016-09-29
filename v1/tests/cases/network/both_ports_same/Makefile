.PHONY: test clean

test: depends
	docker run --rm -p 8081:8081 -t both_ports_same /test.sh 8081 hello &
	sleep 1
	./client.sh both_ports_same hello

depends: Dockerfile test.sh
	docker build -t both_ports_same .

clean:
	docker rmi both_ports_same
