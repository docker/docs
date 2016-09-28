+++
draft = "true"
+++

# Docker UCP Integration Tests

The UCP integration tests are built upon Docker Machine and the Go test
library github.com/stretchr/testify.

The integration tests live in integration/ and wont run by
default with go test, unless you set env vars to enable them.

A full list of integration test environment variables is listed at the
end of this document.

There are currently three modes of running the tests.  Fresh machines
for each suite, reusing a single engine, or "bring your own" UCP.

## Fresh Machines for each Suite

This is the preferred mode, although it is slower.  Each suite gets
a fresh machine so there's no leakage between runs.  It also allows
for multi-node testing.

When running in this mode, you can use parallel execution of tests to
speed things up.  Each suite will run in parallel, and use a unique
machine.

```bash
export MACHINE_CREATE_FLAGS="--kvm-memory 2048 --kvm-disk-size 8000 --kvm-cache-mode unsafe"
export MACHINE_DRIVER=kvm

make integration
```

If you're developing your own tests or want to focus on a specific test-case
you can run them explicitly using something like this

```bash
make integration TEST_FLAGS=-v INTEGRATION_TEST_SCOPE=./integration/install/simple/...
```

To enable both verbose test output and skip lengthy test runs, set `TEST_FLAGS="-v --short"`

To run a single test or a subset of a full test suite, append the `-m` go flag at the INTEGRATION_TEST_SCOPE, as follows:

```bash
make integration TEST_FLAGS=-v INTEGRATION_TEST_SCOPE="./integration/... -m TestWebTests"
```


## Use a Single Local Machine

This mode is faster, but failures can cascade, the tests can't run in
parallel, and multi-node tests wont work.  It is **highly** advised
that you set up a machine using docker-machine, wire up your environment
pointing to it, and **then** run the tests against that machine.  To run
this mode use `make integration MACHINE_LOCAL=1` instead of setting
`MACHINE_DRIVER`


## BYO Orca

Download a bundle as an admin, source it for your environment, then the following
will run the BYO scenarios:

```bash
make integration INTEGRATION_TEST_SCOPE=./integration/byoorca/...
```

## DTR Integration Tests

In order to run the UCP-DTR Integration test located in `./integration/install/dtr`, the following environment variables need to be specified for a target DTR Test Rig: `DTR_URL`, `DTR_USERNAME`, `DTR_PASSWORD`, `DTR_PRIVATE_REPO`.

## UI testing with Mac

Someone else to fill this section in ;-)

## UI testing with KVM

The KVM machine driver uses a private network as the primary network
to communicate with the VM.  As a result the IP address used throughout
the integration suite is not routable externally **or** from a container
running directly on your linux host.  As a result, the automatic detection
logic within the Makefile/suite wont quite work, but fear not.  Just run
selenium on another "fixed" VM on your dev box, and point the integration
tests at it.

```
docker-machine create --driver kvm selenium
(eval $(docker-machine env selenium); make start-selenium)
```

The above will ultimately spit out a `SELENIUM_URL=http://some_ip_here:4444/wd/hub`

Then to run the suite, assuming you've set your `MACHINE_XXX` var's as
described above use something like the following (substitute settings
as appropriate for your scenario - see table below)

```
make integration SELENIUM_URL=http://192.168.42.136:4444/wd/hub \
    INTEGRATION_TEST_SCOPE=./integration/install/simple/... \
    TEST_FLAGS=-v
```


# Commmon usage patterns

By default the tests use your local images, and copy them to the test
machine.  This works well for developer mode, but if you're trying to test
official builds, you'll need to set a few variables.

* **Testing Specific CI Builds**
    ```
export REGISTRY_USERNAME=yourname
read -s REGISTRY_PASSWORD
export REGISTRY_PASSWORD

make integration ORCA_ORG=dockerorcadev PULL_IMAGES=1 TAG=1.1.0-ob1234
```

* **Testing Latest CI Builds**
    ```
export REGISTRY_USERNAME=yourname
read -s REGISTRY_PASSWORD
export REGISTRY_PASSWORD

make integration ORCA_ORG=dockerorcadev PULL_IMAGES=1 TAG=latest
```

* **Testing Latest Customer Builds**
    ```
export REGISTRY_USERNAME=yourname
read -s REGISTRY_PASSWORD
export REGISTRY_PASSWORD

make integration PULL_IMAGES=1 TAG=latest
```

* **Testing specific Customer Builds**
    ```
export REGISTRY_USERNAME=yourname
read -s REGISTRY_PASSWORD
export REGISTRY_PASSWORD

make integration PULL_IMAGES=1 TAG=1.0.1
```

* **Testing PR Builds**
    ```
export REGISTRY_USERNAME=yourname
read -s REGISTRY_PASSWORD
export REGISTRY_PASSWORD

make integration ORCA_ORG=dockerorcadev PULL_IMAGES=1 TAG=1.1.0-pr1234
```

## More details

Testify suite's can't run in parallel within one package, so each
test suite should be in a distinct directory.  If you're using a local
machine driver, and your dev box has limited resources, you may have
resource exhaustion.  You may consider using the "-parallel N" flag to
go test to limit the number of VMs it will create in parallel.

You can tune the create flags for docker-machine, but be careful not to
set the disk too small or it'll run out of space to fit all the images.

Hint: Most local machine drivers (kvm, vbox, etc.) wont work within a
container, so if you use one of those, you'll have to run the tests from
the dev box.

While developing tests, you might want to use the "-v" flag to go test
to see all the log output immediately instead of only after a failure.
You can also set `PRESERVE_TEST_MACHINE` in your environment to a
non-empty value to prevent the virtual machines from being deleted after
the test run so you can inspect the aftermath.


