---
description: Docker Cloud
keywords: Docker, cloud, release, notes
title: Docker Cloud release notes
---

Did you know we also have a [Release notes category](https://forums.docker.com/c/docker-cloud/release-notes) on the Docker Cloud Product forums? Now you do!

## Docker Cloud June 2016 release notes

In the last month we've made many small improvements to the new Docker Cloud UI, made team and organization repos from Hub visible in Docker Cloud, and enabled linking to BitBucket for automated builds.

We've also added significant new features to the [automated builds](builds/automated-build.md) system, including:

- Branch and tag selection
- Dynamic build rules (AKA regex build rules)
- Different hosted builder node sizes

For more details, find the June post in the [release notes category](https://forums.docker.com/c/docker-cloud/release-notes), and as always, we welcome your feedback on the [Docker Cloud Product Forums](https://forums.docker.com/c/docker-cloud).

## Docker Cloud May 2016 release notes

In our May 2016 release, we introduced a new user interface for Docker Cloud. Try it out and share your feedback in the [Docker Cloud Product Forums](https://forums.docker.com/c/docker-cloud)!

### Added

**Docker Cloud Security Scanning** is now available as a beta add-on service for private repositories. 

### Fixed

- **API docs now say CLI instead of bash** in the languages tab. You pointed out that this was confusing, so we fixed it.
- **Removed old references to Tutum** in the documentation.

### Known issues

- **Documentation screen captures** in some cases still reflect the Docker Cloud 1.0 user interface. This will be updated as soon as possible.

Additional [Known issues here](docker-errors-faq.md)

## Docker Cloud 1.0 release notes
**Tutum is now Docker Cloud**. Docker Cloud is a new service by Docker that implements all features previously offered by Tutum plus integration with Docker Hub Registry service and the common Docker ID credentials.

The following release notes document changes since [Tutum v0.19.5](https://support.tutum.co/support/solutions/articles/5000694910-tutum-0-19-5).


### Added

- **Docker Cloud is Generally Available**: all features of Docker Cloud are Generally Available with the exception of the build features which remain in beta.
- **Docker Hub Registry Integration**: all of your Docker Hub image repositories are available and accessible when you login to Docker Cloud. Changes you make to your repositories are reflected in both Docker Hub and Docker Cloud.
- **Autoredeploy from Docker Hub**: services that use a repository stored in the Docker Hub now have the [**autoredeploy** option](apps/auto-redeploy.md) available, which allows automatic redeployments on push without setting up webhooks.
- **Environment variable substitution on CLI**: the `docker-cloud` CLI now substitutes environment variables in stack files, [the same way Docker Compose does it](/compose/compose-file/#variable-substitution:91de898b5f5cdb090642a917d3dedf68).


### Changed

- **Tutum is now Docker Cloud**: Docker Cloud is a new service by Docker that implements all features previously offered by Tutum.
- **Docker ID**: your Docker ID (formerly known as "Docker Hub account") is used to log into Docker Cloud.
- **Environment variables**: the environment variables that are automatically injected into containers that started with `TUTUM_` now start with `DOCKERCLOUD_`.
- **CLI renaming**: the `tutum` CLI has been deprecated and the new Docker Cloud CLI is now called `docker-cloud`. Login credentials are now shared between the `docker` and `docker-cli` CLIs and stored in `~/.docker/config.json`.
- **API domain**: the API domain is now `https://cloud.docker.com` for REST endpoints, and `wss://ws.cloud.docker.com` for websocket endpoints.
- **API endpoints**: the API endpoints have been relocated to a different URI scheme. [Click here for full documentation about the new endpoints](/apidocs/docker-cloud.md).
- **New Python and Go SDKs**: the new **[python-dockercloud](https://github.com/docker/python-dockercloud)** and **[go-dockercloud](https://github.com/docker/go-dockercloud)** SDKs are available to work with the new Docker Cloud APIs.
- **New HAproxy image**: the new `dockercloud/haproxy` repository can be used as a proxy/load balancer for user's applications and will automatically reconfigure itself if configured with API access via API role.
- **Docker Registry**: the Docker registry at `tutum.co` has been deprecated and replaced by the Docker Hub. It requires Docker Engine 1.6 or higher. Repositories are now shared between Docker Cloud and Docker Hub and will appear in both sites.
- **Agent renamed**: the `tutum-agent` has been renamed to `dockercloud-agent`. The installation script is now at `https://get.cloud.docker.com`. Its configuration file is now at `/etc/dockercloud/agent/` and logs are stored at `/var/log/dockercloud/`.
- **Deploy to Docker Cloud button**: the "Deploy to Tutum" button has been renamed to **Deploy to Docker Cloud**. [Click here to learn more](apps/deploy-to-cloud-btn.md).
- **AWS object names**: the names of the objects created by default in AWS have changed: the VPC is now called `dc-vpc` and has a CIDR of `10.78.0.0/16`, the subnets are called `dc-subnet`, the security group is now called `dc-vpc-default`, the internet gateway is now called `dc-gateway` and the route table is now called `dc-route-table`.
- **User endpoints**: the new domain used by node, service and container endpoints is now `dockerapp.io`. Endpoints now do not include the username and use short UUIDs to ensure uniqueness.
- **Community Forums**: the [Docker Cloud forums](https://forums.docker.com/c/docker-cloud) are now the recommended place to get in touch with the community.


### Fixed

- **Overlay network**: we have fixed a memory limit issue on the overlay network containers that was causing containers to not attach to the overlay network under certain circumstances.
- **Scale up trigger**: we have fixed an issue where sometimes containers created by using a "scale up" trigger didn't inherit the service configuration and marked all other containers in the service with the "redeployment needed" flag.

### Known issues

- **Documentation screen captures** in most cases still reflect the Tutum interface and branding. We will update these and refresh the documentation as we go.
- **References to Tutum remain** in the documentation. We will update these and refresh the documentation as we go.
