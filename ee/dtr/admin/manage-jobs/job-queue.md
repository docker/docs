---
title: Job Queue
description: Learn how Docker Trusted Registry runs batch jobs for troubleshooting job-related issues.
keywords: dtr, job queue, job management
---

Docker Trusted Registry (DTR) uses a job queue to schedule batch jobs. Jobs are added to a cluster-wide job queue, and then consumed and executed by a job runner within DTR.

![batch jobs diagram](../../images/troubleshoot-batch-jobs-1.svg)

All DTR replicas have access to the job queue, and have a job runner component
that can get and execute work.

## How it works

When a job is created, it is added to a cluster-wide job queue and enters the `waiting` state.
When one of the DTR replicas is ready to claim the job, it waits a random time of up
to `3` seconds to give every replica the opportunity to claim the task.

A replica claims a job by adding its replica ID to the job. That way, other
replicas will know the job has been claimed. Once a replica claims a job, it adds
that job to an internal queue, which in turn sorts the jobs by their `scheduledAt` time.
Once that happens, the replica updates the job status to `running`, and
starts executing it.

The job runner component of each DTR replica keeps a `heartbeatExpiration`
entry on the database that is shared by all replicas. If a replica becomes
unhealthy, other replicas notice the change and update the status of the failing worker to `dead`.
Also, all the jobs that were claimed by the unhealthy replica enter the `worker_dead` state,
so that other replicas can claim the job.

## Job Types

DTR runs periodic and long-running jobs. The following is a complete list of jobs you can filter for via [the user interface](../manage-jobs/audit-jobs-via-ui.md) or [the API](../manage-jobs/audit-jobs-via-api.md).   

| Job               | Description                                                                                                                                                                                                                                               |
|:------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| gc                | A garbage collection job that deletes layers associated with deleted images.                                                                                                                                                                                 |
| onlinegc                | A garbage collection job that deletes layers associated with deleted images without putting the registry in read-only mode.  |
| onlinegc_metadata                | A garbage collection job that deletes metadata associated with deleted images. |
| onlinegc_joblogs                | A garbage collection job that deletes job logs based on a configured job history setting. |
| metadatastoremigration   | A necessary migration that enables the `onlinegc` feature. |
| sleep             | Used for testing the correctness of the jobrunner. It sleeps for 60 seconds.                                                                                                                                                                           |
| false             | Used for testing the correctness of the jobrunner. It runs the `false` command and immediately fails.                                                                                                                                                 |
| tagmigration      | Used for synchronizing tag and manifest information between the DTR database and the storage backend.                                                                                                                                       |
| bloblinkmigration | A DTR 2.1 to 2.2 upgrade process that adds references for blobs to repositories in the database.                                                                                                                                          |
| license_update    | Checks for license expiration extensions if online license updates are enabled.                                                                                                                                                             |
| scan_check        | An image security scanning job. This job does not perform the actual scanning, rather it spawns `scan_check_single` jobs (one for each layer in the image). Once all of the `scan_check_single` jobs are complete, this job will terminate.                |
| scan_check_single | A security scanning job for a particular layer given by the `parameter: SHA256SUM`. This job breaks up the layer into components and checks each component for vulnerabilities.                                                                            |
| scan_check_all    | A security scanning job that updates all of the currently scanned images to display the latest vulnerabilities.                                                                                                                                            |
| update_vuln_db    | A job that is created to update DTR's vulnerability database. It uses an Internet connection to check for database updates through `https://dss-cve-updates.docker.com/` and updates the `dtr-scanningstore` container if there is a new update available. |
| scannedlayermigration  | A DTR 2.4 to 2.5 upgrade process that restructures scanned image data. |
| push_mirror_tag  | A job that pushes a tag to another registry after a push mirror policy has been evaluated. |
| poll_mirror  | A global cron that evaluates poll mirroring policies. |
| webhook           | A job that is used to dispatch a webhook payload to a single endpoint.                                                                                                                                                                                     |
| nautilus_update_db           | The old name for the `update_vuln_db` job. This may be visible on old log files.                                                                                                                                                                                   |
| ro_registry           | A user-initiated job for manually switching DTR into read-only mode.     |
| tag_pruning           | A job for cleaning up unnecessary or unwanted repository tags which can be configured by repository admins. For configuration options, see [Tag Pruning](../../user/tag-pruning).                                                                                                                                                                      |

## Job Status

Jobs can have one of the following status values:

| Status          | Description                                                                                                                               |
|:----------------|:------------------------------------------------------------------------------------------------------------------------------------------|
| waiting         | Unclaimed job waiting to be picked up by a worker.                                                                              |
| running         | The job is currently being run by the specified `workerID`.                                                                             |
| done            | The job has successfully completed.                                                                                                        |
| error           | The job has completed with errors.                                                                                                         |
| cancel_request  | The status of a job is monitored by the worker in the database. If the job status changes to `cancel_request`, the job is canceled by the worker. |
| cancel          | The job has been canceled and was not fully executed.                                                                                          |
| deleted         | The job and its logs have been removed.                                                                                                        |
| worker_dead     | The worker for this job has been declared `dead` and the job will not continue.                                                            |
| worker_shutdown | The worker that was running this job has been gracefully stopped.                                                                          |
| worker_resurrection | The worker for this job has reconnected to the database and will cancel this job.                                          |

## Where to go next

- [Audit Jobs via Web Interface](audit-jobs-via-ui)
- [Audit Jobs via API](audit-jobs-via-api)
