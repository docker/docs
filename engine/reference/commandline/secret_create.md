---
datafolder: engine-cli
datafile: docker_secret_create
title: docker secret create
---
<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->
{% include cli.md %}

## Examples

### Create a secret

```bash
$ echo <secret> | docker secret create my_secret -
mhv17xfe3gh6xc4rij5orpfds

$ docker secret ls
ID                          NAME                    CREATED                                   UPDATED                                   SIZE
mhv17xfe3gh6xc4rij5orpfds   my_secret               2016-10-27 23:25:43.909181089 +0000 UTC   2016-10-27 23:25:43.909181089 +0000 UTC   1679
```

### Create a secret with a file

```bash
$ docker secret create my_secret ./secret.json
mhv17xfe3gh6xc4rij5orpfds

$ docker secret ls
ID                          NAME                    CREATED                                   UPDATED                                   SIZE
mhv17xfe3gh6xc4rij5orpfds   my_secret               2016-10-27 23:25:43.909181089 +0000 UTC   2016-10-27 23:25:43.909181089 +0000 UTC   1679
```

### Create a secret with labels

```bash
$ docker secret create --label env=dev --label rev=20161102 my_secret ./secret.json
jtn7g6aukl5ky7nr9gvwafoxh

$ docker secret inspect my_secret
[
    {
        "ID": "jtn7g6aukl5ky7nr9gvwafoxh",
        "Version": {
            "Index": 541
        },
        "CreatedAt": "2016-11-03T20:54:12.924766548Z",
        "UpdatedAt": "2016-11-03T20:54:12.924766548Z",
        "Spec": {
            "Name": "my_secret",
            "Labels": {
                "env": "dev",
                "rev": "20161102"
            },
            "Data": null
        },
        "Digest": "sha256:4212a44b14e94154359569333d3fc6a80f6b9959dfdaff26412f4b2796b1f387",
        "SecretSize": 1679
    }
]

```
