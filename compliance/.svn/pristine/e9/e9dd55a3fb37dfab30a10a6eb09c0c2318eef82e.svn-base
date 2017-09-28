---
title: "Identification and authentication"
description: "Identification and authentication reference"
keywords: "standards, compliance, security, 800-53, Identification and authentication"
---

## IA-1 Identification And Authentication Policy And Procedures

#### Description

The organization:
<ol type="a">
<li>Develops, documents, and disseminates to [Assignment: organization-defined personnel or roles]:</li>

<ol type="1">
<li>An identification and authentication policy that addresses purpose, scope, roles, responsibilities, management commitment, coordination among organizational entities, and compliance; and</li>
<li>Procedures to facilitate the implementation of the identification and authentication policy and associated identification and authentication controls; and</li>
</ol>
<li>Reviews and updates the current:</li>

<ol type="1">
<li>Identification and authentication policy [Assignment: organization-defined frequency]; and</li>
<li>Identification and authentication procedures [Assignment: organization-defined frequency].</li>
</ol>
</ol>

#### Control Information

Responsible role(s) - Organization

## IA-2 Identification And Authentication (Organizational Users)

#### Description

The information system uniquely identifies and authenticates organizational users (or processes acting on behalf of organizational users).

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko2u0">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko2u0" class="tab-pane fade in active">
Docker Enterprise Edition can be configured to identify and
authenticate users via it&#39;s integrated support for LDAP. Users and
groups managed within the organization&#39;s LDAP directory service (e.g.
Active Directory) can be synchronized to UCP and DTR on a regular
interval. When a user is removed from the LDAP-backed directory, that
user becomes inactive within UCP and DTR. In addition, UCP and DTR
teams can be mapped to groups synchronized via LDAP. When a user is
added/removed to/from the LDAP group, that same user is automatically
added/removed to/from the UCP and DTR team. Additional information can
be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/external-auth/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/external-auth/</a></li>
</ul>

</div>
</div>

### IA-2 (1) Network Access To Privileged Accounts

#### Description

The information system implements multifactor authentication for network access to privileged accounts.

#### Control Information

Responsible role(s) - Organization

### IA-2 (2) Network Access To Non-Privileged Accounts

#### Description

The information system implements multifactor authentication for network access to non-privileged accounts.

#### Control Information

Responsible role(s) - Organization

### IA-2 (3) Local Access To Privileged Accounts

#### Description

The information system implements multifactor authentication for local access to privileged accounts.

#### Control Information

Responsible role(s) - Organization

### IA-2 (4) Local Access To Non-Privileged Accounts

#### Description

The information system implements multifactor authentication for local access to non-privileged accounts.

#### Control Information

Responsible role(s) - Organization

### IA-2 (5) Group Authentication

#### Description

The organization requires individuals to be authenticated with an individual authenticator when a group authenticator is employed.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Trusted Registry (DTR)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko2ug">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu89idlmg00bko2v0">UCP</a></li>
<li><a data-toggle="tab" data-target="#b70lu89idlmg00bko2vg">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko2ug" class="tab-pane fade in active">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, Docker Trusted
Registry requires individual users to be authenticated in order to
gain access to the system. Any permissions granted to the team(s) that
which the user is a member are subsequently applied.
</div>
<div id="b70lu89idlmg00bko2v0" class="tab-pane fade">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, Universal Control
Plane requires individual users to be authenticated in order to gain
access to the system. Any permissions granted to the team(s) that
which the user is a member are subsequently applied.
</div>
<div id="b70lu89idlmg00bko2vg" class="tab-pane fade">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, Docker Enterprise
Edition requires individual users to be authenticated in order to gain
access to the system. Any permissions granted to the team(s) that
which the user is a member are subsequently applied.
</div>
</div>

### IA-2 (6) Network Access To Privileged Accounts - Separate Device

#### Description

The information system implements multifactor authentication for network access to privileged accounts such that one of the factors is provided by a device separate from the system gaining access and the device meets [Assignment: organization-defined strength of mechanism requirements].

#### Control Information

Responsible role(s) - Organization

### IA-2 (7) Network Access To Non-Privileged Accounts - Separate Device

