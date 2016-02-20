"""
Script that automates trusted pull/pushes on different docker versions.
"""

from __future__ import print_function
from collections import OrderedDict
import atexit
import json
import os
import platform
import pwd
import re
import shutil
import subprocess
from tempfile import mkdtemp
from time import time
import urllib
from urlparse import urljoin

# Configuration for testing

# please give the full path to the binary (or if it's on your path, just the
# binary name) for these if you do not want them downloaded, otherwise these
# can be ignored.  Up to you to make sure you are running the correct daemon
# version.
DOCKERS = {
    "1.8": "docker-1.8.3",
    "1.9": "docker-1.9.1",
    "1.10": "docker",
    "1.10.2": "docker-1.10.2.RC1",
}

# delete any of these if you want to specify the docker binaries yourself
DOWNLOAD_DOCKERS = {
    "1.8": ("https://get.docker.com", "docker-1.8.3"),
    "1.9": ("https://get.docker.com", "docker-1.9.1"),
    "1.10": ("https://get.docker.com", "docker-1.10.1")
}

# please replace with private registry if you want to test against a private
# registry
REGISTRY = "docker.io"

# please enter your username if it does not match your shell username, or set the
# environment variable DOCKER_USERNAME
REGISTRY_USERNAME = os.getenv("DOCKER_USERNAME", pwd.getpwuid(os.getuid())[0])

# what you want the testing repo names to be prefixed with
REPO_PREFIX = "docker_test"

# Assumes default docker config dir
DEFAULT_DOCKER_CONFIG = os.path.expanduser("~/.docker")

# Assumes the test will be run with `python misc/dockertest.py` from
# the root of the notary repo after binaries are built
NOTARY_CLIENT = "bin/notary -c cmd/notary/config.json"

# Assumes the trust server will be run using compose
TRUST_SERVER = "https://notary-server:4443"

# ---- setup ----

def download_docker(download_dir="/tmp"):
    """
    Downloads the relevant docker binaries and sets the docker values
    """
    system = platform.system()
    architecture = "x86_64"
    if platform.architecture()[0] != "64bit":
        architecture = "i386"

    downloadfile = urllib.URLopener()
    for version in DOWNLOAD_DOCKERS:
        domain, binary = DOWNLOAD_DOCKERS[version]
        filename = os.path.join(download_dir, binary)
        if not os.path.isfile(filename):
            url = urljoin(
                domain, "/".join(["builds", system, architecture, binary]))
            print("Downloading", url)
            downloadfile.retrieve(url, filename)

        os.chmod(filename, 0755)
        DOCKERS[version] = filename

def setup():
    """
    Ensure we are set up to run the test
    """
    download_docker()
    # copy the docker config dir over so we don't break anything in real docker
    # config directory
    os.mkdir(_TEMP_DOCKER_CONFIG_DIR)
    # copy any docker creds over so we can push
    configfile = os.path.join(_TEMP_DOCKER_CONFIG_DIR, "config.json")
    shutil.copyfile(os.path.join(DEFAULT_DOCKER_CONFIG, "config.json"), configfile)
    # always clean up the config file so creds aren't left in this temp directory
    atexit.register(os.remove, configfile)
    defaulttlsdir = os.path.join(DEFAULT_DOCKER_CONFIG, "tls")
    tlsdir = os.path.join(_TEMP_DOCKER_CONFIG_DIR, "tls")
    if os.path.exists(tlsdir):
        shutil.copytree(defaulttlsdir, tlsdir)

    # make sure that the cert is in the right place for local notary
    if TRUST_SERVER == "https://notary-server:4443":
        tlsdir = os.path.join(tlsdir, "notary-server:4443")
        if not os.path.isdir(tlsdir):
            try:
                shutil.rmtree(tlsdir)  # in case it's not a directory
            except OSError as ex:
                if "No such file or directory" not in str(ex):
                    raise
            os.makedirs(tlsdir)
        cert = os.path.join(tlsdir, "root-ca.crt")
        if not os.path.isfile(cert):
            shutil.copyfile("fixtures/root-ca.crt", cert)

