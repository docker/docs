---
description: components and formatting examples used in Docker's docs
title: Links
toc_max: 3
---

## Examples

- [External links](https://docker.com) open in a new tab
- [Internal links](links.md) open in the same tab

You can use relative links, using source filenames,
or you can use absolute links for pages as they appear on the final site.

#### Links to auto-generated content

When you link to heading IDs in auto-generated pages, such as CLI reference content,
you won't get any help from your editor in resolving the anchor names. That's
because the pages are generated at build-time and your editor or LSP doesn't know
about them in advance.

## Syntax

```md
[External links](https://docker.com)
[Internal links](links.md)
```
