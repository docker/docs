---
description: >
  Integrate JFrog Artifactory and JFrog Container Registry with Docker Scout
keywords: >
  docker scout, jfrog, artifactory, jcr, integration, image analysis, security,
  cves
title: Artifactory integration
---

{% include scout-early-access.md %}

Integrating Docker Scout with JFrog Artifactory lets you run image analysis
automatically on images in Artifactory registries.

## Local image analysis

You can analyze Artifactory images for vulnerabilities locally using Docker Desktop or the Docker CLI. You first need to authenticate with JFrog Artifactory using the `[Docker login](/engine/reference/commandline/login/)` command. For example:

```bash
docker login {URL}
```

> **Tip**
>
> For cloud-hosted Artifactory you can find the credentials for your Artifactory repository by
> selecting it in the Artifactory UI and then the **Set Me Up** button.
{: .tip }

## Remote image analysis

To automatically analyze images running in remote environments you need to deploy the Docker Scout Artifactory agent. The agent is a
standalone service that analyzes images and uploads the result to Docker Scout.
You can view the results using the
[Docker Scout web UI](https://dso.docker.com/){: target="\_blank" rel="noopener"
}.

### How the agent works

The Docker Scout Artifactory agent is available as an
[image on Docker Hub](https://hub.docker.com/r/docker/artifactory-agent){:
target="\_blank" rel="noopener" }. The agent works by continuously polling
Artifactory for new images. When it finds a new image, it performs the following
steps:

1. Pull the image from Artifactory
2. Analyze the image
3. Upload the analysis result to Docker Scout

The agent records the Software Bill of Materials (SBOM) for the image, and the
SBOMs for all of its base images. The recorded SBOMs include both Operating
System (OS)-level and application-level programs or dependencies that the image
contains.

Additionally, the agent sends the following metadata about the image to Docker Scout:

- The source repository URL and commit SHA for the image
- Build instructions
- Build date
- Tags and digest
- Target platforms
- Layer sizes

The agent never transacts the image
itself, nor any data inside the image, such as code, binaries, and layer blobs.

The agent doesn't detect and analyze pre-existing images. It only analyzes
images that appear in the registry while the agent is running.

### Deploy the agent

This section describes the steps for deploying the Artifactory agent.

#### Prerequisites

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

#### Create the configuration file

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

#### Run the agent

The following example shows how to run the Docker Scout Artifactory agent using
`docker run`. This command creates a bind mount for the directory containing the
JSON configuration file created earlier at `/opt/artifactory-agent/data` inside
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

#### Analyzing pre-existing data

By default the agent detects and analyzes images as they're created and
updated. If you want to use the agent to analyze pre-existing images, you
can use backfill mode. Use the `--backfill-from=TIME` command line option,
where `TIME` is an ISO 8601 formatted time, to run the agent in backfill mode.
If you use this option, the agent analyzes all images pushed between that
time and the current time when the agent starts, then exits.

For example:

```console
$ docker run \
  --mount type=bind,src=/var/opt/artifactory-agent,target=/opt/artifactory-agent/data \
  docker/artifactory-agent:v1 --backfill-from=2022-04-10T10:00:00Z
```

When running a backfill multiple times, the agent won't analyze images that
it's already analyzed. To force re-analysis, provide the `--force` command
line flag.

### View analysis results

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
