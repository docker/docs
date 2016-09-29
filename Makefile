MACOSX_DEPLOYMENT_TARGET?=10.10
REPO_ROOT=$(shell git rev-parse --show-toplevel)
OUTPUT?=$(REPO_ROOT)/v1/mac/build/Docker.app
PROJECT_ROOT?=$(GOPATH)/src/github.com/docker/pinata
CACHE_DIR?=$(REPO_ROOT)/_cache

# Getting version from Info.plist file
# NOTE(aduermael): this won't work on Windows
# It will be updated soon to support both platform
# with tags of the form: "win-.*" and "mac-.*"
plistPath=$(PROJECT_ROOT)/v1/mac/src/docker-app/docker/docker/Info.plist
versionFromPlist=$(shell /usr/libexec/PlistBuddy -c "Print CFBundleShortVersionString" "$(plistPath)" 2> /dev/null)
PARTS=$(subst -, ,$(versionFromPlist))
VERSION=v$(word 1, $(PARTS))-$(word 2, $(PARTS))

# hockeyapp read-only tokens
HA_MAC_TOKEN=bf3c4239192a4511ab54ff5e963d51b1
HA_WIN_TOKEN=64336e7527dc477da596bedfa2804540

# opam flags
OPAMROOT=$(CACHE_DIR)/opam
OPAMDIR=$(REPO_ROOT)/v1/opam
OPAMFLAGS=MACOSX_DEPLOYMENT_TARGET=$(MACOSX_DEPLOYMENT_TARGET) OPAMROOT=$(OPAMROOT) OPAMYES=1 OPAMCOLORS=1 OPAMDIR=$(OPAMDIR) GO15VENDOREXPERIMENT=1
OPAMLIBS=mirage-block-c docker-diagnose osx-daemon osx-hyperkit
OPAMCMDS=osxfs
OPAMSUPPORT=nurse

# TODO: this needs a cleaner solution
BACKENDCMDS=driver.amd64-linux vmnetd osx.hyperkit.linux hyperkit frontend shell #driver.amd64-qemu

LICENSEDIRS=\
	$(PROJECT_ROOT)/v1/opam \
	$(PROJECT_ROOT)/v1/vendor \
	$(PROJECT_ROOT)/v1/uefi \
	$(OPAMROOT) \
	$(PROJECT_ROOT)/v1/cmd/com.docker.hyperkit \
	$(PROJECT_ROOT)/v1/docker_proxy/vendor \
	$(PROJECT_ROOT)/v1/mac/src/docker-app/docker/Carthage/Checkouts \
	$(PROJECT_ROOT)/v1/mac/dependencies/qemu

.PHONY: all depends opam perf OSS-LICENSES dmg dsym-zip clean cacheclean versions go-fmt go-lint go-vet go-test go-depends

all: opam mac
	@

depends: mac-depends qemu-depends opam-depends go-depends moby-depends
	@

clean: opam-clean backend-clean moby-clean
	$(MAKE) -C $(PROJECT_ROOT)/v1/mac clean

cacheclean:
	rm -rf "$(HOME)/.docker-ci-cache/opam"
	rm -rf "$(OPAMROOT)"

OSS-LICENSES:
	$(OPAMFLAGS) v1/opam/opam-licenses $(OPAMCMDS)
	$(MAKE) -C $(PROJECT_ROOT)/v1/cmd/com.docker.hyperkit LICENSE
	$(foreach dir, $(LICENSEDIRS), mkdir -p $(dir);)
	$(PROJECT_ROOT)/v1/mac/scripts/list-licenses $(LICENSEDIRS) > OSS-LICENSES

# opam applications

UPSTREAM=$(shell ls $(OPAMDIR)/repo/packages/upstream | awk -F. '{ print $$1 }')
DEV=$(shell ls $(OPAMDIR)/repo/packages/dev)

