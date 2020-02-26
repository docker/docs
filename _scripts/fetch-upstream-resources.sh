#!/bin/bash

# Fetches upstream resources from docker/docker and docker/distribution
# before handing off the site to Jekyll to build
# Relies on the following environment variables which are usually set by
# the Dockerfile. Uncomment them here to override for debugging

# Helper functino to deal with sed differences between osx and Linux
# See https://stackoverflow.com/a/38595160
sedi () {
    sed --version >/dev/null 2>&1 && sed -i -- "$@" || sed -i "" "$@"
}

# Assume non-local mode until we check for -l
LOCAL=0

while getopts ":hl" opt; do
  case ${opt} in
    l ) LOCAL=1
        echo "Running in local mode"
        break
      ;;
    \? ) echo "Usage: $0 [-h] | -l"
         echo "When running in local mode, operates on the current working directory."
         echo "Otherwise, operates on md_source in the scope of the Dockerfile"
         break
      ;;
  esac
done

# Do some sanity-checking to make sure we are running this from the right place
if ! [ -f _config.yml ]; then
  echo "Could not find _config.yml. We may not be in the right place. Bailing."
  exit 1
fi

# Parse some variables from _config.yml and make them available to this script
# This only finds top-level variables with _version in them that don't have any
# leading space. This is brittle!

while read i; do
  # Store the key as a variable name and the value as the variable value
  varname=$(echo "$i" | sed 's/"//g' | awk -F ':' {'print $1'} | tr -d '[:space:]')
  varvalue=$(echo "$i" | sed 's/"//g' | awk -F ':' {'print $2'} | tr -d '[:space:]')
  echo "Setting \$${varname} to $varvalue"
  declare "$varname=$varvalue"
done < <(cat ./_config.yml |grep '_version:' |grep '^[a-z].*')

# Replace variable in toc.yml with value from above
sedi "s/{{ site.latest_stable_docker_engine_api_version }}/$latest_stable_docker_engine_api_version/g" ./_data/toc.yaml

# Distribution
DISTRIBUTION_SVN_BRANCH="branches/release/2.6"
DISTRIBUTION_BRANCH="release/2.6"

# Directories to get via SVN. We use this because you can't use git to clone just a portion of a repository
svn co https://github.com/docker/docker-ce/tags/v18.03.1-ce/components/cli/docs/extend ./engine/extend || (echo "Failed engine/extend download" && exit 1)
svn co https://github.com/docker/docker-ce/tags/v18.03.1-ce/components/engine/docs/api ./engine/api || (echo "Failed engine/api download" && exit 1) # This will only get you the old API MD files 1.18 through 1.24
svn co https://github.com/docker/distribution/"$DISTRIBUTION_SVN_BRANCH"/docs/spec ./registry/spec || (echo "Failed registry/spec download" && exit 1)
svn co https://github.com/docker/compliance/trunk/docs/compliance ./compliance || (echo "Failed docker/compliance download" && exit 1)

# Get the Library docs
svn co https://github.com/docker-library/docs/trunk ./_samples/library || (echo "Failed library download" && exit 1)
# Remove symlinks to maintainer.md because they break jekyll and we don't use em
find ./_samples/library -maxdepth 9  -type l -delete
# Loop through the README.md files, turn them into rich index.md files
FILES=$(find ./_samples/library -type f -name 'README.md')
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
  cat ${curdir}/front-matter.txt ./_samples/boilerplate.txt > ${curdir}/header.txt
  echo {% raw %} >> ${curdir}/header.txt
  cat ${curdir}/header.txt ${curdir}/README.md > ${curdir}/index.md
  echo {% endraw %} >> ${curdir}/index.md
  rm -rf ${curdir}/front-matter.txt
  rm -rf ${curdir}/header.txt
done

# Get the Engine APIs that are in Swagger
# Be careful with the locations on Github for these
# When you change this you need to make sure to copy the previous
# directory into a new one in the docs git and change the index.html
wget --quiet --directory-prefix=./engine/api/v1.25/ https://raw.githubusercontent.com/docker/docker/v1.13.0/api/swagger.yaml                    || (echo "Failed 1.25 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.26/ https://raw.githubusercontent.com/docker/docker/v17.03.0-ce/api/swagger.yaml                || (echo "Failed 1.26 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.27/ https://raw.githubusercontent.com/docker/docker/v17.03.1-ce/api/swagger.yaml                || (echo "Failed 1.27 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.28/ https://raw.githubusercontent.com/docker/docker/v17.04.0-ce/api/swagger.yaml                || (echo "Failed 1.28 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.29/ https://raw.githubusercontent.com/docker/docker/v17.05.0-ce/api/swagger.yaml                || (echo "Failed 1.29 swagger download" && exit 1)
# New location for swagger.yaml for 17.06+
wget --quiet --directory-prefix=./engine/api/v1.30/ https://raw.githubusercontent.com/docker/docker-ce/v17.06.2-ce/components/engine/api/swagger.yaml || (echo "Failed 1.30 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.31/ https://raw.githubusercontent.com/docker/docker-ce/v17.07.0-ce/components/engine/api/swagger.yaml || (echo "Failed 1.31 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.32/ https://raw.githubusercontent.com/docker/docker-ce/v17.09.1-ce/components/engine/api/swagger.yaml || (echo "Failed 1.32 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.33/ https://raw.githubusercontent.com/docker/docker-ce/v17.10.0-ce/components/engine/api/swagger.yaml || (echo "Failed 1.33 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.34/ https://raw.githubusercontent.com/docker/docker-ce/v17.11.0-ce/components/engine/api/swagger.yaml || (echo "Failed 1.34 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.35/ https://raw.githubusercontent.com/docker/docker-ce/v17.12.1-ce/components/engine/api/swagger.yaml || (echo "Failed 1.35 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.36/ https://raw.githubusercontent.com/docker/docker-ce/v18.02.0-ce/components/engine/api/swagger.yaml || (echo "Failed 1.36 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/v1.37/ https://raw.githubusercontent.com/docker/docker-ce/v18.03.1-ce/components/engine/api/swagger.yaml || (echo "Failed 1.37 swagger download" && exit 1)


# Get a few one-off files that we use directly from upstream
wget --quiet --directory-prefix=./engine/                       https://raw.githubusercontent.com/docker/docker-ce/v18.03.1-ce/components/cli/docs/deprecated.md                    || (echo "Failed engine/deprecated.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/             https://raw.githubusercontent.com/docker/docker-ce/v18.03.1-ce/components/cli/docs/reference/builder.md             || (echo "Failed engine/reference/builder.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/             https://raw.githubusercontent.com/docker/docker-ce/v18.03.1-ce/components/cli/docs/reference/run.md                 || (echo "Failed engine/reference/run.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/commandline/ https://raw.githubusercontent.com/docker/docker-ce/v18.03.1-ce/components/cli/docs/reference/commandline/cli.md     || (echo "Failed engine/reference/commandline/cli.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/commandline/ https://raw.githubusercontent.com/docker/docker-ce/v18.03.1-ce/components/cli/docs/reference/commandline/dockerd.md || (echo "Failed engine/reference/commandline/dockerd.md download" && exit 1)
wget --quiet --directory-prefix=./registry/                     https://raw.githubusercontent.com/docker/distribution/"$DISTRIBUTION_BRANCH"/docs/configuration.md                  || (echo "Failed registry/configuration.md download" && exit 1)

# Remove things we don't want in the build
rm ./registry/spec/api.md.tmpl
rm -rf ./apidocs/cloud-api-source
rm -rf ./tests
rm ./_samples/library/index.md
