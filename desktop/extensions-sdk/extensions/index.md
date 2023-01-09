---
title: Extension release process
description: General steps in how to publish an extension
keywords: Docker, Extensions, sdk, publish
---

This section describes how to make your extension available, and how to make your extension more visible so users can discover it and install it in a single click.

## Release your extension

You have developed your extension and tested it locally. You are now ready to release the extension and make it available for others to install and use (either internally with your team, or more publicly).

Releasing your extension consists of:

- Provide information about your extension: description, screenshots, etc. so users can decide to install your extension
- [Validate](./validate.md) the extension is built in the right format and includes the required information
- Make the extension image available on [Docker Hub](https://hub.docker.com/){: target="_blank" rel="noopener" class="_" }:

See [Package and release your extension](DISTRIBUTION.md) for details about the release process.

## Promote your extension

Once your extension is available on Docker Hub, users who have access to the extension image can install it using the Docker CLI.
However, you might want a better way to make your extension visible and easy to install than copy-pasting CLI commands in a terminal.

### Use a share extension link

You can [generate a share URL](share.md) in order to share your extension within your team, or promote your extension on the internet (website, blogs, social media...). This will allow users to view the extension description and screenshots, and decide to install it in a single click.

### Publish your extension in the Marketplace

You can publish your extension in the Extension marketplace to make it more discoverable. You must [submit your extension](publish.md) if you want to have it published in the marketplace.

## What happens next

### Extension new releases

Once you have released your extension, you can push a new release just by pushing a new version of the extension image, with an incremented tag (still using semver conventions).
Docker extensions published in the marketplace will benefit from update notifications to all Desktop users that have installed the extension. See more details about [new releases and updates](DISTRIBUTION.md#new-releases-and-updates)

### Extension support and user feedback

In addition to provide a description of you extension features, and screenshots, you should also specify additional URLs in [extension labels](labels.md). This will direct users to your website for reporting bugs and feedback, and accessing documentation and support.

{% include extensions-form.md %}

