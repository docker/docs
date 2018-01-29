---
title: Manage webhooks
description: Learn how to create, configure, and test webhooks in Docker Trusted Registry.
keywords: docker, registry, webhooks
---

DTR includes webhooks for common events, such as pushing a new tag or deleting
an image. This lets you build complex CI and CD pipelines from your own DTR
cluster.

The webhook events you can subscribe to are:

**Repository specific events**

- Tag push
- Tag delete
- Manifest push
- Manifest delete
- Security scan completed
- Security scan failed

**Namespace specific events**

- Repo events (created/updated/deleted)

**Global events**

- Security scanner update complete

You need to be at least an admin of the repository or namespace in question to
subscribe to an event for the repository or namespace. A global administrator can subscribe to any event.
For example, a user must be an admin of repository "foo/bar" to subscribe to
its tag push events.

## Subscribing to events

To subscribe to events you must send an API query to
`/api/v0/webhooks` with the following payload:

```
{
  "type": "",
  "key": "",
  "endpoint": "https://example.com"
}
```

The keys in the payload are:

- `type`: The event type. These are listed below.
- `key`: The namespace or repo to subscribe to (eg "foo/bar" to subscribe to
pushes to repo "bar" within namespace "foo")
- `endpoint`: The URL to send the POST payload to.

Normal users **must** supply a "key" to scope a particular webhook event to
a repository or namespace. System administrators can choose to omit this,
meaning a webhook payload will be sent for **every** repository or namespace.

### Event types

**Repository specific events**

- Tag push                (type `TAG_PUSH`)
- Tag delete              (type `TAG_DELETE`)
- Manifest push           (type `MANIFEST_PUSH`)
- Manifest delete         (type `MANIFEST_DELETE`)
- Security scan completed (type `SCAN_COMPLETED`)
- Security scan failed    (type `SCAN_FAILED`)

**Namespace specific events**

- Repo created/updated/deleted (type `REPO_EVENT`)

**Global events**

- Security scanner update complete (type `SCANNER_UPDATE_COMPLETED`)


### Receiving a payload

Whenever a suitable action occurs DTR will send a POST request to the given
endpoint with a JSON encoded payload. The payload will always have the
following wrapper:

```
{
  "type": "...",
  "createdAt": "2012-04-23T18:25:43.511Z",
  "contents": {...}
}
```

`type` refers to the event type being sent, and `contents` refers to
the payload of the event itself. Each event is different, therefore the
structure of the JSON object in `contents` will change depending on the event
type. These are listed in a section below.

### Testing payload subscriptions

Before subscribing to a payload you can view and test your endpoints using
fake data. To send a test payload, fire a `POST` request to
`/api/v0/webhooks/test` with the following payload:

```
{
  "type": "...",
  "endpoint": "https://www.example.com/"
}
```

Change `type` to the event type that you want to receive. DTR will then send
an example payload to the endpoint specified. The example
payload sent is always the same.

## Content structure

Comments (`// here`) are added for documentation only; they are not
present in POST payloads.

### Repository event content structure

**Tag push**

```
{
  "namespace": "",    // (string) repository's namespace name
  "repository": "",   // (string) repository name
  "tag": "",          // (string) the name of the tag just pushed
  "digest": "",       // (string) sha256 digest of the manifest the tag points to (eg. "sha256:0afb...")
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar:tag)
  "os": "",           // (string) the OS for the tag's manifest
  "architecture": "", // (string) the architecture for the tag's manifest
  "author": "",       // (string) the username of the person who pushed the tag
  "pushedAt": "",     // (string) JSON encoded timestamp of when the push occurred
}
```

**Tag delete**

```
{
  "namespace": "",    // (string) repository's namespace name
  "repository": "",   // (string) repository name
  "tag": "",          // (string) the name of the tag just deleted
  "digest": "",       // (string) sha256 digest of the manifest the tag points to (eg. "sha256:0afb...")
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar:tag)
  "os": "",           // (string) the OS for the tag's manifest
  "architecture": "", // (string) the architecture for the tag's manifest
  "author": "",       // (string) the username of the person who deleted the tag
  "deletedAt": "",     // (string) JSON encoded timestamp of when the delete occurred
}
```
**Manifest push**

```
{
  "namespace": "",    // (string) repository's namespace name
  "repository": "",   // (string) repository name
  "digest": "",       // (string) sha256 digest of the manifest (eg. "sha256:0afb...")
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar@sha256:0afb...)
  "os": "",           // (string) the OS for the manifest
  "architecture": "", // (string) the architecture for the manifest
  "author": "",       // (string) the username of the person who pushed the manifest
  "pushedAt": "",     // (string) JSON encoded timestamp of when the push occurred
}
```

**Manifest delete**

