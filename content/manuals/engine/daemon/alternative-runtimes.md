---
title: Alternative container runtimes
description: |
  Docker Engine uses runc as the default container runtime, but you
  can specify alternative runtimes using the CLI or by configuring
  the daemon
keywords: engine, runtime, containerd, runtime v2, shim
aliases:
  - /engine/alternative-runtimes/
---

Docker Engine uses containerd for managing the container lifecycle,
which includes creating, starting, and stopping containers.
By default, containerd uses runc as its container runtime.

## What runtimes can I use?

You can use any runtime that implements the containerd 
[shim API](https://github.com/containerd/containerd/blob/main/core/runtime/v2/README.md).
Such runtimes ship with a containerd shim, and you can use them without any
additional configuration. See [Use containerd shims](#use-containerd-shims).

Examples of runtimes that implement their own containerd shims include:

- [Wasmtime](https://wasmtime.dev/)
- [gVisor](https://github.com/google/gvisor)
- [Kata Containers](https://katacontainers.io/)

You can also use runtimes designed as drop-in replacements for runc. Such
runtimes depend on the runc containerd shim for invoking the runtime binary.
You must manually register such runtimes in the daemon configuration.

[youki](https://github.com/youki-dev/youki)
is one example of a runtime that can function as a runc drop-in replacement.
Refer to the [youki example](#youki) explaining the setup.

## Use containerd shims

containerd shims let you use alternative runtimes without having to change the
configuration of the Docker daemon. To use a containerd shim, install the shim
binary on `PATH` on the system where the Docker daemon is running.

To use a shim with `docker run`, specify the fully qualified name of the
runtime as the value to the `--runtime` flag:

```console
$ docker run --runtime io.containerd.kata.v2 hello-world
```

### Use a containerd shim without installing on PATH

You can use a shim without installing it on `PATH`, in which case you need to
register the shim in the daemon configuration as follows:

```json
{
  "runtimes": {
    "foo": {
      "runtimeType": "/path/to/containerd-shim-foobar-v1"
    }
  }
}
```

To use the shim, specify the name that you assigned to it:

```console
$ docker run --runtime foo hello-world
```

### Configure shims

If you need to pass additional configuration for a containerd shim, you can
use the `runtimes` option in the daemon configuration file.

1. Edit the daemon configuration file by adding a `runtimes` entry for the
   shim you want to configure.

   - Specify the fully qualified name for the runtime in `runtimeType` key
   - Add your runtime configuration under the `options` key

   ```json
   {
     "runtimes": {
       "gvisor": {
         "runtimeType": "io.containerd.runsc.v1",
         "options": {
           "TypeUrl": "io.containerd.runsc.v1.options",
           "ConfigPath": "/etc/containerd/runsc.toml"
         }
       }
     }
   }
   ```

2. Reload the daemon's configuration.

   ```console
   # systemctl reload docker
   ```

3. Use the customized runtime using the `--runtime` flag for `docker run`.

   ```console
   $ docker run --runtime gvisor hello-world
   ```

For more information about the configuration options for containerd shims, see
[Configure containerd shims](/reference/cli/dockerd.md#configure-containerd-shims).

## Examples

The following examples show you how to set up and use alternative container
runtimes with Docker Engine.

- [youki](#youki)
- [Wasmtime](#wasmtime)

### youki

youki is a container runtime written in Rust.
youki claims to be faster and use less memory than runc,
making it a good choice for resource-constrained environments.

youki functions as a drop-in replacement for runc, meaning it relies on the
runc shim to invoke the runtime binary. When you register runtimes acting as
runc replacements, you configure the path to the runtime executable, and
optionally a set of runtime arguments. For more information, see
[Configure runc drop-in replacements](/reference/cli/dockerd.md#configure-runc-drop-in-replacements).

To add youki as a container runtime:

1. Install youki and its dependencies.

   For instructions, refer to the
   [official setup guide](https://youki-dev.github.io/youki/user/basic_setup.html).

2. Register youki as a runtime for Docker by editing the Docker daemon
   configuration file, located at `/etc/docker/daemon.json` by default.

   The `path` key should specify the path to wherever you installed youki.

   ```console
   # cat > /etc/docker/daemon.json <<EOF
   {
     "runtimes": {
       "youki": {
         "path": "/usr/local/bin/youki"
       }
     }
   }
   EOF
   ```

3. Reload the daemon's configuration.

   ```console
   # systemctl reload docker
   ```

Now you can run containers that use youki as a runtime.

```console
$ docker run --rm --runtime youki hello-world
```

### Wasmtime

{{< summary-bar feature_name="Wasmtime" >}}

Wasmtime is a
[Bytecode Alliance](https://bytecodealliance.org/)
project, and a Wasm runtime that lets you run Wasm containers.
Running Wasm containers with Docker provides two layers of security.
You get all the benefits from container isolation,
plus the added sandboxing provided by the Wasm runtime environment.

To add Wasmtime as a container runtime, follow these steps:

1. Turn on the [containerd image store](/manuals/engine/storage/containerd.md)
   feature in the daemon configuration file.

   ```json
   {
     "features": {
       "containerd-snapshotter": true
     }
   }
   ```

2. Restart the Docker daemon.

   ```console
   # systemctl restart docker
   ```

3. Install the Wasmtime containerd shim on `PATH`.

   The following command Dockerfile builds the Wasmtime binary from source
   and exports it to `./containerd-shim-wasmtime-v1`.

   ```console
   $ docker build --output . - <<EOF
   FROM rust:latest as build
   RUN cargo install \
       --git https://github.com/containerd/runwasi.git \
       --bin containerd-shim-wasmtime-v1 \
       --root /out \
       containerd-shim-wasmtime
   FROM scratch
   COPY --from=build /out/bin /
   EOF
   ```

   Put the binary in a directory on `PATH`.

   ```console
   $ mv ./containerd-shim-wasmtime-v1 /usr/local/bin
   ```

Now you can run containers that use Wasmtime as a runtime.

```console
$ docker run --rm \
 --runtime io.containerd.wasmtime.v1 \
 --platform wasi/wasm32 \
 michaelirwin244/wasm-example
```

## Related information

- To learn more about the configuration options for container runtimes,
  see [Configure container runtimes](/reference/cli/dockerd.md#configure-container-runtimes).
- You can configure which runtime that the daemon should use as its default.
  Refer to [Configure the default container runtime](/reference/cli/dockerd.md#configure-the-default-container-runtime).
