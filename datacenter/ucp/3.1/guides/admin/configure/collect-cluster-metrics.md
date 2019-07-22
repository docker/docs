---
description: Collecting UCP cluster metrics with Prometheus
keywords: prometheus, metrics, ucp
title: Collect UCP cluster metrics with Prometheus
redirect_from:
- /engine/admin/prometheus/
---

[Prometheus](https://prometheus.io/) is an open-source systems monitoring and
alerting toolkit. You can configure Docker as a Prometheus target. This topic
shows you how to configure Docker, set up Prometheus to run as a Docker
container, and monitor your Docker instance using Prometheus.

In UCP 3.0, Prometheus servers were standard containers. In UCP 3.1, Prometheus runs as a Kubernetes deployment. By default, this will be a DaemonSet that runs on every manager node. One benefit of this change is you can set the DaemonSet to not schedule on any nodes, which effectively disables Prometheus if you don’t use the UCP web interface.

The data is stored locally on disk for each Prometheus server, so data is not replicated on new managers or if you schedule Prometheus to run on a new node. Metrics are not kept longer than 24 hours.

Events, logs, and metrics are sources of data that provide observability of your cluster. Metrics monitors numerical data values that have a time-series component. There are several sources from which metrics can be derived, each providing different kinds of meaning for a business and its applications.

The Docker EE platform provides a base set of metrics that gets you running and into production without having to rely on external or 3rd party tools. Docker strongly encourages the use of additional monitoring to provide more comprehensive visibility into your specific Docker environment, but recognizes the need for a basic set of metrics built into the product. The following are examples of these metrics:

## Business metrics ##

These are high-level aggregate metrics that typically combine technical, financial, and organizational data to create metrics for business leaders of the IT infrastructure. Some examples of business metrics might be:
  - Company or division-level application downtime    
  - Aggregate resource utilization    
  - Application resource demand growth

## Application metrics ##

These are metrics about domain of APM tools like AppDynamics or DynaTrace and provide metrics about the state or performance of the application itself.
  - Service state metrics
  - Container platform metrics
  - Host infrastructure metrics

Docker EE 2.1 does not collect or expose application level metrics. 

The following are metrics Docker EE 2.1 collects, aggregates, and exposes:

## Service state metrics ##

These are metrics about the state of services running on the container platform. These types of metrics have very low cardinality, meaning the values are typically from a small fixed set of possibilities, commonly binary.
  - Application health
  - Convergence of K8s deployments and Swarm services
  - Cluster load by number of services or containers or pods

Web UI disk usage metrics, including free space, only reflect the Docker managed portion of the filesystem: `/var/lib/docker`. To monitor the total space available on each filesystem of a UCP worker or manager, you must deploy a third party monitoring solution to monitor the operating system.

## Deploy Prometheus on worker nodes

Universal Control Plane deploys Prometheus by default on the manager nodes to provide a built-in metrics backend. For cluster sizes over 100 nodes or for use cases where scraping metrics from the Prometheus instances are needed, we recommend that you deploy Prometheus on dedicated worker nodes in the cluster.

To deploy Prometheus on worker nodes in a cluster:

1. Begin by sourcing an admin bundle.

2. Verify that ucp-metrics pods are running on all managers.

    ```
    $ kubectl -n kube-system get pods -l k8s-app=ucp-metrics -o wide
    NAME                READY     STATUS    RESTARTS   AGE       IP              NODE
    ucp-metrics-hvkr7   3/3       Running   0          4h        192.168.80.66   3a724a-0
    ```

3. Add a Kubernetes node label to one or more workers.  Here we add a label with key "ucp-metrics" and value "" to a node with name "3a724a-1".

    ```
    $ kubectl label node 3a724a-1 ucp-metrics=
    node "test-3a724a-1" labeled
    ```
    
     > **SELinux Prometheus Deployment for UCP 3.1.0, 3.1.1, and 3.1.2**
     >
     > If you are using SELinux, you must label your `ucp-node-certs` directories properly on your worker nodes before you move the ucp-metrics workload to them. To run ucp-metrics on a worker node, update the `ucp-node-certs` label by running `sudo chcon -R system_u:object_r:container_file_t:s0 /var/lib/docker/volumes/ucp-node-certs/_data`.

4. Patch the ucp-metrics DaemonSet's nodeSelector using the same key and value used for the node label. This example shows the key "ucp-metrics" and the value "".

    ```
    $ kubectl -n kube-system patch daemonset ucp-metrics --type json -p '[{"op": "replace", "path": "/spec/template/spec/nodeSelector", "value": {"ucp-metrics": ""}}]'
    daemonset "ucp-metrics" patched
    ```

5. Observe that ucp-metrics pods are running only on the labeled workers.

    ```
    $ kubectl -n kube-system get pods -l k8s-app=ucp-metrics -o wide
    NAME                READY     STATUS        RESTARTS   AGE       IP              NODE
    ucp-metrics-88lzx   3/3       Running       0          12s       192.168.83.1    3a724a-1
    ucp-metrics-hvkr7   3/3       Terminating   0          4h        192.168.80.66   3a724a-0
    ```

## Configure external Prometheus to scrape metrics from UCP

To configure your external Prometheus server to scrape metrics from Prometheus in UCP:

1. Begin by sourcing an admin bundle.

2. Create a Kubernetes secret containing your bundle’s TLS material.

    ```
    (cd $DOCKER_CERT_PATH && kubectl create secret generic prometheus --from-file=ca.pem --from-file=cert.pem --from-file=key.pem)
    ```

3. Create a Prometheus deployment and ClusterIP service using YAML as follows.

   On AWS with Kube’s cloud provider configured, you can replace `ClusterIP` with `LoadBalancer` in the service YAML then access the service through the load balancer. If running Prometheus external to UCP, change the following domain for the inventory container in the Prometheus deployment from `ucp-controller.kube-system.svc.cluster.local` to an external domain to access UCP from the Prometheus node.

    ```
    kubectl apply -f - <<EOF
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: prometheus
    data:
      prometheus.yaml: |
        global:
          scrape_interval: 10s
        scrape_configs:
        - job_name: 'ucp'
          tls_config:
            ca_file: /bundle/ca.pem
            cert_file: /bundle/cert.pem
            key_file: /bundle/key.pem
            server_name: proxy.local
          scheme: https
          file_sd_configs:
          - files:
            - /inventory/inventory.json
    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: prometheus
    spec:
      replicas: 2
      selector:
        matchLabels:
          app: prometheus
      template:
        metadata:
          labels:
            app: prometheus
        spec:
          containers:
          - name: inventory
            image: alpine
            command: ["sh", "-c"]
            args:
            - apk add --no-cache curl &&
              while :; do
                curl -Ss --cacert /bundle/ca.pem --cert /bundle/cert.pem --key /bundle/key.pem --output /inventory/inventory.json https://ucp-controller.kube-system.svc.cluster.local/metricsdiscovery;
                sleep 15;
              done
            volumeMounts:
            - name: bundle
              mountPath: /bundle
            - name: inventory
              mountPath: /inventory
          - name: prometheus
            image: prom/prometheus
            command: ["/bin/prometheus"]
            args:
            - --config.file=/config/prometheus.yaml
            - --storage.tsdb.path=/prometheus
            - --web.console.libraries=/etc/prometheus/console_libraries
            - --web.console.templates=/etc/prometheus/consoles
            volumeMounts:
            - name: bundle
              mountPath: /bundle
            - name: config
              mountPath: /config
            - name: inventory
              mountPath: /inventory
          volumes:
          - name: bundle
            secret:
              secretName: prometheus
          - name: config
            configMap:
              name: prometheus
          - name: inventory
            emptyDir:
              medium: Memory
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: prometheus
    spec:
      ports:
      - port: 9090
        targetPort: 9090
      selector:
        app: prometheus
      sessionAffinity: ClientIP
    EOF
    ```

4. Determine the service ClusterIP.

    ```
    $ kubectl get service prometheus
    NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
    prometheus   ClusterIP   10.96.254.107   <none>        9090/TCP   1h
    ```

5. Forward port 9090 on the local host to the ClusterIP. The tunnel created does not need to be kept alive and is only intended to expose the Prometheus UI.

    ```
    ssh -L 9090:10.96.254.107:9090 ANY_NODE
    ```

6. Visit `http://127.0.0.1:9090` to explore the UCP metrics being collected by Prometheus.