```
{
  "namespace": "",    // (string) repository's namespace name
  "repository": "",   // (string) repository name
  "digest": "",       // (string) sha256 digest of the manifest (eg. "sha256:0afb...")
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar@sha256:0afb...)
  "os": "",           // (string) the OS for the manifest
  "architecture": "", // (string) the architecture for the manifest
  "author": "",       // (string) the username of the person who deleted the manifest
  "deletedAt": "",    // (string) JSON encoded timestamp of when the delete occurred
}
```

**Security scan completed**

```
{
  "namespace": "",    // (string) repository's namespace name
  "repository": "",   // (string) repository name
  "tag": "",          // (string) the name of the tag scanned
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar:tag)
  "scanSummary": {
    "namespace": "",          // (string) repository's namespace name
    "repository": "",         // (string) repository name
    "tag": "",                // (string) the name of the tag just pushed
    "critical": 0,            // (int) number of critical issues, where CVSS >= 7.0
    "major": 0,               // (int) number of major issues, where CVSS >= 4.0 && CVSS < 7
    "minor": 0,               // (int) number of minor issues, where CVSS > 0 && CVSS < 4.0
    "last_scan_status": 0,    // (int) enum; see scan status section
    "check_completed_at": "", // (string) JSON encoded timestamp of when the scan completed
  }
}
```

**Security scan failed**

```
{
  "namespace": "",    // (string) repository's namespace name
  "repository": "",   // (string) repository name
  "tag": "",          // (string) the name of the tag scanned
  "imageName": "",    // (string) the fully-qualified image name including DTR host used to pull the image (eg. 10.10.10.1/foo/bar@sha256:0afb...)
  "error": "",        // (string) the error that occurred whilst scanning
}
```

### Namespace specific event structure

**Repository event (created/updated/deleted)**

```
{
  "namespace": "",    // (string) repository's namespace name
  "repository": "",   // (string) repository name
  "event": "",        // (string) enum: "REPO_CREATED", "REPO_DELETED" or "REPO_UPDATED"
  "author": "",       // (string) the name of the user which authored the event
  "data": {}          // (object) when updating or creating a repo this follows the same format as an API response from /api/v0/repositories/{namespace}/{repository}
}
```

### Global event structure

**Security scanner update complete**

```
{
  "scanner_version": "",
  "scanner_updated_at": "", // (string) JSON encoded timestamp of when the scanner updated
  "db_version": 0,          // (int) newly updated database version
  "db_updated_at": "",      // (string) JSON encoded timestamp of when the database updated
  "success": <true|false>   // (bool) whether the update was successful
  "replicas": {             // (object) a map keyed by replica ID containing update information for each replica
    "replica_id": {
      "db_updated_at": "",  // (string) JSON encoded time of when the replica updated
      "version": "",        // (string) version updated to
      "replica_id": ""      // (string) replica ID
    },
    ...
  }
}
```

### Security scan status enums


- 0: Failed. An error occurred checking an image's layer
- 1: Unscanned. The image has not yet been scanned
- 2: Scanning: Scanning in progress
- 3: Pending: The image will be scanned when a worker is available
- 4: Scanned: The image has been scanned but vulnerabilities have not yet been checked
- 5: Checking: The image is being checked for vulnerabilities
- 6: Completed: The image has been fully security scanned


## Viewing and managing existing subscriptions

### Viewing all subscriptions

To view existing subscriptions make a GET request to `/api/v0/webhooks`. As
a normal user (ie. not a system administrator), this will show all of your
current subscriptions across every namespace and repository. As a system
administrator this will show **every** webhook within DTR.

The API response will be in the following format:

```
[
  {
    "id": "",        // (string): UUID of the webhook subscription
    "type": "",      // (string): webhook event type
    "key": "",       // (string): the individual resource this subscription is scoped to
    "endpoint": "",  // (string): the endpoint in which the payload will be sent to via a POST request
    "authorID": "",  // (string): the ID of the user which created this subscription
    "createdAt": "", // (string): JSON encoded datetime when this subscription was created
  },
  ...
]
```

For more information view the API documentation.

### Viewing subscriptions for a particular resource

You can also view subscriptions for a given resource that you are an
administrator of. For example, if you have admin rights to the repository
"foo/bar" you can view all subscriptions (even other people's) from a
particular API endpoint. These endpoints are:

- `/api/v0/repositories/{namespace}/{repository}/webhooks`: GET to view all
webhooks for a repository resource's events
- `/api/v0/repositories/{namespace}/webhooks`: GET to view all webhooks for a
namespace resource's events

### Deleting a subscription

To delete a webhook subscription send a `DELETE` request to
`/api/v0/webhooks/{id}`, replacing `{id}` with the webhook subscription ID
which you would like to delete.

Only a system administrator or an administrator for the resource which the
payload subscribes to can delete a subscription. As a normal user, you can only
delete subscriptions for repositories which you administer.