# ---- tests ----

_TEMPDIR = mkdtemp(prefix="docker-version-test")
_TEMP_DOCKER_CONFIG_DIR = os.path.join(_TEMPDIR, "docker-config-dir")
_TRUST_DIR = os.path.join(_TEMP_DOCKER_CONFIG_DIR, "trust")


_ENV = os.environ.copy()
_ENV.update({
    # enable content trust and use our own server
    "DOCKER_CONTENT_TRUST_SERVER": TRUST_SERVER,
    "DOCKER_CONTENT_TRUST": "1",

    # environment variables that notary uses
    "NOTARY_ROOT_PASSPHRASE": "randompass",
    "NOTARY_TARGETS_PASSPHRASE": "randompass",
    "NOTARY_SNAPSHOT_PASSPHRASE": "randompass",

    # environment variables used by current version of docker
    "DOCKER_CONTENT_TRUST_ROOT_PASSPHRASE": "randompass",
    "DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE": "randompass",

    # environment variables used by docker 1.8
    "DOCKER_CONTENT_TRUST_OFFLINE_PASSPHRASE": "randompass",
    "DOCKER_CONTENT_TRUST_TAGGING_PASSPHRASE": "randompass",

    # do not use the default docker config directory
    "DOCKER_CONFIG": _TEMP_DOCKER_CONFIG_DIR
})

_DIGEST_REGEX = re.compile(r"\b[dD]igest: sha256:([0-9a-fA-F]+)\b")
_SIZE_REGEX = re.compile(r"\bsize: ([0-9]+)\b")
_PULL_A_REGEX = re.compile(
    r"Pull \(\d+ of \d+\): .+:(.+)@sha256:([0-9a-fA-F]+)")


def clear_tuf():
    """
    Removes the trusted certificates and TUF metadata in ~/.docker/trust
    """
    try:
        shutil.rmtree(os.path.join(_TRUST_DIR, "trusted_certificates"))
        shutil.rmtree(os.path.join(_TRUST_DIR, "tuf"))
    except OSError as ex:
        if "No such file or directory" not in str(ex):
            raise

def clear_keys():
    """
    Removes the TUF keys in trust directory, since the key format changed
    between versions and can cause problems if testing newer docker versions
    before testing older docker versions.
    """
    try:
        shutil.rmtree(os.path.join(_TRUST_DIR, "private"))
    except OSError as ex:
        if "No such file or directory" not in str(ex):
            raise


def run_cmd(cmd, fileoutput):
    """
    Takes a string command, runs it, and returns the output even if it fails.
    """
    print("$ " + cmd)
    fileoutput.write("$ {0}\n".format(cmd))
    try:
        output = subprocess.check_output(cmd.split(), stderr=subprocess.STDOUT,
                                         env=_ENV)
    except subprocess.CalledProcessError as ex:
        print(ex.output)
        fileoutput.write(ex.output)
        raise
    else:
        if output:
            print(output)
            fileoutput.write(output)
        return output
    finally:
        print()
        fileoutput.write("\n")


def rmi(fout, docker_version, image, tag):
    """
    Ensures that an image is no longer available locally to docker.
    """
    try:
        run_cmd(
            "{0} rmi {1}:{2}".format(DOCKERS[docker_version], image, tag),
            fout)
    except subprocess.CalledProcessError as ex:
        if "could not find image" not in str(ex):
            raise

def assert_equality(actual, expected):
    """
    Assert equality, print nice message
    """
    assert actual == expected, "\nGot     : {0}\nExpected: {1}".format(
        repr(actual), repr(expected))


