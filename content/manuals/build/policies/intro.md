---
title: Introduction to build policies
linkTitle: Introduction
description: Get started with writing and evaluating build policies
weight: 10
---

Build policies let you validate the inputs to your Docker builds before they
run. This tutorial walks you through creating your first policy, teaching the
Rego basics you need along the way.

## What you'll learn

By the end of this tutorial, you'll understand:

- How to create and organize policy files
- Basic Rego syntax and patterns
- How to write policies that validate URLs, checksums, and images
- How policies evaluate during builds

## Prerequisites

- Buildx version 0.xy or later <!-- TODO: update version -->
- Basic familiarity with Dockerfiles and building images

## How policies work

When you build an image, buildx resolves all the inputs your Dockerfile
references: base images from `FROM` instructions, files from `ADD` or `COPY`,
and Git repositories. Before running the build, buildx evaluates your policies
against these inputs. If any input violates a policy, the build fails before
any instructions execute.

Policies are written in Rego, a declarative language designed for expressing
rules and constraints. You don't need to know Rego to get started - this
tutorial teaches you what you need.

## Create your first policy

Create a new directory for this tutorial and add a Dockerfile:

```console
$ mkdir policy-tutorial
$ cd policy-tutorial
```

Create a `Dockerfile` that downloads a file with `ADD`:

```dockerfile
FROM scratch
ADD https://example.com/index.html /index.html
```

Now create a policy file. Policies use the `.rego` extension and live alongside
your Dockerfile. Create `Dockerfile.rego`:

```rego {title="Dockerfile.rego"}
package docker

default allow := false

allow if input.local
allow if {
  input.http.host == "example.com"
}

decision := {"allow": allow}
```

Save this file as `Dockerfile.rego` in the same directory as your Dockerfile.

Let's break down what this policy does:

- `package docker` - All build policies must start with this package declaration
- `default allow := false` - Start with deny-by-default for security
- `allow if input.local` - The first rule allows any local files (your build context)
- `allow if { input.http.host == "example.com" }` - The second rule allows HTTP downloads from `example.com`
- `decision := {"allow": allow}` - The final decision object tells buildx whether to allow or deny the input

This policy says: "Only allow local files and HTTP downloads from
`example.com`". Any other inputs cause the build to fail.

### About `input.local`

You'll see `allow if input.local` in nearly every policy. This rule allows
local file access, which includes your build context (the `.` directory) and
importantly, the Dockerfile itself. Without this rule, buildx can't read your
Dockerfile to start the build.

Even builds that don't reference any files from the build context need
`input.local` because the Dockerfile is a local file. The policy evaluates
before the build starts, and denying local access means denying access to the
Dockerfile.

In rare cases, you might want stricter local file policies - for example,
checking `input.local.name` to restrict which directories can be used as build
context. But most policies simply allow all local access.

## Run a build with your policy

Build the image with policy evaluation enabled:

```console
$ docker build .
```

The build succeeds because the URL in your Dockerfile matches the policy. Now
try changing the URL in your Dockerfile to something else:

```dockerfile
FROM scratch
ADD https://api.github.com/users/octocat /user.json
```

Build again:

```console
$ docker build .
```

This time the build fails with a policy violation. The `api.github.com`
hostname doesn't match the rule in your policy, so buildx rejects it before
running any build steps.

## Understand Rego rules

Rego policies are built from rules. A rule defines when something is allowed.
The basic pattern is:

```rego
allow if {
    condition_one
    condition_two
    condition_three
}
```

All conditions must be true for the rule to match. Think of it as an AND
operation - the URL must match AND the checksum must match AND any other
conditions you specify.

You can have multiple `allow` rules in one policy. If any rule matches, the
input is allowed:

```rego
# Allow downloads from example.com
allow if {
    input.http.host == "example.com"
}

# Also allow downloads from api.github.com
allow if {
    input.http.host == "api.github.com"
}
```

This works like OR - the input can match the first rule OR the second rule.

## Access input fields

The `input` object gives you access to information about build inputs. The
structure depends on the input type:

- `input.http` - Files downloaded with `ADD https://...`
- `input.image` - Container images from `FROM` or `COPY --from`
- `input.git` - Git repositories from `ADD git://...` or build context
- `input.local` - Local file context

For HTTP downloads, you can access:

| Key                 | Description                        | Example                          |
| ------------------- | ---------------------------------- | -------------------------------- |
| `input.http.url`    | The full URL                       | `https://example.com/index.html` |
| `input.http.schema` | The protocol (HTTP/HTTPS)          | `https`                          |
| `input.http.host`   | The hostname                       | `example.com`                    |
| `input.http.path`   | The URL path, including parameters | `/index.html`                    |

Update your policy to require HTTPS:

```rego
package docker

default allow := false

allow if {
    input.http.host == "example.com"
    input.http.schema == "https"
}

decision := {"allow": allow}
```

Now the policy requires both the hostname to be `example.com` and the protocol
to be HTTPS. HTTP URLs (without TLS) would fail the policy check.

## Pattern matching and strings

Rego provides [built-in functions] for pattern matching. Use `startswith()` to
match URL prefixes:

[built-in functions]: https://www.openpolicyagent.org/docs/policy-language#built-in-functions

```rego
allow if {
    startswith(input.http.url, "https://example.com/")
}
```

This allows any URL that starts with `https://example.com/`.

Use `regex.match()` for complex patterns:

```rego
allow if {
    regex.match(`^https://example\.com/.+\.json$`, input.http.url)
}
```

This matches URLs that:

- Start with `https://example.com/`
- End with `.json`
- Have at least one character between the domain and extension

## Policy file location

Policy files live adjacent to the Dockerfile they validate, using the naming
pattern `<dockerfile-name>.rego`:

```text
project/
├── Dockerfile           # Main Dockerfile
├── Dockerfile.rego      # Policy for Dockerfile
├── lint.Dockerfile      # Linting Dockerfile
└── lint.Dockerfile.rego # Policy for lint.Dockerfile
```

When you build, buildx automatically loads the corresponding policy file:

```console
$ docker buildx build -f Dockerfile .        # Loads Dockerfile.rego
$ docker buildx build -f lint.Dockerfile .   # Loads lint.Dockerfile.rego
```

## Next steps

You now understand how to write basic build policies for HTTP resources. To
continue learning:

- Learn [Image validation](./validate-images.md) to validate container images
  from `FROM` instructions
- Learn [Git validation](./validate-git.md) to validate Git repositories used
  in builds
- See [Example policies](./examples.md) for copy-paste-ready policies covering
  common scenarios
- Read the [Input reference](./inputs.md) for all available input fields
- Check the [Built-in functions](./built-ins.md) for signature verification,
  attestations, and other security checks
