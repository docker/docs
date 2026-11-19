# {{ .Title }}

{{ .Description }}

{{ range hugo.Data.navigation.categories }}
- [{{ .label }}]({{ .landing | absURL }}): {{ .description }}
{{ end }}
