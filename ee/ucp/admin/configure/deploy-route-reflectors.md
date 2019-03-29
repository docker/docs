---
title: Improve network performance with Route Reflectors
description: Learn how to deploy Calico Route Reflectors to improve performance
  of Kubernetes networking
keywords: cluster, node, label, certificate, SAN
---

UCP uses Calico as the default Kubernetes networking solution. Calico is
configured to create a BGP mesh between all nodes in the cluster.

As you add more nodes to the cluster, networking performance starts decreasing.
If your cluster has more than 100 nodes, you should reconfigure Calico to use
Route Reflectors instead of a node-to-node mesh.

This article guides you in deploying Calico Route Reflectors in a UCP cluster.
UCP running on Microsoft Azure uses Azure SDN instead of Calico for
multi-host networking.
If your UCP deployment is running on Azure, you don't need to configure it this
way.

## Before you begin

For production-grade systems, you should deploy at least two Route Reflectors,
each running on a dedicated node. These nodes should not be running any other
workloads.

If Route Reflectors are running on a same node as other workloads, swarm ingress
and NodePorts might not work in these workloads.

## Choose dedicated nodes

Start by tainting the nodes, so that no other workload runs there. [Configure
your CLI with a UCP client bundle](/ee/ucp/user-access/cli/), and for each dedicated node, run:

```
kubectl taint node <node-name> \
  com.docker.ucp.kubernetes.calico/route-reflector=true:NoSchedule
```

Then add labels to those nodes, so that you can target them when deploying the
Route Reflectors. For each dedicated node, run:

```
kubectl label nodes <node-name> \
  com.docker.ucp.kubernetes.calico/route-reflector=true
```

## Deploy the Route Reflectors

Create a `calico-rr.yaml` file with the following content:

```
kind: DaemonSet
apiVersion: extensions/v1beta1
metadata:
  name: calico-rr
  namespace: kube-system
  labels:
    app: calico-rr
spec:
  updateStrategy:
    type: RollingUpdate
  selector:
    matchLabels:
      k8s-app: calico-rr
  template:
    metadata:
      labels:
        k8s-app: calico-rr
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      tolerations:
        - key: com.docker.ucp.kubernetes.calico/route-reflector
          value: "true"
          effect: NoSchedule
      hostNetwork: true
      containers:
        - name: calico-rr
          image: calico/routereflector:v0.6.1
          env:
            - name: ETCD_ENDPOINTS
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: etcd_endpoints
            - name: ETCD_CA_CERT_FILE
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: etcd_ca
            # Location of the client key for etcd.
            - name: ETCD_KEY_FILE
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: etcd_key # Location of the client certificate for etcd.
            - name: ETCD_CERT_FILE
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: etcd_cert
            - name: IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - mountPath: /calico-secrets
              name: etcd-certs
          securityContext:
            privileged: true
      nodeSelector:
        com.docker.ucp.kubernetes.calico/route-reflector: "true"
      volumes:
      # Mount in the etcd TLS secrets.
        - name: etcd-certs
          secret:
            secretName: calico-etcd-secrets
  ```

Then, deploy the DaemonSet using:

```
kubectl create -f calico-rr.yaml
```

## Configure calicoctl

To reconfigure Calico to use Route Reflectors instead of a node-to-node mesh,
you'll need to tell `calicoctl` where to find the etcd key-value store managed
by UCP.  From a CLI with a UCP client bundle, create a shell alias to start
`calicoctl` using the `{{ page.ucp_org }}/ucp-dsinfo` image:

```
UCP_VERSION=$(docker version --format {% raw %}'{{index (split .Server.Version "/") 1}}'{% endraw %})
alias calicoctl="\
docker run -i --rm \
  --pid host \
  --net host \
  -e constraint:ostype==linux \
  -e ETCD_ENDPOINTS=127.0.0.1:12378 \
  -e ETCD_KEY_FILE=/ucp-node-certs/key.pem \
  -e ETCD_CA_CERT_FILE=/ucp-node-certs/ca.pem \
  -e ETCD_CERT_FILE=/ucp-node-certs/cert.pem \
  -v /var/run/calico:/var/run/calico \
  -v ucp-node-certs:/ucp-node-certs:ro \
  {{ page.ucp_org }}/ucp-dsinfo:${UCP_VERSION} \
  calicoctl \
"
```

