---
title: Reset a user password
description: Learn how to recover your Docker Datacenter credentials
keywords: ucp, authentication
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
---

{% if include.ui %}
{% if include.version=="ucp-3.0" %}
This topic is under construction.

{% elsif include.version=="ucp-2.2" %}

If you have administrator credentials to UCP, you can reset the password of
other users.

If that user is being managed using an LDAP service, you need to change the
user password on that system. If the user account is managed using UCP,
log in with administrator credentials to the UCP web UI, navigate to
the **Users** page, and choose the user whose password you want to change.
In the details pane, click **Configure** and select **Security** from the
dropdown.

![](../images/recover-a-user-password-1.png){: .with-border}

Update the user's password and click **Save**.

If you're an administrator and forgot your password, you can ask other users
with administrator credentials to change your password.
If you're the only administrator, use **ssh** to log in to a manager
node managed by UCP, and run:

```none
{% raw %}
docker exec -it ucp-auth-api enzi \
  "$(docker inspect --format '{{ index .Args 0 }}' ucp-auth-api)" \
  passwd -i
{% endraw %}
```

{% endif %}
{% endif %}
