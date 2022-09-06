---
title: Kubernetes admission controllers
description:
keywords:
---

{% include atomist/disclaimer.md %}

This page describes how to set up the Atomist-Kubernetes integration for using a
validating admission control webhook. This integration is a prerequisite for the
following use cases:

- Deployment policies: protecting your cluster from running vulnerable images.
- Deployment tracking: keeping track of what images are currently deployed in
  your environments.

## How the integration works

This integration deploys, among other things, a Clojure web service in your
cluster. The web service interacts with the Atomist webhook to retrieve and
upload data about images.

Data about images is retrieved if you are using
[deployment policies](../configure/deployment-policies.md) to control whether
images get deployed to Kubernetes based on policy checks.

The service will also upload data to Atomist. The payload that it sends
describes what images are deployed in the cluster, across all namespaces. This
helps you keep track of your deployed images, in what environment (staging,
production), and lets you see a delta between candidate images and what you are
currently running. Go to [deployment tracking](./deploys.md) to learn more.

## Configuration

The following parameters, defined as environment variables, are used to
configure the service:

- `ATOMIST_URL`: the webhook URL to use when talking back to Atomist. This can
  be found via the
  [Integrations page](https://dso.atomist.com/r/auth/integrations) on the
  Atomist website.
- `ATOMIST_WORKSPACE`: the ID of your Atomist workspace. This can be found via
  the [Integrations page](https://dso.atomist.com/r/auth/integrations) on the
  Atomist website.
- `ATOMIST_APIKEY`: key used to authenticate with the Atomist API. This needs to
  have sufficient privileges against the workspace above. API keys can be
  generated via the
  [Integrations page](https://dso.atomist.com/r/auth/integrations) on the
  Atomist website.
- `CLUSTER_NAME`: a name which is used to describe the cluster. This name will
  be used to track which clusters an image has been deployed to as it progresses
  so a name like `staging` or `production` is a likely value.

## Installation steps

Below are instructions for two methods of installing admission controllers in
Kubernetes:

- [Using flux](#install-using-flux)
- [Manually](#install-manually)

### Install using Flux

To install the controller using [Flux](https://fluxcd.io), start by creating a
`GitRepository` resource, pointing to the GitHub repository containing the
Kubernetes resources that the controller requires. Put that someplace where Flux
can find it.

```yaml
# gitrepo.yml
apiVersion: source.toolkit.fluxcd.io/v1beta1
kind: GitRepository
metadata:
  name: adm-ctrl
  namespace: flux-system
spec:
  interval: 30s
  ref:
    branch: main
  url: https://github.com/atomisthq/adm-ctrl
```

Using that source we can create a `Kustomization` which will allow us to pull in
the resources (from the `resources/k8s/controller` directory of the repository)
required by the controller. We'll want to customize the `CLUSTER_NAME`
environment variable in the controller deployment so we can use Kustomize to do
that. This file will also be the place where we specify which controller image
we are running.

```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
kind: Kustomization
metadata:
  name: adm-ctrl
  namespace: flux-system
spec:
  targetNamespace: atomist
  interval: 10m0s
  decryption:
    provider: sops
    secretRef:
      name: sops-gpg
  sourceRef:
    kind: GitRepository
    name: adm-ctrl
  path: ./resources/k8s/controller
  prune: true
  patches:
    - patch: |-
        apiVersion: apps/v1
        kind: Deployment
        metadata:
          name: policy-controller
          namespace: atomist
        spec:
          template:
            spec:
              containers:
                - name: controller
                  env:
                    - name: CLUSTER_NAME
                      value: production
      target:
        kind: Deployment
        name: policy-controller
  images:
    - newTag: v4-5-ga51c3ee
      name: atomist/adm-ctrl
```

For this example we're going to encode the remaining three environment variables
listed above (`ATOMIST_URL`, `ATOMIST_WORKSPACE`, `ATOMIST_APIKEY`) into a
single secret using [sops](https://fluxcd.io/docs/guides/mozilla-sops/). Once
that secret file has been created somewhere in the repo we'll need a
`kustomization.yaml` alongside it to let Flux know about it.

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - secret.yaml
```

We can then create another `Kustomization` file which will pull in all resources
in the directory we put the lat file. For this example that happens to be in
`./adm-ctrl/production` but it can, of course, be anywhere relevant to the
layout of your Flux repository.

```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
kind: Kustomization
metadata:
  name: adm-ctrl-resources
  namespace: flux-system
spec:
  interval: 10m0s
  decryption:
    provider: sops
    secretRef:
      name: sops-gpg
  sourceRef:
    kind: GitRepository
    name: flux-system
  path: ./adm-ctrl/production
  prune: true
```

Now you can commit these changes to your Flux repository and have the various
controllers pick up the changes and create all the necessary resources.

## Install manually

If you prefer to have complete control over the resources which are created in
your cluster then you can choose to install the various resources in the
`resources/k8s` directory of this repo manually. You can use
[Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/) to
generate the files to `kubectl apply` if that works for you.

### Prerequisites

1.  [Kustomize][kustomize] (tested with `v4.4.1`)
2.  [kubectl][kubectl] (tested with `v1.21.1`)
3.  kubectl must be authenticated and the current context should be set to the
    cluster that you will be updating

[kustomize]: https://kubectl.docs.kubernetes.io/installation/kustomize/
[kubectl]: https://kubectl.docs.kubernetes.io/installation/kubectl/

### Fork this repository

[Fork this repository](https://github.com/atomisthq/adm-ctrl/fork).

This repository contains a set of base Kubernetes that you'll adapt using
Kustomize.

### API key and endpoint URL

Create an overlay for your cluster. Choose a cluster name and then create a new
overlay directory.

```bash
CLUSTER_NAME=replacethis
mkdir -p resources/k8s/overlays/${CLUSTER_NAME}
```

Create a file named `resources/k8s/overlays/${CLUSTER_NAME}/endpoint.env`.

The file will start out like this.

```properties
apiKey=<replace this>
url=<replace this>
```

The `apiKey` and `url` should be filled in with your workspace values. Find
these in the [Integrations page](https://dso.atomist.com/r/auth/integrations) of
the Atomist app and replace them in the file.

### Create certificates for the admission controller

The communication between the api-server and the admission controller will be
over HTTPS. This will be configured by running 3 Kubernetes jobs in the cluster.

- `policy-controller-cert-create` job: this job will create SSL certificates and
  store them in a secret named `policy-controller-admission-cert` in the atomist
  namespace
- `policy-controller-cert-path` job: this will patch the admission controller
  webhook with the ca cert (so that the api-server will trust the new
  policy-controller)
- `keystore-create` job: this will read the SSL certificates created by the
  policy-controller-cert-create job and create a key store for the
  policy-controller HTTP server. The key store is also stored in a secret named
  `keystore` in the atomist namespace.

You can do steps 1 and 3 now.

```bash
# creates roles and service account for running jobs
kustomize build resources/k8s/certs | kubectl apply -f -
kubectl apply -f resources/k8s/jobs/create.yaml
kubectl apply -f resources/k8s/jobs/keystore_secret.yaml
```

### Update Kubernetes cluster

This procedure will create a service account, a cluster role binding, two
secrets, a service, and a deployment. All will be created in a new namespaced
called `atomist`.

![controller diagram](./docs/controller.png)

Use the same overlay that you created above
(`resources/k8s/overlays/${CLUSTER_NAME}`). Copy in a template
`kustomization.yaml` file.

```bash
cp resources/templates/default_controller.yaml resources/k8s/overlays/${CLUSTER_NAME}/kustomization.yaml
```

This `kustomization.yaml` file will permit you to change the `CLUSTER_NAME`
environment variable. In the initial copy of the file, the value will be
`"default"`, but it should be changed to the name of your cluster. This change
is made to the final line in your new `kustomization.yaml` file.

```yaml
resources:
  - ../../controller
secretGenerator:
  - name: endpoint
    behavior: merge
    envs:
      - endpoint.env
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

Deploy the admission controller into the current Kubernetes context using the
command shown below.

```bash
kustomize build resources/k8s/overlays/sandbox | kubectl apply -f -
```

At this point, the admission controller will be running but the cluster will not
be routing any admission control requests to it. Create a configuration to start
sending admission control requests to the controller using the following script.

```bash
# skip the kube-system namespace
kubectl label namespace kube-system policy-controller.atomist.com/webhook=ignore
# validating webhook configuration
kubectl apply -f resources/k8s/admission/admission.yaml
# finally, patch the admission webhook with the ca certificate generated earlier
kubectl apply -f resources/k8s/jobs/patch.yaml
```

## Enable image check policy

```bash
kubectl annotate namespace production policy-controller.atomist.com/policy=enabled
```

Disable policy on a namespace by removing the annotation or setting it to
something other than `enabled`.

```bash
kubectl annotate namespace production policy-controller.atomist.com/policy-
```
