+++
title ="Release Notes"
description="Docker Universal Control Plane"
[menu.main]
weight="-99"
+++


# UCP Release Notes

The latest release is 0.8.  Consult with your Docker sales engineer for the
release notes of earlier versions.

## Version 0.8

The following notes apply to this release:

### LDAP/AD integration

You can now configure UCP to use an LDAP or Active Directory service
for authentication.  When logged in with an admin account, go to the
Settings page and select "LDAP" from the Auth Method pull-down.

### DTR integration

You can now configure UCP to connect to a Docker Trusted Registry version
1.4.3 or newer.

### Teams and ACL

Teams can be set up to map to LDAP/AD groups, or managed entirely
within UCP.  Labels can then be set up on resources, and access can be
granted to those labels.

### Multi-host networking

The UCP bootstrapping tool now contains a utility for viewing and
configuring daemon configuraion.  After deploying your controllers
and replica nodes, you can enable multi-host networking with the
`engine-discovery` command.  For more usage information, run
`docker run --rm docker/ucp engine-discovery --help`

### UI

- Refined look and feel
- Teams UI
- LDAP/AD configuration UI
- Collapseable navigation bar


### Misc

- Now requires engine 1.10.0-rc1 or newer
- Etcd updated to 2.2.4
- Swarm 1.1.0-RC2
