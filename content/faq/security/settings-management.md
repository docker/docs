---
description: Find the answers to common security related FAQs
keywords: Docker, Docker Hub, Docker Desktop secuirty FAQs, secuirty, platform
title: Settings Management FAQs
---

### Can I prevent our developers from running privileged containers?

Yes, by enabling [Enhanced Container Isolation](../../desktop/hardened-desktop/enhanced-container-isolation/_index.md) and locking that setting with [Settings Management](../../desktop/hardened-desktop/settings-management/_index.md). This is only available for Docker Business customers.

### How can I restrict access to certain settings so users cannot change things (e.g. enabling Kubernetes service, turning on send usage statistics, turning on experimental features etc.)? 

Yes, you can do that with [Settings Management](../../desktop/hardened-desktop/settings-management/_index.md).

### Can I prevent a developer from enabling the unsafe “Expose daemon on tcp://localhost:2375 without TLS” option on Windows?

Yes, you can do that with [Settings Management](../../desktop/hardened-desktop/settings-management/_index.md).

### Can I restrict the write access to `settings.json` to prevent modification by our developers?  

This would crash the application. For Docker Business customers however administrators can lock the Docker Desktop settings through the Settings Management feature by deploying an `admin-settings.json` file.