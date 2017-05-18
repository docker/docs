---
description: Install Machine command-line completion
keywords: machine, docker, orchestration, cli, reference
title: Command-line completion
---

Docker Machine comes with [command completion](http://en.wikipedia.org/wiki/Command-line_completion)
for the bash and zsh shell.

## Installing Command Completion

### Bash

Make sure bash completion is installed. If you are using a current version of Linux in a non-minimal installation, bash completion should be available.

On a Mac, install with `brew install bash-completion`.

Place the completion script in `/etc/bash_completion.d/` as follows:

*   On a Mac:

    ```none
    curl -L https://raw.githubusercontent.com/docker/docker/master/contrib/completion/bash/docker > `brew --prefix`/etc/bash_completion.d/docker
    ```

*   On a standard Linux installation:

    ```none
    curl -L https://raw.githubusercontent.com/docker/docker/master/contrib/completion/bash/docker > /etc/bash_completion.d/docker
    ```

Completion will be available upon next login.


### Zsh

Place the completion script in your `/path/to/zsh/completion`, using e.g. `~/.zsh/completion/`

    mkdir -p ~/.zsh/completion
    curl -L https://raw.githubusercontent.com/docker/docker/master/contrib/completion/zsh/_docker > ~/.zsh/completion/_docker-machine

Include the directory in your `$fpath`, e.g. by adding in `~/.zshrc`

    fpath=(~/.zsh/completion $fpath)

Make sure `compinit` is loaded or do it by adding in `~/.zshrc`

    autoload -Uz compinit && compinit -i

Then reload your shell

    exec $SHELL -l


<!--[metadata]>
## Available completions

**TODO**
<![end-metadata]-->
