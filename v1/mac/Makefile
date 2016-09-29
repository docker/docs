currentDirectory              := $(dir $(makefilePath))
currentDirectory              := $(patsubst %/,%,${currentDirectory})
buildOutputDirectory          := ${currentDirectory}/build

# Projects Subdirectories
projectSubDirectories         := src/docker-app

# ignoring accidental filesystem filename collisions with rule names
.PHONY : all build $(projectSubDirectories) version

output = > buildlog.txt

all: checkenv build

build:
	@./scripts/make.bash

# check environment
checkenv:
ifndef GOPATH
	$(warning GOPATH is not set. type:)
	$(warning source scripts/setenv)
	$(warning to define a default location under $$HOME/go)
	$(error Aborting)
endif

# install build dependencies
depends:
	@echo "Setting up dependencies ..."
	@echo "Setting up go dependencies ..."
	@$(MAKE) -C ../../ depends ${output}
	@for dir in ${projectSubDirectories}; do \
		$(MAKE) -C $$dir depends ${output}; \
	done

# uninstall build dependencies
cleandepends:
	@echo "Cleaning up dependencies ..."
	@for dir in ${projectSubDirectories}; do \
		$(MAKE) -C $$dir cleandepends ${output}; \
	done

# clean build & dependencies
superclean: cleandepends clean

# clean build
clean:
	@./scripts/make.bash -c
	@for dir in ${projectSubDirectories}; do \
		$(MAKE) -C $$dir clean ${output}; \
	done

# remake
re:
	@./scripts/make.bash -cb

# dev mode
dev:
	@for dir in ${projectSubDirectories}; do \
		$(MAKE) -C $$dir dev ${output}; \
	done

# run
run: all
	@./scripts/make.bash -r

# sign
sign:
	@./scripts/make.bash -s

# package
package:
	@./scripts/make.bash -p

# publish
publish:
	@./scripts/make.bash -u
