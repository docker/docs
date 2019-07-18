---
title: "System and communications protection"
description: "System and communications protection reference"
keywords: "standards, compliance, security, 800-53, System and communications protection"
---

## SC-1 System And Communications Protection Policy And Procedures

#### Description

The organization:
<ol type="a">
<li>Develops, documents, and disseminates to [Assignment: organization-defined personnel or roles]:</li>

<ol type="1">
<li>A system and communications protection policy that addresses purpose, scope, roles, responsibilities, management commitment, coordination among organizational entities, and compliance; and</li>
<li>Procedures to facilitate the implementation of the system and communications protection policy and associated system and communications protection controls; and</li>
</ol>
<li>Reviews and updates the current:</li>

<ol type="1">
<li>System and communications protection policy [Assignment: organization-defined frequency]; and</li>
<li>System and communications protection procedures [Assignment: organization-defined frequency].</li>
</ol>
</ol>

#### Control Information

Responsible role(s) - Organization

## SC-2 Application Partitioning

#### Description

The information system separates user functionality (including user interface services) from information system management functionality.

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
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caeklg">DTR</a></li>
<li><a data-toggle="tab" data-target="#bb2j0dhludq000caekm0">UCP</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caeklg" class="tab-pane fade in active">
Docker Trusted Registry is made up of a number of backend services
that provide for both user functionality (including user interface
services) and system management functionality. Each of these services
operates independently of one another. Additional information can be
found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/dtr/2.3/guides/architecture/">https://docs.docker.com/datacenter/dtr/2.3/guides/architecture/</a></li>
<li><a href="https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Docker_EE_Best_Practices_and_Design_Considerations#Docker_Trusted_Registry">https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Docker_EE_Best_Practices_and_Design_Considerations#Docker_Trusted_Registry</a></li>
</ul>

</div>
<div id="bb2j0dhludq000caekm0" class="tab-pane fade">
Universal Control Plane is made up of a number of backend services
that provide for both user functionality (including user interface
services) and system management functionality. Each of these services
operates independently of one another. Additional information can be
found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/architecture/">https://docs.docker.com/datacenter/ucp/2.2/guides/architecture/</a></li>
<li><a href="https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Docker_EE_Best_Practices_and_Design_Considerations#Universal_Control_Plane">https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Docker_EE_Best_Practices_and_Design_Considerations#Universal_Control_Plane</a></li>
</ul>

</div>
</div>

### SC-2 (1) Interfaces For Non-Privileged Users

#### Description

The information system prevents the presentation of information system management-related functionality at an interface for non-privileged users.

#### Control Information

Responsible role(s) - Organization

## SC-3 Security Function Isolation

#### Description

The information system isolates security functions from nonsecurity functions.

#### Control Information

Responsible role(s) - Organization

### SC-3 (1) Hardware Separation

#### Description

The information system utilizes underlying hardware separation mechanisms to implement security function isolation.

#### Control Information

Responsible role(s) - Organization

### SC-3 (2) Access / Flow Control Functions

#### Description

The information system isolates security functions enforcing access and information flow control from nonsecurity functions and from other security functions.

#### Control Information

Responsible role(s) - Organization

### SC-3 (3) Minimize Nonsecurity Functionality

#### Description

The organization minimizes the number of nonsecurity functions included within the isolation boundary containing security functions.

#### Control Information

Responsible role(s) - Organization

### SC-3 (4) Module Coupling And Cohesiveness

#### Description

The organization implements security functions as largely independent modules that maximize internal cohesiveness within modules and minimize coupling between modules.

#### Control Information

Responsible role(s) - Organization

### SC-3 (5) Layered Structures

#### Description

The organization implements security functions as a layered structure minimizing interactions between layers of the design and avoiding any dependence by lower layers on the functionality or correctness of higher layers.

#### Control Information

Responsible role(s) - Organization

## SC-4 Information In Shared Resources

#### Description

The information system prevents unauthorized and unintended information transfer via shared system resources.

#### Control Information

Responsible role(s) - Organization

### SC-4 (2) Periods Processing

#### Description

The information system prevents unauthorized information transfer via shared resources in accordance with [Assignment: organization-defined procedures] when system processing explicitly switches between different information classification levels or security categories.

#### Control Information

Responsible role(s) - Organization

## SC-5 Denial Of Service Protection

#### Description

