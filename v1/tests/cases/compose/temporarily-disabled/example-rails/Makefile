# Based on https://docs.docker.com/compose/rails/
.PHONY: test clean

test: Dockerfile Gemfile Makefile database.yml docker-compose.yml
	touch Gemfile.lock
	docker-compose run --rm web rails new . --force --database=postgresql --skip-bundle
	sed -ie "s/# gem 'therubyracer/gem 'therubyracer/" Gemfile
	docker-compose build --force-rm
	cp database.yml config
	docker-compose up -d
	docker-compose run web rake db:create
	sleep 5
	curl -L http://docker.local:3000 -o - | grep -q "aboard"
	docker rm examplerails_web_run_1
	docker-compose down --rmi all
	rm -rf Gemfile.lock Gemfilee README.rdoc Rakefile app bin config config.ru
	rm -rf db lib log public test tmp vendor Gemfile
	git checkout Gemfile

clean:
	docker-compose down --rmi all
	rm -rf Gemfile.lock Gemfilee README.rdoc Rakefile app bin config config.ru
	rm -rf db lib log public test tmp vendor Gemfile
	git checkout Gemfile
