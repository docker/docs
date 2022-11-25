---
description: components and formatting examples used in Docker's docs
title: Tooltips
toc_max: 3
---
Tooltips are not visible on mobile devices or touchscreens, so don't rely on them as the
only way to communicate important info.

## Examples

<button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="left" title="Tooltip on left">Tooltip on left</button>

  <button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="top" title="Tooltip on top">Tooltip on top</button>

  <button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="bottom" title="Tooltip on bottom">Tooltip on bottom</button>

  <button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="right" title="Tooltip on right">Tooltip on right</button>

This is a paragraph that has a tooltip. We position it to the left so it doesn't align with the middle top of the paragraph (that looks weird).
{:data-toggle="tooltip" data-placement="left" title="Markdown tooltip example"}

<a href="/contribute/components/tooltips" target="_blank" rel="noopener" class="_"><span class="badge badge-info" data-toggle="tooltip" data-placement="right" title="Open the test page (in a new window)">Test</span></a>

## HTML 

```html
<button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="left" title="Tooltip on left">Tooltip on left</button>

  <button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="top" title="Tooltip on top">Tooltip on top</button>

  <button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="bottom" title="Tooltip on bottom">Tooltip on bottom</button>

  <button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="right" title="Tooltip on right">Tooltip on right</button>

This is a paragraph that has a tooltip. We position it to the left so it doesn't align with the middle top of the paragraph (that looks weird).
{:data-toggle="tooltip" data-placement="left" title="Markdown tooltip example"}

<a href="/contribute/components" target="_blank" rel="noopener" class="_"><span class="badge badge-info" data-toggle="tooltip" data-placement="right" title="Open the test page (in a new window)">Test</span></a>
```