The information system protects against or limits the effects of the following types of denial of service attacks: [Assignment: organization-defined types of denial of service attacks or references to sources for such information] by employing [Assignment: organization-defined security safeguards].

#### Control Information

Responsible role(s) - Organization

### SC-5 (1) Restrict Internal Users

#### Description

The information system restricts the ability of individuals to launch [Assignment: organization-defined denial of service attacks] against other information systems.

#### Control Information

Responsible role(s) - Organization

### SC-5 (2) Excess Capacity / Bandwidth / Redundancy

#### Description

The information system manages excess capacity, bandwidth, or other redundancy to limit the effects of information flooding denial of service attacks.

#### Control Information

Responsible role(s) - Organization

### SC-5 (3) Detection / Monitoring

#### Description

The organization:
<ol type="a">
<li>Employs [Assignment: organization-defined monitoring tools] to detect indicators of denial of service attacks against the information system; and</li>
<li>Monitors [Assignment: organization-defined information system resources] to determine if sufficient resources exist to prevent effective denial of service attacks.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

## SC-6 Resource Availability

#### Description

The information system protects the availability of resources by allocating [Assignment: organization-defined resources] by [Selection (one or more); priority; quota; [Assignment: organization-defined security safeguards]].

#### Control Information

Responsible role(s) - Organization

## SC-7 Boundary Protection

#### Description

The information system:
<ol type="a">
<li>Monitors and controls communications at the external boundary of the system and at key internal boundaries within the system;</li>
<li>Implements subnetworks for publicly accessible system components that are [Selection: physically; logically] separated from internal organizational networks; and</li>
<li>Connects to external networks or information systems only through managed interfaces consisting of boundary protection devices arranged in accordance with an organizational security architecture.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SC-7 (3) Access Points

#### Description

The organization limits the number of external network connections to the information system.

#### Control Information

Responsible role(s) - Organization

### SC-7 (4) External Telecommunications Services

#### Description

The organization:
<ol type="a">
<li>Implements a managed interface for each external telecommunication service;</li>
<li>Establishes a traffic flow policy for each managed interface;</li>
<li>Protects the confidentiality and integrity of the information being transmitted across each interface;</li>
<li>Documents each exception to the traffic flow policy with a supporting mission/business need and duration of that need; and</li>
<li>Reviews exceptions to the traffic flow policy [Assignment: organization-defined frequency] and removes exceptions that are no longer supported by an explicit mission/business need.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SC-7 (5) Deny By Default / Allow By Exception

#### Description

The information system at managed interfaces denies network communications traffic by default and allows network communications traffic by exception (i.e., deny all, permit by exception).

#### Control Information

Responsible role(s) - Organization

### SC-7 (7) Prevent Split Tunneling For Remote Devices

#### Description

The information system, in conjunction with a remote device, prevents the device from simultaneously establishing non-remote connections with the system and communicating via some other connection to resources in external networks.

#### Control Information

Responsible role(s) - Organization

### SC-7 (8) Route Traffic To Authenticated Proxy Servers

#### Description

The information system routes [Assignment: organization-defined internal communications traffic] to [Assignment: organization-defined external networks] through authenticated proxy servers at managed interfaces.

#### Control Information

Responsible role(s) - Organization

### SC-7 (9) Restrict Threatening Outgoing Communications Traffic

#### Description

The information system:
<ol type="a">
<li>Detects and denies outgoing communications traffic posing a threat to external information systems; and</li>
<li>Audits the identity of internal users associated with denied communications.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SC-7 (10) Prevent Unauthorized Exfiltration

#### Description

The organization prevents the unauthorized exfiltration of information across managed interfaces.

#### Control Information

Responsible role(s) - Organization

### SC-7 (11) Restrict Incoming Communications Traffic

#### Description

The information system only allows incoming communications from [Assignment: organization-defined authorized sources] to be routed to [Assignment: organization-defined authorized destinations].

#### Control Information

Responsible role(s) - Organization

### SC-7 (12) Host-Based Protection

#### Description

The organization implements [Assignment: organization-defined host-based boundary protection mechanisms] at [Assignment: organization-defined information system components].

#### Control Information

Responsible role(s) - Organization

### SC-7 (13) Isolation Of Security Tools / Mechanisms / Support Components

#### Description

The organization isolates [Assignment: organization-defined information security tools, mechanisms, and support components] from other internal information system components by implementing physically separate subnetworks with managed interfaces to other components of the system.

#### Control Information

