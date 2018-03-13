---
title: Recover a user password
description: Learn how to recover your Docker Datacenter credentials
keywords: docker, ucp, authentication
---

If you have administrator credentials to UCP, you can reset the password of
other users.

If that user is being managed using an LDAP service, you need to change the
user password on that system. If the user account is managed using UCP,
log in with administrator credentials to the **UCP web UI**, navigate to
the **User Management** tab, and choose the user whose password you want to change.

![](../../images/recover-a-user-password-1.png){: .with-border}

If you're an administrator and forgot your password, you can ask other users
with administrator credentials to change your password.
If you're the only administrator, use **ssh** to log in to a manager
node managed by UCP, and run:

{% raw %}
```none
docker exec -it ucp-auth-api enzi \
  "$(docker inspect --format '{{ index .Args 0 }}' ucp-auth-api)" \
  passwd -i
```
{% endraw %}
