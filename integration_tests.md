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

cd controller
go test -p ${PARALLEL_COUNT} ./integration/...
```


While developing tests, you might want to use the "-v" flag to go test
to see all the log output immediately instead of only after a failure.
You can also set PRESERVE\_TEST\_MACHINE in your environment to a
non-empty value to prevent the virtual machines from being deleted after
the test run so you can inspect the aftermath.


The tests can hammer on your disks pretty hard.  If you're testing
locally, you might see slightly less impact on your system by running
the tests under ionice, although the brunt of the I/O load is still in
the VMs.
