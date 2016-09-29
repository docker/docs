.PHONY: test clean

test:
	rm -f /tmp/touch_a_file
	docker run --rm -v /tmp:/tmp busybox touch /tmp/touch_a_file
	ls -l /tmp/touch_a_file

depends:
	docker pull busybox

clean:
