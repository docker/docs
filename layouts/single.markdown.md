{{- partial "markdown/render.md" (dict "page" . "content" (partial "markdown/body/single.md" .)) -}}
{{ define "_partials/markdown/body/single.md" -}}
# {{ .Title }}

{{ partial "markdown/content.html" . }}
{{ end -}}
