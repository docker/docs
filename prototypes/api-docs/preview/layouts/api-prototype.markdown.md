{{- $data := index hugo.Data "api-prototype" -}}
{{- $api := dict -}}{{- range $data.apis -}}{{- if eq .id $.Params.apiID -}}{{- $api = . -}}{{- end -}}{{- end -}}
# {{ .Title }}

Prototype: converted specifications contain documented assumptions pending product review.

{{ if eq .Params.view "catalog" }}
Choose a Docker HTTP API:
{{ range $data.apis }}
- [{{ .title }} API {{ .version }}]({{ .url }}): {{ len .operations }} operations; {{ .connection }} connection
{{ end }}
{{ else }}
[API catalog](/api-prototype/) · [{{ $api.title }} overview]({{ $api.url }}) · [Product manual]({{ ref . $api.manual }}) · [Source package](/api-prototype/sources/{{ $api.id }}/source.zip)

API version: {{ $api.version }}

{{ if eq .Params.view "overview" }}
## Overview

{{ partial "api-prototype/description.html" (dict "text" $api.description "api" $api) }}
## {{ if eq $api.connection "unix" }}Connecting to {{ $api.title }}{{ else }}Connecting to the {{ $api.title }} API{{ end }}

{{ $api.auth }}
{{ range $api.servers }}
Server: `{{ .url }}`
{{ end }}
{{ if eq $api.connection "unix" }}
```console
curl --unix-socket /var/run/docker.sock http://localhost/v{{ $api.version }}/version
```
{{ end }}
{{ range $api.guides }}
{{ $guide := site.GetPage (index (split . "#") 0) }}
- [{{ $guide.Title }}]({{ ref $ . }})
{{ end }}

{{ range (partial "api-prototype/overview-tags.html" $api) }}
## {{ .summary }}

{{ partial "api-prototype/description.html" (dict "text" .description "api" $api) }}
{{ end }}
## Operations
{{ range $api.operations }}
- [{{ .method }} {{ .path }}]({{ .url }}): {{ .summary }}
{{ end }}
## Schemas
{{ range $api.schemas }}
- [{{ .name }}]({{ .url }})
{{ end }}
## Source review

Owner: {{ $api.owner }}. Digest: `{{ $api.digest }}`.

[Migration ledger](/api-prototype/sources/{{ $api.id }}/migration.json) · [Validation exceptions](/api-prototype/sources/{{ $api.id }}/exceptions.json)
{{ else if eq .Params.view "operation" }}
{{ range $api.operations }}{{ if eq .id $.Params.operationID }}
`{{ .method }} {{ .path }}`

{{ partial "api-prototype/description.html" (dict "text" .description "api" $api) }}
{{ if .deprecated }}
Deprecated operation.
{{ end }}
## Connection and access

{{ $api.auth }}
{{ if and (eq $api.product "engine") (where .parameters "name" "X-Registry-Auth") }}
`X-Registry-Auth` delegates registry credentials and does not authenticate the daemon caller.
{{ end }}
{{ range .servers }}
Server: `{{ .url }}`
{{ end }}
Effective security: alternatives are OR; schemes within an alternative are AND. An empty array declares no HTTP authentication requirement.

```json
{{ .security | jsonify (dict "indent" "  ") }}
```
## Example request

Replace placeholders and provide the required credentials or request body.

```console
{{ .curl }}
```
{{ range .curlNotes }}
{{ . }}
{{ end }}
## Parameters
{{ range .parameters }}
### {{ .name }}

Location: {{ .in }}. Required: {{ if .required }}yes{{ else }}no{{ end }}.

{{ partial "api-prototype/description.html" (dict "text" .description "api" $api) }}

```json
{{ . | jsonify (dict "indent" "  ") }}
```
{{ end }}
## Request and responses
{{ range .variants }}
### {{ .direction }} {{ .status }} {{ .media }}

{{ partial "api-prototype/description.html" (dict "text" .description "api" $api) }}
{{ if not .media }}
No response content is declared.
{{ end }}
{{ if isset . "schema" }}
Schema:

```json
{{ .schema | jsonify (dict "indent" "  ") }}
```
{{ end }}
{{ if isset . "itemSchema" }}
Stream item schema:

```json
{{ .itemSchema | jsonify (dict "indent" "  ") }}
```
{{ end }}
{{ with .headers }}
Headers:

```json
{{ . | jsonify (dict "indent" "  ") }}
```
{{ end }}
{{ range .examples }}
{{ .name }}:
{{ if eq .valid false }}
Source example does not satisfy its schema; owner review is required.
{{ end }}
```json
{{ .value | jsonify (dict "indent" "  ") }}
```
{{ end }}
{{ end }}
## Complete operation contract

```json
{{ .raw | jsonify (dict "indent" "  ") }}
```

## Referenced schemas
{{ range .references }}
- {{ if .url }}[{{ .ref }}]({{ .url }}){{ else }}`{{ .ref }}`{{ end }}
{{ end }}
{{ end }}{{ end }}
{{ else if eq .Params.view "schema" }}
{{ range $api.schemas }}{{ if eq .name $.Params.schemaName }}
Schema constraints and annotations:

```json
{{ .schema | jsonify (dict "indent" "  ") }}
```
{{ end }}{{ end }}
{{ end }}
{{ end }}
