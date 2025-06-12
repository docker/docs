---
title: Troubleshoot Docker Offload
linktitle: Troubleshoot
weight: 999
description: Learn how to troubleshoot issues with Docker Offload.
keywords: cloud, troubleshooting, cloud mode, Docker Desktop, cloud builder, usage
tags: [Troubleshooting]
---

Docker Offload requires:

- Authentication
- An active internet connection
- No restrictive proxy or firewall blocking traffic to Docker Cloud
- Beta access to Docker Offload
- Docker Desktop 4.43 or later

Docker Desktop uses Offload to run both builds and containers in the cloud.
If builds or containers are failing to run, falling back to local, or reporting
session errors, use the following steps to help resolve the issue.

1. Use the following command to check if the connection is active:

   ```console
   $ docker offload status
   ```

2. To get more information, run the following command:

   ```console
   $ docker offload diagnose
   ```

3. If you're not connected, start a new session:

   ```console
   $ docker offload start
   ```

4. Verify authentication with `docker login`.

5. If needed, you can sign out and then sign in again:

   ```console
   $ docker logout
   $ docker login
   ```

6. Verify your usage and billing. For more information, see [Docker Offload usage](/offload/usage/).