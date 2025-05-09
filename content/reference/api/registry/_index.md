---
title: Docker Registry HTTP API V2
linktitle: Registry API
description: "Specification for the Registry API."
keywords: registry, on-prem, images, tags, repository, distribution, api, advanced
---

## Introduction

The _Docker Registry HTTP API_ is the protocol to facilitate distribution of
images to the Docker Engine. It interacts with instances of the Docker
Registry, which is a service to manage information about Docker images and
enable their distribution. The specification covers the operation of version 2
of this API, known as _Docker Registry HTTP API V2_.

### Scope

This specification covers the URL layout and protocols of the interaction
between Docker Registry and Docker core. This will affect the Docker core
registry API and the rewrite of docker-registry. Docker Registry
implementations may implement other API endpoints, but they are not covered by
this specification.

This includes the following features:

- Namespace-oriented URI Layout
- PUSH/PULL registry server for V2 image manifest format
- Resumable layer PUSH support
- V2 Client library implementation



## Overview

This section covers client flows and details of the API endpoints. The URI
layout of the new API is structured to support a rich authentication and
authorization model by leveraging namespaces. All endpoints will be prefixed
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
2. If a repository  name has two or more path components, they must be
   separated by a forward slash ("/").
3. The total length of a repository name, including slashes, must be less than
   256 characters.

These name requirements _only_ apply to the registry API and should accept a
superset of what is supported by other Docker ecosystem components.

All endpoints should support aggressive HTTP caching, compression and range
headers, where appropriate. The new API attempts to leverage HTTP semantics
where possible but may break from standards to implement targeted features.

