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
{{- with .NextInSection }}
next:
  title: {{ .Title }}
  url: {{ .Permalink }}
{{- end }}
{{- with .PrevInSection }}
prev:
  title: {{ .Title }}
  url: {{ .Permalink }}
{{- end }}
---

{{ .RenderShortcodes }}
