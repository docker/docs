# {{ .Title }}

{{ .Description }}

{{ range hugo.Data.navigation.references }}{{ $label := .label }}{{ with site.GetPage .root }}{{ if ne .Params.sitemap false }}- [{{ $label }}]({{ partial "navigation/url.html" . | absURL }}): {{ .Description }}
{{ end }}{{ end }}{{ end }}
