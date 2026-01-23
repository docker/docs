{{- $specURL := urls.Parse (printf "/%s%s.yaml" .File.Dir .File.ContentBaseName) -}}
# {{ .Title }}

{{ .Content }}

**OpenAPI Specification:** [{{ .Title }} API Spec]({{ $specURL.String | absURL }})

This page provides interactive API documentation. For the machine-readable OpenAPI specification, see the link above.
