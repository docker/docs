"""
"vendors" notary into docker and runs integration tests - then builds the
docker client binary with an API version compatible with the existing
daemon

Usage:
python docker-integration-test.py <directory to build notary>
"""
from __future__ import print_function
import inspect
import os
import re
import shutil
import subprocess
import sys

def this_script_location():
    """
    Returns the absolute path to directory where this script lives, so that
    we don't depend on the CWD for looking for things.
    """
    script_filename = inspect.getfile(inspect.currentframe())
    return os.path.dirname(os.path.abspath(script_filename))

def fake_vendor(docker_dir="docker"):
    """
    "vendors" notary into docker by copying all of notary into the docker
    vendor directory - also appending several lines into the Dockerfile because
    it pulls down notary from github and builds the binaries
    """
    docker_dir = os.path.abspath(os.path.expanduser(docker_dir))
    docker_notary_relpath = "vendor/src/github.com/docker/notary"
    notary_loc = os.path.dirname(this_script_location())  # this script is in misc
    docker_notary_abspath = os.path.join(docker_dir, docker_notary_relpath)

    print("copying notary ({0}) into {1}".format(notary_loc, docker_notary_abspath))

    def ignore_dirs(walked_dir, _):
        """
        Don't vendor everything, particularly not the docker directory
        recursively, if it happened to be in the notary directory
        """
        if walked_dir == notary_loc:
            return [".git", ".cover", "docs", "bin"]
        elif walked_dir == os.path.join(notary_loc, "fixtures"):
            return ["compatibility"]
        elif walked_dir == os.path.dirname(docker_dir):  # don't recursively copy!
            return [os.path.basename(docker_dir)]
        return []

    if os.path.exists(docker_notary_abspath):
        shutil.rmtree(docker_notary_abspath)
    shutil.copytree(
        notary_loc, docker_notary_abspath, symlinks=True, ignore=ignore_dirs)

    # hack this because docker/docker's Dockerfile checks out a particular version of notary
    # based on a tag or SHA, and we want to build based on what was vendored in
    with open(os.path.join(docker_dir, "Dockerfile"), 'a+') as dockerfile:
        dockerfile.write(
            "\n"
            "RUN set -x && "
            "GOPATH=$(pwd)/vendor/src/github.com/docker/notary/Godeps/_workspace:$GOPATH "
            "go build -o /usr/local/bin/notary-server github.com/docker/notary/cmd/notary-server &&"
            "GOPATH=$(pwd)/vendor/src/github.com/docker/notary/Godeps/_workspace:$GOPATH "
            "go build -o /usr/local/bin/notary github.com/docker/notary/cmd/notary")

    # hack the makefile so that we tag the built image as something else so we
    # don't interfere with any other docker test builds
    with open(os.path.join(docker_dir, "Makefile"), 'r') as makefile:
        makefiletext = makefile.read()

    with open(os.path.join(docker_dir, "Makefile"), 'wb') as makefile:
        image_name = os.getenv("DOCKER_TEST_IMAGE_NAME", "notary-docker-vendor-test")
        text = re.sub("^DOCKER_IMAGE := .+$", "DOCKER_IMAGE := {0}".format(image_name),
                      makefiletext, 1, flags=re.M)
        makefile.write(text)

def run_integration_test(docker_dir="docker"):
    """
    Presumes that the fake vendoring has already happened - this runs the
    integration tests.
    """
    docker_dir = os.path.abspath(os.path.expanduser(docker_dir))
    env = os.environ.copy()
    env["TESTFLAGS"] = '-check.f DockerTrustSuite*'
    subprocess.check_call(
        "make test-integration-cli".split(), cwd=docker_dir, env=env)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("\nUsage: python {0} <docker source directory>\n".format(sys.argv[0]))
        sys.exit(1)
    fake_vendor(sys.argv[1])
    run_integration_test(sys.argv[1])
