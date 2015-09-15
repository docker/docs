
+++
title = "Release notes"
description = "Release notes for Docker Trusted Registry"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry,  release"]
[menu.main]
parent="smn_dhe"
weight=100
+++



# Release Notes: Docker Trusted Registry & Commercially Supported Docker Engine

## Prior Versions

These notes refer to the current and immediately prior releases of Docker Trusted Registry  and the commercially supported Docker Engine. For notes on older versions of these, see the [prior release notes archive]({{< relref "prior-release-notes.md" >}}).

## Docker Trusted Registry

### Docker Trusted Registry 1.3.3
(18 September 2015)

This release corrects the following issues in Docker Trusted Registry 1.3.2

* Fixed an issue related to LDAP integration for users of Oracle Virtual Directory.
* Corrected an issue where Docker Trusted Registry would not accept a given certificate if the configured domain was only in the Subject Alternative Names (SANs) field and not in the Common Name (CN) field of the certificate.

### Docker Trusted Registry 1.3.2
(16 September 2015)

This release addresses the following change in Docker Trusted Registry 1.3.2 and is only available to customers who purchased DTR through Amazon Web Services (AWS) Marketplace.

* Docker Trusted Registry (DTR) now supports Amazon Web
Services (AWS) Integrated Billing. Previously, AWS users were required to
separately purchase a DTR license from Docker. AWS users can try DTR
out-of-the-box.

### Docker Trusted Registry 1.3.1
(31 August 2015)

This release corrects the following issues in Docker Trusted Registry 1.3.0

* The dashboard page was calculating incorrect stats.
* LDAP group sync failed to handle paginated results for extremely large groups.
* The repo delete endpoint returned incorrect error codes under certain conditions.

### Docker Trusted Registry 1.3.0
(26 August 2015)

This release addresses a few bugs and issues in Docker Trusted Registry 1.2.0 and introduces some new features and functionality, including:

* A completely new user-interface for the Admin application brings Docker Trusted Registry in line with other Docker products and provides greater ease-of-use.

* A new Accounts & Repos API provides new fine-grained role-based access control down to the per-repo level. See the [API's documentation](/docker-trusted-registry/api/) for more information.

* Improvements to the handling of configuration changes so that fewer restarts are required.

* Multiple security improvements and bug fixes.

### Docker Trusted Registry 1.2.0
(23 July 2015)

This release adds CentOS support and addresses a few bugs and issues in Docker Trusted Registry 1.1.0:

* Fixes an issue where for certain configurations of Docker Trusted Registry, proxy configuration settings and variables were not being passed to all Docker Trusted Registry containers and thus were not being respected.
* Documentation links in the UI now point to correct docs.
* Generated support info bundles have been scrubbed to remove highly sensitive data.
* Certifies support for CentOS 7.1.

## Commercially Supported Docker Engine

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
