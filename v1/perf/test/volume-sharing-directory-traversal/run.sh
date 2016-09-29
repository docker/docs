#!/bin/bash

mkdir -p volumes
for test in $(ls testcases); do
  echo $test
  (mkdir -p volumes/$test      \
     && cd volumes/$test       \
     && ../../create_hierarchy ../../testcases/$test)
  docker run -v $(pwd)/volumes:/volumes volume-sharing-directory-traversal /find.sh $test
  mv volumes/*$test.time logs/
done
