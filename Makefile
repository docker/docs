.PHONY: dns server-prod-target server-target server-extras base base-tag prod prod-tag push-builders js-prod images images-prod

# Set up make's
dns:
	./containers/dnsmasq/configure_system_dns.sh
hub-deps:
	git clone git@github.com:docker/docker-ux.git ./private-deps/docker-ux
	git clone git@github.com:docker/hub-js-sdk.git ./private-deps/hub-js-sdk
# -> bootstrap-dev
server-target:
	mkdir -p app/.build/public/styles
	cp -R app/img app/.build/public
styles-base:
	cp ./private-deps/docker-ux/dist/styles/main.css ./app/.build/public/styles/main.css
images:
	cp -R ./private-deps/docker-ux/dist/images ./app/.build/public/
docker-font-dev:
	cp -R ./private-deps/docker-ux/dist/fonts ./app/.build/public/
	cp ./app/fonts/* ./app/.build/public/fonts/
	mkdir -p app/.build/public/styles
	cp ./app/styles/font-awesome.min.css ./app/.build/public/styles/font-awesome.min.css

# Circle make's
dev:
	docker build -t bagel/hub-builders-dev .
dev-tag:
	$(shell docker tag bagel/hub-builders-dev:latest bagel/hub-builders-dev:$(shell git rev-parse --verify HEAD))
local:
	docker build -f local.Dockerfile -t bagel/hub-builders-local .
copy-local:
	docker run --name bagel-local -d bagel/hub-builders-local sleep 50s
	docker cp bagel-local:/opt/hub/.build-prod ./.local/
stage:
	docker build -f dockerfiles/Dockerfile-stage-build -t bagel/hub-builders-stage .
copy-stage:
	docker run --name bagel-stage -d bagel/hub-builders-stage sleep 50s
	docker cp bagel-stage:/opt/hub/.build-prod ./.stage/
prod:
	docker build -f dockerfiles/Dockerfile-prod-build -t bagel/hub-builders-prod .
base-prod-tag:
	$(shell docker tag bagel/hub-builders-prod:latest bagel/hub-builders-prod:$(shell git rev-parse --verify HEAD))
copy-prod:
	docker run --name bagel-prod -d bagel/hub-builders-prod sleep 50s
	docker cp bagel-prod:/opt/hub/.build-prod .

# Dockerfile make's
server-prod-target:
	rm -rf .build-prod
	mkdir -p .build-prod
server-extras:
	cp app-server/package.json .build-prod/package.json
	cp app-server/favicons/favicon-dev.ico .build-prod/favicon.ico
	cp app-server/Dockerfile .build-prod/Dockerfile
js-prod:
	ENV=production webpack --production --config _webpack/webpack.prod.config.js
	ENV=production webpack --production --config _webpack/webpack.server.config.js
js-stage:
	ENV=staging webpack --production --config _webpack/webpack.prod.config.js
	ENV=staging webpack --production --config _webpack/webpack.server.config.js
js-local:
	ENV=local webpack --production --config _webpack/webpack.prod.config.js
	ENV=local webpack --production --config _webpack/webpack.server.config.js
images-prod:
	cp -R ./private-deps/docker-ux/dist/images .build-prod/public/
docker-font-prod:
	cp -R ./private-deps/docker-ux/dist/fonts .build-prod/public
	cp -R ./app/fonts/* .build-prod/public/fonts/
	mkdir -p app/.build-prod/public/styles
	cp ./app/styles/font-awesome.min.css .build-prod/public/styles/font-awesome.min.css
styles-base-prod:
	cp ./private-deps/docker-ux/dist/styles/main.css .build-prod/public/styles/main.css
stats-dir:
	mkdir -p /stats/css
css-stats:
	/opt/hub/node_modules/.bin/cssstats file /opt/hub/.build-prod/public/styles/$(shell cat /tmp/.client-js-hash) > /stats/css-stats.json

# Unused make commands
# Universe commands are no longer used as we now have the universe branch
dev-test-jest:
	docker build -f dockerfiles/Dockerfile-builders-dev-jest -t bagel/hub-builders-dev-jest .
push-builders:
	docker push bagel/hub-builders-prod
	docker push bagel/hub-builders-dev
prod-tag:
	$(shell docker tag bagel/hub-prod:latest bagel/hub-prod:$(shell git rev-parse --verify HEAD))
universe:
#	[ ! "${$(npm -v):0:1}" == "2" ] && echo "please \"npm install -g npm\" to get npm3'" && exit 1
	rm -rf node_modules
	npm install --production
	docker build -f dockerfiles/milky-way -t bagel/milky-way .
	docker build -f dockerfiles/universe -t bagel/universe .
push-universe:
	$(shell docker tag bagel/milky-way:latest bagel/milky-way:$(shell git rev-parse --verify HEAD))
	$(shell docker tag bagel/universe:latest bagel/universe:$(shell git rev-parse --verify HEAD))
	docker push bagel/milky-way
	docker push bagel/universe
new-universe:
	sed -i '.bak' "s/universe:[a-z0-9]*$$/universe:${UNIVERSE_TAG}/" Dockerfile
	sed -i '.bak' "s/universe:[a-z0-9]*$$/universe:${UNIVERSE_TAG}/" local.Dockerfile
	sed -i '.bak' "s/universe:[a-z0-9]*$$/universe:${UNIVERSE_TAG}/" dockerfiles/*
	sed -i '.bak' "s/milky-way:[a-z0-9]*$$/milky-way:${UNIVERSE_TAG}/" dockerfiles/*
