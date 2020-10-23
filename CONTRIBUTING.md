## Contributing

We value your documentation contributions, and we want to make it as easy
as possible to work in this repository. One of the first things to decide is
which branch to base your work on. If you get confused, just ask and we will
help. If a reviewer realizes you have based your work on the wrong branch, we'll
let you know so that you can rebase it.

>**Note**: To contribute code to Docker projects, see the
[Contribution guidelines](/opensource/index.md).

### Quickstart

If you spot a problem while reading the documentation and want to try to fix it
yourself, click the **Edit this page** link at the bottom of that page. The
page will open in the Github editor, which means you don't need to know a lot
about Git, or even about Markdown.

When you save, you will be prompted to create a fork if you don't already have
one, and to create a branch in your fork and submit the pull request. We hope
you give it a try!

### Overall doc improvements

Most commits will be made against the `master` branch. This includes:

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

## Collaborate on a pull request

Unless the PR author specifically disables it, you can push commits into another
contributor's PR. You can do it from the command line by adding and fetching
their remote, checking out their branch, and adding commits to it. Even easier,
you can add commits from the Github web UI, by clicking the pencil icon for a
given file in the **Files** view.

If a PR consists of multiple small addendum commits on top of a more significant
one, the commit will usually be "squash-merged", so that only one commit is
merged in. On occasion this is not appropriate and all commits will be kept
separate when merging.

## Pull request guidelines

Help us review your PRs more quickly by following these guidelines.

- Try not to touch a large number of files in a single PR if possible.

- Don't change whitespace or line wrapping in parts of a file you are not
  editing for other reasons. Make sure your text editor is not configured to
  automatically reformat the whole file when saving.

- A Netlify test runs for each PR that is against one of our long-lived
  branches like `master` and the `vnext` branches, and deploys the result of
  your PR to a staging site. The URL will be available at the bottom of the PR
  in the **Conversation** view. Check the staging site for problems and fix them
  if necessary. Reviewers will check the staging site too.

If you can think of other ways we could streamline the review process, let us
know.

## Style guide

Docker does not currently maintain a style guide. Use your best judgment, and
try to follow the example set by the existing documentation.
