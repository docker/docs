---
title: "System and information integrity"
description: "System and information integrity reference"
keywords: "standards, compliance, security, 800-53, System and information integrity"
---

## SI-1 System And Information Integrity Policy And Procedures

#### Description

The organization:
<ol type="a">
<li>Develops, documents, and disseminates to [Assignment: organization-defined personnel or roles]:</li>

<ol type="1">
<li>A system and information integrity policy that addresses purpose, scope, roles, responsibilities, management commitment, coordination among organizational entities, and compliance; and</li>
<li>Procedures to facilitate the implementation of the system and information integrity policy and associated system and information integrity controls; and</li>
</ol>
<li>Reviews and updates the current:</li>

<ol type="1">
<li>System and information integrity policy [Assignment: organization-defined frequency]; and</li>
<li>System and information integrity procedures [Assignment: organization-defined frequency].</li>
</ol>
</ol>

#### Control Information

Responsible role(s) - Organization

## SI-2 Flaw Remediation

#### Description

The organization:
<ol type="a">
<li>Identifies, reports, and corrects information system flaws;</li>
<li>Tests software and firmware updates related to flaw remediation for effectiveness and potential side effects before installation;</li>
<li>Installs security-relevant software and firmware updates within [Assignment: organization-defined time period] of the release of the updates; and</li>
<li>Incorporates flaw remediation into the organizational configuration management process.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-2 (1) Central Management

#### Description

The organization centrally manages the flaw remediation process.

#### Control Information

Responsible role(s) - Organization

### SI-2 (2) Automated Flaw Remediation Status

#### Description

The organization employs automated mechanisms [Assignment: organization-defined frequency] to determine the state of information system components with regard to flaw remediation.

#### Control Information

Responsible role(s) - Organization

### SI-2 (3) Time To Remediate Flaws / Benchmarks For Corrective Actions

#### Description

The organization:
<ol type="a">
<li>Measures the time between flaw identification and flaw remediation; and</li>
<li>Establishes [Assignment: organization-defined benchmarks] for taking corrective actions.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-2 (5) Automatic Software / Firmware Updates

#### Description

The organization installs [Assignment: organization-defined security-relevant software and firmware updates] automatically to [Assignment: organization-defined information system components].

#### Control Information

Responsible role(s) - Organization

### SI-2 (6) Removal Of Previous Versions Of Software / Firmware

#### Description

The organization removes [Assignment: organization-defined software and firmware components] after updated versions have been installed.

#### Control Information

Responsible role(s) - Organization

## SI-3 Malicious Code Protection

#### Description

The organization:
<ol type="a">
<li>Employs malicious code protection mechanisms at information system entry and exit points to detect and eradicate malicious code;</li>
<li>Updates malicious code protection mechanisms whenever new releases are available in accordance with organizational configuration management policy and procedures;</li>
<li>Configures malicious code protection mechanisms to:</li>

<ol type="1">
<li>Perform periodic scans of the information system [Assignment: organization-defined frequency] and real-time scans of files from external sources at [Selection (one or more); endpoint; network entry/exit points] as the files are downloaded, opened, or executed in accordance with organizational security policy; and</li>
<li>[Selection (one or more): block malicious code; quarantine malicious code;  send alert to administrator; [Assignment: organization-defined action]] in response to malicious code detection; and</li>
</ol>
<li>Addresses the receipt of false positives during malicious code detection and eradication and the resulting potential impact on the availability of the information system.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-3 (1) Central Management

#### Description

The organization centrally manages malicious code protection mechanisms.

#### Control Information

Responsible role(s) - Organization

### SI-3 (2) Automatic Updates

#### Description

The information system automatically updates malicious code protection mechanisms.

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
<li class="active"><a data-toggle="tab" data-target="#b70lu8hidlmg00bko3ng">Engine</a></li>
</ul>

<div class="tab-content">
<div id="b70lu8hidlmg00bko3ng" class="tab-pane fade in active">
Docker Enterprise Edition packages for supported underlying operating
systems can only be obtained from Docker, Inc. The Docker EE
repositories from which Docker EE packages are obtained are protected
with official GPG keys. Each Docker package is also validated with a
signature definition.
</div>
</div>

### SI-3 (4) Updates Only By Privileged Users

#### Description

