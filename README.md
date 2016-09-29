# Pinata: an experimental standalone Docker client

|  | pr | master | rc | beta | stable |
|---|---|---|---|---|---|
| macOS | [latest](https://download-stage.docker.com/mac/pr/Docker.dmg) | [latest](https://download-stage.docker.com/mac/master/Docker.dmg) | [latest](https://download-stage.docker.com/mac/rc/Docker.dmg) | [latest](https://download.docker.com/mac/beta/Docker.dmg)  | [latest](https://download.docker.com/mac/stable/Docker.dmg)  |
| Windows | [latest](https://download-stage.docker.com/win/pr/InstallDocker.msi) | [latest](https://download-stage.docker.com/win/master/InstallDocker.msi) | [latest](https://download-stage.docker.com/win/test/InstallDocker.msi)  | [latest](https://download.docker.com/win/beta/InstallDocker.msi)  | [latest](https://download.docker.com/win/stable/InstallDocker.msi)  |
*( if you get Access Denied errors, it means nothing has been published to this channel yet )*

To list all the versions: http://omakase.omakase.e57b9b5b.svc.dockerapp.io/

This is an experimental project to develop a new client for Docker,
separately from the daemon or any other backend component.

By maintaining a standalone client, the goal is to:

1. Allow for more rapid iteration on client functionality.
2. Improve compatibility between different versions of the client and daemon.
3. Add more features to the client without bloating the daemon-side components.
4. Pave the way to simplifying the daemon code base, improving its
quality and making its maintenance easier.

## Versioning

The release cycle respects the following convention: `X-Y[-Z]` where:

- `X` is the version of the docker engine used as a base for the build. The build can be modified during the build process to fit better into the use-case of `Docker.app` (ie. it won't usually be a drop-in replacement, but we will try to upstream our patches as quickly as possible).

- `Y` is an arbitrary string that we can use to define a version of `Docker.app`, independently of the release cycle of docker engine.

- `Z` indicates the build channel (dev, test, master, release). `Z` is empty for releases.

For instance the first beta release of pinata has the version: `1.9.1-beta1`. While on master channel (one build for each PR merged), it has the version: `1.9.1-beta1-master`.

On OS X, the version is defined in XCode project's Info.plist file (key: `CFBundleShortVersionString`). There's also a build number, associated with `CFBundleVersion` key (set by CI).

In Xcode project, the version should use a suffix like `-dev` (`1.9.1-beta1-dev`). That suffix will be replaced/removed by CI.

## INSTALL

### Through HockeyApp

For Docker for Mac see the [Docker.app installation guide](https://github.com/docker/pinata/blob/master/v1/docs/content/mackit/getting-started.md) and for Docker for Windows see the [Docker installation guide](https://github.com/docker/pinata/blob/master/v1/docs/content/winkit/getting-started.md)

### OSX Build

Check that your `GOPATH` is correctly set-up and clone this repository in
`$GOPATH/src/github.com/docker/pinata`.

#### Dependencies

As prerequisites, you need to have `Xcode`, `homebrew` and `go` installed.

To minimize build times, the dependencies are cached with this command

You only need to run it once or when an external dependency was updated

At the root of this repository, type:

```
make depends
```

When you add a new go dependency, add it in the `GO_DEPS` variable of the toplevel
`Makefile`.

#### Build

After a successful `make depends`, type:

```
make
```

If you are asked for the password to the `dockerbuilder` keychain, it is
`docker4all`.

#### Run

After a successful `make depends` and `make`, type:

```
make run
```

You will see the logs on stdout

#### Install

First, make sure you have uninstalled any previous installation of
pinata with:

```
v1/mac/uninstall
```

Then, install with:

```
v1/mac/build/Docker.app/Contents/MacOS/docker-installer
```

#### Tests

You can run the tests by running:

```
make test
```
This is currently Mac only.



### Windows Build

[![Build status](https://ci.appveyor.com/api/projects/status/fpa7neeotor31bdh/branch/master?svg=true)](https://ci.appveyor.com/project/Pinata/pinata/branch/master)

Latest msi builds :
 * On [Master](https://download-stage.docker.com/win/master/InstallDocker.msi) channel.
 * On [Test](https://download-stage.docker.com/win/test/InstallDocker.msi) channel.

Check that your `GOPATH` is correctly set-up and clone this repository in
`$GOPATH/src/github.com/docker/pinata`.

#### Dependencies

Install:

- Go 1.6
- [Visual Studio 2015](https://www.visualstudio.com/en-us/products/vs-2015-product-editions.aspx).  The app builds with the free Community edition but the licensing for that edition doesn't allow its use for commercial, closed source work.

Once you installed the above, open a powershell.

#### Build

The main build is driven by the `please.ps1` powershell script in the `win`
sub-directory.

```
cd <pinata_dir>/win
./please.ps1 package
```

will clean the build directory and build a new package (installer) in
the `build` sub-directory.


```
cd <pinata_dir>/win
./please.ps1 build
```

will build a new `Docker.exe` file but not the installer.

### Troubleshooting

If you have an issue, please report it to the
[bugtracker](https://github.com/docker/pinata/issues) with the output
of:

```
pinata diagnose
```
This is currently Mac only.
