# Regression tests

This directory contains both a generic, cross-platform regression test
framework and regression tests specific to the Docker for Mac and
Windows products.

The generic regression test framework is implemented in Python with
the majority of the code contained in the `./lib` sub-directory. This
directory should *not* contain any Docker for Mac and Windows specific
code.

The Docker for Mac and Windows specific code and test cases are
implemented in the `./cases` sub-directory.


# Prerequisites

On a Mac, this should pretty much just work out of the box.  However,
for development you might want to consider installing `flake8`, a
python style checker.
```
brew install python
pip install flake8
```

On Windows, you need to have python installed as well as some form of
`bash`. The tests must be run from a `bash` shell and `bash.exe` must
be in your path. The simplest way is to install via
[chocolatey](https://chocolatey.org/):

```
choco install git python
```

See Section on Machine setup below for more details.


# Quickstart

The regression test framework allows running tests on a local host (or
inside a VM) as well as against a suitably configured remote host.

To run tests locally, simply execute the `rt-local run` command.  It
will executed all the test cases in the `./cases` sub-directory.

To list all current tests run `rt-local list`, or to get a one line
summary for each test use `rt-local info`.

When running tests, by default a line per test is printed on the
console with a pass/fail indication.  Detailed logs, by default, are
stored in `./_results/<UUID>/`. In that directory, `TESTS.log`
contains detailed logs of all tests, `TESTS.csv` contains a line per
test and `SUMMARY.csv` contains a one line summary of the all tests
run. The directory also contains a log file for each tests, with the
same contents as `TESTS.log`.

If you prefer a bit more information in the log files use:
```
rt-local -x run
```
This executes the tests with `-x` and thus logs all commands executed.

For a CI system, where the output is displayed on a web page use:
```
rt-local -x -vvv run
```
This prints the same information logged to the log file to the console.


To run test on a remote host use the `rt-host`. You must specify the type of the remote host (`osx` or `win`) as well as the address. Optionally, you can supply user name and password.
```
./rt-host -t win -r 172.16.10.131 -- <rt-local arguments>
```
This:
- copies the test framework and tests from the local host to the remote host.
- executes the test on remote host
- copies the results back to the local host


# User's guide

Tests are located in the top-level project directory (default
`./cases`). Within a project tests may be grouped by placing them in a
common sub-directory (forming a _Test Group_).  The order in which
tests and test groups are executed is determined by alphabetical order
and one may prefix a directory with a number to force a order. Each
test is given a unique name based on the place in the directory
hierarchy where it is located.  For the naming of tests, any
prefix-numbers, used to control the order of execution, are
stripped. For example, a test located in:
```
cases/backend/compose/010_django_example/
```
will be named:
```
backend.compose.django_example
```

While traversing the directory tree, directories starting with `_` are
ignored.


## Run control

Apart from the simple command line examples given in the Quickstart
section above, the regression test framework provides fine grained
control over which tests are run.  The primary tool for this are
_labels_.

A test may define that a specific label "foobar" must be present for
it to be run or, by prefixing the label name with a '!', that the test
should not be run if a label is defined. The labels are defined with
`test.sh` (see below).

The regression test framework is completely agnostic to which labels
are used for a given set of test cases, though it defines a number of
labels based on the OS it is executing on. See the output of `rt-local
list` on your host.

A good strategy for using labels is to not define any labels for tests
which should always be executed (e.g. on every CI run).  Then use
labels to control execution for tests which are specific for a given
platform (e.g. `osx` or `win` for OS X and Windows installer tests).
If a particualr test is known not to run on a given platform you can
use, e.g. `!win`, to indicate it.  Finally, define separate labels,
e.g., long running tests or extensive tests which should only be run
on release candidates.

You can control labels for executing tests by using the `-l` flag.
For example, if you have some long running tests which you do not wish
to execute on every CI run and have them marked with a label `long`,
then you can execute them with:
```
./rt-local -l long run
```

You can see which tests would get executed using the `-l` flag for the
`list` command as well:
```
./rt-local -l long list
```

In addition to control which tests are run via labels it is also
possible to specify an individual test or a group name on the command
line to just run these test (subject to labels).  Here are two
examples:
```
./rt-local run pinata.backend.volumes.hardlink
./rt-local run pinata.backend.volumes
```
The first runs a single test, while the second is running all tests
within the volumes group. Note, that this is currently implemented as
a simple prefix match, so, if you have tests such as `foo.bar` and
`foo.bar_baz` and use `./rt-local run foo.bar`, it will execute both
`foo.bar` and `foo.bar_baz`.


## Writing tests

Tests are simple scripts which return `0` on success and a non-zero
code on failure.  A special return code (`253` or `RT_TEST_CANCEL`)
can be use to indicate that the test was cancelled (for whatever
reason).  Each test must be located in its own sub-directory (together
with any files it may require).

Currently, test is a simple shell script called `test.sh`. In the
future we will also support Python and Powershell (for Windows only
tests).

There is a template `test.sh` in `./etc/templates/test.sh` which can be
used for writing tests. It contains a number of special comments
(`SUMMARY`, `LABELS`, `REPEAT`, and, `AUTHOR`) which are used by the
regression test framework. The `SUMMARY` line should contain a *short*
summary of what the test does. The `LABEL` is a (optional) list of
labels to control when a test should be executed.  `AUTHOR` should
contain name and email address of the test author.  If a test is from
multiple authors, multiple `AUTHOR` lines can be used.  Finally, the
`REPEAT` line may contain a single number to indicate that a test
should be executed multiple times.  The `REPEAT` line may also contain
`<label>:<number>` entries to runtest multiple times if a label is
present.

A few guidelines for writing tests:

- A test should always clean up whatever is created during test
  execution. This includes containers, docker images, and files. The
  template contains a `clean_up()` function which can be used for this
  purpose.

- An individual tests should not rely on a artefact left behind by a
  previous test (even if one can control the order in which tests are
  executed)

- Tests should be self contained, i.e. they should not rely on any
  files outside the directory they are located in.

The regression test framework currently passes the following
environment variables into a test script:

- `RT_ROOT`: Points to the root of the test framework.  This can be
  used, e.g. to source common shell functions from
  `${RT_ROOT}/lib/lib.sh`

- `RT_PROJECT_ROOT`: Points to the root of the project. This can be
  used, e.g. to source common shell functions defined by a project.

- `RT_OS`, `RT_OS_VER`: OS and OS version information. `RT_OS` is one
  of `osx` or `win`.

- `RT_RESULTS`: Points to the directory where results data is stored.
  Can be used in conjunction with `RT_TEST_NAME` to store additional
  data, e.g., benchmark results.

- `RT_TEST_NAME`: Name of the test being run. A test implementation
  can use this in conjunction with `RT_RESULTS` to store additional
  data, e.g., benchmark results.

- `RT_LABELS`: A colon separated list of labels defined when the test
  is run. Can be used by the test to trigger different behaviour based
  on a Labels presence, e.g., run longer for release
  tests. `./lib/lib.sh` provides a shell function, `rt_label_set`, to
  check if a label is set.

Users can specify additional environment variables using the `-e` or
`--env` command line option to `rt-local`.  This may be useful for
scenarios where `rt-local` is executed remotely.


### General utilities for writing tests

Under `lib/utils` are a number of small, cross-platform utilities
which may be useful when writing tests (the environment variable `RT_UTILS` points to the directory):

- `rt-filegen`: A small standalone utility to create a file of fixed
  size with random content.

- `rt-filerandgen`: A variation of the above. The filesize is with a
  set maximum.

- `rt-filemd5`: Returns the MD5 checksum of a file.

- `rt-crexec`: A utility to execute multiple commands concurrently in a
  random order.

- `rt-urltest`: A utility like curl, with retry and optional `grep` for a
  keyword on the result.


## Creating a new Test Group

Any directory containing sub-directories with tests under the
top-level project directory forms a test group.  The execution of a
group can be customised by defining a `group.sh` file inside the group
directory.  Like with tests, the `group.sh` may provide a `SUMMARY`
and a set of `LABELS`.

If a group contains a `group.sh` file, it is executed with the `init`
argument before the first test if the groups is executed, and with the
`deinit` argument after the last test of the group was executed.  Test
writers can thus place group specific initialisation code into
`group.sh`. `group.sh` is executed with the same environment variables
set as `test.sh` scripts.

There is a template `group.sh` file in `./etc/templates/group.sh`.

The top-level `group.sh` should also create a `VERSION.txt` file in
`RT_RESULTS`, containing some form of version information.  If the
tests are run against a local build, this could be the git sha value,
or when run as part of CI the version of the build being tested etc.


# Machine setup
WIP

## Windows 10 machines

When installing Window 10 make sure you:

- Select custom settings and disable all the options (basically
  disable windows calling home).
- Don't use a Microsoft account but create a local one with username
  `docker` and password `Pass30rd!`
- Enable Remote Desktop (I use the UI for that)
- Install VMware tools if this is a VM

The remaining setup is done in an elevated powershell. You can basically cut and paste this. The steps performed are:
- Install `choco` (kinda `apt-get` for windows)
- Install git, bash, and python (and python2)
- Disable powersaving modes (no remote access when sleeping)
- Enforce that the network is on a "Private" network
- Enable Windows Remote Management
- Disable the UAC (prompt when running commands elevated)
- Enable Hyper-V (this will reboot your system)

```
Set-ExecutionPolicy -ExecutionPolicy Unrestricted -Force
iwr https://chocolatey.org/install.ps1 -UseBasicParsing | iex
choco install -y git python python2

powercfg -hibernate off
powercfg -change -monitor-timeout-ac 0
powercfg -change -standby-timeout-ac 0
powercfg -change -hibernate-timeout-ac 0

$networkListManager = [Activator]::CreateInstance([Type]::GetTypeFromCLSID([Guid]"{DCB00C01-570F-4A9B-8D69-199FDBA5723B}"))
$connections = $networkListManager.GetNetworkConnections()

# Set network location to Private for all networks
$connections | % {$_.GetNetwork().SetCategory(1)}

winrm quickconfig -force
winrm set winrm/config/client/auth '@{Basic="true"}'
winrm set winrm/config/service/auth '@{Basic="true"}'
winrm set winrm/config/client '@{AllowUnencrypted="true"}'
winrm set winrm/config/service '@{AllowUnencrypted="true"}'

New-ItemProperty -Path HKLM:Software\Microsoft\Windows\CurrentVersion\policies\system -Name EnableLUA -PropertyType DWord -Value 0 -Force

Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All
```
If you create a VM image, take a snapshot now.

**Make sure the first `python` in your path is python 3.x**.  There
seem to be some subtle differences in the way the logging over a
socket work between Python2 and Python3 and with Python2 the socket
may not get closed properly. Using
[Chocolatey](https://chocolatey.org), Python2 is currently installed
in `C:\tools\python2` while Python3 is placed in
`C:\ProgramData\Chocolatey\bin\`.

# Design/Internals

The regression test framework is written in Python with the bulk
implemented by the `rt` module, located in `./lib/python/rt`.  The
main entry point for running local tests is implemented in
`./lib/python/rt/local.py`.

The core of the framework is implemented in the `./lib/python/base.py`
file. This file provided two classes, `Test` and `Group` which handle
the traversal, naming and execution of tests.


## Results

Apart from the log files (see below) and the output on the console,
the regression test framework also produces two CSV files. `TESTS.csv`
contains a line per test and `SUMMARY.csv` contains a one line summary
of all tests run.

The CSV files are linked by a UUID, which can be thought of as a key.
The structure should make it very easy to store the test results in a
simple database containing two tables, one with test summaries and
another with all test results.


## Logging

For logging we utilise the standard Python logging package.  We use
three logging backends:

- A file based one which writes pretty much any output to a log file.
- A file based one per tests which logs the same as the above but on a
  per test basis.
- A console based logger for printing progress on the terminal.

The file based loggers are pretty standard, while the console based
logger is customised to provided colourful output.

We define several additional log levels:

- Two log-levels to capture stdout/stderr from external commands which
  are executed.  Using separate log-levels allows us to annotate the
  output accordingly.

- Separate log levels for test results. The summary of passed, failed
  and skipped tests are logged to different log levels. This allows
  finer control of what and how the results are displayed on the
  console.

The "console" logger can optionally be redirected over a socket to a TCP server. `lib/python/rt/log.py` contains a sample server.

The logging customisation is implemented in `lib/python/rt/log.py`


## Remote execution on Windows

To implement `rt-host` to control windows machines, we utilise Window
Remote Management (WinRM). A python module
[pywinrm](https://github.com/diyan/pywinrm) is used for most of the
heavily lifting. Most of the `rt-host` code is implemented in
`lib/python/rt/host.py`. [This page](http://www.hurryupandwait.io/blog/understanding-and-troubleshooting-winrm-connection-and-authentication-a-thrill-seekers-guide-to-adventure)
has a lot of nice details on configuring WinRM on the Window machine.

WinRM has three issues:
- It only allows execution of commands, but offers no file transfer.
- One only gets the standard output and error streams after the
  command execution finished.
- At least the python binding to not return when a long running remote
  command exits.  I don't know how long long-running is, but basically
  a test run of something like 20mins length will not return.

To address the first issue, we basically start a web server on demand
on the local host which under stands `GET` and `POST` requests.  To
transfer single files these are simply transferred by invoking `curl`
on the remote host.  If we copy directories we `tar` the directory up
before transferring and then extract them.

The second issue is addressed by logging what would normally be logged to stdout to a TCP server using the network logger briefly described above.

The final issue is addressed by starting the `rt-local` command in the
background if we run tests and wait till the log server exits.

# Virtual Machine Management

In order to create Windows and OSX VM's for testing, we have created
`rt-vm`. This tool will interact with a VMware Server (ESXi or
vCenter) or with a local VMware Workstation or Fusion Pro instance and
is capable of deploying the virtual machine templates in the
`./etc/vm-templates` directory.

These templates have a strict naming convention that is:

```
test-<os>-<version>-template
```

VM's that are created from these templates are named:

```
test-<os>-<version>-<uuid>
```

You need to manually create templates first. See
`etc/vm-templates/Makefile` for available targets.  Creating templates
may take quite a while...

## Configuration

To configure this tool, you must first run `./rt-vm init`. This will create a JSON configuration file called `vm.json`. You must configure which driver you want to use, configuration details for that driver, and information about the OS's and OS Versions available.

## VM lifecycle management

Virtual Machine creation, using `./rt-vm create <os> <version>` creates a "Clone" of a the VM template and starts it. We clone a template so we can be sure that we have fresh state each time we start our tests.

Once created, you can use the `./rt-vm start|stop|delete` commands to manage it's lifecycle. If you wish to create snapshots, you may do so from within VMWare itself. To revert to a snapshot, you can use `./rt-vm revert <name> <snapshot>` command.

## TODO

- Add `put|fetch` commands for adding files too and retrieving files from the VM
- Add `run` command for remote command execution
- Implement a `hyper-v` driver to support local testing on Windows
