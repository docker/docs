<!--[metadata]>
+++
title ="Universal Control Plane"
keywords= ["universal, control, plane, ucp"]
description="Docker Universal Control Plane"
[menu.main]
identifier="mn_ucp"
+++
<![end-metadata]-->

# Docker Universal Control Plane

Docker Universal Control Plane (UCP) is an enterprise on premise solution that
enables IT operations teams to deploy and manage their Dockerized applications
in production. It gives developers and administrators the agility and
portability they need to manage Docker resources, all from within the enterprise
firewall.

## Deploy, manage, and monitor Docker Engine resources

Universal Control Plane can be deployed to any private infrastructure and public
cloud including Microsoft Azure, Digital Ocean, Amazon Web Services, and
SoftLayer.

Once deployed, UCP uses Docker Swarm to create and manage clusters, tested up to
10,000 nodes deployed in any private data center or public cloud provider.

## Secure communications

UCP has built in in security, and integration with existing LDAP/AD for
authentication and role based access control. Optionally, you can use its native
integration with Docker Trusted registry. The integration with Docker Trusted
Registry allows enterprises to leverage Docker Content Trust (Notary in the
open source world), a built in security tool for signing images.


Universal Control Plane is the only tool on
the market that comes comes with Docker Content Trust directly
out of the box. With these integrations Universal Control Plane
gives enterprise IT security teams the necessary control over their
environment and application content.

## Control user access

Security is top of mind for many enterprise IT operations teams. Docker UCP
integrates with existing tools like LDAP/AD for user authentication and its
integration with Docker Trusted Registry. This integration enables enterprises
to control the entire build, ship and run workflow in a secure fashion.

Within Docker UCP, you can do a local set up for account information, or you can
do centralized authentication by linking UCP with your LDAP or Active Directory.
The integration with Docker Trusted Registry also means that you can use Docker
Content Trust, to sign your images and ensure that they are not altered in
anyway and are safe for use within your organization. Users can pull images from
Docker Trusted Registry into Docker UCP and not have to worry about their
security.

For RBAC, within UCP you can view the roles of existing accounts as well as see
the roles that they have within UCP. The granular role based access control
allows you to control who can access certain images. This drastically reduces
organizational risk within enterprises.

## Where to go next

* If you are interested in evaluating UCP on your laptop, you can use the [evaluation installation and quickstart](evaluation-install.md).

* Technical users and managers can get detailed explanation of the UCP architecture and requirements from the [plan a production installation](plan-production-install.md) page. The step-by-step [production installation](plan-production-install.md) is also available.

* To learn more about controlling users in UCP, see [Manage and authorize UCP users](manage/monitor-manage-users.md).

* If you have used UCP in our BETA program, be sure to read the [release notes](release_notes.md).
