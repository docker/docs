---
title: Use Pod Security Policies in UCP
description: Learn how to use Pod Security Policies to lock down Kubernetes as part of Universal Control Plane.
keywords: UCP, Kubernetes, psps, pod security policies
redirect_from:
---
Pod Security Policies (PSPs) are cluster-level resources which are enabled by default in Docker Universal Control Plane (UCP) 3.2. See [Pod Security Policy](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) for an explanation of this Kubernetes concept.


There are two default PSPs in UCP: a `privileged` policy
and an `unprivileged` policy. Administrators of the cluster can enforce additional
policies and apply them to users and teams for further control of what runs in the
Kubernetes cluster. This guide describes the two default policies, and
provides two example use cases for custom policies.

## Kubernetes Role Based Access Control (RBAC)

To interact with PSPs, a user will need to be granted access to
the `PodSecurityPolicy` object in Kubernetes RBAC. If the user is a `UCP Admin`,
or has been granted the `ClusterRole`: `cluster-admin` for all namespaces, then
the user can  already manipulate PSPs. Additionally, a normal
user can interact with policies if a UCP admin creates the following
`ClusterRole` and `ClusterRoleBinding`: 

```
$ cat <<EOF | kubectl create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: psp-admin
rules:
- apiGroups:
  - extensions
  resources:
  - podsecuritypolicies
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
EOF

$ USER=jeff

$ cat <<EOF | kubectl create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: psp-admin:$USER
roleRef:
  kind: ClusterRole
  name: psp-admin
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: User
  name: $USER
EOF
```

## Default pod security policies in UCP

By default, there are two policies defined within UCP, `privileged` and
`unprivileged`. Additionally, there is a `ClusterRoleBinding` that gives every
single user access to the privileged policy. This is for backward compatibility
after an upgrade. By default, any user can create any pod.

