---
title: Custom Azure Roles
description: Learn how to create custom RBAC roles to run Docker Enterprise on Azure.
keywords: Universal Control Plane, UCP, install, Docker Enterprise, Azure, Swarm
---

>{% include enterprise_label_shortform.md %}

## Overview

This document describes how to create Azure custom roles to deploy Docker Enterprise resources.

## Deploy a Docker Enterprise Cluster into a single resource group

A [resource group](https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-overview#resource-groups) is a container that holds resources for an Azure solution. These resources are the virtual machines (VMs), networks, and storage accounts associated with the swarm.

To create a custom, all-in-one role with permissions to deploy a Docker Enterprise cluster into a single resource group:

1. Create the role permissions JSON file.
```bash
{
  "Name": "Docker Platform All-in-One",
  "IsCustom": true,
  "Description": "Can install and manage Docker platform.",
  "Actions": [
    "Microsoft.Authorization/*/read",
    "Microsoft.Authorization/roleAssignments/write",
    "Microsoft.Compute/availabilitySets/read",
    "Microsoft.Compute/availabilitySets/write",
    "Microsoft.Compute/disks/read",
    "Microsoft.Compute/disks/write",
    "Microsoft.Compute/virtualMachines/extensions/read",
    "Microsoft.Compute/virtualMachines/extensions/write",
    "Microsoft.Compute/virtualMachines/read",
    "Microsoft.Compute/virtualMachines/write",
    "Microsoft.Network/loadBalancers/read",
    "Microsoft.Network/loadBalancers/write",
    "Microsoft.Network/loadBalancers/backendAddressPools/join/action",
    "Microsoft.Network/networkInterfaces/read",
    "Microsoft.Network/networkInterfaces/write",
    "Microsoft.Network/networkInterfaces/join/action",
    "Microsoft.Network/networkSecurityGroups/read",
    "Microsoft.Network/networkSecurityGroups/write",
    "Microsoft.Network/networkSecurityGroups/join/action",
    "Microsoft.Network/networkSecurityGroups/securityRules/read",
    "Microsoft.Network/networkSecurityGroups/securityRules/write",
    "Microsoft.Network/publicIPAddresses/read",
    "Microsoft.Network/publicIPAddresses/write",
    "Microsoft.Network/publicIPAddresses/join/action",
    "Microsoft.Network/virtualNetworks/read",
    "Microsoft.Network/virtualNetworks/write",
    "Microsoft.Network/virtualNetworks/subnets/read",
    "Microsoft.Network/virtualNetworks/subnets/write",
    "Microsoft.Network/virtualNetworks/subnets/join/action",
    "Microsoft.Resources/subscriptions/resourcegroups/read",
    "Microsoft.Resources/subscriptions/resourcegroups/write",
    "Microsoft.Security/advancedThreatProtectionSettings/read",
    "Microsoft.Security/advancedThreatProtectionSettings/write",
    "Microsoft.Storage/*/read",
    "Microsoft.Storage/storageAccounts/listKeys/action",
    "Microsoft.Storage/storageAccounts/write"
  ],
  "NotActions": [],
  "AssignableScopes": [
    "/subscriptions/6096d756-3192-4c1f-ac62-35f1c823085d"
  ]
}
```
2. Create the Azure RBAC role.
```bash
az role definition create --role-definition all-in-one-role.json
```

## Deploy Docker Enterprise compute resources

Compute resources act as servers for running containers. 

To create a custom role to deploy Docker Enterprise compute resources only:

1. Create the role permissions JSON file.
```bash
{
  "Name": "Docker Platform",
  "IsCustom": true,
  "Description": "Can install and run Docker platform.",
  "Actions": [
    "Microsoft.Authorization/*/read",
    "Microsoft.Authorization/roleAssignments/write",
    "Microsoft.Compute/availabilitySets/read",
    "Microsoft.Compute/availabilitySets/write",
    "Microsoft.Compute/disks/read",
    "Microsoft.Compute/disks/write",
    "Microsoft.Compute/virtualMachines/extensions/read",
    "Microsoft.Compute/virtualMachines/extensions/write",
    "Microsoft.Compute/virtualMachines/read",
    "Microsoft.Compute/virtualMachines/write",
    "Microsoft.Network/loadBalancers/read",
    "Microsoft.Network/loadBalancers/write",
    "Microsoft.Network/networkInterfaces/read",
    "Microsoft.Network/networkInterfaces/write",
    "Microsoft.Network/networkInterfaces/join/action",
    "Microsoft.Network/publicIPAddresses/read",
    "Microsoft.Network/virtualNetworks/read",
    "Microsoft.Network/virtualNetworks/subnets/read",
    "Microsoft.Network/virtualNetworks/subnets/join/action",
    "Microsoft.Resources/subscriptions/resourcegroups/read",
    "Microsoft.Resources/subscriptions/resourcegroups/write",
    "Microsoft.Security/advancedThreatProtectionSettings/read",
    "Microsoft.Security/advancedThreatProtectionSettings/write",
    "Microsoft.Storage/storageAccounts/read",
    "Microsoft.Storage/storageAccounts/listKeys/action",
    "Microsoft.Storage/storageAccounts/write"
  ],
  "NotActions": [],
  "AssignableScopes": [
    "/subscriptions/6096d756-3192-4c1f-ac62-35f1c823085d"
  ]
}
```
2. Create the Docker Platform RBAC role.
```bash
az role definition create --role-definition platform-role.json
```

## Deploy Docker Enterprise network resources

Network resources are services inside your cluster. These resources can include virtual networks, security groups, address pools, and gateways. 

To create a custom role to deploy Docker Enterprise network resources only:

1. Create the role permissions JSON file.
```bash
{
  "Name": "Docker Networking",
  "IsCustom": true,
  "Description": "Can install and manage Docker platform networking.",
  "Actions": [
    "Microsoft.Authorization/*/read",
    "Microsoft.Network/loadBalancers/read",
    "Microsoft.Network/loadBalancers/write",
    "Microsoft.Network/loadBalancers/backendAddressPools/join/action",
    "Microsoft.Network/networkInterfaces/read",
    "Microsoft.Network/networkInterfaces/write",
    "Microsoft.Network/networkInterfaces/join/action",
    "Microsoft.Network/networkSecurityGroups/read",
    "Microsoft.Network/networkSecurityGroups/write",
    "Microsoft.Network/networkSecurityGroups/join/action",
    "Microsoft.Network/networkSecurityGroups/securityRules/read",
    "Microsoft.Network/networkSecurityGroups/securityRules/write",
    "Microsoft.Network/publicIPAddresses/read",
    "Microsoft.Network/publicIPAddresses/write",
    "Microsoft.Network/publicIPAddresses/join/action",
    "Microsoft.Network/virtualNetworks/read",
    "Microsoft.Network/virtualNetworks/write",
    "Microsoft.Network/virtualNetworks/subnets/read",
    "Microsoft.Network/virtualNetworks/subnets/write",
    "Microsoft.Network/virtualNetworks/subnets/join/action",
    "Microsoft.Resources/subscriptions/resourcegroups/read",
    "Microsoft.Resources/subscriptions/resourcegroups/write"
  ],
  "NotActions": [],
  "AssignableScopes": [
    "/subscriptions/6096d756-3192-4c1f-ac62-35f1c823085d"
  ]
}
```
2. Create the Docker Networking RBAC role.
```bash
az role definition create --role-definition networking-role.json
```

## Where to go next
* [Azure Container Instances documentation](https://docs.microsoft.com/en-us/azure/container-instances/)
* [docker/ucp overview](https://docs.docker.com/reference/ucp/3.2/cli/)
* [Universal Control Plane overview](https://docs.docker.com/ee/ucp/)
