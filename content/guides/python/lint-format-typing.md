---
title: Linting, formatting, and type checking for Python
linkTitle: Linting and typing
weight: 25
keywords: Python, linting, formatting, type checking, ruff, pyright
description: Learn how to set up linting, formatting and type checking for your Python application.
aliases:
  - /language/python/lint-format-typing/
---

## Prerequisites

Complete [Develop your app](develop.md).

## Overview

Linting, formatting, and type checking are automated ways to catch bugs,
enforce style, and spot type errors before code runs. Running them on every
commit, in CI, and in your editor catches problems early when they're cheap
to fix.

In this section, you'll configure three tools for your Python application.
Ruff handles linting and formatting in a single fast pass. Pyright statically
checks your code for type errors. Pre-commit hooks run both of these
automatically before each Git commit so problems never reach your remote
branch.

You can run these tools from your local machine or inside a Docker Hardened
Image container with a bind mount, so no local Python installation is
required.

## Linting and formatting with Ruff

Ruff is an extremely fast Python linter and formatter written in Rust. It replaces multiple tools like flake8, isort, and black with a single unified tool.

Create a `pyproject.toml` file in your `python-docker-example` directory:

{{< files name="python-docker-example" >}}

{{< file path="pyproject.toml" status="new" >}}
```toml
[tool.ruff]
target-version = "py312"

[tool.ruff.lint]
select = [
    "E",  # pycodestyle errors
    "W",  # pycodestyle warnings
    "F",  # pyflakes
    "I",  # isort
    "B",  # flake8-bugbear
    "C4",  # flake8-comprehensions
    "UP",  # pyupgrade
    "ARG001", # unused arguments in functions
]
ignore = [
    "E501",  # line too long, handled by black
    "B008",  # do not perform function calls in argument defaults
    "W191",  # indentation contains tabs
    "B904",  # Allow raising exceptions without from e, for HTTPException
]
```
{{< /file >}}

{{< /files >}}

### Using Ruff

Choose how to run Ruff. The **Local** tab uses a Python installation on your
machine. The **Container** tab runs Ruff inside a Docker Hardened Image
with your project directory bind-mounted in, so no local Python install is
required.

{{< tabs >}}
{{< tab name="Local" >}}

Install Ruff:

```console
$ pip install ruff
```

If you're using a virtual environment, make sure it is activated so the `ruff`
command is available.

Run these commands to check and format your code:

```console
# Check for errors
$ ruff check .

# Automatically fix fixable errors
$ ruff check --fix .

# Format code
$ ruff format .
```

{{< /tab >}}
{{< tab name="Container" >}}

Run Ruff inside the DHI development image with your project directory mounted in:

```console
# Check for errors
$ docker run --rm -v $PWD:/app -w /app dhi.io/python:3.12-dev \
  sh -c "pip install --quiet ruff && ruff check ."

# Automatically fix fixable errors
$ docker run --rm -v $PWD:/app -w /app dhi.io/python:3.12-dev \
  sh -c "pip install --quiet ruff && ruff check --fix ."

# Format code
$ docker run --rm -v $PWD:/app -w /app dhi.io/python:3.12-dev \
  sh -c "pip install --quiet ruff && ruff format ."
```

> [!NOTE]
>
> On Windows, replace `$PWD` with `${PWD}` in PowerShell or `%cd%` in
> Command Prompt.

{{< /tab >}}
{{< /tabs >}}

## Type checking with Pyright

Pyright is a fast static type checker for Python that works well with modern Python features.

Update `pyproject.toml` to add the Pyright configuration at the bottom.

{{< files name="python-docker-example" >}}

{{< file path="pyproject.toml" status="modified" hl_lines="21-25" >}}
```toml
[tool.ruff]
target-version = "py312"

[tool.ruff.lint]
select = [
    "E",  # pycodestyle errors
    "W",  # pycodestyle warnings
    "F",  # pyflakes
    "I",  # isort
    "B",  # flake8-bugbear
    "C4",  # flake8-comprehensions
    "UP",  # pyupgrade
    "ARG001", # unused arguments in functions
]
ignore = [
    "E501",  # line too long, handled by black
    "B008",  # do not perform function calls in argument defaults
    "W191",  # indentation contains tabs
    "B904",  # Allow raising exceptions without from e, for HTTPException
]

[tool.pyright]
typeCheckingMode = "strict"
pythonVersion = "3.12"
exclude = [".venv"]
```
{{< /file >}}

{{< /files >}}

### Running Pyright

As with Ruff, choose how to run Pyright. The **Local** tab uses a Python
installation on your machine. The **Container** tab runs Pyright inside
a Docker Hardened Image with your project directory bind-mounted in, so no
local Python install is required.

{{< tabs >}}
{{< tab name="Local" >}}

Install Pyright and run it:

```console
$ pip install pyright
$ pyright
```

{{< /tab >}}
{{< tab name="Container" >}}

Run Pyright inside the DHI development image. The command also installs your project's
dependencies so Pyright can fully resolve all imports:

```console
$ docker run --rm -v $PWD:/app -w /app dhi.io/python:3.12-dev \
  sh -c "pip install --quiet -r requirements.txt pyright && pyright"
```

> [!NOTE]
>
> On Windows, replace `$PWD` with `${PWD}` in PowerShell or `%cd%` in
> Command Prompt.

{{< /tab >}}
{{< /tabs >}}

## Setting up pre-commit hooks

Pre-commit hooks run checks automatically before each commit on your local
machine. Because the hooks fire when you run `git commit` from your host,
this step requires a local Python installation. If you don't have Python
locally, skip this section. The same checks run in CI in the
[next topic](configure-github-actions.md).

Create a `.pre-commit-config.yaml` file in your `python-docker-example`
directory to set up Ruff hooks:

{{< files name="python-docker-example" >}}

{{< file path=".pre-commit-config.yaml" status="new" >}}
```yaml
repos:
  - repo: https://github.com/charliermarsh/ruff-pre-commit
    rev: v0.2.2
    hooks:
      - id: ruff
        args: [--fix]
      - id: ruff-format
```
{{< /file >}}

{{< /files >}}

To install and use:

```console
$ pip install pre-commit
$ pre-commit install
$ git commit -m "Test commit"  # Automatically runs checks
```

## Summary

In this section, you learned how to:

- Configure and use Ruff for linting and formatting
- Set up Pyright for static type checking
- Automate checks with pre-commit hooks

These tools help maintain code quality and catch errors early in development.

## Next steps

- [Configure GitHub Actions](configure-github-actions.md) to run these checks automatically
- Customize linting rules to match your team's style preferences
- Explore advanced type checking features
