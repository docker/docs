---
title: Kit capability support
linkTitle: Capability support
description: Check which v3 kit capabilities Docker Sandboxes supports locally and in the cloud, their minimum sbx versions, and restrictions that affect kit authors.
keywords: sandboxes, sbx, kits, capabilities, compatibility, cloud
weight: 37
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Use this reference to check whether a v3 kit can run on your version of Docker
Sandboxes and whether it needs a local sandbox. Follow each specification reference
for its settings and examples in the Docker Sandbox Kit Specification.

## Support by capability

Version numbers identify the first `sbx` release that supports the
capability, subject to its restrictions. Limited support means that Docker
Sandboxes applies only part of the capability. Declaration only means that
Docker Sandboxes accepts the entry without using its commands.

All capability names have the `com.docker.sandbox/` prefix. Not verified
means Docker Sandboxes support has not been confirmed for that capability.

{{% kit_capabilities.inline %}}
{{- $catalog := hugo.Data.sandbox_kit_capabilities -}}
{{- $support := hugo.Data.sandbox_kit_support -}}
{{- $repository := "https://github.com/docker/sandbox-kit-spec" -}}
{{- if or (ne $catalog.format_version 1) (ne $catalog.schema_version "3") (not (findRE `^v3\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$` $catalog.version)) (not (findRE `^[0-9a-f]{40}$` $catalog.revision)) -}}
  {{- errorf "kit-capabilities: invalid catalog format, schema, version, or revision: %s" .Position -}}
{{- end -}}
{{- if or (not (reflect.IsSlice $catalog.capabilities)) (not $catalog.capabilities) -}}
  {{- errorf "kit-capabilities: catalog must contain capabilities: %s" .Position -}}
{{- end -}}
{{- $ids := newScratch -}}
{{- range $catalog.capabilities -}}
  {{- if or (not (findRE `^com\.docker\.sandbox/[a-z0-9][a-z0-9-]*@[1-9][0-9]*$` .id)) ($ids.Get .id) -}}
    {{- errorf "kit-capabilities: invalid or duplicate catalog ID %q: %s" .id $.Position -}}
  {{- end -}}
  {{- $ids.Set .id true -}}
  {{- if or (ne (printf "%T" .description) "string") (not (strings.TrimSpace .description)) -}}
    {{- errorf "kit-capabilities: missing description for %q: %s" .id $.Position -}}
  {{- end -}}
  {{- if or (ne .spec_url (printf "%s/blob/%s/docs/spec/capabilities/%s.md" $repository $catalog.version .id)) (ne .schema_url (printf "%s/blob/%s/schema/capabilities/%s.schema.json" $repository $catalog.version .id)) -}}
    {{- errorf "kit-capabilities: unpinned specification links for %q: %s" .id $.Position -}}
  {{- end -}}
{{- end -}}
{{- if or (ne $support.format_version 1) (not (reflect.IsSlice $support.capabilities)) -}}
  {{- errorf "kit-capabilities: invalid docs-owned support format: %s" .Position -}}
{{- end -}}
{{- $supportByID := newScratch -}}
{{- range $entry := $support.capabilities -}}
  {{- $id := $entry.id -}}
  {{- if not ($ids.Get $id) -}}
    {{- errorf "kit-capabilities: orphaned support ID %q needs review: %s" $id $.Position -}}
  {{- end -}}
  {{- if $supportByID.Get $id -}}
    {{- errorf "kit-capabilities: duplicate support ID %q: %s" $id $.Position -}}
  {{- end -}}
  {{- if or (not (reflect.IsMap $entry)) (ne (len $entry) 4) -}}
    {{- errorf "kit-capabilities: support ID %q needs id, local, cloud, and restrictions: %s" $id $.Position -}}
  {{- end -}}
  {{- $supportByID.Set $id $entry -}}
  {{- range $field := slice "local" "cloud" -}}
    {{- $value := index $entry $field -}}
    {{- if or (ne (printf "%T" $value) "string") (not (findRE `^(Not verified|Unsupported|Declaration only|(Limited: )?v[0-9]+\.[0-9]+\.[0-9]+)$` $value)) -}}
      {{- errorf "kit-capabilities: malformed %s support for %q: %s" $field $id $.Position -}}
    {{- end -}}
  {{- end -}}
  {{- if or (ne (printf "%T" $entry.restrictions) "string") (not (strings.TrimSpace $entry.restrictions)) (findRE `[\r\n]` $entry.restrictions) -}}
    {{- errorf "kit-capabilities: malformed restrictions for %q: %s" $id $.Position -}}
  {{- end -}}
{{- end }}

