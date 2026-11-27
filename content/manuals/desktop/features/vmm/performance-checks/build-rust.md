---
title: Time a Rust compilation in an image build
linkTitle: Rust build
description: Measure how long docker build takes to compile ripgrep with cargo build --release.
keywords: docker desktop, performance, benchmark, docker build, rust, cargo, cpu, ripgrep
weight: 60
---

This check runs a real compilation the way many Docker builds do: it copies a
source tree into an image and runs `cargo build --release`. Unlike the other
build checks, the work is mostly CPU time inside the engine's virtual machine, so
the result depends on the number and speed of the cores the engine has.

The program is `ripgrep` at a pinned commit, built with a pinned Rust image. The
result is the wall-clock time of the single `docker build` command on your
machine, with nothing cached.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/_index.md).

> [!WARNING]
>
> Step 3 runs `docker buildx prune -a -f`, which deletes your entire local build
> cache. Don't run this check on a machine where you need that cache.

## Set up

```console
$ mkdir -p ~/perfcheck/rust-build && cd ~/perfcheck/rust-build

$ git init --bare -q ripgrep-mirror.git
$ git -C ripgrep-mirror.git remote add origin https://github.com/BurntSushi/ripgrep.git
$ git -C ripgrep-mirror.git fetch -q --depth=1 origin 4649aa9700619f94cf9c66876e9549d83420e16c
$ git -C ripgrep-mirror.git update-ref refs/heads/pinned FETCH_HEAD

$ docker pull rust:1.83-slim-bookworm
```

Pulling the Rust image keeps its download out of the timed build.

## Run the check

1. Create a fresh checkout. This step isn't timed.

   ```console
   $ rm -rf fixture && mkdir fixture
   $ git clone -q --no-local --no-checkout ripgrep-mirror.git fixture/repo
   $ git -C fixture/repo checkout -q --detach 4649aa9700619f94cf9c66876e9549d83420e16c
   $ cd fixture/repo
   ```

2. Write the build file. The `# syntax=` line must be the first line.

   ```console
   $ cat > Dockerfile.rust <<'EOF'
   # syntax=docker/dockerfile:1
   FROM rust:1.83-slim-bookworm
   WORKDIR /src
   COPY . .
   RUN cargo build --release --locked
   EOF
   ```

3. Empty the build cache so the compilation has nothing to reuse. This step
   isn't timed.

   ```console
   $ docker buildx prune -a -f
   ```

4. Run the timed command:

   ```console
   $ time docker build --no-cache --pull=false \
     -f Dockerfile.rust -t perfcheck-rust-build:ripgrep .
   ```

   Record the elapsed time.

   Because the cache is empty and the build uses `--no-cache`, the build can
   fetch the small helper image named by the `# syntax=` line. Don't download it
   in advance. It's part of the cost of a build from nothing.

5. Check that the run is valid. The compiler must have produced the program:

   ```console
   $ docker run --rm --entrypoint /src/target/release/rg perfcheck-rust-build:ripgrep --version
   ripgrep 14.1.1
   ```

6. Clean up:

   ```console
   $ cd ../..
   $ rm -rf fixture
   ```
