---
title: Build release notes
weight: 120
description: Learn about the new features, bug fixes, and breaking changes for the newest Buildx release
keywords: build, buildx, buildkit, release notes
tags: [Release notes]
toc_max: 2
---

This page contains information about the new features, improvements, and bug
fixes in [Docker Buildx](https://github.com/docker/buildx).

## 0.29.0

{{< release-date date="2025-09-30" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.29.0).

### New

- New `--progress=none` option has been added. This is similar to `--progress=quiet`, but it does not print the image ID after image result export. [docker/buildx#3431](https://github.com/docker/buildx/pull/3431)
- Compose compatibility has been updated to v2.9.0.

### Enhancements

- `imagetools create` command now supports `--platform` option to create final image only for specified platforms. The inline attestation for the specified platforms are also kept in the final image. [docker/buildx#3430](https://github.com/docker/buildx/pull/3430)
- DAP debugger can now show the correct file explorer data when the debugger stops because of a build error. [docker/buildx#3410](https://github.com/docker/buildx/pull/3410)
- When building from a Git URL, buildx now optionally supports resolution of the context data on the client side. Git repository is still cloned on the server side, but this can help in cases where one can't be sure what version of Git URL resolution the server side supports. [docker/buildx#3415](https://github.com/docker/buildx/pull/3415)

### Bug fixes

- Fix DAP debugger location resolution when there are multiple build steps with the same BuildKit digest. [docker/buildx#3408](https://github.com/docker/buildx/pull/3408)

## 0.28.0

{{< release-date date="2025-09-03" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.28.0).

### New

- When building with Dockerfile version 1.18.0 or later, you can now use new Git URLs with query options for build context and named contexts in the `build` and `bake` command. [dockerfile/1.18.0](/manuals/build/buildkit/dockerfile-release-notes.md#1180)

### Enhancements

- Add formatting options to the `buildx du` command for custom and machine-readable output. [docker/buildx#3377](https://github.com/docker/buildx/pull/3377)
- Kubernetes driver now supports `env.<key>` driver opts [docker/buildx#3373](https://github.com/docker/buildx/pull/3373)
- Add support for `BUILDKIT_SYNTAX` build argument when BuildKit has a Dockerfile frontend disabled. [docker/buildx#3385](https://github.com/docker/buildx/pull/3385)

### Bug fixes

- Fix failing early when trying to export index annotations with moby exporter. [docker/buildx#3384](https://github.com/docker/buildx/pull/3384)
- Fix possible errors on Windows from symlink handling [docker/buildx#3386](https://github.com/docker/buildx/pull/3386)

## 0.27.0

{{< release-date date="2025-08-20" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.27.0).

### New

- Compose compatibility has been updated to v2.8.1. [docker/buildx#3337](https://github.com/docker/buildx/pull/3337)

### Enhancements

- DAP: Exec shell now restarts with the new container when execution resumes and pauses again. [docker/buildx#3341](https://github.com/docker/buildx/pull/3341)
- DAP: Add `File Explorer` section to variables to inspect filesystem state. [docker/buildx#3327](https://github.com/docker/buildx/pull/3327)
- DAP: Change Dockerfile step order to match more closely with user expectations. [docker/buildx#3325](https://github.com/docker/buildx/pull/3325)
- DAP: Improve determination of the proper parent. [docker/buildx#3366](https://github.com/docker/buildx/pull/3366)
- DAP: Dockerfile nested in the context is now supported. [docker/buildx#3371](https://github.com/docker/buildx/pull/3371)
- Build name shown in history can now be overridden with `BUILDKIT_BUILD_NAME` build argument. [docker/buildx#3330](https://github.com/docker/buildx/pull/3330)
- Bake now supports `homedir()` function. [docker/buildx#3351](https://github.com/docker/buildx/pull/3351)
- Bake default for empty Dockerfile defaults to `Dockerfile` to match the behavior of `build` command. [docker/buildx#3347](https://github.com/docker/buildx/pull/3347)
- Bake supports `pull` and `no_cache` fields for Compose files. [docker/buildx#3352](https://github.com/docker/buildx/pull/3352)
- Sanitize the names of `additional_contexts` from Compose files when building with Bake. [docker/buildx#3361](https://github.com/docker/buildx/pull/3361)

### Bug fixes

- Fix missing WSL libraries in `docker-container` driver when GPU device is requested. [docker/buildx#3320](https://github.com/docker/buildx/pull/3320)

## 0.26.1

{{< release-date date="2025-07-22" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.26.1).

### Bug fixes

- Fix regression when validating compose files with Bake. [docker/buildx#3329](https://github.com/docker/buildx/pull/3329)

## 0.26.0

{{< release-date date="2025-07-21" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.26.0).

### New

- New experimental version of the DAP debugger has been added with a new `dap build` helper command. The new feature can be tried with the [DockerDX VSCode extension](https://github.com/docker/vscode-extension). [docker/buildx#3235](https://github.com/docker/buildx/pull/3235)
- Compose compatibility has been updated to v2.7.1. [docker/buildx#3282](https://github.com/docker/buildx/pull/3282)

### Enhancements

- Bake command now supports pattern-matching target names with wildcards. [docker/buildx#3280](https://github.com/docker/buildx/pull/3280)
- Bake command now supports setting files through environment variable `BUILDX_BAKE_FILE`. [docker/buildx#3242](https://github.com/docker/buildx/pull/3242)
- Bake now ignores unrelated fields when parsing and validating compose files. [docker/buildx#3292](https://github.com/docker/buildx/pull/3292)
- `history` commands will automatically bootstrap the builder. [docker/buildx#3300](https://github.com/docker/buildx/pull/3300)
- Add SLSA v1 support to `history inspect` command. [docker/buildx#3245](https://github.com/docker/buildx/pull/3245)
- Kubernetes driver option `buildkit-root-volume-memory` to use memory mount for the root volume. [docker/buildx#3253](https://github.com/docker/buildx/pull/3253)

### Bug fixes

- Fix possible error from `imagetools` commands when accessing registries that don't return content length. [docker/buildx#3316](https://github.com/docker/buildx/pull/3316)
- Fix duplicated command descriptions from help output. [docker/buildx#3298](https://github.com/docker/buildx/pull/3298)
- Fix `history inspect attachment` to not require an argument. [docker/buildx#3264](https://github.com/docker/buildx/pull/3264)
- Fix resolving environment variables from `.env` file when building compose files with Bake. [docker/buildx#3275](https://github.com/docker/buildx/pull/3275), [docker/buildx#3276](https://github.com/docker/buildx/pull/3276), [docker/buildx#3322](https://github.com/docker/buildx/pull/3322)

## 0.25.0

{{< release-date date="2025-06-17" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.25.0).

### New

- Bake now supports defining `extra-hosts`. [docker/buildx#3234](https://github.com/docker/buildx/pull/3234)

### Enhancements

- Add support for bearer token auth. [docker/buildx#3233](https://github.com/docker/buildx/pull/3233)
- Add custom exit codes for internal, resource, and canceled errors in commands. [docker/buildx#3214](https://github.com/docker/buildx/pull/3214)
- Show variable type when using `--list=variables` with Bake. [docker/buildx#3207](https://github.com/docker/buildx/pull/3207)
- Consider typed, value-less variables to have `null` value in Bake. [docker/buildx#3198](https://github.com/docker/buildx/pull/3198)
- Add support for multiple IPs in extra hosts configuration. [docker/buildx#3244](https://github.com/docker/buildx/pull/3244)
- Support for updated SLSA V1 provenance in `buildx history` commands. [docker/buildx#3245](https://github.com/docker/buildx/pull/3245)
- Add support for `RegistryToken` configuration in imagetools commands. [docker/buildx#3233](https://github.com/docker/buildx/pull/3233)

### Bug fixes

- Fix `keep-storage` flag deprecation notice for `prune` command. [docker/buildx#3216](https://github.com/docker/buildx/pull/3216)

## 0.24.0

{{< release-date date="2025-05-21" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.24.0).

### Enhancements

- New `type` attribute added to the `variable` block in Bake to allow explicit typing of variables. [docker/buildx#3167](https://github.com/docker/buildx/pull/3167), [docker/buildx#3189](https://github.com/docker/buildx/pull/3189), [docker/buildx#3198](https://github.com/docker/buildx/pull/3198)
- New `--finalize` flag added to the `history export` command to finalize build traces before exporting. [docker/buildx#3152](https://github.com/docker/buildx/pull/3152)
- Compose compatibility has been updated to v2.6.3. [docker/buildx#3191](https://github.com/docker/buildx/pull/3191), [docker/buildx#3171](https://github.com/docker/buildx/pull/3171)

### Bug fixes

- Fix issue where some builds may leave behind temporary files after completion. [docker/buildx#3133](https://github.com/docker/buildx/pull/3133)
- Fix wrong image ID returned when building with Docker when containerd-snapshotter is enabled. [docker/buildx#3136](https://github.com/docker/buildx/pull/3136)
- Fix possible panic when using empty `call` definition with Bake. [docker/buildx#3168](https://github.com/docker/buildx/pull/3168)
- Fix possible malformed Dockerfile path with Bake on Windows. [docker/buildx#3141](https://github.com/docker/buildx/pull/3141)
- Fix current builder not being available in JSON output for `ls` command. [docker/buildx#3179](https://github.com/docker/buildx/pull/3179)
- Fix OTEL context not being propagated to Docker daemon. [docker/buildx#3146](https://github.com/docker/buildx/pull/3146)

## 0.23.0

{{< release-date date="2025-04-15" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.23.0).

### New

- New `buildx history export` command allows exporting the build record into a bundle that can be imported to [Docker Desktop](/desktop/). [docker/buildx#3073](https://github.com/docker/buildx/pull/3073)

### Enhancements

- New `--local` and `--filter` flags allow filtering history records in `buildx history ls`. [docker/buildx#3091](https://github.com/docker/buildx/pull/3091)
- Compose compatibility has been updated to v2.6.0. [docker/buildx#3080](https://github.com/docker/buildx/pull/3080), [docker/buildx#3105](https://github.com/docker/buildx/pull/3105)
- Support CLI environment variables in standalone mode. [docker/buildx#3087](https://github.com/docker/buildx/pull/3087)

### Bug fixes

- Fix `--print` output for Bake producing output with unescaped variables that could cause build errors later. [docker/buildx#3097](https://github.com/docker/buildx/pull/3097)
- Fix `additional_contexts` field not working correctly when pointing to another service. [docker/buildx#3090](https://github.com/docker/buildx/pull/3090)
- Fix empty validation block crashing the Bake HCL parser. [docker/buildx#3101](https://github.com/docker/buildx/pull/3101)

## 0.22.0

{{< release-date date="2025-03-18" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.22.0).

### New

- New command `buildx history import` lets you import build records into Docker Desktop for further debugging in the [Build UI](/desktop/use-desktop/builds/). This command requires [Docker Desktop](/desktop/) to be installed. [docker/buildx#3039](https://github.com/docker/buildx/pull/3039)

### Enhancements

- History records can now be opened by offset from the latest in `history inspect`, `history logs` and `history open` commands (e.g. `^1`). [docker/buildx#3049](https://github.com/docker/buildx/pull/3049), [docker/buildx#3055](https://github.com/docker/buildx/pull/3055)
- Bake now supports the `+=` operator to append when using `--set` for overrides. [docker/buildx#3031](https://github.com/docker/buildx/pull/3031)
- Docker container driver adds GPU devices to the container if available. [docker/buildx#3063](https://github.com/docker/buildx/pull/3063)
- Annotations can now be set when using overrides with Bake. [docker/buildx#2997](https://github.com/docker/buildx/pull/2997)
- NetBSD binaries are now included in the release. [docker/buildx#2901](https://github.com/docker/buildx/pull/2901)
- The `inspect` and `create` commands now return an error if a node fails to boot. [docker/buildx#3062](https://github.com/docker/buildx/pull/3062)

### Bug fixes

- Fix double pushing with Docker driver when the containerd image store is enabled. [docker/buildx#3023](https://github.com/docker/buildx/pull/3023)
- Fix multiple tags being pushed for `imagetools create` command. Now only the final manifest pushes by tag. [docker/buildx#3024](https://github.com/docker/buildx/pull/3024)

## 0.21.0

{{< release-date date="2025-02-19" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.21.0).

### New

- New command `buildx history trace` lets you inspect traces of a build in a Jaeger UI-based viewer and compare one trace with another. [docker/buildx#2904](https://github.com/docker/buildx/pull/2904)

### Enhancements

- The history inspection command `buildx history inspect` now supports custom formatting with `--format` flag and JSON formatting for machine-readable output. [docker/buildx#2964](https://github.com/docker/buildx/pull/2964)
- Support for CDI device entitlement in build and bake. [docker/buildx#2994](https://github.com/docker/buildx/pull/2994)
- Supported CDI devices are now shown in the builder inspection. [docker/buildx#2983](https://github.com/docker/buildx/pull/2983)
- When using [GitHub Cache backend `type=gha`](cache/backends/gha.md), the URL of the Version 2 or API is now read from the environment and sent to BuildKit. Version 2 backend requires BuildKit v0.20.0 or later. [docker/buildx#2983](https://github.com/docker/buildx/pull/2983), [docker/buildx#3001](https://github.com/docker/buildx/pull/3001)

### Bug fixes

- Avoid unnecessary warnings and prompts when using `--progress=rawjson`. [docker/buildx#2957](https://github.com/docker/buildx/pull/2957)
- Fix regression with debug shell sometimes not working correctly on `--on=error`. [docker/buildx#2958](https://github.com/docker/buildx/pull/2958)
- Fix possible panic errors when using an unknown variable in the Bake definition. [docker/buildx#2960](https://github.com/docker/buildx/pull/2960)
- Fix invalid duplicate output on JSON format formatting of `buildx ls` command. [docker/buildx#2970](https://github.com/docker/buildx/pull/2970)
- Fix bake handling cache imports with CSV string containing multiple registry references. [docker/buildx#2944](https://github.com/docker/buildx/pull/2944)
- Fix issue where error from pulling BuildKit image could be ignored. [docker/buildx#2988](https://github.com/docker/buildx/pull/2988)
- Fix race on pausing progress on debug shell. [docker/buildx#3003](https://github.com/docker/buildx/pull/3003)

## 0.20.1

{{< release-date date="2025-01-23" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.20.1).

### Bug fixes

- Fix `bake --print` output after missing some attributes for attestations. [docker/buildx#2937](https://github.com/docker/buildx/pull/2937)
- Fix allowing comma-separated image reference strings for cache import and export values. [docker/buildx#2944](https://github.com/docker/buildx/pull/2944)

## 0.20.0

{{< release-date date="2025-01-20" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.20.0).

> [!NOTE]
>
> This version of buildx enables filesystem entitlement checks for `buildx bake`
> command by default. If your Bake definition needs to read or write files
> outside your current working directory, you need to allow access to these
> paths with `--allow fs=<path|*>`. On the terminal, you can also interactively
> approve these paths with the provided prompt. Optionally, you can disable
> these checks by setting `BUILDX_BAKE_ENTITLEMENTS_FS=0`. This validation
> produced a warning in Buildx v0.19.0+, but starting from current release it
> produces an error. For more information, see the [reference documentation](/reference/cli/docker/buildx/bake.md#allow).

### New

- New `buildx history` command has been added that allows working with build records of completed and running builds. You can use these commands to list, inspect, remove your builds, replay the logs of already completed builds, and quickly open your builds in Docker Desktop Build UI for further debugging. This is an early version of this command and we expect to add more features in the future releases. [#2891](https://github.com/docker/buildx/pull/2891), [#2925](https://github.com/docker/buildx/pull/2925)

### Enhancements

- Bake: Definition now supports new object notation for the fields that previously required CSV strings as inputs (`attest`, `output`, `cache-from`, `cache-to`, `secret`, `ssh`). [docker/buildx#2758](https://github.com/docker/buildx/pull/2758), [docker/buildx#2848](https://github.com/docker/buildx/pull/2848), [docker/buildx#2871](https://github.com/docker/buildx/pull/2871), [docker/buildx#2814](https://github.com/docker/buildx/pull/2814)
- Bake: Filesystem entitlements now error by default. To disable this behavior, you can set `BUILDX_BAKE_ENTITLEMENTS_FS=0`. [docker/buildx#2875](https://github.com/docker/buildx/pull/2875)
- Bake: Infer Git authentication token from remote files to build request. [docker/buildx#2905](https://github.com/docker/buildx/pull/2905)
- Bake: Add support for `--list` flag to list targets and variables. [docker/buildx#2900](https://github.com/docker/buildx/pull/2900), [docker/buildx#2907](https://github.com/docker/buildx/pull/2907)
- Bake: Update lookup order for default definition files to load the files with "override" suffix later. [docker/buildx#2886](https://github.com/docker/buildx/pull/2886)

### Bug fixes

- Bake: Fix entitlements check for default SSH socket. [docker/buildx#2898](https://github.com/docker/buildx/pull/2898)
- Bake: Fix missing default target in group's default targets. [docker/buildx#2863](https://github.com/docker/buildx/pull/2863)
- Bake: Fix named context from target platform matching. [docker/buildx#2877](https://github.com/docker/buildx/pull/2877)
- Fix missing documentation for quiet progress mode. [docker/buildx#2899](https://github.com/docker/buildx/pull/2899)
- Fix missing last progress from loading layers. [docker/buildx#2876](https://github.com/docker/buildx/pull/2876)
- Validate BuildKit configuration before creating a builder. [docker/buildx#2864](https://github.com/docker/buildx/pull/2864)

### Packaging

- Compose compatibility has been updated to v2.4.7. [docker/buildx#2893](https://github.com/docker/buildx/pull/2893), [docker/buildx#2857](https://github.com/docker/buildx/pull/2857), [docker/buildx#2829](https://github.com/docker/buildx/pull/2829)

## 0.19.1

{{< release-date date="2024-11-27" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.19.1).

### Bug fixes

- Reverted the change in v0.19.0 that added new object notation for the fields
  that previously required CSV strings in Bake definition. This enhancement was
  reverted because of backwards incompatibility issues were discovered in some
  edge cases. This feature has now been postponed to the v0.20.0 release.
  [docker/buildx#2824](https://github.com/docker/buildx/pull/2824)

## 0.19.0

{{< release-date date="2024-11-27" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.19.0).

### New

- Bake now requires you to allow filesystem entitlements when your build needs
  to read or write files outside of your current working directory.
  [docker/buildx#2796](https://github.com/docker/buildx/pull/2796),
  [docker/buildx#2812](https://github.com/docker/buildx/pull/2812).

  To allow filesystem entitlements, use the `--allow fs.read=<path>` flag for
  the `docker buildx bake` command.

  This feature currently only reports a warning when using a local Bake
  definition, but will start to produce an error starting from the v0.20
  release. To enable the error in the current release, you can set
  `BUILDX_BAKE_ENTITLEMENTS_FS=1`.

### Enhancements

- Bake definition now supports new object notation for the fields that previously required CSV strings as inputs. [docker/buildx#2758](https://github.com/docker/buildx/pull/2758)

  > [!NOTE]
  > This enhancement was reverted in [v0.19.1](#0191) due to a bug.

- Bake definition now allows defining validation conditions to variables. [docker/buildx#2794](https://github.com/docker/buildx/pull/2794)
- Metadata file values can now contain JSON array values. [docker/buildx#2777](https://github.com/docker/buildx/pull/2777)
- Improved error messages when using an incorrect format for labels. [docker/buildx#2778](https://github.com/docker/buildx/pull/2778)
- FreeBSD and OpenBSD artifacts are now included in the release. [docker/buildx#2774](https://github.com/docker/buildx/pull/2774), [docker/buildx#2775](https://github.com/docker/buildx/pull/2775), [docker/buildx#2781](https://github.com/docker/buildx/pull/2781)

### Bug fixes

- Fixed an issue with printing Bake definitions containing empty Compose networks. [docker/buildx#2790](https://github.com/docker/buildx/pull/2790).

### Packaging

- Compose support has been updated to v2.4.4. [docker/buildx#2806](https://github.com/docker/buildx/pull/2806) [docker/buildx#2780](https://github.com/docker/buildx/pull/2780).

## 0.18.0

{{< release-date date="2024-10-31" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.18.0).

### New

- The `docker buildx inspect` command now displays BuildKit daemon configuration options set with a TOML file. [docker/buildx#2684](https://github.com/docker/buildx/pull/2684)
- The `docker buildx ls` command output is now more compact by default by compacting the platform list. A new `--no-trunc` option can be used for the full list. [docker/buildx#2138](https://github.com/docker/buildx/pull/2138), [docker/buildx#2717](https://github.com/docker/buildx/pull/2717)
- The `docker buildx prune` command now supports new `--max-used-space` and `--min-free-space` filters with BuildKit v0.17.0+ builders. [docker/buildx#2766](https://github.com/docker/buildx/pull/2766)

### Enhancements

- Allow capturing of CPU and memory profiles with `pprof` using the [`BUILDX_CPU_PROFILE`](/manuals/build/building/variables.md#buildx_cpu_profile) and [`BUILDX_MEM_PROFILE`](/manuals/build/building/variables.md#buildx_mem_profile) environment variables. [docker/buildx#2746](https://github.com/docker/buildx/pull/2746)
- Maximum Dockerfile size from standard input has increased. [docker/buildx#2716](https://github.com/docker/buildx/pull/2716), [docker/buildx#2719](https://github.com/docker/buildx/pull/2719)
- Memory allocations have been reduced. [docker/buildx#2724](https://github.com/docker/buildx/pull/2724), [docker/buildx#2713](https://github.com/docker/buildx/pull/2713)
- The `--list-targets` and `--list-variables` flags for `docker buildx bake` no longer require initialization of the builder. [docker/buildx#2763](https://github.com/docker/buildx/pull/2763)

### Bug fixes

- Check warnings now print the full filepath to the offending Dockerfile, relative to the current working directory. [docker/buildx#2672](https://github.com/docker/buildx/pull/2672)
- Fallback images for the `--check` and `--call` options have been updated to correct references. [docker/buildx#2705](https://github.com/docker/buildx/pull/2705)
- Fix issue with the build details link not showing in experimental mode. [docker/buildx#2722](https://github.com/docker/buildx/pull/2722)
- Fix validation issue with invalid target linking for Bake. [docker/buildx#2700](https://github.com/docker/buildx/pull/2700)
- Fix missing error message when running an invalid command. [docker/buildx#2741](https://github.com/docker/buildx/pull/2741)
- Fix possible false warnings for local state in `--call` requests. [docker/buildx#2754](https://github.com/docker/buildx/pull/2754)
- Fix potential issues with entitlements when using linked targets in Bake. [docker/buildx#2701](https://github.com/docker/buildx/pull/2701)
- Fix possible permission issues when accessing local state after running Buildx with `sudo`. [docker/buildx#2745](https://github.com/docker/buildx/pull/2745)

### Packaging

- Compose compatibility has been updated to v2.4.1. [docker/buildx#2760](https://github.com/docker/buildx/pull/2760)

## 0.17.1

{{< release-date date="2024-09-13" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.17.1).

### Bug fixes

- Do not set `network.host` entitlement flag automatically on builder creation
  for the `docker-container` and `kubernetes` drivers if the entitlement is set
  in the [BuildKit configuration file](/manuals/build/buildkit/toml-configuration.md). [docker/buildx#2685]
- Do not print the `network` field with `docker buildx bake --print` when empty. [docker/buildx#2689]
- Fix telemetry socket path under WSL2. [docker/buildx#2698]

[docker/buildx#2685]: https://github.com/docker/buildx/pull/2685
[docker/buildx#2689]: https://github.com/docker/buildx/pull/2689
[docker/buildx#2698]: https://github.com/docker/buildx/pull/2698

## 0.17.0

{{< release-date date="2024-09-10" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.17.0).

### New

- Add `basename`, `dirname` and `sanitize` functions to Bake. [docker/buildx#2649]
- Enable support for Bake entitlements to allow privileged operations during builds. [docker/buildx#2666]

### Enhancements

- Introduce CLI metrics tracking for Bake commands. [docker/buildx#2610]
- Add `--debug` to all build commands. Previously, it was only available on the top-level `docker` and `docker buildx` commands. [docker/buildx#2660]
- Allow builds from stdin for multi-node builders. [docker/buildx#2656]
- Improve `kubernetes` driver initialization. [docker/buildx#2606]
- Include target name in the error message when building multiple targets with Bake. [docker/buildx#2651]
- Optimize metrics handling to reduce performance overhead during progress tracking. [docker/buildx#2641]
- Display the number of warnings after completing a rule check. [docker/buildx#2647]
- Skip build ref and provenance metadata for frontend methods. [docker/buildx#2650]
- Add support for setting network mode in Bake files (HCL and JSON). [docker/buildx#2671]
- Support the `--metadata-file` flag when set along the `--call` flag. [docker/buildx#2640]
- Use shared session for local contexts used by multiple Bake targets. [docker/buildx#2615], [docker/buildx#2607], [docker/buildx#2663]

### Bug fixes

- Improve memory management to avoid unnecessary allocations. [docker/buildx#2601]

### Packaging updates

- Compose support has been updated to v2.1.6. [docker/buildx#2547]

[docker/buildx#2547]: https://github.com/docker/buildx/pull/2547/
[docker/buildx#2601]: https://github.com/docker/buildx/pull/2601/
[docker/buildx#2606]: https://github.com/docker/buildx/pull/2606/
[docker/buildx#2607]: https://github.com/docker/buildx/pull/2607/
[docker/buildx#2610]: https://github.com/docker/buildx/pull/2610/
[docker/buildx#2615]: https://github.com/docker/buildx/pull/2615/
[docker/buildx#2640]: https://github.com/docker/buildx/pull/2640/
[docker/buildx#2641]: https://github.com/docker/buildx/pull/2641/
[docker/buildx#2647]: https://github.com/docker/buildx/pull/2647/
[docker/buildx#2649]: https://github.com/docker/buildx/pull/2649/
[docker/buildx#2650]: https://github.com/docker/buildx/pull/2650/
[docker/buildx#2651]: https://github.com/docker/buildx/pull/2651/
[docker/buildx#2656]: https://github.com/docker/buildx/pull/2656/
[docker/buildx#2660]: https://github.com/docker/buildx/pull/2660/
[docker/buildx#2663]: https://github.com/docker/buildx/pull/2663/
[docker/buildx#2666]: https://github.com/docker/buildx/pull/2666/
[docker/buildx#2671]: https://github.com/docker/buildx/pull/2671/

## 0.16.2

{{< release-date date="2024-07-25" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.16.2).

### Bug fixes

- Fix possible "bad file descriptor" error when exporting local cache to NFS volume [docker/buildx#2629](https://github.com/docker/buildx/pull/2629/)

## 0.16.1

{{< release-date date="2024-07-18" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.16.1).

### Bug fixes

- Fix possible panic due to data race in `buildx bake --print` command [docker/buildx#2603](https://github.com/docker/buildx/pull/2603/)
- Improve messaging about using `--debug` flag to inspect build warnings [docker/buildx#2612](https://github.com/docker/buildx/pull/2612/)

## 0.16.0

{{< release-date date="2024-07-11" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.16.0).

### New

- Bake command now supports `--call` and `--check` flags and `call` attribute in target definitions for selecting custom frontend methods. [docker/buildx#2556](https://github.com/docker/buildx/pull/2556/), [docker/buildx#2576](https://github.com/docker/buildx/pull/2576/)
- {{< badge color=violet text=Experimental >}} Bake now supports `--list-targets` and `--list-variables` flags for inspecting the definition and possible configuration options for your project. [docker/buildx#2556](https://github.com/docker/buildx/pull/2556/)
- Bake definition variables and targets supports new `description` attribute for defining text-based description that can be inspected using e.g. `--list-targets` and `--list-variables`. [docker/buildx#2556](https://github.com/docker/buildx/pull/2556/)
- Bake now supports printing warnings for build check violations. [docker/buildx#2501](https://github.com/docker/buildx/pull/2501/)

### Enhancements

- The build command now ensures that multi-node builds use the same build reference for each node. [docker/buildx#2572](https://github.com/docker/buildx/pull/2572/)
- Avoid duplicate requests and improve the performance of remote driver. [docker/buildx#2501](https://github.com/docker/buildx/pull/2501/)
- Build warnings can now be saved to the metadata file by setting the `BUILDX_METADATA_WARNINGS=1` environment variable. [docker/buildx#2551](https://github.com/docker/buildx/pull/2551/), [docker/buildx#2521](https://github.com/docker/buildx/pull/2521/), [docker/buildx#2550](https://github.com/docker/buildx/pull/2550/)
- Improve message of the `--check` flag when no warnings are detected. [docker/buildx#2549](https://github.com/docker/buildx/pull/2549/)

### Bug fixes

- Fix support for multi-type annotations during build. [docker/buildx#2522](https://github.com/docker/buildx/pull/2522/)
- Fix a regression where possible inefficient transfer of files would occur when switching projects due to incremental transfer reuse. [docker/buildx#2558](https://github.com/docker/buildx/pull/2558/)
- Fix incorrect default load for chained Bake targets. [docker/buildx#2583](https://github.com/docker/buildx/pull/2583/)
- Fix incorrect `COMPOSE_PROJECT_NAME` handling in Bake. [docker/buildx#2579](https://github.com/docker/buildx/pull/2579/)
- Fix index annotations support for multi-node builds. [docker/buildx#2546](https://github.com/docker/buildx/pull/2546/)
- Fix capturing provenance metadata for builds from remote context. [docker/buildx#2560](https://github.com/docker/buildx/pull/2560/)

### Packaging updates

- Compose support has been updated to v2.1.3. [docker/buildx#2547](https://github.com/docker/buildx/pull/2547/)

## 0.15.1

{{< release-date date="2024-06-18" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.15.1).

### Bug fixes

- Fix missing build error and exit code for some validation requests with `--check`. [docker/buildx#2518](https://github.com/docker/buildx/pull/2518/)
- Update fallback image for `--check` to Dockerfile v1.8.1. [docker/buildx#2538](https://github.com/docker/buildx/pull/2538/)

## 0.15.0

{{< release-date date="2024-06-11" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.15.0).

### New

- New `--call` option allows setting evaluation method for a build, replacing the previous experimental `--print` flag. [docker/buildx#2498](https://github.com/docker/buildx/pull/2498/), [docker/buildx#2487](https://github.com/docker/buildx/pull/2487/), [docker/buildx#2513](https://github.com/docker/buildx/pull/2513/)

  In addition to the default `build` method, the following methods are implemented by Dockerfile frontend:

  - [`--call=check`](/reference/cli/docker/buildx/build.md#check): Run validation routines for your build configuration. For more information about build checks, see [Build checks](/manuals/build/checks.md)
  - [`--call=outline`](/reference/cli/docker/buildx/build.md#call-outline): Show configuration that would be used by current build, including all build arguments, secrets, SSH mounts, etc., that your build would use.
  - [`--call=targets`](/reference/cli/docker/buildx/build.md#call-targets): Show all available targets and their descriptions.

- New `--prefer-index` flag has been added to the `docker buildx imagetools create` command to control the behavior of creating image out of one single-platform image manifest. [docker/buildx#2482](https://github.com/docker/buildx/pull/2482/)
- The [`kubernetes` driver](/manuals/build/builders/drivers/kubernetes.md) now supports a `timeout` option for configuring deployment timeout. [docker/buildx#2492](https://github.com/docker/buildx/pull/2492/)
- New metrics definitions have been added for build warning types. [docker/buildx#2482](https://github.com/docker/buildx/pull/2482/), [docker/buildx#2507](https://github.com/docker/buildx/pull/2507/)
- The [`buildx prune`](/reference/cli/docker/buildx/prune.md) and [`buildx du`](/reference/cli/docker/buildx/du.md) commands now support negative and prefix filters. [docker/buildx#2473](https://github.com/docker/buildx/pull/2473/)
- Building Compose files with Bake now supports passing SSH forwarding configuration. [docker/buildx#2445](https://github.com/docker/buildx/pull/2445/)
- Fix issue with configuring the `kubernetes` driver with custom TLS certificates. [docker/buildx#2454](https://github.com/docker/buildx/pull/2454/)
- Fix concurrent kubeconfig access when loading nodes. [docker/buildx#2497](https://github.com/docker/buildx/pull/2497/)

### Packaging updates

- Compose support has been updated to v2.1.2. [docker/buildx#2502](https://github.com/docker/buildx/pull/2502/), [docker/buildx#2425](https://github.com/docker/buildx/pull/2425/)

## 0.14.0

{{< release-date date="2024-04-18" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.14.0).

### Enhancements

- Add support for `--print=lint` (experimental).
  [docker/buildx#2404](https://github.com/docker/buildx/pull/2404),
  [docker/buildx#2406](https://github.com/docker/buildx/pull/2406)
- Fix JSON formatting for custom implementations of print sub-requests in frontends.
  [docker/buildx#2374](https://github.com/docker/buildx/pull/2374)
- Provenance records are now set when building with `--metadata-file`.
  [docker/buildx#2280](https://github.com/docker/buildx/pull/2280)
- Add [Git authentication support](./bake/remote-definition.md#remote-definition-in-a-private-repository) for remote definitions.
  [docker/buildx#2363](https://github.com/docker/buildx/pull/2363)
- New `default-load` driver option for the `docker-container`, `remote`, and `kubernetes` drivers to load build results to the Docker Engine image store by default.
  [docker/buildx#2259](https://github.com/docker/buildx/pull/2259)
- Add `requests.ephemeral-storage`, `limits.ephemeral-storage` and `schedulername` options to the [`kubernetes` driver](/manuals/build/builders/drivers/kubernetes.md).
  [docker/buildx#2370](https://github.com/docker/buildx/pull/2370),
  [docker/buildx#2415](https://github.com/docker/buildx/pull/2415)
- Add `indexof` function for `docker-bake.hcl` files.
  [docker/buildx#2384](https://github.com/docker/buildx/pull/2384)
- OpenTelemetry metrics for Buildx now measure durations of idle time, image exports, run operations, and image transfers for image source operations during build.
  [docker/buildx#2316](https://github.com/docker/buildx/pull/2316),
  [docker/buildx#2317](https://github.com/docker/buildx/pull/2317),
  [docker/buildx#2323](https://github.com/docker/buildx/pull/2323),
  [docker/buildx#2271](https://github.com/docker/buildx/pull/2271)
- Build progress metrics to the OpenTelemetry endpoint associated with the `desktop-linux` context no longer requires Buildx in experimental mode (`BUILDX_EXPERIMENTAL=1`).
  [docker/buildx#2344](https://github.com/docker/buildx/pull/2344)

### Bug fixes

- Fix `--load` and `--push` incorrectly overriding outputs when used with multiple Bake file definitions.
  [docker/buildx#2336](https://github.com/docker/buildx/pull/2336)
- Fix build from stdin with experimental mode enabled.
  [docker/buildx#2394](https://github.com/docker/buildx/pull/2394)
- Fix an issue where delegated traces could be duplicated.
  [docker/buildx#2362](https://github.com/docker/buildx/pull/2362)

### Packaging updates

- Compose support has been updated to [v2.26.1](https://github.com/docker/compose/releases/tag/v2.26.1)
  (via [`compose-go` v2.0.2](https://github.com/compose-spec/compose-go/releases/tag/v2.0.2)).
  [docker/buildx#2391](https://github.com/docker/buildx/pull/2391)

## 0.13.1

{{< release-date date="2024-03-13" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.13.1).

### Bug fixes

- Fix connecting to `docker-container://` and `kube-pod://` style URLs with remote driver. [docker/buildx#2327](https://github.com/docker/buildx/pull/2327)
- Fix handling of `--push` with Bake when a target has already defined a non-image output. [docker/buildx#2330](https://github.com/docker/buildx/pull/2330)

## 0.13.0

{{< release-date date="2024-03-06" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.13.0).

### New

- New `docker buildx dial-stdio` command for directly contacting BuildKit daemon of the configured builder instance. [docker/buildx#2112](https://github.com/docker/buildx/pull/2112)
- Windows container builders can now be created using the `remote` driver and npipe connections. [docker/buildx#2287](https://github.com/docker/buildx/pull/2287)
- Npipe URL scheme is now supported on Windows. [docker/buildx#2250](https://github.com/docker/buildx/pull/2250)
- {{< badge color=violet text=Experimental >}} Buildx can now export OpenTelemetry metrics for build duration and transfer sizes. [docker/buildx#2235](https://github.com/docker/buildx/pull/2235), [docker/buildx#2258](https://github.com/docker/buildx/pull/2258) [docker/buildx#2225](https://github.com/docker/buildx/pull/2225) [docker/buildx#2224](https://github.com/docker/buildx/pull/2224) [docker/buildx#2155](https://github.com/docker/buildx/pull/2155)

### Enhancements

- Bake command now supports defining `shm-size` and `ulimit` values. [docker/buildx#2279](https://github.com/docker/buildx/pull/2279), [docker/buildx#2242](https://github.com/docker/buildx/pull/2242)
- Better handling of connecting to unhealthy nodes with remote driver. [docker/buildx#2130](https://github.com/docker/buildx/pull/2130)
- Builders using the `docker-container` and `kubernetes` drivers now allow `network.host` entitlement by default (allowing access to the container's network). [docker/buildx#2266](https://github.com/docker/buildx/pull/2266)
- Builds can now use multiple outputs with a single command (requires BuildKit v0.13+). [docker/buildx#2290](https://github.com/docker/buildx/pull/2290), [docker/buildx#2302](https://github.com/docker/buildx/pull/2302)
- Default Git repository path is now found via configured tracking branch. [docker/buildx#2146](https://github.com/docker/buildx/pull/2146)
- Fix possible cache invalidation when using linked targets in Bake. [docker/buildx#2265](https://github.com/docker/buildx/pull/2265)
- Fixes for Git repository path sanitization in WSL. [docker/buildx#2167](https://github.com/docker/buildx/pull/2167)
- Multiple builders can now be removed with a single command. [docker/buildx#2140](https://github.com/docker/buildx/pull/2140)
- New cancellation signal handling via Unix socket. [docker/buildx#2184](https://github.com/docker/buildx/pull/2184) [docker/buildx#2289](https://github.com/docker/buildx/pull/2289)
- The Compose spec support has been updated to v2.0.0-rc.8. [docker/buildx#2205](https://github.com/docker/buildx/pull/2205)
- The `--config` flag for `docker buildx create` was renamed to `--buildkitd-config`. [docker/buildx#2268](https://github.com/docker/buildx/pull/2268)
- The `--metadata-file` flag for `docker buildx build` can now also return build reference that can be used for further build debugging, for example, in Docker Desktop. [docker/buildx#2263](https://github.com/docker/buildx/pull/2263)
- The `docker buildx bake` command now shares the same authentication provider for all targets for improved performance. [docker/buildx#2147](https://github.com/docker/buildx/pull/2147)
- The `docker buildx imagetools inspect` command now shows DSSE-signed SBOM and Provenance attestations. [docker/buildx#2194](https://github.com/docker/buildx/pull/2194)
- The `docker buildx ls` command now supports `--format` options for controlling the output. [docker/buildx#1787](https://github.com/docker/buildx/pull/1787)
- The `docker-container` driver now supports driver options for defining restart policy for BuildKit container. [docker/buildx#1271](https://github.com/docker/buildx/pull/1271)
- VCS attributes exported from Buildx now include the local directory sub-paths if they're relative to the current Git repository. [docker/buildx#2156](https://github.com/docker/buildx/pull/2156)
- `--add-host` flag now permits a `=` separator for IPv6 addresses. [docker/buildx#2121](https://github.com/docker/buildx/pull/2121)

### Bug fixes

- Fix additional output when exporting progress with `--progress=rawjson` [docker/buildx#2252](https://github.com/docker/buildx/pull/2252)
- Fix possible console warnings on Windows. [docker/buildx#2238](https://github.com/docker/buildx/pull/2238)
- Fix possible inconsistent configuration merge order when using Bake with many configurations. [docker/buildx#2237](https://github.com/docker/buildx/pull/2237)
- Fix possible panic in the `docker buildx imagetools create` command. [docker/buildx#2230](https://github.com/docker/buildx/pull/2230)

## 0.12.1

{{< release-date date="2024-01-12" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.12.1).

### Bug fixes and enhancements

- Fix incorrect validation of some `--driver-opt` values that could cause a panic and corrupt state to be stored.
  [docker/buildx#2176](https://github.com/docker/buildx/pull/2176)

## 0.12.0

{{< release-date date="2023-11-16" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.12.0).

### New

- New `--annotation` flag for the `buildx build`, and an `annotations` key in the Bake file, that lets you add OCI Annotations to build results.
  [#2020](https://github.com/docker/buildx/pull/2020),
  [#2098](https://github.com/docker/buildx/pull/2098)
- New experimental debugging features, including a new `debug` command and an interactive debugging console.
  This feature currently requires setting `BUILDX_EXPERIMENTAL=1`.
  [#2006](https://github.com/docker/buildx/pull/2006),
  [#1896](https://github.com/docker/buildx/pull/1896),
  [#1970](https://github.com/docker/buildx/pull/1970),
  [#1914](https://github.com/docker/buildx/pull/1914),
  [#2026](https://github.com/docker/buildx/pull/2026),
  [#2086](https://github.com/docker/buildx/pull/2086)

### Bug fixes and enhancements

- The special `host-gateway` IP mapping can now be used with the `--add-host` flag during build.
  [#1894](https://github.com/docker/buildx/pull/1894),
  [#2083](https://github.com/docker/buildx/pull/2083)
- Bake now allows adding local source files when building from remote definition.
  [#1838](https://github.com/docker/buildx/pull/1838)
- The status of uploading build results to Docker is now shown interactively on progress bar.
  [#1994](https://github.com/docker/buildx/pull/1994)
- Error handling has been improved when bootstrapping multi-node build clusters.
  [#1869](https://github.com/docker/buildx/pull/1869)
- The `buildx imagetools create` command now allows adding annotation when creating new images in the registry.
  [#1965](https://github.com/docker/buildx/pull/1965)
- OpenTelemetry build trace delegation from buildx is now possible with Docker and Remote driver.
  [#2034](https://github.com/docker/buildx/pull/2034)
- Bake command now shows all files where the build definition was loaded from on the progress bar.
  [#2076](https://github.com/docker/buildx/pull/2076)
- Bake files now allow the same attributes to be defined in multiple definition files.
  [#1062](https://github.com/docker/buildx/pull/1062)
- Using the Bake command with a remote definition now allows this definition to use local Dockerfiles.
  [#2015](https://github.com/docker/buildx/pull/2015)
- Docker container driver now explicitly sets BuildKit config path to make sure configurations are loaded from same location for both mainline and rootless images.
  [#2093](https://github.com/docker/buildx/pull/2093)
- Improve performance of detecting when BuildKit instance has completed booting.
  [#1934](https://github.com/docker/buildx/pull/1934)
- Container driver now accepts many new driver options for defining the resource limits for BuildKit container.
  [#2048](https://github.com/docker/buildx/pull/2048)
- Inspection commands formatting has been improved.
  [#2068](https://github.com/docker/buildx/pull/2068)
- Error messages about driver capabilities have been improved.
  [#1998](https://github.com/docker/buildx/pull/1998)
- Improve errors when invoking Bake command without targets.
  [#2100](https://github.com/docker/buildx/pull/2100)
- Allow enabling debug logs with environment variables when running in standalone mode.
  [#1821](https://github.com/docker/buildx/pull/1821)
- When using Docker driver the default image resolve mode has been updated to prefer local Docker images for backward compatibility.
  [#1886](https://github.com/docker/buildx/pull/1886)
- Kubernetes driver now allows setting custom annotations and labels to the BuildKit deployments and pods.
  [#1938](https://github.com/docker/buildx/pull/1938)
- Kubernetes driver now allows setting authentication token with endpoint configuration.
  [#1891](https://github.com/docker/buildx/pull/1891)
- Fix possible issue with chained targets in Bake that could result in build failing or local source for a target uploaded multiple times.
  [#2113](https://github.com/docker/buildx/pull/2113)
- Fix issue when accessing global target properties when using the matrix feature of the Bake command.
  [#2106](https://github.com/docker/buildx/pull/2106)
- Fixes for formatting validation of certain build flags
  [#2040](https://github.com/docker/buildx/pull/2040)
- Fixes to avoid locking certain commands unnecessarily while booting builder nodes.
  [#2066](https://github.com/docker/buildx/pull/2066)
- Fix cases where multiple builds try to bootstrap the same builder instance in parallel.
  [#2000](https://github.com/docker/buildx/pull/2000)
- Fix cases where errors on uploading build results to Docker could be dropped in some cases.
  [#1927](https://github.com/docker/buildx/pull/1927)
- Fix detecting capabilities for missing attestation support based on build output.
  [#1988](https://github.com/docker/buildx/pull/1988)
- Fix the build for loading in Bake remote definition to not show up in build history records.
  [#1961](https://github.com/docker/buildx/pull/1961),
  [#1954](https://github.com/docker/buildx/pull/1954)
- Fix errors when building Compose files using the that define profiles with Bake.
  [#1903](https://github.com/docker/buildx/pull/1903)
- Fix possible time correction errors on progress bar.
  [#1968](https://github.com/docker/buildx/pull/1968)
- Fix passing custom cgroup parent to builds that used the new controller interface.
  [#1913](https://github.com/docker/buildx/pull/1913)

### Packaging

- Compose support has been updated to 1.20, enabling "include" functionality when using the Bake command.
  [#1971](https://github.com/docker/buildx/pull/1971),
  [#2065](https://github.com/docker/buildx/pull/2065),
  [#2094](https://github.com/docker/buildx/pull/2094)

## 0.11.2

{{< release-date date="2023-07-18" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.11.2).

### Bug fixes and enhancements

- Fix a regression that caused buildx to not read the `KUBECONFIG` path from the instance store.
  [docker/buildx#1941](https://github.com/docker/buildx/pull/1941)
- Fix a regression with result handle builds showing up in the build history incorrectly.
  [docker/buildx#1954](https://github.com/docker/buildx/pull/1954)

## 0.11.1

{{< release-date date="2023-07-05" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.11.1).

### Bug fixes and enhancements

- Fix a regression for bake where services in profiles would not be loaded.
  [docker/buildx#1903](https://github.com/docker/buildx/pull/1903)
- Fix a regression where `--cgroup-parent` option had no effect during build.
  [docker/buildx#1913](https://github.com/docker/buildx/pull/1913)
- Fix a regression where valid docker contexts could fail buildx builder name
  validation. [docker/buildx#1879](https://github.com/docker/buildx/pull/1879)
- Fix a possible panic when terminal is resized during the build.
  [docker/buildx#1929](https://github.com/docker/buildx/pull/1929)

## 0.11.0

{{< release-date date="2023-06-13" >}}

The full release notes for this release are available
[on GitHub](https://github.com/docker/buildx/releases/tag/v0.11.0).

### New

- Bake now supports [matrix builds](/manuals/build/bake/reference.md#targetmatrix).
  The new matrix field on `target` lets you create multiple similar targets to
  remove duplication in bake files. [docker/buildx#1690](https://github.com/docker/buildx/pull/1690)
- New experimental `--detach` flag for running builds in detached mode.
  [docker/buildx#1296](https://github.com/docker/buildx/pull/1296),
  [docker/buildx#1620](https://github.com/docker/buildx/pull/1620),
  [docker/buildx#1614](https://github.com/docker/buildx/pull/1614),
  [docker/buildx#1737](https://github.com/docker/buildx/pull/1737),
  [docker/buildx#1755](https://github.com/docker/buildx/pull/1755)
- New experimental [debug monitor mode](https://github.com/docker/buildx/blob/v0.11.0-rc1/docs/guides/debugging.md)
  that lets you start a debug session in your builds.
  [docker/buildx#1626](https://github.com/docker/buildx/pull/1626),
  [docker/buildx#1640](https://github.com/docker/buildx/pull/1640)
- New [`EXPERIMENTAL_BUILDKIT_SOURCE_POLICY` environment variable](./building/variables.md#experimental_buildkit_source_policy)
  for applying a BuildKit source policy file.
  [docker/buildx#1628](https://github.com/docker/buildx/pull/1628)

### Bug fixes and enhancements

- `--load` now supports loading multi-platform images when the containerd image
  store is enabled.
  [docker/buildx#1813](https://github.com/docker/buildx/pull/1813)
- Build progress output now displays the name of the builder being used.
  [docker/buildx#1177](https://github.com/docker/buildx/pull/1177)
- Bake now supports detecting `compose.{yml,yaml}` files.
  [docker/buildx#1752](https://github.com/docker/buildx/pull/1752)
- Bake now supports new compose build keys `dockerfile_inline` and `additional_contexts`.
  [docker/buildx#1784](https://github.com/docker/buildx/pull/1784)
- Bake now supports replace HCL function.
  [docker/buildx#1720](https://github.com/docker/buildx/pull/1720)
- Bake now allows merging multiple similar attestation parameters into a single
  parameter to allow overriding with a single global value.
  [docker/buildx#1699](https://github.com/docker/buildx/pull/1699)
- Initial support for shell completion.
  [docker/buildx#1727](https://github.com/docker/buildx/pull/1727)
- BuildKit versions now correctly display in `buildx ls` and `buildx inspect`
  for builders using the `docker` driver.
  [docker/buildx#1552](https://github.com/docker/buildx/pull/1552)
- Display additional builder node details in buildx inspect view.
  [docker/buildx#1440](https://github.com/docker/buildx/pull/1440),
  [docker/buildx#1854](https://github.com/docker/buildx/pull/1874)
- Builders using the `remote` driver allow using TLS without proving its own
  key/cert (if BuildKit remote is configured to support it)
  [docker/buildx#1693](https://github.com/docker/buildx/pull/1693)
- Builders using the `kubernetes` driver support a new `serviceaccount` option,
  which sets the `serviceAccountName` of the Kubernetes pod.
  [docker/buildx#1597](https://github.com/docker/buildx/pull/1597)
- Builders using the `kubernetes` driver support the `proxy-url` option in the
  kubeconfig file.
  [docker/buildx#1780](https://github.com/docker/buildx/pull/1780)
- Builders using the `kubernetes` are now automatically assigned a node name if
  no name is explicitly provided.
  [docker/buildx#1673](https://github.com/docker/buildx/pull/1673)
- Fix invalid path when writing certificates for `docker-container` driver on Windows.
  [docker/buildx#1831](https://github.com/docker/buildx/pull/1831)
- Fix bake failure when remote bake file is accessed using SSH.
  [docker/buildx#1711](https://github.com/docker/buildx/pull/1711),
  [docker/buildx#1734](https://github.com/docker/buildx/pull/1734)
- Fix bake failure when remote bake context is incorrectly resolved.
  [docker/buildx#1783](https://github.com/docker/buildx/pull/1783)
- Fix path resolution of `BAKE_CMD_CONTEXT` and `cwd://` paths in bake contexts.
  [docker/buildx#1840](https://github.com/docker/buildx/pull/1840)
- Fix mixed OCI and Docker media types when creating images using
  `buildx imagetools create`.
  [docker/buildx#1797](https://github.com/docker/buildx/pull/1797)
- Fix mismatched image id between `--iidfile` and `-q`.
  [docker/buildx#1844](https://github.com/docker/buildx/pull/1844)
- Fix AWS authentication when mixing static creds and IAM profiles.
  [docker/buildx#1816](https://github.com/docker/buildx/pull/1816)

## 0.10.4

{{< release-date date="2023-03-06" >}}

{{% include "buildx-v0.10-disclaimer.md" %}}

### Bug fixes and enhancements

- Add `BUILDX_NO_DEFAULT_ATTESTATIONS` as alternative to `--provenance false`. [docker/buildx#1645](https://github.com/docker/buildx/issues/1645)
- Disable dirty Git checkout detection by default for performance. Can be enabled with `BUILDX_GIT_CHECK_DIRTY` opt-in. [docker/buildx#1650](https://github.com/docker/buildx/issues/1650)
- Strip credentials from VCS hint URL before sending to BuildKit. [docker/buildx#1664](https://github.com/docker/buildx/issues/1664)

## 0.10.3

{{< release-date date="2023-02-16" >}}

{{% include "buildx-v0.10-disclaimer.md" %}}

### Bug fixes and enhancements

- Fix reachable commit and warnings on collecting Git provenance info. [docker/buildx#1592](https://github.com/docker/buildx/issues/1592), [docker/buildx#1634](https://github.com/docker/buildx/issues/1634)
- Fix a regression where docker context was not being validated. [docker/buildx#1596](https://github.com/docker/buildx/issues/1596)
- Fix function resolution with JSON bake definition. [docker/buildx#1605](https://github.com/docker/buildx/issues/1605)
- Fix case where original HCL bake diagnostic is discarded. [docker/buildx#1607](https://github.com/docker/buildx/issues/1607)
- Fix labels not correctly set with bake and compose file. [docker/buildx#1631](https://github.com/docker/buildx/issues/1631)

## 0.10.2

{{< release-date date="2023-01-30" >}}

{{% include "buildx-v0.10-disclaimer.md" %}}

### Bug fixes and enhancements

- Fix preferred platforms order not taken into account in multi-node builds. [docker/buildx#1561](https://github.com/docker/buildx/issues/1561)
- Fix possible panic on handling `SOURCE_DATE_EPOCH` environment variable. [docker/buildx#1564](https://github.com/docker/buildx/issues/1564)
- Fix possible push error on multi-node manifest merge since BuildKit v0.11 on
  some registries. [docker/buildx#1566](https://github.com/docker/buildx/issues/1566)
- Improve warnings on collecting Git provenance info. [docker/buildx#1568](https://github.com/docker/buildx/issues/1568)

## 0.10.1

{{< release-date date="2023-01-27" >}}

{{% include "buildx-v0.10-disclaimer.md" %}}

### Bug fixes and enhancements

- Fix sending the correct origin URL as `vsc:source` metadata. [docker/buildx#1548](https://github.com/docker/buildx/issues/1548)
- Fix possible panic from data-race. [docker/buildx#1504](https://github.com/docker/buildx/issues/1504)
- Fix regression with `rm --all-inactive`. [docker/buildx#1547](https://github.com/docker/buildx/issues/1547)
- Improve attestation access in `imagetools inspect` by lazily loading data. [docker/buildx#1546](https://github.com/docker/buildx/issues/1546)
- Correctly mark capabilities request as internal. [docker/buildx#1538](https://github.com/docker/buildx/issues/1538)
- Detect invalid attestation configuration. [docker/buildx#1545](https://github.com/docker/buildx/issues/1545)
- Update containerd patches to fix possible push regression affecting
  `imagetools` commands. [docker/buildx#1559](https://github.com/docker/buildx/issues/1559)

## 0.10.0

{{< release-date date="2023-01-10" >}}

{{% include "buildx-v0.10-disclaimer.md" %}}

### New

- The `buildx build` command supports new `--attest` flag, along with
  shorthands `--sbom` and `--provenance`, for adding attestations for your
  current build. [docker/buildx#1412](https://github.com/docker/buildx/issues/1412)
  [docker/buildx#1475](https://github.com/docker/buildx/issues/1475)
  - `--attest type=sbom` or `--sbom=true` adds [SBOM attestations](/manuals/build/metadata/attestations/sbom.md).
  - `--attest type=provenance` or `--provenance=true` adds [SLSA provenance attestation](/manuals/build/metadata/attestations/slsa-provenance.md).
  - When creating OCI images, a minimal provenance attestation is included
    with the image by default.
- When building with BuildKit that supports provenance attestations Buildx will
  automatically share the version control information of your build context, so
  it can be shown in provenance for later debugging. Previously this only
  happened when building from a Git URL directly. To opt-out of this behavior
  you can set `BUILDX_GIT_INFO=0`. Optionally you can also automatically define
  labels with VCS info by setting `BUILDX_GIT_LABELS=1`.
  [docker/buildx#1462](https://github.com/docker/buildx/issues/1462),
  [docker/buildx#1297](https://github.com/docker/buildx),
  [docker/buildx#1341](https://github.com/docker/buildx/issues/1341),
  [docker/buildx#1468](https://github.com/docker/buildx),
  [docker/buildx#1477](https://github.com/docker/buildx/issues/1477)
- Named contexts with `--build-context` now support `oci-layout://` protocol
  for initializing the context with a value of a local OCI layout directory.
  E.g. `--build-context stagename=oci-layout://path/to/dir`. This feature
  requires BuildKit v0.11.0+ and Dockerfile 1.5.0+. [docker/buildx#1456](https://github.com/docker/buildx/issues/1456)
- Bake now supports [resource interpolation](bake/inheritance.md#reusing-single-attribute-from-targets)
  where you can reuse the values from other target definitions. [docker/buildx#1434](https://github.com/docker/buildx/issues/1434)
- Buildx will now automatically forward `SOURCE_DATE_EPOCH` environment variable
  if it is defined in your environment. This feature is meant to be used with
  updated [reproducible builds](https://github.com/moby/buildkit/blob/master/docs/build-repro.md)
  support in BuildKit v0.11.0+. [docker/buildx#1482](https://github.com/docker/buildx/issues/1482)
- Buildx now remembers the last activity for a builder for better organization
  of builder instances. [docker/buildx#1439](https://github.com/docker/buildx/issues/1439)
- Bake definition now supports null values for [variables](bake/reference.md#variable) and [labels](bake/reference.md#targetlabels)
  for build arguments and labels to use the defaults set in the Dockerfile.
  [docker/buildx#1449](https://github.com/docker/buildx/issues/1449)
- The [`buildx imagetools inspect` command](/reference/cli/docker/buildx/imagetools/inspect.md)
  now supports showing SBOM and Provenance data.
  [docker/buildx#1444](https://github.com/docker/buildx/issues/1444),
  [docker/buildx#1498](https://github.com/docker/buildx/issues/1498)
- Increase performance of `ls` command and inspect flows.
  [docker/buildx#1430](https://github.com/docker/buildx/issues/1430),
  [docker/buildx#1454](https://github.com/docker/buildx/issues/1454),
  [docker/buildx#1455](https://github.com/docker/buildx/issues/1455),
  [docker/buildx#1345](https://github.com/docker/buildx/issues/1345)
- Adding extra hosts with [Docker driver](/manuals/build/builders/drivers/docker.md) now supports
  Docker-specific `host-gateway` special value. [docker/buildx#1446](https://github.com/docker/buildx/issues/1446)
- [OCI exporter](exporters/oci-docker.md) now supports `tar=false` option for
  exporting OCI format directly in a directory. [docker/buildx#1420](https://github.com/docker/buildx/issues/1420)

### Upgrades

- Updated the Compose Specification to 1.6.0. [docker/buildx#1387](https://github.com/docker/buildx/issues/1387)

### Bug fixes and enhancements

- `--invoke` can now load default launch environment from the image metadata. [docker/buildx#1324](https://github.com/docker/buildx/issues/1324)
- Fix container driver behavior in regards to UserNS. [docker/buildx#1368](https://github.com/docker/buildx/issues/1368)
- Fix possible panic in Bake when using wrong variable value type. [docker/buildx#1442](https://github.com/docker/buildx/issues/1442)
- Fix possible panic in `imagetools inspect`. [docker/buildx#1441](https://github.com/docker/buildx/issues/1441)
  [docker/buildx#1406](https://github.com/docker/buildx/issues/1406)
- Fix sending empty `--add-host` value to BuildKit by default. [docker/buildx#1457](https://github.com/docker/buildx/issues/1457)
- Fix handling progress prefixes with progress groups. [docker/buildx#1305](https://github.com/docker/buildx/issues/1305)
- Fix recursively resolving groups in Bake. [docker/buildx#1313](https://github.com/docker/buildx/issues/1313)
- Fix possible wrong indentation on multi-node builder manifests. [docker/buildx#1396](https://github.com/docker/buildx/issues/1396)
- Fix possible panic from missing OpenTelemetry configuration. [docker/buildx#1383](https://github.com/docker/buildx/issues/1383)
- Fix `--progress=tty` behavior when TTY is not available. [docker/buildx#1371](https://github.com/docker/buildx/issues/1371)
- Fix connection error conditions in `prune` and `du` commands. [docker/buildx#1307](https://github.com/docker/buildx/issues/1307)

## 0.9.1

{{< release-date date="2022-08-18" >}}

### Bug fixes and enhancements

- The `inspect` command now displays the BuildKit version in use. [docker/buildx#1279](https://github.com/docker/buildx/issues/1279)
- Fixed a regression when building Compose files that contain services without a
  build block. [docker/buildx#1277](https://github.com/docker/buildx/issues/1277)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.9.1).

## 0.9.0

{{< release-date date="2022-08-17" >}}

### New

- Support for a new [`remote` driver](/manuals/build/builders/drivers/remote.md) that you can use
  to connect to any already running BuildKit instance.
  [docker/buildx#1078](https://github.com/docker/buildx/issues/1078),
  [docker/buildx#1093](https://github.com/docker/buildx/issues/1093),
  [docker/buildx#1094](https://github.com/docker/buildx/issues/1094),
  [docker/buildx#1103](https://github.com/docker/buildx/issues/1103),
  [docker/buildx#1134](https://github.com/docker/buildx/issues/1134),
  [docker/buildx#1204](https://github.com/docker/buildx/issues/1204)
- You can now load Dockerfile from standard input even when the build context is
  coming from external Git or HTTP URL. [docker/buildx#994](https://github.com/docker/buildx/issues/994)
- Build commands now support new the build context type `oci-layout://` for loading
  [build context from local OCI layout directories](/reference/cli/docker/buildx/build.md#source-oci-layout).
  Note that this feature depends on an unreleased BuildKit feature and builder
  instance from `moby/buildkit:master` needs to be used until BuildKit v0.11 is
  released. [docker/buildx#1173](https://github.com/docker/buildx/issues/1173)
- You can now use the new `--print` flag to run helper functions supported by the
  BuildKit frontend performing the build and print their results. You can use
  this feature in Dockerfile to show the build arguments and secrets that the
  current build supports with `--print=outline` and list all available
  Dockerfile stages with `--print=targets`. This feature is experimental for
  gathering early feedback and requires enabling `BUILDX_EXPERIMENTAL=1`
  environment variable. We plan to update/extend this feature in the future
  without keeping backward compatibility. [docker/buildx#1100](https://github.com/docker/buildx/issues/1100),
  [docker/buildx#1272](https://github.com/docker/buildx/issues/1272)
- You can now use the new `--invoke` flag to launch interactive containers from
  build results for an interactive debugging cycle. You can reload these
  containers with code changes or restore them to an initial state from the
  special monitor mode. This feature is experimental for gathering early
  feedback and requires enabling `BUILDX_EXPERIMENTAL=1` environment variable.
  We plan to update/extend this feature in the future without enabling backward
  compatibility.
  [docker/buildx#1168](https://github.com/docker/buildx/issues/1168),
  [docker/buildx#1257](https://github.com/docker/buildx),
  [docker/buildx#1259](https://github.com/docker/buildx/issues/1259)
- Buildx now understands environment variable `BUILDKIT_COLORS` and `NO_COLOR`
  to customize/disable the colors of interactive build progressbar. [docker/buildx#1230](https://github.com/docker/buildx/issues/1230),
  [docker/buildx#1226](https://github.com/docker/buildx/issues/1226)
- `buildx ls` command now shows the current BuildKit version of each builder
  instance. [docker/buildx#998](https://github.com/docker/buildx/issues/998)
- The `bake` command now loads `.env` file automatically when building Compose
  files for compatibility. [docker/buildx#1261](https://github.com/docker/buildx/issues/1261)
- Bake now supports Compose files with `cache_to` definition. [docker/buildx#1155](https://github.com/docker/buildx/issues/1155)
- Bake now supports new builtin function `timestamp()` to access current time. [docker/buildx#1214](https://github.com/docker/buildx/issues/1214)
- Bake now supports Compose build secrets definition. [docker/buildx#1069](https://github.com/docker/buildx/issues/1069)
- Additional build context configuration is now supported in Compose files via `x-bake`. [docker/buildx#1256](https://github.com/docker/buildx/issues/1256)
- Inspecting builder now shows current driver options configuration. [docker/buildx#1003](https://github.com/docker/buildx/issues/1003),
  [docker/buildx#1066](https://github.com/docker/buildx/issues/1066)

### Updates

- Updated the Compose Specification to 1.4.0. [docker/buildx#1246](https://github.com/docker/buildx/issues/1246),
  [docker/buildx#1251](https://github.com/docker/buildx/issues/1251)

### Bug fixes and enhancements

- The `buildx ls` command output has been updated with better access to errors
  from different builders. [docker/buildx#1109](https://github.com/docker/buildx/issues/1109)
- The `buildx create` command now performs additional validation of builder parameters
  to avoid creating a builder instance with invalid configuration. [docker/buildx#1206](https://github.com/docker/buildx/issues/1206)
- The `buildx imagetools create` command can now create new multi-platform images
  even if the source subimages are located on different repositories or
  registries. [docker/buildx#1137](https://github.com/docker/buildx/issues/1137)
- You can now set the default builder config that is used when creating
  builder instances without passing custom `--config` value. [docker/buildx#1111](https://github.com/docker/buildx/issues/1111)
- Docker driver can now detect if `dockerd` instance supports initially
  disabled Buildkit features like multi-platform images. [docker/buildx#1260](https://github.com/docker/buildx/issues/1260),
  [docker/buildx#1262](https://github.com/docker/buildx/issues/1262)
- Compose files using targets with `.` in the name are now converted to use `_`
  so the selector keys can still be used in such targets. [docker/buildx#1011](https://github.com/docker/buildx/issues/1011)
- Included an additional validation for checking valid driver configurations. [docker/buildx#1188](https://github.com/docker/buildx/issues/1188),
  [docker/buildx#1273](https://github.com/docker/buildx/issues/1273)
- The `remove` command now displays the removed builder and forbids removing
  context builders. [docker/buildx#1128](https://github.com/docker/buildx/issues/1128)
- Enable Azure authentication when using Kubernetes driver. [docker/buildx#974](https://github.com/docker/buildx/issues/974)
- Add tolerations handling for kubernetes driver. [docker/buildx#1045](https://github.com/docker/buildx/issues/1045)
  [docker/buildx#1053](https://github.com/docker/buildx/issues/1053)
- Replace deprecated seccomp annotations with `securityContext` in the `kubernetes` driver.
  [docker/buildx#1052](https://github.com/docker/buildx/issues/1052)
- Fix panic on handling manifests with nil platform. [docker/buildx#1144](https://github.com/docker/buildx/issues/1144)
- Fix using duration filter with `prune` command. [docker/buildx#1252](https://github.com/docker/buildx/issues/1252)
- Fix merging multiple JSON files on Bake definition. [docker/buildx#1025](https://github.com/docker/buildx/issues/1025)
- Fix issues with implicit builder created from Docker context had invalid
  configuration or dropped connection. [docker/buildx#1129](https://github.com/docker/buildx/issues/1129)
- Fix conditions for showing no-output warning when using named contexts. [docker/buildx#968](https://github.com/docker/buildx/issues/968)
- Fix duplicating builders when builder instance and docker context have the
  same name. [docker/buildx#1131](https://github.com/docker/buildx/issues/1131)
- Fix printing unnecessary SSH warning logs. [docker/buildx#1085](https://github.com/docker/buildx/issues/1085)
- Fix possible panic when using an empty variable block with Bake JSON
  definition. [docker/buildx#1080](https://github.com/docker/buildx/issues/1080)
- Fix image tools commands not handling `--builder` flag correctly. [docker/buildx#1067](https://github.com/docker/buildx/issues/1067)
- Fix using custom image together with rootless option. [docker/buildx#1063](https://github.com/docker/buildx/issues/1063)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.9.0).

## 0.8.2

{{< release-date date="2022-04-04" >}}

### Updates

- Update Compose spec used by `buildx bake` to v1.2.1 to fix parsing ports definition. [docker/buildx#1033](https://github.com/docker/buildx/issues/1033)

### Bug fixes and enhancements

- Fix possible crash on handling progress streams from BuildKit v0.10. [docker/buildx#1042](https://github.com/docker/buildx/issues/1042)
- Fix parsing groups in `buildx bake` when already loaded by a parent group. [docker/buildx#1021](https://github.com/docker/buildx/issues/1021)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.8.2).

## 0.8.1

{{< release-date date="2022-03-21" >}}

### Bug fixes and enhancements

- Fix possible panic on handling build context scanning errors. [docker/buildx#1005](https://github.com/docker/buildx/issues/1005)
- Allow `.` on Compose target names in `buildx bake` for backward compatibility. [docker/buildx#1018](https://github.com/docker/buildx/issues/1018)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.8.1).

## 0.8.0

{{< release-date date="2022-03-09" >}}

### New

- Build command now accepts `--build-context` flag to [define additional named build contexts](/reference/cli/docker/buildx/build/#build-context)
  for your builds. [docker/buildx#904](https://github.com/docker/buildx/issues/904)
- Bake definitions now support [defining dependencies between targets](bake/contexts.md)
  and using the result of one target in another build.
  [docker/buildx#928](https://github.com/docker/buildx/issues/928),
  [docker/buildx#965](https://github.com/docker/buildx/issues/965),
  [docker/buildx#963](https://github.com/docker/buildx/issues/963),
  [docker/buildx#962](https://github.com/docker/buildx/issues/962),
  [docker/buildx#981](https://github.com/docker/buildx/issues/981)
- `imagetools inspect` now accepts `--format` flag allowing access to config
  and buildinfo for specific images. [docker/buildx#854](https://github.com/docker/buildx/issues/854),
  [docker/buildx#972](https://github.com/docker/buildx/issues/972)
- New flag `--no-cache-filter` allows configuring build, so it ignores cache
  only for specified Dockerfile stages. [docker/buildx#860](https://github.com/docker/buildx/issues/860)
- Builds can now show a summary of warnings sets by the building frontend. [docker/buildx#892](https://github.com/docker/buildx/issues/892)
- The new build argument `BUILDKIT_INLINE_BUILDINFO_ATTRS` allows opting-in to embed
  building attributes to resulting image. [docker/buildx#908](https://github.com/docker/buildx/issues/908)
- The new flag `--keep-buildkitd` allows keeping BuildKit daemon running when removing a builder
  - [docker/buildx#852](https://github.com/docker/buildx/issues/852)

### Bug fixes and enhancements

- `--metadata-file` output now supports embedded structure types. [docker/buildx#946](https://github.com/docker/buildx/issues/946)
- `buildx rm` now accepts new flag `--all-inactive` for removing all builders
  that are not currently running. [docker/buildx#885](https://github.com/docker/buildx/issues/885)
- Proxy config is now read from Docker configuration file and sent with build
  requests for backward compatibility. [docker/buildx#959](https://github.com/docker/buildx/issues/959)
- Support host networking in Compose. [docker/buildx#905](https://github.com/docker/buildx/issues/905),
  [docker/buildx#880](https://github.com/docker/buildx/issues/880)
- Bake files can now be read from stdin with `-f -`. [docker/buildx#864](https://github.com/docker/buildx/issues/864)
- `--iidfile` now always writes the image config digest independently of the
  driver being used (use `--metadata-file` for digest). [docker/buildx#980](https://github.com/docker/buildx/issues/980)
- Target names in Bake are now restricted to not use special characters. [docker/buildx#929](https://github.com/docker/buildx/issues/929)
- Image manifest digest can be read from metadata when pushed with `docker`
  driver. [docker/buildx#989](https://github.com/docker/buildx/issues/989)
- Fix environment file handling in Compose files. [docker/buildx#905](https://github.com/docker/buildx/issues/905)
- Show last access time in `du` command. [docker/buildx#867](https://github.com/docker/buildx/issues/867)
- Fix possible double output logs when multiple Bake targets run same build
  steps. [docker/buildx#977](https://github.com/docker/buildx/issues/977)
- Fix possible errors on multi-node builder building multiple targets with
  mixed platform. [docker/buildx#985](https://github.com/docker/buildx/issues/985)
- Fix some nested inheritance cases in Bake. [docker/buildx#914](https://github.com/docker/buildx/issues/914)
- Fix printing default group on Bake files. [docker/buildx#884](https://github.com/docker/buildx/issues/884)
- Fix `UsernsMode` when using rootless container. [docker/buildx#887](https://github.com/docker/buildx/issues/887)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.8.0).

## 0.7.1

{{< release-date date="2021-08-25" >}}

### Fixes

- Fix issue with matching exclude rules in `.dockerignore`. [docker/buildx#858](https://github.com/docker/buildx/issues/858)
- Fix `bake --print` JSON output for current group. [docker/buildx#857](https://github.com/docker/buildx/issues/857)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.7.1).

## 0.7.0

{{< release-date date="2021-11-10" >}}

### New features

- TLS certificates from BuildKit configuration are now transferred to build
  container with `docker-container` and `kubernetes` drivers. [docker/buildx#787](https://github.com/docker/buildx/issues/787)
- Builds support `--ulimit` flag for feature parity. [docker/buildx#800](https://github.com/docker/buildx/issues/800)
- Builds support `--shm-size` flag for feature parity. [docker/buildx#790](https://github.com/docker/buildx/issues/790)
- Builds support `--quiet` for feature parity. [docker/buildx#740](https://github.com/docker/buildx/issues/740)
- Builds support `--cgroup-parent` flag for feature parity. [docker/buildx#814](https://github.com/docker/buildx/issues/814)
- Bake supports builtin variable `BAKE_LOCAL_PLATFORM`. [docker/buildx#748](https://github.com/docker/buildx/issues/748)
- Bake supports `x-bake` extension field in Compose files. [docker/buildx#721](https://github.com/docker/buildx/issues/721)
- `kubernetes` driver now supports colon-separated `KUBECONFIG`. [docker/buildx#761](https://github.com/docker/buildx/issues/761)
- `kubernetes` driver now supports setting Buildkit config file with `--config`. [docker/buildx#682](https://github.com/docker/buildx/issues/682)
- `kubernetes` driver now supports installing QEMU emulators with driver-opt. [docker/buildx#682](https://github.com/docker/buildx/issues/682)

### Enhancements

- Allow using custom registry configuration for multi-node pushes from the
  client. [docker/buildx#825](https://github.com/docker/buildx/issues/825)
- Allow using custom registry configuration for `buildx imagetools` command. [docker/buildx#825](https://github.com/docker/buildx/issues/825)
- Allow booting builder after creating with `buildx create --bootstrap`. [docker/buildx#692](https://github.com/docker/buildx/issues/692)
- Allow `registry:insecure` output option for multi-node pushes. [docker/buildx#825](https://github.com/docker/buildx/issues/825)
- BuildKit config and TLS files are now kept in Buildx state directory and
  reused if BuildKit instance needs to be recreated. [docker/buildx#824](https://github.com/docker/buildx/issues/824)
- Ensure different projects use separate destination directories for
  incremental context transfer for better performance. [docker/buildx#817](https://github.com/docker/buildx/issues/817)
- Build containers are now placed on separate cgroup by default. [docker/buildx#782](https://github.com/docker/buildx/issues/782)
- Bake now prints the default group with `--print`. [docker/buildx#720](https://github.com/docker/buildx/issues/720)
- `docker` driver now dials build session over HTTP for better performance. [docker/buildx#804](https://github.com/docker/buildx/issues/804)

### Fixes

- Fix using `--iidfile` together with a multi-node push. [docker/buildx#826](https://github.com/docker/buildx/issues/826)
- Using `--push` in Bake does not clear other image export options in the file. [docker/buildx#773](https://github.com/docker/buildx/issues/773)
- Fix Git URL detection for `buildx bake` when `https` protocol was used. [docker/buildx#822](https://github.com/docker/buildx/issues/822)
- Fix pushing image with multiple names on multi-node builds. [docker/buildx#815](https://github.com/docker/buildx/issues/815)
- Avoid showing `--builder` flags for commands that don't use it. [docker/buildx#818](https://github.com/docker/buildx/issues/818)
- Unsupported build flags now show a warning. [docker/buildx#810](https://github.com/docker/buildx/issues/810)
- Fix reporting error details in some OpenTelemetry traces. [docker/buildx#812](https://github.com/docker/buildx/issues/812)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.7.0).

## 0.6.3

{{< release-date date="2021-08-30" >}}

### Fixes

- Fix BuildKit state volume location for Windows clients. [docker/buildx#751](https://github.com/docker/buildx/issues/751)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.6.3).

## 0.6.2

{{< release-date date="2021-08-21" >}}

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.6.2).

### Fixes

- Fix connection error showing up in some SSH configurations. [docker/buildx#741](https://github.com/docker/buildx/issues/741)

## 0.6.1

{{< release-date date="2021-07-30" >}}

### Enhancements

- Set `ConfigFile` to parse compose files with Bake. [docker/buildx#704](https://github.com/docker/buildx/issues/704)

### Fixes

- Duplicate progress env var. [docker/buildx#693](https://github.com/docker/buildx/issues/693)
- Should ignore nil client. [docker/buildx#686](https://github.com/docker/buildx/issues/686)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.6.1).

## 0.6.0

{{< release-date date="2021-07-16" >}}

### New features

- Support for OpenTelemetry traces and forwarding Buildx client traces to
  BuildKit. [docker/buildx#635](https://github.com/docker/buildx/issues/635)
- Experimental GitHub Actions remote cache backend with `--cache-to type=gha`
  and `--cache-from type=gha`. [docker/buildx#535](https://github.com/docker/buildx/issues/535)
- New `--metadata-file` flag has been added to build and Bake command that
  allows saving build result metadata in JSON format. [docker/buildx#605](https://github.com/docker/buildx/issues/605)
- This is the first release supporting Windows ARM64. [docker/buildx#654](https://github.com/docker/buildx/issues/654)
- This is the first release supporting Linux Risc-V. [docker/buildx#652](https://github.com/docker/buildx/issues/652)
- Bake now supports building from remote definition with local files or
  another remote source as context. [docker/buildx#671](https://github.com/docker/buildx/issues/671)
- Bake now allows variables to reference each other and using user functions
  in variables and vice-versa.
  [docker/buildx#575](https://github.com/docker/buildx/issues/575),
  [docker/buildx#539](https://github.com/docker/buildx/issues/539),
  [docker/buildx#532](https://github.com/docker/buildx/issues/532)
- Bake allows defining attributes in the global scope. [docker/buildx#541](https://github.com/docker/buildx/issues/541)
- Bake allows variables across multiple files. [docker/buildx#538](https://github.com/docker/buildx/issues/538)
- New quiet mode has been added to progress printer. [docker/buildx#558](https://github.com/docker/buildx/issues/558)
- `kubernetes` driver now supports defining resources/limits. [docker/buildx#618](https://github.com/docker/buildx/issues/618)
- Buildx binaries can now be accessed through [buildx-bin](https://hub.docker.com/r/docker/buildx-bin)
  Docker image. [docker/buildx#656](https://github.com/docker/buildx/issues/656)

### Enhancements

- `docker-container` driver now keeps BuildKit state in volume. Enabling
  updates with keeping state. [docker/buildx#672](https://github.com/docker/buildx/issues/672)
- Compose parser is now based on new [compose-go parser](https://github.com/compose-spec/compose-go)
  fixing support for some newer syntax. [docker/buildx#669](https://github.com/docker/buildx/issues/669)
- SSH socket is now automatically forwarded when building an ssh-based git URL. [docker/buildx#581](https://github.com/docker/buildx/issues/581)
- Bake HCL parser has been rewritten. [docker/buildx#645](https://github.com/docker/buildx/issues/645)
- Extend HCL support with more functions. [docker/buildx#491](https://github.com/docker/buildx/issues/491)
  [docker/buildx#503](https://github.com/docker/buildx/issues/503)
- Allow secrets from environment variables. [docker/buildx#488](https://github.com/docker/buildx/issues/488)
- Builds with an unsupported multi-platform and load configuration now fail fast. [docker/buildx#582](https://github.com/docker/buildx/issues/582)
- Store Kubernetes config file to make buildx builder switchable. [docker/buildx#497](https://github.com/docker/buildx/issues/497)
- Kubernetes now lists all pods as nodes on inspection. [docker/buildx#477](https://github.com/docker/buildx/issues/477)
- Default Rootless image has been set to `moby/buildkit:buildx-stable-1-rootless`. [docker/buildx#480](https://github.com/docker/buildx/issues/480)

### Fixes

- `imagetools create` command now correctly merges JSON descriptor with old one. [docker/buildx#592](https://github.com/docker/buildx/issues/592)
- Fix building with `--network=none` not requiring extra security entitlements. [docker/buildx#531](https://github.com/docker/buildx/issues/531)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.6.0).

## 0.5.1

{{< release-date date="2020-12-15" >}}

### Fixes

- Fix regression on setting `--platform` on `buildx create` outside
  `kubernetes` driver. [docker/buildx#475](https://github.com/docker/buildx/issues/475)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.5.1).

## 0.5.0

{{< release-date date="2020-12-15" >}}

### New features

- The `docker` driver now supports the `--push` flag. [docker/buildx#442](https://github.com/docker/buildx/issues/442)
- Bake supports inline Dockerfiles. [docker/buildx#398](https://github.com/docker/buildx/issues/398)
- Bake supports building from remote URLs and Git repositories. [docker/buildx#398](https://github.com/docker/buildx/issues/398)
- `BUILDX_CONFIG` env var allow users to have separate buildx state from
  Docker config. [docker/buildx#385](https://github.com/docker/buildx/issues/385)
- `BUILDKIT_MULTI_PLATFORM` build arg allows to force building multi-platform
  return objects even if only one `--platform` specified. [docker/buildx#467](https://github.com/docker/buildx/issues/467)

### Enhancements

- Allow `--append` to be used with `kubernetes` driver. [docker/buildx#370](https://github.com/docker/buildx/issues/370)
- Build errors show error location in source files and system stacktraces
  with `--debug`. [docker/buildx#389](https://github.com/docker/buildx/issues/389)
- Bake formats HCL errors with source definition. [docker/buildx#391](https://github.com/docker/buildx/issues/391)
- Bake allows empty string values in arrays that will be discarded. [docker/buildx#428](https://github.com/docker/buildx/issues/428)
- You can now use the Kubernetes cluster config with the `kubernetes` driver. [docker/buildx#368](https://github.com/docker/buildx/issues/368)
  [docker/buildx#460](https://github.com/docker/buildx/issues/460)
- Creates a temporary token for pulling images instead of sharing credentials
  when possible. [docker/buildx#469](https://github.com/docker/buildx/issues/469)
- Ensure credentials are passed when pulling BuildKit container image. [docker/buildx#441](https://github.com/docker/buildx/issues/441)
  [docker/buildx#433](https://github.com/docker/buildx/issues/433)
- Disable user namespace remapping in `docker-container` driver. [docker/buildx#462](https://github.com/docker/buildx/issues/462)
- Allow `--builder` flag to switch to default instance. [docker/buildx#425](https://github.com/docker/buildx/issues/425)
- Avoid warn on empty `BUILDX_NO_DEFAULT_LOAD` config value. [docker/buildx#390](https://github.com/docker/buildx/issues/390)
- Replace error generated by `quiet` option by a warning. [docker/buildx#403](https://github.com/docker/buildx/issues/403)
- CI has been switched to GitHub Actions.
  [docker/buildx#451](https://github.com/docker/buildx/issues/451),
  [docker/buildx#463](https://github.com/docker/buildx/issues/463),
  [docker/buildx#466](https://github.com/docker/buildx/issues/466),
  [docker/buildx#468](https://github.com/docker/buildx/issues/468),
  [docker/buildx#471](https://github.com/docker/buildx/issues/471)

### Fixes

- Handle lowercase Dockerfile name as a fallback for backward compatibility. [docker/buildx#444](https://github.com/docker/buildx/issues/444)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.5.0).

## 0.4.2

{{< release-date date="2020-08-22" >}}

### New features

- Support `cacheonly` exporter. [docker/buildx#337](https://github.com/docker/buildx/issues/337)

### Enhancements

- Update `go-cty` to pull in more `stdlib` functions. [docker/buildx#277](https://github.com/docker/buildx/issues/277)
- Improve error checking on load. [docker/buildx#281](https://github.com/docker/buildx/issues/281)

### Fixes

- Fix parsing json config with HCL. [docker/buildx#280](https://github.com/docker/buildx/issues/280)
- Ensure `--builder` is wired from root options. [docker/buildx#321](https://github.com/docker/buildx/issues/321)
- Remove warning for multi-platform iidfile. [docker/buildx#351](https://github.com/docker/buildx/issues/351)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.4.2).

## 0.4.1

{{< release-date date="2020-05-01" >}}

### Fixes

- Fix regression on flag parsing. [docker/buildx#268](https://github.com/docker/buildx/issues/268)
- Fix using pull and no-cache keys in HCL targets. [docker/buildx#268](https://github.com/docker/buildx/issues/268)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.4.1).

## 0.4.0

{{< release-date date="2020-04-30" >}}

### New features

- Add `kubernetes` driver. [docker/buildx#167](https://github.com/docker/buildx/issues/167)
- New global `--builder` flag to override builder instance for a single command. [docker/buildx#246](https://github.com/docker/buildx/issues/246)
- New `prune` and `du` commands for managing local builder cache. [docker/buildx#249](https://github.com/docker/buildx/issues/249)
- You can now set the new `pull` and `no-cache` options for HCL targets. [docker/buildx#165](https://github.com/docker/buildx/issues/165)

### Enhancements

- Upgrade Bake to HCL2 with support for variables and functions. [docker/buildx#192](https://github.com/docker/buildx/issues/192)
- Bake now supports `--load` and `--push`. [docker/buildx#164](https://github.com/docker/buildx/issues/164)
- Bake now supports wildcard overrides for multiple targets. [docker/buildx#164](https://github.com/docker/buildx/issues/164)
- Container driver allows setting environment variables via `driver-opt`. [docker/buildx#170](https://github.com/docker/buildx/issues/170)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.4.0).

## 0.3.1

{{< release-date date="2019-09-27" >}}

### Enhancements

- Handle copying unix sockets instead of erroring. [docker/buildx#155](https://github.com/docker/buildx/issues/155)
  [moby/buildkit#1144](https://github.com/moby/buildkit/issues/1144)

### Fixes

- Running Bake with multiple Compose files now merges targets correctly. [docker/buildx#134](https://github.com/docker/buildx/issues/134)
- Fix bug when building a Dockerfile from stdin (`build -f -`).
  [docker/buildx#153](https://github.com/docker/buildx/issues/153)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.3.1).

## 0.3.0

{{< release-date date="2019-08-02" >}}

### New features

- Custom `buildkitd` daemon flags. [docker/buildx#102](https://github.com/docker/buildx/issues/102)
- Driver-specific options on `create`. [docker/buildx#122](https://github.com/docker/buildx/issues/122)

### Enhancements

- Environment variables are used in Compose files. [docker/buildx#117](https://github.com/docker/buildx/issues/117)
- Bake now honors `--no-cache` and `--pull`. [docker/buildx#118](https://github.com/docker/buildx/issues/118)
- Custom BuildKit config file. [docker/buildx#121](https://github.com/docker/buildx/issues/121)
- Entitlements support with `build --allow`. [docker/buildx#104](https://github.com/docker/buildx/issues/104)

### Fixes

- Fix bug where `--build-arg foo` would not read `foo` from environment. [docker/buildx#116](https://github.com/docker/buildx/issues/116)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.3.0).

## 0.2.2

{{< release-date date="2019-05-30" >}}

### Enhancements

- Change Compose file handling to require valid service specifications. [docker/buildx#87](https://github.com/docker/buildx/issues/87)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.2.2).

## 0.2.1

{{< release-date date="2019-05-25" >}}

### New features

- Add `BUILDKIT_PROGRESS` env var. [docker/buildx#69](https://github.com/docker/buildx/issues/69)
- Add `local` platform. [docker/buildx#70](https://github.com/docker/buildx/issues/70)

### Enhancements

- Keep arm variant if one is defined in the config. [docker/buildx#68](https://github.com/docker/buildx/issues/68)
- Make dockerfile relative to context. [docker/buildx#83](https://github.com/docker/buildx/issues/83)

### Fixes

- Fix parsing target from compose files. [docker/buildx#53](https://github.com/docker/buildx/issues/53)

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.2.1).

## 0.2.0

{{< release-date date="2019-04-25" >}}

### New features

- First release

For more details, see the complete release notes in the [Buildx GitHub repository](https://github.com/docker/buildx/releases/tag/v0.2.0).
