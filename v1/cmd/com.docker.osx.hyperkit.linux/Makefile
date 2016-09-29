CODE_SIGN_IDENTITY?=please set CODE_SIGN_IDENTITY to your developer key email address

_build/root/Contents/MacOS/com.docker.osx.hyperkit.linux: com.docker.osx.hyperkit.linux
	mkdir -p _build/root/Contents/MacOS
	cp com.docker.osx.hyperkit.linux _build/root/Contents/MacOS/com.docker.osx.hyperkit.linux
	dylibbundler -od -b \
	  -x _build/root/Contents/MacOS/com.docker.osx.hyperkit.linux \
	  -d _build/root/Contents/Resources/lib \
	  -p @executable_path/../Resources/lib

# We must compile with -ldflags="-s" to omit
# DWARF info on OSX when compiling with the
# 1.5 or 1.6 (or ...) toolchain. Otherwise the resulting binary
# will be malformed once we codesign it and
# unable to execute.
# See https://github.com/golang/go/issues/11887#issuecomment-126117692.
# FLAGS=-ldflags="-s"
FLAGS=-ldflags="-s '-extldflags=-sectcreate __TEXT __info_plist ./Info.plist'"

.PHONY: com.docker.osx.hyperkit.linux
com.docker.osx.hyperkit.linux:
	go build $(FLAGS)

# sign executable with entitlements to make it part of the "group.com.docker" OS X app group
.PHONY: sign
sign:
	@codesign -s "$(CODE_SIGN_IDENTITY)" --entitlements ./docker.entitlements ./_build/root/Contents/MacOS/com.docker.osx.hyperkit.linux

.PHONY: signcheck
signcheck:
	@codesign -v ./_build/root/Contents/MacOS/com.docker.osx.hyperkit.linux

.PHONY: clean
clean:
	rm -f com.docker.osx.hyperkit.linux
