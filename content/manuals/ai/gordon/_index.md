---
title: Ask Gordon
description: Streamline your workflow with Docker's AI-powered assistant in Docker Desktop and CLI.
weight: 10
params:
  sidebar:
    badge:
      color: blue
      text: Beta
    group: AI
aliases:
 - /desktop/features/gordon/
---

{{< summary-bar feature_name="Ask Gordon" >}}

Ask Gordon is your personal AI assistant embedded in Docker Desktop and the
Docker CLI. It's designed to streamline your workflow and help you make the most
of the Docker ecosystem.

## Key features

Ask Gordon provides AI-powered assistance in Docker tools. It can:

- Improve Dockerfiles
- Run and troubleshoot containers
- Interact with your images and code
- Find vulnerabilities or configuration issues
- Migrate a Dockerfile to use [Docker Hardened Images](/manuals/dhi/_index.md)

It understands your local environment, including source code, Dockerfiles, and
images, to provide personalized and actionable guidance.

Ask Gordon remembers conversations, allowing you to switch topics more easily.

Ask Gordon is not enabled by default, and is not
production-ready. You may also encounter the term "Docker AI" as a broader
reference to this technology.

> [!NOTE]
>
> Ask Gordon is powered by Large Language Models (LLMs). Like all
> LLM-based tools, its responses may sometimes be inaccurate. Always verify the
> information provided.

### What data does Gordon access?

When you use Ask Gordon, the data it accesses depends on your query:

- Local files: If you use the `docker ai` command, Ask Gordon can access files
  and directories in the current working directory where the command is
  executed. In Docker Desktop, if you ask about a specific file or directory in
  the **Ask Gordon** view, you'll be prompted to select the relevant context.
- Local images: Gordon integrates with Docker Desktop and can view all images in
  your local image store. This includes images you've built or pulled from a
  registry.

