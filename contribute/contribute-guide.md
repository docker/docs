---
title: Contributing guidelines
description: Guidelines for contributing to Docker's docs
keywords: contribute, guide, style guide
---

The live docs are published from the `main` branch. Therefore, you must create pull requests against the `main` branch for all content updates. This includes:

- Conceptual and task-based information
- Restructuring / rewriting
- Doc bug fixing
- Fixing typos, broken links, and any grammar errors

There are two ways to contribute a pull request to the docs repository:

1. You can click **Edit this page** option  in the right column of every page on [https://docs.docker.com/](/).

    This opens the GitHub editor, which means you don't need to know a lot about Git, or even about Markdown. When you save, Git prompts you to create a fork if you don't already have one, and to create a branch in your fork and submit the pull request.

2. Fork the [docs GitHub repository]({{ site.repo }}). Suggest changes or add new content on your local branch, and submit a pull request (PR) to the `main` branch.

    This is the manual, more advanced version of clicking 'Edit this page' on a published docs page. Initiating a docs changes in a PR from your own branch gives you more flexibility, as you can submit changes to multiple pages or files under a single pull request, and even create new topics.

    For a demo of the components, tags, Markdown syntax, styles, etc we use at [https://docs.docker.com/](/), see the Useful components section.

## Important files

Here’s a list of some of the important files:

- `/_data/toc.yaml` defines the left-hand navigation for the docs
- `/js/docs.js` defines most of the docs-specific JS such as the table of contents (ToC) generation and menu syncing
- `/css/style.scss` defines the docs-specific style rules
- `/_layouts/docs.html` is the HTML template file, which defines the header and footer, and includes all the JS/CSS that serves the docs content

### Files not edited here

Some files and directories are maintained in the upstream repositories. You can find a list of such files in `_config_production.yml`. Pull requests against these files will be rejected.

## Pull request guidelines

Help us review your PRs more quickly by following these guidelines.

- Try not to touch a large number of files in a single PR if possible.
- Don't change whitespace or line wrapping in parts of a file you are not editing for other reasons. Make sure your text editor is not configured to
  automatically reformat the whole file when saving.
- We highly recommend that you [build](#build-and-preview-the-docs-locally) and [test](#test-the-docs-locally) the docs locally before submitting a PR. 
- A Netlify test runs for each PR that is against the `main` branch, and deploys the result of your PR to a staging site. The URL will be available at in the **Conversation** tab. Check the staging site to verify how your changes look and fix issues, if necessary.

### Collaborate on a pull request

Unless the PR author specifically disables it, you can push commits into another
contributor's PR. You can do it from the command line by adding and fetching
their remote, checking out their branch, and adding commits to it. Even easier,
you can add commits from the Github web UI, by clicking the pencil icon for a
given file in the **Files** view.

If a PR consists of multiple small addendum commits on top of a more significant
one, the commit will usually be "squash-merged", so that only one commit is
merged into the `main` branch. In some scenarios where a squash and merge isn't appropriate, all commits are kept separate when merging.

### Per-PR staging on GitHub

A Netlify test runs for each PR created against the `main` branch and deploys the result of your PR to a staging site. When the site builds successfully, you will see a comment in the **Conversation** tab in the PR stating **Deploy Preview for docsdocker ready!**. Click the **Browse the preview** URL and check the staging site to verify how your changes look and fix issues, if necessary. Reviewers also check the staged site before merging the PR to protect the integrity of the docs site.

## Build and preview the docs locally

On your local machine, clone the docs repository:

```bash
git clone {{ site.repo }}.git
cd docs
```

Then, build and run the documentation using [Docker Compose](../compose/index.md)

```bash
docker compose up -d --build
```

> You need Docker Compose to build and run the docs locally. Docker Compose is included with [Docker Desktop](../desktop/index.md).
> If you don't have Docker Desktop installed, follow the [instructions](../compose/install/index.md) to install Docker Compose.

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
local builds, and builds that are deployed to [docs.docker.com](/):

- search auto-completion, and generation of `js/metadata.json`
- Google Analytics
- page feedback
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

### Test the docs locally

We use a command-line tool called [vale](https://vale.sh/) to check the style and help you find
errors in your writing. We highly recommend that you use vale to lint your documentation before
you submit your pull request.

You can run vale as a stand-alone tool using the command-line, or you can integrate it with
your editor to get real-time feedback on your writing.

To get started, follow the [vale installation instructions](https://vale.sh/docs/vale-cli/installation/)
for your operating system. To enable the vale integration for your editor, install the relevant plugin:

- [VS Code](https://github.com/errata-ai/vale-vscode)
- [Neovim](https://github.com/jose-elias-alvarez/null-ls.nvim/blob/main/doc/BUILTINS.md#vale)
- [Emacs](https://github.com/tpeacock19/flymake-vale)
- [Jetbrains](https://vale.sh/docs/integrations/jetbrains/)

The vale rules that implement the Docker style guide are included in the Docker docs repository,
in the `.github/vale` directory. Vale will automatically apply these rules when invoked in this
repository.

