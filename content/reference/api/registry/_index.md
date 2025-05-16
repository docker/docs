---
title: Supported registry API for Docker Hub
linktitle: Registry API
description: "Supported registry API endpoints."
keywords: registry, on-prem, images, tags, repository, distribution, api, advanced
---

## Introduction

Docker Hub is an OCI-compliant registry, which means it adheres to the open
standards defined by the Open Container Initiative (OCI) for distributing
container images. This ensures compatibility with a wide range of tools and
platforms in the container ecosystem. For full details, refer to the [OCI
Distribution
Specification](https://github.com/opencontainers/distribution-spec).

This topic provides a reference of selected registry API endpoints that are
supported by Docker Hub. It does not cover the entire OCI Distribution
Specification.

## Overview

All endpoints will be prefixed
by the API version and the repository name:

```text
/v2/<name>/
```

For example, an API endpoint that will work with the `library/ubuntu`
repository, the URI prefix will be:

```text
/v2/library/ubuntu/
```

This scheme provides rich access control over various operations and methods
using the URI prefix and HTTP methods that can be controlled in variety of
ways.

Classically, repository names have always been two path components where each
path component is less than 30 characters. The V2 registry API does not
enforce this. The rules for a repository name are as follows:

1. A repository name is broken up into _path components_. A component of a
   repository name must be at least one lowercase, alpha-numeric characters,
   optionally separated by periods, dashes or underscores. More strictly, it
   must match the regular expression `[a-z0-9]+(?:[._-][a-z0-9]+)*`.
2. If a repository  name has two or more path components, they must be separated
   by a forward slash ("/").
3. The total length of a repository name, including slashes, must be less than
   256 characters.

These name requirements only apply to the registry API and should accept a
superset of what is supported by other Docker ecosystem components.

## Pulling an image

An "image" is a combination of a JSON manifest and individual layer files. The
process of pulling an image centers around retrieving these two components.

The first step in pulling an image is to retrieve the manifest. For reference,
the relevant manifest fields for the registry are the following:

| Field    | Description                                    |
|----------|------------------------------------------------|
| name      | The name of the image.                         |
| tag       | The tag for this version of the image.         |
| fsLayers  | A list of layer descriptors (including digest). |
| signature | A JWS used to verify the manifest content.      |

When the manifest is in hand, the client must verify the signature to ensure
the names and layers are valid. Once confirmed, the client will then use the
digests to download the individual layers. Layers are stored in as blobs in
the V2 registry API, keyed by their digest.

### Pulling an image manifest

The image manifest can be fetched with the following URL:

```text
GET /v2/<name>/manifests/<reference>
```

The `name` and `reference` parameter identify the image and are required. The
reference may include a tag or digest.

The client should include an Accept header indicating which manifest content
types it supports. In a successful response, the Content-Type
header will indicate which manifest type is being returned.

A `404 Not Found` response will be returned if the image is unknown to the
registry. If the image exists and the response is successful, the image
manifest will be returned.

```text
    {
       "name": <name>,
       "tag": <tag>,
       "fsLayers": [
          {
             "blobSum": <digest>
          },
          ...
        ]
       ],
       "history": <v1 images>,
       "signature": <JWS>
    }
```

The client should verify the returned manifest signature for authenticity
before fetching layers.

#### Existing manifests

The image manifest can be checked for existence with the following URL:

```text
HEAD /v2/<name>/manifests/<reference>
```

The `name` and `reference` parameter identify the image and are required. The
reference may include a tag or digest.

A `404 Not Found` response will be returned if the image is unknown to the
registry. If the image exists and the response is successful the response will
be as follows:

```text
200 OK
Content-Length: <length of manifest>
Docker-Content-Digest: <digest>
```

### Pulling a Layer

Layers are stored in the blob portion of the registry, keyed by digest.
Pulling a layer is carried out by a standard HTTP request. The URL is as
follows:

```text
GET /v2/<name>/blobs/<digest>
```

Access to a layer will be gated by the `name` of the repository but is
identified uniquely in the registry by `digest`.

This endpoint may issue a 307 (302 for <HTTP 1.1) redirect to another service
for downloading the layer and clients should be prepared to handle redirects.

This endpoint should support aggressive HTTP caching for image layers. Support
for Etags, modification dates and other cache control headers should be
included. To allow for incremental downloads, `Range` requests should be
supported, as well.

## Pushing an image

Pushing an image works in the opposite order as a pull. After assembling the
image manifest, the client must first push the individual layers. When the
layers are fully pushed into the registry, the client should upload the signed
manifest.

The details of each step of the process are covered in the following sections.

### Pushing a layer

All layer uploads use two steps to manage the upload process. The first step
starts the upload in the registry service, returning a URL to carry out the
second step. The second step uses the upload URL to transfer the actual data.
Uploads are started with a POST request which returns a URL that can be used
to push data and check upload status.

The `Location` header will be used to communicate the upload location after
each request. While it won't change in the this specification, clients should
use the most recent value returned by the API.

#### Starting an upload

To begin the process, a POST request should be issued in the following format:

```text
POST /v2/<name>/blobs/uploads/
```

The parameters of this request are the image namespace under which the layer
will be linked. Responses to this request are covered in the following sections.

#### Existing layers

The existence of a layer can be checked via a `HEAD` request to the blob store
API. The request should be formatted as follows:

```text
HEAD /v2/<name>/blobs/<digest>
```

If the layer with the digest specified in `digest` is available, a `200 OK`
response will be received, with no actual body content (this is according to
HTTP specification). The response will look as follows:

```text
200 OK
Content-Length: <length of blob>
Docker-Content-Digest: <digest>
```

When this response is received, the client can assume that the layer is
already available in the registry under the given name and should take no
further action to upload the layer. Note that the binary digests may differ
for the existing registry layer, but the digests will be guaranteed to match.

#### Uploading the layer

If the POST request is successful, a `202 Accepted` response will be returned
with the upload URL in the `Location` header:

```text
202 Accepted
Location: /v2/<name>/blobs/uploads/<uuid>
Range: bytes=0-<offset>
Content-Length: 0
Docker-Upload-UUID: <uuid>
```

The rest of the upload process can be carried out with the returned URL,
called the "Upload URL" from the `Location` header. All responses to the
upload URL, whether sending data or getting status, will be in this format.
Though the URI format (`/v2/<name>/blobs/uploads/<uuid>`) for the `Location`
header is specified, clients should treat it as an opaque URL and should never
try to assemble it. While the `uuid` parameter may be an actual UUID, this
proposal imposes no constraints on the format and clients should never impose
any.

If clients need to correlate local upload state with remote upload state, the
contents of the `Docker-Upload-UUID` header should be used. Such an id can be
used to key the last used location header when implementing resumable uploads.

#### Upload Progress

The progress and chunk coordination of the upload process will be coordinated
through the `Range` header. While this is a non-standard use of the `Range`
header, there are examples of [similar approaches](https://developers.google.com/youtube/v3/guides/using_resumable_upload_protocol) in APIs with heavy use.
For an upload that just started, for an example with a 1000 byte layer file,
the `Range` header would be as follows:

```text
Range: bytes=0-0
```

To get the status of an upload, issue a GET request to the upload URL:

```text
GET /v2/<name>/blobs/uploads/<uuid>
Host: <registry host>
```

The response will be similar to the previous, except will return 204 status:

```text
204 No Content
Location: /v2/<name>/blobs/uploads/<uuid>
Range: bytes=0-<offset>
Docker-Upload-UUID: <uuid>
```

Note that the HTTP `Range` header byte ranges are inclusive and that will be
honored, even in non-standard use cases.

#### Monolithic upload

A monolithic upload is simply a chunked upload with a single chunk and may be
favored by clients that would like to avoided the complexity of chunking. To
carry out a "monolithic" upload, one can simply put the entire content blob to
the provided URL:

```text
PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
Content-Length: <size of layer>
Content-Type: application/octet-stream

<Layer Binary Data>
```

The "digest" parameter must be included with the `PUT` request. See the
[Completed Upload](#completed-upload) section for details on the parameters
and expected responses.

#### Chunked upload

To carry out an upload of a chunk, the client can specify a range header and
only include that part of the layer file:

```text
PATCH /v2/<name>/blobs/uploads/<uuid>
Content-Length: <size of chunk>
Content-Range: <start of range>-<end of range>
Content-Type: application/octet-stream

<Layer Chunk Binary Data>
```

There is no enforcement on layer chunk splits other than that the server must
receive them in order. The server may enforce a minimum chunk size. If the
server cannot accept the chunk, a `416 Requested Range Not Satisfiable`
response will be returned and will include a `Range` header indicating the
current status:

```text
416 Requested Range Not Satisfiable
Location: /v2/<name>/blobs/uploads/<uuid>
Range: 0-<last valid range>
Content-Length: 0
Docker-Upload-UUID: <uuid>
```

If this response is received, the client should resume from the "last valid
range" and upload the subsequent chunk. A 416 will be returned under the
following conditions:

- Invalid Content-Range header format
- Out of order chunk: the range of the next chunk must start immediately after
  the "last valid range" from the previous response.

When a chunk is accepted as part of the upload, a `202 Accepted` response will
be returned, including a `Range` header with the current upload status:

```text
202 Accepted
Location: /v2/<name>/blobs/uploads/<uuid>
Range: bytes=0-<offset>
Content-Length: 0
Docker-Upload-UUID: <uuid>
```

#### Completed upload

For an upload to be considered complete, the client must submit a `PUT`
request on the upload endpoint with a digest parameter. If it is not provided,
the upload will not be considered complete. The format for the final chunk
will be as follows:

```text
PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
Content-Length: <size of chunk>
Content-Range: <start of range>-<end of range>
Content-Type: application/octet-stream

<Last Layer Chunk Binary Data>
```

Optionally, if all chunks have already been uploaded, a `PUT` request with a
`digest` parameter and zero-length body may be sent to complete and validate
the upload. Multiple "digest" parameters may be provided with different
digests. The server may verify none or all of them but _must_ notify the
client if the content is rejected.

When the last chunk is received and the layer has been validated, the client
will receive a `201 Created` response:

```text
201 Created
Location: /v2/<name>/blobs/<digest>
Content-Length: 0
Docker-Content-Digest: <digest>
```

The `Location` header will contain the registry URL to access the accepted
layer file. The `Docker-Content-Digest` header returns the canonical digest of
the uploaded blob which may differ from the provided digest. Most clients may
ignore the value but if it is used, the client should verify the value against
the uploaded blob data.

##### Digest parameter

The "digest" parameter is designed as an opaque parameter to support
verification of a successful transfer. For example, an HTTP URI parameter
might be as follows:

```text
sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b
```

Given this parameter, the registry will verify that the provided content does
match this digest.

#### Canceling an upload

An upload can be canceled by issuing a DELETE request to the upload endpoint.
The format will be as follows:

```text
DELETE /v2/<name>/blobs/uploads/<uuid>
```

After this request is issued, the upload UUID will no longer be valid and the
registry server will dump all intermediate data. While uploads will time out
if not completed, clients should issue this request if they encounter a fatal
error but still have the ability to issue an HTTP request.

#### Cross repository blob mount

A blob may be mounted from another repository that the client has read access
to, removing the need to upload a blob already known to the registry. To issue
a blob mount instead of an upload, a POST request should be issued in the
following format:

```text
POST /v2/<name>/blobs/uploads/?mount=<digest>&from=<repository name>
Content-Length: 0
```

If the blob is successfully mounted, the client will receive a `201 Created`
response:

```text
201 Created
Location: /v2/<name>/blobs/<digest>
Content-Length: 0
Docker-Content-Digest: <digest>
```

The `Location` header will contain the registry URL to access the accepted
layer file. The `Docker-Content-Digest` header returns the canonical digest of
the uploaded blob which may differ from the provided digest. Most clients may
ignore the value but if it is used, the client should verify the value against
the uploaded blob data.

If a mount fails due to invalid repository or digest arguments, the registry
will fall back to the standard upload behavior and return a `202 Accepted` with
the upload URL in the `Location` header:

```text
202 Accepted
Location: /v2/<name>/blobs/uploads/<uuid>
Range: bytes=0-<offset>
Content-Length: 0
Docker-Upload-UUID: <uuid>
```

This behavior is consistent with older versions of the registry, which do not
recognize the repository mount query parameters.

Note: a client may issue a HEAD request to check existence of a blob in a source
repository to distinguish between the registry not supporting blob mounts and
the blob not existing in the expected repository.

#### Errors

If an 502, 503 or 504 error is received, the client should assume that the
download can proceed due to a temporary condition, honoring the appropriate
retry mechanism. Other 5xx errors should be treated as terminal.

If there is a problem with the upload, a 4xx error will be returned indicating
the problem. After receiving a 4xx response (except 416, as called out previously),
the upload will be considered failed and the client should take appropriate
action.

Note that the upload URL will not be available forever. If the upload UUID is
unknown to the registry, a `404 Not Found` response will be returned and the
client must restart the upload process.

### Pushing an image manifest

Once all of the layers for an image are uploaded, the client can upload the
image manifest. An image can be pushed using the following request format:

```text
    PUT /v2/<name>/manifests/<reference>
    Content-Type: <manifest media type>

    {
       "name": <name>,
       "tag": <tag>,
       "fsLayers": [
          {
             "blobSum": <digest>
          },
          ...
        ]
       ],
       "history": <v1 images>,
       "signature": <JWS>,
       ...
    }
```

The `name` and `reference` fields of the response body must match those
specified in the URL. The `reference` field may be a "tag" or a "digest". The
content type should match the type of the manifest being uploaded.

If there is a problem with pushing the manifest, a relevant 4xx response will
be returned with a JSON error message.

If one or more layers are unknown to the registry, `BLOB_UNKNOWN` errors are
returned. The `detail` field of the error response will have a `digest` field
identifying the missing blob. An error is returned for each unknown blob. The
response format is as follows:

```text
    {
        "errors:" [{
                "code": "BLOB_UNKNOWN",
                "message": "blob unknown to registry",
                "detail": {
                    "digest": <digest>
                }
            },
            ...
        ]
    }
```


## Deleting an image

To delete a manifest, perform a DELETE request to a path in the following
format: `/v2/<name>/manifests/<digest>`.

`<name>` is the namespace of the repository, and `<digest>` is the digest of the
manifest to be deleted. Upon success, the registry must respond with a 202
Accepted code. If the repository does not exist, the response must return 404
Not Found. If manifest deletion is disabled, the registry must respond with
either a 400 Bad Request or a 405 Method Not Allowed.

In the event that a manifest cannot be deleted because it is referenced by
another manifest, it will return 403.

Once deleted, a GET to '`/v2/<name>/manifests/<digest>` and any tag pointing to
that digest will return a 404.

When deleting an image manifest that contains a subject field, and the referrers
API returns a 404, clients SHOULD:

1. Pull the referrers list using the referrers tag schema.
2. Remove the descriptor entry from the array of manifests that references the
   deleted manifest.
3. Push the updated referrers list using the same referrers tag schema. The
   client may use conditional HTTP requests to prevent overwriting an referrers
   list that has changed since it was first pulled.

When deleting a manifest that has an associated referrers tag schema, clients
may also delete the referrers tag when it returns a valid image index.