---
description: components and formatting examples used in Docker's docs
title: Useful components and formatting examples
toc_max: 3
---

This page contains the components, tags, and styles, we use for the docs. We explain the code behind the published page and demo the effects.
For components and controls we are using [Bootstrap](https://getbootstrap.com)

> Note
>
> Check the raw Markdown file to see how the HTML or Markdown is formatted.

## Useful components

### Tabs

Here are some tabs:

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

To add Markdown inside the content of the `<div>`, just add `markdown="1"` as an attribute in the HTML for the `<div>` so that Kramdown renders it.

#### Synchronizing multiple tab groups

Consider an example where you have something like one tab per language, and
you have multiple tab sets like this on a page. You might want them to all
toggle together. We have Javascript that loads on every page that lets you
do this by setting the `data-group` attributes to be the same. The
`data-target` attributes still need to point to unique div IDs.

In the following example, selecting `Go` or `Python` in one tab set toggles the
other tab set to match.

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

### Cards

In a Bootstrap row, your columns need to add up to 12. Here are three cards in
a row, each of which takes up 1/3 (4/12) of the row.

<div class="row">
  <div class="panel col-xs-12 col-md-4">This takes up 1/3 of the row unless the screen is small,
then it takes up the whole row.</div>
  <div class="panel col-xs-12 col-md-4">This takes up 1/3 of the row unless the screen is small,
then it takes up the whole row.</div>
  <div class="panel col-xs-12 col-md-4">This takes up 1/3 of the row unless the screen is small,
then it takes up the whole row.</div>
</div>

### Expand/Collapse accordions

This
implementation makes use of the `.panel-heading` classes in
`_utilities.scss.md`,
along with [FontAwesome icons](http://fontawesome.io/cheatsheet/){: target="_blank" rel="noopener" class="_" }
<i class="fa fa-caret-down" aria-hidden="true"></i> (fa-caret-down) and
<i class="fa fa-caret-up" aria-hidden="true"></i> (fa-caret-up).

<div class="panel panel-default">
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample1" style="cursor: pointer">
    Docker hello-world example
    <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample1">
<pre><code>
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
2db29710123e: Pull complete
Digest: sha256:cc15c5b292d8525effc0f89cb299f1804f3a725c8d05e158653a563f15e4f685
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:

1. The Docker client contacted the Docker daemon.
2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
3. The Docker daemon created a new container from that image which runs the
4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:

$ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
</code></pre>
    </div>
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample2"  style="cursor: pointer"> Another Sample <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample2">
<p>
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.</p>
 </div>
</div>

Adding `block` to the `div` class `collapse` gives you some padding around the
sample content. This works nicely for standard text. If you have a code sample,
the padding renders as white space around the code block grey background. If we
don't like this effect, we can remove `block` for code samples.

The `style="cursor: pointer"` tag enables the expand/collapse functionality to
work on mobile. (You can use the [Xcode iPhone simulator](https://developer.apple.com/library/content/documentation/IDEs/Conceptual/iOS_Simulator_Guide/GettingStartedwithiOSSimulator/GettingStartedwithiOSSimulator.html#//apple_ref/doc/uid/TP40012848-CH5-SW4){: target="_blank" rel="noopener" class="_" } to test on mobile.)

> Note
>
>Make sure `data-target`'s and `id`'s match, and are unique.
>
>For each drop-down, the value for `data-target` and
`collapse` `id` must match, and id's must be unique per page. In this example,
we name these `collapseSample1` and `collapseSample2`.

### Badges

You can have <span class="badge badge-info">badges</span>. You can also have
<span class="badge badge-warning">yellow badges</span> or
<span class="badge badge-danger">red badges</span>.

#### Badges as links

You can make a badge a link. Wrap the `<span>` with an `<a>` (not the other way
around) so that the text on the badge is still white.

```html
<a href="/contribute/overview/" target="_blank" rel="noopener" class="_"><span class="badge badge-info">Test</span></a>
```

<a href="/contribute/overview/" target="_blank" rel="noopener" class="_"><span class="badge badge-info" >Test</span></a>

### Tooltips

To add a tooltip to any element, set `data-toggle="tooltip"` and set a `title`.
Hovering over the element with the mouse pointer will make it visible. Tooltips
are not visible on mobile devices or touchscreens, so don't rely on them as the
only way to communicate important info.

```html
<span class="badge badge-info" data-toggle="tooltip" data-placement="right" title="Open the test page">Test</span>
```
<button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="left" title="Tooltip on left">Tooltip on left</button>

  <button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="top" title="Tooltip on top">Tooltip on top</button>

  <button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="bottom" title="Tooltip on bottom">Tooltip on bottom</button>

  <button type="button" class="btn btn-default" data-toggle="tooltip" data-placement="right" title="Tooltip on right">Tooltip on right</button>

You can optionally set the `data-placement` attribute to `top`, `bottom`,
`middle`, `center`, `left`, or `right`. Only set it if not setting it causes
layout issues.

You don't have to use HTML. You can also set these attributes using Markdown.

```markdown
This is a paragraph that has a tooltip. We position it to the left so it doesn't align with the middle top of the paragraph (that looks weird).
{:data-toggle="tooltip" data-placement="left" title="Markdown tooltip example"}
```

This is a paragraph that has a tooltip. We position it to the left so it doesn't align with the middle top of the paragraph (that looks weird).
{:data-toggle="tooltip" data-placement="left" title="Markdown tooltip example"}

You can also put tooltips on badges with a link.

<a href="/contribute/components" target="_blank" rel="noopener" class="_"><span class="badge badge-info" data-toggle="tooltip" data-placement="right" title="Open the test page (in a new window)">Test</span></a>


## Formatting examples

### Links

It is best practice if [a link opens in a new window](https://docker.com/){: target="_blank" rel="noopener" class="_" }

You can also have [a markdown link to a custom target ID](#formatting-examples)

#### Links to auto-generated content

An example of a link to an auto-generated reference page that we pull in during docs builds:
[/engine/reference/builder/#env](/engine/reference/builder/#env).

  - If you can't find a reference page in the `docker.github.io`
  GitHub repository, but see it out on `docs.docker.com`, you can
  surmise that it's probably auto-generated from the codebase.
  (FYI, to view the Markdown source for the file, just click
  **Edit this page** on `docs.docker.com`. But don't use that URL in your docs.)

  - Go to the file in a web browser, grab everything after the domain name
  from the URL, and use that as the link in your docs file.

  - Keep in mind that this link doesn't resolve until you merge the PR and
  your docs are published on [docs.docker.com](/).

### Images

- A small cute image: ![a small cute image](/images/footer_moby_icon.png)

- A small cute image that is a link. The extra newline here makes it not show
  inline:

  [![a small cute image](/images/footer_moby_icon.png)](https://www.docker.com/)

- A big wide image: ![a pretty wide image](/images/banner_image_24512.png)

- The same as above but using HTML: <img src="/images/banner_image_24512.png" alt="a pretty wide image using HTML"/>

### Videos

To embed a YouTube video on a docs page, open the video on YouTube, click
**Share** > **Embed** and then copy the code displayed.

For example, the video embedded on the Get Started page has the following code:

```html
<iframe width="560" height="315" src="https://www.youtube-nocookie.com/embed/iqqDU2crIEQ?start=30" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
```

You can also add a link to a YouTube video like this:

[Docker 101: Introduction to Docker](https://www.youtube.com/watch?v=V9IJj4MzZBc "Docker 101: Introduction to Docker"){:target="_blank" rel="noopener" class="_"}

To make the `.png` shown above, first take a screen snap of the YouTube video
you want to use, then use a graphics app to overlay a play button onto the
image.

For the overlay, you can use the play button at
[/docker-cloud/images/](https://github.com/docker/docker.github.io/tree/master/docker-cloud/images).

### Lists

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

### Tables

| Permission level                                                         | Access                                                       |
|:-------------------------------------------------------------------------|:-------------------------------------------------------------|
| **Bold** or _italic_ within a table cell. Next cell is empty on purpose. |                                                              |
|                                                                          | Previous cell is empty. A `--flag` in mono text.             |
| Read                                                                     | Pull                                                         |
| Read/Write                                                               | Pull, push                                                   |
| Admin                                                                    | All of the above, plus update description, create, and delete |

The alignment of the cells in the source doesn't really matter. The ending pipe
character is optional (unless the last cell is supposed to be empty). The header
row and separator row are optional.

### Call outs

We support these broad categories of call outs:

- Notes (no Liquid tag required)
- Important, which use the `{: .important}` tag
- Warning , which use the `{: .warning}` tag


Examples are shown in the following sections.

#### Note

A standard note is formatted like this:

```markdown
> Handling transient errors
>
> Note the way the `get_hit_count` function is written. This basic retry
> loop lets us attempt our request multiple times if the redis service is
> not available. This is useful at startup while the application comes
> online, but also makes our application more resilient if the Redis
> service needs to be restarted anytime during the app's lifetime. In a
> cluster, this also helps handling momentary connection drops between
> nodes.
```

A note renders as follows:

  > **Note**
  >
  > Note the way the `get_hit_count` function is written. This basic retry
  > loop lets us attempt our request multiple times if the redis service is
  > not available. This is useful at startup while the application comes
  > online, but also makes our application more resilient if the Redis
  > service needs to be restarted anytime during the app's lifetime. In a
  > cluster, this also helps handling momentary connection drops between
  > nodes.


#### Important

Add the `important` class to your blockquotes if you want to tell users to be careful about something:

The 'important' class renders as follows:

> **Important**
>
> Treat access tokens like your password and keep them secret. Store your
> tokens securely (for example, in a credential manager).
{: .important}

#### Warning

Use the `warning` class to let people know this is dangerous or they should pay close attention to this part of the road before moving on:

The 'warning' class renders as follows:

> **Warning**
>
> Removing Volumes
>
> By default, named volumes in your compose file are NOT removed when running
> `docker-compose down`. If you want to remove the volumes, you will need to add
> the `--volumes` flag.
>
> The Docker Dashboard does _not_ remove volumes when you delete the app stack.
{: .warning}

### Includes

If you want to reuse content in multiple pages you can use the liquid formatting to add it to your doc. 
For example:

```liquid
{% include upgrade-cta.html%}
```

Inserts all the content in `upgrade-cta.html` into your document. 

### Code blocks

Rouge provides lots of different code block "hints". If you leave off the hint,
it tries to guess and sometimes gets it wrong. These are just a few hints that
we use often.

#### Raw, no highlighting

The raw markup is needed to keep Liquid from interpreting the things with double
braces as templating language.

{% raw %}
```none
none with raw
$ some command with {{double braces}}
$ some other command
```
{% endraw %}

#### Raw, Bash

{% raw %}
```bash
bash with raw
$ some command with {{double braces}}
$ some other command
```
{% endraw %}

#### Bash

```bash
$ echo "deb https://packages.docker.com/1.12/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
```

#### Go

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

#### PowerShell

```powershell
Install-Module DockerMsftProvider -Force
Install-Package Docker -ProviderName DockerMsftProvider -Force
[System.Environment]::SetEnvironmentVariable("DOCKER_FIPS", "1", "Machine")
Expand-Archive docker-18.09.1.zip -DestinationPath $Env:ProgramFiles -Force
```

#### Python

```python
return html.format(name=os.getenv('NAME', "world"), hostname=socket.gethostname(), visits=visits)
```

#### Ruby

```ruby
docker_service 'default' do
  action [:create, :start]
end
```

#### JSON

```json
"server": {
  "http_addr": ":4443",
  "tls_key_file": "./fixtures/notary-server.key",
  "tls_cert_file": "./fixtures/notary-server.crt"
}
```

#### HTML

```html
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
</head>
</html>
```

#### Markdown

```markdown
# Hello
```

#### ini

```ini
[supervisord]
nodaemon=true

[program:sshd]
command=/usr/sbin/sshd -D
```

#### Dockerfile

```dockerfile
# syntax=docker/dockerfile:1

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

#### YAML

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
