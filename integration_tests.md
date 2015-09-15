# Orca Integration Tests

The orca integration tests are built upon Docker Machine and the Go test
library github.com/stretchr/testify.

The integration tests live in controller/integration/ and wont run by
default with the script/test utility.

Since it is slow to create virtual machines, the tests are structured to
attempt to run in parallel as much as possible.  Testify suite's can't
run in parallel within one package, so each test suite should be in a
distinct directory.  If you're using a local machine driver, and your
dev box has limited resources, you may have resource exhaustion.
You may consider using the "-parallel N" flag to go test to limit the
number of VMs it will create in parallel.

You can tune the create flags for docker-machine, but be careful not to
set the disk too small or it'll run out of space to fit all the images.

Here's a typical scenario for running the integration tests.  Substitute
your favorite machine driver and tuning settings:

```bash
export MACHINE_CREATE_FLAGS="--kvm-memory 2048 --kvm-disk-size 5000"
export MACHINE_DRIVER=kvm
PARALLEL_COUNT=$(grep "^processor" /proc/cpuinfo | wc -l)

make TEST_TIMEOUT=30m TEST_PARALLEL=2 integration
```

Hint: Most local machine drivers (kvm, vbox, etc.) wont work within a
container, so if you use one of those, you'll have to run the tests from
the dev box.


While developing tests, you might want to use the "-v" flag to go test
to see all the log output immediately instead of only after a failure.
You can also set PRESERVE\_TEST\_MACHINE in your environment to a
non-empty value to prevent the virtual machines from being deleted after
the test run so you can inspect the aftermath.

By default, the tests will copy the non-official images (orca itself,
the bootstrapper, etc.) and pull the official images (swarm, etc.) If you
instead want it to copy (not pull) all the images, set COPY\_ALL\_IMAGES
to a non-empty value.


The tests can hammer on your disks pretty hard.  If you're testing
locally, you might see slightly less impact on your system by running
the tests under ionice, although the brunt of the I/O load is still in
the VMs.

## AWS Example

Hint: If you're using your own AWS setup, make sure the Docker Machine
security group allows port 443 from anywhere, otherwise the tests will
fail to access Orca.  Also make sure that all ports are allowed from the
internal network for join to work.  Our common service account already
has this set up.

Normal reserved instances take ~3 minutes to spin up.  You may be able to
run the tests within the default 10 minute timeout, but it's pushing it.
If you use spot instances, they take ~6 minutes to spin up, so you
definitely need more than 10 minutes on the go test timeout.

```bash

### Set the AWS environment variables from lastpass here...

export MACHINE_CREATE_FLAGS="--amazonec2-request-spot-instance --amazonec2-spot-price 0.05 --amazonec2-instance-type m1.small --amazonec2-ami ami-e91605d9"
export MACHINE_DRIVER=amazonec2

# without local volume mounts, requires full build each iteration
./script/run make TEST_TIMEOUT=30m TEST_PARALLEL=2 all image integration

# with local volume mounts (faster incremental startup)
./script/run_inc make TEST_TIMEOUT=30m TEST_PARALLEL=2 integration
```

Depending on the size of your system (or VM), you might experiment
with different parallel values to see what works the fastest without
causing failures.  By default the parallel count is the number of CPUs.