#### Description

The information system implements multifactor authentication for network access to non-privileged accounts such that one of the factors is provided by a device separate from the system gaining access and the device meets [Assignment: organization-defined strength of mechanism requirements].

#### Control Information

Responsible role(s) - Organization

### IA-2 (8) Network Access To Privileged Accounts - Replay Resistant

#### Description

The information system implements replay-resistant authentication mechanisms for network access to privileged accounts.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko300">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko300" class="tab-pane fade in active">
Docker Enterprise Edition integrates with LDAP for authenticating
users to an external directory service. You should configure your
external directory service for ensuring that you are protected against
replay attacks.
</div>
</div>

### IA-2 (9) Network Access To Non-Privileged Accounts - Replay Resistant

#### Description

The information system implements replay-resistant authentication mechanisms for network access to non-privileged accounts.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko30g">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko30g" class="tab-pane fade in active">
Docker Enterprise Edition integrates with LDAP for authenticating
users to an external directory service. You should configure your
external directory service for ensuring that you are protected against
replay attacks.
</div>
</div>

### IA-2 (10) Single Sign-On

#### Description

The information system provides a single sign-on capability for [Assignment: organization-defined information system accounts and services].

#### Control Information

Responsible role(s) - Organization

### IA-2 (11) Remote Access  - Separate Device

#### Description

The information system implements multifactor authentication for remote access to privileged and non-privileged accounts such that one of the factors is provided by a device separate from the system gaining access and the device meets [Assignment: organization-defined strength of mechanism requirements].

#### Control Information

Responsible role(s) - Organization

### IA-2 (12) Acceptance Of Piv Credentials

#### Description

The information system accepts and electronically verifies Personal Identity Verification (PIV) credentials.

#### Control Information

Responsible role(s) - Organization

### IA-2 (13) Out-Of-Band Authentication

#### Description

The information system implements [Assignment: organization-defined out-of-band authentication] under [Assignment: organization-defined conditions].

#### Control Information

Responsible role(s) - Organization

## IA-3 Device Identification And Authentication

#### Description

The information system uniquely identifies and authenticates [Assignment: organization-defined specific and/or types of devices] before establishing a [Selection (one or more): local; remote; network] connection.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Trusted Registry (DTR)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko310">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu89idlmg00bko31g">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu89idlmg00bko320">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko310" class="tab-pane fade in active">
Docker Trusted Registry replicas reside on Universal Control Plane
worker nodes. In order for UCP worker nodes to join a Universal
Control Plane cluster, they must be identified and authenticated via a
worker token. Additional Docker Trusted Registry replicas can only be
added after a UCP administrator user has authenticated in to the UCP
cluster and when mutual TLS authentication between the UCP worker and
manager nodes has been established. Additional information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/dtr/2.3/guides/admin/install/#step-7-join-replicas-to-the-cluster">https://docs.docker.com/datacenter/dtr/2.3/guides/admin/install/#step-7-join-replicas-to-the-cluster</a></li>
</ul>

</div>
<div id="b70lu89idlmg00bko31g" class="tab-pane fade">
In order for other Docker EE engine nodes to be able to join a
cluster managed by Universal Control Plane, they must be identified
and authenticated via either a manager or worker token. Use of the
token includes trust on first use mutual TLS.
</div>
<div id="b70lu89idlmg00bko320" class="tab-pane fade">
In order for nodes to join a Universal Control Plane cluster, they
must be identified and authenticated via either a manager or worker
token. Additional information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/scale-your-cluster/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/scale-your-cluster/</a></li>
</ul>

</div>
</div>

### IA-3 (1) Cryptographic Bidirectional Authentication

#### Description

The information system authenticates [Assignment: organization-defined specific devices and/or types of devices] before establishing [Selection (one or more): local; remote; network] connection using bidirectional authentication that is cryptographically based.

#### Control Information

Responsible role(s) - Organization

### IA-3 (3) Dynamic Address Allocation

#### Description

