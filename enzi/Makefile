MAKEFILE_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

ORG?=docker
REPO?=enzi
TAG?=dev

UCP_VERSION?=$(shell cd ../../; make print-UCP_VERSION)

image:
	./build.sh
	docker build $(DOCKER_BUILD_FLAGS) --build-arg "UCP_VERSION=$(UCP_VERSION)" -t $(ORG)/$(REPO):$(TAG) .

example-service:
	docker run -t --rm \
		-v $(MAKEFILE_DIR):/go/src/github.com/docker/orca/enzi \
		-w /go/src/github.com/docker/orca/enzi \
		golang:1.6.2-alpine go build -o example-service/example-service github.com/docker/orca/enzi/example-service

	docker build $(DOCKER_BUILD_FLAGS) -t $(ORG)/$(REPO)-example-service:$(TAG) ./example-service

.PHONY: image example-service
