---
description: Docker Store overview
keywords: Docker, docker, store, purchase images
title: Docker Store overview
---

**The Docker Store is now generally available!**

You can [learn more about publishing](https://success.docker.com/Store),
or [apply to be a publisher](https://store.docker.com/publisher/signup).

-----------------

## What is Docker Store?

### Get software distributed as Docker images

For developers and IT, the Docker Store is the place to find the best
trusted commercial and free software distributed as Docker Images.

### Publish, distribute, and sell your Dockerized content

For publishers, Docker Store is the best way for you to distribute and sell your
Dockerized content. Publish your software through the Docker Store to experience
the benefits below:

* Access to Docker’s large and growing customer-base. Docker has
experienced rapid adoption, and is wildly popular in dev-ops
environments. Docker users have pulled images over four billion times,
and they are increasingly turning to the Docker Store as the canonical
source for high-quality, curated content.

* Customers can try or buy your software, right from your product listing. Your
content is accessible for installation, trial, and purchase from the Docker
Store and the Docker CLI.

* Use our licensing support.  We can limit access to your software to
a) logged-in users, b) users who have purchased a license, or
c) all Docker users.  We’ll help you manage and control your distribution.

* We'll handle checkout.  You don’t have to set up your own digital e-commerce
site when you sell your content through the Docker Store.  We'll even help you
set pricing—and you can forget about the rest.

* Seamless updates and upgrades.  We tell customers when your content
has upgrades or updates available, right inside their Docker host product.

* It’s a win-win for our platform and publishers: great content improves our
ecosystem, and our flexible platform helps you bring your content to market.

* Achieve the Docker Certified quality mark.  Publisher container images and
plugins that meet the quality, security, and support criteria of the program
will display a “Docker Certified” badge within the Docker Store and external
marketing.

### Publisher distribution models

The Docker Store welcomes free and open-source content, as well as software sold
directly by publishers.  We support the following commercial models:

#### Paid-via-Docker content

This is content for which customers transact via Docker, as described in the
publisher agreement.  Paid-via-Docker content includes both software than can be
deployed on a host, as well as software that runs in the cloud and can be
accessed by the customer via an ambassador container (containerized cloud
services, for example).

#### Free content

Free content is provided free-of-charge, and customers may pull it from the
Docker Hub either at their discretion or upon license acceptance, at the
publisher’s discretion.  You agree that you will not charge customers for any
Free Content by making it available for purchase outside of the Docker Store.

## FAQs for Docker Certified Publishers

### What is the Docker Certified program?

Docker Certified Container images and plugins are meant to differentiate high
quality content on Docker Store. Customers can consume Certified Containers with
confidence knowing that both Docker and the publisher will stand behind the
solution.  Further details can be found in the [Docker Partner Program Guide](https://www.docker.com/partnerprogramguide){: target="_blank" class="_"}.

### What are the benefits of Docker Certified?

Docker Store will promote Docker Certified Containers and Plugins running on
Docker Certified Infrastructure trusted and high quality content. With over 8B
image pulls and access to Docker’s large customer base, a publisher can
differentiate their content by certifying their images and plugins. With a
revenue share agreement, Docker can be a channel for your content.   The Docker
Certified badge can also be listed alongside external references to your
product.

### How will the Docker Certified Container image be listed on Docker Store?

These images are differentiated from other images on store through a
certification badge. A user can search specifically for CI’s by limiting their
search parameters to show only certified content.

![certified content example](images/FAQ-certified-content.png)

### Is certification optional or required to be listed on Store?

Certification is recommended for most commercial and supported container images.
Free, community, and other commercial (non-certified) content may also be listed
on Docker Store.

![certified content example](images/FAQ-types-of-certified-content.png)

### How will support be handled?

All Docker Certified Container images and plugins running on Docker Certified
Infrastructure come with SLA based support provided by the publisher and Docker.
Normally, a customer contacts the publisher for container and application level
issues.   Likewise, a customer will contact Docker for Docker Edition support.
In the case where a customer calls Docker (or vice versa) about an issue on the
application, Docker will advise the customer about the publisher support process
and will perform a handover directly to the publisher if required.  TSAnet is
required for exchange of support tickets between the publisher and Docker.

### How does a publisher apply to the Docker Certified program?

Start by applying to be a [Docker Technology
Partner](https://goto.docker.com/partners){: target="_blank" class="_"}

Requires acceptance of partnership agreement for completion

Identify commercial content that can be listed on Store and includes a support
offering

Test your image against the Docker CS Engine 1.12+ or on a Docker Certified
Infrastructure version 17.03 and above  (Plugins must run on 17.03 and above)

Submit your image for Certification through the publisher portal. Docker will
scan the image and work with you to address vulnerabilities.  Docker will also
conduct a best practices review of the image.

Be a [TSAnet](https://www.tsanet.org/){: target="_blank" class="_"} member or join the Docker Limited Group.

Upon completion of Certification criteria, and acceptance by Docker, Publisher’s
product page will be updated to reflect Certified status.
Is there a fee to join the program?

In the future, Docker may charge a small annual listing fee. This is waived for
the initial period.

### What is the difference between Official Images and Docker Certified?

Many Official images will transition to the Docker Certified program and will be
maintained and updated by the original owner of the software. Docker will
continue to maintain of some base OS images and language frameworks.  

### How will certification of plugins be handled?

Docker Certification program recognizes the need to apply special scrutiny and
testing to containers that access system level interfaces like storage volumes
and networking.   Docker identifies these special containers as “Plugins” which
require additional testing by the publisher or Docker.  These plugins employ the
V2 Plugin Architecture that was first made available in 1.12 (experimental) and
now available in Docker Enterprise Edition 17.03

## FAQs on using Docker Store

### How do I find content?

Type a search in the search bar. Click one of the suggested matches, or press
`Enter` to run a full search.

![](images/store-search.png)

The search returns any results that match in the image name, description, or
publisher name. If you run a complete search, you can also limit your results by
category.

You can also click **Browse** from the top menu to see all of the images
available in the Store, and filter them by category.

![](images/store-browse.png)

### How do I get a Docker image from the Store?

Once you find an image you want, click **Get Image** to agree to the end-user
agreement, then use the `docker pull` command from the image's Store listing to
download it.

![](images/store-get.png)

Some images may require that you accept the end user agreement or terms of
service before you can pull them, and paid and subscription images may require
that you provide billing information if you have not already done so.

Once you've accepted the terms and provided billing information, you'll see a
link to your list of subscriptions, and the `docker pull` command for the image.
Copy this and paste it into your command shell.

![](images/store-pullcmd.png)

### What types of images are available?

You can download two types of images from the Docker Store:

* **Docker Verified images**. (Recommended) These images are verified
by Docker, have a high level of security, and generally subscribe to
Docker best practices.

* **Community/Hub images**. When you choose this option, you see
images directly from Docker Hub. These images are not verified by Docker.

### What version of an image do I need?

In many cases there will be multiple versions of an image available. Some
versions may offer smaller parent image sizes, or address specific security
vulnerabilities.

To see a list of an image's versions, click **View all versions**.