For detail on individual endpoints, see the [Detail](#detail)
section.

### Errors

Actionable failure conditions, covered in detail in their relevant sections,
are reported as part of 4xx responses, in a JSON response body. One or more
errors will be returned in the following format:

```text
    {
        "errors:" [{
                "code": <error identifier>,
                "message": <message describing condition>,
                "detail": <unstructured>
            },
            ...
        ]
    }
```

The `code` field will be a unique identifier, all caps with underscores by
convention. The `message` field will be a human readable string. The optional
`detail` field may contain arbitrary JSON data providing information the
client can use to resolve the issue.

While the client can take action on certain error codes, the registry may add
new error codes over time. All client implementations should treat unknown
error codes as `UNKNOWN`, allowing future error codes to be added without
breaking API compatibility. For the purposes of the specification error codes
will only be added and never removed.

For a complete account of all error codes, see the [Errors](#errors-2)
section.

### API version check

A minimal endpoint, mounted at `/v2/` will provide version support information
based on its response statuses. The request format is as follows:

```text
    GET /v2/
```

If a `200 OK` response is returned, the registry implements the V2(.1)
registry API and the client may proceed safely with other V2 operations.
Optionally, the response may contain information about the supported paths in
the response body. The client should be prepared to ignore this data.

If a `401 Unauthorized` response is returned, the client should take action
based on the contents of the "WWW-Authenticate" header and try the endpoint
again. Depending on access control setup, the client may still have to
authenticate against different resources, even if this check succeeds.

If `404 Not Found` response status, or other unexpected status, is returned,
the client should proceed with the assumption that the registry does not
implement V2 of the API.

When a `200 OK` or `401 Unauthorized` response is returned, the
"Docker-Distribution-API-Version" header should be set to "registry/2.0".
Clients may require this header value to determine if the endpoint serves this
API. When this header is omitted, clients may fallback to an older API version.

### Content digests

This API design is driven heavily by [content addressability](http://en.wikipedia.org/wiki/Content-addressable_storage).
The core of this design is the concept of a content addressable identifier. It
uniquely identifies content by taking a collision-resistant hash of the bytes.
Such an identifier can be independently calculated and verified by selection
of a common _algorithm_. If such an identifier can be communicated in a secure
manner, one can retrieve the content from an insecure source, calculate it
independently and be certain that the correct content was obtained. Put simply,
the identifier is a property of the content.

To disambiguate from other concepts, this identifier is called a _digest_. A
_digest_ is a serialized hash result, consisting of a _algorithm_ and _hex_
portion. The _algorithm_ identifies the methodology used to calculate the
digest. The _hex_ portion is the hex-encoded result of the hash.

A _digest_ is defined as a string to match the following grammar:

```text
digest      := algorithm ":" hex
algorithm   := /[A-Fa-f0-9_+.-]+/
hex         := /[A-Fa-f0-9]+/
```

Some examples of _digests_ include the following:

| digest                                                                  | description                |
|-------------------------------------------------------------------------|----------------------------|
| sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b | Common sha256 based digest |

While the _algorithm_ does let one implement a wide variety of
algorithms, compliant implementations should use sha256. Heavy processing of
input before calculating a hash is discouraged to avoid degrading the
uniqueness of the _digest_ but some canonicalization may be performed to
ensure consistent identifiers.

Let's use a simple example in pseudo-code to demonstrate a digest calculation:

```text
let C = 'a small string'
let B = sha256(C)
let D = 'sha256:' + EncodeHex(B)
let ID(C) = D
```

In the previous example, you have bytestring `C` passed into a function, `SHA256`, that returns a
bytestring `B`, which is the hash of `C`. `D` gets the algorithm concatenated
with the hex encoding of `B`. We then define the identifier of `C` to `ID(C)`
as equal to `D`. A digest can be verified by independently calculating `D` and
comparing it with identifier `ID(C)`.

#### Digest header

To provide verification of HTTP content, any response may include a
`Docker-Content-Digest` header. This will include the digest of the target
entity returned in the response. For blobs, this is the entire blob content. For
manifests, this is the manifest body without the signature content, also known
as the JWS payload. Note that the commonly used canonicalization for digest
calculation may be dependent on the mediatype of the content, such as with
manifests.

The client may choose to ignore the header or may verify it to ensure content
integrity and transport security. This is most important when fetching by a
digest. To ensure security, the content should be verified against the digest
used to fetch the content. At times, the returned digest may differ from that
used to initiate a request. Such digests are considered to be from different
_domains_, meaning they have different values for _algorithm_. In such a case,
the client may choose to verify the digests in both domains or ignore the
server's digest. To maintain security, the client _must_ always verify the
content against the _digest_ used to fetch the content.


### Pulling an image

An "image" is a combination of a JSON manifest and individual layer files. The
process of pulling an image centers around retrieving these two components.

The first step in pulling an image is to retrieve the manifest. For reference,
the relevant manifest fields for the registry are the following:

| field     | description                                    |
|-----------|------------------------------------------------|
| name      | The name of the image.                         |
| tag       | The tag for this version of the image.         |
| fsLayers  | A list of layer descriptors (including digest) |
| signature | A JWS used to verify the manifest content      |


When the manifest is in hand, the client must verify the signature to ensure
the names and layers are valid. Once confirmed, the client will then use the
digests to download the individual layers. Layers are stored in as blobs in
the V2 registry API, keyed by their digest.

#### Pulling an image manifest

The image manifest can be fetched with the following URL:

```text
GET /v2/<name>/manifests/<reference>
```

The `name` and `reference` parameter identify the image and are required. The
reference may include a tag or digest.

The client should include an Accept header indicating which manifest content
types it supports. For more details on the manifest formats and their content
types, see
[manifest-v2-2.md](https://github.com/distribution/distribution/blob/main/docs/content/spec/manifest-v2-2.md).
In a successful response, the Content-Type header will indicate which manifest
type is being returned.

A `404 Not Found` response will be returned if the image is unknown to the
registry. If the image exists and the response is successful, the image
manifest will be returned.

```json
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

##### Existing manifests

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

#### Pulling a layer

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

### Pushing an image

Pushing an image works in the opposite order as a pull. After assembling the
image manifest, the client must first push the individual layers. When the
layers are fully pushed into the registry, the client should upload the signed
manifest.

The details of each step of the process are covered in the following sections.

#### Pushing a layer

All layer uploads use two steps to manage the upload process. The first step
starts the upload in the registry service, returning a URL to carry out the
second step. The second step uses the upload URL to transfer the actual data.
Uploads are started with a POST request which returns a URL that can be used
to push data and check upload status.

The `Location` header will be used to communicate the upload location after
each request. While it won't change in the this specification, clients should
use the most recent value returned by the API.

##### Starting an upload

To begin the process, a POST request should be issued in the following format:

```text
POST /v2/<name>/blobs/uploads/
```

The parameters of this request are the image namespace under which the layer
will be linked. Responses to this request are covered in the following sections.

##### Existing layers

The existence of a layer can be checked via a `HEAD` request to the blob store
API. The request should be formatted as follows:

```text
HEAD /v2/<name>/blobs/<digest>
```

If the layer with the digest specified in `digest` is available, a 200 OK
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

##### Uploading the layer

If the POST request is successful, a `202 Accepted` response will be returned
with the upload URL in the `Location` header:

```text
202 Accepted
Location: /v2/<name>/blobs/uploads/<uuid>
Range: bytes=0-<offset>
Content-Length: 0
Docker-Upload-UUID: <uuid>
```

The rest of the upload process can be carried out with the returned url,
called the "Upload URL" from the `Location` header. All responses to the
upload URL, whether sending data or getting status, will be in this format.
Though the URI format (`/v2/<name>/blobs/uploads/<uuid>`) for the `Location`
header is specified, clients should treat it as an opaque url and should never
try to assemble it. While the `uuid` parameter may be an actual UUID, this
proposal imposes no constraints on the format and clients should never impose
any.

If clients need to correlate local upload state with remote upload state, the
contents of the `Docker-Upload-UUID` header should be used. Such an id can be
used to key the last used location header when implementing resumable uploads.

##### Upload progress

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

##### Monolithic upload

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

The "digest" parameter must be included with the PUT request. See the
[Completed upload](#completed-upload) section for details on the parameters
and expected responses.

##### Chunked upload

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

##### Completed upload

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

###### Digest parameter

The "digest" parameter is designed as an opaque parameter to support
verification of a successful transfer. For example, an HTTP URI parameter
might be as follows:

```text
sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b
```

Given this parameter, the registry will verify that the provided content does
match this digest.

##### Canceling an upload

An upload can be canceled by issuing a DELETE request to the upload endpoint.
The format will be as follows:

```text
DELETE /v2/<name>/blobs/uploads/<uuid>
```

After this request is issued, the upload `uuid` will no longer be valid and the
registry server will dump all intermediate data. While uploads will time out
if not completed, clients should issue this request if they encounter a fatal
error but still have the ability to issue an HTTP request.

##### Cross repository blob mount

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

##### Errors

If an 502, 503 or 504 error is received, the client should assume that the
download can proceed due to a temporary condition, honoring the appropriate
retry mechanism. Other 5xx errors should be treated as terminal.

If there is a problem with the upload, a 4xx error will be returned indicating
the problem. After receiving a 4xx response (except 416, as called out
previously), the upload will be considered failed and the client should take
appropriate action.

Note that the upload URL will not be available forever. If the upload UUID is
unknown to the registry, a `404 Not Found` response will be returned and the
client must restart the upload process.

#### Pushing an Image Manifest

Once all of the layers for an image are uploaded, the client can upload the
image manifest. An image can be pushed using the following request format:

```json
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
content type should match the type of the manifest being uploaded, as specified
in [manifest-v2-2.md](https://github.com/distribution/distribution/blob/main/docs/content/spec/manifest-v2-2.md).

If there is a problem with pushing the manifest, a relevant 4xx response will
be returned with a JSON error message. See the
[PUT manifest](#put-manifest) section for details on possible error codes that
may be returned.

If one or more layers are unknown to the registry, `BLOB_UNKNOWN` errors are
returned. The `detail` field of the error response will have a `digest` field
identifying the missing blob. An error is returned for each unknown blob. The
response format is as follows:

```json
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

### Deleting an image

To delete a manifest, perform a DELETE request to a path in the following format: `/v2/<name>/manifests/<digest>`

<name> is the namespace of the repository, and <digest> is the digest of the manifest to be deleted. Upon success, the registry must respond with a 202 Accepted code. If the repository does not exist, the response must return 404 Not Found. If manifest deletion is disabled, the registry must respond with either a 400 Bad Request or a 405 Method Not Allowed.

In the event that a manifest cannot be deleted because it is referenced by another manifest, it will return 403.

Once deleted, a GET to '`/v2/<name>/manifests/<digest>` and any tag pointing to that digest will return a 404.

When deleting an image manifest that contains a subject field, and the referrers API returns a 404, clients should:

1. Pull the referrers list using the referrers tag schema.
2. Remove the descriptor entry from the array of manifests that references the deleted manifest.
3. Push the updated referrers list using the same referrers tag schema. The client may use conditional HTTP requests to prevent overwriting an referrers list that has changed since it was first pulled.

When deleting a manifest that has an associated referrers tag schema, clients may also delete the referrers tag when it returns a valid image index.

## Detail

The behavior of the endpoints are covered in detail in this section, organized
by route and entity. All aspects of the request and responses are covered,
including headers, parameters and body formats. Examples of requests and their
corresponding responses, with success and failure, are enumerated.

> [!NOTE]
>
> The sections on endpoint detail are arranged with an example
> request, a description of the request, followed by information about that
> request.

A list of methods and URIs are covered in the following table:

|Method|Path|Entity|Description|
|------|----|------|-----------|
| GET | `/v2/` | Base | Check that the endpoint implements Docker Registry API V2. |
| GET | `/v2/<name>/tags/list` | Tags | Fetch the tags under the repository identified by `name`. |
| GET | `/v2/<name>/manifests/<reference>` | Manifest | Fetch the manifest identified by `name` and `reference` where `reference` can be a tag or digest. A `HEAD` request can also be issued to this endpoint to obtain resource information without receiving all data. |
| PUT | `/v2/<name>/manifests/<reference>` | Manifest | Put the manifest identified by `name` and `reference` where `reference` can be a tag or digest. |
| DELETE | `/v2/<name>/manifests/<reference>` | Manifest | Delete the manifest identified by `name` and `reference`. Note that a manifest can _only_ be deleted by `digest`. |
| GET | `/v2/<name>/blobs/<digest>` | Blob | Retrieve the blob from the registry identified by `digest`. A `HEAD` request can also be issued to this endpoint to obtain resource information without receiving all data. |
| DELETE | `/v2/<name>/blobs/<digest>` | Blob | Delete the blob identified by `name` and `digest` |
| POST | `/v2/<name>/blobs/uploads/` | Initiate Blob Upload | Initiate a resumable blob upload. If successful, an upload location will be provided to complete the upload. Optionally, if the `digest` parameter is present, the request body will be used to complete the upload in a single request. |
| GET | `/v2/<name>/blobs/uploads/<uuid>` | Blob Upload | Retrieve status of upload identified by `uuid`. The primary purpose of this endpoint is to resolve the current status of a resumable upload. |
| PATCH | `/v2/<name>/blobs/uploads/<uuid>` | Blob Upload | Upload a chunk of data for the specified upload. |
| PUT | `/v2/<name>/blobs/uploads/<uuid>` | Blob Upload | Complete the upload specified by `uuid`, optionally appending the body as the final chunk. |
| DELETE | `/v2/<name>/blobs/uploads/<uuid>` | Blob Upload | Cancel outstanding upload processes, releasing associated resources. If this is not called, the unfinished uploads will eventually timeout. |
| GET | `/v2/_catalog` | Catalog | Retrieve a sorted, JSON list of repositories available in the registry. |


The detail for each endpoint is covered in the following sections.

### Errors

The error codes encountered via the API are enumerated in the following table:

| Code                    | Message                                        | Description                                                                                                                                                                                                                                                                                         |
|-------------------------|------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `BLOB_UNKNOWN`          | blob unknown to registry                       | This error may be returned when a blob is unknown to the registry in a specified repository. This can be returned with a standard get or if a manifest references an unknown layer during upload.                                                                                                   |
| `BLOB_UPLOAD_INVALID`   | blob upload invalid                            | The blob upload encountered an error and can no longer proceed.                                                                                                                                                                                                                                     |
| `BLOB_UPLOAD_UNKNOWN`   | blob upload unknown to registry                | If a blob upload has been canceled or was never started, this error code may be returned.                                                                                                                                                                                                           |
| `DIGEST_INVALID`        | provided digest did not match uploaded content | When a blob is uploaded, the registry will check that the content matches the digest provided by the client. The error may include a detail structure with the key "digest", including the invalid digest string. This error may also be returned when a manifest includes an invalid layer digest. |
| `MANIFEST_BLOB_UNKNOWN` | blob unknown to registry                       | This error may be returned when a manifest blob is  unknown to the registry.                                                                                                                                                                                                                        |
| `MANIFEST_INVALID`      | manifest invalid                               | During upload, manifests undergo several checks ensuring validity. If those checks fail, this error may be returned, unless a more specific error is included. The detail will contain information the failed validation.                                                                           |
| `MANIFEST_UNKNOWN`      | manifest unknown                               | This error is returned when the manifest, identified by name and tag is unknown to the repository.                                                                                                                                                                                                  |
| `MANIFEST_UNVERIFIED`   | manifest failed signature verification         | During manifest upload, if the manifest fails signature verification, this error will be returned.                                                                                                                                                                                                  |
| `NAME_INVALID`          | invalid repository name                        | Invalid repository name encountered either during manifest validation or any API operation.                                                                                                                                                                                                         |
| `NAME_UNKNOWN`          | repository name not known to registry          | This is returned if the name used during an operation is unknown to the registry.                                                                                                                                                                                                                   |
| `SIZE_INVALID`          | provided length did not match content length   | When a layer is uploaded, the provided size will be checked against the uploaded content. If they do not match, this error will be returned.                                                                                                                                                        |
| `TAG_INVALID`           | manifest tag did not match URI                 | During a manifest upload, if the tag in the manifest does not match the URI tag, this error will be returned.                                                                                                                                                                                       |
| `UNAUTHORIZED`          | authentication required                        | The access controller was unable to authenticate the client. Often this will be accompanied by a Www-Authenticate HTTP response header indicating how to authenticate.                                                                                                                              |
| `DENIED`                | requested access to the resource is denied     | The access controller denied access for the operation on a resource.                                                                                                                                                                                                                                |
| `UNSUPPORTED`           | The operation is unsupported.                  | The operation was unsupported due to a missing implementation or invalid set of parameters.                                                                                                                                                                                                         |

### Base

Base V2 API route. Typically, this can be used for lightweight version checks and to validate registry authentication.

#### GET base

Check that the endpoint implements Docker Registry API V2.

```text
GET /v2/
Host: <registry host>
Authorization: <scheme> <token>
```

The following parameters should be specified on the request:

|Name|Kind|Description|
|----|----|-----------|
|`Host`|header|Standard HTTP Host Header. Should be set to the registry host.|
|`Authorization`|header|An RFC7235 compliant authorization header.|

###### On Success: OK

```text
200 OK
```

The API implements V2 protocol and is accessible.

###### On Failure: Not Found

```text
404 Not Found
```

The registry does not implement the V2 API.

###### On Failure: Authentication Required

```text
401 Unauthorized
WWW-Authenticate: <scheme> realm="<realm>", ..."
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client is not authenticated.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`WWW-Authenticate`|An RFC7235 compliant authentication challenge header.|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated in the following:

|Code|Message|Description|
|----|-------|-----------|
| `UNAUTHORIZED` | authentication required | The access controller was unable to authenticate the client. Often this will be accompanied by a Www-Authenticate HTTP response header indicating how to authenticate. |

###### On Failure: Too Many Requests

```text
429 Too Many Requests
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client made too many requests within a time interval.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated in the following:

|Code|Message|Description|
|----|-------|-----------|
| `TOOMANYREQUESTS` | too many requests | Returned when a client attempts to contact a service too many times |

### Manifest

Create, update, delete and retrieve manifests.

#### GET manifest

Fetch the manifest identified by `name` and `reference` where `reference` can be a tag or digest. A `HEAD` request can also be issued to this endpoint to obtain resource information without receiving all data.

```text
GET /v2/<name>/manifests/<reference>
Host: <registry host>
Authorization: <scheme> <token>
```

The following parameters should be specified on the request:

|Name|Kind|Description|
|----|----|-----------|
|`Host`|header|Standard HTTP Host Header. Should be set to the registry host.|
|`Authorization`|header|An RFC7235 compliant authorization header.|
|`name`|path|Name of the target repository.|
|`reference`|path|Tag or digest of the target manifest.|

###### On Success: OK

```text
200 OK
Docker-Content-Digest: <digest>
Content-Type: <media type of manifest>

{
   "name": <name>,
   "tag": <tag>,
   "fsLayers": [
      {
         "blobSum": "<digest>"
      },
      ...
    ]
   ],
   "history": <v1 images>,
   "signature": <JWS>
}
```

The manifest identified by `name` and `reference`. The contents can be used to identify and resolve resources required to run the specified image.

The following headers will be returned with the response:

|Name|Description|
|----|-----------|
|`Docker-Content-Digest`|Digest of the targeted content for the request.|

###### On Failure: Bad Request

```text
400 Bad Request
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The name or reference was invalid.

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `NAME_INVALID` | invalid repository name | Invalid repository name encountered either during manifest validation or any API operation. |
| `TAG_INVALID` | manifest tag did not match URI | During a manifest upload, if the tag in the manifest does not match the URI tag, this error will be returned. |

###### On Failure: Authentication Required

```text
401 Unauthorized
WWW-Authenticate: <scheme> realm="<realm>", ..."
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client is not authenticated.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`WWW-Authenticate`|An RFC7235 compliant authentication challenge header.|
|`Content-Length`|Length of the JSON response body.|


The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `UNAUTHORIZED` | authentication required | The access controller was unable to authenticate the client. Often this will be accompanied by a Www-Authenticate HTTP response header indicating how to authenticate. |

###### On Failure: No Such Repository Error

```text
404 Not Found
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The repository is not known to the registry.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|


The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `NAME_UNKNOWN` | repository name not known to registry | This is returned if the name used during an operation is unknown to the registry. |

###### On Failure: Access Denied

```text
403 Forbidden
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client does not have required access to the repository.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|


The error codes that may be included in the response body are enumerated in the following:

|Code|Message|Description|
|----|-------|-----------|
| `DENIED` | requested access to the resource is denied | The access controller denied access for the operation on a resource. |

###### On Failure: Too Many Requests

```text
429 Too Many Requests
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client made too many requests within a time interval.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated in the following:

|Code|Message|Description|
|----|-------|-----------|
| `TOOMANYREQUESTS` | too many requests | Returned when a client attempts to contact a service too many times |

#### PUT manifest

Put the manifest identified by `name` and `reference` where `reference` can be a tag or digest.

```text
PUT /v2/<name>/manifests/<reference>
Host: <registry host>
Authorization: <scheme> <token>
Content-Type: <media type of manifest>

{
   "name": <name>,
   "tag": <tag>,
   "fsLayers": [
      {
         "blobSum": "<digest>"
      },
      ...
    ]
   ],
   "history": <v1 images>,
   "signature": <JWS>
}
```
The following parameters should be specified on the request:

|Name|Kind|Description|
|----|----|-----------|
|`Host`|header|Standard HTTP Host Header. Should be set to the registry host.|
|`Authorization`|header|An RFC7235 compliant authorization header.|
|`name`|path|Name of the target repository.|
|`reference`|path|Tag or digest of the target manifest.|

###### On Success: Created

```text
201 Created
Location: <url>
Content-Length: 0
Docker-Content-Digest: <digest>
```

The manifest has been accepted by the registry and is stored under the specified `name` and `tag`.

The following headers will be returned with the response:

|Name|Description|
|----|-----------|
|`Location`|The canonical location url of the uploaded manifest.|
|`Content-Length`|The `Content-Length` header must be zero and the body must be empty.|
|`Docker-Content-Digest`|Digest of the targeted content for the request.|

###### On Failure: Invalid Manifest

```text
400 Bad Request
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The received manifest was invalid in some way, as described by the error codes. The client should resolve the issue and retry the request.

The error codes that may be included in the response body are enumerated in the following:

|Code|Message|Description|
|----|-------|-----------|
| `NAME_INVALID` | invalid repository name | Invalid repository name encountered either during manifest validation or any API operation. |
| `TAG_INVALID` | manifest tag did not match URI | During a manifest upload, if the tag in the manifest does not match the URI tag, this error will be returned. |
| `MANIFEST_INVALID` | manifest invalid | During upload, manifests undergo several checks ensuring validity. If those checks fail, this error may be returned, unless a more specific error is included. The detail will contain information the failed validation. |
| `MANIFEST_UNVERIFIED` | manifest failed signature verification | During manifest upload, if the manifest fails signature verification, this error will be returned. |
| `BLOB_UNKNOWN` | blob unknown to registry | This error may be returned when a blob is unknown to the registry in a specified repository. This can be returned with a standard get or if a manifest references an unknown layer during upload. |

###### On Failure: Authentication Required

```text
401 Unauthorized
WWW-Authenticate: <scheme> realm="<realm>", ..."
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client is not authenticated.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`WWW-Authenticate`|An RFC7235 compliant authentication challenge header.|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `UNAUTHORIZED` | authentication required | The access controller was unable to authenticate the client. Often this will be accompanied by a Www-Authenticate HTTP response header indicating how to authenticate. |

###### On Failure: No Such Repository Error

```text
404 Not Found
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The repository is not known to the registry.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `NAME_UNKNOWN` | repository name not known to registry | This is returned if the name used during an operation is unknown to the registry. |

###### On Failure: Access Denied

```text
403 Forbidden
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client does not have required access to the repository.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated in the following:

|Code|Message|Description|
|----|-------|-----------|
| `DENIED` | requested access to the resource is denied | The access controller denied access for the operation on a resource. |

###### On Failure: Too Many Requests

```text
429 Too Many Requests
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client made too many requests within a time interval.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `TOOMANYREQUESTS` | too many requests | Returned when a client attempts to contact a service too many times |

###### On Failure: Missing Layer(s)

```text
400 Bad Request
Content-Type: application/json; charset=utf-8

{
    "errors:" [{
            "code": "BLOB_UNKNOWN",
            "message": "blob unknown to registry",
            "detail": {
                "digest": "<digest>"
            }
        },
        ...
    ]
}
```

One or more layers may be missing during a manifest upload. If so, the missing layers will be enumerated in the error response.

The error codes that may be included in the response body are enumerated in the following:

|Code|Message|Description|
|----|-------|-----------|
| `BLOB_UNKNOWN` | blob unknown to registry | This error may be returned when a blob is unknown to the registry in a specified repository. This can be returned with a standard get or if a manifest references an unknown layer during upload. |

###### On Failure: Not allowed

```text
405 Method Not Allowed
```

Manifest put is not allowed because the registry is configured as a pull-through cache or for some other reason

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `UNSUPPORTED` | The operation is unsupported. | The operation was unsupported due to a missing implementation or invalid set of parameters. |

#### DELETE manifest

Delete the manifest identified by `name` and `reference`. Note that a manifest can _only_ be deleted by `digest`.

```text
DELETE /v2/<name>/manifests/<reference>
Host: <registry host>
Authorization: <scheme> <token>
```

The following parameters should be specified on the request:

|Name|Kind|Description|
|----|----|-----------|
|`Host`|header|Standard HTTP Host Header. Should be set to the registry host.|
|`Authorization`|header|An RFC7235 compliant authorization header.|
|`name`|path|Name of the target repository.|
|`reference`|path|Tag or digest of the target manifest.|

###### On Success: Accepted

```text
202 Accepted
```

###### On Failure: Invalid Name or Reference

```text
400 Bad Request
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The specified `name` or `reference` were invalid and the delete was unable to proceed.

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `NAME_INVALID` | invalid repository name | Invalid repository name encountered either during manifest validation or any API operation. |
| `TAG_INVALID` | manifest tag did not match URI | During a manifest upload, if the tag in the manifest does not match the uri tag, this error will be returned. |

###### On Failure: Authentication Required

```text
401 Unauthorized
WWW-Authenticate: <scheme> realm="<realm>", ..."
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client is not authenticated.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`WWW-Authenticate`|An RFC7235 compliant authentication challenge header.|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `UNAUTHORIZED` | authentication required | The access controller was unable to authenticate the client. Often this will be accompanied by a Www-Authenticate HTTP response header indicating how to authenticate. |

###### On Failure: No Such Repository Error

```text
404 Not Found
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The repository is not known to the registry.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `NAME_UNKNOWN` | repository name not known to registry | This is returned if the name used during an operation is unknown to the registry. |

###### On Failure: Access Denied

```text
403 Forbidden
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client does not have required access to the repository.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `DENIED` | requested access to the resource is denied | The access controller denied access for the operation on a resource. |

###### On Failure: Too Many Requests

```text
429 Too Many Requests
Content-Length: <length>
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The client made too many requests within a time interval.

The following headers will be returned on the response:

|Name|Description|
|----|-----------|
|`Content-Length`|Length of the JSON response body.|

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `TOOMANYREQUESTS` | too many requests | Returned when a client attempts to contact a service too many times |

###### On Failure: Unknown Manifest

```text
404 Not Found
Content-Type: application/json; charset=utf-8

{
	"errors:" [
	    {
            "code": <error code>,
            "message": "<error message>",
            "detail": ...
        },
        ...
    ]
}
```

The specified `name` or `reference` are unknown to the registry and the delete was unable to proceed. Clients can assume the manifest was already deleted if this response is returned.

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `NAME_UNKNOWN` | repository name not known to registry | This is returned if the name used during an operation is unknown to the registry. |
| `MANIFEST_UNKNOWN` | manifest unknown | This error is returned when the manifest, identified by name and tag is unknown to the repository. |

###### On Failure: Not allowed

```text
405 Method Not Allowed
```

Manifest delete is not allowed because the registry is configured as a pull-through cache or `delete` has been disabled.

The error codes that may be included in the response body are enumerated below:

|Code|Message|Description|
|----|-------|-----------|
| `UNSUPPORTED` | The operation is unsupported. | The operation was unsupported due to a missing implementation or invalid set of parameters. |