{{- partial "markdown/render.md" (dict "page" . "content" (partial "markdown/body/installation.md" .)) -}}
{{ define "_partials/markdown/body/installation.md" -}}
# {{ .Title }}

{{ .Description }}

{{ range .Pages }}- [{{ .LinkTitle }}]({{ .Permalink }}): {{ .Description }}
{{ end }}
{{ end -}}
