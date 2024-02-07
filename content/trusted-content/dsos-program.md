---
description: Learn about the Docker-Sponsored Open Source Program and how it works
title: Docker-Sponsored Open Source Program
keywords: docker hub, hub, insights, analytics, open source, Docker sponsored, program
aliases:
  - /docker-hub/dsos-program/
---

[Docker Sponsored Open Source images](https://hub.docker.com/search?q=&image_filter=open_source) are published and maintained by open-source projects sponsored by Docker through the program.

Images that are part of this program have a special badge on Docker Hub making it easier for users to identify projects that Docker has verified as trusted, secure, and active open-source projects.

![Docker-Sponsored Open Source badge](images/trusted-content/sponsored-badge-iso.png)

## For content publishers

The Docker-Sponsored Open Source (DSOS) Program provides several features and benefits to non-commercial open source developers.

The program grants the following perks to eligible projects:

- Repository logo
- Verified Docker-Sponsored Open Source badge
- Insights and analytics
- Access to [Docker Scout](#docker-scout) for software supply chain management
- Removal of rate limiting for developers
- Improved discoverability on Docker Hub

These benefits are valid for one year and publishers can renew annually if the project still meets the program requirements. Program members and all users pulling public images from the project namespace get access to unlimited pulls and unlimited egress.

### Repository logo

DSOS organizations can upload custom images for individual repositories on Docker Hub.
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
3. Select the upload logo button, represented by a camera icon
   ({{< inline-image src="images/trusted-content/upload_logo_sm.png" alt="camera icon" >}})
   overlaying the current repository logo.
4. In the dialog that opens, select the PNG image that you want to upload to
   set it as the logo for the repository.

#### Remove the logo

Select the **Clear** button ({{< inline-image src="images/trusted-content/clear_logo_sm.png"
alt="clear button" >}}) to remove a logo.

Removing the logo makes the repository default to using the organization logo, if set, or the following default logo if not.

![Default logo which is a 3D grey cube](images/trusted-content/default_logo_sm.png)

### Verified Docker-Sponsored Open Source badge

Docker verifies that developers can trust images with this badge on Docker Hub as an active open source project.

![Fluent org with a Docker-Sponsored Open Source badge](images/trusted-content/sponsored-badge.png)

### Insights and analytics

The [insights and analytics](/docker-hub/publish/insights-analytics) service provides usage metrics for how
the community uses Docker images, granting insight into user behavior.

The usage metrics show the number of image pulls by tag or by digest, and breakdowns by
geolocation, cloud provider, client, and more.

![The insights and analytics tab on the Docker Hub website](images/trusted-content/insights-and-analytics-tab.png)

You can select the time span for which you want to view analytics data. You can also export the data in either a summary or raw format.

### Docker Scout

DSOS projects can enable Docker Scout on up to 100 repositories for free. Docker
Scout provides automatic image analysis, policy evaluation for improved supply
chain management, integrations with third-party systems like CI platforms and
source code management, and more.

You can enable Docker Scout on a per-repository basis. For information about
how to use this product, refer to the [Docker Scout documentation](/scout/).

### Who's eligible for the Docker-Sponsored Open Source program?

To qualify for the program, a publisher must share the project namespace in public repositories, meet [the Open Source Initiative definition](https://opensource.org/docs/osd), and be in active development with no pathway to commercialization.

Find out more by heading to the
[Docker-Sponsored Open Source Program](https://www.docker.com/community/open-source/application/) application page.
