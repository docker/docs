---
title: Time many small file operations in a shared folder
linkTitle: Small files
description: Measure how long it takes to create, stat, read, and delete 10,000 small files in a folder shared from your machine.
keywords: docker desktop, performance, benchmark, file sharing, small files, virtiofs, bind mount
weight: 90
---

Installing dependencies, unpacking an archive, or running a code generator in
your project folder creates, checks, reads, and removes many small files through
the engine's file sharing. This check does the same in four steps on a folder
shared from your machine:

1. Create 10 000 one-line files.
2. Look up the size of each file.
3. Read every file.
4. Delete every file.

Each step runs in its own short-lived container with the folder shared in as
`-v "$PWD:/bind"`. The result is the wall-clock time of all four `docker run`
commands together on your machine, including container startup.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/_index.md).

## Set up

Pull the image so its download stays out of the timing:

```console
$ mkdir -p ~/perfcheck/small-files && cd ~/perfcheck/small-files

$ docker pull alpine:3.24.1
```

## Run the check

1. Start from an empty shared folder. This step isn't timed.

   ```console
   $ rm -rf shared && mkdir shared && cd shared
   ```

2. Run the timed command:

   ```console
   $ time (
     docker run --rm -v "$PWD:/bind" alpine:3.24.1 sh -c 'cd /bind && i=1; while [ $i -le 10000 ]; do echo $i > f$i; i=$((i+1)); done' &&
     docker run --rm -v "$PWD:/bind" alpine:3.24.1 sh -c 'find /bind -type f -exec stat -c "%s" {} + > /dev/null' &&
     docker run --rm -v "$PWD:/bind" alpine:3.24.1 sh -c 'find /bind -type f -exec cat {} + > /dev/null' &&
     docker run --rm -v "$PWD:/bind" alpine:3.24.1 sh -c 'find /bind -type f -delete'
   )
   ```

   Record the elapsed time.

3. Check that the run is valid. The timed block must finish without an error,
   and the last step must have removed every file:

   ```console
   $ ls -A | wc -l
   0
   ```

4. Clean up:

   ```console
   $ cd ..
   $ rm -rf shared
   ```
