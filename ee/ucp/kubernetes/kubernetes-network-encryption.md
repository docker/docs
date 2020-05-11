---
title: Kubernetes Network Encryption
description: Learn how to configure network encryption in Kubernetes
keywords: ucp, cli, administration, kubectl, Kubernetes, security, network, ipsec, ipip, esp, calico
---

>{% include enterprise_label_shortform.md %}

Docker Enterprise Edition provides data-plane level IPSec network encryption to securely encrypt application 
traffic in a Kubernetes cluster. This secures application traffic within a cluster when running in untrusted 
infrastructure or environments. It is an optional feature of UCP that is enabled by deploying the SecureOverlay 
components on Kubernetes when using the default Calico driver for networking configured for IPIP tunneling 
(the default configuration).

Kubernetes network encryption is enabled by two components in UCP: the SecureOverlay Agent and SecureOverlay 
Master. The agent is deployed as a per-node service that manages the encryption state of the data plane. The 
agent controls the IPSec encryption on Calico’s IPIP tunnel traffic between different nodes in the Kubernetes 
cluster. The master is the second component, deployed on a UCP manager node, which acts as the key management 
process that configures and periodically rotates the encryption keys.

Kubernetes network encryption uses AES Galois Counter Mode (AES-GCM) with 128-bit keys by default. Encryption 
is not enabled by default and requires the SecureOverlay Agent and Master to be deployed on UCP to begin 
encrypting traffic within the cluster. It can be enabled or disabled at any time during the cluster lifecycle. 
However, it should be noted that it can cause temporary traffic outages between pods during the first few minutes 
of traffic enabling/disabling. When enabled, Kubernetes pod traffic between hosts is encrypted at the IPIP tunnel 
interface in the UCP host.

![Kubernetes Network Encryption](/ee/images/kubernetes-network-encryption.png)

## Requirements

Kubernetes network encryption is supported for the following platforms:
* Docker Enterprise 2.1+ (UCP 3.1+)
* Kubernetes 1.11+
* On-premise, AWS, GCE
* Azure is not supported for network encryption as encryption utilizes Calico’s IPIP tunnel
* Only supported when using UCP’s default Calico CNI plugin
* Supported on all Docker Enterprise supported Linux OSes

## Configuring MTUs

Before deploying the SecureOverlay components, ensure that Calico is configured so that the IPIP tunnel 
MTU maximum transmission unit (MTU), or the largest packet length that the container will allow, leaves sufficient headroom for the encryption overhead.   Encryption adds 26 bytes of overhead, but every IPSec 
packet size must be a multiple of 4 bytes.  IPIP tunnels require 20 bytes of encapsulation overhead.  The IPIP 
tunnel interface MTU must be no more than "EXTMTU - 46 - ((EXTMTU - 46) modulo 4)", where EXTMTU is the minimum MTU 
of the external interfaces.   An IPIP MTU of 1452 should generally be safe for most deployments. 

Changing UCP's MTU requires updating the UCP configuration.  This process is described [here](../admin/configure/ucp-configuration-file.md).  

Update the following values to the new MTU:

     [cluster_config]
      ...
      calico_mtu = "1452"
      ipip_mtu = "1452"
      ...

## Configuring SecureOverlay

SecureOverlay allows you to enable IPSec network encryption in Kubernetes. Once the cluster nodes’ MTUs are properly configured, deploy the SecureOverlay components using the SecureOverlay YAML file to UCP.

Beginning with UCP 3.2.4, you can configure SecureOverlay in two ways:
* Using the UCP configuration file or
* Using the SecureOverlay YAML file 

### UCP configuration file

Add `secure-overlay` to the UCP configuration file. Set this option to `true` to enable IPSec network encryption. The default is `false`. See [cluster_config options](../admin/configure/ucp-configuration-file.md#cluster_config-table-required) for more information.

### SecureOverlay YAML file

First, [download the SecureOverlay YAML file.](ucp-secureoverlay.yml)

Next, issue the following command from any machine with the properly configured kubectl environment and the proper UCP bundle's credentials:

```
$ kubectl apply -f ucp-secureoverlay.yml
```

Run this command at cluster installation time before starting any workloads.

To remove the encryption from the system, issue the following command:

```
$ kubectl delete -f ucp-secureoverlay.yml
```
