# {{ .Title }}
{{ .RenderShortcodes }}
{{ range where .Pages "Permalink" "ne" "" }}
- [{{ .Title }}]({{ .Permalink }})
{{ end }}
