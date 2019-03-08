---
title: Create a service account for a Kubernetes app
description: Learn how to use a service account to give a Kubernetes workload access to cluster resources.
keywords: UCP, Docker EE, Kubernetes, authorization, access control, grant
---

Kubernetes enables access control for workloads by providing service accounts.
A service account represents an identity for processes that run in a pod.
When a process is authenticated through a service account, it can contact the
API server and access cluster resources. If a pod doesn't have an assigned
service account, it gets the `default` service account.
Learn about [managing service accounts](https://v1-8.docs.kubernetes.io/docs/admin/service-accounts-admin/).

In Docker EE, you give a service account access to cluster resources by
creating a grant, the same way that you would give access to a user or a team.
Learn how to [grant access to cluster resources](../authorization/index.md).

In this example, you create a service account and a grant that could be used
for an NGINX server.

## Create the Kubernetes namespace

A Kubernetes user account is global, but a service account is scoped to a
namespace, so you need to create a namespace before you create a service
account.

1.  Navigate to the **Namespaces** page and click **Create**.
2.  In the **Object YAML** editor, append the following text.
    ```yaml
    metadata:
      name: nginx
    ```
3.  Click **Create**.
4.  In the **nginx** namespace, click the **More options** icon,
    and in the context menu, select **Set Context**, and click **Confirm**.

    ![](../images/create-service-account-1.png){: .with-border}

5.  Click the **Set context for all namespaces** toggle and click **Confirm**.

## Create a service account

Create a service account named `nginx-service-account` in the `nginx`
namespace.

1.  Navigate to the **Service Accounts** page and click **Create**.
2.  In the **Namespace** dropdown, select **nginx**.
3.  In the **Object YAML** editor, paste the following text.
    ```yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: nginx-service-account
    ```
3.  Click **Create**.

    ![](../images/create-service-account-2.png){: .with-border}

## Create a grant

To give the service account access to cluster resources, create a grant with
`Restricted Control` permissions.

1.  Navigate to the **Grants** page and click **Create Grant**.
2.  In the left pane, click **Resource Sets**, and in the **Type** section,
    click **Namespaces**.
3.  Select the **nginx** namespace.
4.  In the left pane, click **Roles**. In the **Role** dropdown, select
    **Restricted Control**.
5.  In the left pane, click **Subjects**, and select **Service Account**.

    > Service account subject type
    >
    > The **Service Account** option in the **Subject Type** section appears only
    > when a Kubernetes namespace is present.
    {: .important}

6.  In the **Namespace** dropdown, select **nginx**, and in the
    **Service Account** dropdown, select **nginx-service-account**.
7.  Click **Create**.

    ![](../images/create-service-account-3.png){: .with-border}

Now `nginx-service-account` has access to all cluster resources that are
assigned to the `nginx` namespace.

## Where to go next

- [Deploy an ingress controller for a Kubernetes app](deploy-ingress-controller.md)