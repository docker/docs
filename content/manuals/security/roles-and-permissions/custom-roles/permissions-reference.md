---
title: Custom role permissions reference
linkTitle: Permissions reference
description: >-
  Permissions available for Docker custom roles across organization management,
  Docker Hub, billing, licensing, AI Governance, Docker Hardened Images,
  Docker Offload, and Docker Cloud Sandboxes.
keywords: >-
  custom roles, custom role permissions, Docker, Docker Hub, organization
  management, billing, licensing, licenses, AI Governance, access tokens, SSO,
  SCIM, OIDC, DHI, Docker Offload, Cloud Sandboxes, security
weight: 20
aliases:
  - /enterprise/security/roles-and-permissions/custom-roles/permissions-reference/
---

{{< summary-bar feature_name="Custom roles" >}}

Custom roles use permissions from organization management, Docker Hub,
billing, licensing, AI Governance, Docker Hardened Images, Docker
Offload, and Docker Cloud Sandboxes. Use the following tables to [create
or edit a custom role](manage.md).

{{% custom_role_permissions.inline %}}
{{- $data := hugo.Data.custom_role_permissions -}}
{{- if not $data.groups -}}
  {{- errorf "custom-role-permissions: no groups found in data/custom_role_permissions.yaml: %s" .Position -}}
{{- end -}}
{{- range $group := $data.groups -}}
  {{- if not $group.name -}}
    {{- errorf "custom-role-permissions: group without a name: %s" $.Position -}}
  {{- end -}}
  {{- if $group.relations }}

## {{ $group.name }}

| Permission | Description |
| :--------- | :---------- |
{{ range $group.relations -}}
  {{- if or (not .name) (not .description) -}}
    {{- errorf "custom-role-permissions: relation %q in group %q needs a name and description: %s" .name $group.name $.Position -}}
  {{- end -}}
| {{ replace .name "|" "\\|" }} | {{ replace .description "|" "\\|" }} |
{{ end -}}
  {{- end -}}
{{- end }}
{{% /custom_role_permissions.inline %}}

## Next steps

- [Manage custom roles](manage.md): Create, assign, and delete custom
  roles
- [Core roles and permissions](/manuals/security/roles-and-permissions/core-roles.md):
  Compare built-in Member, Editor, and Owner permissions
