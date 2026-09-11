# API reference presentation

## Direction

Use the Docker Docs visual language for a reference that readers scan repeatedly:
paths, parameter names, response codes, and examples should be legible at a glance.
The endpoint signature is the strongest accent. Other borders group related
information; avoid decorative panels, shadows, and animation.

## Tokens

- Paper: `#ffffff`; subtle surface: `#f9f9fa`.
- Text: `#2c333f`; secondary text: `#566581`.
- Divider: `#e7eaef`; action blue: `#0d4df2`.
- Dark mode uses the site's dark surface and gray/blue tokens.
- `Roboto Flex` for prose and navigation; the site's `Roboto Mono` for code.
- Type sizes: 32px page title, 20px sections, 15px prose and field headings,
  14px navigation and controls, 13px code and secondary information.
- Prose uses a 1.7 line height and a maximum measure of 72 characters.

## Layout

Keep the three-column reference, left aligned, with the request example aligned
to the opening description. Use a single-column reading flow below the desktop
breakpoint, retaining the example near the beginning of the operation.

```text
Product navigation | Breadcrumbs              Version / reference links
                   | Operation title
                   | Method and endpoint
                   | Description + reference | Request example
                   | Parameters              | Copy action
                   | Responses               | Source details
```

An alternative with a full-width description placed examples too far below the
title on operations with long introductions. A panel around every section would
also compete with the catalog's cards. Keep reference sections open and reserve
the panel treatment for executable examples.

## Review against the brief

Retain the catalog's structure and the site's fonts, header, and page shell.
Replace oversized heading weights, undersized code, uppercase decorative labels,
and uneven navigation spacing with a restrained type scale. Align HTTP methods
in a fixed navigation column so long operation names wrap coherently. Keep
source-review details available without giving them the prominence of a warning
on every page. Native controls, visible focus, and static reference content
remain part of the design.

## Repository guidance proposal

Suggested addition to `AGENTS.md`: verify font-family names against the
`@font-face` declarations in `assets/css/style.css`. The registered code font
family is `Roboto Mono`; the shared `roboto flex mono` alias falls back to a
system font. The prototype uses the registered family explicitly.
