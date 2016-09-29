.PHONY: test clean

test: depends
	docker run --rm -p 8082:8081 -t both_ports_different /test.sh 8081 hello &
	sleep 1
	./client.sh both_ports_different hello

depends: Dockerfile test.sh
	docker build -t both_ports_different .

clean:
	docker rmi both_ports_different
