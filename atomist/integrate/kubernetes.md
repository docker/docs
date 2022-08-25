---
title: Kubernetes admission controllers
description:
keywords:
---

{% include atomist/disclaimer.md %}

## Prerequisites

1.  [`kustomize`][kustomize] (tested with `v4.4.1`)
2.  [`kubectl`][kubectl] (tested with `v1.21.1`)
3.  `kubectl` must be authenticated and the current context should be set to the
    cluster that you will be updating

[kustomize]: https://kubectl.docs.kubernetes.io/installation/kustomize/
[kubectl]: https://kubectl.docs.kubernetes.io/installation/kubectl/

## Fork repository

[Fork this repository](https://github.com/atomisthq/adm-ctrl/fork).

This repository contains a set of base Kubernetes resources that you'll adapt
using `kustomize`.

### API key and endpoint URL

Create an overlay for your cluster. Choose a cluster name and then create a new
overlay directory.

```
$ CLUSTER_NAME=replacethis
$ mkdir -p resources/k8s/overlays/${CLUSTER_NAME}
```

Create a file named `resources/k8s/overlays/${CLUSTER_NAME}/endpoint.env`.

```properties
apiKey=<replace this>
url=<replace this>
team=<replace this>
```

The `apiKey` and `url` must be filled in with your values from your Atomist
workspace. You can find these values in the
[integrations tab](https://dso.atomist.com/r/auth/integrations).

You'll also need your `id` for your Atomist `team`. This is the nine character
value that you'll find at the top of
[this page](https://dso.atomist.com/r/auth/policies).

![workspace id](./img/kubernetes/settings.png)

### Update Kubernetes cluster

This procedure will create a service account, a cluster role binding, two
secrets, a service, and a deployment. All will be created in a new namespaced
called `atomist`.

![controller diagram](./img/kubernetes/controller.png)

Use the same overlay that you created above
(`resources/k8s/overlays/${CLUSTER_NAME}`). Copy in a template
kustomization.yaml file.

```
$ cp resources/templates/default_controller.yaml resources/k8s/overlays/${CLUSTER_NAME}/kustomization.yaml
```

This Kustomization file requires one edit. The last line shows a patch with a
`value` of `"default"`. This should be updated to be the name of your cluster.

```yaml
resources:
  - ../../controller
secretGenerator:
  - name: endpoint
    behavior: merge
    envs:
      - endpoint.env
images:
  - name: atomist/adm-ctrl
    newTag: v4
patchesJson6902:
  - target:
      group: apps
      version: v1
      kind: Deployment
      name: policy-controller
    patch: |-
      - op: replace
        path: /spec/template/spec/containers/0/env/2/value
        value: "default"
```

Deploy the admission controller into the current Kubernetes context by running
the following script.

```bash
# creates roles and service account for running jobs
kustomize build resources/k8s/certs | kubectl apply -f -
# create SSL certs for your new admission controller - stores them as secrets in the atomist namespace of your cluster
kubectl apply -f resources/k8s/jobs/create.yaml
# creates a secret to store the keystore for your new admission controller
kubectl apply -f resources/k8s/jobs/keystore_secret.yaml
# install admission controller pod
kustomize build resources/k8s/overlays/${CLUSTER_NAME} | kubectl apply -f -
# skip the kube-system namespace
kubectl label namespace kube-system policy-controller.atomist.com/webhook=ignore
# validating webhook configuration
kubectl apply -f resources/k8s/admission/admission.yaml
# patch the admission webhook with the ca certificate generated earlier
kubectl apply -f resources/k8s/jobs/patch.yaml
```

## Enable image policy

The admission controller will still be admitting all pods. Reject pods with
un-checked images by enabling the policy one namespace at a time. For example,
start verifying that new pods in namespace `production` must have passed all
necessary rules by annotating the namespace with the annotation
`policy-controller.atomist.com/policy`.

```bash
$ kubectl annotate namespace production policy-controller.atomist.com/policy=enabled
```

Disable admission control on a namespace by removing the annotation or setting
it to something other than `enabled`.

```bash
$ kubectl annotate namespace production policy-controller.atomist.com/policy-
```

[dynamic-admission-control]:
  https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
