# Identify the build by a pre-release semantic version string
# Examples:
#   "dev" for dev builds <-- default unless you explicitly set this
#   "rc1" for a Release Candidate 1
#   "ob1234" for CI official build 1234
#   "pr4567" for Github pull request 4567
PRE_RELEASE?=dev
DOCKER?=docker

CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
SHELL=/bin/sh
UCP_VERSION=$(shell egrep "Version[ ]+=" version/version.go | cut -f2 -d= | cut -f2 -d\")
export UCP_VERSION

ifeq "$(strip $(PRE_RELEASE))" ""
TAG?=$(UCP_VERSION)
else
TAG?=$(UCP_VERSION)-$(PRE_RELEASE)
endif
ORG?=docker

# XXX Fix this in the integration framework to use ORG instead of ORCA_ORG
ORCA_ORG?=$(ORG)
export ORCA_ORG

DEV_IMAGE_NAME?=ucp-controller-dev
COMMIT=`git rev-parse --short HEAD`
BUILDTIME=$(shell date -u)
MEDIA_SRCS=$(shell find controller/static/ -type f \
				-not -path "controller/static/dist/*" \
				-not -path "controller/static/semantic/node_modules/*" \
				-not -path "controller/static/node_modules/*")
export TAG

SRC_DIR:=$(shell pwd)
BIN_DIR:=$(SRC_DIR)/bin
export SRC_DIR

SELENIUM_VERSION=2.52.0
SELENIUM_URL?=$(shell if [ -n "$$DOCKER_HOST" ] && docker inspect selenium-ucp-hub >/dev/null 2> /dev/null; then echo "http://$$(echo $$DOCKER_HOST|cut -f3 -d/|cut -f1 -d: ):4444/wd/hub"; elif $(DOCKER) inspect selenium-ucp-hub >/dev/null 2> /dev/null ; then echo "http://$$($(DOCKER) inspect selenium-ucp-hub | jq -r '.[0].NetworkSettings.IPAddress'):4444/wd/hub" ; fi)
export SELENIUM_URL

# TODO clean up all the cruft at the top-level so it's in one place!
CONTROLLER_ROOTS:=controller auth dockerhub pkg project registry utils Godeps version
CONTROLLER_SRC:=$(shell find $(CONTROLLER_ROOTS) -name \*.go -print)
BOOTSTRAP_ROOTS:=bootstrap auth dockerhub pkg project registry utils Godeps version
BOOTSTRAP_SRC:=$(shell find $(BOOTSTRAP_ROOTS) -name \*.go -print)
PROXY_ROOTS:=proxy auth dockerhub pkg project registry utils Godeps version
PROXY_SRC:=$(shell find $(PROXY_ROOTS) -name \*.go -print)
# XXX this isn't totally accurate, but close enough
TOOLS_ROOTS:=tools integration/utils version
TOOLS_SRC:=$(shell find $(TOOLS_ROOTS) -name \*.go -print)

TOOL_NAMES:=$(shell cd tools; ls)
TOOLS:=$(addprefix $(BIN_DIR)/,$(TOOL_NAMES))

DEP_IMAGE_DIRS=$(shell cd images ; find * -maxdepth 1 -type d )

# Override test settings as needed - stress tests can take a long time!
TEST_TIMEOUT ?= 120m
# Our tests can be quite I/O intensive when running locally so err on the side of non-parallel
TEST_PARALLEL ?= 1
# flags for go test. Add "-v" to turn on verbose test output, add "--short" to skip lengthy tests
TEST_FLAGS ?=
MACHINE_PREFIX?=$(USER)
# Narrow the scope as needed
INTEGRATION_TEST_SCOPE ?= "./integration/..."
# selenium hub server
export MACHINE_PREFIX

TEST_MACHINE=$(MACHINE_PREFIX)OrcaTest-00000

INTEGRATION_LOG_FILE=$(MACHINE_PREFIX)integration.log
export INTEGRATION_LOG_FILE

all: media build image bundle

# Build the individual binaries based on source level dependencies
build: $(BIN_DIR)/controller $(BIN_DIR)/ucp-tool $(BIN_DIR)/ucp-proxy $(TOOLS)


$(BIN_DIR)/controller: $(CONTROLLER_SRC)
	@echo "Building bin/$(notdir $@)"
	@cd controller && godep go build -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -X github.com/docker/orca/version.GitCommit=$(COMMIT) -X github.com/docker/orca/version.PreRelease=$(PRE_RELEASE) -X \"github.com/docker/orca/version.BuildTime=$(BUILDTIME)\"" -o $(BIN_DIR)/controller .

# TODO - wire up proper dependencies for these...
$(TOOLS): $(TOOLS_SRC)
	@echo "Building bin/$(notdir $@)"
	@(cd tools/$(notdir $@) && godep go build -o $@ . )

$(BIN_DIR)/ucp-tool: $(BOOTSTRAP_SRC)
	@echo "Building bin/$(notdir $@)"
	@cd bootstrap && godep go build -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -X github.com/docker/orca/version.GitCommit=$(COMMIT) -X github.com/docker/orca/version.PreRelease=$(PRE_RELEASE) -X \"github.com/docker/orca/version.BuildTime=$(BUILDTIME)\"" -o $(BIN_DIR)/ucp-tool .

$(BIN_DIR)/ucp-proxy: $(PROXY_SRC)
	@echo "Building bin/$(notdir $@)"
	@cd proxy && godep go build -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -X github.com/docker/orca/version.GitCommit=$(COMMIT) -X github.com/docker/orca/version.PreRelease=$(PRE_RELEASE)" -o $(BIN_DIR)/ucp-proxy .


# Build the images based on dependencies to the binaries
# TODO - deps still aren't quite right - it *always* builds - figure out why so it does incremental builds properly
image: /var/run/docker.sock $(ORG)/ucp-controller $(ORG)/ucp-auth $(ORG)/ucp-proxy $(ORG)/ucp depimages

bundle:
	@if [ "$(ORG)" = "dockerorcadev" ] ; then \
	    IMAGES=$$($(DOCKER) run --rm $(ORG)/ucp:$(TAG) images --list --image-version dev: |tr '\n' ' ') ; \
	    echo "$(DOCKER) save $(ORG)/ucp:$(TAG) $${IMAGES} > ucp_images_$(TAG).tar" ; \
	    $(DOCKER) save $(ORG)/ucp:$(TAG) $${IMAGES} > ucp_images_$(TAG).tar ; \
	else \
	    IMAGES=$$($(DOCKER) run --rm $(ORG)/ucp:$(TAG) images --list |tr '\n' ' ') ; \
	    echo "$(DOCKER) save $(ORG)/ucp:$(TAG) $${IMAGES} > ucp_images_$(TAG).tar" ; \
	    $(DOCKER) save $(ORG)/ucp:$(TAG) $${IMAGES} > ucp_images_$(TAG).tar ; \
	fi
	@rm -f ucp_images_$(TAG).tar.gz
	gzip -v ucp_images_$(TAG).tar

# TODO - Any way we can wire up the FROM dependency too?
$(ORG)/ucp-controller: media $(BIN_DIR)/controller Dockerfile.controller
	@echo "Building image $@:$(TAG)"
	@$(DOCKER) build $(DOCKER_BUILD_FLAGS) --build-arg "UCP_VERSION=$(UCP_VERSION)" -t $(ORG)/ucp-controller:$(TAG) -f Dockerfile.controller .

$(ORG)/ucp-auth:
	@echo "Building image $@:$(TAG)"
	(cd enzi && $(MAKE) image ORG=$(ORG) REPO=ucp-auth TAG=$(TAG))

$(ORG)/ucp-proxy: $(BIN_DIR)/ucp-proxy Dockerfile.proxy
	@echo Building image $@:$(TAG)
	$(DOCKER) build $(DOCKER_BUILD_FLAGS) -t $(ORG)/ucp-proxy:$(TAG) -f Dockerfile.proxy .

$(ORG)/ucp: $(BIN_DIR)/ucp-tool Dockerfile.bootstrap
	@echo Building image $@:$(TAG)
	$(DOCKER) build $(DOCKER_BUILD_FLAGS) -t $(ORG)/ucp:$(TAG) -f Dockerfile.bootstrap .

# TODO if these stick around long-term would be nice to wire up proper dependencies for incremental builds
depimages:
	@for i in $(DEP_IMAGE_DIRS) ; do \
	    echo "building $(ORG)/ucp-$${i}"; \
	    (cd images/$${i} && $(MAKE) ORG=$(ORG) TAG=$(TAG) ) || exit 1; \
	done


freshen:
	@echo "Pulling base layers to get fresh content"
	docker pull $$(grep "^FROM" Dockerfile.controller | cut -f2 -d' ')
	@for i in $(DEP_IMAGE_DIRS) ; do \
	    (cd images/$${i} && docker pull $$(grep "^FROM" Dockerfile | cut -f2 -d' ') ) || exit 1; \
	done

clean:
	@rm -rf controller/static/dist && \
		rm -rf controller/static/semantic/node_modules && \
		rm -f controller/controller && \
		rm -rf $(BIN_DIR) && \
		rm -f controller/ca-certificates.crt \
		rm -f ucp_images*.tar.gz \
		rm -f integration.log

semantic:
	@(cd controller/static/semantic && \
		npm install && \
		./node_modules/gulp/bin/gulp.js build)

media: controller/static/dist/.bundle_timestamp
controller/static/dist/.bundle_timestamp: $(MEDIA_SRCS)
	@echo "Building media"
	@(cd controller/static && \
		./node_modules/webpack/bin/webpack.js --config webpack.config.release.js) && \
		touch controller/static/dist/.bundle_timestamp

/var/run/docker.sock:
	$(error You must run your container with "-v /var/run/docker.sock:/var/run/docker.sock")

test:
	@echo "Beginning unit run with the following settings:"
	@echo "TEST_FLAGS=$(TEST_FLAGS)  	# \"-v\" verbose output, \"--short\" to skip lengthy integration tests"
	@echo ""
	@unset MACHINE_DRIVER && godep go test $(TEST_FLAGS) $$(go list ./... | grep -v '/vendor/' | grep -v '/integration/')

test-cov:
	@echo "Beginning unit run with coverage calculation and the following settings:"
	@echo "TEST_PARALLEL=$(TEST_PARALLEL)"
	@echo "TEST_VERBOSE=$(TEST_VERBOSE)  # Set to \"-v\" for verbose output"
	@echo ""
	@unset MACHINE_DRIVER && ./script/coverage.sh $(TEST_PARALLEL) $(TEST_VERBOSE)

integration:
	@echo "Beginning integration run with the following settings:"
	@echo "DOCKER_HOST=$(DOCKER_HOST)"
	@echo "INTEGRATION_TEST_SCOPE=$(INTEGRATION_TEST_SCOPE)"
	@echo "MACHINE_CREATE_FLAGS=$(MACHINE_CREATE_FLAGS)"
	@echo "MACHINE_DRIVER=$(MACHINE_DRIVER)"
	@echo "MACHINE_FIXUP_COMMAND=$(MACHINE_FIXUP_COMMAND)"
	@echo "MACHINE_LOCAL=$(MACHINE_LOCAL)  	    	# If non-empty triggers local test run"
	@echo "MACHINE_PREFIX=$(MACHINE_PREFIX)         # Set to differentiate the OrcaTest machine names"
	@echo "ORCA_ORG=$(ORCA_ORG)  			# defaults to docker"
	@echo "PULL_IMAGES=$(PULL_IMAGES)  		    # if non-empty will pull using registry env vars"
	@echo "PURGE_MACHINES=$(PURGE_MACHINES)  	    # if non-empty will purge any lingering test machines at end of run"
	@echo "PRESERVE_TEST_MACHINE=$(PRESERVE_TEST_MACHINE)		# If non-empty preserve machines after run"
	@echo "REGISTRY_USERNAME=$(REGISTRY_USERNAME)  	# and friends required for pulling"
	@echo "SELENIUM_URL=$(SELENIUM_URL)"
	@echo "STRESS_OBJECT_COUNT=$(STRESS_OBJECT_COUNT)  		# override stress test defaults (varies)"
	@echo "SWARM_IMAGE=$(SWARM_IMAGE)  		# override the standard swarm image for testing"
	@echo "TAG=$(TAG)"
	@echo "TEST_PARALLEL=$(TEST_PARALLEL)"
	@echo "TEST_TIMEOUT=$(TEST_TIMEOUT)"
	@echo "TEST_FLAGS=$(TEST_FLAGS)  			# \"-v\" verbose output, \"--short\" to skip lengthy integration tests"
	@echo "USE_TEST_MACHINE=$(USE_TEST_MACHINE)"
	@echo ""
	@if [ -n "$(LOG_DIR)" ] ; then (cd $(LOG_DIR); cat /dev/null > $(INTEGRATION_LOG_FILE) ); else cat /dev/null > $(INTEGRATION_LOG_FILE); fi
	@if [ -n "$${USE_TEST_MACHINE}" ]; then \
	    (eval $$(docker-machine env $(TEST_MACHINE) ); \
	     export MACHINE_LOCAL=1 ; \
	     godep go test -timeout $(TEST_TIMEOUT) -p $(TEST_PARALLEL) $(TEST_FLAGS) $(INTEGRATION_TEST_SCOPE) ); \
	elif [ -n "$${MACHINE_LOCAL}" -o  -n "$${MACHINE_DRIVER}" -o -n "$${DOCKER_HOST}" ]; then \
	    godep go test -timeout $(TEST_TIMEOUT) -p $(TEST_PARALLEL) $(TEST_FLAGS) $(INTEGRATION_TEST_SCOPE); \
	    ret=$$?; \
	    if [ -n "$${PURGE_MACHINES}" ]; then \
	        for m in $$(docker-machine ls | cut -f1 -d' ' | grep "$${MACHINE_PREFIX}OrcaTest-") ; do \
	            echo "Purging left-over test machine $${m}"; \
	            docker-machine rm -f $${m}; \
                done; \
	    fi; \
	    exit $${ret} ;\
	else \
	    echo "ERROR: You must set MACHINE_DRIVER, MACHINE_LOCAL or DOCKER_HOST for integration tests"; \
	    /bin/false; \
	fi

create-test-machine:
	if [ -n "$${MACHINE_DRIVER}" ]; then \
	    docker-machine create --driver $(MACHINE_DRIVER) $(MACHINE_CREATE_FLAGS) $(TEST_MACHINE) ; \
	else \
	    echo "ERROR: You must set MACHINE_DRIVER, to create a machine"; \
	    /bin/false; \
	fi
	if [ -n "$${MACHINE_FIXUP_COMMAND}" ]; then \
	    echo "Fixing the machine by running $${MACHINE_FIXUP_COMMAND}" ; \
	    docker-machine ssh $(TEST_MACHINE) $${MACHINE_FIXUP_COMMAND} ; \
	fi
	@echo "Run the following to use the machine"
	@echo 'eval $$(docker-machine env $$(make print-TEST_MACHINE))'

load-test-machine:
	ORG=$(ORG) $(SRC_DIR)/script/copy_orca_images_machine $(TEST_MACHINE)
	(eval $$(docker-machine env $(TEST_MACHINE) ); docker pull busybox)

clean-test-machine:
	docker-machine rm -y $(TEST_MACHINE)


start-selenium: start-selenium-hub start-selenium-chrome start-selenium-firefox
	@echo "SELENIUM_URL=$(SELENIUM_URL)"

stop-selenium: stop-selenium-hub stop-selenium-chrome stop-selenium-firefox

start-selenium-hub:
	@$(DOCKER) inspect selenium-ucp-hub > /dev/null 2>&1; \
	if [ $$? = 1 ]; then \
	    echo " -> starting selenium hub container"; \
	    $(DOCKER) run -d --name selenium-ucp-hub -p 4444:4444 selenium/hub:"${SELENIUM_VERSION}" > /dev/null; \
	    sleep 3; \
	fi

stop-selenium-hub:
	@$(DOCKER) inspect selenium-ucp-hub > /dev/null 2>&1; \
	if [ $$? = 0 ]; then \
	    $(DOCKER) rm -fv selenium-ucp-hub > /dev/null; \
	fi

start-selenium-firefox:
	@$(DOCKER) inspect selenium-ucp-firefox > /dev/null 2>&1; \
	if [ $$? = 1 ]; then \
	    echo " -> starting selenium firefox container"; \
	    $(DOCKER) run -P -d --name selenium-ucp-firefox --link selenium-ucp-hub:hub selenium/node-firefox-debug:"${SELENIUM_VERSION}" > /dev/null; \
	    sleep 3; \
	fi

stop-selenium-firefox:
	@$(DOCKER) inspect selenium-ucp-firefox > /dev/null 2>&1; \
	if [ $$? = 0 ]; then \
	    $(DOCKER) rm -fv selenium-ucp-firefox > /dev/null; \
	fi

start-selenium-chrome:
	@$(DOCKER) inspect selenium-ucp-chrome > /dev/null 2>&1; \
	if [ $$? = 1 ]; then \
	    echo " -> starting selenium chrome container"; \
	    $(DOCKER) run -P -d --name selenium-ucp-chrome --link selenium-ucp-hub:hub selenium/node-chrome-debug:"${SELENIUM_VERSION}" > /dev/null; \
	    sleep 3; \
	fi

stop-selenium-chrome:
	@$(DOCKER) inspect selenium-ucp-chrome > /dev/null 2>&1; \
	if [ $$? = 0 ]; then \
	    $(DOCKER) rm -fv selenium-ucp-chrome > /dev/null; \
	fi

release: all image push

push: /var/run/docker.sock
	@echo ""
	@echo "About to push $(TAG) images to $(ORG)..."
	@echo ""
	@if [ -n "$${REGISTRY_USERNAME}" -a -n "$${REGISTRY_PASSWORD}" ]; then \
	    $(DOCKER) login -u ${REGISTRY_USERNAME} -p ${REGISTRY_PASSWORD}; \
	fi
	@sleep 2
	$(DOCKER) push $(ORG)/ucp-controller:$(TAG)
	$(DOCKER) push $(ORG)/ucp-auth:$(TAG)
	@for i in $(DEP_IMAGE_DIRS) ; do \
	    echo "Pushing $(ORG)/ucp-$${i}"; \
	    (cd images/$${i} && $(MAKE) ORG=$(ORG) TAG=$(TAG) push) || exit 1; \
	done
	$(DOCKER) push $(ORG)/ucp:$(TAG)

print-%: ; @echo $($*)

.PHONY: all build clean semantic media image test test-cov integration release push release start-selenium stop-selenium start-selenium-hub start-selenium-chrome stop-selenium-hub stop-selenium-chrome $(ORG)/ucp-controller $(ORG)/ucp bundle create-test-machine clean-test-machine load-test-machine freshen

