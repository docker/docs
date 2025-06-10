---
title: Troubleshoot unassociated machines
description: Learn how to troubleshoot common unassociated account issues.
keywords: unassociated machines, unassociated accounts, troubleshoot
tags: [Troubleshooting]
toc_max: 2
---

If you experience issues with unassociated machine management, refer to the
following solutions.

## Machine incorrectly identified as belonging to your organization

### Possible causes

- Docker's machine identification algorithm incorrectly associated the machine
with your organization based on registry usage patterns
- A contractor or temporary user accessed your organization's registries from
a personal machine
- Shared or public registries created false associations

### Affected environments

- All Docker Desktop versions
- All operating systems

### Solution

Docker can add incorrectly identified machines to an ignore list to prevent
future appearances.

[Contact Docker Support](https://hub.docker.com/support/contact) and provide:

- The machine ID
- The reason for why the machine doesn't belong to your organization

## Users cannot sign in to Docker Desktop after enforcement

### Error message

```txt
Sign-in required by your organization
```

### Possible causes

- User is running an outdated version of Docker Desktop that doesn't support
sign-in enforcement
- Network connectivity issues preventing authentication
- User is attempting to sign in with an incorrect email address

### Affected environments

- Docker Desktop versions before 4.37
- Networks with restricted internet access
- Corporate firewalls blocking Docker authentication services

### Solution

1. Verify the user is running Docker Desktop version 4.37 or later.
2. If not, update to the latest version.
3. Ensure network access to Docker's authentication services:
    - https://login.docker.com
    - https://auth.docker.io
4. Confirm the user is signing in with their work email address.

If issues persist, temporarily disable enforcement for that specific machine
while troubleshooting.

## Machine remains in unassociated list after user signs in

### Possible causes

- Auto-provisioning is not enabled for the user's email domain
- The user signed in with a personal email address instead of their work email
- There's a delay in the data refresh cycle

### Affected environments

- Organizations without domain auto-provisioning enabled
- All Docker Desktop versions

### Solution

Recommended solution:

1. Check if the user appears in your organization's member list
2. If not visible, go to Settings > General > Unassociated machines
3. Look for the machine showing an email address
4. Select the machine and choose Add to organization

Alternative solution:

1. Enable auto-provisioning for your verified domains
2. Ask the user to sign in again with their work email address
3. The user will be automatically added to your organization

## Unassociated machines count seems inaccurate

### Possible causes

- Docker Desktop instances are shared between multiple users
- Users have multiple Docker Desktop installations (personal and work machines)
- Data collection limitations due to network restrictions or opt-outs

### Affected environments

- Shared workstations or virtual desktop infrastructure (VDI)
- Air-gapped or restricted network environments
- Organizations with users who have opted out of telemetry

### Solution

Review the machine list to identify patterns:

- Multiple recent activities from the same machine ID may indicate sharing
- Consider the registry access patterns shown in the details
- For shared machines, enforce sign-in and add users as they authenticate
- For air-gapped environments, consider implementing centralized Docker Desktop
configuration

> [!NOTE]
>
> Docker achieves approximately 97% accuracy in machine identification.
A ~3% variance is expected and normal."

## Sign-in enforcement not working for some machines

### Possible causes

- Machines are running Docker Desktop versions that don't support enforcement
- Users haven't restarted Docker Desktop since enforcement was enabled
- Network issues preventing the enforcement check

### Affected environments

- Docker Desktop versions before 4.37
- All operating systems

### Steps to replicate

1. Enable sign-in enforcement for a machine
2. User opens Docker Desktop

- Expected result: Sign-in prompt appears
- Actual result: No prompt, Docker Desktop works normally

### Solution

1. Verify the machine is running Docker Desktop 4.37 or later
2. Ask the user to restart Docker Desktop completely
3. Check that the machine ID matches the one in your enforcement list
4. If the issue persists, disable and re-enable enforcement for that specific
machine

## Auto-provisioning not working after sign-in enforcement

### Possible causes

- Domain auto-provisioning is not enabled
- User signed in with an unverified domain
- Organization has reached its seat limit

### Affected environments

- Organizations without verified domains
- Organizations at seat capacity

### Solution

Recommended solution:

Verify domain auto-provisioning is enabled:

1. Go to Settings > Security > Domain management
2. Ensure the user's email domain is verified and auto-provisioning is enabled

Check organization seat usage:

1. If at capacity, purchase additional seats or remove inactive users
2. Manually add the user if auto-provisioning cannot be enabled

Alternative solution:

1. Set up Single Sign-On (SSO) for automatic user provisioning
2. Enable Just-in-Time (JIT) provisioning through your SSO configuration