opam-depends:
	@brew install opam || true &> /dev/null
	@brew install dylibbundler || true &> /dev/null
	@$(OPAMFLAGS) $(OPAMDIR)/opam-boot
	@$(OPAMFLAGS) opam update -u
	@$(OPAMFLAGS) opam install depext
	@$(OPAMFLAGS) opam depext $(UPSTREAM) $(DEV) &> /dev/null
	@$(OPAMFLAGS) opam install $(UPSTREAM) $(DEV)

opam-lib-clean-%s:
	$(OPAMFLAGS) $(MAKE) -C $(PROJECT_ROOT)/v1/$* clean

opam-cmd-clean-%s:
	$(OPAMFLAGS) $(MAKE) -C $(PROJECT_ROOT)/v1/cmd/com.docker.$* clean

opam-support-clean-%:
	$(OPAMFLAGS) $(MAKE) -C $(PROJECT_ROOT)/support/$* clean

opam-clean: $(OPAMLIBS:%=opam-lib-clean-%s) $(OPAMCMDS:%=opam-cmd-clean-%s) $(OPAMSUPPORT:%=opam-support-clean-%)
	@

opam-lib-%:
	cd $(PROJECT_ROOT)/v1/$* && $(OPAMFLAGS) ./build.sh

opam-cmd-%:
	cd $(PROJECT_ROOT)/v1/cmd/com.docker.$* && $(OPAMFLAGS) ./build.sh

opam-support-%:
	cd $(PROJECT_ROOT)/support/$* && $(OPAMFLAGS) ./build.sh

opam: $(OPAMLIBS:%=opam-lib-%) $(OPAMCMDS:%=opam-cmd-%) $(OPAMSUPPORT:%=opam-support-%) OSS-LICENSES
	@

# backend

backend-cmd-clean-%:
	cd $(PROJECT_ROOT)/v1/cmd/com.docker.$* && $(OPAMFLAGS) $(MAKE) clean

backend-lib-clean-%:
	cd $(PROJECT_ROOT)/v1/$* && $(OPAMFLAGS) $(MAKE) clean

backend-clean: $(BACKENDCMDS:%=backend-cmd-clean-%)
	@

backend-cmd-%:
	cd $(PROJECT_ROOT)/v1/cmd/com.docker.$* && $(OPAMFLAGS) $(MAKE) CACHE_DIR=$(CACHE_DIR)

backend-cmd-vmnetd: opam
backend-cmd-hyperkit: backend-cmd-vmnetd

backend: $(BACKENDCMDS:%=backend-cmd-%)
	@

# moby
moby-depends:
	go get -u github.com/justincormack/regextract

moby:
	cd $(PROJECT_ROOT)/v1/moby && make

moby-clean:
	cd $(PROJECT_ROOT)/v1/moby && make clean

# mac app

mac-depends:
	cd $(PROJECT_ROOT)/v1/mac/scripts && ./make.bash -dy

mac: opam backend moby docker-release
	cd $(PROJECT_ROOT)/v1/mac/scripts && ./make.bash -cby

dmg:
	cd $(PROJECT_ROOT)/v1/mac/scripts && ./make-dmg

dsym-zip:
	cd $(PROJECT_ROOT)/v1/mac/scripts && ./make-dsym-zip

# run Docker.app

dev: opam mac
	rm -rf "$(PROJECT_ROOT)/v1/mac/build"
	rm -rf "$(PROJECT_ROOT)/v1/mac/src/docker-app/build"
	cd $(PROJECT_ROOT)/v1/mac/src/docker-app && make dev

# open Docker.app .xcodeproj

run:
	$(PROJECT_ROOT)/v1/mac/build/Docker.app/Contents/MacOS/Docker

backend-run:
	@$(PROJECT_ROOT)/v1/cmd/com.docker.shell/com.docker.shell -debug -bundle $(PROJECT_ROOT)/v1/mac/build/Docker.app

# tests
lint: go-fmt go-lint go-vet
	# lint test scripts
	brew install shellcheck
	find tests/cases -type f | xargs -L1 file -I | grep 'text/x-shellscript' | cut -f1 -d":" | xargs -L1 shellcheck -e SC2129,SC1090,SC2039

