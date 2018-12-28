---
title: Check Notary audit logs
description: When you push signed images, Docker Trusted Registry keeps audit
  logs for the changes made to the image metadata. Learn how to view these logs.
keywords: registry, monitor, troubleshoot
---

Docker Content Trust (DCT) keeps audit logs of changes made to trusted repositories.
Every time you push a signed image to a repository, or delete trust data for a
repository, DCT logs that information.

These logs are only available from the DTR API.

## Get an authentication token

To access the audit logs you need to authenticate your requests using an
authentication token. You can get an authentication token for all repositories,
or one that is specific to a single repository.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab1">Global</a></li>
  <li><a data-toggle="tab" data-target="#tab2">Repository-specific</a></li>
</ul>
<div class="tab-content">
  <div id="tab1" class="tab-pane fade in active" markdown="1">

```bash
curl --insecure --silent \
--user <user>:<password> \
"https://<dtr-url>/auth/token?realm=dtr&service=dtr&scope=registry:catalog:*"
```

  <hr></div>
  <div id="tab2" class="tab-pane fade" markdown="1">

```bash
curl --insecure --silent \
--user <user>:<password> \
"https://<dtr-url>/auth/token?realm=dtr&service=dtr&scope=repository:<dtr-url>/<repository>:pull"
```

  <hr></div>
</div>

DTR returns a JSON file with a token, even when the user doesn't have access
to the repository to which they requested the authentication token. This token
doesn't grant access to DTR repositories.

The JSON file returned has the following structure:


```json
{
  "token": "<token>",
  "access_token": "<token>",
  "expires_in": "<expiration in seconds>",
  "issued_at": "<time>"
}
```

[Learn more about the authentication API](/registry/spec/auth/token/).

## Changefeed API

Once you have an authentication token you can use the following endpoints to
get audit logs:

| URL                                                | Description                               | Authorization             |
|:---------------------------------------------------|:------------------------------------------|:--------------------------|
| `GET /v2/_trust/changefeed`                        | Get audit logs for all repositories.      | Global scope token        |
| `GET /v2/<dtr-url>/<repository>/_trust/changefeed` | Get audit logs for a specific repository. | Repository-specific token |

Both endpoints have the following query string parameters:

<table>
  <thead>
    <tr>
      <th>Field name</th>
      <th>Required</th>
      <th>Type</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>change_id</code></td>
      <td>Yes</td>
      <td>String</td>
      <td>
        <p>A non-inclusive starting change ID from which to start returning results. This will
        typically be the first or last change ID from the previous page of records
        requested, depending on which direction your are paging in.</p>
        <p>The value <code>0</code> indicates records should be returned starting from the beginning
        of time.</p>
        <p>The value <code>1</code> indicates records should be returned starting from the most
        recent record. If <code>1</code> is provided, the implementation will also assume the
        records value is meant to be negative, regardless of the given sign.</p>
      </td>
    </tr>
    <tr>
      <td><code>records</code></td>
      <td>Yes</td>
      <td>Signed integer</td>
      <td>
      <p>The number of records to return. A negative value indicates the number
      of records preceding the <code>change_id</code> should be returned.
      Records are always returned sorted from oldest to newest. </p>
      </td>
    </tr>
  </tbody>
</table>

### Response

The response is a JSON like:

```json
{
  "count": 1,
  "records": [
    {
      "ID": "0a60ec31-d2aa-4565-9b74-4171a5083bef",
      "CreatedAt": "2017-11-06T18:45:58.428Z",
      "GUN": "dtr.example.org/library/wordpress",
      "Version": 1,
      "SHA256": "a4ffcae03710ae61f6d15d20ed5e3f3a6a91ebfd2a4ba7f31fc6308ec6cc3e3d",
      "Category": "update"
    }
  ]
}
```

Below is the description for each of the fields in the response:

