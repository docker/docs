---
title: "Audit and accountability"
description: "Audit and accountability reference"
keywords: "standards, compliance, security, 800-53, Audit and accountability"
---

## AU-1 Audit And Accountability Policy And Procedures

#### Description

The organization:
<ol type="a">
<li>Develops, documents, and disseminates to [Assignment: organization-defined personnel or roles]:</li>

<ol type="1">
<li>An audit and accountability policy that addresses purpose, scope, roles, responsibilities, management commitment, coordination among organizational entities, and compliance; and</li>
<li>Procedures to facilitate the implementation of the audit and accountability policy and associated audit and accountability controls; and</li>
</ol>
<li>Reviews and updates the current:</li>

<ol type="1">
<li>Audit and accountability policy [Assignment: organization-defined frequency]; and</li>
<li>Audit and accountability procedures [Assignment: organization-defined frequency].</li>
</ol>
</ol>

#### Control Information

Responsible role(s) - Organization

## AU-2 Audit Events

#### Description

The organization:
<ol type="a">
<li>Determines that the information system is capable of auditing the following events: [Assignment: organization-defined auditable events];</li>
<li>Coordinates the security audit function with other organizational entities requiring audit-related information to enhance mutual support and to help guide the selection of auditable events;</li>
<li>Provides a rationale for why the auditable events are deemed to be adequate to support after-the-fact investigations of security incidents; and</li>
<li>Determines that the following events are to be audited within the information system: [Assignment: organization-defined audited events (the subset of the auditable events defined in AU-2 a.) along with the frequency of (or situation requiring) auditing for each identified event].</li>
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
<td>Docker EE system<br/>service provider corporate<br/>service provider hybrid<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>service provider hybrid<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1h0">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1hg">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1i0">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1h0" class="tab-pane fade in active">
All of the event types indicated by this control are logged by a
combination of the backend ucp-controller service within Universal
Control Plane and the backend services that make up Docker Trusted
Registry. Additional documentation can be found at the following
resource:


<ul>
<li><a href="https://docs.docker.com/datacenter/dtr/2.3/guides/admin/monitor-and-troubleshoot/">https://docs.docker.com/datacenter/dtr/2.3/guides/admin/monitor-and-troubleshoot/</a></li>
<li><a href="https://docs.docker.com/datacenter/dtr/2.3/guides/architecture/#dtr-internal-components">https://docs.docker.com/datacenter/dtr/2.3/guides/architecture/#dtr-internal-components</a></li>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/architecture/#ucp-internal-components">https://docs.docker.com/datacenter/ucp/2.2/guides/architecture/#ucp-internal-components</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1hg" class="tab-pane fade">
Both Universal Control Plane and Docker Trusted Registry backend
service containers, all of which reside on Docker Enterprise Edition,
log all of the event types indicated by this control (as explained by
their component narratives). These and other application containers
that reside on Docker Enterprise Edition can be configured to log data
via an appropriate Docker logging driver. Instructions for configuring
logging drivers can be found at the following resource:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1i0" class="tab-pane fade">
All of the event types indicated by this control are logged by the
backend ucp-controller service within Universal Control Plane. In
addition, each container created on a Universal Control Plane cluster
logs event data. Supporting documentation for configuring UCP logging
can be referenced at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-2 (3) Reviews And Updates

#### Description

The organization reviews and updates the audited events [Assignment: organization-defined frequency].

#### Control Information

Responsible role(s) - Organization

## AU-3 Content Of Audit Records

#### Description

The information system generates audit records containing information that establishes what type of event occurred, when the event occurred, where the event occurred, the source of the event, the outcome of the event, and the identity of any individuals or subjects associated with the event.

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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Authentication and Authorization Service (eNZi)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1ig">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1j0">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1jg">UCP</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1k0">eNZi</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1ig" class="tab-pane fade in active">
Docker Trusted Registry generates all of the audit record information
indicated by this control. A sample audit event has been provided
below:

