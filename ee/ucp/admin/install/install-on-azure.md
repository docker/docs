---
title: Install UCP on Azure
description: Learn how to install Docker Universal Control Plane in a Microsoft Azure environment.
keywords: Universal Control Plane, UCP, install, Docker EE, Azure, Kubernetes
---

Docker Universal Control Plane (UCP) closely integrates with Microsoft Azure for its Kubernetes Networking 
and Persistent Storage feature set. UCP deploys the Calico CNI provider. In Azure,
the Calico CNI leverages the Azure networking infrastructure for data path 
networking and the Azure IPAM for IP address management. There are 
infrastructure prerequisites required prior to UCP installation for the 
Calico / Azure integration.

## Docker UCP Networking

Docker UCP configures the Azure IPAM module for Kubernetes to allocate IP
addresses for Kubernetes pods. The Azure IPAM module requires each Azure virtual
machine which is part of the Kubernetes cluster to be configured with a pool of IP
addresses.

There are two options for provisoning IPs for the Kubernetes cluster on Azure:

- _An automated mechanism provided by UCP which allows for IP pool configuration and maintenance
  for standalone Azure virtual machines._ This service runs within the
  `calico-node` daemonset and provisions 128 IP addresses for each
  node by default. For information on customizing this value, see [Adjusting the IP count value](#adjusting-the-ip-count-value).
- _Manual provision of additional IP address for each Azure virtual machine._ This
  could be done through the Azure Portal, the Azure CLI `$ az network nic ip-config create`,
  or an ARM template. You can find an example of an ARM template
  [here](#manually-provision-ip-address-as-part-of-an-azure-virtual-machine-scale-set).

## Azure Prerequisites 

You must meet the following infrastructure prerequisites in order 
to successfully deploy Docker UCP on Azure:

- All UCP Nodes (Managers and Workers) need to be deployed into the same Azure
  Resource Group. The Azure Networking components (Virtual Network, Subnets,
  Security Groups) could be deployed in a second Azure Resource Group.
- The Azure Virtual Network and Subnet must be appropriately sized for your
  environment, as addresses from this pool will be consumed by Kubernetes Pods.
  For more information, see [Considerations for IPAM
  Configuration](#considerations-for-ipam-configuration).
- All UCP worker and manager nodes need to be attached to the same Azure
  Subnet.
- The Azure Virtual Machine Object Name needs to match the Azure Virtual Machine
  Computer Name and the Node Operating System's Hostname which is the FQDN of
  the host, including domain names. Note that this requires all characters to be in lowercase.
- An Azure Service Principal with `Contributor` access to the Azure Resource
  Group hosting the UCP Nodes. This Service principal will be used by Kubernetes
  to communicate with the Azure API. The Service Principal ID and Secret Key are
  needed as part of the UCP prerequisites. If you are using a separate Resource
  Group for the networking components, the same Service Principal will need
  `Network Contributor` access to this Resource Group.

UCP requires the following information for the installation:

- `subscriptionId` - The Azure Subscription ID in which the UCP 
objects are being deployed. 
- `tenantId` - The Azure Active Directory Tenant ID in which the UCP 
objects are being deployed. 
- `aadClientId` - The Azure Service Principal ID.
- `aadClientSecret` - The Azure Service Principal Secret Key.

### Azure Configuration File

For Docker UCP to integrate with Microsoft Azure,each UCP node in your cluster
needs an Azure configuration file, `azure.json`. Place the file within
`/etc/kubernetes`. Since the config file is owned by `root`, set its permissions 
to `0644` to ensure the container user has read access.

The following is an example template for `azure.json`. Replace `***` with real values, and leave the other
parameters as is.

```
{
    "cloud":"AzurePublicCloud", 
    "tenantId": "***",
    "subscriptionId": "***",
    "aadClientId": "***",
    "aadClientSecret": "***",
    "resourceGroup": "***",
    "location": "****",
    "subnetName": "/****",
    "securityGroupName": "****",
    "vnetName": "****",
    "cloudProviderBackoff": false,
    "cloudProviderBackoffRetries": 0,
    "cloudProviderBackoffExponent": 0,
    "cloudProviderBackoffDuration": 0,
    "cloudProviderBackoffJitter": 0,
    "cloudProviderRatelimit": false,
    "cloudProviderRateLimitQPS": 0,
    "cloudProviderRateLimitBucket": 0,
    "useManagedIdentityExtension": false,
    "useInstanceMetadata": true
}
```

There are some optional parameters for Azure deployments:

- `primaryAvailabilitySetName` - The Worker Nodes availability set.
- `vnetResourceGroup` - The Virtual Network Resource group, if your Azure Network objects live in a 
seperate resource group.
- `routeTableName` - If you have defined multiple Route tables within
an Azure subnet.

See [Kubernetes' azure.go](https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/azure/azure.go) for more details on this configuration file.

## Considerations for IPAM Configuration

The subnet and the virtual network associated with the primary interface of the
Azure virtual machines need to be configured with a large enough address
prefix/range. The number of required IP addresses depends on the workload and
the number of nodes in the cluster.

For example, in a cluster of 256 nodes, make sure that the address space of the subnet and the
virtual network can allocate at least 128 * 256 IP addresses, in order to run a maximum of 128 pods
concurrently on a node. This would be ***in addition to*** initial IP allocations to virtual machine 
NICs (network interfaces) during Azure resource creation.

Accounting for IP addresses that are allocated to NICs during virtual machine bring-up, set
the address space of the subnet and virtual network to `10.0.0.0/16`. This
ensures that the network can dynamically allocate at least 32768 addresses,
plus a buffer for initial allocations for primary IP addresses.

> Azure IPAM, UCP, and Kubernetes
> 
> The Azure IPAM module queries an Azure virtual machine's metadata to obtain
> a list of IP addresses which are assigned to the virtual machine's NICs. The
> IPAM module allocates these IP addresses to Kubernetes pods. You configure the
> IP addresses as `ipConfigurations` in the NICs associated with a virtual machine
> or scale set member, so that Azure IPAM can provide them to Kubernetes when
> requested.
{: .important}

## Manually provision IP address pools as part of an Azure virtual machine scale set

Configure IP Pools for each member of the virtual machine scale set during provisioning by
associating multiple `ipConfigurations` with the scale set’s
`networkInterfaceConfigurations`. Here's an example `networkProfile`
configuration for an ARM template that configures pools of 32 IP addresses
for each virtual machine in the virtual machine scale set.

```json
"networkProfile": {
  "networkInterfaceConfigurations": [
    {
      "name": "[variables('nicName')]",
      "properties": {
        "ipConfigurations": [
          {
            "name": "[variables('ipConfigName1')]",
            "properties": {
              "primary": "true",
              "subnet": {
                "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'), '/subnets/', variables('subnetName'))]"
              },
              "loadBalancerBackendAddressPools": [
                {
                  "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/loadBalancers/', variables('loadBalancerName'), '/backendAddressPools/', variables('bePoolName'))]"
                }
              ],
              "loadBalancerInboundNatPools": [
                {
                  "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/loadBalancers/', variables('loadBalancerName'), '/inboundNatPools/', variables('natPoolName'))]"
                }
              ]
            }
          },
          {
            "name": "[variables('ipConfigName2')]",
            "properties": {
              "subnet": {
                "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'), '/subnets/', variables('subnetName'))]"
              }
            }
          }
          .
          .
          .
          {
            "name": "[variables('ipConfigName32')]",
            "properties": {
              "subnet": {
                "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'), '/subnets/', variables('subnetName'))]"
              }
            }
          }
        ],
        "primary": "true"
      }
    }
  ]
}
```

## UCP Installation 

### Adjust the IP Count Value

If you have manually attached additional IP addresses to the Virtual Machines
(via an ARM Template, Azure CLI or Azure Portal) or you want to reduce the
number of IP Addresses automatically provisioned by UCP from the default of 128
addresses, you can alter the `azure_ip_count` variable in the UCP
Configuration file before installation. If you are happy with 128 addresses per
Virtual Machine, proceed to [installing UCP](#install-ucp).

Once UCP has been installed, the UCP [configuration
file](../configure/ucp-configuration-file/) is managed by UCP and populated with
all of the cluster configuration data, such as AD/LDAP information or networking
configuration. As there is no Universal Control Plane deployed yet, we are able
to stage a [configuration file](../configure/ucp-configuration-file/) just
containing the Azure IP Count value. UCP will populate the rest of the cluster
variables during and after the installation.

Below are some example configuration files with just the `azure_ip_count`
variable defined. These 3-line files can be preloaded into a Docker Swarm prior
to installing UCP in order to override the default `azure_ip_count` value of 128 IP
addresses per node. See [UCP configuration file](../configure/ucp-configuration-file/)
to learn more about the configuration file, and other variables that can be staged pre-install.

> Note: Do not set the `azure_ip_count` to a value of less than 6 if you have not
> manually provisioned additional IP addresses for each Virtual Machine. The UCP
> installation will need at least 6 IP addresses to allocate to the core UCP components 
> that run as Kubernetes pods. That is in addition to the Virtual
> Machine's private IP address.

If you have manually provisioned additional IP addresses for each Virtual
Machine, and want to disallow UCP from dynamically provisioning IP
addresses for you, then your UCP configuration file would be: 

```
$ vi example-config-1
[cluster_config]
  azure_ip_count = "0"
```

If you want to reduce the IP addresses dynamically allocated from 128 to a
custom value, then your UCP configuration file would be: 

```
$ vi example-config-2
[cluster_config]
  azure_ip_count = "20" # This value may be different for your environment
```
See [Considerations for IPAM
Configuration](#considerations-for-ipam-configuration) to calculate an
appropriate value.

To preload this configuration file prior to installing UCP:

1. Copy the configuration file to a Virtual Machine that you wish to become a UCP Manager Node. 
  
2. Initiate a Swarm on that Virtual Machine.

    ```
    $ docker swarm init
    ```

3. Upload the configuration file to the Swarm, by using a [Docker Swarm Config](/engine/swarm/configs/). 
This Swarm Config will need to be named `com.docker.ucp.config`.
    ```
    $ docker config create com.docker.ucp.config <local-configuration-file>
    ```

4. Check that the configuration has been loaded succesfully.
    ```
    $ docker config list
    ID                          NAME                                                      CREATED             UPDATED
    igca3q30jz9u3e6ecq1ckyofz   com.docker.ucp.config                                     1 days ago          1 days ago
    ```

5. You are now ready to [install UCP](#install-ucp). As you have already staged
   a UCP configuration file, you will need to add `--existing-config` to the
   install command below.  

If you need to adjust this value post-installation, see [instructions](../configure/ucp-configuration-file/)
on how to download the UCP configuration file, change the value, and update the configuration via the API. 
If you reduce the value post-installation, existing virtual machines will not be
reconciled, and you will have to manually edit the IP count in Azure.  

### Install UCP

Run the following command to install UCP on a manager node. The `--pod-cidr`
option maps to the IP address range that you have configured for the Azure
subnet, and the `--host-address` maps to the private IP address of the master
node. Finally if you have set the [Ip Count
Value](#adjusting-the-ip-count-value) you will need to add `--existing-config`
to the install command below.

> Note: The `pod-cidr` range must match the Azure Virtual Network's Subnet
> attached the hosts. For example, if the Azure Virtual Network had the range
> `172.0.0.0/16` with Virtual Machines provisioned on an Azure Subnet of
> `172.0.1.0/24`, then the Pod CIDR should also be `172.0.1.0/24`.

```bash
docker container run --rm -it \
  --name ucp \
  --volume /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} install \
  --host-address <ucp-ip> \
  --pod-cidr <ip-address-range> \
  --cloud-provider Azure \
  --interactive
```
