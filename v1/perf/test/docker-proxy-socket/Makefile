.PHONY: test clean

test: results.dat

results.dat: logs/baseline.time logs/test-up.time logs/test-down.time
	./results.sh > $@

copy-up/data:
	dd if=/dev/zero of=copy-up/data bs=1048576 count=100

copy-up/nada: nada.c
	cc -o copy-up/nada nada.c

logs/baseline.time: baseline/Dockerfile
	mkdir -p logs
	sh -c 'docker rm socket9p-build-baseline || true'
	sh -c 'time docker build -t socket9p-build-baseline baseline' 2> logs/baseline.time

logs/test-up.time: copy-up/Dockerfile copy-up/data copy-up/nada
	mkdir -p logs
	ls -lh copy-up/data
	sh -c 'docker rm socket9p-copy-up-perf || true'
	sh -c 'time docker build -t socket9p-copy-up-perf copy-up' 2> logs/test-up.time

logs/test-down.time: logs/test-up.time
	mkdir -p logs
	sh -c 'docker rm socket9p-copy-down-perf || true'
	docker create --name socket9p-copy-down-perf socket9p-copy-up-perf nada
	sh -c 'time docker cp socket9p-copy-down-perf:/data data' 2> logs/test-down.time

clean:
	docker rm socket9p-copy-down-perf || true
	docker rmi socket9p-copy-up-perf || true
	docker rmi socket9p-build-baseline || true
	rm -f copy-up/data copy-up/nada
	rm -f data
	rm -f logs/*.time
	rm -f results.dat
