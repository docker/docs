

# This repo is for preview/dev purposes only

The docs team is in the process of migrating to this repo. During this time we're also
converting from a Hugo-based doc system to a Jekyll-based doc system. That means there
are lots of content errors and so forth at the preview URL. We're on it.

While this repo is not the source of truth, it's important to us that it be viewable
so contributors can see what we're doing, and so that they can migrate open pull requests
from other repos such as `docker/docker` that make doc changes, into this repo.

## Timeline of migration

- During the week of Monday, September 26th, any existing docs PRs need to be migrated over or merged.
- We’ll do one last “pull” from the various docs repos on Wednesday, September 28th, at which time the docs/ folders in the various repos will be emptied.
- Between the 28th and full cutover, the docs team will be testing the new repo and making sure all is well across every page.
- Full cutover (production is drawing from the new repo, new docs work is pointed at the new repo, dissolution of old docs/ folders) is complete on Monday, October 3rd.

# Docs @ Docker

Welcome to the repo for our documentation. This is the source for the URL
served at docs.docker.com.

Feel free to send us pull requests and file issues. Our docs are completely
open source and we deeply appreciate contributions from our community!

## Staging

You have three options:

1. (Most performant, slowest setup) Clone this repo, [install the GitHub Pages Ruby gem](https://help.github.com/articles/setting-up-your-github-pages-site-locally-with-jekyll/), then run `jekyll serve` from within the directory.
2. (Slower, fast setup) Clone this repo, and from within the directory, run:
   `docker run -ti --rm -v "$PWD":/docs -p 4000:4000 docs/docstage`
3. (Edit entirely in the browser, no local clone) Fork this repo in GitHub, change your fork's repository name to `YOUR_GITHUB_USERNAME.github.io`, and make any changes.

In the first two options, the site will be staged at `http://localhost:4000` (unless Jekyll is behaving in some non-default way).

In the third option, the site will be viewable at `http://YOUR_GITHUB_USERNAME.github.io`, about a minute after your first change is merged into your fork.

## Important files

- `/_data/toc.yaml` defines the left-hand navigation for the docs
- `/js/menu.js` defines most of the docs-specific JS such as TOC generation and menu syncing
- `/css/documentation.css` defines the docs-specific style rules
- `/_layouts/docs.html` is the HTML template file, which defines the header and footer, and includes all the JS/CSS that serves the docs content

## Relative linking for GitHub viewing

Feel free to link to `../foo.md` so that the docs are readable in GitHub, but keep in mind that Jekyll templating notation
`{% such as this %}` will render in raw text and not be processed. In general it's best to assume the docs are being read
directly on docs.docker.com.

## Style guide

If you have questions about how to write for Docker's documentation, please see
the [style guide](https://docs.docker.com/opensource/doc-style/). The style guide provides
guidance about grammar, syntax, formatting, styling, language, or tone. If
something isn't clear in the guide, please submit an issue to let us know or
submit a pull request to help us improve it.

### Generate the man pages

For information on generating man pages (short for manual page), see the README.md
document in [the man page directory](https://github.com/docker/docker/tree/master/man)
in this project.

## Copyright and license

Code and documentation copyright 2016 Docker, inc, released under the Apache 2.0 license.
