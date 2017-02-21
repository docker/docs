---
description: Your Docker subscription gives you access to prioritized support. You
  can file tickets via email, your the support portal.
keywords: Docker, support, help
title: Get support
product_name: Universal Control Plane
product_version: 2.1
---

Your Docker Data Center, or Universal Control Plane subscription gives you
access to prioritized support. The service levels depend on your subscription.

If you need help, you can file a ticket via:

* [Email](mailto:support@docker.com)
* [Docker support page](https://support.docker.com/)

Be sure to use your company email when filing tickets.

## Download a support dump

Docker Support engineers may ask you to provide a UCP support dump, which is an
archive that contains UCP system logs and diagnostic information. To obtain a
support dump:

1. Log into the UCP UI with an administrator account.

2. On the top-right menu, **click your username**, and choose **Support Dump**.
   An archive will be downloaded by your browser after a brief time interval.

If the user interface is not accessible, you may perform the following number of
steps instead to obtain a single-node version of the support dump:

1. Obtain direct CLI access to the docker daemon on a UCP manager node.

2. Run the CLI support tool with the following command:
	```bash
	$ docker run --rm \
	--name ucp \
	-v /var/run/docker.sock:/var/run/docker.sock \
	{{ page.docker_image }} \
	support > docker-support.tgz
	```
