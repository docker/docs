---
title: Publish your extension
description: General steps in how to publish an extension
keywords: Docker, Extensions, sdk, publish
---

This section provides a how-to guide, general information on publishing your extension to the Extensions Marketplace.

For information on how Docker Extensions are packaged and distributed, see [Packaging, distribution and API dependencies](DISTRIBUTION.md)

## How to publish your extension 

To publish the extension you need to upload the Docker image to [DockerHub](https://hub.docker.com/){: target="_blank" rel="noopener" class="_" }:

1. Tag the previous image to prepend the account owner at the beginning of the image name:

    ```console
    $ docker tag <name-of-your-extension> owner/<name-of-your-extension>
    ```

    > Note
    > 
    > For Docker Extensions images to be listed in Docker Desktop, they must be approved by Docker and the tags must follow semantic versioning, e.g: `0.0.1`.
    > 
    > See [distribution and new releases](DISTRIBUTION.md#distribution-and-new-releases) for more information.
    > 
    > See [semver.org](https://semver.org/){:target="_blank" rel="noopener" class="_"} to learn more about semantic versioning.

2. Push the image to Docker Hub:

    ```console
    $ docker push owner/<name-of-your-extension>
    ```

    > Having trouble pushing the image?
    >  
    > Ensure you are logged into DockerHub. Otherwise, run `docker login` to authenticate.

## Submit your extension to be published in the Extensions Marketplace

Docker Desktop displays published extensions in [the Extensions Marketplace](https://hub.docker.com/search?q=&type=extension){: target="_blank" rel="noopener" class="_" }. The Extensions Marketplace is a curated space where developers can discover extensions to improve their developer experience and upload their own extension to share with the world. 

If you want your extension to be published in the Marketplace, you can submit your extension [here](https://www.docker.com/products/extensions/submissions/){: target="_blank" rel="noopener" class="_" }. 

All extensions submitted to the Extension Marketplace are reviewed and approved by our team before listing. This review process ensures a level of trust, security, and quality for developers using Docker Extensions and allows for extension developers to get feedback.

### Before you submit

Ensure your extension has followed the guidelines outlined in this section before submitting for your extension for review. We highly encourage you to check our guidelines as not doing so may considerably impact the duration of the approval process. 

These guidelines do not replace our terms of service or guarantee approval. As the Extension Marketplace continues adding new features for both Extension users and publishers, expect that your extension should be maintained over time to ensure it stays available in the Marketplace.

#### Guidelines:
- Test your extension for crashes, bugs, and performance issues
- Test your extension with potential users
- Ensure that you’ve ran our [validation checks](../build/build-install.md)
- Review our [design guidelines](../design/design-guidelines.md)
- Ensure the [UI styling](../design/overview.md) is in line with Docker Desktop guidelines
- Ensure your extensions support both light and dark mode
- Consider the needs of both new and existing users of your extension
- Test your extension on various platforms (Mac, Windows, Linux)
- Read our [Terms of Service](https://www.docker.com/legal/extensions_marketplace_developer_agreement/){: target="_blank" rel="noopener" class="_" }

### After you submit

Once you’ve submitted your extension, here is what you can expect from the review process:

- Timing: Extensions are reviewed by us manually. Although we strive for having your submission approved as soon as possible, bear in mind this is a manual process to ensure extensions meet high standards. If your extension is complex, if it does not follow our guidelines, or if you did not complete the submission form properly, it may require more time to properly consider your extension.
- Rejections: Docker strives to review extensions for consideration fairly and consistently. We will do our best to provide adequate and actionable feedback for you so that we can reconsider publishing your extension after you’ve made appropriate changes. If your extension has been rejected, you can communicate directly with us.

## What's next
Find more information about:
- [The `metadata.json` file](METADATA.md)
- [Labels in your `Dockerfile`](labels.md)
- [Distributing your extension](DISTRIBUTION.md)
- [Building extensions for multiple architectures](multi-arch.md)
