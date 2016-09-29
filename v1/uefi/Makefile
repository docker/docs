# By default we use the UEFI.fd which commited. To rebuild use "make rebuild"

.PHONY: rebuild
rebuild: UEFI.fd LICENSE

UEFI.fd LICENSE: Dockerfile
	docker build -t xhyve-uefi:build .
	docker run --rm xhyve-uefi:build cat /UEFI.fd > UEFI.fd
	docker run --rm xhyve-uefi:build cat /LICENSE > LICENSE

.PHONY: clean
clean:
	# XXX git reset UEFI.fd?
	docker images -q xhyve-uefi:build | xargs -r docker rmi -f
