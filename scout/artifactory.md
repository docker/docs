---
description: >
  Integrate JFrog Artifactory and JFrog Container Registry with Docker Scout
keywords: >
  docker scout, jfrog, artifactory, jcr, integration, image analysis, security,
  cves
title: Artifactory integration
---

> **Note**
>
> Docker Scout is an [early access](../release-lifecycle.md#early-access-ea)
> product.
>
> If you're interested in this integration for your organization and want to
> learn more, get in touch by filling out the contact form on the
> [Docker Scout product page](https://docker.com/products/docker-scout){:
> target="\_blank" rel="noopener" }.

Integrating Docker Scout with JFrog Artifactory lets you run image analysis
automatically on images in your Artifactory registries.

This integration is made possible by a monitoring agent. The agent is a
standalone service that analyzes images and uploads the result to Docker Scout.
You can view the results using the
[Docker Scout web UI](https://dso.docker.com/){: target="\_blank" rel="noopener"
}.

## How it works

The Docker Scout Artifactory agent is available as an
[image on Docker Hub](https://hub.docker.com/r/docker/artifactory-agent){:
target="\_blank" rel="noopener" }. The agent works by continuously polling
Artifactory for new images. When it finds a new image, it performs the following
steps:

1. Pull the image from Artifactory
2. Analyze the image
3. Upload the analysis result to Docker Scout

The agent records the Software Bill of Material (SBOM) for the image, and the
SBOMs for all of its base images. The recorded SBOMs include both Operating
System (OS)-level and application-level programs or dependencies that the image
contains.

Additionally, metadata about the image itself is also recorded:

- The source repository for the image
- Build instructions
- Build date
- Tags and digest
- Target platforms
- Layer sizes

The agent sends this data to Docker Scout. The agent never transacts the image
itself, nor any data inside the image, such as code, binaries, and layer blobs.

The agent doesn't detect and analyze pre-existing images. It only analyzes
images that appear in the registry while the agent is running.

## Deploy the agent

This section describes the steps for deploying the Artifactory agent.

### Prerequisites

Before you deploy the agent, ensure that you meet the prerequisites:

- The server where you host the agent can access the following resources over
  the network:
  - Your JFrog Artifactory instance
  - `hub.docker.com`, port 443, for authenticating with Docker
  - `api.dso.docker.com`, port 443, for transacting data to Docker Scout
- The server isn't behind a proxy
- The registries are Docker V2 registries. V1 registries aren't supported.

The agent supports all versions of JFrog Artifactory and JFrog Container
Registry.

### Create the configuration file

You configure the agent using a JSON file. The agent expects the configuration
file to be in `/opt/artifactory-agent/data/config.json` on startup.

The configuration file includes the following properties:

| Property                    | Description                                                                     |
| --------------------------- | ------------------------------------------------------------------------------- |
| `agent_id`                  | Unique identifier for the agent.                                                |
| `docker.organization_name`  | Name of the Docker organization.                                                |
| `docker.username`           | Username of the admin user in the Docker organization.                          |
| `docker.pat`                | Personal access token of the admin user with read and write permissions.        |
| `artifactory.base_url`      | Base URL of the Artifactory instance.                                           |
| `artifactory.username`      | Username of the Artifactory user with read permissions that the agent will use. |
| `artifactory.password`      | Password or API token for the Artifactory user.                                 |
| `artifactory.image_filters` | Optional: List of repositories and images to analyze.                           |

If you don't specify any repositories in `artifactory.image_filters`, the agent
runs image analysis on all images in your Artifactory instance.

The following snippet shows a sample configuration:

```json
{
  "agent_id": "acme-prod-agent",
  "docker": {
    "organization_name": "acme",
    "username": "mobythewhale",
    "pat": "dckr_pat__dsaCAs_xL3kNyupAa7dwO1alwg"
  },
  "artifactory": [
    {
      "base_url": "https://acme.jfrog.io",
      "username": "acmeagent",
      "password": "hayKMvFKkFp42RAwKz2K",
      "image_filters": [
        {
          "repository": "dev-local",
          "images": ["internal/repo1", "internal/repo2"]
        },
        {
          "repository": "prod-local",
          "images": ["staging/repo1", "prod/repo1"]
        }
      ]
    }
  ]
}
```

Create a configuration file and save it somewhere on the server where you plan
to run the agent. For example, `/var/opt/artifactory-agent/config.json`.

### Run the agent

The following example shows how to run the Docker Scout Artifactory agent using
`docker run`. This command creates a bind mount for the directory containing the
JSON configuration file created earlier to `/opt/artifactory-agent/data` inside
the container. Make sure the mount path you use is the directory containing the
`config.json` file.

<!-- prettier-ignore -->
> **Important**
>
> Use the `v1` tag of the Artifactory agent image. Don't use the `latest` tag as
> doing so may incur breaking changes.
{: .important }

```console
$ docker run \
  --mount type=bind,src=/var/opt/artifactory-agent,target=/opt/artifactory-agent/data \
  docker/artifactory-agent:v1
```

## View analysis results

You can view the image analysis results in the Docker Scout web UI.

1. Go to [Docker Scout web UI](https://dso.docker.com).
2. Sign in using your Docker ID.

   Once signed in, you're taken to the **Images** page. This page displays the
   repositories in your organization connected to Docker Scout.

3. Select the image in the list.
4. Select the tag.

When you have selected a tag, you're taken to the vulnerability report for that
tag. Here, you can select if you want to view all vulnerabilities in the image,
or vulnerabilities introduced in a specific layer. You can also filter
vulnerabilities by severity, and whether or not there's a fix version available.
