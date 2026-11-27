{{- partial "markdown/render.md" (dict "page" . "content" (partial "markdown/body/list.md" .)) -}}
{{ define "_partials/markdown/body/list.md" -}}
# {{ .Title }}

{{ partial "markdown/content.html" . }}
{{ end -}}
