---
title: Publish your extension in the Marketplace
description: Docker extension distribution
keywords: Docker, extensions, publish
---

## Submit your extension to the Marketplace

Docker Desktop displays published extensions in the Extensions Marketplace on [Docker Desktop](https://open.docker.com/extensions/marketplace) and [Docker Hub](https://hub.docker.com/search?q=&type=extension).
The Extensions Marketplace is a space where developers can discover extensions to improve their developer experience and propose their own extension to be available for all Desktop users.

Whenever you are [ready to publish](./DISTRIBUTION.md) your extension in the Marketplace, you can [self-publish your extension](https://github.com/docker/extensions-submissions/issues/new?assignees=&labels=&template=1_automatic_review.yaml&title=%5BSubmission%5D%3A+)

> **Note**
>
> As the Extension Marketplace continues to add new features for both Extension users and publishers, you are expected
> to maintain your extension over time to ensure it stays available in the Marketplace.

> **Important**
>
> The Docker manual review process for extensions is paused at the moment. Submit your extension through the [automated submission process](https://github.com/docker/extensions-submissions/issues/new?assignees=&labels=&template=1_automatic_review.yaml&title=%5BSubmission%5D%3A+)
{ .important }
### Before you submit

Before you submit your extension, it must pass the [validation](validate.md) checks.

It is highly recommended that your extension follows the guidelines outlined in this section before submitting your
extension. If you request a review from the Docker Extensions team and have not followed the guidelines, the review process may take longer. 

These guidelines don't replace Docker's terms of service or guarantee approval:
- Review the [design guidelines](../design/design-guidelines.md)
- Ensure the [UI styling](../design/index.md) is in line with Docker Desktop guidelines
- Ensure your extensions support both light and dark mode
- Consider the needs of both new and existing users of your extension
- Test your extension with potential users
- Test your extension for crashes, bugs, and performance issues
- Test your extension on various platforms (Mac, Windows, Linux)
- Read the [Terms of Service](https://www.docker.com/legal/extensions_marketplace_developer_agreement/)

#### Validation process

Submitted extensions go through an automated validation process. If all the validation checks pass successfully, the extension is
published on the Marketplace and accessible to all users within a few hours.
It is the fastest way to get developers the tools they need and to get feedback from them as you work to
evolve/polish your extension.

> **Important**
>
> Docker Desktop caches the list of extensions available in the Marketplace for 12 hours. If you don't see your
> extension in the Marketplace, you can restart Docker Desktop to force the cache to refresh.
{ .important }
