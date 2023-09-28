---
title: Get started with Policy Evaluation in Docker Scout
keywords: scout, supply chain, vulnerabilities, packages, cves, policy
description: |
  Policies in Docker Scout let you define supply chain rules and thresholds
  for your artifacts, and track how your artifacts perform against those
  requirements over time
---

> **Beta**
>
> Policy Evaluation is a [Beta](/release-lifecycle/#beta) feature of Docker
> Scout. This feature is available to organizations participating in the
> limited preview program for policies.
>
> If you're interested in trying out this feature, reach out using the form on
> the [Docker Scout product page](https://docker.com/products/docker-scout)
{ .restricted }

In software supply chain management, maintaining the security and reliability
of artifacts is a top priority. Policy Evaluation in Docker Scout introduces a
layer of control, on top of existing analysis capabilities. It lets you define
supply chain rules for your artifacts, and helps you track how your artifacts
perform, relative to your rules and thresholds, over time.

Learn how you can use Policy Evaluation to ensure that your artifacts align
with established best practices.

## How it works

When you activate Docker Scout for a repository, images that you push are
[automatically analyzed](../image-analysis.md). The analysis gives you insights
about the composition of your images, including what packages they contain and
what vulnerabilities they're exposed to. Policy Evaluation builds on top of the
image analysis feature, interpreting the analysis results against the rules
defined by policies.

A policy defines one or more criteria that your artifacts should fulfill. For
example, one of the default policies in Docker Scout is the **Critical
vulnerabilities** policy, which proclaims that your artifacts must not contain
any critical vulnerabilities. If an artifact contains one or more
vulnerabilities with a critical severity, that artifact fails the evaluation.

In Docker Scout, policies are designed to help you ratchet forward your
security and supply chain stature. Where other tools focus on providing a pass
or fail status, Docker Scout policies visualizes how small, incremental changes
affect policy status, even when your artifacts don't meet the policy
requirements (yet). By tracking how the fail gap changes over time, you more
easily see whether your artifact is improving or deteriorating relative to
policy.

Policies don't necessarily have to be related to application security and
vulnerabilities. You can use policies to measure and track other aspects of
supply chain management as well, such as base image dependencies and
open-source licenses.

## Default policies

Docker Scout ships the following three out-of-the-box policies:

- [Critical and high vulnerabilities with fixes](#critical-and-high-vulnerabilities-with-fixes)
- [Critical vulnerabilities](#critical-vulnerabilities)
- [Packages with GPL3+ licenses](#packages-with-gpl3-licenses)

These policies are turned on by default for Scout-enabled repositories. There's
currently no way to turn off or configure these policies.

### Critical and high vulnerabilities with fixes

This policy requires that your artifacts aren't exposed to known
vulnerabilities with a critical or high severity, and where there's a fix
version available. Essentially, this means that there's an easy fix that you
can deploy for images that fail this policy: upgrade the vulnerable package to
a version containing a fix for the vulnerability.

This policy only flags vulnerabilities that were published more than 30
days ago, with the rationale that newly discovered vulnerabilities
shouldn't cause your evaluations to fail until you've had a chance to
address them.

This policy is unfulfilled if an artifact is affected by one or more critical-
or high-severity vulnerability, where a fix version is available.

### Critical vulnerabilities

This policy requires that your artifacts contain no known critical
vulnerabilities. The policy is unfulfilled if your artifact contains one or
more critical vulnerabilities.

This policy flags all critical vulnerabilities, whether or not there's a fix
version available.

### Packages with GPL3+ licenses

This policy requires that your artifacts don't contain packages distributed
under a GPL3+ [copyleft](https://en.wikipedia.org/wiki/Copyleft) license.

This policy is unfulfilled if your artifacts contain one or more packages with
a violating license.
