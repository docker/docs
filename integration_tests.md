# Orca Integration Tests

The orca integration tests are built upon Docker Machine and the Go test
library github.com/stretchr/testify.

The integration tests live in controller/integration/ and wont run by
default with go test, unless you set env vars to enable them.

A full list of integration test environment variables is listed at the
end of this document.

There are currently two modes of running the tests.

## Fresh Machines for each Suite

This is the preferred mode, although it is slower.  Each suite gets
a fresh machine so there's no leakage between runs.  It also allows
for multi-node testing.

When running in this mode, you can use parallel execution of tests to
speed things up.  Each suite will run in parallel, and use a unique
machine.

```bash
export MACHINE_CREATE_FLAGS="--kvm-memory 2048 --kvm-disk-size 5000"
export MACHINE_DRIVER=kvm
PARALLEL_COUNT=3

make TEST_TIMEOUT=30m TEST_PARALLEL=2 integration
```

If you're developing your own tests, you can run them explicitly using
go test directly with something like the following

```bash
go test -v -p 1 ./controller/integration/MYTEST/...
```


## Use a Single Local Machine

This mode is faster, but failures can cascade and the tests can't
run in parallel.  It is **highly** advised that you set up a machine
using docker-machine, wire up your environment pointing to it, and
**then** run the tests against that machine.  To run this mode use
`MACHINE_LOCAL=1 make integration` or if you prefer to have more control
over running the tests, `export MACHINE_LOCAL=1 && go test -v -p 1
./controller/integration/...`


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
fail to access Orca.  Also make sure that all ports are allowed from the
internal network for join to work.  Our common service account already
has this set up.

Normal reserved instances take ~3 minutes to spin up.  You may be able to
run the tests within the default 10 minute timeout, but it's pushing it.
If you use spot instances, they take ~6 minutes to spin up, so you
definitely need more than 10 minutes on the go test timeout.


### Spot Environment Variables

```bash

### Set the AWS environment variables from lastpass here...

export MACHINE_CREATE_FLAGS="--amazonec2-request-spot-instance --amazonec2-spot-price 0.05 --amazonec2-instance-type m1.small --amazonec2-ami ami-e91605d9 --engine-install-url https://test.docker.com"
export MACHINE_DRIVER=amazonec2
```

### On-demand Environment Variables

```bash

### Set the AWS environment variables from lastpass here...

export MACHINE_CREATE_FLAGS="--amazonec2-instance-type m1.small --amazonec2-ami ami-e91605d9 --engine-install-url https://test.docker.com"
export MACHINE_DRIVER=amazonec2
```

### Running the tests

Use one of the following

* Without local volume mounts, requires full build each iteration
    * `./script/run make TEST_TIMEOUT=30m TEST_PARALLEL=2 all image integration`
* With local volume mounts (faster incremental startup)
    * `./script/run_inc make TEST_TIMEOUT=30m TEST_PARALLEL=2 all image integration`
* Run from the local system (not a container)
    * `make TEST_TIMEOUT=30m TEST_PARALLEL=2 all image integration`
* Run the raw tests (useful when you want to limit scope to specific tests - can be combined with the run scripts)
    * `go test -v -timeout 30m -p 1 ./controller/integration/...`


Depending on the size of your system (or VM), you might experiment
with different parallel values to see what works the fastest without
causing failures.  By default the parallel count is the number of CPUs.

# Environment Variables

**Variable** | **Description**
-------------|-----------------
`MACHINE_DRIVER` | Passed to docker-machine during create as `--driver`
`MACHINE_CREATE_FLAGS` | Additional flags passed during the create operation
`TEST_TIMEOUT` | Used by Makefile to set the go test `-timeout` value
`TEST_PARALLEL` | Used by Makefile to set the go test `-parallel` value (default is # of CPUs)
`PRESERVE_TEST_MACHINE` | Doesn't remove the test machine after the run.  Can be useful for troubleshooting test failures during test development
`MACHINE_LOCAL` | Set to a non-empty string to make the tests run against the "local" docker engine (based on `DOCKER_HOST` and friends) This will disable all dual-node tests, and may lead to cascading failures.  You should make sure Orca isn't already installed on the local engine
`TAG` | Specify which image tag to use for the orca images
`PULL_ALL_IMAGES` | If set to non-zero, pull the images (if not present).  If not set, copy them from the local system
`REGISTRY_USERNAME` | When pulling set this to a user that has permission to pull the official images
`REGISTRY_PASSWORD` | When pulling set this to a user that has permission to pull the official images
`REGISTRY_EMAIL` | When pulling set this to a user that has permission to pull the official images