Responsible role(s) - Organization

### SC-7 (14) Protects Against Unauthorized Physical Connections

#### Description

The organization protects against unauthorized physical connections at [Assignment: organization-defined managed interfaces].

#### Control Information

Responsible role(s) - Organization

### SC-7 (15) Route Privileged Network Accesses

#### Description

The information system routes all networked, privileged accesses through a dedicated, managed interface for purposes of access control and auditing.

#### Control Information

Responsible role(s) - Organization

### SC-7 (16) Prevent Discovery Of Components / Devices

#### Description

The information system prevents discovery of specific system components composing a managed interface.

#### Control Information

Responsible role(s) - Organization

### SC-7 (17) Automated Enforcement Of Protocol Formats

#### Description

The information system enforces adherence to protocol formats.

#### Control Information

Responsible role(s) - Organization

### SC-7 (18) Fail Secure

#### Description

The information system fails securely in the event of an operational failure of a boundary protection device.

#### Control Information

Responsible role(s) - Organization

### SC-7 (19) Blocks Communication From Non-Organizationally Configured Hosts

#### Description

The information system blocks both inbound and outbound communications traffic between [Assignment: organization-defined communication clients] that are independently configured by end users and external service providers.

#### Control Information

Responsible role(s) - Organization

### SC-7 (20) Dynamic Isolation / Segregation

#### Description

The information system provides the capability to dynamically isolate/segregate [Assignment: organization-defined information system components] from other components of the system.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekmg">Engine</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekmg" class="tab-pane fade in active">
Docker Enterprise Edition is designed to run application containers
whose content can be nonely isolated/segregated from other
application containers within the same node/cluster. This is
accomplished by way of Linux kernel primitives and various security
profiles that can be applied to the underlying host OS. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/engine/security/security/">https://docs.docker.com/engine/security/security/</a></li>
<li><a href="https://docs.docker.com/engine/userguide/networking/overlay-security-model/">https://docs.docker.com/engine/userguide/networking/overlay-security-model/</a></li>
<li><a href="https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Securing_Docker_EE_and_Security_Best_Practices#Engine_and_Node_Security">https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Securing_Docker_EE_and_Security_Best_Practices#Engine_and_Node_Security</a></li>
</ul>

</div>
</div>

### SC-7 (21) Isolation Of Information System Components

#### Description

The organization employs boundary protection mechanisms to separate [Assignment: organization-defined information system components] supporting [Assignment: organization-defined missions and/or business functions].

#### Control Information

Responsible role(s) - Organization

### SC-7 (22) Separate Subnets For Connecting To Different Security Domains

#### Description

The information system implements separate network addresses (i.e., different subnets) to connect to systems in different security domains.

#### Control Information

Responsible role(s) - Organization

### SC-7 (23) Disable Sender Feedback On Protocol Validation Failure

#### Description

The information system disables feedback to senders on protocol format validation failure.

#### Control Information

Responsible role(s) - Organization

## SC-8 Transmission Confidentiality And Integrity

#### Description

The information system protects the [Selection (one or more): confidentiality; integrity] of transmitted information.

#### Control Information

Responsible role(s) - Organization

### SC-8 (1) Cryptographic Or Alternate Physical Protection

#### Description

The information system implements cryptographic mechanisms to [Selection (one or more): prevent unauthorized disclosure of information; detect changes to information] during transmission unless otherwise protected by [Assignment: organization-defined alternative physical safeguards].

#### Control Information

Responsible role(s) - Organization

### SC-8 (2) Pre / Post Transmission Handling

#### Description

The information system maintains the [Selection (one or more): confidentiality; integrity] of information during preparation for transmission and during reception.

#### Control Information

Responsible role(s) - Organization

### SC-8 (3) Cryptographic Protection For Message Externals

#### Description

The information system implements cryptographic mechanisms to protect message externals unless otherwise protected by [Assignment: organization-defined alternative physical safeguards].

#### Control Information

Responsible role(s) - Organization

### SC-8 (4) Conceal / Randomize Communications

#### Description

The information system implements cryptographic mechanisms to conceal or randomize communication patterns unless otherwise protected by [Assignment: organization-defined alternative physical safeguards].

#### Control Information

Responsible role(s) - Organization

## SC-10 Network Disconnect

#### Description

The information system terminates the network connection associated with a communications session at the end of the session or after [Assignment: organization-defined time period] of inactivity.

