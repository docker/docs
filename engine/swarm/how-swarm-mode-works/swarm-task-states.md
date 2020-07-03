---
title: Swarm task states
description: Learn about tasks that are scheduled on your swarm.
keywords: swarm, task, service
redirect_from:
- /datacenter/ucp/2.2/guides/admin/monitor-and-troubleshoot/troubleshoot-task-state/
---

Docker lets you create services, which can start tasks. A service is a
description of a desired state, and a task does the work. Work is scheduled on
swarm nodes in this sequence:

1.  Create a service by using `docker service create`.
2.  The request goes to a Docker manager node.
3.  The Docker manager node schedules the service to run on particular nodes.
4.  Each service can start multiple tasks.
5.  Each task has a life cycle, with states like `NEW`, `PENDING`, and `COMPLETE`.

Tasks are execution units that run once to completion. When a task stops, it
isn't executed again, but a new task may take its place.

Tasks advance through a number of states until they complete or fail. Tasks are
initialized in the `NEW` state. The task progresses forward through a number of
states, and its state doesn't go backward. For example, a task never goes from
`COMPLETE` to `RUNNING`.

Tasks go through the states in the following order:

| Task state  | Description                                                                                                 |
| ----------- | ----------------------------------------------------------------------------------------------------------- |
| `NEW`       | The task was initialized.                                                                                   |
| `PENDING`   | Resources for the task were allocated.                                                                      |
| `ASSIGNED`  | Docker assigned the task to nodes.                                                                          |
| `ACCEPTED`  | The task was accepted by a worker node. If a worker node rejects the task, the state changes to `REJECTED`. |
| `PREPARING` | Docker is preparing the task.                                                                               |
| `STARTING`  | Docker is starting the task.                                                                                |
| `RUNNING`   | The task is executing.                                                                                      |
| `COMPLETE`  | The task exited without an error code.                                                                      |
| `FAILED`    | The task exited with an error code.                                                                         |
| `SHUTDOWN`  | Docker requested the task to shut down.                                                                     |
| `REJECTED`  | The worker node rejected the task.                                                                          |
| `ORPHANED`  | The node was down for too long.                                                                             |
| `REMOVE`    | The task is not terminal but the associated service was removed or scaled down.                             |

## View task state

Run `docker service ps <service-name>` to get the state of a task. The
`CURRENT STATE` field shows the task's state and how long it's been
there.

```bash
$ docker service ps webserver
ID             NAME              IMAGE    NODE        DESIRED STATE  CURRENT STATE            ERROR                              PORTS
owsz0yp6z375   webserver.1       nginx    UbuntuVM    Running        Running 44 seconds ago
j91iahr8s74p    \_ webserver.1   nginx    UbuntuVM    Shutdown       Failed 50 seconds ago    "No such container: webserver.…"
7dyaszg13mw2    \_ webserver.1   nginx    UbuntuVM    Shutdown       Failed 5 hours ago       "No such container: webserver.…"
```

## Where to go next

- [Learn about swarm tasks](https://github.com/docker/swarmkit/blob/master/design/task_model.md)
