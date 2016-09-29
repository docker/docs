.PHONY: clean all

all: com.docker.hyperkit

com.docker.hyperkit:
	opam config exec -- $(MAKE) -C hyperkit PRI_ADDR_PREFIX=* CONNECT_SOCKET_NAME=@connect
	cp hyperkit/build/com.docker.hyperkit.sym com.docker.hyperkit

# FreeBSD style(9) requires copyright to be within a /*- ... */ block
LICENSE:
	find hyperkit -type f | xargs awk '/^\/\*-/{p=1;print FILENAME ":";print;next} p&&/^.*\*\//{print;print "";p=0};p' > $@

clean:
	@rm -f com.docker.hyperkit
	@rm -f LICENSE
	$(MAKE) -C hyperkit clean
