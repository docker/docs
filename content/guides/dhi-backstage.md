---
title: Secure a Backstage application with Docker Hardened Images
description: Secure a Backstage developer portal using Docker Hardened Images, covering native module compilation, Socket Firewall protection, and distroless runtime images.
summary: Learn how to secure a Backstage developer portal using Docker Hardened Images (DHI), handle native module compilation with better-sqlite3, add Socket Firewall protection during dependency installation, and produce a distroless runtime image using DHI customizations.
keywords: docker hardened images, dhi, backstage, CNCF, developer portal, node.js, native modules, sqlite, better-sqlite3, distroless, socket firewall, dhictl, multi-stage build
tags: ["Docker Hardened Images", "dhi"]
params:
  proficiencyLevel: Intermediate
  time: 45 minutes
  prerequisites:
    - Docker Desktop or Docker Engine with BuildKit enabled
    - A Docker Hub account authenticated with docker login and docker login dhi.io
    - A Backstage project created with @backstage/create-app
    - Basic familiarity with multi-stage Dockerfiles and Node.js native modules
---

This guide shows how to secure a Backstage application using Docker Hardened Images (DHI). Backstage is a CNCF open source developer portal used by thousands of organizations to manage their software catalogs, templates, and developer tooling.

By the end of this guide, you'll have a Backstage container image that is distroless, runs as a non-root user by default, and has dramatically fewer CVEs than the standard `node:24-trixie-slim` base image while still supporting the native module compilation that Backstage requires.

## Prerequisites

- Docker Desktop or Docker Engine with BuildKit enabled
- A Docker Hub account authenticated with `docker login` and `docker login dhi.io`
- A Backstage project created with `@backstage/create-app`

## Why Backstage needs customization

The DHI migration examples cover applications where you can swap the base image and everything works. Backstage is different. It uses `better-sqlite3` and other packages that compile native Node.js modules at install time, which means the build stage needs `g++`, `make`, `python3`, and `sqlite-dev` — none of which are in the base `dhi.io/node` image. The runtime image only needs the shared library (`sqlite-libs`) that the compiled native module links against.

This is a common pattern. Any Node.js application that depends on native addons (such as `bcrypt`, `sharp`, `sqlite3`, or `node-canvas`) faces the same challenge. The approach in this guide applies to all of them.

## Step 1: Examine the original Dockerfile

The official Backstage documentation recommends a multi-stage Dockerfile using `node:24-trixie-slim` (Debian). A typical setup looks like this:

```dockerfile
# Stage 1 - Create yarn install skeleton layer
FROM node:24-trixie-slim AS packages
WORKDIR /app
COPY backstage.json package.json yarn.lock ./
COPY .yarn ./.yarn
COPY .yarnrc.yml ./
COPY packages packages
COPY plugins plugins
RUN find packages \! -name "package.json" -mindepth 2 -maxdepth 2 \
    -exec rm -rf {} \+

# Stage 2 - Install dependencies and build packages
FROM node:24-trixie-slim AS build
ENV PYTHON=/usr/bin/python3
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends python3 g++ build-essential && \
    rm -rf /var/lib/apt/lists/*
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*
USER node
WORKDIR /app
COPY --from=packages --chown=node:node /app .
RUN --mount=type=cache,target=/home/node/.cache/yarn,sharing=locked,uid=1000,gid=1000 \
    yarn install --immutable
COPY --chown=node:node . .
RUN yarn tsc
RUN yarn --cwd packages/backend build
RUN mkdir packages/backend/dist/skeleton packages/backend/dist/bundle \
    && tar xzf packages/backend/dist/skeleton.tar.gz \
       -C packages/backend/dist/skeleton \
    && tar xzf packages/backend/dist/bundle.tar.gz \
       -C packages/backend/dist/bundle

# Stage 3 - Build the actual backend image
FROM node:24-trixie-slim
ENV PYTHON=/usr/bin/python3
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends python3 g++ build-essential && \
    rm -rf /var/lib/apt/lists/*
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*
USER node
WORKDIR /app
COPY --from=build --chown=node:node /app/.yarn ./.yarn
COPY --from=build --chown=node:node /app/.yarnrc.yml ./
COPY --from=build --chown=node:node /app/backstage.json ./
COPY --from=build --chown=node:node /app/yarn.lock \
     /app/package.json \
     /app/packages/backend/dist/skeleton/ ./
RUN --mount=type=cache,target=/home/node/.cache/yarn,sharing=locked,uid=1000,gid=1000 \
    yarn workspaces focus --all --production
COPY --from=build --chown=node:node /app/packages/backend/dist/bundle/ ./
CMD ["node", "packages/backend", "--config", "app-config.yaml"]
```

