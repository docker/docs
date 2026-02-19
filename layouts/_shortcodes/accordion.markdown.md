{{ $title := .Get "title" -}}
{{ $body := .InnerDeindent -}}
**{{ $title }}**

{{ $body }}
