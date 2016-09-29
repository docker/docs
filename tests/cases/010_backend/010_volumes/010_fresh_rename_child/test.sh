#!/bin/sh
# SUMMARY: Test that new node children are not created with stale parent paths
# LABELS:
# REPEAT:
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    rm -r gah
}
trap clean_up EXIT

# Test code goes here
mkdir -p gah/wtf/grrr
echo "broken?" > gah/wtf/grrr/foo
docker run --rm -v "$PWD/gah:/gah" alpine sh -c "mv /gah/wtf/grrr /gah/sadness && cat /gah/sadness/foo"

exit 0
