---
title: About Docker CE
description: Lists the installation methods
keywords: docker, installation, install, docker ce, docker ee, docker editions, stable, edge
redirect_from:
- /installation/
- /engine/installation/linux/
- /engine/installation/linux/frugalware/
- /engine/installation/frugalware/
- /engine/installation/linux/other/
- /engine/installation/linux/archlinux/
- /engine/installation/linux/cruxlinux/
- /engine/installation/linux/gentoolinux/
- /engine/installation/linux/docker-ce/
- /engine/installation/linux/docker-ee/
- /engine/installation/
---

Docker Community Edition (CE) is ideal for developers and small
teams looking to get started with Docker and experimenting with container-based
apps. Docker CE has two update channels, **stable** and **edge**:

* **Stable** gives you reliable updates every quarter.
* **Edge** gives you new features every month.

For more information about Docker CE, see
[Docker Community Edition](https://www.docker.com/community-edition/){: target="_blank" class="_" }.

## Supported platforms

Docker CE is available on multiple platforms, on cloud and on-premises.
Use the following tables to choose the best installation path for you.

### Desktop

{% assign green-check = '![yes](/install/images/green-check.svg){: style="height: 14px; margin: 0 auto"}' %}

| Platform                                                                    |      x86_64       |
|:----------------------------------------------------------------------------|:-----------------:|
| [Docker for Mac (macOS)](/docker-for-mac/install.md)                        | {{ green-check }} |
| [Docker for Windows (Microsoft Windows 10)](/docker-for-windows/install.md) | {{ green-check }} |


### Cloud

{% assign green-check = '![yes](/install/images/green-check.svg){: style="height: 14px; margin: 0 auto"}' %}


| Platform                                | Docker Community Edition |
|:----------------------------------------|:------------------------:|
| [Amazon Web Services](/docker-for-aws/) |    {{ green-check }}     |
| [Microsoft Azure](/docker-for-azure/)   |    {{ green-check }}     |


### Server

{% assign green-check = '![yes](/install/images/green-check.svg){: style="height: 14px; margin: 0 auto"}' %}
{% assign install-prefix-ce = '/install/linux/docker-ce' %}

| Platform                                    | x86_64 / amd64                                         | ARM                                                    | ARM64 / AARCH64                                        | IBM Power (ppc64le)                                    | IBM Z (s390x)                                          |
|:--------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|
| [CentOS]({{ install-prefix-ce }}/centos.md) | [{{ green-check }}]({{ install-prefix-ce }}/centos.md) |                                                        | [{{ green-check }}]({{ install-prefix-ce }}/centos.md) |                                                        |                                                        |
| [Debian]({{ install-prefix-ce }}/debian.md) | [{{ green-check }}]({{ install-prefix-ce }}/debian.md) | [{{ green-check }}]({{ install-prefix-ce }}/debian.md) | [{{ green-check }}]({{ install-prefix-ce }}/debian.md) |                                                        |                                                        |
| [Fedora]({{ install-prefix-ce }}/fedora.md) | [{{ green-check }}]({{ install-prefix-ce }}/fedora.md) |                                                        |                                                        |                                                        |                                                        |
| [Ubuntu]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) |


## Time-based release schedule

Starting with Docker 17.03, Docker uses a time-based release schedule.

- Docker CE Edge releases generally happen monthly.
- Docker CE Stable releases generally happen quarterly, with patch releases as
  needed.

### Updates, and patches

- A given Docker CE Stable release receives patches and updates for **one
  month after the next Docker CE Stable release**.
- A given Docker CE Edge release does not receive any patches or updates after
  a subsequent Docker CE Edge or Stable release.

## Get started

After setting up Docker, you can learn the basics over at
[Getting started with Docker](/get-started/).