The organization:
<ol type="a">
<li>Standardizes dynamic address allocation lease information and the lease duration assigned to devices in accordance with [Assignment: organization-defined lease information and lease duration]; and</li>
<li>Audits lease information when assigned to a device.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### IA-3 (4) Device Attestation

#### Description

The organization ensures that device identification and authentication based on attestation is handled by [Assignment: organization-defined configuration management process].

#### Control Information

Responsible role(s) - Organization

## IA-4 Identifier Management

#### Description

The organization manages information system identifiers by:
<ol type="a">
<li>Receiving authorization from [Assignment: organization-defined personnel or roles] to assign an individual, group, role, or device identifier;</li>
<li>Selecting an identifier that identifies an individual, group, role, or device;</li>
<li>Assigning the identifier to the intended individual, group, role, or device;</li>
<li>Preventing reuse of identifiers for [Assignment: organization-defined time period]; and</li>
<li>Disabling the identifier after [Assignment: organization-defined time period of inactivity].</li>
</ol>

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko32g">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko32g" class="tab-pane fade in active">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to prevent the reuse of user identifiers for a
specified period of time. Refer to your directory service&#39;s
documentation for configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to prevent the reuse of user identifiers for a
specified period of time. Refer to your directory service&#39;s
documentation for configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to prevent the reuse of user identifiers for a
specified period of time. Refer to your directory service&#39;s
documentation for configuring this.
</div>
</div>

### IA-4 (1) Prohibit Account Identifiers As Public Identifiers

#### Description

The organization prohibits the use of information system account identifiers that are the same as public identifiers for individual electronic mail accounts.

#### Control Information

Responsible role(s) - Organization

### IA-4 (2) Supervisor Authorization

#### Description

The organization requires that the registration process to receive an individual identifier includes supervisor authorization.

#### Control Information

Responsible role(s) - Organization

### IA-4 (3) Multiple Forms Of Certification

#### Description

The organization requires multiple forms of certification of individual identification be presented to the registration authority.

#### Control Information

Responsible role(s) - Organization

### IA-4 (4) Identify User Status

#### Description

The organization manages individual identifiers by uniquely identifying each individual as [Assignment: organization-defined characteristic identifying individual status].

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko330">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko330" class="tab-pane fade in active">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to uniquely identify each individual according to
the requirements of this control. Refer to your directory service&#39;s
documentation for configuring this.
</div>
</div>

### IA-4 (5) Dynamic Management

#### Description

The information system dynamically manages identifiers.

#### Control Information

Responsible role(s) - Organization

### IA-4 (6) Cross-Organization Management

#### Description

The organization coordinates with [Assignment: organization-defined external organizations] for cross-organization management of identifiers.

#### Control Information

Responsible role(s) - Organization

### IA-4 (7) In-Person Registration

#### Description

The organization requires that the registration process to receive an individual identifier be conducted in person before a designated registration authority.

#### Control Information

Responsible role(s) - Organization

## IA-5 Authenticator Management

#### Description

The organization manages information system authenticators by:
<ol type="a">
<li>Verifying, as part of the initial authenticator distribution, the identity of the individual, group, role, or device receiving the authenticator;</li>
<li>Establishing initial authenticator content for authenticators defined by the organization;</li>
<li>Ensuring that authenticators have sufficient strength of mechanism for their intended use;</li>
<li>Establishing and implementing administrative procedures for initial authenticator distribution, for lost/compromised or damaged authenticators, and for revoking authenticators;</li>
<li>Changing default content of authenticators prior to information system installation;</li>
<li>Establishing minimum and maximum lifetime restrictions and reuse conditions for authenticators;</li>
<li>Changing/refreshing authenticators [Assignment: organization-defined time period by authenticator type];</li>
<li>Protecting authenticator content from unauthorized disclosure and modification;</li>
<li>Requiring individuals to take, and having devices implement, specific security safeguards to protect authenticators; and</li>
<li>Changing authenticators for group/role accounts when membership to those accounts changes.</li>
</ol>

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko33g">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko33g" class="tab-pane fade in active">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to establish initial authenticator content according
to the requirements of this control. Refer to your directory service&#39;s
documentation for configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to enforce strength requirements for authenticators
according to the requirements of this control. Refer to your directory
service&#39;s documentation for configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to distribute, redistribute, and revoke
authenticators according to the requirements of this control. Refer to
your directory service&#39;s documentation for configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to change default authenticator content according to
the requirements of this control. Refer to your directory service&#39;s
documentation for configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to set minimum and maximum lifetime restrictions and
reuse conditions for authenticators according to the requirements of
this control. Refer to your directory service&#39;s documentation for
configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to refresh authenticators at a regular cadence
according to the requirements of this control. Refer to your directory
service&#39;s documentation for configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to protect authenticator content from unauthorized
disclosure or modification according to the requirements of this
control. Refer to your directory service&#39;s documentation for
configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to implement specific security safeguards to protect
authentications according to the requirements of this control. Refer
to your directory service&#39;s documentation for configuring this.The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to change authenticators for group or role accounts
when membership to those groups or roles changes  according to the
requirements of this control. Refer to your directory service&#39;s
documentation for configuring this.
</div>
</div>

