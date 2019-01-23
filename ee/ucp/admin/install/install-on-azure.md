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

Docker UCP configures the Azure IPAM module for Kubernetes to allocate
IP addresses to Kubernetes pods.  The Azure IPAM module requires each Azure
virtual machine that's part of the Kubernetes cluster to be configured with a pool of
IP addresses.

There are two options for provisoning IPs for the Kubernetes cluster on Azure
- Docker UCP provides an automated mechanism to configure and maintain IP pools 
  for standalone Azure virtual machines. This service runs within the calico-node daemonset 
  and by default will provision 128 IP address for each node. This value can be 
  configured through the `azure_ip_count`in the UCP 
  [configuration file](../configure/ucp-configuration-file) before or after the 
  UCP installation. Note that if this value is reduced post-installation, existing 
  virtual machines will not be reconciled, and you will have to manually edit the IP count
  in Azure. 
- Manually provision additional IP address for each Azure virtual machine. This could be done
  as part of an Azure Virtual Machine Scale Set through an ARM template. You can find an example [here](#set-up-ip-configurations-on-an-azure-virtual-machine-scale-set). 
  Note that the `azure_ip_count` value in the UCP 
  [configuration file](../configure/ucp-configuration-file) will need to be set
  to 0, otherwise UCP's IP Allocator service will provision the IP Address on top of 
  those you have already provisioned.

## Azure Prerequisites 

You must meet these infrastructure prerequisites in order 
to successfully deploy Docker UCP on Azure

- All UCP Nodes (Managers and Workers) need to be deployed into the same 
Azure Resource Group. The Azure Networking (Vnets, Subnets, Security Groups) 
components could be deployed in a second Azure Resource Group.
- All UCP Nodes (Managers and Workers) need to be attached to the same 
Azure Subnet.
- All UCP (Managers and Workers) need to be tagged in Azure with the 
`Orchestrator` tag. Note the value for this tag is the Kubernetes version number
in the format `Orchestrator=Kubernetes:x.y.z`. This value may change in each 
UCP release. To find the relevant version please see the UCP 
[Release Notes](../../release-notes). For example for UCP 3.1.0 the tag 
would be `Orchestrator=Kubernetes:1.11.2`. 
- The Azure Virtual Machine Object Name needs to match the Azure Virtual Machine 
Computer Name and the Node Operating System's Hostname. Note this applies to the 
FQDN of the host including domain names. 
- An Azure Service Principal with `Contributor` access to the Azure Resource 
Group hosting the UCP Nodes. Note, if using a separate networking Resource 
Group the same Service Principal will need `Network Contributor` access to this 
Resource Group.

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

Use the following command to install UCP on the manager node.
The `--pod-cidr` option maps to the IP address range that you configured for
the subnets in the previous sections, and the `--host-address` maps to the
IP address of the master node.

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

#### Additional Notes

- The Kubernetes `pod-cidr` must match the Azure Vnet of the hosts. 
