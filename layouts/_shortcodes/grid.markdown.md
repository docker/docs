{{ range index .Page.Params (.Get "items" | default "grid") }}
- [{{ .title }}]({{ .link }}): {{ .description }}
{{ end }}
