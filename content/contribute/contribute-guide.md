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

1. You can select the **Edit this page** option in the right column of every page on [https://docs.docker.com/](/).

    This opens the GitHub editor, which means you don't need to know a lot about Git, or even about Markdown. When you save, GitHub prompts you to create a fork if you don't already have one, and to create a branch in your fork and submit the pull request.

<!-- markdownlint-disable-next-line -->
2. Fork the [docs GitHub repository]({{% param "repo" %}}). Suggest changes or add new content on your local branch, and submit a pull request (PR) to the `main` branch.

    This is the manual, more advanced version of selecting 'Edit this page' on a published docs page. Initiating a docs change in a PR from your own branch gives you more flexibility, as you can submit changes to multiple pages or files under a single pull request, and even create new topics.

    For a demo of the components, tags, Markdown syntax, styles, etc. used in [https://docs.docker.com/](/), see the components section.

## Important files

Here’s a list of some of the important files:

- `/content` contains all the pages. The site uses filesystem-based routing, so the filepath of the source files correspond to the url subpath.
- `/data/toc.yaml` defines the left-hand navigation for the docs.
- `/layouts` contains the html templates used to generate the HTML pages.

CLI reference documentation is generated from code, and is not maintained in
this repository. (They're only vendored here, either manually in `/data` or
automatically in `/_vendor`.) To update CLI reference docs, refer to the
corresponding repository containing the CLI source code.

## Pull request guidelines

Help us review your PRs more quickly by following these guidelines.

- Try not to touch a large number of files in a single PR if possible.
- Don't change whitespace or line wrapping in parts of a file you aren't editing for other reasons. Make sure your text editor isn't configured to
  automatically reformat the whole file when saving.
- We highly recommend that you [build](#build-and-preview-the-docs-locally) and [test](#test-the-docs-locally) the docs locally before submitting a PR. 
- A Netlify test runs for each PR that is against the `main` branch, and deploys the result of your PR to a staging site. The URL will be available in the **Conversation** tab. Check the staging site to verify how your changes look and fix issues, if necessary.

### Collaborate on a pull request

Unless the PR author specifically disables it, you can push commits into another
contributor's PR. You can do it from the command line by adding and fetching
their remote, checking out their branch, and adding commits to it. Even easier,
you can add commits from the GitHub web UI, by clicking the pencil icon for a
given file in the **Files** view.

If a PR consists of multiple small addendum commits on top of a more significant
one, the commit will usually be "squash-merged", so that only one commit is
merged into the `main` branch. In some scenarios where a squash and merge isn't appropriate, all commits are kept separate when merging.

### Per-PR staging on GitHub

A Netlify test runs for each PR created against the `main` branch and deploys the result of your PR to a staging site. When the site builds successfully, you will see a comment in the **Conversation** tab in the PR stating **Deploy Preview for docsdocker ready!**. Click the **Browse the preview** URL and check the staging site to verify how your changes look and fix issues, if necessary. Reviewers also check the staged site before merging the PR to protect the integrity of the docs site.

## Build and preview the docs locally

On your local machine, clone the docs repository:

```console
$ git clone {{% param "repo" %}}.git
$ cd docs
```

Then, build and run the documentation using [Docker Compose](../compose/_index.md):

```console
$ docker compose watch
```

> **Note**
>
>You need Docker Compose to build and run the docs locally. Docker Compose is included with [Docker Desktop](../desktop/_index.md). If you don't have Docker Desktop installed, follow the [instructions](../compose/install/index.md) to install Docker Compose.

When the container is built and running, visit [http://localhost:1313](http://localhost:1313) in your web browser to view the docs.

[Compose `watch`](../compose/file-watch.md) causes your
running container to rebuild itself automatically when you make changes to your
content files.

To stop the development container:

1. In your terminal, press `<Ctrl+C>` to exit the file watch mode of Compose.
2. Stop the Compose services with the `docker compose down` command.

### Test the docs locally

For testing the documentation, we use: 
- a command-line tool called [vale](https://vale.sh/) to check the style and help you find errors in your writing.
- [wjdp/htmltest](https://github.com/wjdp/htmltest) for proofing HTML markup and checking for broken links

You can use Buildx to run all the tests locally before you create your pull
request:

```console
$ docker buildx bake validate
```

You can also integrate `vale` with your text editor to get real-time feedback
on your writing. Refer to the relevant plugin for your editor:

- [VS Code](https://github.com/chrischinchilla/vale-vscode)
- [Neovim](https://github.com/neovim/nvim-lspconfig/blob/master/doc/server_configurations.md#vale_ls)
- [Emacs](https://github.com/tpeacock19/flymake-vale)
- [Jetbrains](https://vale.sh/docs/integrations/jetbrains/)
- [Sublime Text](https://github.com/errata-ai/LSP-vale-ls)

The `vale` rules that implement the Docker style guide are included in the
Docker docs repository, in the `.github/vale` directory.
