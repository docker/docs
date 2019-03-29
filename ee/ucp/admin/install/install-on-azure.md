---
title: Install UCP on Azure
description: Learn how to install Docker Universal Control Plane in a Microsoft Azure environment.
keywords: Universal Control Plane, UCP, install, Docker EE, Azure, Kubernetes
---

Docker UCP closely integrates into Microsoft Azure for its Kubernetes Networking 
and Persistent Storage feature set. UCP deploys the Calico CNI provider. In Azure
the Calico CNI leverages the Azure networking infrastructure for data path 
networking and the Azure IPAM for IP address management. There are 
infrastructure prerequisites that are required prior to UCP installation for the 
Calico / Azure integration.

## Docker UCP Networking

Docker UCP configures the Azure IPAM module for Kubernetes to allocate IP
addresses to Kubernetes pods.  The Azure IPAM module requires each Azure virtual
machine that's part of the Kubernetes cluster to be configured with a pool of IP
addresses.

There are two options for provisoning IPs for the Kubernetes cluster on Azure
- Docker UCP provides an automated mechanism to configure and maintain IP pools
  for standalone Azure virtual machines. This service runs within the
  calico-node daemonset and by default will provision 128 IP address for each
  node. For information on customising this value see [adjusting the ip count value](#adjusting-the-ip-count-value) for more information.
- Manually provision additional IP address for each Azure virtual machine. This
  could be done in the Azure Portal, via the Azure CLI `$ az network nic
  ip-config create` or through an ARM template. You can find an example of an
  ARM template
  [here](#manually-provision-ip-address-as-part-of-an-azure-virtual-machine-scale-set).

## Azure Prerequisites 

You must meet the following infrastructure prerequisites in order 
to successfully deploy Docker UCP on Azure:

- All UCP Nodes (Managers and Workers) need to be deployed into the same Azure
  Resource Group. The Azure Networking components (Virtual Network, Subnets,
  Security Groups)  could be deployed in a second Azure Resource Group.
- The Azure Virtual Network and Subnet must be appropriately sized for your environment, as
  addresses from this pool will be consumed by Kubernetes Pods. For more
  information, see [Considerations for IPAM
  Configuration](#considerations-for-ipam-configuration).
- All UCP worker and manager nodes need to be attached to the same Azure
  Subnet.
- All UCP workers and managers need to have the `Orchestrator` tag in Azure with the Kubernetes version as the value, following this format: `Orchestrator=Kubernetes:x.y.z`. This value may change in
  each UCP release. To find the relevant version please see the UCP [Release
  Notes](../../release-notes). For example for UCP `3.1.0` the tag would be
  `Orchestrator=Kubernetes:1.11.2`. 
- The Azure Virtual Machine Object Name needs to match the Azure Virtual Machine
  Computer Name and the Node Operating System's Hostname (This is the FQDN of
  the host, including domain names). Note this value is case sensitive, and all
  values should be in lowercase.
- An Azure Service Principal with `Contributor` access to the Azure Resource
  Group hosting the UCP Nodes. If you are using a separate Resource Group for
  the networking components, the same Service Principal will need `Network
  Contributor` access to this Resource Group.

UCP requires the following information for the installation:

- `subscriptionId` - The Azure Subscription ID in which the UCP 
objects are being deployed. 
- `tenantId` - The Azure Active Directory Tenant ID in which the UCP 
objects are being deployed. 
- `aadClientId` - The Azure Service Principal ID
- `aadClientSecret` - The Azure Service Principal Secret Key

### Azure Configuration File

For Docker UCP to integrate into Microsoft Azure, you need to place an Azure 
configuration file within each UCP node in your cluster, at 
`/etc/kubernetes/azure.json`. The `azure.json` file needs 0644 permissions.

See the template below. Note entries that do not contain `****` should not be 
changed.

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

There are some optional values for Azure deployments:

- `"primaryAvailabilitySetName": "****",` - The Worker Nodes availability set.
- `"vnetResourceGroup": "****",` - If your Azure Network objects live in a 
seperate resource group.
- `"routeTableName": "****",` - If you have defined multiple Route tables within
an Azure subnet.

More details on this configuration file can be found 
[here](https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/azure/azure.go).

## Considerations for IPAM Configuration

The subnet and the virtual network associated with the primary interface of
the Azure virtual machines need to be configured with a large enough address prefix/range. 
The number of required IP addresses depends on the number of pods running
on each node and the number of nodes in the cluster.

For example, in a cluster of 256 nodes, to run a maximum of 128 pods
concurrently on a node, make sure that the address space of the subnet and the
virtual network can allocate at least 128 * 256 IP addresses, _in addition to_
initial IP allocations to virtual machine NICs during Azure resource creation.

Accounting for IP addresses that are allocated to NICs during virtual machine bring-up, set
the address space of the subnet and virtual network to 10.0.0.0/16. This
ensures that the network can dynamically allocate at least 32768 addresses,
plus a buffer for initial allocations for primary IP addresses.

> Azure IPAM, UCP, and Kubernetes
> 
> The Azure IPAM module queries an Azure virtual machine's metadata to obtain
> a list of IP addresses that are assigned to the virtual machine's NICs. The
> IPAM module allocates these IP addresses to Kubernetes pods. You configure the
> IP addresses as `ipConfigurations` in the NICs associated with a virtual machine
> or scale set member, so that Azure IPAM can provide them to Kubernetes when
> requested.
{: .important}

## Manually provision IP address as part of an Azure virtual machine scale set

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

## Install UCP 

### Adjusting the IP Count Value

If you have manually attached additional IP addresses to the Virtual Machines (via an ARM Template, Azure CLI or Azure Portal) or you want to reduce the number of IP Addresses automatically provisioned by UCP from the default of 128 addresses, then you should alter the IP Count Value in the UCP Configuration file before installation. If you are happy with 128 addresses per Virtual Machine, proceed to [Installing UCP](#installing-ucp).

If you have manually attached additional IP addresses then your UCP configuration file would be: 

```
$ cat example-config-1
[cluster_config]
  azure_ip_count = "0"
```

If you want to reduce the IP Addresses dynamically allocated from 128 to a
custom value (see [Considerations for IPAM
Configuration](#considerations-for-ipam-configuration) to calculate an
appropriate value) then your UCP configuration file would be: 

```
$ cat example-config-2
[cluster_config]
  azure_ip_count = "20" # Note this value may be different for your environment
```

To pre-load this configuration file prior to a UCP installation, connect to the Virtual Machine where you would you like to install UCP and initiate a Swarm. 

```
$ docker swarm init
```

Then preload the UCP configuration file as a Swarm Config.

```
$ docker config create com.docker.ucp.config example-config-<n>
```

Now you are ready to proceed with the install. 

If you need to adjust this value post installation then this can be changed via
the UCP [configuration file](../configure/ucp-configuration-file/). However if
this value is reduced post-installation, existing virtual machines will not be
reconciled, and you will have to manually edit the IP count in Azure.  

### Installing UCP

Use the following command to install UCP on the manager node. The `--pod-cidr`
option maps to the IP address range that you have configured for the Azure
subnet, and the `--host-address` maps to the private IP address of the master
node.

> Note: The `pod-cidr` range must match the Azure Virtual Network's Subnet
> attached the hosts. For example if the Azure Virtual Network had the range
> `172.0.0.0/16` with Virtual Machines provisioned on an Azure Subnet of
> `172.0.1.0/24`, then the Pod CIDR should also be `172.0.1.0/24`

```bash
docker container run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} install \
  --host-address <ucp-ip> \
  --pod-cidr <ip-address-range> \
  --cloud-provider Azure \
  --interactive
```
