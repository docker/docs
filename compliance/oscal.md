---
description: OSCAL compliance guidance for Docker Enterprise
keywords: standards, compliance, security, benchmark
title: OSCAL (Experimental)
---

The Open Security Controls Assessment Language (OSCAL) is a set of standardized XML- and JSON-based formats developed 
by the National Institute of Standards and Technology (NIST) for reporting, implementing, and 
assessing system security controls. Docker Enterprise contains OSCAL integrations for automating  
attestation and reporting requirements of your container platform. Use the OSCAL integrations to automatically:

* Assess system compliance.
* Correct system configurations that have drifted and are non-compliant.
* Download Docker’s officially signed and maintained OSCAL-formatted attestation artifacts that 
describe how Docker’s products adhere to and can be configured to meet compliance standards and hardening benchmarks.

These OSCAL auditing and reporting capabilities are exposed via OSCAL API endpoints built into UCP and  
enabled by default. Capabilities also include the majority of the recommendations in the CIS Kubernetes Benchmark, along with a 
subset of the CIS Docker Benchmark checks that are currently performed by the external `docker-bench-security` tool. 

> **Note:** Docker App is an experimental feature.

     - Experimental features provide early access to future product functionality. These 
     features are only intended for testing and feedback only as they may change between 
     releases without warning or can be removed entirely from a future release. Experimental 
     features must not be used in production environments. Docker does not offer support 
     for experimental features. For more information, see 
     [Experimental features](https://success.docker.com/article/experimental-features).
     - Not all of the CIS Docker Benchmark checks have been incorporated into the OSCAL functions 
     in UCP at the time of this “tech preview” release. Continue to use the `docker-bench-security tool` 
     in concert with the UCP endpoints.   

In terms of the data they expose, the OSCAL API endpoints build upon one another as they traverse the Docker Enterprise 
stack. At the highest stack layer, with the entire Docker Enterprise platform (Engine, UCP, DTR), 
the greatest number of security controls are available. For example:

* Verify FIPS mode is enabled on all Engine nodes.
* Verify LDAP is enabled and configured in UCP.
* Verify vulnerability scanning is enabled in DTR. 

**LIMITATIONS**

If an Engine node goes down or a node health check fails, this is reported as a risk for any security controls 
that execute checks against the Engine nodes. For example, if FIPS mode must be enabled on all nodes in the 
cluster per the requirements of a particular security control, and a cluster node goes down, the cluster node failure
is reported for that security control. However, no action is taken on the failed node. 


**DISCLAIMER**

This content is provided for informational purposes only and has not been vetted by any 
third-party security assessors. You are solely responsible for developing, implementing, and managing your 
applications and/or subscriptions running on your own platform in compliance with applicable laws, regulations, 
and contractual obligations. The documentation is provided "as-is" and without any warranty of any kind, whether 
express, implied or statutory, and Docker, Inc. expressly disclaims all warranties for non-infringement, 
merchantability or fitness for a particular purpose.
