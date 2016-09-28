ORG?=dockerorcadev
PROJ?=ucp-rethinkdb
# TODO - need to decide how we're handling upstream tags - using latest for now
#TAG?=$(shell cd ../../; make print-TAG)
TAG=latest

DOCKER?=docker

image:
	# Build the ucp-rethinkdb-build image. Set the VERSION build arg to the
	# version of RethinkDB you wish to build.
	$(DOCKER) build --build-arg VERSION=2.3.4 $(DOCKER_BUILD_FLAGS) -t ucp-rethinkdb-build -f Dockerfile.build .

	# Create a container from that image, copy out the built RethinkDB binary
	# into the local bin directory, and remove the container.
	$(DOCKER) rm -f ucp-rethinkdb-build || true
	$(DOCKER) create --name ucp-rethinkdb-build ucp-rethinkdb-build
	$(DOCKER) cp ucp-rethinkdb-build:/usr/local/bin/rethinkdb .
	$(DOCKER) rm ucp-rethinkdb-build

	# Build the minimal image.
	$(DOCKER) build $(DOCKER_BUILD_FLAGS) -t $(ORG)/$(PROJ):$(TAG) -f Dockerfile.min .

push:
	$(DOCKER) push $(ORG)/$(PROJ):$(TAG)
