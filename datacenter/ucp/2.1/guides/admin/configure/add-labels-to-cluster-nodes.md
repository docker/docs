---
title: Add labels to cluster nodes
description: Learn how to add metadata to cluster nodes, that can be used to specify constraints when deploying services.
keywords: Docker, cluster, nodes, labels
---

After deploying UCP you can add labels to your nodes. Labels are metadata that
you can use to organize nodes.
You can also use these labels as deployment constraints for your services.

When deploying a service, you can specify constraints, so that the service only
gets scheduled on a node that has a label that fulfills all the constraints
you specify.

As an example, you can apply labels based on their role in the development
lifecycle, or the hardware resources they have.

![](../../images/add-labels-to-cluster-nodes-1.svg)

## Apply labels to a node

Log in with administrator credentials in the **UCP web UI**, navigate to the
**Nodes** page, and choose the node you want to apply labels to.

Click the **Add label** button, and add one or more labels to the node.

![](../../images/add-labels-to-cluster-nodes-2.png){: .with-border}

Once you're done, click **Save Changes**.

You can also do this from the CLI by running:

```none
docker node update --label-add <key>=<value> <node-id>
```

## Add constraint to a service

When deploying a service, you can specify constraints, so that the service only
gets scheduled on a node that has a label that fulfills all the constraints
you specify.

In this example, when users deploy a service they can add constraints for the
service to be scheduled on nodes that have an SSD storage.

You can add deployment constraints to your docker-stack.yml file or when you are creating a stack:
![](../../images/use-contraints-in-stack-deployement.png)

Or you can apply them when you are creating a service.

![](../../images/add-constraint-to-service.png)

You can check if a service have deployment constraints, navigate to the **Services** page,
and choose the service your want to check.
Once you're done, click **Scheduling**.

![](../../images/check-contraints.png)

From here, you can add or remove deployment constraint too.

## Where to go next

* [Store logs in an external system](store-logs-in-an-external-system.md)