#### Control Information

Responsible role(s) - Organization

## SC-11 Trusted Path

#### Description

The information system establishes a trusted communications path between the user and the following security functions of the system: [Assignment: organization-defined security functions to include at a minimum, information system authentication and re-authentication].

#### Control Information

Responsible role(s) - Organization

### SC-11 (1) Logical Isolation

#### Description

The information system provides a trusted communications path that is logically isolated and distinguishable from other paths.

#### Control Information

Responsible role(s) - Organization

## SC-12 Cryptographic Key Establishment And Management

#### Description

The organization establishes and manages cryptographic keys for required cryptography employed within the information system in accordance with [Assignment: organization-defined requirements for key generation, distribution, storage, access, and destruction].

#### Control Information

Responsible role(s) - Organization

### SC-12 (1) Availability

#### Description

The organization maintains availability of information in the event of the loss of cryptographic keys by users.

#### Control Information

Responsible role(s) - Organization

### SC-12 (2) Symmetric Keys

#### Description

The organization produces, controls, and distributes symmetric cryptographic keys using [Selection: NIST FIPS-compliant; NSA-approved] key management technology and processes.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekn0">Engine</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekn0" class="tab-pane fade in active">
Docker Enterprise Edition can be installed on the following operating
systems: CentOS 7.1&#43;, Red Hat Enterprise Linux 7.0&#43;, Ubuntu 14.04
LTS&#43;, SUSE Linux Enterprise 12&#43; and Windows Server 2016&#43;. In order to
meet the requirements of this control, reference the chosen operating
system&#39;s documentation to ensure it is configured in FIPS mode.
</div>
</div>

### SC-12 (3) Asymmetric Keys

#### Description

The organization produces, controls, and distributes asymmetric cryptographic keys using [Selection: NSA-approved key management technology and processes; approved PKI Class 3 certificates or prepositioned keying material; approved PKI Class 3 or Class 4 certificates and hardware security tokens that protect the user�s private key].

#### Control Information

Responsible role(s) - Organization

## SC-13 Cryptographic Protection

#### Description

The information system implements [Assignment: organization-defined cryptographic uses and type of cryptography required for each use] in accordance with applicable federal laws, Executive Orders, directives, policies, regulations, and standards.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekng">Engine</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekng" class="tab-pane fade in active">
Docker Enterprise Edition can be installed on the following operating
systems: CentOS 7.1&#43;, Red Hat Enterprise Linux 7.0&#43;, Ubuntu 14.04
LTS&#43;, SUSE Linux Enterprise 12&#43; and Windows Server 2016&#43;. In order to
meet the requirements of this control, reference the chosen operating
system&#39;s documentation to ensure it is configured in FIPS mode.
</div>
</div>

## SC-15 Collaborative Computing Devices

#### Description

The information system:
<ol type="a">
<li>Prohibits remote activation of collaborative computing devices with the following exceptions: [Assignment: organization-defined exceptions where remote activation is to be allowed]; and</li>
<li>Provides an explicit indication of use to users physically present at the devices.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SC-15 (1) Physical Disconnect

#### Description

The information system provides physical disconnect of collaborative computing devices in a manner that supports ease of use.

#### Control Information

Responsible role(s) - Organization

### SC-15 (3) Disabling / Removal In Secure Work Areas

#### Description

The organization disables or removes collaborative computing devices from [Assignment: organization-defined information systems or information system components] in [Assignment: organization-defined secure work areas].

#### Control Information

Responsible role(s) - Organization

### SC-15 (4) Explicitly Indicate Current Participants

#### Description

The information system provides an explicit indication of current participants in [Assignment: organization-defined online meetings and teleconferences].

#### Control Information

Responsible role(s) - Organization

## SC-16 Transmission Of Security Attributes

#### Description

The information system associates [Assignment: organization-defined security attributes] with information exchanged between information systems and between system components.

#### Control Information

Responsible role(s) - Organization

### SC-16 (1) Integrity Validation

#### Description

The information system validates the integrity of transmitted security attributes.

#### Control Information

Responsible role(s) - Organization

## SC-17 Public Key Infrastructure Certificates

#### Description

The organization issues public key certificates under an [Assignment: organization-defined certificate policy] or obtains public key certificates from an approved service provider.

#### Control Information

Responsible role(s) - Organization

## SC-18 Mobile Code

#### Description