### IA-5 (1) Password-Based Authentication

#### Description

The information system, for password-based authentication:
<ol type="a">
<li>Enforces minimum password complexity of [Assignment: organization-defined requirements for case sensitivity, number of characters, mix of upper-case letters, lower-case letters, numbers, and special characters, including minimum requirements for each type];</li>
<li>Enforces at least the following number of changed characters when new passwords are created: [Assignment: organization-defined number];</li>
<li>Stores and transmits only cryptographically-protected passwords;</li>
<li>Enforces password minimum and maximum lifetime restrictions of [Assignment: organization-defined numbers for lifetime minimum, lifetime maximum];</li>
<li>Prohibits password reuse for [Assignment: organization-defined number] generations; and</li>
<li>Allows the use of a temporary password for system logons with an immediate change to a permanent password.</li>
</ol>

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko340">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko340" class="tab-pane fade in active">
An external directory service integrated with Docker Enterprise
Edition via LDAP can be configured to enforce minimum password
complexity requirements. Refer to your directory service&#39;s
documentation for configuring this.An external directory service integrated with Docker Enterprise
Edition via LDAP can be configured to enforce the requirement to
change at least one character when changing passwords according to the
requirements of this control. Refer to your directory service&#39;s
documentation for configuring this.An external directory service integrated with Docker Enterprise
Edition via LDAP can be configured to store and transmit
cryptographically protected passwords according to the requirements of
this control. Refer to your directory service&#39;s documentation for
configuring this.An external directory service integrated with Docker Enterprise
Edition via LDAP can be configured to enforce the required minimum and
maximum lifetime restrictions according to the requirements of this
control. Refer to your directory service&#39;s documentation for
configuring this.An external directory service integrated with Docker Enterprise
Edition via LDAP can be configured to enforce the required number of
generations before password reuse according to the requirements of
this control. Refer to your directory service&#39;s documentation for
configuring this.An external directory service integrated with Docker Enterprise
Edition via LDAP can be configured to enforce the requirement to
change initial/temporary passwords upon first login according to the
requirements of this control. Refer to your directory service&#39;s
documentation for configuring this.
</div>
</div>

### IA-5 (2) Pki-Based Authentication

#### Description

