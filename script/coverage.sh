#!/bin/bash
# Generate test coverage statistics for Go packages.
#
# Works around the fact that `go test -coverprofile` currently does not work
# with multiple packages, see https://code.google.com/p/go/issues/detail?id=6909

WORKDIR=.cover
PROFILE="$WORKDIR/cover.out"
MODE=count

testflags="$@"

generate_cover_data() {
    rm -rf "$WORKDIR"
    mkdir "$WORKDIR"

    for pkg in "$@"; do
		# Do not unit test packages with the string "integration" in their package name
		if [[ $pkg == *"integration"* ]]
		then
			continue
		fi
        f="$WORKDIR/$(echo $pkg | tr / -).cover"
        godep go test -p $testflags -test.short -covermode="$MODE" -coverprofile="$f" "$pkg"
    done

    echo "mode: $MODE" >"$PROFILE"
    grep -h -v "^mode:" "$WORKDIR"/*.cover >>"$PROFILE"
}

show_cover_report() {
    godep go tool cover -${1}="$PROFILE"
}

generate_cover_data $(godep go list ./...)
show_cover_report func
show_cover_report html
