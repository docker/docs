---
title: Secrets top-level element
keywords: compose, compose specification
fetch_remote:
  line_start: 2
  line_end: -1
---

Secrets provide a more secure way of getting sensitive information in to your application's services, so you don't have to rely on using environment variables. If you’re injecting passwords and API keys as environment variables, you risk unintentional information exposure. Environment variables are often available to all processes, and it can be difficult to track access. They can also be printed in logs when debugging errors without your knowledge. Using Secrets mitigates these risks.