<table>
  <thead>
    <tr>
      <th>Field name</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>count</code></td>
      <td>The number of records returned.</td>
    </tr>
    <tr>
      <td><code>ID</code></td>
      <td>The ID of the change record. Should be used in the <code>change_id</code>
      field of requests to provide a non-exclusive starting index. It should be treated
      as an opaque value that is guaranteed to be unique within an instance of
      notary.</td>
    </tr>
    <tr>
      <td><code>CreatedAt</code></td>
      <td>The time the change happened.</td>
    </tr>
    <tr>
      <td><code>GUN</code></td>
      <td>The DTR repository that was changed.</td>
    </tr>
    <tr>
      <td><code>Version</code></td>
      <td>
        <p>The version that the repository was updated to. This increments
        every time there's a change to the trust repository.</p>
        <p>This is always <code>0</code> for events representing trusted data being
        removed from the repository.</p>
        </td>
    </tr>
    <tr>
      <td><code>SHA256</code></td>
      <td>
        <p>The checksum of the timestamp being updated to. This can be used
        with the existing notary APIs to request said timestamp.</p>
        <p>This is always an empty string for events representing trusted data
        being removed from the repository</p>
      </td>
    </tr>
    <tr>
      <td><code>Category</code></td>
      <td>The kind of change that was made to the trusted repository. Can be
      <code>update</code>, or <code>deletion</code>.</td>
    </tr>
  </tbody>
</table>


The results only include audit logs for events that happened more than 60
seconds ago, and are sorted from oldest to newest.

Even though the authentication API always returns a token, the changefeed
API validates if the user has access to see the audit logs or not:

* If the user is an admin they can see the audit logs for any repositories,
* All other users can only see audit logs for repositories they have read access.


## Example - Get audit logs for all repositories

Before going through this example, make sure that you:

* Are a DTR admin user.
* Configured your [machine to trust DTR](../../user/access-dtr/index.md).
* Created the `library/wordpress` repository.
* Installed `jq`, to make it easier to parse the JSON responses.

```bash
# Pull an image from Docker Hub
docker pull wordpress:latest

# Tag that image
docker tag wordpress:latest <dtr-url>/library/wordpress:1

# Log into DTR
docker login <dtr-url>

# Push the image to DTR and sign it
DOCKER_CONTENT_TRUST=1 docker push <dtr-url>/library/wordpress:1

# Get global-scope authorization token, and store it in TOKEN
export TOKEN=$(curl --insecure --silent \
--user '<user>:<password>' \
'https://<dtr-url>/auth/token?realm=dtr&service=dtr&scope=registry:catalog:*' | jq --raw-output .token)

# Get audit logs for all repositories and pretty-print it
# If you pushed the image less than 60 seconds ago, it's possible
# That DTR doesn't show any events. Retry the command after a while.
curl --insecure --silent \
--header "Authorization: Bearer $TOKEN" \
"https://<dtr-url>/v2/_trust/changefeed?records=10&change_id=0" | jq .
```

## Example - Get audit logs for a single repository

Before going through this example, make sure that you:

* Configured your [machine to trust DTR](../../user/access-dtr/index.md).
* Created the `library/nginx` repository.
* Installed `jq`, to make it easier to parse the JSON responses.

```bash
# Pull an image from Docker Hub
docker pull nginx:latest

# Tag that image
docker tag nginx:latest <dtr-url>/library/nginx:1

# Log into DTR
docker login <dtr-url>

# Push the image to DTR and sign it
DOCKER_CONTENT_TRUST=1 docker push <dtr-url>/library/nginx:1

# Get global-scope authorization token, and store it in TOKEN
export TOKEN=$(curl --insecure --silent \
--user '<user>:<password>' \
'https://<dtr-url>/auth/token?realm=dtr&service=dtr&scope=repository:<dtr-url>/<repository>:pull' | jq --raw-output .token)

# Get audit logs for all repositories and pretty-print it
# If you pushed the image less than 60 seconds ago, it's possible that
# Docker Content Trust won't show any events. Retry the command after a while.
curl --insecure --silent \
--header "Authorization: Bearer $TOKEN" \
"https://<dtr-url>/v2/<dtr-url>/<dtr-repo>/_trust/changefeed?records=10&change_id=0" | jq .
```

