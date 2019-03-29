---
title: Configure Docker Desktop Enterprise on Mac
description: Learn about Docker Desktop Enterprise
keywords: Docker EE, Windows, Mac, Docker Desktop, Enterprise
---

<<<<<<< HEAD
This page contains information on how system administrators can configure Docker Desktop Enterprise (DDE) settings, specify and lock configuration parameters to create a standardized development environment on Mac operating systems.

## Environment configuration (administrators only)

The administrator configuration file allows you to customize and standardize your Docker Desktop environment across the organization.

When you install Docker Desktop Enterprise, a configuration file with default values is installed at the following location. Do not change the location of the `admin-settings.json` file.
=======
# Environment configuration on Mac (administrators only)

The administrator configuration file allows you to customize and standardize your Docker Desktop environment across the organization.

When you install Docker Desktop Enterprise, a configuration file with default values is installed in, and must remain in, the following location:
>>>>>>> 1013: Move desktop ent content to docs-private

`/Library/Application Support/Docker/DockerDesktop/admin-settings.json`

To edit `admin-settings.json`, you must have sudo access privileges.

<<<<<<< HEAD
### Syntax for `admin-settings.json`

1. `configurationFileVersion`: This must be the first parameter listed in `admin-settings.json`. It specifies the version of the configuration file format and must not be changed.
=======
## Syntax for `admin-settings.json`

1. `configurationFileVersion`: This must be the first parameter listed in `admin-settings.json`. It specifies the version of the configuration file format and must be set to 1.
>>>>>>> 1013: Move desktop ent content to docs-private

