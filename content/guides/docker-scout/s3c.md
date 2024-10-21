---
title: Software supply chain security
description: Learn about software supply chain security (S3C), what it means, and why it is important.
keywords: docker scout, secure, software, supply, chain, security, sssc, sscs, s3c
aliases:
  - /scout/concepts/s3c/
weight: 30
---

{{< youtube-embed YzNK6E7APv0 >}}

The term "software supply chain" refers to the end-to-end process of developing
and delivering software, from the development to deployment and maintenance.
Software supply chain security, or "S3C" for short, is the practice for
protecting the components and processes of the supply chain.

S3C is a fundamental change in how organizations approach software security.
Traditionally in the software industry, security and compliance has been mostly
an afterthought, left to the software delivery or release phase. With S3C,
security is integrated into the entire software development lifecycle, from the
inner loop of development and testing, to the outer loop of shipping and
monitoring.

Following industry best practices for software supply chain conduct is
important because it helps organizations protect their software from security
threats, compliance risks, and other vulnerabilities. Implementing a software
supply chain security framework improves visibility, collaboration, and
traceability of a project across stakeholders. This helps organizations detect,
respond to, and remediate threats more effectively.

## Securing the software supply chain

Building a secure software supply chain involves several key steps, such as:

- Identify the software components and dependencies you use to build and run
  your applications.
- Automate security testing throughout the software development lifecycle.
- Monitor your software supply chain for security threats.
- Implement security policies that govern how software is built, and the
  components it contains.

Managing the software supply chain is a complex task, especially in the modern
day where software is built using multiple components from different sources.
Organizations need to have a clear understanding of the software components
they use, and the security risks associated with them.

## How Docker Scout is different

Docker Scout is a platform designed to help organizations secure their software
supply chain. It provides tools and services for identifying and managing
software assets and policies, and automated remediation of security threats.

Unlike traditional security tools that focus on scheduled, point-in-time scans
at specific stages in the software development lifecycle, Docker Scout uses a
modern event-driven model that spans the entire software supply chain. This
means that when a new vulnerability affecting your images is disclosed, your
updated risk assessment is available within seconds, and earlier in the
development process.

Docker Scout works by analyzing the composition of your images to create a
Software Bill of Materials (SBOM). The SBOM is cross-referenced against the
security advisories to identify CVEs that affect your images. Docker Scout
integrates with [over 20 different security
advisories](/manuals/scout/deep-dive/advisory-db-sources.md), and updates its
vulnerability database in real-time. This ensures that your security posture is
represented using the latest available information.

<div id="scout-lp-survey-anchor"></div>
