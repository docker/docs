---
title: Migrate to Compose V2
description: How to migrate from Compose V1 to V2
keywords: compose, upgrade, migration, v1, v2, docker compose vs docker-compose
aliases:
- /compose/compose-v2/
- /compose/cli-command-compatibility/
---

From July 2023 Compose V1 stopped receiving updates. It’s also no longer available in new releases of Docker Desktop.

Compose V2, which was first released in 2020, is included with all currently supported versions of Docker Desktop. It offers an improved CLI experience, improved build performance with BuildKit, and continued new-feature development.

## How do I switch to Compose V2?

The easiest and recommended way is to make sure you have the latest version of [Docker Desktop](../desktop/release-notes.md), which bundles the Docker Engine and Docker CLI platform including Compose V2.

With Docker Desktop, Compose V2 is always accessible as `docker compose`.
Additionally, the **Use Compose V2** setting is turned on by default, which provides an alias from `docker-compose`.

For manual installs on Linux, you can get Compose V2 by either:
- [Using Docker's repository](install/linux.md#install-using-the-repository) (recommended)
- [Downloading and installing manually](install/linux.md#install-the-plugin-manually)

## What are the differences between Compose V1 and Compose V2?

### `docker-compose` vs `docker compose`

Unlike Compose V1, Compose V2 integrates into the Docker CLI platform and the recommended command-line syntax is `docker compose`.

The Docker CLI platform provides a consistent and predictable set of options and flags, such as the `DOCKER_HOST` environment variable or the `--context` command-line flag.

This change lets you use all of the shared flags on the root `docker` command.
For example, `docker --log-level=debug --tls compose up` enables debug logging from the Docker Engine as well as ensuring that TLS is used for the connection.

> **Tip**
>
> Update scripts to use Compose V2 by replacing the hyphen (`-`) with a space, using `docker compose` instead of `docker-compose`.
{ .tip }

### Service container names

Compose generates container names based on the project name, service name, and scale/replica count.

In Compose V1, an underscore (`_`) was used as the word separator.
In Compose V2, a hyphen (`-`) is used as the word separator.

Underscores aren't valid characters in DNS hostnames.
By using a hyphen instead, Compose V2 ensures service containers can be accessed over the network via consistent, predictable hostnames.
 
For example, running the Compose command `-p myproject up --scale=1 svc` results in a container named `myproject_svc_1` with Compose V1 and a container named `myproject-svc-1` with Compose V2.

> **Tip**
>
>In Compose V2, the global `--compatibility` flag or `COMPOSE_COMPATIBILITY` environment variable preserves the Compose V1 behavior to use underscores (`_`) as the word separator.
As this option must be specified for every Compose V2 command run, it's recommended that you only use this as a temporary measure while transitioning to Compose V2.
{ .tip }

### Command-line flags and subcommands

Compose V2 supports almost all Compose V1 flags and subcommands, so in most cases, it can be used as a drop-in replacement in scripts.

#### Unsupported in V2

The following were deprecated in Compose V1 and aren't supported in Compose V2:
* `docker-compose scale`. Use `docker compose up --scale` instead.
* `docker-compose rm --all`

#### Different in V2

The following behave differently between Compose V1 and V2:

|                         | Compose V1                                                       | Compose V2                                                                    |
|-------------------------|------------------------------------------------------------------|-------------------------------------------------------------------------------|
| `--compatibility`       | Deprecated. Migrates YAML fields based on legacy schema version. | Uses `_` as word separator for container names instead of `-` to match V1.    |
| `ps --filter KEY-VALUE` | Undocumented. Allows filtering by arbitrary service properties.  | Only allows filtering by specific properties, e.g. `--filter=status=running`. |

### Environment variables

Environment variable behavior in Compose V1 wasn't formally documented and behaved inconsistently in some edge cases.

For Compose V2, the [Environment variables](/compose/environment-variables/) section covers both [precedence](/compose/environment-variables/envvars-precedence) as well as [`.env` file interpolation](environment-variables/variable-interpolation.md) and includes many examples covering tricky situations such as escaping nested quotes.

Check if:
- Your project uses multiple levels of environment variable overrides, for example `.env` file and `--env` CLI flags.
- Any `.env` file values have escape sequences or nested quotes.
- Any `.env` file values contain literal `$` signs in them. This is common with PHP projects.
- Any variable values use advanced expansion syntax, for example `${VAR:?error}`.

> **Tip**
>
> Run `docker compose config` on the project to preview the configuration after Compose V2 has performed interpolation to
verify that values appear as expected.
>
> Maintaining backwards compatibility with Compose V1 is typically achievable by ensuring that literal values (no
interpolation) are single-quoted and values that should have interpolation applied are double-quoted.
{ .tip }

## What does this mean for my projects that use Compose V1?

For most projects, switching to Compose V2 requires no changes to the Compose YAML or your development workflow.

It's recommended that you adapt to the new preferred way of running Compose V2, which is to use `docker compose` instead of `docker-compose`.
This provides additional flexibility and removes the requirement for a `docker-compose` compatibility alias. 

However, Docker Desktop continues to support a `docker-compose` alias to redirect commands to `docker compose` for convenience and improved compatibility with third-party tools and scripts.

## Is there anything else I need to know before I switch?

### Migrating running projects

In both V1 and V2, running `up` on a Compose project recreates service containers as necessary to reach the desired state based on comparing the actual state in the Docker Engine to the resolved project configuration including Compose YAML, environment variables, and command-line flags.

Because Compose V1 and V2 [name service containers differently](#service-container-names), running `up` using V2 the first time on a project with running services originally launched by V1, results in service containers being recreated with updated names.

Note that even if `--compatibility` flag is used to preserve the V1 naming style, Compose still needs to recreate service containers originally launched by V1 the first time `up` is run by V2 to migrate the internal state.

### Using Compose V2 with Docker-in-Docker

Compose V2 is now included in the [Docker official image on Docker Hub](https://hub.docker.com/_/docker).

Additionally, a new [docker/compose-bin image on Docker Hub](https://hub.docker.com/r/docker/compose-bin) packages the latest version of Compose V2 for use in multi-stage builds.

## Can I still use Compose V1 if I want to?

Yes. You can still download and install Compose V1 packages, but you won't get support from Docker if anything breaks.

>**Warning**
>
> The final Compose V1 release, version 1.29.2, was May 10, 2021. These packages haven't received any security updates since then. Use at your own risk. 
{ .warning }

## Additional Resources

- [docker-compose V1 on PyPI](https://pypi.org/project/docker-compose/1.29.2/)
- [docker/compose V1 on Docker Hub](https://hub.docker.com/r/docker/compose)
- [docker-compose V1 source on GitHub](https://github.com/docker/compose/releases/tag/1.29.2)
