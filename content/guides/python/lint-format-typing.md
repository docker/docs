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

In this section, you'll learn how to set up code quality tools for your Python application. This includes:

- Linting and formatting with Ruff
- Static type checking with Pyright
- Automating checks with pre-commit hooks

## Linting and formatting with Ruff

Ruff is an extremely fast Python linter and formatter written in Rust. It replaces multiple tools like flake8, isort, and black with a single unified tool.

Create a `pyproject.toml` file:

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

### Using Ruff

Run these commands to check and format your code:

```bash
# Check for errors
ruff check .

# Automatically fix fixable errors
ruff check --fix .

# Format code
ruff format .
```

## Type checking with Pyright

Pyright is a fast static type checker for Python that works well with modern Python features.

Add `Pyright` configuration in `pyproject.toml`:

```toml
[tool.pyright]
typeCheckingMode = "strict"
pythonVersion = "3.12"
exclude = [".venv"]
```

### Running Pyright

To check your code for type errors:

```bash
pyright
```

## Setting up pre-commit hooks

Pre-commit hooks automatically run checks before each commit. The following `.pre-commit-config.yaml` snippet sets up Ruff:

```yaml
  https: https://github.com/charliermarsh/ruff-pre-commit
  rev: v0.2.2
  hooks:
    - id: ruff
      args: [--fix]
    - id: ruff-format
```

To install and use:

```bash
pre-commit install
git commit -m "Test commit"  # Automatically runs checks
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