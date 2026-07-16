# {{ .Title }}

| Name | Description |
|------|-------------|
{{- range hugo.Data.samples.samples }}
{{- if in .services $.Params.service }}
| [{{ .title }}]({{ .url }}) | {{ chomp .description }} |
{{- end }}
{{- end }}
