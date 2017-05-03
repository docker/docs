.PHONY: docs build release

CONTAINER_NAME=api-build
export CONTAINER_NAME

docs: build
	docker run --rm -it -p 4567:4567 dockercloud/api-docs

build:
	git submodule init && git submodule update
	docker build -t dockercloud/api-docs .

release: build
	docker rm -f $(CONTAINER_NAME) || true
	docker run --name $(CONTAINER_NAME) -it \
		dockercloud/api-docs middleman build
	docker cp $(CONTAINER_NAME):/app/build ../apidocs/
	rm -r ../apidocs/docker-cloud
	mv ../apidocs/build ../apidocs/docker-cloud
	cp ../apidocs/docker-cloud/index.html ../apidocs/layouts/single.html
	# defuse text that looks like a go template
	grep -rl '{{' ../apidocs/* | xargs sed -i~ "s/{{/{ {/g"
	rm ../apidocs/docker-cloud/includes/node.html~
	rm ../apidocs/docker-cloud/includes/stack.html~
	rm ../apidocs/docker-cloud/index.html~
	rm ../apidocs/layouts/single.html~
	docker rm $(CONTAINER_NAME)
