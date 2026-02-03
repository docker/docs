{{ $src := .Get "src" }}
{{ $alt := .Get "alt" }}
{{ $title := .Get "title" }}
![{{ $alt }}]({{ $src }}{{ with $title }} "{{ . }}"{{ end }})