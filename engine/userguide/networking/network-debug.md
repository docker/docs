---
description: Learn to use the built-in network debugger to debug overlay networking problems
keywords: network, troubleshooting, debug
title: Debug overlay or swarm networking issues
---

Docker CE 17.12 and higher introduce a network debugging tool designed to help
debug issues with overlay networks and swarm services running on Linux hosts.
When enabled, a network diagnostic server listens on the specified port and
provides diagnostic information. The network debugging tool should only be
started to debug specific issues, and should not be left running all the time.

Information about networks is stored in a database, which can be examined using
the API. The database consists of one table per network.

The Docker API exposes endpoints to query and control the network debugging
tool. CLI integration is provided as a preview, but the implementation is not
yet considered stable and commands and options may change without notice.

## Enable the diagnostic tool

The tool currently only works on Docker hosts running on Linux. Repeat these
steps for each node participating in the swarm.

1.  Set the `network-diagnostic-port` to a port which is free on the Docker
    host, in the `/etc/docker/daemon.json` configuration file.

    ```json
    “network-diagnostic-port”: <port>
    ```

2.  Get the process ID (PID) of the `dockerd` process. It is the second field in
    the output, and is typically a number from 2 to 6 digits long.

    ```bash
    $ ps aux |grep dockerd | grep -v grep
    ```

3.  Reload the Docker configuration without restarting Docker, by sending the
    `HUP` signal to the PID you found in the previous step.

    ```bash
    kill -HUP <pid-of-dockerd>
    ```

A message like the following will appear in the Docker host logs:

```none
Starting the diagnose server listening on <port> for commands
```

## Disable the diagnostic tool

Repeat these steps for each node participating in the swarm.

1.  Remove the `network-diagnostic-port` key from the `/etc/docker/daemon.json`
    configuration file.

2.  Get the process ID (PID) of the `dockerd` process. It is the second field in
    the output, and is typically a number from 2 to 6 digits long.

    ```bash
    $ ps aux |grep dockerd | grep -v grep
    ```

3.  Reload the Docker configuration without restarting Docker, by sending the
    `HUP` signal to the PID you found in the previous step.

    ```bash
    kill -HUP <pid-of-dockerd>
    ```

A message like the following will appear in the Docker host logs:

```none
Disabling the diagnose server
```

## Access the diagnostic tool's API

The network diagnostic tool exposes its own RESTful API. To access the API,
send a HTTP request to the port where the tool is listening. The following
commands assume the tool is listening on port 2000.

Examples are not given for every endpoint.

### Get help

```bash
$ curl localhost:2000/help

OK
/updateentry
/getentry
/gettable
/leavenetwork
/createentry
/help
/clusterpeers
/ready
/joinnetwork
/deleteentry
/networkpeers
/
/join
```

### Join or leave the swarm

```bash
$ curl localhost:2000/join?members=ip1,ip2,...
```

```bash
$ curl localhost:2000/leave?members=ip1,ip2,...
```

### Join or leave a network

```bash
$ curl localhost:2000/joinnetwork?nid=<network id>
```

```bash
$ curl localhost:2000/leavenetwork?nid=<network id>
```

### List cluster peers

```bash
$ curl localhost:2000/clusterpeers
```

### List nodes connected to a given network

```bash
$ curl localhost:2000/networkpeers?nid=<network id>
```

### Dump database tables

```bash
$ curl localhost:2000/gettable?nid=<network id>&tname=<table name>
```

### Interact with a specific database table

```bash
$ curl localhost:2000/<method>?nid=<network id>&tname=<table name>&key=<key>[&value=<value>]
```

## Access the diagnostic tool's CLI

The CLI is provided as a preview and is not yet stable. Commands or options may
change at any time.

The CLI executable is called `debugClient` and is made available using a
standalone container.

The following flags are supported:

| Flag          | Description                                     |
|---------------|-------------------------------------------------|
| -c <string>   | Command to run. One of `sd` or `overlay`.       |
| -ip <string>  | The IP address to query. Defaults to 127.0.0.1. |
| -net <string> | The target network ID.                          |
| -port <int>   | The target port.                                |
| -v            | Enable verbose output.                          |

### Access the CLI

The CLI is provided as a container that needs to run using privileged mode.

1.  To run the container, use a command like the following:

    ```bash
    $ docker container run --name net-diagnostic -d --privileged --network host fcrisciani/network-diagnostic
    ```

2.  Connect to the container using `docker container attach <container-ID>`,
    and start the server using the following command:

    ```bash
    $ kill -HUP 1
    ```

3.  If you have not already done so, join the Docker host to the swarm, then
    run the diagnostic CLI within the container.

    ```bash
    $ ./diagnosticClient <flags>...
    ```

4.  When finished debugging, stop the container.

### Examples

#### Dump the Service discovery table and verify the node ownership

**Standalone network:**

```bash
$ debugClient -c sd -v -net n8a8ie6tb3wr2e260vxj8ncy4
```

**Overlay network:**

```bash
$ debugClient -port 2001 -c overlay -v -net n8a8ie6tb3wr2e260vxj8ncy4
```


