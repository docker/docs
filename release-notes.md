+++
title = "Release notes"
description = "Release notes for Docker Trusted Registry"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, release notes"]
[menu.main]
parent="smn_dhe"
weight=100
+++


# Release Notes: Docker Trusted Registry & Commercially Supported Docker Engine

This document describes the latest changes, additions, known issues, and fixes for both the Docker Trusted Registry and the commercially supported Docker .

## Prior Versions

These notes refer to the current and immediately prior releases of Docker
Trusted Registry and the commercially supported Docker . For notes on
older versions, see the [prior release notes archive](prior-release-notes.md).

# Docker Trusted Registry 1.4.2
(21 December 2015)

Release notes for the Trusted Registry contain the following sections:

* Additional storage backend
* Fixed or updated with this release

## Additional storage backend
This release introduces using Openstack Swift as a storage backend. Refer to the [configuration documentation](configuration.md) for details on the Swift driver.  

## Fixed or updated with this release
This release addresses the following issues in Docker Trusted Registry 1.4.1.

* Updated the registry from version 2.2.0 to 2.2.1 to ensure that the backend storage Swift driver works correctly.

* Added a link to the release notes from the Trusted Registry UI Updates page.

* The Trusted Registry UI now displays a warning message if an administrator had not yet enabled authentication.

* Corrected an issue where if Trusted Registry administrators were also in a
read-only global team, then they would not see the Trusted Registry admin user
interface.

* Corrected an issue where uploading an image that took longer than five minutes resulted in an authentication failure.

* Fixed several issues that caused the Trusted Registry log files to record excessive unnecessary information.

* Fixed inconsistencies in the garbage collection cron job scheduling between the API and the Trusted Registry UI. This includes:

    * Fixed an issue that if you ran the garbage collection process, the timestamp would not correctly display your last run.

    * Improved the garbage collection cron format. Now a Trusted Registry admin can schedule cron jobs using additional fields such as `@yearly`, `@annually`, `@monthly`, `@weekly`, `@daily`, `@midnight`, and `@hourly`.

    * Fixed an issue where if a Trusted Registry admin did not schedule a garbage collection process, then the Trusted Registry server would log an error.

# Docker Trusted Registry 1.4.1
(24 November 2015)

## Fixed with this release
This release addresses the following issues in Docker Trusted Registry 1.4.0.

* Trusted Registry administrators previously could not pull unlisted repositories in any authorization mode.

* When using LDAP authentication, only users with lowercase letters, numbers, underscores, periods, and hyphens in their usernames in the LDAP server were
synchronized to the Trusted Registry user database. The Trusted Registry now
synchronizes users with usernames containing uppercase letters. If this affects
your organization, perform a LDAP sync from the Trusted Registry UI. Navigate to
Settings > Auth to perform the sync.

* Fixed an issue where Trusted Registry administrators could not list all repositories in the registries. To list them, you must use the `catalog` API using a `bash` shell. The following example lists repositories in a Trusted Registry located at my.dtr.host where the user `admin` has password `password`.

  ```
  bash -c 'host=vagrant.host admin=admin password=password token=$(curl -u $admin:$password -k "https://$host/auth/token?service=$host&scope=registry:catalog:*" | python2 -c "import json,sys;obj=json.load(sys.stdin);print obj[\"token\"]") && curl -k -H "Authorization: Bearer $token" "https://$host/v2/_catalog"'
  ```


# Docker Trusted Registry 1.4
(12 November 2015)

Release notes for the Trusted Registry contain the following sections:

* New features
* Fixed with this release
* Other corrected issues
* Known issues

### New features
This release introduces the following new features. For additional information
on these features, refer to the documentation or when appropriate, the API
documentation.

* Image deletion and garbage collection

  * You can now delete an image in the registry's image index. This step of marking an unwanted image is called a soft delete. Refer to the [documentation](soft-garbage.md).

  * Administrators can use the dashboard or API to configure a task to regularly reclaim the disk space taken up by deleted images. Refer to the [documentation](soft-garbage.md).

