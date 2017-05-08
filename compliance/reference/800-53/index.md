---
title: NIST 800-53 Control Reference
---

All of the NIST 800-53 Rev. 4 controls applicable to Docker Enterprise Edition can be referenced in this section. For generating your own security documentation using the OpenControl-formatted version of these controls, please refer to our compliance repository at [https://github.com/docker/compliance](https://github.com/docker/compliance).

The controls have been broken out by family and each control's origin is mapped to one of the following:

| Control origination                                                  | Definition                                                                                                     |
|:------------------------------------------------------|:---------------------------------------------------------------------------------------------------------------|
| Organization                                          | A control satisfied by the organization.                                                                       |
| Docker Enterprise Edition (EE) system                      | A control specific to the Docker Enterprise Edition system itself                                              |
| Docker Enterprise Edition (EE) application                 | A control specific to an application(s) hosted on Docker Enterprise Edition                                    |
| Inherited from pre-existing Provisional Authorization | A control that is inherited from another provider system that has already received a Provisional Authorization |

The following Docker EE system components are referenced by the controls:

- Docker EE Engine
- Universal Control Plane (UCP)
- Docker Trusted Registry (DTR)
- Authentication and Authorization Service (eNZi)
- Docker Security Scanning (DSS)

In addition, each control is assigned one of the following implementation statuses:

| Implementation status                                                  | Definition                                                                                    |
|:------------------------------------------------------|:---------------------------------------------------------------------------------------------------------------|
| Complete               | The control is fully in place and meets all requirements |
| Partial      | The control is only partially in place or does not meet all requirements. A plan for achieving full implementation should be included in the Plan of Action & Milestone documentation |
| Planned                    | The control is not in place. A plan for achieving full implementation should be included int he Plan of Action & Milestone documentation |
| None             | The control is not applicable within the environment. A description of why the requirement does not apply should be included |

Control narratives that include an `[Assignment: ...]` block should be substituted by your organization's requirements or by the FedRAMP requirements.
