---
title: Improve network performance with route reflectors
description: Learn how to deploy Calico route reflectors to improve the performance of Kubernetes networking.
keywords: cluster, node, label, certificate, SAN
---

>{% include enterprise_label_shortform.md %}

Docker Universal Control Plane (UCP) uses [Calico](https://docs.projectcalico.org/v3.11/introduction/) as the default Kubernetes networking solution. Calico is configured to create a Border Gateway Protocol (BGP) mesh between all nodes in the cluster.

As you add more nodes to the cluster, networking performance starts decreasing. In fact, if your cluster has more than 100 nodes, you should reconfigure Calico to use route reflectors instead of a node-to-node mesh.

Route reflectors are useful in large scale deployments to reduce the number of BGP connections needed for correct and complete route propagation.

> **Note**
>
> UCP running on Microsoft Azure uses Azure Software Defined Networking (SDN) instead of Calico for multi-host networking. If your UCP deployment is running on Azure, you don't need to configure it this way. See [Install UCP on Azure](https://docs.docker.com/ee/ucp/admin/install/cloudproviders/install-on-azure/) for more information.

## Before you begin

For production-grade systems, deploy at least two route reflectors, each running on its own node. These nodes should be dedicated solely to the purpose, too, because Swarm ingress and NodePorts may not function on any other workloads running on a route reflector node.

## Choose dedicated nodes

1. [Taint the nodes](https://docs.docker.com/ee/ucp/admin/configure/restrict-services-to-worker-nodes/) to ensure that they are unable to run other workloads. To do this, configure the CLI with a [UCP client bundle](/ee/ucp/user-access/cli/).

2. For each dedicated node:

    ```
    kubectl taint node <node-name> \
    com.docker.ucp.kubernetes.calico/route-reflector=true:NoSchedule
    ```
3. Add labels to each of the dedicated nodes:

    ```
    kubectl label nodes <node-name> \
    com.docker.ucp.kubernetes.calico/route-reflector=true
    ```

## Deploy the route reflectors

1. Create a `calico-rr.yaml` file with the following content:

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

2. Deploy the DaemonSet:

    ```
    kubectl create -f calico-rr.yaml
    ```

## Configure calicoctl

To reconfigure Calico to use route reflectors instead of a node-to-node mesh, `calicoctl` needs to be able to locate the etcd key-value store managed by UCP. From a CLI configured with a UCP client bundle, create a shell alias to start `calicoctl` using the `{{ page.ucp_org }}/ucp-dsinfo` image:

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

After configuring `calicoctl`, check the current Calico BGP configuration:

```
calicoctl get bgpconfig
```

If you don't see any configuration listed, create one:

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

This action creates a new configuration with node-to-node mesh BGP disabled.

If you have a configuration, and `meshenabled` is set to `true`:

1. Update your configuration:

    ```
    calicoctl get bgpconfig --output yaml > bgp.yaml
    ```

2. Edit the `bgp.yaml` file, changing `nodeToNodeMeshEnabled` to `false`. 
3. Update the Calico configuration:

    ```
    calicoctl replace -f - < bgp.yaml
    ```

## Configure Calico to use route reflectors

To configure Calico to use route reflectors, you first need to get the Autonomous System (AS) number for the network. To get the AS number for your network:

```
calicoctl get nodes --output=wide
```

Using the AS number, create the Calico configuration by customizing and running the following snipped for each route reflector:

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
* `IP_RR` is the IP address of the node on which the route reflector pod is deployed.
* `AS_NUMBER` is the same `AS number` for your nodes.

You can learn more about this configuration in the
[Calico documentation](https://docs.projectcalico.org/v3.1/usage/routereflector/calico-routereflector).

## Stop calico-node pods

1. Manually delete any `calico-node` pods that are running on nodes dedicated to the running of route reflectors, as this will ensure that there are no instances in which pods and route reflectors are running on the same node.

2. Using your UCP client bundle:

```
# Find the Pod name
kubectl -n kube-system \
  get pods --selector k8s-app=calico-node -o wide | \
  grep <node-name>

# Delete the Pod
kubectl -n kube-system delete pod <pod-name>
```

## Validate peers

1. Verify that the `calico-node` pods running on other nodes are peering with the route reflector. 
2. From a CLI configured with a UCP client bundle, use a Swarm affinity filter to run `calicoctl node status` on any node running `calico-node`:

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

The following is a sample output.

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
## Where to go next
* [Monitor the cluster status](https://docs.docker.com/ee/ucp/admin/monitor-and-troubleshoot/)
