---
description: API Roles
keywords: API, Services, roles
redirect_from:
- /docker-cloud/feature-reference/api-roles/
title: Service API roles
notoc: true
---

You can configure a service so that it can access the Docker Cloud API. When you
grant API access to a service, its containers receive a token through an
environment variable, which is used to query the Docker Cloud API.

Docker Cloud has a "full access" role which when granted allows any operation
to be performed on the API. You can enable this option on the **Environment variables** screen of the Service wizard, or [specify it in your service's stackfile](stack-yaml-reference.md#roles). When enabled, Docker Cloud generates an authorization token for the
service's containers which is stored in an environment variable called
`DOCKERCLOUD_AUTH`.

Use this variable to set the `Authorization` HTTP header when calling
Docker Cloud's API:

```bash
$ curl -H "Authorization: $DOCKERCLOUD_AUTH" -H "Accept: application/json" https://cloud.docker.com/api/app/v1/service/
```

You can use this feature with Docker Cloud's [automatic environment variables](service-links.md), to let your application inside a container read and perform operations using Docker Cloud's API.

```bash
$ curl -H "Authorization: $DOCKERCLOUD_AUTH" -H "Accept: application/json" $WEB_DOCKERCLOUD_API_URL
```

For example, you can use information retrieved using the API to read the linked
endpoints, and use them to reconfigure a proxy container.

See the [API documentation](/apidocs/docker-cloud.md) for more information on the different API operations available.
