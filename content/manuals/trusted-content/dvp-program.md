---
description: Learn what the Docker Verified Publisher Program is and how it works
title: Docker Verified Publisher Program
aliases:
- /docker-hub/publish/publish/
- /docker-hub/publish/customer_faq/
- /docker-hub/publish/publisher_faq/
- /docker-hub/publish/certify-images/
- /docker-hub/publish/certify-plugins-logging/
- /docker-hub/publish/trustchain/
- /docker-hub/publish/byol/
- /docker-hub/publish/publisher-center-migration/
- /docker-hub/publish/
- /docker-hub/publish/repository-logos/
- /docker-hub/dvp-program/
---

[The Docker Verified Publisher Program](https://hub.docker.com/search?q=&image_filter=store) provides high-quality images from commercial publishers verified by Docker.

These images help development teams build secure software supply chains, minimizing exposure to malicious content early in the process to save time and money later.

Images that are part of this program have a special badge on Docker Hub making it easier for users to identify projects that Docker has verified as high-quality commercial publishers.

![Docker-Sponsored Open Source badge](./images/verified-publisher-badge-iso.png)

## For content publishers

The Docker Verified Publisher Program (DVP) provides several features and benefits to Docker
Hub publishers. The program grants the following perks based on participation tier:

- Repository logo
- Verified publisher badge
- Priority search ranking in Docker Hub
- Insights and analytics
- Vulnerability analysis
- Additional Docker Business seats
- Removal of rate limiting for developers
- Co-marketing opportunities

### Repository logo

DVP organizations can upload custom images for individual repositories on Docker Hub.
This lets you override the default organization-level logo on a per-repository basis.

Only a user with administrative access (owner or team member with administrator permission)
over the repository can change the repository logo.

#### Image requirements

- The supported filetypes for the logo image are JPEG and PNG.
- The minimum allowed image size in pixels is 120×120.
- The maximum allowed image size in pixels is 1000×1000.
- The maximum allowed image file size is 5MB.

#### Set the repository logo

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Go to the page of the repository that you want to change the logo for.
3. Select the upload logo button, represented by a camera icon ({{< inline-image
   src="./images/upload_logo_sm.png" alt="camera icon" >}}) overlaying the
current repository logo.
4. In the dialog that opens, select the PNG image that you want to upload to
   set it as the logo for the repository.

#### Remove the logo

Select the **Clear** button ({{< inline-image src="images/clear_logo_sm.png"
alt="clear button" >}}) to remove a logo.

Removing the logo makes the repository default to using the organization logo, if set, or the following default logo if not.

![Default logo which is a 3D grey cube](images/default_logo_sm.png)

### Verified publisher badge

Images that are part of this program have a badge on Docker Hub making it easier for developers
to identify projects that Docker has verified as high quality publishers and with content they can trust.

![Docker, Inc. org with a verified publisher badge](./images/verified-publisher-badge.png)

### Insights and analytics

The [insights and analytics](/docker-hub/publish/insights-analytics) service provides usage metrics for how
the community uses Docker images, granting insight into user behavior.

The usage metrics show the number of image pulls by tag or by digest, and breakdowns by
geolocation, cloud provider, client, and more.

![The insights and analytics tab on the Docker Hub website](./images/insights-and-analytics-tab.png)

You can select the time span for which you want to view analytics data. You can also export the data in either a summary or raw format.

### Vulnerability analysis

[Docker Scout](/scout/) provides automatic vulnerability analysis
for DVP images published to Docker Hub.
Scanning images ensures that the published content is secure, and proves to
developers that they can trust the image.

You can enable analysis on a per-repository
basis. For more about using this feature, see [Basic vulnerability scanning](/docker-hub/vulnerability-scanning/).

### Who's eligible to become a verified publisher?

Any independent software vendor who distributes software on Docker Hub can join
the Verified Publisher Program. Find out more by heading to the
[Docker Verified Publisher Program](https://www.docker.com/partners/programs) page.
