---
description: Find the answers to common security related FAQs
keywords: Docker, Docker Hub, Docker Desktop secuirty FAQs, secuirty, platform, Docker Scout, admin, security
title: General security FAQs
---

### How do I report a vulnerability?

If you’ve discovered a security vulnerability in Docker, we encourage you to report it responsibly. Report security issues to security@docker.com so that they can be quickly addressed by our team.

### How are passwords managed when SSO isn't used? 

Passwords are encrypted and salt-hashed. If you use application-level passwords instead of SSO, you are responsible for ensuring that your employees know how to pick strong passwords, don't share passwords, and don't reuse passwords across multiple systems. 

### Does Docker require password resets when SSO isn't used? 

Passwords aren't required to be periodically reset. NIST no longer recommends password resets as part of best practice.

### Does Docker lockout users after failed sign-ins? 

Docker Hub’s global setting for system lockout is after 10 failed sign in attempts in a period of 5 minutes, and the lockout duration is 5 minutes. The same global policy applies to authenticated Docker Desktop users and Docker Scout, both of which use Docker Hub for authentication.

### Do you support physical MFA with YubiKeys? 

You can configure this through SSO using your IdP. Check with your IdP if they support physical MFA.

### How are sessions managed and do they expire?

Docker Desktop uses tokens to manage sessions after a user signs in. Docker Desktop signs you out after 90 days, or 30 days of inactivity.

In Docker Hub, you need to re-authenticate after 24 hours. If users are authenticating using SSO, the default session timeout for the IdP is respected.

Custom settings per organization for sessions aren't supported.

### How does Docker attribute downloads to us and what data is used to classify or verify the user is part of our organization? 

Docker Desktop downloads are linked to a specific organization by the user's email containing the customer's domain. Additionally, we use IP addresses to correlate users with organizations.

### How do you attribute that number of downloads to us from IP data if most of our engineers work from home and aren’t allowed to use VPNs? 

We attribute users and their IP addresses to domains using 3rd party data enrichment software, where our provider analyzes activity from public and private data sources related to that specific IP address, then uses that activity to identify the domain and map it to the IP address.

Some users authenticate by signing in to Docker Desktop and joining their domain's Docker organization, which allows us to map them with a much higher degree of accuracy and report on direct feature usage for you. We highly encourage you to get your users authenticated so we can provide you with the most accurate data.

### How does Docker distinguish between employee users and contractor users? 

Organizations set up in Docker use verified domains and any team member with an email domain other than what's verified is noted as a "Guest" in that organization.

### How long are Docker Hub logs available? 

Docker provides various types of audit logs and log retention varies. For example, Docker Hub Activity logs are available for 90 days. You are responsible for exporting logs or setting up drivers to their own internal systems.  

### Can I export a list of all users with their assigned roles and privileges and if so, in what format?

Using the [Export Members](../../admin/organization/members.md#export-members) feature, you can export to CSV a list of your organization's users with role and team information. 

### How does Docker Desktop handle and store authentication information?

Docker Desktop utilizes the host operating system's secure key management for handling and storing authentication tokens necessary for authenticating with image registries. On macOS, this is [Keychain](https://support.apple.com/guide/security/keychain-data-protection-secb0694df1a/web); on Windows, this is [Security and Identity API via Wincred](https://learn.microsoft.com/en-us/windows/win32/api/wincred/); and on Linux, this is [Pass](https://www.passwordstore.org/). 

### How does Docker Hub secure passwords in storage and in transit? 

This is applicable only when using Docker Hub's application-level password versus SSO/SAML. When using SSO, Docker Hub doesn't store passwords. Application-level passwords are hashed in storage (SHA-256) and encrypted in transit (TLS).

### How do we de-provision access to CLI users who use personal access tokens instead of our IdP? We use SSO but not SCIM

If SCIM isn't enabled, you have to manually remove PAT users from the organization in our system. Using SCIM automates this.

### What metadata is collected from container images that Scout analyzes?

For information about the metadata stored by Docker Scout, see [Data handling](../../scout/data-handling.md).

### To which portions of the host filesystem do containers have read and write access? Can containers running as root gain access to admin-owned files or directories on the host? 

File sharing (bind mount from the host filesystem) uses a user-space crafted file server (running in `com.docker.backend` as the user running Docker Desktop), so containers can’t gain any access that the user on the host doesn’t already have.

### How are Extensions within the Marketplace vetting for security prior to placement? 

Security vetting for extensions is on our roadmap however this vetting isn't currently done. 

Extensions are not covered as part of Docker’s Third-Party Risk Management Program.

### Can I disable private repos in my organization via a setting to make sure nobody is pushing images into Docker Hub? 

No. With [Registry Access Management](../../security/for-admins/registry-access-management.md) (RAM), administrators can ensure that their developers using Docker Desktop only access allowed registries. This is done through the Registry Access Management dashboard on Docker Hub. 

