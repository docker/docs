---
title: DTR cache configuration reference
description: Learn about the different configuration options for DTR caches.
keywords: DTR, cache
---

DTR caches are based on Docker Registry, and use the same configuration
file format.
[Learn more about the configuration options](/registry/configuration.md).

The DTR cache extends the Docker Registry configuration file format by
introducing a new middleware called `downstream` that has three configuration
options: `blobttl`, `upstreams`, and `cas`:

```none
# Settings that you would include in a
# Docker Registry configuration file followed by

middleware:
  registry:
      - name: downstream
        options:
          blobttl: 24h
          upstreams:
            - <Externally-reachable address for upstream registry or content cache in format scheme://host:port>
          cas:
            - <Absolute path to next-hop upstream registry or content cache CA certificate in the container's filesystem>
```

Below you can find the description for each parameter, specific to DTR caches.

<table>
  <tr>
    <th>Parameter</th>
    <th>Required</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>
      <code>blobttl</code>
    </td>
    <td>
      no
    </td>
    <td>
The TTL for blobs in the cache. This field takes a positive integer and an optional suffix indicating the unit of time. If
this field is configured, "storage.delete.enabled" must be configured to true. Possible units are:
      <ul>
        <li><code>ns</code> (nanoseconds)</li>
        <li><code>us</code> (microseconds)</li>
        <li><code>ms</code> (milliseconds)</li>
        <li><code>s</code> (seconds)</li>
        <li><code>m</code> (minutes)</li>
        <li><code>h</code> (hours)</li>
      </ul>
    If you omit the suffix, the system interprets the value as nanoseconds.
    </td>
  </tr>
  <tr>
    <td>
      <code>cas</code>
    </td>
    <td>
      no
    </td>
    <td>
      An optional list of absolute paths to PEM-encoded CA certificates of upstream registries or content caches.
    </td>
  </tr>
<tr>
  <td>
    <code>upstreams</code>
  </td>
  <td>
    yes
  </td>
  <td>
      A list of externally-reachable addresses for upstream registries of content caches. If more than one host is specified, it will pull from registries in round-robin order.
  </td>
</tr>
</table>
