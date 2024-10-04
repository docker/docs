---
title: Testing
description: Test your Docker setup.
weight: 30
---

Ensure supported version of Docker Desktop
CAUTION: This step could affect the experience for users on older versions of Docker Desktop.  Existing users may have older versions of Docker Desktop that are no longer supported or are out of date.  It is highly recommended that everyone update to a supported version.  We recommend using a MDM solution to manage the version of Docker Desktop for users.  Users may also get Docker Desktop directly from Docker or through a company software portal.  In any of these cases it's important that the users are upgraded to a supported Docker Desktop version.

SSO and SCIM testing
If you want to use SCIM for further automation of provisioning and deprovisioning of users, there are some additional configurations required by your identity provider team.  Please see here for a list of settings.  Once all of the configuration is done, it is time for testing of the SSO connection, group mapping, provisioning, and SCIM (if configured).  SSO testing can be done by logging into Docker Desktop or Docker Hub with the email address associated with a Docker account that also belongs to the domain that was verified.  Users that log in using their Docker usernames will continue to be unaffected by the SSO/SCIM setup. NOTE: Some users may need CLI based logins to Docker Hub, and for this they will need a personal access token (PAT).  Please see here for more details. 

Test Registry/Image Access Management
CAUTION: This step will affect any existing users signing into your Docker organization.  Please communicate with your users before completing this step.  If you are planning to use Registry Access Management (RAM) and/or Image Access Management (IAM), configure the settings in the Docker admin portal.  Please see here for RAM details, and here for the video walkthrough.  Please see here for the IAM details, and here for the video walkthrough.

Deploy settings and enforce sign in to test group
Deploy the Docker settings and enforce sign in to a small group of test users via MDM.  Have this group test their developer workflows with containers using Docker Desktop and Hub to confirm all settings and enforce sign in are working as expected.  

Test Build Cloud capabilities
Have one of your Docker Desktop testers connect to the cloud builder you created and do a build.  See here for more details.

Verify Scout monitoring of repositories
Check the scout.docker.com portal to verify the data and trending for the repositories enabled.