The information system updates malicious code protection mechanisms only when directed by a privileged user.

#### Control Information

Responsible role(s) - Organization

### SI-3 (6) Testing / Verification

#### Description

The organization:
<ol type="a">
<li>Tests malicious code protection mechanisms [Assignment: organization-defined frequency] by introducing a known benign, non-spreading test case into the information system; and</li>
<li>Verifies that both detection of the test case and associated incident reporting occur.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-3 (7) Nonsignature-Based Detection

#### Description

The information system implements nonsignature-based malicious code detection mechanisms.

#### Control Information

Responsible role(s) - Organization

### SI-3 (8) Detect Unauthorized Commands

#### Description

The information system detects [Assignment: organization-defined unauthorized operating system commands] through the kernel application programming interface at [Assignment: organization-defined information system hardware components] and [Selection (one or more): issues a warning; audits the command execution; prevents the execution of the command].

#### Control Information

Responsible role(s) - Organization

### SI-3 (9) Authenticate Remote Commands

#### Description

The information system implements [Assignment: organization-defined security safeguards] to authenticate [Assignment: organization-defined remote commands].

#### Control Information

Responsible role(s) - Organization

### SI-3 (10) Malicious Code Analysis

#### Description

The organization:
<ol type="a">
<li>Employs [Assignment: organization-defined tools and techniques] to analyze the characteristics and behavior of malicious code; and</li>
<li>Incorporates the results from malicious code analysis into organizational incident response and flaw remediation processes.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

## SI-4 Information System Monitoring

#### Description

The organization:
<ol type="a">
<li>Monitors the information system to detect:</li>

<ol type="1">
<li>Attacks and indicators of potential attacks in accordance with [Assignment: organization-defined monitoring objectives]; and</li>
<li>Unauthorized local, network, and remote connections;</li>
</ol>
<li>Identifies unauthorized use of the information system through [Assignment: organization-defined techniques and methods];</li>
<li>Deploys monitoring devices:</li>

<ol type="1">
<li>Strategically within the information system to collect organization-determined essential information; and</li>
<li>At ad hoc locations within the system to track specific types of transactions of interest to the organization;</li>
</ol>
<li>Protects information obtained from intrusion-monitoring tools from unauthorized access, modification, and deletion;</li>
<li>Heightens the level of information system monitoring activity whenever there is an indication of increased risk to organizational operations and assets, individuals, other organizations, or the Nation based on law enforcement information, intelligence information, or other credible sources of information;</li>
<li>Obtains legal opinion with regard to information system monitoring activities in accordance with applicable federal laws, Executive Orders, directives, policies, or regulations; and</li>
<li>Provides [Assignment: organization-defined information system monitoring information] to [Assignment: organization-defined personnel or roles] [Selection (one or more): as needed; [Assignment: organization-defined frequency]].</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-4 (1) System-Wide Intrusion Detection System

#### Description

The organization connects and configures individual intrusion detection tools into an information system-wide intrusion detection system.

#### Control Information

Responsible role(s) - Organization

### SI-4 (2) Automated Tools For Real-Time Analysis

#### Description

The organization employs automated tools to support near real-time analysis of events.

#### Control Information

Responsible role(s) - Organization

### SI-4 (3) Automated Tool Integration

#### Description

The organization employs automated tools to integrate intrusion detection tools into access control and flow control mechanisms for rapid response to attacks by enabling reconfiguration of these mechanisms in support of attack isolation and elimination.

#### Control Information

Responsible role(s) - Organization

### SI-4 (4) Inbound And Outbound Communications Traffic

#### Description

The information system monitors inbound and outbound communications traffic [Assignment: organization-defined frequency] for unusual or unauthorized activities or conditions.

#### Control Information

Responsible role(s) - Organization

### SI-4 (5) System-Generated Alerts

#### Description

The information system alerts [Assignment: organization-defined personnel or roles] when the following indications of compromise or potential compromise occur: [Assignment: organization-defined compromise indicators].

#### Control Information

Responsible role(s) - Organization

### SI-4 (7) Automated Response To Suspicious Events

#### Description

The information system notifies [Assignment: organization-defined incident response personnel (identified by name and/or by role)] of detected suspicious events and takes [Assignment: organization-defined least-disruptive actions to terminate suspicious events].

