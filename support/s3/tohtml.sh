#!/bin/sh -ex

docker build -q -f Dockerfile.tohtml -t tohtml .
docker run -v `pwd`/html:/html tohtml sh -c "rm -rf /html/*"

VERSION=2
ARGS="$ARGS -v `pwd`/formatter:/formatter"
ARGS="$ARGS -v `pwd`/html:/html"

if [ -z "$1" ]; then
  echo Processing all users
  docker run $ARGS \
    -v `pwd`/data/$VERSION:/logs \
    -it tohtml /formatter/process-all-users.sh /logs /html
else
  docker run $ARGS \
    -v `pwd`/data/$VERSION/$1:/logs \
    -it tohtml /formatter/process-one-user.sh /logs /html $1
fi
