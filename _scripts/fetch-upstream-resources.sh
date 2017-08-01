#!/bin/bash

# Fetches upstream resources from docker/docker and docker/distribution
# before handing off the site to Jekyll to build
# Relies on the following environment variables which are usually set by
# the Dockerfile. Uncomment them here to override for debugging

# Engine stable
ENGINE_SVN_BRANCH="branches/17.06"
ENGINE_BRANCH="17.06"
ENGINE_EDGE_BRANCH="17.07"

# Distribution
DISTRIBUTION_SVN_BRANCH="branches/release/2.6"
DISTRIBUTION_BRANCH="release/2.6"

# Directories to get via SVN. We use this because you can't use git to clone just a portion of a repository
svn co https://github.com/docker/docker-ce/"$ENGINE_SVN_BRANCH"/components/cli/docs/extend md_source/engine/extend || (echo "Failed engine/extend download" && exit -1)
svn co https://github.com/docker/docker-ce/"$ENGINE_SVN_BRANCH"/components/engine/docs/api md_source/engine/api || (echo "Failed engine/api download" && exit -1) # This will only get you the old API MD files 1.18 through 1.24
svn co https://github.com/docker/distribution/"$DISTRIBUTION_SVN_BRANCH"/docs/spec md_source/registry/spec || (echo "Failed registry/spec download" && exit -1)

# Get the Library docs
svn co https://github.com/docker-library/docs/trunk md_source/_samples/library || (echo "Failed library download" && exit -1)
# Remove symlinks to maintainer.md because they break jekyll and we don't use em
find md_source/_samples/library -maxdepth 9  -type l -delete
# Loop through the README.md files, turn them into rich index.md files
FILES=$(find md_source/_samples/library -type f -name 'README.md')
for f in $FILES
do
  curdir=$(dirname "${f}")
  justcurdir="${curdir##*/}"
  if [ -e ${curdir}/README-short.txt ]
  then
    # shortrm=$(<${curdir}/README-short.txt)
    shortrm=$(cat ${curdir}/README-short.txt)
  fi
  echo "Adding front-matter to ${f} ..."
  echo --- >> ${curdir}/front-matter.txt
  echo title: "${justcurdir}" >> ${curdir}/front-matter.txt
  echo keywords: library, sample, ${justcurdir} >> ${curdir}/front-matter.txt
  echo repo: "${justcurdir}" >> ${curdir}/front-matter.txt
  echo layout: docs >> ${curdir}/front-matter.txt
  echo permalink: /samples/library/${justcurdir}/ >> ${curdir}/front-matter.txt
  echo redirect_from: >> ${curdir}/front-matter.txt
  echo - /samples/${justcurdir}/ >> ${curdir}/front-matter.txt
  echo description: \| >> ${curdir}/front-matter.txt
  echo \ \ ${shortrm} >> ${curdir}/front-matter.txt
  echo --- >> ${curdir}/front-matter.txt
  echo >> ${curdir}/front-matter.txt
  echo ${shortrm} >> ${curdir}/front-matter.txt
  echo >> ${curdir}/front-matter.txt
  if [ -e ${curdir}/github-repo ]
  then
    # gitrepo=$(<${curdir}/github-repo)
    gitrepo=$(cat ${curdir}/github-repo)
    echo >> ${curdir}/front-matter.txt
    echo GitHub repo: \["${gitrepo}"\]\("${gitrepo}"\)\{: target="_blank"\} >> ${curdir}/front-matter.txt
    echo >> ${curdir}/front-matter.txt
  fi
  cat ${curdir}/front-matter.txt md_source/_samples/boilerplate.txt > ${curdir}/header.txt
  echo {% raw %} >> ${curdir}/header.txt
  cat ${curdir}/header.txt ${curdir}/README.md > ${curdir}/index.md
  echo {% endraw %} >> ${curdir}/index.md
  rm -rf ${curdir}/front-matter.txt
  rm -rf ${curdir}/header.txt
done