GOPACKAGES = $(eval GOPACKAGES := $(shell cd $(PROJECT_ROOT)/v1 && go list -e ./... | grep -v vendor | grep -v moby))$(GOPACKAGES)


go-depends:
	go get -u github.com/golang/lint/golint

go-fmt:
	@for pkg in $(GOPACKAGES) ; do \
		echo "gofmt $${pkg##*pinata/} ..." ;\
		cd $(PROJECT_ROOT)/$${pkg##*pinata/} ;\
		test -z "$$(gofmt -s -l . 2>&1 | grep -v ^vendor/ | tee /dev/stderr)" || exit 1 ;\
	done

go-lint:
	@for pkg in $(GOPACKAGES) ; do \
		echo "golint $${pkg##*pinata/} ..." ;\
		cd $(PROJECT_ROOT)/$${pkg##*pinata/} ;\
		test -z "$$(golint . 2>&1 | grep -v ^vendor/ | tee /dev/stderr)" || exit 1 ;\
	done

go-vet:
	@cd $(PROJECT_ROOT) && go vet $(GOPACKAGES)

go-test:
	@cd $(PROJECT_ROOT) && for pkg in $(GOPACKAGES) ; do \
		echo "testing $$pkg ..." ;\
		go test -race -v $$pkg ;\
	done

test-depends: opam
	cd $(PROJECT_ROOT)/v1/tests && $(OPAMFLAGS) ./build.sh

test: lint test-depends go-test
	(cd $(PROJECT_ROOT)/tests && ./rt-local -l nostart,checkout -v -x run)

# test-dmg  also tests the dmg - it's assumed that `make dmg` was performed first first
test-dmg: lint test-depends go-test
	(cd $(PROJECT_ROOT)/tests && ./rt-local -l installer,checkout -v -x run)

fulltest:
	(cd $(PROJECT_ROOT)/tests && ./rt-local -l nostart,release,checkout -v -x run)
	PINATA_APP_PATH=$(OUTPUT) $(PROJECT_ROOT)/v1/tests/pinata-rt test -e

perf:
	make -C $(PROJECT_ROOT)/v1/perf

# qemu
QEMUV = 2.4.1
export QEMUV
qemu-depends:
	@mkdir -p $(CACHE_DIR)
	@cd $(PROJECT_ROOT)/v1/cmd/com.docker.driver.amd64-qemu && make depends CACHE_DIR=$(CACHE_DIR)

# upload to HockeyApp

upload:
	@cd $(PROJECT_ROOT)/v1/mac/scripts && ./make.bash -uy

release:
	rm -rf "$(PROJECT_ROOT)/v1/mac/build"
	rm -rf "$(HOME)/.docker-ci-cache"
	rm -rf "$(CACHE_DIR)"
	make depends
	make
	make test
	git tag $(VERSION) -a -m "Release $(VERSION)"
	git push upstream $(VERSION)

versions:
	@echo git tag name: $(VERSION)
	@echo Xcode project version \(Info.plist\): $(versionFromPlist)
	@echo Changelog: $(shell head -n 1 CHANGELOG | cut -f 3 -d" ")
	@echo docker-diagnose: $(shell cat v1/docker-diagnose/src/dockerCli.ml | grep check_version)

release-to-rc:
	@echo "Releasing latest builds to RC (the newest unreleased build will also be downloaded)"
	docker-release --channel rc --arch mac --build latest publish
	docker-release --channel rc --arch win --build latest publish

# docker-release build
docker-release:
	cd $(PROJECT_ROOT)/v1/docker-release && make

# helpful targets for development
logwatch:
	syslog -w -F '$$Time $$Host $$(Sender)[$$(Facility)][$$(PID)]\n<$$((Level)(str))>: $$Message' \
	-k Sender  Seq Docker -o \
	-k Sender  Seq docker -o \
	-k Message Seq Docker -o \
	-k Message Seq docker
