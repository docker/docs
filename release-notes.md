+++
title = "Release notes"
description = "Release notes for Docker Trusted Registry"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, release notes"]
[menu.main]
parent="smn_dhe"
weight=100
+++


# Release Notes: Docker Trusted Registry & Commercially Supported Docker Engine

This document describes the latest changes, additions, and fixes.

## Prior Versions

These notes refer to the current and immediately prior releases of Docker
Trusted Registry and the commercially supported Docker Engine. For notes on
older versions of these, see the [prior release notes archive](prior-release-notes.md).


## Docker Trusted Registry 1.4
(12 November 2015)

### New features
This release introduces the following new features. For additional information
on these features, refer to the documentation or when appropriate, the API
documentation.

* Image deletion and garbage collection

  * You can now delete an image in the registry's image index. This step of marking an unwanted image is called a soft delete. Refer to the [documentation](soft-garbage.md).

  * Administrators can use the dashboard or API to configure a task to regularly reclaim the disk space taken up by deleted images. Refer to the [documentation](soft-garbage.md).

* Repositories, Account Management, and interactive API UIs

  * Set up, and manage user accounts, teams, organizations, and repositories from either APIs or through the Trusted Registry user interface. Refer to either the API documentation or the [documentation](accounts.md) for performing tasks in the UI.

  * Search, browse, and discover images created by other users through either APIs or through the Trusted Registry user interface.

  * Users, depending on their roles, can access account information through the Trusted Registry user interface. Refer to the [documentation](accounts.md) for details.

  * View new API documentation through the Trusted Registry user interface. As before, you can also view it from the [documentation section](https://docs.docker.com/docker-trusted-registry/).

* New APIs

  * New APIs for accessing repositories, account management, indexing, searching, and reindexing.

  * You can also view an API and using Swagger UI, click the "Try it out" button to perform the action. This might be useful if you need to reindex.

* Different repository behavior. You must explicitly create (or have it performed for you if you don't have the  correct permissions) a repository before pushing to it. This behavior is different than how you would perform this in an unsecured free and open-source software (FOSS) registry.

* New experimental feature. Docker Trusted Registry now integrates with Docker Content Trust using Notary. This is an experimental feature that is available with this release. See the [configuration documentation](configuration.md).

### Fixed with this release
This release corrects the following issues in Docker Trusted Registry 1.3.3

#### LDAP Configuration

* Performance for LDAP user authenticaiton has been significantly increased, reducing the number of required LDAP requests to only a single BIND request to authenticate a user.

* The "Read-Write Search Filter" and "Read-Only Search Filter" fields have been deprecated. You can now create organization accounts and teams in the Trusted Registry to allow for more fine grained access control. Team member lists can be synced with a group in LDAP.

* An "Admin Password" is now required. Use this password to login as the
user admin in case the Trusted Registry is unable to authenticate you using
your LDAP server. This account can be used to login to the Trusted Registry and correct identity and authentication settings.

* Users on your LDAP server are now synced to the Trusted Registry's local
database using your configured "User Search Filter". Objects in LDAP that match
this filter and have a valid "User Login Attribute" are created as a local
user with the "User Login Attribute" as their username. Only these users are
able to login to Docker Trusted Registry.

* The "Admin LDAP DN" must now be specified to identify the group object on your LDAP server. This should be synced to the system administrators list. The "Admin
Group Member Attribute" should be set to the name of the attribute on this group
object which corresponds to the Distinguished Name of the group member objects.
This setting deprecates the old "Admin Search Filter" field.

#### Other Issues/Fixes

* Corrected an issue where you could switch from none to managed authentication
without creating an administrator.

* Added a "rootdirectory" parameter to the S3 storage option.

### Docker Trusted Registry 1.3.3
(18 September 2015) (amended: 2 November 2015)

This release corrects the following issues in Docker Trusted Registry 1.3.2

* Fixed an issue related to LDAP integration for users of Oracle Virtual Directory.

* Corrected an issue where Docker Trusted Registry would not accept a given certificate if the configured domain was only in the Subject Alternative Names (SANs) field and not in the Common Name (CN) field of the certificate.

* Docker discovered an issue in which the tokens used in authorization caused a
break in certain deployments that utilized a load balancer in front of multiple
Trusted Registry instances to achieve high availability. Docker regrets any
inconvenience this may have caused you and is working on a future fix.

## Commercially Supported Docker Engine

Commercially Supported (CS) Docker Engine is a packaged release that identifies
a release of Docker Engine for which you can receive support from Docker or one
of its partners. This release is functionally equivalent to the corresponding
Docker Engine release that it references. However, a CS release also includes
back-ported fixes (security-related and priority defects) from the open source.
It incorporates defect fixes that you can use in environments where new features
cannot be adopted as quickly for consistency and compatibility reasons.  

### CS Docker Engine 1.9.0
(12 November 2015)

Highlighted Feature Summary:

* Network Management and Plugins: Networks are now first class objects that can
 be listed, created, deleted, inspected and connected to or disconnected from a
 container. They can be manipulated outside of the container themselves and are
 fully manageable on its own lifecycle. Network functionality can also be
 extended using plugins.

* Docker now provides support for the in-box Overlay (for cross-host networking)
 and Bridge network plugins. You can find more information about how to manage
 networks and using  network plugins in the documentation.

* Volume Management and Plugins: Volumes also become discrete, manageable objects in Docker. Volumes can be listed, created, deleted, and inspected.
Similar to networks, they have their own managed lifecycle outside of the
container. Plugins allow others to write and extend the functionality of volumes
or provide integration with other types of storage.

* The in-box volume driver is included and supported. You can find more information about how to manage volumes  and using  volume plugins in the
documentation.

* Docker Content Trust: Content trust gives you the ability to both verify the integrity and the publisher of all the data received from a registry over any channel. Content Trust is currently only supported using Docker Hub notary servers.

* Updated the release cadence of CS Docker Engine. Starting with this version, Docker supports **every** major release of Docker Engine from open
source with three releases under support at one time. This means you’ll be able
to take advantage of the latest and greatest features and you won’t have to wait
for a supported release to take advantage of a specific feature.

Refer to the detailed list of all changes since the release of CS Engine 1.6.
https://github.com/docker/docker/releases.

### CS Docker Engine 1.6.2-cs7
(12 October 2015)

As part of our ongoing security efforts, <a href="http://blog.docker.com/2015/10/security-release-docker-1-8-3-1-6-2-cs7" target="_blank">a vulnerability was discovered</a> that
affects the way content is stored and retrieved within the Docker Engine and CS
Docker Engine. Today we are releasing a security update that fixes this
issue in both Docker Engine 1.8.3 and CS Docker Engine 1.6.2-cs7. The <a
href="https://github.com/docker/docker/blob/master/CHANGELOG.md#161-2015-10-12"
target="_blank">change log for Docker Engine 1.8.3</a> has a complete list of
all the changes incorporated into both the open source and commercially
supported releases.

We recommend that users upgrade to CS Docker Engine 1.6.2-cs7.
If you are unable to upgrade to CS Docker Engine 1.6.2-cs7  right away, remember to only pull content from trusted sources.

To keep up to date on all the latest Docker Security news, make sure you check
out our [Security page](http://www.docker.com/docker-security), subscribe to our
mailing list, or find us in #docker-security.

### CS Docker Engine 1.6.2-cs6
(23 July 2015)

* Certifies support for CentOS 7.1.

### CS Docker Engine 1.6.2-cs5
(21 May 2015)

For customers running Docker Engine on [supported versions of RedHat Enterprise
Linux](https://www.docker.com/enterprise/support/) with [SELinux
enabled](
https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/7/html/SELinux_Users_and_Administrators_Guide/sect-Security-Enhanced_Linux-Working_with_SELinux-Enabling_and_Disabling_SELinux.html
), the `docker build` and `docker run`
commands will not have DNS host name resolution and bind-mounted volumes may
not be accessible.
As a result, customers with SELinux will be unable to use hostname-based network
access in either `docker build` or `docker run`, nor will they be able to
`docker run` containers
that use `--volume` or `-v` bind-mounts (with an incorrect SELinux label) in
their environment. By installing Docker Engine 1.6.2-cs5, customers can use Docker as intended on RHEL with SELinux enabled.

For example, you see will failures like:

```
[root@dtr ~]# docker -v
Docker version 1.6.0-cs2, build b8dd430
[root@dtr ~]# ping dtr.home.org.au
PING dtr.home.org.au (10.10.10.104) 56(84) bytes of data.
64 bytes from dtr.home.gateway (10.10.10.104): icmp_seq=1 ttl=64 time=0.663 ms
^C
--- dtr.home.org.au ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1001ms
rtt min/avg/max/mdev = 0.078/0.370/0.663/0.293 ms
[root@dtr ~]# docker run --rm -it debian ping dtr.home.org.au
ping: unknown host
[root@dtr ~]# docker run --rm -it debian cat /etc/resolv.conf
cat: /etc/resolv.conf: Permission denied
[root@dtr ~]# docker run --rm -it debian apt-get update
Err http://httpredir.debian.org jessie InRelease

Err http://security.debian.org jessie/updates InRelease

Err http://httpredir.debian.org jessie-updates InRelease

Err http://security.debian.org jessie/updates Release.gpg
  Could not resolve 'security.debian.org'
Err http://httpredir.debian.org jessie Release.gpg
  Could not resolve 'httpredir.debian.org'
Err http://httpredir.debian.org jessie-updates Release.gpg
  Could not resolve 'httpredir.debian.org'
[output truncated]

```

or when running a `docker build`:

```
[root@dtr ~]# docker build .
Sending build context to Docker daemon 11.26 kB
Sending build context to Docker daemon
Step 0 : FROM fedora
 ---> e26efd418c48
Step 1 : RUN yum install httpd
 ---> Running in cf274900ea35

One of the configured repositories failed (Fedora 21 - x86_64),
and yum doesn't have enough cached data to continue. At this point the only
safe thing yum can do is fail. There are a few ways to work "fix" this:

[output truncated]
```

**Affected Versions**: All previous versions of Docker Engine when SELinux
is enabled.

Docker **highly recommends** that all customers running previous versions of
Docker Engine update to this release.

#### **How to workaround this issue**

Customers who choose not to install this update have two options. The
first option is to disable SELinux. This is *not recommended* for production
systems where SELinux is typically required.

The second option is to pass the following parameter in to `docker run`.

  	     --security-opt=label:type:docker_t

This parameter cannot be passed to the `docker build` command.

#### **Upgrade notes**

When upgrading, make sure you stop Docker Trusted Registry first, perform the Engine upgrade, and
then restart Docker Trusted Registry.

If you are running with SELinux enabled, previous Docker Engine releases allowed
you to bind-mount additional volumes or files inside the container as follows:

		$ docker run -it -v /home/user/foo.txt:/foobar.txt:ro <imagename>

In the 1.6.2-cs5 release, you must ensure additional bind-mounts have the correct
SELinux context. For example, if you want to mount `foobar.txt` as read-only
into the container, do the following to create and test your bind-mount:

1. Add the `z` option to the bind mount when you specify `docker run`.

		$ docker run -it -v /home/user/foo.txt:/foobar.txt:ro,z <imagename>

2. Exec into your new container.

	For example, if your container is `bashful_curie`, open a shell on the
	container:

		$ docker exec -it bashful_curie bash

3. Use `cat` to check the permissions on the mounted file.

		$ cat /foobar.txt
		the contents of foobar appear

	If you see the file's contents, your mount succeeded. If you receive a
	`Permission denied` message and/or the `/var/log/audit/audit.log` file on
	your Docker host contains an AVC Denial message, the mount did not succeed.

		type=AVC msg=audit(1432145409.197:7570): avc:  denied  { read } for  pid=21167 comm="cat" name="foobar.txt" dev="xvda2" ino=17704136 scontext=system_u:system_r:svirt_lxc_net_t:s0:c909,c965 tcontext=unconfined_u:object_r:user_home_t:s0 tclass=file

	Recheck your command line to make sure you passed in the `z` option.
