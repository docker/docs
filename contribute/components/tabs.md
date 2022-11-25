---
description: components and formatting examples used in Docker's docs
title: Tabs
toc_max: 3
---

## Examples

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab3">TAB 1 HEADER</a></li>
  <li><a data-toggle="tab" data-target="#tab4">TAB 2 HEADER</a></li>
</ul>
<div class="tab-content">
<div id="tab3" class="tab-pane fade in active" markdown="1">

##### A Markdown header

- list item 1
- list item 2
<hr>
</div>
<div id="tab4" class="tab-pane fade" markdown="1">

##### Another Markdown header

- list item 3
- list item 4
<hr>
</div>
</div>

#### Synchronizing multiple tab groups

Consider an example where you have something like one tab per language, and
you have multiple tab sets like this on a page. You might want them to all
toggle together. In the following example, selecting `Go` or `Python` in one tab set toggles the
other tab set to match. We have Javascript that loads on every page that lets you
do this by setting the `data-group` attributes to be the same. The
`data-target` attributes still need to point to unique div IDs.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#python" data-group="python">Python</a></li>
</ul>
<div class="tab-content">
  <div id="go" class="tab-pane fade in active">Go content here<hr></div>
  <div id="python" class="tab-pane fade in">Python content here<hr></div>
</div>

And some content between the two sets, just for fun...

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#go-2" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#python-2" data-group="python">Python</a></li>
</ul>
<div class="tab-content">
  <div id="go-2" class="tab-pane fade in active">Go content here<hr></div>
  <div id="python-2" class="tab-pane fade in">Python content here<hr></div>
</div>

## HTML

```html
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab3">TAB 1 HEADER</a></li>
  <li><a data-toggle="tab" data-target="#tab4">TAB 2 HEADER</a></li>
</ul>
<div class="tab-content">
<div id="tab3" class="tab-pane fade in active" markdown="1">

##### A Markdown header

- list item 1
- list item 2
<hr>
</div>
<div id="tab4" class="tab-pane fade" markdown="1">

##### Another Markdown header

- list item 3
- list item 4
<hr>
</div>
</div>
```

```html
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#python" data-group="python">Python</a></li>
</ul>
<div class="tab-content">
  <div id="go" class="tab-pane fade in active">Go content here<hr></div>
  <div id="python" class="tab-pane fade in">Python content here<hr></div>
</div>

And some content between the two sets, just for fun...

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#go-2" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#python-2" data-group="python">Python</a></li>
</ul>
<div class="tab-content">
  <div id="go-2" class="tab-pane fade in active">Go content here<hr></div>
  <div id="python-2" class="tab-pane fade in">Python content here<hr></div>
```