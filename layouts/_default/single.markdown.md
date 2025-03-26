{{- $slug := .File.TranslationBaseName -}}
{{- $outputPath := printf "%s.md" $slug -}}
{{- .Page.Store.Set "RelativePermalink" (print "/" $outputPath) -}}

{{ .RawContent }}

{{ .RenderShortcodes }}