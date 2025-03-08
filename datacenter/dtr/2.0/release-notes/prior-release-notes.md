---
description: Archived release notes for Docker Trusted Registry
keywords: docker, documentation, about, technology, understanding, enterprise, hub, registry, Docker Trusted Registry, release
redirect_from:
- /docker-trusted-registry/prior-release-notes/
- /docker-trusted-registry/release-notes/prior-release-notes/
title: DTR release notes archive
---

This document contains the release notes for all versions of Docker Trusted
Registry.

## Version 1.4.3
(22 February 2016)

The Trusted Registry is supported on SUSE Linux Enterprise 12 OS.

This release addresses the following issues in Docker Trusted Registry 1.4.2.

* Improved the Trusted Registry UI response when performing certain operations with a large set of users.

* Created a new Trusted Registry screen where image tags in a repository are displayed. This fixed the issue where long image tags were truncated in the UI.

* You can now download the Trusted Registry for offline installation. Refer to the documentation.

* Corrected an issue where if the Trusted Registry was set to a non default port, users couldn’t push images to it.

* Improved LDAP configuration. There are now additional user search filters in the Trust Registry UI. The location is Settings > Auth. Select LDAP authentication method. The filters are:

    * `UsernameAttrIsEmail`
    * `ScopeOneLevel`

* Fixed an issue where the Trusted Registry correctly updates team members after an LDAP sync. This removed duplication of users if they were moved to a different team.

