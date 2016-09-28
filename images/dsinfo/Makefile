PROJ=ucp-dsinfo
TAG?=$(shell cd ../../; make print-TAG)
ORG?=docker
DOCKER?=docker

image:
	$(DOCKER) build $(DOCKER_BUILD_FLAGS) -t $(ORG)/$(PROJ):$(TAG) .

push:
	$(DOCKER) push $(ORG)/$(PROJ):$(TAG)