{&#34;level&#34;:&#34;info&#34;,&#34;license_key&#34;:&#34;123456789123456789123456789&#34;,&#34;msg&#34;:&#34;eNZi:Password
based auth
suceeded&#34;,&#34;remote_addr&#34;:&#34;192.168.33.1:55905&#34;,&#34;time&#34;:&#34;2016-11-09T22:41:01Z&#34;,&#34;type&#34;:&#34;auth
ok&#34;,&#34;username&#34;:&#34;dockeruser&#34;}
</div>
<div id="b70lu81idlmg00bko1j0" class="tab-pane fade">
Both Universal Control Plane and Docker Trusted Registry are
pre-configured to take advantage of Docker Enterprise Edition&#39;s
built-in logging mechanisms. A sample audit event recorded by Docker
Enterprise Edition has been provided below:

{&#34;level&#34;:&#34;info&#34;,&#34;license_key&#34;:&#34;123456789123456789123456789&#34;,&#34;msg&#34;:&#34;eNZi:Password
based auth
suceeded&#34;,&#34;remote_addr&#34;:&#34;192.168.33.1:55905&#34;,&#34;time&#34;:&#34;2016-11-09T22:41:01Z&#34;,&#34;type&#34;:&#34;auth
ok&#34;,&#34;username&#34;:&#34;dockeruser&#34;}

Additional documentation can be referenced at the following resource:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1jg" class="tab-pane fade">
Universal Control Plane generates all of the audit record information
indicated by this control. A sample audit event has been provided
below:

{&#34;level&#34;:&#34;info&#34;,&#34;license_key&#34;:&#34;123456789123456789123456789&#34;,&#34;msg&#34;:&#34;eNZi:Password
based auth
suceeded&#34;,&#34;remote_addr&#34;:&#34;192.168.33.1:55905&#34;,&#34;time&#34;:&#34;2016-11-09T22:41:01Z&#34;,&#34;type&#34;:&#34;auth
ok&#34;,&#34;username&#34;:&#34;dockeruser&#34;}
</div>
<div id="b70lu81idlmg00bko1k0" class="tab-pane fade">
Docker Enterprise Edition generates all of the audit record
information indicated by this control. A sample audit event has been
provided below:

{&#34;level&#34;:&#34;info&#34;,&#34;license_key&#34;:&#34;123456789123456789123456789&#34;,&#34;msg&#34;:&#34;eNZi:Password
based auth
suceeded&#34;,&#34;remote_addr&#34;:&#34;192.168.33.1:55905&#34;,&#34;time&#34;:&#34;2016-11-09T22:41:01Z&#34;,&#34;type&#34;:&#34;auth
ok&#34;,&#34;username&#34;:&#34;dockeruser&#34;}
</div>
</div>

### AU-3 (1) Additional Audit Information

#### Description

The information system generates audit records containing the following additional information: [Assignment: organization-defined additional, more detailed information].

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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1kg">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1l0">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1lg">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1kg" class="tab-pane fade in active">
Universal Control Plane can be configured to log data to a remote
logging stack, which in turn, sends the Docker Trusted Registry
backend container audit records to the remote logging stack. The
logging stack can subsequently be used to interpolate the information
defined by this control from the logged audit records. Additional
information can be found at the following resource:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1l0" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack. The logging
stack can subsequently be used to interpolate the information defined
by this control from the logged audit records. Additional
documentation can be found at the following resource:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1lg" class="tab-pane fade">
Universal Control Plane can be configured to log data to a remote
logging stack. The logging stack can subsequently be used to
interpolate the information defined by this control from the logged
audit records. Additional documentation can be found at the following
resource:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-3 (2) Centralized Management Of Planned Audit Record Content

#### Description

The information system provides centralized management and configuration of the content to be captured in audit records generated by [Assignment: organization-defined information system components].

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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1m0">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1mg">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1n0">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1m0" class="tab-pane fade in active">
Universal Control Plane can be configured to log data to a remote
logging stack, which in turn, sends the Docker Trusted Registry
backend container audit records to the remote logging stack. The
logging stack can subsequently be used to interpolate the information
defined by this control from the logged audit records. Additional
information can be found at the following resource:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1mg" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack. The logging
stack can subsequently be used to interpolate the information defined
by this control from the logged audit records. Additional
documentation can be found at the following resource:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1n0" class="tab-pane fade">
Universal Control Plane can be configured to log data to a remote
logging stack. The logging stack can subsequently be used to
interpolate the information defined by this control from the logged
audit records. Additional documentation can be found at the following
resource:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

## AU-4 Audit Storage Capacity

#### Description

The organization allocates audit record storage capacity in accordance with [Assignment: organization-defined audit record storage requirements].

#### Control Information

Responsible role(s) - Organization

### AU-4 (1) Transfer To Alternate Storage

#### Description

The information system off-loads audit records [Assignment: organization-defined frequency] onto a different system or media than the system being audited.

#### Control Information

Responsible role(s) - Organization

## AU-5 Response To Audit Processing Failures

#### Description

The information system:
<ol type="a">
<li>Alerts [Assignment: organization-defined personnel or roles] in the event of an audit processing failure; and</li>
<li>Takes the following additional actions: [Assignment: organization-defined actions to be taken (e.g., shut down information system, overwrite oldest audit records, stop generating audit records)].</li>
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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider system specific<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1ng">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1o0">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1og">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1ng" class="tab-pane fade in active">
Universal Control Plane can be configured to log data to a remote
logging stack, which in turn, sends the Docker Trusted Registry
backend container audit records to the remote logging stack. The
logging stack can subsequently be configured to alert individuals in
the event of log processing failures. Additional information can be
found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1o0" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack. The logging
stack can be used to interpolate the information defined by this
control and also be configured to alert on any audit processing
failures. Additional information can be found at the following
resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1og" class="tab-pane fade">
Universal Control Plane can be configured to log data to a remote
logging stack. The logging stack can subsequently be configured to
alert individuals in the event of log processing failures. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-5 (1) Audit Storage Capacity

#### Description

The information system provides a warning to [Assignment: organization-defined personnel, roles, and/or locations] within [Assignment: organization-defined time period] when allocated audit record storage volume reaches [Assignment: organization-defined percentage] of repository maximum audit record storage capacity.

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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1p0">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1pg">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1q0">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1p0" class="tab-pane fade in active">
Universal Control Plane can be configured to log data to a remote
logging stack, which in turn, sends the Docker Trusted Registry
backend container audit records to the remote logging stack. The
logging stack can subsequently be configured to warn the organization
when the allocated log storage is full. Additional information can be
found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1pg" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack. The logging
stack can subsequently be configured to warn the organization when the
allocated log storage is full. Additional information can be found at
the following resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1q0" class="tab-pane fade">
Universal Control Plane can be configured to log data to a remote
logging stack. The logging stack can subsequently be configured to
warn the organization when the allocated log storage is full.
Additional information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-5 (2) Real-Time Alerts

#### Description

The information system provides an alert in [Assignment: organization-defined real-time period] to [Assignment: organization-defined personnel, roles, and/or locations] when the following audit failure events occur: [Assignment: organization-defined audit failure events requiring real-time alerts].

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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1qg">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1r0">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1rg">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1qg" class="tab-pane fade in active">
Universal Control Plane can be configured to log data to a remote
logging stack, which in turn, sends the Docker Trusted Registry
backend container audit records to the remote logging stack. The
logging stack can subsequently be configured to warn the organization
when audit log failures occur. Additional information can be found at
the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1r0" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack.  The
logging stack can subsequently be configured to warn the organization
when audit log failures occur. Additional information can be found at
the following resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1rg" class="tab-pane fade">
Universal Control Plane can be configured to log data to a remote
logging stack. The logging stack can subsequently be configured to
warn the organization when audit log failures occur. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-5 (3) Configurable Traffic Volume Thresholds

#### Description

The information system enforces configurable network communications traffic volume thresholds reflecting limits on auditing capacity and [Selection: rejects; delays] network traffic above those thresholds.

#### Control Information

Responsible role(s) - Organization

### AU-5 (4) Shutdown On Failure

#### Description

The information system invokes a [Selection: full system shutdown; partial system shutdown; degraded operational mode with limited mission/business functionality available] in the event of [Assignment: organization-defined audit failures], unless an alternate audit capability exists.

#### Control Information

Responsible role(s) - Organization

## AU-6 Audit Review, Analysis, And Reporting

#### Description

The organization:
<ol type="a">
<li>Reviews and analyzes information system audit records [Assignment: organization-defined frequency] for indications of [Assignment: organization-defined inappropriate or unusual activity]; and</li>
<li>Reports findings to [Assignment: organization-defined personnel or roles].</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### AU-6 (1) Process Integration

#### Description

The organization employs automated mechanisms to integrate audit review, analysis, and reporting processes to support organizational processes for investigation and response to suspicious activities.

#### Control Information

Responsible role(s) - Organization

### AU-6 (3) Correlate Audit Repositories

#### Description

The organization analyzes and correlates audit records across different repositories to gain organization-wide situational awareness.

#### Control Information

Responsible role(s) - Organization

### AU-6 (4) Central Review And Analysis

#### Description

The information system provides the capability to centrally review and analyze audit records from multiple components within the system.

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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1s0">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1sg">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1t0">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1s0" class="tab-pane fade in active">
Universal Control Plane can be configured to log data to a remote
logging stack, which in turn, sends the Docker Trusted Registry
backend container audit records to the remote logging stack. The
organization can subsequently centrally review and analyze all of the
Docker EE audit records. Additional information can be found at the
following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1sg" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack. The
organization can subsequently centrally review and analyze all of the
Docker EE audit records. Additional information can be found at the
following resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1t0" class="tab-pane fade">
Universal Control Plane can be configured to log data to a remote
logging stack. The organization can subsequently centrally review and
analyze all of the Docker EE audit records. Additional information can
be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-6 (5) Integration / Scanning And Monitoring Capabilities

#### Description

The organization integrates analysis of audit records with analysis of [Selection (one or more): vulnerability scanning information; performance data; information system monitoring information; [Assignment: organization-defined data/information collected from other sources]] to further enhance the ability to identify inappropriate or unusual activity.

#### Control Information

Responsible role(s) - Organization

### AU-6 (6) Correlation With Physical Monitoring

#### Description

The organization correlates information from audit records with information obtained from monitoring physical access to further enhance the ability to identify suspicious, inappropriate, unusual, or malevolent activity.

#### Control Information

Responsible role(s) - Organization

### AU-6 (7) Permitted Actions

#### Description

The organization specifies the permitted actions for each [Selection (one or more): information system process; role; user] associated with the review, analysis, and reporting of audit information.

#### Control Information

Responsible role(s) - Organization

### AU-6 (8) Full Text Analysis Of Privileged Commands

#### Description

The organization performs a full text analysis of audited privileged commands in a physically distinct component or subsystem of the information system, or other information system that is dedicated to that analysis.

#### Control Information

Responsible role(s) - Organization

### AU-6 (9) Correlation With Information From Nontechnical Sources

#### Description

The organization correlates information from nontechnical sources with audit information to enhance organization-wide situational awareness.

#### Control Information

Responsible role(s) - Organization

### AU-6 (10) Audit Level Adjustment

#### Description

The organization adjusts the level of audit review, analysis, and reporting within the information system when there is a change in risk based on law enforcement information, intelligence information, or other credible sources of information.

#### Control Information

Responsible role(s) - Organization

## AU-7 Audit Reduction And Report Generation

#### Description

The information system provides an audit reduction and report generation capability that:
<ol type="a">
<li>Supports on-demand audit review, analysis, and reporting requirements and after-the-fact investigations of security incidents; and</li>
<li>Does not alter the original content or time ordering of audit records.</li>
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
<td></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1tg">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1u0">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1ug">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1tg" class="tab-pane fade in active">
Universal Control Plane can be configured to log data to a remote
logging stack, which in turn, sends the Docker Trusted Registry
backend container audit records to the remote logging stack. The
logging stack can subsequently be used to facilitate the audit
reduction and report generation requirements of this control.
Additional information can be found at the following resources:

The underlying operating system chosen to support Docker Trusted
Registry should be certified to ensure that logs are not altered
during generation and transmission to a remote logging stack.
<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1u0" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack. The logging
stack can subsequently be used to facilitate the audit reduction and
report generation requirements of this control. Additional information
can be found at the following resources:

The underlying operating system chosen to support Docker Enterprise
Edition should be certified to ensure that logs are not altered during
generation and transmission to a remote logging stack.
<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1ug" class="tab-pane fade">
Universal Control Plane can be configured to log data to a remote
logging stack. The logging stack can subsequently be used to
facilitate the audit reduction and report generation requirements of
this control. Additional information can be found at the following
resources:

The underlying operating system chosen to support Universal Control
Plane should be certified to ensure that logs are not altered during
generation and transmission to a remote logging stack.
<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-7 (1) Automatic Processing

#### Description

The information system provides the capability to process audit records for events of interest based on [Assignment: organization-defined audit fields within audit records].

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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko1v0">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko1vg">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko200">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko1v0" class="tab-pane fade in active">
Universal Control Plane can be configured to log data to a remote
logging stack, which in turn, sends the Docker Trusted Registry
backend container audit records to the remote logging stack. The
logging stack can subsequently be configured to parse information by
organization-defined audit fields. Additional information can be found
at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko1vg" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack. The logging
stack can subsequently be configured to parse information by
organization-defined audit fields. Additional information can be found
at the following resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko200" class="tab-pane fade">
Universal Control Plane can be configured to log data to a remote
logging stack. The logging stack can subsequently be configured to
parse information by organization-defined audit fields. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-7 (2) Automatic Sort And Search

#### Description

The information system provides the capability to sort and search audit records for events of interest based on the content of [Assignment: organization-defined audit fields within audit records].

#### Control Information

Responsible role(s) - Organization

## AU-8 Time Stamps

#### Description

The information system:
<ol type="a">
<li>Uses internal system clocks to generate time stamps for audit records; and</li>
<li>Records time stamps for audit records that can be mapped to Coordinated Universal Time (UTC) or Greenwich Mean Time (GMT) and meets [Assignment: organization-defined granularity of time measurement].</li>
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
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko20g">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko210">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko21g">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko20g" class="tab-pane fade in active">
Docker Trusted Registry uses the system clock of the underlying
operating system on which it runs. This behavior cannot be modified.The underlying operating system on which Docker Trusted Registry runs
should be configured such that its system clock uses Coordinated
Universal Time (UTC) as indicated by this control. Refer to the
operating system&#39;s instructions for doing so.
</div>
<div id="b70lu81idlmg00bko210" class="tab-pane fade">
Docker Enterprise Edition uses the system clock of the underlying
operating system on which it runs. This behavior cannot be modified.The underlying operating system on which Docker Enterprise Edition
runs should be configured such that its system clock uses Coordinated
Universal Time (UTC) as indicated by this control. Refer to the
operating system&#39;s instructions for doing so.
</div>
<div id="b70lu81idlmg00bko21g" class="tab-pane fade">
Universal Control Plane uses the system clock of the underlying
operating system on which it runs. This behavior cannot be modified.The underlying operating system on which Universal Control Plane runs
should be configured such that its system clock uses Coordinated
Universal Time (UTC) as indicated by this control. Refer to the
operating system&#39;s instructions for doing so.
</div>
</div>

### AU-8 (1) Synchronization With Authoritative Time Source

#### Description

The information system:
<ol type="a">
<li>Compares the internal information system clocks [Assignment: organization-defined frequency] with [Assignment: organization-defined authoritative time source]; and</li>
<li>Synchronizes the internal system clocks to the authoritative time source when the time difference is greater than [Assignment: organization-defined time period].</li>
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
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko220">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko22g">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko230">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko220" class="tab-pane fade in active">
The underlying operating system on which Docker Trusted Registry runs
should be configured such that its system clock compares itself with
an authoritative time source as indicated by this control. This can be
accomplished by utilizing the Network Time Protocol (NTP). Refer to
the operating system&#39;s instructions for doing so.The underlying operating system on which Docker Trusted Registry runs
should be configured such that its system clock synchronizes itself to
an authoritative time source as defined by part (a) of this control
any time the time difference exceeds that of the organization-defined
time period. This can be accomplished by utilizing the Network Time
Protocol (NTP). Refer to the operating system&#39;s instructions for doing
so.
</div>
<div id="b70lu81idlmg00bko22g" class="tab-pane fade">
The underlying operating system on which Docker Enterprise Edition runs should
be configured such that its system clock compares itself with an
authoritative time source as indicated by this control. This can be
accomplished by utilizing the Network Time Protocol (NTP). Refer to
the operating system&#39;s instructions for doing so.The underlying operating system on which Docker Enterprise Edition
runs should be configured such that its system clock synchronizes
itself to an authoritative time source as defined by part (a) of this
control any time the time difference exceeds that of the
organization-defined time period. This can be accomplished by
utilizing the Network Time Protocol (NTP). Refer to the operating
system&#39;s instructions for doing so.
</div>
<div id="b70lu81idlmg00bko230" class="tab-pane fade">
The underlying operating system on which Universal Control Plane runs
should be configured such that its system clock compares itself with
an authoritative time source as indicated by this control. This can be
accomplished by utilizing the Network Time Protocol (NTP). Refer to
the operating system&#39;s instructions for doing so.The underlying operating system on which Universal Control Plane runs
should be configured such that its system clock synchronizes itself to
an authoritative time source as defined by part (a) of this control
any time the time difference exceeds that of the organization-defined
time period. This can be accomplished by utilizing the Network Time
Protocol (NTP). Refer to the operating system&#39;s instructions for doing
so.
</div>
</div>

### AU-8 (2) Secondary Authoritative Time Source

#### Description

The information system identifies a secondary authoritative time source that is located in a different geographic region than the primary authoritative time source.

#### Control Information

Responsible role(s) - Organization

## AU-9 Protection Of Audit Information

#### Description

The information system protects audit information and audit tools from unauthorized access, modification, and deletion.

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
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko23g">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko240">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko24g">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko23g" class="tab-pane fade in active">
By default, Docker Trusted Registry is configured to use the
underlying logging capabilities of Docker Enterprise Edition. As such,
on the underlying Linux operating system, only root and sudo users and
users that have been added to the &#39;docker&#39; group have the ability to
access the logs generated by UCP backend service containers. In
addition, only UCP Administrator users can change the logging endpoint
of the system should it be decided that logs be sent to a remote
logging stack. In this case, the organization is responsible for
configuring the remote logging stack per the provisions of this
control.
</div>
<div id="b70lu81idlmg00bko240" class="tab-pane fade">
On the underlying Linux operating system supporting Docker Enterprise
Edition, only root and sudo users and users that have been added to
the &#34;docker&#34; group have the ability to access the logs generated by
UCP backend service containers. Should the organization decide to
configure Docker Enterprise Edition to use a logging driver other than
the default json-file driver, the organization is subsequently
responsible for configuring the chosen logging stack per the
provisions of this control. In addition, for Linux operating systems
supporting Docker Enterprise Edition that use the systemd daemon, it
is imperative that the Journal is secured per the requirements of this
control. The same applies for Linux operating systems supporting
Docker Enterprise Edition that instead use upstart. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko24g" class="tab-pane fade">
By default, Universal Control Plane is configured to use the
underlying logging capabilities of Docker Enterprise Edition. As such,
on the underlying Linux operating system, only root and sudo users and
users that have been added to the &#39;docker&#39; group have the ability to
access the logs generated by UCP backend service containers. In
addition, only UCP Administrator users can change the logging endpoint
of the system should it be decided that logs be sent to a remote
logging stack. In this case, the organization is responsible for
configuring the remote logging stack per the provisions of this
control.
</div>
</div>

### AU-9 (1) Hardware Write-Once Media

#### Description

The information system writes audit trails to hardware-enforced, write-once media.

#### Control Information

Responsible role(s) - Organization

### AU-9 (2) Audit Backup On Separate Physical Systems / Components

#### Description

The information system backs up audit records [Assignment: organization-defined frequency] onto a physically different system or system component than the system or component being audited.

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
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko250">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko25g">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko260">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko250" class="tab-pane fade in active">
Docker Trusted Registry resides as an Application on a Universal
Control Plane cluster, and can be configured to send logs to a remote
logging stack. The logging stack can subsequently be configured to
back up audit records per the schedule defined by this control.
Additional information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko25g" class="tab-pane fade">
Docker Enterprise Edition can be configured to use a logging driver
that can subsequently meet the backup requirements of this control.
Additional information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko260" class="tab-pane fade">
Universal Control Plane can be configured to send logs to a remote
logging stack. The logging stack can subsequently be configured to
back up audit records per the schedule defined by this control.
Additional information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-9 (3) Cryptographic Protection

#### Description

The information system implements cryptographic mechanisms to protect the integrity of audit information and audit tools.

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
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko26g">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko270">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko27g">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko26g" class="tab-pane fade in active">
Docker Trusted Registry resides as an Application on a Universal
Control Plane cluster, and can be configured to send logs to a remote
logging stack. The logging stack can subsequently be configured to
meet the encryption mechanisms required by this control. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko270" class="tab-pane fade">
Docker Enterprise Edition can be configured to use a logging driver
that can subsequently meet the encryption mechanisms required by this
control. Additional information can be found at the following
resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko27g" class="tab-pane fade">
Universal Control Plane can be configured to send logs to a remote
logging stack. The logging stack can subsequently be configured to
meet the encryption mechanisms required by this control. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-9 (4) Access By Subset Of Privileged Users

#### Description

The organization authorizes access to management of audit functionality to only [Assignment: organization-defined subset of privileged users].

#### Control Information

Responsible role(s) - Organization

### AU-9 (5) Dual Authorization

#### Description

The organization enforces dual authorization for [Selection (one or more): movement; deletion] of [Assignment: organization-defined audit information].

#### Control Information

Responsible role(s) - Organization

### AU-9 (6) Read Only Access

#### Description

The organization authorizes read-only access to audit information to [Assignment: organization-defined subset of privileged users].

#### Control Information

Responsible role(s) - Organization

## AU-10 Non-Repudiation

#### Description

The information system protects against an individual (or process acting on behalf of an individual) falsely denying having performed [Assignment: organization-defined actions to be covered by non-repudiation].

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
<td>complete<br/></td>
<td>Docker EE system<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko280">Engine</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko280" class="tab-pane fade in active">
Docker Enterprise Edition includes functionality known as Docker
Content Trust which allows one to cryptographically sign Docker
images. It enforces client-side signing and verification of image tags
and provides the ability to use digital signatures for data sent to
and received from Docker Trusted Registry. This ultimately provides
one with the ability to verify both the integrity and the publisher of
all data received from DTR over any channel. With Docker Content
Trust, an organization can enforce signature verification of all
content and prohibit unsigned and unapproved content from being
manipulated; thus supproting the non-repudiation requirements of this
control. Additional information can be found at the following
resources:


<ul>
<li><a href="https://docs.docker.com/engine/security/trust/content_trust/">https://docs.docker.com/engine/security/trust/content_trust/</a></li>
</ul>

</div>
</div>

### AU-10 (1) Association Of Identities

#### Description

The information system:
<ol type="a">
<li>Binds the identity of the information producer with the information to [Assignment: organization-defined strength of binding]; and</li>
<li>Provides the means for authorized individuals to determine the identity of the producer of the information.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### AU-10 (2) Validate Binding Of Information Producer Identity

#### Description

The information system:
<ol type="a">
<li>Validates the binding of the information producer identity to the information at [Assignment: organization-defined frequency]; and</li>
<li>Performs [Assignment: organization-defined actions] in the event of a validation error.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### AU-10 (3) Chain Of Custody

#### Description

The information system maintains reviewer/releaser identity and credentials within the established chain of custody for all information reviewed or released.

#### Control Information

Responsible role(s) - Organization

### AU-10 (4) Validate Binding Of Information Reviewer Identity

#### Description

The information system:
<ol type="a">
<li>Validates the binding of the information reviewer identity to the information at the transfer or release points prior to release/transfer between [Assignment: organization-defined security domains]; and</li>
<li>Performs [Assignment: organization-defined actions] in the event of a validation error.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

## AU-11 Audit Record Retention

#### Description

The organization retains audit records for [Assignment: organization-defined time period consistent with records retention policy] to provide support for after-the-fact investigations of security incidents and to meet regulatory and organizational information retention requirements.

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
<td>Docker EE system<br/>service provider corporate<br/>service provider hybrid<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>service provider hybrid<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko28g">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko290">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko29g">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko28g" class="tab-pane fade in active">
The organization will be responsible for meeting the requirements of
this control. To assist with these requirements, Docker Trusted
Registry resides as an Application on a Universal Control Plane
cluster, and as such, can be configured to send logs to a remote
logging stack. This logging stack can subsequently be configured to
retain logs for the duration required by this control. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko290" class="tab-pane fade">
The organization will be responsible for meeting the requirements of
this control. To assist with these requirements, Docker Enterprise
Edition can be configured to use a logging driver that stores data in
a location for the duration specified by this control. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko29g" class="tab-pane fade">
The organization will be responsible for meeting the requirements of
this control. To assist with these requirements, Universal Control
Plane can be configured to send logs to a remote logging stack. This
logging stack can subsequently be configured retain logs for the
duration required by this control. Additional information can be found
at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-11 (1) Long-Term Retrieval Capability

#### Description

The organization employs [Assignment: organization-defined measures] to ensure that long-term audit records generated by the information system can be retrieved.

#### Control Information

Responsible role(s) - Organization

## AU-12 Audit Generation

#### Description

The information system:
<ol type="a">
<li>Provides audit record generation capability for the auditable events defined in AU-2 a. at [Assignment: organization-defined information system components];</li>
<li>Allows [Assignment: organization-defined personnel or roles] to select which auditable events are to be audited by specific components of the information system; and</li>
<li>Generates audit records for the events defined in AU-2 d. with the content defined in AU-3.</li>
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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko2a0">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko2ag">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko2b0">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko2a0" class="tab-pane fade in active">
All of the event types indicated by AU-2 a. are logged by a
combination of the backend services within Universal Control Plane and
Docker Trusted Registry. The underlying Linux operating system
supporting DTR can be configured to audit Docker-specific events with
the auditd daemon. Refer to the specific Linux distribution in use for
instructions on configuring this service. Additional information can
be found at the following resources:

Using auditd on the Linux operating system supporting DTR, the
organization can configure audit rules to select which Docker-specific
events are to be audited. Refer to the specific Linux distribution in
use for instructions on configuring this service.
<ul>
<li><a href="https://docs.docker.com/datacenter/dtr/2.3/guides/admin/monitor-and-troubleshoot/">https://docs.docker.com/datacenter/dtr/2.3/guides/admin/monitor-and-troubleshoot/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko2ag" class="tab-pane fade">
Both Universal Control Plane and Docker Trusted Registry backend
service containers, all of which reside on Docker Enterprise Edition,
log all of the event types indicated by this AU-2 a. These and other
application containers that reside on Docker Enterprise Edition can be
configured to log data via an appropriate Docker logging driver. The
underlying Linux operating system supporting Docker Enterprise Edition
can be configured to audit Docker-specific events with the auditd
daemon. Refer to the specific Linux distribution in use for
instructions on configuring this service. Additional information can
be found at the following resources:

Using auditd on the Linux operating system supporting CS Docker
Engine, the organization can configure audit rules to select which
Docker-specific events are to be audited. Refer to the specific Linux
distribution in use for instructions on configuring this service.
<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko2b0" class="tab-pane fade">
All of the event types indicated by AU-2 a. are logged by the backend
ucp-controller service within Universal Control Plane. In addition,
each container created on a Universal Control Plane cluster logs event
data. The underlying Linux operating system supporting UCP can be
configured to audit Docker-specific events with the auditd daemon.
Refer to the specific Linux distribution in use for instructions on
configuring this service. Additional information can be found at the
following resources:

Using auditd on the Linux operating system supporting UCP, the
organization can configure audit rules to select which Docker-specific
events are to be audited. Refer to the specific Linux distribution in
use for instructions on configuring this service.
<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-12 (1) System-Wide / Time-Correlated Audit Trail

#### Description

The information system compiles audit records from [Assignment: organization-defined information system components] into a system-wide (logical or physical) audit trail that is time-correlated to within [Assignment: organization-defined level of tolerance for the relationship between time stamps of individual records in the audit trail].

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
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>Docker EE system<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko2bg">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko2c0">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko2cg">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko2bg" class="tab-pane fade in active">
Docker Trusted Registry resides as an Application on a Universal
Control Plane cluster, and as such, can be configured to send logs to
a remote logging stack. This logging stack can subsequently be used to
compile audit records in to a system-wide audit trail that is
time-correlated per the requirements of this control. Additional
information can be found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko2c0" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack. This
logging stack can subsequently be used to compile audit records in to
a system-wide audit trail that is time-correlated per the requirements
of this control. Additional information can be found at the following
resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko2cg" class="tab-pane fade">
Universal Control Plane can be configured to send logs to a remote
logging stack. This logging stack can subsequently be used to compile
audit records in to a system-wide audit trail that is time-correlated
per the requirements of this control. Additional information can be
found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

### AU-12 (2) Standardized Formats

#### Description

The information system produces a system-wide (logical or physical) audit trail composed of audit records in a standardized format.

#### Control Information

Responsible role(s) - Organization

### AU-12 (3) Changes By Authorized Individuals

#### Description

The information system provides the capability for [Assignment: organization-defined individuals or roles] to change the auditing to be performed on [Assignment: organization-defined information system components] based on [Assignment: organization-defined selectable event criteria] within [Assignment: organization-defined time thresholds].

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
<td>service provider hybrid<br/>shared<br/></td>
</tr>
<tr>
<td>Docker Enterprise Edition Engine</td>
<td>complete<br/></td>
<td>service provider hybrid<br/>shared<br/></td>
</tr>
<tr>
<td>Universal Control Plane (UCP)</td>
<td>complete<br/></td>
<td>service provider hybrid<br/>shared<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu81idlmg00bko2d0">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko2dg">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu81idlmg00bko2e0">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu81idlmg00bko2d0" class="tab-pane fade in active">
Docker Trusted Registry resides as an Application on a Universal
Control Plane cluster, and as such, can be configured to send logs to
a remote logging stack. This logging stack can subsequently be used to
meet the requirements of this control. Additional information can be
found at the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko2dg" class="tab-pane fade">
Docker Enterprise Edition can be configured with various logging
drivers to send audit events to an external logging stack. This
logging stack can subsequently be used to meet the requirements of
this control. Additional information can be found at the following
resources:


<ul>
<li><a href="https://docs.docker.com/engine/admin/logging/overview/">https://docs.docker.com/engine/admin/logging/overview/</a></li>
</ul>

</div>
<div id="b70lu81idlmg00bko2e0" class="tab-pane fade">
Universal Control Plane can be configured to send logs to a remote
logging stack. This logging stack can subsequently be used to meet the
requirements of this control. Additional information can be found at
the following resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/">https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/store-logs-in-an-external-system/</a></li>
</ul>

</div>
</div>

## AU-13 Monitoring For Information Disclosure

#### Description

The organization monitors [Assignment: organization-defined open source information and/or information sites] [Assignment: organization-defined frequency] for evidence of unauthorized disclosure of organizational information.

#### Control Information

Responsible role(s) - Organization

### AU-13 (1) Use Of Automated Tools

#### Description

The organization employs automated mechanisms to determine if organizational information has been disclosed in an unauthorized manner.

#### Control Information

Responsible role(s) - Organization

### AU-13 (2) Review Of Monitored Sites

#### Description

The organization reviews the open source information sites being monitored [Assignment: organization-defined frequency].

#### Control Information

Responsible role(s) - Organization

## AU-14 Session Audit

#### Description

The information system provides the capability for authorized users to select a user session to capture/record or view/hear.

#### Control Information

Responsible role(s) - Organization

### AU-14 (1) System Start-Up

#### Description

The information system initiates session audits at system start-up.

#### Control Information

Responsible role(s) - Organization

### AU-14 (2) Capture/Record And Log Content

#### Description

The information system provides the capability for authorized users to capture/record and log content related to a user session.

#### Control Information

Responsible role(s) - Organization

### AU-14 (3) Remote Viewing / Listening

#### Description

The information system provides the capability for authorized users to remotely view/hear all content related to an established user session in real time.

#### Control Information

Responsible role(s) - Organization

## AU-15 Alternate Audit Capability

#### Description

The organization provides an alternate audit capability in the event of a failure in primary audit capability that provides [Assignment: organization-defined alternate audit functionality].

#### Control Information

Responsible role(s) - Organization

## AU-16 Cross-Organizational Auditing

#### Description

The organization employs [Assignment: organization-defined methods] for coordinating [Assignment: organization-defined audit information] among external organizations when audit information is transmitted across organizational boundaries.

#### Control Information

Responsible role(s) - Organization

### AU-16 (1) Identity Preservation

#### Description

The organization requires that the identity of individuals be preserved in cross-organizational audit trails.

#### Control Information

Responsible role(s) - Organization

### AU-16 (2) Sharing Of Audit Information

#### Description

The organization provides cross-organizational audit information to [Assignment: organization-defined organizations] based on [Assignment: organization-defined cross-organizational sharing agreements].

#### Control Information

Responsible role(s) - Organization

