---
description: Docker-Sponsored Open Source Program
title: Docker-Sponsored Open Source Program
keywords: docker hub, hub, insights, analytics, open source, Docker sponsored, program
---

[Docker Sponsored Open Source images](https://hub.docker.com/search?q=&image_filter=open_source){:target="_blank" rel="noopener" class="_"} are published and maintained by open-source projects sponsored by Docker through the program.

Images that are part of this program have a special badge on Docker Hub making it easier for users to identify projects that Docker has verified as trusted, secure, and active open-source projects.

![Docker-Sponsored Open Source badge](images/sponsored-badge-iso.png)

## For content publishers

The Docker-Sponsored Open Source (DSOS) Program provides several features and benefits to non-commercial open source developers.

The program grants the following perks to eligible projects:

- Repository logo
- Verified Docker-Sponsored Open Source badge
- Insights and analytics
- Vulnerability analysis
- Removal of rate limiting for developers
- Improved discoverability on Docker Hub

These benefits are valid for one year and publishers can renew annually if the project still meets the program requirements. Program members, and all users pulling public images from the project namespace get access to unlimited pulls and unlimited egress.

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

1. Sign in to Docker Hub.
2. Go to the page of the repository that you want to change the logo for.
3. Select the upload logo button, represented by a camera icon
   (![Camera icon](images/upload_logo_sm.png){: .inline height="22px" })
   overlaying the current repository logo.
4. In the dialog that opens, select the PNG image that you want to upload to
   set it as the logo for the repository.

#### Remove the logo

Remove a logo using the clear button (![Clear button](images/clear_logo_sm.png){: .inline width="20px" }).

Removing the logo makes the repository fallback to using the organization logo, if set, and the default logo if not.

![Default logo which is a 3D grey cube](images/default_logo_sm.png){: width="50px" }

### Verified Docker-Sponsored Open Source badge

Docker verifies that developers can trust images with this badge on Docker Hub as an active open source project.

![Fluent org with a Docker-Sponsored Open Source badge](images/sponsored-badge.png)

### Insights and analytics

The [insights and analytics](/docker-hub/publish/insights-analytics){:
target="blank" rel="noopener" class=""} service provides usage metrics for how
the community uses Docker images, granting insight into user behavior.

The usage metrics show the number of image pulls by tag or by digest, and breakdowns by
geolocation, cloud provider, client, and more.

![The insights and analytics tab on the Docker Hub website](images/insights-and-analytics-tab.png)

You can use the view to select the time span you want to view analytics data and export the data in
either a summary or raw format.

### Vulnerability analysis

[Docker Scout](/scout/){:
target="blank" rel="noopener" class=""} provides automatic vulnerability analysis
for DVP images published to Docker Hub.
Scanning images ensures that the published content is secure, and proves to
developers that they can trust the image.

Analysis is enabled on a per-repository
basis, refer to [vulnerability scanning](/docker-hub/vulnerability-scanning/){:
target="blank" rel="noopener" class=""} for more information about how to use
it.

> **Note**
>
> Content publishers in the Docker-Sponsored Open Source Program receive 3 free
> Docker Team Seats

### Who's eligible for the Docker-Sponsored Open Source program?

To qualify for the program, a publisher must share the project namespace in public repositories, meet [the Open Source Initiative definition](https://opensource.org/docs/osd), and be in active development with no pathway to commercialization.

Find out more by heading to the
[Docker-Sponsored Open Source Program](https://www.docker.com/community/open-source/application/#){:target="_blank"
rel="noopener" class="_"} application page.
