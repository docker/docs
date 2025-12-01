---
title: Examples
description: Get inspiration from agent examples
keywords: [ai, agent, cagent]
weight: 40
---

Get inspiration from the following agent examples.
See more examples in the [cagent GitHub repository](https://github.com/docker/cagent/tree/main/examples).

## Development team

{{% cagent-example.inline "dev-team.yaml" %}}
{{- $example := .Get 0 }}
{{- $baseUrl := "https://raw.githubusercontent.com/docker/cagent/refs/heads/main/examples" }}
{{- $url := fmt.Printf "%s/%s" $baseUrl $example }}
{{- with resources.GetRemote $url }}
{{ $data := .Content | transform.Unmarshal }}

```yaml {collapse=true}
{{ .Content }}
```

{{ end }}
{{% /cagent-example.inline %}}

## Go developer

{{% cagent-example.inline "gopher.yaml" /%}}

## Technical blog writer

{{% cagent-example.inline "blog.yaml" /%}}
