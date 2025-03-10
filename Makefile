# Directory containing the content to validate. Default is "content".
# It can be overridden by setting the CONTENT_DIR environment variable.
# Example: CONTENT_DIR=content/manuals/compose make vale
CONTENT_DIR := $(or $(CONTENT_DIR), content)

# Docker image to use for vale.
VALE_IMAGE := jdkato/vale:latest

.PHONY: vale
vale: ## run vale
	docker run --rm -v $(PWD):/docs \
		-w /docs \
		-e PIP_BREAK_SYSTEM_PACKAGES=1 \
		$(VALE_IMAGE) $(CONTENT_DIR)

.PHONY: validate
validate: ## run validations
	docker buildx bake validate

.PHONY: vendor
vendor: ## vendor hugo modules
	./hack/vendor
