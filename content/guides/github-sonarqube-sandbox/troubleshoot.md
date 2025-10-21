---
title: Troubleshoot code quality workflows
linkTitle: Troubleshooting
summary: Resolve common issues with E2B sandboxes, MCP server connections, and GitHub/SonarQube integration.
description: Solutions for MCP tools not loading, authentication errors, permission issues, workflow timeouts, and other common problems when building code quality workflows with E2B.
weight: 30
---

This page covers common issues you might encounter when building code quality
workflows with E2B sandboxes and MCP servers, along with their solutions.

If you're experiencing problems not covered here, check the
[E2B documentation](https://e2b.dev/docs).

## MCP tools not available

Issue: Claude reports `I don't have any MCP tools available`.

Solution:

1. Verify you're using the authorization header:

    ```plaintext
    --header "Authorization: Bearer ${mcpToken}"
    ```

2. Check you're waiting for MCP initialization:

    ```javascript
    await new Promise(resolve => setTimeout(resolve, 1000));
    ```

3. Ensure credentials are in both `envs` and `mcp` configuration.
4. Verify your API tokens are valid and have proper scopes.

## GitHub tools work but SonarQube doesn't

Issue: GitHub MCP tools load but SonarQube tools don't appear.

Solution: SonarQube MCP server requires GitHub to be configured simultaneously.
Always include both servers in your sandbox configuration, even if you're only
testing one.

## Claude can’t access private repositories

Issue: “I don’t have access to that repository”.

Solution:

1. Verify your GitHub token has `repo` scope (not just `public_repo`).
2. Test with a public repository first.
3. Ensure the repository owner and name are correct in your `.env`.

## Workflow times out or runs too long

Issue: Workflow doesn’t complete or Claude credits run out.

Solutions:

1. Use `timeoutMs: 0` for complex workflows to allow unlimited time.
2. Break complex workflows into smaller, focused tasks.
3. Monitor your Anthropic API credit usage.
4. Add checkpoints in prompts: “After each step, show progress before continuing”.