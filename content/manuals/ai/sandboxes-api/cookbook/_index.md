---
title: Docker Sandboxes SDK cookbook
linkTitle: Cookbook
description: Build agent workflows with the Docker Sandboxes SDK for TypeScript.
keywords: cloud sandboxes, sandboxes sdk, typescript, agent kits
weight: 60
params:
  sidebar:
    groups:
      - Get started
      - Working in a sandbox
      - Packaging and state
      - Security and policy
      - Connect and configure
      - Requests and responses
      - Long-running work
---

Use these recipes to add file transfers, processes, storage, and other sandbox
operations to your application. Each recipe shows the relevant SDK calls and
an expandable complete TypeScript example.

Start with [Get started](../get-started.md) to run your first sandbox, or
[install the SDK](../sdks.md) to use these examples in an existing project.
For HTTP operations and request fields, see the
[API reference](/reference/api/sandboxes/index.md).

## Get started

Authenticate and launch a kit before exploring individual SDK operations.

- [Authenticate to Docker](connect-to-cloud-with-a-bearer-token.md): Choose interactive sign-in, a PAT, or your own token provider.
- [Run a complete example](run-a-complete-example.md): Sign in, print a greeting from a sandbox, and clean up.
- [Create your first sandbox](create-your-first-sandbox.md): Launch a kit, run a command, and choose a cleanup pattern.
- [Run agents with kits](add-tools-with-kits.md): Choose a bundled kit and run an agent with its provider credential.
- [Run your first command](run-your-first-command.md): Get a sandbox handle and collect a command's result.
- [Delete a cloud sandbox](delete-a-cloud-sandbox.md): Delete a sandbox and confirm cleanup.
- [Name a sandbox and find it again](name-a-sandbox-and-find-it-again.md): Save a resource name and read the sandbox again.

## Working in a sandbox

Use a running sandbox for commands, project files, and web applications.

- [Run an interactive shell in a cloud sandbox](run-an-interactive-shell-in-a-cloud-sandbox.md): Send terminal input and reconnect to output.
- [Copy a file into a cloud sandbox](copy-a-file-into-a-cloud-sandbox.md): Create directories and send text or binary files.
- [Read files out of a sandbox](read-files-out-of-a-sandbox.md): Read, download, move, and remove files.
- [Expose a port from a cloud sandbox](expose-a-port-from-a-cloud-sandbox.md): Expose an HTTP application and withdraw access.
- [Stop and restart a sandbox](stop-and-restart-a-sandbox.md): Stop execution while keeping the sandbox's disk.
- [What your workload starts with](what-your-workload-starts-with.md): Inspect environment settings and override one command.
- [Run something that produces real output](run-something-that-produces-real-output.md): Choose captured results or streamed progress.

## Packaging and state

Keep an environment or its state for later work.

- [Register and manage an image](register-and-manage-an-image.md): Register image content for reuse.
- [Snapshot and fork a sandbox](snapshot-and-fork-a-sandbox.md): Save state and restore it into another sandbox.

## Security and policy

Give agents the access they need without embedding credentials in application code.

- [Control what a sandbox can reach](control-what-a-sandbox-can-reach.md): Attach policies and inspect allowed or blocked traffic.
- [Manage cloud secrets](manage-cloud-secrets.md): Store, rotate, and delete workload credentials.
- [Give a sandbox an MCP gateway](give-a-sandbox-an-mcp-gateway.md): Configure external tools for a kit and manage its gateway.
- [Add a Docker credential for cloud sandboxes](add-a-docker-credential-for-cloud-sandboxes.md): Store a Docker credential for workload use.
- [Get a stored secret into a sandbox](get-a-stored-secret-into-a-sandbox.md): Attach a stored provider credential at sandbox creation.
- [Manage network policies](manage-network-policies.md): Manage reusable personal policy definitions.

## Connect and configure

Choose a different image, control lifetime, or connect additional storage and clients.

- [Run your own container image](run-your-own-container-image.md): Create a sandbox from your own registry image.
- [Keep a cloud sandbox running](keep-a-cloud-sandbox-running.md): Choose a lifetime and renew its remaining duration.
- [Attach persistent storage](attach-persistent-storage.md): Keep data independently of a sandbox.
- [Get an SSH certificate](get-an-ssh-certificate.md): Obtain a short-lived certificate for your public key.
- [Get image pull URLs](get-image-pull-urls.md): Read image references for an OCI client.
- [Let a stopped sandbox resume on demand](let-a-stopped-sandbox-resume-on-demand.md): Configure request-triggered startup.

## Requests and responses

Handle request options and failures deliberately.

- [Handle errors and degradation](handle-errors-and-degradation.md): Separate request errors, wait failures, and command results.
- [Retry without creating duplicates](retry-without-creating-duplicates.md): Keep retries from duplicating work.
- [Page through and filter lists](page-through-and-filter-lists.md): Iterate complete collections and filter results.
- [Send values the API accepts](send-values-the-api-accepts.md): Preserve omitted settings and use SDK-native values.
- [Work within the limits](work-within-the-limits.md): Bound concurrent work and honor retry delays.

## Long-running work

Continue work across lost connections or coordinate several sandboxes.

- [Find a process you lost track of](find-a-process-you-lost-track-of.md): Find an existing process instead of starting another.
- [Clone a cloud sandbox](clone-a-cloud-sandbox.md): Copy selected configuration into a fresh sandbox.
- [Run work across many sandboxes](run-work-across-many-sandboxes.md): Run work across selected sandboxes with bounded concurrency.
- [Recover when the endpoint moves](recover-when-the-endpoint-moves.md): Refresh a handle before opening a new connection.
