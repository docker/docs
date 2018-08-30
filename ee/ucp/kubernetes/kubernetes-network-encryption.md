---
title: Kubernetes Network Encryption
description: Learn how to configure network encryption in Kubernetes
keywords: ucp, cli, administration, kubectl, Kubernetes, security, network, ipsec, ipip, esp, calico
---

Docker Enterprise provides data-plane level IPSec network encryption to securely encrypt application traffic in a Kubernetes cluster. This secures application traffic within a cluster when running in untrusted infrastructure or environments. It is an optional feature of UCP that is enabled by deploying the Secure Overlay components on Kuberenetes when using the default Calico driver for networking configured for IPIP tunnelling (the default configuration).
