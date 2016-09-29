.PHONY: test clean

test:
	docker-compose stop
	docker-compose up

depends: docker-compose.yml
	docker-compose pull

clean:
	docker-compose stop