#### Control Information

Responsible role(s) - Organization

### SI-4 (9) Testing Of Monitoring Tools

#### Description

The organization tests intrusion-monitoring tools [Assignment: organization-defined frequency].

#### Control Information

Responsible role(s) - Organization

### SI-4 (10) Visibility Of Encrypted Communications

#### Description

The organization makes provisions so that [Assignment: organization-defined encrypted communications traffic] is visible to [Assignment: organization-defined information system monitoring tools].

#### Control Information

Responsible role(s) - Organization

### SI-4 (11) Analyze Communications Traffic Anomalies

#### Description

The organization analyzes outbound communications traffic at the external boundary of the information system and selected [Assignment: organization-defined interior points within the system (e.g., subnetworks, subsystems)] to discover anomalies.

#### Control Information

Responsible role(s) - Organization

### SI-4 (12) Automated Alerts

#### Description

The organization employs automated mechanisms to alert security personnel of the following inappropriate or unusual activities with security implications: [Assignment: organization-defined activities that trigger alerts].

#### Control Information

Responsible role(s) - Organization

### SI-4 (13) Analyze Traffic / Event Patterns

#### Description

The organization:
<ol type="a">
<li>Analyzes communications traffic/event patterns for the information system;</li>
<li>Develops profiles representing common traffic patterns and/or events; and</li>
<li>Uses the traffic/event profiles in tuning system-monitoring devices to reduce the number of false positives and the number of false negatives.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-4 (14) Wireless Intrusion Detection

#### Description

The organization employs a wireless intrusion detection system to identify rogue wireless devices and to detect attack attempts and potential compromises/breaches to the information system.

#### Control Information

Responsible role(s) - Organization

### SI-4 (15) Wireless To Wireline Communications

#### Description

The organization employs an intrusion detection system to monitor wireless communications traffic as the traffic passes from wireless to wireline networks.

#### Control Information

Responsible role(s) - Organization

### SI-4 (16) Correlate Monitoring Information

#### Description

The organization correlates information from monitoring tools employed throughout the information system.

#### Control Information

Responsible role(s) - Organization

### SI-4 (17) Integrated Situational Awareness

#### Description

The organization correlates information from monitoring physical, cyber, and supply chain activities to achieve integrated, organization-wide situational awareness.

#### Control Information

Responsible role(s) - Organization

### SI-4 (18) Analyze Traffic / Covert Exfiltration

#### Description

The organization analyzes outbound communications traffic at the external boundary of the information system (i.e., system perimeter) and at [Assignment: organization-defined interior points within the system (e.g., subsystems, subnetworks)] to detect covert exfiltration of information.

#### Control Information

Responsible role(s) - Organization

### SI-4 (19) Individuals Posing Greater Risk

#### Description

The organization implements [Assignment: organization-defined additional monitoring] of individuals who have been identified by [Assignment: organization-defined sources] as posing an increased level of risk.

#### Control Information

Responsible role(s) - Organization

### SI-4 (20) Privileged Users

#### Description

The organization implements [Assignment: organization-defined additional monitoring] of privileged users.

#### Control Information

Responsible role(s) - Organization

### SI-4 (21) Probationary Periods

#### Description

The organization implements [Assignment: organization-defined additional monitoring] of individuals during [Assignment: organization-defined probationary period].

#### Control Information

Responsible role(s) - Organization

### SI-4 (22) Unauthorized Network Services

#### Description

The information system detects network services that have not been authorized or approved by [Assignment: organization-defined authorization or approval processes] and [Selection (one or more): audits; alerts [Assignment: organization-defined personnel or roles]].

#### Control Information

Responsible role(s) - Organization

### SI-4 (23) Host-Based Devices

#### Description

The organization implements [Assignment: organization-defined host-based monitoring mechanisms] at [Assignment: organization-defined information system components].

#### Control Information

Responsible role(s) - Organization

### SI-4 (24) Indicators Of Compromise

#### Description

The information system discovers, collects, distributes, and uses indicators of compromise.

#### Control Information

Responsible role(s) - Organization

## SI-5 Security Alerts, Advisories, And Directives

#### Description

