---
title: Deploy a DTR cache with Kubernetes
description: Deploy a DTR cache to allow users in remote geographical locations to pull images faster.
keywords: DTR, cache, kubernetes
---

This example guides you through deploying a DTR cache, assuming that you've got
a DTR deployment up and running. The below guide has been tested on 
Universal Control Plane 3.1, however it should work on any Kubernetes Cluster 
1.8 or higher. 

The DTR cache is going to be deployed as a Kubernetes Deployment, so that 
Kubernetes automatically takes care of scheduling and restarting the service if
something goes wrong.

We'll manage the cache configuration using a Kubernetes Config Map, and the TLS
certificates using Kubernetes secrets. This allows you to manage the 
configurations securely and independently of the node where the cache is 
actually running.

## Prepare the cache deployment

At the end of this exercise you should have the following file structure on your
workstation:

```
├── dtrcache.yaml        # Yaml file to deploy cache with a single command
├── config.yaml          # The cache configuration file
└── certs
    ├── cache.cert.pem   # The cache public key certificate, including any intermediaries
    ├── cache.key.pem    # The cache private key
    └── dtr.cert.pem     # DTR CA certificate
```

### Create the DTR Cache certificates

The DTR cache will be deployed with a TLS endpoint. For this you will need to 
generate a TLS ceritificate and key from a certificate authority. The way you 
expose the DTR Cache will change the SANs required for 
this certificate.

