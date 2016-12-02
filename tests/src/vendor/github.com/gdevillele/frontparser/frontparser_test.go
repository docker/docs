package frontparser

const validFrontmatterSample string = `
---
title: ceci est un titre
template: fancy.tmpl
---

# Markdown title

This is some content.

## A Subtitle

Some more content.
`

// missing heading delimiter
const invalidFrontmatterSample1 string = `
---
`

func TestPositiveFrontmatterDetection() {

}

func TestNegativeFrontmatterDetection() {

}
