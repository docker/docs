---
datafolder: engine-cli
datafile: docker_secret_inspect
title: docker secret inspect
---
<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->
{% include cli.md %}

## Examples

### Inspecting a secret by name or ID

You can inspect a secret, either by its *name*, or *ID*

For example, given the following secret:

```bash
$ docker secret ls
ID                          NAME                    CREATED                                   UPDATED
mhv17xfe3gh6xc4rij5orpfds   secret.json             2016-10-27 23:25:43.909181089 +0000 UTC   2016-10-27 23:25:43.909181089 +0000 UTC
```

```bash
$ docker secret inspect secret.json
[
    {
        "ID": "mhv17xfe3gh6xc4rij5orpfds",
            "Version": {
            "Index": 1198
        },
        "CreatedAt": "2016-10-27T23:25:43.909181089Z",
        "UpdatedAt": "2016-10-27T23:25:43.909181089Z",
        "Spec": {
            "Name": "secret.json"
        }
    }
]
```

### Formatting secret output

You can use the --format option to obtain specific information about a
secret. The following example command outputs the creation time of the
secret.

```bash
{% raw %}
$ docker secret inspect --format='{{.CreatedAt}}' mhv17xfe3gh6xc4rij5orpfds
2016-10-27 23:25:43.909181089 +0000 UTC
{% endraw %}
```