def pull(fout, docker_version, image, tag, expected_sha):
    """
    Pulls an image using docker, and asserts that the sha is correct.  Make
    sure it is untagged first.
    """
    clear_tuf()
    rmi(fout, docker_version, image, tag)
    output = run_cmd("{0} pull {1}:{2}".format(DOCKERS[docker_version],
                                               image, tag),
                     fout)
    sha = _DIGEST_REGEX.search(output).group(1)
    assert_equality(sha, expected_sha)


def push(fout, docker_version, image, tag):
    """
    Tags an image with the docker version and pushes it.  Returns the sha and
    expected size.
    """
    clear_tuf()

    # tag image with the docker version
    run_cmd(
        "{0} tag -f alpine {1}:{2}".format(DOCKERS[docker_version], image, tag),
        fout)

    # push!
    output = run_cmd("{0} push {1}:{2}".format(DOCKERS[docker_version],
                                               image, tag),
                     fout)
    sha = _DIGEST_REGEX.search(output).group(1)
    size = _SIZE_REGEX.search(output).group(1)

    # list
    targets = notary_list(fout, image)
    for target in targets:
        if target[0] == tag:
            assert_equality(target, [tag, sha, size, "targets"])
    return sha, size


def notary_list(fout, repo):
    """
    Calls notary list on the repo and returns a list of lists of tags, shas,
    sizes, and roles.
    """
    clear_tuf()
    output = run_cmd(
        "{0} -d {1} list {2}".format(NOTARY_CLIENT, _TRUST_DIR, repo), fout)
    lines = output.strip().split("\n")
    assert len(lines) >= 3, "not enough targets"
    return [line.strip().split() for line in lines[2:]]


def test_pull_a(fout, docker_version, image, expected_tags):
    """
    Pull -A on an image and ensure that all the expected tags are present
    """
    clear_tuf()
    # remove every image possible
    for tag in expected_tags:
        rmi(fout, docker_version, image, tag)

    # pull -a
    output = run_cmd(
        "{0} pull -a {1}".format(DOCKERS[docker_version], image), fout)
    pulled_tags = _PULL_A_REGEX.findall(output)

    assert_equality(len(pulled_tags), len(expected_tags))
    for tag, info in expected_tags.iteritems():
        found = [pulled for pulled in pulled_tags if pulled[0] == tag]
        assert found
        assert_equality(found[0][1], info["sha"])


def test_push(tempdir, docker_version, image, tag="", allow_push_failure=False,
              do_after_first_push=None):
    """
    Tests a push of an image by pushing with this docker version, and asserting
    that all the other docker versions can pull it.
    """
    if not tag:
        tag = docker_version

    filename = os.path.join(
        tempdir, "{0}_{1}_push_{2}").format(time(), docker_version, tag)

    with open(filename, 'wb') as fout:
        try:
            sha, size = push(fout, docker_version, image, tag=tag)
        except subprocess.CalledProcessError:
            if allow_push_failure:
                return {"push": "failed, but that was expected"}
            raise

        return_val = {
            "push": {
                "sha": sha,
                "size": size
            }
        }

        if do_after_first_push is not None:
            do_after_first_push(fout, image)

        for ver in DOCKERS:
            try:
                pull(fout, ver, image, tag, sha)
            except subprocess.CalledProcessError:
                print("pulling {0}:{1} with {2} (expected hash {3}) failed".format(
                    image, tag, ver, sha))
                raise
            else:
                return_val["push"][ver] = "pull succeeded"

        return return_val


