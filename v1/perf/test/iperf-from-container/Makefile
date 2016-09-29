
.PHONY: test build depends

test: depends build
	mkdir -p logs
	./run.sh > logs/stdout
	./results.sh > results.dat

build: Dockerfile
	docker build -f Dockerfile -t nettest:iperf .

clean:
	docker images -q nettest:iperf | xargs docker rmi -f
	rm -rf logs
	rm -f results.dat

depends:
	which iperf || brew install iperf
