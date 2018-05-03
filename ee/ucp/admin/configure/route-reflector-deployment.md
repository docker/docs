# Steps to deploy Route Reflector on UCP 3.0

UCP 3.0 ships with Calico for Kubernetes networking. Calico use BGP as the control plane to distribute Pod routes between the node in the UCP cluster. Calcio BGP is setup as mesh between the nodes. For large scale production-grade deployments it is recommended to deploy Route Reflectors to avoid node to node BGP mesh. This document provides detailed steps to deploy Route Reflectors in UCP cluster.

Note: 
1) If you are using UCP 3.0 on Azure deploying route reflectors is not required. The networking control plane is handled
by Azure SDN and not calico.

2) Nodes marked for RR should not run any non system workloads. These nodes will not be gauranteed to work for ingress (swarm) or nodePort(kubernetes)

3) It is recommended to deploy Route Reflectors if the cluster scale exceeds 100 nodes.

4) It is recommended to deploy atleast 2 Route Reflectors in the cluster. Identify the nodes where the route reflector needs to be deployed considering that these nodes will not be available to the scheduler to schedule any workloads. Choose the orchestrator as Kubernetes for these nodes.


1)  On one of the nodes in Docker EE cluster, download the calicoctl binary from 
    ```
    curl --output /usr/bin/calicoctl -O -L https://github.com/projectcalico/calicoctl/releases/download/v3.1.1/calicoctl --output /usr/bin/calicoctl
    chmod +x /usr/bin/calicoctl
    ```
    
2) `calicoctl` will need to be configured with the etcd information. There are 2 options:

    a) You can set the environment variables
    
       ```
       export ETCD_ENDPOINTS=127.0.0.1:12379 
       export ETCD_KEY_FILE=/var/lib/docker/volumes/ucp-node-certs/_data/key.pem 
       export ETCD_CA_CERT_FILE=/var/lib/docker/volumes/ucp-node-certs/_data/ca.pem 
       export ETCD_CERT_FILE=/var/lib/docker/volumes/ucp-node-certs/_data/cert.pem 
       ```
       
    b) You can create a configuration file.
       `calicoctl` by default looks at `/etc/calico/calicoctl.cfg` for the configuration. For custom 
       configuration file location, use `--config` with `calicoctl`
       
       `calicoctl.cfg:`
       ```
       apiVersion: projectcalico.org/v3
       kind: CalicoAPIConfig
       metadata:
       spec:
         datastoreType: "etcdv3"
         etcdEndpoints: "127.0.0.1:12379"
         etcdKeyFile: "/var/lib/docker/volumes/ucp-node-certs/_data/key.pem"
         etcdCertFile: "/var/lib/docker/volumes/ucp-node-certs/_data/cert.pem"
         etcdCACertFile: "/var/lib/docker/volumes/ucp-node-certs/_data/ca.pem"
       ```
   
3) Identify the nodes where the route reflector needs to be deployed. Taint the nodes to ensure only the Route Reflector        pods that tolerate the taint can be scheduled on these nodes.

   In this example node `ubuntu-0` is being tainted. Use the client bundle and run:
   
   ```kubectl taint node ubuntu-0 com.docker.ucp.kubernetes.calico/route-reflector=true:NoSchedule```
   
   This will need to be done on every node where the calico route reflector pod will need to be deployed.

4) Add labels to the same node. The labels will be used to for node placement by the scheduler to deploy the calico route      reflector pods.
   Use the client bundle and run:

   ```kubectl label nodes ubuntu-0 com.docker.ucp.kubernetes.calico/route-reflector=true ```

   This will need to be done on every node where the calico route reflector pod will need to be deployed.

5) Deploy the calico route reflector Daemonset. The calico-rr pods will be deployed on all nodes that have the taint
   and the node labels.
   
   `calico-rr.yaml`
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
  kubectl create -f calico-rr.yaml
  ```
  
6) Use `calicoctl` to get the current bgp configuration and modify it to turn off node to node bgp mesh.

   ```calicoctl get bgpconfig -o yaml > bgp.yaml```
   
   Set the nodeToNodeMeshEnabled to false in the spec.
  
   ` bgp.yaml:`
   ```
   ....
   spec:
     asNumber: 63400
     logSeverityScreen: Info
     nodeToNodeMeshEnabled: false
   ....
   ```
   Replace the bgp configuation with the modified `bgp.yaml`
   
   ```
   calicoctl replace -f bgp.yaml
   ```
   
   If the bgpConfig object is empty, create the bgp config object.
   ```
   cat << EOF | calicoctl create -f -
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
   
7) Create a bgp RR peer configuration. You can create configurations for each RR

   ```
   calicoctl create -f - << EOF
   apiVersion: projectcalico.org/v3
   kind: BGPPeer
   metadata:
     name: bgppeer-global
   spec:
     peerIP: <IP_RR> #This will be the IP of node where the calico route reflector pod will be deployed.
     asNumber: <AS_NUM> #Use the same asNumber from the configuration above step 6).
   EOF
   ```
   The bgp RR configuration needs to be added for every Route Reflector that will be deployed.
   Refer to https://docs.projectcalico.org/v3.1/usage/routereflector/calico-routereflector for more information.


8) If you have calico-node pod running on the nodes marked for route reflector in step 3). You will need to manually delete the calico-node pod. We recommend this to avoid running both calico route reflector and calico-node pods on the same node.
   
   ```
   kubectl get pods -n kube-system -o wide | grep ubuntu-0
   calico-node-t4lwt                        2/2       Running   0          3h        172.31.20.89     ubuntu-0
   
   kubectl delete pod -n kube-system calico-node-t4lwt
   ```
  
9) You can verify the calico-node from other nodes peering with the route reflector are status by downloading calicoctl        bundle on other nodes and checking for `calicoctl node status`
   ```
   sudo ./calicoctl node status
   Calico process is running.

   IPv4 BGP status
   +--------------+-----------+-------+----------+-------------+
   | PEER ADDRESS | PEER TYPE | STATE |  SINCE   |    INFO     |
   +--------------+-----------+-------+----------+-------------+
   | 172.31.24.86 | global    | up    | 23:10:04 | Established |
   +--------------+-----------+-------+----------+-------------+

   IPv6 BGP status
   No IPv6 peers found.
   ```

