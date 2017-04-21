# Docs @ Docker

Welcome to the repo for our documentation. This is the source for
[https://docs.docker.com/](https://docs.docker.com/).

Feel free to send us pull requests and file issues. Our docs are completely
open source and we deeply appreciate contributions from our community!

## Providing feedback

We really want your feedback, and we've made it easy. You can edit, rate, or
file an issue at the bottom of every page on
[https://docs.docker.com/](https://docs.docker.com/).

**Please only file issues about the documentation in this repository.** One way
to think about this is that you should file a bug here if your issue is that you
don't see something that should be in the docs, or you see something incorrect
or confusing in the docs.

- If your problem is a general question about how to configure or use Docker,
  ask in [https://forums.docker.com](https://forums.docker.com) instead.

- If you have an idea for a new feature or behavior change in a specific aspect
  of Docker, or have found a bug in part of Docker, please file that issue in
  the project's code repository.

## Contributing

We value your documentation contributions, and we want to make it as easy
as possible to work in this repository. One of the first things to decide is
which branch to base your work on. If you get confused, just ask and we will
help. If a reviewer realizes you have based your work on the wrong branch, we'll
let you know so that you can rebase it.

>**Note**: To contribute code to Docker projects, see the
[Contribution guidelines](https://docs.docker.com/opensource/project/who-written-for/).

### Files not edited here

Files and directories listed in the `path:` keys in
[`.NOT_EDITED_HERE.yaml`](.NOT_EDITED_HERE.yaml) are maintained in other
repositories and should not be edited in this one. Pull requests against these
files will be rejected. Make your edits to the files in the repository and path
in the `source:` key in the YAML file.

### Overall doc improvements

Most commits will be made against the `master` branch. This include:

- Conceptual and task-based information not specific to new features
- Restructuring / rewriting
- Doc bug fixing
- Typos and grammar errors

One quirk of this project is that the `master` branch is where the live docs are
published from, so upcoming features can't be documented there. See
[Specific new features for a project](#specific-new-features-for-a-project)
for how to document upcoming features. These feature branches will be periodically
merged with `master`, so don't worry about fixing typos and documentation bugs
there.

>Do you enjoy creating graphics? Good graphics are key to great documentation,
and we especially value contributions in this area.

### Specific new features for a project

Our docs cover many projects which release at different times. **If, and only if,
your pull request relates to a currently unreleased feature of a project, base
your work on that project's `vnext` branch.** These branches were created by
cloning `master` and then importing a project's `master` branch's docs into it
(at the time of the migration), in a way that preserved the commit history. When
a project has a release, its `vnext` branch will be merged into `master` and your
work will be visible on [https://docs.docker.com/](https://docs.docker.com/).

The following `vnext` branches currently exist:

- **[vnext-engine](https://github.com/docker/docker.github.io/tree/vnext-engine):**
  docs for upcoming features in the [docker/docker](https://github.com/moby/moby/)
  project

- **[vnext-compose](https://github.com/docker/docker.github.io/tree/vnext-compose):**
  docs for upcoming features in the [docker/compose](https://github.com/docker/compose/)
  project

- **[vnext-distribution](https://github.com/docker/docker.github.io/tree/vnext-distribution):**
  docs for upcoming features in the [docker/distribution](https://github.com/docker/distribution/)
  project

- **[vnext-swarm](https://github.com/docker/docker.github.io/tree/vnext-swarm):**
  docs for upcoming features in the [docker/swarm](https://github.com/docker/swarm/)
  project

- **[vnext-toolbox](https://github.com/docker/docker.github.io/tree/vnext-toolbox):**
  docs for upcoming features in the [docker/toolbox](https://github.com/docker/toolbox/)
  project

## Per-PR staging on GitHub

For every PR against `master` and all the long-lived branches, a staged version
of the site is built using Netlify. If the site builds, you will see
**deploy/netlify — Deploy preview ready**. Otherwise, you will see an error.
Click **Details** to review the staged site or the errors that prevented it from
building. Review the staged site and amend your commit if necessary. Reviewers
will also check the staged site before merging the PR, to protect the integrity
of [https://docs.docker.com/](https://docs.docker.com/).

## Staging the docs

You have three options:

1.  On your local machine, clone this repo and run our staging container:

    ```bash
    git clone https://github.com/docker/docker.github.io.git
    cd docker.github.io
    docker-compose up
    ```

    If you haven't got Docker Compose installed,
    [follow these installation instructions](https://docs.docker.com/compose/install/).

    The container runs in the background and incrementally rebuilds the site each
    time a file changes. You can keep your browser open to http://localhost:4000/
    and refresh to see your changes. The container runs in the foreground, but
    you can use `CTRL+C` to get the command prompt back. To stop the container,
    issue the following command:

    ```bash
    docker-compose down
    ```

2.  Install Jekyll and GitHub Pages on your local machine.

    a. Clone this repo by running:

       ```bash
       git clone https://github.com/docker/docker.github.io.git
       ```

    b. Install Ruby 2.3 or later as described in [Installing Ruby]
    (https://www.ruby-lang.org/en/documentation/installation/).

    c. Install Bundler:

       ```bash
       gem install bundler
       ```

    d. If you use Ubuntu, install packages required for the Nokogiri HTML
       parser:

       ```bash
       sudo apt-get install ruby-dev zlib1g-dev liblzma-dev
       ```

    e. Install Jekyll and other required dependencies:

       ```bash
       bundle install
       ```

       >**Note**: You may have to install some packages manually.   

    f. Change the directory to `docker.github.io`.

    g. Use the `jekyll serve` command to continuously build the HTML output.

    The `jekyll serve` process runs in the foreground, and starts a web server
    running on http://localhost:4000/ by default. To stop it, use `CTRL+C`.
    You can continue working in a second terminal and Jekyll will rebuild the
    website incrementally. Refresh the browser to preview your changes.

3.  Use Github Pages, with or without a local clone. Fork this repo in GitHub,
    change your fork's repository name to `YOUR_GITHUB_USERNAME.github.io`, and
    make changes to the Markdown files in your `master` branch. Browse to
    https://\<YOUR_GITHUB_USERNAME\>.github.io/ to preview the changes.

## Important files

- `/_data/toc.yaml` defines the left-hand navigation for the docs
- `/js/menu.js` defines most of the docs-specific JS such as TOC generation and menu syncing
- `/css/documentation.css` defines the docs-specific style rules
- `/_layouts/docs.html` is the HTML template file, which defines the header and footer, and includes all the JS/CSS that serves the docs content

## Relative linking for GitHub viewing

Feel free to link to `../foo.md` so that the docs are readable in GitHub, but keep in mind that Jekyll templating notation
`{% such as this %}` will render in raw text and not be processed. In general it's best to assume the docs are being read
directly on [https://docs.docker.com/](https://docs.docker.com/).

## Style guide

If you have questions about how to write for Docker's documentation, please see
the [style guide](https://docs.docker.com/opensource/doc-style/). The style guide provides
guidance about grammar, syntax, formatting, styling, language, or tone. If
something isn't clear in the guide, please submit an issue to let us know or
submit a pull request to help us improve it.

### Testing changes and practical guidance

If you want to test a style change, or if you want to see how to achieve a
particular outcome with Markdown, Bootstrap, JQuery, or something else, have
a look at `test.md` (which renders in the site at `/test/`).

### Per-page front-matter

The front-matter of a given page is in a section at the top of the Markdown
file that starts and ends with three hyphens. It includes YAML content. The
following keys are supported. The title, description, and keywords are required.

| Key                    | Required  | Description                             |
|------------------------|-----------|-----------------------------------------|
| title                  | yes       | The page title. This is added to the HTML output as a `<h1>` level header. |
| description            | yes       | A sentence that describes the page contents. This is added to the HTML metadata. |
| keywords               | yes       | A comma-separated list of keywords. These are added to the HTML metadata. |
| redirect_from          | no        | A YAML list of pages which should redirect to THIS page. At build time, each page listed here is created as a HTML stub containing a 302 redirect to this page. |
| notoc                  | no        | Either `true` or `false`. If `true`, no in-page TOC is generated for the HTML output of this page. Defaults to `false`. Appropriate for some landing pages that have no in-page headings.|
| toc_min                | no        | Ignored if `notoc` is set to `true`. The minimum heading level included in the in-page TOC. Defaults to `2`, to show `<h2>` headings as the minimum. |
| toc_max                | no        | Ignored if `notoc` is set to `false`. The maximum heading level included in the in-page TOC. Defaults to `3`, to show `<h3>` headings. Set to the same as `toc_min` to only show `toc_min` level of headings. |
| tree                   | no        | Either `true` or `false`. Set to `false` to disable the left-hand site-wide navigation for this page. Appropriate for some pages like the search page or the 404 page. |
| no_ratings             | no        | Either `true` or `false`. Set to `true` to disable the page-ratings applet for this page. Defaults to `false`. |

The following is an example of valid (but contrived) page metadata. The order of
the metadata elements in the front-matter is not important.

```liquid
---
description: Instructions for installing Docker on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/ubuntulinux/
- /installation/ubuntulinux/
- /engine/installation/linux/ubuntulinux/
title: Get Docker for Ubuntu
toc_min: 1
toc_max: 6
tree: false
no_ratings: true
---
```

### Creating tabs

The use of tabs, as on pages like https://docs.docker.com/engine/api/, requires
the use of HTML. The tabs use Bootstrap CSS/JS, so refer to those docs for more
advanced usage. For a basic horizontal tab set, copy/paste starting from this
code and implement from there. Keep an eye on those `href="#id"` and `id="id"`
references as you rename, add, and remove tabs.

```
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab1">TAB 1 HEADER</a></li>
  <li><a data-toggle="tab" data-target="#tab2">TAB 2 HEADER</a></li>
</ul>
<div class="tab-content">
  <div id="tab1" class="tab-pane fade in active">TAB 1 CONTENT</div>
  <div id="tab2" class="tab-pane fade">TAB 2 CONTENT</div>
</div>
```

For more info and a few more permutations, see `test.md`.

### Running in-page Javascript

If you need to run custom Javascript within a page, and it depends upon JQuery
or Bootstrap, make sure the `<script>` tags are at the very end of the page,
after all the content. Otherwise the script may try to run before JQuery and
Bootstrap JS are loaded.

> **Note**: In general, this is a bad idea.

## Copyright and license

Code and documentation copyright 2017 Docker, inc, released under the Apache 2.0 license.
