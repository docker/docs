---
description: Archived release notes for commercially supported Docker Engine
keywords: docker, documentation, about, technology, understanding, enterprise, hub, registry, release, commercially supported Docker Engine
redirect_from:
- /docker-trusted-registry/cse-prior-release-notes/
- /docker-trusted-registry/cs-engine/release-notes/prior-release-notes/
- /cs-engine/release-notes/prior-release-notes/
title: Release notes archive for Commercially Supported Docker Engine.
---

This document contains the previous versions of the commercially supported
Docker Engine release notes. It includes issues, fixes, and new features.

Refer to the [detailed list](https://github.com/moby/moby/releases) of all changes since the release of CS Engine 1.10.3-cs3

## CS Engine 1.10.3-cs4
(12 Jan 2017)

Bumps RunC version to address CVE-2016-9962.

## CS Engine 1.10.3-cs3
(25 April 2016)

This release addresses the following issue:

A vulnerability in the Go standard runtime libraries allowed a maliciously crafted client certificate to be used to cause an infinite loop in a TLS server. This can lead to a Denial of Service against the Docker Engine if it is deployed such that it uses TLS client certificate authentication. This vulnerability has been fixed in this release. We consider this a low-impact issue, due to complexity of attack. Customers should consider upgrading if their deployed Docker Engines are exposed to potentially malicious network attackers.

This issue is resolved by using Go runtime v1.5.4 which was released to address this vulnerability

* https://github.com/moby/moby/pull/21977
* https://github.com/moby/moby/pull/21987

## CS Engine 1.10.3-cs2
(18 March 2016)

Bug fix release picking up changes from Docker 1.10.3 release.

Refer to the [detailed list](https://github.com/moby/moby/releases/tag/v1.10.3) of all changes since the release of CS Engine 1.10.2-cs1

## CS Engine 1.10.2-cs1
(22 February 2016)

In this release the CS Engine is supported on SUSE Linux Enterprise 12 OS.

Refer to the [detailed list](https://github.com/moby/moby/releases) of all changes since the release of CS Engine 1.9.1.

## CS Engine 1.9.1-cs3
(6 January 2016)

This release addresses the following issues:

* The commercially supported Engine 1.9.1-cs3 now supports multi-host networking
for all the kernels that the base CS Engine is supported on.

>**Note**: Centos 7 has its firewall enabled by default and it prevents the VXLAN tunnel from communicating. If this applies to you, then after installing the CS Engine, execute the following command in the Linux host:

    sudo firewall-cmd --zone=public --permanent --add-port=4789/udp


* Corrected an issue where Docker didn't remove the Masquerade NAT rule from `iptables` when the network was removed. This caused the gateway address to be
incorrectly propagated as the source address of a connection.

* Fixed an issue where if the daemon started multiple containers concurrently, then the `/etc/hosts` files were incompletely populated. This issue occurred randomly.

* Corrected an issue where the same IP address for different Docker containers resulted in network connection inconsistencies. Now each container has a separate IP address.

* Corrected an issue where the IPv6 gateway was not created when using custom networks although the network had a configured gateway.

* Fixed an issue where users might have experienced a panic error if the  daemon was started with the `—cluster-store` option, but without the `—cluster-advertise` option.

## CS Engine 1.9.1-cs2
(4 December 2015)

Starting with this release, upgrading minor versions, for example, from 1.9.0 to 1.9.1, is faster and easier.

You can refer to the detailed list of all changes since the release of CS Engine
1.9.0
https://github.com/moby/moby/releases.

## CS Engine 1.9.0
(12 November 2015)

Highlighted feature summary:

* Network Management and Plugins. Networks are now first class objects that can be listed, created, deleted, inspected, and connected to or disconnected from a
container. They can be manipulated outside of the container themselves and are
fully manageable on its own lifecycle. You can also use plugins to extend
network functionality.

* Docker, Inc. now provides support for the in-box Overlay (for cross-host networking) and Bridge network plugins. You can find more information about how
to manage networks and using network plugins in the [documentation](/engine/userguide/networking/index.md).

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

Refer to the [detailed list](https://github.com/moby/moby/releases) of all changes since the release of CS Engine 1.6.

## CS Engine 1.6.2-cs7
(12 October 2015)

As part of our ongoing security efforts, <a href="http://blog.docker.com/2015/10/security-release-docker-1-8-3-1-6-2-cs7" target="_blank">a vulnerability was discovered</a> that affects the way content
is stored and retrieved within the Docker Engine and CS Docker Engine. Today we
are releasing a security update that fixes this issue in both Docker Engine 1.8.3 and CS Docker Engine 1.6.2-cs7. The <a href="https://github.com/moby/moby/blob/master/CHANGELOG.md#161-2015-10-12" target="_blank">change log for Docker Engine 1.8.3</a> has a complete list of all the changes incorporated into both the open source and commercially
supported releases.

We recommend that users upgrade to CS Docker Engine 1.6.2-cs7. If you are unable
to upgrade to CS Docker Engine 1.6.2-cs7  right away, remember to only pull
content from trusted sources.

To keep up to date on all the latest Docker Security news, make sure you check
out our [Security page](http://www.docker.com/docker-security), subscribe to our mailing list, or find us in #docker-security.

## CS Docker Engine 1.6.2-cs6
(23 July 2015)

Certifies support for CentOS 7.1.

## CS Docker Engine 1.6.2-cs5
(21 May 2015)

For customers running Docker Engine on [supported versions of RedHat Enterprise Linux](https://www.docker.com/enterprise/support/) with SELinux enabled, the `docker build` and `docker run` commands will not have DNS host name resolution and bind-mounted volumes may not be accessible. As a result, customers with
SELinux will be unable to use hostname-based network access in either `docker build` or `docker run`, nor will they be able to `docker run` containers that use `--volume` or `-v` bind-mounts (with an incorrect SELinux label) in their environment. By installing Docker Engine 1.6.2-cs5, customers can use Docker as intended on RHEL with SELinux enabled.

For example, you see will failures such as:

```bash
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

```bash
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

**Affected Versions**: All previous versions of Docker Engine when SELinux is
 enabled.

Docker **highly recommends** that all customers running previous versions of Docker Engine update to this release.

### **How to workaround this issue**

Customers who choose not to install this update have two options. The first
option is to disable SELinux. This is *not recommended* for production systems
where SELinux is typically required.

The second option is to pass the following parameter in to `docker run`.

    --security-opt=label:type:docker_t

This parameter cannot be passed to the `docker build` command.

### **Upgrade notes**

When upgrading, make sure you stop Docker Trusted Registry first, perform the Engine upgrade, and then restart Docker Trusted Registry.

If you are running with SELinux enabled, previous Docker Engine releases allowed
you to bind-mount additional volumes or files inside the container as follows:

    $ docker run -it -v /home/user/foo.txt:/foobar.txt:ro <imagename>

In the 1.6.2-cs5 release, you must ensure additional bind-mounts have the
correct SELinux context. For example, if you want to mount `foobar.txt` as
read-only into the container, do the following to create and test your
bind-mount:

1.  Add the `z` option to the bind mount when you specify `docker run`.

    ```bash
		$ docker run -it -v /home/user/foo.txt:/foobar.txt:ro,z <imagename>
    ```

2.  Exec into your new container.

	  For example, if your container is `bashful_curie`, open a shell on the
	  container:

    ```bash
		$ docker exec -it bashful_curie bash
    ```

3.  Use `cat` to check the permissions on the mounted file.

    ```bash
		$ cat /foobar.txt
		the contents of foobar appear
    ```

If you see the file's contents, your mount succeeded. If you receive a
`Permission denied` message and/or the `/var/log/audit/audit.log` file on your
Docker host contains an AVC Denial message, the mount did not succeed.

    type=AVC msg=audit(1432145409.197:7570): avc:  denied  { read } for  pid=21167 comm="cat" name="foobar.txt" dev="xvda2" ino=17704136 scontext=system_u:system_r:svirt_lxc_net_t:s0:c909,c965 tcontext=unconfined_u:object_r:user_home_t:s0 tclass=file

Recheck your command line to make sure you passed in the `z` option.


## CS Engine 1.6.2-cs4
(13 May 2015)

Fix mount regression for `/sys`.

## CS Engine 1.6.1-cs3
(11 May 2015)

Docker Engine version 1.6.1 has been released to address several vulnerabilities
and is immediately available for all supported platforms. Users are advised to
upgrade existing installations of the Docker Engine and use 1.6.1 for new installations.

It should be noted that each of the vulnerabilities allowing privilege escalation
may only be exploited by a malicious Dockerfile or image.  Users are advised to
run their own images and/or images built by trusted parties, such as those in
the official images library.

Send any questions to security@docker.com.


### **[CVE-2015-3629](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2015-3629) Symlink traversal on container respawn allows local privilege escalation**

Libcontainer version 1.6.0 introduced changes which facilitated a mount namespace
breakout upon respawn of a container. This allowed malicious images to write
files to the host system and escape containerization.

Libcontainer and Docker Engine 1.6.1 have been released to address this
vulnerability. Users running untrusted images are encouraged to upgrade Docker Engine.

Discovered by Tõnis Tiigi.


### **[CVE-2015-3627](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2015-3627) Insecure opening of file-descriptor 1 leading to privilege escalation**

The file-descriptor passed by libcontainer to the pid-1 process of a container
has been found to be opened prior to performing the chroot, allowing insecure
open and symlink traversal. This allows malicious container images to trigger
a local privilege escalation.

Libcontainer and Docker Engine 1.6.1 have been released to address this
vulnerability. Users running untrusted images are encouraged to upgrade
Docker Engine.

Discovered by Tõnis Tiigi.

### **[CVE-2015-3630](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2015-3630) Read/write proc paths allow host modification & information disclosure**

Several paths underneath /proc were writable from containers, allowing global
system manipulation and configuration. These paths included `/proc/asound`,
`/proc/timer_stats`, `/proc/latency_stats`, and `/proc/fs`.

By allowing writes to `/proc/fs`, it has been noted that CIFS volumes could be
forced into a protocol downgrade attack by a root user operating inside of a
container. Machines having loaded the timer_stats module were vulnerable to
having this mechanism enabled and consumed by a container.

We are releasing Docker Engine 1.6.1 to address this vulnerability. All
versions up to 1.6.1 are believed vulnerable. Users running untrusted
images are encouraged to upgrade.

Discovered by Eric Windisch of the Docker Security Team.

### **[CVE-2015-3631](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2015-3631) Volume mounts allow LSM profile escalation**

By allowing volumes to override files of `/proc` within a mount namespace, a user
could specify arbitrary policies for Linux Security Modules, including setting
an unconfined policy underneath AppArmor, or a `docker_t` policy for processes
managed by SELinux. In all versions of Docker up until 1.6.1, it is possible for
malicious images to configure volume mounts such that files of proc may be overridden.

We are releasing Docker Engine 1.6.1 to address this vulnerability. All versions
up to 1.6.1 are believed vulnerable. Users running untrusted images are encouraged
to upgrade.

Discovered by Eric Windisch of the Docker Security Team.

### **AppArmor policy improvements**

The 1.6.1 release also marks preventative additions to the AppArmor policy.
Recently, several CVEs against the kernel have been reported whereby mount
namespaces could be circumvented through the use of the `sys_mount` syscall from
inside of an unprivileged Docker container. In all reported cases, the
AppArmor policy included in libcontainer and shipped with Docker has been
sufficient to deflect these attacks. However, we have deemed it prudent to
proactively tighten the policy further by outright denying the use of the
`sys_mount` syscall.

Because this addition is preventative, no CVE-ID is requested.

## CS Engine 1.6.0-cs2
(23 Apr 2015)

First release, see the [Docker Engine 1.6.0 Release notes](/v1.6/release-notes/)
  for more details.
