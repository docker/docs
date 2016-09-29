.PHONY: test clean

test: depends
	docker run --rm -t forward_dns /test.sh docker.io MX 8.8.8.8 aspmx.l.google.com

depends: Dockerfile test.sh
	docker build -t forward_dns .

clean:
	docker rmi forward_dns
