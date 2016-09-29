
.PHONY: all list clean
all:
	@echo pre-flight checks starting
	./bin/pre-flight-checks
	@echo pre-flight checks complete
	./bin/run-tests

list:
	./bin/list-tests -v

clean:
	./bin/clean-tests
