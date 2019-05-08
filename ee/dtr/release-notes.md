---
title: DTR release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Trusted Registry
keywords: docker trusted registry, whats new, release notes
toc_min: 1
toc_max: 2
redirect_from:
  - /datacenter/dtr/2.4/guides/release-notes/
  - /datacenter/dtr/2.5/guides/release-notes/
---

Here you can learn about new features, bug fixes, breaking changes, and
known issues for each DTR version.

You can then use [the upgrade instructions](admin/upgrade.md),
to upgrade your installation to the latest release.

* [Version 2.6](#version-26)
* [Version 2.5](#version-25)
* [Version 2.4](#version-24)

# Version 2.6

## 2.6.6 
(2019-5-6)

### Security

* Refer to [DTR image vulnerabilities](https://success.docker.com/article/dtr-image-vulnerabilities) for details regarding actions to be taken, timeline, and any status updates/issues/recommendations.

### Enhancements

* DTR now supports an option to keep your tag metadata when switching storage backends via the API. This is similar to the `--storage-migrated` option when performing an NFS reconfiguration via `docker run docker/dtr reconfigure --nfs-url ...`. (docker/dhe-deploy#10246)
    - To use this option, first write your current storage settings to a JSON file via `curl ... /api/v0/admin/settings/registry > storage.json`.
    - Next, add `keep_metadata: true` as a top-level key in the JSON you just created and modify it to contain your new storage settings.
    - Finally, update your Registry settings with your modified JSON file via `curl -X PUT .../api/v0/admin/settings/registry -d @storage.json`.
    
### Bug Fixes

* Fixed an issue where replica version was inferred from DTR volume labels. (docker/dhe-deploy#10266)

### Security
* Bumped the Golang version for DTR to 1.12.4. (docker/dhe-deploy#10290)
* Bumped the Alpine version of the base image to 3.9. (docker/dhe-deploy#10290)

### Known issues

* Docker Engine Enterprise Edition (Docker EE) Upgrade
  * There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before `18.09` to version `18.09` or greater. For DTR-specific changes, see [2.5 to 2.6 upgrade](/ee/dtr/admin/upgrade/#25-to-26-upgrade).
* Web Interface
  * Poll mirroring for Docker plugins such as `docker/imagefs` is currently broken. (docker/dhe-deploy #9490)
  * When viewing the details of a scanned image tag, the header may display a different vulnerability count from the layer details. (docker/dhe-deploy #9474)
  * In order to set a tag limit for pruning purposes, immutability must be turned off for a repository. This limitation is not clear in the **Repository Settings** view. (docker/dhe-deploy #9554)
* Webhooks
  * When configured for "Image promoted from repository" events, a webhook notification is triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
* System
  * When upgrading from `2.5` to `2.6`, the system will run a `metadatastoremigration` job after a successful upgrade. This is necessary for online garbage collection. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

## 2.6.5
(2019-4-11)

### Bug Fixes
* Fixed a bug where the web interface was not rendering for non-admin users.
* Removed `Users` tab from the side navigation [#10222](https://github.com/docker/dhe-deploy/pull/10222)

### Known issues

* Docker Engine Enterprise Edition (Docker EE) Upgrade
  * There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before `18.09` to version `18.09` or greater. For DTR-specific changes, see [2.5 to 2.6 upgrade](/ee/dtr/admin/upgrade/#25-to-26-upgrade).
* Web Interface
  * Poll mirroring for Docker plugins such as `docker/imagefs` is currently broken. (docker/dhe-deploy #9490)
  * When viewing the details of a scanned image tag, the header may display a different vulnerability count from the layer details. (docker/dhe-deploy #9474)
  * In order to set a tag limit for pruning purposes, immutability must be turned off for a repository. This limitation is not clear in the **Repository Settings** view. (docker/dhe-deploy #9554)
* Webhooks
  * When configured for "Image promoted from repository" events, a webhook notification is triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
* System
  * When upgrading from `2.5` to `2.6`, the system will run a `metadatastoremigration` job after a successful upgrade. This is necessary for online garbage collection. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

## 2.6.4 
(2019-3-28)

### Enhancements

* Added `--storage-migrated` option to reconfigure with migrated content when moving content to a new NFS URL. (ENGDTR-794)
* Added a job log status filter which allows users to exclude jobs that are not currently ***running***. (docker/dhe-deploy #10077)

### Bug Fixes

* If you have a repository in DTR 2.4 with manifest lists enabled, `docker pull` would fail on images that have been pushed to the repository after you upgrade to 2.5 and opt into garbage collection. This also applied when upgrading from 2.5 to 2.6. The issue has been fixed in DTR 2.6.4. (ENGDTR-330 and docker/dhe-deploy #10105) 

### Known issues

* Docker Engine Enterprise Edition (Docker EE) Upgrade
  * There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before `18.09` to version `18.09` or greater. For DTR-specific changes, see [2.5 to 2.6 upgrade](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

* Web Interface
  * Poll mirroring for Docker plugins such as `docker/imagefs` is currently broken. (docker/dhe-deploy #9490)
  * When viewing the details of a scanned image tag, the header may display a different vulnerability count from the layer details. (docker/dhe-deploy #9474)
  * In order to set a tag limit for pruning purposes, immutability must be turned off for a repository. This limitation is not clear in the **Repository Settings** view. (docker/dhe-deploy #9554)
* Webhooks
  * When configured for "Image promoted from repository" events, a webhook notification is triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
* System
  * When upgrading from `2.5` to `2.6`, the system will run a `metadatastoremigration` job after a successful upgrade. This is necessary for online garbage collection. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

## 2.6.3

(2019-2-28)

### Changelog

* Bump the Golang version that is used to build DTR to version 1.11.5. (docker/dhe-deploy#10060)

### Bug Fixes

* Users with read-only permissions can no longer see the README edit button for a repository. (docker/dhe-deploy#10056)

### Known issues

* Docker Engine Enterprise Edition (Docker EE) Upgrade
  * There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before `18.09` to version `18.09` or greater. For DTR-specific changes, see [2.5 to 2.6 upgrade](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

* Web Interface
  * Poll mirroring for Docker plugins such as `docker/imagefs` is currently broken. (docker/dhe-deploy #9490)
  * When viewing the details of a scanned image tag, the header may display a different vulnerability count from the layer details. (docker/dhe-deploy #9474)
  * In order to set a tag limit for pruning purposes, immutability must be turned off for a repository. This limitation is not clear in the **Repository Settings** view. (docker/dhe-deploy #9554)
  * Changing your S3 settings through the web interface will lead to erased metadata (ENGDTR-793). See [Restore to Cloud Storage](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#restoretocloudstorage) for Docker's recommended recovery strategy.

* CLI
  * When reconfiguring and restoring DTR, specifying `--nfs-storage-url` will assume you are switching to a fresh storage backend and will wipe your existing tags (ENGDTR-794). See [Reconfigure Using a Local NFS Volume](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#reconfigureusingalocalnfsvolume) and [Restore to a Local NFS Volume](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#restoretoalocalnfsvolume) for Docker's recommended recovery strategies.

* Webhooks
  * When configured for "Image promoted from repository" events, a webhook notification is triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)

* System
  * When upgrading from `2.5` to `2.6`, the system will run a `metadatastoremigration` job after a successful upgrade. This is necessary for online garbage collection. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

## 2.6.2 

(2019-1-29)

### Bug Fixes

* Fixed a bug where scanning Windows images were stuck in Pending state. (docker/dhe-deploy #9969)

### Known issues

* Docker Engine Enterprise Edition (Docker EE) Upgrade
  * There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before `18.09` to version `18.09` or greater. For DTR-specific changes, see [2.5 to 2.6 upgrade](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

* Web Interface
  * Users with read-only permissions to a repository can edit the repository README but their changes will not be saved. Only repository admins should have the ability to [edit the description](/ee/dtr/admin/manage-users/permission-levels/#team-permission-levels) of a repository. (docker/dhe-deploy #9677)
  * Poll mirroring for Docker plugins such as `docker/imagefs` is currently broken. (docker/dhe-deploy #9490)
  * When viewing the details of a scanned image tag, the header may display a different vulnerability count from the layer details. (docker/dhe-deploy #9474)
  * In order to set a tag limit for pruning purposes, immutability must be turned off for a repository. This limitation is not clear in the **Repository Settings** view. (docker/dhe-deploy #9554)
  * Changing your S3 settings through the web interface will lead to erased metadata (ENGDTR-793). See [Restore to Cloud Storage](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#restoretocloudstorage) for Docker's recommended recovery strategy.

* CLI
  * When reconfiguring and restoring DTR, specifying `--nfs-storage-url` will assume you are switching to a fresh storage backend and will wipe your existing tags (ENGDTR-794). See [Reconfigure Using a Local NFS Volume](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#reconfigureusingalocalnfsvolume) and [Restore to a Local NFS Volume](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#restoretoalocalnfsvolume) for Docker's recommended recovery strategies.

* Webhooks
  * When configured for "Image promoted from repository" events, a webhook notification is triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)

* System
  * When upgrading from `2.5` to `2.6`, the system will run a `metadatastoremigration` job after a successful upgrade. This is necessary for online garbage collection. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

## 2.6.1 

(2019-01-09)

### Bug Fixes

* Fixed a bug where notary signing data was not being backed up properly (docker/dhe-deploy #9862)
* Allow a cluster to go from 2 replicas to 1 without forcing removal (docker/dhe-deploy #9840)
* Fixed a race condition in initialization of the scan vulnerability database (docker/dhe-deploy #9907)

### Changelog
* GoLang version bump to 1.11.4.

### Known issues

* Docker Engine Enterprise Edition (Docker EE) Upgrade
  * There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before `18.09` to version `18.09` or greater. For DTR-specific changes, see [2.5 to 2.6 upgrade](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

* Web Interface
  * Users with read-only permissions to a repository can edit the repository README but their changes will not be saved. Only repository admins should have the ability to [edit the description](/ee/dtr/admin/manage-users/permission-levels/#team-permission-levels) of a repository. (docker/dhe-deploy #9677)
  * Poll mirroring for Docker plugins such as `docker/imagefs` is currently broken. (docker/dhe-deploy #9490)
  * When viewing the details of a scanned image tag, the header may display a different vulnerability count from the layer details. (docker/dhe-deploy #9474)
  * In order to set a tag limit for pruning purposes, immutability must be turned off for a repository. This limitation is not clear in the **Repository Settings** view. (docker/dhe-deploy #9554)
  * Changing your S3 settings through the web interface will lead to erased metadata (ENGDTR-793). See [Restore to Cloud Storage](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#restoretocloudstorage) for Docker's recommended recovery strategy.

* CLI
  * When reconfiguring and restoring DTR, specifying `--nfs-storage-url` will assume you are switching to a fresh storage backend and will wipe your existing tags (ENGDTR-794). See [Reconfigure Using a Local NFS Volume](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#reconfigureusingalocalnfsvolume) and [Restore to a Local NFS Volume](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#restoretoalocalnfsvolume) for  Docker's recommended recovery strategies.

* Webhooks
  * When configured for "Image promoted from repository" events, a webhook notification is triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)

* System
  * When upgrading from `2.5` to `2.6`, the system will run a `metadatastoremigration` job after a successful upgrade. This is necessary for online garbage collection. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

## 2.6.0 

(2018-11-08)

### New Features

* Web Interface
  * Online garbage collection is no longer an experimental feature. Users can now write to DTR and push images during garbage collection. [Learn about garbage collection](/ee/dtr/admin/configure/garbage-collection/).
  * Repository admins can now enable tag pruning for every repository that they manage by adding a pruning policy or setting a tag limit. [Learn about tag pruning](/ee/dtr/user/tag-pruning).
  * Users can now review and audit repository events on the web interface with the addition of the **Activity** tab on each repository. [Learn about repository event audits](/ee/dtr/user/audit-repository-events/).
  * DTR admins can now enable auto-deletion of repository events based on specified conditions. [Learn about repository event auto-deletion](/ee/dtr/admin/configure/auto-delete-repo-events/).
  * DTR admins can now review and audit jobs on the web interface with the addition of **Job Logs** within System settings. [Learn about job audits on the web interface](/ee/dtr/admin/manage-jobs/audit-jobs-via-ui/).
  * DTR admins can now enable auto-deletion of job logs based on specified conditions. [Learn about job log auto-deletion](/ee/dtr/admin/manage-jobs/auto-delete-job-logs/).
  * Users can now mirror images from another Docker Trusted or Docker Hub registry using the web interface. [Learn about pull mirroring](/ee/dtr/user/promotion-policies/pull-mirror).

* CLI
  * To support NFS v4, users can now pass additional options such as `--async-nfs` and `--nfs-options` when installing or reconfiguring NFS for external storage. See [docker/dtr install](/reference/dtr/2.6/cli/install) and [docker/dtr reconfigure](/reference/dtr/2.6/cli/reconfigure) for more details.
  * When installing and restoring DTR from an existing backup, users are now required to specify a storage flag: `--dtr-use-default-storage`, `--dtr-storage-volume`, or `--nfs-storage-url`. This ensures recovery of the configured storage setting when the backup was created. See [docker/dtr restore](/reference/dtr/2.6/cli/restore) for more details.

* API
  * Security admins can now export vulnerability scans to CSV via the `GET /api/v0/imagescan/scansummary/repositories/{namespace}/{reponame}/{tag}/export` endpoint. Specify `text/csv` as an Accept request HTTP header.
  * Repository admins can now interact with repository pruning policies using the following endpoints:
   * `GET /api/v0/repositories/{namespace}/{reponame}/pruningPolicies` 
   * `POST /api/v0/repositories/{namespace}/{reponame}/pruningPolicies`
   * `GET /api/v0/repositories/{namespace}/{reponame}/pruningPolicies/test`
   * `GET /api/v0/repositories/{namespace}/{reponame}/pruningPolicies/{pruningpolicyid}`
   * `GET /api/v0/repositories/{namespace}/{reponame}/pruningPolicies/{pruningpolicyid}`
   * `PUT /api/v0/repositories/{namespace}/{reponame}/pruningPolicies/{pruningpolicyid}`
   * `DELETE /api/v0/repositories/{namespace}/{reponame}/pruningPolicies/{pruningpolicyid}`

   See [Docker Trusted Registry API](../../reference/dtr/2.6/api/) for endpoint details and example usage. Alternatively, you can log in to the DTR web interface and select **API** from the bottom left navigation pane.

### Known issues

* Docker Engine Enterprise Edition (Docker EE) Upgrade
  * There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before `18.09` to version `18.09` or greater. For DTR-specific changes, see [2.5 to 2.6 upgrade](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

* Web Interface
  * Users with read-only permissions to a repository can edit the repository README but their changes will not be saved. Only repository admins should have the ability to [edit the description](/ee/dtr/admin/manage-users/permission-levels/#team-permission-levels) of a repository. (docker/dhe-deploy #9677)
  * Poll mirroring for Docker plugins such as `docker/imagefs` is currently broken. (docker/dhe-deploy #9490)
  * When viewing the details of a scanned image tag, the header may display a different vulnerability count from the layer details. (docker/dhe-deploy #9474)
  * In order to set a tag limit for pruning purposes, immutability must be turned off for a repository. This limitation is not clear in the **Repository Settings** view. (docker/dhe-deploy #9554)
  * Changing your S3 settings through the web interface will lead to erased metadata (ENGDTR-793). See [Restore to Cloud Storage](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#restoretocloudstorage) for Docker's recommended recovery strategy.

* CLI
  * When reconfiguring and restoring DTR, specifying `--nfs-storage-url` will assume you are switching to a fresh storage backend and will wipe your existing tags (ENGDTR-794). See [Reconfigure Using a Local NFS Volume](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#reconfigureusingalocalnfsvolume) and [Restore to a Local NFS Volume](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#restoretoalocalnfsvolume) for Docker's recommended recovery strategies.

* Webhooks
  * When configured for "Image promoted from repository" events, a webhook notification is triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)

* System
  * When upgrading from `2.5` to `2.6`, the system will run a `metadatastoremigration` job after a successful upgrade. This is necessary for online garbage collection. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](/ee/dtr/admin/upgrade/#25-to-26-upgrade).

### Deprecations

* API
  * `GET /api/v0/imagescan/repositories/{namespace}/{reponame}/{tag}` is deprecated in favor of `GET /api/v0/imagescan/scansummary/repositories/{namespace}/{reponame}/{tag}`. 
  * The following endpoints have been removed since online garbage collection will take care of these operations: 
    * `DELETE /api/v0/accounts/{namespace}/repositories` 
    * `DELETE /api/v0/repositories/{namespace}/{reponame}/manifests/{reference}`
  * The `enableManifestLists` field on the `POST /api/v0/repositories/{namespace}` endpoint will be removed in DTR 2.7. See [Deprecation Notice](deprecation-notice) for more details.


# Version 2.5


> **Important DTR Upgrade Information**
> If you have manifest lists enabled on any of your repositories:
>
> Upgrade path from 2.5.x to 2.6: Upgrade directly to 2.6.4.

## 2.5.11

(2019-05-06)

### Security

* Bumped the Golang version for DTR to 1.12.4. (docker/dhe-deploy #10301)
* Bumped the Alpine version of the base image to 3.9. (docker/dhe-deploy #10301)
* Bumped Python dependencies to address vulnerabilities. (docker/dhe-deploy #10308 and #10311)

### Bug Fixes

* Fixed an issue where read / write permissions were used when copying files into containers. (docker/dhe-deploy #10207)
* Fixed an issue where non-admin users could not access their repositories from the Repositories page on the web interface. (docker/dhe-deploy #10294)

### Known Issues

* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
  * When configured for "Image promoted from repository" events, a webhook notification will be triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

## 2.5.10

(2019-3-28)

### Bug Fixes

* If you have a repository in DTR 2.4 with manifest lists enabled, `docker pull` used to fail on images that were pushed to the repository after you upgraded to 2.5 and opted into garbage collection. This has been fixed in 2.5.10. (docker/dhe-deploy#10106)

### Known Issues
* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
  * When configured for "Image promoted from repository" events, a webhook notification will be triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

## 2.5.9

(2019-2-28)

### Changelog

* Bump the Golang version that is used to build DTR to version 1.10.8. (docker/dhe-deploy#10071)

### Known Issues
* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
  * When configured for "Image promoted from repository" events, a webhook notification will be triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

* System
  * When opting into online garbage collection, the system will run a `metadatastoremigration` job after a successful upgrade. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](../../v18.03/ee/dtr/admin/configure/garbage-collection/#metadata-store-migration).

## 2.5.8
 
(2019-1-29)

### Bug Fixes

* Fixed an issue that prevented vulnerability updates from running if they were previously interrupted. (docker/dhe-deploy #9958)

### Known Issues
* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
  * When configured for "Image promoted from repository" events, a webhook notification will be triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

* System
  * When opting into online garbage collection, the system will run a `metadatastoremigration` job after a successful upgrade. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](../../v18.03/ee/dtr/admin/configure/garbage-collection/#metadata-store-migration).

## 2.5.7 

(2019-01-09)

### Bug Fixes

* Fixed a bug where manifest lists were being appended to existing manifests lists when pushed. (docker/dhe-deploy #9811)
* Updated GoRethink library to avoid potential lock contention. (docker/dhe-deploy #9812)
* Fixed a bug where notary signing data was not being backed up properly. (docker/dhe-deploy #9851)

### Changelog
* GoLang version bump to 1.10.7.

### Known Issues
* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
  * When configured for "Image promoted from repository" events, a webhook notification will be triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

* System
  * When opting into online garbage collection, the system will run a `metadatastoremigration` job after a successful upgrade. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](../../v18.03/ee/dtr/admin/configure/garbage-collection/#metadata-store-migration).

## 2.5.6 

(2018-10-25)

### Bug Fixes
* Fixed a bug where Windows images could not be promoted. (docker/dhe-deploy#9215)
* Removed Python3 from base image. (docker/dhe-deploy#9219)
* Added CSP (docker/dhe-deploy#9366)
* Included foreign layers in scanned images. (docker/dhe-deploy#9488)
* Added dotnet.marsu to nautilus base image. (docker/dhe-deploy#9503)
* Backported ManifestList fixes. (docker/dhe-deploy#9547)
* Removed support sidebar link and associated content. (docker/dhe-deploy#9411)

### Known Issues
* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
  * When configured for "Image promoted from repository" events, a webhook notification will be triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

* System
  * When opting into online garbage collection, the system will run a `metadatastoremigration` job after a successful upgrade. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](../../v18.03/ee/dtr/admin/configure/garbage-collection/#metadata-store-migration).

## 2.5.5 

(2018-8-30)

### Bug Fixes

* Fixed bug where repository tag list UI was not loading after a tag migration.
* Fixed bug to enable poll mirroring with Windows images.
* The RethinkDB image has been patched to remove unused components with known vulnerabilities including the RethinkCLI. To get an equivalent interface, run RethinkCLI from a separate image using `docker run -it --rm --net dtr-ol -v dtr-ca-$REPLICA_ID:/ca dockerhubenterprise/rethinkcli:v2.3.0 $REPLICA_ID`.

### Known Issues
* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
  * When configured for "Image promoted from repository" events, a webhook notification will be triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

* System
  * When opting into online garbage collection, the system will run a `metadatastoremigration` job after a successful upgrade. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](../../v18.03/ee/dtr/admin/configure/garbage-collection/#metadata-store-migration).

## 2.5.3 

(2018-6-21)

### New Features

* Allow users to adjust DTR log levels for alternative logging solutions.

### Bug Fixes

* Fixed URL redirect to release notes.
* Prevent OOM during garbage collection by reading less data into memory at a time.
* Fixed issue where worker capacities wouldn't update on minor version upgrades.

### Known Issues
* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
  * When configured for "Image promoted from repository" events, a webhook notification will be triggered twice during an image promotion when scanning is enabled on a repository. (docker/dhe-deploy #9685)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

* System
  * When opting into online garbage collection, the system will run a `metadatastoremigration` job after a successful upgrade. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](../../v18.03/ee/dtr/admin/configure/garbage-collection/#metadata-store-migration).


## 2.5.2 

(2018-5-21)

### Bug fixes

* Fixed a problem where promotion policies based on scanning results would not be executed correctly.

### Known issues

* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

* System
  * When opting into online garbage collection, the system will run a `metadatastoremigration` job after a successful upgrade. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](../../v18.03/ee/dtr/admin/configure/garbage-collection/#metadata-store-migration).

## 2.5.1 

(2018-5-17)

### New features

* Headers added to all API and registry responses to improve security (enforce HTST, XSS Protection, prevent MIME sniffing).

### Bug fixes

* Allow for AlibabaCloud as storage backend.
* Fix a problem that made pulling images from Google Cloud fail when DTR was configured to redirect requests.
* Avoid sending redundant webhooks and fix inaccurate repository pull/push counts when manifest lists are pushed.
* Several fixes of common workflows when the experimental online garbage collection is enabled, including:
  * Support scanning.
  * Adding event stream items for online garbage collection activity like layers being deleted.
  * Fix failing repositories promotion policies.
  * Fix inaccurate pull/push counts.
* Some internationalization fixes.
* Fix a bug causing poll mirroring from Docker Hub to fail under certain conditions.
* Copy existing scan results to new target repository when an image is promoted.
* Address an issue causing scan results to not be available for images with long names.
* Remove a race condition in which repositories deleted during tagmigration were causing tagmigration to fail.
* Enhancements to the mirroring interface including:
  * Fixed URL for the destination repository.
  * Option to skip TLS verification when testing mirroring.
  
  ### Known issues

* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

* System
  * When opting into online garbage collection, the system will run a `metadatastoremigration` job after a successful upgrade. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](../../v18.03/ee/dtr/admin/configure/garbage-collection/#metadata-store-migration).

## 2.5.0 

(2018-4-17)

### New features

* You can now configure DTR to automatically create a new repository when
users push to a repository in their personal namespace that doesn't exist yet.
This makes the behavior of DTR consistent with Docker Hub. By default this
setting is disabled, so that DTR continues behaving the same way after an upgrade.
[Learn about creating repositories on push](admin/configure/allow-creation-on-push.md).
* You can create push mirroring policies to automatically push an image to
another DTR deployment or Docker Hub, when the image complies with a policy
of your choice.
[Learn about push mirroring](user/promotion-policies/push-mirror.md).
* You can configure a repository in a DTR deployment to mirror a repository
in a different DTR deployment by constantly monitoring it and pulling new
images when they are available.
[Learn about pull mirroring](user/promotion-policies/pull-mirror.md).
* Added the `emergency-repair` command to the DTR CLI tool. This allows you to
recover your DTR cluster from a loss of quorum and is an alternative to
restoring from a backup.
[Learn about the emergency-repair command](admin/disaster-recovery/repair-a-cluster.md).
* Users can now create access tokens that can be used to authenticate in the
DTR API without providing their credentials.
[Learn about access tokens](user/access-tokens.md).
* You can now configure DTR to run garbage collection jobs without putting DTR
into read-only mode. This is still experimental.
[Learn about garbage collection](admin/configure/garbage-collection.md).
* Administrators can hide vulnerabilities in given image layers if they
know that the vulnerability has been fixed.
[Learn how to override vulnerability reports](user/manage-images/override-a-vulnerability.md)
* You can now connect one DTR deployment to multiple UCPs, allowing you to
use Docker Content Trust in a seamless way between multiple UCPs.
* Added new endpoints to the DTR API to query the results of the Vulnerability
scanner:
  * `/api/v0/imagescan/scansummary/repositories/{namespace}/{reponame}/{tag}` returns
  the scanning summary for a given tag.
  * `/api/v0/imagescan/scansummary/cve/{cve}` gets the scan summary by CVE.
  * `/api/v0/imagescan/scansummary/layer/{layerid}` gets the scan summary by layer SHA.
  * `/api/v0/imagescan/scansummary/license/{license}` gets the scan summary by
  license type.
  * `/api/v0/imagescan/scansummary/component/{component}` get the scan summary by
  component.
* The API endpoint `/api/v0/repositories/{namespace}/{reponame}/manifests/{reference}`
has been deprecated. Use `/api/v0/repositories/{namespace}/{reponame}/tags/{tag}`
instead.

### Bug fixes

* Web Interface
  * Several improvements to the web interface to make it more stable
* User accounts
  * When a user changes their password they are automatically logged out.
* Vulnerability scanner
  * Fixed problem causing errors when trying to view scanning information when
an image has not been scanned yet.
* docker/dtr tool
  * When using `docker/dtr reconfigure --log-host`, you now need to also
specify `--log-protocol`.
  * You can now tune the RethinkDB cache size for improved performance. Use the
  `--replica-rethinkdb-cache-mb` option available on install, join, or reconfigure.
* Misc
  * Removed support for manifest schema v1. This doesn't  affect users.

### Known issues

* Web Interface
  * The web interface shows "This repository has no tags" in repositories where tags
  have long names. As a workaround, reduce the length of the name for the
  repository and tag.
  * When deleting a repository with signed images, the DTR web interface no longer
  shows instructions on how to delete trust data.
  * There's no web interface support to update mirroring policies when rotating the TLS
  certificates used by DTR. Use the API instead.
  * The web interface for promotion policies is currently broken if you have a large number
  of repositories.
  * Clicking "Save & Apply" on a promotion policy doesn't work.
* Webhooks
  * There is no webhook event for when an image is pulled.
  * HTTPS webhooks do not go through HTTPS proxy when configured. (docker/dhe-deploy #9492)
* Online garbage collection
  * The events API won't report events when tags and manifests are deleted.
  * The events API won't report blobs deleted by the garbage collection job.
* Docker EE Advanced features
  * Scanning any new push after metadatastore migration will not yet work.
  * Pushes to repos with promotion policies (repo as source) are broken when an
  image has a layer over 100MB.
  * On upgrade the scanningstore container may restart with this error message:
  FATAL:  database files are incompatible with server

* System
  * When opting into online garbage collection, the system will run a `metadatastoremigration` job after a successful upgrade. If the three system attempts fail, you will have to retrigger the `metadatastoremigration` job manually. [Learn about manual metadata store migration](../../v18.03/ee/dtr/admin/configure/garbage-collection/#metadata-store-migration).

# Version 2.4

> **Important DTR Upgrade Information**
> If you have manifest lists enabled on any of your repositories:
>
> Upgrade path from 2.4.x to 2.5: Do not opt into garbage collection, or directly upgrade to 2.5.10 if you need to opt into > garbage collection.
> Upgrade path from 2.5.x to 2.6: Upgrade directly to 2.6.4.

## 2.4.12

(2019-05-06)

### Security

* Bumped the Golang version for DTR to 1.12.4. [docker/dhe-deploy #10303](https://github.com/docker/dhe-deploy/pull/10303)
* Bumped Python dependencies to address vulnerabilities. [docker/dhe-deploy#10309](https://github.com/docker/dhe-deploy/pull/10309)

## 2.4.11

(2019-4-11)

### Changelog

* Bumped the Golang version that is used to build DTR to version 1.11.5. [docker/dhe-deploy#10155](https://github.com/docker/dhe-deploy/pull/10155)

## 2.4.10

(2019-2-28)

### Changelog

* Bump the Golang version that is used to build DTR to version 1.10.8. (docker/dhe-deploy#10068)

**Known issues**

* Backup uses too much memory and can cause out of memory issues for large databases.
* The `--nfs-storage-url` option uses the system's default NFS version instead
of testing the server to find which version works.


## Version 2.4.8

(2019-01-29)

### Changelog
* GoLang version bump to 1.10.6.

**Known issues**

* Backup uses too much memory and can cause out of memory issues for large databases.
* The `--nfs-storage-url` option uses the system's default NFS version instead
of testing the server to find which version works.


## Version 2.4.7

(2018-10-25)

### Bug Fixes
* Added CSP (Content Security Policy). (docker/dhe-deploy#9367 and docker/dhe-deploy#9584)
* Fixed critical vulnerability in RethinkDB. (docker/dhe-deploy#9574)

### Changelog
* Patched security vulnerabilities in the load balancer.
* Patch packages and base OS to eliminate and address some critical vulnerabilities in DTR dependencies.

**Known issues**

* Backup uses too much memory and can cause out of memory issues for large databases.
* The `--nfs-storage-url` option uses the system's default NFS version instead
of testing the server to find which version works.

## Version 2.4.6

(2018-07-26)

### Bug Fixes
* Fixed bug where repository tag list UI was not loading after a tag migration.
* The RethinkDB image has been patched to remove unused components with known vulnerabilities including the rethinkcli. To get an equivalent interface please run the rethinkcli from a separate image using `docker run -it --rm --net dtr-ol -v dtr-ca-$REPLICA_ID:/ca dockerhubenterprise/rethinkcli $REPLICA_ID`.

**Known issues**

* Backup uses too much memory and can cause out of memory issues for large databases.
* The `--nfs-storage-url` option uses the system's default NFS version instead
of testing the server to find which version works.

## Version 2.4.5

(2018-06-21)

**New Features**

* Allow users to adjust DTR log levels for alternative logging solutions.

**Bug Fixes**

* Prevent OOM during garbage collection by reading less data into memory at a time.

**Known issues**

* Backup uses too much memory and can cause out of memory issues for large databases.
* The `--nfs-storage-url` option uses the system's default NFS version instead
of testing the server to find which version works.

## Version 2.4.4

(2018-05-17)

**New features**

* Headers added to all API and registry responses to improve security (enforce HTST, XSS Protection, prevent MIME sniffing).

**Bug fixes**

* Fixed a problem that made pulling images from Google Cloud fail when DTR was configured to redirect requests.
* Remove a race condition in which repos deleted during tagmigration were causing tagmigration to fail.
* Reduce noise in the jobrunner logs by changing some of the more detailed messages to debug level.
* Eliminate a race condition in which webhook for license updates doesn't fire.

**Known issues**

* Backup uses too much memory and can cause out of memory issues for large databases.
* The `--nfs-storage-url` option uses the system's default NFS version instead
of testing the server to find which version works.

## Version 2.4.3 

(2018-03-19)

**Security notice**

* Dependencies updated to consume upstream CVE patches.

## Version 2.4.2 

(2018-02-13)

**Security notice**

The log driver is now disabled for containers started by backup and HA cluster
join operations. This is a critical security fix for customers that rely on
Docker Trusted Registry 2.2, 2.3 and 2.4 with a log driver to capture logs from
all containers across the platform.

Caution is advised when applying this update, make sure you redeploy DTR, and in
the process you will create new credentials because the previous ones were
potentially disclosed due to the vulnerability.

Use the `--log-driver=none` option for `docker run` when running a DTR backup, HA
cluster join or dumpcerts.

## 2.4.1 

(2017-11-20)

**Bug fixes**

* Fixed a bug that cause certain vulnerabilities to not be found during scanning.
* Increased speed of lock expiration in case of failed joins.
* Fixed notification when toggling active status of webhooks.
* Speed up detection of dead jobrunners.
* Fixed a bug where garbage collection ran in a suboptimal mode if scheduled as
a cron from the UI.
* Fixed a potential issue with the way we untar files in uploads of the
vulnerability database.
* Fixed scanning issue with some windows images.
* Fixed a bug with not backing up repository team permissions correctly.

**General improvements**

* Improved resilience of garbage collection.
* Improved logging of garbage collection.
* Improved memory usage during backup.
* Improved error handling when uploading invalid vulnerability databases.
* Improve resilience of DTR join operations.
* Hide secrets on storage config pages.

**Deprecations**

* The `api/v0/imagescan/layer/{layerid}` endpoint is deprecated, and will be
removed in DTR 2.5. You can use the
`/api/v0/imagescan/repositories/{namespace}/{reponame}/{tag}` endpoint instead.

**Known issues**

* Backup uses too much memory and can cause out of memory issues for large databases.
* The `--nfs-storage-url` option uses the system's default NFS version instead
of testing the server to find which version works.


## DTR 2.4.0 

(2017-11-2)

**New features**

* Upgraded to Swagger 2.0 and Swagger UI 3.0.
* DTR can now be deployed on IBM Z (s390x architecture).
* Updated the `docker/dtr-rethink` images to include `rethinkcli` for easier troubleshooting.
* Notary now allows you to see audit logs using the `/v2/_trust/changefeed`, and
`/v2/<repository>/_trust/changefeed` endpoints.

**Bug fixes**

* When setting up periodic garbage collection, it used to run in a different,
less thorough mode than when run manually. Now garbage collection always runs in
the correct mode.
* Fix error when garbage collecting manifest lists.
* Fixed issue when reconfiguring DTR from a non-local storage to NFS, causing the
change to not be persisted.
* Backported Docker Distribution race fix. [2299](https://github.com/docker/distribution/pull/2299)
* Reduced unnecessary logs in Jobrunner.
* Other general reliability improvements.

**Known issues**

* Backup uses too much memory and can cause out of memory issues for large databases.
* The `--nfs-storage-url` option uses the system's default NFS version instead
of testing the server to find which version works.

## Earlier versions

- [DTR 2.3 release notes](/datacenter/dtr/2.3/guides/release-notes.md)
- [DTR 2.2 release notes](/datacenter/dtr/2.2/guides/release-notes/index.md)
- [DTR 2.1 release notes](/datacenter/dtr/2.1/guides/release-notes.md)
- [DTR 2.0 release notes](/datacenter/dtr/2.0/release-notes/index.md)