The information system, for PKI-based authentication:
<ol type="a">
<li>Validates certifications by constructing and verifying a certification path to an accepted trust anchor including checking certificate status information;</li>
<li>Enforces authorized access to the corresponding private key;</li>
<li>Maps the authenticated identity to the account of the individual or group; and</li>
<li>Implements a local cache of revocation data to support path discovery and validation in case of inability to access revocation information via the network.</li>
</ol>

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Trusted Registry (DTR)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko34g">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu89idlmg00bko350">UCP</a></li>
<li><a data-toggle="tab" data-target="#b70lu89idlmg00bko35g">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko34g" class="tab-pane fade in active">
Docker Trusted Registry includes a Docker volume which holds the root
key material for the DTR root CA that issues certificats. In addition
Universal Control Plane contains two, built-in root certificate
authorities. One CA is used for signing client bundles generated by
users. The other CA is used for TLS communication between UCP cluster
nodes. Should you choose to use certificates signed by an external CA,
in order to successfully authenticate in to the system, those
certificates must include a root CA public certificate, a service
certificate and any intermediate CA public certificates (in addition
to SANs for all addresses used to reach the UCP controller), and a
private key for the server. When adding DTR replicas, the UCP nodes on
which they&#39;re installed are authenticated to the cluster via the
appropriate built-in CA.Access to Docker Trusted Registry is only granted when a user has a
valid certificate bundle. This is enforced with the public/private key
pair included with the user&#39;s certificate bundle in Universal Control
Plane.Only after a client bundle has been generated or an existing public
key has been added for a particular user is that user able to execute
commands against Docker Trusted Registry. This bundle maps the
authenticated identity to that of the user&#39;s profile in Universal
Control Plane.When a client bundle has been generated or an existing public key has
been added for a particular Universal Control Plane user which
subsequently grants that user access to Docker Trusted Registry, it is
attached to that user&#39;s Universal Control Plane profile. Bundles/keys
can be revoked by an Administrator or the user themselves. The
cluster&#39;s internal certificates can also be revoked and updated.
Additional information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/use-your-own-tls-certificates/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/use-your-own-tls-certificates/</a></li>
<li><a href="https://docs.docker.com/datacenter/dtr/2.3/guides/admin/configure/use-your-own-tls-certificates/">https://docs.docker.com/datacenter/dtr/2.3/guides/admin/configure/use-your-own-tls-certificates/</a></li>
</ul>

</div>
<div id="b70lu89idlmg00bko350" class="tab-pane fade">
Universal Control Plane contains two, built-in root certificate
authorities. One CA is used for signing client bundles generated by
users. The other CA is used for TLS communication between UCP cluster
nodes. Should you choose to use certificates signed by an external CA,
in order to successfully authenticate in to the system, those
certificates must include a root CA public certificate, a service
certificate and any intermediate CA public certificates (in addition
to SANs for all addresses used to reach the UCP controller), and a
private key for the server.Access to a Universal Control Plane cluster is only granted when a
user has a valid certificate bundle. This is enforced with the
public/private key pair included with the user&#39;s certificate bundle.Only after a client bundle has been generated or an existing public
key has been added for a particular user is that user able to execute
commands against the Universal Control Plane cluster. This bundle maps
the authenticated identity to that of the user.When a client bundle has been generated or an existing public key has
been added for a particular Universal Control Plane user, it is
attached to that user&#39;s profile. Bundles/keys can be revoked by an
Administrator or the user themselves. The cluster&#39;s internal
certificates can also be revoked and updated. Additional information
can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/use-your-own-tls-certificates/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/use-your-own-tls-certificates/</a></li>
</ul>

</div>
<div id="b70lu89idlmg00bko35g" class="tab-pane fade">
All users within a Docker Enterprise Edition cluster can create a
client certificate bundle for authenticating in to the cluster from
the Docker client tooling. When a user attempts to authenticate in to
the Docker cluster, the system validates the certificates per the
requirements of this control.All users within a Docker Enterprise Edition cluster can create a
client certificate bundle for authenticating in to the cluster from
the Docker client tooling. When a user attempts to authenticate in to
the Docker cluster, the system enforces authorized access to the
corresponding private key per the requirements of this control.All users within a Docker Enterprise Edition cluster can create a
client certificate bundle for authenticating in to the cluster from
the Docker client tooling. When a user attempts to authenticate in to
the Docker cluster, the system maps the authenticated identity to the
account of the individual or group per the requirements of this
control.All users within a Docker Enterprise Edition cluster can create a
client certificate bundle for authenticating in to the cluster from
the Docker client tooling. When a user attempts to authenticate in to
the Docker cluster, it is up to the underlying operating system
hosting Docker Enterprise Edition to ensure that it implements a local
cache of revocation data per the requirements of this control.
</div>
</div>

### IA-5 (3) In-Person Or Trusted Third-Party Registration

#### Description

The organization requires that the registration process to receive [Assignment: organization-defined types of and/or specific authenticators] be conducted [Selection: in person; by a trusted third party] before [Assignment: organization-defined registration authority] with authorization by [Assignment: organization-defined personnel or roles].

