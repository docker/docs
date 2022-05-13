ifneq (, $(BUILDX_BIN))
	export BUILDX_CMD = $(BUILDX_BIN)
else ifneq (, $(shell docker buildx version))
	export BUILDX_CMD = docker buildx
else ifneq (, $(shell which buildx))
	export BUILDX_CMD = $(which buildx)
else
	$(error "Buildx is required: https://github.com/docker/buildx#installing")
endif

BUILDX_REPO ?= https://github.com/docker/buildx.git
BUILDX_REF ?= master

# Generate YAML docs from remote bake definition
# Usage BUILDX_REF=v0.7.0 make buildx-yaml
buildx-yaml:
	$(eval $@_TMP_OUT := $(shell mktemp -d -t docs-output.XXXXXXXXXX))
	DOCS_FORMATS=yaml $(BUILDX_CMD) bake --set "*.output=$($@_TMP_OUT)" "$(BUILDX_REPO)#$(BUILDX_REF)" update-docs
	rm -rf ./_data/buildx/*
	cp -R "$($@_TMP_OUT)"/out/reference/*.yaml ./_data/buildx/
	rm -rf $($@_TMP_OUT)/*

# Build website and output to _site folder
release:
	rm -rf _site
	$(BUILDX_CMD) bake release

# Vendor Gemfile.lock
vendor:
	$(BUILDX_CMD) bake vendor

# Deploy website and run it through Docker compose
# Available in your browser at http://localhost:4000
deploy:
	docker compose up --build

.PHONY: buildx-yaml release vendor deploy
