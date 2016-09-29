### Preparing a Mac release

This document describes the Mac-specific part of the release process. See
[RELEASE.md](RELEASE.md) for the overall process.

### Update metadata

- Run `make versions` to verify that all version strings are
  correct. Manual check:

    - Verify that the version up to date, of the form
      `<docker engine version>-<arbitrary string>-dev` (see
      `README.md > Versioning `). (like "1.9.1-beta2-dev"). The version is
      stored in `CFBundleShortVersionString` in
      `v1/mac/src/docker-app/docker/docker/Info.plist` and the xcode
      project can be opened with `make dev`.

    - Verify that the right version is in
      [docker-diagnose](../v1/docker-diagnose/src/diagnose.ml)

    - Verify that the right version is in the [CHANGELOG](../CHANGELOG)

- Update the date and add final entries in the CHANGELOG.


### Run tests

Verify that the tests pass:

```
make test
```

If brave enough, verify that the long-running tests also work:

```
make fulltest
```

Additional manual tests (optional):

1. Verify that both auto-update and download from dmg works by
  - uninstalling current version (use button in settings)
  - download previous master build from hockeyapp and install from dmg
  - auto update to latest master

2. Verify that Docker.app is still running after a reboot.

3. Build and run some test containers to see that everything works as expected.

See also [MAC-TESTING.md](MAC-TESTING.md)

### Tag the release

See [RELEASE.md](RELEASE.md) for details.

### Verify and test build

See [MAC-TESTING.md](MAC-TESTING.md) for pre-release test plan.
