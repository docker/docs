---
title: Audit Jobs via the API
description: Learn how Docker Trusted Registry runs batch jobs for job-related troubleshooting.
keywords: dtr, troubleshoot, audit, job logs, jobs, api
redirect_from: /ee/dtr/admin/monitor-and-troubleshoot/troubleshoot-batch-jobs/
---


## Overview

This covers troubleshooting batch jobs via the API and was introduced in DTR 2.2. Starting in DTR 2.6, admins have the ability to [audit jobs](audit-jobs-via-ui.md) using the web interface. 

## Prerequisite
   * [Job Queue](job-queue.md)

### Job capacity

Each job runner has a limited capacity and will not claim jobs that require a
higher capacity. You can see the capacity of a job runner via the
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
`scan` and 1 `scanCheck`. Next, review the list of available jobs:

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

If worker `000000000000` notices the jobs
in `waiting` state above, then it will be able to pick up jobs `0` and `2` since it has the capacity
for both. Job `1` will have to wait until the previous scan job, `0`, is completed. The job queue will then look like: 

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
You can get a list of jobs via the `GET /api/v0/jobs/` endpoint. Each job
looks like:

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
The JSON fields of interest here are:

* `id`: The ID of the job
* `workerID`: The ID of the worker in a DTR replica that is running this job
* `status`: The current state of the job
* `action`: The type of job the worker will actually perform
* `capacityMap`: The available capacity a worker needs for this job to run


### Cron jobs

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

The `schedule` field uses a cron expression following the `(seconds) (minutes) (hours) (day of month) (month) (day of week)` format. For example, `57 54 3 * * *` with cron ID `48875b1b-5006-48f5-9f3c-af9fbdd82255` will be run at `03:54:57` on any day of the week or the month, which is `2017-02-22T03:54:57Z` in the example JSON response above.

## Where to go next

- [Enable auto-deletion of job logs](./auto-delete-job-logs.md)
