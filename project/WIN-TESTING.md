This is a template for tests which have to be performed before a Windows release.

- [ ] Run regression tests with `release,installer` labels like this: `./rt-local -vxl release,installer` -e D4X_INSTALLER_URL=https://download-stage.docker.com/win/master/InstallDocker.msi run`
  - [ ] On Windows 10 Pro Build 10586 (aka 1511)
  - [ ] On Windows 10 Pro Build 14363 (aka 1607 aka aka RS1 Anniversary Release)

The quickest way to download other binaries is via:
```
Import-Module BitsTransfer
$url = "https://download-stage.docker.com/.../InstallDocker.msi"
Start-BitsTransfer -Source $url -Destination InstallDocker.msi
```

The following steps currently have to be performed manually:
**Make sure you start the app from the desktop icon, and not via the CLI!!!**
- Test auto-update.
  - Download a previous release from the test channel via [this page](http://omakase.omakase.e57b9b5b.svc.dockerapp.io/)
  - Install it
  - Wait for the update dialogue, verify the version, and update
  - [ ] Verify that the version after restart

- Verify versions
  - Verify there is a whale on the Desktop. Check if it says Beta
  - [ ] Until #4613 is closed, verify with `docker version` that Beta release use experimental and that stable release use non-experimental.
  - [ ] Verify the version in the `About Box`
  - [ ] Verify `docker-machine version`
  - [ ] Verify `docker-compose version`
  - [ ] Verify that the path was added to a Powershell, cmd, and bash window

- Test settings UI (has to be executed in order):
  - [ ] Shared Drive: Make sure you have a D drive. Share it, unshare it and share it. Verify with `docker run --rm -v d:/:/data alpine ls /data` that the sharing/unsharing worked. **Leave the drives shared**.
  - [ ] Advanced: change the CPU and Memory settings a couple of times. Verify with `Get-VMProcessor MobyLinuxVM` and `Get-VMMemory MobyLinuxVM` that the VM has the new settings. **Leave it at a non-default value**.
  - Network: Change the Subnet address to `10.0.76.0` and the mask to `255.255.255.248` and the DNS to static `8.8.4.4`. Verify:
    - [ ] `ipconfig` that the ip address/netmask has changed
    - [ ] `docker run --rm -v c:/:/data alpine ls /data` verifies that mounting still works
    - [ ] `docker run --rm alpine wget http://www.google.com` verifies that DNS and networking still works
  - **Need a test for proxy settings**
  - Docker Daemon:
    - [ ] Verify the URL points to the damon config file section 
    - [ ] Change the debug setting to `true`. Verify using `docker info` (look for `Debug Mode (server): true`)
  - Diagnose and Feedback
    - [ ] Verify the documentation link points to the documentation
    - [ ] Verify the issue link to Github works
    - [ ] Verify the `Logs` link works
    - [ ] Click on `Upload diagnostics` and verify that the diagnostics can be retrieved using the `Dockerfile.fetch`.
  - Exit the App and Start again. Verify that the settings are still what they were before the restart and verify that they are actually used:
    - [ ] `docker images` should be empty
    - [ ] `ipconfig` that the ip address/netmask
    - [ ] `docker run --rm -v c:/:/data alpine ls /data` verifies that mounting still works
    - [ ] `docker run --rm alpine wget http://www.google.com` verifies that DNS and networking still 
    - [ ] `docker info` for debug settings
  - Reset: Reset to factory and check, e.g.:
    - [ ] Shared Drive: Make sure they are not ticked
    - [ ] Shared Drive: Verify that `docker run --rm -v c:/:/data alpine ls 
/data` is empty. Make sure it pulls in alpine.
    - [ ] Advanced: Check that CPUs is reset and memory set to 2GB
    - [ ] Advanced: Verify with `Get-VMProcessor MobyLinuxVM` and `Get-VMMemory MobyLinuxVM`
    - [ ] Network: Make sure the IP address is `10.0.75.0` and mask is `255.255.255.0` and DNS is unticked.
    - [ ] Docker Daemon: Make sure `Debug` is set to `false`
    - [ ] Docker Daemon: Check `docker info` for debug setting


- Test with a user using a Microsoft Live ID for login
  - [ ] Verify sharing/unsharing via the UI works

- Test with a user with a space in the User name
  - [ ] Verify sharing and/unsharing via the UI works

- Mixpanel: Requires logging in to mixpanel:
  - [ ] verify that mixpanel events in the final release binary report the correct version string (due to #4925)
