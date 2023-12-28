---
description: components and formatting examples used in Docker's docs
title: Code blocks
toc_max: 3
---

Rouge provides lots of different code block "hints". If you leave off the hint,
it tries to guess and sometimes gets it wrong. These are just a few hints that
we use often.

## Variables

If your example contains a placeholder value that's subject to change,
use the format `<[A-Z_]+>` for the placeholder value: `<MY_VARIABLE>`

```none
export name=<MY_NAME>
```

This syntax is reserved for variable names, and will cause the variable to
be rendered in a special color and font style.

## Interactive code

> **Experimental**
>
> This component is experimental. The API, appearance, and features may change.
{ .experimental }

This feature is realized using [codapi](https://codapi.org/). You can make code
blocks executable by setting an `interactive` attribute on the code block.

````md
```bash {interactive=true}
function greet() {
  echo "Hello $1"
}

greet "World"
```
````

A **Run** button appears below the code block.

```bash {interactive=true}
function greet() {
  echo "Hello $1"
}

greet "World"
```

You can make the code blocks editable by setting `editor=true`. This adds an
**Edit** button to the right of the **Run** button.

```bash {interactive=true,editor=true}
function greet() {
  echo "Hello $1"
}

greet "World"
```

### Known limitations

Syntax highlighting

: Focusing on an editable code block causes syntax highlighting to go away.
  This happens because syntax highlighting is created using HTML elements
  injected at build-time, representing the syntax tree detected by the lexer.

  When focused, code blocks are re-injected with the plain text code, without
  the wrapper elements representing the syntax tree.

  For better compatibility with editable code blocks, we should consider moving
  off the built-in `chroma` generator to something like `highlight.js` instead,
  which handles syntax highlighting on the fly, at runtime.

Languages

: This prototype uses codapi cloud. The available execution sandboxes and their
  configurations cover only a narrow set of the code examples in our docs.
  Notably, there's no sandbox which includes the Docker runtime and CLI, as well
  as other fundamental tools like `curl`.

Formatting

: Many of our examples weren't originally written and formatted as runnable
  examples. For example, `console` blocks often contain a mix of input/output
  lines, which would either require re-formatting or some form of pre-processing
  before it's sent to an execution environment, filtering output lines and
  stripping prefixes.

Unsupported features

: codapi supports features that aren't covered in this prototype, including:

  - [Code cells](https://github.com/nalgeon/codapi-js/blob/main/docs/html.md#code-cells)
  - [Custom actions](https://github.com/nalgeon/codapi-js/blob/main/docs/html.md#custom-actions)
  - [Files](https://github.com/nalgeon/codapi-js/blob/main/docs/html.md#files)
  - [Templates](https://github.com/nalgeon/codapi-js/blob/main/docs/html.md#templates)

## Bash

Use the `bash` language code block when you want to show a Bash script:

```bash
#!/usr/bin/bash
echo "deb https://packages.docker.com/1.12/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
```

If you want to show an interactive shell, use `console` instead.
In cases where you use `console`, make sure to add a dollar character
for the user sign:

```console
$ echo "deb https://packages.docker.com/1.12/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
```

## Go

```go
incoming := map[string]interface{}{
    "asdf": 1,
    "qwer": []interface{}{},
    "zxcv": []interface{}{
        map[string]interface{}{},
        true,
        int(1e9),
        "tyui",
    },
}
```

## PowerShell

```powershell
Install-Module DockerMsftProvider -Force
Install-Package Docker -ProviderName DockerMsftProvider -Force
[System.Environment]::SetEnvironmentVariable("DOCKER_FIPS", "1", "Machine")
Expand-Archive docker-18.09.1.zip -DestinationPath $Env:ProgramFiles -Force
```

## Python

```python
return html.format(name=os.getenv('NAME', "world"), hostname=socket.gethostname(), visits=visits)
```

## Ruby

```ruby
docker_service 'default' do
  action [:create, :start]
end
```

## JSON

```json
"server": {
  "http_addr": ":4443",
  "tls_key_file": "./fixtures/notary-server.key",
  "tls_cert_file": "./fixtures/notary-server.crt"
}
```

#### HTML

```html
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
</head>
</html>
```

## Markdown

```markdown
# Hello
```

If you want to include a triple-fenced code block inside your code block,
you can wrap your block in a quadruple-fenced code block:

`````markdown
````markdown
# Hello

```go
log.Println("did something")
```
````
`````

## ini

```ini
[supervisord]
nodaemon=true

[program:sshd]
command=/usr/sbin/sshd -D
```

## Dockerfile

```dockerfile
# syntax=docker/dockerfile:1

FROM ubuntu

RUN apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8

RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list

RUN apt-get update && apt-get install -y python-software-properties software-properties-common postgresql-9.3 postgresql-client-9.3 postgresql-contrib-9.3

# Note: The official Debian and Ubuntu images automatically ``apt-get clean``
# after each ``apt-get``

USER postgres

RUN    /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/9.3/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/9.3/main/postgresql.conf

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

CMD ["/usr/lib/postgresql/9.3/bin/postgres", "-D", "/var/lib/postgresql/9.3/main", "-c", "config_file=/etc/postgresql/9.3/main/postgresql.conf"]
```

## YAML

```yaml
authorizedkeys:
  image: dockercloud/authorizedkeys
  deployment_strategy: every_node
  autodestroy: always
  environment:
    - AUTHORIZED_KEYS=ssh-rsa AAAAB3Nsomelongsshkeystringhereu9UzQbVKy9o00NqXa5jkmZ9Yd0BJBjFmb3WwUR8sJWZVTPFL
  volumes:
    /root:/user:rw
```
