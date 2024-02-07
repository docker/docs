---
description: Language-specific guides overview
keywords: guides, docker, language, node, java, python, go, golang, .net
title: Language-specific guides overview
toc_min: 1
toc_max: 2
---

The language-specific guides walk you through the process of:
* Containerizing language-specific applications
* Setting up a development environment
* Configuring a CI/CD pipeline
* Deploying an application locally using Kubernetes

In addition to the language-specific modules, Docker documentation also provides guidelines to build images and efficiently manage your development environment. For more information, refer to the following topics:

* [Best practices for writing Dockerfiles](../develop/develop-images/dockerfile_best-practices.md)
* [Docker development best practices](../develop/dev-best-practices.md)
* [Build images with BuildKit](../build/buildkit/index.md#getting-started)
* [Build with Docker](../build/guide/_index.md)

## Language-specific guides

Learn how to containerize your applications and start developing using Docker. Choose one of the following languages to get started.

{{< languages.inline >}}
<div class="grid grid-cols-3 auto-rows-fr sm:flex-col sm:h-auto gap-4">
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/nodejs/"><img class="m-auto rounded" src="{{ (resources.Get "images/language/nodejs.webp").Permalink }}" alt="Develop with Node"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/python/"><img class="m-auto rounded" src="{{ (resources.Get "images/language/python.webp").Permalink }}" alt="Develop with Python"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/java/"><img class="m-auto rounded" src="{{ (resources.Get "images/language/java.webp").Permalink }}" alt="Develop with Java"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/golang/"><img class="m-auto rounded" src="{{ (resources.Get "images/language/golang.webp").Permalink }}" alt="Develop with Go"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/dotnet/"><img class="m-auto rounded" src="{{ (resources.Get "images/language/c-sharp.webp").Permalink }}" alt="Develop with C#"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/rust/"><img class="m-auto rounded" src="{{ (resources.Get "images/language/rust-logo.webp").Permalink }}" alt="Develop with Rust"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/php/"><img class="m-auto rounded" src="{{ (resources.Get "images/language/php-logo.webp").Permalink }}" alt="Develop with PHP"></a>
    </div>
</div>
{{< /languages.inline >}}