The organization:
<ol type="a">
<li>Defines acceptable and unacceptable mobile code and mobile code technologies;</li>
<li>Establishes usage restrictions and implementation guidance for acceptable mobile code and mobile code technologies; and</li>
<li>Authorizes, monitors, and controls the use of mobile code within the information system.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SC-18 (1) Identify Unacceptable Code / Take Corrective Actions

#### Description

The information system identifies [Assignment: organization-defined unacceptable mobile code] and takes [Assignment: organization-defined corrective actions].

#### Control Information

Responsible role(s) - Organization

### SC-18 (2) Acquisition / Development / Use

#### Description

The organization ensures that the acquisition, development, and use of mobile code to be deployed in the information system meets [Assignment: organization-defined mobile code requirements].

#### Control Information

Responsible role(s) - Organization

### SC-18 (3) Prevent Downloading / Execution

#### Description

The information system prevents the download and execution of [Assignment: organization-defined unacceptable mobile code].

#### Control Information

Responsible role(s) - Organization

### SC-18 (4) Prevent Automatic Execution

#### Description

The information system prevents the automatic execution of mobile code in [Assignment: organization-defined software applications] and enforces [Assignment: organization-defined actions] prior to executing the code.

#### Control Information

Responsible role(s) - Organization

### SC-18 (5) Allow Execution Only In Confined Environments

#### Description

The organization allows execution of permitted mobile code only in confined virtual machine environments.

#### Control Information

Responsible role(s) - Organization

## SC-19 Voice Over Internet Protocol

#### Description

The organization:
<ol type="a">
<li>Establishes usage restrictions and implementation guidance for Voice over Internet Protocol (VoIP) technologies based on the potential to cause damage to the information system if used maliciously; and</li>
<li>Authorizes, monitors, and controls the use of VoIP within the information system.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

## SC-20 Secure Name / Address Resolution Service (Authoritative Source)

#### Description

The information system:
<ol type="a">
<li>Provides additional data origin authentication and integrity verification artifacts along with the authoritative name resolution data the system returns in response to external name/address resolution queries; and</li>
<li>Provides the means to indicate the security status of child zones and (if the child supports secure resolution services) to enable verification of a chain of trust among parent and child domains, when operating as part of a distributed, hierarchical namespace.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SC-20 (2) Data Origin / Integrity

#### Description

The information system provides data origin and integrity protection artifacts for internal name/address resolution queries.

#### Control Information

Responsible role(s) - Organization

## SC-21 Secure Name / Address Resolution Service (Recursive Or Caching Resolver)

#### Description

The information system requests and performs data origin authentication and data integrity verification on the name/address resolution responses the system receives from authoritative sources.

#### Control Information

Responsible role(s) - Organization

## SC-22 Architecture And Provisioning For Name / Address Resolution Service

#### Description

The information systems that collectively provide name/address resolution service for an organization are fault-tolerant and implement internal/external role separation.

#### Control Information

Responsible role(s) - Organization

## SC-23 Session Authenticity

#### Description

The information system protects the authenticity of communications sessions.

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
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caeko0">DTR</a></li>
<li><a data-toggle="tab" data-target="#bb2j0dhludq000caekog">Engine</a></li>
<li><a data-toggle="tab" data-target="#bb2j0dhludq000caekp0">UCP</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caeko0" class="tab-pane fade in active">
All remote access sessions to Docker Trusted Registry are protected
with Transport Layer Security (TLS) 1.2 with the AES GCM cipher. This
is included at both the HTTPS application layer for access to the DTR
user interface and for command-line based connections to the registry.
In addition to this, all communication to DTR is enforced by way of
two-way mutual TLS authentication.
</div>
<div id="bb2j0dhludq000caekog" class="tab-pane fade">
All remote access sessions to Docker Enterprise Edition are protected
with Transport Layer Security (TLS) 1.2 with the AES GCM cipher. In
addition to this, all communication to and between Docker Enterprise
Editions is enforced by way of two-way mutual TLS authentication.
</div>
<div id="bb2j0dhludq000caekp0" class="tab-pane fade">
All remote access sessions to Universal Control Plane are protected
with Transport Layer Security (TLS) 1.2 with the AES GCM cipher. This
is included at both the HTTPS application layer for access to the UCP
user interface and for command-line based connections to the cluster.
In addition to this, all communication to UCP is enforced by way of
two-way mutual TLS authentication.
</div>
</div>

