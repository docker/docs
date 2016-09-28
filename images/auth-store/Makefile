PROJ=ucp-auth-store
TAG?=$(shell cd ../../; make print-TAG)
ORG?=docker
DOCKER?=docker
FROM_IMAGE=$(shell grep "^FROM" Dockerfile | cut -f2 -d' ')
# Note: Upstream version pinned in the Dockerfile

image:
	$(DOCKER) inspect $(FROM_IMAGE) > /dev/null || $(DOCKER) pull $(FROM_IMAGE)
	$(DOCKER) build $(DOCKER_BUILD_FLAGS) -t $(ORG)/$(PROJ):$(TAG) .

push:
	$(DOCKER) push $(ORG)/$(PROJ):$(TAG)
