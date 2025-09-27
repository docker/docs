---
title: Personal access tokens
linkTitle: Personal access tokens
description: Create and manage personal Docker access tokens for secure CLI authentication and automation
keywords: personal access tokens, PAT, docker cli authentication, docker hub security, programmatic access
weight: 10
aliases:
 - /docker-hub/access-tokens/
 - /security/for-developers/access-tokens/
---

Personal access tokens (PATs) provide a secure alternative to passwords for Docker CLI authentication. Use PATs to authenticate automated systems, CI/CD pipelines, and development tools without exposing your Docker Hub password.

## Key benefits

PATs offer significant security advantages over password authentication:

- Enhanced security: Investigate token usage, disable suspicious tokens, and prevent administrative actions that could compromise your account if your system is compromised.
- Better automation: Issue multiple tokens for different integrations, each with specific permissions, and revoke them independently when no longer needed.
- Two-factor authentication compatibility: Required when you have two-factor authentication turned on, providing secure CLI access without bypassing 2FA protection.
- Usage tracking: Monitor when and how tokens are used to identify potential security issues or unused automation.

## Who should use personal access tokens?

Use PATs for these common scenarios:

- Development workflows: Authenticate Docker CLI during local development
- CI/CD pipelines: Automate image builds and deployments in continuous integration systems
- Automation scripts: Push and pull images in automated deployment or backup scripts
- Development tools: Integrate Docker Hub access with IDEs, container management tools, or monitoring systems
- Two-factor authentication: Required for CLI access when 2FA is turned on

> [!NOTE]
>
> For organization-wide automation, consider [organization access tokens](/manuals/enterprise/security/access-tokens.md) which aren't tied to individual user accounts.

## Create a personal access token

> [!IMPORTANT]
>
> Treat access tokens like passwords and keep them secure. Store tokens in credential managers and never commit them to source code repositories.

To create a personal access token:

1. Sign in to [Docker Home](https://app.docker.com/).
1. Select your avatar in the top-right corner and from the drop-down menu select **Account settings**.
1. Select **Personal access tokens**.
1. Select **Generate new token**.
1. Configure your token:
   - **Description:** Use a descriptive name that indicates the token's purpose
   - **Expiration date:** Set an expiration date based on your security policies
   - **Access permissions:** **Read**, **Write**, or **Delete**.
1. Select **Generate**. Copy the token that appears on the screen and save it. You won't be able to retrieve the token once you exit the screen.

## Use personal access tokens

Sign in to the Docker CLI using your personal access token:

```console
$ docker login --username <YOUR_USERNAME>
Password: [paste your PAT here]
```

When prompted for a password, enter your personal access token instead of your Docker Hub password.

## Modify personal access tokens

> [!NOTE]
>
> You can't edit the expiration date on an existing personal access token. You must create a new PAT if you need to set a new expiration date.

You can rename, activate, deactivate, or delete a token as needed. You can manage your tokens in your account settings.

1. Sign in to [Docker Home](https://app.docker.com/login).
1. Select your avatar in the top-right corner and from the drop-down menu select **Account settings**.
1. Select **Personal access tokens**.
      - This page shows an overview of all your tokens, and lists if the token was generated manually or if it was
   [auto-generated](#auto-generated-tokens). You can also view the scope of the
   tokens, which tokens are activate and inactive, when they were created, when
   they were last used, and their expiration date.
1. Select the actions menu on the far right of a token row, then select **Deactivate** or **Activate**, **Edit**, or **Delete** to modify the token.
1. After editing the token, select **Save token**.

## Auto-generated tokens

Docker Desktop automatically creates authentication tokens when you sign in, with these characteristics:

- Automatic creation: Generated when you sign in to Docker Desktop
- Full permissions: Include Read, Write, and Delete access
- Session-based: Automatically removed when Docker Desktop session expires
- Account limits: Up to 5 auto-generated tokens per account
- Automatic cleanup: Older tokens are deleted when new ones are created

You can manually delete auto-generated tokens if needed, but they'll be recreated when you use Docker Desktop.

## Fair use policy

When using personal access tokens, be aware that excessive token creation may result in throttling or additional charges. Docker reserves the right to impose restrictions on accounts with excessive PAT usage to ensure fair resource allocation and maintain service quality.

Best practices for fair use include:

- Reuse tokens across similar use cases instead of creating many single-purpose tokens
- Delete unused tokens regularly
- Use [organization access tokens](/manuals/enterprise/security/access-tokens.md) for organization-wide automation
- Monitor token usage to identify optimization opportunities
