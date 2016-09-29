#!/bin/sh

PORT=$1
RESULT=$2

echo $RESULT | nc -l -p $PORT
