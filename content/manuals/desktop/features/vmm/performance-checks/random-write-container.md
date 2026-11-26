---
title: Time random writes to a large file in a container
linkTitle: Random writes, container
description: Measure how long 25,000 random 4 KB writes to a 512 MB file take on a container's own storage.
keywords: docker desktop, performance, benchmark, fio, container storage, random write
weight: 130
---

A database or cache inside a container updates small pieces of its files in no
particular order. This check makes 25 000 writes of 4 KB each at random
locations in a 512 MB file on the container's own storage, using `fio`. The
writes aren't forced to disk one at a time, which matches how most applications
write.

The result is the wall-clock time of the single `docker run` command on your
machine, including the time `fio` takes to create the file.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/).

## Set up

```console
$ mkdir -p ~/perfcheck/randwrite-container && cd ~/perfcheck/randwrite-container

$ printf 'FROM alpine:3.24.1\nRUN apk add --no-cache fio=~3.41\n' | docker build -t perfcheck/fio:3.41 -

$ docker run --rm perfcheck/fio:3.41 fio --version
fio-3.41
```

## Run the check

1. Run the timed command:

   ```console
   $ time docker run --rm perfcheck/fio:3.41 \
     fio --name=randwrite --filename=/tmp/test.img --rw=randwrite --bs=4k \
         --ioengine=psync --number_ios=25000 --size=512M --invalidate=1
   ```

   Record the elapsed time.

2. Check that the run is valid. The `issued rwts:` line must start with
   `total=0,25000,0`: no reads and 25 000 writes. Docker removes the container
   and its file when it exits, so there's nothing to clean up.
