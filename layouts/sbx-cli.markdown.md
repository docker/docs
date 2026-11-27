{{- partial "markdown/render.md" (dict "page" . "content" (partial "markdown/body/sbx-cli.md" .)) -}}
{{ define "_partials/markdown/body/sbx-cli.md" -}}
{{- $data := index hugo.Data.sbx_cli .Params.datafile -}}
# {{ .Title }}

{{ with $data.synopsis }}**Description:** {{ strings.TrimSpace . }}{{ end }}

{{ with $data.usage }}**Usage:** `{{ . }}`{{ end }}

{{ if $data.experimental }}> [!NOTE]
> **Experimental**
>
> This command is experimental. Experimental features are intended for testing and feedback as their functionality or design may change between releases without warning or can be removed entirely in a future release.
{{ end }}

{{ with $data.description }}## Description

{{ . }}
{{ end }}

{{ if eq .Kind "section" }}## Commands

| Command | Description |
|---------|-------------|
{{ range .Pages }}{{ if .Params.datafile }}{{ $child := index hugo.Data.sbx_cli .Params.datafile }}| [`{{ .Title }}`]({{ .Permalink }}) | {{ if $child.experimental }}**experimental** {{ end }}{{ strings.TrimSpace $child.synopsis }} |
{{ end }}{{ end }}
{{ end }}

{{ with $data.options }}{{ $opts := where . "name" "ne" "help" }}{{ with $opts }}## Options

| Option | Default | Description |
|--------|---------|-------------|
{{ range . }}| {{ with .shorthand }}`-{{ . }}`, {{ end }}`--{{ .name }}` | {{ with .default_value }}{{ $skipDefault := `[],false,` }}{{ cond (in $skipDefault .) "" (printf "`%s`" .) }}{{ end }} | {{ if .experimental }}**experimental** {{ end }}{{ with .usage }}{{ strings.Replace (strings.Replace (strings.TrimSpace .) "\n" "<br>") "|" "\\|" }}{{ end }} |
{{ end }}
{{ end }}{{ end }}

{{ with $data.inherited_options }}{{ $opts := where . "name" "ne" "help" }}{{ with $opts }}## Global options

| Option | Default | Description |
|--------|---------|-------------|
{{ range . }}| {{ with .shorthand }}`-{{ . }}`, {{ end }}`--{{ .name }}` | {{ with .default_value }}{{ $skipDefault := `[],false,` }}{{ cond (in $skipDefault .) "" (printf "`%s`" .) }}{{ end }} | {{ if .experimental }}**experimental** {{ end }}{{ with .usage }}{{ strings.Replace (strings.Replace (strings.TrimSpace .) "\n" "<br>") "|" "\\|" }}{{ end }} |
{{ end }}
{{ end }}{{ end }}

{{ with $data.example }}## Examples

{{ $dedented := replaceRE `(?m)^ {2,6}` "" . }}{{ printf "```console\n%s\n```" (strings.TrimSpace $dedented) }}
{{ end }}
{{ end -}}
