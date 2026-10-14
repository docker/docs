{{ if or (not (isset .Params "type")) (not (isset .Params "default")) }}
  {{ errorf "setting-metadata requires type and default: %s" .Position }}
{{ end }}
- Type: {{ .Get "type" }}
- Default: {{ .Get "default" }}
{{ with .Get "env" }}- Environment variable: `{{ . }}`
{{ end }}
