{{- /*
  Markdown rendering of an OpenAPI reference page.

  Mirrors the HTML api-reference layout: unmarshals the colocated YAML
  Page Resource into a plain Hugo map and walks the spec to produce a
  flat markdown document. Used by the "View Markdown" / "Copy Markdown"
  actions and the .md alternate link.
*/ -}}
{{- $specRes := .Resources.GetMatch "*.yaml" -}}
{{- $api := $specRes | transform.Unmarshal -}}
{{- $methods := slice "get" "post" "put" "patch" "delete" -}}
{{- $paths := $api.paths -}}
# {{ .Title }}

{{ with .Description }}{{ . }}{{ end }}

{{ with $api.info }}
{{- with .version }}- **Version**: `{{ . }}`
{{ end -}}
{{ end -}}
{{- range $api.servers }}- **Base URL**: `{{ .url }}`
{{ end }}
- **OpenAPI specification**: [`{{ path.Base $specRes.RelPermalink }}`]({{ $specRes.Permalink }})

{{ with $api.info }}{{ with .description }}{{ . }}{{ end }}{{ end }}

{{ with $api.components }}{{ with .securitySchemes }}
## Authentication

{{ range $name, $scheme := . }}
**`{{ $name }}`**{{ with $scheme.type }} — type: `{{ . }}`{{ end }}{{ with $scheme.scheme }}, scheme: `{{ . }}`{{ end }}{{ with $scheme.bearerFormat }}, bearer format: `{{ . }}`{{ end }}

{{ with $scheme.description }}{{ . }}{{ end }}
{{ end -}}
{{ end }}{{ end }}

{{- range $api.tags }}
## {{ .name | title }}

{{ with .description }}{{ . }}{{ end }}

{{- $tagName := .name -}}
{{- range $path, $item := $paths -}}
  {{- $sharedParams := slice -}}
  {{- with index $item "parameters" -}}
    {{- range . -}}
      {{- $p := partial "api-ref/resolve.html" (dict "api" $api "node" .) -}}
      {{- $sharedParams = $sharedParams | append $p -}}
    {{- end -}}
  {{- end -}}
  {{- range $methods -}}
    {{- $method := . -}}
    {{- with index $item $method -}}
      {{- $op := . -}}
      {{- if in $op.tags $tagName }}
### `{{ upper $method }}` `{{ $path }}`

{{ with $op.summary }}**{{ . }}**{{ end }}

{{ with $op.description }}{{ . }}{{ end }}

{{- $params := $sharedParams -}}
{{- with index $op "parameters" -}}
  {{- range . -}}
    {{- $p := partial "api-ref/resolve.html" (dict "api" $api "node" .) -}}
    {{- $params = $params | append $p -}}
  {{- end -}}
{{- end -}}
{{ if $params }}
**Parameters**

{{ range $params -}}
- `{{ .name }}` ({{ .in }}{{ with .schema }}{{ with .type }}, {{ . }}{{ end }}{{ end }}{{ if .required }}, required{{ end }}){{ with .description }} — {{ . | strings.TrimSpace }}{{ end }}
{{ end }}
{{- end }}

{{- with $op.requestBody -}}
  {{- $body := partial "api-ref/resolve.html" (dict "api" $api "node" .) }}
**Request body**{{ with $body.description }} — {{ . | strings.TrimSpace }}{{ end }}

{{ with index $body.content "application/json" -}}
{{- template "api-ref-md-schema-link" (dict "schema" .schema) }}
{{ range $exName, $ex := .examples }}{{ with $ex.value }}
```json
{{ . | jsonify (dict "indent" "  ") }}
```
{{ end }}{{ end }}
{{- end }}
{{- end }}

{{- with $op.responses }}
**Responses**

{{ range $code, $resp := . -}}
{{- $r := partial "api-ref/resolve.html" (dict "api" $api "node" $resp) }}
`{{ $code }}` — {{ with $r.description }}{{ . | strings.TrimSpace }}{{ end }}
{{ with index $r.content "application/json" }}
{{ template "api-ref-md-schema-link" (dict "schema" .schema) }}
{{ range $exName, $ex := .examples }}{{ with $ex.value }}
```json
{{ . | jsonify (dict "indent" "  ") }}
```
{{ end }}{{ end }}
{{- end }}
{{ end }}
{{ end -}}

      {{ end -}}
    {{- end -}}
  {{- end -}}
{{- end }}
{{- end }}

{{- with $api.components }}{{ with .schemas }}
## Schemas

{{ range $name, $schema := . }}
### `{{ $name }}`

{{ with $schema.description }}{{ . | strings.TrimSpace }}{{ end }}

{{ with $schema.enum -}}
Enum: {{ range $i, $v := . }}{{ if $i }}, {{ end }}`{{ $v }}`{{ end }}
{{ end }}

{{- with $schema.properties }}
{{- $required := $schema.required | default slice }}
| Property | Type | Required | Description |
| -------- | ---- | -------- | ----------- |
{{ range $propName, $prop := . -}}
| `{{ $propName }}` | {{ template "api-ref-md-type" $prop }} | {{ if in $required $propName }}yes{{ else }}no{{ end }} | {{ with $prop.description }}{{ . | strings.TrimSpace | strings.ReplaceRE "\\s*\\n\\s*" " " }}{{ end }} |
{{ end }}
{{ end -}}

{{ end -}}
{{ end }}{{ end }}

{{- /* ── Helpers ───────────────────────────────────────────────────────── */ -}}

{{- define "api-ref-md-type" -}}
  {{- $s := . -}}
  {{- if reflect.IsMap $s -}}
    {{- with index $s "$ref" -}}
      {{- $parts := split . "/" -}}
      {{- $name := index $parts (sub (len $parts) 1) -}}
      `{{ $name }}`
    {{- else -}}
      {{- with $s.type -}}
        {{- if eq . "array" -}}
          `array<`{{ template "api-ref-md-type" $s.items }}`>`
        {{- else -}}
          `{{ . }}`
        {{- end -}}
      {{- else -}}
        {{- if $s.allOf -}}`object`
        {{- else if or $s.anyOf $s.oneOf -}}`any`
        {{- else -}}—{{- end -}}
      {{- end -}}
    {{- end -}}
  {{- else -}}—{{- end -}}
{{- end -}}

{{- define "api-ref-md-schema-link" -}}
  {{- $schema := .schema -}}
  {{- $indent := .indent | default "" -}}
  {{- if reflect.IsMap $schema -}}
    {{- with index $schema "$ref" -}}
      {{- $parts := split . "/" -}}
      {{- $name := index $parts (sub (len $parts) 1) -}}
{{ $indent }}Schema: `{{ $name }}`
    {{- end -}}
  {{- end -}}
{{- end -}}