The organization:
<ol type="a">
<li>Receives information system security alerts, advisories, and directives from [Assignment: organization-defined external organizations] on an ongoing basis;</li>
<li>Generates internal security alerts, advisories, and directives as deemed necessary;</li>
<li>Disseminates security alerts, advisories, and directives to: [Selection (one or more): [Assignment: organization-defined personnel or roles]; [Assignment: organization-defined elements within the organization]; [Assignment: organization-defined external organizations]]; and</li>
<li>Implements security directives in accordance with established time frames, or notifies the issuing organization of the degree of noncompliance.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-5 (1) Automated Alerts And Advisories

#### Description

The organization employs automated mechanisms to make security alert and advisory information available throughout the organization.

#### Control Information

Responsible role(s) - Organization

## SI-6 Security Function Verification

#### Description

The information system:
<ol type="a">
<li>Verifies the correct operation of [Assignment: organization-defined security functions];</li>
<li>Performs this verification [Selection (one or more): [Assignment: organization-defined system transitional states]; upon command by user with appropriate privilege; [Assignment: organization-defined frequency]];</li>
<li>Notifies [Assignment: organization-defined personnel or roles] of failed security verification tests; and</li>
<li>[Selection (one or more): shuts the information system down; restarts the information system; [Assignment: organization-defined alternative action(s)]] when anomalies are discovered.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-6 (2) Automation Support For Distributed Testing

#### Description

The information system implements automated mechanisms to support the management of distributed security testing.

#### Control Information

Responsible role(s) - Organization

### SI-6 (3) Report Verification Results

#### Description

The organization reports the results of security function verification to [Assignment: organization-defined personnel or roles].

#### Control Information

Responsible role(s) - Organization

## SI-7 Software, Firmware, And Information Integrity

#### Description

The organization employs integrity verification tools to detect unauthorized changes to [Assignment: organization-defined software, firmware, and information].

#### Control Information

Responsible role(s) - Organization

### SI-7 (1) Integrity Checks

#### Description

The information system performs an integrity check of [Assignment: organization-defined software, firmware, and information] [Selection (one or more): at startup; at [Assignment: organization-defined transitional states or security-relevant events]; [Assignment: organization-defined frequency]].

#### Control Information

Responsible role(s) - Organization

### SI-7 (2) Automated Notifications Of Integrity Violations

#### Description

The organization employs automated tools that provide notification to [Assignment: organization-defined personnel or roles] upon discovering discrepancies during integrity verification.

#### Control Information

Responsible role(s) - Organization

### SI-7 (3) Centrally-Managed Integrity Tools

#### Description

The organization employs centrally managed integrity verification tools.

#### Control Information

Responsible role(s) - Organization

### SI-7 (5) Automated Response To Integrity Violations

#### Description

The information system automatically [Selection (one or more): shuts the information system down; restarts the information system; implements [Assignment: organization-defined security safeguards]] when integrity violations are discovered.

#### Control Information

Responsible role(s) - Organization

### SI-7 (6) Cryptographic Protection

#### Description

The information system implements cryptographic mechanisms to detect unauthorized changes to software, firmware, and information.

#### Control Information

Responsible role(s) - Organization

### SI-7 (7) Integration Of Detection And Response

#### Description

The organization incorporates the detection of unauthorized [Assignment: organization-defined security-relevant changes to the information system] into the organizational incident response capability.

#### Control Information

Responsible role(s) - Organization

### SI-7 (8) Auditing Capability For Significant Events

#### Description

The information system, upon detection of a potential integrity violation, provides the capability to audit the event and initiates the following actions: [Selection (one or more): generates an audit record; alerts current user; alerts [Assignment: organization-defined personnel or roles]; [Assignment: organization-defined other actions]].

#### Control Information

Responsible role(s) - Organization

### SI-7 (9) Verify Boot Process

#### Description

The information system verifies the integrity of the boot process of [Assignment: organization-defined devices].

#### Control Information

Responsible role(s) - Organization

### SI-7 (10) Protection Of Boot Firmware

#### Description

The information system implements [Assignment: organization-defined security safeguards] to protect the integrity of boot firmware in [Assignment: organization-defined devices].

#### Control Information

Responsible role(s) - Organization

### SI-7 (11) Confined Environments With Limited Privileges

#### Description

The organization requires that [Assignment: organization-defined user-installed software] execute in a confined physical or virtual machine environment with limited privileges.

#### Control Information

