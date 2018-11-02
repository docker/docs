---
title: Mirror images from another registry
description: Learn how to set up a repository to poll for changes in another registry and automatically pull new images from it.
keywords: registry, promotion, mirror
---

Docker Trusted Registry allows you to set up a mirror of a repository by
constantly polling it and pulling new image tags as they are pushed. This ensures your images are replicated across different registries for high availability. It also makes it easy to create a development pipeline that allows different
users access to a certain image without giving them access to everything in the remote registry.

![pull mirror](../../images/pull-mirror-1.svg)

To mirror a repository, start by [creating a repository](../manage-images/index.md)
in the DTR deployment that will serve as your mirror. Previously, you were only able to set up pull mirroring from the API. Starting in DTR 2.6, you can also mirror and pull from a remote DTR or Docker Hub repository. 

## Pull mirroring on the web interface

To get started:

1.  Navigate to `https://<dtr-url>` and log in with your UCP credentials. 

2.  Select **Repositories** on the left navigation pane, and then click on the name of the repository that you want to view. Note that you will have to click on the repository name following the `/` after the specific namespace for your repository.

3.  Select the **Mirrors** tab and click **New mirror policy**.

4. In the ***New Mirror*** page, specify the following details:
   * Mirror direction: Choose "Pull from remote registry"
   * Registry type: You can choose between **Docker Trusted Registry** and **Docker Hub**. If you choose DTR, enter your DTR URL. Otherwise, **Docker Hub** defaults to `https://index.docker.io`.
   * Username and Password or access token: Your credentials in the remote repository you wish to poll from
   * Repository: Enter the `namespace` and the `repository name after the `/`.
   * Show advanced settings: Enter the TLS details for the remote repository or check `Skip TLS verification`.
    ![](../../images/pull-mirror-1.png){: .img-fluid .with-border}

5. Click **Connect**.

> Known Issues
>
> For issues related to pull mirroring, see [DTR 2.6.0 Release Notes](../../release-notes).


## Pull mirroring on the API
The easiest way to interact with the DTR API is to use the interactive documentation
available from the web interface. Click **API** from the bottom left navigation pane.

Search for the endpoint:

```
POST /api/v0/repositories/{namespace}/{reponame}/pollMirroringPolicies
```

Click **Try it out** and enter your HTTP request details. `namespace` and `reponame` refer
to the repository that will be the mirror. The other fields refer to the remote repository to poll from and the credentials to use. As a best practice, use a service account just for this purpose. Instead of providing the password for that account, you should pass an
[authentication token](../access-tokens.md)

If the Docker Trusted or Hub registry to mirror images from is using self-signed certificates or
certificates signed by your own certificate authority, you also need to provide
the public key certificate for that certificate authority.
You can get it by accessing `https://<dtr-domain>/ca`.

Click **Execute**. On success, the API returns an `HTTP 201` response. This means
that the repository will be polling the source repository every couple of minutes.

## Where to go next

* [Mirror images to another registry](push-mirror.md)
