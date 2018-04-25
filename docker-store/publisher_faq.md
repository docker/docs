---
description: Docker Store frequently asked questions
keywords: Docker, docker, store, purchase images
title: Docker Store Publisher FAQs
---

## Certification program

### What is the certification program for images and plugins, and what are some benefits?

The Docker Certification program for Infrastructure, Images, and Plugins is
designed for both technology partners and enterprise customers to recognize
high-quality Containers and Plugins, provide collaborative support, and ensure
compatibility with Docker EE. Docker Certification is aligned to the available
Docker EE infrastructure and gives enterprises a trusted way to run more
technology in containers with support from both Docker and the publisher. The
[Docker Technology Partner guide](https://www.docker.com/partners/partner-program#/technology_partner)
explains the Technology Partner program and the Docker Certification Program for
Infrastructure, Images, and Plugins in more detail.

## Publisher signup and approval

### How do I get started with the publisher signup and approval process?

Start by applying to be a Docker Technology Partner at https://goto.docker.com/partner and click on "Publisher".

* Requires acceptance of partnership agreement for completion
* Identify content that can be listed on Store and includes a support offering
* Test your image against Docker Certified Infrastructure version 17.03 and
above (Plugins must run on 17.03 and above).
* Submit your image for Certification through the publisher portal. Docker
scans the image and work with you to address vulnerabilities. Docker also
conducts a best practices review of the image.
* Be a TSAnet member or join the Docker Limited Group.
* Upon completion of Certification criteria, and acceptance by Docker,
Publisher’s product page is updated to reflect Certified status.

### What is the Docker Store Publisher Program application timeline?

1-2 weeks.

### Can we have a group of people work on the same product and publish to Store? (This replicates our internal workflow where more than one person is working on Dockerizing our product.)

Yes. You can submit your content as a team.

## Product submission

### How long does it typically take to have an image approved?

2 Weeks.

### Once a product is published, what is the process for pushing a new build (1.2, 1.3)? Will we simply edit the same product, adding the newly tagged repos?

Edit the same product and update with the newly tagged repos.

### On the Information page, organization details are required. Do we need to fill those in again for every product we publish, or are they carried over? And if we change them for a later image publish, are they updated for all images published by our organization?

Organization details need to be filled in only once. Updating organization info
once updates this for all images published by your organization.

### On the page for another vendor’s product on Docker store, I see the following chunks of data: How do these fields map to the following that are required in the publish process?

#### Fields I see

* Description
* License
* Feedback
* Contributing Guidelines
* Documentation

#### Fields in the publish process

* Product description
* Support link
* Documentation link
* Screenshots
* Tier description
* Installation instructions

*Description* maps to *Product description* in the publish process.
*License* maps to *Support Link* in the publish process.
*Documentation* maps to *Documentation Link* in the publish process.
*Feedback* is provided via customer reviews. https://store.docker.com/images/node?tab=reviews is an example.
*Tier Description* is what you see once users get entitled to a plan. For instance, in https://store.docker.com/images/openmaptiles-openstreetmap-maps/plans/f1fc533a-76f0-493a-80a1-4e0a2b38a563?tab=instructions `A detailed street map of any place on a planet. Evaluation and non-production use. Production use license available separately` is what this publisher entered in the Tier description
*Installation instructions* is documentation on installing your software. In this case the documentation is just `Just launch the container and the map is going to be available on port 80 - ready-to-use - with instructions and list of available styles.` (We recommend more details for any content thats a certification candidate).

### How can I remove a submission? I don’t want to currently have this image published as it is missing several information.

If you would like your submission removed, let us know by contacting us at
publisher-support@docker.com.

### Can publishers publish multi-container apps?

Yes. Publishers can provide multiple images and add a compose file in the
install instructions to describe how the multi-container app can be used. For
now, we recommend asking publishers to look at this example from Microsoft
https://store.docker.com/images/mssql-server-linux where they have Supported
Tags listed in the Install instructions (you don't necessarily need to list it
in the readme).

### Regarding source repo tags: it says not to use “latest”. However, if we want users to be able to download the images without specifying a tag, then presumably an image tagged “latest” is required. So how do we go about that?

You can not submit "latest" tags via the certification/store publish workflow.
The reason we do this is so that users are aware of the exact version they
download. To make the user experience easy we have a copy widget that users can
use to copy the pull command and paste in their command line. Here is a
[screenshot](https://user-images.githubusercontent.com/2453622/32354702-1bec633a-bfe8-11e7-9f80-a02c26b1b10c.png)
to provide additional clarity.

### I have two plans, can I use the same repository but different tags for the two plans?

We expect publishers to use a different repository for each plan. If a user is entitled to a plan in your product, the user is entitled to all tags in the relevant. 
For instance, if you have a `Developer` Plan, that is mapped to repositories store/`mynamespace`/`myrepo1`:`mytag1`, another plan (say `Production`) **should** map to a different repository. 
**_Any user that is entitled to the `Developer` plan will be able to pull all tags in store/`mynamespace`/`myrepo1`_**. 

## Licensing, terms and conditions, and pricing

### What options are presented to users to pull an image?

We provide users the following options to access your software
* logged-in users.
* users who have accepted ToS
* all users (including users without Docker Identity)
Here is a [screenshot](https://user-images.githubusercontent.com/2453622/32067299-00cf1210-ba83-11e7-89f8-15deed6fef62.png) to describe how publishers can update the options provided to customers.

### If something is published as a free tier, for subscribed users only, does a user need to explicitly click Accept on the license terms for which we provide the link before they can download the image?
Yes

### Do you have a license enforcement system for docker images sold on store? How are they protected, once they have been downloaded? What happens if a customer stop paying for the image I am selling after, let's say, 2 months?

We provide the following licensing option to customers:
* Bring your own License or BYOL.

The expectation is that the publisher would take care of License Keys within the
container. The License Key itself can be presented to the customer via Docker
Store. We expect the Publisher to build short circuits into the container, so
the container stops running once the License Key expires. Once a customer
cancels, or if the customer subscription expires - the customer cannot
download updates from the Store.

If a user cancels their subscription, they cannot download updates
from the Store. The container may continue running. If you have a licensing
scheme built into the container, the licensing scheme can be a forcing function
and stop the container. (_We do not build anything into the container, it is up to the publisher_).

### How does a customer transition from a Trial to a Paid subscription? Question assumes these are two separate pulls from Store, or can they just drop in a license via Store?

Publisher can provide two different tokens or let customers use the same token
and internally map the customer to a paid plan vs a free trial.

### What are Docker Store pricing plans like? Can I have metered pricing?

As a publisher you can charge a subscription fee every month in USD. The amount
is determined by you. We are working on other pricing options. If you have
feedback about pricing, send us an email at publisher-support@docker.com

### As a publisher, I have not setup any payment account. How does money get to me if my commercial content gets purchased by customers?

We (Docker) cut you a check post a revenue share. Your Docker Store Vendor
Agreement should cover specifics.

### How does Docker handle Export control? Can individual countries be specified if differing from Docker's list of embargoed countries?

We provide export control via blacklisting several countries, IPs and users
based on the national export compliance database. Any export control we do is
across all products, we do not selectively blacklist versions and products for
specific groups. Send us an email at publisher-support if you have questions

## Analytics

### Where can I view customer insights?

Analytics reports are only available to Publishers with Certified or Commercial
Content. Go to https://store.docker.com/publisher/center and click on "Actions"
for the product you'd like to view analytics for. Here is a
[screenshot](https://user-images.githubusercontent.com/2453622/32352202-6e87ce6e-bfdd-11e7-8fb0-08fe5a3e8930.png).

### How do metrics differentiate between Free and Paid subscribers?

The Analytics reports contain information about the Subscriber and the
relevant product plan. You can identify subscribers for each plan
for each product.

### Can I preview my submission before publishing?

Yes. You can preview your submission including the image you've submitted, the look and feel of the detail page and any markdown descriptions you might have.

Here are a few screenshots that illustrate the preview experience for markdown content.
Product Description preview [screenshot](https://user-images.githubusercontent.com/2453622/32344591-9cd6b456-bfc4-11e7-9505-1f7e8235f812.png).
Install instructions description preview [screenshot](https://user-images.githubusercontent.com/2453622/32344592-9cf2e234-bfc4-11e7-9e60-d773b62eae07.png).

## Other FAQs

### Can a publisher respond to a review of their product?

Yes

### Can I have a publish by date for my content?

Not yet. Potential ETA Q2 2018.
