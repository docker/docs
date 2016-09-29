CODE_SIGN_IDENTITY?=please set CODE_SIGN_IDENTITY to your developer key email address

# We must compile with -ldflags="-s" to omit
# DWARF info on OSX when compiling with the
# 1.5 or 1.6 (or ...) toolchain. Otherwise the resulting binary
# will be malformed once we codesign it and
# unable to execute.
# See https://github.com/golang/go/issues/11887#issuecomment-126117692.
FLAGS=-ldflags="-s '-extldflags=-sectcreate __TEXT __info_plist ./Info.plist'"

all:
	go build $(FLAGS)

# sign executable with entitlements to make it part of the "group.com.docker" OS X app group
.PHONY: sign
sign:
	@codesign -s "$(CODE_SIGN_IDENTITY)" --entitlements ./docker.entitlements ./com.docker.frontend

.PHONY: signcheck
signcheck:
	@codesign -v ./com.docker.frontend

clean:
	rm -f com.docker.frontend

