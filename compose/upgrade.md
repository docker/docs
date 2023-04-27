---
title: Upgrade Compose
description: how to upgrade Compose V1 to V2 and FAQs
keywords: compose, upgrade, v1, v2
---

From the end of June 2023 Compose V1 won’t be supported anymore and will be removed from all Docker Desktop versions.

For detailed information on the history of Compose and what Compose V1 includes, see the [Evolution of Compose](compose-v2/index.md)

## What are the functional differences between Compose V1 and Compose V2?

Compose V2 integrates compose functions into the Docker platform, continuing to support most of the previous `docker-compose` features and flags. You can run Compose V2 by replacing the hyphen (`-`) with a space, using `docker compose`, instead of `docker-compose`.

The `compose` command in the Docker CLI supports most of the `docker-compose` commands and flags. It is expected to be a drop-in replacement for `docker-compose`. 

If you see any Compose functionality that is not available in the `compose` command, create an issue in the [Compose](https://github.com/docker/compose/issues){:target="_blank" rel="noopener" class="_"} GitHub repository, so we can prioritize it.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab1">Commands not yet implemented</a></li>
  <li><a data-toggle="tab" data-target="#tab2">Flags not be implemented</a></li>
  <li><a data-toggle="tab" data-target="#tab3">New commands in Compose v2</a></li>
</ul>
<div class="tab-content">
<div id="tab1" class="tab-pane fade in active" markdown="1">

The following commands have not been implemented yet, and may be implemented at a later time.
Let us know if these commands are a higher priority for your use cases.

`compose build --memory`: This option is not yet supported by BuildKit. The flag is currently supported, but is hidden to avoid breaking existing Compose usage. It does not have any effect.

<hr>
</div>
<div id="tab2" class="tab-pane fade" markdown="1">

The list below includes the flags that we are not planning to support in Compose in the Docker CLI,
either because they are already deprecated in `docker-compose`, or because they are not relevant for Compose in the Docker CLI.

* `compose ps --filter KEY-VALUE` Not relevant due to its complicated usage with the `service` command and also because it is not documented properly in `docker-compose`.
* `compose rm --all` Deprecated in docker-compose.
* `compose scale` Deprecated in docker-compose (use `compose up --scale` instead)

Global flags:

* `--compatibility` has been resignified Docker Compose V2. This now means that in the command running V2 will behave as V1 used to do.
  * One difference is in the word separator on container names. V1 used to use `_` as separator while V2 uses `-` to keep the names more hostname friendly. So when using `--compatibility` Docker 
    Compose should use `_` again. Just make sure to stick to one of them otherwise Docker Compose will not be able to recognize the container as an instance of the service.
<hr>
</div>
<div id="tab3" class="tab-pane fade" markdown="1">

#### Copy

The `cp` command is intended to copy files or folders between service containers and the local filesystem.  
This command is a bidirectional command, we can copy **from** or **to** the service containers.

Copy a file from a service container to the local filesystem:

```console
$ docker compose cp my-service:~/path/to/myfile ~/local/path/to/copied/file
```

We can also copy from the local filesystem to all the running containers of a service:

```console
$ docker compose cp --all ~/local/path/to/source/file my-service:~/path/to/copied/file
```


#### List

The ls command is intended to list the Compose projects. By default, the command only lists the running projects, 
we can use flags to display the stopped projects, to filter by conditions and change the output to `json` format for example.

```console
$ docker compose ls --all --format json
[{"Name":"dockergithubio","Status":"exited(1)","ConfigFiles":"/path/to/docs/docker-compose.yml"}]
```

### Use `--project-name` with Compose commands

With the GA version of Compose, you can run some commands:
- outside of directory containing the project compose file
- or without specifying the path of the Compose with the `--file` flag
- or without specifying the project directory with the `--project-directory` flag

When a compose project has been loaded once, we can just use the `-p` or `--project-name` to reference it:

```console
$ docker compose -p my-loaded-project restart my-service
```

This option works with the `start`, `stop`, `restart` and `down` commands.

### Config command

The config command is intended to show the configuration used by Docker Compose to run the actual project after normalization and templating. The resulting output might contain superficial differences in formattting and style.
For example, some fields in the Compose Specification support both short and a long format so the output structure might not match the input structure but is guaranteed to be semantically equivalent.

Similarly, comments in the source file are not preserved.

In the example below we can see the config command expanding the `ports` section:

docker-compose.yml:
```yaml
services:
  web:
    # default to latest but allow overriding the tag
    image: nginx:${TAG-latest}
    ports:
      - 80:80
```
With `$ docker compose config` the output turns into:
```yaml
name: docs-example
services:
  web:
    image: nginx:stable-alpine
    networks:
      default: null
    ports:
    - mode: ingress
      target: 80
      published: "80"
      protocol: tcp
networks:
  default:
    name: basic_default
```

The result above is a full size configuration of what will be used by Docker Compose to run the project.
<hr>
</div>
</div>

## What does this mean for my projects that use Compose V1?

For almost all projects, moving to Compose V2 requires no changes to the Compose YAML or your workflow.

We recommend that you adapt to the new preferred way of running Compose v2, which is to use `docker compose` instead of `docker-compose`.

However, since the introduction of Compose V2, Docker Desktop provides a shim/alias to redirect `docker-compose` commands to `docker compose` when the **Use Compose V2** setting is selected.

Docker Desktop will continue to provide this alias by default for the foreseeable future.

## How do I upgrade?

The easiest and recommended way to upgrade to Compose V2, is to make sure you have the latest version of [Docker Desktop](../desktop/release-notes.md). Docker Desktop includes Compose V2 along with Docker Engine and Docker CLI which are Compose prerequisites.

If you already have Docker Engine and Docker CLI installed, you can install the Compose plugin from the command line, by either:
- [Using Docker's repository](linux.md#install-using-the-repository)
- [Downloading and installing manually](linux.md#install-the-plugin-manually)
>This is only available on Linux
{: .important}

## Is there anything else I need to know before I upgrade?

Projects that rely on variable interpolation should take extra care during upgrade. The behavior in Compose v1 was not formally documented and behaved inconsistently.

We have provided documentation showing the supported syntaxes and numerous examples that cover tricky edge cases such as escaping quotes.

#### Ask yourself…

- Do any values have nested quotes or escapes?
- Do any values contain literal `$` signs in them? (This is common with PHP projects!)
- Do any values use advanced expansion syntax? (e.g. `${VAR:?error}`)

#### Recommended Action

Run `docker compose config` on the project to preview the configuration after Compose has performed interpolation to verify that values appear as expected.

#### Backwards Compatibility

Maintaining backwards compatibility with v1 is *********typically********* achievable by ensuring that literal values (no interpolation) are single-quoted and values that should have interpolation applied are double-quoted.

## Will I still be able to use Compose V1 if I really want to?

Yes. You can still download and install Compose V1 packages, but you won't get support from Docker if anything breaks. 

>**Warning**
>
> The final Compose v1 release (v1.29.2) was May 10, 2021. These packages haven't received any security updates since then. Use at your own risk. 
{: .warning}

