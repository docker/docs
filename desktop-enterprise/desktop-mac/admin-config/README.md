---
title: "Environment configuration on Mac"
summary: task
visibleto: employees            # logged in with any Docker ID that has @docker.com email as primary email for Docker ID
author: paige.hargrave
platform:
   - mac
tags:
   - installing                 # pick from kbase list: https://github.com/docker/kbase/blob/master/README.md#prerequisites
---

Environment configuration on Mac (administrators only)
------------------------------------------------------

The administrator configuration file allows you to customize and standardize your Docker Desktop environment across the organization. When you install Docker Desktop Enterprise, a configuration file with default values is installed in, and must remain in, the following location:

When you install Docker Desktop Enterprise, a configuration file with default values is installed in, and must remain in, the following location:

`/Library/Application Support/Docker/DockerDesktop/admin-settings.json`

To edit `admin-settings.json`, you must have sudo access privileges. 

#### Syntax for `admin-settings.json`:

1.  `configurationFileVersion`: This must be the first parameter listed in `admin-settings.json`. It specifies the version of the configuration file format and must be set to 1.

2.  A nested list of configuration parameters, each of which contains a minimum of
    the following two settings:

-   `locked`: If set to `true`, users without elevated access privileges are not able to edit this setting
    from the UI or by directly editing the `admin-settings.json` file. If set to `false`, users without elevated access privileges can change this setting from the UI or by directly editing
    `admin-settings.json`. If this setting is omitted, the default value is `false'???

-   `value`: Specifies the value of the parameter. The default value, contained in the initial `admin-settings.json` file that is installed with Docker Desktop Enterprise, is used when first running Docker Desktop Enterprise and after a reset to factory defaults. If this setting is omitted, the default value is used.

#### Parameters and settings
The following `admin-settings.json` code and table provide the required syntax and descriptions for parameters and values:

```json
{
  "configurationFileVersion": 1,
  
  "dockerCliOptions": {
    "stackOrchestrator": {
      "locked": false,
      "value": "swarm"
    }
  },

  "proxy": {
    "locked": false,
    "value": {
      "http": "http://proxy.docker.com:8080",
      "https": "https://proxy.docker.com:8080",
      "exclude": "docker.com,github.com"
    }
  },
"linuxVM": {
    "cpus": {
      "locked": false,
      "value": 2
    },
    "memoryMiB": {
      "locked": false,
      "value": 2048
    },
    "swapMiB": {
      "locked": false,
      "value": 1024
    },
    "diskPath": {
      "locked": false,
      "value": null
    },
    "diskSizeMiB": {
      "locked": false,
      "value": 65536
    },
    "filesharingDirectories": {
      "locked":false, 
      "value":["/Users", ...]
    }
    "dockerDaemonOptions": {
      "experimental": {
         "locked": false,
         "value": true
      }
    }
  },
  
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
  },
  
}
```

| Parameter                        | Description                      |
| :--------------------------------- | :--------------------------------- |
| `configurationFileVersion`        | Specifies the version of the configuration file format.    |
| `dockerCliOptions`                | Specifies key-value pairs in the user's `%HOME%\\.docker\\config.json` file. In the sample code provided, the orchestration for docker stack commands is set to `swarm` rather than `kubernetes`. |   
| `proxy`                          | The `http` setting specifies the HTTP proxy setting. The `https` setting specifies the HTTPS proxy setting. The `exclude` setting specifies a comma-separated list of hosts and domains to bypass the proxy. **Warning:** This parameter should be locked after being set: `locked: "true"`.                 |             
| `linuxVM`                         | Parameters and settings related to the Linux VM - grouped together in this example for convenience.          |
| `cpus`                            | Specifies the default number of virtual CPUs for the VM. If the physical machine has only 1 core, the default value is set to 1.    |
| `memoryMiB`                       | Specifies the amount of memory in MiB (1 MiB = 1048576 bytes) allocated for the VM.|
| `swapMiB`                         | Specifies the amount of memory in MiB (1 MiB = 1048576 bytes) allocated for the swap file.                |
| `diskPath`                        |  **Warning:** Do not lock this parameter as it can potentially break the version pack switch.    |
| `diskSizeMiB`                     | Specifies the amount of disk storage in MiB (1 MiB = 1048576 bytes) allocated for images and containers.                       |
| `filesharingDirectories`          | The host folders that users can bind-mount in containers.         |
| `dockerDaemonOptions`             | Overrides the options in the linux daemon config file. For more information, see [the Docker engine reference](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-configuration-file)        |
| (End of `linuxVM` section.)       |                                   |
| `kubernetes`                      | Parameters and settings related to kubernetes options - grouped together here for convenience.                  |
| `enabled`                         | If `locked` is set to `true`, the k8s cluster starts when Docker Desktop Enterprise is  started.                          |
| `showSystemContainers`            | If true, displays k8s internal containers when running docker commands such as `docker ps`.     |
| `podNetworkCIDR`                  | This is currently unimplemented. `locked` must be set to true.     |
| `serviceCIDR`                     | This is currently unimplemented. `locked` must be set to true.     |
| (End of `kubernetes` section.)    |                                   |
