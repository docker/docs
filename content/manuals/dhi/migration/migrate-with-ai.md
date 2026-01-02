--- 
title: Migrate using Docker's AI-powered assistant
linktitle: AI-assisted migration
description: Use Docker's AI-powered assistant to automatically migrate your Dockerfile to Docker Hardened Images
weight: 15
keywords: ai assistant, migrate dockerfile, docker hardened images, automated migration
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
---

{{< summary-bar feature_name="Ask Gordon DHI migration" >}}

You can use Docker's AI-powered assistant to automatically migrate your
Dockerfile to use Docker Hardened Images (DHI).

1. Ensure Docker's AI-powered assistant is [enabled](/manuals/ai/gordon.md#enable-ask-gordon).
2. In the terminal, navigate to the directory containing your Dockerfile.
3. Start a conversation with the assistant:
   ```bash
   docker ai
   ```
4. Type:
   ```console
   "Migrate my dockerfile to DHI"
   ```
5. Follow the conversation with the assistant. The assistant will edit your Dockerfile, so when
   it requests access to the filesystem and more, type `yes` to allow the assistant to proceed.

When the migration is complete, you see a success message:

```text
The migration to Docker Hardened Images (DHI) is complete. The updated Dockerfile
successfully builds the image, and no vulnerabilities were detected in the final image.
The functionality and optimizations of the original Dockerfile have been preserved.
```

> [!IMPORTANT]
>
> As with any AI tool, you must verify the assistant's edits and test your image.
