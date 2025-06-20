---
title: Monitor activity
description: Use Docker features to monitory user activity and compliance with your organization settings.
weight: 50
---

In hardened environments, it’s not enough to configure secure defaults. You
also need ongoing visibility into how Docker is being used, where settings may
drift, and whether your container environments meet compliance requirements.

This module walks you through how to monitor Docker organization activity,
audit Desktop settings across your fleet, and integrate with external tooling
like SIEM or Slack.

## Prerequisites

Before you begin, ensure you have:

- A Docker Business subscription
- Organization owner access to your Docker organization
- Docker Desktop deployed across managed machines
- Optional. Docker Scout enabled for image analysis and SBOM indexing

## Step one: Review activity logs

Docker automatically tracks high-level organizational activity such as:

- User sign-ins
- Team and role changes
- Repository actions
- SSO enforcement status
- Domain verification events

To view logs:

1. Go to the [Docker Admin Console](https://app.docker.com/admin)
2. Select your organization.
3. Navigate to **Activity Logs**.

You can search by event type or user to trace changes across your org.

## Step two: Monitor Desktop settings compliance

If you're using centralized settings via `admin-settings.json` or the Admin
Console, you can audit compliance across your fleet.

To view compliance reports:

1. In the Admin Console, go to **Settings management**.
2. Open the **Reporting** tab to see which machines are:
    - Compliant with enforced settings
    - Out of sync or missing required controls

## Step three: Set up Docker Scout for image visibility

Use [Docker Scout](https://docs.docker.com/scout/) to track security posture at
the container image level. Scout supports:

- Software Bill of Materials (SBOM) indexing
- Vulnerability scanning
- Policy enforcement
- Exceptions and remediation tracking

You can integrate Scout with:

- GitHub Actions
- GitLab CI/CD
- Jenkins
- Azure DevOps
- Artifactory, ECR, ACR, and more

To start, visit the [Docker Scout integrations overview](https://docs.docker.com/scout/integrations/).

## Step four: Enable alerts and external integrations

For real-time visibility, consider integrating Docker logs and insights with:

- Slack: Docker Scout supports alerting via Slack for policy violations and
vulnerability reports
- SIEM tools: Export activity logs or Scout scan results into tools like
Splunk or Sentinel
- Webhook-based integrations: Set up Docker Hub [webhooks](https://docs.docker.com/docker-hub/repos/manage/webhooks/) for image pull/push notifications

## Best practices

- Review activity logs regularly (weekly or during incident response).
- Monitor settings compliance to detect drift across endpoints.
- Enable SBOM indexing and scan enforcement via Docker Scout.
- Push logs and alerts into your broader monitoring and alerting systems.
- Use webhook or CI integrations to track image updates and policy violations
in real time.
