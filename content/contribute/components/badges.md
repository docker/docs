---
description: components and formatting examples used in Docker's docs
title: Badges
toc_max: 3
---

We use badges to highlight new content in our docs. 

### Examples

{{< badge color=blue text="blue badge" >}}
{{< badge color=amber text="amber badge" >}}
{{< badge color=red text="red badge" >}}
{{< badge color=green text="green badge" >}}
{{< badge color=violet text="violet badge" >}}

You can also make a badge a link.

[{{< badge color="blue" text="badge with a link" >}}](../_index.md)

### Usage guidelines

We use the violet badge to highlight new content in our docs: {{< badge color=violet text="New" >}} 

Best practice is to use this badge for no longer than 2 months post release of the feature.

### Markup

```go
{{</* badge color=amber text="amber badge" */>}}
[{{</* badge color="blue" text="badge with a link" */>}}](../overview.md)
```
