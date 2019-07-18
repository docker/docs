---
title: "Risk assessment"
description: "Risk assessment reference"
keywords: "standards, compliance, security, 800-53, Risk assessment"
---

## RA-1 Risk Assessment Policy And Procedures

#### Description

The organization:
<ol type="a">
<li>Develops, documents, and disseminates to [Assignment: organization-defined personnel or roles]:</li>

<ol type="1">
<li>A risk assessment policy that addresses purpose, scope, roles, responsibilities, management commitment, coordination among organizational entities, and compliance; and</li>
<li>Procedures to facilitate the implementation of the risk assessment policy and associated risk assessment controls; and</li>
</ol>
<li>Reviews and updates the current:</li>

<ol type="1">
<li>Risk assessment policy [Assignment: organization-defined frequency]; and</li>
<li>Risk assessment procedures [Assignment: organization-defined frequency].</li>
</ol>
</ol>

#### Control Information

Responsible role(s) - Organization

## RA-2 Security Categorization

#### Description

The organization:
<ol type="a">
<li>Categorizes information and the information system in accordance with applicable federal laws, Executive Orders, directives, policies, regulations, standards, and guidance;</li>
<li>Documents the security categorization results (including supporting rationale) in the security plan for the information system; and</li>
<li>Ensures that the authorizing official or authorizing official designated representative reviews and approves the security categorization decision.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

## RA-3 Risk Assessment

#### Description

The organization:
<ol type="a">
<li>Conducts an assessment of risk, including the likelihood and magnitude of harm, from the unauthorized access, use, disclosure, disruption, modification, or destruction of the information system and the information it processes, stores, or transmits;</li>
<li>Documents risk assessment results in [Selection: security plan; risk assessment report; [Assignment: organization-defined document]];</li>
<li>Reviews risk assessment results [Assignment: organization-defined frequency];</li>
<li>Disseminates risk assessment results to [Assignment: organization-defined personnel or roles]; and</li>
<li>Updates the risk assessment [Assignment: organization-defined frequency] or whenever there are significant changes to the information system or environment of operation (including the identification of new threats and vulnerabilities), or other conditions that may impact the security state of the system.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

## RA-5 Vulnerability Scanning

#### Description

The organization:
<ol type="a">
<li>Scans for vulnerabilities in the information system and hosted applications [Assignment: organization-defined frequency and/or randomly in accordance with organization-defined process] and when new vulnerabilities potentially affecting the system/applications are identified and reported;</li>
<li>Employs vulnerability scanning tools and techniques that facilitate interoperability among tools and automate parts of the vulnerability management process by using standards for:</li>

<ol type="1">
<li>Enumerating platforms, software flaws, and improper configurations;</li>
<li>Formatting checklists and test procedures; and</li>
<li>Measuring vulnerability impact;</li>
</ol>
<li>Analyzes vulnerability scan reports and results from security control assessments;</li>
<li>Remediates legitimate vulnerabilities [Assignment: organization-defined response times] in accordance with an organizational assessment of risk; and</li>
<li>Shares information obtained from the vulnerability scanning process and security control assessments with [Assignment: organization-defined personnel or roles] to help eliminate similar vulnerabilities in other information systems (i.e., systemic weaknesses or deficiencies).</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### RA-5 (1) Update Tool Capability

#### Description

The organization employs vulnerability scanning tools that include the capability to readily update the information system vulnerabilities to be scanned.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Security Scanning (DSS)</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Docker Trusted Registry (DTR)</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekg0">DSS</a></li>
<li><a data-toggle="tab" data-target="#bb2j0dhludq000caekgg">DTR</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekg0" class="tab-pane fade in active">
To assist the orgnization in meeting the requirements of this
control, the Docker Security Scanning (DSS) component of Docker
Trusted Registry (DTR) that is included with the Docker Enterprise
Edition Advanced tier can be used to scan Docker images for
vulnerabilities against known vulnerability databases. Scans can be
triggered either manually or when Docker images are pushed to DTR.
</div>
<div id="bb2j0dhludq000caekgg" class="tab-pane fade">
The Docker Security Scanning tool allows for the scanning of Docker
images in Docker Trusted Registry against the Common Vulnerabilities
and Exposures (CVE) dictionary.
</div>
</div>

### RA-5 (2) Update By Frequency / Prior To New Scan / When Identified

#### Description

