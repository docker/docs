+++
title = "Trusted Registry release notes"
description = "Docker Trusted Registry release notes "
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, release notes, Docker Trusted Registry"]
[menu.main]
parent="workw_dtr"
weight=100
+++


# Docker Trusted Registry release notes

This document describes the latest changes, additions, known issues, and fixes
for the Docker Trusted Registry.

## Prior versions

These notes refer to the current and immediately prior releases of Docker
Trusted Registry. For notes on older versions, see the [prior release notes archive](prior-release-notes.md).

# Docker Trusted Registry 1.4.3
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

* Previously, if you started the Trusted Registry 1.4.2 with CS Engine 1.7.0 onwards, it might not start because `docker` might start the Trusted Registry containers in an order that makes [links impossible to create](https://github.com/docker/docker/issues/17118). Using CS Engine 1.9 and later with the latest Trusted Registry includes creation of a custom network that allows all containers to connect to each other without links. This means that every time the Trusted Registry starts up, there should be no error.

    Also, when you upgrade CS Engine from 1.6 to 1.9 and the Trusted Registry admin server starts, it checks if it's running with links enabled. If that happens, the Trusted Registry restarts everything, creating the new network if necessary and removing the links, replacing them with a custom "dtr" network.

* This version deprecates the search API exposed at the /api/v0/index/search
endpoint. Starting on DTR 2.0 this API endpoint will no longer be available.

* On DTR 2.0 the autocomplete API exposed at the /api/v0/index/autocomplete
endpoint, will start returning a new result structure.

* On DTR 2.0 all API endpoints that return UserAccess and RepoUserAccess objects
are going to return objects with different properties.

# Docker Trusted Registry 1.4.2
(21 December 2015)

Release notes contain the following sections:

* Additional storage backend
* Fixed or updated with this release

## Additional storage backend
This release introduces using Openstack Swift as a storage backend. Refer to the [configuration documentation](configure/configuration.md) for details on the Swift driver.  

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

* New experimental feature. Docker Trusted Registry now integrates with Docker Content Trust using Notary. This is an experimental feature that is available with this release. See the [configuration documentation](configure/configuration.md).

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
