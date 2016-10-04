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
	docker cp $(CONTAINER_NAME):/app/build ../cloud-api-output/
	rm -r ../cloud-api-output/docker-cloud
	mv ../cloud-api-output/build ../cloud-api-output/docker-cloud
	cp ../cloud-api-output/docker-cloud/index.html ../cloud-api-output/layouts/single.html
	# defuse text that looks like a go template
	grep -rl '{{' ../cloud-api-output/* | xargs sed -i~ "s/{{/{ {/g"
	rm ../cloud-api-output/docker-cloud/includes/node.html~
	rm ../cloud-api-output/docker-cloud/includes/stack.html~
	rm ../cloud-api-output/docker-cloud/index.html~
	rm ../cloud-api-output/layouts/single.html~
	docker rm $(CONTAINER_NAME)
