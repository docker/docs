---
title: Set the user's session timeout
description: Learn how to set the session timeout for users and other session properties.
keywords: UCP, authorization, authentication, security, session, timeout
---

Docker Universal Control Plane enables setting properties of user sessions,
like session timeout and number of concurrent sessions.

To configure UCP login sessions, go to the UCP web interface, navigate to the
**Admin Settings** page and click **Authentication & Authorization**.

![](../../images/authentication-authorization.png)

## Login session controls

|          Field          |                                                                                                                                                                                                                                             Description                                                                                                                                                                                                                                             |
| :---------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Lifetime Minutes          | The initial lifetime of a login session, starting from the time UCP generates the session. When this time expires, UCP invalidates the active session. To establish a new session, the user must authenticate again. The default is 60 minutes with a minimum of 10 minutes.                                                                                                                                                                                                                                                                             |
| Renewal Threshold Minutes | The time by which UCP extends an active session before session expiration. UCP extends the session by the number of minutes specified in **Lifetime Minutes**. The threshold value can't be greater than **Lifetime Minutes**. The default extension is 20 minutes. To specify that no sessions are extended, set the threshold value to zero. This may cause users to be logged out unexpectedly while using the UCP web interface. The maximum threshold is 5 minutes less than **Lifetime Minutes**. |
| Per User Limit          | The maximum number of simultaneous logins for a user. If creating a new session exceeds this limit, UCP deletes the least recently used session. Every time you use a session token, the server marks it with the current time (`lastUsed` metadata).  When you create a new session that would put you over the per user limit, the session with the oldest `lastUsed` time is deleted. This is not necessarily the oldest session. To disable this limit, set the value to zero. The default limit is 10 sessions.                                                                                                                                                                                                                                                                                                       |
