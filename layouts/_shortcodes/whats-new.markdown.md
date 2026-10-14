{{- $data := index hugo.Data "whats-new" -}}
{{- with $data.items }}

## What's new

{{ range . }}
{{ printf "- [%s: %s](%s): %s (%s)\n" .product .title .url .description (.published | time.Format "Jan 2, 2006") -}}
{{ end }}
{{- end }}
