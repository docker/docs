---
description: Migrate content from the Publisher Center to Docker Hub
keywords: publish, store, publisher center, verified image,
title: Migrate content from the Publisher Center to Docker Hub
redirect_from:
- /docker-hub/publish/store-migration/
---

The [Publisher Center](https://hub.docker.com/publisher/center){: target="_blank" rel="noopener" class="_"} (previously known as Docker Store) was used by Docker trusted partners until recently to publish their verified content on Docker Hub. Docker Hub and the Publisher Center serve similar purposes. Our new streamlined model uses public repositories as the universal mechanism for image publishing, and our goal is to create a unified and seamless experience for our users.

Additionally, we’ve moved away from a payment-for-images model and are instead focusing on improving the quality of the content provided on Docker Hub through our Trusted Content programs such as [Docker Official Images](../official_images.md), [Docker Verified Publishers](../publish/index.md), and the [Docker Open Source Program](https://www.docker.com/community/open-source/application/){: target="_blank" rel="noopener" class="_"}. Therefore, Docker has decided to retire the Publisher Center in stages and invest resources to improve Docker Hub.

Docker will issue brown-out periods before fully retiring the Publisher Center. The first brown-out period is on **April 21 2022 4:00 pm - 10:00 pm UTC**. You can confirm if you are affected during the brown-out periods by sending us an email at `publisher-support@docker.com`.

## What’s happening to the content in the Publisher Center?

Content in the Publisher Center (Docker Store) is currently placed in a read-only mode, and is no longer shown with a Verified Publisher badge in Docker Hub. This means, users will no longer be able to check out, subscribe to products, or pull images that were previously pulled from the Publisher Center.

A majority of publishers have already migrated their images from the  Publisher Center into Docker Hub. These images will have no break in service and users will be able to access the content through Docker Hub going forward.

All content in the Publisher Center will be eventually removed from Docker Hub. We recommend that you log into the Publisher Center and verify whether any of the images you are using were pulled from a Docker Store image. If so, contact the publisher for an alternative to those images.

If you are a publisher, we recommend that you migrate content to public repositories owned by your organization as soon as possible to retain your existing repositories and to continue publishing your content through Docker Hub.

## How do I migrate my content as a publisher?

The following sections provide recommendations for publishers to migrate content from the Publisher Center to Docker Hub:

### Download data from the Publisher Center

1. Take a snapshot of the data you’d like to retain. If there are any reports that you’d like to keep, log into your Docker Hub account as an owner and  click [Publisher Center](https://hub.docker.com/publisher/center){: target="_blank" rel="noopener" class="_"}.
2. Select **Analytics**. You’ll now see the available data for your product. Choose the timeframe you need and capture that data before you migrate.
3. Complete the following steps for each repository that you own:

   1. Ensure you pull an image for each tag and in every repository as a backup.
   2. Take a screenshot and copy the suggested docker pull syntax, the image name and tags, and any description/terms you’d like to maintain.
   3. Ensure you’ve saved your publisher logo and product logos for future use.
   4. After you’ve completed all the steps above, ensure you have all the information that you need for each repository. This helps you to audit your images, assess the usage, and finalize the versions that you’d like to publish on Docker Hub.

### Create your new Docker organization

1. Create a new namespace/organization in [Docker Hub](https://hub.docker.com/orgs){: target="_blank" rel="noopener" class="_"}.
For step-by-step instructions, see [Teams and Organizations](../orgs.md). For information about creating repositories, see [Repositories](../repos.md).
2. As part of creating a new organization, we recommend that you create a secure group alias or email login as a backup administrator login.
3. When you create a new organization, we recommend that you use a name similar to the one you used in the Publisher Center.
4. Edit your new Docker Hub namespace settings to include the Gravatar image `https://hub.docker.com/orgs/<namespace>/settings/general`. You may need to create a new Gravatar based on a logo you saved earlier.

> **Note:**
>
> The image that you add to your namespace appears on all images and repositories.

### Upload all your images

1. Create individual repositories in the organization. For more information, see [Create new repositories](../repos.md#creating-repositories). Note that the repository descriptions are more streamlined in Docker Hub than in the Publisher Center. Also, the Buy Now and Licensing options do not exist anymore as these features have been deprecated.
2. In many cases, publishers choose to refer to and link to the licensing terms in the repository description to clarify what the users may be agreeing to. You may want to consider creating the individual repositories as a Private repository so you can test the settings before making it available externally.
3. Push each image and tag to the corresponding repository.

### Join the Docker Verified Publisher program

Services offered through the Docker Verified Publishers (DVP) program replace the Publisher Center.

After you’ve pushed the required images to each repository in your organization, let us know by sending us an email at `publisher-support@docker.com`. We will work with you to onboard you to your DVP program.

Docker Verified Publisher partners enjoy benefits such as:

- Removal of rate limiting on all repos in the DVP partners’ organization, providing a premium user experience. This means all Docker users, irrespective of their Docker subscription status, have access to unlimited pulls for partner’s images.
- **Docker Verified Publisher** badging on partner organizations and repos, indicating the trusted content and verified source. This is part of Docker’s Secure Software Supply Chain initiative.
- Priority search ranking in Docker Hub.
- Inclusion as one of two trusted sources in the image access controls included in the Docker Business subscription tier, bringing essential security and management capabilities to larger Docker customers.
- Regular reporting to track key partner repo metrics such as pull requests, unique IP addresses, and more.
- Co-marketing opportunities including social shares, amplifying blog posts, the exclusive right to sponsor DockerCon, etc.
- And more!

[Contact us](https://www.docker.com/partners/programs/){: target="_blank" rel="noopener" class="_"} to learn more about the Docker Verified Publisher program.

## Support

If you have any questions or concerns, contact us at `publisher-support@docker.com`.
