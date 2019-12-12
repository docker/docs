---
title: Configure Docker Assemble
description: Installing Docker Assemble
keywords: Assemble, Docker Enterprise, plugin, Spring Boot, .NET, c#, F#
---

Although you don’t need to configure anything to build a project using Docker Assemble, you may wish to  override the defaults, and in some cases, add fields that weren’t automatically detected from the project file. To support this, Docker Assemble allows you to add a file `docker-assemble.yaml` to the root of your project. The settings you provide in the `docker-assemble.yaml` file overrides any auto-detection and can themselves be overridden by command-line arguments

The `docker-assemble.yaml` file is in YAML syntax and has the following informal schema:

- `version`: (_string_) mandatory, must contain `0.2.0`

- `image`: (_map_) contains options related to the output image.

    - `platforms`: (_list of strings_) lists the possible platforms which can be built (for example, `linux/amd64`, `windows/amd64`). The default is determined automatically from the project type and content. Note that by default Docker Assemble will build only for `linux/amd64` unless `--push` is used. See [Building Multiplatform images](/assemble/images/#multi-platform-images).

    - `ports`: (_list of strings_) contains ports to expose from a container running the image. e.g `80/tcp` or `8080`. Default is to automatically determine the set of ports to expose where possible. To disable this and export no ports specify a list containing precisely one element of `none`.

    - `labels`: (_map_) contains labels to write into the image as `key`-`value` (_string_) pairs.

    - `repository-namespace`: (_string_) the registry and path component of the desired output image. e.g. `docker.io/library` or `docker.io/user`.

    - `repository-name`: (_string_) the name of the specific image within `repository-namespace`. Overrides any name derived from the build system specific configuration.

    - `tag`: (_string_) the default tag to use. Overrides and version/tag derived from the build system specific configuration.

    - `healthcheck`: (_map_) describes how to check a container running the image is healthy.

        - `kind`: (_string_) sets the type of Healthcheck to perform. Valid values are `none`, `simple-tcpport-open` and `springboot`. See [Health checks](/assemble/images/#health-checks).

        - `interval`: (_duration_) the time to wait between checks.

        - `timeout`: (_duration_) the time to wait before considering the check to have hung.

        - `start-period`: (_duration_) period for the container to initialize before the retries starts to count down

        - `retries`: (_integer_) number of consecutive failures needed to consider a container as unhealthy.

- `springboot`: (_map_) if this is a Spring Boot project then contains related configuration options.

    - `enabled`: (_boolean_) true if this is a springboot project.

    - `java-version`: (_string_) configures the Java version to use. Valid options are `8` and `10`.

    - `build-image`: (_string_) sets a custom base build image

    - `runtime-images` (_map_) sets a custom base runtime image by platform. For valid keys, refer to the **Spring Boot** section in [Custom base images](/assemble/images/#custom-base-images).

- `aspnetcore`: (_map_) if this is an ASP.NET Core project then contains related configuration options.

    - `enabled`: (_boolean_) true if this is an ASP.NET Core project.

    - `version`: (_string_) configures the ASP.NET Core version to use. Valid options are `1.0`, `1.1`, `2.0` and `2.1`.

    - `build-image`: (_string_) sets a custom base build image

    - `runtime-images` (_map_) sets a custom base runtime image by platform. For valid keys, refer to the **ASP.NET Core** section in [Custom base images](/assemble/images/#custom-base-images).

> **Notes:**
>
> - The only mandatory field in `docker-assemble.yaml` is `version`. All other parameters are optional.
>
> - At most one of `dotnet` or `springboot` can be present in the yaml file.
>
> - Fields of type duration are integers with nanosecond granularity. However the following units of time are supported: `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. For example, `25s`.

Each setting in the configuration file has a command line equivalent which can be used with the `-o/--option` argument, which takes a `KEY=VALUE` string where `KEY` is constructed by joining each element of the YAML hierarchy with a period (.).

For example,  the `image → repository-namespace` key in the YAML becomes `-o image.repository-namespace=NAME` on the command line and `springboot → enabled` becomes `-o springboot.enabled=BOOLEAN`.

The following convenience aliases take precedence over the `-o/--option` equivalents:

- `--namespace` is an alias for `image.repository-namespace`;

- `--name` corresponds to `image.repository-name`;

- `--tag` corresponds to `image.tag`;

- `--label` corresponds to `image.labels` (can be used multiple times);

- `--port` corresponds to `image.ports` (can be used multiple times)
