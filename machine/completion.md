---
description: Install Machine command-line completion
keywords: machine, docker, orchestration, cli, reference
title: Command-line completion
---

Docker Machine comes with [command completion](http://en.wikipedia.org/wiki/Command-line_completion)
for the bash and zsh shell.

## Installing Command Completion

### Bash

Make sure bash completion is installed. If you are using a current version of
Linux in a non-minimal installation, bash completion should be available.

On a Mac, install with `brew install bash-completion`.

Place the completion script in `/etc/bash_completion.d/` as follows:

*   On a Mac:

    ```shell
    sudo curl -L https://raw.githubusercontent.com/docker/machine/v{{site.machine_version}}/contrib/completion/bash/docker-machine.bash -o `brew --prefix`/etc/bash_completion.d/docker-machine
    ```

*   On a standard Linux installation:

    ```shell
    sudo curl -L https://raw.githubusercontent.com/docker/machine/v{{site.machine_version}}/contrib/completion/bash/docker-machine.bash -o /etc/bash_completion.d/docker-machine
    ```

Completion is available upon next login.


### Zsh

Place the completion script in your a `completion` file within the ZSH
configuration directory, such as `~/.zsh/completion/`.

```shell
mkdir -p ~/.zsh/completion
curl -L https://raw.githubusercontent.com/docker/machine/v{{site.machine_version}}/contrib/completion/zsh/_docker-machine > ~/.zsh/completion/_docker-machine
```

Include the directory in your `$fpath`, by adding a like the following to the
`~/.zshrc` configuration file.

```shell
fpath=(~/.zsh/completion $fpath)
```

Make sure `compinit` is loaded or do it by adding in `~/.zshrc`:

```shell
autoload -Uz compinit && compinit -i
```

Then reload your shell:

```shell
exec $SHELL -l
```

## Available completions

Depending on what you typed on the command line so far, it completes:

- commands and their options
- container IDs and names
- image repositories and image tags
- file paths

## Where to go next

* [Get started with a local VM](get-started.md)
* [Machine command-line reference](reference/index.md)
* [Machine drivers](drivers/index.md)
* [Machine concepts and help](concepts.md)
