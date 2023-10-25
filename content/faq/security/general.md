---
description: Security FAQs
keywords: Docker, Docker Hub, Docker Desktop secuirty FAQs, secuirty, platform
title: Secuirty FAQs
---

## How does Docker Desktop handle and store authentication information?

Docker Desktop utilizes the host operating system's secure key management for handling and storing authentication tokens necessary for authenticating with image registries. On macOS, this is [Keychain](https://support.apple.com/guide/security/keychain-data-protection-secb0694df1a/web); on Windows, this is [Security and Identity API via Wincred](https://learn.microsoft.com/en-us/windows/win32/api/wincred/); and on Linux, this is [Pass](https://www.passwordstore.org/). 