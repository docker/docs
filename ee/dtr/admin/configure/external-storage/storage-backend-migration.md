---
title: Switch storage backends
description: Storage backend migration for Docker Trusted Registry
keywords: dtr, storage drivers, local volume, NFS, Azure, S3,
---

Starting in DTR 2.6, switching storage backends initializes a new metadata store and erases your existing tags. This helps facilitate online garbage collection, which has been introduced in 2.5 as an experimental feature. In earlier versions, DTR would subsequently start a `tagmigration` job to rebuild tag metadata from the file layout in the image layer store. This job has been discontinued for DTR 2.5.x (with garbage collection) and DTR 2.6, as your storage backend could get out of sync with your DTR metadata, like your manifests and existing repositories. As best practice, DTR storage backends and metadata should always be moved, backed up, and restored together.

## DTR 2.6.4 and above

In DTR 2.6.4, a new flag, `--storage-migrated`, [has been added to `docker/dtr reconfigure`](/reference/dtr/2.6/cli/reconfigure/) which lets you indicate the migration status of your storage data during a reconfigure. If you are not worried about losing your existing tags, you can skip the recommended steps below and [perform a reconfigure](/reference/dtr/2.6/cli/reconfigure/).

### Best practice for data migration

Docker recommends the following steps for your storage backend and metadata migration:

1. Disable garbage collection by selecting "Never" under **System > Garbage Collection**, so blobs referenced in the backup that you create continue to exist. See [Garbage collection](/ee/dtr/admin/configure/garbage-collection/) for more details. Make sure to keep it disabled while you're performing the metadata backup and migrating your storage data.

    ![](/ee/dtr/images/garbage-collection-0.png){: .img-fluid .with-border}

2. [Back up your existing metadata](/ee/dtr/admin/disaster-recovery/create-a-backup/#back-up-dtr-metadata). See [docker/dtr backup](/reference/dtr/2.6/cli/backup/) for CLI command description and options. 

3. Migrate the contents of your current storage backend to the new one you are switching to. For example, upload your current storage data to your new NFS server.

4. [Restore DTR from your backup](/ee/dtr/admin/disaster-recovery/restore-from-backup/) and specify your new storage backend. See [docker/dtr destroy](/reference/dtr/2.6/cli/destroy/) and [docker/dtr restore](/reference/dtr/2.6/cli/backup/) for CLI command descriptions and options.

5. With DTR restored from your backup and your storage data migrated to your new backend, garbage collect any dangling blobs using the following API request:

        ```bash
        curl -u <username>:$TOKEN -X POST "https://<dtr-url>/api/v0/jobs" -H "accept: application/json" -H "content-type: application/json" -d "{ \"action": \"onlinegc_blobs\" }"
        ``` 
        On success, you should get a `202 Accepted` response with a job `id` and other related details.
   
This ensures any blobs which are not referenced in your previously created backup get destroyed.
    
### Alternative option for data migration

- If you have a long maintenance window, you can skip some steps from above and do the following:

    1. Put DTR in "read-only" mode using the following API request:
   
        ```bash
        curl -u <username>:$TOKEN -X POST "https://<dtr-url>/api/v0/meta/settings" -H "accept: application/json" -H "content-type: application/json" -d "{ \"readOnlyRegistry\": true }"
        ``` 
        On success, you should get a `202 Accepted` response.

    2. Migrate the contents of your current storage backend to the new one you are switching to. For example, upload your current storage data to your new NFS server.

    3. [Reconfigure DTR](/reference/dtr/2.6/cli/reconfigure) while specifying the `--storage-migrated` flag to preserve your existing tags. 


## DTR 2.6.0-2.6.4 and DTR 2.5 (with experimental garbage collection)

Make sure to [perform a backup](/ee/dtr/admin/disaster-recovery/create-a-backup/#back-up-dtr-data) before you change your storage backend when running DTR 2.5 (with online garbage collection) and 2.6.0-2.6.3. If you encounter an issue with lost tags, refer to the following resources:
  * For changes to reconfigure and restore options in DTR 2.6, see [docker/dtr reconfigure](/reference/dtr/2.6/cli/reconfigure/) and [docker/dtr restore](/reference/dtr/2.6/cli/restore). 
  * For Docker's recommended recovery strategies, see [DTR 2.6 lost tags after reconfiguring storage](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage).
  * For NFS-specific changes, see [Use NFS](nfs.md). 
  * For S3-specific changes, see [Learn how to configure DTR with Amazon S3](s3.md).

Upgrade to [DTR 2.6.4](#dtr-264-and-above) and follow [best practice for data migration](#best-practice-for-data-migration) to avoid the wiped tags issue when moving from one NFS serverto another. 

## Where to go next

- [Use NFS](nfs.md)
- [Use S3](s3.md)
- CLI reference pages
  - [docker/dtr install](/reference/dtr/2.6/cli/install/)
  - [docker/dtr reconfigure](/reference/dtr/2.6/cli/reconfigure/)
