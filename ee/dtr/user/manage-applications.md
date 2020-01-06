---
title: Manage applications
description: Learn how to manage applications in Docker Trusted Registry.
keywords: DTR, trusted registry, Docker apps
---

>{% include enterprise_label_shortform.md %}

With the introduction of [the experimental `app` plugin](/engine/reference/commandline/app/) to the Docker CLI, DTR has been enhanced to include application management. In DTR 2.7, you can push an app to your DTR repository and have an application be clearly distinguished from [individual and multi-architecture container images](/ee/dtr/user/manage-images/pull-and-push-images/#push-an-image/), as well as [plugins](/engine/reference/commandline/plugin_push/). When you push an application to DTR, you see two image tags:

| Image | Tag | Type | Under the hood |
|-------|-----|------|----------------|
| Invocation | `<app_tag>-invoc` | Container image represented by OS and architecture (e.g. `linux amd64`) | Uses Docker Engine. The Docker daemon is responsible for building and pushing the image. |
| Application with bundled components | `<app_tag>` | Application | Uses the app client to build and push the image. `docker app` is experimental on the Docker client. | 

Notice the app-specific tags, `app` and `app-invoc`, with scan results for the bundled components in the former and the invocation image in the latter. To view the scanning results for the bundled components, click "View Details" next to the `app` tag.

![](/ee/dtr/images/manage-applications-1.png){: .with-border}

Click on the image name or digest to see the vulnerabilities for that specific image.

![](/ee/dtr/images/manage-applications-2.png){: .with-border}

## Parity with existing repository and image features

The following repository and image management events also apply to applications:

- [Creation](/app/working-with-app/#initialize-and-deploy-a-new-docker-app-project-from-scratch/)
- [DTR pushes](/app/working-with-app/#push-the-app-to-dtr) 
- [Vulnerability scans](/ee/dtr/user/manage-images/scan-images-for-vulnerabilities/)
- [Vulnerability overrides](/ee/dtr/user/manage-images/override-a-vulnerability/) 
- [Deletion](/ee/dtr/user/manage-images/delete-images/)
- [Immutable tags](/ee/dtr/user/manage-images/prevent-tags-from-being-overwritten/)
- [Promotion policies](/ee/dtr/user/promotion-policies/)

### Limitations

- You cannot sign an application since the Notary signer cannot sign [OCI (Open Container Initiative)](https://github.com/opencontainers/image-spec/blob/master/spec.md) indices.
- Scanning-based policies do not take effect until after all images bundled in the application have been scanned. 
- Docker Content Trust (DCT) does not work for applications and multi-arch images, which are the same under the hood.

## Troubleshooting tips

### x509 certificate errors

```bash
fixing up "35.165.223.150/admin/lab-words:0.1.0" for push: failed to resolve "35.165.223.150/admin/lab-words:0.1.0-invoc", push the image to the registry before pushing the bundle: failed to do request: Head https://35.165.223.150/v2/admin/lab-words/manifests/0.1.0-invoc: x509: certificate signed by unknown authority
```

#### Workaround

Check that your DTR has been configured with your TLS certificate's Fully Qualified Domain Name (FQDN). See [Configure DTR](/ee/dtr/admin/install/#step-5-configure-dtr) for more details. For `docker app` testing purposes, you can pass the `--insecure-registries` option for pushing an application`.

```bash
docker app push hello-world --tag 35.165.223.150/admin/lab-words:0.1.0 --insecure-registries 35.165.223.150
35.165.223.150/admin/lab-words:0.1.0-invoc
Successfully pushed bundle to 35.165.223.150/admin/lab-words:0.1.0. Digest is sha256:bd1a813b6301939fa46e617f96711e0cca1e4065d2d724eb86abde6ef7b18e23.
```

## Known Issues

See [DTR 2.7 Release Notes - Known Issues](/ee/dtr/release-notes/#270) for known issues related to applications in DTR.

