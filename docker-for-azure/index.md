---
description: Setup & Prerequisites
keywords: azure, microsoft, iaas, tutorial
title: Docker for Azure Setup & Prerequisites
redirect_from:
- /engine/installation/azure/
---

## Prerequisites

- Access to an Azure account with admin privileges
- SSH key that you want to use when accessing your completed Docker install on Azure

## Configuration

Docker for Azure is installed with an Azure template that configures Docker in swarm-mode, running on VMs backed by a custom VHD. There are two ways you can deploy Docker for Azure. You can use the Azure Portal (browser based), or use the Azure CLI. Both have the following configuration options.

### Configuration options

#### Manager Count
The number of Managers in your swarm. You can pick either 1, 3 or 5 managers. We only recommend 1 manager for testing and dev setups. There are no failover guarantees with 1 manager — if the single manager fails the swarm will go down as well. Additionally, upgrading single-manager swarms is not currently guaranteed to succeed.

We recommend at least 3 managers, and if you have a lot of workers, you should pick 5 managers.

#### Manager VM size
The VM type for your manager nodes. The larger your swarm, the larger the VM size you should use.

#### Worker VM size
The VM type for your worker nodes.

#### Worker Count
The number of workers you want in your swarm (1-100).

### Service Principal

To set up Docker for Azure, a [Service Principal](https://azure.microsoft.com/en-us/documentation/articles/active-directory-application-objects/) is required. Docker for Azure uses the principal to operate Azure APIs as you scale up and down or deploy apps on your swarm. Docker provides a containerized helper-script to help create the Service Principal:

    docker run -ti docker4x/create-sp-azure sp-name
    ...
    Your access credentials =============================
    AD App ID:       <app-id>
    AD App Secret:   <secret>
    AD Tenant ID:   <tenant-id>

If you have multiple Azure subscriptions, make sure you're creating the Service Principal with subscription ID that you shared with Docker when signing up for the beta.

`sp-name` is the name of the authentication app that the script creates with Azure. The name is not important, simply choose something you'll recognize in the Azure portal.

If the script fails, it's typically because your Azure user account doesn't have sufficient privileges. Contact your Azure administrator.

When setting up the ARM template, you will be prompted for the App ID (a UUID) and the app secret.

### SSH Key

Docker for Azure uses SSH for accessing the Docker swarm once it's deployed. During setup, you will be prompted for a SSH public key. If you don't have a SSH key, you can generate one with `puttygen` or `ssh-keygen`. You only need the public key component to set up Docker for Azure. Here's how to get the public key from a .pem file:

    ssh-keygen -y -f my-key.pem

### Installing with the CLI
You can also invoke the Docker for Azure template from the Azure CLI:

Here is an example of how to use the CLI. Make sure you populate all of the parameters and their values:
```
$ azure group create  --name DockerGroup --location centralus --deployment-name docker.template --template-file <templateurl>
```
