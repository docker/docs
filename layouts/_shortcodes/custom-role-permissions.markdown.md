{{- /*
  Renders the custom role permission tables from data/custom_role_permissions.yaml,
  which is synced from docker/iam-services by the "Sync custom role permissions"
  workflow. Call it with {{% custom-role-permissions %}} so the headings are
  parsed as Markdown and appear in the page's table of contents.
*/ -}}
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
