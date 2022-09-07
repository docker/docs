---
title: Deployment tracking
description:
keywords:
---

{% include atomist/disclaimer.md %}

By integrating Atomist with a runtime environment, you can track vulnerabilities
for deployed containers. This gives you contexts for whether security debt is
increasing or decreasing.

There are several options for how you could implement deployment tracking:

- Invoking the API directly
- Adding it as a step in your continuous deployment pipeline
- For Kubernetes, creating [admission controllers](./kubernetes.md)

## API

Each Atomist workspace exposes an API endpoint. Submitting a POST request to the
endpoint updates Atomist about what image you are running in your environments.
This will let you compare scan results for images you build against images of
containers running in staging or production.

You can find the API endpoint URL on the
[Integrations](https://dso.docker.com/r/auth/integrations) page. Using this API
requires an API key.

The most straight-forward usage is to post to this endpoint using a webhook.
When deploying a new image, submit an automated POST request (using `curl`, for
example) as part of your deployment pipeline.

```bash
$ curl <api-endpoint-url> \\
  -X POST \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer <api-token>" \\
  -d '{"image": {"url": "<image-url>@<sha256-digest>"}}'
```

### Parameters

The API supports the following parameters in the request body:

```json
{
  "image": {
    "url": "string",
    "name": "string"
  },
  "environment": {
    "name": "string"
  },
  "platform": {
    "os": "string",
    "architecture": "string",
    "variant": "string"
  }
}
```

| Parameter               | Mandatory | Default    | Description                                                                                                                           |
| ----------------------- | :-------: | ---------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| `image.url`             |    Yes    |            | Fully qualified reference name of the image, plus version (digest). You **must** specify the image version by digest.                 |
| `image.name`            |    No     |            | Optional identifier. If you deploy many containers from the same image in any one environment, each instance must have a unique name. |
| `environment.name`      |    No     | `deployed` | Use custom environment names to track different image versions in environments, like `staging` and `production`                       |
| `platform.os`           |    No     | `linux`    | Image operating system.                                                                                                               |
| `platform.architecture` |    No     | `amd64`    | Instruction set architecture.                                                                                                         |
| `platform.variant`      |    No     |            | Optional variant label.                                                                                                               |

