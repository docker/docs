---
title: Completion
weight: 10
description: Set up your shell to get autocomplete for Docker commands and flags
keywords: cli, shell, fish, bash, zsh, completion, options
aliases:
  - /config/completion/
---

You can generate a shell completion script for the Docker CLI using the `docker
completion` command. The completion script gives you word completion for
commands, flags, and Docker objects (such as container and volume names) when
you hit `<Tab>` as you type into your terminal.

You can generate completion scripts for the following shells:

- [Bash](#bash)
- [Zsh](#zsh)
- [fish](#fish)

## Bash

To get Docker CLI completion with Bash, you first need to install the
`bash-completion` package which contains a number of Bash functions for shell
completion.

```bash
# Install using APT:
sudo apt install bash-completion

# Install using Homebrew (Bash version 4 or later):
brew install bash-completion@2
# Homebrew install for older versions of Bash:
brew install bash-completion

# With pacman:
sudo pacman -S bash-completion
```

After installing `bash-completion`, source the script in your shell
configuration file (in this example, `.bashrc`):

```bash
# On Linux:
cat <<EOT >> ~/.bashrc
if [ -f /etc/bash_completion ]; then
    . /etc/bash_completion
fi
EOT

# On macOS / with Homebrew:
cat <<EOT >> ~/.bash_profile
[[ -r "$(brew --prefix)/etc/profile.d/bash_completion.sh" ]] && . "$(brew --prefix)/etc/profile.d/bash_completion.sh"
EOT
```

And reload your shell configuration:

```console
$ source ~/.bashrc
```

Now you can generate the Bash completion script using the `docker completion` command:

```console
$ mkdir -p ~/.local/share/bash-completion/completions
$ docker completion bash > ~/.local/share/bash-completion/completions/docker
```

## Zsh

The Zsh [completion system](http://zsh.sourceforge.net/Doc/Release/Completion-System.html)
takes care of things as long as the completion can be sourced using `FPATH`.

If you use Oh My Zsh, you can install completions without modifying `~/.zshrc`
by storing the completion script in the `~/.oh-my-zsh/completions` directory.

```console
$ mkdir -p ~/.oh-my-zsh/completions
$ docker completion zsh > ~/.oh-my-zsh/completions/_docker
```

If you're not using Oh My Zsh, store the completion script in a directory of
your choice and add the directory to `FPATH` in your `.zshrc`.

```console
$ mkdir -p ~/.docker/completions
$ docker completion zsh > ~/.docker/completions/_docker
```

```console
$ cat <<"EOT" >> ~/.zshrc
FPATH="$HOME/.docker/completions:$FPATH"
autoload -Uz compinit
compinit
EOT
```

## Fish

fish shell supports a [completion system](https://fishshell.com/docs/current/#tab-completion) natively.
To activate completion for Docker commands, copy or symlink the completion script to your fish shell `completions/` directory:

```console
$ mkdir -p ~/.config/fish/completions
$ docker completion fish > ~/.config/fish/completions/docker.fish
```
