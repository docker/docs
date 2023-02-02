---
description: Known issues for Mac
keywords: mac, troubleshooting, known issues, Docker Desktop
title: Known issues for Docker Desktop on Mac
---
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab3">Known issues for Mac with Intel chip</a></li>
  <li><a data-toggle="tab" data-target="#tab4">Known issues for Mac with Apple silicon</a></li>
</ul>
<div class="tab-content">
<div id="tab3" class="tab-pane fade in active" markdown="1">
<br>
* IPv6 is not (yet) supported on Docker Desktop.

* On Apple silicon in native `arm64` containers, older versions of `libssl` such as `debian:buster`, `ubuntu:20.04`, and `centos:8` will segfault when connected to some TLS servers, for example, `curl https://dl.yarnpkg.com`. The bug is fixed in newer versions of `libssl` in `debian:bullseye`, `ubuntu:21.04`, and `fedora:35`.

* The Mac Activity Monitor reports that Docker is using twice the amount of memory than it is actually using. This is due to a bug in MacOS. We have written [a detailed report](https://docs.google.com/document/d/17ZiQC1Tp9iH320K-uqVLyiJmk4DHJ3c4zgQetJiKYQM/edit?usp=sharing){:target="_blank" rel="noopener" class="_"} on this.

* You might encounter errors when using `docker-compose up` with Docker Desktop
  (`ValueError: Extra Data`). We've identified this is likely related to data
  and/or events being passed all at once rather than one by one, so sometimes
  the data comes back as 2+ objects concatenated and causes an error.

* Force-ejecting the `.dmg` after running `Docker.app` from it can cause the
  whale icon to become unresponsive, Docker tasks to show as not responding in
  the Activity Monitor, and for some processes to consume a large amount of CPU
  resources. Reboot and restart Docker to resolve these issues.

* Docker does not auto-start on login even when it is enabled in **Settings**. This is related to a
  set of issues with Docker helper, registration, and versioning.

* Docker Desktop uses the `HyperKit` hypervisor
  (https://github.com/docker/hyperkit){:target="_blank" rel="noopener" class="_"} in macOS 10.10 Yosemite and higher. If
  you are developing with tools that have conflicts with `HyperKit`, such as
  [Intel Hardware Accelerated Execution Manager
  (HAXM)](https://software.intel.com/en-us/android/articles/intel-hardware-accelerated-execution-manager/){:target="_blank" rel="noopener" class="_"},
  the current workaround is not to run them at the same time. You can pause
  `HyperKit` by quitting Docker Desktop temporarily while you work with HAXM.
  This allows you to continue work with the other tools and prevent `HyperKit`
  from interfering.

* If you are working with applications like [Apache
  Maven](https://maven.apache.org/){:target="_blank" rel="noopener" class="_"} that expect settings for `DOCKER_HOST` and
  `DOCKER_CERT_PATH` environment variables, specify these to connect to Docker
  instances through Unix sockets. For example:

  ```console
  $ export DOCKER_HOST=unix:///var/run/docker.sock
  ```

* <a name="bind-mounted-dirs"></a> There are a number of issues with the performance of directories bind-mounted
  into containers. In particular, writes of small blocks, and traversals of large
  directories are currently slow. Additionally, containers that perform large
  numbers of directory operations, such as repeated scans of large directory
  trees, may suffer from poor performance. Applications that behave in this way
  include:

  - `rake`
  - `ember build`
  - Symfony
  - Magento
  - Zend Framework
  - PHP applications that use [Composer](https://getcomposer.org){:target="_blank" rel="noopener" class="_"} to install
    dependencies in a `vendor` folder

  As a workaround for this behavior, you can put vendor or third-party library
  directories in Docker volumes, perform temporary file system operations
  outside of bind mounts, and use third-party tools like Unison or `rsync` to
  synchronize between container directories and bind-mounted directories. We are
  actively working on performance improvements using a number of different
  techniques.  To learn more, see the [topic on our roadmap](https://github.com/docker/roadmap/issues/7){: target="_blank" rel="noopener" class="_" }.

<hr>
</div>  
<div id="tab4" class="tab-pane fade" markdown="1">
<br>
- Some command line tools do not work when Rosetta 2 is not installed.
  - The old version 1.x of `docker-compose`. We recommend that you use Compose V2 instead. Either type `docker compose` or enable the **Use Docker Compose V2** option in the [General preferences tab](../settings/mac.md#general).
  - The `docker scan` command and the underlying `snyk` binary.
  - The `docker-credential-ecr-login` credential helper.
- Some images do not support the ARM64 architecture. You can add `--platform linux/amd64` to run (or build) an Intel image using emulation.

   However, attempts to run Intel-based containers on Apple silicon machines under emulation can crash as qemu sometimes fails to run the container. In addition, filesystem change notification APIs (`inotify`) do not work under qemu emulation. Even when the containers do run correctly under emulation, they will be slower and use more memory than the native equivalent.

   In summary, running Intel-based containers on Arm-based machines should be regarded as "best effort" only. We recommend running arm64 containers on Apple silicon machines whenever possible, and encouraging container authors to produce arm64, or multi-arch, versions of their containers. We expect this issue to become less common over time, as more and more images are rebuilt [supporting multiple architectures](https://www.docker.com/blog/multi-arch-build-and-images-the-simple-way/){:target="_blank" rel="noopener" class="_"}.
- `ping` from inside a container to the Internet does not work as expected.  To test the network, we recommend using `curl` or `wget`. See [docker/for-mac#5322](https://github.com/docker/for-mac/issues/5322#issuecomment-809392861){:target="_blank" rel="noopener" class="_"}.
- Users may occasionally experience data drop when a TCP stream is half-closed.
<hr>
</div>
</div>