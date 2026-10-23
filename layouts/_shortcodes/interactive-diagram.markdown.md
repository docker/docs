{{- $src := .Get "src" -}}
{{- $filePath := path.Join .Page.File.Dir $src -}}
{{- if not (fileExists $filePath) -}}
  {{- errorf "interactive-diagram shortcode: file %q not found: %s" $filePath .Position -}}
{{- end -}}
{{- $diagram := readFile $filePath | transform.Unmarshal -}}
{{ $diagram.title }}

{{ $diagram.description }}

{{ range $index, $step := $diagram.steps -}}
{{ add $index 1 }}. {{ $step.label }}: {{ $step.body }}
{{ end -}}
