---
description: Learn how to describe and optimize your Docker Hub repositories for better discoverability.
keywords: Docker Hub, Hub, repository information, repository discoverability, best practices
title: Repository information
toc_max: 3
weight: 40
aliases:
- /docker-hub/repos/categories/
---

Each repository can include a description, an overview, and categories to help
users understand its purpose and usage. Adding clear repository information
ensures that others can find your images and use them effectively.

You can only modify the repository information of repositories that aren't
archived. If a repository is archived, you must unarchive it to modify the
information. For more details, see [Unarchive a repository](../archive.md#unarchive-a-repository).

## Repository description

The description appears in search results when using the `docker search` command
and in the search results on Docker Hub.

Consider the following repository description best practices.

- Summarize the purpose. Clearly state what the image does in a concise and
  specific manner. Make it clear if it's for a particular application, tool, or
  platform, or has a distinct use case.
- Highlight key features or benefits. Briefly mention the primary benefits or
  unique features that differentiate the image. Examples include high
  performance, ease of use, lightweight build, or compatibility with specific
  frameworks or operating systems.
- Include relevant keywords. Use keywords that users may search for to increase
  visibility, such as technology stacks, use cases, or environments.
- Keep it concise. The description can be a maximum of 100 characters. Aim to
  stay within one or two sentences for the description to ensure it's easy to
  read in search results. Users should quickly understand the image's value.
- Focus on the audience. Consider your target audience (developers, system
  administrators, etc.) and write the description to address their needs
  directly.

Following these practices can help make the description more engaging and
effective in search results, driving more relevant traffic to your repository.

### Add or update a repository description

1. Sign in to [Docker Hub](https://hub.docker.com).

2. Select **Repositories**.

   A list of your repositories appears.

3. Select a repository.

   The **General** page for the repository appears.

4. Select the pencil icon under the description field.

5. Specify a description.

   The description can be up to 100 characters long.

6. Select **Update**.

## Repository overview

An overview describes what your image does and how to run it. It displays in the
public view of your repository when the repository has at least one image. If
automated builds are enabled, the overview will be synced from the source code
repository's `README.md` file on each successful build.

Consider the following repository overview best practices.

- Describe what the image is, the features it offers, and why it should be used.
  Can include examples of usage or the team behind the project.
- Explain how to get started with running a container using the image. You can
  include a minimal example of how to use the image in a Dockerfile.
- List the key image variants and tags to use them, as well as use cases for the
  variants.
- Link to documentation or support sites, communities, or mailing lists for
  additional resources.
- Provide contact information for the image maintainers.
- Include the license for the image and where to find more details if needed.

### Add or update a repository overview

1. Sign in to [Docker Hub](https://hub.docker.com).

2. Select **Repositories**.

   A list of your repositories appears.

3. Select a repository.

   The **General** page for the repository appears.

4. Under **Repository overview**, select **Edit** or **Add overview**.

   The **Write** and **Preview** tabs appear.

5. Under **Write**, specify your repository overview.

   You can use basic Markdown and use the **Preview** tab to preview the formatting.

6. Select **Update**.

## Repository categories

You can tag Docker Hub repositories with categories, representing the primary
intended use cases for your images. This lets users more easily find and
explore content for the problem domain that they're interested in.

### Available categories

The Docker Hub content team maintains a curated list of categories.

{{% include "hub-categories.md" %}}

### Auto-generated categories

> [!NOTE]
>
> Auto-generated categories only apply to Docker Verified Publishers and
> Docker-Sponsored Open Source program participants.

For repositories that pre-date the Categories feature in Docker Hub,
categories have been automatically generated and applied, using OpenAI, based
on the repository title and description.

As an owner of a repository that has been auto-categorized, you can manually
edit the categories if you think they're inaccurate. See [Manage categories for
a repository](#manage-categories-for-a-repository).

The auto-generated categorization was a one-time effort to help seed categories
onto repositories created before the feature existed. Categories are not
assigned to new repositories automatically.

### Manage categories for a repository

You can tag a repository with up to three categories.

To edit the categories of a repository:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **Repositories**.

   A list of your repositories appears.

3. Select a repository.

   The **General** page for the repository appears.

4. Select the pencil icon under the description field.
5. Select the categories you want to apply.
6. Select **Update**.

If you're missing a category, use the
[Give feedback link](https://docker.qualtrics.com/jfe/form/SV_03CrMyAkCWVylKu)
to let us know what categories you'd like to see.