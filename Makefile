.PHONY: maintainers

maintainers:
	@docker build --rm --force-rm -t docker/maintainers .
	@docker run --rm -v $(CURDIR):/root/maintainers docker/maintainers