def test_docker_version(docker_version, repo_name="", do_after_first_push=None):
    """
    Initialize a repo with one docker version.  Test that all other docker
    versions against that repo (both pulling and pushing).
    """
    if not repo_name:
        repo_name = "repo_by_{0}".format(docker_version)
    tempdir = os.path.join(_TEMPDIR, repo_name)
    os.makedirs(tempdir)
    image = "{0}/{1}/{2}_{3}-{4}".format(
        REGISTRY, REGISTRY_USERNAME, REPO_PREFIX, repo_name, time())

    result = OrderedDict([
        (docker_version, test_push(tempdir, docker_version, image,
                                   do_after_first_push=do_after_first_push))
    ])

    # push again if we did something after the first push
    if do_after_first_push:
        tag = docker_version + "_push_again"
        result[tag] = test_push(
            tempdir, docker_version, image, tag=tag,
            # 1.8.x and 1.9.x might fail to push again after snapshot rotation
            # or delegation manipulation
            allow_push_failure=re.compile(r"1\.[0-9](\.\d+)?$").search(docker_version))

    for ver in DOCKERS:
        if ver != docker_version:
            # 1.8.x and 1.9.x will fail to push if the repo was created by
            # a more recent docker, since the key format has changed, or if a
            # snapshot rotation or delegation has occurred
            can_fail = (
                (do_after_first_push or
                 re.compile(r"1\.[1-9][0-9](\.\d+)?$").search(docker_version)) and
                re.compile(r"1\.[0-9](\.\d+)?$").search(ver))

            result[ver] = test_push(tempdir, ver, image, allow_push_failure=can_fail)

    # find all the successfully pushed tags
    expected_tags = {}
    for ver in result:
        if isinstance(result[ver]["push"], dict):
            expected_tags[ver] = result[ver]["push"]

    with open(os.path.join(tempdir, "pull_a"), 'wb') as fout:
        for ver in DOCKERS:
            try:
                test_pull_a(fout, ver, image, expected_tags)
            except subprocess.CalledProcessError:
                result[ver]["pull-a"] = "failed"
            else:
                result[ver]["pull-a"] = "success"

    with open(os.path.join(tempdir, "notary_list"), 'wb') as fout:
        targets = notary_list(fout, image)
        assert_equality(len(targets), len(expected_tags))
        for tag, info in expected_tags.iteritems():
            found = [target for target in targets if target[0] == tag]
            assert found
            assert_equality(
                found[0][1:],
                [info["sha"], info["size"], "targets"])

        result["list"] = "listed expected targets successfully"

    with open(os.path.join(tempdir, "result.json"), 'wb') as fout:
        json.dump(result, fout, indent=2)

    return result


def rotate_to_server_snapshot(fout, image):
    """
    Uses the notary client to rotate the snapshot key to be server-managed.
    """
    run_cmd(
        "{0} -d {1} key rotate {2} -t snapshot -r".format(
            NOTARY_CLIENT, _TRUST_DIR, image),
        fout)
    run_cmd(
        "{0} -d {1} publish {2}".format(NOTARY_CLIENT, _TRUST_DIR, image),
        fout)


def test_all_docker_versions():
    """
    Initialize a repo with each docker version, and test that other docker
    versions can read/write to it.
    """
    print("Output files at", _TEMPDIR)
    results = OrderedDict()
    for docker_version in DOCKERS:
        clear_keys()

        # test with just creating a regular repo
        result = test_docker_version(docker_version)
        print("\nRepo created with docker {0}:".format(docker_version))
        print(json.dumps(result, indent=2))
        results[docker_version] = result

        # do snapshot rotation after creating the repo, and see if it's still ok
        repo_name = "repo_by_{0}_snapshot_rotation".format(docker_version)
        result = test_docker_version(
            docker_version, repo_name=repo_name,
            do_after_first_push=rotate_to_server_snapshot)

        print("\nRepo created with docker {0} and snapshot key rotated:"
              .format(docker_version))
        print(json.dumps(result, indent=2))
        results[docker_version + "_snapshot_rotation"] = result

    with open(os.path.join(_TEMPDIR, "total_results.json"), 'wb') as fout:
        json.dump(results, fout, indent=2)

    print("\n\nFinal results:")
    results["output_dir"] = _TEMPDIR
    print(json.dumps(results, indent=2))


if __name__ == "__main__":
    setup()
    test_all_docker_versions()
