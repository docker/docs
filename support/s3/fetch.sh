#!/bin/sh

set -e

# $1 is the platform (mac|win)
# $2 is the UUID from the bug report
# $3 is the optional report ID

path=""
dst=""
bucket="docker-pinata-support"

if [ -z "${1}" ]; then
	echo "No platform provided"
	exit 1
fi

case "${1}" in
	win)
		path="incoming/w1"
		;;
	mac)
		path="incoming/2"
		;;
	*)
		echo "${1} is not a valid platform"
		exit 1
esac

if [ -z "${2}" ]; then
	echo "No UUID provided"
	exit 1
fi

path="${path}/${2}"
dst="/logs/${2}"

if [ -n "${3}" ]; then
	path="${path}/${3}"
	dst="${dst}/${3}"	
fi

aws s3 cp "s3://${bucket}/${path}" "${dst}" --recursive --quiet 