### SC-23 (1) Invalidate Session Identifiers At Logout

#### Description

The information system invalidates session identifiers upon user logout or other session termination.

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
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekpg">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekpg" class="tab-pane fade in active">
Docker Enterprise Edition invalidates session identifiers upon user
logout per the requirements of this control.
</div>
</div>

### SC-23 (3) Unique Session Identifiers With Randomization

#### Description

The information system generates a unique session identifier for each session with [Assignment: organization-defined randomness requirements] and recognizes only session identifiers that are system-generated.

#### Control Information

Responsible role(s) - Organization

### SC-23 (5) Allowed Certificate Authorities

#### Description

The information system only allows the use of [Assignment: organization-defined certificate authorities] for verification of the establishment of protected sessions.

#### Control Information

Responsible role(s) - Organization

## SC-24 Fail In Known State

#### Description

The information system fails to a [Assignment: organization-defined known-state] for [Assignment: organization-defined types of failures] preserving [Assignment: organization-defined system state information] in failure.

#### Control Information

Responsible role(s) - Organization

## SC-25 Thin Nodes

#### Description

The organization employs [Assignment: organization-defined information system components] with minimal functionality and information storage.

#### Control Information

Responsible role(s) - Organization

## SC-26 Honeypots

#### Description

The information system includes components specifically designed to be the target of malicious attacks for the purpose of detecting, deflecting, and analyzing such attacks.

#### Control Information

Responsible role(s) - Organization

## SC-27 Platform-Independent Applications

#### Description

The information system includes: [Assignment: organization-defined platform-independent applications].

#### Control Information

Responsible role(s) - Organization

## SC-28 Protection Of Information At Rest

#### Description

The information system protects the [Selection (one or more): confidentiality; integrity] of [Assignment: organization-defined information at rest].

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekq0">Engine</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekq0" class="tab-pane fade in active">
All remote access sessions to Docker Enterprise Edition are protected
with Transport Layer Security (TLS) 1.2 with the AES GCM cipher. In
addition to this, all communication to/from and between Docker
Enterprise Edition nodes is enforced by way of two-way mutual TLS
authentication. All Swarm Mode manager nodes in a Docker Enterprise
Edition cluster store state metadata and user secrets encrypted at
rest using the AES GCM cipher.
</div>
</div>

### SC-28 (1) Cryptographic Protection

#### Description

The information system implements cryptographic mechanisms to prevent unauthorized disclosure and modification of [Assignment: organization-defined information] on [Assignment: organization-defined information system components].

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
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>none<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekqg">DTR</a></li>
<li><a data-toggle="tab" data-target="#bb2j0dhludq000caekr0">Engine</a></li>
<li><a data-toggle="tab" data-target="#bb2j0dhludq000caekrg">UCP</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekqg" class="tab-pane fade in active">
All remote access sessions to Docker Trusted Registry are protected
with Transport Layer Security (TLS) 1.2 with the AES GCM cipher. This
is included at both the HTTPS application layer for access to the DTR
user interface and for command-line based connections to the registry.
In addition to this, all communication to DTR is enforced by way of
two-way mutual TLS authentication.
</div>
<div id="bb2j0dhludq000caekr0" class="tab-pane fade">
All remote access sessions to Docker Enterprise Edition are protected
with Transport Layer Security (TLS) 1.2 with the AES GCM cipher. In
addition to this, all communication to and between Docker Enterprise
Editions is enforced by way of two-way mutual TLS authentication.
</div>
<div id="bb2j0dhludq000caekrg" class="tab-pane fade">
All remote access sessions to Universal Control Plane are protected
with Transport Layer Security (TLS) 1.2 with the AES GCM cipher. This
is included at both the HTTPS application layer for access to the UCP
user interface and for command-line based connections to the cluster.
In addition to this, all communication to UCP is enforced by way of
two-way mutual TLS authentication.
</div>
</div>

### SC-28 (2) Off-Line Storage

#### Description

The organization removes from online storage and stores off-line in a secure location [Assignment: organization-defined information].

#### Control Information

Responsible role(s) - Organization

## SC-29 Heterogeneity

#### Description

The organization employs a diverse set of information technologies for [Assignment: organization-defined information system components] in the implementation of the information system.

#### Control Information

Responsible role(s) - Organization

### SC-29 (1) Virtualization Techniques

#### Description

