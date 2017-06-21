---
description: Smoketest page
title: Testing page
hide_from_sitemap: true
---

# Heading 1

Most pages don't actually have a H1 heading. The page title from the metadata is
automatically inserted.

## Heading 2

This is the highest heading included in the right-nav. To include more heading
levels, set `toc_min: 1` in the page-s front-matter. You can go all the way to
6, but if `toc_min` is geater than `toc_max` then no headings will show.

### Heading 3

This is the lowest heading included in the right-nav, by default. To include
more heading levels, set `toc_max: 4` in the page's front-matter. You can go all
the way to 6.

#### Heading 4

This heading is not included in the right-nav. To include it set `toc_max: 4` in
the page's front-matter.

##### Heading 5

This heading is not included in the right-nav. To include it set `toc_max: 5` in
the page's front-matter.

###### Heading 6

This is probably too many headings. Try to avoid it.

This heading is not included in the right-nav. To include it set `toc_max: 6` in
the page's front-matter.

## Typography

Plain block of text.

Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.

**Inline text styles**:

- **bold**
- _italic_
- ***bold italic***
- ~~strikethrough~~
- <u>underline</u>
- _<u>underline italic</u>_
- **<u>underline bold</u>**
- ***<u>underline bold italic</u>***
- `monospace text`
- **`monospace bold`**

## Links and images

### Links

