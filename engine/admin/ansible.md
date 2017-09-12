---
description: Installation and using Docker via Ansible
keywords: ansible, installation, usage, docker, documentation
title: Use Ansible
---

> **Note**:
> Please note this is a community contributed installation path.

## Requirements

To use this guide you'll need a working installation of
[Ansible](https://www.ansible.com/) version 2.1.0 or later.

Requirements on the host that will execute the module:

```
python >= 2.6
docker-py >= 1.7.0
Docker API >= 1.20
```

## Installation

The `docker_container` module is a core module, and will ship with
Ansible by default.

## Usage

Task example that pulls the latest version of the `nginx` image and
runs a container. Bind address and ports are in the example defined
as [a variable](https://docs.ansible.com/ansible/playbooks_variables.html).

```
---
- name: nginx container
  docker:
    name: nginx
    image: nginx
    state: reloaded
    ports:
    - "{{ nginx_bind_address }}:{{ nginx_port }}:{{ nginx_port }}"
    cap_drop: all
    cap_add:
      - setgid
      - setuid
    pull: always
    restart_policy: on-failure
    restart_policy_retry: 3
    volumes:
      - /some/nginx.conf:/etc/nginx/nginx.conf:ro
  tags:
    - docker_container
    - nginx
...
```

## Documentation

The documentation for the `docker_container` module is present at
[docs.ansible.com](https://docs.ansible.com/ansible/docker_container_module.html).

Documentation covering Docker images, networks, and services is also present
at [docs.ansible.com](https://docs.ansible.com/ansible/list_of_cloud_modules.html#docker).
