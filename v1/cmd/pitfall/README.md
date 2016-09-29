pitfall
=======

Pinata Integration Testing For ALL

## How to use

```
NAME:
   pitfall - Pinata Integration Tests For All!

USAGE:
   pitfall [global options] path_to_tests

VERSION:
   1.0.0

COMMANDS:
GLOBAL OPTIONS:
   --driver, -d "esx"			Driver esx|fusion [$PITFALL_DRIVER]
   --os, -o "osx"			OS under test [$PITFALL_OS]
   --os-version, -x "10.11"		OS Version under test [$PITFALL_OS_VERSION]
   --esx-host, -e "172.16.1.10"		ESXi Host [$PITFALL_ESX_HOST]
   --esx-user, -u "root"		ESXi User [$PITFALL_ESX_USER]
   --esx-pass, -p "slartibartfast"	ESXi Password [$PITFALL_ESX_PASS]
   --esx-datastore, -d "datastore2"	ESXi Datastore [$PITFALL_ESX_DATASTORE]
   --build, -b 				Build Number [$PITFALL_BUILD_ID]
   --help, -h				show help
   --version, -v			print the version
```
## What does it do?

- Returns a pre-created VM to a known clean snapshot
- Installs the supplied build number
- Runs the integration tests
- Saves the results!

## Why does it exist?

Because it is a handy tool that can be used with DataKit to provide some awesome CI!

## What does an image need to be compatible

- Username: `Docker`, Password: `containyourself`
- Automatic login for the Docker User
- Passwordless Sudo
- Screensaver disabled
- Power Saving disabled
- Remote Login (SSH) enabled
- Virtualization Passthrough Features Enabled
- At least 4096MB RAM
- password authentication must be enabled in sshd (Yosemite has only keyboard-interactive on by default)

## Building Pitfall with Docker

Run `make` from the `pitfall` directory. This will produce `pitfall-linux` and `pitfall-darwin` binaries. The container image `pitfall:build` will also be 
be available to run `pitfall` commands.

`make clean` removes the container images.

## Building Pitfall manually

To build pitfall manually, move to /v1 (so `vendor` is in the cwd) and run

```
go build ./cmd/pitfall
```

## TODO

- Implement proper Windows Support
- Finishing implementing the "CloneVM" functions so we can use Linked Clones for better parallelization
- Allow the user to supply a `patch` flag for testing OSX 10.11.3 or OSX 10.11.4 - this will require a sanely named snapshot
- Save the results of `pinata-diagnose` too
