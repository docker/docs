# windows-installer

Installation [instructions](https://docs.docker.com/installation/windows/) available on the Docker documentation site.

## What is included:

- [msys-git](http://msysgit.github.io/) for tools like `OpenSSH` and `BASH`
- [VirtualBox](https://www.virtualbox.org)
- [Boot2Docker-cli management tool](https://github.com/boot2docker/boot2docker-cli)
- [Boot2Docker ISO](https://github.com/boot2docker/boot2docker)
- [Docker Client for Windows](https://github.com/docker/docker)

## Why Inno Setup?

I've chosen to make a simple installer using [Inno Setup](http://www.jrsoftware.org/)
because that is what the [msysGit](http://git-scm.com/) installer is built with.

(It also happens that I've used Inno Setup before, so I can make something faster.)

Making a simple Wix for the Boot2Docker-cli should be simple, and this can then be
used in this all-in-one installer too.

## Maintenance

See `MAINTENANCE.md` for instructions on how to update, bundle and compile the
Boot2Docker Windows Installer.