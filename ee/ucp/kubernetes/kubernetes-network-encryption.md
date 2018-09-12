---
title: Kubernetes Network Encryption
description: Learn how to configure network encryption in Kubernetes
keywords: ucp, cli, administration, kubectl, Kubernetes, security, network, ipsec, ipip, esp, calico
---

Docker Enterprise provides data-plane level IPSec network encryption to securely encrypt application traffic in 
a Kubernetes cluster. This secures application traffic within a cluster when running in untrusted infrastructure 
or environments. It is an optional feature of UCP that is enabled by deploying the Secure Overlay components on 
Kuberenetes when using the default Calico driver for networking configured for IPIP tunnelling (the default configuration).

Kubernetes network encryption is enabled by two components in UCP: the SecureOverlay Agent and SecureOverlay Master. 
The agent is deployed as a per-node service that manages the encryption state of the data plane. The agent controls 
the IPSec encryption on Calico’s IPIP tunnel traffic between different nodes in the Kubernetes cluster. The master 
is the second component, which acts as the key management process that configures and periodically rotates the 
encryption keys.

Kubernetes network encryption uses AES-GCM with 128-bit keys (by default) and encrypts traffic between pods residing 
on different nodes. Encryption is not enabled by default and requires the SecureOverlay Agent and Master to be deployed 
on UCP to begin encrypting traffic within the cluster. It can be enabled or disabled at any time during the cluster 
lifecycle. However, note that enabling or disabling traffic can cause temporary inter-pod traffic outages during the 
first few minutes of encryption reconfiguration. When enabled, Kubernetes pod traffic between hosts is encrypted at 
the IPIP tunnel interface in the UCP host.

## Requirements

Kubernetes Network Encryption is supported for the following platforms:
* Docker Enterprise 2.1+ (UCP 3.1+)
* Kubernetes 1.11+
* On-prem, AWS, GCE (*Azure is not supported for network encryption as encryption utilizes Calico’s IPIP tunnel which is not supported in Azure)
* Only supported when using UCP’s default Calico CNI plugin
* Supported on all Docker Enterprise supported Linux OSes

## Configuring SecureOverlay

Once the cluster nodes’ MTUs are properly configured, deploy the SecureOverlay components using the following YAML file to UCP:

```
######################
# Cluster role for key management jobs
######################
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: ucp-secureoverlay-mgr
rules:
  - apiGroups: [""]
    resources:
      - secrets
    verbs:
      - get
      - update
---
######################
# Cluster role binding for key management jobs
######################
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: ucp-secureoverlay-mgr
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ucp-secureoverlay-mgr
subjects:
- kind: ServiceAccount
  name: ucp-secureoverlay-mgr
  namespace: kube-system
---
######################
# Service account for key management jobs
######################
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ucp-secureoverlay-mgr
  namespace: kube-system
---
######################
# Cluster role for secure overlay per-node agent
######################
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: ucp-secureoverlay-agent
rules:
  - apiGroups: [""]
    resources:
      - nodes
    verbs:
      - get
      - list
      - watch
---
######################
# Cluster role binding for secure overlay per-node agent
######################
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: ucp-secureoverlay-agent
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ucp-secureoverlay-agent
subjects:
- kind: ServiceAccount
  name: ucp-secureoverlay-agent
  namespace: kube-system
---
######################
# Service account secure overlay per-node agent
######################
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ucp-secureoverlay-agent
  namespace: kube-system
---
######################
# K8s secret of current key configuration
######################
apiVersion: v1
kind: Secret
metadata:
  name: ucp-secureoverlay
  namespace: kube-system
type: Opaque
data:
  keys: ""
---
######################
# DaemonSet for secure overlay per-node agent
######################
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: ucp-secureoverlay-agent
  namespace: kube-system
  labels:
    k8s-app: ucp-secureoverlay-agent
spec:
  selector:
    matchLabels:
      k8s-app: ucp-secureoverlay-agent
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: ucp-secureoverlay-agent
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      hostNetwork: true
      priorityClassName: system-node-critical
      terminationGracePeriodSeconds: 10
      serviceAccountName: ucp-secureoverlay-agent
      containers:
      - name: ucp-secureoverlay-agent
        image: docker/ucp-secureoverlay-agent:3.1.0
        securityContext:
          capabilities:
            add: ["NET_ADMIN"]
        env:
        - name: MY_NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        volumeMounts:
        - name: ucp-secureoverlay
          mountPath: /etc/secureoverlay/
          readOnly: true
      volumes:
      - name: ucp-secureoverlay
        secret:
          secretName: ucp-secureoverlay
---
######################
# Deployment for manager of the whole cluster (to rotate keys)
######################
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ucp-secureoverlay-mgr
  namespace: kube-system
spec:
  selector:
    matchLabels:
      app: ucp-secureoverlay-mgr
  replicas: 1
  template:
    metadata:
      name: ucp-secureoverlay-mgr
      namespace: kube-system
      labels:
        app: ucp-secureoverlay-mgr
    spec:
      serviceAccountName: ucp-secureoverlay-mgr
      restartPolicy: Always
      containers:
      - name: ucp-secureoverlay-mgr
        image: docker/ucp-secureoverlay-mgr:3.1.0
```

After one downloads the YAML file, run the following command from any machine with the properly configured kubectl environment and the proper UCP bundle's credentials: 

```
$ kubectl apply -f ucp-secureoverlay.yml
```

Run this command at cluster installation time before starting any workloads.

To remove the encryption from the system, issue the command:

```
$ kubectl delete -f ucp-secureoverlay.yml
```
