There are several versions of the Compose file format – 1, 2, 2.1 and 3. For
details on versions and how to upgrade, see
[Versioning](compose-versioning.md#versioning) and
[Upgrading](compose-versioning.md#upgrading).

This table shows which Compose file versions support specific Docker releases.

| **Compose file format** | **Docker Engine release** |
|  -------------------    |    ------------------     |
|      3.3                |       17.06.0+            |
|      3.0 - 3.2          |       1.13.0+             |
|      2.1                |       1.12.0+             |
|      2.0                |       1.10.0+             |
|      1.0                |       1.9.1.+             |

In addition to Compose file format versions shown in the table, the Compose
product itself is on a release schedule, as shown in [Compose
releases](https://github.com/docker/compose/releases/), but file format versions
do not necessairly increment with each release. For example, Compose file format
3.0 was introduced in [Compose release
1.10.0](https://github.com/docker/compose/releases/tag/1.10.0), and versioned
gradually in subsequent releases.

> Looking for more detail on Docker and Compose compatibility?
>
> We recommend keeping up-to-date with newer releases as much as possible.
However, if you are using an older version of Docker and want to determine which
Compose release is compatible, please refer to the [Compose release
notes](https://github.com/docker/compose/releases/). Each set of release notes
gives finer-tuned detail on which versions of Docker Engine are supported, along
with compatible Compose file format versions. (See also, the discussion in
[issue #3404](Compatibility matrix between docker-compose and docker versions)
on GitHub.)
{: .note-vanilla}
