---
description: components and formatting examples used in Docker's docs
title: Links
toc_max: 3
---
## Examples

It is best practice if [a link opens in a new window](https://docker.com/){: target="_blank" rel="noopener" class="_" }

#### Links to auto-generated content

An example of a link to an auto-generated reference page that we pull in during docs builds:
[/engine/reference/builder/#env](/engine/reference/builder/#env).

  - If you can't find a reference page in the `docker.github.io`
  GitHub repository, but see it out on `docs.docker.com`, you can
  surmise that it's probably auto-generated from the codebase.
  (FYI, to view the Markdown source for the file, just click
  **Edit this page** on `docs.docker.com`. But don't use that URL in your docs.)

  - Go to the file in a web browser, grab everything after the domain name
  from the URL, and use that as the link in your docs file.

  - Keep in mind that this link doesn't resolve until you merge the PR and
  your docs are published on [docs.docker.com](/).

- It is best practice to avoid the use of absolute links when linking to other docs pages. Otherwise broken links may not be picked up. 

## HTML

```html

It is best practice if [a link opens in a new window](https://docker.com/){: target="_blank" rel="noopener" class="_" }

You can also have [a markdown link to a custom target ID](#formatting-examples)

An example of a link to an auto-generated reference page that we pull in during docs builds:
[/engine/reference/builder/#env](/engine/reference/builder/#env).

```