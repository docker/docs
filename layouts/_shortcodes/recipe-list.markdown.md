{{ range where .Page.RegularPages "Params.sidebar.group" (.Get "group") }}
- [{{ .LinkTitle }}]({{ .Permalink }}): {{ .Description }}
{{ end }}
