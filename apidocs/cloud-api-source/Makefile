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
	docker cp $(CONTAINER_NAME):/app/build ../
	rm -r ../docker-cloud
	mv ../build ../docker-cloud
	cp ../docker-cloud/index.html ../layouts/single.html
	# defuse text that looks like a go template
	# commenting out until we know that we're not using Hugo anywhere & can remove
#	grep -rl '{ {' ../docker-cloud/* | xargs sed -i~ "s/{ {/{ {/g"
	#rm source/includes/node.md~
	#rm source/includes/stack.md~
	#rm ../source/index.html~
	#rm ../source/layouts/single.html~
	docker rm $(CONTAINER_NAME)
