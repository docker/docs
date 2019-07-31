---
title: Manage webhooks via the API
description: Learn how to create, configure, and test webhooks for DTR using the API.
keywords: dtr, webhooks, api, registry
---

## Prerequisite

See [Webhook types](/ee/dtr/admin/manage-webhooks/index.md/#webhook-types) for a list of events you can trigger notifications for via the API.

## API Base URL

Your DTR hostname serves as the base URL for your API requests.

## Swagger API explorer

From the DTR web interface, click **API** on the bottom left navigation pane to explore the API resources and endpoints.  Click **Execute** to send your API request.

## API requests via curl

You can use [curl](https://curl.haxx.se/docs/manpage.html) to send HTTP or HTTPS API requests. Note that you will have to specify `skipTLSVerification: true` on your request in order to test the webhook endpoint over HTTP.

### Example curl request

```bash
curl -u test-user:$TOKEN -X POST "https://dtr-example.com/api/v0/webhooks" -H "accept: application/json" -H "content-type: application/json" -d "{ \"endpoint\": \"https://webhook.site/441b1584-949d-4608-a7f3-f240bdd31019\", \"key\": \"maria-testorg/lab-words\", \"skipTLSVerification\": true, \"type\": \"TAG_PULL\"}"
```

### Example JSON response

```json
{
  "id": "b7bf702c31601efb4796da59900ddc1b7c72eb8ca80fdfb1b9fecdbad5418155",
  "type": "TAG_PULL",
  "key": "maria-testorg/lab-words",
  "endpoint": "https://webhook.site/441b1584-949d-4608-a7f3-f240bdd31019",
  "authorID": "194efd8e-9ee6-4d43-a34b-eefd9ce39087",
  "createdAt": "2019-05-22T01:55:20.471286995Z",
  "lastSuccessfulAt": "0001-01-01T00:00:00Z",
  "inactive": false,
  "tlsCert": "",
  "skipTLSVerification": true
}
```

## Subscribe to events

To subscribe to events, send a `POST` request to
`/api/v0/webhooks` with the following JSON payload:

### Example usage

```
{
  "type": "TAG_PUSH",
  "key": "foo/bar",
  "endpoint": "https://example.com"
}
```

The keys in the payload are:

- `type`: The event type to subcribe to.
- `key`: The namespace/organization or repo to subscribe to. For example, "foo/bar" to subscribe to
pushes to the "bar" repository within the namespace/organization "foo".
- `endpoint`: The URL to send the JSON payload to.

Normal users **must** supply a "key" to scope a particular webhook event to
a repository or a namespace/organization. DTR admins can choose to omit this,
meaning a POST event notification of your specified type will be sent for all DTR repositories and namespaces.

### Receive a payload

Whenever your specified event type occurs, DTR will send a POST request to the given
endpoint with a JSON-encoded payload. The payload will always have the
following wrapper:

```
{
  "type": "...",
  "createdAt": "2012-04-23T18:25:43.511Z",
  "contents": {...}
}
```

- `type` refers to the event type received at the specified subscription endpoint.  
- `contents` refers to the payload of the event itself. Each event is different, therefore the
structure of the JSON object in `contents` will change depending on the event
type. See [Content structure](#content-structure) for more details.

### Test payload subscriptions

Before subscribing to an event, you can view and test your endpoints using
fake data. To send a test payload, send `POST` request to
`/api/v0/webhooks/test` with the following payload:

```
{
  "type": "...",
  "endpoint": "https://www.example.com/"
}
```

Change `type` to the event type that you want to receive. DTR will then send
an example payload to your specified endpoint. The example
payload sent is always the same.

## Content structure

Comments after (`//`) are for informational purposes only, and the example payloads have been clipped for brevity.

### Repository event content structure

**Tag push**

```
{
  "namespace": "",    // (string) namespace/organization for the repository
  "repository": "",   // (string) repository name
  "tag": "",          // (string) the name of the tag just pushed
  "digest": "",       // (string) sha256 digest of the manifest the tag points to (eg. "sha256:0afb...")
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar:tag)
  "os": "",           // (string) the OS for the tag's manifest
  "architecture": "", // (string) the architecture for the tag's manifest
  "author": "",       // (string) the username of the person who pushed the tag
  "pushedAt": "",     // (string) JSON-encoded timestamp of when the push occurred
  ...
}
```

**Tag delete**

```
{
  "namespace": "",    // (string) namespace/organization for the repository
  "repository": "",   // (string) repository name
  "tag": "",          // (string) the name of the tag just deleted
  "digest": "",       // (string) sha256 digest of the manifest the tag points to (eg. "sha256:0afb...")
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar:tag)
  "os": "",           // (string) the OS for the tag's manifest
  "architecture": "", // (string) the architecture for the tag's manifest
  "author": "",       // (string) the username of the person who deleted the tag
  "deletedAt": "",     // (string) JSON-encoded timestamp of when the delete occurred
  ...
}
```
**Manifest push**

```
{
  "namespace": "",    // (string) namespace/organization for the repository
  "repository": "",   // (string) repository name
  "digest": "",       // (string) sha256 digest of the manifest (eg. "sha256:0afb...")
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar@sha256:0afb...)
  "os": "",           // (string) the OS for the manifest
  "architecture": "", // (string) the architecture for the manifest
  "author": "",       // (string) the username of the person who pushed the manifest
  ...
}
```

**Manifest delete**

```
{
  "namespace": "",    // (string) namespace/organization for the repository
  "repository": "",   // (string) repository name
  "digest": "",       // (string) sha256 digest of the manifest (eg. "sha256:0afb...")
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar@sha256:0afb...)
  "os": "",           // (string) the OS for the manifest
  "architecture": "", // (string) the architecture for the manifest
  "author": "",       // (string) the username of the person who deleted the manifest
  "deletedAt": "",    // (string) JSON-encoded timestamp of when the delete occurred
  ...
}
```

**Security scan completed**

```
{
  "namespace": "",    // (string) namespace/organization for the repository
  "repository": "",   // (string) repository name
  "tag": "",          // (string) the name of the tag scanned
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar:tag)
  "scanSummary": {
    "namespace": "",          // (string) repository's namespace/organization name
    "repository": "",         // (string) repository name
    "tag": "",                // (string) the name of the tag just pushed
    "critical": 0,            // (int) number of critical issues, where CVSS >= 7.0
    "major": 0,               // (int) number of major issues, where CVSS >= 4.0 && CVSS < 7
    "minor": 0,               // (int) number of minor issues, where CVSS > 0 && CVSS < 4.0
    "last_scan_status": 0,    // (int) enum; see scan status section
    "check_completed_at": "", // (string) JSON-encoded timestamp of when the scan completed
    ...
  }
}
```

**Security scan failed**

```
{
  "namespace": "",    // (string) namespace/organization for the repository
  "repository": "",   // (string) repository name
  "tag": "",          // (string) the name of the tag scanned
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar@sha256:0afb...)
  "error": "",        // (string) the error that occurred while scanning
  ...
}
```

### Namespace-specific event structure

**Repository event (created/updated/deleted)**

```
{
  "namespace": "",    // (string) repository's namespace/organization name
  "repository": "",   // (string) repository name
  "event": "",        // (string) enum: "REPO_CREATED", "REPO_DELETED" or "REPO_UPDATED"
  "author": "",       // (string) the name of the user responsible for the event
  "data": {}          // (object) when updating or creating a repo this follows the same format as an API response from /api/v0/repositories/{namespace}/{repository}
}
```

### Global event structure

**Security scanner update complete**

```
{
  "scanner_version": "",
  "scanner_updated_at": "", // (string) JSON-encoded timestamp of when the scanner updated
  "db_version": 0,          // (int) newly updated database version
  "db_updated_at": "",      // (string) JSON-encoded timestamp of when the database updated
  "success": <true|false>   // (bool) whether the update was successful
  "replicas": {             // (object) a map keyed by replica ID containing update information for each replica
    "replica_id": {
      "db_updated_at": "",  // (string) JSON-encoded time of when the replica updated
      "version": "",        // (string) version updated to
      "replica_id": ""      // (string) replica ID
    },
    ...
  }
}
```

### Security scan status codes


- 0: **Failed**. An error occurred checking an image's layer
- 1: **Unscanned**. The image has not yet been scanned
- 2: **Scanning**. Scanning in progress
- 3: **Pending**: The image will be scanned when a worker is available
- 4: **Scanned**: The image has been scanned but vulnerabilities have not yet been checked
- 5: **Checking**: The image is being checked for vulnerabilities
- 6: **Completed**: The image has been fully security scanned


## View and manage existing subscriptions

### View all subscriptions

To view existing subscriptions, send a `GET` request to `/api/v0/webhooks`. As
a normal user (i.e. not a DTR admin), this will show all of your
current subscriptions across every namespace/organization and repository. As a DTR
admin, this will show **every** webhook configured for your DTR.

The API response will be in the following format:

```
[
  {
    "id": "",        // (string): UUID of the webhook subscription
    "type": "",      // (string): webhook event type
    "key": "",       // (string): the individual resource this subscription is scoped to
    "endpoint": "",  // (string): the endpoint to send POST event notifications to
    "authorID": "",  // (string): the user ID resposible for creating the subscription
    "createdAt": "", // (string): JSON-encoded datetime when the subscription was created
  },
  ...
]
```

For more information, [view the API documentation](/reference/dtr/{{site.dtr_version}}/api/).

### View subscriptions for a particular resource

You can also view subscriptions for a given resource that you are an
admin of. For example, if you have admin rights to the repository
"foo/bar", you can view all subscriptions (even other people's) from a
particular API endpoint. These endpoints are:

- `GET /api/v0/repositories/{namespace}/{repository}/webhooks`: View all
webhook subscriptions for a repository
- `GET /api/v0/repositories/{namespace}/webhooks`: View all webhook subscriptions for a
namespace/organization

### Delete a subscription

To delete a webhook subscription, send a `DELETE` request to
`/api/v0/webhooks/{id}`, replacing `{id}` with the webhook subscription ID
which you would like to delete.

Only a DTR admin or an admin for the resource with the event subscription can delete a subscription. As a normal user, you can only
delete subscriptions for repositories which you manage.

## Where to go next

- [Manage jobs](/ee/dtr/admin/manage-jobs/job-queue/)