Run this image and inspect what's available inside the container:

```console
docker build -t backstage:init .
docker run -d \
    -e APP_CONFIG_backend_database_client='better-sqlite3' \
    -e APP_CONFIG_backend_database_connection=':memory:' \
    -e APP_CONFIG_auth_providers_guest_dangerouslyAllowOutsideDevelopment='true' \
    -p 7007:7007 \
    -u 1000 \
    --cap-drop=ALL \
    --read-only \
    --tmpfs /tmp \
    backstage:init
```

This works, but the runtime container has a shell, a package manager, and yarn. None of these are needed to run Backstage. Run `docker exec` to see what's accessible inside:

```console
docker exec -it <container-id> sh
$ cat /etc/shells
# /etc/shells: valid login shells
/bin/sh
/usr/bin/sh
/bin/bash
/usr/bin/bash
/bin/rbash
/usr/bin/rbash
/usr/bin/dash
$ yarn --version
4.12.0
$ dpkg --version
dpkg version 1.22.11 (arm64).
$ whoami
node
$ id
uid=1000(node) gid=1000(node) groups=1000(node)
```

The `node:24-trixie-slim` image ships with three shells (`dash`, `bash`, and `rbash`), a package manager (`dpkg`), and `yarn`. Each of these tools increases the attack surface. An attacker who gains access to this container could use them for lateral movement across your infrastructure.

## Step 2: Switch the build stages to DHI

Replace all three stages with DHI equivalents. DHI Node.js images are available in both 
Alpine and Debian variants. This guide uses the Alpine variant (`dhi.io/node:24-alpine3.23`) 
because it produces a smaller image. If you need to stay on Debian for compatibility reasons, 
use `dhi.io/node:24-bookworm` and keep `apt-get` instead of `apk`.

```dockerfile
# Stage 1: prepare packages
FROM --platform=$BUILDPLATFORM dhi.io/node:24-alpine3.23-dev AS packages
WORKDIR /app
COPY backstage.json package.json yarn.lock ./
COPY .yarn ./.yarn
COPY .yarnrc.yml ./
COPY packages packages
COPY plugins plugins
RUN find packages \! -name "package.json" -mindepth 2 -maxdepth 2 \
    -exec rm -rf {} \+

# Stage 2: build the application
FROM --platform=$BUILDPLATFORM dhi.io/node:24-alpine3.23-dev AS build
ENV PYTHON=/usr/bin/python3
RUN apk add --no-cache g++ make python3 sqlite-dev && \
    rm -rf /var/lib/apk/lists/*
WORKDIR /app
COPY --from=packages --chown=node:node /app .
RUN --mount=type=cache,target=/home/node/.cache/yarn,sharing=locked,uid=1000,gid=1000 \
    yarn install --immutable
COPY --chown=node:node . .
RUN yarn tsc
RUN yarn --cwd packages/backend build
RUN mkdir packages/backend/dist/skeleton packages/backend/dist/bundle \
    && tar xzf packages/backend/dist/skeleton.tar.gz \
       -C packages/backend/dist/skeleton \
    && tar xzf packages/backend/dist/bundle.tar.gz \
       -C packages/backend/dist/bundle

# Final Stage: create the runtime image
FROM dhi.io/node:24-alpine3.23-dev
ENV PYTHON=/usr/bin/python3
RUN apk add --no-cache g++ make python3 sqlite-dev && \
    rm -rf /var/lib/apk/lists/*
WORKDIR /app
COPY --from=build --chown=node:node /app/.yarn ./.yarn
COPY --from=build --chown=node:node /app/.yarnrc.yml ./
COPY --from=build --chown=node:node /app/backstage.json ./
COPY --from=build --chown=node:node /app/yarn.lock \
     /app/package.json \
     /app/packages/backend/dist/skeleton/ ./
RUN --mount=type=cache,target=/home/node/.cache/yarn,sharing=locked,uid=1000,gid=1000 \
    yarn workspaces focus --all --production \
    && rm -rf "$(yarn cache clean)"
COPY --from=build --chown=node:node /app/packages/backend/dist/bundle/ ./
CMD ["node", "packages/backend", "--config", "app-config.yaml"]
```

