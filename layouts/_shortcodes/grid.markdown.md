{{ range index .Page.Params (.Get "items" | default "grid") }}
- {{ if .link }}[{{ .title }}]({{ .link }}){{ else }}{{ .title }}{{ end }}: {{ .description }}
{{ end }}
