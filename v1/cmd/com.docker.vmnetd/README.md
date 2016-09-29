### com.docker.vmnetd

FIXME: update that README

To build the com.docker.vmnetd, a `com.apple.framework.vmnet` proxy service,
first locate a suitable kernel and
ramdisk (e.g. by opening a boot2docker.iso). Put these files in

- resources/vmlinuz64
- resources/initrd.img

Next build the package:

```
make build/final.pkg
```

To give it a whirl, install it as normal e.g.

```
open build/final.pkg
```

Next, run your VM:

```
/opt/docker/bin/docker-vmd.linux
```

Testing the proxy with and without signing
==========================================

To build the proxy:

```
$ cat > Makefile.inc <<EOT
SIGNING_ID="my apple ID"
EOT

$ make
```

This should build

- `build/proxy`: a signed binary
- `build/proxy.unsigned`: an unsigned version

Testing
-------

The `--test` argument attempts to open vmnet.

The signed binary, as a user:
```
$ ./build/proxy --test
Attempting to initialise vmnet:
Failed to initialise vmnet: take a look in Console.app for
com.apple.framework.vmnet
```

The signed binary, with sudo:
```
$ sudo ./build/proxy --test
Password:
$ echo $?
132
```
-- it crashed.

The unsigned binary, as a user:

```
$ ./build/proxy.unsigned --test
Attempting to initialise vmnet:
Failed to initialise vmnet: take a look in Console.app for
com.apple.framework.vmnet
```

The unsigned binary, with sudo:

```
$ sudo ./build/proxy.unsigned --test
Attempting to initialise vmnet:
Successfully opened vmnet.
```
