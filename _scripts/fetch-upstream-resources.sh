#!/bin/bash

# Fetches upstream resources from docker/docker and docker/distribution
# before handing off the site to Jekyll to build
# Relies on the following environment variables which are usually set by
# the Dockerfile. Uncomment them here to override for debugging

# Engine stable
ENGINE_SVN_BRANCH="branches/17.03.x"
ENGINE_BRANCH="17.03.x"

# Distribution
DISTRIBUTION_SVN_BRANCH="branches/release/2.6"
DISTRIBUTION_BRANCH="release/2.6"

# Directories to get via SVN. We use this because you can't use git to clone just a portion of a repository
svn co https://github.com/docker/docker/"$ENGINE_SVN_BRANCH"/docs/extend md_source/engine/extend || (echo "Failed engine/extend download" && exit -1)
svn co https://github.com/docker/docker/"$ENGINE_SVN_BRANCH"/docs/api md_source/engine/api || (echo "Failed engine/api download" && exit -1) # This will only get you the old API MD files 1.18 through 1.24
svn co https://github.com/docker/distribution/"$DISTRIBUTION_SVN_BRANCH"/docs/spec md_source/registry/spec || (echo "Failed registry/spec download" && exit -1)


# Get the Engine APIs that are in Swagger
# Be careful with the locations on Github for these
wget -O md_source/engine/api/v1.25/swagger.yaml https://raw.githubusercontent.com/docker/docker/v1.13.0/api/swagger.yaml || (echo "Failed 1.25 swagger download" && exit -1)
wget -O md_source/engine/api/v1.26/swagger.yaml https://raw.githubusercontent.com/docker/docker/v17.03.0-ce/api/swagger.yaml || (echo "Failed 1.26 swagger download" && exit -1)
wget -O md_source/engine/api/v1.27/swagger.yaml https://raw.githubusercontent.com/docker/docker/v17.03.1-ce/api/swagger.yaml || (echo "Failed 1.27 swagger download" && exit -1)

# Get the Edge API Swagger (we only keep the latest one of these
# When you change this you need to make sure to copy the previous
# directory into a new one in the docs git and change the index.html
wget -O md_source/engine/api/v1.28/swagger.yaml https://raw.githubusercontent.com/docker/docker/v17.04.0-ce/api/swagger.yaml || (echo "Failed 1.28 swagger download or the 1.28 directory doesn't exist" && exit -1)
wget -O md_source/engine/api/v1.29/swagger.yaml https://raw.githubusercontent.com/moby/moby/17.05.x/api/swagger.yaml || (echo "Failed 1.29 swagger download or the 1.29 directory doesn't exist" && exit -1)


# Get a few one-off files that we use directly from upstream
wget -O md_source/engine/reference/builder.md https://raw.githubusercontent.com/docker/docker/"$ENGINE_BRANCH"/docs/reference/builder.md || (echo "Failed engine/reference/builder.md download" && exit -1)
wget -O md_source/engine/reference/run.md https://raw.githubusercontent.com/docker/docker/"$ENGINE_BRANCH"/docs/reference/run.md || (echo "Failed engine/reference/run.md download" && exit -1)
wget -O md_source/edge/engine/reference/run.md https://raw.githubusercontent.com/docker/docker/v17.05.0-ce/docs/reference/run.md || (echo "Failed engine/reference/run.md download" && exit -1)
wget -O md_source/engine/reference/commandline/cli.md https://raw.githubusercontent.com/docker/docker/"$ENGINE_BRANCH"/docs/reference/commandline/cli.md || (echo "Failed engine/reference/commandline/cli.md download" && exit -1)
wget -O md_source/engine/deprecated.md https://raw.githubusercontent.com/docker/docker/"$ENGINE_BRANCH"/docs/deprecated.md || (echo "Failed engine/deprecated.md download" && exit -1)
wget -O md_source/registry/configuration.md https://raw.githubusercontent.com/docker/distribution/"$DISTRIBUTION_BRANCH"/docs/configuration.md || (echo "Failed registry/configuration.md download" && exit -1)

# Remove things we don't want in the build
rm md_source/registry/spec/api.md.tmpl
rm -rf md_source/apidocs/cloud-api-source
rm -rf md_source/tests
