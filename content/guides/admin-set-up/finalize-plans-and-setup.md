---
title: Finalize plans and begin setup
description:
weight: 20
---

Create SSO Connection
Once the domain is verified, the next step is to create the SSO connection.  This will involve your identity provider team to configure the identity groups and help set up the SSO connection.  Note that this step of creating the SSO connection will not affect the Docker Desktop user experience, and you will be able to test before enforcing SSO for all users.  The steps in the process are located here.  

Finalize baseline configuration settings
Come to agreement between your Docker organization owner and your Development lead on the settings to be configured as part of the Docker baseline.  This should include the enforce sign in configuration for your Docker organization.

Manage Organizations
If you have more than one organization, it’s recommended that you either consolidate them into one organization or use the account hierarchy feature to manage multiple organizations.  Please work with the CS and implementation teams to make this happen.

Finalize security configuration settings
Come to agreement between your Infosec representative, Docker organization owner, and Development lead on the security features/settings to be preset as part of your Docker baseline configuration.

Send finalized settings files to MDM team
Once all of the settings have been entered to the files that need to be distributed, pass the files to your MDM team to package up.  It’s highly recommended that the next step in week 3 is a test distribution to a small number of Docker Desktop users to verify the functionality works as expected.

Set up free tier Docker product entitlements included in the subscription
Set up the cloud builder for free monthly minutes in Docker Build Cloud, and up to three repositories to monitor via Docker Scout.  Please note that your free entitlements stop when your limits are exceeded so there is no fear of a surprise cost overage.  The instructions on setting up the cloud builder are located on build.docker.com and there is a video walkthrough here, and the instructions on adding a repository for scout monitoring is here for Docker Hub repositories, and here for integration to other image registries.