Responsible role(s) - Organization

### SI-7 (12) Integrity Verification

#### Description

The organization requires that the integrity of [Assignment: organization-defined user-installed software] be verified prior to execution.

#### Control Information

Responsible role(s) - Organization

### SI-7 (13) Code Execution In Protected Environments

#### Description

The organization allows execution of binary or machine-executable code obtained from sources with limited or no warranty and without the provision of source code only in confined physical or virtual machine environments and with the explicit approval of [Assignment: organization-defined personnel or roles].

#### Control Information

Responsible role(s) - Organization

### SI-7 (14) Binary Or Machine Executable Code

#### Description

The organization:
<ol type="a">
<li>Prohibits the use of binary or machine-executable code from sources with limited or no warranty and without the provision of source code; and</li>
<li>Provides exceptions to the source code requirement only for compelling mission/operational requirements and with the approval of the authorizing official.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-7 (15) Code Authentication

#### Description

The information system implements cryptographic mechanisms to authenticate [Assignment: organization-defined software or firmware components] prior to installation.

#### Control Information

Responsible role(s) - Organization

### SI-7 (16) Time Limit On Process Execution W/O Supervision

#### Description

The organization does not allow processes to execute without supervision for more than [Assignment: organization-defined time period].

#### Control Information

Responsible role(s) - Organization

## SI-8 Spam Protection

#### Description

The organization:
<ol type="a">
<li>Employs spam protection mechanisms at information system entry and exit points to detect and take action on unsolicited messages; and</li>
<li>Updates spam protection mechanisms when new releases are available in accordance with organizational configuration management policy and procedures.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-8 (1) Central Management

#### Description

The organization centrally manages spam protection mechanisms.

#### Control Information

Responsible role(s) - Organization

### SI-8 (2) Automatic Updates

#### Description

The information system automatically updates spam protection mechanisms.

#### Control Information

Responsible role(s) - Organization

### SI-8 (3) Continuous Learning Capability

#### Description

The information system implements spam protection mechanisms with a learning capability to more effectively identify legitimate communications traffic.

#### Control Information

Responsible role(s) - Organization

## SI-10 Information Input Validation

#### Description

The information system checks the validity of [Assignment: organization-defined information inputs].

#### Control Information

Responsible role(s) - Organization

### SI-10 (1) Manual Override Capability

#### Description

The information system:
<ol type="a">
<li>Provides a manual override capability for input validation of [Assignment: organization-defined inputs];</li>
<li>Restricts the use of the manual override capability to only [Assignment: organization-defined authorized individuals]; and</li>
<li>Audits the use of the manual override capability.</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-10 (2) Review / Resolution Of Errors

#### Description

The organization ensures that input validation errors are reviewed and resolved within [Assignment: organization-defined time period].

#### Control Information

Responsible role(s) - Organization

### SI-10 (3) Predictable Behavior

#### Description

The information system behaves in a predictable and documented manner that reflects organizational and system objectives when invalid inputs are received.

#### Control Information

Responsible role(s) - Organization

### SI-10 (4) Review / Timing Interactions

#### Description

The organization accounts for timing interactions among information system components in determining appropriate responses for invalid inputs.

#### Control Information

Responsible role(s) - Organization

### SI-10 (5) Restrict Inputs To Trusted Sources And Approved Formats

#### Description

The organization restricts the use of information inputs to [Assignment: organization-defined trusted sources] and/or [Assignment: organization-defined formats].

#### Control Information

Responsible role(s) - Organization

## SI-11 Error Handling

#### Description

The information system:
<ol type="a">
<li>Generates error messages that provide information necessary for corrective actions without revealing information that could be exploited by adversaries; and</li>
<li>Reveals error messages only to [Assignment: organization-defined personnel or roles].</li>
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
<li class="active"><a data-toggle="tab" data-target="#b70lu8hidlmg00bko3o0">DTR</a></li>
<li><a data-toggle="tab" data-target="#b70lu8hidlmg00bko3og">Engine</a></li>
<li><a data-toggle="tab" data-target="#b70lu8hidlmg00bko3p0">UCP</a></li>
</ul>

