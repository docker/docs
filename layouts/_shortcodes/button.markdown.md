{{ $text := .Get "text" -}}
{{ $url := .Get "url" -}}
[{{ $text }}]({{ $url }})
