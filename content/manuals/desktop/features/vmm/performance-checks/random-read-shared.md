---
title: Time random reads from a large file in a shared folder
linkTitle: Random reads, shared folder
description: Measure how long 25,000 random 4 KB reads from a 1 GB file take on a folder shared from your machine.
keywords: docker desktop, performance, benchmark, fio, file sharing, virtiofs, bind mount, random read
weight: 110
---

A database, search index, or cache in your project folder reads small pieces of
a large file in no particular order. This check measures that through the
engine's file sharing. It makes 25 000 reads of 4 KB each at random locations in a
1 GB file shared from your machine, using `fio`.

You write the file before timing starts, so only the reads are measured. The
result is the wall-clock time of the single `docker run` command on your
machine.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/).

## Set up

Build an image with a pinned `fio` so the installation stays out of the timing:

```console
$ mkdir -p ~/perfcheck/randread-shared && cd ~/perfcheck/randread-shared

$ printf 'FROM alpine:3.24.1\nRUN apk add --no-cache fio=~3.41\n' | docker build -t perfcheck/fio:3.41 -

$ docker run --rm perfcheck/fio:3.41 fio --version
fio-3.41
```

## Run the check

1. Write a fresh 1 GB file into an empty shared folder. This step isn't timed.

   ```console
   $ rm -rf shared && mkdir shared && cd shared
   $ docker run --rm -v "$PWD:/work" -w /work alpine:3.24.1 \
     sh -c 'dd if=/dev/zero of=test.img bs=1M count=1024 conv=fsync'
   ```

2. Run the timed command:

   ```console
   $ time docker run --rm -v "$PWD:/work" -w /work perfcheck/fio:3.41 \
     fio --name=randread --filename=test.img --rw=randread --bs=4k \
         --ioengine=psync --number_ios=25000 --invalidate=1
   ```

   Record the elapsed time.

3. Check that the run is valid. In the `fio` summary, the `issued rwts:` line
   must start with `total=25000,0,0`: 25 000 reads and no writes.

4. Clean up:

   ```console
   $ cd ..
   $ rm -rf shared
   ```
