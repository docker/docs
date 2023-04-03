---
description: File watch automatically updates running services as you work
keywords: compose, file watch, experimental 
title: Automatically update services with file watch
---

> **Note**
>
> The Compose file watch feature is currently [Experimental](../release-lifecycle.md).

Use `watch` to automatically update your running Compose services as you edit and save your code. 

`watch` adheres to the following file path rules:

* All paths are relative to the build context
* Directories are watched recursively
* Glob patterns are not supported
* Rules from `.dockerignore` apply
  * Use `include` / `exclude` to override
  * Temporary/backup files for common IDEs (Vim, Emacs, JetBrains, & more) are ignored automatically
* `.git` directories are ignored automatically

## Configuration

The `watch` attribute defines a list of rules that control automatic service updates based on local file changes.

Each rule requires, a `path` pattern and `action` to take when a modification is detected. There are two possible actions for `watch` and depending on
the `action`, additional fields might be accepted or required. 

### `action`

#### Sync

If `action` is set to `sync`, Compose makes sure any changes made to files on your host automatically match with the corresponding files within the service container.

Sync is ideal for frameworks that support "Hot Reload" or equivalent functionality.

More generally, they can be used in place of bind mounts for many development use cases.

#### Rebuild

If `action` is set to `rebuild`, Compose automatically builds a new image with BuildKit and replaces the running service container.

The behavior is the same as running `docker compose up --build <svc>`.

Rebuild is ideal for compiled languages or as fallbacks for modifications to particular files that require a full
image rebuild (e.g. `package.json`).

>**Tip**
>
> Optimize your `Dockerfile` for speedy
incremental rebuilds with [image layer caching](/build/cache)
and [multi-stage builds](/build/building/multi-stage/).
{: .tip}

### `path` and `target`

The `target` field controls how the path is mapped into the container.

For `path: ./app/html` and a change to `./app/html/index.html`:

* `target: /app/html` -> `/app/html/index.html`
* `target: /app/static` -> `/app/static/index.html`
* `target: /assets` -> `/assets/index.html`

## Example

```yaml
services:
  app:
    build: ./app
    x-develop:
      watch:
        - action: sync
          path: ./app/html
          target: /app/html
        - action: rebuild
          path: ./app/requirements.txt
```

## Use `watch`

1. Add `watch` sections to one or more services in `compose.yaml`.
2. Launch a Compose project with `docker compose up --build --wait`.
3. Run `docker compose alpha watch` to start the file watch mode.
4. Edit service source files using your preferred IDE or editor.

>**Tip**
>
> Looking for a sample project to test things out? Check
out [`dockersamples/avatars`](https://github.com/dockersamples/avatars) for a demonstration of Compose `watch`.
{: .tip}

## Feedback

We are actively looking for feedback on this feature. Give feedback or report any bugs you may find in the [Compose Specification repository](https://github.com/compose-spec/compose-spec/pull/253).