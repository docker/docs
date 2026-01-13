---
title: {{ .Title }}
url: {{ .Permalink }}
{{- range .Ancestors }}
  {{- if and (not .IsHome) .Permalink }}
parent:
  title: {{ .Title }}
  url: {{ .Permalink }}
  {{- break }}
  {{- end }}
{{- end }}
{{- if .Ancestors }}
breadcrumbs:
{{- range .Ancestors.Reverse }}
  {{- if and (not .IsHome) .Permalink }}
  - title: {{ .Title }}
    url: {{ .Permalink }}
  {{- end }}
{{- end }}
  - title: {{ .Title }}
    url: {{ .Permalink }}
{{- end }}
{{- with .NextInSection }}
next:
  title: {{ .Title }}
  url: {{ .Permalink }}
{{- end }}
{{- with .PrevInSection }}
prev:
  title: {{ .Title }}
  url: {{ .Permalink }}
{{- end }}
{{- $specURL := urls.Parse (printf "/%s%s.yaml" .File.Dir .File.ContentBaseName) }}
openapi_spec: {{ $specURL.String | absURL }}
---

{{ .Content }}

**OpenAPI Specification:** [{{ .Title }} API Spec]({{ $specURL.String | absURL }})

This page provides interactive API documentation. For the machine-readable OpenAPI specification, see the link above.
