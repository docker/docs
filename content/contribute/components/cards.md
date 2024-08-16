---
description: components and formatting examples used in Docker's docs
title: Cards
toc_max: 3
grid:
- title: Docker Desktop
  description: Docker on your Desktop.
  icon: install_desktop
  link: /desktop/
- title: Docker Engine
  description: Vrrrrooooommm
  icon: developer_board
  link: /engine/
- title: Docker Build
  description: Clang bang
  icon: build
  link: /build/
- title: Docker Compose
  description: Figgy!
  icon: account_tree
  link: /compose/
- title: Docker Hub
  description: so much content, oh wow
  icon: hub
  link: /docker-hub/
---

Cards can be added to a page using the `card` shortcode.
The parameters for this shortcode are:

| Parameter   | Description                                                          |
| ----------- | -------------------------------------------------------------------- |
| title       | The title of the card                                                |
| icon        | The icon slug of the card                                            |
| image       | Use a custom image instead of an icon (mutually exclusive with icon) |
| link        | (Optional) The link target of the card, when clicked                 |
| description | A description text, in Markdown                                      |

> [!NOTE]
>
> There's a known limitation with the Markdown description of cards,
> in that they can't contain relative links, pointing to other .md documents.
> Such links won't render correctly. Instead, use an absolute link to the URL
> path of the page that you want to link to.
>
> For example, instead of linking to `../install/linux.md`, write:
> `/engine/install/linux/`.

## Example

{{< card
  title="Get your Docker on"
  icon=favorite
  link=https://docs.docker.com/
  description="Build, share, and run your apps with Docker" >}}

## Markup

```go
{{</* card
  title="Get your Docker on"
  icon=favorite
  link=https://docs.docker.com/
  description="Build, share, and run your apps with Docker"
*/>}}
```

### Grid

There's also a built-in `grid` shortcode that generates a... well, grid... of cards.
The grid size is 3x3 on large screens, 2x2 on medium, and single column on small.

{{< grid >}}

Grid is a separate shortcode from `card`, but it implements the same component under the hood.
The markup you use to insert a grid is slightly unconventional. The grid shortcode takes no arguments.
All it does is let you specify where on a page you want your grid to appear.

```go
{{</* grid */>}}
```

The data for the grid is defined in the front matter of the page, under the `grid` key, as follows:

```yaml
# front matter section of a page
title: some page
grid:
  - title: "Docker Engine"
    description: Vrrrrooooommm
    icon: "developer_board"
    link: "/engine/"
  - title: "Docker Build"
    description: Clang bang
    icon: "build"
    link: "/build/"
```
