# Regression tests for Docker for Mac and Windows

The way the tests are structured they should support a number of
different use cases:

- Developer testing: Run tests against a local build before checking in.
- Continuous integration (CI): Run some tests on every PR
- Continuous integration (CI): Run more tests ideally on every master
  build, or least as frequent as possible.
- Release engineering: Run a lot more tests longer

To support these different use-case we need to control:
- which tests are run. This is achieved via use of
  labels (see below).
- against what build-artefacts tests are run.

## Known issues and current limitations

- On a Mac, currently you can only run tests either against your local
  build or a local install. We don't yet support downloading a
  installer and test the install procedure.  Running against your
  local build is the default. When running tests, it will stop teh
  Docker app currently running, start the App from your build
  directory and use the binaries from your build directory. If you
  specify the `nostart` label, tests will be run against whatever
  version of Docker is currently running.

## Labels

We define a number of additional labels to control test execution:

- `master`: Enables some tests which run longer and may not be run
  during normal CI but we should run as often as possible on master
  builds.

- `release`: Enables even more tests, should be used on release
  candidates.

- `installer` also, installs from installer.

- `nostart`: Run tests assuming that app is already running.

- `benchmarks`: This label runs benchmarks and most of the other tests
  are disabled.  The `release` label also runs the benchmarks.


## Common command line invocations

Running tests is controlled primarily by labels, but also, to a lesser
extend by environment variables. This section provides some common use
cases and some sample command line invocations.

### OS X: Developer running tests against a local build

This can be done simply by calling `make run`.

Alternatively, if you already have docker running, use:
```
./rt-local -vx -l nostart run
```

### Windows: Developer running against a local build

You have to have build a installer for this, i.e., `.\please.ps
package`. Then you can run the tests with:
```
 ./rt-local -vx -l installer,master run
```
Note, that this will stop and un-install any Docker for Windows app
you have and install your local build.

You may also want to set `D4X_PASSWORD` to your password, so that
volume sharing tests work.


### Windows CI

To test a Installer build by CI you can specify it via the
`D4X_INSTALLER_URL` environment variable. This will download the
latest master build, install it and run tests on it:

```
./rt-local -e D4X_INSTALLER_URL=https://download-stage.docker.com/win/master/InstallDocker.msi -vx -l installer,master run
```

Again, you may want to set `D4X_PASSWORD` and you may also want to
change `master` to `release` to run more extensive tests.


## Writing tests

Also see the more general section on writing tests in the top-level RT
directory.

Tests should ideally be written to work for Docker for Mac and
Windows. If a test is specific to either platform, please use the
`osx` or `win` labels provided by the RT framework.

A few more guideline for writing tests:

- All `test.sh` must source `${RT_PROJECT_ROOT}/_lib/lib.sh` several
  variables are set that way.

- All tests sharing temporary files between host and a container must
  use the `D4X_LOCAL_TMPDIR` environment variable to refer to the host
  side location.  `/tmp` is very different on Windows.

The library in in `${RT_PROJECT_ROOT}/_lib/lib.sh` provides a set of
common shell functions which should be usable for both Docker for Mac
and Docker for Windows.  These are:

- `d4x_app_installed()`: Returns True if the Application is installed.

- `d4x_app_running()`: Returns True if the application/backend is running.

- `d4x_app_start()`: Start the application.

- `d4x_backend_start()`: If the application allows on to start a
  backend only, ie without UI, use this function.

- `d4x_app_stop()`: Stop the application/backend.
