# Docker Toolbox for Windows

Installation [instructions](http://docs.docker.com/windows/started/) available on the Docker documentation site.

## Why Inno Setup?

I've chosen to make a simple installer using [Inno Setup](http://www.jrsoftware.org/)
because that is what the [msysGit](http://git-scm.com/) installer is built with.

(It also happens that I've used Inno Setup before, so I can make something faster.)

Making a simple Wix for the Boot2Docker-cli should be simple, and this can then be
used in this all-in-one installer too.

## Maintenance

See `MAINTENANCE.md` for instructions on how to update, bundle and compile the
Docker Toolbox Windows Installer.

## License
Docker Toolbox code is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/docker/toolbox/blob/master/LICENSE) for the full license text. 

Docker Toolbox Logo and all other related Docker artwork © Docker, Inc. 2015.  All rights reserved; not licensed for third party use.
