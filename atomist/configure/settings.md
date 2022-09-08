---
title: Settings
description:
keywords:
---

This page describes the configurable settings for Atomist. Enabling any of these
settings instructs Atomist to carry out an action whenever a specific Git event
occurs. These features require that you
[install the Atomist GitHub app](/atomist/integrate/github/#connect-to-github){:
target="blank" rel="noopener" class=""} in your GitHub organization.

To view and manage these settings, go to the
[settings page](https://dso.docker.com/r/auth/policies){: target="blank"
rel="noopener" class=""} on the Atomist website.

## New image vulnerabilities

Scan container images for new critical and high-severity vulnerabilities
introduced via pull requests. New vulnerabilities display as a GitHub status
check on the pull request.

## Base image tags

Pin base image tags to digests in Dockerfiles and check for supported tags on
Docker official images. Automatically creates a pull request pinning the
Dockerfile to the latest digest for the base image tag used.
