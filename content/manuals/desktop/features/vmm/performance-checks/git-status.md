---
title: Time the first git status in a shared repository
linkTitle: First git status
description: Measure how long the first git status takes in a container on a large repository shared from your machine.
keywords: docker desktop, performance, benchmark, git status, file sharing, virtiofs, bind mount
weight: 10
---

This check measures the first `git status` inside a container, on a freshly
created checkout of PostHog shared from your machine. The checkout has about
47 000 files and is 589 MB. The result is the duration of the single `git status`
command in milliseconds, measured inside the container. Nothing reads the files
through the container before the timed command.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/).

## Get the source

Download the pinned commit once.

```console
$ mkdir -p ~/perfcheck/git-status && cd ~/perfcheck/git-status

$ git init --bare -q posthog-mirror.git
$ git -C posthog-mirror.git remote add origin https://github.com/PostHog/posthog.git
$ git -C posthog-mirror.git fetch -q --depth=1 origin 4fa8e10ae2e988fd64ce93c8408856100af16c4f
$ git -C posthog-mirror.git update-ref refs/heads/pinned FETCH_HEAD
```

Confirm the commit:

```console
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

2. Wait, then refresh the Git index on your machine. If the files have the same
   timestamp as the index, Git re-reads their contents, which changes what the
   check measures.

   ```console
   $ sleep 2
   $ (cd fixture/repo && git update-index --really-refresh)
   ```

3. Start a container with the checkout shared in, and install Git. Neither
   command reads the checkout.

   ```console
   $ docker run -d --rm --name perfcheck-git-status \
     -v "$PWD/fixture/repo:/work" -w /work \
     alpine:3.24.1 sleep infinity

   $ docker exec perfcheck-git-status apk add --no-cache git=~2.54.0 coreutils
   ```

4. Confirm the pinned version, because a different Git version does a different
   amount of work:

   ```console
   $ docker exec perfcheck-git-status git --version
   git version 2.54.0
   ```

5. Run the timed command. The `-c` options turn off Git's optional speed-ups so
   your global Git configuration can't change the result.

   ```console
   $ docker exec perfcheck-git-status sh -c '
     set -e
     start=$(date +%s%3N)
     git -c safe.directory=/work \
         -c core.preloadindex=true \
         -c core.fsmonitor=false \
         -c core.untrackedCache=false \
         status --porcelain=v1 --untracked-files=all > /tmp/status.out
     end=$(date +%s%3N)
     printf "git_status_ms=%s\n" "$((end - start))"
     if [ -s /tmp/status.out ]; then echo "INVALID: the checkout is not clean"; fi
   '
   ```

   The value of `git_status_ms` is the result.

6. Check that the run is valid. The command must exit without an error and
   print no `INVALID` line. A non-empty status means either the checkout was
   modified or the file sharing reports file information in a way that makes Git
   think files changed. Discard the run.

7. Clean up:

   ```console
   $ docker rm -f perfcheck-git-status
   $ rm -rf fixture
   ```

## Results that aren't valid

- You ran `git status` twice on the same checkout. The second run uses a warm
  index and is a different measurement.
- Your machine has a non-default Git configuration, such as a global
  `core.autocrlf`, a file system monitor, or a clean or smudge filter.
- An editor or file indexer opened the checkout.
