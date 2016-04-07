+++
title = "Auth configuration"
description = "Authentication configuration for Docker Trusted Registry"
keywords = ["docker, documentation, about, technology, understanding, configuration, auth,  enterprise, hub, registry"]
[menu.main]
parent="workw_dtr_configure"
weight=6
+++

# Configure your auth settings

Use the Auth screen to add users and control their access to the Trusted
Registry. This document explains the three authentication methods and how to
import users into the Trusted Registry through the UI.

1. From the Trusted Registry dashboard, navigate to Settings > Auth.
2. Use the drop down menu to select either Managed or LDAP. The screen refreshes
reflecting your choice.

Note that once you start using a particular method, you need to stick with your choice.

There are three authentication methods:

* [None](#no-authentication-none)
* [Managed](#managed-authentication)
* [LDAP](#ldap-authentication)

![Auth settings page</admin/settings#auth>](../images/admin-settings-auth.png)

> **Note**: If you have issues logging into the Docker Trusted Registry admin web interface after changing the authentication settings, you may need to [troubleshoot DTR](../monitor-troubleshoot/troubleshoot.md).

## No authentication (None)

No or `None` authentication means that everyone can access your Trusted Registry
web administration site. This is the default setting when you first install the
Trusted Registry. One of your first configuration tasks is to switch your
authentication to either managed or LDAP so you can create the Trusted Registry
administrator. Until you do, you can't create repos, nor push or pull images.
See the [configuration overview](configuration.md) to manually create an admin account.

## Managed authentication

With `Managed` authentication, the Trusted Registry admin can manually control users' access by setting username/password pairs. The admin can then [use the API](http://docs.docker.com/apidocs/v1.3.3/) to give these users global "admin", "read-write" or "read-only" privileges while assigning them organization, team, or user repository access. Note that you can **only** set the global role of `Admin - all repositories` though the UI for the admin. The global roles for `Read-write - all repositories` and `Read-only - all repositories` have been deprecated.

When you create users and assign their roles through the API, you do not need
to assign those users roles using the Trusted Registry admin UI.

1. Choose the appropriate button to either add one user, or to upload a CSV file containing username, password pairs, and selection boxes for "admin",
"read-write", and "read-only" roles.
2. Click Save.

If you make an error, or need to remove a user, you can also delete them from this screen.

## LDAP authentication

Use LDAP authentication to integrate your Trusted Registry into your
organization's existing LDAP user and authentication database. To improve the performance of the Trusted Registry's Access Control Lists,
User and Group membership data is synced into Docker Trusted Registry's database
at a configurable *LDAP Sync Interval*. User passwords are not transferred
during syncing. The Trusted Registry defers to the LDAP server to validate
username/password pairs.

LDAP syncing creates new users that that do not already exist in the Trusted Registry. Any existing users that are not found by the LDAP sync are marked as inactive and not deleted. You can also sync team membership with the LDAP group. This is performed after you have finished configuring your settings.

Because connecting to LDAP involves existing infrastructure external to the
Trusted Registry and Docker, you need to gather the details required to
configure the Trusted Registry for your organization's particular LDAP
implementation.

### Add additional users through user sets

In the User Set section, you can add parameters to further refine your LDAP
integration. Clicking Add User Set, displays additional User Sets in the UI.
There is no limit to adding User Sets. Use this additional User Set to target
other users that may be located in different nodes of your organization. When
the next LDAP sync occurs, the Trusted Registry will find all the different sets
of users. The set of Trusted Registry users is the union of all of those sets.

Each of the parameters are explained in the LDAP Configuration options section
in this document.

### Test your sync

You can test that you have the correct LDAP server information by connecting to
the LDAP server from inside a Docker container running on the same server as
your Docker Trusted Registry:

If the LDAP server is configured to use *StartTLS*, then you need to
add `-Z` to the `ldapsearch` following command example.

```
docker run --rm -it svendowideit/ldapsearch -h <LDAP Server hostname> -b <User Base DN> -D <Search User DN> -w <Search User Password>
```

The result of this query should be a (very) long list. If you get an
authentication error, then the details you have are not sufficient. Contact
your organization's LDAP team.

The *User Login Attribute* key setting must match the field used in the LDAP
server for the user's login-name. On OpenLDAP, it's generally `uid`, and on
Microsoft Active Directory servers, it's `sAMAccountName`. The `ldapsearch`
output should allow you to confirm which setting you need.

![LDAP authentication settings page</admin/settings#auth>](../images/admin-settings-authentication-ldap.png)

### LDAP Configuration options

* **Admin Password**: *required*, use this password to login as the user `admin` in case Docker Trusted Registry is unable to authenticate you using your LDAP server. This account may be used to login to the Trusted Registry and correct identity and authentication settings.
* **LDAP Server URL**: *required*, defaults to null, LDAP server URL (for example, - ldap://example.com).
* **Use StartTLS**: defaults to unchecked, check to enable StartTLS.
* **User Base DN**: *required*, defaults to null, user base DN in the form (for example, - dc=example,dc=com).
* **User Name Attribute Is Email**: if your user names in your LDAP server are email addresses, this will replace the @ sign with _ when storing the user since @ signs are not supported by the open source registry user naming scheme.
* **Scope One Level**: this value is used to indicate searching all entries one level under the base DN, but not including the base DN and not including any entries under that one level under the base DN.

![LDAP Scope One Level</configure/settings#auth>](../images/ldap-scope-down.png)

* **User Login Attribute**: *required*, defaults to null, user login attribute (for example - uid or sAMAccountName).
* **Search User DN**: *required*, defaults to null, search user DN (for example,   domain\username).
* **Search User Password**: *required*, defaults to null, search user password.
* **LDAP Sync Interval**: *required*, defaults to 1h0m0s, sets the interval for Docker Trusted Registry to sync with the LDAP database.
* **User Search Filter**: users on your LDAP server are synced to Docker Trusted Registry's local database using this search filter. Objects in LDAP that match
this filter and have a valid "User Login Attribute" are created as a local user
with the "User Login Attribute" as their username. Only these users are able to
login to the Trusted Registry.
* **Admin LDAP DN**: *required*, this field is used to identify the group object on your LDAP server which is synced to the system administrators list.
* **Admin Group Member Attribute**: *required*, this value matches the name of the attribute on this group object which corresponds to the Distinguished Name
of the group member objects.

### Confirm login with current configuration

Test your current LDAP configuration before saving it by entering a test
username and password.   Click Try Login. If the login succeeds, your
configuration is working.

## See also

* [Configure DTR](config-general.md)
* [Troubleshoot DTR](../monitor-troubleshoot/troubleshoot.md)