2. A nested list of configuration parameters, each of which contains a minimum of the following two settings:

  - `locked`: If set to `true`, users without elevated access privileges are not able to edit this setting from the UI or by directly editing the `settings.json` file (the `settings.json` file stores the user's preferences). If set to `false`, users without elevated access privileges can change this setting from the UI or by directly editing
<<<<<<< HEAD
  `settings.json`. If this setting is omitted, the default value is `false`.

  - `value`: Specifies the value of the parameter. Docker Desktop Enterprise uses the value when first started and after a reset to factory defaults. If this setting is omitted, a default value that is built into the application is used.

### Parameters and settings
=======
  `settings.json`. If this setting is omitted, the default value is `false'.

  - `value`: Specifies the value of the parameter. Docker Desktop Enterprise uses the value when first started and after a reset to factory defaults. If this setting is omitted, a default value that is built into the application is used.

## Parameters and settings
>>>>>>> 1013: Move desktop ent content to docs-private

The following `admin-settings.json` code and table provide the required syntax and descriptions for parameters and values:

```json
{
  "configurationFileVersion": 1,
<<<<<<< HEAD
   "analyticsEnabled": {
    "locked": false,
    "value": false
   },
=======

   "analyticsEnabled": {
	 "locked": false,
	 "value": false
   },

>>>>>>> 1013: Move desktop ent content to docs-private
  "dockerCliOptions": {
    "stackOrchestrator": {
      "locked": false,
      "value": "swarm"
    }
  },
<<<<<<< HEAD
=======
  "versionPacks": {
    "allowUserInstall": {
      "value": true
    }
  },
>>>>>>> 1013: Move desktop ent content to docs-private
  "proxy": {
    "locked": false,
    "value": {
      "http": "http://proxy.docker.com:8080",
      "https": "https://proxy.docker.com:8080",
      "exclude": "docker.com,github.com"
    }
  },
"linuxVM": {
<<<<<<< HEAD
  "cpus": {
    "locked": false,
    "value": 2
=======
    "cpus": {
      "locked": false,
      "value": 2
>>>>>>> 1013: Move desktop ent content to docs-private
    },
    "memoryMiB": {
      "locked": false,
      "value": 2048
    },
    "swapMiB": {
      "locked": false,
      "value": 1024
    },
    "diskSizeMiB": {
      "locked": false,
      "value": 65536
    },
    "dataFolder" : {
      "value" : "/Users/...",
      "locked" : false
    },
    "filesharingDirectories": {
      "locked":false,
      "value":["/Users", "..."]
    },
    "dockerDaemonOptions": {
      "experimental": {
         "locked": false,
         "value": true
      }
    }
  },
<<<<<<< HEAD
=======

>>>>>>> 1013: Move desktop ent content to docs-private
  "kubernetes": {

    "enabled": {
      "locked": false,
      "value": false
    },
    "showSystemContainers": {
      "locked": false,
      "value": false
    },
    "podNetworkCIDR": {
      "locked": false,
      "value": null
    },
    "serviceCIDR": {
      "locked": false,
      "value": null
    }
  }
<<<<<<< HEAD
=======

>>>>>>> 1013: Move desktop ent content to docs-private
}
```

Parameter values and descriptions for environment configuration on Mac:

| Parameter                        | Description                      |
| :--------------------------------- | :--------------------------------- |
| `configurationFileVersion`        | Specifies the version of the configuration file format.    |
| `analyticsEnabled`                | If `value` is true, allow Docker Desktop Enterprise to sends diagnostics, crash reports, and usage data. This information helps Docker improve and troubleshoot the application.                |
<<<<<<< HEAD
| `dockerCliOptions`                | Specifies key-value pairs in the user's `~/.docker/config.json` file. In the sample code provided, the orchestration for docker stack commands is set to `swarm` rather than `kubernetes`. |
=======
| `dockerCliOptions`                | Specifies key-value pairs in the user's `%HOME%\\.docker\\config.json` file. In the sample code provided, the orchestration for docker stack commands is set to `swarm` rather than `kubernetes`. |
| `versionPacks`                    | Parameters and settings related to version packs - grouped together here for convenience. |
| `allowUserInstall`                | If true, users are able to install new version packs. If false, only the admin can install new version packs. |
>>>>>>> 1013: Move desktop ent content to docs-private
| `proxy`                          | The `http` setting specifies the HTTP proxy setting. The `https` setting specifies the HTTPS proxy setting. The `exclude` setting specifies a comma-separated list of hosts and domains to bypass the proxy. **Warning:** This parameter should be locked after being set: `locked: "true"`.                 |
| `linuxVM`                         | Parameters and settings related to the Linux VM - grouped together in this example for convenience.          |
| `cpus`                            | Specifies the default number of virtual CPUs for the VM. If the physical machine has only 1 core, the default value is set to 1.    |
| `memoryMiB`                       | Specifies the amount of memory in MiB (1 MiB = 1048576 bytes) allocated for the VM.|
| `swapMiB`                         | Specifies the amount of memory in MiB (1 MiB = 1048576 bytes) allocated for the swap file.                |
| `dataFolder`                      | Specifies the directory containing the VM disk files.  |
| `diskSizeMiB`                     | Specifies the amount of disk storage in MiB (1 MiB = 1048576 bytes) allocated for images and containers.                       |
| `filesharingDirectories`          | The host folders that users can bind-mount in containers.         |
| `dockerDaemonOptions`             | Overrides the options in the linux daemon config file. For more information, see [Docker engine reference](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-configuration-file).        |
| (End of `linuxVM` section.)       |                                   |
| `kubernetes`                      | Parameters and settings related to kubernetes options - grouped together here for convenience.                  |
<<<<<<< HEAD
| `enabled`                         | If `locked` is set to `true`, the Kubernetes cluster starts when Docker Desktop Enterprise is  started.                          |
| `showSystemContainers`            | If true, displays Kubernetes internal containers when running docker commands such as `docker ps`.     |
=======
| `enabled`                         | If `locked` is set to `true`, the k8s cluster starts when Docker Desktop Enterprise is  started.                          |
| `showSystemContainers`            | If true, displays k8s internal containers when running docker commands such as `docker ps`.     |
>>>>>>> 1013: Move desktop ent content to docs-private
| `podNetworkCIDR`                  | This is currently unimplemented. `locked` must be set to true.     |
| `serviceCIDR`                     | This is currently unimplemented. `locked` must be set to true.     |
| (End of `kubernetes` section.)    |                                   |