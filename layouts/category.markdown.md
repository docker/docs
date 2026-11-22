# {{ .Title }}

{{ .Description }}

{{ range hugo.Data.navigation.categories }}{{ if eq .id $.Params.category }}{{ range .groups }}
## {{ .label }}

{{ range .entries }}{{ $path := .page | default "" }}{{ $label := .label | default "" }}
{{ with .collection }}{{ $product := index hugo.Data.navigation.products . }}{{ $path = $product.root }}{{ $label = $product.label }}{{ end }}
{{ with site.GetPage $path }}{{ if ne .Params.sitemap false }}- [{{ $label | default .LinkTitle }}]({{ .Permalink }}): {{ .Description }}
{{ end }}{{ end }}{{ end }}{{ end }}{{ end }}{{ end }}
