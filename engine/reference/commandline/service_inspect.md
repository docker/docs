---
datafolder: engine-cli
datafile: docker_service_inspect
title: docker service inspect
---
<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->
{% include cli.md %}

## Examples

### Inspecting a service  by name or ID

You can inspect a service, either by its *name*, or *ID*

For example, given the following service;

```bash
$ docker service ls
ID            NAME   MODE        REPLICAS  IMAGE
dmu1ept4cxcf  redis  replicated  3/3       redis:3.0.6
```

Both `docker service inspect redis`, and `docker service inspect dmu1ept4cxcf`
produce the same result:

```bash
$ docker service inspect redis
[
    {
        "ID": "dmu1ept4cxcfe8k8lhtux3ro3",
        "Version": {
            "Index": 12
        },
        "CreatedAt": "2016-06-17T18:44:02.558012087Z",
        "UpdatedAt": "2016-06-17T18:44:02.558012087Z",
        "Spec": {
            "Name": "redis",
            "TaskTemplate": {
                "ContainerSpec": {
                    "Image": "redis:3.0.6"
                },
                "Resources": {
                    "Limits": {},
                    "Reservations": {}
                },
                "RestartPolicy": {
                    "Condition": "any",
                    "MaxAttempts": 0
                },
                "Placement": {}
            },
            "Mode": {
                "Replicated": {
                    "Replicas": 1
                }
            },
            "UpdateConfig": {},
            "EndpointSpec": {
                "Mode": "vip"
            }
        },
        "Endpoint": {
            "Spec": {}
        }
    }
]
```

```bash
$ docker service inspect dmu1ept4cxcf
[
    {
        "ID": "dmu1ept4cxcfe8k8lhtux3ro3",
        "Version": {
            "Index": 12
        },
        ...
    }
]
```

### Inspect a service using pretty-print

You can print the inspect output in a human-readable format instead of the default
JSON output, by using the `--pretty` option:

```bash
$ docker service inspect --pretty frontend
ID:		c8wgl7q4ndfd52ni6qftkvnnp
Name:		frontend
Labels:
 - org.example.projectname=demo-app
Service Mode:	REPLICATED
 Replicas:		5
Placement:
UpdateConfig:
 Parallelism:	0
ContainerSpec:
 Image:		nginx:alpine
Resources:
Endpoint Mode:  vip
Ports:
 Name =
 Protocol = tcp
 TargetPort = 443
 PublishedPort = 4443
```

You can also use `--format pretty` for the same effect.


### Finding the number of tasks running as part of a service

The `--format` option can be used to obtain specific information about a
service. For example, the following command outputs the number of replicas
of the "redis" service.

```bash
{% raw %}
$ docker service inspect --format='{{.Spec.Mode.Replicated.Replicas}}' redis
10
{% endraw %}
```
