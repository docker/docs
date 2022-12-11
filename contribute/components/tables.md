---
description: components and formatting examples used in Docker's docs
title: Tables
toc_max: 3
---

## Example

### Basic table

| Permission level                                                         | Access                                                       |
|:-------------------------------------------------------------------------|:-------------------------------------------------------------|
| **Bold** or _italic_ within a table cell. Next cell is empty on purpose. |                                                              |
|                                                                          | Previous cell is empty. A `--flag` in mono text.             |
| Read                                                                     | Pull                                                         |
| Read/Write                                                               | Pull, push                                                   |
| Admin                                                                    | All of the above, plus update description, create, and delete |

### Feature-support table

{% assign yes = '![yes](/assets/images/green-check.svg){: .inline style="height: 14px; margin: 0 auto"}' %}

| Platform                | x86_64 / amd64         | 
|:------------------------|:-----------------------|
| Ubuntu     | {{ yes }} |
| Debian     | {{ yes }} |
| Fedora     | {{ yes }} |


## Markdown

### Basic table

```
| Permission level                                                         | Access                                                       |
|:-------------------------------------------------------------------------|:-------------------------------------------------------------|
| **Bold** or _italic_ within a table cell. Next cell is empty on purpose. |                                                              |
|                                                                          | Previous cell is empty. A `--flag` in mono text.             |
| Read                                                                     | Pull                                                         |
| Read/Write                                                               | Pull, push                                                   |
| Admin                                                                    | All of the above, plus update description, create, and delete |
```
The alignment of the cells in the source doesn't really matter. The ending pipe
character is optional (unless the last cell is supposed to be empty). The header
row and separator row are optional.

### Feature-support table

Before you add the table you need to add:

{% highlight liquid %}
{% raw %} {% assign yes = '![yes](/assets/images/green-check.svg){% endraw %}{% raw %}{: .inline style="height: 14px; margin: 0 auto"}' %}{% endraw %}
{% endhighlight %}

```
| Platform                | x86_64 / amd64         | 
|:------------------------|:-----------------------|
| Ubuntu     | {% raw %}{{ yes }}{% endraw %} |
| Debian     | {% raw %}{{ yes }}{% endraw %} |
| Fedora     | {% raw %}{{ yes }}{% endraw %} |
```