Build and tag this version:

```console
docker build -t backstage:dhi-dev .
```

> [!NOTE]
>
> The `-dev` variant includes a shell and package manager, which is why `apk add` works. Backstage requires `python3` and native build tools in the runtime image because `yarn workspaces focus --all --production` recompiles native modules during the production install. This is specific to Backstage's build process — most Node.js applications can use the standard (non-dev) DHI runtime variant without additional packages.

The DHI images come with attestations that the original `node:24-trixie-slim` images don't have. Check what's attached:

```console
docker scout attest list dhi.io/node:24-alpine3.23
```

DHI images ship with 15 attestations including CycloneDX SBOM, SLSA provenance, OpenVEX, Scout health reports, secret scans, virus/malware reports, and an SLSA verification summary.

## Step 3: Add Socket Firewall protection

DHI provides `-sfw` (Socket Firewall) variants for Node.js images. Socket Firewall intercepts `npm` and `yarn` commands during the build to detect and block malicious packages before they execute install scripts.

To enable Socket Firewall, change the `-dev` tags to `-sfw-dev` in all three stages. The SFW version of the Dockerfile:

```dockerfile
# Stage 1: prepare packages
FROM --platform=$BUILDPLATFORM dhi.io/node:24-alpine3.23-sfw-dev AS packages
WORKDIR /app
COPY backstage.json package.json yarn.lock ./
COPY .yarn ./.yarn
COPY .yarnrc.yml ./
COPY packages packages
COPY plugins plugins
RUN find packages \! -name "package.json" -mindepth 2 -maxdepth 2 \
    -exec rm -rf {} \+

# Stage 2: build the packages
FROM --platform=$BUILDPLATFORM dhi.io/node:24-alpine3.23-sfw-dev AS build-packages
ENV PYTHON=/usr/bin/python3
RUN apk add --no-cache g++ make python3 sqlite-dev && \
    rm -rf /var/lib/apk/lists/*
WORKDIR /app
COPY --from=packages --chown=node:node /app .
RUN --mount=type=cache,target=/home/node/.cache/yarn,sharing=locked,uid=1000,gid=1000 \
    yarn install --immutable
COPY --chown=node:node . .
RUN yarn tsc
RUN yarn --cwd packages/backend build
RUN mkdir packages/backend/dist/skeleton packages/backend/dist/bundle \
    && tar xzf packages/backend/dist/skeleton.tar.gz \
       -C packages/backend/dist/skeleton \
    && tar xzf packages/backend/dist/bundle.tar.gz \
       -C packages/backend/dist/bundle

# Final Stage: create the runtime image
FROM dhi.io/node:24-alpine3.23-sfw-dev
ENV PYTHON=/usr/bin/python3
RUN apk add --no-cache g++ make python3 sqlite-dev && \
    rm -rf /var/lib/apk/lists/*
WORKDIR /app
COPY --from=build-packages --chown=node:node /app/.yarn ./.yarn
COPY --from=build-packages --chown=node:node /app/.yarnrc.yml ./
COPY --from=build-packages --chown=node:node /app/backstage.json ./
COPY --from=build-packages --chown=node:node /app/yarn.lock \
     /app/package.json \
     /app/packages/backend/dist/skeleton/ ./
RUN --mount=type=cache,target=/home/node/.cache/yarn,sharing=locked,uid=1000,gid=1000 \
    yarn workspaces focus --all --production \
    && rm -rf "$(yarn cache clean)"
COPY --from=build-packages --chown=node:node /app/packages/backend/dist/bundle/ ./
CMD ["node", "packages/backend", "--config", "app-config.yaml"]
```

