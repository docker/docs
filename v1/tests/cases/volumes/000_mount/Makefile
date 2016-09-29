.PHONY: test depends clean

CMD=echo -e '\n************\n' > sep && \
    mount > mtab && \
    cat sep mtab sep && \
    apk update && \
    apk add util-linux && \
    echo -e '\n**** VM ****\n' && \
    nsenter -t 1 -m cat /etc/mtab && \
    cat sep && \
    nsenter -t 1 ps aucx && \
    cat sep && \
    grep osxfs mtab
VOLS=-v /tmp:/tmp
TEST_MOUNT=docker run --privileged --pid=host --rm $(VOLS) alpine:3.3 sh -c "$(CMD)"

test: depends
	$(TEST_MOUNT)

depends:
	@

clean:
	@
