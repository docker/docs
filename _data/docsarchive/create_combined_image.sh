#!/bin/bash

ARCHIVE_BRANCHES=("v1.4" "v1.5" "v1.6" "v1.7" "v1.8" "v1.9" "v1.10" "v1.11" "master")
CURRENT_VERSION="v1.12"

rm -rf ./combined
mkdir -p ./combined

echo "Copying current version"
containerid=$(docker run -d docs:latest);
docker cp $containerid:/usr/share/nginx/html/ ./combined/;
docker rm -fv $containerid;

for BRANCH in ${ARCHIVE_BRANCHES[@]}; do
	if [ "$BRANCH" == "master" ]
	then
		# for now, skip current version, as we don't have an archive link for it
		continue
	fi

	VERSION="$BRANCH"
	BASEURL="$VERSION/"
	
	echo "Collecting files for $VERSION"
	mkdir -p ./combined/html/${BASEURL};
	containerid=$(docker run -d docs:${VERSION});
	docker cp $containerid:/usr/share/nginx/html/ ./combined/html/${BASEURL};
	docker rm -fv $containerid;
	echo "Finished copying $VERSION"
done

echo "Building combined version"
docker build -t docs:all -f Dockerfile.combined .

echo "Removing content"
rm -rf ./combined
