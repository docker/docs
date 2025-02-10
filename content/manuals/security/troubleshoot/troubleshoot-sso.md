---
title: Troubleshoot SSO
description: Troubleshoot common SSO issues
keywords: sso, admin, security
---

If you experience issues with single sign-on (SSO), refer to the following solutions.

## Groups are not formatted correctly

#### Error message
When this issue occurs, the following error message is common:
```bash
sso error: Some of the groups assigned to the user are not formatted as '<organization name>:<team name>'. Directory groups will be ignored and user will be provisioned into the default organization and team.
```

#### Possible causes
The following causes may create this issue:
- Incorrect group name formatting in your identity provider (IdP): Docker requires groups to follow the format `<organization>:<team>`. If the groups assigned to a user do not follow this format, they will be ignored.
- Non-matching groups between IdP and Docker organization: If a group in your IdP does not have a corresponding team in Docker, it will not be recognized, and the user will be placed in the default organization and team.

#### Affected environments
The following environments can be affected by this issue:
- Docker single sign-on setup using IdPs such as Okta or Azure AD
- Organizations using group-based role assignments in Docker

#### Steps to replicate
Use the following steps to replicate this issue:
1. Attempt to sign in to Docker using SSO.
2. The user is assigned groups in the IdP but does not get placed in the expected Docker team.
3. Review Docker logs or IdP logs to find the error message.

#### Solution
The recommended solution is to update group names in your IdP:
1. Go to your IdP's group management section.
2. Check the groups assigned to the affected user.
3. Ensure each group follows the required format:
```bash
<organization>:<team>
```
4. Update any incorrectly formatted groups to match this pattern.
5. Save changes and retry signing in with SSO.