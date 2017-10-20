---
datafolder: engine-cli
datafile: docker_info
title: docker info
---
<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/cli
-->
{% if page.datafolder contains '-edge' %}
  {% include edge_only.md section="cliref" %}
{% endif %}
{% include cli.md datafolder=page.datafolder datafile=page.datafile %}

## Warnings about kernel support

If your operating system does not enable certain capabilities, you may see
warnings such as one of the following, when you run `docker info`:

```none
WARNING: Your kernel does not support swap limit capabilities. Limitation discarded.
```

```none
WARNING: No swap limit support
```

You can ignore these warnings unless you actually need the ability to
[limit these resources](/engine/admin/resource_constraints.md), in which case you
should consult your operating system's documentation for enabling them.
[Learn more](/engine/installation/linux/linux-postinstall.md#your-kernel-does-not-support-cgroup-swap-limit-capabilities).
