# Contributing to Docker Documentation

We deeply value documentation contributions from the Docker community. We'd like to make it as easy
as possible for you to work in this repository. The documentation for Docker is published at [https://docs.docker.com/](https://docs.docker.com/). The following sections guide you through the process of contributing to Docker documentation.

## Table of Contents

- [Contribution guidelines](#contribution-guidelines)
  - [Files not edited here](#files-not-edited-here)
  - [Important files](#important-files)
  - [Per-page front-matter](#per-page-front-matter)
- [Pull request guidelines](#pull-request-guidelines)
- [Collaborate on a pull request](#collaborate-on-a-pull-request)
- [Per-PR staging on GitHub](#per-pr-staging-on-github)
- [Relative linking for GitHub viewing](#relative-linking-for-github-viewing)
- [Testing changes and practical guidance](#testing-changes-and-practical-guidance)
  - [Creating tabs](#creating-tabs)
  - [Running in-page Javascript](#running-in-page-javascript)
  - [Images](#images)
  - [Style guide](#style-guide)
- [Build and preview the docs locally](#build-and-preview-the-docs-locally)
  - [Build the docs with deployment features enabled](#build-the-docs-with-deployment-features-enabled)
- [Copyright and license](#copyright-and-license)

## Contribution guidelines

The live docs are published from the `master` branch. Therefore, you must create pull requests against the `master` branch for all content updates. This includes:

- Conceptual and task-based information
- Restructuring / rewriting
- Doc bug fixing
- Fixing typos, broken links, and any grammar errors

There are two ways to contribute a pull request to the docs repository:

1. You can click **Edit this page** option  in the right column of every page on [https://docs.docker.com/](https://docs.docker.com/).

    This opens the GitHub editor, which means you don't need to know a lot about Git, or even about Markdown. When you save, Git prompts you to create a fork if you don't already have one, and to create a branch in your fork and submit the pull request.

2. Fork the [docs GitHub repository](https://github.com/docker/docker.github.io). Suggest changes or add new content on your local branch, and submit a pull request (PR) to the `master` branch.

    This is the manual, more advanced version of clicking 'Edit this page' on a published docs page. Initiating a docs changes in a PR from your own branch gives you more flexibility, as you can submit changes to multiple pages or files under a single pull request, and even create new topics.

    For a demo of the components, tags, Markdown syntax, styles, etc we use at [https://docs.docker.com/](https://docs.docker.com/), see [test.md](/test.md).

### Important files

Here’s a list of some of the important files:

- `/_data/toc.yaml` defines the left-hand navigation for the docs
- `/js/docs.js` defines most of the docs-specific JS such as the table of contents (ToC) generation and menu syncing
- `/css/style.scss` defines the docs-specific style rules
- `/_layouts/docs.html` is the HTML template file, which defines the header and footer, and includes all the JS/CSS that serves the docs content

### Files not edited here

Some files and directories are maintained in the upstream repositories. You can find a list of such files in [`_config_production.yml`](_config_production.yml#L52). Pull requests against these files will be rejected.

### Per-page front-matter

The front-matter of a given page is in a section at the top of the Markdown
file that starts and ends with three hyphens. It includes YAML content. The
following keys are supported. The title, description, and keywords are required.

| Key                    | Required  | Description                             |
|------------------------|-----------|-----------------------------------------|
| title                  | yes       | The page title. This is added to the HTML output as a `<h1>` level header. |
| description            | yes       | A sentence that describes the page contents. This is added to the HTML metadata. |
| keywords               | yes       | A comma-separated list of keywords. These are added to the HTML metadata. |
| redirect_from          | no        | A YAML list of pages which should redirect to the current page. At build time, each page listed here is created as an HTML stub containing a 302 redirect to this page. |
| notoc                  | no        | Either `true` or `false`. If `true`, no in-page TOC is generated for the HTML output of this page. Defaults to `false`. Appropriate for some landing pages that have no in-page headings.|
| toc_min                | no        | Ignored if `notoc` is set to `true`. The minimum heading level included in the in-page TOC. Defaults to `2`, to show `<h2>` headings as the minimum. |
| toc_max                | no        | Ignored if `notoc` is set to `false`. The maximum heading level included in the in-page TOC. Defaults to `3`, to show `<h3>` headings. Set to the same as `toc_min` to only show `toc_min` level of headings. |
| skip_read_time         | no        | Set to `true` to disable the 'Estimated reading time' banner for this page. |
| sitemap                | no        | Exclude the page from indexing by search engines. When set to `false`, the page is excluded from `sitemap.xml`, and a `<meta name="robots" content="noindex"/>` header is added to the page. |

Here's an example of a valid (but contrived) page metadata. The order of
the metadata elements in the front-matter isn't important.

```liquid
---
description: Instructions for installing Docker Engine on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
title: Install Docker Engine on Ubuntu
redirect_from:
- /ee/docker-ee/ubuntu/
- /engine/installation/linux/docker-ce/ubuntu/
- /engine/installation/linux/docker-ee/ubuntu/
- /engine/installation/linux/ubuntu/
- /engine/installation/linux/ubuntulinux/
- /engine/installation/ubuntulinux/
- /install/linux/docker-ce/ubuntu/
- /install/linux/docker-ee/ubuntu/
- /install/linux/ubuntu/
- /installation/ubuntulinux/
toc_max: 4
---
```

## Pull request guidelines

Help us review your PRs more quickly by following these guidelines.

- Try not to touch a large number of files in a single PR if possible.
- Don't change whitespace or line wrapping in parts of a file you are not
  editing for other reasons. Make sure your text editor is not configured to
  automatically reformat the whole file when saving.
- We highly recommend that you build the docs locally and verify your changes before submitting a PR. See the section [Build and preview the docs locally](#build-and-preview-the-docs-locally).
- A Netlify test runs for each PR that is against the `master` branch, and deploys the result of your PR to a staging site. The URL will be available at in the **Conversation** tab. Check the staging site to verify how your changes look and fix issues, if necessary.
- Use relative linking to link to other topics. See [Relative linking for GitHub viewing](#relative-linking-for-GitHub-viewing).

## Collaborate on a pull request

Unless the PR author specifically disables it, you can push commits into another
contributor's PR. You can do it from the command line by adding and fetching
their remote, checking out their branch, and adding commits to it. Even easier,
you can add commits from the Github web UI, by clicking the pencil icon for a
given file in the **Files** view.

If a PR consists of multiple small addendum commits on top of a more significant
one, the commit will usually be "squash-merged", so that only one commit is
merged into the `master` branch. In some scenarios where a squash and merge isn't appropriate, all commits are kept separate when merging.

## Per-PR staging on GitHub

A Netlify test runs for each PR created against the `master` branch and deploys the result of your PR to a staging site. When the site builds successfully, you will see a comment in the **Conversation** tab in the PR stating **Deploy Preview for docsdocker ready!**. Click the **Browse the preview** URL and check the staging site to verify how your changes look and fix issues, if necessary. Reviewers also check the staged site before merging the PR to protect the integrity of the docs site.

## Relative linking for GitHub viewing

Feel free to link to `../foo.md` so that the docs are readable in GitHub, but keep in mind that Jekyll templating notation
`{% such as this %}` renders in raw text and will not be processed. In general, it's best to assume the docs are being read
directly on [https://docs.docker.com/](https://docs.docker.com/).

## Testing changes and practical guidance

If you want to test a style change, or if you want to see how to achieve a
particular outcome with Markdown, Bootstrap, JQuery, or something else, have
a look at `test.md` (which renders in the site at `/test/`).

### Creating tabs

The use of tabs, as on pages like [https://docs.docker.com/engine/api/](/engine/api/), requires the use of HTML. The tabs use Bootstrap CSS/JS, so refer to those docs for more advanced usage. For a basic horizontal tab set, copy/paste starting from this code and implement from there. Keep an eye on those `href="#id"` and `id="id"` references as you rename, add, and remove tabs.

```html
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab1">TAB 1 HEADER</a></li>
  <li><a data-toggle="tab" data-target="#tab2">TAB 2 HEADER</a></li>
</ul>
<div class="tab-content">
  <div id="tab1" class="tab-pane fade in active">TAB 1 CONTENT</div>
  <div id="tab2" class="tab-pane fade">TAB 2 CONTENT</div>
</div>
```

For more info and a few more permutations, see [test.md](/test.md).

### Running in-page Javascript

If you need to run custom Javascript within a page, and it depends upon JQuery
or Bootstrap, make sure the `<script>` tags are at the very end of the page,
after all the content. Otherwise, the script may try to run before JQuery and
Bootstrap JS are loaded.

### Images

Don't forget to remove images that are no longer used. Keep the images sorted
in the local `images/` directory, with names that naturally group related images
together in alphabetical order. For example,  use `settings-file-share.png`
and `settings-proxies.png` instead of `file-share-settings.png` and
`proxies-settings.png`.

You may also use numbers, especially in the case of a
sequence, for example, `run-only-the-images-you-trust-1.svg`,
`run-only-the-images-you-trust-2.png`, `run-only-the-images-you-trust-3.png`.
When applicable, capture windows rather than the rectangular regions. This
eliminates unpleasant background and saves the editors the need to crop images.

On Mac, capture windows without shadows. To do this, when  you press
<kbd>Command-Shift-4</kbd>, press <kbd>Option</kbd> while clicking on the window. To disable shadows completely, run:

```bash
$ defaults write com.apple.screencapture disable-shadow -bool TRUE
$ killall SystemUIServer  # restart it.
```

You can restore the shadows later using `-bool FALSE`.

To keep the Git repository light, _please_ compress the images
(losslessly). On Mac, you can use [ImageOptim](https://imageoptim.com) for
instance. Be sure to compress the images **before** adding them to the
repository. Compressing images after adding them to the repository actually worsens the impact on the Git repo (however, ut still optimizes the bandwidth during browsing).

### Style guide

Docker does not currently maintain a style guide. Follow the examples set by the existing documentation and use your best judgment.

## Build and preview the docs locally

On your local machine, clone the docs repository:

```bash
git clone --recursive https://github.com/docker/docker.github.io.git
cd docker.github.io
```

Then, build and run the documentation using [Docker Compose](https://docs.docker.com/compose/)

```bash
docker compose up -d --build
```

> You need Docker Compose to build and run the docs locally. Docker Compose is included with [Docker Desktop](https://docs.docker.com/desktop/).
> If you don't have Docker Desktop installed, follow the [instructions](https://docs.docker.com/compose/install/) to install Docker Compose.

When the container is built and running, visit [http://localhost:4000](http://localhost:4000) in your web browser to view the docs.

To rebuild the docs after you made changes, run the `docker compose up` command
again. This rebuilds the docs, and updates the container with your changes:

```bash
docker compose up -d --build
```

To stop the staging container, use the `docker compose down` command:

```bash
docker compose down
```

### Build the docs with deployment features enabled

The default configuration for local builds of the documentation disables some
features to allow for a shorter build-time. The following options differ between
local builds, and builds that are deployed to [docs.docker.com](https://docs.docker.com/):

- search auto-completion, and generation of `js/metadata.json`
- Google Analytics
- page ratings
- `sitemap.xml` generation
- minification of stylesheets (`css/style.css`)
- adjusting "edit this page" links for content in other repositories

If you want to contribute in these areas, you can perform a "production" build
locally. To preview the documentation with deployment features enabled, set the `JEKYLL_ENV` environment variable when building the documentation:

```bash
JEKYLL_ENV=production docker compose up --build
```

When the container is built and running, visit [http://localhost:4000](http://localhost:4000) in your web browser to view the docs.

To rebuild the docs after you make changes, repeat the steps above.

## Copyright and license

Copyright 2013-2021 Docker, inc, released under the Apache 2.0 license.