<div class="tab-content">
<div id="b70lu8hidlmg00bko3o0" class="tab-pane fade in active">
All error messages generated via the configured logging mechanism of
Docker Trusted Registry are displayed such that they meet the
requirements of this control. Only users that are authorized the
appropriate level of access can view these error messages.
</div>
<div id="b70lu8hidlmg00bko3og" class="tab-pane fade">
All error messages generated via the logging mechanisms of the Docker
Enterprise Edition engine are displayed such that they meet the
requirements of this control. Only users that are authorized the
appropriate level of access can view these error messages.
</div>
<div id="b70lu8hidlmg00bko3p0" class="tab-pane fade">
All error messages generated via the configured logging mechanism of
Universal Control Plane are displayed such that they meet the
requirements of this control. Only users that are authorized the
appropriate level of access can view these error messages.
</div>
</div>

## SI-12 Information Handling And Retention

#### Description

The organization handles and retains information within the information system and information output from the system in accordance with applicable federal laws, Executive Orders, directives, policies, regulations, standards, and operational requirements.

#### Control Information

Responsible role(s) - Organization

## SI-13 Predictable Failure Prevention

#### Description

The organization:
<ol type="a">
<li>Determines mean time to failure (MTTF) for [Assignment: organization-defined information system components] in specific environments of operation; and</li>
<li>Provides substitute information system components and a means to exchange active and standby components at [Assignment: organization-defined MTTF substitution criteria].</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-13 (1) Transferring Component Responsibilities

#### Description

The organization takes information system components out of service by transferring component responsibilities to substitute components no later than [Assignment: organization-defined fraction or percentage] of mean time to failure.

#### Control Information

Responsible role(s) - Organization

### SI-13 (3) Manual Transfer Between Components

#### Description

The organization manually initiates transfers between active and standby information system components [Assignment: organization-defined frequency] if the mean time to failure exceeds [Assignment: organization-defined time period].

#### Control Information

Responsible role(s) - Organization

### SI-13 (4) Standby Component Installation / Notification

#### Description

The organization, if information system component failures are detected:
<ol type="a">
<li>Ensures that the standby components are successfully and transparently installed within [Assignment: organization-defined time period]; and</li>
<li>[Selection (one or more): activates [Assignment: organization-defined alarm]; automatically shuts down the information system].</li>
</ol>

#### Control Information

Responsible role(s) - Organization

### SI-13 (5) Failover Capability

#### Description

The organization provides [Selection: real-time; near real-time] [Assignment: organization-defined failover capability] for the information system.

#### Control Information

Responsible role(s) - Organization

## SI-14 Non-Persistence

#### Description

The organization implements non-persistent [Assignment: organization-defined information system components and services] that are initiated in a known state and terminated [Selection (one or more): upon end of session of use; periodically at [Assignment: organization-defined frequency]].

#### Control Information

Responsible role(s) - Organization

### SI-14 (1) Refresh From Trusted Sources

#### Description

The organization ensures that software and data employed during information system component and service refreshes are obtained from [Assignment: organization-defined trusted sources].

#### Control Information

Responsible role(s) - Organization

## SI-15 Information Output Filtering

#### Description

The information system validates information output from [Assignment: organization-defined software programs and/or applications] to ensure that the information is consistent with the expected content.

#### Control Information

Responsible role(s) - Organization

## SI-16 Memory Protection

#### Description

The information system implements [Assignment: organization-defined security safeguards] to protect its memory from unauthorized code execution.

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
<td>service provider hybrid<br/></td>
</tr>
</table>

#### Implementation Details

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#b70lu8hidlmg00bko3pg">Engine</a></li>
</ul>

<div class="tab-content">
<div id="b70lu8hidlmg00bko3pg" class="tab-pane fade in active">
Docker Enterprise Edition can be installed on the following operating
systems: CentOS 7.1&#43;, Red Hat Enterprise Linux 7.0&#43;, Ubuntu 14.04
LTS&#43;, SUSE Linux Enterprise 12&#43; and Windows Server 2016&#43;. In order to
meet the requirements of this control, reference the chosen operating
system&#39;s security documentation for information regarding the
protection of memory from unauthorized code execution.
</div>
</div>

## SI-17 Fail-Safe Procedures

#### Description

The information system implements [Assignment: organization-defined fail-safe procedures] when [Assignment: organization-defined failure conditions occur].

#### Control Information

Responsible role(s) - Organization

