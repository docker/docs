---
title: Docker Agentic Platform release notes
linkTitle: Release notes
description: Review feature updates, behavior changes, and fixes in Docker Agentic Platform.
keywords: docker agentic platform, release notes, updates, fixes
weight: 90
---

## September 24, 2026

This marks the experimental public release of Docker Agentic Platform.
Features and behavior may change. To begin, [activate your subscription](signup.md).

- Added a **Kits** catalog with search, curated and community filters, and links
  to Docker Hub. Selecting **Run** opens the sandbox launcher with the kit
  selected.
- Added public kit references in the launcher. Both the kit and its base image
  must be public.
- Added Hermes and Antigravity as curated kits. Hermes uses Anthropic or
  OpenAI credentials; Antigravity uses a Google credential.
- Updated the launcher with separate credential and GitHub controls. The
  launcher reuses saved keys and tokens that the kit supports without requiring
  you to select them again. Optional credentials can be deselected.
- Added an agent prompt and authoring documentation link on **Kits**, and a
  compact **Recent sandboxes** list on **New**.
- Added custom secrets for services outside the built-in provider list.
- Added platform selection and a **Restart** timer action for accounts with
  these controls enabled. Resuming a sandbox starts a fresh lifecycle timer.
- Added support for multiple shell tabs and port publishing from the
  **Connect** panel.
- Added a VS Code install link and Codex CLI support in **Add to your client**.
- Fixed always-applied policies to stay selected in the launcher.
- Improved sandbox quota errors with guidance on stopping or deleting
  sandboxes to free capacity.
- Updated sandbox creation to show provisioning progress and resume tracking
  after a page reload.

## September 14, 2026

- Added image attachments in Claude Code sandboxes. Paste or drop an image
  into the terminal to include it in your prompt.

## September 1, 2026

- Added options to save and select credentials in the sandbox launcher.

## August 28, 2026

- Fixed **Usage & billing** to show the active billing cycle and accrued
  compute costs for plans that renew mid-month.

## August 27, 2026

- Added Groq and xAI credentials under **Secrets** for OpenCode sandboxes.

## August 26, 2026

Run agents and tools in isolated cloud sandboxes with Docker Agentic Platform.

The initial release includes:

- Predefined sandbox types for Claude Code, Codex, OpenCode, Copilot, and
  Gemini CLI
- A terminal for working with your agent or shell
- Pause, resume, and delete controls for sandboxes
- Automatic sandbox stop or deletion after a timer from 1 to 24 hours
- Options to connect and authorize predefined MCP servers or add custom
  servers by URL
- API keys and tokens saved under **Secrets**, with their values kept outside
  sandboxes
- Network policies combining read-only kit rules with the read-only **Open**
  and **Balanced** presets or custom user policies. Deny rules take precedence
  over allow rules.
- A choice of compute sizes, billed by the second while the sandbox runs