#### Control Information

Responsible role(s) - Organization

### IA-5 (4) Automated Support  For Password Strength Determination

#### Description

The organization employs automated tools to determine if password authenticators are sufficiently strong to satisfy [Assignment: organization-defined requirements].

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko360">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko360" class="tab-pane fade in active">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP can be
configured with automation to ensure that password authenticators meet
strength requirements as defined by this control. Refer to your
directory service&#39;s documentation for configuring this.
</div>
</div>

### IA-5 (5) Change Authenticators Prior To Delivery

#### Description

The organization requires developers/installers of information system components to provide unique authenticators or change default authenticators prior to delivery/installation.

#### Control Information

Responsible role(s) - Organization

### IA-5 (6) Protection Of Authenticators

#### Description

The organization protects authenticators commensurate with the security category of the information to which use of the authenticator permits access.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko36g">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko36g" class="tab-pane fade in active">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to protect authenticators as required by this
control. Refer to your directory service&#39;s documentation for
configuring this.
</div>
</div>

### IA-5 (7) No Embedded Unencrypted Static Authenticators

#### Description

The organization ensures that unencrypted static authenticators are not embedded in applications or access scripts or stored on function keys.

#### Control Information

Responsible role(s) - Organization

### IA-5 (8) Multiple Information System Accounts

#### Description

The organization implements [Assignment: organization-defined security safeguards] to manage the risk of compromise due to individuals having accounts on multiple information systems.

#### Control Information

Responsible role(s) - Organization

### IA-5 (9) Cross-Organization Credential Management

#### Description

The organization coordinates with [Assignment: organization-defined external organizations] for cross-organization management of credentials.

#### Control Information

Responsible role(s) - Organization

### IA-5 (10) Dynamic Credential Association

#### Description

The information system dynamically provisions identities.

#### Control Information

Responsible role(s) - Organization

### IA-5 (11) Hardware Token-Based Authentication

#### Description

The information system, for hardware token-based authentication, employs mechanisms that satisfy [Assignment: organization-defined token quality requirements].

#### Control Information

Responsible role(s) - Organization

### IA-5 (12) Biometric-Based Authentication

#### Description

The information system, for biometric-based authentication, employs mechanisms that satisfy [Assignment: organization-defined biometric quality requirements].

#### Control Information

Responsible role(s) - Organization

### IA-5 (13) Expiration Of Cached Authenticators

#### Description

The information system prohibits the use of cached authenticators after [Assignment: organization-defined time period].

#### Control Information

Responsible role(s) - Organization

### IA-5 (14) Managing Content Of Pki Trust Stores

#### Description

The organization, for PKI-based authentication, employs a deliberate organization-wide methodology for managing the content of PKI trust stores installed across all platforms including networks, operating systems, browsers, and applications.

#### Control Information

Responsible role(s) - Organization

### IA-5 (15) Ficam-Approved Products And Services

#### Description

The organization uses only FICAM-approved path discovery and validation products and services.

#### Control Information

Responsible role(s) - Organization

## IA-6 Authenticator Feedback

#### Description

The information system obscures feedback of authentication information during the authentication process to protect the information from possible exploitation/use by unauthorized individuals.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Trusted Registry (DTR)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko370">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu89idlmg00bko37g">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko370" class="tab-pane fade in active">
Docker Trusted Registry obscures all feedback of authentication
information during the authentication process. This includes both
authentication via the web UI and the CLI.
</div>
<div id="b70lu89idlmg00bko37g" class="tab-pane fade">
Universal Control Plane obscures all feedback of authentication
information during the authentication process. This includes both
authentication via the web UI and the CLI.
</div>
</div>

## IA-7 Cryptographic Module Authentication

#### Description

