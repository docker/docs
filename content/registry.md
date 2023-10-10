---
title: Registry
description: Registry documentation has moved
keywords: registry, distribution
aliases:
  - /registry/compatibility/
  - /registry/deploying/
  - /registry/deprecated/
  - /registry/garbage-collection/
  - /registry/help/
  - /registry/insecure/
  - /registry/introduction/
  - /registry/notifications/
  - /registry/recipes/
  - /registry/recipes/apache/
  - /registry/recipes/nginx/
  - /registry/recipes/osx-setup-guide/
  - /registry/spec/
  - /registry/spec/api/
  - /registry/spec/auth/
  - /registry/spec/auth/jwt/
  - /registry/spec/auth/oauth/
  - /registry/spec/auth/scope/
  - /registry/spec/auth/token/
  - /registry/spec/deprecated-schema-v1/
  - /registry/spec/implementations/
  - /registry/spec/json/
  - /registry/spec/manifest-v2-1/
  - /registry/spec/manifest-v2-2/
  - /registry/spec/menu/
  - /registry/storage-drivers/
  - /registry/storage-drivers/azure/
  - /registry/storage-drivers/filesystem/
  - /registry/storage-drivers/gcs/
  - /registry/storage-drivers/inmemory/
  - /registry/storage-drivers/oss/
  - /registry/storage-drivers/s3/
  - /registry/storage-drivers/swift/
---

Registry, the open source implementation for storing and distributing container
images and other content, has been donated to the CNCF. Registry now goes under
the name of Distribution, and the documentation has moved to
[distribution.github.io/distribution][docs].

For reference documentation on the API protocol that Distribution implements,
see [Registry HTTP API][api].

For documentation related to authentication, see:

- [Token authentication specification][token]
- [OAuth 2.0 token authentication][oauth2]
- [JWT authentication][jwt]
- [Token scope and access][scope]

For information about image manifests, see:

- [Image Manifest Version 2, Schema 2][schema2]
- [Image Manifest Version 2, Schema 1][schema1] (deprecated)

[spec]: https://github.com/opencontainers/distribution-spec
[docs]: https://distribution.github.io/distribution/
[api]: https://distribution.github.io/distribution/spec/api/
[oauth2]: https://distribution.github.io/distribution/spec/auth/oauth/
[jwt]: https://distribution.github.io/distribution/spec/auth/jwt/
[token]: https://distribution.github.io/distribution/spec/auth/token/
[scope]: https://distribution.github.io/distribution/spec/auth/scope/
[schema2]: https://distribution.github.io/distribution/spec/manifest-v2-2/
[schema1]: https://distribution.github.io/distribution/spec/deprecated-schema-v1/
