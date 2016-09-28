PROJ=ucp-compose

ORG=docker
TAG=$(shell cd ../../; make print-TAG)
DOCKER?=docker

image:
	    $(DOCKER) build $(DOCKER_BUILD_FLAGS) -t $(ORG)/$(PROJ):$(TAG) .

push:
	    $(DOCKER) push $(ORG)/$(PROJ):$(TAG)
