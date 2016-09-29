.PHONY: test clean

test:
	docker run --privileged -t alpine:latest /usr/sbin/ntpd -d -n -q -p 192.168.65.1

clean:
