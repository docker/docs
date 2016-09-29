# Based on https://docs.docker.com/compose/django/
# The sed scripts are a bit hacky but we need to replace
# the DATABASE section
.PHONY: test clean

test: Dockerfile Makefile db.txt docker-compose.yml requirements.txt
	docker-compose run --rm web django-admin.py startproject composeexample .
	sed -ie '/^DATABASES/{n;N;N;N;d;}' composeexample/settings.py
	sed -ie '/^DATABASES/r db.txt' composeexample/settings.py
	docker-compose up -d
	sleep 5
	curl -L http://docker.local:8000 -o - | grep -q "Congratulations"
	docker-compose down --rmi all
	rm -rf composeexample
	rm -f manage.py

clean:
	docker-compose down --rmi all
	rm -rf composeexample
	rm -f manage.py
