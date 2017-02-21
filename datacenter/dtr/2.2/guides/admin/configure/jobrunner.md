---
title: Jobrunner
description: Learn about the inner-workings of the jobrunner container in the DTR workflow.
keywords: docker, job, runner
---

The jobrunner container is a DTR mechanism that:

1. Consumes jobs from a cluster-wide jobs queue 
2. Performs the work of the given action

There is one jobrunner container per replica. 

## Jobrunner Workflow
[//]: # (uncomment once diagrams are complete) The following diagram depicts the behavior of the jobrunner:

[//]: # (Placeholder for jobrunner diagram. @sarahpark will work on the diagram)

When a job is scheduled (see [Job Actions](#job-actions) below) it is put onto the 
cluster-wide jobs queue with an initial status of `waiting`. When a jobrunner worker 
is available to pick up the job, it will claim the it (i.e: the workerID will be set 
to the replicaID of the jobrunner container that claimed the job) and set the job 
status to `running`. The worker will carry out the job and then set the appropriate 
status when complete (see [Job Statuses](#job-statuses) below).

[//]: # (Placeholder for jobrunner scheduling. @sarahpark will work on the diagram)

Each jobrunner has an internal queue of `waiting` jobs sorted by their `scheduledAt`
time (from earliest to latest). When a worker is free to claim the next job, it claims
it after a delay of up to 3 seconds. This delay is imposed so that each available worker
has a chance to compete for the job. The worker that was successfully able to claim the
job will set the job's `workerID` and all other workers will drop the job from their
internal queue. If a job cannot be claimed due to capacity limits (see [Job Capacities](#job-capacities))
then it is placed into a separate queue and will go through the claiming process
when the worker has enough free capacity for it.

Jobrunners monitor each other's `heartbeatExpiration` in the workers table. When a worker 
see that another worker hasn't updated its expiration in a long time, it sets the 
dead worker's status to `dead` and its jobs to `worker_dead`. If the dead worker is able to
reconnect to the database and notices that it's jobs have been set to `worker_dead`,
it sets those job statuses to `worker_resurrection` and cancels them.

### Job Actions
The available job actions are:

- `gc`
: Garbage collection deletes layers associated with deleted images.
- `sleep`
: Sleep is used to test the correctness of the jobrunner. It sleeps for 60 seconds.
- `false`
: False is used to test the correctness of the jobrunner. It runs the `false` command and immediately fails.
- `tagmigration`
: Tag migration is used to sync tag and manifest information from the blobstore into the database.
This information is used to for information in the API, UI, and also for GC.
- `bloblinkmigration`
: bloblinkmigration is a 2.1 to 2.1 upgrade process that adds references for blobs to repositories in the database.
- `license_update`
: License update checks for license expiration extensions if online license updates are enabled.
- `nautilus_scan_check`
: An image security scanning job. This job does not perform the actual scanning, rather it
spawns `nautilus_scan_check_single` jobs (one for each layer in the image). Once
all of the `nautilus_scan_check_single` jobs are complete, this job will terminate.
- `nautilus_scan_check_single`
: A security scanning job for a particular layer given by the `parameter: SHA256SUM`. This job
breaks up the layer into components and checks each component for vulnerabilities
(see [Security Scanning](../../user/manage-images/scan-images-for-vulnerabilities.md)).
- `nautilus_update_db`
: A job that is created to update DTR's vulnerability database. It uses an
Internet connection to check for database updates through `https://dss-cve-updates.docker.com/` and
updates the dtr-scanningstore container if there is a new update available (see [Set up vulnerability scans](set-up-vulnerability-scans.md)). 
- `webhook`
: A job that is used to dispatch a webhook payload to a single endpoint

#### Job Capacities
As mentioned above, each jobrunner container acts as one worker that can carry out these actions.
The number of a particular action a worker can carry out is defined by it's capacity which can be 
seen in the `GET /api/v0/workers` endpoint. For example the workers entry may look like this:

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

This means that the worker with the replica ID `000000000000` has a capacity of 1 `scan` and 1
`scanCheck`. A job may have a `capacityMap` field which dictates how much capacity a worker
must have available for the job to be executed.

For example, if we take the above worker's `capacityMap` and the following jobs:

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

Our worker will be able to pick up job id `0` and `2` since it has the capacity for both,
while id `1` will have to wait until the previous scan job is complete:

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

## Job Statuses

Jobs can have the following statuses:

- `waiting`
: the job is unclaimed and waiting to be picked up by a worker
- `running`
: the worker defined by `workerID` is currently running the job
- `done`
: the job has successfully completed
- `error`
: the job has completed with errors
- `cancel_request`
: the worker monitors the job statuses in the database. If the status for a job changes
to `cancel_request`, the worker will cancel the job
- `cancel`
: the job has been cancelled and not fully executed
- `deleted`
: the job and logs have been removed  
- `worker_dead`
: the worker for this job has been declared `dead` and the job will not continue
- `worker_shutdown`
: the worker that was running this job has been gracefully stopped 
- `worker_resurrection`
: the worker for this job has reconnected to the database and will cancel these jobs

## Troubleshooting a Job

An entry for a job can look like:

```json
{
	"id": "1fcf4c0f-ff3b-471a-8839-5dcb631b2f7b",
	"retryFromID": "1fcf4c0f-ff3b-471a-8839-5dcb631b2f7b",
	"workerID": "000000000000",
	"status": "done",
	"scheduledAt": "2017-02-17T01:09:47.771Z",
	"lastUpdated": "2017-02-17T01:10:14.117Z",
	"action": "nautilus_scan_check_single",
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

- `id`: the ID of the job itself
- `workerID`: the ID of the jobrunner worker (synonymous with the DTR replica ID) that is running this job
- `status`: the current state of the job (see [Job Statuses](#job-statuses))
- `action`: what job the worker will actually perform (see [Job Actions](#job-actions))
- `capacityMap`: the available "capacity" a worker needs for this job to run (see [Job Capacities](#job-capacities))

You can view the logs of a particular job by hitting the `GET /api/v0/jobs/{jobID}/logs` endpoint
with the job's `id` as `{jobID}`.

## Cron jobs

Several of the jobs listed in [Job Actions](#job-actions) have been set to run on a
recurring schedule. You can view these jobs with the `GET /api/v0/crons` endpoint which
will return a list similar to this example:

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
      "action": "nautilus_update_db",
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

The `schedule` is simlar to the style of a typical Unix crontab format:
`"second minute hour day month year"`. This determines the next time the `action` will
take place.
