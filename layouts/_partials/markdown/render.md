{{/* Bodies use named text partials because Hugo 0.169.0 decorators HTML-escape their inner blocks. */}}
{{- $page := .page -}}
{{- $content := .content -}}
{{- $probe := $content -}}
{{- if $page.Params.datafolder -}}
  {{/* CLI option tables contain literal syntax that is not a Goldmark document. */}}
  {{- $data := index hugo.Data.cli $page.Params.datafolder $page.Params.datafile -}}
  {{- $probe = printf "%s\n\n%s\n\n%s" ($data.short | default "") ($data.long | default "") ($data.examples | default "") -}}
  {{- range where ($data.options | default slice) "hidden" false -}}
    {{- with .description -}}
      {{- $option := $page.RenderString . -}}
    {{- end -}}
  {{- end -}}
{{- end -}}
{{- $rendered := $page.RenderString $probe -}}
{{- $destinations := dict "links" ($page.Store.Get "markdown-links") "images" ($page.Store.Get "markdown-images") -}}
{{- with $page.OutputFormats.Get "markdown" -}}
  {{- $file := printf "__markdown-links%s.json" .RelPermalink -}}
  {{- $resource := resources.FromString $file ($destinations | jsonify) -}}
  {{- $published := $resource.RelPermalink -}}
{{- end -}}
{{- return $content -}}