For example:

  - If you are deploying the DTR Cache with an 
    [Ingress Object](https://kubernetes.io/docs/concepts/services-networking/ingress/) 
    you will need to use an external DTR cache address which resolves to your 
    ingress controller as part of your certificate. 
  - If you are exposing the DTR cache through a Kubernetes
    [Cloud Provider](https://kubernetes.io/docs/concepts/services-networking/#loadbalancer)
    then you will need the external Loadbalancer address as part of your 
    certificate. 
  - If you are exposing the DTR Cache through a 
    [Node Port](https://kubernetes.io/docs/concepts/services-networking/#nodeport)
     or a Host Port you will need to use a node's FQDN as a SAN in your 
     certificate.

On your workstation, create a directory called `certs`. Within it place the 
newly created certificate `cache.cert.pem` and key `cache.key.pem` for your DTR 
cache. Also place the certificate authority (including any intermedite 
certificate authorities) of the certificate from your DTR deployment. This could 
be sourced from the main DTR deployment using curl. 

```
$ curl -s https://<dtr-fqdn>/ca -o certs/dtr.cert.pem`.
```

### Create the DTR Config

The DTR Cache will take its configuration from a file mounted into the container. 
Below is an example configuration file for the DTR Cache. This yaml should be 
customised for your environment with the relevant external dtr cache, worker 
node or external loadbalancer FQDN.

With this configuration, the cache fetches image layers from DTR and keeps a 
local copy for 24 hours. After that, if a user requests that image layer, the 
cache will fetch it again from DTR.

The cache, by default, is configured to store image data inside its container.
Therefore if something goes wrong with the cache service, and Kubernetes deploys
a new pod, cached data is not persisted. Data will not be lost as it is still 
stored in the primary DTR. You can 
[customize the storage parameters](/registry/configuration/#storage),
if you want the cached images to be backended by persistent storage.

> **Note**: Kubernetes Peristent Volumes or Persistent Volume Claims would have to be
> used to provide persistent backend storage capabilities for the cache.

```
cat > config.yaml <<EOF
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
  host: https://<external-fqdn-dtrcache> # Could be DTR Cache / Loadbalancer / Worker Node external FQDN
  tls:
    certificate: /certs/cache.cert.pem
    key: /certs/cache.key.pem
middleware:
  registry:
      - name: downstream
        options:
          blobttl: 24h
          upstreams:
            - https://<dtr-url> # URL of the Main DTR Deployment
          cas:
            - /certs/dtr.cert.pem
EOF
```

See [Configuration Options](/registry/configuration/#list-of-configuration-options) for a full list of registry configuration options.

### Define Kubernetes Resources

The Kubernetes Manifest file to deploy the DTR Cache is independent of how you
choose to expose the DTR cache within your environment. The below example has
been tested to work on Universal Control Plane 3.1, however it should work on
any Kubernetes Cluster 1.8 or higher. 

```
cat > dtrcache.yaml <<EOF
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
            - /config/config.yaml
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

## Create Kubernetes Resources

At this point you should have a file structure on your workstation which looks 
like this:

```
├── dtrcache.yaml        # Yaml file to deploy cache with a single command
├── config.yaml          # The cache configuration file
└── certs
    ├── cache.cert.pem   # The cache public key certificate
    ├── cache.key.pem    # The cache private key
    └── dtr.cert.pem     # DTR CA certificate
```

You will also need the `kubectl` command line tool configured to talk to your 
Kubernetes cluster, either through a Kubernetes Config file or a [Universal
Control Plane client bundle](/ee/ucp/user-access/kubectl/).

First we will create a Kubernetes namespace to logically separate all of our 
DTR cache components.

```
$ kubectl create namespace dtr
```

Create the Kubernetes Secrets, containing the DTR cache TLS certificates, and a 
Kubernetes ConfigMap containing the DTR cache configuration file. 

```
$ kubectl -n dtr create secret generic dtr-certs \
    --from-file=certs/dtr.cert.pem \
    --from-file=certs/cache.cert.pem \
    --from-file=certs/cache.key.pem

$ kubectl -n dtr create configmap dtr-cache-config \
    --from-file=config.yaml
```

Finally create the Kubernetes Deployment. 

```
$ kubectl create -f dtrcache.yaml
```

You can check if the deployment has been successful by checking the running pods
in your cluster: `kubectl -n dtr get pods `

If you need to troubleshoot your deployment, you can use 
`kubectl -n dtr describe pods <pods>` and / or `kubectl -n dtr logs <pods>`. 

### Exposing the DTR Cache

For external access to the DTR cache we need to expose the Cache Pods to the 
outside world. In Kubernetes there are multiple ways for you to expose a service,
dependent on your infrastructure and your environment. For more information,
see [Publishing services - service types
](https://kubernetes.io/docs/concepts/services-networking/#publishing-services-service-types) on the Kubernetes docs.
It is important though that you are consistent in exposing the cache through the 
same interface you created a certificate for [previously](#create-the-dtr-cache-certificates).
Otherwise the TLS certificate may not be valid through this alternative 
interface.

> #### DTR Cache Exposure 
>
> You only need to expose your DTR cache through ***one*** external interface.

#### NodePort

The first example exposes the DTR cache through **NodePort**. In this example you would 
have added a worker node's FQDN to the TLS Certificate in [step 1](#create-the-dtr-cache-certificates).
Here you will be accessing the DTR cache through an exposed port on a worker 
node's FQDN.

```
cat > dtrcacheservice.yaml <<EOF
apiVersion: v1
kind: Service
metadata:
  name: dtr-cache
  namespace: dtr
spec:
  type: NodePort
  ports:
  - name: https
    port: 443
    targetPort: 443
    protocol: TCP
  selector:
    app: dtr-cache
EOF

kubectl create -f dtrcacheservice.yaml
```

To find out which port the DTR cache has been exposed on, you will need to run:

```
$ kubectl -n dtr get services 
```

You can test that your DTR cache is externally reachable by using `curl` to hit 
the API endpoint, using both a worker node's external address, and the **NodePort**.

```
curl -X GET https://<workernodefqdn>:<nodeport>/v2/_catalog
{"repositories":[]}
```

#### Ingress Controller

This second example will expose the DTR cache through an **ingress** object. In 
this example you will need to create a DNS rule in your environment that will 
resolve a DTR cache external FQDN address to the address of your ingress 
controller. You should have also specified the same DTR cache external FQDN
address within the DTR cache certificate in [step 1](#create-the-dtr-cache-certificates).

> Note an ingress controller is a prerequisite for this example. If you have not 
> deployed an ingress controller on your cluster, see [Layer 7 Routing for UCP](/ee/ucp/kubernetes/layer-7-routing). This 
> ingress controller will also need to support SSL passthrough.

```
cat > dtrcacheservice.yaml <<EOF
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
    - <external-dtr-cache-fqdn> # Replace this value with your external DTR Cache address
  rules:
  - host: <external-dtr-cache-fqdn> # Replace this value with your external DTR Cache address
    http:
      paths:
      - backend:
          serviceName: dtr-cache
          servicePort: 443
EOF

kubectl create -f dtrcacheservice.yaml
```

You can test that your DTR cache is externally reachable by using curl to hit 
the API endpoint. The address should be the one you have defined above in the 
serivce definition file.

```
curl -X GET https://external-dtr-cache-fqdn/v2/_catalog
{"repositories":[]}
```

## Next Steps 

[Integrate your cache into DTR and configure users](simple#register-the-cache-with-dtr)
