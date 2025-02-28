.PHONY: validate
validate: ## run validations
	docker buildx bake validate

.PHONY: vendor
vendor: ## vendor hugo modules
	./hack/vendor