Specification descriptions and references come from the
[Docker Sandbox Kit Specification {{ $catalog.version }}]({{ $repository }}/releases/tag/{{ $catalog.version }}).
Docker Sandboxes support is tracked separately.

{{ range $catalog.capabilities -}}
  {{- $shortID := strings.TrimPrefix "com.docker.sandbox/" .id -}}
  {{- $entry := ($supportByID.Get .id) | default (dict "local" "Not verified" "cloud" "Not verified" "restrictions" "Docker Sandboxes support has not been verified.") }}

### {{ $shortID }}

{{ .description }}

| Local sandboxes | Cloud sandboxes |
| :-------------- | :-------------- |
| {{ $entry.local }} | {{ $entry.cloud }} |
{aria-label="Runtime support for {{ $shortID }}"}

{{ $entry.restrictions }}

Specification reference: [{{ $shortID }}]({{ .spec_url }}).

{{ end -}}
{{% /kit_capabilities.inline %}}

## Port address families

For v3 `port@1` entries, `family` selects the address families used to publish
the port: `ipv4`, `ipv6`, or `dual` for both. If no kit specifies a family,
the port defaults to IPv4, even with an explicit `transport: tcp`. Use
`family: dual` when the service needs both IPv4 and IPv6 clients.

When kits declare the same port and transport, an entry without `family`
uses the family specified by another kit in the set. Kits that specify
conflicting families for that port cannot be used together.

Cloud sandboxes support IPv4 TCP ports only. They reject required `ipv6` and
`dual` entries and skip optional ones. For capability syntax and
examples, follow the specification reference for [port@1](#port1).

## Capability groups

Docker Sandboxes supports capability groups in `sbx` v0.48.0 and later.
A group is selected only when the runtime can provide every capability in it.
If a capability is unavailable, an optional group is skipped in full, including
its network rules, setup commands, generated files, and agent instructions.
An unavailable capability in a required group prevents sandbox creation.

Availability also depends on the host: Git identity requires configured
`user.name` and `user.email`, and SSH-agent access requires an available
agent with forwarding enabled. See the upstream
[capability group definition](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/SPEC-v3.md#711-capability-groups)
for syntax and composition rules.

## Capabilities that require sandbox creation

Include kits that need the following settings when you create the sandbox.
`sbx kit add` cannot apply these settings to an existing sandbox:

| Capability | Settings that require creation |
| --- | --- |
| `credential@1` | Credential routing and authentication |
| `kit-registry@1` | Access to the kit registry |
| `long-running@1` | Keeping the sandbox running after sessions disconnect |
| `port@1` | Published ports |
| `resources@1` | CPU and memory allocation |
| `volume@1` | Block storage and tmpfs mounts |
| `lifecycle@1` | Generated `files` require sandbox creation. Install and startup hooks can be added with `sbx kit add`. |
| `network-policy@1` and `network-policy@2` | Host-level deny rules require sandbox creation. Allow rules can be added with `sbx kit add`, as can HTTP method and path rules for `network-policy@2`. |

`sbx kit add` rejects kits that declare credentials, published ports,
resource settings, volumes, generated lifecycle files, or host-level deny
rules. To use those kits, remove and recreate the sandbox with the kits
included. Save any work you need from the sandbox before removing it.

Skill discovery paths (`agent-skills@1`) and bundled skills (`agent-skill@1`)
can be added with `sbx kit add`: the command replaces the sandbox container
and applies the updated skill configuration.
