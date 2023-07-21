---
description: components and formatting examples used in Docker's docs
title: Links
toc_max: 3
---

## Internal links

Internal links should open [in the same window](/)

Use the path to the final permalink for the file, without the `.md` extension and a `/` at the end. For example, [to link to this file](/contribute/components/links), use `[](/contribute/components/links)`.

You can also have [a markdown link to a sub-heading or custom target ID](#external-links)

### HTML

```html
Internal links should open [in the same window](/)

Use the path to the final permalink for the file, without the `.md` extension and a `/` at the end. For example, [to link to this file](/contribute/components/links), use `[](/contribute/components/links)`.

You can also have [a markdown link to a sub-heading or custom target ID](#external-links)
```

## External links

External links should open [in a new window](https://docker.com/){: target="_blank" rel="noopener" class="_" }, have a `/` at the end of the URL.

### HTML

```html
External links should open [in a new window](https://docker.com/){: target="_blank" rel="noopener" class="_" }
```

### Links to auto-generated content

An example of a link to an auto-generated reference page that we pull in during docs builds:
[/engine/reference/builder/#env](/engine/reference/builder/#env).

  - If you can't find a reference page in the GitHub repository, but see it
  out on `docs.docker.com`, you can surmise that it's probably auto-generated 
  from the codebase. (FYI, to view the Markdown source for the file, just select
  **Edit this page** on `docs.docker.com`. But don't use that URL in your docs.)

  - Go to the file in a web browser, grab everything after the domain name
  from the URL, and use that as the link in your docs file.

  - Keep in mind that this link doesn't resolve until you merge the PR and
  your docs are published on [docs.docker.com](/).

- It is best practice to avoid the use of absolute links when linking to other docs pages. Otherwise broken links may not be picked up. 

## HTML

```html
An example of a link to an auto-generated reference page that we pull in during docs builds:
[/engine/reference/builder/#env](/engine/reference/builder/#env).
```