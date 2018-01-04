---
title: Reset a user password
description: Learn how to recover your Docker Enterprise Edition credentials.
keywords: ucp, authentication, password
ui_tabs:
- version: ucp-3.0
  orlower: true
---
{% if include.version=="ucp-3.0" %}

## Change user passwords

Docker EE administrators can reset user passwords managed in UCP:

1. Log in to UCP with administrator credentials.
2. Click **Users** under **User Management**.
3. Select the user whose password you want to change.
4. Select **Configure > Security**.
5. Enter the new password, confirm, and **Save**.

Users passwords managed with an LDAP service must be changed on the LDAP server.

![](../images/recover-a-user-password-1.png){: .with-border}

## Change administrator passwords

Administrators who need a password change can ask another administrator for help
or use **ssh** to log in to a manager node managed by Docker EE and run:

```none
{% raw %}
docker exec -it ucp-auth-api enzi \
  "$(docker inspect --format '{{ index .Args 0 }}' ucp-auth-api)" \
  passwd -i
{% endraw %}
```

{% endif %}


