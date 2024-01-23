---
description: components and formatting examples used in Docker's docs
title: Images
toc_max: 3
---

## Example

- A small image: ![a small image](/assets/images/footer_moby_icon.png)

- Large images occupy the full width of the reading column by default:

  ![a pretty wide image](/assets/images/banner_image_24512.png)

- Image size can be set using query parameters: `?h=<height>&w=<width>`

  ![a pretty wide image](/assets/images/banner_image_24512.png?w=100&h=50)

- Image with a border, also set with a query parameter: `?border=true`

  ![a small image](/assets/images/footer_moby_icon.png?border=true)

- Image that works in both light and dark mode:
  {{< dynamic-image lightSrc="/assets/images/dd-light.webp" darkSrc="/assets/images/dd-dark.webp" >}}

## HTML and Markdown

```markdown
- A small image: ![a small image](/assets/images/footer_moby_icon.png)

- Large images occupy the full width of the reading column by default:

  ![a pretty wide image](/assets/images/banner_image_24512.png)

- Image size can be set using query parameters: `?h=<height>&w=<width>`

  ![a pretty wide image](/assets/images/banner_image_24512.png?w=100&h=50)

- Image with a border, also set with a query parameter: `?border=true`

  ![a small image](/assets/images/footer_moby_icon.png?border=true)

- Image that works in both light and dark mode:
  {{</* dynamic-image lightSrc="/assets/images/dd-light.webp" darkSrc="/assets/images/dd-dark.webp" */>}}
```
