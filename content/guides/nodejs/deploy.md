---
title: Deploy your Node.js application
linkTitle: Deploy your app
weight: 50
keywords: deploy, kubernetes, node, node.js, production
description: Learn how to deploy your containerized Node.js application to Kubernetes with production-ready configuration
aliases:
  - /language/nodejs/deploy/
  - /guides/language/nodejs/deploy/
---

## Prerequisites

- Complete all the previous sections of this guide, starting with [Containerize a Node.js application](containerize.md).
- [Turn on Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

## Overview

In this section, you'll learn how to deploy your containerized Node.js application to Kubernetes using Docker Desktop. This deployment uses production-ready configurations including security hardening, auto-scaling, persistent storage, and high availability features.

You'll deploy a complete stack including:

- Node.js Todo application with 3 replicas.
- PostgreSQL database with persistent storage.
- Auto-scaling based on CPU and memory usage.
- Ingress configuration for external access.
- Security settings.

## Create a Kubernetes deployment file

Create a new file called `nodejs-sample-kubernetes.yaml` in your project root:

```yaml
# ========================================
# Node.js Todo App - Kubernetes Deployment
# ========================================

apiVersion: v1
kind: Namespace
metadata:
  name: todoapp
  labels:
    app: todoapp

---
# ========================================
# ConfigMap for Application Configuration
# ========================================
apiVersion: v1
kind: ConfigMap
metadata:
  name: todoapp-config
  namespace: todoapp
data:
  NODE_ENV: 'production'
  ALLOWED_ORIGINS: 'https://yourdomain.com'
  POSTGRES_HOST: 'todoapp-postgres'
  POSTGRES_PORT: '5432'
  POSTGRES_DB: 'todoapp'
  POSTGRES_USER: 'todoapp'

---
# ========================================
# Secret for Database Credentials
# ========================================
apiVersion: v1
kind: Secret
metadata:
  name: todoapp-secrets
  namespace: todoapp
type: Opaque
data:
  postgres-password: dG9kb2FwcF9wYXNzd29yZA== # base64 encoded "todoapp_password"

---
# ========================================
# PostgreSQL PersistentVolumeClaim
# ========================================
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: todoapp
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: standard

---
# ========================================
# PostgreSQL Deployment
# ========================================
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todoapp-postgres
  namespace: todoapp
  labels:
    app: todoapp-postgres
spec:
  replicas: 1
  selector:
    matchLabels:
      app: todoapp-postgres
  template:
    metadata:
      labels:
        app: todoapp-postgres
    spec:
      containers:
        - name: postgres
          image: postgres:16-alpine
          ports:
            - containerPort: 5432
              name: postgres
          env:
            - name: POSTGRES_DB
              valueFrom:
                configMapKeyRef:
                  name: todoapp-config
                  key: POSTGRES_DB
            - name: POSTGRES_USER
              valueFrom:
                configMapKeyRef:
                  name: todoapp-config
                  key: POSTGRES_USER
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: todoapp-secrets
                  key: postgres-password
          volumeMounts:
            - name: postgres-storage
              mountPath: /var/lib/postgresql/data
          livenessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - todoapp
                - -d
                - todoapp
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - todoapp
                - -d
                - todoapp
            initialDelaySeconds: 5
            periodSeconds: 5
      volumes:
        - name: postgres-storage
          persistentVolumeClaim:
            claimName: postgres-pvc

---
# ========================================
# PostgreSQL Service
# ========================================
apiVersion: v1
kind: Service
metadata:
  name: todoapp-postgres
  namespace: todoapp
  labels:
    app: todoapp-postgres
spec:
  type: ClusterIP
  ports:
    - port: 5432
      targetPort: 5432
      protocol: TCP
      name: postgres
  selector:
    app: todoapp-postgres

---
# ========================================
# Application Deployment
# ========================================
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todoapp-deployment
  namespace: todoapp
  labels:
    app: todoapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: todoapp
  template:
    metadata:
      labels:
        app: todoapp
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      containers:
        - name: todoapp
          image: ghcr.io/your-username/docker-nodejs-sample:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 3000
              name: http
              protocol: TCP
          env:
            - name: NODE_ENV
              valueFrom:
                configMapKeyRef:
                  name: todoapp-config
                  key: NODE_ENV
            - name: ALLOWED_ORIGINS
              valueFrom:
                configMapKeyRef:
                  name: todoapp-config
                  key: ALLOWED_ORIGINS
            - name: POSTGRES_HOST
              valueFrom:
                configMapKeyRef:
                  name: todoapp-config
                  key: POSTGRES_HOST
            - name: POSTGRES_PORT
              valueFrom:
                configMapKeyRef:
                  name: todoapp-config
                  key: POSTGRES_PORT
            - name: POSTGRES_DB
              valueFrom:
                configMapKeyRef:
                  name: todoapp-config
                  key: POSTGRES_DB
            - name: POSTGRES_USER
              valueFrom:
                configMapKeyRef:
                  name: todoapp-config
                  key: POSTGRES_USER
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: todoapp-secrets
                  key: postgres-password
          livenessProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 5
            periodSeconds: 5
          resources:
            requests:
              memory: '256Mi'
              cpu: '250m'
            limits:
              memory: '512Mi'
              cpu: '500m'
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop:
                - ALL

---
# ========================================
# Application Service
# ========================================
apiVersion: v1
kind: Service
metadata:
  name: todoapp-service
  namespace: todoapp
  labels:
    app: todoapp
spec:
  type: ClusterIP
  ports:
    - name: http
      port: 80
      targetPort: 3000
      protocol: TCP
  selector:
    app: todoapp

---
# ========================================
# Ingress for External Access
# ========================================
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: todoapp-ingress
  namespace: todoapp
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: 'letsencrypt-prod'
spec:
  tls:
    - hosts:
        - yourdomain.com
      secretName: todoapp-tls
  rules:
    - host: yourdomain.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: todoapp-service
                port:
                  number: 80

---
# ========================================
# HorizontalPodAutoscaler
# ========================================
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: todoapp-hpa
  namespace: todoapp
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: todoapp-deployment
  minReplicas: 1
  maxReplicas: 5
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80

---
# ========================================
# PodDisruptionBudget
# ========================================
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: todoapp-pdb
  namespace: todoapp
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: todoapp
```

## Configure the deployment

Before deploying, you need to customize the deployment file for your environment:

1. **Image reference**: Replace `your-username` with your GitHub username or Docker Hub username:

   ```yaml
   image: ghcr.io/your-username/docker-nodejs-sample:latest
   ```

2. **Domain name**: Replace `yourdomain.com` with your actual domain in two places:

   ```yaml
   # In ConfigMap
   ALLOWED_ORIGINS: "https://yourdomain.com"

   # In Ingress
   - host: yourdomain.com
   ```

3. **Database password** (optional): The default password is already base64 encoded. To change it:

   ```console
   $ echo -n "your-new-password" | base64
   ```

   Then update the Secret:

   ```yaml
   data:
     postgres-password: <your-base64-encoded-password>
   ```

4. **Storage class**: Adjust based on your cluster (current: `standard`)

## Understanding the deployment

The deployment file creates a complete application stack with multiple components working together.

### Architecture

The deployment includes:

- **Node.js application**: Runs 3 replicas of your containerized Todo app
- **PostgreSQL database**: Single instance with 10Gi of persistent storage
- **Services**: Kubernetes services handle load balancing across application replicas
- **Ingress**: External access through an ingress controller with SSL/TLS support

### Security

The deployment uses several security features:

- Containers run as a non-root user (UID 1001)
- Read-only root filesystem prevents unauthorized writes
- Linux capabilities are dropped to minimize attack surface
- Sensitive data like database passwords are stored in Kubernetes secrets

### High availability

To keep your application running reliably:

- Three application replicas ensure service continues if one pod fails
- Pod disruption budget maintains at least one available pod during updates
- Rolling updates allow zero-downtime deployments
- Health checks on the `/health` endpoint ensure only healthy pods receive traffic

### Auto-scaling

The Horizontal Pod Autoscaler scales your application based on resource usage:

- Scales between 1 and 5 replicas automatically
- Triggers scaling when CPU usage exceeds 70%
- Triggers scaling when memory usage exceeds 80%
- Resource limits: 256Mi-512Mi memory, 250m-500m CPU per pod

### Data persistence

PostgreSQL data is stored persistently:

- 10Gi persistent volume stores database files
- Database initializes automatically on first startup
- Data persists across pod restarts and updates

## Deploy your application

### Step 1: Deploy to Kubernetes

Deploy your application to the local Kubernetes cluster:

```console
$ kubectl apply -f nodejs-sample-kubernetes.yaml
```

You should see output confirming all resources were created:

```shell
namespace/todoapp created
secret/todoapp-secrets created
configmap/todoapp-config created
persistentvolumeclaim/postgres-pvc created
deployment.apps/todoapp-postgres created
service/todoapp-postgres created
deployment.apps/todoapp-deployment created
service/todoapp-service created
ingress.networking.k8s.io/todoapp-ingress created
poddisruptionbudget.policy/todoapp-pdb created
horizontalpodautoscaler.autoscaling/todoapp-hpa created
```

### Step 2: Verify the deployment

Check that your deployments are running:

```console
$ kubectl get deployments -n todoapp
```

Expected output:

```shell
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
todoapp-deployment   3/3     3            3           30s
todoapp-postgres     1/1     1            1           30s
```

Verify your services are created:

```console
$ kubectl get services -n todoapp
```

Expected output:

```shell
NAME               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
todoapp-service    ClusterIP   10.111.101.229   <none>        80/TCP     45s
todoapp-postgres   ClusterIP   10.111.102.130   <none>        5432/TCP   45s
```

Check that persistent storage is working:

```console
$ kubectl get pvc -n todoapp
```

Expected output:

```shell
NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
postgres-pvc   Bound    pvc-12345678-1234-1234-1234-123456789012   10Gi       RWO            standard       1m
```

### Step 3: Access your application

For local testing, use port forwarding to access your application:

```console
$ kubectl port-forward -n todoapp service/todoapp-service 8080:80
```

Open your browser and visit [http://localhost:8080](http://localhost:8080) to see your Todo application running in Kubernetes.

### Step 4: Test the deployment

Test that your application is working correctly:

1. **Add some todos** through the web interface
2. **Check application pods**:

   ```console
   $ kubectl get pods -n todoapp -l app=todoapp
   ```

3. **View application logs**:

   ```console
   $ kubectl logs -f deployment/todoapp-deployment -n todoapp
   ```

4. **Check database connectivity**:

   ```console
   $ kubectl get pods -n todoapp -l app=todoapp-postgres
   ```

5. **Monitor auto-scaling**:
   ```console
   $ kubectl describe hpa todoapp-hpa -n todoapp
   ```

### Step 5: Clean up

When you're done testing, remove the deployment:

```console
$ kubectl delete -f nodejs-sample-kubernetes.yaml
```

## Summary

You've deployed your containerized Node.js application to Kubernetes. You learned how to:

- Create a comprehensive Kubernetes deployment file with security hardening
- Deploy a multi-tier application (Node.js + PostgreSQL) with persistent storage
- Configure auto-scaling, health checks, and high availability features
- Test and monitor your deployment locally using Docker Desktop's Kubernetes

Your application is now running in a production-like environment with enterprise-grade features including security contexts, resource management, and automatic scaling.

---

## Related resources

Explore official references and best practices to sharpen your Kubernetes deployment workflow:

- [Kubernetes documentation](https://kubernetes.io/docs/home/) – Learn about core concepts, workloads, services, and more.
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md) – Use Docker Desktop's built-in Kubernetes support for local testing and development.
- [`kubectl` CLI reference](https://kubernetes.io/docs/reference/kubectl/) – Manage Kubernetes clusters from the command line.
- [Kubernetes Deployment resource](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) – Understand how to manage and scale applications using Deployments.
- [Kubernetes Service resource](https://kubernetes.io/docs/concepts/services-networking/service/) – Learn how to expose your application to internal and external traffic.
