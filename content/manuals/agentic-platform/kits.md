---
title: Kits
description: Browse curated and community kits or launch a public sandbox kit in Docker Agentic Platform.
keywords: docker agentic platform, kits, community kits, custom kits, docker hub, sandboxes
weight: 15
aliases:
  - /agentic-platform/concepts/kits/
  - /agentic-platform/guides/manage-kits/
---

A kit defines a sandbox's base image, agent, setup, network rules, and the
credentials it can use. Choose a kit from the **Kits** catalog or enter a public
kit reference to run an agent in Docker Agentic Platform.

Launching a kit from the Console requires both the kit and its base image to
be public. To run a kit from a private Docker Hub repository, use the Docker
Sandboxes CLI. The **Kits** page provides the command to copy.

## Browse and run a kit

The catalog includes kits curated by Docker and community kits published on
Docker Hub. Curated kits include agents such as Claude Code, Codex, Antigravity,
and Hermes, as well as a Shell kit for working without a pre-installed agent.

1. Open **Kits** in the Console.
2. Search by name or description. Use **Curated kits** or **Community kits** to
   filter the catalog, or **All kits** to see both.
3. Select **Run** on a kit to open the sandbox launcher with that kit selected.
4. Review its credentials, network policies, tools, compute size, and
   timer, then select **Run** in the launcher to create the sandbox.

Use a kit's Hub link to open its repository on Docker Hub. For the launch steps,
see [Get started](get-started.md#start-a-sandbox).

## Run a kit by reference

To use a public kit that you already know:

1. Open **New** and open the kit picker.
2. Select **add a public kit**.
3. Enter the kit's registry reference, such as `myorg/my-kit:1.0`. You can omit
   the `docker.io/` prefix for Docker Hub references.
4. Wait for the kit's name and supported credentials to appear. If you see an
   error, check the reference and confirm that the kit is public.
5. Configure the sandbox and select **Run**.

Enter a reference to a sandbox kit. You can't launch a sandbox from a container
image reference or a mixin kit. Mixins add capabilities to other kits, and the
Console doesn't support combining kits.

The credentials available in the launcher depend on the kit. A custom kit must
declare the credentials it needs, including GitHub credentials for private
repository access. See [Secrets](secrets.md).

## Kit size limits

Custom kits can contain up to 64 MiB of registry content, measured as the total
size of the config blob and compressed layers declared in the kit manifest.
The kit's file payload must also fit within 256 MiB when loaded for launch.
These limits apply to the kit artifact, not its referenced base image.

If the launcher reports "That kit is larger than this platform accepts", the
reference exceeds the 64 MiB limit. Check that the reference points to a
sandbox kit rather than a container image. Put large dependencies in the
base image and keep the kit's bundled files small.

## Create your own kit

To ask your coding agent to help create a kit, copy the prompt from **Kits** and
paste it into your agent. Follow the **Kit authoring docs** link for instructions,
or see [Kit authoring](/manuals/ai/sandboxes/customize/author/_index.md).

After publishing the kit to a public registry, enter its reference in the
launcher. Its base image must also be public.
