---
title: Time the first text search in a shared repository
linkTitle: First text search
description: Measure how long a ripgrep search takes in a container on a large repository shared from your machine.
keywords: docker desktop, performance, benchmark, ripgrep, file sharing, virtiofs, bind mount
weight: 20
---

This check measures one `ripgrep` search for a fixed string across a freshly
created checkout of PostHog shared in from your machine. The checkout has about
47 000 files and is 589 MB. Nothing reads the files through the container
before the search.

The result is the duration of the single search command in milliseconds,
measured inside the container. `ripgrep` chooses its own thread count. The
parallelism your engine gives a container is part of what you measure.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/_index.md).

## Get the source

```console
$ mkdir -p ~/perfcheck/code-search && cd ~/perfcheck/code-search

$ git init --bare -q posthog-mirror.git
$ git -C posthog-mirror.git remote add origin https://github.com/PostHog/posthog.git
$ git -C posthog-mirror.git fetch -q --depth=1 origin 4fa8e10ae2e988fd64ce93c8408856100af16c4f
$ git -C posthog-mirror.git update-ref refs/heads/pinned FETCH_HEAD

$ git -C posthog-mirror.git rev-parse refs/heads/pinned
4fa8e10ae2e988fd64ce93c8408856100af16c4f
```

## Run the check

1. Create a fresh checkout. This step isn't timed.

   ```console
   $ rm -rf fixture && mkdir fixture
   $ git clone -q --no-local --no-checkout posthog-mirror.git fixture/repo
   $ git -C fixture/repo checkout -q --detach 4fa8e10ae2e988fd64ce93c8408856100af16c4f
   ```

2. Start a container with the checkout shared in, and install `ripgrep`.
   Neither command reads the checkout.

   ```console
   $ docker run -d --rm --name perfcheck-code-search \
     -v "$PWD/fixture/repo:/work" -w /work \
     alpine:3.24.1 sleep infinity

   $ docker exec perfcheck-code-search apk add --no-cache ripgrep=~15.1.0 coreutils
   ```

3. Confirm the version. Different `ripgrep` releases search at different
   speeds.

   ```console
   $ docker exec perfcheck-code-search sh -c 'rg --version | head -1'
   ripgrep 15.1.0
   ```

4. Run the timed command:

   ```console
   $ docker exec perfcheck-code-search sh -c '
     set -e
     start=$(date +%s%3N)
     rg --no-config --hidden --glob "!.git/**" --fixed-strings posthog . >/dev/null
     end=$(date +%s%3N)
     printf "search_ms=%s\n" "$((end - start))"
   '
   ```

   The value of `search_ms` is the result.

5. Check that the run is valid. Run the same search again, without timing it,
   and compare the output with the expected values:

   ```console
   $ docker exec perfcheck-code-search sh -c '
     rg --no-config --hidden --glob "!.git/**" --fixed-strings posthog . > /tmp/search.out
     printf "lines=%s\n"  "$(wc -l < /tmp/search.out | tr -d " ")"
     printf "digest=%s\n" "$(sed "s/\r$//" /tmp/search.out | LC_ALL=C sort | sha256sum | cut -d " " -f 1)"
   '
   ```

   Expected output:

   ```text
   lines=177556
   digest=7b8e6be6abc6421718afda2e33a16b59b137f300895b75a2cea42ce98a9de996
   ```

   Any other output means you searched a different set of files or used a
   different `ripgrep` version. The duration isn't comparable.

6. Clean up:

   ```console
   $ docker rm -f perfcheck-code-search
   $ rm -rf fixture
   ```
