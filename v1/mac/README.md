Before you start
================

This guide assumes you are building on OSX.

Install docker on your local machine
------------------------------------

The build will need the `docker` CLI so install it somehow. For example
you could install the [Docker toolbox](http://docs.docker.com/mac/started/)
or use homebrew.

Build the app
=============

The build will require `sudo` in order to start the vmnet proxy (if one
is not already running). On my laptop `sudo` requires a password if it
hasn't been used in a while which tends to cause confusion when run from
a script. I tend to run a `sudo ls` before starting
the build, although I could run `visudo` and remove the password requirement.

Build everything:

```
sudo ls # check that sudo works
./v1/mac/make -cbsv
```

Note it might ask you questions if you haven't got homebrew setup.

Note that it does a decent job of incremental builds: the first build will be
slow but subsequent builds will be faster.

Note that part of that build actually used `hyperkit` and `boot2docker` to
recompile a Linux kernel and a `boot2docker` disk image -- how cool is that??

Run the app
===========

Remove any vestiges of previous versions you might have running:
```
./v1/mac/uninstall
```

Then look at the build output:

```
open ./v1/mac/build/OSX-Release
```

Double-click on the whale. You should get prompted for your admin password
(first time only).

All the services should be installed and starting.

Run docker in xhyve
-------------------

You can tell when xhyve has started because the `~/.xhyve/env` file
appears. It should also hack your `.bash\_profile` to automatically source
the `env` file for the benefit of future terminals. Either launch a fresh
terminal or run

```
. ~/.xhyve/env
```

Once the environment is setup, you should be able to run commands like:

```
docker pull busybox
docker run -t -i busybox sh
```

Troubleshooting
===============

SMJobBless
----------

Did you get a strange error mentioning `SMJobBless`? This happens if the
signing key baked into the property lists isn't exactly the same as the
one used to sign the binaries.

Run the verification tool:

```
$ cd v1/mac
$ ./scripts/verifySign 
com.docker.vmnetd requires docker-installer to be signed with:
  "Developer ID Application: Docker Inc (9BNSXJN65R)"
com.docker.installer is actually signed with:
  "Developer ID Application: Docker Inc (9BNSXJN65R)"
OK so far
com.docker.installer requires com.docker.vmnetd to be signed with
  "Developer ID Application: Docker Inc (9BNSXJN65R)"
com.docker.vmnetd is actually signed with:
  "Developer ID Application: Docker Inc (9BNSXJN65R)"
Everything looks ok
```

The most common way this happens is if you forget to run the `-s` ("signing")
stage of the `BuildDocker` build.

Hyperkit
--------

Is your `com.docker.hyperkit` binary not running?

This commonly happens when `hyperkit` fails to connect to one of the necessary
services (`vmnet` or `irmin9p`). Open `Console.app` and search for `hyperkit`:
this should point you in the right direction.

