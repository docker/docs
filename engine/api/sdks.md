---
title: SDKs for Docker Engine API
description: Client libraries for the Docker Engine API.
keywords: API, SDK, library, Docker, index, registry, REST, documentation, clients, C#, Erlang, Go, Groovy, Java, JavaScript, Perl, PHP, Python, Ruby, Rust, Scala
redirect_from:
  - /engine/api/client-libraries/
  - /engine/reference/api/remote_api_client_libraries/
  - /reference/api/remote_api_client_libraries/
---

The Docker SDKs allow you to build applications that can control and manage the Docker Engine. They are interfaces for the [Docker Engine API](index.md), but also contain a number of tools to make it easier to work with the API.

There are official libraries available in Python and Go, and there are a number of community supported libraries for other languages.

## Python

The Docker SDK for Python is available on the Python Package Index (PyPI), and can be installed with PIP:

    $ pip install docker

To see how to start using it, [head to the getting started guide](getting-started.md).

For a full reference, see the [Docker SDK for Python documentation](https://docker-py.readthedocs.io).

## Go

The Docker SDK for Go is a package inside the Docker Engine repository. To use it, you import it:

{% highlight go %}
import "github.com/docker/docker/client"
{% endhighlight %}

To see how to start using it, [head to the getting started guide](getting-started.md).

[A full reference is available on GoDoc.](https://godoc.org/github.com/moby/moby/client)

## Other languages

There a number of community supported libraries available for other languages. They have not been tested by the Docker maintainers for compatibility, so if you run into any issues, file them with the library maintainers.

| Language      | Library |
| ------------- |---------|
| C             | [libdocker](https://github.com/danielsuo/libdocker) |
| C#            | [Docker.DotNet](https://github.com/ahmetalpbalkan/Docker.DotNet) |
| C++           | [lasote/docker_client](https://github.com/lasote/docker_client) |
| Dart          | [bwu_docker](https://github.com/bwu-dart/bwu_docker) |
| Erlang        | [erldocker](https://github.com/proger/erldocker) |
| Gradle        | [gradle-docker-plugin](https://github.com/gesellix/gradle-docker-plugin) |
| Groovy        | [docker-client](https://github.com/gesellix/docker-client) |
| Haskell       | [docker-hs](https://github.com/denibertovic/docker-hs) |
| HTML (Web Components) | [docker-elements](https://github.com/kapalhq/docker-elements) |
| Java          | [docker-client](https://github.com/spotify/docker-client) |
| Java          | [docker-java](https://github.com/docker-java/docker-java) |
| NodeJS        | [dockerode](https://github.com/apocas/dockerode) |
| Perl          | [Eixo::Docker](https://github.com/alambike/eixo-docker) |
| PHP           | [Docker-PHP](https://github.com/docker-php/docker-php) |
| Ruby          | [docker-api](https://github.com/swipely/docker-api) |
| Rust          | [docker-rust](https://github.com/abh1nav/docker-rust) |
| Rust          | [shiplift](https://github.com/softprops/shiplift) |
| Scala         | [tugboat](https://github.com/softprops/tugboat) |
| Scala         | [reactive-docker](https://github.com/almoehi/reactive-docker) |
