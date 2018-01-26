---
title: Troubleshoot batch jobs
description: Learn how Docker Trusted Registry run batch jobs, so that you can troubleshoot when something goes wrong
keywords: dtr, troubleshoot
---

DTR uses a job queue to schedule batch jobs. A job is placed on this work queue,
and a job runner component of DTR consumes work from this cluster-wide job
queue and executes it.

![batch jobs diagram](../../images/troubleshoot-batch-jobs-1.svg)

All DTR replicas have access to the job queue, and have a job runner component
that can get and execute work.

## How it works

When a job is created, it is added to a cluster-wide job queue with the
`waiting` status.
When one of the DTR replicas is ready to claim, it waits a random time of up
to 3 seconds, giving the opportunity to every replica to claim the task.

A replica gets a job by adding it's replica ID to the job. That way, other
replicas know the job has been claimed. Once a replica claims a job it adds
it to an internal queue of jobs, that is sorted by their `scheduledAt` time.
When that time happens, the replica updates the job status to `running`, and
starts executing it.

The job runner component of each DTR replica keeps an `heartbeatExpiration`
entry on the database shared by all replicas. If a replica becomes
unhealthy, other replicas notice this and update that worker status to `dead`.
Also, all the jobs that replica had claimed are updated to the status `worker_dead`,
so that other replicas can claim the job.

## Job types

DTR has several types of jobs.

| Job               | Description                                                                                                                                                                                                                                               |
|:------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| gc                | Garbage collection job that deletes layers associated with deleted images                                                                                                                                                                                 |
| sleep             | Sleep is used to test the correctness of the jobrunner. It sleeps for 60 seconds                                                                                                                                                                          |
| false             | False is used to test the correctness of the jobrunner. It runs the `false` command and immediately fails                                                                                                                                                 |
| tagmigration      | Tag migration is used to synchronize tag and manifest information between the DTR database and the storage backend.                                                                                                                                       |
| bloblinkmigration | bloblinkmigration is a 2.1 to 2.1 upgrade process that adds references for blobs to repositories in the database                                                                                                                                          |
| license_update    | License update checks for license expiration extensions if online license updates are enabled                                                                                                                                                             |
| scan_check        | An image security scanning job. This job does not perform the actual scanning, rather it spawns `scan_check_single` jobs (one for each layer in the image). Once all of the `scan_check_single` jobs are complete, this job will terminate                |
| scan_check_single | A security scanning job for a particular layer given by the `parameter: SHA256SUM`. This job breaks up the layer into components and checks each component for vulnerabilities                                                                            |
| scan_check_all    | A security scanning job that updates all of the currently scanned images to display the latest vulnerabilities                                                                                                                                            |
| update_vuln_db    | A job that is created to update DTR's vulnerability database. It uses an Internet connection to check for database updates through `https://dss-cve-updates.docker.com/` and updates the `dtr-scanningstore` container if there is a new update available |
| webhook           | A job that is used to dispatch a webhook payload to a single endpoint                                                                                                                                                                                     |

## Job status

Jobs can be in one of the following status:

| Status          | Description                                                                                                                               |
|:----------------|:------------------------------------------------------------------------------------------------------------------------------------------|
| waiting         | The job is unclaimed and waiting to be picked up by a worker                                                                              |
| running         | The worker defined by `workerID` is currently running the job                                                                             |
| done            | The job has successfully completed                                                                                                        |
| error           | The job has completed with errors                                                                                                         |
| cancel_request  | The worker monitors the job statuses in the database. If the status for a job changes to `cancel_request`, the worker will cancel the job |
| cancel          | The job has been cancelled and not fully executed                                                                                         |
| deleted         | The job and logs have been removed                                                                                                        |
| worker_dead     | The worker for this job has been declared `dead` and the job will not continue                                                            |
| worker_shutdown | The worker that was running this job has been gracefully stopped                                                                          |
| worker          | resurrection| The worker for this job has reconnected to the database and will cancel these jobs                                          |

