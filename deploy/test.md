---
title: Manage logs
description: |
  The reason you would do this is X, Y, and Z.

  This can be a multiline description but should probably `be brief`.
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
cli_tabs:
- version: docker-cli-linux
- version: docker-cli-win
- version: kubectl
next_steps:
- path: /engine/install
  title: Install Docker
- path: /get-started/
  title: Get Started with Docker
---
{% if include.ui %}
To do this foobar task, you'll want to flip the switch under **Tasks > Foobar**,
enter your Lorem Ipsum value for {{ site.tablabels[tab.version] }}, then
click **Save**.

{% if include.version=="ucp-3.0" %}
![Image number 1](https://docs.docker.com/datacenter/ucp/2.2/guides/images/monitor-ucp-0.png)
{% elsif include.version=="ucp-2.2" %}
![Image number 2](https://docs.docker.com/datacenter/ucp/2.2/guides/images/monitor-ucp-1.png)
{% endif %}
{% endif %}

{% if include.cli %}
The command line workflow is essentially the same across the various CLIs.
First you enumerate the services on the node of choice, then you run the
`foobar` command.

{% if include.version=="docker-cli-linux" %}
```bash
$ docker stack deploy -c test.yml smokestack
```
{% elsif include.version=="docker-cli-win" %}
```powershell
docker stack deploy -c test.yml smokestack
```
{% elsif include.version=="kubectl" %}
```bash
$ kubectl get pod -f ./pod.yaml
```
{% endif %}
{% endif %}
