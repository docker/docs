CODE_SIGN_IDENTITY?=please set CODE_SIGN_IDENTITY to your developer key email address

DOCKER_VMNETD?=com.docker.vmnetd
INSTALLER_BUNDLE?=com.docker.docker-installer

BIN_DIR?=/opt/docker/bin/
LOG_DIR?=/opt/docker/var/log

GIT_VERSION := $(shell git describe --abbrev=6 --dirty --always --tags)

ifeq ($V, 1)
	VERBOSE =
else
	VERBOSE = @
endif

include config.mk
-include Makefile.inc

SRC := \
  src/proxy.c \
	src/main.c


OBJ := $(SRC:src/%.c=build/%.o)
DEP := $(OBJ:%.o=%.d)
INC := -Iinclude

CFLAGS += -DVERSION=\"$(GIT_VERSION)\"
LDFLAGS += -mmacosx-version-min=10.10 \
	   -sectcreate __TEXT __info_plist build/$(DOCKER_VMNETD).info.plist \
	   -sectcreate __TEXT __launchd_plist build/$(DOCKER_VMNETD).plist \
	   -framework CoreFoundation \
	   -framework SystemConfiguration

all: build/$(DOCKER_VMNETD).plist $(DOCKER_VMNETD) build/proxy build/client build/client_v1

.PHONY: clean all
.SUFFIXES:

-include $(DEP)

utils: build/utils.o

build:
	@mkdir -p build

$(DOCKER_VMNETD): build/$(DOCKER_VMNETD)
	rm -f $(DOCKER_VMNETD)
	ln -s build/$(DOCKER_VMNETD)

build/%.o: src/%.c
	@echo cc $<
	@mkdir -p $(dir $@)
	$(VERBOSE) $(ENV) $(CC) $(CFLAGS) $(INC) $(DEF) -Wno-format-nonliteral -MMD -MT $@ -MF build/$*.d -o $@ -c $<

.PHONY: build/CODE_SIGN_IDENTITY
build/CODE_SIGN_IDENTITY:
	if [ ! -e build/CODE_SIGN_IDENTITY ] || [ "$(shell cat build/CODE_SIGN_IDENTITY)" != "${CODE_SIGN_IDENTITY}" ]; then \
	  echo "${CODE_SIGN_IDENTITY}" > build/CODE_SIGN_IDENTITY;                  \
	fi

build/client: build/apple_utils.o build/utils.o build/commit.o build/client.o build/protocol.o 
	ld -framework CoreFoundation  -framework SystemConfiguration -o build/client build/apple_utils.o build/utils.o build/commit.o build/protocol.o build/client.o -arch x86_64 -lc

build/client_v1: build/apple_utils.o build/utils.o build/client_v1.o build/protocol.o
	ld -framework CoreFoundation  -framework SystemConfiguration -o build/client_v1 build/apple_utils.o build/utils.o build/commit.o build/protocol.o build/client_v1.o -arch x86_64 -lc

build/proxy: build/apple_utils.o build/utils.o build/commit.o build/proxy.o build/main.o build/protocol.o
	$(VERBOSE) $(ENV) $(LD) $(LDFLAGS) -o build/proxy.sym build/apple_utils.o build/utils.o build/commit.o build/protocol.o build/proxy.o build/main.o
	$(VERBOSE) $(ENV) $(DSYM) build/proxy.sym -o build/proxy.dSYM
	$(VERBOSE) $(ENV) $(STRIP) build/proxy.sym -o $@

build/$(DOCKER_VMNETD): build/apple_utils.o build/utils.o build/commit.o build/proxy.o build/main_launchd.o build/commit.o build/protocol.o build/$(DOCKER_VMNETD).info.plist build/$(DOCKER_VMNETD).plist
	$(VERBOSE) $(ENV) $(LD) $(LDFLAGS) -o build/$(DOCKER_VMNETD).sym build/apple_utils.o build/utils.o build/commit.o build/protocol.o build/proxy.o build/main_launchd.o
	$(VERBOSE) $(ENV) $(DSYM) build/$(DOCKER_VMNETD).sym -o build/$(DOCKER_VMNETD).dSYM
	$(VERBOSE) $(ENV) $(STRIP) build/$(DOCKER_VMNETD).sym -o $@

build/$(DOCKER_VMNETD).plist: $(DOCKER_VMNETD).plist.in
	mkdir -p build
	sed 's|@BINARY@|${BIN_DIR}${DOCKER_VMNETD}|g' $(DOCKER_VMNETD).plist.in | \
	sed 's|@LOG@|${LOG_DIR}|g' > build/$(DOCKER_VMNETD).plist

build/$(DOCKER_VMNETD).info.plist: $(DOCKER_VMNETD).info.plist.in build/CODE_SIGN_IDENTITY
	sed 's|@CODE_SIGN_IDENTITY@|$(shell cat build/CODE_SIGN_IDENTITY)|g' $(DOCKER_VMNETD).info.plist.in | \
	sed 's|@DOCKER_VMNETD_BUNDLE@|${DOCKER_VMNETD}|g' | \
	sed 's|@INSTALLER_BUNDLE@|${INSTALLER_BUNDLE}|g'  > \
	build/$(DOCKER_VMNETD).info.plist

.PHONY: src/commit.c
src/commit.c:
	sed "s|@SHA@|$(shell git show --pretty=%H --no-patch)|g" src/commit.c.in > src/commit.c

clean:
	@rm -rf build
	@rm -f $(DOCKER_VMNETD) src/commit.c
