---
title: Publish your extension
description: General steps in how to publish an extension
keywords: Docker, Extensions, sdk, publish
---

To publish the extension, upload the Docker image to [DockerHub](https://hub.docker.com/).

Tag the previous image to prepend the account owner at the beginning of the image name:

`docker tag <name-of-your-extension> owner/<name-of-your-extension>`

Push the image to DockerHub:

`docker push owner/<name-of-your-extension>`

> Note
> 
> For Docker Extensions images to be listed in Docker Desktop, they must be approved by Docker and the tags must follow semantic versioning, e.g: `0.0.1`.
> 
> See [distribution and new releases](https://docs.docker.com/desktop/extensions-sdk/extensions/DISTRIBUTION/#distribution-and-new-releases) for more information.
> 
> See [semver.org](https://semver.org/) to learn more about semantic versioning.
> 

> Having trouble pushing the image?
> 
> Ensure you are logged into DockerHub. Otherwise, run `docker login` to authenticate.
> 

## Publish your extension in the Marketplace

Docker Desktop displays published extensions in the Extensions Marketplace.
If you want your extension to be published in the Marketplace, you can submit your extension [here](https://www.docker.com/products/extensions/submissions/). We'll review your submission and provide feedback if changes are needed before we can validate and publish it to make it available to all Docker Desktop users.

## Clean up

To remove the extension, run:

`docker extension rm <name-of-your-extension>`

## What's next
Find more information about:
- [The `metadata.json` file](METADATA.md)
- [Labels in your `Dockerfile`](labels.md)
- [Distributing your extension](DISTRIBUTION.md)
