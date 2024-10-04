---
title: Deploy 
description: Deploy your Docker setup across your company.
weight: 40
---

Enforce SSO
CAUTION: This step will affect any existing users signing into your Docker organization.  Please communicate with your users and carefully read and follow the list of instructions in the admin UI before confirming this step!  Enforcing SSO means that anyone who has a Docker profile with an email address that matches your verified domain MUST log in using your SSO connection.  Make sure the Identity provider groups associated with your SSO connection cover all the developer groups that you want to have access to the Docker Subscription.

Deploy configuration settings and enforce sign in to users
CAUTION: This step will affect all existing users of Docker Desktop.  Please communicate with your users before taking this step, and ensure IT and MDM teams are ready for any unexpected issues to arise.  Have the MDM team deploy the configuration files for Docker to all users.  

Congratulations, you have successfully completed the admin implementation process for Docker!  