* Repositories, Account Management, and interactive API UIs

  * Set up, and manage user accounts, teams, organizations, and repositories from either APIs or through the Trusted Registry user interface. Refer to either the API documentation or the [documentation](accounts.md) for performing tasks in the UI.

  * Search, browse, and discover images created by other users through either APIs or through the Trusted Registry UI.

  * Users, depending on their roles, can access account information through the Trusted Registry UI. Refer to the [documentation](accounts.md) for details.

  * View new API documentation through the Trusted Registry UI. You can also view this [documentation](https://docs.docker.com/docker-trusted-registry/) from Docker, Inc. docs section.

* New APIs

  * There are new APIs for accessing repositories, account management, indexing, searching, and reindexing.

  * You can also view an API and using the Swagger UI, click the Try it out button to perform the action. This might be useful, for example, if you need to reindex.

* Different repository behavior. A repository must first exist before you can push an image to it. This means you must explicitly create (or have it performed for you if you don't have the correct permissions) a repository. This behavior is different than how you would perform this in a free and open-source software registry.

* New experimental feature. Docker Trusted Registry now integrates with Docker Content Trust using Notary. This is an experimental feature that is available with this release. See the [configuration documentation](configuration.md).

### Fixed with this release
This release corrects the following issues in Docker Trusted Registry 1.3.3.

#### LDAP Configuration

* Performance for LDAP user authentication has been significantly improved, reducing the number of required LDAP requests to only a single BIND request to authenticate a user.

* The "Read-Write Search Filter" and "Read-Only Search Filter" fields have been deprecated. You can create organization accounts and teams in the Trusted
Registry to allow for more fine grained access control. Team member lists can be
synced with a group in LDAP.

* The system requires an "Admin Password". Use this password to log in as the
user admin in case the Trusted Registry is unable to authenticate you using
your LDAP server. This account can be used to log in to the Trusted Registry and manage identity and authentication settings.

* Users on your LDAP server are now synced to the Trusted Registry's local
database using your configured "User Search Filter". Objects in LDAP that match
this filter and have a valid "User Login Attribute" are created as a local
user with the "User Login Attribute" as their username. Only these users are
able to log in to the Trusted Registry.

* The "Admin LDAP DN" must be specified to identify the group object on your LDAP server. This should be synced to the system administrators list. The
"Admin Group Member Attribute" should be set to the name of the attribute on
this group object which corresponds to the Distinguished Name of the group
member objects. This setting deprecates the old "Admin Search Filter" field.

#### Other corrected issues

* Corrected an issue where you could switch from none to managed authentication
without creating an administrator.

* Added a "rootdirectory" parameter to the S3 storage option.

#### Known issues

* Organization owners are unable to delete a repository from the UI. You can still delete a repository through the API and a system administrator can still
delete a repository from the UI.

* When using LDAP authentication, only users with valid usernames in the LDAP server are synchronized to the Trusted Registry user database. A valid username
only contains lowercase letters, numbers, underscores, periods, hyphens, and
begins and ends with an alphanumeric character.

* After upgrading to the Trusted Registry 1.4.0, users will not be able to access images from the Trusted Registry if they were previously created without
using namespaces. Docker recommends upgrading to version 1.4.1. After upgrading,
for any repository that was previously pushed to the Trusted Registry, that is
now inaccessible, a Trusted Registry administrator must make it accessible by
the following steps:

  1. If the repository name is not of the form "namespace/repository", then an  administrator must pull all of that repository's tags. For example, you might have an image called `devops_nginx`. The following example shows how you would pull it from a Trusted Registry instance located at my.dtr.host.

        ```
        sudo docker pull --all-tags my.dtr.host/devops_nginx
        ```

  2. Create a new repository. In the Trusted Registry dashboard, navigate to Repositories > New repository.

  3. Select the account that you want to associate to that repository, and enter a repository name in that field and save. If you do not see the account name you wanted to use, then create a new organization or user first. For the example `devops_nginx`, you could use `devops` as the organization and `nginx` as the repository name.

  4. Next, in a `bash` shell, retag all the tags of that repository as seen in the following example:

        ```
        for tag in `sudo docker images | grep my.dtr.host/devops_nginx | awk '{print $2}'`
        do sudo docker tag my.dtr.host/devops_nginx:$tag my.dtr.host/devops/nginx:$tag
        done
        ```
  5. Push the newly tagged version back to the Trusted Registry as seen in the following example:

        `sudo docker push my.dtr.host/devops/nginx`

## Commercially Supported Docker Engine

The commercially Supported (CS) Docker Engine is a packaged release that
identifies a release of Docker Engine for which you can receive support from
Docker or one of its partners. This release is functionally equivalent to the
corresponding Docker Engine release that it references. However, a commercially
supported release also includes back-ported fixes (security-related and priority
defects) from the open source. It incorporates defect fixes that you can use in
environments where new features cannot be adopted as quickly for consistency and
compatibility reasons.  

### Commercially Supported Docker Engine 1.9.1-cs3
(6 January 2016)

This release addresses the following issues:

* The commercially supported Engine 1.9.1-cs3 now supports multi-host networking
for all the kernels that the base CS Engine is supported on.

>**Note**: Centos 7 has its firewall enabled by default and it prevents the VXLAN tunnel from communicating. If this applies to you, then after installing the CS Engine, execute the following command in the Linux host:

`sudo firewall-cmd --zone=public --permanent --add-port=4789/udp`


* Corrected an issue where Docker didn't remove the Masquerade NAT rule from `iptables` when the network was removed. This caused the gateway address to be
incorrectly propagated as the source address of a connection.

* Fixed an issue where if the daemon started multiple containers concurrently, then the `/etc/hosts` files were incompletely populated. This issue occurred randomly.

* Corrected an issue where the same IP address for different Docker containers resulted in network connection inconsistencies. Now each container has a separate IP address.

* Corrected an issue where the IPv6 gateway was not created when using custom networks although the network had a configured gateway.

* Fixed an issue where users might have experienced a panic error if the  daemon was started with the `—cluster-store` option, but without the `—cluster-advertise` option.

### Commercially Supported Docker Engine 1.9.1-cs2
(4 December 2015)

Starting with this release, upgrading minor versions, for example, from 1.9.0 to 1.9.1, is faster and easier. See the [upgrade](install/upgrade.md) documentation for details.

You can refer to the detailed list of all changes since the release of CS Engine
1.9.0
https://github.com/docker/docker/releases.

### Commercially Supported Docker Engine 1.9.0
(12 November 2015)

Highlighted feature summary:

* Network Management and Plugins. Networks are now first class objects that can be listed, created, deleted, inspected, and connected to or disconnected from a
container. They can be manipulated outside of the container themselves and are
fully manageable on its own lifecycle. You can also use plugins to extend
network functionality.

* Docker, Inc. now provides support for the in-box Overlay (for cross-host networking) and Bridge network plugins. You can find more information about how
to manage networks and using network plugins in the [documentation](https://docs.docker.com/engine/userguide/networking/dockernetworks/).

* Volume Management and Plugins. Volumes also become discrete, manageable objects in Docker. Volumes can be listed, created, deleted, and inspected.
Similar to networks, they have their own managed lifecycle outside of the
container. Plugins allow others to write and extend the functionality of volumes
or provide integration with other types of storage.

* The in-box volume driver is included and supported. You can find more information about how to manage volumes  and using  volume plugins in the
documentation.

* Docker Content Trust. Use Content Trust to both verify the integrity and the publisher of all the data received from a registry over any channel. Content Trust is currently only supported using Docker Hub notary servers.

* Updated the release cadence of the CS Docker Engine. Starting with this version, Docker supports **every** major release of Docker Engine from open
source with three releases under support at one time. This means you’ll be able
to take advantage of the latest and greatest features and you won’t have to wait
for a supported release to take advantage of a specific feature.

Refer to the [detailed list](https://github.com/docker/docker/releases) of all changes since the release of CS Engine 1.6.
