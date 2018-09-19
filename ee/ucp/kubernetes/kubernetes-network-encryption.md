---
title: Kubernetes Network Encryption
description: Learn how to configure network encryption in Kubernetes
keywords: ucp, cli, administration, kubectl, Kubernetes, security, network, ipsec, ipip, esp, calico
---

Docker Enterprise Edition provides data-plane level IPSec network encryption to securely encrypt application 
traffic in a Kubernetes cluster. This secures application traffic within a cluster when running in untrusted 
infrastructure or environments. It is an optional feature of UCP that is enabled by deploying the Secure Overlay 
components on Kuberenetes when using the default Calico driver for networking configured for IPIP tunnelling 
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

Kubernetes Network Encryption is supported for the following platforms:
* Docker Enterprise 2.1+ (UCP 3.1+)
* Kubernetes 1.11+
* On-premise, AWS, GCE
* Azure is not supported for network encryption as encryption utilizes Calico’s IPIP tunnel
* Only supported when using UCP’s default Calico CNI plugin
* Supported on all Docker Enterprise supported Linux OSes

## Configuring MTUs

Before deploying the SecureOverlay components one must ensure that Calico is configured so that the IPIP tunnel 
MTU leaves sufficient headroom for the encryption overhead.   Encryption adds 26 bytes of overhead but every IPSec 
packet size must be a multiple of 4 bytes.  IPIP tunnels require 20 bytes of encapsulation overhead.  So the IPIP 
tunnel interface MTU must be no more than "EXTMTU - 46 - ((EXTMTU - 46) modulo 4)" where EXTMTU is the minimum MTU 
of the external interfaces.   An IPIP MTU of 1452 should generally be safe for most deployments. 

Changing UCP's MTU requires updating the UCP configuration.  This process is described [here](/ee/ucp/admin/configure/ucp-configuration-file).  

The user must update the following values to the new MTU:

     [cluster_config]
      ...
      calico_mtu = "1452"
      ipip_mtu = "1452"
      ...

## Configuring SecureOverlay

Once the cluster nodes’ MTUs are properly configured, deploy the SecureOverlay components using the Secure Overlay YAML file to UCP.

[Download the Secure Overlay YAML file here.](ucp-secureoverlay.yml)

After downloading the YAML file, run the following command from any machine with the properly configured kubectl environment and the proper UCP bundle's credentials:

```
$ kubectl apply -f ucp-secureoverlay.yml
```

Run this command at cluster installation time before starting any workloads.

To remove the encryption from the system, issue the command:

```
$ kubectl delete -f ucp-secureoverlay.yml
```
