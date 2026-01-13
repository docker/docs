---
title: {{ .Title }}
url: {{ .Permalink }}
{{- range .Ancestors }}
  {{- if and (not .IsHome) .Permalink }}
parent:
  title: {{ .Title }}
  url: {{ .Permalink }}
  {{- break }}
  {{- end }}
{{- end }}
{{- if .Ancestors }}
breadcrumbs:
{{- range .Ancestors.Reverse }}
  {{- if and (not .IsHome) .Permalink }}
  - title: {{ .Title }}
    url: {{ .Permalink }}
  {{- end }}
{{- end }}
  - title: {{ .Title }}
    url: {{ .Permalink }}
{{- end }}
{{- $children := where .Pages "Permalink" "ne" "" }}
{{- if $children }}
children:
{{- range $children }}
  - title: {{ .Title }}
    url: {{ .Permalink }}
    {{- with .Description }}
    description: {{ . }}
    {{- end }}
{{- end }}
{{- end }}
---

{{ .RenderShortcodes }}