- [a markdown link](https://docker.com/)
https://github.com/docker/docker.github.io/tree/master/docker-cloud/images
- [a markdown link that opens in a new window](https://docker.com/){: target="_blank" class="_" }
  (the `class="_"` trick prevents Atom from italicizing the whole rest of the file until it encounters another underscore.)

- <a href="https://docker.com/">an HTML link</a>

- <a href="https://docker.com/" target="_blank" class="_">an HTML link that opens in a new window</a>

- A link to a Github PR in `docker/docker`: {% include github-pr.md pr=28199 %}

- A link to a Github PR in `docker/docker.github.io`: {% include github-pr.md repo=docker.github.io pr=9999 %}

(you can also specify `org=foo` to use a Github organization other than Docker).


### Images

- A small cute image: ![a small cute image](/images/footer_moby_icon.png)

- A small cute image that is a link. The extra newline here makes it not show
  inline:

  [![a small cute image](/images/footer_moby_icon.png)](https://www.docker.com/)

- A big wide image: ![a pretty wide image](/images/banner_image_24512.png)

- The same as above but using HTML: <img src="/images/banner_image_24512.png" alt="a pretty wide image using HTML"/>

[Some Bootstrap image classes](https://v4-alpha.getbootstrap.com/content/images/)
might be interesting. You can use them with Markdown or HTML images.

- An image using the Bootstrap "thumbnail" class: ![an image as a thumbnail](/images/footer_moby_icon.png){: class="img-thumbnail" }

- The same one, but using HTML: <img class="img-thumbnail" src="/images/footer_moby_icon.png" alt="a pretty wide image as a thumbnail, using HTML"/>

## Videos

You can add a link to a YouTube video like this:

[![Deploying Swarms on Microsoft Azure with Docker Cloud](/docker-cloud/cloud-swarm/images/video-azure-docker-cloud.png)](https://www.youtube.com/watch?v=LlpyiGAVBVg "Deploying Swarms on Microsoft Azure with Docker Cloud"){:target="_blank" class="_"}

To make the `.png` shown above, first take a screen snap of the YouTube video
you want to use, then use a graphics app to overlay a play button onto the
image.

For the overlay, you can use the play button at
[/docker-cloud/images/](https://github.com/docker/docker.github.io/tree/master/docker-cloud/images).

## Lists

- Bullet list item 1
- Bullet list item 2
* Bullet list item 3

1.  Numbered list item 1. Two spaces between the period and the first letter
    helps with alignment.

2.  Numbered list item 2. Let's put a note in it.

    >**Note**: We did it!

3.  Numbered list item 3 with a code block in it. You need the blank line before
    the code block happens.

    ```bash
    $ docker run hello-world
    ```

4.  Numbered list item 4 with a bullet list inside it and a numbered list
    inside that.

    - Sub-item 1
    - Sub-item 2
      1.  Sub-sub-item 1
      2.  Sub-sub-item-2 with a table inside it because we like to party!
          Indentation is super important.

          |Header 1 | Header 2 |
          |---------|----------|
          | Thing 1 | Thing 2  |
          | Thing 3 | Thing 4  |


## Tables

Some tables in markdown and html.

| Permission level                                                         | Access                                                       |
|:-------------------------------------------------------------------------|:-------------------------------------------------------------|
| **Bold** or _italic_ within a table cell. Next cell is empty on purpose. |                                                              |
|                                                                          | Previous cell is empty. A `--flag` in mono text.             |
| Read                                                                     | Pull                                                         |
| Read/Write                                                               | Pull, push                                                   |
| Admin                                                                    | All of the above, plus update description, create and delete |

The alignment of the cells in the source doesn't really matter. The ending pipe
character is optional (unless the last cell is supposed to be empty). The header
row and separator row are optional.

If you need block-level HTML within your table cells, such as multiple
paragraphs, lists, sub-tables, etc, then you need to make a HTML table.
This is also the case if you need to use rowspans or colspans. Try to avoid
setting styles directly on your tables! If you set the width on a `<td>`, you
only need to do it on the first one. If you have a `<th>`, set it there.

> **Note**: If you need to have **markdown** in a **HTML** table, add
> `markdown="1"` to the HTML for the `<td>` cells that contain the Markdown.

<table>
  <tr>
    <th width="50%">Left channel</th>
    <th>Right channel</th>
  </tr>
  <tr>
  <td>This is some test text. <br><br>This is more <b>text</b> on a new line. <br><br>Lorem ipsum dolor <tt>sit amet</tt>, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
    </td>
    <td>This is some more text about the right hand side. There is a <a href="https://github.com/moby/moby/tree/master/experimental" target="_blank" class="_">link here to the Docker Experimental Features README</a> on GitHub.<br><br>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</td>
  </tr>
  <tr>
  <td>
  <p><a class="button outline-btn" href="/">Go to the docs!</a></p>
  <p><a href="/"><font color="#BDBDBD" size="-1">It is dark here. You are likely to be eaten by a grue.</font></a></p>
  </td>
  <td>
  <p><a class="button outline-btn" href="/">Go to the docs!</a></p>
  <p><a href="/"><font color="#BDBDBD" size="-1">It is dark here. You are likely to be eaten by a grue.</font></a></p>
  </td>
  </tr>
</table>

## Glossary links and content

The glossary source lives in the documentation repository
[docker.github.io](https://github.com/docker/docker.github.io) in
`_data/glossary.yaml`. The glossary publishes to
[https://docs.docker.com/glossary/](https://docs.docker.com/glossary/).

To update glossary content, edit `_data/glossary.yaml`.

To link to a glossary term, link to `glossary.md?term=YourGlossaryTerm` (for
example, [swarm](glossary.md?term=swarm)).

## Mixing Markdown and HTML

You can use <b>span-level</b> HTML tags within Markdown.

You can use a `<br />` tag to impose an extra newline like here:<br />

You can use entities like `&nbsp;` to keep a&nbsp;bunch&nbsp;of&nbsp;words&nbsp;together&nbsp;.

<center>
You can use block-level HTML tags too. This paragraph is centered.
</center>

Keep reading for more examples, such as creating tabbed content within the
page or displaying content as "cards".

## Jekyll / Liquid tricks

This paragraph is centered and colored green by setting CSS directly on the element.
**Even though you can do this and it's sometimes the right way to go, remember that if
we re-skin the site, any inline styles will need to be dealt with manually!**
{: style="text-align:center; color: green" }

{% assign my-text="foo" %}

The Liquid assignment just before this line fills in the following token {{ my-text }}.
This will be effective for the rest of this file unless the token is reset.

{% capture my-other-text %}foo{% endcapture %}
Here is another way: {{ my-other-text }}

You can nest captures within each other to represent more complex logic with Liquid.

### Liquid variable scope

- Things set in the top level of `_config.yml` are available as site variables, like `{{ site.debug }}`.
- Things set in the page's metadata, either via the defaults in `_config.yml` or per page, are available as page variables, like `{{ page.title }}`.
- In-line variables set via `assign` or `capture` are available for the remainder of the page after they are set.
- If you include a file, you can pass key-value pairs at the same time. These are available as include variables, like `{{ include.toc_min }}`.

## Bootstrap and CSS tricks

### Tabs

Here are some tabs:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab1">TAB 1 HEADER</a></li>
  <li><a data-toggle="tab" data-target="#tab2">TAB 2 HEADER</a></li>
</ul>
<div class="tab-content">
  <div id="tab1" class="tab-pane fade in active">TAB 1 CONTENT</div>
  <div id="tab2" class="tab-pane fade">TAB 2 CONTENT</div>
</div>

You need to adjust the `id` and `data-target` values to match your use case.

If you have Markdown inside the content of the `<div>`, just add `markdown="1"`
as an attribute in the HTML for the `<div>` and Kramdown will render it.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab3">TAB 1 HEADER</a></li>
  <li><a data-toggle="tab" data-target="#tab4">TAB 2 HEADER</a></li>
</ul>
<div class="tab-content">
<div id="tab3" class="tab-pane fade in active" markdown="1">
#### A Markdown header

- list item 1
- list item 2
</div>
<div id="tab4" class="tab-pane fade" markdown="1">
#### Another Markdown header

- list item 3
- list item 4
</div>
</div>

#### Synchronizing multiple tab groups

Consider an example where you have something like one tab per language, and
you have multiple tab sets like this on a page. You might want them to all
toggle together. We have Javascript that loads on every page that lets you
do this by setting the `data-group` attributes to be the same. Note that the
`data-target` attributes still need to point to unique div IDs.

In this example, selecting `Go` or `Python` in one tab set will toggle the
other tab set to match.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#go" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#python" data-group="python">Python</a></li>
</ul>
<div class="tab-content">
  <div id="go" class="tab-pane fade in active">Go content here</div>
  <div id="python" class="tab-pane fade in">Python content here</div>
</div>

And some content between the two sets, just for fun...

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#go-2" data-group="go">Go</a></li>
  <li><a data-toggle="tab" data-target="#python-2" data-group="python">Python</a></li>
</ul>
<div class="tab-content">
  <div id="go-2" class="tab-pane fade in active">Go content here</div>
  <div id="python-2" class="tab-pane fade in">Python content here</div>
</div>


### Cards

In a Bootstrap row, your columns need to add up to 12. Here are three cards in
a row, each of which takes up 1/3 (4/12) of the row. You need a couple `<br />`s
to clear the row before.<br /><br />


<div class="row">
  <div class="panel col-xs-12 col-md-4">This will take up 1/3 of the row unless the screen is small,
then it will take up the whole row.</div>
  <div class="panel col-xs-12 col-md-4">This will take up 1/3 of the row unless the screen is small,
then it will take up the whole row.</div>
  <div class="panel col-xs-12 col-md-4">This will take up 1/3 of the row unless the screen is small,
then it will take up the whole row.</div>
</div>

### Columnar text

You can use the CSS `column-count` to flow your text into multiple columns.
You need a couple `<br />`s to clear the row before.<br /><br />

<div style="column-count: 3">
This example uses a HTML div. This example uses a HTML div. This example uses a HTML div.
This example uses a HTML div. This example uses a HTML div. This example uses a HTML div.
This example uses a HTML div. This example uses a HTML div. This example uses a HTML div.
This example uses a HTML div. This example uses a HTML div. This example uses a HTML div.
This example uses a HTML div. This example uses a HTML div. This example uses a HTML div.
This example uses a HTML div. This example uses a HTML div. This example uses a HTML div.
This example uses a HTML div. This example uses a HTML div. This example uses a HTML div.
This example uses a HTML div. This example uses a HTML div. This example uses a HTML div.
This example uses a HTML div. This example uses a HTML div. This example uses a HTML div.
</div>

This example does it with Markdown. You can't have any blank lines or it will
break the Markdown block up. This example does it with Markdown. You can't have any blank lines or it will
break the Markdown block up. This example does it with Markdown. You can't have any blank lines or it will
break the Markdown block up. This example does it with Markdown. You can't have any blank lines or it will
break the Markdown block up. This example does it with Markdown. You can't have any blank lines or it will
break the Markdown block up. This example does it with Markdown. You can't have any blank lines or it will
break the Markdown block up. This example does it with Markdown. You can't have any blank lines or it will
break the Markdown block up. This example does it with Markdown. You can't have any blank lines or it will
break the Markdown block up. This example does it with Markdown. You can't have any blank lines or it will
break the Markdown block up.
{: style="column-count: 3 "}

## Running in-page Javascript

If you need to run custom Javascript within a page, and it depends upon JQuery
or Bootstrap, make sure the `<script>` tags are at the very end of the page,
after all the content. Otherwise the script may try to run before JQuery and
Bootstrap JS are loaded.

> **Note**: In general, this is a bad idea.

## Admonitions (notes)

Current styles for admonitions in
[`_scss/_notes.scss`](https://github.com/docker/docker.github.io/blob/master/_scss/_notes.scss)
support two broad categories of admonitions: those with prefixed text (**Note:**,
**Important:**, **Warning**) and those with prefixed icons.

The new styles (with icons) are defined in a way that will not impact
admonitions formatted with the original styles (prefixed text), so notes in your
published documents won't be adversely affected.

Examples of both styles are shown below.

### Examples (original styles, prefix words)

Admonitions with prefixed text use the following class tags, as shown in the examples.

* **Note:** No class tag is needed for standard notes.
* **Important:** Use the `important` class.
* **Warning:** Use the `warning` class.


> **Note**: This is a note using the old note style

> **Note**: This is a note using
> the old style and has multiple lines, but a single paragraph

> Pssst, wanna know something?
>
> You include a small description here telling users to be on the lookout

> It's not safe out there, take this Moby with you
>
> Add the `important` class to your blockquotes if you want to tell users
 to be careful about something.
{: .important}

> Ouch, don't do that!
>
> Use the `warning` class to let people know this is dangerous or they
 should pay close attention to this part of the road.
>
> You can also add more paragraphs here if your explanation is
 super complex.
{: .warning}

>**This is a crazy note**
>
> This note has tons of content in it:
>
> - List item 1
> - List item 2
>
> |Table column 1  | Table column 2 |
> |----------------|----------------|
> | Row 1 column 1 | Row 2 column 2 |
> | Row 2 column 1 | Row 2 column 2 |
>
> And another sentence to top it all off.

### Examples with FontAwesome icons

>  Pssst, wanna know something?
>
> You include a small description here telling users to be on the lookout
>
> This is an example of a note using the `{: .note-vanilla}` tag to get an icon instead of a "Note" prefix, and write your own note title.
{: .note-vanilla}


> It's not safe out there, take this Moby with you
>
> Use `{: .important-vanilla}` after your important to get an "important" icon.
{: .important-vanilla}

> Ouch, don't touch that hot Docker engine!
>
> Use `{: .warning-vanilla}` after your warning to get an icon instead of a "Warning" prefix.
>
> You can also add more paragraphs here if your explanation is
 super complex.
{: .warning-vanilla}

### Examples with both prefixed word and icon

The current CSS also supports this kind of of admonition.

> **Notes**
>
> * This is a note about a thing.
>
> *  This is another note about the same thing.
{: .note-vanilla}

## Code blocks

Rouge provides lots of different code block "hints". If you leave off the hint,
it tries to guess and sometimes gets it wrong. These are just a few hints that
we use often.

### Raw, no highlighting

The raw markup is needed to keep Liquid from interperting the things with double
braces as templating language.

```none
none with raw
{% raw %}
$ some command with {{double braces}}
$ some other command
{% endraw %}
```

### Raw, Bash

```bash
bash with raw
{% raw %}
$ some command with {{double braces}}
$ some other command
{% endraw %}
```

### Bash

```bash
$ echo "deb https://packages.docker.com/1.12/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
```

### GO

```go
incoming := map[string]interface{}{
    "asdf": 1,
    "qwer": []interface{}{},
    "zxcv": []interface{}{
        map[string]interface{}{},
        true,
        int(1e9),
        "tyui",
    },
}
```

### Python

```python
return html.format(name=os.getenv('NAME', "world"), hostname=socket.gethostname(), visits=visits)
```

### Ruby

```ruby
docker_service 'default' do
  action [:create, :start]
end
```

### JSON

Warning: Syntax highlighting breaks easily for JSON if the code you present is
not a valid JSON document. Try running your snippet through [this
linter](http://jsonlint.com/) to make sure it's valid, and remember: there is no
syntax for comments in JSON!

```json
"server": {
  "http_addr": ":4443",
  "tls_key_file": "./fixtures/notary-server.key",
  "tls_cert_file": "./fixtures/notary-server.crt"
}
```

### HTML

```html
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
</head>
</html>
```

### Markdown

```md
[![Deploy to Docker Cloud](https://files.cloud.docker.com/images/deploy-to-dockercloud.svg)](https://cloud.docker.com/stack/deploy/?repo=<repo_url>)
```

### ini

```ini
[supervisord]
nodaemon=true

[program:sshd]
command=/usr/sbin/sshd -D
```

### Dockerfile

To enable syntax highlighting for Dockerfiles, use the `conf` lexer, for now.
In the future, native Dockerfile support is coming to Rouge.

```conf
#
# example Dockerfile for https://docs.docker.com/examples/postgresql_service/
#

FROM ubuntu

RUN apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8

RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list

RUN apt-get update && apt-get install -y python-software-properties software-properties-common postgresql-9.3 postgresql-client-9.3 postgresql-contrib-9.3

# Note: The official Debian and Ubuntu images automatically ``apt-get clean``
# after each ``apt-get``

USER postgres

RUN    /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/9.3/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/9.3/main/postgresql.conf

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

CMD ["/usr/lib/postgresql/9.3/bin/postgres", "-D", "/var/lib/postgresql/9.3/main", "-c", "config_file=/etc/postgresql/9.3/main/postgresql.conf"]
```

### YAML

```yaml
authorizedkeys:
  image: dockercloud/authorizedkeys
  deployment_strategy: every_node
  autodestroy: always
  environment:
    - AUTHORIZED_KEYS=ssh-rsa AAAAB3Nsomelongsshkeystringhereu9UzQbVKy9o00NqXa5jkmZ9Yd0BJBjFmb3WwUR8sJWZVTPFL
  volumes:
    /root:/user:rw
```
