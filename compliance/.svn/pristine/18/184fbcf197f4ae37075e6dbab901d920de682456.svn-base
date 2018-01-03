---
title: NIST 800-53 control reference
---

All of the NIST 800-53 Rev. 4 controls applicable to Docker Enterprise Edition can be referenced in this section. For generating your own security documentation using the OpenControl-formatted version of these controls, please refer to our compliance repository at [https://github.com/docker/compliance](https://github.com/docker/compliance).

The controls have been broken out by family and each control's origin is mapped to one of the following:

|Control Origination|Definition|Example|
|-------------------|----------|-------|
|Service provider corporate|A control that originates from agency's corporate network|DNS from the corporate network provides address resolution services for the information system and the service offering|
|Docker EE system|A control specific to Docker EE|Docker EE LDAP configuration|
|Service provider hybrid|A control that makes use of both corporate controls and additional controls specific to Docker EE|There are scans of the corporate network infrastructure; scans of Docker images via DTR would be included|
|Configured by customer|A control where the Docker EE end-user's application needs to apply a configuration in order to meet the control requirement|User profiles, policy/audit configurations, enable/disabling key switches (e.g., enable/disable http or https, etc), entering an IP range specific to the end-user's organization are configurable by the customer|
|Provided by customer|A control where the Docker EE end-user's application needs to provide additional hardware or software in order to meet the control requirement|The customer provides a SAML SSO solution to implement two-factor authentication|
|Shared|A control that is managed and implemented partially by the Docker EE system and partially by the Docker EE end-user|Security awareness training must be conducted by both the Docker EE operators and end-users|
|Inherited from pre-existing Provisional Authorization|A control that is inherited from another CSP system that has already received a Provisional Authorization|Docker EE inherites PE controls from an IaaS provider|

The following Docker EE system components are referenced by the controls:

- Docker EE Engine
- Universal Control Plane (UCP)
- Docker Trusted Registry (DTR)
- Authentication and Authorization Service (eNZi)
- Docker Security Scanning (DSS)

In addition, each control is assigned one or more of the following implementation statuses:

|Implementation status|Definition|
|---------------------|----------|
|Complete|The control is fully in place and meets all requirements|
|Partial|The control is only partially in place or does not meet all requirements. A plan for achieving full implementation should be included in the Plan of Action & Milestone documentation|
|Planned|The control is not in place. A plan for achieving full implementation should be included in the Plan of Action & Milestone documentation|
|None|The control is not applicable within the environment. A description of why the requirement does not apply should be included|

Control narratives that include an `[Assignment: ...]` block should be substituted by your organization's requirements or by the FedRAMP requirements.
