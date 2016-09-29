.PHONY: test clean

test:
	rm -f /tmp/mode0
	touch /tmp/mode0
	chmod 0 /tmp/mode0
	docker run --rm -v /tmp:/tmp busybox chmod 600 /tmp/mode0
	rm -f /tmp/mode0

depends:
	docker pull busybox

clean:
