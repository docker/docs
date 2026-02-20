{{- $data := "" }}
{{- if .Params.datafolder }}
  {{- $data = index (index site.Data .Params.datafolder) .Params.datafile }}
{{- else }}
  {{- $data = index site.Data .Params.datafile }}
{{- end -}}
# {{ .Title }}

{{ with $data.short }}**Description:** {{ . }}{{ end }}

{{ with $data.usage }}**Usage:** `{{ . }}`{{ end }}

{{ with $data.aliases }}{{ $aliases := strings.Replace . (printf "%s, " $.Title) "" }}**Aliases:** {{ range $i, $alias := (strings.Split $aliases ", ") }}{{ if $i }}, {{ end }}`{{ $alias }}`{{ end }}{{ end }}

{{ .Content }}

{{ if $data.deprecated }}> [!WARNING]
> **Deprecated**
>
> This command is deprecated. It may be removed in a future Docker version.
{{ end }}

{{ if or $data.experimental $data.experimentalcli }}> [!NOTE]
> **Experimental**
>
> This command is experimental. Experimental features are intended for testing and feedback as their functionality or design may change between releases without warning or can be removed entirely in a future release.
{{ end }}

{{ with $data.kubernetes }}**Orchestrator:** Kubernetes{{ end }}
{{ with $data.swarm }}**Orchestrator:** Swarm{{ end }}

{{ with $data.long }}## Description

{{ . }}
{{ end }}

{{ with $data.options }}{{ $opts := where . "hidden" false }}{{ with $opts }}## Options

| Option | Default | Description |
|--------|---------|-------------|
{{ range . }}{{ $short := .shorthand }}{{ $long := .option }}| {{ with $short }}`-{{ . }}`, {{ end }}`--{{ $long }}` | {{ with .default_value }}{{ $skipDefault := `[],map[],false,0,0s,default,'',""` }}{{ cond (in $skipDefault .) "" (printf "`%s`" .) }}{{ end }} | {{ with .min_api_version }}API {{ . }}+{{ end }}{{ with .deprecated }} **Deprecated**{{ end }}{{ with .experimental }} **experimental (daemon)**{{ end }}{{ with .experimentalcli }} **experimental (CLI)**{{ end }}{{ with .kubernetes }} **Kubernetes**{{ end }}{{ with .swarm }} **Swarm**{{ end }}{{ if .description }} {{ strings.Replace .description "\n" "<br>" }}{{ end }} |
{{ end }}
{{ end }}{{ end }}

{{ with $data.examples }}## Examples

{{ . }}
{{ end }}

{{ if eq .Kind "section" }}## Subcommands

| Command | Description |
|---------|-------------|
{{ range .Pages }}{{ if and .Params.datafolder .Params.datafile }}{{ $subdata := index (index site.Data .Params.datafolder) .Params.datafile }}| [`{{ .Title }}`]({{ .Permalink }}) | {{ $subdata.short }} |
{{ end }}{{ end }}
{{ end }}
