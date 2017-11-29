---
title: Reset a user password
description: Learn how to recover your Docker Datacenter credentials.
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

## User passwords

Docker EE administrators can reset user passwords managed in the Docker EE UI:

1. Log in with administrator credentials to the Docker EE UI.
2. Click **Users** under **User Management**.
3. Select the user whose password you want to change.
4. Select **Configure > Security**.
5. Enter the new password, confirm, and **Save**.

Users passwords managed with an LDAP service must be changed on the LDAP server.

![](../images/recover-a-user-password-1.png){: .with-border}

## Administrator passwords

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
