---
title: Assign roles and permissions
description:
weight: 20
---
Managing roles and permissions in Docker Business is about more than just granting access—it’s about balancing security, collaboration, and operational efficiency. This guide offers a practical framework for navigating roles and permissions, providing actionable steps, use cases, and pathways to relevant Docker documentation for deeper insight.

Why Roles and Permissions Matter
Roles and permissions are the backbone of secure access in Docker Business. They help you:

Maintain control over sensitive settings.
Simplify collaboration within teams.
Prevent unauthorized access to critical resources.
Instead of overwhelming users with technical details, this guide highlights how to use roles and permissions strategically and points you to the tools and configurations available in Docker.

→ Discover the foundational principles of roles and permissions: Roles and Permissions Overview

Getting Started with Roles
Docker Business roles are simple yet powerful tools for managing user access. Each role—Admin, Developer, Viewer—has specific permissions tailored to different organizational needs.

The Power of Role Templates
Think of roles as templates. Instead of manually configuring access for every user:

Assign Admins for organizational or team-level oversight.
Use the Developer role for users actively building and deploying projects.
Default contractors or external stakeholders to Viewer, ensuring they can observe without modifying resources.
→ Find out which role fits your use case: Docker Roles Documentation

How to Use Teams for Simplified Access Control
While roles define access, teams structure it. Teams let you group users logically—by project, department, or function—and apply permissions collectively. This approach reduces errors and makes scaling permissions straightforward.

Making Teams Work for You
Align Teams with Projects: Create dedicated teams for individual projects and assign specific resources (e.g., repositories, containers) to those teams.
Cross-Team Collaboration: For multi-department initiatives, grant overlapping access to shared resources while keeping other teams isolated.
Dynamic Scaling: Add or remove users from teams as roles evolve, ensuring minimal disruption.
→ Learn to create and manage teams effectively: Docker Team Management

Scaling Permissions with Owners and Admins
Organizational roles like Owner and Admin provide a broader layer of oversight. Think of them as the framework within which teams operate.

When to Use Owners and Admins
Assign Owner roles sparingly to manage billing, organizational setup, and global security.
Delegate Admin roles for day-to-day team and resource management, ensuring a distributed yet secure approach.
→ Explore the responsibilities of Owners and Admins: Owners and Admins Explained

Practical Scenarios for Roles and Permissions
Scenario 1: New Project Kickoff
Goal: Quickly onboard a team of developers for a new project.
Solution: Create a team, assign the Developer role, and map the team to the relevant repositories. Use granular permissions to ensure they only access necessary resources.

Scenario 2: Adding External Auditors
Goal: Provide temporary access to external auditors without compromising security.
Solution: Assign auditors to a Viewer-only team with access restricted to a read-only repository.

Scenario 3: Delegating Team Oversight
Goal: Empower team leads to manage their own teams while maintaining organizational security.
Solution: Assign the Admin role to team leads, allowing them to oversee permissions and users within their domain.

→ Find more best practices for managing roles and permissions: Roles and Permissions Use Cases

Actionable Tips for Managing Permissions
Audit Regularly: Use Docker’s activity logging tools to monitor role changes and permissions updates.
→ Explore activity logging

Automate Assignments: Use group-based permissions to streamline access control, reducing manual errors.
→ Learn about group mapping

Document Your Strategy: Maintain an internal record of roles and permissions, especially for cross-functional teams and external collaborators.



Best Practices for Role Assignment
Principle of Least Privilege: Assign the minimum level of access needed for users to perform their tasks. For example:

Developers working on staging environments might not need admin-level access.
Contractors may require Viewer access only.
Team-Based Assignments: Use teams to group users and assign roles collectively rather than individually. This simplifies management as your organization scales.

