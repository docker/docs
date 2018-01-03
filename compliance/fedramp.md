---
description: FedRAMP compliance guidance for Docker Enterprise Edition
keywords: standards, compliance, security
title: FedRAMP
---

The [Federal Risk and Authorization Management Program (FedRAMP)](https://www.fedramp.gov/) is a U.S. Federal Government-wide program that provides for a standardized approach to security assessment and authorization. Federal agencies that choose to leverage cloud services must ensure that they're utilizing FedRAMP-authorized providers. The security controls FedRAMP requires a provider to adhere to are a subset of the controls documented by NIST Special Publication 800-53. As with the baselines set within NIST 800-53 (low, moderate and high), FedRAMP also incorporates these same baselines in its authorization process. In addition, when agencies deploy systems (like Docker Enterprise Edition) on top of these providers, they must acquire an Authority to Operate (ATO) for those system that are in line with those agencies' own security procedures.

It is important to note that Docker, Inc is not a cloud service provider. While Docker does offer various SaaS-hosted services, which include Docker Hub, Docker Store and Docker Cloud, these services are *not* FedRAMP provisionally authorized. However, Docker's Enterprise product stack can be installed on top of compute services offered by a number of FedRAMP provisionally-authorized infrastructure-as-a-service (IaaS) providers. Examples include Microsoft Azure Government and Amazon Web Services GovCloud. Agencies can subsequently inherit the FedRAMP controls already satisfied by those providers and can combine those controls with the NIST 800-53 controls applicable to Docker Enterprise Edition and that which are documented on our site in order to gain an ATO for Docker Enterprise Edition.

Refer to the [NIST 800-53](/compliance/nist/800_53/) section for more information on the applicable NIST 800-53 controls.
