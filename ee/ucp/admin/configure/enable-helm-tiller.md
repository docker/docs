---
title: Enable Helm and Tiller with UCP
description: Learn how to modify service accounts to enable Helm and Tiller to operate with UCP.
keywords: Helm, ucp, Tiller, Kubernetes, service accounts, Kubernetes
---

To use [Helm and Tiller](https://helm.sh/) with UCP, you must modify the `kube-system` default service account to define the necessary roles. Enter the following `kubectl` commands in this order:

```
kubectl create rolebinding default-view --clusterrole=view --serviceaccount=kube-system:default --namespace=kube-system

kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
```

## Using Helm

For more information about using Helm, see [Using Helm - Role-Based Access Control](https://docs.helm.sh/using_helm/#role-based-access-control).
