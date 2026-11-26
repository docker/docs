{{- partial "markdown/render.md" (dict "page" . "content" (partial "markdown/body/category-directory.md" .)) -}}
{{ define "_partials/markdown/body/category-directory.md" -}}
# {{ .Title }}

{{ .Description }}

{{ range hugo.Data.navigation.categories }}
- [{{ .label }}]({{ .landing | absURL }}): {{ .description }}
{{ end }}
{{ end -}}
