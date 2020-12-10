---
description: Pricing & Pull Rate Limiting FAQs
keywords: Docker, Docker Hub, pricing FAQs, pull rate limiting FAQs, subscription, platform
title: Pricing & Pull Rate Limiting FAQs
---

## Pricing FAQs

As of May 14, 2020, Docker announced [new pricing changes](/docker-hub/billing/) for Docker Hub subscriptions.

### What plans and pricing changes did Docker announce on May 14?

Docker announced the following plans and pricing changes.

Immediately available for Individuals and Teams:
* Free plans will continue to be available for both individuals and development teams that include unlimited public repos.
* NEW Pro plan for individuals with unlimited private repos, unlimited public repos, and unlimited collaborators starting at $5/month with an Annual subscription.
* NEW Team plan for development teams with unlimited private repos and unlimited public repos starting at $25 per month for the first 5 users and $7 per month for each user thereafter with an Annual subscription. The Team plan offers advanced collaboration and management tools, including organization and team management with role-based access controls.

### Are prices listed in US dollars?

All prices reflect US dollars.

### How does Team pricing work?

Team starts at $25 per month for the first 5 users and $7 per month for each user thereafter with an Annual subscription.

### How will the new pricing plan impact existing Docker Hub customers?

Legacy customers have until January 31, 2021 to switch to the new pricing plans. For more information on how to move to the new pricing plans, view the [Resource Consumption Updates FAQ](https://www.docker.com/pricing/resource-consumption-updates).

### How can I compare which features are in each plan?

You can see pricing and a full list of features for each product at [Docker Pricing](https://www.docker.com/pricing){: target="blank" rel="noopener" class=""}.

### What is the difference between the legacy plans and the newly announced plans?

The legacy plans were based on a private repository/parallel autobuild pricing model. The new Pro and Team plans are now based on a per-seat pricing model. Both Pro and Team offer unlimited private repos. The Free plan offers unlimited public repositories at no cost per month.

### If I am an existing paid Docker Hub customer, when do I need to change my plan?

If you are an existing Docker subscriber you have until January 31st 2021 to move to the new pricing plans mentioned above. You can convert to either a Monthly plan or an Annual plan.

### If I am an existing subscriber but I don’t do anything by the due date, what will happen?

If no action is taken by the date, Docker will automatically upgrade you to an equivalent plan.

### Will my price per month increase or decrease?

Depending on your configuration you may find it more economical to move to one of the new pricing plans available. For Teams, the key factor affecting a price increase or decrease is the total number of seats you need to support your organization.

### What are the differences between the Team and Free Team plans?

For details on the differences for each plan option, please see [Docker Billing](index.md).

### Do collaborator limits differ between the Free and Pro plans?

The Free plan includes unlimited collaborators for public repositories and 0 collaborators for private repositories. Pro includes unlimited collaborators for public repositories and 1 unique collaborator for private repositories. Note that a collaborator on a private repository within a Pro plan can only be used with a service account.

For details on the differences between plans, see [Docker Pricing](https://www.docker.com/pricing){: target="blank" rel="noopener" class=""}.

For more information on, see [Service accounts](/docker-hub/repos/#service-accounts) page.

### How can I create a new Docker Hub account?

You can create a new account at [Docker Pricing](https://www.docker.com/pricing){: target="blank" rel="noopener" class=""} where you can choose a plan for Individuals or a plan for Teams.

### How do I upgrade to a Pro plan from a legacy individual plan?

Upgrading your legacy plan to a Pro plan offers you unlimited public repos, unlimited private repos, and unlimited collaborators. For information on how to upgrade to a Pro plan from a legacy plan, see [Upgrade to a Pro plan](/docker-hub/billing/upgrade/#upgrade-to-a-pro-plan).

### How do I upgrade to a Team plan from a legacy organization plan?

Upgrading your legacy plan to a Team plan offers you unlimited private repos, unlimited teams, and 3 parallel builds. For information on how to upgrade to a Team plan from a legacy (per-repository) plan, see [Upgrade to a Team plan](/docker-hub/billing/upgrade/#upgrade-to-a-team-plan).

### How do downgrades from a Pro or Team plan work?

When you downgrade your Pro or Team plan, changes are applied at the end of your billing cycle. For example, if you are currently on a Team plan which is billed on the 8th of every month and you choose to downgrade to a Free Team plan on the 15th, your Team plan will be active until 7th of the following month. You will be transferred to a Free Team plan on the 8th of the following month. Please refer to the next two questions below for instructions on how to downgrade your account.

### How do I downgrade from a Team plan to a Free Team plan?

Before you downgrade to a Free plan, ensure that your organization details are updated to reflect features available in the Free plan. For example, you may need to reduce the number of team members and convert any private repositories to public repositories. For information on what’s included in the Free plan, see [Docker Billing](/docker-hub/billing/).

For information on how to downgrade from a Team plan to a Free plan, see [Downgrade from Team to a Free plan](/docker-hub/billing/downgrade/#downgrade-from-team-to-a-free-plan).

### How do I downgrade from Pro to a Free plan?

Before you downgrade to a Free plan, ensure that your account organization details are updated to reflect features available in the Free plan. For example, you may need to convert any private repositories to public repositories. For information on what’s included in the Free plan, see [Docker Billing](/docker-hub/billing/).

For information on how to downgrade from Pro to a Free plan, see [Downgrade from Pro to a Free plan](/docker-hub/billing/downgrade/#downgrade-from-pro-to-a-free-plan).

### How many seats do I need for my organization if I select a Team plan?

Your organization's number of paid seats must equal the number of organization members. If you're not using all of your organization's paid seats, you can downgrade to pay for fewer seats. Organization owners and members each fill a seat. If you've sent a pending invitation to a prospective organization member, the invitation will fill a seat.

### How do I add or remove paid seats for my team members?

For information on how to add paid seats to a monthly plan for your organization, see [Add seats to a monthly plan](/docker-hub/billing/add-seats/).

For information on how to remove paid seats from a monthly plan for your organization, see [Remove seats from a monthly plan](/docker-hub/billing/remove-seats/).

### I want to run an automated agent that makes container requests on behalf of my organization. Which license do I need?

Automated agents or service accounts that make container image requests of Docker Hub must be licensed under a Docker Team subscription.

### How do I add a member to a team in my organization?

For information on how to add a member to a team, see [Add a member to a team](/docker-hub/orgs/#add-a-member-to-a-team).

### Do you offer annual subscriptions?

Yes! Both the Pro and Team plans are available via annual subscriptions and include a discount from the monthly subscription price.

You can view the annual subscription pricing for each product at [Docker Pricing](https://www.docker.com/pricing){: target="blank" rel="noopener" class=""}.

### Are legacy plans eligible for annual pricing?

No, you must move to the new Pro or Team plans to take advantage of the Annual pricing.

### Can Docker Hub continue to be used for open source projects?

Yes, Docker will continue to offer a Free plan with unlimited public repositories and unlimited collaborators for public repositories at no cost per month. Docker is committed to supporting the broader open source communities.

For questions about open source projects and pricing please visit the [Open Source Community Application](https://www.docker.com/community/open-source/application){: target="blank" rel="noopener" class=""}.

### I have more questions regarding the new pricing, who can I reach out to?

If you have any questions about how the new seat-based pricing impacts you that are not covered in the FAQ, please reach out and a member of the Docker team will get back to you.

## Pull Rate Limiting FAQs

### What are the Docker Terms of Service?

The Docker Terms of Service is an agreement between you and Docker that governs your use of Docker products and services. Click on the link to view the full [Docker Terms of Service](https://www.docker.com/legal){: target="blank" rel="noopener" class=""}.

### When do the Docker Terms of Service take effect?

The updates to the [Docker Terms of Service](https://www.docker.com/legal) are effective immediately.

### What changes is Docker making to the Terms of Service?

The most notable changes are in Section 2.5. To see all of the changes, we recommend you read the full [Docker Terms of Service](https://www.docker.com/legal).

### I read that images are going to expire on November 1? What is going on here?

A planned policy for image expiration enforcement for inactive images for free accounts was announced to take effect November 1, 2020. After hearing feedback from customers, this policy has been placed on hold until mid 2021. More information is available in [this blog post](https://www.docker.com/blog/docker-hub-image-retention-policy-delayed-and-subscription-updates/){: target="blank" rel="noopener" class=""}.

### What are the rate limits for pulling Docker images from the Docker Hub Registry?

Rate limits for Docker image pulls are based on the account type of the user requesting the image - not the account type of the image’s owner. These are defined on the [pricing](https://www.docker.com/pricing) page.

The highest entitlement a user has, based on their personal account and any orgs they belong to, will be used. Unauthenticated pull requests are “anonymous” and will be rate limited based on IP address rather than user ID. For more information on authenticating image pulls, please see [How do I authenticate pull requests](/docker-hub/download-rate-limit/#how-do-i-authenticate-pull-requests).

### Are there rate limits on Docker Hub container image requests?

Rate limits for unauthenticated and free Docker Hub usage went into effect on November 1, 2020. More information about this is available in [this blog post](https://www.docker.com/blog/docker-hub-image-retention-policy-delayed-and-subscription-updates/){: target="blank" rel="noopener" class=""}.

### How is a pull request defined for purposes of rate limiting?

A pull request is up to two GET requests to the registry URL path `/v2/*/manifests/*`.

This accounts for the fact that container pull requests for multi-arch images require a manifest list to be downloaded followed by the actual image manifest for the required architecture. HEAD requests are not counted.

Note that all pull requests, including ones for images you already have, are counted by this method. This is the trade-off for not counting individual layers.

### Can I use a Docker Pro subscription to run a service account that makes unlimited container pull requests to Docker Hub on behalf of an organization?

Docker Pro subscriptions are designed for use by a single individual. An automated service account (such as an account used for a teams CI) needs to must be licensed as part of a Docker Team subscription.

### Can I run a local mirror of Docker Hub?

Please see [https://docs.docker.com/registry/recipes/mirror/](https://docs.docker.com/registry/recipes/mirror/) to run [docker/distribution](https://github.com/docker/distribution){: target="blank" rel="noopener" class=""} as a mirror. Because distribution uses HEAD manifest requests to check if cached content has changed, it will not count towards the rate limit. Note that initial requests for Hub images that are not yet in cache will be counted.

### Are image layers counted?

No. Because we are limiting on manifest requests, the number of layers (blob requests) related to a pull is unlimited at this time. This is a change from our previous policy based on community feedback in order to be more user-friendly, so users do not need to count layers on each image they may be using.

### Are anonymous (unauthenticated) pulls rate-limited based on IP address?

Yes. Pull rates are limited based on individual IP address (e.g., for anonymous users: 100 pulls per 6 hours per IP address). See [pricing](https://www.docker.com/pricing){: target="blank" rel="noopener" class=""} for more details.

### Are pulls from logged-in accounts rate-limited based on IP address?

No, limits on pull requests that are [authenticated with an account](/docker-hub/download-rate-limit/#how-do-i-authenticate-pull-requests) are based on that account and not on the IP. Free accounts are limited to 200 pulls per 6 hour period. Paid accounts are subject to fair usage as outlined in [Docker Pricing](https://www.docker.com/pricing){: target="blank" rel="noopener" class=""}.

### Will I be rate-limited if I’m logged in but someone else on my IP hits the anonymous limit?

No, users who are logged in and [authenticating their pull requests](/docker-hub/download-rate-limit/#how-do-i-authenticate-pull-requests) will be limited based on their user account only. If an anonymous user on your IP is limited, you will not be affected as long as your requests are authenticated and you have not hit your own limit.

### Does it matter what image I am pulling?

No, all images are treated equally. The limits are entirely based on the account level of the user doing the pull, not the account level of the repository owner.

### How will these limits be adjusted?

We will be monitoring these limits closely and making sure they are appropriate for common use cases within their tier. In particular, free and anonymous limits should always support normal single-developer workflows. Adjustments will be made as needed in that spirit. You may also [submit input](https://github.com/docker/roadmap/projects/1){: target="blank" rel="noopener" class=""} on the limits.

### What about CI systems where pulls will be anonymous?

We recognize there are some circumstances where many pulls will be made that can not be authenticated. For example, cloud CI providers may host builds based on PRs to open source projects. The project owners may be unable to securely use their project’s Docker Hub credentials to authenticate pulls in this scenario, and the scale of these providers would likely trigger the anonymous rate limits. We will unblock these scenarios as necessary and continue iterating on our rate limiting mechanisms to improve the experience, in cooperation with these providers. Please reach out to [support@docker.com](mailto:support@docker.com) if you are encountering issues.

### Will Docker offer dedicated plans for open source projects?

Yes, as part of Docker’s commitment to the open source community, we have created a new program to support free use of Docker in open source projects. [See this blog post](https://www.docker.com/blog/expanded-support-for-open-source-software-projects/){: target="blank" rel="noopener" class=""} for more information.

To apply for an open source plan, complete our application at [https://www.docker.com/community/open-source/application](https://www.docker.com/community/open-source/application){: target="blank" rel="noopener" class=""}.

For more information, view our blogs on [inactive image retention](https://www.docker.com/blog/scaling-dockers-business-to-serve-millions-more-developers-storage){: target="blank" rel="noopener" class=""} and [pull rate limits](https://www.docker.com/blog/scaling-docker-to-serve-millions-more-developers-network-egress){: target="blank" rel="noopener" class=""}. Alternatively, check out [Docker Pricing](https://www.docker.com/pricing){: target="blank" rel="noopener" class=""}

<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "mainEntity": [
    {
      "@type": "Question",
      "name": "What plans and pricing changes did Docker announce on May 14?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Docker announced the following plans and pricing changes.\n\nImmediately available for Individuals and Teams:\n\nFree plans will continue to be available for both individuals and development teams that include unlimited public repos.\nNEW Pro plan for individuals with unlimited private repos, unlimited public repos, and unlimited collaborators starting at $5/month with an Annual subscription.\nNEW Team plan for development teams with unlimited private repos and unlimited public repos starting at $25 per month for the first 5 users and $7 per month for each user thereafter with an Annual subscription. The Team plan offers advanced collaboration and management tools, including organization and team management with role-based access controls."
      }
    },
    {
      "@type": "Question",
      "name": "Are prices listed in US dollars?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "All prices reflect US dollars."
      }
    },
    {
      "@type": "Question",
      "name": "How does Team pricing work?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Team starts at $25 per month for the first 5 users and $7 per month for each user thereafter with an Annual subscription."
      }
    },
    {
      "@type": "Question",
      "name": "How will the new pricing plan impact existing Docker Hub customers?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Legacy customers have until January 31, 2021 to switch to the new pricing plans. For more information on how to move to the new pricing plans, view the Resource Consumption Updates FAQ."
      }
    },
    {
      "@type": "Question",
      "name": "How can I compare which features are in each plan?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "You can see pricing and a full list of features for each product at https://www.docker.com/pricing."
      }
    },
    {
      "@type": "Question",
      "name": "What is the difference between the legacy plans and the newly announced plans?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "The legacy plans were based on a private repository/parallel autobuild pricing model. The new Pro and Team plans are now based on a per-seat pricing model. Both Pro and Team offer unlimited private repos. The Free plan offers unlimited public repositories at no cost per month."
      }
    },
    {
      "@type": "Question",
      "name": "If I am an existing paid Docker Hub customer, when do I need to change my plan?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "If you are an existing Docker subscriber you have until January 31st 2021 to move to the new pricing plans mentioned above. You can convert to either a Monthly plan or an Annual plan."
      }
    },
    {
      "@type": "Question",
      "name": "If I am an existing subscriber but I don’t do anything by the due date, what will happen?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "If no action is taken by the date, Docker will automatically upgrade you to an equivalent plan."
      }
    },
    {
      "@type": "Question",
      "name": "Will my price per month increase or decrease?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Depending on your configuration you may find it more economical to move to one of the new pricing plans available. For Teams, the key factor affecting a price increase or decrease is the total number of seats you need to support your organization."
      }
    },
    {
      "@type": "Question",
      "name": "What are the differences between the Team and Free Team plans?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "For details on the differences for each plan option, please see https://docs.docker.com/docker-hub/billing/."
      }
    },
    {
      "@type": "Question",
      "name": "Do collaborator limits differ between the Free and Pro plans?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "The Free plan includes unlimited collaborators for public repositories and 0 collaborators for private repositories. Pro includes unlimited collaborators for public repositories and 1 unique collaborator for private repositories. Note that a collaborator on a private repository within a Pro plan can only be used with a service account.\n\nFor details on the differences between plans, see https://www.docker.com/pricing. \\ For more information on service accounts visit our docs page."
      }
    },
    {
      "@type": "Question",
      "name": "How can I create a new Docker Hub account?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "You can create a new account at https://hub.docker.com/pricing where you can choose a plan for Individuals or a plan for Teams."
      }
    },
    {
      "@type": "Question",
      "name": "How do I upgrade to a Pro plan from a legacy individual plan?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Upgrading your legacy plan to a Pro plan offers you unlimited public repos, unlimited private repos, and unlimited collaborators. For information on how to upgrade to a Pro plan from a legacy plan, see Upgrade to a Pro plan."
      }
    },
    {
      "@type": "Question",
      "name": "How do I upgrade to a Team plan from a legacy organization plan?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Upgrading your legacy plan to a Team plan offers you unlimited private repos, unlimited teams, and 3 parallel builds. For information on how to upgrade to a Team plan from a legacy (per-repository) plan, see Upgrade to a Team plan."
      }
    },
    {
      "@type": "Question",
      "name": "How do downgrades from a Pro or Team plan work?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "When you downgrade your Pro or Team plan, changes are applied at the end of your billing cycle. For example, if you are currently on a Team plan which is billed on the 8th of every month and you choose to downgrade to a Free Team plan on the 15th, your Team plan will be active until 7th of the following month. You will be transferred to a Free Team plan on the 8th of the following month. Please refer to the next two questions below for instructions on how to downgrade your account."
      }
    },
    {
      "@type": "Question",
      "name": "How do I downgrade from a Team plan to a Free Team plan?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Before you downgrade to a Free plan, ensure that your organization details are updated to reflect features available in the Free plan. For example, you may need to reduce the number of team members and convert any private repositories to public repositories. For information on what’s included in the Free plan, see https://docs.docker.com/docker-hub/billing/.\n\nFor information on how to downgrade from a Team plan to a Free plan, see Downgrade from Team to a Free plan."
      }
    },
    {
      "@type": "Question",
      "name": "How do I downgrade from Pro to a Free plan?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Before you downgrade to a Free plan, ensure that your account organization details are updated to reflect features available in the Free plan. For example, you may need to convert any private repositories to public repositories. For information on what’s included in the Free plan, see https://docs.docker.com/docker-hub/billing/.\n\nFor information on how to downgrade from Pro to a Free plan, see Downgrade from Pro to a Free plan."
      }
    },
    {
      "@type": "Question",
      "name": "How many seats do I need for my organization if I select a Team plan?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Your organization’s number of paid seats must equal the number of organization members. If you’re not using all of your organization’s paid seats, you can downgrade to pay for fewer seats. Organization owners and members each fill a seat. If you’ve sent a pending invitation to a prospective organization member, the invitation will fill a seat."
      }
    },
    {
      "@type": "Question",
      "name": "How do I add or remove paid seats for my team members?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "For information on how to add paid seats to a monthly plan for your organization, see Add seats to a monthly plan.\n\nFor information on how to remove paid seats from a monthly plan for your organization, see Remove seats from a monthly plan."
      }
    },
    {
      "@type": "Question",
      "name": "I want to run an automated agent that makes container requests on behalf of my organization. Which license do I need?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Automated agents or service accounts that make container image requests of Docker Hub must be licensed under a Docker Team subscription."
      }
    },
    {
      "@type": "Question",
      "name": "How do I add a member to a team in my organization?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "For information on how to add a member to a team, see Add a member to a team."
      }
    },
    {
      "@type": "Question",
      "name": "Do you offer annual subscriptions?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Yes! Both the Pro and Team plans are available via annual subscriptions and include a discount from the monthly subscription price.\n\nYou can view the annual subscription pricing for each product at https://www.docker.com/pricing."
      }
    },
    {
      "@type": "Question",
      "name": "Are legacy plans eligible for annual pricing?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "No, you must move to the new Pro or Team plans to take advantage of the Annual pricing."
      }
    },
    {
      "@type": "Question",
      "name": "Can Docker Hub continue to be used for open source projects?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Yes, Docker will continue to offer a Free plan with unlimited public repositories and unlimited collaborators for public repositories at no cost per month. Docker is committed to supporting the broader open source communities.\n\nFor questions about open source projects and pricing please visit the Open Source Community Application."
      }
    },
    {
      "@type": "Question",
      "name": "I have more questions regarding the new pricing, who can I reach out to?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "If you have any questions about how the new seat-based pricing impacts you that are not covered in the FAQ, please reach out and a member of the Docker team will get back to you."
      }
    },
    {
      "@type": "Question",
      "name": "What are the Docker Terms of Service?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "The Docker Terms of Service is an agreement between you and Docker that governs your use of Docker products and services. Click on the link to view the full Docker Terms of Service."
      }
    },
    {
      "@type": "Question",
      "name": "When do the Docker Terms of Service take effect?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "The updates to the Docker Terms of Service are effective immediately."
      }
    },
    {
      "@type": "Question",
      "name": "What changes is Docker making to the Terms of Service?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "The most notable changes are in Section 2.5. To see all of the changes, we recommend you read the full Docker Terms of Service."
      }
    },
    {
      "@type": "Question",
      "name": "I read that images are going to expire on November 1? What is going on here?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "A planned policy for image expiration enforcement for inactive images for free accounts was announced to take effect November 1, 2020. After hearing feedback from customers, this policy has been placed on hold until mid 2021. More information is available in this blog post."
      }
    },
    {
      "@type": "Question",
      "name": "What are the rate limits for pulling Docker images from the Docker Hub Registry?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Rate limits for Docker image pulls are based on the account type of the user requesting the image - not the account type of the image’s owner. These are defined on the pricing page.\n\nThe highest entitlement a user has, based on their personal account and any orgs they belong to, will be used. Unauthenticated pull requests are “anonymous” and will be rate limited based on IP address rather than user ID. For more information on authenticating image pulls, please see this docs page."
      }
    },
    {
      "@type": "Question",
      "name": "Are there rate limits on Docker Hub container image requests?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Rate limits for unauthenticated and free Docker Hub usage went into effect on November 1, 2020. More information about this is available in this blog post."
      }
    },
    {
      "@type": "Question",
      "name": "How is a pull request defined for purposes of rate limiting?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "A pull request is up to two GET requests to the registry URL path /v2/*/manifests/*.\n\nThis accounts for the fact that container pull requests for multi-arch images require a manifest list to be downloaded followed by the actual image manifest for the required architecture. HEAD requests are not counted.\n\nNote that all pull requests, including ones for images you already have, are counted by this method. This is the trade-off for not counting individual layers."
      }
    },
    {
      "@type": "Question",
      "name": "Can I use a Docker Pro subscription to run a service account that makes unlimited container pull requests to Docker Hub on behalf of an organization?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Docker Pro subscriptions are designed for use by a single individual. An automated service account (such as an account used for a teams CI) needs to must be licensed as part of a Docker Team subscription."
      }
    },
    {
      "@type": "Question",
      "name": "Can I run a local mirror of Docker Hub?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Please see https://docs.docker.com/registry/recipes/mirror/ to run docker/distribution as a mirror. Because distribution uses HEAD manifest requests to check if cached content has changed, it will not count towards the rate limit. Note that initial requests for Hub images that are not yet in cache will be counted."
      }
    },
    {
      "@type": "Question",
      "name": "Are image layers counted?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "No. Because we are limiting on manifest requests, the number of layers (blob requests) related to a pull is unlimited at this time. This is a change from our previous policy based on community feedback in order to be more user-friendly, so users do not need to count layers on each image they may be using."
      }
    },
    {
      "@type": "Question",
      "name": "Are anonymous (unauthenticated) pulls rate-limited based on IP address?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Yes. Pull rates are limited based on individual IP address (e.g., for anonymous users: 100 pulls per 6 hours per IP address). See pricing for more details."
      }
    },
    {
      "@type": "Question",
      "name": "Are pulls from logged-in accounts rate-limited based on IP address?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "No, limits on pull requests that are authenticated with an account are based on that account and not on the IP. Free accounts are limited to 200 pulls per 6 hour period. Paid accounts are subject to fair usage as outlined in Docker pricing."
      }
    },
    {
      "@type": "Question",
      "name": "Will I be rate-limited if I’m logged in but someone else on my IP hits the anonymous limit?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "No, users who are logged in and authenticating their pull requests will be limited based on their user account only. If an anonymous user on your IP is limited, you will not be affected as long as your requests are authenticated and you have not hit your own limit."
      }
    },
    {
      "@type": "Question",
      "name": "Does it matter what image I am pulling?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "No, all images are treated equally. The limits are entirely based on the account level of the user doing the pull, not the account level of the repository owner."
      }
    },
    {
      "@type": "Question",
      "name": "How will these limits be adjusted?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "We will be monitoring these limits closely and making sure they are appropriate for common use cases within their tier. In particular, free and anonymous limits should always support normal single-developer workflows. Adjustments will be made as needed in that spirit. You may also submit input on the limits."
      }
    },
    {
      "@type": "Question",
      "name": "What about CI systems where pulls will be anonymous?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "We recognize there are some circumstances where many pulls will be made that can not be authenticated. For example, cloud CI providers may host builds based on PRs to open source projects. The project owners may be unable to securely use their project’s Docker Hub credentials to authenticate pulls in this scenario, and the scale of these providers would likely trigger the anonymous rate limits. We will unblock these scenarios as necessary and continue iterating on our rate limiting mechanisms to improve the experience, in cooperation with these providers. Please reach out to support@docker.com if you are encountering issues."
      }
    },
    {
      "@type": "Question",
      "name": "Will Docker offer dedicated plans for open source projects?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Yes, as part of Docker’s commitment to the open source community, we have created a new program to support free use of Docker in open source projects. See this blog post for more information.\n\nTo apply for an open source plan, complete our application at:\n\nhttps://www.docker.com/community/open-source/application.\n\nFor more information view our blogs on inactive image retention and pull rate limits.\n\nOr check out: https://www.docker.com/pricing"
      }
    }
  ]
}
</script>


