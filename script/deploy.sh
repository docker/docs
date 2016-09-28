#!/bin/bash

# Try several times to do the build (sometimes network or other issues causes
# it to fail)
for i in $(seq 0 4); do
    echo "Trying build $i..."
    make clean
    make
    if [[ $? -eq 0 ]]; then
        cp dist/* $CIRCLE_ARTIFACTS

        if [ ! -z "$CIRCLE_TAG" ]; then
            curl --header "Content-Type: application/json" \
                 --data "{\"build_parameters\": {\"CIRCLE_TAG\": \"$CIRCLE_TAG\", \"TOOLBOX_BUILD_NUM\": $CIRCLE_BUILD_NUM, \"TOOLBOX_ARTIFACTS\": \"$CIRCLE_ARTIFACTS\"}}" \
                 --request POST "https://circleci.com/api/v1/project/docker/toolbox-release/tree/master?circle-token=$CIRCLE_TOKEN"
        fi

        exit 0
    fi
done

exit 1
