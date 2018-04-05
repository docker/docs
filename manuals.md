---
title: Руководства по продуктам и инструментам
notoc: true
---

The Docker platform is comprised of a family of tools and products. After
learning the general principles of the Docker workflow under [Guides](/), you
can find the documentation for these tools and products here.

Платформа Docker состоит из семейства инструментов и продуктов. После изучения общего принципа рабочего процесса Докера в разделе [Guides](/), вы можете найти документацию для этих инструментов и продуктов здесь.

## Поддерживаемые платформы

Docker CE и EE доступны на нескольких платформах, в облаках и на местах. 
Используйте следующие таблицы, чтобы выбрать лучший путь установки для вас.

### Рабочий стол

{% include docker_desktop_matrix.md %}

### Облако

{% include docker_cloud_matrix.md %}

Смотрите также [Docker Cloud](#docker-cloud) для инструкций по установке
Digital Ocean, Packet, SoftLayer, или Bring Your Own Cloud.

### Сервер

{% include docker_platform_matrix.md %}

## Инструменты

Бесплатные загружаемые материалы, которые помогают вашему устройству использовать контейнеры Docker.

| Инструмент                                | Описание                                                                                                |
|:------------------------------------------|:-------------------------------------------------------------------------------------------------------|
| [Docker Compose](/compose/overview/)      | Позволяет определять, создавать и запускать многоконтейнерные приложения                                |
| [Docker Machine](/machine/overview/)      | Предоставляет управление докционированными хостами                                                      |
| [Docker Notary](/notary/getting_started/) | Позволяет подписывать контейнеры для включения Trust Docker Content Trust                              |
| [Docker Registry](/registry/)             | Прог-обеспечение Docker Hub и Docker Store хранит и распространяет реестр контейнеров
             |

## Продукты

Коммерческие продукты Docker, которые превращают ваше контейнерное решение в 
готовое к производству.

| Продукт                                                     | Описание                                                                                                       |
|:-------------------------------------------------------------|:------------------------------------------------------------------------------------------------------------------|
| [Docker Cloud](/docker-cloud/)                               | Manages multi-container applications and host resources running on a cloud provider (such as Amazon Web Services) |
| [Universal Control Plane (UCP)](/datacenter/ucp/2.2/guides/) | Manages your Docker swarm on-premise, or on the cloud                                                             |
| [Docker Trusted Registry (DTR)](/datacenter/dtr/2.3/guides/) | Securely stores and scans your Docker images                                                                      |
| [Docker Store](/docker-store/)                               | Public, Docker-hosted registry that distributes free and paid images from various publishers                      |

## Заменённые продукты и инструменты

* [Docker Hub](/docker-hub/) - Superseded by Docker Store and Docker Cloud
* [Docker Swarm](/swarm/overview/) - Functionality folded directly into native Docker, no longer a standalone tool
* [Docker Toolbox](/toolbox/overview/) - Superseded by Docker for Mac and Windows
