#!/bin/bash

export DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
export HUBBOSS_CONFIG=$DIR/hubboss-config.ini

if [ "$#" -ne 2 ]; then
  echo "Usage: $(basename $0) <from-version> <to-version>" 
  echo "Example: $(basename $0) 5.0.0 6.0.0" 
  exit 1
fi

$DIR/aws_stage.sh -p mercury-ui -from $1 -to $2 -webs 1 -workers 0 -yes