The tests can hammer on your disks pretty hard.  If you're testing
locally, you might see slightly less impact on your system by running
the tests under ionice, although the brunt of the I/O load is still in
the VMs.

## AWS Examples

Hint: If you're using your own AWS setup, make sure the Docker Machine
security group allows port 443 from anywhere, otherwise the tests will
fail to access UCP.  Also make sure that all ports are allowed from the
internal network for join to work.  Our common service account already
has this set up.

### Environment Variables

The easiest way to set up AWS-specific environment variables is to use the
matrix helper, which is used to provide human-readable mappings for common
configurations.

```bash
### Set the AWS environment variables from lastpass here...

# engine=oss1.12 needs jenkins credentials, set them here if you're using it
read -s JENKINS_PASS
export JENKINS_USER=myusername JENKINS_PASS

eval $(script/matrix_helper platform=aws-ubuntu15.10 engine=oss1.12)
```

The `matrix_helper` will throw an error if you used an item that depends on an
existing environment variable that is unset (currently only for jenkins
credentials for `engine=oss1.12`).

Look at `script/matrix_params.json` for the mappings. Specify either zero or one
for each type (engine, platform, swarm).

### Running the tests

Use one of the following

* Without local volume mounts, requires full build each iteration
    * `./script/run make TEST_TIMEOUT=30m TEST_PARALLEL=2 all image integration`
* With local volume mounts (faster incremental startup)
    * `./script/run_inc make TEST_TIMEOUT=30m TEST_PARALLEL=2 all image integration`
* Run from the local system (not a container)
    * `make TEST_TIMEOUT=30m TEST_PARALLEL=2 all image integration`


Depending on the size of your system (or VM), you might experiment
with different parallel values to see what works the fastest without
causing failures.  By default the parallel count is the number of CPUs.

### docker-machine race condition

There's a race condition in docker-machine which may come into play with some
of the tests that start multiple nodes concurrently. If machine has never
created a node before, it will need to generate some certs, however there is no
locking in machine, so there's a race condition that the various cert files may
not match.

To resolve this, create a node with docker-machine prior to running the tests
and ensure that `$HOME/.docker/machine/certs` exists.

When running within a container, using `script/run` or `script/run_inc`, these
scripts attempt to work around this by keeping a persistent volume
(`ucp-controller-dev-dockerdir` by default) with the `.docker` directory. As
long as you run a non-parallel test initially, then the certs will be properly
created.

This race condition causes this type of error message:

```
INFO[0212] Orca Images and friends: [dockerorcadev/ucp:latest busybox:latest docker/ucp-cfssl:1.0.0 docker/ucp-dsinfo:1.0.0 docker/ucp-controller:1.0.0 docker/ucp-swarm:1.0.0 docker/ucp-etcd:1.0.0 docker/ucp-proxy:1.0.0] 
INFO[0212] Transfering dockerorcadev/ucp:latest to 52.37.133.85:2376 
--- FAIL: TestAcceptance (212.39s)
        Error Trace:    suite.go:33
                        suite.go:68
                        suite.go:62
                        acceptance_test.go:352
        Error:          Expected nil, but got: &errors.errorString{s:"Failed to save image to destination: dockerorcadev/ucp:latest: Post https://52.37.133.85:2376/v1.15/images/load: remote error: bad certificate"}
```

# Environment Variables

**Variable** | **Description**
-------------|-----------------
`INTEGRATION_TEST_SCOPE` | What scope of tests to run
`MACHINE_DRIVER` | Passed to docker-machine during create as `--driver`
`MACHINE_CREATE_FLAGS` | Additional flags passed during the create operation
`TEST_TIMEOUT` | Used by Makefile to set the go test `-timeout` value
`TEST_PARALLEL` | Used by Makefile to set the go test `-parallel` value (default is # of CPUs)
`PRESERVE_TEST_MACHINE` | Doesn't remove the test machine after the run.  Can be useful for troubleshooting test failures during test development
`MACHINE_FIXUP_COMMAND` | Command to run on test machine after `docker-machine` creates it (optional), try something like `curl https://test.docker.com | sudo sh`
`MACHINE_LOCAL` | Set to a non-empty string to make the tests run against the "local" docker engine (based on `DOCKER_HOST` and friends) This will disable all dual-node tests, and may lead to cascading failures.  You should make sure UCP isn't already installed on the local engine
`TAG` | Specify which image tag to use for the UCP images
`ORCA_ORG` | If set, use this org instead of the default 'dockerorca' (typically set to 'dockerorcadev' to test CI builds
`PULL_IMAGES` | If set to non-empty, pull the images (if not present).  If not set, copy them from the local system
`REGISTRY_USERNAME` | When pulling set this to a user that has permission to pull the official images
`REGISTRY_PASSWORD` | When pulling set this to a user that has permission to pull the official images
`REGISTRY_EMAIL` | When pulling set this to a user that has permission to pull the official images
`STRESS_OBJECT_COUNT` | If set, tunes the number of objects created during stress tests
`SWARM_IMAGE` | If set, pull/transfer the given image and tag it so that will be used by the integration tests
`STRESS_OBJECT_COUNT` | If set, tunes the number of objects created during stress tests
`SWARM_IMAGE` | If set, pull/transfer the given image and tag it so that will be used by the integration tests
`DTR_URL` | The FQDN or URL of a Docker Trusted Registry Test Rig
`DTR_USERNAME` | The username to be used for UCP-DTR Integration Tests
`DTR_PASSWORD` | The password to be used for UCP-DTR Integration Tests
`DTR_PRIVATE_REPO` | A pre-populated private repository in the DTR Test Rig, such as `admin/hello`
