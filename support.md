+++
title ="Work with Docker Support"
description="Docker Universal Control Plane"
[menu.main]
weight="2"
+++


# Work with Docker Support

Your UCP purchase includes a contract with Docker Support. Depending on your request, Support may ask you to provide a data transfer or "dump" of information from you system. UCP supports generating support dumps across the entire Swarm cluster. The dump operation use the public `dsinfo` image developed by Docker Support.

To generate a dump for your cluster:

1. Log into as a user UCP with the administrator privileges.

2. Click **Admin** > **Support Dump**.

    UCP builds a `docker-support-<datestamp>-<timestap>.zip` file and downloads it to your browser's default download folder.
