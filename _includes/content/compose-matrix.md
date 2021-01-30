This table shows which Compose file versions support specific Docker releases.

| **Compose file format** | **Docker Engine release** |
|  -------------------    |    ------------------     |
|  Compose specification  |       19.03.0+            |
|      3.8                |       19.03.0+            |
|      3.7                |       18.06.0+            |
|      3.6                |       18.02.0+            |
|      3.5                |       17.12.0+            |
|      3.4                |       17.09.0+            |
|      3.3                |       17.06.0+            |
|      3.2                |       17.04.0+            |
|      3.1                |       1.13.1+             |
|      3.0                |       1.13.0+             |
|      2.4                |       17.12.0+            |
|      2.3                |       17.06.0+            |
|      2.2                |       1.13.0+             |
|      2.1                |       1.12.0+             |
|      2.0                |       1.10.0+             |
|      1.0                |       1.9.1.+             |

In addition to Compose file format versions shown in the table, the Compose
itself is on a release schedule, as shown in [Compose
releases](https://github.com/docker/compose/releases/), but file format versions
do not necessarily increment with each release. For example, Compose file format
3.0 was first introduced in [Compose release
1.10.0](https://github.com/docker/compose/releases/tag/1.10.0), and versioned
gradually in subsequent releases.

The latest Compose file format is defined by the [Compose Specification](https://github.com/compose-spec/compose-spec/blob/master/spec.md){:target="_blank" rel="noopener" class="_"} and is implemented by Docker Compose **1.27.0+**.