# Get the Engine APIs that are in Swagger
# Be careful with the locations on Github for these
wget -O md_source/engine/api/v1.25/swagger.yaml https://raw.githubusercontent.com/docker/docker/v1.13.0/api/swagger.yaml || (echo "Failed 1.25 swagger download" && exit -1)
wget -O md_source/engine/api/v1.26/swagger.yaml https://raw.githubusercontent.com/docker/docker/v17.03.0-ce/api/swagger.yaml || (echo "Failed 1.26 swagger download" && exit -1)
wget -O md_source/engine/api/v1.27/swagger.yaml https://raw.githubusercontent.com/docker/docker/v17.03.1-ce/api/swagger.yaml || (echo "Failed 1.27 swagger download" && exit -1)

# Get the Edge API Swagger (we only keep the latest one of these
# When you change this you need to make sure to copy the previous
# directory into a new one in the docs git and change the index.html
wget -O md_source/engine/api/v1.28/swagger.yaml https://raw.githubusercontent.com/docker/docker/v17.04.0-ce/api/swagger.yaml || (echo "Failed 1.28 swagger download or the 1.28 directory doesn't exist" && exit -1)
wget -O md_source/engine/api/v1.29/swagger.yaml https://raw.githubusercontent.com/docker/docker/17.05.x/api/swagger.yaml || (echo "Failed 1.29 swagger download or the 1.29 directory doesn't exist" && exit -1)
# New location for swagger.yaml for 17.06
wget -O md_source/engine/api/v1.30/swagger.yaml https://raw.githubusercontent.com/docker/docker-ce/17.06/components/engine/api/swagger.yaml || (echo "Failed 1.30 swagger download or the 1.30 directory doesn't exist" && exit -1)

# Get dockerd.md for stable and edge, from upstream
wget -O md_source/engine/reference/commandline/dockerd.md https://raw.githubusercontent.com/docker/docker-ce/"$ENGINE_BRANCH"/components/cli/docs/reference/commandline/dockerd.md || (echo "Failed to fetch stable dockerd.md" && exit -1)
wget -O md_source/edge/engine/reference/commandline/dockerd.md https://raw.githubusercontent.com/docker/docker-ce/"$ENGINE_EDGE_BRANCH"/components/cli/docs/reference/commandline/dockerd.md || (echo "Failed to fetch edge dockerd.md" && exit -1)

# Get a few one-off files that we use directly from upstream
wget -O md_source/engine/reference/builder.md https://raw.githubusercontent.com/docker/docker-ce/"$ENGINE_BRANCH"/components/cli/docs/reference/builder.md || (echo "Failed engine/reference/builder.md download" && exit -1)
wget -O md_source/engine/reference/run.md https://raw.githubusercontent.com/docker/docker-ce/"$ENGINE_BRANCH"/components/cli/docs/reference/run.md || (echo "Failed engine/reference/run.md download" && exit -1)
# Adjust this one when Edge != Stable
wget -O md_source/edge/engine/reference/run.md https://raw.githubusercontent.com/docker/docker-ce/"$ENGINE_BRANCH"/components/cli/docs/reference/run.md || (echo "Failed engine/reference/run.md download" && exit -1)
wget -O md_source/engine/reference/commandline/cli.md https://raw.githubusercontent.com/docker/docker-ce/"$ENGINE_BRANCH"/components/cli/docs/reference/commandline/cli.md || (echo "Failed engine/reference/commandline/cli.md download" && exit -1)
wget -O md_source/engine/deprecated.md https://raw.githubusercontent.com/docker/docker-ce/"$ENGINE_BRANCH"/components/cli/docs/deprecated.md || (echo "Failed engine/deprecated.md download" && exit -1)
wget -O md_source/registry/configuration.md https://raw.githubusercontent.com/docker/distribution/"$DISTRIBUTION_BRANCH"/docs/configuration.md || (echo "Failed registry/configuration.md download" && exit -1)

# Remove things we don't want in the build
rm md_source/registry/spec/api.md.tmpl
rm -rf md_source/apidocs/cloud-api-source
rm -rf md_source/tests
rm md_source/_samples/library/index.md
