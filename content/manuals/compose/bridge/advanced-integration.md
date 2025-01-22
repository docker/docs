---
title: Advanced integration
linkTitle: Advanced
weight: 30
description: Learn about how Compose Bridge can function a kubectl plugin
keywords: kubernetes, compose, compose bridge, plugin, advanced
---

{{< summary-bar feature_name="Compose bridge" >}}

Compose Bridge can also function as a `kubectl` plugin, allowing you to integrate its capabilities directly into your Kubernetes command-line operations. This integration simplifies the process of converting and deploying applications from Docker Compose to Kubernetes.

## Use `compose-bridge` as a `kubectl` plugin

To use the `compose-bridge` binary as a `kubectl` plugin, you need to make sure that the binary is available in your PATH and the name of the binary is prefixed with `kubectl-`. 

1. Rename or copy the `compose-bridge` binary to `kubectl-compose_bridge`:

    ```console
    $ mv /path/to/compose-bridge /usr/local/bin/kubectl-compose_bridge
    ```

2. Ensure that the binary is executable:
    
    ```console
    $ chmod +x /usr/local/bin/kubectl-compose_bridge
    ```

3. Verify that the plugin is recognized by `kubectl`:

    ```console
    $ kubectl plugin list
    ```

    In the output, you should see `kubectl-compose_bridge`.

4. Now you can use `compose-bridge` as a `kubectl` plugin:

    ```console
   $ kubectl compose-bridge [command]
    ```

Replace `[command]` with any `compose-bridge` command you want to use.
