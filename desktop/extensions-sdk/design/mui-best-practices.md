---
title: MUI best practices
description: Guidelines for using MUI to maximise compatibility with Docker Desktop
keywords: Docker, extensions, mui, theme, theming, material-ui, material
---

This article assumes you're following our recommended practice by using our [Material UI theme](https://www.npmjs.com/package/@docker/docker-mui-theme). Following the steps below maximises compatibility with Docker Desktop and minimises the work you need to do as an extension author. They should be considered supplementary to the non-MUI-specific guidelines found in the [UI Styling overview](overview.md).

## Assume the theme can change at any time

Resist the temptation to fine-tune your UI with precise colors, offsets and font sizings to make it look as attractive as possible. Any specializations you make today will be relative to the current MUI theme, and may look worse when the theme changes. Any part of the theme might change without warning, including (but not limited to):

-  The font, or font sizes
-  Border thicknesses or styles
-  Colors:
   -  Our palette members (e.g. `red-100`) could change their RGB values
   -  The semantic colors (e.g. `error`, `primary`, `textPrimary`, etc) could be changed to use a different member of our palette
   -  Background colors (e.g. those of the page, or of dialogs) could change
-  Spacings:
   -  The size of the basic unit of spacing,(exposed via `theme.spacing`. For instance, we may allow users to customize the density of the UI
   -  The default spacing between paragraphs or grid items

The best way to build your UI, so that it’s robust against future theming changes, is to:

-  Override the default styling as little as possible.
-  Use semantic typography. e.g. use `Typography`s or `Link`s with appropriate `variant`s instead of using typographical HTML elements (`<a>`, `<p>`, `<h1>`, etc) directly.
-  Use canned sizes. e.g. use `size="small"` on buttons, or `fontSize="small"` on icons, instead of specifying sizes in pixels.
-  Prefer semantic colors. e.g. use `error` or `primary` over explicit color codes.
-  Write as little CSS as possible. Write semantic markup instead. For example, if you want to space out paragraphs of text, use the `paragraph` prop on your `Typography` instances. If you want to space out something else, use a `Stack` or `Grid` with the default spacing.
-  Use visual idioms you’ve seen in the Docker Desktop UI, since these are the main ones we’ll test any theme changes against.

## When you go custom, centralize it

Sometimes you’ll need a piece of UI that doesn’t exist in our design system. If so, we recommend that you first reach out to us. We may already have something in our internal design system, or we may be able to expand our design system to accommodate your use case.

If you still decide to build it yourself after contacting us, try and define the new UI in a reuseable fashion. If you define your custom UI in just one place, it’ll make it easier to change in the future if our core theme changes. You could use:

-  A new `variant` of an existing component - see [MUI docs](https://mui.com/material-ui/customization/theme-components/#creating-new-component-variants)
-  A MUI mixin (a freeform bundle of reuseable styling rules defined inside a theme)
-  A new [reuseable component](https://mui.com/material-ui/customization/how-to-customize/#2-reusable-component)

Some of the above options require you to extend our MUI theme. See the MUI documentation on [theme composition](https://mui.com/material-ui/customization/theming/#nesting-the-theme).

## What's next?

- Take a look at our [UI styling guide](overview.md).
- Learn how to [publish your extension](../extensions/index.md).
