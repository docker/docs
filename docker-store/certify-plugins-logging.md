---
description: Run certification tests against your images
keywords: Docker, docker, store, certified content, logging
title: Certify Docker logging plugins
---

## Introduction

Content that qualifies as **Docker Certified** must conform to best practices and pass certain baseline tests.

Docker Store lets you publish certified images as well as plugins for logging, volumes, and networks. You must certify your own _images and logging plugins_ with the `inspect` tools as explained in these docs. Currently, Docker Store certifies your volume and networking plugins for you upon submission.

This page explains how publishers can successfully test their **Docker logging plugins**. Also available: [Certify your Docker images](certify-images).

> Content that requires a non-certified infrastructure environment cannot be published as certified.

## Certify your logging plugins

You must use the tool, `inspectDockerLoggingPlugin`, to certify your content for publication on Docker Store by ensuring that your Docker logging plugins conform to best practices.

The `inspectDockerLoggingPlugin` command verifies that your Docker logging plugin can be installed and works on Docker Enterprise Edition. It also runs a container from an official Docker image of `alpine:latest` and outputs the contents of a file named `quotes.txt` (included in this repo). In sum, the `inspectDockerLoggingPlugin` command:

- Inspects and displays the Docker logging plugin.

- Installs the Docker logging plugin on Docker EE.

- Runs a Docker service container with the Docker logging plugin, reads a file named `quotes.txt`, echos its contents to `stdout`, and logs the file's content.

- Displays the container logs and compares it to `quotes.txt`. If they match, the test is successful.

The syntax for running a specific logging plugin is: `docker container run --log-driver`.

No parameters are passed to the logging plugin. If parameters are required for the Docker logging plugin to work correctly, then a custom test script must be written and used. The default `docker container run` command is:

  ```
  docker container run -it --log-driver xxxxxxxxxxxxxxxxxxxxx \
  --volume \"$(pwd)/quotes.txt:/quotes.txt\" alpine:latest \
  sh -c 'cat /quotes.txt;sleep 20
  ```

The custom script must log the contents of the `quotes.txt` file. It should also cleanup (remove the container and docker image). Refer to the `--test-script` command argument in the command help.

### Docker container logs

