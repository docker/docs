{{- partial "markdown/render.md" (dict "page" . "content" (partial "markdown/body/home.md" .)) -}}
{{ define "_partials/markdown/body/home.md" -}}
# Docker Documentation

{{ range hugo.Data.navigation.categories }}

## {{ .label }}

{{ .description }}

{{ range .groups }}{{ range .entries }}
{{ $path := .page | default "" }}{{ $label := .label | default "" }}
{{ with .collection }}{{ $product := index hugo.Data.navigation.products . }}{{ $path = $product.root }}{{ $label = $product.label }}{{ end }}
{{ with site.GetPage $path }}{{ if ne .Params.sitemap false }}- [{{ $label | default .LinkTitle }}]({{ .Permalink }})
{{ end }}{{ end }}{{ end }}{{ end }}
{{ end }}

## Installation

[Get Docker]({{ "/get-docker/" | absURL }})
{{ end -}}
