### Mac test plan

- [ ] Basic installation: Uninstall any running Docker.app's and then install the build

Run release tests on the build, e.g. by replacing BUILD with the correct version in the following bash script:

```
# make sure v1/mac/build is empty first, otherwise that build will be tested
export BUILD="1.12.xx.xx"
export URL="https://download-stage.docker.com/mac/rc/$BUILD/Docker.dmg"
./rt-local -l release,installer -vv -e "D4X_INSTALLER_URL=${URL}" -x run
```

VMs can now be created easily with `./rt-vm` from `tests/` (requires access to Cambridge office network).

Pinata-CI builds for tags will appear here: https://rod.cam.docker.com:8443

- [ ] OS X 10.10 release tests
- [ ] OS X 10.11 release tests
- [ ] OS X 10.12 release tests

- [ ] Auto update: Uninstall build, then install previous test build and auto update to the current one. If a test build is not available, pick a master build close to previous release and auto update to current release.
- [ ] Auto update on 10.12: Install previous master build on OSX 10.12 and upgrade to latest version to verify that auto updater works
- [ ] Check that events appear in Mixpanel

- [ ] Verify that experimental features are enabled/disabled (beta=enabled, stable=disabled)

UI tests (optional, but should be tested if there is a relevant change)

- [ ] Test that migration from Toolbox works as expected
- [ ] Test that filesharing UI works
- [ ] Test that insecure registries and registry mirrors work

On release binary:
- [ ] Verify that Info.plist auto updates from the correct appcast endpoint.
- [ ] Run `./rt-local -l release -v -x run` a last time to check that tests are still green
- [ ] Verify that `docker-diagnose` works on release build (after #4869)
- [ ] Test moving between channels, see below

#### Channel tests

If releasing beta:

1. install current stable then attempt beta installation over top
2. install current stable over top of testing beta

Both should result in warnings/errors and a functional installation of the new install afterward.

If releasing stable:

1. install current beta then attempt stable installation over top
2. install current beta over top of testing stable

Both should result in warnings/errors and a functional installation of the new install afterward.
