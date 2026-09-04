{{- /*
  Markdown output of the files shortcode (used for the *.md
  alternative output format that LLMs consume).

  Renders the project name and each file as a fenced code block with
  its path as a label. The file.html child populates the same .Store
  this template reads from.
*/ -}}
{{- with .Inner }}{{/* trigger child shortcodes */}}{{ end -}}
{{- $name := trim (.Get "name") " " -}}
{{- $files := .Store.Get "files" -}}

**`{{ $name }}/`**
{{ range $files }}
`{{ .path }}`{{ with .status }} ({{ . }}){{ end }}:

```{{ .lang }}
{{ .content }}
```
{{ end }}
