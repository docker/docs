---
datafolder: engine-cli
datafile: docker_container_cp
title: docker container cp
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

Suppose a container has finished producing some output as a file it saves
to somewhere in its filesystem. This could be the output of a build job or
some other computation. You can copy these outputs from the container to a
location on your local host.

If you want to copy the `/tmp/foo` directory from a container to the
existing `/tmp` directory on your host. If you run `docker container cp` in your `~`
(home) directory on the local host:

    $ docker container cp compassionate_darwin:tmp/foo /tmp

Docker creates a `/tmp/foo` directory on your host. Alternatively, you can omit
the leading slash in the command. If you execute this command from your home
directory:

    $ docker container cp compassionate_darwin:tmp/foo tmp

If `~/tmp` does not exist, Docker will create it and copy the contents of
`/tmp/foo` from the container into this new directory. If `~/tmp` already
exists as a directory, then Docker will copy the contents of `/tmp/foo` from
the container into a directory at `~/tmp/foo`.

When copying a single file to an existing `LOCALPATH`, the `docker container cp` command
will either overwrite the contents of `LOCALPATH` if it is a file or place it
into `LOCALPATH` if it is a directory, overwriting an existing file of the same
name if one exists. For example, this command:

    $ docker container cp sharp_ptolemy:/tmp/foo/myfile.txt /test

If `/test` does not exist on the local machine, it will be created as a file
with the contents of `/tmp/foo/myfile.txt` from the container. If `/test`
exists as a file, it will be overwritten. Lastly, if `/test` exists as a
directory, the file will be copied to `/test/myfile.txt`.

Next, suppose you want to copy a file or folder into a container. For example,
this could be a configuration file or some other input to a long running
computation that you would like to place into a created container before it
starts. This is useful because it does not require the configuration file or
other input to exist in the container image.

If you have a file, `config.yml`, in the current directory on your local host
and wish to copy it to an existing directory at `/etc/my-app.d` in a container,
this command can be used:

    $ docker container cp config.yml myappcontainer:/etc/my-app.d

If you have several files in a local directory `/config` which you need to copy
to a directory `/etc/my-app.d` in a container:

    $ docker container cp /config/. myappcontainer:/etc/my-app.d

The above command will copy the contents of the local `/config` directory into
the directory `/etc/my-app.d` in the container.

Finally, if you want to copy a symbolic link into a container, you typically
want to  copy the linked target and not the link itself. To copy the target, use
the `-L` option, for example:

    $ ln -s /tmp/somefile /tmp/somefile.ln
    $ docker container cp -L /tmp/somefile.ln myappcontainer:/tmp/

This command copies content of the local `/tmp/somefile` into the file
`/tmp/somefile.ln` in the container. Without `-L` option, the `/tmp/somefile.ln`
preserves its symbolic link but not its content.