* Previously, if you started the Trusted Registry 1.4.2 with CS Engine 1.7.0 onwards, it might not start because `docker` might start the Trusted Registry containers in an order that makes [links impossible to create](https://github.com/moby/moby/issues/17118). Using CS Engine 1.9 and later with the latest Trusted Registry includes creation of a custom network that allows all containers to connect to each other without links. This means that every time the Trusted Registry starts up, there should be no error.

    Also, when you upgrade CS Engine from 1.6 to 1.9 and the Trusted Registry admin server starts, it checks if it's running with links enabled. If that happens, the Trusted Registry restarts everything, creating the new network if necessary and removing the links, replacing them with a custom "dtr" network.

* This version deprecates the search API exposed at the /api/v0/index/search
endpoint. Starting on DTR 2.0 this API endpoint will no longer be available.

* On DTR 2.0 the autocomplete API exposed at the /api/v0/index/autocomplete
endpoint, will start returning a new result structure.

* On DTR 2.0 all API endpoints that return UserAccess and RepoUserAccess objects
are going to return objects with different properties.

## Version 1.4.2
(21 December 2015)

Release notes contain the following sections:

* Additional storage backend
* Fixed or updated with this release

**Additional storage backend**

This release introduces using Openstack Swift as a storage backend. Refer to the [configuration documentation](../configure/config-storage.md) for details on the Swift driver.

**Fixed or updated with this release**
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

## Version 1.4.1
(24 November 2015)

**Fixed with this release**
This release addresses the following issues in Docker Trusted Registry 1.4.0.

* Trusted Registry administrators previously could not pull unlisted repositories in any authorization mode.

* When using LDAP authentication, only users with lowercase letters, numbers, underscores, periods, and hyphens in their usernames in the LDAP server were
synchronized to the Trusted Registry user database. The Trusted Registry now
synchronizes users with usernames containing uppercase letters. If this affects
your organization, perform a LDAP sync from the Trusted Registry UI. Navigate to
Settings > Auth to perform the sync.

* Fixed an issue where Trusted Registry administrators could not list all repositories in the registries. To list them, you must use the `catalog` API using a `bash` shell. The following example lists repositories in a Trusted Registry located at my.dtr.host where the user `admin` has password `password`.

```none
  bash -c 'host=vagrant.host admin=admin password=password token=$(curl -u $admin:$password -k "https://$host/auth/token?service=$host&scope=registry:catalog:*" | python2 -c "import json,sys;obj=json.load(sys.stdin);print obj[\"token\"]") && curl -k -H "Authorization: Bearer $token" "https://$host/v2/_catalog"'
```


## Version 1.4.3
(12 November 2015)

Release notes for the Trusted Registry contain the following sections:

* New features
* Fixed with this release
* Other corrected issues
* Known issues

**New features**
This release introduces the following new features. For additional information
on these features, refer to the documentation or when appropriate, the API
documentation.

* Image deletion and garbage collection

  * You can now delete an image in the registry's image index. This step of marking an unwanted image is called a soft delete.

  * Administrators can use the dashboard or API to configure a task to regularly reclaim the disk space taken up by deleted images.

* Repositories, Account Management, and interactive API UIs

  * Set up, and manage user accounts, teams, organizations, and repositories from either APIs or through the Trusted Registry user interface.

  * Search, browse, and discover images created by other users through either APIs or through the Trusted Registry UI.

  * Users, depending on their roles, can access account information through the Trusted Registry UI.

  * View new API documentation through the Trusted Registry UI. You can also view this documentation from Docker, Inc. docs section.

* New APIs

  * There are new APIs for accessing repositories, account management, indexing, searching, and reindexing.

  * You can also view an API and using the Swagger UI, click the Try it out button to perform the action. This might be useful, for example, if you need to reindex.

* Different repository behavior. A repository must first exist before you can push an image to it. This means you must explicitly create (or have it performed for you if you don't have the correct permissions) a repository. This behavior is different than how you would perform this in a free and open-source software registry.

* New experimental feature. Docker Trusted Registry now integrates with Docker Content Trust using Notary. This is an experimental feature that is available with this release.

**Fixed with this release**
This release corrects the following issues in Docker Trusted Registry 1.3.3.

**LDAP Configuration**

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

**Other corrected issues**

* Corrected an issue where you could switch from none to managed authentication
without creating an administrator.

* Added a "rootdirectory" parameter to the S3 storage option.

**Known issues**

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


## Version 1.3.3
(18 September 2015) (amended: 2 November 2015)

This release corrects the following issues in Docker Trusted Registry 1.3.2

* Fixed an issue related to LDAP integration for users of Oracle Virtual Directory.

* Corrected an issue where Docker Trusted Registry would not accept a given certificate if the configured domain was only in the Subject Alternative Names
(SANs) field and not in the Common Name (CN) field of the certificate.

* Docker, Inc. discovered an issue in which the tokens used in authorization caused a break in certain deployments that utilized a load balancer in front of
multiple Trusted Registry instances to achieve high availability. We regret any
inconvenience this may have caused you and is working on a future fix.

## Version 1.3.2
(16 September 2015)

This release addresses the following change in Docker Trusted Registry 1.3.2 and is only available to customers who purchased DTR through Amazon Web Services (AWS) Marketplace.

* Docker Trusted Registry (DTR) now supports Amazon Web
Services (AWS) Integrated Billing. Previously, AWS users were required to
separately purchase a DTR license from Docker. AWS users can try DTR
out-of-the-box.

## Version 1.3.1
(31 August 2015)

This release corrects the following issues in Docker Trusted Registry 1.3.0

* The dashboard page was calculating incorrect stats.
* LDAP group sync failed to handle paginated results for extremely large groups.
* The repo delete endpoint returned incorrect error codes under certain conditions.

## Version 1.3.0
(26 August 2015)

This release addresses a few bugs and issues in Docker Trusted Registry 1.2.0 and introduces some new features and functionality, including:

* A completely new user-interface for the Admin application brings Docker Trusted Registry in line with other Docker products and provides greater ease-of-use.

* A new Accounts & Repos API provides new fine-grained role-based access control down to the per-repo level. See the [API's documentation](/apidocs/v1.3.3/) for more information.

* Improvements to the handling of configuration changes so that fewer restarts are required.

* Multiple security improvements and bug fixes.

## Version 1.2.0
(23 July 2015)

This release adds CentOS support and addresses a few bugs and issues in Docker Trusted Registry 1.1.0:

* Fixes an issue where for certain configurations of Docker Trusted Registry, proxy configuration settings and variables were not being passed to all Docker Trusted Registry containers and thus were not being respected.
* Documentation links in the UI now point to correct docs.
* Generated support info bundles have been scrubbed to remove highly sensitive data.
* Certifies support for CentOS 7.1.

## Version 1.1.0
(23 June 2015)

This release of Docker Trusted Registry (formerly DHE) adds major integration
with the AWS and Azure marketplaces, giving customers a smoother installation
path. Docker Trusted Registry 1.1 also adds finer-grained permissions and
improvements and additions to the UI and logging. Bugs in LDAP/AD integration
have also been remediated, improving the stability and usability of Docker
Trusted Registry. See below for specifics.

**New Features**

* New, more granular, [roles for users](../user-management/index.md). Docker Trusted Registry users can now be assigned different levels of access
(admin, r/w, r/o) to the repositories. **Important:** Existing Docker Trusted
Registry users should make sure to see the note [below](prior-release-notes.md#dhe-1-0-upgrade-warning) regarding migrating users before upgrading.
* A new storage status indicator for storage space. The dashboard now shows used and available storage space for supported storage drivers.
* A new diagnostics tool gathers and bundles Docker Trusted Registry logs, system information, container
information, and other configuration settings for use by Docker support or as a
backup.
* Performance and reliability improvements to the S3 storage backend.
* Docker Trusted Registry images are now available on the Amazon AWS and Microsoft Azure marketplaces.

**Fixes**

The following notable issues have been remediated:

* Fixed an issue that caused Docker Trusted Registry logins to fail if some LDAP servers were unreachable.
* Fixed a resource leak in Docker Trusted Registry storage.

**DHE 1.0 Upgrade Warning**

If you currently use DHE 1.0, you **must** follow the
[upgrading instructions](https://forums.docker.com/t/upgrading-docker-hub-enterprise-to-docker-trusted-registry/1925)
in our support Knowledge Base. These instructions show you how to modify
existing authentication data and storage volume settings to move to Docker
Trusted Registry. Automatic upgrading has been disabled for DHE users because of
these issues.

## Version 1.0.1
(11 May 2015)

- Addresses compatibility issue with 1.6.1 CS Docker Engine

## Version 1.0.0
(23 Apr 2015)

- First release
