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