Best practices require Docker logging plugins to support the [ReadLogs API](/engine/extend/plugins_logging/#logdriverreadlogs) so that the logs can be retrieved with the `docker container logs` command. If the `ReadLogs` API is not supported, a custom script is needed to retrieve the logs and print them to `stdout`. Refer to the `--get-logs-script` command argument in the command help.

### Prerequisites

Your Docker EE installation must be running on the server used to verify your submissions. If necessary, request entitlement to a specific [Docker Enterprise Edition](https://store.docker.com/editions/enterprise/docker-ee-trial).

- Docker EE (on the server for verifying submissions)
- git client
- inspectDockerLoggingPlugin tool

### Set up testing environment

There are three steps: (1) install git, (2) configure credentials, and (3) configure endpoints.

1.  Install git (required for `inspectDockerLoggingPlugin`):

    **Ubuntu**

    ```bash
    sudo apt-get update -qq
    sudo apt-get install git -y
    ```

    **RHEL/CentOS**

    ```bash
    sudo yum makecache fast
    sudo yum install git -y
    ```

    **Windows**

    To download and install git for Windows: <https://git-scm.com/download/win>.

2.  Configure your Docker Registry credentials by either _defining environment variables_ **or** _passing them as arguments_ to `inspectDockerLoggingPlugin`.

    a.  Define environment variables for registry credentials, `DOCKER_USER` and `DOCKER_PASSWORD`:

    **Linux**

    ```bash
    export DOCKER_USER="my_docker_registry_user_account"
    export DOCKER_PASSWORD="my_docker_registry_user_account_password"
    ```

    **Windows command prompt**

    ```bash
    set DOCKER_USER="my_docker_registry_user_account"
    set DOCKER_PASSWORD="my_docker_registry_user_account_password"
    ```

    **Windows powershell**

    ```bash
    $env:DOCKER_USER="my_docker_registry_user_account"
    $env:DOCKER_PASSWORD="my_docker_registry_user_account_password"
    ```

    b.  Pass arguments to `inspectDockerLoggingPlugin` (or be prompted for them):

    ```
    --docker-user
    --docker-password
    ```

3.  Configure endpoints (and override default values) by either _defining environment variables_ **or** _passing them as arguments_ to `inspectDockerLoggingPlugin`.

    By default, `inspectDockerLoggingPlugin` uses these two endpoints to communicate with the Docker Hub Registry:

    - Registry Authentication Endpoint: **https://auth.docker.io**
    - Registry API Endpoint: **https://registry-1.docker.io**

    You may want to use your private registry for initial testing and override the defaults.

    a.  Define environment variables, `DOCKER_REGISTRY_AUTH_ENDPOINT` and  `DOCKER_REGISTRY_API_ENDPOINT`:

    **Linux or MacOS**

    ```bash
    export DOCKER_REGISTRY_AUTH_ENDPOINT="https://my_docker_registry_authentication_endpoint"
    export DOCKER_REGISTRY_API_ENDPOINT="https://my_docker_registry_api_enpoint"
    ```

    **Windows command prompt**

    ```bash
    set DOCKER_REGISTRY_AUTH_ENDPOINT="https://my_docker_registry_authentication_endpoint"
    set DOCKER_REGISTRY_API_ENDPOINT="https://my_docker_registry_api_enpoint"
    ```

    **Windows powershell**

    ```bash
    $env:DOCKER_REGISTRY_AUTH_ENDPOINT="https://my_docker_registry_authentication_endpoint"
    $env:DOCKER_REGISTRY_API_ENDPOINT="https://my_docker_registry_api_enpoint"
    ```

    b.  Pass your endpoints as arguments to `inspectDockerLoggingPlugin`:

    ```
    --docker-registry-auth-endpoint
    --docker-registry-api-endpoint
    ```

### Syntax

1.  Download `inspectDockerLoggingPlugin` command.

    | OS/Architecture | Download Link |
    |-----------------|------------------|
    | Windows/X86  | [https://s3.amazonaws.com/store-logos-us-east-1/certification/windows/inspectDockerLoggingPlugin.exe](https://s3.amazonaws.com/store-logos-us-east-1/certification/windows/inspectDockerLoggingPlugin.exe) |
    | Linux/X86 | [https://s3.amazonaws.com/store-logos-us-east-1/certification/linux/inspectDockerLoggingPlugin](https://s3.amazonaws.com/store-logos-us-east-1/certification/linux/inspectDockerLoggingPlugin) |
    | Linux/IBMZ | [https://s3.amazonaws.com/store-logos-us-east-1/certification/zlinux/inspectDockerLoggingPlugin](https://s3.amazonaws.com/store-logos-us-east-1/certification/zlinux/inspectDockerLoggingPlugin) |
    | Linux/IBMPOWER | [https://s3.amazonaws.com/store-logos-us-east-1/certification/power/inspectDockerLoggingPlugin](https://s3.amazonaws.com/store-logos-us-east-1/certification/power/inspectDockerLoggingPlugin) |

2.  Set permissions on `inspectDockerLoggingPlugin` so that it is executable:

    ```
    chmod u+x inspectDockerLoggingPlugin
    ```

3.  Get the product ID from the plan page you'd like to reference for the certification test. Make sure the checkbox is checked and the plan is saved first.

    ![product ID](images/store-product-id.png)

    ```none
    Inspects a Docker logging plugin to see if it conforms to best practices.

    Syntax: inspectDockerLoggingPlugin [options] dockerLoggingPlugin

    Options:
      -docker-password string
        	 Docker Password.  This overrides the DOCKER_PASSWORD environment variable.
      -docker-registry-api-endpoint string
        	 Docker Registry API Endpoint. This overrides the DOCKER_REGISTRY_API_ENDPOINT environment variable. (default "https://registry-1.docker.io")
      -docker-registry-auth-endpoint string
        	 Docker Registry Authentication Endpoint. This overrides the DOCKER_REGISTRY_AUTH_ENDPOINT environment variable. (default "https://auth.docker.io")
      -docker-user string
        	 Docker User ID.  This overrides the DOCKER_USER environment variable.
      -get-logs-script string
        	 An optional custom script used to retrieve the logs.
      -help
        	 Help on the command.
      -html
        	 Generate HTML output.
      -json
        	 Generate JSON output.
      -product-id string
        	 Optional Product identifier from Docker Store for this plugin. Please include it when you want the output sent to docker store for certification.
      -test-script string
        	 An optional custom script used to test the Docker logging plugin. The script gets passed 1 parameter - the Docker logging plugin name.
      -verbose
        	 Displays more verbose output.

    dockerLoggingPlugin
  	The Docker logging plugin to inspect. This argument is required.
    ```

## Inspection Output

By default, `inspectDockerLoggingPlugin` displays output locally to `stdout` (the default), JSON, and HTML. You can also upload output to Docker Store, which is recommended for admnistrator verification.

-  **Upload to Docker Store** (by entering `product-id` at the commandline).

-  **Send message to `stdout`**. This is the default.

-  **JSON sent to `stdout`**. Use the `--json` option to override and replace the messages sent to `stdout`.

-  **HTML local file**. Use the `--html` option to generate an HTML report. Both `--json` and `--html` can be specified at the same time.

## Inspection Examples

* [Inspect a Docker logging plugin with messages sent to stdout](#inspect-logging-plugin-stdout)
* [Inspect a Docker logging plugin with JSON output](#inspect-logging-plugin-json)
* [Inspect a Docker logging plugin with HTML output](#inspect-logging-plugin-html)
* [Send data to API endpoint on external server](#send-data-to-api-endpoint-on-external-server)

<a name="inspect-logging-plugin-stdout">

### Inspect a Docker logging plugin with messages sent to stdout

#### To inspect the Docker logging plugin "gforghetti/docker-log-driver-test:latest", and upload the result to Docker Store (leave out the `-product-id` parameter if you are just testing):

```
gforghetti:~:$ ./inspectDockerLoggingPlugin -product-id=<store-product-id> gforghetti/docker-log-driver-test:latest
```
#### Output:

```
**************************************************************************************************************************************************************************************************
* Docker logging plugin: gforghetti/docker-log-driver-test:latest
**************************************************************************************************************************************************************************************************

**************************************************************************************************************************************************************************************************
* Step #1 Inspecting the Docker logging plugin: gforghetti/docker-log-driver-test:latest ...
**************************************************************************************************************************************************************************************************
Passed:   Docker logging plugin image gforghetti/docker-log-driver-test:latest has been inspected.

**************************************************************************************************************************************************************************************************
* Step #2 Docker logging plugin information
**************************************************************************************************************************************************************************************************
+-------------------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Docker logging plugin:  | gforghetti/docker-log-driver-test:latest                                                                                                                             |
| Description:            | jsonfilelog as plugin                                                                                                                                                |
| Documentation:          | -                                                                                                                                                                    |
| Digest:                 | sha256:870c81b9b8872956eec83311707e52ccd4af40683ba64aac3c1700d252a07624                                                                                              |
| Base layer digest:      | sha256:8b0c5cbf1339dacef5f56717567aeee37d9e4f196f0874457d46c01592a30d70                                                                                              |
| Docker version:         | 17.06.0-ce                                                                                                                                                           |
| Interface Socket:       | jsonfile.sock                                                                                                                                                        |
| Interface Socket Types: | docker.logdriver/1.0                                                                                                                                                 |
| IpcHost:                | false                                                                                                                                                                |
| PidHost:                | false                                                                                                                                                                |
| Entrypoint:             | /usr/bin/docker-log-driver                                                                                                                                           |
| WorkDir:                | /tmp                                                                                                                                                                 |
| User:                   |                                                                                                                                                                      |
+-------------------------+----------------------------------------------------------------------------------------------------------------------------------------------------------------------+

**************************************************************************************************************************************************************************************************
* Step #3 Checking to see if the Docker logging plugin is already installed.
**************************************************************************************************************************************************************************************************
The Docker logging plugin gforghetti/docker-log-driver-test:latest needs to be installed.

**************************************************************************************************************************************************************************************************
* Step #4 Installing the Docker logging plugin gforghetti/docker-log-driver-test:latest ...
**************************************************************************************************************************************************************************************************
Passed:   Docker logging plugin gforghetti/docker-log-driver-test:latest has been installed successfully.

**************************************************************************************************************************************************************************************************
* Step #5 Testing the Docker logging plugin: gforghetti/docker-log-driver-test:latest ...
**************************************************************************************************************************************************************************************************
Starting a Docker container to test the docker logging plugin gforghetti/docker-log-driver-test:latest

**************************************************************************************************************************************************************************************************
* Step #6 Retrieving the Docker Logs ...
**************************************************************************************************************************************************************************************************
Retrieving the Docker logs using the "docker container logs c0bc7de01b155203d1ab72a1496979c4715c9b9f81c7f5a1a37ada3aa333e1d5" command

**************************************************************************************************************************************************************************************************
* Step #7 Verifying that the contents retrieved matches what was sent to the Docker Logging plugin.
**************************************************************************************************************************************************************************************************
Passed:   Docker logging plugin Test was successful.

**************************************************************************************************************************************************************************************************
* Step #8 Removing the Docker container and any associated volumes.
**************************************************************************************************************************************************************************************************
Passed:   Docker container and any associated volumes removed.

**************************************************************************************************************************************************************************************************
* Step #9 Removing the Docker logging plugin
**************************************************************************************************************************************************************************************************
Passed:   Docker logging plugin gforghetti/docker-log-driver-test:latest was removed.

**************************************************************************************************************************************************************************************************
* Summary of the inspection for the Docker logging plugin: gforghetti/docker-log-driver-test:latest
**************************************************************************************************************************************************************************************************

Report Date: Tue Sep 19 12:40:11 2017
Operating System: Operating System: MacOS darwin Version: 10.12.6
Docker version 17.09.0-ce-rc2, build 363a3e7


Passed:   Docker logging plugin image gforghetti/docker-log-driver-test:latest has been inspected.
Passed:   Docker logging plugin gforghetti/docker-log-driver-test:latest has been installed successfully.
Passed:   Docker logging plugin Test was successful.
Passed:   Docker container and any associated volumes removed.
Passed:   Docker logging plugin gforghetti/docker-log-driver-test:latest was removed.

The inspection of the Docker logging plugin gforghetti/docker-log-driver-test:latest has completed.

If -product-id is specified on command line:
**************************************************************************************************************************************************************************************************
* Step #10 Upload the test result to Docker Store.
**************************************************************************************************************************************************************************************************
Passed:   The test results are uploaded to Docker Store.

gforghetti:~/$
```

<a name="inspect-logging-plugin-json">

### Inspect a Docker logging plugin with JSON Output

#### To inspect the  Docker logging plugin `gforghetti/docker-log-driver-test:latest` with JSON Output:

```
gforghetti:~:$ ./inspectDockerLoggingPlugin --json gforghetti/docker-log-driver-test:latest | jq
```

Note: The output was piped to the **jq** command to display it "nicely".

#### Output:

```
{
  "Date": "Tue Sep 19 12:41:49 2017",
  "SystemOperatingSystem": "Operating System: MacOS darwin Version: 10.12.6",
  "SystemDockerVersion": "Docker version 17.09.0-ce-rc2, build 363a3e7",
  "DockerLogginPlugin": "gforghetti/docker-log-driver-test:latest",
  "Description": "jsonfilelog as plugin",
  "Documentation": "-",
  "DockerLoggingPluginDigest": "sha256:870c81b9b8872956eec83311707e52ccd4af40683ba64aac3c1700d252a07624",
  "BaseLayerImageDigest": "sha256:8b0c5cbf1339dacef5f56717567aeee37d9e4f196f0874457d46c01592a30d70",
  "DockerVersion": "17.06.0-ce",
  "Entrypoint": "/usr/bin/docker-log-driver",
  "InterfaceSocket": "jsonfile.sock",
  "InterfaceSocketTypes": "docker.logdriver/1.0",
  "WorkDir": "/tmp",
  "User": "",
  "IpcHost": false,
  "PidHost": false,
  "Errors": 0,
  "Warnings": 0,
  "HTMLReportFile": "",
  "VulnerabilitiesScanURL": "",
  "Results": [
    {
      "Status": "Passed",
      "Message": "Docker logging plugin image gforghetti/docker-log-driver-test:latest has been inspected."
    },
    {
      "Status": "Passed",
      "Message": "Docker logging plugin gforghetti/docker-log-driver-test:latest has been installed successfully."
    },
    {
      "Status": "Passed",
      "Message": "Docker logging plugin Test was successful."
    },
    {
      "Status": "Passed",
      "Message": "Docker container and any associated volumes removed."
    },
    {
      "Status": "Passed",
      "Message": "Docker logging plugin gforghetti/docker-log-driver-test:latest was removed."
    }
  ]
}
gforghetti:~/$
```

<a name="inspect-logging-plugin-html">

### Inspect a Docker logging plugin with HTML output

#### To inspect the  Docker logging plugin `gforghetti/docker-log-driver-test:latest` with HTML output:

```
gforghetti:~:$ ./inspectDockerLoggingPlugin --html gforghetti/docker-log-driver-test:latest
```

#### Output:

Note: The majority of the stdout message output has been intentionally omitted below.

```
The inspection of the Docker logging plugin cpuguy83/docker-logdriver-test:latest has completed.
An HTML report has been generated in the file cpuguy83-docker-logdriver-test-latest_inspection_report.html
gforghetti:~/$
```

![HTML Output Image](images/gforghetti-log-driver-latest_inspection_report.html.png)

<a name="send-to-http-api-endpoint">

### Send data to API endpoint on external server

#### Introduction

The **http_api_endpoint** is an HTTP Server that can be used to test docker logging plugins that do not support the read log api and instead send data to an API Endpoint running on an external server.
The [Sumo Logic Logging Plugin](https://store.docker.com/plugins/sumologic-logging-plugin) is one example.

You can configure those docker logging plugins to send their logging data to the **http_api_endpoint** HTTP Server for testing the plugin and then code a script to retrieve the logs using a curl command.

#### Download

* [Linux/X86](https://s3.amazonaws.com/store-logos-us-east-1/certification/linux/http_api_endpoint)
* [Windows/X86](https://s3.amazonaws.com/store-logos-us-east-1/certification/windows/http_api_endpoint.exe)

#### Syntax

```
./http_api_endpoint [options]
```

Options:

 * **--port**    (The port for the **http_api_endpoint** HTTP Server to listen on. Defaults to port 80)
 * **--debug**   (write debugging information)
 * **--help**    (display the command help)

#### Using and testing the **http_api_endpoint** HTTP Server

The **curl** command can be used to test and use the **http_api_endpoint** HTTP Server.

* Issue the following curl command to send new data to the **http_api_endpoint**:

  ```
  # DATA='Hello World!'
  # curl -s -X POST -d "${DATA}" http://127.0.0.1:80
  ```

  Note: if any data was previously sent, it will be replaced.

* Issue the following curl command to send data to the **http_api_endpoint** and append that data to the already collected data:

  ```
  # DATA='Hello World!'
  # curl -s -X POST -d "${DATA}" http://127.0.0.1:80
  ```

* Issue the following curl command to retrieve the data from the http_api_endpoint:

  ```
  # curl -s -X GET http://127.0.0.1:80
  ```
  ```
  Hello World!
  ```

* Issue the following curl command to erase any data currently collected by the http_api_endpoint:

  ```
  # curl -s -X DELETE http://127.0.0.1:80
  ```

* To Terminate:

  ```
  # curl -s http://127.0.0.1:80/EXIT
  ```

#### Example of using the http_api_endpoint HTTP Server for a Logging Plugin

##### Script to run a container to test the Logging Plugin

```
# cat test_new_plugin.sh
```
```
#!/usr/bin/env bash

#######################################################################################################################################
# This bash script tests a Docker logging plugin that does not support the read log api and instead sends data to an API Endpoint running on an external server.
#
#######################################################################################################################################
# Docker Inc.
#######################################################################################################################################

#######################################################################################################################################
# Make sure the Docker logging plugin was specified on the command line
#######################################################################################################################################
DOCKER_LOGGING_PLUGIN=$1
if [[ -z $DOCKER_LOGGING_PLUGIN ]]; then
    printf 'You must specify the Docker Loggin Plugin!\n'
    exit 1
fi

HTTP_API_ENDPOINT='http://localhost:80'

#######################################################################################################################################
# Check to make sure the http_api_endpoint HTTP Server is running
#######################################################################################################################################
curl -s -X POST "${HTTP_API_ENDPOINT}"
if [[ $? -ne 0 ]]; then
    printf 'Unable to connect to the HTTP API Endpoint: '"${HTTP_API_ENDPOINT}"'!\n'
    exit 1
fi

#######################################################################################################################################
# Run a alpine container with the plugin and send data to it
#######################################################################################################################################
docker container run \
--rm \
--log-driver="${DOCKER_LOGGING_PLUGIN}" \
--log-opt sumo-url="${HTTP_API_ENDPOINT}" \
--log-opt sum-sending-interval=5s \
--log-opt sumo-compress=false \
--volume $(pwd)/quotes.txt:/quotes.txt \
alpine:latest \
sh -c 'cat /quotes.txt;sleep 10'

exit $?
```

##### Script to retrieve the logging data from the http_api_endpoint HTTP Server

```
# cat get_plugin_logs.sh
```
```
#!/usr/bin/env sh

#######################################################################################################################################
# This bash script retrieves any data logged to the http_api_endpoint HTTP Server.
#######################################################################################################################################
# Docker Inc.
#######################################################################################################################################

curl -s -X GET http://127.0.0.1:80

```

##### To test the Docker logging plugin

```
./inspectDockerLoggingPlugin --verbose --html --test-script ./test_plugin.sh --get-logs-script ./get_plugin_logs.sh myNamespace/docker-logging-driver:1.0.2
```
