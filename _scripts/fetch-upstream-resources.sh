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
find . -name ".svn" -print0 | xargs -0 /bin/rm -rf

# Get a few one-off files that we use directly from upstream
wget --quiet --directory-prefix=./engine/                       "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/deprecated.md"                    || (echo "Failed engine/deprecated.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/             "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/reference/builder.md"             || (echo "Failed engine/reference/builder.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/             "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/reference/run.md"                 || (echo "Failed engine/reference/run.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/commandline/ "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/reference/commandline/cli.md"     || (echo "Failed engine/reference/commandline/cli.md download" && exit 1)
wget --quiet --directory-prefix=./engine/reference/commandline/ "https://raw.githubusercontent.com/docker/cli/${ENGINE_BRANCH}/docs/reference/commandline/dockerd.md" || (echo "Failed engine/reference/commandline/dockerd.md download" && exit 1)
wget --quiet --directory-prefix=./registry/                     "https://raw.githubusercontent.com/docker/distribution/${DISTRIBUTION_BRANCH}/docs/configuration.md"  || (echo "Failed registry/configuration.md download" && exit 1)

# Remove things we don't want in the build
rm -f ./engine/extend/cli_plugins.md # the cli plugins api is not a stable API, and not included in the TOC for that reason.
rm -f ./registry/spec/api.md.tmpl
