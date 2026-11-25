{{- partial "markdown/render.md" (dict "page" . "content" (partial "markdown/body/reference-directory.md" .)) -}}
{{ define "_partials/markdown/body/reference-directory.md" -}}
# {{ .Title }}

{{ .Description }}

{{ range partialCached "reference-directory/groups.html" . "reference-directory" }}

## {{ .title }}

{{ range .items }}- [{{ .title }}]({{ .link | absURL }}): {{ .description }}
{{ end }}{{ end }}
{{ end -}}
