all:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) images $(DOIT_VARS)'

include Makefile.variable

# the prints are for back-compat with the jenkins build images script
print-DOCKER_HUB_ORG:
	@shared/inpython3-notty.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) print-DOCKER_HUB_ORG $(DOIT_VARS) 2>/dev/null | grep OUTPUT: | cut -c9-'

print-BOOTSTRAP_NAME:
	@shared/inpython3-notty.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) print-BOOTSTRAP_NAME $(DOIT_VARS) 2>/dev/null | grep OUTPUT: | cut -c9-'

print-VERSION:
	@shared/inpython3-notty.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) print-VERSION $(DOIT_VARS) 2>/dev/null | grep OUTPUT: | cut -c9-'

api:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) api_img $(DOIT_VARS)'

nginx:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) nginx_img $(DOIT_VARS)'

rethink:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) rethink_img $(DOIT_VARS)'

registry:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) registry_img $(DOIT_VARS)'

bootstrap:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) bootstrap_img $(DOIT_VARS)'

notary_server:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) notary_server_img $(DOIT_VARS)'

notary_signer:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) notary_signer_img $(DOIT_VARS)'

jobrunner:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) jobrunner_img $(DOIT_VARS)'

api_ui:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) api_ui $(DOIT_VARS)'

moshpit:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) moshpit_img $(DOIT_VARS)'

run_moshpit_dropper:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) run_moshpit_dropper $(DOIT_VARS)'

clean:
	shared/inpython3.sh 'doit forget $(DOIT_DB_PARAMS)'
	$(RM) -r $(MAKEFILE_DIR)/coverage $(MAKEFILE_DIR)/go-build-cache
	docker rm -f golistcontainer || true

makedoitimg:
	cd shared; docker build -t dockerhubenterprise/doit:withscript -f DoitDockerfile .

gen-mocks:
	./gen_mocks.sh

makeandpushdoitimg: makedoitimg
	docker push dockerhubenterprise/doit:withscript

push: all
	eval "$$(shared/inpython3-notty.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) print-VERSION $(DOIT_VARS) 2>/dev/null | grep OUTPUT: | cut -c9-')"; \
		for reponame in $(shell shared/inpython3-notty.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) print_components $(DOIT_VARS) 2>/dev/null | grep OUTPUT: | cut -c9-') ; do \
		  (docker push $$reponame:$$VERSION || exit $$?) && \
		  ([ ! -z "$(PUSH_NO_LATEST)" ] || docker push $$reponame:latest || exit $$?) \
		done

tar:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) tar $(DOIT_VARS)'

pull:
	eval "$$(shared/inpython3-notty.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) print-VERSION $(DOIT_VARS) 2>/dev/null | grep OUTPUT: | cut -c9-')"; \
		for reponame in $(shell shared/inpython3-notty.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) print_components $(DOIT_VARS) 2>/dev/null | grep OUTPUT: | cut -c9-') ; do \
			docker pull $$reponame:$$VERSION; \
		done

fmt:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) gofmt $(DOIT_VARS)'

test-ginkgo:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) test_ginkgo $(DOIT_VARS)'

test:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) test $(DOIT_VARS)'

build-integration:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) build_integration $(DOIT_VARS)'

test-integration:
	env > $(MAKEFILE_DIR)/integration/local.env
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) test_integration $(DOIT_VARS)'

# Default test settings
# Running the tests requires go install github.com/onsi/ginkgo/...
export TEST_TIMEOUT ?= 120m
export TEST_ARGS ?= "-timeout 100m -v ./ha-integration"
export UCP_REPO ?= docker
export UCP_TAG ?= 1.0.1
export DTR_REPO ?= dockerhubenterprise
export DTR_TAG ?= latest
ha-integration:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) $(DOIT_DB_PARAMS) update_constants $(DOIT_VARS)'
	@echo "Beginning HA integration test with the following settings:"
	@echo "NUM_MACHINES=$(NUM_MACHINES)				# Number of HA nodes and instances to provision"
	@echo "TEST_ARGS=$(TEST_ARGS)"
	@echo "MACHINE_DRIVER=$(MACHINE_DRIVER)"
	@echo "MACHINE_CREATE_FLAGS=$(MACHINE_CREATE_FLAGS)"
	@echo "MACHINE_PREFIX=$(MACHINE_PREFIX)         # Set to differentiate the DTRTest machine names"
	@echo "PULL_IMAGES=$(PULL_IMAGES)  		   # if non-empty will pull using registry env vars"
	@echo "PURGE_MACHINES=$(PURGE_MACHINES)  	    # if non-empty will purge any lingering test machines at end of run"
	@echo "REGISTRY_USERNAME=$(REGISTRY_USERNAME)  	# required for pulling"
	@echo "REGISTRY_PASSWORD=<omitted>"
	@echo "REGISTRY_EMAIL=$(REGISTRY_EMAIL)"
	@echo "UCP_REPO=$(UCP_REPO)"
	@echo "UCP_TAG=$(UCP_TAG)"
	@echo "DTR_REPO=$(DTR_REPO)"
	@echo "DTR_TAG=$(DTR_TAG)"
	@echo "FORCE_PURGE=$(FORCE_PURGE)"	# if you set this all containers and volumes will be deleted and the deamon will be reset to it's unconfigured state
	@echo "USE_PRIVATE_IP=$(USE_PRIVATE_IP)			# uses private ips to perform the setup"
	@echo ""
	if [ -n "$${NUM_MACHINES}" -o  -n "$${MACHINE_DRIVER}" ]; then \
	    GO15VENDOREXPERIMENT=1 ginkgo $(TEST_ARGS); \
	    ret=$$?; \
	    if [ -n "$${PURGE_MACHINES}" ]; then \
	        for m in $$(docker-machine ls | cut -f1 -d' ' | grep "$${MACHINE_PREFIX}-DTRTest-") ; do \
	            echo "Purging left-over test machine $${m}"; \
	            docker-machine rm -f $${m}; \
                done; \
	    fi; \
	    exit $${ret} ;\
	else \
	    echo "ERROR: You must set MACHINE_DRIVER and NUM_MACHINES for HA integration tests"; \
	    /bin/false; \
	fi

local-deploy:
	./local_deploy.sh

# Generates swagger API docs for docs.docker.com
api-docs:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) api_docs_gen $(DOIT_VARS)'

# Generates docs for all commands in the docker/dtr image
# Set the GA release channel before running this to not document hidden commands
cli-docs:
	shared/inpython3.sh 'doit --verbosity 2 $(DOIT_PARAMS) clidocs $(DOIT_VARS)'

# Generate and serve docs as they're available on docs.docker.com
docs:
	$(MAKE) -C docs docs

.PHONY: all bootstrap api api_ui nginx jobrunner notary_server notary_signer rethink registry moshpit clean push tar pull test test-integration coverage docs make-coverage-image api-docs ha-integration makedoitimg makeandpushdoitimg
