vm-templates
============

## Acknowledgements

Adapted from https://github.com/timsutton/osx-vm-templates
and https://github.com/StefanScherer/packer-windows

## Building


All templates:

`make`

Just one:

`make machines/test-osx-10.10-template`

## Advanced Usage

### ISOs

ISO can be obtained from Amazon S3 by issuing a `make get-isos`

For Windows, you can obtain ISO's using an MSDN subscription

For Mac, you must download the appropriate installer from the App Store and then run:

```
sudo ./prepare_iso/prepare_iso.sh "/path/to/Install OS X Mountain Lion.app" iso
```

This will generate a `.dmg` in the `iso` directory.

### Packer Templates

There are two generic packer JSON configurations. `packer/osx.json` and `packer/win.json`
Each of these files contains the necessary config to build OSX and Windows VM's for
use with the regression test framework

Each of these files takes a number of parameters that allow us to build different OS
versions from the same configuration. These are passed to packer using the `-var-file`
option. This options allows variables to be passed through from a JSON file.
See: `packer/osx-10.10.json` or `packer/win-10-pro-10586` as en example

### Adding a new version of an OS

If you'd like to add a new version of an OS the process is as follows.

1) Download (and in the case of OSX, generate) the ISO required for OS install
2) Copy an existing JSON variables file and change as appropriate
3) Ensure the ISO and SHA1 in the variables file are correct
4) Add the new version to the Makefile
5) Upload the ISO and the new machine to S3

### Adding a new OS

If the OS you are adding differs greatly from the existing supported OS's, you will have
to create a new JSON file. It's **HIGHLY** suggested that you look to projects like
[Bento](https://github.com/chef/bento) or one of the projects that we currently borrow
from rather than re-inventing the wheel

### Uploading to VMware

Once you have created your template, you may use it locally with `./rt-vm`.
If you want to upload to an ESX server, `./rt-vm upload` will handle this for you.

`./rt-vm upload etc/vm-templates/test-osx-10.10-template datastore-on-10`

The test infrastructure is very dependent on where the VM templates are placed, so please
be sure to consult with someone who understands the current design before attempting this.

It should be noted that CI should be "paused" during template upload

