---
title: Time random reads from a large file in a container
linkTitle: Random reads, container
description: Measure how long 25,000 random 4 KB reads from a 1 GB file take on a container's own storage.
keywords: docker desktop, performance, benchmark, fio, container storage, random read
weight: 120
---

This check makes the same random reads as [random reads in a shared
folder](random-read-shared.md), but on the container's own storage. It
represents a database or cache that keeps its data inside the container. `fio`
writes a 1 GB file, then makes 25 000 reads of 4 KB each at random locations.

Writing the file is part of the timed command, as it was for the published
number. The result is the wall-clock time of the single `docker run` command on
your machine.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/).

## Set up

```console
$ mkdir -p ~/perfcheck/randread-container && cd ~/perfcheck/randread-container

$ printf 'FROM alpine:3.24.1\nRUN apk add --no-cache fio=~3.41\n' | docker build -t perfcheck/fio:3.41 -

$ docker run --rm perfcheck/fio:3.41 fio --version
fio-3.41
```

## Run the check

1. Run the timed command:

   ```console
   $ time docker run --rm perfcheck/fio:3.41 \
     fio --name=randread --filename=/tmp/test.img --rw=randread --bs=4k \
         --ioengine=psync --number_ios=25000 --size=1G --invalidate=1
   ```

   Record the elapsed time.

2. Check that the run is valid. The `issued rwts:` line must start with
   `total=25000,0,0`. Docker removes the container and its file when it exits,
   so there's nothing to clean up.