The organization updates the information system vulnerabilities scanned [Selection (one or more): [Assignment: organization-defined frequency]; prior to a new scan; when new vulnerabilities are identified and reported].

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Security Scanning (DSS)</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekh0">DSS</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekh0" class="tab-pane fade in active">
To assist the orgnization in meeting the requirements of this
control, the Docker Security Scanning component of Docker Trusted
Registry (DTR) that is included with the Docker Enterprise Edition
Advanced tier compiles a bill of materials (BOM) for each Docker image
that it scans. DSS is also synchronized to an aggregate listing of
known vulnerabilities that is compiled from both the MITRE and NVD CVE
databases. Additional information can be found at the following
resources:


<ul>
<li><a href="https://docs.docker.com/datacenter/dtr/2.3/guides/admin/configure/set-up-vulnerability-scans/">https://docs.docker.com/datacenter/dtr/2.3/guides/admin/configure/set-up-vulnerability-scans/</a></li>
<li><a href="https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Securing_Docker_EE_and_Security_Best_Practices#Image_Scanning">https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Securing_Docker_EE_and_Security_Best_Practices#Image_Scanning</a></li>
</ul>

</div>
</div>

### RA-5 (3) Breadth / Depth Of Coverage

#### Description

The organization employs vulnerability scanning procedures that can identify the breadth and depth of coverage (i.e., information system components scanned and vulnerabilities checked).

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Security Scanning (DSS)</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
<tr>
<td>Docker Trusted Registry (DTR)</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekhg">DSS</a></li>
<li><a data-toggle="tab" data-target="#bb2j0dhludq000caeki0">DTR</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekhg" class="tab-pane fade in active">
To assist the orgnization in meeting the requirements of this
control, the Docker Security Scanning component of Docker Trusted
Registry (DTR) that is included with the Docker Enterprise Edition
Advanced tier identifies vulnerabilities in a Docker image and marks
them against predefined criticality levels; critical major and minor.
</div>
<div id="bb2j0dhludq000caeki0" class="tab-pane fade">
The Docker Security Scanning tool allows for the scanning of Docker
images in Docker Trusted Registry against the Common Vulnerabilities
and Exposures (CVE).&#39; dictionary
</div>
</div>

### RA-5 (4) Discoverable Information

#### Description

The organization determines what information about the information system is discoverable by adversaries and subsequently takes [Assignment: organization-defined corrective actions].

#### Control Information

Responsible role(s) - Organization

### RA-5 (5) Privileged Access

#### Description

The information system implements privileged access authorization to [Assignment: organization-identified information system components] for selected [Assignment: organization-defined vulnerability scanning activities].

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Security Scanning (DSS)</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekig">DSS</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekig" class="tab-pane fade in active">
Only the appropriate users that the organization has provided Docker
Trusted Registry access to are able to view and interpret
vulnerability scan results.
</div>
</div>

### RA-5 (6) Automated Trend Analyses

#### Description

The organization employs automated mechanisms to compare the results of vulnerability scans over time to determine trends in information system vulnerabilities.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Security Scanning (DSS)</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekj0">DSS</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekj0" class="tab-pane fade in active">
For each Docker image pushed to Docker Trusted Registry at a given
time, Docker Security Scaninng retains a list of vulnerabilities
detected. The DTR API can be queried to retrieve the vulnerability
scan results over a period of time for a given Docker image such that
the results can be compared per the requirements of this control.
</div>
</div>

### RA-5 (8) Review Historic Audit Logs

#### Description

The organization reviews historic audit logs to determine if a vulnerability identified in the information system has been previously exploited.

#### Control Information

Responsible role(s) - Docker system

<table>
<tr>
<th>Component</th>
<th>Implementation Status(es)</th>
<th>Control Origin(s)</th>
</tr>
<tr>
<td>Docker Security Scanning (DSS)</td>
<td>none<br/></td>
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#bb2j0dhludq000caekjg">DSS</a></li>
</ul>

<div class="tab-content">
<div id="bb2j0dhludq000caekjg" class="tab-pane fade in active">
Docker Security Scanning maintains a historical bill-of-materials
(BOM) for all Docker images that are scanned. Results of previous
vulnerability scans can be reviewed and audited per the requirements
of this control.
</div>
</div>

### RA-5 (10) Correlate Scanning Information

#### Description

The organization correlates the output from vulnerability scanning tools to determine the presence of multi-vulnerability/multi-hop attack vectors.

#### Control Information

Responsible role(s) - Organization

## RA-6 Technical Surveillance Countermeasures Survey

#### Description

The organization employs a technical surveillance countermeasures survey at [Assignment: organization-defined locations] [Selection (one or more): [Assignment: organization-defined frequency]; [Assignment: organization-defined events or indicators occur]].

#### Control Information

Responsible role(s) - Organization

