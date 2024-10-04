---
title: Source file conventions
description: How new .md files should be formatted
keywords: source file, contribute, style guide
weight: 30
---

## File name

When you create a new .md file for new content, make sure: 
- File names are as short as possible
- Try to keep the file name to one word or two words
- Use a dash to separate words. For example:
  - `add-seats.md`  and `remove-seats.md`.
  - `multiplatform-images` preferred to `multi-platform-images`.

## Front matter

The front matter of a given page is in a section at the top of the Markdown
file that starts and ends with three hyphens. It includes YAML content. The
following keys are supported. The title, description, and keywords are required.

| Key             | Required | Description                                                                                                                                                                                                    |
|-----------------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| title           | yes      | The page title. This is added to the HTML output as a `<h1>` level header.                                                                                                                                     |
| description     | yes      | A sentence that describes the page contents. This is added to the HTML metadata. It’s not rendered on the page.                                                                                                |
| keywords        | yes      | A comma-separated list of keywords. These are added to the HTML metadata.                                                                                                                                      |
| aliases         | no       | A YAML list of pages which should redirect to the current page. At build time, each page listed here is created as an HTML stub containing a 302 redirect to this page.                                        |
| notoc           | no       | Either `true` or `false`. If `true`, no in-page TOC is generated for the HTML output of this page. Defaults to `false`. Appropriate for some landing pages that have no in-page headings.                      |
| toc_min         | no       | Ignored if `notoc` is set to `true`. The minimum heading level included in the in-page TOC. Defaults to `2`, to show `<h2>` headings as the minimum.                                                           |
| toc_max         | no       | Ignored if `notoc` is set to `false`. The maximum heading level included in the in-page TOC. Defaults to `3`, to show `<h3>` headings. Set to the same as `toc_min` to only show `toc_min` level of headings.  |
| sitemap         | no       | Exclude the page from indexing by search engines. When set to `false`, the page is excluded from `sitemap.xml`, and a `<meta name="robots" content="noindex"/>` header is added to the page.                   |
| sidebar.reverse | no       | This parameter for section pages changes the sort order of the pages in that section. Pages that would normally appear at the top, by weight or by title, will instead appear near the bottom, and vice versa. |
| sidebar.goto    | no       | Set this to change the URL that the sidebar should point to for this entry. See [pageless sidebar entries](#pageless-sidebar-entries).                                                                         |
| sidebar.badge   | no       | Set this to add a badge to the sidebar entry for this page. This param option consists of two fields: `badge.text` and `badge.color`.                                                                          |

Here's an example of a valid (but contrived) page metadata. The order of
the metadata elements in the front matter isn't important.

```liquid
---
description: Instructions for installing Docker Engine on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
title: Install Docker Engine on Ubuntu
aliases:
- /ee/docker-ee/ubuntu/
- /engine/installation/linux/docker-ce/ubuntu/
- /engine/installation/linux/docker-ee/ubuntu/
- /engine/installation/linux/ubuntu/
- /engine/installation/linux/ubuntulinux/
- /engine/installation/ubuntulinux/
- /install/linux/docker-ce/ubuntu/
- /install/linux/docker-ee/ubuntu/
- /install/linux/ubuntu/
- /installation/ubuntulinux/
toc_max: 4
---
```

## Body

The body of the page (with the exception of keywords) starts after the front matter.

### Text length

Splitting long lines (preferably up to 80 characters) can make it easier to provide feedback on small chunks of text.

## Pageless sidebar entries

If you want to add an entry to the sidebar, but you want the link to point somewhere else, you can use the `sidebar.goto` parameter.
This is useful in combination with `build.render` set to `always`, which creates a pageless entry in the sidebar that links to another page.

```md
---
title: Dummy sidebar link
build:
  render: never
sidebar:
  goto: /some/other/page/
weight: 30
---
```
