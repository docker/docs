#!/bin/sh
# SUMMARY: Test perl -i in-place edit
# LABELS: !win
# AUTHOR: David Sheets <david.sheets@docker.com>

# Windows filesystems don't support inplace edit (see http://perldoc.perl.org/perldiag.html - Can't do inplace edit without backup)

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    rm "${D4X_LOCAL_TMPDIR}/perl_i_test" || true
    rmdir "${D4X_LOCAL_TMPDIR}/perl_i_test" || true
}
trap clean_up EXIT

docker run --rm -w /tmp -v "${D4X_LOCAL_TMPDIR}:/tmp" alpine sh -c \
       "apk update && apk add perl && echo 'one two' > perl_i_test && \
        perl -i -p -e 's/two/three/g' perl_i_test && \
        [ \"\`cat perl_i_test\`\" = 'one three' ]"

echo "/tmp/perl_i_test:"
cat "${D4X_LOCAL_TMPDIR}/perl_i_test"

[ "$(cat "${D4X_LOCAL_TMPDIR}/perl_i_test")" = 'one three' ]
[ $? -ne 0 ] && exit 1

exit 0
