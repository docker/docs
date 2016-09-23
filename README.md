# Docs @ Docker

Welcome to the repo for our documentation. This is the source for the URL
served at docs.docker.com.

Feel free to send us pull requests and file issues. Our docs are completely
open source and we deeply appreciate contributions from our community!

## Staging

You have three options:

1. (Most performant, slowest setup) Clone this repo, [install the GitHub Pages Ruby gem](https://help.github.com/articles/setting-up-your-github-pages-site-locally-with-jekyll/), then run `jekyll serve` from within the directory.
2. (Slower, fast setup) Clone this repo, and from within the directory, run:
   `docker run -ti --rm -v "$PWD":/docs -p 4000:4000 johndmulhausen/docs`
3. (Edit entirely in the browser, no local clone) [Fork this repo in GitHub](https://github.com/docker/docker.github.io#fork-destination-box), change your fork's repository name to YOUR_GITHUB_USERNAME.github.io, and make any changes.

In the first two options, the site will be staged at `http://localhost:4000` (unless Jekyll is behaving in some non-default way).

In the third option, the site will be viewable at `http://YOUR_GITHUB_USERNAME.github.io`, about a minute after your first change is merged into your fork.

## Important files

- `/_data/toc.yaml` defines the left-hand navigation for the docs
- `/_includes/tree.html` spits out the left-hand navigation's HTML, based on the contents of `/_data/toc.yaml`
- `/_layouts/docs.html` is the HTML template file, which defines the header and footer, and includes all the JS/CSS that serves the docs content

## Relative linking for GitHub viewing

Feel free to link to `../foo.md` so that the docs are readable in GitHub, but keep in mind that Jekyll templating notation
`{% such as this %}` will render in raw text and not be processed. In general it's best to assume the docs are being read
directly on docs.docker.com.
