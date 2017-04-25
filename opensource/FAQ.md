---
description: Overview of contributing
keywords: open, source, contributing, overview
title: FAQ for contributors
---

This section contains some frequently asked questions and tips for
troubleshooting problems in your code contribution.

- [How do I set my signature?](FAQ.md#how-do-i-set-my-signature)
- [How do I track changes from the docker repo upstream?](FAQ.md#how-do-i-track-changes-from-the-docker-repo-upstream)
- [How do I format my Go code?](FAQ.md#how-do-i-format-my-go-code)
- [What is the pre-pull request checklist?](FAQ.md#what-is-the-pre-pull-request-checklist)
- [How should I comment my code?](FAQ.md#how-should-i-comment-my-code)
- [How do I rebase my feature branch?](FAQ.md#how-do-i-rebase-my-feature-branch)

## How do I set my signature

1.  Change to the root of your `docker-fork` repository.

    ```
    $ cd docker-fork
    ```

2.  Set your `user.name` for the repository.

    ```
    $ git config --local user.name "FirstName LastName"
    ```

3.  Set your `user.email` for the repository.

    ```
    $ git config --local user.email "emailname@mycompany.com"
    ```

## How do I track changes from the docker repo upstream

Set your local repo to track changes upstream, on the `docker` repository.

1.  Change to the root of your Docker repository.

    ```
    $ cd docker-fork
    ```

2.  Add a remote called `upstream` that points to `docker/docker`.

    ```
    $ git remote add upstream https://github.com/moby/moby.git
    ```

## How do I format my Go code

Run `gofmt -s -w filename.go` on each changed file before committing your changes.
Most [editors have plug-ins](https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins) that do the formatting automatically.

## What is the pre-pull request checklist

* Sync and cleanly rebase on top of Docker's `master`; do not mix multiple branches
  in the pull request.

* Squash your commits into logical units of work using
  `git rebase -i` and `git push -f`.

* If your code requires a change to tests or documentation, include code, test,
and documentation changes in the same commit as your code; this ensures a
revert would remove all traces of the feature or fix.

* Reference each issue in your pull request description (`#XXXX`).

## How should I comment my code?

The Go blog wrote about code comments, it is <a href="http://goo.gl/fXCRu"
target="_blank">a single page explanation</a>. A summary follows:

- Comments begin with two forward `//` slashes.
- To document a type, variable, constant, function, or even a package, write a
regular comment directly preceding the elements declaration, with no intervening blank
line.
- Comments on package declarations should provide general package documentation.
- For packages that need large amounts of introductory documentation: the
package comment is placed in its own file.
- Subsequent lines of text are considered part of the same paragraph; you must
leave a blank line to separate paragraphs.
-  Indent pre-formatted text relative to the surrounding comment text (see gob's doc.go for an example).
- URLs are converted to HTML links; no special markup is necessary.

## How do I rebase my feature branch?

Always rebase and squash your commits before making a pull request.

1.  Fetch any of the last minute changes from `docker/docker`.

    ```
    $ git fetch upstream master
    ```

3.  Start an interactive rebase.

    ```
    $ git rebase -i upstream/master
    ```

4.  Rebase opens an editor with a list of commits.

    ```
    pick 1a79f55 Tweak some of images
    pick 3ce07bb Add a new line
    ```

    If you run into trouble, `git --rebase abort` removes any changes and gets you
back to where you started.

4.  Replace the `pick` keyword with `squash` on all but the first commit.

    ```
    pick 1a79f55 Tweak some of images
    squash 3ce07bb Add a new line
    ````

    After closing the file, `git` opens your editor again to edit the commit
    message.

5.  Edit and save your commit message.

    ```
    $ git commit -s
    ```

    Make sure your message includes your signature.

8.  Push any changes to your fork on GitHub, using the `-f` option to
force the previous change to be overwritten.

    ```
    $ git push -f origin my-keen-feature
    ```

## How do I update vendor package from upstream ?

1.  If you are not using the development container, download the
    [vndr](https://github.com/LK4D4/vndr) vendoring tool. The `vndr`
    tool is included in the development container.
    
2.  Edit the package version in `vendor.conf` to use the package you want to use, such as
    `github.com/gorilla/mux`.
    
3.  Run `vndr <package-name>`. For example:

    ```bash
    vndr github.com/gorilla/mux
    ```
