# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)

GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

.PHONY: all clean binary build

all: clean binary

clean:
	@echo "+ $@"
	@rm -rf "${PREFIX}/bin/garant"

${PREFIX}/bin/garant:
	@echo "+ $@"
	@go build -o $@ ./cmd/garant

binary: ${PREFIX}/bin/garant
	@echo "+ $@"

build:
	@echo "+ $@"
	@docker build -t garant:${GIT_BRANCH} .

minbuild: build
	@echo "+ $@"
	$(eval C_ID := $(shell docker create garant:${GIT_BRANCH}))
	$(eval TMP_DIR := $(shell TMPDIR=. mktemp -dt garant))
	@docker cp ${C_ID}:/go/bin/garant ${TMP_DIR}
	@docker rm ${C_ID}
	@cp ./Dockerfile_min ${TMP_DIR}/Dockerfile
	@cp ./cmd/garant/config.yml ${TMP_DIR}
	@cp ./cmd/garant/signing_key.json ${TMP_DIR}
	@docker build -t garant:${GIT_BRANCH}_min ${TMP_DIR}
	@rm -rf ${TMP_DIR}
