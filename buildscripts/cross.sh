#!/usr/bin/env bash

GOARCH="amd64"

if [[ "${NOTARY_BUILDTAGS}" == *pkcs11* ]]; then
	export CGO_ENABLED=1
else
	export CGO_ENABLED=0
fi


for os in "$@"; do
	export GOOS="${os}"

	if [[ "${GOOS}" == "darwin" ]]; then
		export CC="o64-clang"
		export CXX="o64-clang++"
		# -ldflags=-s:  see https://github.com/golang/go/issues/11994
		export LDFLAGS="${GO_LDFLAGS} -ldflags=-s"
	else
		unset CC
		unset CXX
		LDFLAGS="${GO_LDFLAGS}"
	fi

	mkdir -p "${NOTARYDIR}/cross/${GOOS}/${GOARCH}";
	go build \
		-o "${NOTARYDIR}/cross/${GOOS}/${GOARCH}/notary" \
		-a \
		-tags "${NOTARY_BUILDTAGS}" \
		${LDFLAGS} \
		./cmd/notary;
done