## Job capacity

Each job runner has a limited capacity and doesn't claim jobs that require an
higher capacity. You can see the capacity of a job runner using the
`GET /api/v0/workers` endpoint:

```json
{
  "workers": [
    {
      "id": "000000000000",
      "status": "running",
      "capacityMap": {
        "scan": 1,
        "scanCheck": 1
      },
      "heartbeatExpiration": "2017-02-18T00:51:02Z"
    }
  ]
}
```

This means that the worker with replica ID `000000000000` has a capacity of 1
`scan` and 1 `scanCheck`. If this worker notices that the following jobs
are available:

```json
{
  "jobs": [
    {
      "id": "0",
      "workerID": "",
      "status": "waiting",
      "capacityMap": {
        "scan": 1
      }
    },
    {
       "id": "1",
       "workerID": "",
       "status": "waiting",
       "capacityMap": {
         "scan": 1
       }
    },
    {
     "id": "2",
      "workerID": "",
      "status": "waiting",
      "capacityMap": {
        "scanCheck": 1
      }
    }
  ]
}
```

Our worker can pick up job id `0` and `2` since it has the capacity
for both, while id `1` needs to wait until the previous scan job is complete:

```json
{
  "jobs": [
    {
      "id": "0",
      "workerID": "000000000000",
      "status": "running",
      "capacityMap": {
        "scan": 1
      }
    },
    {
       "id": "1",
       "workerID": "",
       "status": "waiting",
       "capacityMap": {
         "scan": 1
       }
    },
    {
     "id": "2",
      "workerID": "000000000000",
      "status": "running",
      "capacityMap": {
        "scanCheck": 1
      }
    }
  ]
}
```

## Troubleshoot jobs

You can get the list of jobs, using the `GET /api/v0/jobs/` endpoint. Each job
looks like this:

```json
{
	"id": "1fcf4c0f-ff3b-471a-8839-5dcb631b2f7b",
	"retryFromID": "1fcf4c0f-ff3b-471a-8839-5dcb631b2f7b",
	"workerID": "000000000000",
	"status": "done",
	"scheduledAt": "2017-02-17T01:09:47.771Z",
	"lastUpdated": "2017-02-17T01:10:14.117Z",
	"action": "scan_check_single",
	"retriesLeft": 0,
	"retriesTotal": 0,
	"capacityMap": {
      	  "scan": 1
	},
	"parameters": {
      	  "SHA256SUM": "1bacd3c8ccb1f15609a10bd4a403831d0ec0b354438ddbf644c95c5d54f8eb13"
	},
	"deadline": "",
	"stopTimeout": ""
}
```

The fields of interest here are:

* `id`: the ID of the job
* `workerID`: the ID of the worker in a DTR replica that is running this job
* `status`: the current state of the job
* `action`: what job the worker will actually perform
* `capacityMap`: the available capacity a worker needs for this job to run


## Cron jobs

Several of the jobs performed by DTR are run in a recurrent schedule. You can
see those jobs using the `GET /api/v0/crons` endpoint:


```json
{
  "crons": [
    {
      "id": "48875b1b-5006-48f5-9f3c-af9fbdd82255",
      "action": "license_update",
      "schedule": "57 54 3 * * *",
      "retries": 2,
      "capacityMap": null,
      "parameters": null,
      "deadline": "",
      "stopTimeout": "",
      "nextRun": "2017-02-22T03:54:57Z"
    },
    {
      "id": "b1c1e61e-1e74-4677-8e4a-2a7dacefffdc",
      "action": "update_db",
      "schedule": "0 0 3 * * *",
      "retries": 0,
      "capacityMap": null,
      "parameters": null,
      "deadline": "",
      "stopTimeout": "",
      "nextRun": "2017-02-22T03:00:00Z"
    }
  ]
}
```

The `schedule` uses a Unix crontab syntax.
