---
title: License your installation
description: Learn how to license your Docker Trusted Registry installation.
keywords: dtr, install, license
---

>{% include enterprise_label_shortform.md %}

By default, Docker Trusted Registry (DTR) automatically uses the same license file applied to
your Universal Control Plane (UCP). In the following scenarios, you need to
manually apply a license to your DTR:

* Major version upgrade
* License expiration


## Download your license

Visit Docker Hub's [Enterprise Trial page](https://hub.docker.com/editions/enterprise/docker-ee-trial) to start your one-month trial. After signing up, you should receive a confirmation email with a link to your subscription page. You can find your **License Key** in the **Resources** section of the Docker Enterprise Setup Instructions page.  

![](/ee/dtr/images/license-1.png){: .with-border}

Click "License Key" to download your license.

## License your installation

After downloading your license key, navigate to `https://<dtr-url>` and log in with your credentials.
Select **System** from the left navigation pane, and click *Apply new license* to upload your license
key.

![](/ee/dtr/images/license-2.png){: .with-border}

Within **System > General** under the **License** section, you should see the tier, date of expiration, and ID for your license. 

## Where to go next

- [Use your own TLS certificates](use-your-own-tls-certificates)
- [Enable single sign-on](enable-single-sign-on)