The organization employs virtualization techniques to support the deployment of a diversity of operating systems and applications that are changed [Assignment: organization-defined frequency].

#### Control Information

Responsible role(s) - Organization

## SC-30 Concealment And Misdirection

#### Description

The organization employs [Assignment: organization-defined concealment and misdirection techniques] for [Assignment: organization-defined information systems] at [Assignment: organization-defined time periods] to confuse and mislead adversaries.

#### Control Information

Responsible role(s) - Organization

### SC-30 (2) Randomness

#### Description

The organization employs [Assignment: organization-defined techniques] to introduce randomness into organizational operations and assets.

#### Control Information

Responsible role(s) - Organization

### SC-30 (3) Change Processing / Storage Locations

#### Description

The organization changes the location of [Assignment: organization-defined processing and/or storage] [Selection: [Assignment: organization-defined time frequency]; at random time intervals]].

#### Control Information

Responsible role(s) - Organization

### SC-30 (4) Misleading Information

#### Description

The organization employs realistic, but misleading information in [Assignment: organization-defined information system components] with regard to its security state or posture.

#### Control Information

Responsible role(s) - Organization

### SC-30 (5) Concealment Of System Components

#### Description

The organization employs [Assignment: organization-defined techniques] to hide or conceal [Assignment: organization-defined information system components].

#### Control Information

Responsible role(s) - Organization

## SC-31 Covert Channel Analysis

#### Description

The organization:
<ol type="a">
<li>Performs a covert channel analysis to identify those aspects of communications within the information system that are potential avenues for covert [Selection (one or more): storage; timing] channels; and</li>
<li>Estimates the maximum bandwidth of those channels.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SC-31 (1) Test Covert Channels For Exploitability

#### Description

The organization tests a subset of the identified covert channels to determine which channels are exploitable.

#### Control Information

Responsible role(s) - Organization

### SC-31 (2) Maximum Bandwidth

#### Description

The organization reduces the maximum bandwidth for identified covert [Selection (one or more); storage; timing] channels to [Assignment: organization-defined values].

#### Control Information

Responsible role(s) - Organization

### SC-31 (3) Measure Bandwidth In Operational Environments

#### Description

The organization measures the bandwidth of [Assignment: organization-defined subset of identified covert channels] in the operational environment of the information system.

#### Control Information

Responsible role(s) - Organization

## SC-32 Information System Partitioning

#### Description

The organization partitions the information system into [Assignment: organization-defined information system components] residing in separate physical domains or environments based on [Assignment: organization-defined circumstances for physical separation of components].

#### Control Information

Responsible role(s) - Organization

## SC-34 Non-Modifiable Executable Programs

#### Description

The information system at [Assignment: organization-defined information system components]:
<ol type="a">
<li>Loads and executes the operating environment from hardware-enforced, read-only media; and</li>
<li>Loads and executes [Assignment: organization-defined applications] from hardware-enforced, read-only media.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SC-34 (1) No Writable Storage

#### Description

The organization employs [Assignment: organization-defined information system components] with no writeable storage that is persistent across component restart or power on/off.

#### Control Information

Responsible role(s) - Organization

### SC-34 (2) Integrity Protection / Read-Only Media

#### Description

The organization protects the integrity of information prior to storage on read-only media and controls the media after such information has been recorded onto the media.

#### Control Information

Responsible role(s) - Organization

### SC-34 (3) Hardware-Based Protection

#### Description

The organization:
<ol type="a">
<li>Employs hardware-based, write-protect for [Assignment: organization-defined information system firmware components]; and</li>
<li>Implements specific procedures for [Assignment: organization-defined authorized individuals] to manually disable hardware write-protect for firmware modifications and re-enable the write-protect prior to returning to operational mode.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

## SC-35 Honeyclients

#### Description

The information system includes components that proactively seek to identify malicious websites and/or web-based malicious code.

#### Control Information

Responsible role(s) - Organization

## SC-36 Distributed Processing And Storage

#### Description

The organization distributes [Assignment: organization-defined processing and storage] across multiple physical locations.

#### Control Information

Responsible role(s) - Organization

### SC-36 (1) Polling Techniques

#### Description

The organization employs polling techniques to identify potential faults, errors, or compromises to [Assignment: organization-defined distributed processing and storage components].

#### Control Information

Responsible role(s) - Organization

## SC-37 Out-Of-Band Channels

#### Description

