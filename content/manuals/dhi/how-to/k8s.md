---
title: Use a Docker Hardened Image in Kubernetes
linktitle: Use an image in Kubernetes
description: Learn how to use Docker Hardened Images in Kubernetes deployments.
keywords: use hardened image, kubernetes, k8s
weight: 35
---

{{< summary-bar feature_name="Docker Hardened Images" >}}

## Authentication

To be able to use Docker Hardened Images in Kubernetes, you need to create a 
Kubernetes secret for pulling images from your mirror or internal registry.

> [!NOTE]
>
> You need to create this secret in each Kubernetes namespace that uses a DHI.

Create a secret using a Personal Access Token (PAT). Ensure the token has at least
read-only access to private repositories. For Docker Hub replace `<registry server>`
with `docker.io`.

```console
$ kubectl create -n <kubernetes namespace> secret docker-registry <secret name> --docker-server=<registry server> \
        --docker-username=<registry user> --docker-password=<access token> \
        --docker-email=<registry email>
```

To tests the secrets use the following command:

```console
kubectl apply --wait -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: dhi-test
  namespace: <kubernetes namespace>
spec:
  containers:
  - name: test
    image: <your-namespace>/dhi-bash:5
    command: [ "sh", "-c", "echo 'Hello from DHI in Kubernetes!'" ]
  imagePullSecrets:
  - name: <secret name>
EOF
```

Get the status of the pod by running:

```console
$ kubectl get -n <kubernetes namespace> pods/dhi-test
```

The command should return the following result:

```console
NAME       READY   STATUS      RESTARTS     AGE
dhi-test   0/1     Completed   ...          ...
```

If instead, the result is the following, there might be an issue with your secret.

```console
NAME       READY   STATUS         RESTARTS   AGE
dhi-test   0/1     ErrImagePull   0          ...
```

Verify the output of the pod by running, which should return `Hello from DHI in Kubernetes!`

```console
kubectl logs -n <kubernetes namespace> pods/dhi-test
```

After a successful test, the test pod can be deleted with the following command:

```console
$ kubectl delete -n <kubernetes namespace> pods/dhi-test
```
