.PHONY: test

test:
	rm -rf /tmp/foo
	docker run --rm -v /tmp:/tmp alpine dd if=/dev/zero bs=8k count=1 of=/tmp/foo
	docker run --rm -v /tmp:/tmp alpine dd if=/dev/zero bs=16k count=1 of=/tmp/foo
	docker run --rm -v /tmp:/tmp alpine dd if=/dev/zero bs=32k count=1 of=/tmp/foo
	docker run --rm -v /tmp:/tmp alpine dd if=/dev/zero bs=64k count=1 of=/tmp/foo
	docker run --rm -v /tmp:/tmp alpine dd if=/dev/zero bs=128k count=1 of=/tmp/foo
