---
description: components and formatting examples used in Docker's docs
title: Accordions
toc_max: 3
---

## Example

<i class="fa fa-caret-down" aria-hidden="true"></i> (fa-caret-down) and
<i class="fa fa-caret-up" aria-hidden="true"></i> (fa-caret-up).

<div class="panel panel-default">
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample1" style="cursor: pointer">
    Docker hello-world example
    <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample1">
<pre>
<code>
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
2db29710123e: Pull complete
Digest: sha256:cc15c5b292d8525effc0f89cb299f1804f3a725c8d05e158653a563f15e4f685
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:

1. The Docker client contacted the Docker daemon.
2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
3. The Docker daemon created a new container from that image which runs the
4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:

$ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
</code>
</pre>
</div>

<div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample2"  style="cursor: pointer"> Another Sample <i class="chevron fa fa-fw"></i></div>
<div class="collapse block" id="collapseSample2">
<p>
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.</p>
 </div>
</div>

## HTML

```
<div class="panel panel-default">
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample1" style="cursor: pointer">
    Docker hello-world example
    <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample1">
<pre>
<code>
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
2db29710123e: Pull complete
Digest: sha256:cc15c5b292d8525effc0f89cb299f1804f3a725c8d05e158653a563f15e4f685
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:

1. The Docker client contacted the Docker daemon.
2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
3. The Docker daemon created a new container from that image which runs the
4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:

$ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
</code>
</pre>
</div>

<div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample2"  style="cursor: pointer"> Another Sample <i class="chevron fa fa-fw"></i></div>
<div class="collapse block" id="collapseSample2">
<p>
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.</p>
 </div>
</div>
```

This implementation makes use of the `.panel-heading` classes in
`_utilities.scss.md`,
along with [FontAwesome icons](http://fontawesome.io/cheatsheet/){: target="_blank" rel="noopener" class="_" }

> Note
>
>Make sure `data-target`'s and `id`'s match, and are unique.
>
>For each drop-down, the value for `data-target` and
`collapse` `id` must match, and id's must be unique per page. In this example,
we name these `collapseSample1` and `collapseSample2`.

Adding `block` to the `div` class `collapse` gives you some padding around the
sample content. This works nicely for standard text. If you have a code sample,
the padding renders as white space around the code block grey background.

The `style="cursor: pointer"` tag enables the expand/collapse functionality to
work on mobile. (You can use the [Xcode iPhone simulator](https://developer.apple.com/library/content/documentation/IDEs/Conceptual/iOS_Simulator_Guide/GettingStartedwithiOSSimulator/GettingStartedwithiOSSimulator.html#//apple_ref/doc/uid/TP40012848-CH5-SW4){: target="_blank" rel="noopener" class="_" } to test on mobile.)