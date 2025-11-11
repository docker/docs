{{ .Title }}
{{ .RenderShortcodes }}
{{ range .Pages }}
- [{{ .Title }}](https://docs.docker.com{{ .RelPermalink }})
{{ end }}
