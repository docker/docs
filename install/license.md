<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/license/"]
title = "License DTR"
description = "Learn how to license your Docker Trusted Registry installation."
keywords = ["docker, dtr, install, license"]
[menu.main]
parent="workw_dtr_install"
weight=31
+++
<![end-metadata]-->


# License DTR

By default, you don't need to license your Docker Trusted Registry. When
installing DTR, it automatically starts using the same license file used on
your Docker Universal Control Plane cluster.

However, there are some situations when you have to manually license your
DTR installation.:

* During an upgrade to a new major version;
* When your current license expires.


## Download your license

When your new license is issued, you can download it on **Docker Hub**. Navigate
to your **Profile settings**, and click the
[Licenses tab](https://hub.docker.com/account/licenses/).

![](../images/get-license-2.png)


## License your installation

Once you've downloaded the license file, you can apply it to your DTR
installation. Navigate to the **DTR web app**, and then go to the **Settings
page**.

![](../images/license-1.png)

Click the **Apply new license** button, and upload your new license file.


## Where to go next

* [Install DTR](install-dtr.md)
* [Install DTR offline](install-dtr-offline.md)
