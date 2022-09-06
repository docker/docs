---
title: Settings
description:
keywords:
---

This page describes the configurable settings for Atomist. Enabling any of these
settings instructs Atomist to carry out an action whenever a specific git event
occurs. These features require that you
[install the Atomist GitHub app](/atomist/integrate/github/#connect-to-github){:
target="blank" rel="noopener" class=""} in your GitHub organization.

To view and manage these settings, go to the
[settings page](https://dso.docker.com/r/auth/policies){: target="blank"
rel="noopener" class=""} on the Atomist website.

## New image vulnerabilities

Scan container images for new critical and high-severity vulnerabilities
introduced via pull requests. New vulnerabilities are displayed as a GitHub
status check on the pull request.

## Dockerfile best practices

Avoid common Dockerfile misconfigurations and usage that can cause security and
operational problems.
[Best practice](/develop/develop-images/dockerfile_best-practices/){:
target="blank" rel="noopener" class=""} violations are shown in GitHub checks on
commits.

## Base image tags

Pin base image tags to digests in Dockerfiles and check for supported tags on
Docker official images. Automatically creates a pull request pinning the
Dockerfile to the latest digest for the base image tag being used.

## Secret scanning

Prevent leaking API keys, access tokens, passwords and other sensitive data by
keeping them out of your codebase. Secret scanning detects and alerts you when
secrets are committed in your code and configuration in a GitHub repository. It
helps prevent secrets from being exposed by adding a failed GitHub Check when a
secret is detected.

Secrets for the following services are automatically detected:

- AWS
- Facebook
- Google
- Mailchimp
- Mailgun
- PayPal
- Picatic API
- Square
- Stripe
- Twilio
- Twitter

This setting is extendable. Enable detection for secrets of additional service
by adding regular expression patterns.
