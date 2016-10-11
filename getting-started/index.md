---
title: "Getting Started, Part 1: Orientation and Setup"
---

# Getting Started, Part 1: Orientation and Setup

{% include content/docker_elevator_pitch.md %}

## What we'll be covering

This tutorial will create a simple application that runs in a cluster, so you
get a sense of how to build distributed applications with the Docker platform.
We will achieve this in the following steps:

1. Get set up and oriented, on this page.
2. [Create a "Hello World" application that identifies its environment](part2.md)
3. [Hook up a visitor counter](part3.md)
4. [Scale our app as if it were very high traffic, by setting up a cluster in
   production](part4.md)

The application itself is very simple so that you are not too distracted by
what the code is doing. After all, the value of Docker is in how it can build,
ship, and run applications; it's totally agnostic as to what your application
actually does.

## Setup

Before we get started, make sure your system has the latest version of Docker
installed.

[Install Docker](/engine/installation/index.md){: class="button darkblue-btn"}

> Note: If you're in Linux, you'll want to install
  [Docker Toolbox](../toolbox/index.md) so you get Docker Compose.

## Let's go!

If you understand that container images package application code and their
dependencies all together in a portable deliverable, and your environment has
Docker installed, let's move on!

[On to "Getting Started, Part 2: Creating and Building Your App" >>](part2.md){: class="button darkblue-btn"}