Build this version:

```console
docker build -t backstage:dhi-sfw-dev .
```

When you build, you'll see Socket Firewall messages in the build output: `Protected by Socket Firewall` for any `yarn` and `npm` commands executed in the Dockerfile or in the running containers.

> [!TIP]
>
> The `-sfw-dev` variant is larger (1.9 GB versus 1.72 GB) because Socket Firewall adds monitoring tooling. The security benefit during `yarn install` outweighs the size increase.

## Step 4: Remove the shell and the package manager with DHI customizations

The previous steps still use the `-dev` or `-sfw-dev` variant as the runtime image, which includes a shell and package manager. DHI customizations let you start from the base (non-dev) image — which has no shell and no package manager — and add only the runtime libraries and language runtimes your application needs.

> [!IMPORTANT]
>
> When creating a customization, only add what your application needs at runtime:
>
> - **System packages** - add shared libraries (such as `sqlite-libs`) and
>   language runtimes from the DHI catalog (such as `python-3.14`).
>   Do not add build tools (such as `g++`, `make`, or `python3` from Alpine).
> - **Build tools** - keep these in the `-dev` build stage only. Never add them
>   to the runtime customization.
>
> Language runtimes installed from the DHI hardened package feed are patched and
> tracked in the image SBOM, which is why they are acceptable as system packages.
> Build tools from Alpine or Debian package feeds are not hardened and should
> never appear in the runtime image.

For Backstage, the runtime image needs:

- **sqlite-libs** - the shared library that the compiled `better-sqlite3` native module links against (added as a system package).
- **Python** - if your Backstage plugins or configuration require Python at runtime. Added as the `python-3.14` system package from the DHI catalog. Unlike `python3` installed via `apk`, this package is patched by Docker and tracked in the image SBOM.

Docker will continuously build with SLSA Level 3 compliance and patch these customized images within the guaranteed SLA for CVE patching.

To create the customization, use one of the following methods.

{{< tabs >}}
{{< tab name="Docker Hub UI" >}}

After you mirror the Node.js DHI repository to your organization's namespace:

1. Open the mirrored Node.js repository in Docker Hub.
2. Select **Customize** and choose the `node:24-alpine3.23` tag.
3. Under **Packages**, add `sqlite-libs` and `python-3.14`.
4. Create the customization.

For more information, see [Customize an image](/dhi/how-to/customize/).

{{< /tab >}}
{{< tab name="dhictl CLI" >}}

`dhictl` is Docker's command-line tool for managing Docker Hardened Images. It lets you browse the DHI catalog, mirror images, and create customizations directly from your terminal. You can integrate `dhictl` into CI/CD pipelines and infrastructure-as-code workflows. You can install `dhictl` as a standalone binary or as a Docker CLI plugin (`docker dhi`); for installation instructions, see [Use the DHI CLI](/dhi/how-to/cli/).

Rather than writing the customization YAML by hand, use `dhictl` to scaffold a starting point:

```console
dhictl customization prepare --org YOUR_ORG node 24-alpine3.23 \
    --destination YOUR_ORG/dhi-node \
    --name "backstage" \
    --tag-suffix "_backstage" \
    --output node-backstage.yaml
```

Edit the generated file to add the runtime libraries:

```yaml
name: backstage

source: dhi/node
tag_definition_id: node/alpine-3.23/24

destination: YOUR_ORG/dhi-node
tag_suffix: _backstage

platforms:
  - linux/amd64
  - linux/arm64

contents:
  packages:
    - sqlite-libs
    - python-3.14

accounts:
  root: true
  runs-as: node
  users:
    - name: node
      uid: 1000
  groups:
    - name: node
      gid: 1000

```

Then create the customization:

```console
dhictl customization create --org YOUR_ORG node-backstage.yaml
```

