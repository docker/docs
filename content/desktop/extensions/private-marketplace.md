---
description: Extensions
keywords: Docker Extensions, Docker Desktop, Linux, Mac, Windows,
title: Configure an extensions private marketplace
---

This page shows administrators how to configure a private extensions marketplace with a curated list of extensions for your Docker Desktop users.

Extensions private marketplace is designed specifically for organizations who don’t give developers root access to their machines.

## Prerequisites

- [Download and install Docker Desktop 4.25.0 or later](https://docs.docker.com/desktop/release-notes/).
- Extensions private marketplace is available to Docker Business customers only.

Private marketplace is based on admin settings; see detailed docs about admin settings [here](https://docs.docker.com/desktop/hardened-desktop/settings-management/configure/).

## Initializing the private marketplace

Create a folder locally where you will create the content that will be deployed to the developers’ machines:

```bash
mkdir myMarketplace
cd myMarketplace
```

Initialize the configuration files for your marketplace:

{{< tabs >}}
{{< tab name="Mac" >}}

```bash
/Applications/Docker.app/Contents/Resources/bin/extension-admin init
```

{{< /tab >}}
{{< tab name="Windows" >}}

```bash
C:\Program Files\Docker\Docker\resources\bin\extension-admin init
```

{{< /tab >}}
{{< tab name="Linux" >}}

```bash
/opt/docker-desktop/extension-admin init
```

{{< /tab >}}
{{< /tabs >}}

This creates 2 files:

- `admin-settings.json`, which activates the private marketplace feature once it’s been applied to Docker Desktop on your developers’ machines,
- `extensions.txt`, which determines which extensions will be listed in your private marketplace.

## Changing private marketplace behaviour

The generated `admin-settings.json` file includes various settings you can modify.

Each setting has a `value` that you can set, including a `locked` field that allows you to lock the setting and make it unchangeable by your developers.

- `extensionsEnabled` enables Docker extensions.
- `extensionsPrivateMarketplace` activates the private marketplace and ensures Docker Desktop connect to content defined and controlled by the administrator instead of the public Docker marketplace.
- `onlyMarketplaceExtensions` allows or blocks users from installing other extensions by using the command line. Teams developing new extensions must have this setting unlocked (`"locked": false`) to install and test extensions being developed.
- `extensionsPrivateMarketplaceAdminContactURL` defines a contact link for users to request new extensions in the private marketplace. If empty then no link will be shown to your developers on Docker Desktop, otherwise this can be either an HTTP link or a “mailto:” link, e.g.:

  ```json
  "extensionsPrivateMarketplaceAdminContactURL": {
    "locked": true,
    "value": "mailto:admin@acme.com"
  }
  ```

You can find more information about the `admin-settings.json` file [here](https://docs.docker.com/desktop/hardened-desktop/settings-management/configure/).

## List allowed extensions for the private marketplace

The generated `extensions.txt` file defines the list of extensions that will be available in your private marketplace.

Each line in the file is an allowed extension, following one of these formats:

- `org/repo:tag` e.g. `docker/disk-usage-extension:0.2.8`,
- `org/repo` e.g. `docker/disk-usage-extension`, in which case the latest semver tag available for the image will be used.

You may also comment out lines with `#` so the extension will be ignored.

> **Note**
>
> This list can include different types of extension images: 
> 
> - extensions from the public marketplace or any public image stored in Docker Hub,
> - extension images stored in Docker Hub as private images,*
> - extension images stored in a private registry.*
> 
> *The administrator and the Docker Desktop users will need to be logged in and have pull access to these images.

> **Note**
> 
> Your developers will only be able to install the version of the extension that you’ve listed.

## Generate the private marketplace

Once the list in `extensions.txt` is ready, you can generate the marketplace:

{{< tabs >}}
{{< tab name="Mac" >}}

```bash
/Applications/Docker.app/Contents/Resources/bin/extension-admin generate
```

{{< /tab >}}
{{< tab name="Windows" >}}

```bash
C:\Program Files\Docker\Docker\resources\bin\extension-admin generate
```

{{< /tab >}}
{{< tab name="Linux" >}}

```bash
/opt/docker-desktop/extension-admin generate
```

{{< /tab >}}
{{< /tabs >}}

This will create an `extension-marketplace` directory and download the marketplace metadata for all the allowed extensions into it.

The marketplace content is generated from extension image info in the [same format as public Docker extensions](https://docs.docker.com/desktop/extensions-sdk/extensions/labels/) (image labels), and will include extension title, description, screenshots, links, etc. 

## Test the private marketplace setup on your Docker Desktop installation

As a best practice, we recommend you to try the private marketplace on your Docker Desktop installation:

{{< tabs >}}
{{< tab name="Mac" >}}

```bash
sudo /Applications/Docker.app/Contents/Resources/bin/extension-admin apply
```

{{< /tab >}}
{{< tab name="Windows (run as Admin)" >}}

```bash
C:\Program Files\Docker\Docker\resources\bin\extension-admin apply
```

{{< /tab >}}
{{< tab name="Linux" >}}

```bash
sudo /opt/docker-desktop/extension-admin apply
```

{{< /tab >}}
{{< /tabs >}}

This will copy the relevant generated files to the location where Docker Desktop reads its configuration files.

Finally, quit and reopen Docker Desktop to make it take your configuration files into account.

You must also sign in with an account attached to your organization as a Docker Business plan to enable the private extension marketplace.

When you click on the “Extensions” tab in the left panel you should see the private marketplace listing only the extensions you put in `extensions.txt`.

![Extensions Private Marketplace](/assets/images/extensions-private-marketplace.webp)

## Distribute the private marketplace to your developers’ machines

Once you’ve confirmed that the private marketplace configuration works, the final step is to distribute the files to the developers’ machines with the MDM software your organization uses e.g. [Jamf](https://www.jamf.com/).

It is assumed that you have the ability to push the `extension-marketplace` folder and `admin-settings.json` file to the locations specified above through a device management software such as [Jamf](https://www.jamf.com/).

It’s also necessary that your developers be logged in to Docker Desktop in order for the private marketplace configuration to take effect. As an administrator, you need to [configure a registry.json to enforce Docker Desktop sign-in](https://docs.docker.com/security/for-admins/configure-sign-in/).