> Note: PSPs do not override security defaults built into the
> UCP RBAC engine for Kubernetes pods. These [Security
> defaults](https://docs.docker.com/ee/ucp/authorization/) prevent non-admin
> users from mounting host paths into pods or starting privileged pods.

```bash
$ kubectl get podsecuritypolicies
NAME           PRIV    CAPS   SELINUX    RUNASUSER   FSGROUP    SUPGROUP   READONLYROOTFS   VOLUMES
privileged     true    *      RunAsAny   RunAsAny    RunAsAny   RunAsAny   false            *
unprivileged   false          RunAsAny   RunAsAny    RunAsAny   RunAsAny   false            *
```

The specification for the `privileged` policy is as follows:

```
  allowPrivilegeEscalation: true
  allowedCapabilities:
  - '*'
  fsGroup:
    rule: RunAsAny
  hostIPC: true
  hostNetwork: true
  hostPID: true
  hostPorts:
  - max: 65535
    min: 0
  privileged: true
  runAsUser:
    rule: RunAsAny
  seLinux:
    rule: RunAsAny
  supplementalGroups:
    rule: RunAsAny
  volumes:
  - '*'
```

The specification for the `unprivileged` policy is as follows:

```
  allowPrivilegeEscalation: false
  allowedHostPaths:
  - pathPrefix: /dev/null
    readOnly: true
  fsGroup:
    rule: RunAsAny
  hostPorts:
  - max: 65535
    min: 0
  runAsUser:
    rule: RunAsAny
  seLinux:
    rule: RunAsAny
  supplementalGroups:
    rule: RunAsAny
  volumes:
  - '*'
```

## Use the unprivileged policy

> Note: When following this guide, if the prompt `$` follows `admin`, the action
> needs to be performed by a user with access to create pod security policies as
> discussed in [the Kubernetes RBAC section](#kubernetes-role-based-access-control). If the prompt `$`
> follows `user`, the UCP account does not need access to the PSP
> object in Kubernetes. The user only needs the ability to create Kubernetes pods.
> 

To switch users from the `privileged` policy to the `unprivileged` policy (or
any custom policy), an admin must first remove the `ClusterRoleBinding` that
links all users and service accounts to the `privileged` policy.

```
admin $ kubectl delete clusterrolebindings ucp:all:privileged-psp-role
```

When the `ClusterRoleBinding` is removed, cluster admins can still deploy pods,
and these pods are deployed with the `privileged` policy. But users or service
accounts are unable to deploy pods, because Kubernetes does not know what pod
security policy to apply.

```bash
user $ kubectl apply -f pod.yaml
Error from server (Forbidden): error when creating "pod.yaml": pods "demopod" is forbidden: unable to validate against any pod security policy: []
```

Therefore, to allow a user or a service account to use the `unprivileged` policy
(or any custom policy), you must create a `RoleBinding` to link that user or
team with the alternative policy. For the `unprivileged` policy, a `ClusterRole`
has already been defined, but has not been attached to a user. 

```bash
# List Existing Cluster Roles
admin $ kubectl get clusterrole | grep psp
privileged-psp-role                                                    3h47m
unprivileged-psp-role                                                  3h47m

# Define which user to apply the ClusterRole too
admin $ USER=jeff

# Create a RoleBinding linking the ClusterRole to the User 
admin $ cat <<EOF | kubectl create -f -
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: unprivileged-psp-role:$USER
  namespace: default
roleRef:
  kind: ClusterRole
  name: unprivileged-psp-role
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: User
  name: $USER
  namespace: default
EOF
```

In the following example, when user "jeff" deploys a basic `nginx` pod, the `unprivileged` policy 
then gets applied.

```bash
user $ cat <<EOF | user create -f -
apiVersion: v1
kind: Pod
metadata:
  name: demopod
spec:
  containers:
    - name:  demopod
      image: nginx
EOF

user $ kubectl get pods
NAME      READY   STATUS    RESTARTS   AGE
demopod   1/1     Running   0          10m
```

To check which PSP is applied to a pod, you can get a detailed
view of the pod spec using the `-o yaml` or `-o json` syntax with `kubectl`. You
can parse JSON output with [jq](https://stedolan.github.io/jq/). 

```bash
user $ kubectl get pods demopod -o json | jq '.metadata.annotations."kubernetes.io/psp"'
"unprivileged"
```

## Reenable the privileged PSP for all users

To revert to the default UCP configuration, in which all UCP users and service
accounts use the `privileged` PSP, recreate the default
`ClusterRoleBinding`:

```bash
admin $ cat <<EOF | user create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ucp:all:privileged-psp-role
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: privileged-psp-role
subjects:
- kind: Group
  name: system:authenticated
  apiGroup: rbac.authorization.k8s.io
- kind: Group
  name: system:serviceaccounts
  apiGroup: rbac.authorization.k8s.io
EOF
```

## PSP examples

UCP admins or users with the [correct
permissions](#kubernetes-role-based-access-control) can create their own custom
policies and attach them to UCP users or teams. This section highlights two
potential use cases for custom PSPs. These two uses cases can
be combined into the same policy. Note there are many more use cases with PSPs not covered 
in this document.

- Preventing containers that start as the Root User

- Applying default seccomp policies to all Kubernetes Pods.

For the full list of parameters that can be configured in a PSP,
see the [Kubernetes
documentation](https://kubernetes.io/docs/concepts/policy/pod-security-policy/).

### Example 1: Use a PSP to enforce "no root users"

A common use case for PSPs is to prevent a user from deploying
containers that run with the root user. A PSP can be created to
enforce this with the parameter `MustRunAsNonRoot`.

```bash
admin $ cat <<EOF | kubectl create -f -
apiVersion: extensions/v1beta1
kind: PodSecurityPolicy
metadata:  
  name: norootcontainers
spec:
  allowPrivilegeEscalation: false
  allowedHostPaths:
  - pathPrefix: /dev/null
    readOnly: true
  fsGroup:
    rule: RunAsAny
  hostPorts:
  - max: 65535
    min: 0
  runAsUser:
    rule: MustRunAsNonRoot
  seLinux:
    rule: RunAsAny
  supplementalGroups:
    rule: RunAsAny
  volumes:
  - '*'
EOF
```

If not done previously, the admin user must remove the `ClusterRoleBinding`
for the `privileged` policy, and then add a new `ClusterRole` and `RoleBinding`
to link a user to the new `norootcontainers` policy.

```bash
# Delete the default privileged ClusterRoleBinding
admin $ kubectl delete clusterrolebindings ucp:all:privileged-psp-role

# Create a ClusterRole Granting Access to the Policy
admin $ cat <<EOF | admin create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: norootcontainers-psp-role
rules:
- apiGroups:
  - policy
  resourceNames:
  - norootcontainers
  resources:
  - podsecuritypolicies
  verbs:
  - use
EOF

# Define a User to attach to the No Root Policy
admin $ USER=jeff

# Create a RoleBinding attaching the User to the ClusterRole
admin $ cat <<EOF | admin create -f -
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: norootcontainers-psp-role:$USER
  namespace: default
roleRef:
  kind: ClusterRole
  name: norootcontainers-psp-role
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: User
  name: $USER
  namespace: default
EOF
```

If a user tries to deploy a pod that runs as a root user, such as the upstream
`nginx` image, this should fail with a `ConfigError`.

```bash
user $ cat <<EOF | user create -f -
apiVersion: v1
kind: Pod
metadata:
  name: demopod
spec:
  containers:
    - name:  demopod
      image: nginx
EOF

user $ kubectl get pods
NAME      READY   STATUS                       RESTARTS   AGE
demopod   0/1     CreateContainerConfigError   0          37s

user $ kubectl describe pods demopod
<..>
 Error: container has runAsNonRoot and image will run as root
```

### Example 2: Use a PSP to apply seccomp policies

A second use case for PSPs is to prevent a user from deploying
containers without a [seccomp
policy](https://docs.docker.com/engine/security/seccomp/). By default,
Kubernetes does not apply a seccomp policy to pods, so a default seccomp policy
could be applied for all pods by a PSP.

```bash
admin $ cat <<EOF | kubectl create -f -
apiVersion: extensions/v1beta1
kind: PodSecurityPolicy
metadata:  
  name: seccomppolicy
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: 'docker/default'
    seccomp.security.alpha.kubernetes.io/defaultProfileName:  'docker/default'
spec:
  allowPrivilegeEscalation: false
  allowedHostPaths:
  - pathPrefix: /dev/null
    readOnly: true
  fsGroup:
    rule: RunAsAny
  hostPorts:
  - max: 65535
    min: 0
  runAsUser:
    rule: RunAsAny
  seLinux:
    rule: RunAsAny
  supplementalGroups:
    rule: RunAsAny
  volumes:
  - '*'
EOF
```

If not done previously, the admin user must remove the `ClusterRoleBinding`
for the `privileged` policy, and then add a new `ClusterRole` and `RoleBinding`
to link a user to the new `applyseccompprofile` policy.

```bash
# Delete the default privileged ClusterRoleBinding
admin $ kubectl delete clusterrolebindings ucp:all:privileged-psp-role

# Create a ClusterRole Granting Access to the Policy
admin $ cat <<EOF | kubectl create -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: applyseccompprofile-psp-role
rules:
- apiGroups:
  - policy
  resourceNames:
  - seccomppolicy
  resources:
  - podsecuritypolicies
  verbs:
  - use
EOF

# Define a User to attach to the No Root Policy
admin $ USER=jeff

# Create a RoleBinding attaching the User to the ClusterRole
admin $ cat <<EOF | kubectl create -f -
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: applyseccompprofile-psp-role:$USER
  namespace: default
roleRef:
  kind: ClusterRole
  name: applyseccompprofile-psp-role
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: User
  name: $USER
  namespace: default
EOF
```

As shown in the following example, if a user tries to deploy an `nginx` pod without applying 
a seccomp policy as the pod metadata, Kubernetes automatically applies a policy for the user. 

```bash
user $ cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Pod
metadata:
  name: demopod
spec:
  containers:
    - name:  demopod
      image: nginx
EOF

user $ kubectl get pods
NAME      READY   STATUS    RESTARTS   AGE
demopod   1/1     Running   0          16s

user $ kubectl get pods demopod -o json | jq '.metadata.annotations."seccomp.security.alpha.kubernetes.io/pod"'
"docker/default"
```
