---
title: Set up your company for success with Docker
linkTitle: Admin set up
summary: Get the most out of Docker by streamlining workflows, standardizing development environments, and ensuring smooth deployments across your company.
description: Learn how to onboard your company and take advantage of all of the Docker products and features.
tags: [admin]
params:
  time: 20 minutes
  image:
  resource_links:
    - title: Overview of Administration in Docker
      url: /admin/
    - title: Single sign-on
      url: /security/for-admins/single-sign-on/
    - title: Enforce sign-in
      url: /security/for-admins/enforce-sign-in/
    - title: Roles and permissions
      url: /security/for-admins/roles-and-permissions/
    - title: Settings Management
      url: /security/for-admins/hardened-desktop/settings-management/
    - title: Registry Access Management
      url: /security/for-admins/hardened-desktop/registry-access-management/
    - title: Image Access Management
      url: /security/for-admins/hardened-desktop/image-access-management/
    - title: Docker subscription information
      url: "https://www.docker.com/pricing/"
---

Docker's tools provide a scalable, secure platform that empowers your
developers to create, ship, and run applications faster. As an administrator,
you can streamline workflows, standardize development environments, and ensure
smooth deployments across your organization.

By configuring Docker products to suit your company's needs, you can optimize
performance, simplify user management, and maintain control over resources.
This guide helps you set up and configure Docker products to maximize
productivity and success for your team while meeting compliance and security
policies.

## Who’s this for?

- Administrators responsible for managing Docker environments within their
  organization
- IT leaders looking to streamline development and deployment workflows
- Teams aiming to standardize application environments across multiple users
- Organizations seeking to optimize their use of Docker products for greater
  scalability and efficiency
- Organizations with a
  [Docker Business subscription](https://www.docker.com/pricing/)

## What you’ll learn

- Why signing into your company's Docker organization provides access to usage
  data and enhanced functionality
- How to standardize Docker Desktop versions and settings to create a consistent
  baseline for all users, while allowing flexibility for advanced developers
- Strategies for implementing Docker's security configurations to meet company
  IT and software development security requirements without hindering developer productivity

## Features covered

This guide covers the following Docker features:

- [Organizations](/manuals/admin/organization/_index.md): The core structure
  for managing your Docker environment, grouping users, teams, and image
  repositories. Your organization was created with your subscription and is
  managed by one or more owners. Users signed into the organization are
  assigned seats based on the purchased subscription.
- [Enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md):
  By default, Docker Desktop doesn't require sign-in. You can configure
  settings to enforce this and ensure your developers sign in to your
  Docker organization.
- [SSO](/manuals/enterprise/security/single-sign-on/_index.md): Without SSO,
  user management in a Docker organization is manual. Setting
  up an SSO connection between your identity provider and Docker ensures
  compliance with your security policy and automates user provisioning. Adding
  SCIM further automates user provisioning and de-provisioning.
- General and security settings: Configuring key settings ensures smooth
  onboarding and usage of Docker products within your environment. You can also
  enable security features based on your company's specific security needs.

## Who needs to be involved

- Docker organization owner: Must be involved in the process and is required
  for several key steps
- DNS team: Needed during the SSO setup to verify the company domain
- MDM team: Responsible for distributing Docker-specific configuration files to
  developer machines
- Identity Provider team: Required for configuring the identity provider and
  establishing the SSO connection during setup
- Development lead: A development lead with knowledge of Docker configurations
  to help establish a baseline for developer settings
- IT team: An IT representative familiar with company desktop policies to
  assist with aligning Docker configuration to those policies
- Infosec: A security team member with knowledge of company development
  security policies to help configure security features
- Docker testers: A small group of developers to test the new settings and
  configurations before full deployment

## Tools integration

This guide covers integration with:

- Okta
- Entra ID SAML 2.0
- Azure Connect (OIDC)
- MDM solutions like Intune
