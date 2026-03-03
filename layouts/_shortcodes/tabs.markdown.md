{{ with .Inner }}{{/* don't do anything, just call it */}}{{ end -}}
{{ range (.Store.Get "tabs") -}}
**{{ .name }}**

{{ .content }}
{{- end -}}
