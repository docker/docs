#!/bin/sh

# Fetches upstream resources from docker/docker and docker/distribution
# before handing off the site to Jekyll to build
# Relies on the "ENGINE_BRANCH" and "DISTRIBUTION_BRANCH" environment variables,
# which are usually set by the Dockerfile.
: "${ENGINE_BRANCH?No release branch set for docker/docker and docker/cli}"
: "${DISTRIBUTION_BRANCH?No release branch set for docker/distribution}"

# Translate branches for use by svn
engine_svn_branch="branches/${ENGINE_BRANCH}"
if [ "${engine_svn_branch}" = "branches/master" ]; then
	engine_svn_branch=trunk
fi
distribution_svn_branch="branches/${DISTRIBUTION_BRANCH}"
if [ "${distribution_svn_branch}" = "branches/master" ]; then
	distribution_svn_branch=trunk
fi

# Directories to get via SVN. We use this because you can't use git to clone just a portion of a repository
svn co "https://github.com/docker/cli/${engine_svn_branch}/docs/extend"              ./engine/extend || (echo "Failed engine/extend download" && exit 1)
svn co "https://github.com/docker/docker/${engine_svn_branch}/docs/api"              ./engine/api    || (echo "Failed engine/api download" && exit 1)
svn co "https://github.com/docker/distribution/${distribution_svn_branch}/docs/spec" ./registry/spec || (echo "Failed registry/spec download" && exit 1)
svn co "https://github.com/mirantis/compliance/trunk/docs/compliance"                ./compliance    || (echo "Failed docker/compliance download" && exit 1)

# Cleanup svn directories
find . -name .svn -exec rm -rf '{}' \;

# Get the Engine APIs that are in Swagger
# Add a new engine/api/<version>.md file to add a new API version page.
# @TODO stop fetching individual files, onces all API docs are unified upstream in moby/moby
wget --quiet --directory-prefix=./engine/api/ https://raw.githubusercontent.com/docker/docker-ce/v17.06.2-ce/components/engine/api/swagger.yaml; mv ./engine/api/swagger.yaml ./engine/api/v1.30.yaml || (echo "Failed 1.30 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/ https://raw.githubusercontent.com/docker/docker-ce/v17.07.0-ce/components/engine/api/swagger.yaml; mv ./engine/api/swagger.yaml ./engine/api/v1.31.yaml || (echo "Failed 1.31 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/ https://raw.githubusercontent.com/docker/docker-ce/v17.09.1-ce/components/engine/api/swagger.yaml; mv ./engine/api/swagger.yaml ./engine/api/v1.32.yaml || (echo "Failed 1.32 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/ https://raw.githubusercontent.com/docker/docker-ce/v17.10.0-ce/components/engine/api/swagger.yaml; mv ./engine/api/swagger.yaml ./engine/api/v1.33.yaml || (echo "Failed 1.33 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/ https://raw.githubusercontent.com/docker/docker-ce/v17.11.0-ce/components/engine/api/swagger.yaml; mv ./engine/api/swagger.yaml ./engine/api/v1.34.yaml || (echo "Failed 1.34 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/ https://raw.githubusercontent.com/docker/docker-ce/v17.12.1-ce/components/engine/api/swagger.yaml; mv ./engine/api/swagger.yaml ./engine/api/v1.35.yaml || (echo "Failed 1.35 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/ https://raw.githubusercontent.com/docker/docker-ce/v18.02.0-ce/components/engine/api/swagger.yaml; mv ./engine/api/swagger.yaml ./engine/api/v1.36.yaml || (echo "Failed 1.36 swagger download" && exit 1)
wget --quiet --directory-prefix=./engine/api/ https://raw.githubusercontent.com/docker/docker-ce/v18.03.1-ce/components/engine/api/swagger.yaml; mv ./engine/api/swagger.yaml ./engine/api/v1.37.yaml || (echo "Failed 1.37 swagger download" && exit 1)

# Get a few one-off files that we use directly from upstream
wget --quiet --directory-prefix=./engine/                       "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/deprecated.md"                    || (echo "Failed engine/deprecated.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/             "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/reference/builder.md"             || (echo "Failed engine/reference/builder.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/             "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/reference/run.md"                 || (echo "Failed engine/reference/run.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/commandline/ "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/reference/commandline/cli.md"     || (echo "Failed engine/reference/commandline/cli.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/commandline/ "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/reference/commandline/dockerd.md" || (echo "Failed engine/reference/commandline/dockerd.md download" && exit 1)
wget --quiet --directory-prefix=./registry/                     "https://raw.githubusercontent.com/docker/distribution/${DISTRIBUTION_BRANCH}/docs/configuration.md"  || (echo "Failed registry/configuration.md download" && exit 1)

# Remove things we don't want in the build
rm ./registry/spec/api.md.tmpl
