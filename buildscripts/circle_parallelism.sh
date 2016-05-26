 #!/usr/bin/env bash

 case $CIRCLE_NODE_INDEX in
 0) docker run --rm -e NOTARY_BUILDTAGS=pkcs11 notary_client make vet lint fmt misspell
 	docker run --rm -e NOTARY_BUILDTAGS=pkcs11 --env-file buildscripts/env.list --user notary notary_client bash -c "make ci && codecov"
	;;
 1) docker run --rm -e NOTARY_BUILDTAGS=none notary_client make vet lint fmt misspell
 	docker run --rm -e NOTARY_BUILDTAGS=none --env-file buildscripts/env.list --user notary notary_client bash -c "make ci && codecov"
	;;
 2) SKIPENVCHECK=1 make TESTDB=mysql integration
 	;;
 3) SKIPENVCHECK=1 make TESTDB=rethink integration
 	;;
 esac
