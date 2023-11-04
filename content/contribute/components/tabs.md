---
description: components and formatting examples used in Docker's docs
title: Tabs
toc_max: 3
---

The tabs component consists of two shortcodes:

- `{{</* tabs */>}}`
- `{{</* tab name="name of the tab" */>}}`

The `{{</* tabs */>}}` shortcode is a parent, component, wrapping a number of `tabs`.
Each `{{</* tab */>}}` is given a name using the `name` attribute.

You can optionally specify a `group` attribute for the `tabs` wrapper to indicate
that a tab section should belong to a group of tabs. See [Groups](#groups).

## Example

{{< tabs >}}
{{< tab name="JavaScript">}}

```js
console.log("hello world")
```

{{< /tab >}}
{{< tab name="Go">}}

```go
fmt.Println("hello world")
```

{{< /tab >}}
{{< /tabs >}}

## Markup

````markdown
{{</* tabs */>}}
{{</* tab name="JavaScript" */>}}

```js
console.log("hello world")
```

{{</* /tab */>}}
{{</* tab name="Go" */>}}

```go
fmt.Println("hello world")
```

{{</* /tab */>}}
{{</* /tabs */>}}
````

## Groups

You can optionally specify a tab group on the `tabs` shortcode.
Doing so will synchronize the tab selection for all of the tabs that belong to the same group.

### Tab group example

The following example shows two tab sections belonging to the same group.

{{< tabs group="code" >}}
{{< tab name="JavaScript">}}

```js
console.log("hello world")
```

{{< /tab >}}
{{< tab name="Go">}}

```go
fmt.Println("hello world")
```

{{< /tab >}}
{{< /tabs >}}

{{< tabs group="code" >}}
{{< tab name="JavaScript">}}

```js
const res = await fetch("/users/1")
```

{{< /tab >}}
{{< tab name="Go">}}

```go
resp, err := http.Get("/users/1")
```

{{< /tab >}}
{{< /tabs >}}
