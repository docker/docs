1. Ensure Gordon is [enabled](/manuals/ai/gordon.md#enable-ask-gordon).
1. In Gordon's Toolbox, ensure Gordon's [Developer MCP Toolkit is enabled](/manuals/ai/gordon/mcp/built-in-tools.md#configuration).
1. In the terminal, navigate to the directory containing your Dockerfile.
1. Start a conversation with Gordon:
   ```bash
   docker ai
   ```
1. Type:
   ```console
   "Migrate my dockerfile to DHI"
   ```
1. Follow the conversation with Gordon. Gordon will edit your Dockerfile, so when
   it requests access to the filesystem and more, type `yes` to allow Gordon to proceed.

   > [!NOTE]
   > To learn more about Gordon's data retention and the data it
   > can access, see [Gordon](/manuals/ai/gordon.md#what-data-does-gordon-access).

When the migration is complete, you see a success message:

```text
The migration to Docker Hardened Images (DHI) is complete. The updated Dockerfile
successfully builds the image, and no vulnerabilities were detected in the final image.
The functionality and optimizations of the original Dockerfile have been preserved.
```

> [!IMPORTANT]
> As with any AI tool, you must verify Gordon's edits and test your image.