## Disable node-to-node BGP mesh

Now that you've configured `calicoctl`, you can check the current Calico BGP
configuration:

```
calicoctl get bgpconfig
```

If you don't see any configuration listed, create one by running:

```
calicoctl create -f - <<EOF
apiVersion: projectcalico.org/v3
kind: BGPConfiguration
metadata:
  name: default
spec:
  logSeverityScreen: Info
  nodeToNodeMeshEnabled: false
  asNumber: 63400
EOF
```

This creates a new configuration with node-to-node mesh BGP disabled.
If you have a configuration, and `meshenabled` is set to `true`, update your
configuration:

```
calicoctl get bgpconfig --output yaml > bgp.yaml
```

Edit the `bgp.yaml` file, updating `nodeToNodeMeshEnabled` to `false`. Then
update Calico configuration by running:

```
calicoctl replace -f - < bgp.yaml
```

## Configure Calico to use Route Reflectors

To configure Calico to use the Route Reflectors you need to know the AS number
for your network first. For that, run:

```
calicoctl get nodes --output=wide
```

Now that you have the AS number, you can create the Calico configuration.
For each Route Reflector, customize and run the following snippet:

```
calicoctl create -f - << EOF
apiVersion: projectcalico.org/v3
kind: BGPPeer
metadata:
  name: bgppeer-global
spec:
  peerIP: <IP_RR>
  asNumber: <AS_NUMBER>
EOF
```

Where:
* `IP_RR` is the IP of the node where the Route Reflector pod is deployed.
* `AS_NUMBER` is the same `AS number` for your nodes.

You can learn more about this configuration in the
[Calico documentation](https://docs.projectcalico.org/v3.1/usage/routereflector/calico-routereflector).

## Stop calico-node pods

If you have `calico-node` pods running on the nodes dedicated for running the
Route Reflector, manually delete them. This ensures that you don't have them
both running on the same node.

Using your UCP client bundle, run:

```
# Find the Pod name
kubectl -n kube-system \
  get pods --selector k8s-app=calico-node -o wide | \
  grep <node-name>

# Delete the Pod
kubectl -n kube-system delete pod <pod-name>
```

## Validate peers

Now you can check that `calico-node` pods running on other nodes are peering
with the Route Reflector.  From a CLI with a UCP client bundle, use a Swarm affinity filter to run `calicoctl node
status` on any node running `calico-node`:

```
UCP_VERSION=$(docker version --format {% raw %}'{{index (split .Server.Version "/") 1}}'{% endraw %})
docker run -i --rm \
  --pid host \
  --net host \
  -e affinity:container=='k8s_calico-node.*' \
  -e ETCD_ENDPOINTS=127.0.0.1:12378 \
  -e ETCD_KEY_FILE=/ucp-node-certs/key.pem \
  -e ETCD_CA_CERT_FILE=/ucp-node-certs/ca.pem \
  -e ETCD_CERT_FILE=/ucp-node-certs/cert.pem \
  -v /var/run/calico:/var/run/calico \
  -v ucp-node-certs:/ucp-node-certs:ro \
  {{ page.ucp_org }}/ucp-dsinfo:${UCP_VERSION} \
  calicoctl node status
```

You should see something like:

```
IPv4 BGP status
+--------------+-----------+-------+----------+-------------+
| PEER ADDRESS | PEER TYPE | STATE |  SINCE   |    INFO     |
+--------------+-----------+-------+----------+-------------+
| 172.31.24.86 | global    | up    | 23:10:04 | Established |
+--------------+-----------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.
```
