---
title: Share an application from a cloud sandbox
linkTitle: Share an app
description: Build a public sample in a cloud sandbox, share its HTTPS URL, and close the port when you finish.
keywords: docker agentic platform, cloud sandboxes, workflows, share an app
weight: 10
---

Run the Welcome to Docker application in a cloud sandbox and share a preview
without setting up the project on your computer. Provisional duration: 20 minutes;
this estimate needs measurement during review.

## Before you begin

You need a browser and an active
[Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access).
This public sample needs no GitHub token, model API key, or MCP tools.
Use non-sensitive data: the preview URL is public, and Docker account
authentication doesn't protect it. Add authentication to your application
before sharing private content or operations that change data.

Docker bills compute while the sandbox runs, including when you disconnect.
Review the compute price and timer before launching. Choose **Stop** as the
timeout action to retain files; **Delete** removes them. A stopped application
can't serve visitors. See [Sandbox lifecycle](/manuals/agentic-platform/sandboxes.md#manage-the-lifecycle).

## 1. Launch a shell

1. Open the [Console](https://agentic-platform.docker.com/) and select **+ New sandbox**.
2. Choose the **shell** kit. Deselect saved credentials and MCP tools.
3. Select **Balanced** for GitHub, Docker Hub, and npm access. Deselect **Open**
   if selected; combining it with Balanced doesn't restrict Open's access.
4. Review the compute size, **Run for** duration, and **When the time is up**
   action, then select **Run**.

The browser terminal opens a Bash shell inside the sandbox. Local files aren't
mounted there. See [Network policies](/manuals/agentic-platform/policies.md)
if your account applies additional restrictions.

## 2. Build and run the sample

In the sandbox terminal, clone the public project into the sandbox filesystem:

```console
$ cd /home/agent/workspace
$ git clone https://github.com/docker/welcome-to-docker.git project
$ cd project
$ docker build -t welcome-to-docker .
$ docker run -d --name welcome-to-docker -p 3000:3000 welcome-to-docker
```

The build should finish successfully, and the last command should return a
container ID. Check the running container:

```console
$ docker ps --filter name=welcome-to-docker
```

Expect the container to be running with port `3000` published in the sandbox.
This container port mapping alone doesn't create the public cloud URL.

## 3. Open and share the preview

1. On the sandbox detail page, select **Connect**.
2. Under **Connect via a Port**, enter `3000` and select **Open Port**.
3. Open the generated HTTPS URL in a separate private browser session where
   you aren't signed in to Docker. Expect the Welcome to Docker application.
4. Share the URL only after checking the page and confirming it contains no
   private data. Anyone with the URL can access the exposed application.

The service runs inside the sandbox. Closing your browser leaves it running
until you stop it or the timer action takes effect.

## 4. Withdraw access and save work

Use the close action next to port `3000` in **Connect**. In the separate browser
session, reload the URL and confirm that the application is no longer reachable.
Don't share the URL further if withdrawal hasn't succeeded; stop the sandbox
and investigate before continuing.

Stop the container in the sandbox terminal:

```console
$ docker stop welcome-to-docker
```

The project remains at `/home/agent/workspace/project`. To keep changes outside
the sandbox, [copy the files](/manuals/ai/sandboxes/cloud/usage.md#transfer-files)
or push them to a repository you control using an appropriately scoped
[GitHub credential](/manuals/agentic-platform/secrets.md#github-credential).
Verify your saved copy before deleting the sandbox; files don't sync automatically.

## 5. Clean up

In **Sandboxes**, pause the sandbox and confirm its paused status to stop compute
charges while retaining files. Resuming starts a fresh timer with the original
duration. Delete the sandbox after saving anything you need; deletion removes
files stored only there.

For CLI port controls, timers, and file transfers, see
[Use cloud sandboxes](/manuals/ai/sandboxes/cloud/usage.md).
