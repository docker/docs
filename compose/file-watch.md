---
description: File watch automatically updates running services as you work
keywords: compose, file watch, experimental 
title: Automatically update services during development with Compose file watch
---

Watch mode monitors files and directories related to Compose services.

On change, multiple actions are available:

* **Sync**: create, overwrite, and delete corresponding files within the service container
* **Rebuild**: build a new image with BuildKit and replace the running service container

### Example `watch` configuration

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

## Feature stability
This feature is still [**experimental**](/release-lifecycle/).

In Compose YAML, any fields that start with `x-` do not have the same compatibility guarantees.

We are actively looking for feedback on the `x-develop` section schema and may make breaking changes if necessary.

## File Paths

* All paths are relative to the build context
* Directories are watched recursively
* Glob patterns are not supported
* Rules from `.dockerignore` apply
  * Use `include` / `exclude` to override
  * Temporary/backup files for common IDEs (Vim, Emacs, JetBrains, & more) are ignored automatically
  * `.git` directories are ignored automatically

## Usage

> 👀 Looking for a sample project to test things out? Check
out [`dockersamples/avatars`](https://github.com/dockersamples/avatars) for a whimsical demonstration of Compose file
watch.

1. Add `watch` sections to one or more services in `compose.yaml`
2. Launch Compose project with `docker compose up --build --wait`
3. Run `docker compose alpha watch` to start file watch mode
4. Start hacking! Edit service source files using your preferred IDE or editor
5. Compose automatically syncs files and rebuilds containers as needed

## Configuration
The `watch` field defines a list of rules that control automatic service updates based on local file changes.

For each rule, a `path` pattern and `action` to take when a modification is detected are required. Depending on
the `action`, additional fields might be accepted or required.

### Syncing files
Use `action: sync` to configure Compose to automatically copy files from your host into a running service on file
modification.

Sync rules are ideal for frameworks that support "Hot Reload" or equivalent functionality.

More generally, they can be used in place of bind mounts for many development use cases.

#### Path mapping
The `target` field controls how the path is mapped into the container.

For `path: ./app/html` and a change to `./app/html/index.html`:

* `target: /app/html` -> `/app/html/index.html`
* `target: /app/static` -> `/app/static/index.html`
* `target: /assets` -> `/assets/index.html`

### Rebuilding services
Use `action: rebuild` to configure Compose to automatically trigger a new image build and replace a running service on
file modification.

The behavior is equivalent to running `docker compose up --build <svc>` (but automatic).

> Optimize your `Dockerfile` for speedy
incremental rebuilds with [image layer caching](/build/cache)
and [multi-stage builds](/build/building/multi-stage/).

Rebuild rules are ideal for compiled languages or as fallbacks for modifications to particular files that require a full
image rebuild (e.g. `package.json`).