Monitor the build progress:

```console
dhictl customization build list --org YOUR_ORG YOUR_ORG/dhi-node "backstage"
```

Docker builds the customized image on its secure infrastructure and publishes it as `YOUR_ORG/dhi-node:24-alpine3.23_backstage`.

> [!NOTE]
>
> If your Backstage configuration does not require Python at runtime, you can omit the `python-3.14` from the packages list. The `sqlite-libs` package alone is sufficient to run Backstage with `better-sqlite3`.

{{< /tab >}}
{{< /tabs >}}

### Update the Dockerfile

Update only the final stage of your Dockerfile to use the customized image:

```dockerfile
# Final Stage: create the runtime image
FROM YOUR_ORG/dhi-node:24-alpine3.23_backstage
WORKDIR /app
COPY --from=build --chown=node:node /app/node_modules ./node_modules
COPY --from=build --chown=node:node /app/packages/backend/dist/bundle/ ./
CMD ["node", "packages/backend", "--config", "app-config.yaml"]
```

Build this version:

```console
docker build -t backstage:dhi .
```

Since the customization includes only runtime libraries and OCI artifacts — no build tools, no package manager, no shell — the resulting image is distroless:

```console
docker run --rm YOUR_ORG/dhi-node:24-alpine3.23_backstage sh -c "echo hello"
docker: Error response from daemon: ... exec: "sh": executable file not found in $PATH
```

With the Enterprise customization:

- The runtime image is distroless — no shell, no package manager.
- Docker automatically rebuilds your customized image when the base Node.js image or any of its packages receive a security patch.
- The full chain of trust is maintained, including SLSA Build Level 3 provenance.
- Both the Node.js and Python runtimes are tracked in the image SBOM.

Confirm the container no longer has shell access:

```console
docker exec -it <container-id> sh
OCI runtime exec failed: exec failed: unable to start container process: ...
```

Use [Docker Debug](/dhi/how-to/debug/) if you need to troubleshoot a running distroless container.

> [!NOTE]
>
> If your organization requires FIPS/STIG compliant images, that's also an option in DHI Enterprise.

## Step 5: Verify the results

Compare the DHI-based image against the original using Docker Scout:

```console
docker scout compare backstage:dhi \
    --to backstage:init \
    --platform linux/amd64 \
    --ignore-unchanged
```

A typical comparison across the approaches shows results similar to the following:

| Metric | Original | DHI -dev | DHI -sfw-dev | Enterprise |
|--------|----------|----------|--------------|------------|
| Disk usage | 1.61 GB | 1.72 GB | 1.9 GB | 1.49 GB |
| Content size | 268 MB | 288 MB | 328 MB | 247 MB |
| Shell in runtime | Yes | Yes | Yes | No |
| Package manager | Yes | Yes | Yes | No |
| Non-root default | No | No | No | Yes |
| Socket Firewall | No | No | Yes (build) | Yes (build) / No (runtime) |
| SLSA provenance | No | Base only | Base only | Full (Level 3) |

> [!NOTE]
>
> The `-sfw-dev` variant is larger because Socket Firewall adds monitoring tooling to the image. The additional size is in the build stages, and the security benefit during `yarn install` outweighs the size increase.

For a more thorough assessment, scan with multiple tools:

```console
trivy image backstage:dhi
grype backstage:dhi
docker scout quickview backstage:dhi
```

Different scanners detect different issues. Running all three gives you the most complete view of your security posture.

## What's next

- [Customize an image](/dhi/how-to/customize/) — complete reference on the Enterprise customization UI.
- [Create and build a DHI](/dhi/how-to/build/) — learn how to write a DHI definition file, build images locally.
- [Use the DHI CLI](/dhi/how-to/cli/) — manage DHI images, mirrors, and customizations from the command line.
- [Migrate to DHI](/dhi/migration/) — for applications that work with standard DHI images without additional packages.
- [Compare images](/dhi/how-to/compare/) — evaluate security improvements between your original and hardened images.
- [Docker Debug](/dhi/how-to/debug/) — troubleshoot distroless containers that have no shell.
