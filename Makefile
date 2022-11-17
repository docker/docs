ifneq (, $(BUILDX_BIN))
	export BUILDX_CMD = $(BUILDX_BIN)
else ifneq (, $(shell docker buildx version))
	export BUILDX_CMD = docker buildx
else ifneq (, $(shell which buildx))
	export BUILDX_CMD = $(which buildx)
else
	$(error "Buildx is required: https://github.com/docker/buildx#installing")
endif

# Build website and output to _site folder
release:
	rm -rf _site
	$(BUILDX_CMD) bake release

# Vendor Gemfile.lock
vendor:
	$(BUILDX_CMD) bake vendor
	
# Run all validators
validate:
	$(BUILDX_CMD) bake validate

# Check for broken links
htmlproofer:
	$(BUILDX_CMD) bake htmlproofer

# Lint tool for markdown files
mdl:
	$(BUILDX_CMD) bake mdl

# Deploy website and run it through Docker compose
# Available in your browser at http://localhost:4000
deploy:
	docker compose up --build

# Used in a Dev Environment container
watch:
	bundle install
	bundle exec jekyll serve --watch --config _config.yml --disable-disk-cache

.PHONY: buildx-yaml release vendor htmlproofer mdl deploy watch
