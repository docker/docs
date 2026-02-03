# {{ .Title }}

| Name | Description |
|------|-------------|
{{- range site.Data.samples.samples }}
{{- if in .services $.Params.service }}
| [{{ .title }}]({{ .url }}) | {{ chomp .description }} |
{{- end }}
{{- end }}
