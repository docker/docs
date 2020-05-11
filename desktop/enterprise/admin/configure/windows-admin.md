---
title: Configure Docker Desktop Enterprise on Windows
description: Learn about Docker Desktop Enterprise configuration
keywords: Docker Admin, Windows, Docker Desktop, Enterprise
redirect_from:
- /ee/desktop/admin/configure/windows-admin/
---

This page contains information on how system administrators can configure Docker Desktop Enterprise (DDE) settings, specify and lock configuration parameters to create a standardized development environment on Windows operating systems.

## Environment configuration (administrators only)

The administrator configuration file allows you to customize and standardize your Docker Desktop environment across the organization.

When you install Docker Desktop Enterprise, a configuration file with default values is installed at the following location. Do not change the location of the `admin-settings.json` file.

`%ProgramData%\DockerDesktop\admin-settings.json`

which defaults to:

`C:\ProgramData\DockerDesktop\admin-settings.json`

You must have administrator access privileges to edit `admin-settings.json`.

### Syntax for `admin-settings.json`

1. `configurationFileVersion`: This must be the first parameter listed in `admin-settings.json`. It specifies the version of the configuration file format and must not be changed.

2. A nested list of configuration parameters, each of which contains a minimum of
    the following two settings:

- `locked`: If set to `true`, users without elevated access privileges are not able to edit this setting
    from the UI or by directly editing the `settings.json` file (the `settings.json` file stores the user's preferences). If set to `false`, users without elevated access privileges can change this setting from the UI or by directly editing
    `settings.json`. If this setting is omitted, the default value is `false'.

- `value`: Specifies the value of the parameter. Docker Desktop Enterprise uses the value when first started and after a reset to factory defaults. If this setting is omitted, a default value that is built into the application is used.

### Parameters and settings

The following `admin-settings.json` code and table provide the required syntax and descriptions for parameters and values:

```json
{
  "configurationFileVersion": 2,
  "engine": {
    "locked": false,
    "value": "linux"
  },
  "analyticsEnabled": {
    "locked": false,
    "value": false
  },
  "exposeDockerAPIOnTCP2375": {
    "locked": false,
    "value": false
  },
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
    "dataFolder": {
      "locked": false,
      "value": null
    },
    "diskSizeMiB": {
      "locked": false,
      "value": 65536
    },
    "hypervCIDR": {
      "locked": false,
      "value": "10.0.75.0/28"
    },
    "vpnkitCIDR": {
      "locked": false,
      "value": "192.168.65.0/28"
    },
    "useDnsForwarder": {
      "locked": false,
      "value": true
    },
    "dns": {
      "locked": false,
      "value": "8.8.8.8"
    },
    "dockerDaemonOptions": {
      "experimental": {
        "locked": false,
        "value": true
      }
    }
  },
  "windows": {
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
  "template": {
    "defaultOrg": {
      "value": "myorg",
      "locked": true
    },
    "defaultRegistry": {
      "value": "mydtr:5000",
      "locked": true
    },
    "repositories": {
      "value": [
        "https://one/library.yaml",
        "https://two/library.yaml",
        "https://three/library.yaml"
      ],
      "locked": true
    }
  },
  "filesharingDirectories": {
    "locked": false,
    "value": [
      "%USERPROFILE%"
    ]
  }
}
```

Parameter values and descriptions for environment configuration on Windows:

| Parameter                        | Description                      |
| :--------------------------------- | :--------------------------------- |
| `configurationFileVersion`        | Specifies the version of the configuration file format.    |
| `engine`                          | Specifies the default Docker engine to be used. `linux` specifies the Linux engine. `windows` specifies the Windows engine.               |
| `analyticsEnabled`                | If `value` is true, allow Docker Desktop Enterprise to sends diagnostics, crash reports, and usage data. This information helps Docker improve and troubleshoot the application.                |
| `exposeDockerAPIOnTCP2375`        | Exposes Docker API on a specified port. In this example, setting 'locked' to `true` exposes the Docker API on port 2375. **Warning:** This is unauthenticated and should only be enabled if protected by suitable firewall rules.|
| `dockerCliOptions`                | Specifies key-value pairs in the user's `%HOME%\.docker\config.json` file. In the sample code provided, the orchestration for docker stack commands is set to `swarm` rather than `kubernetes`. |
| `proxy`                           | The `http` setting specifies the HTTP proxy setting. The `https` setting specifies the HTTPS proxy setting. The `exclude` setting specifies a comma-separated list of hosts and domains to bypass the proxy. **Warning:** This parameter should be locked after being set: `locked: "true"`.                 |
| `linuxVM`                         | Parameters and settings related to the Linux VM - grouped together in this example for convenience.          |
| `cpus`                            | Specifies the default number of virtual CPUs for the VM. If the physical machine has only 1 core, the default value is set to 1.    |
| `memoryMiB`                       | Specifies the amount of memory in MiB (1 MiB = 1048576 bytes) allocated for the VM.
| `swapMiB`                         | Specifies the amount of memory in MiB (1 MiB = 1048576 bytes) allocated for the swap file.                |
| `dataFolder`                      | Specifies the root folder where Docker Desktop should put VM disk files.    |
| `diskSizeMiB`                     | Specifies the amount of disk storage in MiB (1 MiB = 1048576 bytes) allocated for images and containers.                       |
| `hypervCIDR`                      | Specifies the subnet used for Hyper-V networking. The chosen subnet must not conflict with other resources on your network.                          |
| `vpnkitCIDR`                      | Specifies the subnet used for VPNKit networking and drive sharing. The chosen subnet must not conflict with other resources on your network.                          |
| `useDnsForwarder`                 | If `value` is set to `true`, this automatically determines the upstream DNS servers based on the host's network adapters.      |
| `dns`                             | If `value` for `useDnsForwarder` is set to `false`, the Linux VM uses the server information in this `value` setting for DNS resolution.                       |
| `dockerDaemonOptions`             | Overrides the options in the Linux daemon config file. For more information, see [Docker engine reference](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-configuration-file).      |
| (End of `linuxVM` section.)       |                                   |
| `windows`                         | Parameters and settings related to the Windows daemon-related options - grouped together in this example for convenience.          |
| `dockerDaemonOptions`             | Overrides the options in the Windows daemon config file. For more information, see [Docker engine reference](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-configuration-file).     |
| (End of `windows` section.)       |                                   |
| `kubernetes`                      | Parameters and settings related to kubernetes options - grouped together here for convenience.                  |
| `enabled`                         | If `locked` is set to `true`, the Kubernetes cluster starts when Docker Desktop Enterprise is started.                          |
| `showSystemContainers`            | If true, displays Kubernetes internal containers when running docker commands such as `docker ps`.     |
| `podNetworkCIDR`                  | This is currently unimplemented. `locked` must be set to true.     |
| `serviceCIDR`                     | This is currently unimplemented. `locked` must be set to true.     |
| (End of `kubernetes` section.)    |                                   |
|`template`|Parameters and settings related to Docker Template and Application Designer - grouped together in this example for convenience. For more information, see [`Docker template config`](/engine/reference/commandline/template_config/).|
|`defaultOrg`| Specifies the default organization to be used in Docker Template and Docker Application Designer. If `locked` is set to `true`, the Kubernetes cluster starts when Docker Desktop Enterprise is started. |
|`defaultRegistry`|Specifies the default registry to be used in Docker Template and Application Designer.|
|`repositories`|Lists the repositories that are allowed.|
| `filesharingDirectories`          | The host folders that users can bind-mount in containers.         |

### File format update

#### From version 1 to 2

Docker Desktop Enterprise 2.3.0.0-ent contains a change in the configuration file format.

If you haven’t made any changes to the `admin-settings.json` file in the previous installation, you can simply delete it and Docker Desktop will re-create it automatically.
Otherwise manual steps are required to update the `admin-settings.json` file.

1. Increase the value of the `configurationFileVersion` field from 1 to 2, i.e. before:
    ```json
   {
      "configurationFileVersion": 1,
      ...
   }
    ```
    after:
    ```json
   {
      "configurationFileVersion": 2,
      ...
   }
    ```

2. Replace the `sharedDrives` and `sharedFolders` fields with a single `filesharingDirectories` field, e.g. before:
    ```json
   {
      ...
      "sharedDrives": {
        "locked": true,
        "value": ["C"]
      },
      "sharedFolders": ["%USERPROFILE%"]
   }
    ```
    after:
    ```json
   {
      ...
      "filesharingDirectories": {
        "locked": true,
        "value": ["C:", "%USERPROFILE%"]
      }
   }
    ```

