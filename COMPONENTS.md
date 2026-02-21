# Docker Documentation Components Guide

This guide explains how to use components, shortcodes, and special features
when writing Docker documentation. For writing style and grammar, see
[STYLE.md](STYLE.md).

## Front matter

Every documentation page requires front matter at the top of the file.

### Required fields

```yaml
---
title: Install Docker Engine on Ubuntu
description: Instructions for installing Docker Engine on Ubuntu
keywords: requirements, apt, installation, ubuntu
---
```

| Field       | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| title       | Page title, used in `<h1>` and `<title>` tag                 |
| description | SEO description (150-160 characters), added to HTML metadata |
| keywords    | Comma-separated keywords for SEO                             |

### Optional fields

| Field           | Description                                                        |
| --------------- | ------------------------------------------------------------------ |
| linkTitle       | Shorter title for navigation and sidebar (if different from title) |
| weight          | Controls sort order in navigation (lower numbers appear first)     |
| aliases         | List of URLs that redirect to this page                            |
| toc_min         | Minimum heading level in table of contents (default: 2)            |
| toc_max         | Maximum heading level in table of contents (default: 3)            |
| notoc           | Set to `true` to disable table of contents                         |
| sitemap         | Set to `false` to exclude from search engine indexing              |
| sidebar.badge   | Add badge to sidebar (see [Badges](#badges))                       |
| sidebar.reverse | Reverse sort order of pages in section                             |
| sidebar.goto    | Override sidebar link URL                                          |

### Example with optional fields

```yaml
---
title: Install Docker Engine on Ubuntu
linkTitle: Install on Ubuntu
description: Instructions for installing Docker Engine on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall
weight: 10
aliases:
  - /engine/installation/linux/ubuntu/
  - /install/linux/ubuntu/
toc_max: 4
params:
  sidebar:
    badge:
      color: blue
      text: Beta
---
```

### Series (guide) pages

Section pages under `content/guides/` automatically use the `series` layout
(via a Hugo cascade in `hugo.yaml`). Series pages support additional front
matter parameters for the metadata card:

```yaml
---
title: Getting started
description: Learn the basics of Docker
summary: |
  A longer summary shown on the series landing page.
params:
  proficiencyLevel: Beginner
  time: 15 minutes
  prerequisites: None
---
```

| Field            | Description                              |
| ---------------- | ---------------------------------------- |
| summary          | Extended description for the series page |
| proficiencyLevel | Skill level (Beginner, Intermediate)     |
| time             | Estimated time to complete               |
| prerequisites    | Prerequisites or "None"                  |

## Shortcodes

Shortcodes are reusable components that add rich functionality to your
documentation.

### Tabs

Use tabs for platform-specific instructions or content variations.

**Example:**

{{< tabs >}}
{{< tab name="Linux" >}}

```console
$ docker run hello-world
```

{{< /tab >}}
{{< tab name="macOS" >}}

```console
$ docker run hello-world
```

{{< /tab >}}
{{< tab name="Windows" >}}

```powershell
docker run hello-world
```

{{< /tab >}}
{{< /tabs >}}

**Syntax:**

```markdown
{{</* tabs */>}}
{{</* tab name="Linux" */>}}
Content for Linux
{{</* /tab */>}}
{{</* tab name="macOS" */>}}
Content for macOS
{{</* /tab */>}}
{{</* /tabs */>}}
```

**Tab groups (synchronized selection):**

Use the `group` attribute to synchronize tab selection across multiple tab
sections:

```markdown
{{</* tabs group="os" */>}}
{{</* tab name="Linux" */>}}
Linux content for first section
{{</* /tab */>}}
{{</* tab name="macOS" */>}}
macOS content for first section
{{</* /tab */>}}
{{</* /tabs */>}}

{{</* tabs group="os" */>}}
{{</* tab name="Linux" */>}}
Linux content for second section
{{</* /tab */>}}
{{</* tab name="macOS" */>}}
macOS content for second section
{{</* /tab */>}}
{{</* /tabs */>}}
```

When a user selects "Linux" in the first section, all tabs in the "os" group
switch to "Linux".

### Accordion

Use accordions for collapsible content like optional steps or advanced
configuration.

**Example:**

{{< accordion title="Advanced configuration" >}}

```yaml
version: "3.8"
services:
  web:
    build: .
    ports:
      - "8000:8000"
```

{{< /accordion >}}

**Syntax:**

```markdown
{{</* accordion title="Advanced configuration" */>}}
Content inside the accordion
{{</* /accordion */>}}
```

### Include

Reuse content across multiple pages using the include shortcode. Include
files must be in the `content/includes/` directory.

**Syntax:**

```markdown
{{</* include "filename.md" */>}}
```

**Example:**

```markdown
{{</* include "install-prerequisites.md" */>}}
```

### Badges

Use badges to highlight new, beta, experimental, or deprecated content.

**Example:**

{{< badge color=blue text="Beta" >}}
{{< badge color=violet text="Experimental" >}}
{{< badge color=green text="New" >}}
{{< badge color=gray text="Deprecated" >}}

**Syntax:**

Inline badge:

```markdown
{{</* badge color=blue text="Beta" */>}}
```

Badge as link:

```markdown
[{{</* badge color=blue text="Beta feature" */>}}](link-to-page.md)
```

Sidebar badge (in front matter):

```yaml
---
title: Page title
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
---
```

**Color options:**

- `blue` - Beta content
- `violet` - Experimental or early access
- `green` - New GA content
- `amber` - Warning or attention
- `red` - Critical
- `gray` - Deprecated

**Usage guidelines:**

- Use badges for no longer than 2 months post-release
- Use violet for experimental/early access features
- Use blue for beta features
- Use green for new GA features
- Use gray for deprecated features

### Summary bars

Summary bars indicate subscription requirements, version requirements, or
administrator-only features at the top of a page.

**Example:**

{{< summary-bar feature_name="Docker Scout" >}}

**Setup:**

1. Add feature to `/data/summary.yaml`:

```yaml
features:
  Docker Scout:
    subscription: Business
    availability: GA
    requires: "Docker Desktop 4.17 or later"
```

2. Add shortcode to page:

```markdown
{{</* summary-bar feature_name="Docker Scout" */>}}
```

**Attributes in summary.yaml:**

| Attribute    | Description                           | Values                                            |
| ------------ | ------------------------------------- | ------------------------------------------------- |
| subscription | Subscription tier required            | All, Personal, Pro, Team, Business                |
| availability | Product development stage             | Experimental, Beta, Early Access, GA, Retired     |
| requires     | Minimum version requirement           | String describing version (link to release notes) |
| for          | Indicates administrator-only features | Administrators                                    |

### Buttons

Create styled buttons for calls to action.

**Syntax:**

```markdown
{{</* button text="Download Docker Desktop" url="/get-docker/" */>}}
```

### Cards

Create card layouts for organizing content.

**Syntax:**

```markdown
{{</* card
  title="Get started"
  description="Learn Docker basics"
  link="/get-started/"
*/>}}
```

### Icons

Use Material Symbols icons in your content.

**Syntax:**

```markdown
{{</* icon name="check_circle" */>}}
```

Browse available icons at
[Material Symbols](https://fonts.google.com/icons).

## Callouts

Use GitHub-style callouts to emphasize important information. Use sparingly -
only when information truly deserves special emphasis.

**Syntax:**

```markdown
> [!NOTE]
> This is a note with helpful context.

> [!TIP]
> This is a helpful suggestion or best practice.

> [!IMPORTANT]
> This is critical information users must understand.

> [!WARNING]
> This warns about potential issues or consequences.

> [!CAUTION]
> This is for dangerous operations requiring extreme care.
```

**When to use each type:**

| Type      | Use for                                                        | Color  |
| --------- | -------------------------------------------------------------- | ------ |
| Note      | Information that may not apply to all readers or is tangential | Blue   |
| Tip       | Helpful suggestions or best practices                          | Green  |
| Important | Issues of moderate magnitude                                   | Yellow |
| Warning   | Actions that may cause damage or data loss                     | Red    |
| Caution   | Serious risks                                                  | Red    |

## Code blocks

Format code with syntax highlighting using triple backticks and language
hints.

### Language hints

Common language hints:

- `console` - Interactive shell with `$` prompt
- `bash` - Bash scripts
- `dockerfile` - Dockerfiles
- `yaml` - YAML files
- `json` - JSON data
- `go`, `python`, `javascript` - Programming languages
- `powershell` - PowerShell commands

**Interactive shell example:**

````markdown
```console
$ docker run hello-world
```
````

**Bash script example:**

````markdown
```bash
#!/usr/bin/bash
echo "Hello from Docker"
```
````

### Variables in code

Use `<VARIABLE_NAME>` format for placeholder values:

````markdown
```console
$ docker tag <IMAGE_NAME> <REGISTRY>/<IMAGE_NAME>:<TAG>
```
````

This syntax renders variables in a special color and font style.

### Highlighting lines

Highlight specific lines with `hl_lines`:

````markdown
```go {hl_lines=[3,4]}
func main() {
    fmt.Println("Hello")
    fmt.Println("This line is highlighted")
    fmt.Println("This line is highlighted")
}
```
````

### Collapsible code blocks

Make long code blocks collapsible:

````markdown
```dockerfile {collapse=true}
# syntax=docker/dockerfile:1
FROM golang:1.21-alpine
RUN apk add --no-cache git
# ... more lines
```
````

## Images

Add images using standard Markdown syntax with optional query parameters for
sizing and borders.

**Basic syntax:**

```markdown
![Alt text description](/assets/images/image.png)
```

**With size parameters:**

```markdown
![Alt text](/assets/images/image.png?w=600&h=400)
```

**With border:**

```markdown
![Alt text](/assets/images/image.png?border=true)
```

**Combined parameters:**

```markdown
![Alt text](/assets/images/image.png?w=600&h=400&border=true)
```

**Best practices:**

- Store images in `/static/assets/images/`
- Use descriptive alt text for accessibility
- Compress images before adding to repository
- Capture only relevant UI, avoid unnecessary whitespace
- Use realistic text, not lorem ipsum
- Remove unused images from repository

## Videos

Embed videos using HTML video tags or platform-specific embeds.

**Local video:**

```html
<video controls width="100%">
  <source src="/assets/videos/demo.mp4" type="video/mp4" />
</video>
```

**YouTube embed:**

```html
<iframe
  width="560"
  height="315"
  src="https://www.youtube.com/embed/VIDEO_ID"
  frameborder="0"
  allow="autoplay; encrypted-media"
  allowfullscreen
>
</iframe>
```

## Links

Use standard Markdown link syntax. For internal links, use relative paths
with source filenames.

**External link:**

```markdown
[Docker Hub](https://hub.docker.com/)
```

**Internal link (same section):**

```markdown
[Installation guide](install.md)
```

**Internal link (different section):**

```markdown
[Get started](/get-started/overview.md)
```

**Link to heading:**

```markdown
[Prerequisites](#prerequisites)
```

**Best practices:**

- Use descriptive link text (around 5 words)
- Front-load important terms
- Avoid generic text like "click here" or "learn more"
- Don't include end punctuation in link text
- Use relative paths for internal links

## Lists

### Unordered lists

```markdown
- First item
- Second item
- Third item
  - Nested item
  - Another nested item
```

### Ordered lists

```markdown
1. First step
2. Second step
3. Third step
   1. Nested step
   2. Another nested step
```

### Best practices

- Limit bulleted lists to 5 items when possible
- Don't add commas or semicolons to list item ends
- Capitalize list items for ease of scanning
- Make all list items parallel in structure
- For nested sequential lists, use lowercase letters (a, b, c)

## Tables

Use standard Markdown table syntax with sentence case for headings.

**Example:**

```markdown
| Feature        | Description                     | Availability |
| -------------- | ------------------------------- | ------------ |
| Docker Compose | Multi-container orchestration   | All          |
| Docker Scout   | Security vulnerability scanning | Business     |
| Build Cloud    | Remote build service            | Pro, Team    |
```

**Best practices:**

- Use sentence case for table headings
- Don't leave cells empty - use "N/A" or "None"
- Align decimals on the decimal point
- Keep tables scannable - avoid dense content

## Quick reference

### File structure

```plaintext
content/
├── manuals/              # Product documentation
│   ├── docker-desktop/
│   ├── docker-hub/
│   └── ...
├── guides/               # Learning guides
├── reference/            # API and CLI reference
└── includes/             # Reusable content snippets
```

### Common patterns

**Platform-specific instructions:**

Use tabs with consistent names: Linux, macOS, Windows

**Optional content:**

Use accordions for advanced or optional information

**Version/subscription indicators:**

Use badges or summary bars

**Important warnings:**

Use callouts (NOTE, WARNING, CAUTION)

**Code examples:**

Use `console` for interactive shells, appropriate language hints for code
