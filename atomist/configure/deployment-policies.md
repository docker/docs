---
title: Deployment policies
description:
keywords:
---

{% include atomist/disclaimer.md %}

A deployment policy is a way to determine the quality of a container image, and
assess whether it is safe to deploy. It's like a test suite, ensuring that the
final image is well-formed and compliant. Images failing to comply with a policy
check are prevented from being deployed to your environments.

Each policy consists of one or more rules. Each rule describes a property
requirement for images. The following are a few example properties that policy
rules can check for:

- Does the image have the required labels?
- Was the image created by a trusted builder?
- Was the image created from a known git commit SHA?
- Has the image been scanned for vulnerabilities?
- Does the image contain any vulnerabilities that are not already present in the
  target branch?

Atomist comes with a set of built-in rules, but you can also define your own
ones using GitHub Checks.

## When to use deployment policies

Deployment policies shine in a GitOps workflow, where new candidate images can
be pulled into a workload once they are ready. The combination of GitOps
controllers, admission controllers, and modular image policy, gives teams the
ability to plug consistent validation into their cloud native delivery.

## Kubernetes integration

Integrating your deployment policies with Kubernetes is done via
[admissions controllers](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/).
To get started with deployment policies for Kubernetes,
[create the Atomist admission controllers](../integrate/kubernetes.md) in your
cluster.
