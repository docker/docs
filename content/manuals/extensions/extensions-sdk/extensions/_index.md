---
title: "Part two: Publish"
description: General steps in how to publish an extension
keywords: Docker, Extensions, sdk, publish
aliases: 
 - /desktop/extensions-sdk/extensions/
weight: 40
---

This section describes how to make your extension available and more visible, so users can discover it and install it with a single click.

## Release your extension

After you have developed your extension and tested it locally, you are ready to release the extension and make it available for others to install and use (either internally with your team, or more publicly).

Releasing your extension consists of:

- Providing information about your extension: description, screenshots, etc. so users can decide to install your extension
- [Validating](validate.md) that the extension is built in the right format and includes the required information
- Making the extension image available on [Docker Hub](https://hub.docker.com/)

See [Package and release your extension](DISTRIBUTION.md) for more details about the release process.

## Promote your extension

Once your extension is available on Docker Hub, users who have access to the extension image can install it using the Docker CLI.

### Use a share extension link

You can also [generate a share URL](share.md) in order to share your extension within your team, or promote your extension on the internet. The share link lets users view the extension description and screenshots.

### Publish your extension in the Marketplace

You can publish your extension in the Extensions Marketplace to make it more discoverable. You must [submit your extension](publish.md) if you want to have it published in the Marketplace.

## What happens next

### New releases

Once you have released your extension, you can push a new release just by pushing a new version of the extension image, with an incremented tag (still using `semver` conventions).
Extensions published in the Marketplace benefit from update notifications to all Desktop users that have installed the extension. For more details, see [new releases and updates](DISTRIBUTION.md#new-releases-and-updates).

### Extension support and user feedback

In addition to providing a description of your extension's features and screenshots, you should also specify additional URLs using [extension labels](labels.md). This direct users to your website for reporting bugs and feedback, and accessing documentation and support.

{{% include "extensions-form.md" %}}
