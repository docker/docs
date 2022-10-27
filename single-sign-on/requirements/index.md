---
description: Single Sign-on requirements
keywords: Single Sign-on, SSO, sign-on, requirements
title: Requirements
---

## Prerequisites

* You must first notify your company about the new SSO login procedures
* Verify that your org members have the latest Docker Desktop version 4.4.2, or later, installed on their machines
* New org members must create a Personal Access Token (PAT) to sign in to the CLI, however existing users can use their username and password during the grace period as specified below
* Confirm that all CI/CD pipelines have replaced their passwords with PATs
* For your service accounts, add your additional domains or enable it in your IdP
* Test SSO using your domain email address and IdP password to successfully sign in and log out of Docker Hub

## Create a Personal Access Token (PAT)

Before you configure SSO for your organization, new members of your organization must [create an access token](../../docker-hub/access-tokens.md) to sign in to the CLI. There is a grace period for existing users, which will expire in the near future. Before the grace period ends, your users will be able to sign in from Docker Desktop CLI using their previous credentials until PATs are mandatory.
In addition, you should add all email addresses to your IdP.
