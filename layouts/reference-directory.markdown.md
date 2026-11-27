{{- partial "markdown/render.md" (dict "page" . "content" (partial "markdown/body/reference-directory.md" .)) -}}
{{ define "_partials/markdown/body/reference-directory.md" -}}
# {{ .Title }}

{{ .Description }}

{{ range hugo.Data.navigation.references }}{{ $label := .label }}{{ with site.GetPage .root }}{{ if ne .Params.sitemap false }}- [{{ $label }}]({{ partial "navigation/url.html" . | absURL }}): {{ .Description }}
{{ end }}{{ end }}{{ end }}
{{ end -}}