To provide accurate responses, Ask Gordon may send relevant files, directories,
or image metadata to the Gordon backend with your query. This data transfer
occurs over the network but is never stored persistently or shared with third
parties. It is used only to process your request and formulate a response. For
details about privacy terms and conditions for Docker AI, review [Gordon's
Supplemental Terms](https://www.docker.com/legal/docker-ai-supplemental-terms/).

All data transferred is encrypted in transit.

### How your data is collected and used

Docker collects anonymized data from your interactions with Ask Gordon to
improve the service. This includes:

- Your queries: Questions you ask Gordon.
- Responses: Answers provided by Gordon.
- Feedback: Thumbs-up and thumbs-down ratings.

To ensure privacy and security:

- Data is anonymized and cannot be traced back to you or your account.
- Docker does not use this data to train AI models or share it with third
  parties.

By using Ask Gordon, you help improve Docker AI's reliability and accuracy for
everyone.

If you have concerns about data collection or usage, you can
[disable](#disable-ask-gordon) the feature at any time.

## Enable Ask Gordon

1. Sign in to your Docker account.
1. Go to the **Beta features** tab in settings.
1. Check the **Enable Docker AI** checkbox.

   The Docker AI terms of service agreement appears. You must agree to the terms
   before you can enable the feature. Review the terms and select **Accept and
   enable** to continue.

1. Select **Apply**.

> [!IMPORTANT]
>
> For Docker Desktop versions 4.41 and earlier, this setting is under the
> **Experimental features** tab on the **Features in development** page.

## Using Ask Gordon

You can access Gordon:

- In Docker Desktop, in the **Ask Gordon** view.
- In the Docker CLI, with the `docker ai` command.

After you enable Docker AI features, you will also see **Ask Gordon** in other
places in Docker Desktop. Whenever you see a button with the **Sparkles** (✨)
icon, you can use it to get contextual support from Ask Gordon.

## Example workflows

Ask Gordon is a general-purpose AI assistant for Docker tasks and workflows. Here
are some things you can try:

- [Troubleshoot a crashed container](#troubleshoot-a-crashed-container)
- [Get help with running a container](#get-help-with-running-a-container)
- [Improve a Dockerfile](#improve-a-dockerfile)
- [Migrate a Dockerfile to DHI](#migrate-a-dockerfile-to-dhi)

For more examples, try asking Gordon directly. For example:

```console
$ docker ai "What can you do?"
```

### Troubleshoot a crashed container

If you start a container with an invalid configuration or command, use Ask Gordon
to troubleshoot the error. For example, try starting a Postgres container without
a database password:

```console
$ docker run postgres
Error: Database is uninitialized and superuser password is not specified.
       You must specify POSTGRES_PASSWORD to a non-empty value for the
       superuser. For example, "-e POSTGRES_PASSWORD=password" on "docker run".

       You may also use "POSTGRES_HOST_AUTH_METHOD=trust" to allow all
       connections without a password. This is *not* recommended.

       See PostgreSQL documentation about "trust":
       https://www.postgresql.org/docs/current/auth-trust.html
```

In the **Containers** view in Docker Desktop, select the ✨ icon next to the
container's name, or inspect the container and open the **Ask Gordon** tab.

### Get help with running a container

If you want to run a specific image but are not sure how, Gordon can help you get
set up:

1. Pull an image from Docker Hub (for example, `postgres`).
1. Open the **Images** view in Docker Desktop and select the image.
1. Select the **Run** button.

In the **Run a new container** dialog, you see a message about **Ask Gordon**.

![Screenshot showing Ask Gordon hint in Docker Desktop.](../../images/gordon-run-ctr.png)

The linked text in the hint is a suggested prompt to start a conversation with
Ask Gordon.

### Improve a Dockerfile

Gordon can analyze your Dockerfile and suggest improvements. To have Gordon
evaluate your Dockerfile using the `docker ai` command:

1. Go to your project directory:

   ```console
   $ cd <path-to-your-project>
   ```

1. Use the `docker ai` command to rate your Dockerfile:

   ```console
   $ docker ai rate my Dockerfile
   ```

Gordon will analyze your Dockerfile and identify opportunities for improvement
across several dimensions:

- Build cache optimization
- Security
- Image size efficiency
- Best practices compliance
- Maintainability
- Reproducibility
- Portability
- Resource efficiency

### Migrate a Dockerfile to DHI

Migrating your Dockerfile to use [Docker Hardened Images](/manuals/dhi/_index.md)
helps you build more secure, minimal, and production-ready containers. DHIs
reduce vulnerabilities, enforce best practices, and simplify compliance, making
them a strong foundation for secure software supply chains.

To request Gordon's help for the migration:

{{% include "gordondhi.md" %}}

## Disable Ask Gordon

### For individual users

If you've enabled Ask Gordon and you want to disable it again:

1. Open the **Settings** view in Docker Desktop.
1. Go to **Beta features**.
1. Clear the **Enable Docker AI** checkbox.
1. Select **Apply**.

### For organizations

To disable Ask Gordon for your entire Docker organization, use [Settings
Management](/manuals/enterprise/security/hardened-desktop/settings-management/_index.md)
and add this property to your `admin-settings.json` file:

```json
{
  "enableDockerAI": {
    "value": false,
    "locked": true
  }
}
```

Or disable all Beta features by setting `allowBetaFeatures` to false:

```json
{
  "allowBetaFeatures": {
    "value": false,
    "locked": true
  }
}
```

## Feedback

<!-- vale Docker.We = NO -->

We value your input on Ask Gordon and encourage you to share your experience.
Your feedback helps us improve and refine Ask Gordon for all users. If you
encounter issues, have suggestions, or simply want to share what you like,
here's how you can get in touch:

- Thumbs-up and thumbs-down buttons

  Rate Ask Gordon's responses using the thumbs-up or thumbs-down buttons in the
  response.

- Feedback survey

  You can access the Ask Gordon survey by following the _Give feedback_ link in
  the **Ask Gordon** view in Docker Desktop, or from the CLI by running the
  `docker ai feedback` command.


