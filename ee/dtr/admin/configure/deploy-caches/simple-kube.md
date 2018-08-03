---
title: Deploy a DTR cache with Kubernetes
description: Deploy a DTR cache to make users in remove geographical locations
  pull images faster.
keywords: DTR, cache, kubernetes
---

This example guides you through deploying a DTR cache, assuming that you've got
a DTR deployment up and running. It also assumes that you've provisioned
a Kubernetes Cluster of 1.8 or higher.

The DTR cache is going to be deployed as a Kubernetes Deployment, so that Kubernetes
automatically takes care of scheduling and restarting the service if
something goes wrong.

We'll manage the cache configuration using a Kubernetes Config Map, and the TLS
certificates using Kubernetes secrets. This allows you to manage the configurations
securely and independently of the node where the cache is actually running.

## Prepare the cache deployment

At the end of this exercise you should have a file structure that looks like this:

```
├── dtrcache.yml        # Yaml file to deploy cache with a single command
├── config.yml          # The cache configuration file
└── certs
    ├── cache.cert.pem  # The cache public key certificate
    ├── cache.key.pem   # The cache private key
    └── dtr.cert.pem    # DTR CA certificate
```

### Create the DTR Cache certicates

The DTR cache will be deployed with a TLS endpoint. For this you will need to generate
a SSL ceritificate and key from a certificate authority. Depending on how you would
 like to expose the DTR Cache it will depend on the SANs required for this certificate.

For example:

  - If you are deploying the DTR Cache with an 
    [Ingress Object](https://kubernetes.io/docs/concepts/services-networking/ingress/) 
    you will need to use the FQDN of your `Ingress Object` as part of your certificate. 
  - If you exposing the DTR cache through a Kubernetes
    [Cloud Provider](https://kubernetes.io/docs/concepts/services-networking/#loadbalancer)
    then you will need the external Loadbalancer address as part of your certificate. 
  - If you are exposing the DTR Cache through a 
    [Node Port](https://kubernetes.io/docs/concepts/services-networking/#nodeport)` or a
    `Host Port]` you will need to use a Node's FQDN as a SAN in your certificate. You could use
    a Kubernetes `NodeScheduler` to pin the node the DTR cache is deployed on.

On your workstation, create a directory called `certs`. Within here place the newly 
created certificate `cache.cert.pem` and key `cache.key.pem` for your DTR cache. Also
 place the certificate authority (including any intermedite certificate authorities)
 of the certificate from your DTR deployment. This could be sourced from curl.
`curl -s https://<dtr-fqdn>/ca -o certs/dtr.cert.pem`.

### Create the DTR Config

This is the configuration file of the registry component of the DTR Cache. This yaml
should be updated with the relevant `cache-fqdn` and `dtr-fqdn`.

```
version: 0.1
log:
  level: info
storage:
  delete:
    enabled: true
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: 0.0.0.0:443
  secret: generate-random-secret
  host: https://<external-fqdn-dtrcache> # Could be Ingress / LB / Host
  tls:
    certificate: /certs/cache.cert.pem
    key: /certs/cache.key.pem
middleware:
  registry:
      - name: downstream
        options:
          blobttl: 24h
          upstreams:
            - https://<dtr-url>
          cas:
            - /certs/dtr.cert.pem
```

With this configuration, the cache fetches image layers from DTR and keeps
a local copy for 24 hours. After that, if a user requests that image layer,
the cache fetches it again from DTR.

The cache is configured to persist data inside its container.
If something goes wrong with the cache service, Docker automatically redeploys a
new container, but previously cached data is not persisted.
You can [customize the storage parameters](/registry/configuration.md#storage),
if you want to store the image layers using a persistent storage backend.

### Create Kubernetes Resources

Create a Kubernetes namespace to logical seperate all of our DTR cache components.

```
$ kubectl create ns dtr
```

Create the Kubernetes Secrets and the Kubernetes ConfigMaps. Note the commands below
will only work if your workstation directory was configured our examples structure 
above. 

```
$ kubectl create secret generic dtr-certs \
    --from-file=certs/dtr.cert.pem \
    --from-file=certs/cache.cert.pem \
    --from-file=certs/cache.key.pem

$ kubectl create configmap dtr-cache-config \
    --from-file=config.yml
```

Create the Kubernetes Deployment. 

```
cat <<EOF | kubectl create -f -
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: dtr-cache
  namespace: dtr
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dtr-cache
  template:
    metadata:
      labels:
        app: dtr-cache
      annotations:
       seccomp.security.alpha.kubernetes.io/pod: docker/default
    spec:
      containers:
        - name: dtr-cache
          image: {{ page.dtr_org }}/{{ page.dtr_repo }}-content-cache:{{ page.dtr_version }}
          command: ["bin/sh"]
          args:
            - start.sh
            - /config/config.yml
          ports:
          - name: https
            containerPort: 443
          volumeMounts:
          - name: dtr-certs
            readOnly: true
            mountPath: /certs/
          - name: dtr-cache-config
            readOnly: true
            mountPath: /config
      volumes:
      - name: dtr-certs
        secret:
          secretName: dtr-certs
      - name: dtr-cache-config
        configMap:
          defaultMode: 0666
          name: dtr-cache-config
EOF
```

You can check if the deployment has been successful with `kubectl get pods -n dtr`.
If you need to troubleshoot your deployment, you can use `kubectl describe -n dtr <pods>`
and / or `kubectl logs -n dtr <pods>`

For external access to the DTR cache you could expose your DTR cache through 
multiple interfaces: Ingress Object, Node Port, Host Port or Loadbalancer. 
The following yaml will expose the DTR cache through an ingress obect.

> Note an ingress controller is a pre-requsite for this example. If you have not deploy
> an ingress controller on your cluster, here are instructions for the Nginx ingress
> controller for [UCP](ucp/kubernetes/layer-7-routing). 

```
cat <<EOF | kubectl create -f -
kind: Service
apiVersion: v1
metadata:
  name: dtr-cache
  namespace: dtr
spec:
  selector:
    app: dtr-cache
  ports:
  - protocol: TCP
    port: 443
    targetPort: 443
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: dtr-cache
  namespace: dtr
  annotations:
    nginx.ingress.kubernetes.io/ssl-passthrough: "true"
    nginx.ingress.kubernetes.io/secure-backends: "true"
spec:
  tls:
  - hosts:
    - <external-dtr-cache-fqdn>
  rules:
  - host: <external-dtr-cache-fqdn>
    http:
      paths:
      - backend:
          serviceName: dtr-cache
          servicePort: 443
```

You can test that your DTR cache is externally reachable by curling the API endpoint.

```
curl -X GET https://<dtr-cache-endpoint>/v2/_catalog
{"repositories":[]}
```


## Next Steps 

[Integrate your cache into DTR and configure users](simple.md#register-the-cache-with-dtr)