The information system implements mechanisms for authentication to a cryptographic module that meet the requirements of applicable federal laws, Executive Orders, directives, policies, regulations, standards, and guidance for such authentication.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Trusted Registry (DTR)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko380">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu89idlmg00bko38g">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko380" class="tab-pane fade in active">
All access to Docker Trusted Registry is protected with Transport
Layer Security (TLS) 1.2 with the AES-GCM cipher. This includes both
SSH access to the individual UCP nodes and CLI-/web-based access to
the UCP management functions with mutual TLS and HTTPS respectively.
</div>
<div id="b70lu89idlmg00bko38g" class="tab-pane fade">
All access to Universal Control Plane is protected with Transport
Layer Security (TLS) 1.2 with the AES GCM cipher. This includes both
SSH access to the individual UCP nodes and CLI-/web-based access to
the UCP management functions with mutual TLS and HTTPS respectively.
</div>
</div>

## IA-8 Identification And Authentication (Non-Organizational Users)

#### Description

The information system uniquely identifies and authenticates non-organizational users (or processes acting on behalf of non-organizational users).

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Trusted Registry (DTR)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko390">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu89idlmg00bko39g">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko390" class="tab-pane fade in active">
Users managed by Docker Trusted Registry can be grouped per the
requirements of the organization and as defined by this control. This
can include groupings for non-organizational users.
</div>
<div id="b70lu89idlmg00bko39g" class="tab-pane fade">
Users managed by Universal Control Plane can be grouped per the
requirements of the organization and as defined by this control. This
can include groupings for non-organizational users.
</div>
</div>

### IA-8 (1) Acceptance Of Piv Credentials From Other Agencies

#### Description

The information system accepts and electronically verifies Personal Identity Verification (PIV) credentials from other federal agencies.

#### Control Information

Responsible role(s) - Organization

### IA-8 (2) Acceptance Of Third-Party Credentials

#### Description

The information system accepts only FICAM-approved third-party credentials.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko3a0">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko3a0" class="tab-pane fade in active">
An external directory service integrated with Docker Enterprise
Edition via LDAP can be configured to meet the FICAM requirements as
indicated by this control. Refer to your directory service&#39;s
documentation for configuring this.
</div>
</div>

### IA-8 (3) Use Of Ficam-Approved Products

#### Description

The organization employs only FICAM-approved information system components in [Assignment: organization-defined information systems] to accept third-party credentials.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko3ag">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko3ag" class="tab-pane fade in active">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to meet the FICAM requirements as indicated by this
control. Refer to your directory service&#39;s documentation for
configuring this.
</div>
</div>

### IA-8 (4) Use Of Ficam-Issued Profiles

#### Description

The information system conforms to FICAM-issued profiles.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu89idlmg00bko3b0">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu89idlmg00bko3b0" class="tab-pane fade in active">
The organization is responsible for meeting the requirements of this
control. To assist with meeting these requirements, an external
directory service integrated with Docker Enterprise Edition via LDAP
can be configured to meet the FICAM requirements as indicated by this
control. Refer to your directory service&#39;s documentation for
configuring this.
</div>
</div>

### IA-8 (5) Acceptance Of Piv-I Credentials

#### Description

The information system accepts and electronically verifies Personal Identity Verification-I (PIV-I) credentials.

#### Control Information

Responsible role(s) - Organization

## IA-9 Service Identification And Authentication

#### Description

The organization identifies and authenticates [Assignment: organization-defined information system services] using [Assignment: organization-defined security safeguards].

#### Control Information

Responsible role(s) - Organization

### IA-9 (1) Information Exchange

#### Description

The organization ensures that service providers receive, validate, and transmit identification and authentication information.

#### Control Information

Responsible role(s) - Organization

### IA-9 (2) Transmission Of Decisions

#### Description

The organization ensures that identification and authentication decisions are transmitted between [Assignment: organization-defined services] consistent with organizational policies.

#### Control Information

Responsible role(s) - Organization

## IA-10 Adaptive Identification And Authentication

#### Description

The organization requires that individuals accessing the information system employ [Assignment: organization-defined supplemental authentication techniques or mechanisms] under specific [Assignment: organization-defined circumstances or situations].

#### Control Information

Responsible role(s) - Organization

## IA-11 Re-Authentication

#### Description

The organization requires users and devices to re-authenticate when [Assignment: organization-defined circumstances or situations requiring re-authentication].

#### Control Information

Responsible role(s) - Organization

