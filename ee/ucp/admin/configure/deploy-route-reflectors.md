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

## Choose dedicated notes

Start by tainting the nodes, so that no other workload runs there. Configure
your CLI with a UCP client bundle, and for each dedicated node, run:

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
you'll need to SSH into a UCP node and download the `calicoctl` tool.

Log in to a UCP node using SSH, and run:

```
sudo curl --location https://github.com/projectcalico/calicoctl/releases/download/v3.1.1/calicoctl \
  --output /usr/bin/calicoctl
sudo chmod +x /usr/bin/calicoctl
```

Now you need to configure `calicoctl` to communicate with the etcd key-value
store managed by UCP. Create a file named `/etc/calico/calicoctl.cfg` with
the following content:

```
apiVersion: projectcalico.org/v3
kind: CalicoAPIConfig
metadata:
spec:
  datastoreType: "etcdv3"
  etcdEndpoints: "127.0.0.1:12378"
  etcdKeyFile: "/var/lib/docker/volumes/ucp-node-certs/_data/key.pem"
  etcdCertFile: "/var/lib/docker/volumes/ucp-node-certs/_data/cert.pem"
  etcdCACertFile: "/var/lib/docker/volumes/ucp-node-certs/_data/ca.pem"
```

## Disable node-to-node BGP mesh

Not that you've configured `calicoctl`, you can check the current Calico BGP
configuration:

```
sudo calicoctl get bgpconfig
```

If you don't see any configuration listed, create one by running:

```
cat << EOF | sudo calicoctl create -f -
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
sudo calicoctl get bgpconfig --output yaml > bgp.yaml
```

Edit the `bgp.yaml` file, updating `nodeToNodeMeshEnabled` to `false`. Then
update Calico configuration by running:

```
sudo calicoctl replace -f bgp.yaml
```

## Configure Calico to use Route Reflectors

To configure Calico to use the Route Reflectors you need to know the AS number
for your network first. For that, run:

```
sudo calicoctl get nodes --output=wide
```

Now that you have the AS number, you can create the Calico configuration.
For each Route Reflector, customize and run the following snippet:

```
sudo calicoctl create -f - << EOF
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
kubectl get pods -n kube-system -o wide | grep <node-name>

# Delete the Pod
kubectl delete pod -n kube-system <pod-name>
```

## Validate peers

Now you can check that other `calico-node` pods running on other nodes are
peering with the Route Reflector:

```
sudo calicoctl node status
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