The organization employs [Assignment: organization-defined out-of-band channels] for the physical delivery or electronic transmission of [Assignment: organization-defined information, information system components, or devices] to [Assignment: organization-defined individuals or information systems].

#### Control Information

Responsible role(s) - Organization

### SC-37 (1) Ensure Delivery / Transmission

#### Description

The organization employs [Assignment: organization-defined security safeguards] to ensure that only [Assignment: organization-defined individuals or information systems] receive the [Assignment: organization-defined information, information system components, or devices].

#### Control Information

Responsible role(s) - Organization

## SC-38 Operations Security

#### Description

The organization employs [Assignment: organization-defined operations security safeguards] to protect key organizational information throughout the system development life cycle.

#### Control Information

Responsible role(s) - Organization

## SC-39 Process Isolation

#### Description

The information system maintains a separate execution domain for each executing process.

#### Control Information

Responsible role(s) - Organization

### SC-39 (1) Hardware Separation

#### Description

The information system implements underlying hardware separation mechanisms to facilitate process separation.

#### Control Information

Responsible role(s) - Organization

### SC-39 (2) Thread Isolation

#### Description

The information system maintains a separate execution domain for each thread in [Assignment: organization-defined multi-threaded processing].

#### Control Information

Responsible role(s) - Organization

## SC-40 Wireless Link Protection

#### Description

The information system protects external and internal [Assignment: organization-defined wireless links] from [Assignment: organization-defined types of signal parameter attacks or references to sources for such attacks].

#### Control Information

Responsible role(s) - Organization

### SC-40 (1) Electromagnetic Interference

#### Description

The information system implements cryptographic mechanisms that achieve [Assignment: organization-defined level of protection] against the effects of intentional electromagnetic interference.

#### Control Information

Responsible role(s) - Organization

### SC-40 (2) Reduce Detection Potential

#### Description

The information system implements cryptographic mechanisms to reduce the detection potential of wireless links to [Assignment: organization-defined level of reduction].

#### Control Information

Responsible role(s) - Organization

### SC-40 (3) Imitative Or Manipulative Communications Deception

#### Description

The information system implements cryptographic mechanisms to identify and reject wireless transmissions that are deliberate attempts to achieve imitative or manipulative communications deception based on signal parameters.

#### Control Information

Responsible role(s) - Organization

### SC-40 (4) Signal Parameter Identification

#### Description

The information system implements cryptographic mechanisms to prevent the identification of [Assignment: organization-defined wireless transmitters] by using the transmitter signal parameters.

#### Control Information

Responsible role(s) - Organization

## SC-41 Port And I/O Device Access

#### Description

The organization physically disables or removes [Assignment: organization-defined connection ports or input/output devices] on [Assignment: organization-defined information systems or information system components].

#### Control Information

Responsible role(s) - Organization

## SC-42 Sensor Capability And Data

#### Description

The information system:
<ol type="a">
<li>Prohibits the remote activation of environmental sensing capabilities with the following exceptions: [Assignment: organization-defined exceptions where remote activation of sensors is allowed]; and</li>
<li>Provides an explicit indication of sensor use to [Assignment: organization-defined class of users].</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SC-42 (1) Reporting To Authorized Individuals Or Roles

#### Description

The organization ensures that the information system is configured so that data or information collected by the [Assignment: organization-defined sensors] is only reported to authorized individuals or roles.

#### Control Information

Responsible role(s) - Organization

### SC-42 (2) Authorized Use

#### Description

The organization employs the following measures: [Assignment: organization-defined measures], so that data or information collected by [Assignment: organization-defined sensors] is only used for authorized purposes.

#### Control Information

Responsible role(s) - Organization

### SC-42 (3) Prohibit Use Of Devices

#### Description

The organization prohibits the use of devices possessing [Assignment: organization-defined environmental sensing capabilities] in [Assignment: organization-defined facilities, areas, or systems].

#### Control Information

Responsible role(s) - Organization

## SC-43 Usage Restrictions

#### Description

The organization:
<ol type="a">
<li>Establishes usage restrictions and implementation guidance for [Assignment: organization-defined information system components] based on the potential to cause damage to the information system if used maliciously; and</li>
<li>Authorizes, monitors, and controls the use of such components within the information system.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

## SC-44 Detonation Chambers

#### Description

The organization employs a detonation chamber capability within [Assignment: organization-defined information system, system component, or location].

#### Control Information

Responsible role(s) - Organization

