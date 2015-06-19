<!--[metadata]>
+++
title = "Admin guide"
description = "Documentation describing administration of Docker Trusted Registry"
keywords = ["docker, documentation, about, technology, hub, registry, enterprise"]
[menu.main]
parent="smn_dhe"
weight=3
+++
<![end-metadata]-->


# Docker Trusted Registry Administrator's Guide

This guide covers tasks and functions an administrator of Docker Trusted Registry
(DTR) will need to know about, such as reporting, logging, system management,
performance metrics, etc.
For tasks DTR users need to accomplish, such as using DTR to push and pull
images, please visit the [User's Guide](./userguide).

## Reporting

### System Health

![System Health page</admin/metrics/>](../assets/admin-metrics.png)

The "System Health" tab displays "hardware" resource utilization and network traffic metrics for the DTR host as well as for each of its contained services. The CPU and RAM usage meters at the top indicate overall resource usage for the host, while detailed time-series charts are provided below for each service.

In addition, if your registry is using a filesystem storage driver, you will see a usage meter indicating used and available space on the storage volume. Third-party storage back-ends are not supported, so if you are using one, this meter will not be displayed.

You can mouse-over the charts or meters to see detailed data points.

Clicking on a service name (i.e., "load_balancer", "admin_server", etc.) will
display the network, CPU, and memory (RAM) utilization data for the specified
service. See below for a
[detailed explanation of the available services](#services).

### Logs

![System Logs page</admin/logs/>](../assets/admin-logs.png)

Click the "Logs" tab to view all logs related to your DTR instance. You will see
log sections on this page for each service in your DTR instance. Older or newer
logs can be loaded by scrolling up or down. See below for a
[detailed explanation of the available services](#services).

DTR's log files can be found on the host in `/usr/local/etc/dtr/logs/`. The
files are limited to a maximum size of 64mb. They are rotated every two weeks,
when the aggregator sends logs to the collection server, or they are rotated if
a logfile would exceed 64mb without rotation. Log files are named `<component
name>-<timestamp at rotation>`, where the "component name" is the service it
provides (`manager`, `admin-server`, etc.).

### Usage statistics and crash reports

During normal use, DTR generates usage statistics and crash reports. This
information is collected by Docker, Inc. to help us prioritize features, fix
bugs, and improve our products. Specifically, Docker, Inc. collects the
following information:

* Error logs
* Crash logs

## Emergency access to DTR

If your authenticated or public access to the DTR web interface has stopped
working, but your DTR admin container is still running, you can add an
[ambassador container](https://docs.docker.com/articles/ambassador_pattern_linking/)
to get temporary unsecure access to it by running:

    $ docker run --rm -it --link docker_trusted_registry_admin_server:admin -p 9999:80 svendowideit/ambassador

> **Note:** This guide assumes that you are a member of the `docker` group,
> or have root privileges. Otherwise, you may need to add `sudo` to the example
> command above.

This will give you access on port `9999` on your DTR server - `http://<dtr-host-ip>:9999/admin/`.

## Services

DTR runs several Docker services which are essential to its reliability and
usability. The following services are included; you can see their details by
running queries on the [System Health](#system-health) and [Logs](#logs) pages:

* `admin_server`: Used for displaying system health, performing upgrades,
configuring settings, and viewing logs.
* `load_balancer`: Used for maintaining high availability by distributing load
to each image storage service (`image_storage_X`).
* `log_aggregator`: A microservice used for aggregating logs from each of the
other services. Handles log persistence and rotation on disk.
* `image_storage_X`: Stores Docker images using the [Docker Registry HTTP API V2](https://github.com/docker/distribution/blob/master/doc/SPEC.md). Typically,
multiple image storage services are used in order to provide greater uptime and
faster, more efficient resource utilization.
* `postgres`: A database service used to host authentication (LDAP) data and other datasets as needed by DTR.

## DTR system management

The `docker/trusted-registry` image is used to control the DTR system. This
image uses the Docker socket to orchestrate the multiple services that comprise
DTR.

     $ sudo bash -c "$(sudo docker run docker/trusted-registry [COMMAND])"

Supported commands are: `install`, `start`, `stop`, `restart`, `pull`, `info`,
`export-settings`, `diagnostics`, `status`, `upgrade`.

> **Note**: `sudo` is needed for `docker/trusted-registry` commands to
> ensure that the Bash script is run with full access to the Docker host.

### `install`

Install DTR.

### `start`

Start DTR containers that are not running.

### `stop`

Stop DTR containers that are running.

### `restart`

Stop and then start the DTR containers.

### `status`

Display the current running status of only the DTR containers.

```
$ sudo bash -c "$(docker run docker/trusted-registry status)"
INFO  [1.1.0-alpha-001472_g8a9ddb4] Attempting to connect to docker engine dockerHost="unix:///var/run/docker.sock"
INFO  [1.1.0-alpha-001472_g8a9ddb4] Running status command
docker_trusted_registry_load_balancer
    Daemon [default (unix:///var/run/docker.sock)]
        Id: 4d6abd5c39acda25e3d3ccf7cc2acf00f32c7786a7e86fb56daf7fd67584ce9f
        Created: 2015-06-16 21:52:53+00:00
        Status: Up 4 minutes
        Image: docker/trusted-registry-nginx:1.1.0-alpha-001472_g8a9ddb4
        Ports:
            tcp://0.0.0.0:443 -> 443
            tcp://0.0.0.0:80 -> 80
        Command:
            nginxWatcher
        Linked To:
            None

docker_trusted_registry_auth_server
    Daemon [default (unix:///var/run/docker.sock)]
        Id: 22d5c1cf988338638dd810bc8111295f71713e81338d16298028122d33eed64a
        Created: 2015-06-16 21:52:46+00:00
...
```

### `info`

Display the version and info for the Docker daemon, and version and image ID's
of DTR.


```
$ sudo bash -c "$(docker run docker/trusted-registry info)"
INFO  [1.1.0-alpha-001472_g8a9ddb4] Attempting to connect to docker engine dockerHost="unix:///var/run/docker.sock"
{
  "DockerEngine": {
    "Version": {
      "ApiVersion": "1.20",
      "Arch": "amd64",
      "GitCommit": "55bdb51",
      "GoVersion": "go1.4.2",
      "KernelVersion": "3.16.0-4-amd64",
      "Os": "linux",
      "Version": "1.6.0"
    },
    "Info": {
      "ID": "QUMM:6SGD:6ZK4:TLJD:LTX7:64Z5:WP4Y:NE3N:TY7P:Y2RR:KVGO:IWRX",
      "Containers": 15,
      "Driver": "btrfs",
      "DriverStatus": [],
      "ExecutionDriver": "native-0.2",
      "Images": 2793,
      "KernelVersion": "3.16.0-4-amd64",
      "OperatingSystem": "Debian GNU/Linux stretch/sid",
      "NCPU": 4,
      "MemTotal": 12305711104,
      "Name": "t440s",
      "Labels": null,
      "Debug": true,
      "NFd": 43,
      "NGoroutines": 85,
      "SystemTime": "2015-06-17T04:24:54.634746915+10:00",
      "NEventsListener": 1,
      "InitPath": "/usr/bin/docker",
      "InitSha1": "",
      "IndexServerAddress": "https://index.docker.io/v1/",
      "MemoryLimit": false,
      "SwapLimit": false,
      "IPv4Forwarding": true,
      "DockerRootDir": "/data/docker",
      "HttpProxy": "",
      "HttpsProxy": "",
      "NoProxy": ""
    }
  },
  "DTR": {
    "Version": "1.1.0-alpha-001472_g8a9ddb4",
    "GitSHA": "8a9ddb4595c3",
    "StorageDriver": "filesystem",
    "AuthDriver": "dtr",
    "ImageIDs": {
      "Garant": "59bc135c362ad7e44743800b037061976210a9cc6aec323c3ea6eb93ebb513ca",
      "Registry": "6aba58d8bbe71b14edd538a20ac98e1279577bbef461ca25fd2794dcb017c1dc",
      "AdminServer": "af4dfb1f386e3e07b612f5f59f08166ce499ef1dfc619d499a42c53c5e424acf",
      "Manager": "3abc65af8385e63d61af40a1393438d0d720e6bf2a60c1b15b7a17a2a0d8965b",
      "LogAggregator": "01da5d7ef561a251c0c63b860a95d55b602cc70347192ef34acd3b1c5bcd317f",
      "Nginx": "631537f98c8876050fae00106c8db424d03e408b27cc14b5eb1fc11abbaba03b"
    },
    "LicenseKeyID": "2Y6QPUBxoYEms6pIysneyum6SZY_QxE9v4zLF8i1wBNZ"
  }
}
```

### `diagnostics`

The `diagnostics` command is used to extract configuration and run time data
about your containers for support purposes.

The output includes the `docker inspect` output for all
containers, running and not, so please check the resulting files for passwords
and other proprietary information before sending it.

`$ sudo bash -c "$(docker run docker/trusted-registry diagnostics)" > diagnostics.zip`

> **Warning:** These diagnostics files may contain secrets that you need to remove
> before passing on - such as raw container log files, azure storage credentials, or passwords that may be
> sent to non-DTR containers using the `docker run -e PASSWORD=asdf` environment variable
> options.

Stream to STDOUT a zip file containing CSDE and DTR configuration, state, and log
files to help the Docker Enterprise support team:

- your Docker host's `ca-certificates.crt`
- `containers/`: the first 20 running, stopped and paused containers `docker inspect`
  information and log files.
- `dockerEngine/`: the Docker daemon's `info` and `version` output
- `dockerState/`: the Docker daemon's container states, image states, daemon log file, and daemon configuration file
- `dtrlogs/`: the DTR container log files
- `manager/`: the DTR `/usr/local/etc/dtr` DTR configuration directory and DTR manager `info` output. See the [export settings section](#export-settings) for more details.
- `sysinfo/`: Host information
- `errors.txt`: errors and warnings encountered while running diagnostics


### `export-settings`

Export the DTR configuration files for backup or diagnostics use.

`$ sudo bash -c "$(docker run docker/trusted-registry export-settings)" > export-settings.tar.gz`

> **Warning:** These diagnostics files may contain secrets that you need to remove
> before passing on - such as azure storage credentials.

Stream to STDOUT a gzipped tar file containing the DTR configuration files from `/usr/local/etc/dtr/`:

- `garant.yml`
- `generatedConfigs/nginx.conf`
- `generatedConfigs/stacker.yml`
- `hub.yml`
- `license.json`
- `ssl/server.pem`
- `storage.yml`

## Client Docker Daemon diagnostics

To debug client Docker daemon communication issues with DTR, we also provide
a diagnostics tool to be run on the client Docker daemon.

> **Warning:** These diagnostics files may contain secrets that you need to remove
> before passing on - such as raw container log files, azure storage credentials, or passwords that may be
> sent to non-DTR containers using the `docker run -e PASSWORD=asdf` environment variable
> options.

You can download and run this tool using the following command:

> **Note:** If you supply an administrator username and password, then the
> `diagnostics` tool will also download some logs and configuration data
> from the remote DTR server.

```
$ wget https://dhe.mycompany.com/admin/bin/diagnostics && chmod +x diagnostics
$ sudo ./diagnostics dhe.mycompany.com > enduserDiagnostics.zip
DTR administrator username (provide empty string if there is no admin server authentication): 
DTR administrator password (provide empty string if there is no admin server authentication): 
WARN  [1.1.0-alpha-001472_g8a9ddb4] Encountered errors running diagnostics errors=[Failed to copy DTR Adminserver's exported settings into ZIP output: "Failed to read next tar header: \"archive/tar: invalid tar header\"" Failed to copy logs from DTR Adminserver into ZIP output: "Failed to read next tar header: \"archive/tar: invalid tar header\"" error running "sestatus": "exit status 127" error running "dmidecode": "exit status 127"]
```

The zip file will contain the following information:

- your local Docker host's `ca-certificates.crt`
- `containers/`: the first 20 running, stopped and paused containers `docker inspect`
  information and log files.
- `dockerEngine/`: the local Docker daemon's `info` and `version` output
- `dockerState/`: the local Docker daemon's container states, image states, log file, and daemon configuration file
- `dtr/`: Remote DTR services information. This directory will only be populated if the user enters a DTR "admin" username and password.
- - `dtr/logs/`: the remote DTR container log files. This directory will only be populated if the user enters a DTR "admin" username and password.
- - `dtr/exportedSettings/`: the DTR manager container's log files and a backup of the `/usr/local/etc/dtr` DTR configuration directory. See the [export settings section](#export-settings) for more details.
- `sysinfo/`: local Host information
- `errors.txt`: errors and warnings encountered while running diagnostics

## Next Steps

For information on installing DTR, take a look at the [Installation instructions](./install.md).
