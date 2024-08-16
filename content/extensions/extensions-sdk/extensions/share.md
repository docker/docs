---
title: Share your extension
description: Share your extension with a share link
keywords: Docker, extensions, share
aliases: 
 - /desktop/extensions-sdk/extensions/share/
---

Once your extension image is accessible on Docker Hub, anyone with access to the image can install the extension.

People can install your extension by typing `docker extension install my/awesome-extension:latest` in to the terminal.

However, this option doesn't provide a preview of the extension before it's installed.

## Createa a share URL

Docker lets you share your extensions using a URL.

When people navigate to this URL, it opens Docker Desktop and displays a preview of your extension in the same way as an extension in the Marketplace. From the preview, users can then select **Install**.

![Navigate to extension link](images/open-share.png)

To generate this link you can either:

- Run the following command:

  ```console
  $ docker extension share my/awesome-extension:0.0.1
  ```

- Once you have installed your extension locally, navigate to the **Manage** tab and select **Share**.

  ![Share button](images/list-preview.png)

> [!NOTE]
>
> Previews of the extension description or screenshots, for example, are created using [extension labels](labels.md).
