# Based on https://docs.docker.com/compose/wordpress/
.PHONY: test clean

test: Dockerfile Makefile docker-compose.yml wp-config.php
	curl https://wordpress.org/latest.tar.gz | tar -xzf -
	cp wp-config.php wordpress/
	docker-compose up -d
	sleep 5
	curl -L http://docker.local:8000 -o - | grep -q "WordPress"
	docker-compose down --rmi all
	rm -rf wordpress

clean:
	docker-compose down --rmi all
	rm -rf wordpress
