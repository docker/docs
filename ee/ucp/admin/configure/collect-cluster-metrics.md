---
description: Collecting UCP cluster metrics with Prometheus
keywords: prometheus, metrics, ucp
title: Collect UCP cluster metrics with Prometheus
redirect_from:
- /engine/admin/prometheus/
---

> Beta disclaimer
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

[Prometheus](https://prometheus.io/) is an open-source systems monitoring and
alerting toolkit. You can configure Docker as a Prometheus target. This topic
shows you how to configure Docker, set up Prometheus to run as a Docker
container, and monitor your Docker instance using Prometheus.

In UCP 3.0, Prometheus servers were standard containers. In UCP 3.1, Prometheus runs as a Kubernetes deployment. By default, this will be a daemonset that runs on every manager node. One benefit of this change is you can set the daemonset to not schedule on any nodes, which effectively disables Prometheus if you don’t use the UCP web interface.

The data is stored locally on disk for each Prometheus server, so data is not replicated on new managers or if you schedule Prometheus to run on a new node. Metrics are not kept longer than 24 hours.

> **Warning**: Upgrading UCP from 3.0.x to 3.1.x causes loss of metrics data.

## Deploy Prometheus on worker nodes

To deploy Prometheus on worker nodes in a cluster:

1. Begin by sourcing an admin bundle.

2. Verify that ucp-metrics pods are running on all managers.

```
$ kubectl -n kube-system get pods -l k8s-app=ucp-metrics -o wide
NAME                READY     STATUS    RESTARTS   AGE       IP              NODE
ucp-metrics-hvkr7   3/3       Running   0          4h        192.168.80.66   3a724a-0
```

3. Add a Kubernetes node label to one or more workers.  Here we add a label with key "ucp-metrics" and value "".

```
$ kubectl label node 3a724a-1 ucp-metrics=
node "test-3a724a-1" labeled
```

4. Patch the ucp-metrics DaemonSet's nodeSelector using the same key and value used for the node label. Here we again use the key “ucp-metrics” and the value “”.

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

## Configure external Prometheus to scrape Prometheus metrics from UCP

To configure your external Prometheus server to scrape metrics from Prometheus in UCP:

1. Begin by sourcing an admin bundle.

2. Create a Kubernetes secret containing your bundle’s TLS material.

```
(cd $DOCKER_CERT_PATH && kubectl create secret generic prometheus --from-file=ca.pem --from-file=cert.pem --from-file=key.pem)
```

3. Apply some YAML to create a Prometheus deployment and ClusterIP service. On AWS with Kube’s cloud provider configured, you could replace `ClusterIP` with `LoadBalancer` in the service YAML and then access the service through the load balancer. If you are running Prometheus external to UCP, change the following domain for inventory container in the Prometheus deployment from ucp-controller.kube-system.svc.cluster.local to a external domain for UCP accessible from the Prometheus node.

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
