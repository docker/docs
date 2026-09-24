---
title: Docker Sandboxes API reference
layout: api-reference-generated
linkTitle: Sandboxes
description: HTTP API reference for Docker Cloud Sandboxes.
keywords: cloud sandboxes, sandboxes api, api reference, openapi
weight: 32
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
---

> [!NOTE]
> The Docker Sandboxes API and SDK are experimental. Features, interfaces,
> and behavior may change.

This reference covers the public HTTP API. For runnable programs, start with
[the SDK guide](/manuals/ai/sandboxes-api/_index.md). [Authenticate to Docker](/manuals/ai/sandboxes-api/authentication.md)
before making requests.

<a href="api.yaml" download>Download the OpenAPI document</a> to use these operations with your own
HTTP tools. The download contains the same paths, request formats, and responses
shown here, including streaming and interactive operations.

Use this API to create sandboxes, run processes, work with files, and manage sandbox resources. Start with a kit through the SDK, or use the HTTP operations here directly.

Management requests use an Authorization bearer token. File and process requests go to the sandbox endpoint in core.endpoint.uri and require a short-lived scoped endpoint credential, not the management token. The SDK obtains that credential for you. Each operation also checks its declared permissions. Direct Unix socket access uses the operating system's socket access controls instead of HTTP authentication. Interactive WebSocket operations describe their supported credential transports separately.

The Docker Cloud API base URL is https://connect.docker.com/sandboxes. Append each /v1 path without removing the base URL's path prefix. For sandbox endpoint operations, use the sandbox endpoint as the base URL instead.

Operations list their success and error responses. A 202 response means work is still in progress; read the resource until it reaches the expected state. Errors contain a stable code, a message, and optional details. Inspect the code as well as the HTTP status.

Send the current ETag in If-Match when an operation requires it. Operations that accept Idempotency-Key use that HTTP header to identify retries: the same key and payload return the first result; a different payload fails. Keys are retained for at least 24 hours.

Lists use pageSize and pageToken. Continue with nextPageToken until it is empty. Supported filters and ordering are listed on each operation. Resource names contain immutable identifiers; displayName is a label. A 64-bit integer is returned as a JSON string; inputs accept a number or a string.

## Operations

60 operations are grouped by resource. Open a request or response schema to see its complete definition. Named schemas link to the definitions at the end of this page.

## Files

Read, write, copy, and manage files inside a sandbox.

| Method | Path | Operation |
| --- | --- | --- |
| GET | `/v1/files` | [`list`](#operation-list-get-v1-files) |
| DELETE | `/v1/files` | [`remove`](#operation-remove-delete-v1-files) |
| GET | `/v1/files/content` | [`readFile`](#operation-readFile-get-v1-files-content) |
| PUT | `/v1/files/content` | [`writeFile`](#operation-writeFile-put-v1-files-content) |
| PUT | `/v1/files/directories` | [`mkdir`](#operation-mkdir-put-v1-files-directories) |
| GET | `/v1/files/download` | [`download`](#operation-download-get-v1-files-download) |
| POST | `/v1/files/move` | [`move`](#operation-move-post-v1-files-move) |
| GET | `/v1/files/stat` | [`stat`](#operation-stat-get-v1-files-stat) |
| POST | `/v1/files/upload` | [`upload`](#operation-upload-post-v1-files-upload) |

### GET `/v1/files` {#operation-list-get-v1-files}

List lists directory entries by path ascending.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesFilesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `path` | query | Yes | `string` | path is the directory path. |
| `pageSize` | query | No | `integer` | page_size is an optional page size. Omitted or zero uses the backend default. Unless the operation states otherwise, page-size limits and handling of larger requests are backend-specific; use the backend support guide. Continue with nextPageToken until it is empty. |
| `pageToken` | query | No | `string` | page_token is an optional continuation token. |
| `orderBy` | query | No | `string` | order_by is a single order field with optional direction; the portable field is path. |


#### Responses

**200** Success

`application/json`: [`FilesListResponse`](#schema-FilesListResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/FilesListResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### DELETE `/v1/files` {#operation-remove-delete-v1-files}

Remove deletes a path. Removing an absent path succeeds. A nonrecursive
removal of a nonempty directory returns FAILED_PRECONDITION unchanged.
A recursive removal that stops partway still answers OK: check
RemoveResponse.failed_path, not the operation status. Removed paths stay removed.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesFilesWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `path` | query | Yes | `string` | path is the sandbox path to remove. |
| `recursive` | query | No | `boolean` | recursive removes directory contents when true. |


#### Responses

**200** Success

`application/json`: [`FilesRemoveResponse`](#schema-FilesRemoveResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/FilesRemoveResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/files/content` {#operation-readFile-get-v1-files-content}

ReadFile returns raw file content over HTTP; Stat exposes file metadata.
Content above the unary message cap is refused, never truncated.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesFilesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `path` | query | Yes | `string` | path is the sandbox path of a regular file. |


#### Responses

**200** Success

Header `Content-Length`: `integer`. Exact byte count of the complete file; this response does not use chunked transfer encoding.

`application/octet-stream`: `Object`

<details>
<summary>200 response schema</summary>

```yaml
contentMediaType: application/octet-stream
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### PUT `/v1/files/content` {#operation-writeFile-put-v1-files-content}

WriteFile accepts raw content with path and mode query parameters.
It replaces the complete target state and carries no request_id.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesFilesWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `path` | query | Yes | `string` | path is the sandbox path. |
| `mode` | query | No | `integer` | mode carries Unix permission bits. Upload applies mode & 0777, or 0644 when that value is zero, without subtracting umask; Download reports stored bits. |


#### Request body

Required.

**Content-Type:** `application/octet-stream`

`Object`

<details>
<summary>Request schema</summary>

```yaml
contentMediaType: application/octet-stream
```

</details>



#### Responses

**200** Success

`application/json`: [`FilesWriteFileResponse`](#schema-FilesWriteFileResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/FilesWriteFileResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### PUT `/v1/files/directories` {#operation-mkdir-put-v1-files-directories}

Mkdir creates a directory. Without parents, a missing parent returns
NOT_FOUND without creating directories. A repeated Mkdir converges when
the leaf already exists as a directory. A converged call leaves
the directory's mode unchanged on every request.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesFilesWrite`.

#### Request body

Required.

**Content-Type:** `application/json`

[`FilesMkdirRequest`](#schema-FilesMkdirRequest)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/FilesMkdirRequest"
```

</details>



#### Responses

**200** Success

`application/json`: [`FilesMkdirResponse`](#schema-FilesMkdirResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/FilesMkdirResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/files/download` {#operation-download-get-v1-files-download}

Download streams file headers, bytes, and per-file errors. A path that
leaves the jail through a symlink, or names a non-regular object, reports
a FileError with a error reason instead of reading through it.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesFilesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `paths` | query | Yes | Array of `string` | paths naming stable regular files or absent entries at canonical absolute paths yield contiguous results in request order, including duplicates. |


#### Responses

**200** Success

Read metadata followed by one continuous raw content part. JSON controls are separate parts with Content-Type application/json; binary content uses application/octet-stream. A final stream-error part carries the shared Error for a whole-stream failure. A clean closing delimiter with no stream-error is required for success; aborted connections are failures.

`multipart/mixed`: 

**Stream format:** Read metadata followed by one continuous raw content part. JSON controls are separate parts with Content-Type application/json; binary content uses application/octet-stream. A final stream-error part carries the shared Error for a whole-stream failure. A clean closing delimiter with no stream-error is required for success; aborted connections are failures.

<details>
<summary>Stream framing</summary>

```yaml
description: Read metadata followed by one continuous raw content part. JSON controls are separate parts with Content-Type application/json; binary content uses application/octet-stream. A final stream-error part carries the shared Error for a whole-stream failure. A clean closing delimiter with no stream-error is required for success; aborted connections are failures.
parts:
  metadata:
    contentType: application/json
    schema:
      $ref: "#/components/schemas/FileHeader"
  content:
    contentType: application/octet-stream
    schema:
      contentMediaType: application/octet-stream
  error:
    contentType: application/json
    schema:
      $ref: "#/components/schemas/FileError"
  stream-error:
    contentType: application/json
    schema:
      $ref: "#/components/schemas/Error"
repeatedItems: true
requiresClosingDelimiter: true
requiresTerminal: false
```

</details>

**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/files/move` {#operation-move-post-v1-files-move}

Move renames or moves a path. The destination is exact: an existing regular
file is replaced, the source basename is never appended, and a directory
destination is refused unchanged. A source or destination that leaves the jail
through a symlink refuses with FAILED_PRECONDITION and a error reason.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesFilesWrite`.

#### Request body

Required.

**Content-Type:** `application/json`

[`FilesMoveRequest`](#schema-FilesMoveRequest)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/FilesMoveRequest"
```

</details>



#### Responses

**200** Success

`application/json`: [`FilesMoveResponse`](#schema-FilesMoveResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/FilesMoveResponse"
```

</details>


**409** The path resolves through a symlink that does not resolve to a target inside the sandbox jail.

The path names an existing socket, FIFO, device, or other non-regular object where a regular file is required.

`from` names a regular file and `to` an existing directory, which `Move` refuses rather than moving the file into it.

`from` names a directory and `to` an existing path, which `Move` refuses rather than replacing.

`from` and `to` are on different filesystems, and `Move` never copies the contents to cross the boundary.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/files/stat` {#operation-stat-get-v1-files-stat}

Stat reads metadata for one path.
A metadata read has a different response from the file-content resource.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesFilesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `path` | query | Yes | `string` | path is the sandbox path. |


#### Responses

**200** Success

`application/json`: [`FilesStatResponse`](#schema-FilesStatResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/FilesStatResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/files/upload` {#operation-upload-post-v1-files-upload}

Upload streams one or more files. The first frame must be FileHeader.
Writing through a symlink that leaves the jail, or onto an existing
non-regular object, refuses with FAILED_PRECONDITION and a error reason.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesFilesWrite`.

#### Request body

Required.

Send ordered metadata/content pairs. Each metadata part is application/json; the following content part is application/octet-stream containing raw bytes, including an empty file. Content-Disposition is form-data with the corresponding part name. Content-Transfer-Encoding is refused. A clean closing multipart delimiter is required; truncated content never completes the current item. Malformed occurrences report invalidArgument with a BadRequest detail identifying the occurrence index and path.

**Content-Type:** `multipart/form-data`

**Stream format:** Send ordered metadata/content pairs. Each metadata part is application/json; the following content part is application/octet-stream containing raw bytes, including an empty file. Content-Disposition is form-data with the corresponding part name. Content-Transfer-Encoding is refused. A clean closing multipart delimiter is required; truncated content never completes the current item. Malformed occurrences report invalidArgument with a BadRequest detail identifying the occurrence index and path.

<details>
<summary>Stream framing</summary>

```yaml
description: Send ordered metadata/content pairs. Each metadata part is application/json; the following content part is application/octet-stream containing raw bytes, including an empty file. Content-Disposition is form-data with the corresponding part name. Content-Transfer-Encoding is refused. A clean closing multipart delimiter is required; truncated content never completes the current item. Malformed occurrences report invalidArgument with a BadRequest detail identifying the occurrence index and path.
partOrder:
  - metadata
  - content
repeatedItems: true
requiresClosingDelimiter: true
```

</details>


#### Responses

**200** Success

`application/json`: [`FilesUploadResponse`](#schema-FilesUploadResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/FilesUploadResponse"
```

</details>


**409** The path resolves through a symlink that does not resolve to a target inside the sandbox jail.

The path names an existing socket, FIFO, device, or other non-regular object where a regular file is required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## Processes

Run commands and interact with processes inside a sandbox.

| Method | Path | Operation |
| --- | --- | --- |
| GET | `/v1/processes` | [`listProcesses`](#operation-listProcesses-get-v1-processes) |
| POST | `/v1/processes` | [`createProcess`](#operation-createProcess-post-v1-processes) |
| POST | `/v1/processes/exec` | [`exec`](#operation-exec-post-v1-processes-exec) |
| GET | `/v1/processes/{process}` | [`getProcess`](#operation-getProcess-get-v1-processes-process) |
| GET | `/v1/processes/{process}/interact` | [`interact`](#operation-interact-get-v1-processes-process-interact) |
| GET | `/v1/processes/{process}/output` | [`readOutput`](#operation-readOutput-get-v1-processes-process-output) |
| POST | `/v1/processes/{process}/signal` | [`signal`](#operation-signal-post-v1-processes-process-signal) |

### GET `/v1/processes` {#operation-listProcesses-get-v1-processes}

ListProcesses lists processes in creation order.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesExec`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `pageSize` | query | No | `integer` | page_size is an optional page size. Omitted or zero uses the backend default. Unless the operation states otherwise, page-size limits and handling of larger requests are backend-specific; use the backend support guide. Continue with nextPageToken until it is empty. |
| `pageToken` | query | No | `string` | page_token is an optional continuation token. |
| `filter` | query | No | `string` | filter is comma-separated exact-match field=value terms; the portable fields are session and state. |


#### Responses

**200** Success

`application/json`: [`ListProcessesResponse`](#schema-ListProcessesResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ListProcessesResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/processes` {#operation-createProcess-post-v1-processes}

CreateProcess starts a durable interactive process and returns its resource name.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesExec`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |


#### Request body

Required.

**Content-Type:** `application/json`

[`CreateProcessRequest`](#schema-CreateProcessRequest)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/CreateProcessRequest"
```

</details>



#### Responses

**201** Success

`application/json`: [`Process`](#schema-Process)

<details>
<summary>201 response schema</summary>

```yaml
$ref: "#/components/schemas/Process"
```

</details>


**409** Creation conflicts with an existing resource (ALREADY_EXISTS).

The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/processes/exec` {#operation-exec-post-v1-processes-exec}

Exec runs one command to completion. A retry after transport failure may run it again.
Action at the collection root: it runs a command to completion instead of creating a
process a caller can address.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesExec`.

#### Request body

Required.

**Content-Type:** `application/json`

[`ExecRequest`](#schema-ExecRequest)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/ExecRequest"
```

</details>



#### Responses

**200** Success

`application/json`: [`ExecResponse`](#schema-ExecResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ExecResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/processes/{process}` {#operation-getProcess-get-v1-processes-process}

GetProcess reads one process by its endpoint-scoped resource name.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesExec`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `process` | path | Yes | `string` | The process id. |


#### Responses

**200** Success

`application/json`: [`Process`](#schema-Process)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Process"
```

</details>


**404** The target is absent or not visible within the caller's scope.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>404 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/processes/{process}/interact` {#operation-interact-get-v1-processes-process-interact}

Interact binds the complete process name from the path to the first Attach frame.
A conflicting name is INVALID_ARGUMENT before any attach.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer` or `sandboxWebSocketBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesExec`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `process` | path | Yes | `string` | The process id. |


#### Responses

**101** The connection becomes a sandboxes.v1 session: one InteractRequest per inbound text frame and one ProcessOutput per outbound one. Path parameters identify the target; the first attach frame must agree. Refusals after upgrade carry an error frame followed by a close code.

**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



#### Interactive connection

Use exactly one sandbox-scoped Authorization bearer credential or, when endpoint credentialTransports advertises webSocketSubprotocol, offer sandboxes.bearer.v1. followed by canonical unpadded base64url of that credential alongside sandboxes.v1. The decoded token is limited to 4096 bytes. The server echoes only sandboxes.v1; mixed or duplicate credentials, query credentials and cookie-only authentication are refused.

Open a WSS endpoint using its scoped credential. Browsers may send cookies; the server ignores them as authority and permits cross-origin explicit bearer authentication. Browser WebSocket errors do not expose pre-upgrade HTTP status or bodies; only post-upgrade Error frames are typed. No reconnect or input replay is automatic.

**Request frames:** [`InteractRequest`](#schema-InteractRequest). **Response frames:** [`ProcessOutput`](#schema-ProcessOutput).

**Reconnect:** Acquire a fresh scoped credential, open a new connection, and send attach with attach.resumeFrom set to the last received chunk.streamSequence. Resume replays output only; never replay stdin or process creation.

<details>
<summary>WebSocket protocol</summary>

```yaml
subprotocol: sandboxes.v1
firstRequestField: attach
resumeRequestField: attach.resumeFrom
resumeResponseField: chunk.streamSequence
bearerSubprotocolPrefix: sandboxes.bearer.v1.
maxBearerTokenBytes: 4096
requestFrame:
  $ref: "#/components/schemas/InteractRequest"
responseFrame:
  $ref: "#/components/schemas/ProcessOutput"
terminalFrame: The last text frame before a close carries the one error body, because a close reason holds at most 123 bytes.
errorFrame:
  $ref: "#/components/schemas/Error"
authentication: Use exactly one sandbox-scoped Authorization bearer credential or, when endpoint credentialTransports advertises webSocketSubprotocol, offer sandboxes.bearer.v1. followed by canonical unpadded base64url of that credential alongside sandboxes.v1. The decoded token is limited to 4096 bytes. The server echoes only sandboxes.v1; mixed or duplicate credentials, query credentials and cookie-only authentication are refused.
browser: Open a WSS endpoint using its scoped credential. Browsers may send cookies; the server ignores them as authority and permits cross-origin explicit bearer authentication. Browser WebSocket errors do not expose pre-upgrade HTTP status or bodies; only post-upgrade Error frames are typed. No reconnect or input replay is automatic.
reconnect: Acquire a fresh scoped credential, open a new connection, and send attach with attach.resumeFrom set to the last received chunk.streamSequence. Resume replays output only; never replay stdin or process creation.
preUpgradeStatus:
  "400": The request has an invalid subprotocol offer, mixed or malformed credentials, an oversized credential, or a query parameter.
  "401": The scoped credential is missing, invalid or expired.
  "403": The scoped credential does not grant this endpoint permission.
  "404": The scoped credential does not address this sandbox endpoint.
closeCode:
  "1000": The service ended the stream, or the caller cancelled it.
  "1002": The handshake negotiated no subprotocol.
  "1003": A frame arrived as binary; every frame is text.
  "1007": A frame was not the public JSON request message.
  "1008": The service refused with a 4xx error, or the endpoint credential expired.
  "1009": A frame exceeded the route’s frame limit.
  "1011": The service failed with a code the front door maps to a 5xx.
```

</details>

### GET `/v1/processes/{process}/output` {#operation-readOutput-get-v1-processes-process-output}

ReadOutput returns one page of retained output after a stream_sequence, the
same cursor Interact resumes from, so a REST client can read what a process printed.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesExec`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `process` | path | Yes | `string` | The process id. |
| `resumeFrom` | query | No | `integer or string` (int64) | resume_from returns output after this stream_sequence; 0 means from the start. Missing output after this point fails with FAILED_PRECONDITION, as for Attach. |
| `maxBytes` | query | No | `integer` | max_bytes is a soft bound on the page's data: whole chunks only, and the first retained chunk always fits so a caller always advances. 0 means the backend's default page. |


#### Responses

**200** Success

`application/json`: [`ReadOutputResponse`](#schema-ReadOutputResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ReadOutputResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/processes/{process}/signal` {#operation-signal-post-v1-processes-process-signal}

Signal sends one signal to an existing process.
Action: a signal is delivered to a running process and leaves no representation to
replace.

**Base URL:** `{endpoint}`. 

`endpoint`: The sandbox API base URL from Sandbox.core.endpoint.uri; preserve its path prefix when appending /v1 routes. HTTP endpoints require a scoped endpoint bearer; direct Unix sockets use socket access control.

**Authentication:** `sandboxBearer`.

Use a scoped endpoint credential, not the management token.

**Required permissions:** `sandboxesExec`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `process` | path | Yes | `string` | The process id. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  signal:
    not:
      enum:
        - unspecified
    description: signal is the signal to send.
    $ref: "#/components/schemas/ProcessSignal"
title: SignalRequest
required:
  - signal
additionalProperties: false
description: SignalRequest sends one signal.
```

</details>



#### Responses

**200** Success

`application/json`: [`SignalResponse`](#schema-SignalResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/SignalResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## Credentials

Exchange credentials for the authenticated owner.

| Method | Path | Operation |
| --- | --- | --- |
| POST | `/v1/identity/exchange` | [`exchangeDockerCredential`](#operation-exchangeDockerCredential-post-v1-identity-exchange) |

### POST `/v1/identity/exchange` {#operation-exchangeDockerCredential-post-v1-identity-exchange}

ExchangeDockerCredential consumes an identity token once for the authenticated owner.
The exchanged credential is stored for that owner; the response is empty.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `credentialsExchange`.

#### Request body

Required.

**Content-Type:** `application/json`

[`IdentityExchangeServiceExchangeDockerCredentialRequest`](#schema-IdentityExchangeServiceExchangeDockerCredentialRequest)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/IdentityExchangeServiceExchangeDockerCredentialRequest"
```

</details>



#### Responses

**200** Success

`application/json`: [`ExchangeCompositionCredentialResponse`](#schema-ExchangeCompositionCredentialResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ExchangeCompositionCredentialResponse"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## Images

Manage images used to create sandboxes.

| Method | Path | Operation |
| --- | --- | --- |
| GET | `/v1/images` | [`listImages`](#operation-listImages-get-v1-images) |
| POST | `/v1/images` | [`createImage`](#operation-createImage-post-v1-images) |
| GET | `/v1/images/{image}` | [`getImage`](#operation-getImage-get-v1-images-image) |
| DELETE | `/v1/images/{image}` | [`deleteImage`](#operation-deleteImage-delete-v1-images-image) |
| GET | `/v1/images/{image}/pull-spec` | [`getImagePullSpec`](#operation-getImagePullSpec-get-v1-images-image-pull-spec) |

### GET `/v1/images` {#operation-listImages-get-v1-images}

ListImages lists images in the caller's owner scope.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `imagesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `pageSize` | query | No | `integer` | page_size is an optional page size. Omitted or zero uses the backend default. Unless the operation states otherwise, page-size limits and handling of larger requests are backend-specific; use the backend support guide. Continue with nextPageToken until it is empty. |
| `pageToken` | query | No | `string` | page_token is an opaque continuation token. |
| `filter` | query | No | `string` | filter is comma-separated exact-match field=value terms; the portable fields are name, display_name, uid and status. |
| `orderBy` | query | No | `string` | order_by is a single order field with optional direction; the portable fields are created_at and name. |


#### Responses

**200** Success

`application/json`: [`ListImagesResponse`](#schema-ListImagesResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ListImagesResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/images` {#operation-createImage-post-v1-images}

CreateImage registers pushed OCI content or disk content derived from a snapshot or sandbox.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `imagesWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |


#### Request body

Required.

**Content-Type:** `application/json`

[`CreateImageRequest`](#schema-CreateImageRequest)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/CreateImageRequest"
```

</details>



#### Responses

**201** Success

Header `Cache-Control`: `string`. 

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Image`](#schema-Image) and `Object`

<details>
<summary>201 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Image"
  - not:
      properties:
        status:
          enum:
            - pending
            - waitingForPush
            - preparing
      required:
        - status
```

</details>


**202** Accepted. The resource is still progressing; read it or follow its events until completion.

Header `Cache-Control`: `string`. 

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Image`](#schema-Image) and `Object`

<details>
<summary>202 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Image"
  - properties:
      status:
        enum:
          - pending
          - waitingForPush
          - preparing
    required:
      - status
```

</details>


**409** Creation conflicts with an existing resource (ALREADY_EXISTS).

The Idempotency-Key was already used with a different payload.

Header `Cache-Control`: `string`. 

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

Header `Cache-Control`: `string`. 

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/images/{image}` {#operation-getImage-get-v1-images-image}

GetImage reads one image by its resource name.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `imagesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `image` | path | Yes | `string` | The image id. |
| `If-None-Match` | header | No | `string` | A quoted entity-tag or comma-separated list, compared weakly with the current resource etag. A match, including a weak form or wildcard *, answers 304 without a body after authorization. Malformed input is treated as no match. |


#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Image`](#schema-Image)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Image"
```

</details>


**304** The resource's etag equals the If-None-Match value, so this response carries no body.

**404** The target is absent or not visible within the caller's scope.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>404 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### DELETE `/v1/images/{image}` {#operation-deleteImage-delete-v1-images-image}

DeleteImage deletes an image; deleting an already absent image succeeds.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `imagesWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `image` | path | Yes | `string` | The image id. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |


#### Responses

**204** Success. The response carries no body.

**409** The target's preparation, capture, or import has not finished; deletion is refused until it settles to a terminal state.

An admitted consumer or live resource holds a content lease on the target; deletion is refused while any lease remains held.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/images/{image}/pull-spec` {#operation-getImagePullSpec-get-v1-images-image-pull-spec}

GetImagePullSpec returns short-lived pull material through registry-transfer support.
Action: the pull spec is derived material about the image, not a sub-resource of it.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `imagesPull`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `image` | path | Yes | `string` | The image id. |


#### Responses

**200** Success

Header `Cache-Control`: `string`. 

`application/json`: [`ImagePullSpec`](#schema-ImagePullSpec)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ImagePullSpec"
```

</details>


**501** This operation is not available on this service.

Header `Cache-Control`: `string`. 

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

Header `Cache-Control`: `string`. 

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## MCP gateways

Configure MCP gateways and authorize upstream servers.

| Method | Path | Operation |
| --- | --- | --- |
| GET | `/v1/mcp-servers/{mcpServer}/authorization` | [`getMcpAuthorization`](#operation-getMcpAuthorization-get-v1-mcp-servers-mcpServer-authorization) |
| POST | `/v1/mcp-servers/{mcpServer}/authorization/authorize` | [`authorizeMcpServer`](#operation-authorizeMcpServer-post-v1-mcp-servers-mcpServer-authorization-authorize) |
| GET | `/v1/sandboxes/{sandbox}/mcp-gateway` | [`getMcpGateway`](#operation-getMcpGateway-get-v1-sandboxes-sandbox-mcp-gateway) |
| POST | `/v1/sandboxes/{sandbox}/mcp-gateway/servers` | [`addMcpGatewayServer`](#operation-addMcpGatewayServer-post-v1-sandboxes-sandbox-mcp-gateway-servers) |
| POST | `/v1/sandboxes/{sandbox}/mcp-gateway/start` | [`startMcpGateway`](#operation-startMcpGateway-post-v1-sandboxes-sandbox-mcp-gateway-start) |
| POST | `/v1/sandboxes/{sandbox}/mcp-gateway/stop` | [`stopMcpGateway`](#operation-stopMcpGateway-post-v1-sandboxes-sandbox-mcp-gateway-stop) |

### GET `/v1/mcp-servers/{mcpServer}/authorization` {#operation-getMcpAuthorization-get-v1-mcp-servers-mcpServer-authorization}

GetMcpAuthorization reads the owner-scoped flow, including its pending URL, under mcp.write.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `mcpWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `mcpServer` | path | Yes | `string` | The mcpServer id. |
| `If-None-Match` | header | No | `string` | A quoted entity-tag or comma-separated list, compared weakly with the current resource etag. A match, including a weak form or wildcard *, answers 304 without a body after authorization. Malformed input is treated as no match. |


#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`McpAuthorization`](#schema-McpAuthorization)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/McpAuthorization"
```

</details>


**304** The resource's etag equals the If-None-Match value, so this response carries no body.

**404** The target is absent or not visible within the caller's scope.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>404 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/mcp-servers/{mcpServer}/authorization/authorize` {#operation-authorizeMcpServer-post-v1-mcp-servers-mcpServer-authorization-authorize}

AuthorizeMcpServer authorizes credentials for one named upstream.
Action: authorization runs an interactive grant whose result is not the request body.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `mcpWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `mcpServer` | path | Yes | `string` | The mcpServer id. |
| `If-Match` | header | No | `string` | Optional current authorization etag. A supplied stale value answers 412 with failedPrecondition; omission permits create or convergence, including forced reauthorization. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  remoteUrl:
    type:
      - string
      - "null"
    format: uri
    description: remote_url binds authorization to a concrete upstream URL when supplied.
  forceReauth:
    type: boolean
    description: |-
      force_reauth starts a fresh flow while existing credentials serve until completion.
      Repeating a pending flow converges without reminting.
title: AuthorizeMcpServerRequest
additionalProperties: false
description: AuthorizeMcpServerRequest starts or refreshes MCP upstream authorization.
```

</details>



#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`McpAuthorization`](#schema-McpAuthorization) and `Object`

<details>
<summary>200 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/McpAuthorization"
  - not:
      properties:
        status:
          enum:
            - pending
      required:
        - status
```

</details>


**202** Accepted. The resource is still progressing; read it or follow its events until completion.

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`McpAuthorization`](#schema-McpAuthorization) and `Object`

<details>
<summary>202 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/McpAuthorization"
  - properties:
      status:
        enum:
          - pending
    required:
      - status
```

</details>


**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/sandboxes/{sandbox}/mcp-gateway` {#operation-getMcpGateway-get-v1-sandboxes-sandbox-mcp-gateway}

GetMcpGateway reads gateway state for the sandbox.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `mcpRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |


#### Responses

**200** Success

`application/json`: [`McpGateway`](#schema-McpGateway)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/McpGateway"
```

</details>


**404** The target is absent or not visible within the caller's scope.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>404 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes/{sandbox}/mcp-gateway/servers` {#operation-addMcpGatewayServer-post-v1-sandboxes-sandbox-mcp-gateway-servers}

AddMcpGatewayServer adds one server to an existing ready gateway.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `mcpWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  server:
    type: string
    minLength: 1
    description: |-
      server is the external configured key; a supplied registration must define this same key.
      A mismatch refuses FAILED_PRECONDITION before changing the gateway.
title: AddMcpGatewayServerRequest
required:
  - server
additionalProperties: false
description: AddMcpGatewayServerRequest adds one server to an existing ready gateway.
```

</details>



#### Responses

**200** Success

`application/json`: [`AddMcpGatewayServerResponse`](#schema-AddMcpGatewayServerResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/AddMcpGatewayServerResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes/{sandbox}/mcp-gateway/start` {#operation-startMcpGateway-post-v1-sandboxes-sandbox-mcp-gateway-start}

StartMcpGateway ensures a gateway exists for the sandbox.
Action: a lifecycle transition is not a representation a caller can PUT.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `mcpWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  servers:
    type: array
    items:
      type: string
      minLength: 1
    description: servers are initial server names for a backend-minted gateway.
  static:
    type: boolean
    description: static pins the requested server set and closes gateway-side discovery.
  gatewayUrl:
    type:
      - string
      - "null"
    format: uri
    description: gateway_url attaches read-only to a pre-existing shareable gateway when supported.
title: StartMcpGatewayRequest
additionalProperties: false
description: |-
  StartMcpGatewayRequest ensures a sandbox gateway exists.
  servers must be empty when gateway_url is set
```

</details>



#### Responses

**200** Success

`application/json`: [`McpGateway`](#schema-McpGateway) and `Object`

<details>
<summary>200 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/McpGateway"
  - not:
      properties:
        state:
          enum:
            - provisioning
      required:
        - state
```

</details>


**202** Accepted. The resource is still progressing; read it or follow its events until completion.

`application/json`: [`McpGateway`](#schema-McpGateway) and `Object`

<details>
<summary>202 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/McpGateway"
  - properties:
      state:
        enum:
          - provisioning
    required:
      - state
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes/{sandbox}/mcp-gateway/stop` {#operation-stopMcpGateway-post-v1-sandboxes-sandbox-mcp-gateway-stop}

StopMcpGateway stops or detaches the gateway and succeeds when already absent.
Action: a lifecycle transition is not a representation a caller can PUT.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `mcpWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |


#### Responses

**200** Success

`application/json`: [`StopMcpGatewayResponse`](#schema-StopMcpGatewayResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/StopMcpGatewayResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## Network policies

Read enforced outbound network policies and their decision logs.

| Method | Path | Operation |
| --- | --- | --- |
| GET | `/v1/network-policies` | [`getNetworkPoliciesInOwnerScope`](#operation-getNetworkPoliciesInOwnerScope-get-v1-network-policies) |
| GET | `/v1/policy-logs` | [`listPolicyLogEntriesInOwnerScope`](#operation-listPolicyLogEntriesInOwnerScope-get-v1-policy-logs) |
| GET | `/v1/sandboxes/{sandbox}/network-policies` | [`getNetworkPolicies`](#operation-getNetworkPolicies-get-v1-sandboxes-sandbox-network-policies) |
| GET | `/v1/sandboxes/{sandbox}/policy-logs` | [`listPolicyLogEntries`](#operation-listPolicyLogEntries-get-v1-sandboxes-sandbox-policy-logs) |

### GET `/v1/network-policies` {#operation-getNetworkPoliciesInOwnerScope-get-v1-network-policies}

GetNetworkPolicies returns effective and exact views of the same installed generation.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `networkPoliciesRead`.

#### Responses

**200** Success

`application/json`: [`GetNetworkPoliciesResponse`](#schema-GetNetworkPoliciesResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/GetNetworkPoliciesResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/policy-logs` {#operation-listPolicyLogEntriesInOwnerScope-get-v1-policy-logs}

ListPolicyLogEntries lists observed policy decisions.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `networkPoliciesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `pageSize` | query | No | `integer` | page_size is an optional page size. Omitted or zero uses the backend default. Unless the operation states otherwise, page-size limits and handling of larger requests are backend-specific; use the backend support guide. Continue with nextPageToken until it is empty. |
| `pageToken` | query | No | `string` | page_token is an opaque continuation token. |
| `filter` | query | No | `string` | filter is comma-separated exact-match field=value terms; the portable field is domain. |
| `since` | query | No | `string` (date-time) | An RFC 3339 timestamp, for example 2026-01-01T12:00:00Z. Fractional seconds may contain up to nine digits. The supported range is 0001-01-01T00:00:00Z through 9999-12-31T23:59:59.999999999Z. |
| `until` | query | No | `string` (date-time) | An RFC 3339 timestamp, for example 2026-01-01T12:00:00Z. Fractional seconds may contain up to nine digits. The supported range is 0001-01-01T00:00:00Z through 9999-12-31T23:59:59.999999999Z. |


#### Responses

**200** Success

`application/json`: [`ListPolicyLogEntriesResponse`](#schema-ListPolicyLogEntriesResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ListPolicyLogEntriesResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/sandboxes/{sandbox}/network-policies` {#operation-getNetworkPolicies-get-v1-sandboxes-sandbox-network-policies}

GetNetworkPolicies returns effective and exact views of the same installed generation.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `networkPoliciesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |


#### Responses

**200** Success

`application/json`: [`GetNetworkPoliciesResponse`](#schema-GetNetworkPoliciesResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/GetNetworkPoliciesResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/sandboxes/{sandbox}/policy-logs` {#operation-listPolicyLogEntries-get-v1-sandboxes-sandbox-policy-logs}

ListPolicyLogEntries lists observed policy decisions.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `networkPoliciesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `pageSize` | query | No | `integer` | page_size is an optional page size. Omitted or zero uses the backend default. Unless the operation states otherwise, page-size limits and handling of larger requests are backend-specific; use the backend support guide. Continue with nextPageToken until it is empty. |
| `pageToken` | query | No | `string` | page_token is an opaque continuation token. |
| `filter` | query | No | `string` | filter is comma-separated exact-match field=value terms; the portable field is domain. |
| `since` | query | No | `string` (date-time) | An RFC 3339 timestamp, for example 2026-01-01T12:00:00Z. Fractional seconds may contain up to nine digits. The supported range is 0001-01-01T00:00:00Z through 9999-12-31T23:59:59.999999999Z. |
| `until` | query | No | `string` (date-time) | An RFC 3339 timestamp, for example 2026-01-01T12:00:00Z. Fractional seconds may contain up to nine digits. The supported range is 0001-01-01T00:00:00Z through 9999-12-31T23:59:59.999999999Z. |


#### Responses

**200** Success

`application/json`: [`ListPolicyLogEntriesResponse`](#schema-ListPolicyLogEntriesResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ListPolicyLogEntriesResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## Sandboxes

Create sandboxes and manage their lifecycle, ports, and SSH access.

| Method | Path | Operation |
| --- | --- | --- |
| GET | `/v1/sandboxes` | [`listSandboxes`](#operation-listSandboxes-get-v1-sandboxes) |
| POST | `/v1/sandboxes` | [`createSandbox`](#operation-createSandbox-post-v1-sandboxes) |
| GET | `/v1/sandboxes/{sandbox}` | [`getSandbox`](#operation-getSandbox-get-v1-sandboxes-sandbox) |
| PATCH | `/v1/sandboxes/{sandbox}` | [`updateSandbox`](#operation-updateSandbox-patch-v1-sandboxes-sandbox) |
| DELETE | `/v1/sandboxes/{sandbox}` | [`deleteSandbox`](#operation-deleteSandbox-delete-v1-sandboxes-sandbox) |
| POST | `/v1/sandboxes/{sandbox}/endpoint-credentials` | [`createEndpointCredential`](#operation-createEndpointCredential-post-v1-sandboxes-sandbox-endpoint-credentials) |
| GET | `/v1/sandboxes/{sandbox}/ports` | [`listPorts`](#operation-listPorts-get-v1-sandboxes-sandbox-ports) |
| POST | `/v1/sandboxes/{sandbox}/ports` | [`createPort`](#operation-createPort-post-v1-sandboxes-sandbox-ports) |
| GET | `/v1/sandboxes/{sandbox}/ports/{port}` | [`getPort`](#operation-getPort-get-v1-sandboxes-sandbox-ports-port) |
| DELETE | `/v1/sandboxes/{sandbox}/ports/{port}` | [`deletePort`](#operation-deletePort-delete-v1-sandboxes-sandbox-ports-port) |
| POST | `/v1/sandboxes/{sandbox}/renew-timeout` | [`renewSandboxTimeout`](#operation-renewSandboxTimeout-post-v1-sandboxes-sandbox-renew-timeout) |
| POST | `/v1/sandboxes/{sandbox}/ssh-certificates` | [`issueSSHCert`](#operation-issueSSHCert-post-v1-sandboxes-sandbox-ssh-certificates) |
| POST | `/v1/sandboxes/{sandbox}/start` | [`startSandbox`](#operation-startSandbox-post-v1-sandboxes-sandbox-start) |
| POST | `/v1/sandboxes/{sandbox}/stop` | [`stopSandbox`](#operation-stopSandbox-post-v1-sandboxes-sandbox-stop) |

### GET `/v1/sandboxes` {#operation-listSandboxes-get-v1-sandboxes}

ListSandboxes lists sandboxes in the caller's owner scope.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `pageSize` | query | No | `integer` | page_size is an optional page size. Omitted or zero uses the backend default. Unless the operation states otherwise, page-size limits and handling of larger requests are backend-specific; use the backend support guide. Continue with nextPageToken until it is empty. |
| `pageToken` | query | No | `string` | page_token is an opaque continuation token. |
| `filter` | query | No | `string` | filter is comma-separated exact-match field=value terms; the portable fields are name, display_name, uid, status and agent. |
| `orderBy` | query | No | `string` | order_by is a single order field with optional direction; the portable fields are created_at and name. |


#### Responses

**200** Success

`application/json`: [`ListSandboxesResponse`](#schema-ListSandboxesResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ListSandboxesResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes` {#operation-createSandbox-post-v1-sandboxes}

CreateSandbox admits a sandbox. Unknown fields at any level are rejected before
replay lookup, reservation, policy binding, or workload effects.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesCreate`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |


#### Request body

Required.

**Content-Type:** `application/json`

[`CreateSandboxRequest`](#schema-CreateSandboxRequest)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/CreateSandboxRequest"
```

</details>



#### Responses

**201** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox) and `Object`

<details>
<summary>201 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Sandbox"
  - not:
      properties:
        core:
          properties:
            status:
              enum:
                - creating
                - starting
          required:
            - status
      required:
        - core
```

</details>


**202** Accepted. The resource is still progressing; read it or follow its events until completion.

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox) and `Object`

<details>
<summary>202 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Sandbox"
  - properties:
      core:
        properties:
          status:
            enum:
              - creating
              - starting
        required:
          - status
    required:
      - core
```

</details>


**409** Creation conflicts with an existing resource (ALREADY_EXISTS).

Kit admission was refused; an available settled report describes the rejected fields.

The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/sandboxes/{sandbox}` {#operation-getSandbox-get-v1-sandboxes-sandbox}

GetSandbox reads one sandbox by its resource name.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `If-None-Match` | header | No | `string` | A quoted entity-tag or comma-separated list, compared weakly with the current resource etag. A match, including a weak form or wildcard *, answers 304 without a body after authorization. Malformed input is treated as no match. |


#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Sandbox"
```

</details>


**304** The resource's etag equals the If-None-Match value, so this response carries no body.

**404** The target is absent or not visible within the caller's scope.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>404 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### PATCH `/v1/sandboxes/{sandbox}` {#operation-updateSandbox-patch-v1-sandboxes-sandbox}

UpdateSandbox changes only the display label under the observed etag.
Its resource name, backing identity and child addresses remain unchanged.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesRename`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `Idempotency-Key` | header | Yes | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  displayName:
    type:
      - string
      - "null"
    maxLength: 64
    minLength: 1
    pattern: ^[a-zA-Z0-9_-]+$
    description: display_name must be present and retains the existing scoped label rules.
title: UpdateSandboxRequest
required:
  - displayName
additionalProperties: false
description: UpdateSandboxRequest changes only the display label under the observed etag.
```

</details>



#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Sandbox"
```

</details>


**409** The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### DELETE `/v1/sandboxes/{sandbox}` {#operation-deleteSandbox-delete-v1-sandboxes-sandbox}

DeleteSandbox returns the deleting resource while work remains, or an empty message when absent.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesDelete`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |
| `force` | query | No | `boolean` | force requests deletion even when the backend would otherwise refuse the state. |


#### Responses

**202** Accepted. The resource is still progressing; read it or follow its events until completion.

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox) and `Object`

<details>
<summary>202 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Sandbox"
  - properties:
      core:
        properties:
          status:
            enum:
              - deleting
        required:
          - status
    required:
      - core
```

</details>


**204** Success. The response carries no body.

**409** The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes/{sandbox}/endpoint-credentials` {#operation-createEndpointCredential-post-v1-sandboxes-sandbox-endpoint-credentials}

CreateEndpointCredential grants short-lived access to one sandbox endpoint.
Issuance checks current authority and never stores a token in the replay ledger.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesCredential`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  permissions:
    type: array
    items:
      $ref: "#/components/schemas/Permission"
      enum:
        - sandboxesExec
        - sandboxesFilesRead
        - sandboxesFilesWrite
    minItems: 1
    uniqueItems: true
    description: "Allowed values: sandboxesExec, sandboxesFilesRead, sandboxesFilesWrite."
  ttl:
    description: ttl defaults to five minutes; supplied durations are clamped to one second through five minutes.
    $ref: "#/components/schemas/Duration"
title: CreateEndpointCredentialRequest
required:
  - permissions
additionalProperties: false
description: CreateEndpointCredentialRequest requests only endpoint permissions the caller holds.
```

</details>



#### Responses

**200** Success

Header `Cache-Control`: `string`. 

`application/json`: [`EndpointCredential`](#schema-EndpointCredential)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/EndpointCredential"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

Header `Cache-Control`: `string`. 

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/sandboxes/{sandbox}/ports` {#operation-listPorts-get-v1-sandboxes-sandbox-ports}

ListPorts lists currently published ports for one sandbox.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |


#### Responses

**200** Success

`application/json`: [`ListPortsResponse`](#schema-ListPortsResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ListPortsResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes/{sandbox}/ports` {#operation-createPort-post-v1-sandboxes-sandbox-ports}

CreatePort makes one sandbox port reachable outside the sandbox.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesPorts`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  port:
    description: port is the requested publication.
    $ref: "#/components/schemas/CreatePortRequestPortInput"
title: CreatePortRequest
required:
  - port
additionalProperties: false
description: CreatePortRequest publishes one sandbox port.
```

</details>



#### Responses

**201** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Port`](#schema-Port)

<details>
<summary>201 response schema</summary>

```yaml
$ref: "#/components/schemas/Port"
```

</details>


**409** Creation conflicts with an existing resource (ALREADY_EXISTS).

The sandbox already publishes this port number.

The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/sandboxes/{sandbox}/ports/{port}` {#operation-getPort-get-v1-sandboxes-sandbox-ports-port}

GetPort reads one published port by its resource name.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `port` | path | Yes | `string` | The port id. |
| `If-None-Match` | header | No | `string` | A quoted entity-tag or comma-separated list, compared weakly with the current resource etag. A match, including a weak form or wildcard *, answers 304 without a body after authorization. Malformed input is treated as no match. |


#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Port`](#schema-Port)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Port"
```

</details>


**304** The resource's etag equals the If-None-Match value, so this response carries no body.

**404** No publication exists for this port resource name.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>404 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### DELETE `/v1/sandboxes/{sandbox}/ports/{port}` {#operation-deletePort-delete-v1-sandboxes-sandbox-ports-port}

DeletePort withdraws one published port. A name that is not published
succeeds, so a retried delete reports the same result as the first.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesPorts`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `port` | path | Yes | `string` | The port id. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |


#### Responses

**204** Success. The response carries no body.

**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes/{sandbox}/renew-timeout` {#operation-renewSandboxTimeout-post-v1-sandboxes-sandbox-renew-timeout}

RenewSandboxTimeout sets a TTL when timeout renewal is supported. The TTL starts at acceptance, and the response returns the new expires_at with the sandbox's timeout configuration preserved.
Action: the TTL restarts at acceptance, so the result is not a value the caller could
PUT.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesLifecycle`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  timeout:
    description: timeout is the new TTL relative to server acceptance time.
    $ref: "#/components/schemas/Duration"
title: RenewSandboxTimeoutRequest
required:
  - timeout
additionalProperties: false
description: RenewSandboxTimeoutRequest sets a new remaining TTL.
```

</details>



#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Sandbox"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes/{sandbox}/ssh-certificates` {#operation-issueSSHCert-post-v1-sandboxes-sandbox-ssh-certificates}

IssueSSHCert signs a caller-held public key for sandbox SSH access.
A create: each call mints one certificate under the sandbox's certificate collection.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesSsh`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  publicKey:
    type: string
    minLength: 1
    description: public_key is the caller-held SSH public key.
  ttl:
    description: ttl is the requested certificate lifetime.
    $ref: "#/components/schemas/Duration"
title: IssueSSHCertRequest
required:
  - publicKey
additionalProperties: false
description: IssueSSHCertRequest asks the backend to sign a caller-held public key.
```

</details>



#### Responses

**201** Success

Header `Cache-Control`: `string`. 

`application/json`: [`IssueSSHCertResponse`](#schema-IssueSSHCertResponse)

<details>
<summary>201 response schema</summary>

```yaml
$ref: "#/components/schemas/IssueSSHCertResponse"
```

</details>


**409** Creation conflicts with an existing resource (ALREADY_EXISTS).

Header `Cache-Control`: `string`. 

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**501** This operation is not available on this service.

Header `Cache-Control`: `string`. 

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

Header `Cache-Control`: `string`. 

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes/{sandbox}/start` {#operation-startSandbox-post-v1-sandboxes-sandbox-start}

StartSandbox starts a stopped sandbox; a backend that preserved memory resumes it.
Action: a lifecycle transition is not a representation a caller can PUT.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesLifecycle`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |


#### Request body

Optional.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties: {}
title: StartSandboxRequest
additionalProperties: false
description: StartSandboxRequest starts or resumes a sandbox.
```

</details>



#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox) and `Object`

<details>
<summary>200 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Sandbox"
  - not:
      properties:
        core:
          properties:
            status:
              enum:
                - starting
          required:
            - status
      required:
        - core
```

</details>


**202** Accepted. The resource is still progressing; read it or follow its events until completion.

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox) and `Object`

<details>
<summary>202 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Sandbox"
  - properties:
      core:
        properties:
          status:
            enum:
              - starting
        required:
          - status
    required:
      - core
```

</details>


**409** The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/sandboxes/{sandbox}/stop` {#operation-stopSandbox-post-v1-sandboxes-sandbox-stop}

StopSandbox stops a running sandbox; a backend advertising memory preservation hibernates it.
Action: a lifecycle transition is not a representation a caller can PUT.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `sandboxesLifecycle`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |


#### Request body

Optional.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties: {}
title: StopSandboxRequest
additionalProperties: false
description: StopSandboxRequest stops a sandbox; a backend advertising memory preservation hibernates it.
```

</details>



#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox) and `Object`

<details>
<summary>200 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Sandbox"
  - not:
      properties:
        core:
          properties:
            status:
              enum:
                - stopping
          required:
            - status
      required:
        - core
```

</details>


**202** Accepted. The resource is still progressing; read it or follow its events until completion.

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox) and `Object`

<details>
<summary>202 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Sandbox"
  - properties:
      core:
        properties:
          status:
            enum:
              - stopping
        required:
          - status
    required:
      - core
```

</details>


**409** The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## Secrets

Store credentials and read their metadata without returning secret material.

| Method | Path | Operation |
| --- | --- | --- |
| GET | `/v1/secrets` | [`listSecrets`](#operation-listSecrets-get-v1-secrets) |
| POST | `/v1/secrets` | [`createSecret`](#operation-createSecret-post-v1-secrets) |
| GET | `/v1/secrets/{secret}` | [`getSecret`](#operation-getSecret-get-v1-secrets-secret) |
| PUT | `/v1/secrets/{secret}` | [`updateSecret`](#operation-updateSecret-put-v1-secrets-secret) |
| DELETE | `/v1/secrets/{secret}` | [`deleteSecret`](#operation-deleteSecret-delete-v1-secrets-secret) |

### GET `/v1/secrets` {#operation-listSecrets-get-v1-secrets}

ListSecrets lists secret metadata in the caller's owner scope.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `secretsRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `pageSize` | query | No | `integer` | page_size is an optional page size. Omitted or zero uses the backend default. Unless the operation states otherwise, page-size limits and handling of larger requests are backend-specific; use the backend support guide. Continue with nextPageToken until it is empty. |
| `pageToken` | query | No | `string` | page_token is an opaque continuation token. |
| `filter` | query | No | `string` | filter is comma-separated exact-match field=value terms; the portable fields are name, display_name and uid. |
| `orderBy` | query | No | `string` | order_by is a single order field with optional direction; the portable fields are created_at and name. |


#### Responses

**200** Success

`application/json`: [`ListSecretsResponse`](#schema-ListSecretsResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ListSecretsResponse"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/secrets` {#operation-createSecret-post-v1-secrets}

CreateSecret stores new secret material.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `secretsWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |


#### Request body

Required.

**Content-Type:** `application/json`

[`CreateSecretRequest`](#schema-CreateSecretRequest)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/CreateSecretRequest"
```

</details>



#### Responses

**201** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Secret`](#schema-Secret)

<details>
<summary>201 response schema</summary>

```yaml
$ref: "#/components/schemas/Secret"
```

</details>


**409** Creation conflicts with an existing resource (ALREADY_EXISTS).

The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/secrets/{secret}` {#operation-getSecret-get-v1-secrets-secret}

GetSecret reads one secret's metadata; material is never returned.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `secretsRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `secret` | path | Yes | `string` | The secret id. |
| `If-None-Match` | header | No | `string` | A quoted entity-tag or comma-separated list, compared weakly with the current resource etag. A match, including a weak form or wildcard *, answers 304 without a body after authorization. Malformed input is treated as no match. |


#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Secret`](#schema-Secret)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Secret"
```

</details>


**304** The resource's etag equals the If-None-Match value, so this response carries no body.

**404** The target is absent or not visible within the caller's scope.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>404 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### PUT `/v1/secrets/{secret}` {#operation-updateSecret-put-v1-secrets-secret}

UpdateSecret replaces secret metadata and material.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `secretsWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `secret` | path | Yes | `string` | The secret id. |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |


#### Request body

Required.

**Content-Type:** `application/json`

[`UpdateSecretBody`](#schema-UpdateSecretBody)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/UpdateSecretBody"
```

</details>



#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Secret`](#schema-Secret)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Secret"
```

</details>


**409** The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### DELETE `/v1/secrets/{secret}` {#operation-deleteSecret-delete-v1-secrets-secret}

DeleteSecret deletes a secret; deleting an already absent secret succeeds.
Its response carries the credential fence receipt, so it answers 200 with a
body rather than the 204 an empty-response delete declares.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `secretsWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `secret` | path | Yes | `string` | The secret id. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |


#### Responses

**200** Success

`application/json`: [`DeleteSecretResponse`](#schema-DeleteSecretResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/DeleteSecretResponse"
```

</details>


**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## Snapshots

Capture sandbox state and restore it into a new sandbox.

| Method | Path | Operation |
| --- | --- | --- |
| POST | `/v1/sandboxes/{sandbox}/snapshots` | [`createSnapshot`](#operation-createSnapshot-post-v1-sandboxes-sandbox-snapshots) |
| GET | `/v1/snapshots` | [`listSnapshots`](#operation-listSnapshots-get-v1-snapshots) |
| GET | `/v1/snapshots/{snapshot}` | [`getSnapshot`](#operation-getSnapshot-get-v1-snapshots-snapshot) |
| DELETE | `/v1/snapshots/{snapshot}` | [`deleteSnapshot`](#operation-deleteSnapshot-delete-v1-snapshots-snapshot) |
| POST | `/v1/snapshots/{snapshot}/restore` | [`restoreSnapshot`](#operation-restoreSnapshot-post-v1-snapshots-snapshot-restore) |

### POST `/v1/sandboxes/{sandbox}/snapshots` {#operation-createSnapshot-post-v1-sandboxes-sandbox-snapshots}

CreateSnapshot checkpoints a running sandbox; a stopped source is rejected.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `snapshotsWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | path | Yes | `string` | The sandbox id. |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  displayName:
    type: string
    not:
      type: string
      enum:
        - .
        - ..
    minLength: 1
    pattern: ^[^/]+$
    description: display_name is a scoped label; the backend assigns the immutable resource ID.
  description:
    type: string
    description: description is caller-supplied text.
  captureMode:
    description: |-
      capture_mode selects disk-only or memory plus disk capture. A value this
      enum does not define is rejected, never coerced to a capture.
    $ref: "#/components/schemas/CaptureMode"
title: CreateSnapshotRequest
required:
  - displayName
additionalProperties: false
description: CreateSnapshotRequest captures one running sandbox.
```

</details>



#### Responses

**201** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Snapshot`](#schema-Snapshot) and `Object`

<details>
<summary>201 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Snapshot"
  - not:
      properties:
        status:
          enum:
            - creating
      required:
        - status
```

</details>


**202** Accepted. The resource is still progressing; read it or follow its events until completion.

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Snapshot`](#schema-Snapshot) and `Object`

<details>
<summary>202 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Snapshot"
  - properties:
      status:
        enum:
          - creating
    required:
      - status
```

</details>


**409** Creation conflicts with an existing resource (ALREADY_EXISTS).

The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/snapshots` {#operation-listSnapshots-get-v1-snapshots}

ListSnapshots lists snapshots in the caller's owner scope.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `snapshotsRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox` | query | No | `string` | (OPTIONAL) sandbox is the complete resource name, in the form sandboxes/{sandbox}. |
| `pageSize` | query | No | `integer` | page_size is an optional page size. Omitted or zero uses the backend default. Unless the operation states otherwise, page-size limits and handling of larger requests are backend-specific; use the backend support guide. Continue with nextPageToken until it is empty. |
| `pageToken` | query | No | `string` | page_token is an opaque continuation token. |
| `filter` | query | No | `string` | filter is comma-separated exact-match field=value terms; the portable fields are name, display_name and uid. |
| `orderBy` | query | No | `string` | order_by is a single order field with optional direction; the portable fields are created_at and name. |


#### Responses

**200** Success

`application/json`: [`ListSnapshotsResponse`](#schema-ListSnapshotsResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ListSnapshotsResponse"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/snapshots/{snapshot}` {#operation-getSnapshot-get-v1-snapshots-snapshot}

GetSnapshot reads one snapshot by its resource name.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `snapshotsRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `snapshot` | path | Yes | `string` | The snapshot id. |
| `If-None-Match` | header | No | `string` | A quoted entity-tag or comma-separated list, compared weakly with the current resource etag. A match, including a weak form or wildcard *, answers 304 without a body after authorization. Malformed input is treated as no match. |


#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Snapshot`](#schema-Snapshot)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Snapshot"
```

</details>


**304** The resource's etag equals the If-None-Match value, so this response carries no body.

**404** The target is absent or not visible within the caller's scope.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>404 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### DELETE `/v1/snapshots/{snapshot}` {#operation-deleteSnapshot-delete-v1-snapshots-snapshot}

DeleteSnapshot deletes a snapshot; deleting an already absent snapshot succeeds.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `snapshotsWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `snapshot` | path | Yes | `string` | The snapshot id. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |


#### Responses

**204** Success. The response carries no body.

**409** The target's preparation, capture, or import has not finished; deletion is refused until it settles to a terminal state.

An admitted consumer or live resource holds a content lease on the target; deletion is refused while any lease remains held.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/snapshots/{snapshot}/restore` {#operation-restoreSnapshot-post-v1-snapshots-snapshot-restore}

RestoreSnapshot forks a new sandbox from a snapshot.
Action: a restore creates a sandbox, so the snapshot in the path is the source and not
the created resource.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `snapshotsRead`, `sandboxesCreate`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `snapshot` | path | Yes | `string` | The snapshot id. |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |


#### Request body

Required.

**Content-Type:** `application/json`

`object`

<details>
<summary>Request schema</summary>

```yaml
type: object
properties:
  displayName:
    type: string
    maxLength: 64
    pattern: ^[a-zA-Z0-9_-]*$
    description: (OPTIONAL) display_name is a scoped label; the backend assigns the immutable resource ID.
title: RestoreSnapshotRequest
additionalProperties: false
description: RestoreSnapshotRequest creates a new sandbox from one snapshot.
```

</details>



#### Responses

**201** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox) and `Object`

<details>
<summary>201 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Sandbox"
  - not:
      properties:
        core:
          properties:
            status:
              enum:
                - creating
                - starting
          required:
            - status
      required:
        - core
```

</details>


**202** Accepted. The resource is still progressing; read it or follow its events until completion.

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Sandbox`](#schema-Sandbox) and `Object`

<details>
<summary>202 response schema</summary>

```yaml
allOf:
  - $ref: "#/components/schemas/Sandbox"
  - properties:
      core:
        properties:
          status:
            enum:
              - creating
              - starting
        required:
          - status
    required:
      - core
```

</details>


**409** Creation conflicts with an existing resource (ALREADY_EXISTS).

The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## Volumes

Manage persistent storage and attach it to sandboxes.

| Method | Path | Operation |
| --- | --- | --- |
| GET | `/v1/volumes` | [`listVolumes`](#operation-listVolumes-get-v1-volumes) |
| POST | `/v1/volumes` | [`createVolume`](#operation-createVolume-post-v1-volumes) |
| GET | `/v1/volumes/{volume}` | [`getVolume`](#operation-getVolume-get-v1-volumes-volume) |
| DELETE | `/v1/volumes/{volume}` | [`deleteVolume`](#operation-deleteVolume-delete-v1-volumes-volume) |

### GET `/v1/volumes` {#operation-listVolumes-get-v1-volumes}

ListVolumes lists volumes in the caller's owner scope.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `volumesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `pageSize` | query | No | `integer` | page_size is an optional page size. Omitted or zero uses the backend default. Unless the operation states otherwise, page-size limits and handling of larger requests are backend-specific; use the backend support guide. Continue with nextPageToken until it is empty. |
| `pageToken` | query | No | `string` | page_token is an opaque continuation token. |
| `filter` | query | No | `string` | filter is comma-separated exact-match field=value terms; the portable fields are name, display_name and uid. |
| `orderBy` | query | No | `string` | order_by is a single order field with optional direction; the portable fields are created_at and name. |


#### Responses

**200** Success

`application/json`: [`ListVolumesResponse`](#schema-ListVolumesResponse)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/ListVolumesResponse"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### POST `/v1/volumes` {#operation-createVolume-post-v1-volumes}

CreateVolume creates a persistent volume.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `volumesWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `Idempotency-Key` | header | No | `string` | Replay key. Reusing it with a different payload fails with failedPrecondition; accepted keys are retained for at least 24 hours. |


#### Request body

Required.

**Content-Type:** `application/json`

[`CreateVolumeRequest`](#schema-CreateVolumeRequest)

<details>
<summary>Request schema</summary>

```yaml
$ref: "#/components/schemas/CreateVolumeRequest"
```

</details>



#### Responses

**201** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Volume`](#schema-Volume)

<details>
<summary>201 response schema</summary>

```yaml
$ref: "#/components/schemas/Volume"
```

</details>


**409** Creation conflicts with an existing resource (ALREADY_EXISTS).

The Idempotency-Key was already used with a different payload.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>409 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### GET `/v1/volumes/{volume}` {#operation-getVolume-get-v1-volumes-volume}

GetVolume reads one volume by its resource name.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `volumesRead`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `volume` | path | Yes | `string` | The volume id. |
| `If-None-Match` | header | No | `string` | A quoted entity-tag or comma-separated list, compared weakly with the current resource etag. A match, including a weak form or wildcard *, answers 304 without a body after authorization. Malformed input is treated as no match. |


#### Responses

**200** Success

Header `ETag`: `string`. The etag of the resource this response carries, which the next mutation of it sends as If-Match.

`application/json`: [`Volume`](#schema-Volume)

<details>
<summary>200 response schema</summary>

```yaml
$ref: "#/components/schemas/Volume"
```

</details>


**304** The resource's etag equals the If-None-Match value, so this response carries no body.

**404** The target is absent or not visible within the caller's scope.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>404 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



### DELETE `/v1/volumes/{volume}` {#operation-deleteVolume-delete-v1-volumes-volume}

DeleteVolume deletes a volume; deleting an already absent volume succeeds.

Returns unimplemented when this operation is unavailable. Support does not grant permission.

**Base URL:** `https://connect.docker.com/sandboxes`. The API base URL. Append the /v1 paths while preserving this URL's path prefix.

**Authentication:** `bearer`.

**Required permissions:** `volumesWrite`.

#### Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `volume` | path | Yes | `string` | The volume id. |
| `If-Match` | header | Yes | `string` | The complete quoted strong entity-tag observed on the resource, including its quotes. Every mutation of an existing resource requires it: an absent value answers 428 and a stale one 412, each carrying failedPrecondition in the body. |


#### Responses

**204** Success. The response carries no body.

**412** The precondition is not the resource's current etag. The body carries failedPrecondition naming etag_mismatch and an EtagMismatch detail with the current value.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>412 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**428** No precondition was sent. The body carries failedPrecondition naming etag_required.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>428 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**501** This operation is not available on this service.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>501 response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>


**default** The structured Error body identifies the failure with a stable code and optional typed details.

`application/json`: [`Error`](#schema-Error)

<details>
<summary>default response schema</summary>

```yaml
$ref: "#/components/schemas/Error"
```

</details>



## Schemas

Open a schema to read its properties and complete definition. Conditional requirements and alternatives are recorded in the complete definition.

<details id="schema-Error">
<summary>Error</summary>

An HTTP refusal, resource failure, or terminal stream error. Clients preserve unknown codes and structured details.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `code` | Yes | [`ErrorCode`](#schema-ErrorCode) or `string` |  |
| `message` | No | `string` |  |
| `details` | No | Array of [`ErrorDetail`](#schema-ErrorDetail) |  |

```yaml
type: object
title: Error
required:
  - code
description: An HTTP refusal, resource failure, or terminal stream error. Clients preserve unknown codes and structured details.
properties:
  code:
    anyOf:
      - $ref: "#/components/schemas/ErrorCode"
      - type: string
  message:
    type: string
  details:
    type: array
    items:
      $ref: "#/components/schemas/ErrorDetail"
examples:
  - code: notFound
    details: []
    message: sandbox not found
```

</details>

<details id="schema-ErrorCode">
<summary>ErrorCode</summary>

Known public failure categories. Error.code also accepts future strings.



```yaml
type: string
enum:
  - canceled
  - unknown
  - invalidArgument
  - deadlineExceeded
  - notFound
  - alreadyExists
  - permissionDenied
  - resourceExhausted
  - failedPrecondition
  - aborted
  - outOfRange
  - unimplemented
  - internal
  - unavailable
  - dataLoss
  - unauthenticated
description: Known public failure categories. Error.code also accepts future strings.
```

</details>

<details id="schema-ErrorDetail">
<summary>ErrorDetail</summary>





```yaml
oneOf:
  - $ref: "#/components/schemas/BadRequestErrorDetail"
  - $ref: "#/components/schemas/DebugInfoErrorDetail"
  - $ref: "#/components/schemas/ErrorInfoErrorDetail"
  - $ref: "#/components/schemas/HelpErrorDetail"
  - $ref: "#/components/schemas/LocalizedMessageErrorDetail"
  - $ref: "#/components/schemas/PreconditionFailureErrorDetail"
  - $ref: "#/components/schemas/QuotaFailureErrorDetail"
  - $ref: "#/components/schemas/RequestInfoErrorDetail"
  - $ref: "#/components/schemas/ResourceInfoErrorDetail"
  - $ref: "#/components/schemas/RetryInfoErrorDetail"
  - $ref: "#/components/schemas/ValidationViolationsErrorDetail"
  - $ref: "#/components/schemas/EtagMismatchErrorDetail"
  - $ref: "#/components/schemas/KitFieldOutcomeReportErrorDetail"
  - $ref: "#/components/schemas/McpConfigurationReportErrorDetail"
  - $ref: "#/components/schemas/UnknownErrorDetail"
discriminator:
  propertyName: "@type"
  mapping:
    type.googleapis.com/google.rpc.BadRequest: "#/components/schemas/BadRequestErrorDetail"
    type.googleapis.com/google.rpc.DebugInfo: "#/components/schemas/DebugInfoErrorDetail"
    type.googleapis.com/google.rpc.ErrorInfo: "#/components/schemas/ErrorInfoErrorDetail"
    type.googleapis.com/google.rpc.Help: "#/components/schemas/HelpErrorDetail"
    type.googleapis.com/google.rpc.LocalizedMessage: "#/components/schemas/LocalizedMessageErrorDetail"
    type.googleapis.com/google.rpc.PreconditionFailure: "#/components/schemas/PreconditionFailureErrorDetail"
    type.googleapis.com/google.rpc.QuotaFailure: "#/components/schemas/QuotaFailureErrorDetail"
    type.googleapis.com/google.rpc.RequestInfo: "#/components/schemas/RequestInfoErrorDetail"
    type.googleapis.com/google.rpc.ResourceInfo: "#/components/schemas/ResourceInfoErrorDetail"
    type.googleapis.com/google.rpc.RetryInfo: "#/components/schemas/RetryInfoErrorDetail"
    type.googleapis.com/buf.validate.Violations: "#/components/schemas/ValidationViolationsErrorDetail"
    type.googleapis.com/docker.sandboxes.v1.EtagMismatch: "#/components/schemas/EtagMismatchErrorDetail"
    type.googleapis.com/docker.sandboxes.v1.KitFieldOutcomeReport: "#/components/schemas/KitFieldOutcomeReportErrorDetail"
    type.googleapis.com/docker.sandboxes.v1.McpConfigurationReport: "#/components/schemas/McpConfigurationReportErrorDetail"
  defaultMapping: "#/components/schemas/UnknownErrorDetail"
```

</details>

<details id="schema-UnknownErrorDetail">
<summary>UnknownErrorDetail</summary>

A structured error detail identified by @type. Its decoded JSON fields appear beside the type URL.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `@type` | Yes | `string` |  |

```yaml
type: object
required:
  - "@type"
additionalProperties: true
description: A structured error detail identified by @type. Its decoded JSON fields appear beside the type URL.
properties:
  "@type":
    type: string
    not:
      enum:
        - type.googleapis.com/google.rpc.BadRequest
        - type.googleapis.com/google.rpc.DebugInfo
        - type.googleapis.com/google.rpc.ErrorInfo
        - type.googleapis.com/google.rpc.Help
        - type.googleapis.com/google.rpc.LocalizedMessage
        - type.googleapis.com/google.rpc.PreconditionFailure
        - type.googleapis.com/google.rpc.QuotaFailure
        - type.googleapis.com/google.rpc.RequestInfo
        - type.googleapis.com/google.rpc.ResourceInfo
        - type.googleapis.com/google.rpc.RetryInfo
        - type.googleapis.com/buf.validate.Violations
        - type.googleapis.com/docker.sandboxes.v1.EtagMismatch
        - type.googleapis.com/docker.sandboxes.v1.KitFieldOutcomeReport
        - type.googleapis.com/docker.sandboxes.v1.McpConfigurationReport
```

</details>

<details id="schema-ValidationFieldPath">
<summary>ValidationFieldPath</summary>

A path to a field, including any parent fields.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `elements` | No | Array of [`ValidationFieldPathElement`](#schema-ValidationFieldPathElement) | `elements` contains each element of the path, starting from the root and recursing downward. |

```yaml
type: object
properties:
  elements:
    type: array
    items:
      $ref: "#/components/schemas/ValidationFieldPathElement"
    description: "`elements` contains each element of the path, starting from the root and recursing downward."
title: FieldPath
description: A path to a field, including any parent fields.
```

</details>

<details id="schema-ValidationFieldPathElement">
<summary>ValidationFieldPathElement</summary>

One field in a validation path. For a map or list, the subscript identifies the selected element.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `fieldNumber` | No | `integer or null` (int32) | `field_number` is the field number this path element refers to. |
| `fieldName` | No | `string or null` | `field_name` contains the field name this path element refers to. This can be used to display a human-readable path even if the field number is unknown. |
| `fieldType` | No | [`ValidationFieldType`](#schema-ValidationFieldType) or `null` | The type of the field that failed validation. |
| `keyType` | No | [`ValidationFieldType`](#schema-ValidationFieldType) or `null` | `key_type` specifies the map key type of this field. This value is useful when traversing unknown fields through wire data: specifically, it allows handling the differences between different integer encodings. |
| `valueType` | No | [`ValidationFieldType`](#schema-ValidationFieldType) or `null` | `value_type` specifies map value type of this field. This is useful if you want to display a value inside unknown fields through wire data. |

```yaml
type: object
allOf:
  - properties:
      fieldNumber:
        type:
          - integer
          - "null"
        format: int32
        description: "`field_number` is the field number this path element refers to."
      fieldName:
        type:
          - string
          - "null"
        description: |-
          `field_name` contains the field name this path element refers to.
          This can be used to display a human-readable path even if the field number is unknown.
      fieldType:
        oneOf:
          - $ref: "#/components/schemas/ValidationFieldType"
          - type: "null"
        description: The type of the field that failed validation.
      keyType:
        oneOf:
          - $ref: "#/components/schemas/ValidationFieldType"
          - type: "null"
        description: |-
          `key_type` specifies the map key type of this field. This value is useful when traversing
          unknown fields through wire data: specifically, it allows handling the differences between
          different integer encodings.
      valueType:
        oneOf:
          - $ref: "#/components/schemas/ValidationFieldType"
          - type: "null"
        description: |-
          `value_type` specifies map value type of this field. This is useful if you want to display a
          value inside unknown fields through wire data.
  - oneOf:
      - type: object
        properties:
          boolKey:
            type: boolean
            description: "`bool_key` specifies a map key of type bool."
        title: bool_key
        required:
          - boolKey
      - type: object
        properties:
          index:
            type:
              - integer
              - string
            format: int64
            description: "`index` specifies a 0-based index into a repeated field."
        title: index
        required:
          - index
      - type: object
        properties:
          intKey:
            type:
              - integer
              - string
            format: int64
            description: "`int_key` specifies a map key of type int32, int64, sint32, sint64, sfixed32 or sfixed64."
        title: int_key
        required:
          - intKey
      - type: object
        properties:
          stringKey:
            type: string
            description: "`string_key` specifies a map key of type string."
        title: string_key
        required:
          - stringKey
      - type: object
        properties:
          uintKey:
            type:
              - integer
              - string
            format: int64
            description: "`uint_key` specifies a map key of type uint32, uint64, fixed32 or fixed64."
        title: uint_key
        required:
          - uintKey
title: FieldPathElement
description: One field in a validation path. For a map or list, the subscript identifies the selected element.
```

</details>

<details id="schema-ValidationViolation">
<summary>ValidationViolation</summary>

A validation failure, with the affected field, rule, and a human-readable message.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `field` | No | [`ValidationFieldPath`](#schema-ValidationFieldPath) or `null` | The field that failed validation, including its parent fields. |
| `rule` | No | [`ValidationFieldPath`](#schema-ValidationFieldPath) or `null` | The validation rule that failed. |
| `ruleId` | No | `string or null` | `rule_id` is the unique identifier of the `Rule` that was not fulfilled. This is the same `id` that was specified in the `Rule` message, allowing easy tracing of which rule was violated. |
| `message` | No | `string or null` | `message` is a human-readable error message that describes the nature of the violation. This can be the default error message from the violated `Rule`, or it can be a custom message that gives more context about the violation. |
| `forKey` | No | `boolean or null` | `for_key` indicates whether the violation was caused by a map key, rather than a value. |

```yaml
type: object
properties:
  field:
    oneOf:
      - $ref: "#/components/schemas/ValidationFieldPath"
      - type: "null"
    description: The field that failed validation, including its parent fields.
  rule:
    oneOf:
      - $ref: "#/components/schemas/ValidationFieldPath"
      - type: "null"
    description: The validation rule that failed.
  ruleId:
    type:
      - string
      - "null"
    description: |-
      `rule_id` is the unique identifier of the `Rule` that was not fulfilled.
      This is the same `id` that was specified in the `Rule` message, allowing easy tracing of which rule was violated.
  message:
    type:
      - string
      - "null"
    description: |-
      `message` is a human-readable error message that describes the nature of the violation.
      This can be the default error message from the violated `Rule`, or it can be a custom message that gives more context about the violation.
  forKey:
    type:
      - boolean
      - "null"
    description: "`for_key` indicates whether the violation was caused by a map key, rather than a value."
title: Violation
description: A validation failure, with the affected field, rule, and a human-readable message.
```

</details>

<details id="schema-ValidationViolations">
<summary>ValidationViolations</summary>

The validation failures returned for a request.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `violations` | No | Array of [`ValidationViolation`](#schema-ValidationViolation) | `violations` is a repeated field that contains all the `Violation` messages corresponding to the violations detected. |

```yaml
type: object
properties:
  violations:
    type: array
    items:
      $ref: "#/components/schemas/ValidationViolation"
    description: "`violations` is a repeated field that contains all the `Violation` messages corresponding to the violations detected."
title: Violations
description: The validation failures returned for a request.
```

</details>

<details id="schema-ValidationViolationsErrorDetail">
<summary>ValidationViolationsErrorDetail</summary>

The validation failures returned for a request.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `violations` | No | Array of [`ValidationViolation`](#schema-ValidationViolation) | `violations` is a repeated field that contains all the `Violation` messages corresponding to the violations detected. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  violations:
    type: array
    items:
      $ref: "#/components/schemas/ValidationViolation"
    description: "`violations` is a repeated field that contains all the `Violation` messages corresponding to the violations detected."
  "@type":
    type: string
    const: type.googleapis.com/buf.validate.Violations
title: ViolationsErrorDetail
description: The validation failures returned for a request.
required:
  - "@type"
```

</details>

<details id="schema-AppliedKit">
<summary>AppliedKit</summary>

AppliedKit records kit content that was applied to a sandbox.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `ref` | No | `string` | ref is the kit reference supplied or resolved by the backend. |
| `kind` | Yes | `string` | kind is normalized to "sandbox" or "mixin". |
| `appliedAt` | No | [`Timestamp`](#schema-Timestamp) | applied_at is when the kit became active. |
| `digest` | No | `string` | digest identifies the exact applied content. |

```yaml
type: object
properties:
  ref:
    type: string
    description: ref is the kit reference supplied or resolved by the backend.
  kind:
    type: string
    enum:
      - sandbox
      - mixin
    description: kind is normalized to "sandbox" or "mixin".
  appliedAt:
    description: applied_at is when the kit became active.
    $ref: "#/components/schemas/Timestamp"
  digest:
    type: string
    description: digest identifies the exact applied content.
title: AppliedKit
required:
  - kind
description: AppliedKit records kit content that was applied to a sandbox.
```

</details>

<details id="schema-EndpointAuthentication">
<summary>EndpointAuthentication</summary>

EndpointAuthentication describes credential issuance without publishing a token.
local sockets cannot advertise bearer credential transports
scoped bearer discovery states the five-minute credential lifetime; local sockets carry no token lifetime

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `scheme` | Yes | [`EndpointAuthenticationScheme`](#schema-EndpointAuthenticationScheme) |  |
| `defaultTtl` | No | [`Duration`](#schema-Duration) | default_ttl and max_ttl describe CreateEndpointCredential's provider policy. |
| `maxTtl` | No | [`Duration`](#schema-Duration) |  |
| `credentialTransports` | No | Array of [`EndpointCredentialTransport`](#schema-EndpointCredentialTransport) | Omission retains Authorization-header clients; browsers require an explicit subprotocol offer. |

```yaml
type: object
properties:
  scheme:
    not:
      enum:
        - unspecified
    $ref: "#/components/schemas/EndpointAuthenticationScheme"
  defaultTtl:
    description: default_ttl and max_ttl describe CreateEndpointCredential's provider policy.
    $ref: "#/components/schemas/Duration"
  maxTtl:
    $ref: "#/components/schemas/Duration"
  credentialTransports:
    type: array
    items:
      $ref: "#/components/schemas/EndpointCredentialTransport"
    uniqueItems: true
    description: Omission retains Authorization-header clients; browsers require an explicit subprotocol offer.
title: EndpointAuthentication
required:
  - scheme
description: |-
  EndpointAuthentication describes credential issuance without publishing a token.
  local sockets cannot advertise bearer credential transports
  scoped bearer discovery states the five-minute credential lifetime; local sockets carry no token lifetime
dependentSchemas:
  scheme:
    allOf:
      - oneOf:
          - required:
              - defaultTtl
              - maxTtl
            properties:
              scheme:
                const: scopedBearer
              defaultTtl:
                pattern: ^300([.]0{1,9})?s$
              maxTtl:
                pattern: ^300([.]0{1,9})?s$
          - properties:
              scheme:
                const: localSocket
            not:
              anyOf:
                - required:
                    - defaultTtl
                - required:
                    - maxTtl
  credentialTransports:
    allOf:
      - if:
          properties:
            scheme:
              const: localSocket
        then:
          properties:
            credentialTransports:
              maxItems: 0
```

</details>

<details id="schema-EndpointAuthenticationScheme">
<summary>EndpointAuthenticationScheme</summary>

Values:
unspecified
scopedBearer
localSocket



```yaml
type: string
title: EndpointAuthenticationScheme
enum:
  - unspecified
  - scopedBearer
  - localSocket
description: |-
  Values:
  unspecified
  scopedBearer
  localSocket
```

</details>

<details id="schema-EndpointCredentialTransport">
<summary>EndpointCredentialTransport</summary>

Values:
unspecified
authorizationHeader
webSocketSubprotocol



```yaml
type: string
title: EndpointCredentialTransport
enum:
  - unspecified
  - authorizationHeader
  - webSocketSubprotocol
description: |-
  Values:
  unspecified
  authorizationHeader
  webSocketSubprotocol
```

</details>

<details id="schema-KitRef">
<summary>KitRef</summary>

KitRef names a kit to apply.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `ref` | Yes | `string` | ref names the kit source. |
| `kind` | No | `string` | kind is "sandbox", "mixin", or empty to defer to the kit's own manifest. |

```yaml
type: object
properties:
  ref:
    type: string
    minLength: 1
    description: ref names the kit source.
  kind:
    type: string
    enum:
      - ""
      - sandbox
      - mixin
    description: kind is "sandbox", "mixin", or empty to defer to the kit's own manifest.
title: KitRef
required:
  - ref
additionalProperties: false
description: KitRef names a kit to apply.
```

</details>

<details id="schema-Permission">
<summary>Permission</summary>

The actions a credential is allowed to perform. Each operation declares its required permissions; service support does not grant permission.



```yaml
type: string
title: Permission
enum:
  - unspecified
  - sandboxesRead
  - sandboxesCreate
  - sandboxesExec
  - sandboxesFilesRead
  - sandboxesFilesWrite
  - imagesPull
  - secretsWrite
  - sandboxesDelete
  - sandboxesLifecycle
  - sandboxesKits
  - sandboxesPorts
  - sandboxesCredential
  - sandboxesSsh
  - imagesRead
  - imagesWrite
  - secretsRead
  - networkPoliciesRead
  - volumesRead
  - volumesWrite
  - snapshotsRead
  - snapshotsWrite
  - mcpRead
  - mcpWrite
  - credentialsExchange
  - sandboxesRename
  - secretsResolve
  - operator
description: The actions a credential is allowed to perform. Each operation declares its required permissions; service support does not grant permission.
```

</details>

<details id="schema-Platform">
<summary>Platform</summary>

Platform is the OCI platform tuple.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `os` | No | `string` | os is the operating system, for example "linux". |
| `architecture` | No | `string` | architecture is the CPU architecture, for example "amd64" or "arm64". |

```yaml
type: object
properties:
  os:
    type: string
    description: os is the operating system, for example "linux".
  architecture:
    type: string
    description: architecture is the CPU architecture, for example "amd64" or "arm64".
title: Platform
additionalProperties: false
description: Platform is the OCI platform tuple.
```

</details>

<details id="schema-Protocol">
<summary>Protocol</summary>

Protocol names the transport protocol for a published port.

Values:
unspecified: unspecified means the backend default, TCP.
tcp: tcp publishes TCP traffic.
udp: udp publishes UDP traffic.
tcp4: tcp4 publishes IPv4 TCP traffic.
tcp6: tcp6 publishes IPv6 TCP traffic.
udp4: udp4 publishes IPv4 UDP traffic.
udp6: udp6 publishes IPv6 UDP traffic.



```yaml
type: string
title: Protocol
enum:
  - unspecified
  - tcp
  - udp
  - tcp4
  - tcp6
  - udp4
  - udp6
description: |-
  Protocol names the transport protocol for a published port.

  Values:
  unspecified: unspecified means the backend default, TCP.
  tcp: tcp publishes TCP traffic.
  udp: udp publishes UDP traffic.
  tcp4: tcp4 publishes IPv4 TCP traffic.
  tcp6: tcp6 publishes IPv6 TCP traffic.
  udp4: udp4 publishes IPv4 UDP traffic.
  udp6: udp6 publishes IPv6 UDP traffic.
```

</details>

<details id="schema-PublishedPort">
<summary>PublishedPort</summary>

PublishedPort reports external reachability for one sandbox port.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `hostIp` | No | `string` | host_ip is the local host bind address when relevant. |
| `hostPort` | No | `integer` | host_port is the local host port when relevant. |
| `sandboxPort` | No | `integer` | sandbox_port is the port inside the sandbox. |
| `protocol` | No | [`Protocol`](#schema-Protocol) | protocol is the published protocol. |
| `url` | No | `string` | url is the required externally reachable address, including host-local addresses. |

```yaml
type: object
properties:
  hostIp:
    type: string
    description: host_ip is the local host bind address when relevant.
  hostPort:
    type: integer
    description: host_port is the local host port when relevant.
  sandboxPort:
    type: integer
    description: sandbox_port is the port inside the sandbox.
  protocol:
    description: protocol is the published protocol.
    $ref: "#/components/schemas/Protocol"
  url:
    type: string
    description: url is the required externally reachable address, including host-local addresses.
title: PublishedPort
description: PublishedPort reports external reachability for one sandbox port.
```

</details>

<details id="schema-Resources">
<summary>Resources</summary>

Resources describes sandbox CPU and memory allocation.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `cpus` | No | `integer or null` | cpus is present only when the caller or backend chooses a CPU count. |
| `memoryMib` | No | `integer or string or null` (int64) | memory_mib is present only when the caller or backend chooses memory. |

```yaml
type: object
properties:
  cpus:
    type:
      - integer
      - "null"
    description: cpus is present only when the caller or backend chooses a CPU count.
  memoryMib:
    type:
      - integer
      - string
      - "null"
    format: int64
    description: memory_mib is present only when the caller or backend chooses memory.
title: Resources
additionalProperties: false
description: Resources describes sandbox CPU and memory allocation.
```

</details>

<details id="schema-SandboxCore">
<summary>SandboxCore</summary>

SandboxCore is the backend-neutral sandbox view. At most one of image, imageRef is set.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `agent` | No | `string` | agent is the named agent profile, empty for raw image/start-command sandboxes. |
| `status` | Yes | [`SandboxStatus`](#schema-SandboxStatus) | status is the coarse lifecycle state. |
| `resources` | No | [`Resources`](#schema-Resources) | resources are the sandbox resources. |
| `platform` | No | [`Platform`](#schema-Platform) | platform is the canonical placed platform. |
| `environment` | No | `object` | environment is non-secret environment metadata. |
| `ports` | No | Array of [`PublishedPort`](#schema-PublishedPort) | ports are the currently published ports. |
| `appliedKits` | No | Array of [`AppliedKit`](#schema-AppliedKit) | applied_kits records kits active in the sandbox. |
| `createdAt` | Yes | [`Timestamp`](#schema-Timestamp) | created_at is the creation time. |
| `endpoint` | No | [`SandboxEndpoint`](#schema-SandboxEndpoint) | endpoint is where process and file services are served for this sandbox. |
| `origin` | No | `string` | origin is set by merged views to identify the backend that owns this sandbox. |
| `labels` | No | `object` | labels are non-secret caller-assigned key/value metadata. |
| `etag` | Yes | `string` | etag identifies the observed version of this resource. It is opaque and strong, changes on every visible change, and is what a mutation sends as its precondition. |
| `image` | No | `string` | image is the managed Image resource name. |
| `imageRef` | No | `string` | image_ref is an external OCI reference, with the existing backend validation. |

```yaml
type: object
allOf:
  - properties:
      agent:
        type: string
        description: agent is the named agent profile, empty for raw image/start-command sandboxes.
      status:
        not:
          enum:
            - unspecified
        description: status is the coarse lifecycle state.
        readOnly: true
        $ref: "#/components/schemas/SandboxStatus"
      resources:
        description: resources are the sandbox resources.
        $ref: "#/components/schemas/Resources"
      platform:
        description: platform is the canonical placed platform.
        $ref: "#/components/schemas/Platform"
      environment:
        type: object
        additionalProperties:
          type: string
          title: value
        description: environment is non-secret environment metadata.
      ports:
        type: array
        items:
          $ref: "#/components/schemas/PublishedPort"
        description: ports are the currently published ports.
      appliedKits:
        type: array
        items:
          $ref: "#/components/schemas/AppliedKit"
        description: applied_kits records kits active in the sandbox.
      createdAt:
        description: created_at is the creation time.
        $ref: "#/components/schemas/Timestamp"
      endpoint:
        description: endpoint is where process and file services are served for this sandbox.
        $ref: "#/components/schemas/SandboxEndpoint"
      origin:
        type: string
        description: origin is set by merged views to identify the backend that owns this sandbox.
      labels:
        type: object
        additionalProperties:
          type: string
          title: value
        description: labels are non-secret caller-assigned key/value metadata.
      etag:
        type: string
        minLength: 1
        pattern: ^"[^"\x00-\x20\x7f]*"$
        description: |-
          etag identifies the observed version of this resource. It is opaque and strong,
          changes on every visible change, and is what a mutation sends as its precondition.
        readOnly: true
  - properties:
      image:
        type: string
        pattern: ^images/[^/]+$
        description: image is the managed Image resource name.
      imageRef:
        type: string
        description: image_ref is an external OCI reference, with the existing backend validation.
    dependentSchemas:
      image:
        not:
          anyOf:
            - required:
                - imageRef
      imageRef:
        not:
          anyOf:
            - required:
                - image
title: SandboxCore
required:
  - status
  - createdAt
  - etag
description: SandboxCore is the backend-neutral sandbox view. At most one of image, imageRef is set.
```

</details>

<details id="schema-SandboxEndpoint">
<summary>SandboxEndpoint</summary>

SandboxEndpoint tells a client where sandbox endpoint services are served.
network endpoints require scoped credentials; Unix sockets use actual socket access controls
sandbox must contain the endpoint's immutable sandbox UID

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `uri` | No | `string` | uri is the endpoint API base; preserve its path prefix when appending /v1 routes. |
| `protocol` | No | [`SandboxEndpointProtocol`](#schema-SandboxEndpointProtocol) | protocol identifies the endpoint transport. |
| `credentialAudience` | No | `string` | credential_audience is server-controlled and required for network endpoints. |
| `sandbox` | Yes | `string` | sandbox names the immutable sandbox incarnation served by this endpoint. |
| `sandboxUid` | Yes | `string` |  |
| `apiVersion` | Yes | `string` | api_version identifies the /v1 endpoint contract before the client connects. |
| `capabilities` | Yes | Array of `string` | capabilities lists supported public operation IDs, independently of caller grants. |
| `authentication` | Yes | [`EndpointAuthentication`](#schema-EndpointAuthentication) |  |

```yaml
type: object
properties:
  uri:
    type: string
    description: uri is the endpoint API base; preserve its path prefix when appending /v1 routes.
  protocol:
    description: protocol identifies the endpoint transport.
    $ref: "#/components/schemas/SandboxEndpointProtocol"
  credentialAudience:
    type: string
    description: credential_audience is server-controlled and required for network endpoints.
  sandbox:
    type: string
    pattern: ^sandboxes/[^/]+$
    description: sandbox names the immutable sandbox incarnation served by this endpoint.
  sandboxUid:
    type:
      - string
    minLength: 1
    readOnly: true
  apiVersion:
    type: string
    description: api_version identifies the /v1 endpoint contract before the client connects.
    const: v1
  capabilities:
    type: array
    items:
      type: string
      pattern: ^[a-z][A-Za-z0-9]*$
    minItems: 1
    uniqueItems: true
    description: capabilities lists supported public operation IDs, independently of caller grants.
  authentication:
    $ref: "#/components/schemas/EndpointAuthentication"
title: SandboxEndpoint
required:
  - sandbox
  - sandboxUid
  - apiVersion
  - capabilities
  - authentication
description: |-
  SandboxEndpoint tells a client where sandbox endpoint services are served.
  network endpoints require scoped credentials; Unix sockets use actual socket access controls
  sandbox must contain the endpoint's immutable sandbox UID
dependentSchemas:
  authentication:
    allOf:
      - required:
          - protocol
        oneOf:
          - required:
              - credentialAudience
            properties:
              protocol:
                const: http
              credentialAudience:
                minLength: 1
              authentication:
                properties:
                  scheme:
                    const: scopedBearer
          - properties:
              protocol:
                const: unixSocket
              credentialAudience:
                const: ""
              authentication:
                properties:
                  scheme:
                    const: localSocket
```

</details>

<details id="schema-SandboxEndpointProtocol">
<summary>SandboxEndpointProtocol</summary>

SandboxEndpointProtocol identifies the transport used by a sandbox endpoint.

Values:
unspecified: unspecified is not a concrete endpoint protocol.
http: http is the REST API over HTTP.
unixSocket: unixSocket is a local unix socket.



```yaml
type: string
title: SandboxEndpointProtocol
enum:
  - unspecified
  - http
  - unixSocket
description: |-
  SandboxEndpointProtocol identifies the transport used by a sandbox endpoint.

  Values:
  unspecified: unspecified is not a concrete endpoint protocol.
  http: http is the REST API over HTTP.
  unixSocket: unixSocket is a local unix socket.
```

</details>

<details id="schema-SandboxStatus">
<summary>SandboxStatus</summary>

SandboxStatus is the backend-neutral lifecycle state.

Values:
unspecified: unspecified is never a concrete known state.
creating: creating means provisioning has not completed.
starting: starting means a stopped sandbox is starting or resuming.
running: running means the sandbox is running.
stopping: stopping means the sandbox is stopping or hibernating.
stopped: stopped means the sandbox is not running and may be started.
failed: failed means the sandbox reached an unrecoverable failure.
deleting: deleting means deletion is accepted; reads continue until removal.
degraded: degraded means runtime contact is lost; recovery or failure may follow.



```yaml
type: string
title: SandboxStatus
enum:
  - unspecified
  - creating
  - starting
  - running
  - stopping
  - stopped
  - failed
  - deleting
  - degraded
description: |-
  SandboxStatus is the backend-neutral lifecycle state.

  Values:
  unspecified: unspecified is never a concrete known state.
  creating: creating means provisioning has not completed.
  starting: starting means a stopped sandbox is starting or resuming.
  running: running means the sandbox is running.
  stopping: stopping means the sandbox is stopping or hibernating.
  stopped: stopped means the sandbox is not running and may be started.
  failed: failed means the sandbox reached an unrecoverable failure.
  deleting: deleting means deletion is accepted; reads continue until removal.
  degraded: degraded means runtime contact is lost; recovery or failure may follow.
```

</details>

<details id="schema-EffectiveCoreLifecycle">
<summary>EffectiveCoreLifecycle</summary>

EffectiveCoreLifecycle is persisted at admission and cannot silently degrade.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `stopMemoryOutcome` | Yes | [`StopMemoryOutcome`](#schema-StopMemoryOutcome) | stop_memory_outcome is the selected stop/start guarantee for this sandbox. |
| `canCaptureMemory` | No | `boolean` | can_capture_memory is this sandbox's ALL snapshot support, independently set. |

```yaml
type: object
properties:
  stopMemoryOutcome:
    not:
      enum:
        - unspecified
    description: stop_memory_outcome is the selected stop/start guarantee for this sandbox.
    $ref: "#/components/schemas/StopMemoryOutcome"
  canCaptureMemory:
    type: boolean
    description: can_capture_memory is this sandbox's ALL snapshot support, independently set.
title: EffectiveCoreLifecycle
required:
  - stopMemoryOutcome
description: EffectiveCoreLifecycle is persisted at admission and cannot silently degrade.
```

</details>

<details id="schema-EffectiveFeatures">
<summary>EffectiveFeatures</summary>

EffectiveFeatures is emitted only through the composition representation.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `lifecycle` | Yes | [`EffectiveCoreLifecycle`](#schema-EffectiveCoreLifecycle) | lifecycle is required even when no optional feature is selected. |
| `outputRetention` | Yes | [`OutputRetention`](#schema-OutputRetention) | output_retention is the persisted per-process reconnect guarantee. |
| `managedVolumes` | No | [`EffectiveManagedVolumes`](#schema-EffectiveManagedVolumes) | managed_volumes records stable volume identities and concrete attachment modes. |
| `timeouts` | No | [`EffectiveTimeouts`](#schema-EffectiveTimeouts) | timeouts records frozen defaults and the currently armed deadline, if any. |

```yaml
type: object
properties:
  lifecycle:
    description: lifecycle is required even when no optional feature is selected.
    $ref: "#/components/schemas/EffectiveCoreLifecycle"
  outputRetention:
    description: output_retention is the persisted per-process reconnect guarantee.
    $ref: "#/components/schemas/OutputRetention"
  managedVolumes:
    description: managed_volumes records stable volume identities and concrete attachment modes.
    $ref: "#/components/schemas/EffectiveManagedVolumes"
  timeouts:
    description: timeouts records frozen defaults and the currently armed deadline, if any.
    $ref: "#/components/schemas/EffectiveTimeouts"
title: EffectiveFeatures
required:
  - lifecycle
  - outputRetention
description: EffectiveFeatures is emitted only through the composition representation.
```

</details>

<details id="schema-FeatureOptions">
<summary>FeatureOptions</summary>

FeatureOptions carries typed admission intent; unknown fields are refused.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `managedVolumes` | No | [`ManagedVolumeOptions`](#schema-ManagedVolumeOptions) | managed_volumes selects persistent resources and guest targets. |
| `timeouts` | No | [`TimeoutOptions`](#schema-TimeoutOptions) | timeouts enables an explicit deadline; omission enables none. |
| `storedSecrets` | No | [`StoredSecretOptions`](#schema-StoredSecretOptions) | stored_secrets selects authorized stored references, never raw material. |

```yaml
type: object
properties:
  managedVolumes:
    description: managed_volumes selects persistent resources and guest targets.
    $ref: "#/components/schemas/ManagedVolumeOptions"
  timeouts:
    description: timeouts enables an explicit deadline; omission enables none.
    $ref: "#/components/schemas/TimeoutOptions"
  storedSecrets:
    description: stored_secrets selects authorized stored references, never raw material.
    $ref: "#/components/schemas/StoredSecretOptions"
title: FeatureOptions
additionalProperties: false
description: FeatureOptions carries typed admission intent; unknown fields are refused.
```

</details>

<details id="schema-OutputRetention">
<summary>OutputRetention</summary>

OutputRetention records minimum reconnect coverage independently of location.
output retention requires a duration or byte floor

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `minDuration` | No | [`Duration`](#schema-Duration) | min_duration retains each process's output for at least this age when set. |
| `minBytes` | No | `integer or string` (int64) | min_bytes retains at least this many recent bytes per process when nonzero. a byte retention floor must be zero or at least 64 KiB |

```yaml
type: object
properties:
  minDuration:
    description: min_duration retains each process's output for at least this age when set.
    $ref: "#/components/schemas/Duration"
  minBytes:
    type:
      - integer
      - string
    format: int64
    description: |-
      min_bytes retains at least this many recent bytes per process when nonzero.
      a byte retention floor must be zero or at least 64 KiB
title: OutputRetention
description: |-
  OutputRetention records minimum reconnect coverage independently of location.
  output retention requires a duration or byte floor
```

</details>

<details id="schema-StopMemoryOutcome">
<summary>StopMemoryOutcome</summary>

StopMemoryOutcome records the stop/start guarantee independently of location.

Values:
unspecified: unspecified is not an advertised or accepted guarantee.
discard: discard retains disk and starts a fresh process tree.
preserve: preserve restores the stopped process tree on start.



```yaml
type: string
title: StopMemoryOutcome
enum:
  - unspecified
  - discard
  - preserve
description: |-
  StopMemoryOutcome records the stop/start guarantee independently of location.

  Values:
  unspecified: unspecified is not an advertised or accepted guarantee.
  discard: discard retains disk and starts a fresh process tree.
  preserve: preserve restores the stopped process tree on start.
```

</details>

<details id="schema-FilesDownloadRequest">
<summary>FilesDownloadRequest</summary>

DownloadRequest asks for one or more paths.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `paths` | Yes | Array of `string` | paths naming stable regular files or absent entries at canonical absolute paths yield contiguous results in request order, including duplicates. |

```yaml
type: object
properties:
  paths:
    type: array
    items:
      type: string
    minItems: 1
    description: |-
      paths naming stable regular files or absent entries at canonical absolute paths
      yield contiguous results in request order, including duplicates.
title: DownloadRequest
required:
  - paths
additionalProperties: false
description: DownloadRequest asks for one or more paths.
```

</details>

<details id="schema-FilesDownloadResponse">
<summary>FilesDownloadResponse</summary>

DownloadResponse is one frame in a download stream. At most one of data, error, header is set.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `data` | No | `string` (byte) | data carries bytes for the current file. |
| `error` | No | [`FileError`](#schema-FileError) | error reports a path-level failure without aborting the whole stream. |
| `header` | No | [`FileHeader`](#schema-FileHeader) | header starts a downloaded file. |

```yaml
type: object
title: DownloadResponse
unevaluatedProperties: false
description: DownloadResponse is one frame in a download stream. At most one of data, error, header is set.
properties:
  data:
    type: string
    format: byte
    description: data carries bytes for the current file.
  error:
    description: error reports a path-level failure without aborting the whole stream.
    $ref: "#/components/schemas/FileError"
  header:
    description: header starts a downloaded file.
    $ref: "#/components/schemas/FileHeader"
dependentSchemas:
  data:
    not:
      anyOf:
        - required:
            - error
        - required:
            - header
  error:
    not:
      anyOf:
        - required:
            - data
        - required:
            - header
  header:
    not:
      anyOf:
        - required:
            - data
        - required:
            - error
```

</details>

<details id="schema-FileError">
<summary>FileError</summary>

FileError reports a path-level file error.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `path` | No | `string` | path is the path that failed. |
| `failure` | Yes | [`Error`](#schema-Error) | failure retains the shared typed code and details for this path. |

```yaml
type: object
properties:
  path:
    type: string
    description: path is the path that failed.
  failure:
    description: failure retains the shared typed code and details for this path.
    $ref: "#/components/schemas/Error"
title: FileError
required:
  - failure
description: FileError reports a path-level file error.
```

</details>

<details id="schema-FileHeader">
<summary>FileHeader</summary>

FileHeader starts a file in an upload or download stream.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `path` | Yes | `string` | path is the sandbox path. |
| `mode` | No | `integer` | mode carries Unix permission bits. Upload applies mode & 0777, or 0644 when that value is zero, without subtracting umask; Download reports stored bits. |

```yaml
type: object
properties:
  path:
    type: string
    minLength: 1
    description: path is the sandbox path.
  mode:
    type: integer
    description: |-
      mode carries Unix permission bits. Upload applies mode & 0777, or 0644 when
      that value is zero, without subtracting umask; Download reports stored bits.
title: FileHeader
required:
  - path
additionalProperties: false
description: FileHeader starts a file in an upload or download stream.
```

</details>

<details id="schema-FileInfo">
<summary>FileInfo</summary>

FileInfo is file metadata.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `path` | No | `string` | path is the sandbox path. |
| `size` | No | `integer or string` (int64) | size is the file size in bytes. |
| `mode` | No | `integer` | mode is the stored Unix permission bits (0777), without type or special bits. |
| `isDir` | No | `boolean` | is_dir reports whether the path is a directory and stays populated on every backend and every caller. type is the richer authority, and type == directory iff is_dir. |
| `type` | No | [`FileType`](#schema-FileType) | type classifies the object, and Stat and List populate it for every caller. Stat and List resolve an in-jail symlink and report the resolved target's type; a link that cannot be resolved in-jail (dangling, looping, or escaping) reports symlink for itself rather than disclosing an out-of-jail target. |
| `linkTarget` | No | `string` | link_target is the literal readlink() text when type is symlink, since that is the only case where nothing else lets the caller learn what the link names; empty otherwise. |

```yaml
type: object
properties:
  path:
    type: string
    description: path is the sandbox path.
  size:
    type:
      - integer
      - string
    format: int64
    description: size is the file size in bytes.
  mode:
    type: integer
    description: mode is the stored Unix permission bits (0777), without type or special bits.
  isDir:
    type: boolean
    description: |-
      is_dir reports whether the path is a directory and stays populated on
      every backend and every caller. type is the richer authority, and
      type == directory iff is_dir.
  type:
    description: |-
      type classifies the object, and Stat and List populate it for every
      caller. Stat and List resolve an in-jail symlink and report
      the resolved target's type; a link that cannot be resolved in-jail
      (dangling, looping, or escaping) reports symlink for itself
      rather than disclosing an out-of-jail target.
    $ref: "#/components/schemas/FileType"
  linkTarget:
    type: string
    description: |-
      link_target is the literal readlink() text when type is symlink,
      since that is the only case where nothing else lets the caller learn what
      the link names; empty otherwise.
title: FileInfo
description: FileInfo is file metadata.
```

</details>

<details id="schema-FileType">
<summary>FileType</summary>

FileType classifies the object at a path.

Values:
unspecified: unspecified means the backend has not classified the object. Stat and List always classify it.
regular: regular is a regular file.
directory: directory is a directory.
symlink: symlink is a symbolic link that could not be resolved to a target inside the jail.
other: other is a socket, FIFO, device, or other non-regular, non-directory, non-symlink object.



```yaml
type: string
title: FileType
enum:
  - unspecified
  - regular
  - directory
  - symlink
  - other
description: |-
  FileType classifies the object at a path.

  Values:
  unspecified: unspecified means the backend has not classified the object. Stat and List always classify it.
  regular: regular is a regular file.
  directory: directory is a directory.
  symlink: symlink is a symbolic link that could not be resolved to a target inside the jail.
  other: other is a socket, FIFO, device, or other non-regular, non-directory, non-symlink object.
```

</details>

<details id="schema-FilesListResponse">
<summary>FilesListResponse</summary>

ListResponse returns directory entries and a next-page token.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `entries` | No | Array of [`FileInfo`](#schema-FileInfo) | entries is the current page. |
| `nextPageToken` | No | `string` | next_page_token is empty when there are no more pages. |

```yaml
type: object
properties:
  entries:
    type: array
    items:
      $ref: "#/components/schemas/FileInfo"
    description: entries is the current page.
  nextPageToken:
    type: string
    description: next_page_token is empty when there are no more pages.
title: ListResponse
description: ListResponse returns directory entries and a next-page token.
```

</details>

<details id="schema-FilesMkdirRequest">
<summary>FilesMkdirRequest</summary>

MkdirRequest creates a directory.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `path` | Yes | `string` | path is the directory path. |
| `mode` | No | `integer` | mode supplies the new leaf's Unix permission bits: mode & 0777, or 0755 when that value is zero, without subtracting umask. |
| `parents` | No | `boolean` | parents creates missing parents when true. |

```yaml
type: object
properties:
  path:
    type: string
    minLength: 1
    description: path is the directory path.
  mode:
    type: integer
    description: |-
      mode supplies the new leaf's Unix permission bits: mode & 0777, or 0755
      when that value is zero, without subtracting umask.
  parents:
    type: boolean
    description: parents creates missing parents when true.
title: MkdirRequest
required:
  - path
additionalProperties: false
description: MkdirRequest creates a directory.
```

</details>

<details id="schema-FilesMkdirResponse">
<summary>FilesMkdirResponse</summary>

MkdirResponse has no fields.



```yaml
type: object
title: MkdirResponse
description: MkdirResponse has no fields.
```

</details>

<details id="schema-FilesMoveRequest">
<summary>FilesMoveRequest</summary>

MoveRequest moves or renames a path.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `from` | Yes | `string` | from is the source path. |
| `to` | Yes | `string` | to is the exact destination path. |

```yaml
type: object
properties:
  from:
    type: string
    minLength: 1
    description: from is the source path.
  to:
    type: string
    minLength: 1
    description: to is the exact destination path.
title: MoveRequest
required:
  - from
  - to
additionalProperties: false
description: MoveRequest moves or renames a path.
```

</details>

<details id="schema-FilesMoveResponse">
<summary>FilesMoveResponse</summary>

MoveResponse has no fields.



```yaml
type: object
title: MoveResponse
description: MoveResponse has no fields.
```

</details>

<details id="schema-FilesReadFileRequest">
<summary>FilesReadFileRequest</summary>

ReadFileRequest names one regular file.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `path` | Yes | `string` | path is the sandbox path of a regular file. |

```yaml
type: object
properties:
  path:
    type: string
    minLength: 1
    description: path is the sandbox path of a regular file.
title: ReadFileRequest
required:
  - path
additionalProperties: false
description: ReadFileRequest names one regular file.
```

</details>

<details id="schema-FilesReadFileResponse">
<summary>FilesReadFileResponse</summary>

ReadFileResponse carries the whole file; HTTP selects data as its raw response body.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `info` | Yes | [`FileInfo`](#schema-FileInfo) | info remains available to the internal adapter; HTTP clients use Stat for metadata. |
| `data` | No | `string` (byte) | data is the complete content; a file above the message cap is refused. |

```yaml
type: object
properties:
  info:
    description: info remains available to the internal adapter; HTTP clients use Stat for metadata.
    $ref: "#/components/schemas/FileInfo"
  data:
    type: string
    format: byte
    description: data is the complete content; a file above the message cap is refused.
title: ReadFileResponse
required:
  - info
additionalProperties: false
description: ReadFileResponse carries the whole file; HTTP selects data as its raw response body.
```

</details>

<details id="schema-FilesRemoveResponse">
<summary>FilesRemoveResponse</summary>

RemoveResponse reports partial recursive progress to every caller: check failed_path,
not the operation status, when Remove stops partway. Effects are never rolled back.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `completedPaths` | No | Array of `string` | completed_paths lists removed paths in traversal order, bounded to the first 1000; completed_count carries the true total when it exceeds that bound. |
| `completedCount` | No | `integer` | completed_count is the total number of paths removed, which may exceed len(completed_paths). |
| `failedPath` | No | `string` | failed_path is the path where a recursive removal stopped; empty on full success. Every path in completed_paths was removed before this one. |

```yaml
type: object
properties:
  completedPaths:
    type: array
    items:
      type: string
    description: |-
      completed_paths lists removed paths in traversal order, bounded to the
      first 1000; completed_count carries the
      true total when it exceeds that bound.
  completedCount:
    type: integer
    description: |-
      completed_count is the total number of paths removed, which may exceed
      len(completed_paths).
  failedPath:
    type: string
    description: |-
      failed_path is the path where a recursive removal stopped; empty on full
      success. Every path in completed_paths was removed before this one.
title: RemoveResponse
description: |-
  RemoveResponse reports partial recursive progress to every caller: check failed_path,
  not the operation status, when Remove stops partway. Effects are never rolled back.
```

</details>

<details id="schema-FilesStatResponse">
<summary>FilesStatResponse</summary>

StatResponse returns file metadata.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `info` | No | [`FileInfo`](#schema-FileInfo) | info is metadata for the path. |

```yaml
type: object
properties:
  info:
    description: info is metadata for the path.
    $ref: "#/components/schemas/FileInfo"
title: StatResponse
description: StatResponse returns file metadata.
```

</details>

<details id="schema-FilesUploadRequest">
<summary>FilesUploadRequest</summary>

UploadRequest is one frame in an upload stream.



```yaml
type: object
oneOf:
  - type: object
    properties:
      data:
        type: string
        format: byte
        description: data appends bytes to the current file.
    title: data
    required:
      - data
  - type: object
    properties:
      header:
        description: |-
          header starts a new file. The first frame, and the first frame after each
          completed file, must be this arm.
        $ref: "#/components/schemas/FileHeader"
    title: header
    required:
      - header
title: UploadRequest
unevaluatedProperties: false
description: UploadRequest is one frame in an upload stream.
```

</details>

<details id="schema-FilesUploadResponse">
<summary>FilesUploadResponse</summary>

UploadResponse reports completed file writes.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `filesWritten` | No | `integer` | files_written is the number of files successfully written. |

```yaml
type: object
properties:
  filesWritten:
    type: integer
    description: files_written is the number of files successfully written.
title: UploadResponse
description: UploadResponse reports completed file writes.
```

</details>

<details id="schema-FilesWriteFileRequest">
<summary>FilesWriteFileRequest</summary>

WriteFileRequest replaces one file.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `header` | Yes | [`FileHeader`](#schema-FileHeader) | header names the path and mode, as an Upload header does. |
| `data` | No | `string` (byte) | data is the complete new content; empty writes an empty file. |

```yaml
type: object
properties:
  header:
    description: header names the path and mode, as an Upload header does.
    $ref: "#/components/schemas/FileHeader"
  data:
    type: string
    format: byte
    description: data is the complete new content; empty writes an empty file.
title: WriteFileRequest
required:
  - header
additionalProperties: false
description: WriteFileRequest replaces one file.
```

</details>

<details id="schema-FilesWriteFileResponse">
<summary>FilesWriteFileResponse</summary>

WriteFileResponse reports the written file.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `info` | Yes | [`FileInfo`](#schema-FileInfo) | info is the file's metadata after the write. |

```yaml
type: object
properties:
  info:
    description: info is the file's metadata after the write.
    $ref: "#/components/schemas/FileInfo"
title: WriteFileResponse
required:
  - info
description: WriteFileResponse reports the written file.
```

</details>

<details id="schema-AttachmentMode">
<summary>AttachmentMode</summary>

AttachmentMode fixes attachment ownership independently of backend placement.

Values:
unspecified: unspecified selects exclusive read-write attachment.
exclusive: exclusive permits one sandbox to attach the volume read-write.



```yaml
type: string
title: AttachmentMode
enum:
  - unspecified
  - exclusive
description: |-
  AttachmentMode fixes attachment ownership independently of backend placement.

  Values:
  unspecified: unspecified selects exclusive read-write attachment.
  exclusive: exclusive permits one sandbox to attach the volume read-write.
```

</details>

<details id="schema-EffectiveManagedVolumes">
<summary>EffectiveManagedVolumes</summary>

EffectiveManagedVolumes records storage that is not deleted with the sandbox.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `attachments` | Yes | Array of [`EffectiveVolumeAttachment`](#schema-EffectiveVolumeAttachment) | attachments retains admission order and stable volume identities. |

```yaml
type: object
properties:
  attachments:
    type: array
    items:
      $ref: "#/components/schemas/EffectiveVolumeAttachment"
    minItems: 1
    description: attachments retains admission order and stable volume identities.
title: EffectiveManagedVolumes
required:
  - attachments
description: EffectiveManagedVolumes records storage that is not deleted with the sandbox.
```

</details>

<details id="schema-EffectiveVolumeAttachment">
<summary>EffectiveVolumeAttachment</summary>

EffectiveVolumeAttachment retains the identity, not a reusable name lookup.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `volume` | No | `string` | volume is the resource name recorded for this reference. |
| `volumeUid` | No | `string or null` | volume_uid retains the observed backing incarnation when the backend provides one. |
| `target` | Yes | `string` | target is the accepted normalized absolute guest path. |
| `mode` | Yes | [`AttachmentMode`](#schema-AttachmentMode) | mode is explicit and survives sandbox stop/start and backend restart. |

```yaml
type: object
properties:
  volume:
    type: string
    description: volume is the resource name recorded for this reference.
  volumeUid:
    type:
      - string
      - "null"
    description: volume_uid retains the observed backing incarnation when the backend provides one.
    readOnly: true
  target:
    type: string
    minLength: 1
    description: target is the accepted normalized absolute guest path.
  mode:
    not:
      enum:
        - unspecified
    description: mode is explicit and survives sandbox stop/start and backend restart.
    $ref: "#/components/schemas/AttachmentMode"
title: EffectiveVolumeAttachment
required:
  - target
  - mode
description: EffectiveVolumeAttachment retains the identity, not a reusable name lookup.
```

</details>

<details id="schema-ManagedVolumeOptions">
<summary>ManagedVolumeOptions</summary>

ManagedVolumeOptions selects existing volumes; it does not create storage.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `attachments` | Yes | Array of [`VolumeAttachment`](#schema-VolumeAttachment) | attachments cannot repeat a volume identity or overlap another mount target. |

```yaml
type: object
properties:
  attachments:
    type: array
    items:
      $ref: "#/components/schemas/VolumeAttachment"
    minItems: 1
    description: attachments cannot repeat a volume identity or overlap another mount target.
title: ManagedVolumeOptions
required:
  - attachments
additionalProperties: false
description: ManagedVolumeOptions selects existing volumes; it does not create storage.
```

</details>

<details id="schema-VolumeAttachment">
<summary>VolumeAttachment</summary>

VolumeAttachment binds an existing owner-scoped volume before guest execution.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `volume` | Yes | `string` | volume is the complete resource name, in the form volumes/{volume}. |
| `target` | Yes | `string` | target is an absolute guest path and cannot overlap another mount target. |
| `mode` | No | [`AttachmentMode`](#schema-AttachmentMode) | mode defaults to exclusive; it never silently enables shared writers. |

```yaml
type: object
properties:
  volume:
    type: string
    pattern: ^volumes/[^/]+$
    description: volume is the complete resource name, in the form volumes/{volume}.
  target:
    type: string
    minLength: 1
    description: target is an absolute guest path and cannot overlap another mount target.
  mode:
    description: mode defaults to exclusive; it never silently enables shared writers.
    $ref: "#/components/schemas/AttachmentMode"
title: VolumeAttachment
required:
  - volume
  - target
additionalProperties: false
description: VolumeAttachment binds an existing owner-scoped volume before guest execution.
```

</details>

<details id="schema-StoredSecretOptions">
<summary>StoredSecretOptions</summary>

StoredSecretOptions resolves owner-scoped references before the workload can run.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `secrets` | Yes | Array of `string` | secrets are complete Secret resource names resolved once at admission. |

```yaml
type: object
properties:
  secrets:
    type: array
    items:
      type: string
      pattern: ^secrets/[^/]+$
    minItems: 1
    description: secrets are complete Secret resource names resolved once at admission.
title: StoredSecretOptions
required:
  - secrets
additionalProperties: false
description: StoredSecretOptions resolves owner-scoped references before the workload can run.
```

</details>

<details id="schema-EffectiveTimeouts">
<summary>EffectiveTimeouts</summary>

EffectiveTimeouts retains the selected behavior and current action deadline.
DELETE and STOP retain a resume timeout; RESTART (and its deprecated alias KEEP) does not use one

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `onTimeout` | Yes | [`TimeoutAction`](#schema-TimeoutAction) | on_timeout remains the accepted disposition through renewal and restart. |
| `autoResume` | No | `boolean` | auto_resume is explicit; endpoint authorization still precedes any wake-up. |
| `initialTimeout` | Yes | [`Duration`](#schema-Duration) | Records the accepted initial TTL, not the recurring interval for RESTART or KEEP. |
| `expiresAt` | No | [`Timestamp`](#schema-Timestamp) | expires_at is present only while a deadline is armed, not before first readiness or while stopped. |
| `resumeTimeout` | No | [`Duration`](#schema-Duration) | Freezes the action's default for later starts; RESTART and KEEP manage their own interval. |

```yaml
type: object
properties:
  onTimeout:
    not:
      enum:
        - unspecified
    description: on_timeout remains the accepted disposition through renewal and restart.
    $ref: "#/components/schemas/TimeoutAction"
  autoResume:
    type: boolean
    description: auto_resume is explicit; endpoint authorization still precedes any wake-up.
  initialTimeout:
    description: Records the accepted initial TTL, not the recurring interval for RESTART or KEEP.
    $ref: "#/components/schemas/Duration"
  expiresAt:
    description: expires_at is present only while a deadline is armed, not before first readiness or while stopped.
    $ref: "#/components/schemas/Timestamp"
  resumeTimeout:
    description: Freezes the action's default for later starts; RESTART and KEEP manage their own interval.
    $ref: "#/components/schemas/Duration"
title: EffectiveTimeouts
required:
  - onTimeout
  - initialTimeout
description: |-
  EffectiveTimeouts retains the selected behavior and current action deadline.
  DELETE and STOP retain a resume timeout; RESTART (and its deprecated alias KEEP) does not use one
```

</details>

<details id="schema-TimeoutAction">
<summary>TimeoutAction</summary>

TimeoutAction selects the lifecycle action taken at the reported deadline.

Values:
unspecified: unspecified is invalid when explicitly supplied.
delete: delete deletes the sandbox at expiry.
stop: stop applies the sandbox's persisted stop guarantee at expiry.
keep: Deprecated alias for the restart action, still accepted on input and output. New clients should use the canonical restart action.
restart: Preserves memory through backend-managed stop/resume cycles at expiry. Requires memory preservation and the always-on entitlement.



```yaml
type: string
title: TimeoutAction
enum:
  - unspecified
  - delete
  - stop
  - keep
  - restart
description: |-
  TimeoutAction selects the lifecycle action taken at the reported deadline.

  Values:
  unspecified: unspecified is invalid when explicitly supplied.
  delete: delete deletes the sandbox at expiry.
  stop: stop applies the sandbox's persisted stop guarantee at expiry.
  keep: Deprecated alias for the restart action, still accepted on input and output. New clients should use the canonical restart action.
  restart: Preserves memory through backend-managed stop/resume cycles at expiry. Requires memory preservation and the always-on entitlement.
```

</details>

<details id="schema-TimeoutOptions">
<summary>TimeoutOptions</summary>

TimeoutOptions selects deadline behavior; presence enables this feature for the sandbox.
an explicit RESTART (or its deprecated alias KEEP) initial timeout must be at least one hour

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `timeout` | No | [`Duration`](#schema-Duration) | timeout is the initial TTL; omission selects the advertised action default. |
| `onTimeout` | No | [`TimeoutAction`](#schema-TimeoutAction) or `null` | on_timeout defaults to DELETE; an explicit UNSPECIFIED is invalid. |
| `autoResume` | No | `boolean or null` | Omission selects the advertised default; false does not disable RESTART or KEEP cycles. |

```yaml
type: object
properties:
  timeout:
    description: timeout is the initial TTL; omission selects the advertised action default.
    $ref: "#/components/schemas/Duration"
  onTimeout:
    oneOf:
      - $ref: "#/components/schemas/TimeoutAction"
      - type: "null"
    not:
      enum:
        - unspecified
    description: on_timeout defaults to DELETE; an explicit UNSPECIFIED is invalid.
  autoResume:
    type:
      - boolean
      - "null"
    description: Omission selects the advertised default; false does not disable RESTART or KEEP cycles.
title: TimeoutOptions
additionalProperties: false
description: |-
  TimeoutOptions selects deadline behavior; presence enables this feature for the sandbox.
  an explicit RESTART (or its deprecated alias KEEP) initial timeout must be at least one hour
```

</details>

<details id="schema-ProcessAttach">
<summary>ProcessAttach</summary>

Attach selects the process and output resume point.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `name` | Yes | `string` | name is the complete resource name, in the form processes/{process}. |
| `resumeFrom` | No | `integer or string` (int64) | resume_from resumes after a received or published stream_sequence; 0 means from the start. Missing output after this point fails with FAILED_PRECONDITION. |

```yaml
type: object
properties:
  name:
    type: string
    pattern: ^processes/[^/]+$
    description: name is the complete resource name, in the form processes/{process}.
  resumeFrom:
    type:
      - integer
      - string
    format: int64
    description: resume_from resumes after a received or published stream_sequence; 0 means from the start. Missing output after this point fails with FAILED_PRECONDITION.
title: Attach
required:
  - name
additionalProperties: false
description: Attach selects the process and output resume point.
```

</details>

<details id="schema-ProcessChunk">
<summary>ProcessChunk</summary>

Chunk carries process output.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `stream` | No | [`StreamType`](#schema-StreamType) | stream identifies stdout, stderr, or PTY output. |
| `data` | No | `string` (byte) | data is output bytes. |
| `streamSequence` | No | `integer or string` (int64) | stream_sequence is one-based and strictly increases per process output stream. |

```yaml
type: object
properties:
  stream:
    description: stream identifies stdout, stderr, or PTY output.
    $ref: "#/components/schemas/StreamType"
  data:
    type: string
    format: byte
    description: data is output bytes.
  streamSequence:
    type:
      - integer
      - string
    format: int64
    description: stream_sequence is one-based and strictly increases per process output stream.
title: Chunk
description: Chunk carries process output.
```

</details>

<details id="schema-CreateProcessRequest">
<summary>CreateProcessRequest</summary>

CreateProcessRequest starts an interactive process.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `cmd` | Yes | Array of `string` | cmd is the argv vector. |
| `env` | No | `object` | env adds environment variables. |
| `workingDir` | No | `string` | working_dir is the process working directory. |
| `user` | No | `string` | user is the user to run as. |
| `pty` | No | [`PtyConfig`](#schema-PtyConfig) | pty requests a pseudo-terminal; absence uses separate stdin/stdout/stderr pipes. |
| `session` | No | `string` | session is an optional opaque session tag stored verbatim; empty means untagged. |
| `label` | No | `string or null` | (OPTIONAL) label is unique within the sandbox; omission retains the server-generated default. |

```yaml
type: object
properties:
  cmd:
    type: array
    items:
      type: string
    minItems: 1
    description: cmd is the argv vector.
  env:
    type: object
    additionalProperties:
      type: string
      title: value
    description: env adds environment variables.
  workingDir:
    type: string
    description: working_dir is the process working directory.
  user:
    type: string
    description: user is the user to run as.
  pty:
    description: pty requests a pseudo-terminal; absence uses separate stdin/stdout/stderr pipes.
    $ref: "#/components/schemas/PtyConfig"
  session:
    type: string
    description: session is an optional opaque session tag stored verbatim; empty means untagged.
  label:
    type:
      - string
      - "null"
    minLength: 1
    description: (OPTIONAL) label is unique within the sandbox; omission retains the server-generated default.
title: CreateProcessRequest
required:
  - cmd
additionalProperties: false
description: CreateProcessRequest starts an interactive process.
```

</details>

<details id="schema-ExecRequest">
<summary>ExecRequest</summary>

ExecRequest runs a command to completion.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `cmd` | Yes | Array of `string` | cmd is the argv vector. |
| `env` | No | `object` | env adds environment variables. |
| `workingDir` | No | `string` | working_dir is the process working directory. |
| `user` | No | `string` | user is the user to run as. |

```yaml
type: object
properties:
  cmd:
    type: array
    items:
      type: string
    minItems: 1
    description: cmd is the argv vector.
  env:
    type: object
    additionalProperties:
      type: string
      title: value
    description: env adds environment variables.
  workingDir:
    type: string
    description: working_dir is the process working directory.
  user:
    type: string
    description: user is the user to run as.
title: ExecRequest
required:
  - cmd
additionalProperties: false
description: ExecRequest runs a command to completion.
```

</details>

<details id="schema-ExecResponse">
<summary>ExecResponse</summary>

ExecResponse returns process output after the command exits.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `exitCode` | No | `integer` (int32) | exit_code is the process exit code. |
| `stdout` | No | `string` (byte) | stdout is captured stdout, up to the capture limit when incomplete. |
| `stderr` | No | `string` (byte) | stderr is captured stderr, up to the capture limit when incomplete. |
| `incomplete` | No | `boolean` | incomplete is true when the combined capture limit forced stdout and/or stderr to discard some of the command's real output; false whenever every byte the command produced was retained, even if that total lands exactly at the limit. exit_code is still the real, final exit code: the command ran to completion, only its captured output was capped. A retry would rerun the command, so this stays a normal successful response rather than an error a client might retry. |

```yaml
type: object
properties:
  exitCode:
    type: integer
    format: int32
    description: exit_code is the process exit code.
  stdout:
    type: string
    format: byte
    description: stdout is captured stdout, up to the capture limit when incomplete.
  stderr:
    type: string
    format: byte
    description: stderr is captured stderr, up to the capture limit when incomplete.
  incomplete:
    type: boolean
    description: |-
      incomplete is true when the combined capture limit forced stdout and/or
      stderr to discard some of the command's real output; false whenever
      every byte the command produced was retained, even if that total lands
      exactly at the limit. exit_code is still the real, final exit code: the
      command ran to completion, only its captured output was capped. A retry
      would rerun the command, so this stays a normal successful response
      rather than an error a client might retry.
title: ExecResponse
description: ExecResponse returns process output after the command exits.
```

</details>

<details id="schema-ProcessExited">
<summary>ProcessExited</summary>

Exited reports process completion.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `exitCode` | No | `integer` (int32) | exit_code is the process exit code. |
| `reason` | No | `string` | reason is backend-readable exit context. |

```yaml
type: object
properties:
  exitCode:
    type: integer
    format: int32
    description: exit_code is the process exit code.
  reason:
    type: string
    description: reason is backend-readable exit context.
title: Exited
description: Exited reports process completion.
```

</details>

<details id="schema-ProcessHeartbeat">
<summary>ProcessHeartbeat</summary>

Heartbeat is an idle keepalive frame.



```yaml
type: object
title: Heartbeat
additionalProperties: false
description: Heartbeat is an idle keepalive frame.
```

</details>

<details id="schema-InteractRequest">
<summary>InteractRequest</summary>

InteractRequest is one frame on the bidirectional interaction stream.



```yaml
type: object
oneOf:
  - type: object
    properties:
      attach:
        description: attach opens or resumes the interaction.
        $ref: "#/components/schemas/ProcessAttach"
    title: attach
    required:
      - attach
  - type: object
    properties:
      closeStdin:
        type: boolean
        description: |-
          close_stdin closes standard input when true; false has no effect, including on a PTY.
          On a PTY, true delivers the terminal's VEOF input through the line discipline and
          does not close the PTY master; subsequent stdin frames remain valid. Terminal mode
          owns interpretation: canonical mode treats VEOF as marking the end of input, while
          raw mode delivers it as an ordinary input byte, so this promises delivery, never
          interpretation. Each true frame on a PTY delivers exactly one VEOF input; unlike the
          pipe path, repeated frames are not a no-op.
    title: close_stdin
    required:
      - closeStdin
  - type: object
    properties:
      resize:
        description: resize changes terminal size; portable-range dimensions on a running process without a PTY return FAILED_PRECONDITION.
        $ref: "#/components/schemas/TerminalSize"
    title: resize
    required:
      - resize
  - type: object
    properties:
      stdin:
        type: string
        format: byte
        description: stdin writes process standard input.
    title: stdin
    required:
      - stdin
title: InteractRequest
unevaluatedProperties: false
description: InteractRequest is one frame on the bidirectional interaction stream.
```

</details>

<details id="schema-ListProcessesResponse">
<summary>ListProcessesResponse</summary>

ListProcessesResponse returns processes and a next-page token.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `processes` | No | Array of [`Process`](#schema-Process) | processes is the current page. |
| `nextPageToken` | No | `string` | next_page_token is empty when there are no more pages. |

```yaml
type: object
properties:
  processes:
    type: array
    items:
      $ref: "#/components/schemas/Process"
    description: processes is the current page.
  nextPageToken:
    type: string
    description: next_page_token is empty when there are no more pages.
title: ListProcessesResponse
description: ListProcessesResponse returns processes and a next-page token.
```

</details>

<details id="schema-ProcessOutput">
<summary>ProcessOutput</summary>

Output is one server frame on an interaction stream. At most one of chunk, exited, heartbeat is set.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `chunk` | No | [`ProcessChunk`](#schema-ProcessChunk) | chunk carries ordered process output. |
| `exited` | No | [`ProcessExited`](#schema-ProcessExited) | exited reports process completion. |
| `heartbeat` | No | [`ProcessHeartbeat`](#schema-ProcessHeartbeat) | heartbeat keeps an idle stream alive and carries no resume position. |

```yaml
type: object
title: Output
unevaluatedProperties: false
description: Output is one server frame on an interaction stream. At most one of chunk, exited, heartbeat is set.
properties:
  chunk:
    description: chunk carries ordered process output.
    $ref: "#/components/schemas/ProcessChunk"
  exited:
    description: exited reports process completion.
    $ref: "#/components/schemas/ProcessExited"
  heartbeat:
    description: heartbeat keeps an idle stream alive and carries no resume position.
    $ref: "#/components/schemas/ProcessHeartbeat"
dependentSchemas:
  chunk:
    not:
      anyOf:
        - required:
            - exited
        - required:
            - heartbeat
  exited:
    not:
      anyOf:
        - required:
            - chunk
        - required:
            - heartbeat
  heartbeat:
    not:
      anyOf:
        - required:
            - chunk
        - required:
            - exited
```

</details>

<details id="schema-Process">
<summary>Process</summary>

Process is the handle returned before interaction begins.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `uid` | No | `string or null` | (IMMUTABLE) uid is the existing process incarnation, also used as the resource name terminal. |
| `pid` | No | `integer` (int32) | pid is the process id inside the sandbox when available. |
| `session` | No | `string` | session is the session tag from the create, empty when untagged. |
| `state` | No | [`ProcessState`](#schema-ProcessState) | state is the process state, concrete in conforming responses. |
| `lastStreamSequence` | No | `integer or string` (int64) | last_stream_sequence is the highest stream_sequence emitted for the process when this report was assembled; 0 before any output. |
| `receiptKind` | No | [`ReceiptKind`](#schema-ReceiptKind) | receipt_kind classifies this response on CreateProcess; every other operation returning a Process, including ListProcesses, always reports unspecified. |
| `exitCode` | No | `integer or null` (int32) | exit_code is the process exit code, set only when state is exited and absent otherwise. |
| `label` | Yes | `string` | (IMMUTABLE) label retains the unique, immutable process label and its generated default. |
| `name` | No | `string` | (IDENTIFIER) name is scoped to the discovered sandbox endpoint API base. |

```yaml
type: object
properties:
  uid:
    type:
      - string
      - "null"
    description: (IMMUTABLE) uid is the existing process incarnation, also used as the resource name terminal.
    readOnly: true
  pid:
    type: integer
    format: int32
    description: pid is the process id inside the sandbox when available.
  session:
    type: string
    description: session is the session tag from the create, empty when untagged.
  state:
    description: state is the process state, concrete in conforming responses.
    $ref: "#/components/schemas/ProcessState"
  lastStreamSequence:
    type:
      - integer
      - string
    format: int64
    description: last_stream_sequence is the highest stream_sequence emitted for the process when this report was assembled; 0 before any output.
  receiptKind:
    description: |-
      receipt_kind classifies this response on CreateProcess; every other operation
      returning a Process, including ListProcesses, always reports unspecified.
    $ref: "#/components/schemas/ReceiptKind"
  exitCode:
    type:
      - integer
      - "null"
    format: int32
    description: exit_code is the process exit code, set only when state is exited and absent otherwise.
  label:
    type: string
    minLength: 1
    description: (IMMUTABLE) label retains the unique, immutable process label and its generated default.
  name:
    type: string
    description: (IDENTIFIER) name is scoped to the discovered sandbox endpoint API base.
title: Process
required:
  - label
description: Process is the handle returned before interaction begins.
```

</details>

<details id="schema-ProcessState">
<summary>ProcessState</summary>

ProcessState is the coarse process lifecycle a process report carries.

Values:
unspecified: unspecified is not a concrete state; conforming responses never use it.
running: running means the process has not exited.
exited: exited means the process has exited.



```yaml
type: string
title: ProcessState
enum:
  - unspecified
  - running
  - exited
description: |-
  ProcessState is the coarse process lifecycle a process report carries.

  Values:
  unspecified: unspecified is not a concrete state; conforming responses never use it.
  running: running means the process has not exited.
  exited: exited means the process has exited.
```

</details>

<details id="schema-PtyConfig">
<summary>PtyConfig</summary>

PtyConfig configures pseudo-terminal creation.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `initialSize` | No | [`TerminalSize`](#schema-TerminalSize) | initial_size requests an initial terminal size; absence selects a backend-defined size. |

```yaml
type: object
properties:
  initialSize:
    description: initial_size requests an initial terminal size; absence selects a backend-defined size.
    $ref: "#/components/schemas/TerminalSize"
title: PtyConfig
additionalProperties: false
description: PtyConfig configures pseudo-terminal creation.
```

</details>

<details id="schema-ReadOutputResponse">
<summary>ReadOutputResponse</summary>

ReadOutputResponse is one page of ordered output.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `chunks` | No | Array of [`ProcessChunk`](#schema-ProcessChunk) | chunks is retained output in stream_sequence order, possibly empty. |
| `nextResumeFrom` | No | `integer or string` (int64) | next_resume_from is the stream_sequence to pass next; equal to resume_from when nothing new was retained. |
| `exited` | No | [`ProcessExited`](#schema-ProcessExited) | exited is set once the process has exited and no retained bytes follow next_resume_from. |

```yaml
type: object
properties:
  chunks:
    type: array
    items:
      $ref: "#/components/schemas/ProcessChunk"
    description: chunks is retained output in stream_sequence order, possibly empty.
  nextResumeFrom:
    type:
      - integer
      - string
    format: int64
    description: |-
      next_resume_from is the stream_sequence to pass next; equal to resume_from
      when nothing new was retained.
  exited:
    description: exited is set once the process has exited and no retained bytes follow next_resume_from.
    $ref: "#/components/schemas/ProcessExited"
title: ReadOutputResponse
description: ReadOutputResponse is one page of ordered output.
```

</details>

<details id="schema-ReceiptKind">
<summary>ReceiptKind</summary>

ReceiptKind classifies whether a CreateProcess response was freshly assembled or
is the retained receipt of an earlier accepted request.

Values:
unspecified: unspecified is not a concrete classification; a backend or a client generated before this field existed reports or reads this. It asserts neither freshness nor replay.
fresh: fresh means this Process was assembled for this specific accepted request; last_stream_sequence reflects live state as of that acceptance.
replay: replay means the backend returned the retained result of an earlier accepted request under the synchronous mutation replay rule; last_stream_sequence is frozen at the value observed during that original acceptance.



```yaml
type: string
title: ReceiptKind
enum:
  - unspecified
  - fresh
  - replay
description: |-
  ReceiptKind classifies whether a CreateProcess response was freshly assembled or
  is the retained receipt of an earlier accepted request.

  Values:
  unspecified: unspecified is not a concrete classification; a backend or a client generated before this field existed reports or reads this. It asserts neither freshness nor replay.
  fresh: fresh means this Process was assembled for this specific accepted request; last_stream_sequence reflects live state as of that acceptance.
  replay: replay means the backend returned the retained result of an earlier accepted request under the synchronous mutation replay rule; last_stream_sequence is frozen at the value observed during that original acceptance.
```

</details>

<details id="schema-ProcessSignal">
<summary>ProcessSignal</summary>

Signal is a process signal.

Values:
unspecified: unspecified is invalid.
term: term requests graceful termination.
kill: kill requests immediate termination.
int: int requests interrupt.
hup: hup requests hangup.
quit: quit requests quit.
usr1: usr1 requests user signal 1.
usr2: usr2 requests user signal 2.



```yaml
type: string
title: Signal
enum:
  - unspecified
  - term
  - kill
  - int
  - hup
  - quit
  - usr1
  - usr2
description: |-
  Signal is a process signal.

  Values:
  unspecified: unspecified is invalid.
  term: term requests graceful termination.
  kill: kill requests immediate termination.
  int: int requests interrupt.
  hup: hup requests hangup.
  quit: quit requests quit.
  usr1: usr1 requests user signal 1.
  usr2: usr2 requests user signal 2.
```

</details>

<details id="schema-SignalResponse">
<summary>SignalResponse</summary>

SignalResponse has no fields.



```yaml
type: object
title: SignalResponse
description: SignalResponse has no fields.
```

</details>

<details id="schema-StreamType">
<summary>StreamType</summary>

StreamType identifies the source of an output chunk.

Values:
unspecified: unspecified is not a concrete output stream.
stdout: stdout is stdout.
stderr: stderr is stderr.
pty: pty is merged PTY output.



```yaml
type: string
title: StreamType
enum:
  - unspecified
  - stdout
  - stderr
  - pty
description: |-
  StreamType identifies the source of an output chunk.

  Values:
  unspecified: unspecified is not a concrete output stream.
  stdout: stdout is stdout.
  stderr: stderr is stderr.
  pty: pty is merged PTY output.
```

</details>

<details id="schema-TerminalSize">
<summary>TerminalSize</summary>

TerminalSize is a terminal rows/columns pair; 1 through 65535 is the portable range.
bounded requires rows and cols in 1..65535

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `rows` | No | `integer` | rows is the terminal row count. |
| `cols` | No | `integer` | cols is the terminal column count. |
| `bounded` | No | `boolean` | bounded requests the portable range enforced on this size, with no opt in to perform first. Unset or false keeps the unenforced behavior. |

```yaml
type: object
properties:
  rows:
    type: integer
    description: rows is the terminal row count.
  cols:
    type: integer
    description: cols is the terminal column count.
  bounded:
    type: boolean
    description: |-
      bounded requests the portable range enforced on this size, with no opt in
      to perform first. Unset or false keeps the unenforced behavior.
title: TerminalSize
additionalProperties: false
description: |-
  TerminalSize is a terminal rows/columns pair; 1 through 65535 is the portable range.
  bounded requires rows and cols in 1..65535
```

</details>

<details id="schema-AddMcpGatewayServerResponse">
<summary>AddMcpGatewayServerResponse</summary>

AddMcpGatewayServerResponse returns the added or already-present server.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `server` | No | [`McpServer`](#schema-McpServer) | server is the gateway server. |

```yaml
type: object
properties:
  server:
    description: server is the gateway server.
    $ref: "#/components/schemas/McpServer"
title: AddMcpGatewayServerResponse
description: AddMcpGatewayServerResponse returns the added or already-present server.
```

</details>

<details id="schema-AttributedRule">
<summary>AttributedRule</summary>

AttributedRule reports a rule and the layer that supplied it.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `network` | No | `string` | network is the network name or CIDR. |
| `layer` | No | [`PolicyLayer`](#schema-PolicyLayer) | layer identifies where the rule came from. |

```yaml
type: object
properties:
  network:
    type: string
    description: network is the network name or CIDR.
  layer:
    description: layer identifies where the rule came from.
    $ref: "#/components/schemas/PolicyLayer"
title: AttributedRule
description: AttributedRule reports a rule and the layer that supplied it.
```

</details>

<details id="schema-CaptureMode">
<summary>CaptureMode</summary>

CaptureMode selects whether captured images include memory.

Values:
unspecified: unspecified means disk-only capture.
all: all captures memory and disk and requires memory-capture support.
disk: disk captures disk only.



```yaml
type: string
title: CaptureMode
enum:
  - unspecified
  - all
  - disk
description: |-
  CaptureMode selects whether captured images include memory.

  Values:
  unspecified: unspecified means disk-only capture.
  all: all captures memory and disk and requires memory-capture support.
  disk: disk captures disk only.
```

</details>

<details id="schema-CreateImageRequest">
<summary>CreateImageRequest</summary>

CreateImageRequest creates a managed image.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `displayName` | Yes | `string` | display_name is a scoped label; the backend assigns the immutable resource ID. |
| `description` | No | `string` | description is caller-supplied text. |
| `startCmd` | No | Array of `string` | start_cmd is used when the image is launched as a raw image sandbox. |
| `readyCmd` | No | Array of `string` | ready_cmd is used to determine readiness for raw image sandboxes. |

```yaml
type: object
allOf:
  - properties:
      displayName:
        type: string
        not:
          type: string
          enum:
            - .
            - ..
        minLength: 1
        pattern: ^[^/]+$
        description: display_name is a scoped label; the backend assigns the immutable resource ID.
      description:
        type: string
        description: description is caller-supplied text.
      startCmd:
        type: array
        items:
          type: string
        description: start_cmd is used when the image is launched as a raw image sandbox.
      readyCmd:
        type: array
        items:
          type: string
        description: ready_cmd is used to determine readiness for raw image sandboxes.
  - oneOf:
      - type: object
        properties:
          fromImage:
            description: from_image registers an entry that will receive pushed OCI content.
            $ref: "#/components/schemas/ImageFromImage"
        title: from_image
        required:
          - fromImage
title: CreateImageRequest
required:
  - displayName
unevaluatedProperties: false
description: CreateImageRequest creates a managed image.
```

</details>

<details id="schema-CreatePortRequestPortInput">
<summary>CreatePortRequestPortInput</summary>



| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `number` | Yes | `integer` (int32) | (IMMUTABLE) number is immutable and unique within the sandbox; duplicates refuse ALREADY_EXISTS. |
| `protocol` | No | `string` | (IMMUTABLE) protocol defaults to TCP when unspecified. |

```yaml
type: object
properties:
  number:
    type: integer
    maximum: 65535
    minimum: 1
    format: int32
    description: (IMMUTABLE) number is immutable and unique within the sandbox; duplicates refuse ALREADY_EXISTS.
  protocol:
    type: string
    title: PortInputProtocol
    enum:
      - unspecified
      - tcp
    description: (IMMUTABLE) protocol defaults to TCP when unspecified.
title: PortInput
required:
  - number
additionalProperties: false
```

</details>

<details id="schema-CreateSandboxRequest">
<summary>CreateSandboxRequest</summary>

CreateSandboxRequest validates typed composition intent before replay and
persists accepted guarantees, distinct enough from legacy create that old
servers cannot silently drop it.
display_name and display_name_prefix are mutually exclusive At most one of image, imageRef is set.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `agent` | No | `string` | agent names an optional agent profile; unsupported profiles fail explicitly. |
| `displayName` | No | `string` | (OPTIONAL) display_name is a scoped label; the backend assigns the immutable resource ID. |
| `displayNamePrefix` | No | `string` | (OPTIONAL) display_name_prefix requests a generated label and excludes display_name. |
| `environment` | No | `object` | environment is non-secret environment metadata. |
| `resources` | No | [`Resources`](#schema-Resources) | resources requests CPU and memory within the selected placement's limits. |
| `pullPolicy` | No | `string` | pull_policy retains shared image-pull semantics and explicit support errors. |
| `kits` | No | Array of [`KitInput`](#schema-KitInput) | kits are the kit inputs to apply before readiness, in composition order. |
| `mcp` | No | [`McpCreateSpec`](#schema-McpCreateSpec) | mcp requests an MCP gateway wired before the workload starts. |
| `networkPolicies` | No | Array of [`NetworkPolicy`](#schema-NetworkPolicy) | network_policies are inline policies materialized in governance and enforced before the first guest process. Omission or an empty list requests no inline policies; unsupported policies fail before create effects. |
| `policyIds` | No | Array of `string` | policy_ids names existing governance policies, never raw policy content. |
| `platform` | No | [`Platform`](#schema-Platform) | platform selects an exact supported, image-compatible tuple when present. explicit platform requires OS and architecture |
| `startupExecution` | No | [`StartupExecution`](#schema-StartupExecution) | startup_execution is ordinary create input and persists for the sandbox lifetime. |
| `labels` | No | `object` | labels are non-secret caller-assigned key/value metadata recorded at admission. This contract has no verb that changes labels afterwards. |
| `features` | No | [`FeatureOptions`](#schema-FeatureOptions) | features carries typed inputs validated before admission can have effects. |
| `image` | No | `string` | image is the managed Image resource name. |
| `imageRef` | No | `string` | image_ref is an external OCI reference, with the existing backend validation. |

```yaml
type: object
allOf:
  - properties:
      agent:
        type: string
        description: agent names an optional agent profile; unsupported profiles fail explicitly.
      displayName:
        type: string
        maxLength: 64
        pattern: ^[a-zA-Z0-9_-]*$
        description: (OPTIONAL) display_name is a scoped label; the backend assigns the immutable resource ID.
      displayNamePrefix:
        type: string
        description: (OPTIONAL) display_name_prefix requests a generated label and excludes display_name.
      environment:
        type: object
        additionalProperties:
          type: string
          title: value
        description: environment is non-secret environment metadata.
      resources:
        description: resources requests CPU and memory within the selected placement's limits.
        $ref: "#/components/schemas/Resources"
      pullPolicy:
        type: string
        title: CreateSandboxRequestPullPolicy
        enum:
          - unspecified
          - always
        description: pull_policy retains shared image-pull semantics and explicit support errors.
      kits:
        type: array
        items:
          $ref: "#/components/schemas/KitInput"
        description: kits are the kit inputs to apply before readiness, in composition order.
      mcp:
        description: mcp requests an MCP gateway wired before the workload starts.
        $ref: "#/components/schemas/McpCreateSpec"
      networkPolicies:
        type: array
        items:
          $ref: "#/components/schemas/NetworkPolicy"
        description: |-
          network_policies are inline policies materialized in governance and enforced before the first guest process.
          Omission or an empty list requests no inline policies; unsupported policies fail before create effects.
      policyIds:
        type: array
        items:
          type: string
          minLength: 1
        description: policy_ids names existing governance policies, never raw policy content.
      platform:
        description: |-
          platform selects an exact supported, image-compatible tuple when present.
          explicit platform requires OS and architecture
        $ref: "#/components/schemas/Platform"
      startupExecution:
        description: |-
          startup_execution is ordinary create input and persists for
          the sandbox lifetime.
        $ref: "#/components/schemas/StartupExecution"
      labels:
        type: object
        maxProperties: 64
        additionalProperties:
          type: string
          title: value
          description: ""
        description: |-
          labels are non-secret caller-assigned key/value metadata recorded at admission.
          This contract has no verb that changes labels afterwards.
      features:
        description: features carries typed inputs validated before admission can have effects.
        $ref: "#/components/schemas/FeatureOptions"
  - properties:
      image:
        type: string
        pattern: ^images/[^/]+$
        description: image is the managed Image resource name.
      imageRef:
        type: string
        description: image_ref is an external OCI reference, with the existing backend validation.
    dependentSchemas:
      image:
        not:
          anyOf:
            - required:
                - imageRef
      imageRef:
        not:
          anyOf:
            - required:
                - image
title: CreateSandboxRequest
unevaluatedProperties: false
description: |-
  CreateSandboxRequest validates typed composition intent before replay and
  persists accepted guarantees, distinct enough from legacy create that old
  servers cannot silently drop it.
  display_name and display_name_prefix are mutually exclusive At most one of image, imageRef is set.
```

</details>

<details id="schema-CreateSecretRequest">
<summary>CreateSecretRequest</summary>

CreateSecretRequest stores new secret material.
service_type must be empty for custom material (server-assigned)

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `displayName` | Yes | `string` | display_name is a scoped label; the backend assigns the immutable resource ID. |
| `serviceType` | No | `string` | service_type identifies the consuming service. |

```yaml
type: object
allOf:
  - properties:
      displayName:
        type: string
        not:
          type: string
          enum:
            - .
            - ..
        minLength: 1
        pattern: ^[^/]+$
        description: display_name is a scoped label; the backend assigns the immutable resource ID.
      serviceType:
        type: string
        description: service_type identifies the consuming service.
  - oneOf:
      - type: object
        properties:
          custom:
            description: |-
              custom stores an opaque value with caller-authored injection metadata.
              CreateSecretRequest.service_type is server-assigned for custom
              material; leave it empty.
            $ref: "#/components/schemas/CustomSecretMaterial"
        title: custom
        required:
          - custom
      - type: object
        properties:
          oauth:
            description: oauth stores OAuth refresh material.
            $ref: "#/components/schemas/OAuthRefreshMaterial"
        title: oauth
        required:
          - oauth
      - type: object
        properties:
          token:
            description: token stores opaque token material.
            $ref: "#/components/schemas/TokenSecretMaterial"
        title: token
        required:
          - token
title: CreateSecretRequest
required:
  - displayName
unevaluatedProperties: false
description: |-
  CreateSecretRequest stores new secret material.
  service_type must be empty for custom material (server-assigned)
```

</details>

<details id="schema-CreateVolumeRequest">
<summary>CreateVolumeRequest</summary>

CreateVolumeRequest creates one volume.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `displayName` | Yes | `string` | display_name is a scoped label; the backend assigns the immutable resource ID. |

```yaml
type: object
properties:
  displayName:
    type: string
    not:
      type: string
      enum:
        - .
        - ..
    minLength: 1
    pattern: ^[^/]+$
    description: display_name is a scoped label; the backend assigns the immutable resource ID.
title: CreateVolumeRequest
required:
  - displayName
additionalProperties: false
description: CreateVolumeRequest creates one volume.
```

</details>

<details id="schema-CredentialFenceReceipt">
<summary>CredentialFenceReceipt</summary>

CredentialFenceReceipt acknowledges one completed generation without granting current use.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `target` | Yes | `string` | target is the complete Sandbox or Secret resource name bound to this receipt. |
| `generation` | Yes | `integer or string` (int64) | generation increases for each accepted material/rule replacement on that target. |

```yaml
type: object
properties:
  target:
    type: string
    pattern: ^(sandboxes|secrets)/[^/]+$
    description: target is the complete Sandbox or Secret resource name bound to this receipt.
  generation:
    exclusiveMinimum: 0
    type:
      - integer
      - string
    format: int64
    description: generation increases for each accepted material/rule replacement on that target.
title: CredentialFenceReceipt
required:
  - target
  - generation
description: CredentialFenceReceipt acknowledges one completed generation without granting current use.
```

</details>

<details id="schema-CustomInjection">
<summary>CustomInjection</summary>

CustomInjection is caller-authored, non-secret metadata for a custom
secret: which exact hosts receive the value, in which header and format,
and how the placeholder is shaped and delivered. Unlike secret material,
reads return it.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `hosts` | Yes | Array of `string` | hosts are exact lowercase DNS names. Wildcards, IP literals, ports, schemes, paths, and hosts a managed integration already serves are rejected. |
| `header` | Yes | `string` | header is the HTTP header the value is written to. Routing, framing, and cookie headers are rejected. |
| `format` | No | `string` | format renders the value with exactly one %s; %% escapes a literal percent, and any other % sequence is rejected. Empty means "%s". |
| `envVar` | No | `string` | env_var optionally names a sandbox environment variable that carries the secret's placeholder (never the value) as a delivery convenience. |
| `placeholderTemplate` | No | `string` | placeholder_template optionally shapes the placeholder for clients that validate credential shape before sending (e.g. "sk-corp-{rand}"). Each {rand} token is expanded server-side; {rand} is the only token. Empty uses the backend default ("sbx-cs-{rand}"). |

```yaml
type: object
properties:
  hosts:
    type: array
    items:
      type: string
      minLength: 1
    maxItems: 10
    minItems: 1
    description: |-
      hosts are exact lowercase DNS names. Wildcards, IP literals, ports,
      schemes, paths, and hosts a managed integration already serves are
      rejected.
  header:
    type: string
    minLength: 1
    description: |-
      header is the HTTP header the value is written to. Routing, framing,
      and cookie headers are rejected.
  format:
    type: string
    description: |-
      format renders the value with exactly one %s; %% escapes a literal
      percent, and any other % sequence is rejected. Empty means "%s".
  envVar:
    type: string
    description: |-
      env_var optionally names a sandbox environment variable that carries the
      secret's placeholder (never the value) as a delivery convenience.
  placeholderTemplate:
    type: string
    description: |-
      placeholder_template optionally shapes the placeholder for clients that
      validate credential shape before sending (e.g. "sk-corp-{rand}"). Each
      {rand} token is expanded server-side; {rand} is the only token. Empty
      uses the backend default ("sbx-cs-{rand}").
title: CustomInjection
required:
  - hosts
  - header
additionalProperties: false
description: |-
  CustomInjection is caller-authored, non-secret metadata for a custom
  secret: which exact hosts receive the value, in which header and format,
  and how the placeholder is shaped and delivered. Unlike secret material,
  reads return it.
```

</details>

<details id="schema-CustomSecretMaterial">
<summary>CustomSecretMaterial</summary>

CustomSecretMaterial is an opaque secret value plus its injection routing.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `value` | Yes | `string` | value is the secret value. |
| `injection` | Yes | [`CustomInjection`](#schema-CustomInjection) | injection is the caller-authored routing. |

```yaml
type: object
properties:
  value:
    type: string
    minLength: 1
    description: value is the secret value.
    writeOnly: true
  injection:
    description: injection is the caller-authored routing.
    $ref: "#/components/schemas/CustomInjection"
title: CustomSecretMaterial
required:
  - value
  - injection
additionalProperties: false
description: CustomSecretMaterial is an opaque secret value plus its injection routing.
```

</details>

<details id="schema-DeleteSecretResponse">
<summary>DeleteSecretResponse</summary>

DeleteSecretResponse acknowledges deletion and optional consumer fencing.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `credentialFence` | No | [`CredentialFenceReceipt`](#schema-CredentialFenceReceipt) | credential_fence remains the original historical receipt after deletion or name reuse. |

```yaml
type: object
properties:
  credentialFence:
    description: credential_fence remains the original historical receipt after deletion or name reuse.
    $ref: "#/components/schemas/CredentialFenceReceipt"
title: DeleteSecretResponse
description: DeleteSecretResponse acknowledges deletion and optional consumer fencing.
```

</details>

<details id="schema-EffectiveNetworkPolicy">
<summary>EffectiveNetworkPolicy</summary>

EffectiveNetworkPolicy is the policy after layering and fail-closed defaults.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `mode` | Yes | [`NetworkPolicyMode`](#schema-NetworkPolicyMode) | mode is the effective mode and is always concrete in conforming responses. |
| `allowNetworks` | No | Array of [`AttributedRule`](#schema-AttributedRule) | allow_networks are effective allowed rules with attribution. |
| `denyNetworks` | No | Array of [`AttributedRule`](#schema-AttributedRule) | deny_networks are effective denied rules with attribution. |

```yaml
type: object
properties:
  mode:
    not:
      enum:
        - unspecified
    description: mode is the effective mode and is always concrete in conforming responses.
    $ref: "#/components/schemas/NetworkPolicyMode"
  allowNetworks:
    type: array
    items:
      $ref: "#/components/schemas/AttributedRule"
    description: allow_networks are effective allowed rules with attribution.
  denyNetworks:
    type: array
    items:
      $ref: "#/components/schemas/AttributedRule"
    description: deny_networks are effective denied rules with attribution.
title: EffectiveNetworkPolicy
required:
  - mode
description: EffectiveNetworkPolicy is the policy after layering and fail-closed defaults.
```

</details>

<details id="schema-EndpointCredential">
<summary>EndpointCredential</summary>

EndpointCredential is returned once with Cache-Control: no-store; never log its token.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `token` | Yes | `string` |  |
| `expireTime` | Yes | [`Timestamp`](#schema-Timestamp) |  |
| `permissions` | Yes | Array of [`Permission`](#schema-Permission) | Allowed values: sandboxesExec, sandboxesFilesRead, sandboxesFilesWrite. |
| `sandbox` | Yes | `string` |  |
| `audience` | Yes | `string` |  |

```yaml
type: object
properties:
  token:
    type: string
    minLength: 1
  expireTime:
    $ref: "#/components/schemas/Timestamp"
  permissions:
    type: array
    items:
      $ref: "#/components/schemas/Permission"
      enum:
        - sandboxesExec
        - sandboxesFilesRead
        - sandboxesFilesWrite
    minItems: 1
    uniqueItems: true
    description: "Allowed values: sandboxesExec, sandboxesFilesRead, sandboxesFilesWrite."
  sandbox:
    type: string
    pattern: ^sandboxes/[^/]+$
  audience:
    type: string
    minLength: 1
title: EndpointCredential
required:
  - token
  - expireTime
  - permissions
  - sandbox
  - audience
description: "EndpointCredential is returned once with Cache-Control: no-store; never log its token."
```

</details>

<details id="schema-EtagMismatch">
<summary>EtagMismatch</summary>

EtagMismatch is the status detail a stale precondition carries, so a client learns the
current version for an explicit conflict-resolution decision.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `currentEtag` | Yes | `string` | current_etag is informational; clients must not silently retry with it. |

```yaml
type: object
properties:
  currentEtag:
    type: string
    minLength: 1
    description: current_etag is informational; clients must not silently retry with it.
title: EtagMismatch
required:
  - currentEtag
description: |-
  EtagMismatch is the status detail a stale precondition carries, so a client learns the
  current version for an explicit conflict-resolution decision.
```

</details>

<details id="schema-EtagMismatchErrorDetail">
<summary>EtagMismatchErrorDetail</summary>

EtagMismatch is the status detail a stale precondition carries, so a client learns the
current version for an explicit conflict-resolution decision.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `currentEtag` | Yes | `string` | current_etag is informational; clients must not silently retry with it. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  currentEtag:
    type: string
    minLength: 1
    description: current_etag is informational; clients must not silently retry with it.
  "@type":
    type: string
    const: type.googleapis.com/docker.sandboxes.v1.EtagMismatch
title: EtagMismatchErrorDetail
required:
  - currentEtag
  - "@type"
description: |-
  EtagMismatch is the status detail a stale precondition carries, so a client learns the
  current version for an explicit conflict-resolution decision.
```

</details>

<details id="schema-ExactDestination">
<summary>ExactDestination</summary>

ExactDestination uses destination.v1 canonical selector values, never legacy glob syntax.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `kind` | Yes | [`ExactDestinationKind`](#schema-ExactDestinationKind) | kind selects a destination.v1 matching rule. |
| `value` | Yes | `string` | value is canonical for its kind and contains no legacy selector syntax. |
| `port` | No | `integer or null` | An absent port matches all numeric ports; zero is invalid when present. |

```yaml
type: object
properties:
  kind:
    not:
      enum:
        - unspecified
    description: kind selects a destination.v1 matching rule.
    $ref: "#/components/schemas/ExactDestinationKind"
  value:
    type: string
    minLength: 1
    description: value is canonical for its kind and contains no legacy selector syntax.
  port:
    type:
      - integer
      - "null"
    maximum: 65535
    minimum: 1
    description: An absent port matches all numeric ports; zero is invalid when present.
title: ExactDestination
required:
  - kind
  - value
description: ExactDestination uses destination.v1 canonical selector values, never legacy glob syntax.
```

</details>

<details id="schema-ExactDestinationKind">
<summary>ExactDestinationKind</summary>

ExactDestinationKind keeps DNS and IP-literal decisions distinct.

Values:
unspecified
dns
dnsLabel
dnsSubtree
ip
cidr



```yaml
type: string
title: ExactDestinationKind
enum:
  - unspecified
  - dns
  - dnsLabel
  - dnsSubtree
  - ip
  - cidr
description: |-
  ExactDestinationKind keeps DNS and IP-literal decisions distinct.

  Values:
  unspecified
  dns
  dnsLabel
  dnsSubtree
  ip
  cidr
```

</details>

<details id="schema-ExactNetworkPolicy">
<summary>ExactNetworkPolicy</summary>

ExactNetworkPolicy is a complete consumer-installed generation; its last node is the root.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `revision` | Yes | `string` | revision binds authority, configuration and resolved conditions to the installed generation. |
| `selectorVersion` | Yes | `string` | selector_version fixes matching semantics across implementations. |
| `nodes` | Yes | Array of [`ExactPolicyNode`](#schema-ExactPolicyNode) | nodes are complete and topologically ordered; the final node is the root. |

```yaml
type: object
properties:
  revision:
    type: string
    minLength: 1
    description: revision binds authority, configuration and resolved conditions to the installed generation.
  selectorVersion:
    type: string
    description: selector_version fixes matching semantics across implementations.
    const: destination.v1
  nodes:
    type: array
    items:
      $ref: "#/components/schemas/ExactPolicyNode"
    maxItems: 256
    minItems: 1
    description: nodes are complete and topologically ordered; the final node is the root.
title: ExactNetworkPolicy
required:
  - revision
  - selectorVersion
  - nodes
description: ExactNetworkPolicy is a complete consumer-installed generation; its last node is the root.
```

</details>

<details id="schema-ExactPolicyConstant">
<summary>ExactPolicyConstant</summary>

ExactPolicyConstant makes both Boolean identities explicit.

Values:
unspecified
true
false



```yaml
type: string
title: ExactPolicyConstant
enum:
  - unspecified
  - "true"
  - "false"
description: |-
  ExactPolicyConstant makes both Boolean identities explicit.

  Values:
  unspecified
  true
  false
```

</details>

<details id="schema-ExactPolicyNode">
<summary>ExactPolicyNode</summary>

ExactPolicyNode is topologically ordered; references must precede this node.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `source` | No | [`ExactPolicySource`](#schema-ExactPolicySource) | source is required on constants and destinations and absent on Boolean operators. |

```yaml
type: object
allOf:
  - properties:
      source:
        description: source is required on constants and destinations and absent on Boolean operators.
        $ref: "#/components/schemas/ExactPolicySource"
  - oneOf:
      - type: object
        properties:
          all:
            description: all intersects every operand.
            $ref: "#/components/schemas/ExactPolicyOperands"
        title: all
        required:
          - all
      - type: object
        properties:
          any:
            description: any unions its operands without discarding enclosing ceilings.
            $ref: "#/components/schemas/ExactPolicyOperands"
        title: any
        required:
          - any
      - type: object
        properties:
          constant:
            not:
              enum:
                - unspecified
            description: constant explicitly allows or denies independently of the destination.
            $ref: "#/components/schemas/ExactPolicyConstant"
        title: constant
        required:
          - constant
      - type: object
        properties:
          destination:
            description: destination matches only its declared DNS or IP kind.
            $ref: "#/components/schemas/ExactDestination"
        title: destination
        required:
          - destination
      - type: object
        properties:
          not:
            type: integer
            description: not negates one earlier node.
        title: not
        required:
          - not
title: ExactPolicyNode
description: ExactPolicyNode is topologically ordered; references must precede this node.
```

</details>

<details id="schema-ExactPolicyOperands">
<summary>ExactPolicyOperands</summary>

ExactPolicyOperands contains earlier node indexes; empty Boolean groups are invalid.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `nodes` | Yes | Array of `integer` | nodes must all precede the enclosing operator. |

```yaml
type: object
properties:
  nodes:
    type: array
    items:
      type: integer
    maxItems: 64
    minItems: 1
    description: nodes must all precede the enclosing operator.
title: ExactPolicyOperands
required:
  - nodes
description: ExactPolicyOperands contains earlier node indexes; empty Boolean groups are invalid.
```

</details>

<details id="schema-ExactPolicyOrigin">
<summary>ExactPolicyOrigin</summary>

ExactPolicyOrigin separates governance from attached and runtime-owned permissions.

Values:
unspecified
organization
owner
attached
kit
gateway
default



```yaml
type: string
title: ExactPolicyOrigin
enum:
  - unspecified
  - organization
  - owner
  - attached
  - kit
  - gateway
  - default
description: |-
  ExactPolicyOrigin separates governance from attached and runtime-owned permissions.

  Values:
  unspecified
  organization
  owner
  attached
  kit
  gateway
  default
```

</details>

<details id="schema-ExactPolicySource">
<summary>ExactPolicySource</summary>

ExactPolicySource identifies the nonsecret authority of a leaf expression.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `origin` | Yes | [`ExactPolicyOrigin`](#schema-ExactPolicyOrigin) | origin distinguishes the authority that supplied this leaf. |
| `id` | Yes | `string` | id identifies that source without exposing credentials or policy content. |

```yaml
type: object
properties:
  origin:
    not:
      enum:
        - unspecified
    description: origin distinguishes the authority that supplied this leaf.
    $ref: "#/components/schemas/ExactPolicyOrigin"
  id:
    type: string
    minLength: 1
    description: id identifies that source without exposing credentials or policy content.
title: ExactPolicySource
required:
  - origin
  - id
description: ExactPolicySource identifies the nonsecret authority of a leaf expression.
```

</details>

<details id="schema-ExchangeCompositionCredentialResponse">
<summary>ExchangeCompositionCredentialResponse</summary>





```yaml
type: object
title: ExchangeCompositionCredentialResponse
```

</details>

<details id="schema-GetNetworkPoliciesResponse">
<summary>GetNetworkPoliciesResponse</summary>

GetNetworkPoliciesResponse binds both views to exact.revision; a partial read fails.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `effective` | Yes | [`EffectiveNetworkPolicy`](#schema-EffectiveNetworkPolicy) |  |
| `exact` | Yes | [`ExactNetworkPolicy`](#schema-ExactNetworkPolicy) |  |

```yaml
type: object
properties:
  effective:
    $ref: "#/components/schemas/EffectiveNetworkPolicy"
  exact:
    $ref: "#/components/schemas/ExactNetworkPolicy"
title: GetNetworkPoliciesResponse
required:
  - effective
  - exact
description: GetNetworkPoliciesResponse binds both views to exact.revision; a partial read fails.
```

</details>

<details id="schema-IdentityExchangeServiceExchangeDockerCredentialRequest">
<summary>IdentityExchangeServiceExchangeDockerCredentialRequest</summary>



| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `idToken` | Yes | `string` | id_token is single-use identity material, never an owner selector. |

```yaml
type: object
properties:
  idToken:
    type: string
    minLength: 1
    description: id_token is single-use identity material, never an owner selector.
    writeOnly: true
title: IdentityExchangeServiceExchangeDockerCredentialRequest
required:
  - idToken
additionalProperties: false
```

</details>

<details id="schema-Image">
<summary>Image</summary>

Image is the detailed managed image view.
uid must equal the immutable resource ID in name

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `uid` | Yes | `string` | (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment. |
| `name` | Yes | `string` | (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name. |
| `description` | No | `string` | description is caller-supplied text. |
| `status` | Yes | [`ImageStatus`](#schema-ImageStatus) | status is the image readiness state. |
| `imageRef` | No | `string` | image_ref is the registry or backend image reference when available. |
| `imageDigest` | No | `string` | image_digest is the content digest when known. |
| `resources` | No | [`Resources`](#schema-Resources) | resources describes the image's preferred runtime resources. |
| `startCmd` | No | Array of `string` | start_cmd is the command used by raw image sandboxes. |
| `readyCmd` | No | Array of `string` | ready_cmd is the readiness command used by raw image sandboxes. |
| `platform` | No | [`Platform`](#schema-Platform) | platform is the image platform. |
| `totalSizeBytes` | No | `integer or string` (int64) | total_size_bytes is the approximate stored image size. |
| `createdAt` | Yes | [`Timestamp`](#schema-Timestamp) | created_at is the creation timestamp. |
| `updatedAt` | Yes | [`Timestamp`](#schema-Timestamp) | updated_at is the last update timestamp. |
| `failure` | No | [`Error`](#schema-Error) | failure is set when status is FAILED. |
| `source` | No | [`ImageSource`](#schema-ImageSource) | source records how the image was created. |
| `captureMode` | No | [`CaptureMode`](#schema-CaptureMode) | capture_mode records whether captured content includes memory. |
| `etag` | Yes | `string` | etag identifies the observed version of this resource. It is opaque and strong, changes on every visible change, and is what a mutation sends as its precondition. |
| `pushTarget` | No | [`ImagePushTarget`](#schema-ImagePushTarget) | push_target appears only in from_image create responses awaiting upload, never reads. |
| `displayName` | No | `string` | display_name is a scoped label and never selects the resource. |

```yaml
type: object
properties:
  uid:
    type:
      - string
    minLength: 1
    description: (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment.
    readOnly: true
  name:
    type: string
    pattern: ^images/[^/]+$
    description: (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name.
  description:
    type: string
    description: description is caller-supplied text.
  status:
    not:
      enum:
        - unspecified
    description: status is the image readiness state.
    readOnly: true
    $ref: "#/components/schemas/ImageStatus"
  imageRef:
    type: string
    description: image_ref is the registry or backend image reference when available.
  imageDigest:
    type: string
    description: image_digest is the content digest when known.
  resources:
    description: resources describes the image's preferred runtime resources.
    $ref: "#/components/schemas/Resources"
  startCmd:
    type: array
    items:
      type: string
    description: start_cmd is the command used by raw image sandboxes.
  readyCmd:
    type: array
    items:
      type: string
    description: ready_cmd is the readiness command used by raw image sandboxes.
  platform:
    description: platform is the image platform.
    $ref: "#/components/schemas/Platform"
  totalSizeBytes:
    type:
      - integer
      - string
    format: int64
    description: total_size_bytes is the approximate stored image size.
  createdAt:
    description: created_at is the creation timestamp.
    $ref: "#/components/schemas/Timestamp"
  updatedAt:
    description: updated_at is the last update timestamp.
    $ref: "#/components/schemas/Timestamp"
  failure:
    description: failure is set when status is FAILED.
    $ref: "#/components/schemas/Error"
  source:
    description: source records how the image was created.
    $ref: "#/components/schemas/ImageSource"
  captureMode:
    description: capture_mode records whether captured content includes memory.
    $ref: "#/components/schemas/CaptureMode"
  etag:
    type: string
    minLength: 1
    pattern: ^"[^"\x00-\x20\x7f]*"$
    description: |-
      etag identifies the observed version of this resource. It is opaque and strong,
      changes on every visible change, and is what a mutation sends as its precondition.
    readOnly: true
  pushTarget:
    description: push_target appears only in from_image create responses awaiting upload, never reads.
    readOnly: true
    $ref: "#/components/schemas/ImagePushTarget"
  displayName:
    type: string
    description: display_name is a scoped label and never selects the resource.
title: Image
required:
  - uid
  - name
  - status
  - createdAt
  - updatedAt
  - etag
description: |-
  Image is the detailed managed image view.
  uid must equal the immutable resource ID in name
```

</details>

<details id="schema-ImageBlobRef">
<summary>ImageBlobRef</summary>

ImageBlobRef describes one pullable blob.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `digest` | No | `string` | digest is the content digest. |
| `mediaType` | No | `string` | media_type is the OCI media type. |
| `size` | No | `integer or string` (int64) | size is the blob size in bytes. |
| `url` | No | `string` | url is a short-lived pull URL. |

```yaml
type: object
properties:
  digest:
    type: string
    description: digest is the content digest.
  mediaType:
    type: string
    description: media_type is the OCI media type.
  size:
    type:
      - integer
      - string
    format: int64
    description: size is the blob size in bytes.
  url:
    type: string
    description: url is a short-lived pull URL.
title: ImageBlobRef
description: ImageBlobRef describes one pullable blob.
```

</details>

<details id="schema-ImageFromImage">
<summary>ImageFromImage</summary>

ImageFromImage registers a target for pushed OCI content.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `resources` | No | [`Resources`](#schema-Resources) | resources describes preferred runtime resources for the resulting image. |
| `platform` | No | [`Platform`](#schema-Platform) | platform is the platform for the pushed content. |
| `captureMode` | No | `string` | capture_mode must be unspecified or disk; pushed OCI content cannot include memory. |
| `platformHint` | No | [`Platform`](#schema-Platform) | platform_hint picks among a multi-platform push when platform is absent. A hint the backend cannot honor is ignored, never an error. |
| `kitStartup` | No | [`KitStartup`](#schema-KitStartup) | Startup commands copied from the source sandbox. Services with startup replay run these commands when a sandbox created from the image boots. Services without startup replay accept valid commands without running them. |

```yaml
type: object
properties:
  resources:
    description: resources describes preferred runtime resources for the resulting image.
    $ref: "#/components/schemas/Resources"
  platform:
    description: platform is the platform for the pushed content.
    $ref: "#/components/schemas/Platform"
  captureMode:
    type: string
    title: ImageFromImageCaptureMode
    enum:
      - unspecified
      - disk
    description: capture_mode must be unspecified or disk; pushed OCI content cannot include memory.
    not:
      enum:
        - all
  platformHint:
    description: |-
      platform_hint picks among a multi-platform push when platform is absent.
      A hint the backend cannot honor is ignored, never an error.
    $ref: "#/components/schemas/Platform"
  kitStartup:
    description: Startup commands copied from the source sandbox. Services with startup replay run these commands when a sandbox created from the image boots. Services without startup replay accept valid commands without running them.
    $ref: "#/components/schemas/KitStartup"
title: ImageFromImage
additionalProperties: false
description: ImageFromImage registers a target for pushed OCI content.
```

</details>

<details id="schema-ImagePullSpec">
<summary>ImagePullSpec</summary>

ImagePullSpec describes short-lived pull material for placement.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `manifestDigest` | No | `string` | manifest_digest is the manifest digest. |
| `manifestMediaType` | No | `string` | manifest_media_type is the manifest media type. |
| `config` | No | [`ImageBlobRef`](#schema-ImageBlobRef) | config is the image config blob. |
| `layers` | No | Array of [`ImageBlobRef`](#schema-ImageBlobRef) | layers are the image layer blobs. |
| `expiresAt` | No | [`Timestamp`](#schema-Timestamp) | expires_at is when the pull material expires. |

```yaml
type: object
properties:
  manifestDigest:
    type: string
    description: manifest_digest is the manifest digest.
  manifestMediaType:
    type: string
    description: manifest_media_type is the manifest media type.
  config:
    description: config is the image config blob.
    $ref: "#/components/schemas/ImageBlobRef"
  layers:
    type: array
    items:
      $ref: "#/components/schemas/ImageBlobRef"
    description: layers are the image layer blobs.
  expiresAt:
    description: expires_at is when the pull material expires.
    $ref: "#/components/schemas/Timestamp"
title: ImagePullSpec
description: ImagePullSpec describes short-lived pull material for placement.
```

</details>

<details id="schema-ImagePushTarget">
<summary>ImagePushTarget</summary>

ImagePushTarget carries short-lived upload credentials.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `imageRef` | No | `string` | image_ref is the registry reference to push. |
| `accessToken` | No | `string` | access_token is the short-lived push credential. |
| `accessTokenExpiresAt` | No | [`Timestamp`](#schema-Timestamp) | access_token_expires_at is the token expiry. |

```yaml
type: object
properties:
  imageRef:
    type: string
    description: image_ref is the registry reference to push.
  accessToken:
    type: string
    description: access_token is the short-lived push credential.
  accessTokenExpiresAt:
    description: access_token_expires_at is the token expiry.
    $ref: "#/components/schemas/Timestamp"
title: ImagePushTarget
description: ImagePushTarget carries short-lived upload credentials.
```

</details>

<details id="schema-ImageSource">
<summary>ImageSource</summary>

ImageSource records how the image was created.

Values:
unspecified: unspecified means the source is unknown or not reported.
image: image means the image came from pushed or imported image bytes.
sandbox: sandbox means the image came from a sandbox capture.



```yaml
type: string
title: ImageSource
enum:
  - unspecified
  - image
  - sandbox
description: |-
  ImageSource records how the image was created.

  Values:
  unspecified: unspecified means the source is unknown or not reported.
  image: image means the image came from pushed or imported image bytes.
  sandbox: sandbox means the image came from a sandbox capture.
```

</details>

<details id="schema-ImageStatus">
<summary>ImageStatus</summary>

ImageStatus is the image readiness state.

Values:
unspecified: unspecified is never a concrete known state.
pending: pending means an asynchronous build has not started or reported progress.
waitingForPush: waitingForPush means the backend is waiting for OCI content upload.
preparing: preparing means the backend is preparing pushed or captured content.
completed: completed means the image is ready to use.
failed: failed means image preparation reached a terminal failure.
deleted: deleted means the backend exposes a retained delete tombstone.



```yaml
type: string
title: ImageStatus
enum:
  - unspecified
  - pending
  - waitingForPush
  - preparing
  - completed
  - failed
  - deleted
description: |-
  ImageStatus is the image readiness state.

  Values:
  unspecified: unspecified is never a concrete known state.
  pending: pending means an asynchronous build has not started or reported progress.
  waitingForPush: waitingForPush means the backend is waiting for OCI content upload.
  preparing: preparing means the backend is preparing pushed or captured content.
  completed: completed means the image is ready to use.
  failed: failed means image preparation reached a terminal failure.
  deleted: deleted means the backend exposes a retained delete tombstone.
```

</details>

<details id="schema-ImageSummary">
<summary>ImageSummary</summary>

ImageSummary is the list view of an image.
uid must equal the immutable resource ID in name

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `uid` | Yes | `string` | (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment. |
| `name` | Yes | `string` | name ends in the immutable backing ID and never changes with display_name. |
| `status` | Yes | [`ImageStatus`](#schema-ImageStatus) | status is the image readiness state. |
| `source` | No | [`ImageSource`](#schema-ImageSource) | source records how the image was created. |
| `totalSizeBytes` | No | `integer or string` (int64) | total_size_bytes is the approximate stored image size. |
| `createdAt` | Yes | [`Timestamp`](#schema-Timestamp) | created_at is the creation timestamp. |
| `updatedAt` | Yes | [`Timestamp`](#schema-Timestamp) | updated_at is the last update timestamp. |
| `imageRef` | No | `string` | image_ref is the registry or backend image reference when available. |
| `imageDigest` | No | `string` | image_digest is the content digest when known. |
| `captureMode` | No | [`CaptureMode`](#schema-CaptureMode) | capture_mode records whether captured content includes memory. |
| `platform` | No | [`Platform`](#schema-Platform) | platform is the image platform. |
| `displayName` | No | `string` | display_name is a scoped label and never selects the resource. |

```yaml
type: object
properties:
  uid:
    type:
      - string
    minLength: 1
    description: (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment.
    readOnly: true
  name:
    type: string
    pattern: ^images/[^/]+$
    description: name ends in the immutable backing ID and never changes with display_name.
  status:
    not:
      enum:
        - unspecified
    description: status is the image readiness state.
    $ref: "#/components/schemas/ImageStatus"
  source:
    description: source records how the image was created.
    $ref: "#/components/schemas/ImageSource"
  totalSizeBytes:
    type:
      - integer
      - string
    format: int64
    description: total_size_bytes is the approximate stored image size.
  createdAt:
    description: created_at is the creation timestamp.
    $ref: "#/components/schemas/Timestamp"
  updatedAt:
    description: updated_at is the last update timestamp.
    $ref: "#/components/schemas/Timestamp"
  imageRef:
    type: string
    description: image_ref is the registry or backend image reference when available.
  imageDigest:
    type: string
    description: image_digest is the content digest when known.
  captureMode:
    description: capture_mode records whether captured content includes memory.
    $ref: "#/components/schemas/CaptureMode"
  platform:
    description: platform is the image platform.
    $ref: "#/components/schemas/Platform"
  displayName:
    type: string
    description: display_name is a scoped label and never selects the resource.
title: ImageSummary
required:
  - uid
  - name
  - status
  - createdAt
  - updatedAt
description: |-
  ImageSummary is the list view of an image.
  uid must equal the immutable resource ID in name
```

</details>

<details id="schema-IssueSSHCertResponse">
<summary>IssueSSHCertResponse</summary>

IssueSSHCertResponse returns signed SSH client material, never a private key.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `certificate` | No | `string` | certificate is the signed SSH certificate. |
| `sshConfig` | No | `string` | ssh_config optionally renders the connection fields as OpenSSH client configuration; the structured fields are canonical. |
| `knownHosts` | No | `string` | known_hosts is host key material that authenticates the SSH endpoint. |
| `expiresAt` | No | [`Timestamp`](#schema-Timestamp) | expires_at is the instant after which the certificate is invalid. |
| `host` | No | `string` | host is the SSH endpoint host to connect to. |
| `port` | No | `integer` | port is the SSH endpoint port to connect to. |
| `username` | No | `string` | username is the SSH username to present; it carries no authorization. |

```yaml
type: object
properties:
  certificate:
    type: string
    description: certificate is the signed SSH certificate.
  sshConfig:
    type: string
    description: |-
      ssh_config optionally renders the connection fields as OpenSSH client
      configuration; the structured fields are canonical.
  knownHosts:
    type: string
    description: known_hosts is host key material that authenticates the SSH endpoint.
  expiresAt:
    description: expires_at is the instant after which the certificate is invalid.
    $ref: "#/components/schemas/Timestamp"
  host:
    type: string
    description: host is the SSH endpoint host to connect to.
  port:
    type: integer
    description: port is the SSH endpoint port to connect to.
  username:
    type: string
    description: username is the SSH username to present; it carries no authorization.
title: IssueSSHCertResponse
description: IssueSSHCertResponse returns signed SSH client material, never a private key.
```

</details>

<details id="schema-KitArtifactInput">
<summary>KitArtifactInput</summary>

KitArtifactInput carries a v2 artifact inline or by staged digest.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `ref` | Yes | [`KitRef`](#schema-KitRef) |  |

```yaml
type: object
allOf:
  - properties:
      ref:
        $ref: "#/components/schemas/KitRef"
  - oneOf:
      - type: object
        properties:
          inline:
            type: string
            format: byte
            description: inline is the artifact in the pinned JSON encoding.
        title: inline
        required:
          - inline
title: KitArtifactInput
required:
  - ref
unevaluatedProperties: false
description: KitArtifactInput carries a v2 artifact inline or by staged digest.
```

</details>

<details id="schema-KitFieldDisposition">
<summary>KitFieldDisposition</summary>

KitFieldDisposition describes a declaration that was not applied.

Values:
unspecified: unspecified is invalid in a report.
skipped: skipped identifies a contract-permitted optional skip.
refused: refused prevents guest or kit effects.



```yaml
type: string
title: KitFieldDisposition
enum:
  - unspecified
  - skipped
  - refused
description: |-
  KitFieldDisposition describes a declaration that was not applied.

  Values:
  unspecified: unspecified is invalid in a report.
  skipped: skipped identifies a contract-permitted optional skip.
  refused: refused prevents guest or kit effects.
```

</details>

<details id="schema-KitFieldOutcome">
<summary>KitFieldOutcome</summary>

KitFieldOutcome identifies one non-applied declaration without exposing its value.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `kitIndex` | No | `integer` | kit_index is the zero-based position in the originating request's kits. |
| `ref` | Yes | `string` | ref preserves the exact requested kit reference. |
| `digest` | No | `string` | digest is the resolved content identity when resolution succeeded. |
| `fieldPath` | Yes | `string` | field_path is an RFC 6901 pointer into the pinned JSON artifact. |
| `disposition` | Yes | [`KitFieldDisposition`](#schema-KitFieldDisposition) | disposition never claims that an accepted declaration has executed. |
| `reason` | Yes | [`KitFieldReason`](#schema-KitFieldReason) | reason is a bounded category, not caller-controlled diagnostic content. |
| `reasonText` | No | `string` | reason_text contains nonsecret explanatory text, never values or command content. |

```yaml
type: object
properties:
  kitIndex:
    exclusiveMaximum: 64
    type: integer
    description: kit_index is the zero-based position in the originating request's kits.
  ref:
    type: string
    description: ref preserves the exact requested kit reference.
  digest:
    type: string
    description: digest is the resolved content identity when resolution succeeded.
  fieldPath:
    type: string
    description: field_path is an RFC 6901 pointer into the pinned JSON artifact.
  disposition:
    not:
      enum:
        - unspecified
    description: disposition never claims that an accepted declaration has executed.
    $ref: "#/components/schemas/KitFieldDisposition"
  reason:
    not:
      enum:
        - unspecified
    description: reason is a bounded category, not caller-controlled diagnostic content.
    $ref: "#/components/schemas/KitFieldReason"
  reasonText:
    type: string
    description: reason_text contains nonsecret explanatory text, never values or command content.
title: KitFieldOutcome
required:
  - ref
  - fieldPath
  - disposition
  - reason
description: KitFieldOutcome identifies one non-applied declaration without exposing its value.
```

</details>

<details id="schema-KitFieldOutcomeReport">
<summary>KitFieldOutcomeReport</summary>

KitFieldOutcomeReport records settled diagnostics on a sandbox or an admission refusal.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `complete` | No | `boolean` | complete means every refused or skipped declaration is accounted for; absence proves nothing. |
| `outcomes` | No | Array of [`KitFieldOutcome`](#schema-KitFieldOutcome) | outcomes follow request index then UTF-8 field-path order and are never truncated. |

```yaml
type: object
properties:
  complete:
    type: boolean
    description: complete means every refused or skipped declaration is accounted for; absence proves nothing.
  outcomes:
    type: array
    items:
      $ref: "#/components/schemas/KitFieldOutcome"
    maxItems: 256
    description: outcomes follow request index then UTF-8 field-path order and are never truncated.
title: KitFieldOutcomeReport
description: KitFieldOutcomeReport records settled diagnostics on a sandbox or an admission refusal.
```

</details>

<details id="schema-KitFieldOutcomeReportErrorDetail">
<summary>KitFieldOutcomeReportErrorDetail</summary>

KitFieldOutcomeReport records settled diagnostics on a sandbox or an admission refusal.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `complete` | No | `boolean` | complete means every refused or skipped declaration is accounted for; absence proves nothing. |
| `outcomes` | No | Array of [`KitFieldOutcome`](#schema-KitFieldOutcome) | outcomes follow request index then UTF-8 field-path order and are never truncated. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  complete:
    type: boolean
    description: complete means every refused or skipped declaration is accounted for; absence proves nothing.
  outcomes:
    type: array
    items:
      $ref: "#/components/schemas/KitFieldOutcome"
    maxItems: 256
    description: outcomes follow request index then UTF-8 field-path order and are never truncated.
  "@type":
    type: string
    const: type.googleapis.com/docker.sandboxes.v1.KitFieldOutcomeReport
title: KitFieldOutcomeReportErrorDetail
description: KitFieldOutcomeReport records settled diagnostics on a sandbox or an admission refusal.
required:
  - "@type"
```

</details>

<details id="schema-KitFieldReason">
<summary>KitFieldReason</summary>

KitFieldReason provides stable, nonsecret diagnostic categories.

Values:
unspecified: unspecified is invalid in a report.
unsupported: unsupported means the backend cannot honor the declaration.
unknownField: unknownField means the pinned schema does not recognize the declaration.
credentialUnavailable: credentialUnavailable means the declared credential cannot be injected.



```yaml
type: string
title: KitFieldReason
enum:
  - unspecified
  - unsupported
  - unknownField
  - credentialUnavailable
description: |-
  KitFieldReason provides stable, nonsecret diagnostic categories.

  Values:
  unspecified: unspecified is invalid in a report.
  unsupported: unsupported means the backend cannot honor the declaration.
  unknownField: unknownField means the pinned schema does not recognize the declaration.
  credentialUnavailable: credentialUnavailable means the declared credential cannot be injected.
```

</details>

<details id="schema-KitInput">
<summary>KitInput</summary>

KitInput is one kit in either shape.



```yaml
type: object
oneOf:
  - type: object
    properties:
      artifact:
        $ref: "#/components/schemas/KitArtifactInput"
    title: artifact
    required:
      - artifact
title: KitInput
unevaluatedProperties: false
description: KitInput is one kit in either shape.
```

</details>

<details id="schema-KitStartup">
<summary>KitStartup</summary>

Startup commands retained when creating an image. Services with startup replay run them at each boot; services without it accept valid commands without running them.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `kits` | Yes | Array of `string` | kits names the kit identities the commands derive from, in composition order. Required. |
| `commands` | Yes | Array of [`KitStartupCommand`](#schema-KitStartupCommand) | commands run in order at sandbox boot, on a backend that replays them. |

```yaml
type: object
properties:
  kits:
    type: array
    items:
      type: string
      minLength: 1
    minItems: 1
    description: |-
      kits names the kit identities the commands derive from, in composition
      order. Required.
  commands:
    type: array
    items:
      $ref: "#/components/schemas/KitStartupCommand"
    minItems: 1
    description: commands run in order at sandbox boot, on a backend that replays them.
title: KitStartup
required:
  - kits
  - commands
additionalProperties: false
description: Startup commands retained when creating an image. Services with startup replay run them at each boot; services without it accept valid commands without running them.
```

</details>

<details id="schema-KitStartupCommand">
<summary>KitStartupCommand</summary>

KitStartupCommand is one startup command in a KitStartup set, run in order
by a backend that replays it at sandbox boot.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `argv` | Yes | Array of `string` | argv is the command and its arguments, exec-style. An argument other than the executable itself may be empty (an intentional empty string, as in `printf %s ""`); only the vector itself is required non-empty. |
| `user` | No | `string` | user runs the command as this user; empty means the backend default. |
| `background` | No | `boolean` | background launches the command detached instead of waiting for it. |

```yaml
type: object
properties:
  argv:
    type: array
    items:
      type: string
    minItems: 1
    description: |-
      argv is the command and its arguments, exec-style. An argument other than
      the executable itself may be empty (an intentional empty string, as in
      `printf %s ""`); only the vector itself is required non-empty.
  user:
    type: string
    description: user runs the command as this user; empty means the backend default.
  background:
    type: boolean
    description: background launches the command detached instead of waiting for it.
title: KitStartupCommand
required:
  - argv
additionalProperties: false
description: |-
  KitStartupCommand is one startup command in a KitStartup set, run in order
  by a backend that replays it at sandbox boot.
```

</details>

<details id="schema-ListImagesResponse">
<summary>ListImagesResponse</summary>

ListImagesResponse returns image results and a next-page token.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `images` | No | Array of [`ImageSummary`](#schema-ImageSummary) | images is the current page. |
| `nextPageToken` | No | `string` | next_page_token is empty when there are no more pages. |

```yaml
type: object
properties:
  images:
    type: array
    items:
      $ref: "#/components/schemas/ImageSummary"
    description: images is the current page.
  nextPageToken:
    type: string
    description: next_page_token is empty when there are no more pages.
title: ListImagesResponse
description: ListImagesResponse returns image results and a next-page token.
```

</details>

<details id="schema-ListPolicyLogEntriesResponse">
<summary>ListPolicyLogEntriesResponse</summary>

ListPolicyLogEntriesResponse returns policy log entries and a next-page token.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `entries` | No | Array of [`PolicyLogEntry`](#schema-PolicyLogEntry) | entries is the current page. |
| `nextPageToken` | No | `string` | next_page_token is empty when there are no more pages. |

```yaml
type: object
properties:
  entries:
    type: array
    items:
      $ref: "#/components/schemas/PolicyLogEntry"
    description: entries is the current page.
  nextPageToken:
    type: string
    description: next_page_token is empty when there are no more pages.
title: ListPolicyLogEntriesResponse
description: ListPolicyLogEntriesResponse returns policy log entries and a next-page token.
```

</details>

<details id="schema-ListPortsResponse">
<summary>ListPortsResponse</summary>

ListPortsResponse returns currently published ports.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `published` | No | Array of [`Port`](#schema-Port) | published is ordered by number ascending. |

```yaml
type: object
properties:
  published:
    type: array
    items:
      $ref: "#/components/schemas/Port"
    description: published is ordered by number ascending.
title: ListPortsResponse
description: ListPortsResponse returns currently published ports.
```

</details>

<details id="schema-ListSandboxesResponse">
<summary>ListSandboxesResponse</summary>

ListSandboxesResponse returns sandbox results and a next-page token.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `sandboxes` | No | Array of [`Sandbox`](#schema-Sandbox) | sandboxes is the current page. |
| `nextPageToken` | No | `string` | next_page_token is empty when there are no more pages. |

```yaml
type: object
properties:
  sandboxes:
    type: array
    items:
      $ref: "#/components/schemas/Sandbox"
    description: sandboxes is the current page.
  nextPageToken:
    type: string
    description: next_page_token is empty when there are no more pages.
title: ListSandboxesResponse
description: ListSandboxesResponse returns sandbox results and a next-page token.
```

</details>

<details id="schema-ListSecretsResponse">
<summary>ListSecretsResponse</summary>

ListSecretsResponse returns secrets and a next-page token.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `secrets` | No | Array of [`Secret`](#schema-Secret) | secrets is the current page. |
| `nextPageToken` | No | `string` | next_page_token is empty when there are no more pages. |

```yaml
type: object
properties:
  secrets:
    type: array
    items:
      $ref: "#/components/schemas/Secret"
    description: secrets is the current page.
  nextPageToken:
    type: string
    description: next_page_token is empty when there are no more pages.
title: ListSecretsResponse
description: ListSecretsResponse returns secrets and a next-page token.
```

</details>

<details id="schema-ListSnapshotsResponse">
<summary>ListSnapshotsResponse</summary>

ListSnapshotsResponse returns snapshot results and a next-page token.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `snapshots` | No | Array of [`SnapshotSummary`](#schema-SnapshotSummary) | snapshots is the current page. |
| `nextPageToken` | No | `string` | next_page_token is empty when there are no more pages. |

```yaml
type: object
properties:
  snapshots:
    type: array
    items:
      $ref: "#/components/schemas/SnapshotSummary"
    description: snapshots is the current page.
  nextPageToken:
    type: string
    description: next_page_token is empty when there are no more pages.
title: ListSnapshotsResponse
description: ListSnapshotsResponse returns snapshot results and a next-page token.
```

</details>

<details id="schema-ListVolumesResponse">
<summary>ListVolumesResponse</summary>

ListVolumesResponse returns volumes and a next-page token.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `volumes` | No | Array of [`Volume`](#schema-Volume) | volumes is the current page. |
| `nextPageToken` | No | `string` | next_page_token is empty when there are no more pages. |

```yaml
type: object
properties:
  volumes:
    type: array
    items:
      $ref: "#/components/schemas/Volume"
    description: volumes is the current page.
  nextPageToken:
    type: string
    description: next_page_token is empty when there are no more pages.
title: ListVolumesResponse
description: ListVolumesResponse returns volumes and a next-page token.
```

</details>

<details id="schema-McpAuthorization">
<summary>McpAuthorization</summary>

McpAuthorization is the current authorization incarnation for one named upstream.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `server` | No | `string` | (OPTIONAL) server is the complete resource name, in the form mcp-servers/{server}. |
| `uid` | No | `string or null` | (IMMUTABLE) uid is the existing immutable backing identity, omitted when none is available. |
| `status` | No | [`McpAuthorizationStatus`](#schema-McpAuthorizationStatus) |  |
| `authorizationUrl` | No | `string` | authorization_url is present only while interaction is pending. |
| `createdAt` | No | [`Timestamp`](#schema-Timestamp) |  |
| `updatedAt` | No | [`Timestamp`](#schema-Timestamp) |  |
| `expiresAt` | No | [`Timestamp`](#schema-Timestamp) | expires_at bounds the pending flow; expiry records FAILED on this incarnation. |
| `failure` | No | [`Error`](#schema-Error) |  |
| `etag` | No | `string` | etag identifies this incarnation and changes with every visible authorization update. |
| `name` | No | `string` | (IDENTIFIER) name is the complete singleton resource name. |

```yaml
type: object
properties:
  server:
    type: string
    pattern: ^(?:mcp-servers/[^/]+)?$
    description: (OPTIONAL) server is the complete resource name, in the form mcp-servers/{server}.
  uid:
    type:
      - string
      - "null"
    description: (IMMUTABLE) uid is the existing immutable backing identity, omitted when none is available.
    readOnly: true
  status:
    readOnly: true
    $ref: "#/components/schemas/McpAuthorizationStatus"
  authorizationUrl:
    type: string
    description: authorization_url is present only while interaction is pending.
    readOnly: true
  createdAt:
    readOnly: true
    $ref: "#/components/schemas/Timestamp"
  updatedAt:
    readOnly: true
    $ref: "#/components/schemas/Timestamp"
  expiresAt:
    description: expires_at bounds the pending flow; expiry records FAILED on this incarnation.
    readOnly: true
    $ref: "#/components/schemas/Timestamp"
  failure:
    readOnly: true
    $ref: "#/components/schemas/Error"
  etag:
    type: string
    description: etag identifies this incarnation and changes with every visible authorization update.
    readOnly: true
  name:
    type: string
    description: (IDENTIFIER) name is the complete singleton resource name.
title: McpAuthorization
description: McpAuthorization is the current authorization incarnation for one named upstream.
```

</details>

<details id="schema-McpAuthorizationStatus">
<summary>McpAuthorizationStatus</summary>

McpAuthorizationStatus reports the current upstream authorization outcome.

Values:
unspecified
pending
authorized
failed



```yaml
type: string
title: McpAuthorizationStatus
enum:
  - unspecified
  - pending
  - authorized
  - failed
description: |-
  McpAuthorizationStatus reports the current upstream authorization outcome.

  Values:
  unspecified
  pending
  authorized
  failed
```

</details>

<details id="schema-McpConfigurationMode">
<summary>McpConfigurationMode</summary>

McpConfigurationMode identifies desired hosting ownership and discovery behavior.

Values:
unspecified: unspecified is not a concrete report mode.
dynamic: dynamic uses a backend-minted gateway with discovery.
static: static uses a backend-minted gateway without discovery.
attached: attached retains a read-only external binding.



```yaml
type: string
title: McpConfigurationMode
enum:
  - unspecified
  - dynamic
  - static
  - attached
description: |-
  McpConfigurationMode identifies desired hosting ownership and discovery behavior.

  Values:
  unspecified: unspecified is not a concrete report mode.
  dynamic: dynamic uses a backend-minted gateway with discovery.
  static: static uses a backend-minted gateway without discovery.
  attached: attached retains a read-only external binding.
```

</details>

<details id="schema-McpConfigurationOutcome">
<summary>McpConfigurationOutcome</summary>

McpConfigurationOutcome describes one requested server in the reported generation.

Values:
unspecified: unspecified is not a concrete outcome.
pending: pending has not settled and cannot establish completeness.
applied: applied is confirmed membership, not upstream authorization.
skipped: skipped requires an explicitly applicable compatibility rule.
failed: failed retains intent that could not be provisioned.



```yaml
type: string
title: McpConfigurationOutcome
enum:
  - unspecified
  - pending
  - applied
  - skipped
  - failed
description: |-
  McpConfigurationOutcome describes one requested server in the reported generation.

  Values:
  unspecified: unspecified is not a concrete outcome.
  pending: pending has not settled and cannot establish completeness.
  applied: applied is confirmed membership, not upstream authorization.
  skipped: skipped requires an explicitly applicable compatibility rule.
  failed: failed retains intent that could not be provisioned.
```

</details>

<details id="schema-McpConfigurationReport">
<summary>McpConfigurationReport</summary>

McpConfigurationReport is also returned in Error.details on accepted provisioning failure.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `generation` | Yes | `string` | generation is an opaque durable fence, never ordered or synthesized by callers. |
| `mode` | Yes | [`McpConfigurationMode`](#schema-McpConfigurationMode) | mode separates minted discovery behavior from read-only attachment. |
| `complete` | No | `boolean` | complete means every active requested server is applied or explicitly permitted to skip. |
| `servers` | No | Array of [`McpConfiguredServer`](#schema-McpConfiguredServer) | servers is the complete active requested set; together with dormant_servers it matches McpGateway.servers. |
| `dormantServers` | No | Array of [`McpServer`](#schema-McpServer) | dormant_servers retains earlier minted intent while attached; it is empty in minted modes. |

```yaml
type: object
properties:
  generation:
    type: string
    description: generation is an opaque durable fence, never ordered or synthesized by callers.
  mode:
    not:
      enum:
        - unspecified
    description: mode separates minted discovery behavior from read-only attachment.
    $ref: "#/components/schemas/McpConfigurationMode"
  complete:
    type: boolean
    description: complete means every active requested server is applied or explicitly permitted to skip.
  servers:
    type: array
    items:
      $ref: "#/components/schemas/McpConfiguredServer"
    maxItems: 256
    description: servers is the complete active requested set; together with dormant_servers it matches McpGateway.servers.
  dormantServers:
    type: array
    items:
      $ref: "#/components/schemas/McpServer"
    maxItems: 256
    description: dormant_servers retains earlier minted intent while attached; it is empty in minted modes.
title: McpConfigurationReport
required:
  - generation
  - mode
description: McpConfigurationReport is also returned in Error.details on accepted provisioning failure.
```

</details>

<details id="schema-McpConfigurationReportErrorDetail">
<summary>McpConfigurationReportErrorDetail</summary>

McpConfigurationReport is also returned in Error.details on accepted provisioning failure.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `generation` | Yes | `string` | generation is an opaque durable fence, never ordered or synthesized by callers. |
| `mode` | Yes | [`McpConfigurationMode`](#schema-McpConfigurationMode) | mode separates minted discovery behavior from read-only attachment. |
| `complete` | No | `boolean` | complete means every active requested server is applied or explicitly permitted to skip. |
| `servers` | No | Array of [`McpConfiguredServer`](#schema-McpConfiguredServer) | servers is the complete active requested set; together with dormant_servers it matches McpGateway.servers. |
| `dormantServers` | No | Array of [`McpServer`](#schema-McpServer) | dormant_servers retains earlier minted intent while attached; it is empty in minted modes. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  generation:
    type: string
    description: generation is an opaque durable fence, never ordered or synthesized by callers.
  mode:
    not:
      enum:
        - unspecified
    description: mode separates minted discovery behavior from read-only attachment.
    $ref: "#/components/schemas/McpConfigurationMode"
  complete:
    type: boolean
    description: complete means every active requested server is applied or explicitly permitted to skip.
  servers:
    type: array
    items:
      $ref: "#/components/schemas/McpConfiguredServer"
    maxItems: 256
    description: servers is the complete active requested set; together with dormant_servers it matches McpGateway.servers.
  dormantServers:
    type: array
    items:
      $ref: "#/components/schemas/McpServer"
    maxItems: 256
    description: dormant_servers retains earlier minted intent while attached; it is empty in minted modes.
  "@type":
    type: string
    const: type.googleapis.com/docker.sandboxes.v1.McpConfigurationReport
title: McpConfigurationReportErrorDetail
required:
  - generation
  - mode
  - "@type"
description: McpConfigurationReport is also returned in Error.details on accepted provisioning failure.
```

</details>

<details id="schema-McpConfiguredServer">
<summary>McpConfiguredServer</summary>

McpConfiguredServer preserves one canonical requested identity and its provisioning outcome.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `server` | Yes | [`McpServer`](#schema-McpServer) | server is the resolved requested server; its type cannot silently change within retained intent. |
| `outcome` | Yes | [`McpConfigurationOutcome`](#schema-McpConfigurationOutcome) | outcome is concrete even when the gateway remains ready after a failed add. |
| `reason` | No | `string` | reason is nonsecret context and never carries credentials or credential-bearing URLs. |

```yaml
type: object
properties:
  server:
    description: server is the resolved requested server; its type cannot silently change within retained intent.
    $ref: "#/components/schemas/McpServer"
  outcome:
    not:
      enum:
        - unspecified
    description: outcome is concrete even when the gateway remains ready after a failed add.
    $ref: "#/components/schemas/McpConfigurationOutcome"
  reason:
    type: string
    description: reason is nonsecret context and never carries credentials or credential-bearing URLs.
title: McpConfiguredServer
required:
  - server
  - outcome
description: McpConfiguredServer preserves one canonical requested identity and its provisioning outcome.
```

</details>

<details id="schema-McpCreateSpec">
<summary>McpCreateSpec</summary>

McpCreateSpec is CreateSandboxRequest's MCP block: the StartMcpGateway
fields without the sandbox ref, wired before the workload starts.
servers must be empty when gateway_url is set

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `servers` | No | Array of `string` | servers are initial server names for a backend-minted gateway. |
| `static` | No | `boolean` | static pins the requested server set and closes gateway-side discovery. |
| `gatewayUrl` | No | `string or null` (uri) | gateway_url attaches read-only to a pre-existing shareable gateway when supported. |

```yaml
type: object
properties:
  servers:
    type: array
    items:
      type: string
      minLength: 1
    description: servers are initial server names for a backend-minted gateway.
  static:
    type: boolean
    description: static pins the requested server set and closes gateway-side discovery.
  gatewayUrl:
    type:
      - string
      - "null"
    format: uri
    description: gateway_url attaches read-only to a pre-existing shareable gateway when supported.
title: McpCreateSpec
additionalProperties: false
description: |-
  McpCreateSpec is CreateSandboxRequest's MCP block: the StartMcpGateway
  fields without the sandbox ref, wired before the workload starts.
  servers must be empty when gateway_url is set
```

</details>

<details id="schema-McpGateway">
<summary>McpGateway</summary>

McpGateway is gateway state for one sandbox.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `url` | No | `string` | url is the gateway endpoint. |
| `state` | No | [`McpGatewayState`](#schema-McpGatewayState) | state is the gateway readiness state. |
| `servers` | No | Array of [`McpServer`](#schema-McpServer) | servers is the accumulated requested set, not live gateway truth. |
| `reused` | No | `boolean` | reused is true when the response reused an existing gateway. |
| `uid` | No | `string or null` | (IMMUTABLE) uid is the existing immutable backing identity, omitted when none is available. |
| `name` | No | `string` | (IDENTIFIER) name is the complete singleton resource name. |

```yaml
type: object
properties:
  url:
    type: string
    description: url is the gateway endpoint.
  state:
    description: state is the gateway readiness state.
    $ref: "#/components/schemas/McpGatewayState"
  servers:
    type: array
    items:
      $ref: "#/components/schemas/McpServer"
    description: servers is the accumulated requested set, not live gateway truth.
  reused:
    type: boolean
    description: reused is true when the response reused an existing gateway.
  uid:
    type:
      - string
      - "null"
    description: (IMMUTABLE) uid is the existing immutable backing identity, omitted when none is available.
    readOnly: true
  name:
    type: string
    description: (IDENTIFIER) name is the complete singleton resource name.
title: McpGateway
description: McpGateway is gateway state for one sandbox.
```

</details>

<details id="schema-McpGatewayState">
<summary>McpGatewayState</summary>

McpGatewayState reports gateway readiness.

Values:
unspecified: unspecified is never a concrete known state.
provisioning: provisioning means the gateway is starting.
ready: ready means the gateway can accept server operations.
failed: failed means gateway setup failed.



```yaml
type: string
title: McpGatewayState
enum:
  - unspecified
  - provisioning
  - ready
  - failed
description: |-
  McpGatewayState reports gateway readiness.

  Values:
  unspecified: unspecified is never a concrete known state.
  provisioning: provisioning means the gateway is starting.
  ready: ready means the gateway can accept server operations.
  failed: failed means gateway setup failed.
```

</details>

<details id="schema-McpServer">
<summary>McpServer</summary>

McpServer projects one external server key configured on a gateway.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `name` | No | `string` | name is the external configured server key, not an API resource name. |
| `type` | No | [`McpServerType`](#schema-McpServerType) | type is the configured server type. |

```yaml
type: object
properties:
  name:
    type: string
    description: name is the external configured server key, not an API resource name.
  type:
    description: type is the configured server type.
    $ref: "#/components/schemas/McpServerType"
title: McpServer
description: McpServer projects one external server key configured on a gateway.
```

</details>

<details id="schema-McpServerType">
<summary>McpServerType</summary>

McpServerType identifies a gateway server integration style.

Values:
unspecified: unspecified is not a concrete server type.
remote: remote is an HTTP remote MCP server.
container: container is a container-backed MCP server.
stdio: stdio is a stdio-backed MCP server.



```yaml
type: string
title: McpServerType
enum:
  - unspecified
  - remote
  - container
  - stdio
description: |-
  McpServerType identifies a gateway server integration style.

  Values:
  unspecified: unspecified is not a concrete server type.
  remote: remote is an HTTP remote MCP server.
  container: container is a container-backed MCP server.
  stdio: stdio is a stdio-backed MCP server.
```

</details>

<details id="schema-NetworkPolicy">
<summary>NetworkPolicy</summary>

NetworkPolicy is an inline policy request materialized in the governance plane.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `mode` | Yes | `string` | mode is required and cannot be UNSPECIFIED. |
| `allowNetworks` | No | Array of `string` | allow_networks are exact network names or CIDRs allowed by this document. |
| `denyNetworks` | No | Array of `string` | deny_networks are exact network names or CIDRs denied by this document. |

```yaml
type: object
properties:
  mode:
    type: string
    title: NetworkPolicyMode
    enum:
      - denyAll
    description: mode is required and cannot be UNSPECIFIED.
    not:
      enum:
        - unspecified
  allowNetworks:
    type: array
    items:
      type: string
    description: allow_networks are exact network names or CIDRs allowed by this document.
  denyNetworks:
    type: array
    items:
      type: string
    description: deny_networks are exact network names or CIDRs denied by this document.
title: NetworkPolicy
required:
  - mode
additionalProperties: false
description: NetworkPolicy is an inline policy request materialized in the governance plane.
```

</details>

<details id="schema-NetworkPolicyMode">
<summary>NetworkPolicyMode</summary>

NetworkPolicyMode is the top-level egress disposition.

Values:
unspecified: unspecified is invalid in policy documents and in conforming responses.
allowAll: allowAll allows egress unless denied by a narrower rule.
denyAll: denyAll denies egress unless allowed by a narrower rule.



```yaml
type: string
title: NetworkPolicyMode
enum:
  - unspecified
  - allowAll
  - denyAll
description: |-
  NetworkPolicyMode is the top-level egress disposition.

  Values:
  unspecified: unspecified is invalid in policy documents and in conforming responses.
  allowAll: allowAll allows egress unless denied by a narrower rule.
  denyAll: denyAll denies egress unless allowed by a narrower rule.
```

</details>

<details id="schema-OAuthConfig">
<summary>OAuthConfig</summary>

OAuthConfig is non-secret OAuth metadata returned on reads.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `tokenEndpoint` | No | `string` | token_endpoint is the OAuth token endpoint. |
| `clientId` | No | `string` | client_id is the OAuth client id. |

```yaml
type: object
properties:
  tokenEndpoint:
    type: string
    description: token_endpoint is the OAuth token endpoint.
  clientId:
    type: string
    description: client_id is the OAuth client id.
title: OAuthConfig
description: OAuthConfig is non-secret OAuth metadata returned on reads.
```

</details>

<details id="schema-OAuthRefreshMaterial">
<summary>OAuthRefreshMaterial</summary>

OAuthRefreshMaterial is OAuth secret material. Unlike OAuthRefreshMaterialV2 it pairs
no bootstrap token with an expiry and stores client_id as the caller supplied it.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `refreshToken` | Yes | `string` | refresh_token is the OAuth refresh token. |
| `initialAccessToken` | No | `string` | initial_access_token is an optional bootstrap access token. |
| `expiresAt` | No | [`Timestamp`](#schema-Timestamp) | expires_at is the access token expiry when supplied. |
| `scopes` | No | Array of `string` | scopes are requested OAuth scopes. |
| `providerAccountId` | No | `string` | provider_account_id identifies the provider account. |
| `idToken` | No | `string` | id_token is an optional OpenID Connect id token. |
| `clientId` | No | `string` | client_id is the OAuth client id. |

```yaml
type: object
properties:
  refreshToken:
    type: string
    minLength: 1
    description: refresh_token is the OAuth refresh token.
    writeOnly: true
  initialAccessToken:
    type: string
    description: initial_access_token is an optional bootstrap access token.
    writeOnly: true
  expiresAt:
    description: expires_at is the access token expiry when supplied.
    $ref: "#/components/schemas/Timestamp"
  scopes:
    type: array
    items:
      type: string
    description: scopes are requested OAuth scopes.
  providerAccountId:
    type: string
    description: provider_account_id identifies the provider account.
  idToken:
    type: string
    description: id_token is an optional OpenID Connect id token.
    writeOnly: true
  clientId:
    type: string
    description: client_id is the OAuth client id.
title: OAuthRefreshMaterial
required:
  - refreshToken
additionalProperties: false
description: |-
  OAuthRefreshMaterial is OAuth secret material. Unlike OAuthRefreshMaterialV2 it pairs
  no bootstrap token with an expiry and stores client_id as the caller supplied it.
```

</details>

<details id="schema-PolicyDecision">
<summary>PolicyDecision</summary>

PolicyDecision is the outcome of an evaluated access attempt.

Values:
unspecified: unspecified means the decision is unknown or not reported.
allowed: allowed means the attempt was allowed.
blocked: blocked means the attempt was blocked.



```yaml
type: string
title: PolicyDecision
enum:
  - unspecified
  - allowed
  - blocked
description: |-
  PolicyDecision is the outcome of an evaluated access attempt.

  Values:
  unspecified: unspecified means the decision is unknown or not reported.
  allowed: allowed means the attempt was allowed.
  blocked: blocked means the attempt was blocked.
```

</details>

<details id="schema-PolicyDomain">
<summary>PolicyDomain</summary>

PolicyDomain identifies a policy subsystem.

Values:
unspecified: unspecified is not a concrete policy domain.
network: network is network egress policy.



```yaml
type: string
title: PolicyDomain
enum:
  - unspecified
  - network
description: |-
  PolicyDomain identifies a policy subsystem.

  Values:
  unspecified: unspecified is not a concrete policy domain.
  network: network is network egress policy.
```

</details>

<details id="schema-PolicyLayer">
<summary>PolicyLayer</summary>

PolicyLayer identifies where an effective rule came from. ORG and OWNER are
the two governance scopes; exactly one governs a caller, so a conforming effective view never mixes them.

Values:
unspecified: unspecified means the layer is unknown or not reported.
org: org is org-level governance policy.
owner: owner is owner-level governance policy.
kit: kit is the sandbox-scoped policy a kit contributes.
attached: attached is a governance policy bound to the sandbox: an inline document or a referenced policy id, materialized and bound by the governance plane at create.
gateway: gateway is an enforced runtime-owned gateway policy contribution, not a backend service-route exemption.
default: default is an explicit installed default-policy contribution, not a substitute for unknown provenance or an owner-level policy.



```yaml
type: string
title: PolicyLayer
enum:
  - unspecified
  - org
  - owner
  - kit
  - attached
  - gateway
  - default
description: |-
  PolicyLayer identifies where an effective rule came from. ORG and OWNER are
  the two governance scopes; exactly one governs a caller, so a conforming effective view never mixes them.

  Values:
  unspecified: unspecified means the layer is unknown or not reported.
  org: org is org-level governance policy.
  owner: owner is owner-level governance policy.
  kit: kit is the sandbox-scoped policy a kit contributes.
  attached: attached is a governance policy bound to the sandbox: an inline document or a referenced policy id, materialized and bound by the governance plane at create.
  gateway: gateway is an enforced runtime-owned gateway policy contribution, not a backend service-route exemption.
  default: default is an explicit installed default-policy contribution, not a substitute for unknown provenance or an owner-level policy.
```

</details>

<details id="schema-PolicyLogEntry">
<summary>PolicyLogEntry</summary>

PolicyLogEntry is an aggregated observed policy decision.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `resource` | No | `string` | resource is the requested destination or resource. |
| `decision` | No | [`PolicyDecision`](#schema-PolicyDecision) | decision is the observed allow or block outcome. |
| `proxyType` | No | `string` | proxy_type is the proxy or enforcement path that observed the decision. |
| `rule` | No | `string` | rule is the matching rule when available. |
| `reason` | No | `string` | reason is backend-readable context. |
| `count` | No | `integer or string` (int64) | count is the number of coalesced observations. |
| `firstSeen` | No | [`Timestamp`](#schema-Timestamp) | first_seen is the first observation timestamp. |
| `lastSeen` | No | [`Timestamp`](#schema-Timestamp) | last_seen is the most recent observation timestamp. |
| `sandbox` | No | `string` | sandbox is the resource name recorded for this reference. |
| `sandboxUid` | No | `string or null` | sandbox_uid retains the observed backing incarnation when the backend provides one. |
| `domain` | No | [`PolicyDomain`](#schema-PolicyDomain) | domain is the policy domain that emitted the entry. |

```yaml
type: object
properties:
  resource:
    type: string
    description: resource is the requested destination or resource.
  decision:
    description: decision is the observed allow or block outcome.
    $ref: "#/components/schemas/PolicyDecision"
  proxyType:
    type: string
    description: proxy_type is the proxy or enforcement path that observed the decision.
  rule:
    type: string
    description: rule is the matching rule when available.
  reason:
    type: string
    description: reason is backend-readable context.
  count:
    type:
      - integer
      - string
    format: int64
    description: count is the number of coalesced observations.
  firstSeen:
    description: first_seen is the first observation timestamp.
    $ref: "#/components/schemas/Timestamp"
  lastSeen:
    description: last_seen is the most recent observation timestamp.
    $ref: "#/components/schemas/Timestamp"
  sandbox:
    type: string
    description: sandbox is the resource name recorded for this reference.
  sandboxUid:
    type:
      - string
      - "null"
    description: sandbox_uid retains the observed backing incarnation when the backend provides one.
    readOnly: true
  domain:
    description: domain is the policy domain that emitted the entry.
    $ref: "#/components/schemas/PolicyDomain"
title: PolicyLogEntry
description: PolicyLogEntry is an aggregated observed policy decision.
```

</details>

<details id="schema-Port">
<summary>Port</summary>

Port is one sandbox port made reachable from outside the sandbox. The same
facts appear inline on SandboxCore.ports for a caller reading the sandbox.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `number` | Yes | `integer` (int32) | (IMMUTABLE) number is immutable and unique within the sandbox; duplicates refuse ALREADY_EXISTS. |
| `protocol` | No | [`Protocol`](#schema-Protocol) | (IMMUTABLE) protocol defaults to TCP when unspecified. |
| `url` | No | `string` | url is the required externally reachable address, including host-local addresses. |
| `etag` | No | `string` | etag identifies the observed version of this port. |
| `uid` | No | `string or null` | (IMMUTABLE) uid is the durable publication incarnation when the backend provides one. |
| `name` | No | `string` | (IDENTIFIER) name identifies the published port number under its sandbox resource name. |

```yaml
type: object
properties:
  number:
    type: integer
    maximum: 65535
    minimum: 1
    format: int32
    description: (IMMUTABLE) number is immutable and unique within the sandbox; duplicates refuse ALREADY_EXISTS.
  protocol:
    description: (IMMUTABLE) protocol defaults to TCP when unspecified.
    $ref: "#/components/schemas/Protocol"
  url:
    type: string
    description: url is the required externally reachable address, including host-local addresses.
    readOnly: true
  etag:
    type: string
    description: etag identifies the observed version of this port.
    readOnly: true
  uid:
    type:
      - string
      - "null"
    description: (IMMUTABLE) uid is the durable publication incarnation when the backend provides one.
    readOnly: true
  name:
    type: string
    description: (IDENTIFIER) name identifies the published port number under its sandbox resource name.
title: Port
required:
  - number
description: |-
  Port is one sandbox port made reachable from outside the sandbox. The same
  facts appear inline on SandboxCore.ports for a caller reading the sandbox.
```

</details>

<details id="schema-RawImageStartup">
<summary>RawImageStartup</summary>

RawImageStartup resolves identity and working directory from the pinned image.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `start` | No | [`StartupArgv`](#schema-StartupArgv) | start replaces the complete launch argv; omission uses the image command. |

```yaml
type: object
properties:
  start:
    description: start replaces the complete launch argv; omission uses the image command.
    $ref: "#/components/schemas/StartupArgv"
title: RawImageStartup
additionalProperties: false
description: RawImageStartup resolves identity and working directory from the pinned image.
```

</details>

<details id="schema-Sandbox">
<summary>Sandbox</summary>

Sandbox carries shared core state and the effective state of its selected features.
uid must equal the immutable resource ID in name

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `core` | Yes | [`SandboxCore`](#schema-SandboxCore) | core is the backend-neutral sandbox view. |
| `effectiveFeatures` | No | [`EffectiveFeatures`](#schema-EffectiveFeatures) | effective_features is emitted only by the composition contract projection. |
| `failure` | No | [`Error`](#schema-Error) | failure is present only while status is FAILED and clears on recovery. |
| `name` | Yes | `string` | (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name. |
| `uid` | Yes | `string` | (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment. |
| `displayName` | No | `string` | display_name is the scoped label; UpdateSandbox changes it without changing identity. |

```yaml
type: object
properties:
  core:
    description: core is the backend-neutral sandbox view.
    $ref: "#/components/schemas/SandboxCore"
  effectiveFeatures:
    description: effective_features is emitted only by the composition contract projection.
    $ref: "#/components/schemas/EffectiveFeatures"
  failure:
    description: failure is present only while status is FAILED and clears on recovery.
    readOnly: true
    $ref: "#/components/schemas/Error"
  name:
    type: string
    pattern: ^sandboxes/[^/]+$
    description: (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name.
  uid:
    type:
      - string
    minLength: 1
    description: (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment.
    readOnly: true
  displayName:
    type: string
    description: display_name is the scoped label; UpdateSandbox changes it without changing identity.
title: Sandbox
required:
  - core
  - name
  - uid
description: |-
  Sandbox carries shared core state and the effective state of its selected features.
  uid must equal the immutable resource ID in name
```

</details>

<details id="schema-Secret">
<summary>Secret</summary>

Secret is returned metadata and never includes secret material.
uid must equal the immutable resource ID in name

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `uid` | Yes | `string` | (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment. |
| `name` | Yes | `string` | (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name. |
| `serviceType` | No | `string` | service_type identifies the consuming service. |
| `type` | No | [`SecretType`](#schema-SecretType) | type identifies the stored material shape. |
| `oauthConfig` | No | [`OAuthConfig`](#schema-OAuthConfig) | oauth_config is returned only for OAuth-shaped secrets. |
| `scopes` | No | Array of `string` | scopes are non-secret OAuth scopes. |
| `createdAt` | Yes | [`Timestamp`](#schema-Timestamp) | created_at is the creation timestamp. |
| `updatedAt` | Yes | [`Timestamp`](#schema-Timestamp) | updated_at is the last update timestamp. |
| `injection` | No | [`CustomInjection`](#schema-CustomInjection) | injection is returned only for custom-shaped secrets. |
| `placeholder` | No | `string` | placeholder is server-generated per custom secret and stable across value rotation; it is not secret material. A client presents it where the credential would go; for the configured hosts the backend replaces the configured header's value at egress regardless of what was sent. |
| `etag` | Yes | `string` | etag identifies the observed version of this resource. It is opaque and strong, changes on every visible change, and is what a mutation sends as its precondition. |
| `displayName` | No | `string` | display_name is a scoped label and never selects the resource. |

```yaml
type: object
properties:
  uid:
    type:
      - string
    minLength: 1
    description: (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment.
    readOnly: true
  name:
    type: string
    pattern: ^secrets/[^/]+$
    description: (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name.
  serviceType:
    type: string
    description: service_type identifies the consuming service.
  type:
    description: type identifies the stored material shape.
    $ref: "#/components/schemas/SecretType"
  oauthConfig:
    description: oauth_config is returned only for OAuth-shaped secrets.
    $ref: "#/components/schemas/OAuthConfig"
  scopes:
    type: array
    items:
      type: string
    description: scopes are non-secret OAuth scopes.
  createdAt:
    description: created_at is the creation timestamp.
    $ref: "#/components/schemas/Timestamp"
  updatedAt:
    description: updated_at is the last update timestamp.
    $ref: "#/components/schemas/Timestamp"
  injection:
    description: injection is returned only for custom-shaped secrets.
    $ref: "#/components/schemas/CustomInjection"
  placeholder:
    type: string
    description: |-
      placeholder is server-generated per custom secret and stable across value
      rotation; it is not secret material. A client presents it where the
      credential would go; for the configured hosts the backend replaces the
      configured header's value at egress regardless of what was sent.
  etag:
    type: string
    minLength: 1
    pattern: ^"[^"\x00-\x20\x7f]*"$
    description: |-
      etag identifies the observed version of this resource. It is opaque and strong,
      changes on every visible change, and is what a mutation sends as its precondition.
    readOnly: true
  displayName:
    type: string
    description: display_name is a scoped label and never selects the resource.
title: Secret
required:
  - uid
  - name
  - createdAt
  - updatedAt
  - etag
description: |-
  Secret is returned metadata and never includes secret material.
  uid must equal the immutable resource ID in name
```

</details>

<details id="schema-SecretType">
<summary>SecretType</summary>

The stored credential shape: token, OAuth refresh material, custom injection, registry credentials, or a credential resolved from an approved host-owned source. Secret material is never returned.



```yaml
type: string
title: SecretType
enum:
  - unspecified
  - token
  - oauthRefresh
  - custom
  - registryV2
  - hostResolved
description: "The stored credential shape: token, OAuth refresh material, custom injection, registry credentials, or a credential resolved from an approved host-owned source. Secret material is never returned."
```

</details>

<details id="schema-Snapshot">
<summary>Snapshot</summary>

Snapshot is the detailed snapshot view.
uid must equal the immutable resource ID in name

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `uid` | Yes | `string` | (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment. |
| `name` | Yes | `string` | (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name. |
| `description` | No | `string` | description is caller-supplied text. |
| `sandbox` | No | `string` | sandbox is the resource name recorded for this reference. |
| `sandboxUid` | No | `string or null` | sandbox_uid retains the observed backing incarnation when the backend provides one. |
| `status` | Yes | [`SnapshotStatus`](#schema-SnapshotStatus) | status is the snapshot readiness state. |
| `totalSizeBytes` | No | `integer or string` (int64) | total_size_bytes is the approximate stored snapshot size. |
| `resources` | No | [`Resources`](#schema-Resources) | resources describes the captured sandbox resources. |
| `platform` | No | [`Platform`](#schema-Platform) | platform is the captured sandbox platform. |
| `createdAt` | Yes | [`Timestamp`](#schema-Timestamp) | created_at is the creation timestamp. |
| `failure` | No | [`Error`](#schema-Error) | failure is set when status is FAILED. |
| `captureMode` | No | [`CaptureMode`](#schema-CaptureMode) | capture_mode is the mode the snapshot was captured with. |
| `etag` | Yes | `string` | etag identifies the observed version of this resource. It is opaque and strong, changes on every visible change, and is what a mutation sends as its precondition. |
| `displayName` | No | `string` | display_name is a scoped label and never selects the resource. |

```yaml
type: object
properties:
  uid:
    type:
      - string
    minLength: 1
    description: (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment.
    readOnly: true
  name:
    type: string
    pattern: ^snapshots/[^/]+$
    description: (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name.
  description:
    type: string
    description: description is caller-supplied text.
  sandbox:
    type: string
    description: sandbox is the resource name recorded for this reference.
  sandboxUid:
    type:
      - string
      - "null"
    description: sandbox_uid retains the observed backing incarnation when the backend provides one.
    readOnly: true
  status:
    not:
      enum:
        - unspecified
    description: status is the snapshot readiness state.
    readOnly: true
    $ref: "#/components/schemas/SnapshotStatus"
  totalSizeBytes:
    type:
      - integer
      - string
    format: int64
    description: total_size_bytes is the approximate stored snapshot size.
  resources:
    description: resources describes the captured sandbox resources.
    $ref: "#/components/schemas/Resources"
  platform:
    description: platform is the captured sandbox platform.
    $ref: "#/components/schemas/Platform"
  createdAt:
    description: created_at is the creation timestamp.
    $ref: "#/components/schemas/Timestamp"
  failure:
    description: failure is set when status is FAILED.
    $ref: "#/components/schemas/Error"
  captureMode:
    description: capture_mode is the mode the snapshot was captured with.
    $ref: "#/components/schemas/CaptureMode"
  etag:
    type: string
    minLength: 1
    pattern: ^"[^"\x00-\x20\x7f]*"$
    description: |-
      etag identifies the observed version of this resource. It is opaque and strong,
      changes on every visible change, and is what a mutation sends as its precondition.
    readOnly: true
  displayName:
    type: string
    description: display_name is a scoped label and never selects the resource.
title: Snapshot
required:
  - uid
  - name
  - status
  - createdAt
  - etag
description: |-
  Snapshot is the detailed snapshot view.
  uid must equal the immutable resource ID in name
```

</details>

<details id="schema-SnapshotStatus">
<summary>SnapshotStatus</summary>

SnapshotStatus is snapshot readiness state.

Values:
unspecified: unspecified is never a concrete known state.
creating: creating means capture has not completed.
ready: ready means the snapshot can be restored.
failed: failed means capture reached a terminal failure.
deleted: deleted means the backend exposes a retained delete tombstone.



```yaml
type: string
title: SnapshotStatus
enum:
  - unspecified
  - creating
  - ready
  - failed
  - deleted
description: |-
  SnapshotStatus is snapshot readiness state.

  Values:
  unspecified: unspecified is never a concrete known state.
  creating: creating means capture has not completed.
  ready: ready means the snapshot can be restored.
  failed: failed means capture reached a terminal failure.
  deleted: deleted means the backend exposes a retained delete tombstone.
```

</details>

<details id="schema-SnapshotSummary">
<summary>SnapshotSummary</summary>

SnapshotSummary is the list view of a snapshot.
uid must equal the immutable resource ID in name

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `uid` | Yes | `string` | (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment. |
| `name` | Yes | `string` | name ends in the immutable backing ID and never changes with display_name. |
| `sandbox` | No | `string` | sandbox is the resource name recorded for this reference. |
| `sandboxUid` | No | `string or null` | sandbox_uid retains the observed backing incarnation when the backend provides one. |
| `status` | Yes | [`SnapshotStatus`](#schema-SnapshotStatus) | status is the snapshot readiness state. |
| `totalSizeBytes` | No | `integer or string` (int64) | total_size_bytes is the approximate stored snapshot size. |
| `createdAt` | Yes | [`Timestamp`](#schema-Timestamp) | created_at is the creation timestamp. |
| `captureMode` | No | [`CaptureMode`](#schema-CaptureMode) | capture_mode is the mode the snapshot was captured with. |
| `displayName` | No | `string` | display_name is a scoped label and never selects the resource. |

```yaml
type: object
properties:
  uid:
    type:
      - string
    minLength: 1
    description: (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment.
    readOnly: true
  name:
    type: string
    pattern: ^snapshots/[^/]+$
    description: name ends in the immutable backing ID and never changes with display_name.
  sandbox:
    type: string
    description: sandbox is the resource name recorded for this reference.
  sandboxUid:
    type:
      - string
      - "null"
    description: sandbox_uid retains the observed backing incarnation when the backend provides one.
    readOnly: true
  status:
    not:
      enum:
        - unspecified
    description: status is the snapshot readiness state.
    $ref: "#/components/schemas/SnapshotStatus"
  totalSizeBytes:
    type:
      - integer
      - string
    format: int64
    description: total_size_bytes is the approximate stored snapshot size.
  createdAt:
    description: created_at is the creation timestamp.
    $ref: "#/components/schemas/Timestamp"
  captureMode:
    description: capture_mode is the mode the snapshot was captured with.
    $ref: "#/components/schemas/CaptureMode"
  displayName:
    type: string
    description: display_name is a scoped label and never selects the resource.
title: SnapshotSummary
required:
  - uid
  - name
  - status
  - createdAt
description: |-
  SnapshotSummary is the list view of a snapshot.
  uid must equal the immutable resource ID in name
```

</details>

<details id="schema-StartupArgv">
<summary>StartupArgv</summary>

StartupArgv preserves literal arguments, including empty positional arguments.
startup argv requires a nonempty executable

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `argv` | Yes | Array of `string` | argv is executed without shell expansion. It carries at least one element, and that first element is the executable and must not be empty. |

```yaml
type: object
properties:
  argv:
    type: array
    items:
      type: string
    minItems: 1
    description: |-
      argv is executed without shell expansion. It carries at least one element,
      and that first element is the executable and must not be empty.
title: StartupArgv
required:
  - argv
additionalProperties: false
description: |-
  StartupArgv preserves literal arguments, including empty positional arguments.
  startup argv requires a nonempty executable
```

</details>

<details id="schema-StartupExecution">
<summary>StartupExecution</summary>

StartupExecution retains the accepted execution promise across fresh boots.
It is ordinary create input.



```yaml
type: object
oneOf:
  - type: object
    properties:
      rawImage:
        description: raw_image executes image commands using OCI identity and directory defaults.
        $ref: "#/components/schemas/RawImageStartup"
    title: raw_image
    required:
      - rawImage
title: StartupExecution
unevaluatedProperties: false
description: |-
  StartupExecution retains the accepted execution promise across fresh boots.
  It is ordinary create input.
```

</details>

<details id="schema-StopMcpGatewayResponse">
<summary>StopMcpGatewayResponse</summary>

StopMcpGatewayResponse may retain the fenced desired configuration without claiming a live gateway.



```yaml
type: object
title: StopMcpGatewayResponse
description: StopMcpGatewayResponse may retain the fenced desired configuration without claiming a live gateway.
```

</details>

<details id="schema-TokenSecretMaterial">
<summary>TokenSecretMaterial</summary>

TokenSecretMaterial is opaque token material.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `value` | Yes | `string` | value is the secret token. |

```yaml
type: object
properties:
  value:
    type: string
    minLength: 1
    description: value is the secret token.
    writeOnly: true
title: TokenSecretMaterial
required:
  - value
additionalProperties: false
description: TokenSecretMaterial is opaque token material.
```

</details>

<details id="schema-Volume">
<summary>Volume</summary>

Volume is persistent storage metadata within its owner and contract scope.
uid must equal the immutable resource ID in name

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `uid` | Yes | `string` | (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment. |
| `name` | Yes | `string` | (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name. |
| `createdAt` | Yes | [`Timestamp`](#schema-Timestamp) | created_at is the creation timestamp. |
| `etag` | Yes | `string` | etag identifies the observed version of this resource. It is opaque and strong, changes on every visible change, and is what a mutation sends as its precondition. |
| `displayName` | No | `string` | display_name is a scoped label and never selects the resource. |

```yaml
type: object
properties:
  uid:
    type:
      - string
    minLength: 1
    description: (IMMUTABLE) uid is the existing backing ID, equal to the resource name's terminal segment.
    readOnly: true
  name:
    type: string
    pattern: ^volumes/[^/]+$
    description: (IDENTIFIER) name ends in the immutable backing ID and never changes with display_name.
  createdAt:
    description: created_at is the creation timestamp.
    $ref: "#/components/schemas/Timestamp"
  etag:
    type: string
    minLength: 1
    pattern: ^"[^"\x00-\x20\x7f]*"$
    description: |-
      etag identifies the observed version of this resource. It is opaque and strong,
      changes on every visible change, and is what a mutation sends as its precondition.
    readOnly: true
  displayName:
    type: string
    description: display_name is a scoped label and never selects the resource.
title: Volume
required:
  - uid
  - name
  - createdAt
  - etag
description: |-
  Volume is persistent storage metadata within its owner and contract scope.
  uid must equal the immutable resource ID in name
```

</details>

<details id="schema-Duration">
<summary>Duration</summary>

Seconds with a trailing s, for example 30s or 1.5s.



```yaml
type: string
pattern: ^-?\d+(\.\d{1,9})?s$
description: Seconds with a trailing s, for example 30s or 1.5s.
examples:
  - 30s
  - 1.5s
```

</details>

<details id="schema-ValidationFieldType">
<summary>ValidationFieldType</summary>





```yaml
type: string
title: Type
enum:
  - TYPE_DOUBLE
  - TYPE_FLOAT
  - TYPE_INT64
  - TYPE_UINT64
  - TYPE_INT32
  - TYPE_FIXED64
  - TYPE_FIXED32
  - TYPE_BOOL
  - TYPE_STRING
  - TYPE_GROUP
  - TYPE_MESSAGE
  - TYPE_BYTES
  - TYPE_UINT32
  - TYPE_ENUM
  - TYPE_SFIXED32
  - TYPE_SFIXED64
  - TYPE_SINT32
  - TYPE_SINT64
```

</details>

<details id="schema-Timestamp">
<summary>Timestamp</summary>

An RFC 3339 timestamp in UTC, for example 2026-01-02T03:04:05Z.



```yaml
type: string
examples:
  - 2023-01-15T01:30:15.01Z
  - 2024-12-25T12:00:00Z
format: date-time
description: An RFC 3339 timestamp in UTC, for example 2026-01-02T03:04:05Z.
```

</details>

<details id="schema-BadRequest">
<summary>BadRequest</summary>

Describes violations in a client request. This error type focuses on the
syntactic aspects of the request.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `fieldViolations` | No | Array of [`BadRequestFieldViolation`](#schema-BadRequestFieldViolation) | Describes all violations in a client request. |

```yaml
type: object
properties:
  fieldViolations:
    type: array
    items:
      $ref: "#/components/schemas/BadRequestFieldViolation"
    description: Describes all violations in a client request.
title: BadRequest
description: |-
  Describes violations in a client request. This error type focuses on the
  syntactic aspects of the request.
```

</details>

<details id="schema-BadRequestFieldViolation">
<summary>BadRequestFieldViolation</summary>

A message type used to describe a single bad request field.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `field` | No | `string` | A path that leads to a field in the request body. The value will be a sequence of dot-separated identifiers that identify a protocol buffer field. Consider the following: message CreateContactRequest { message EmailAddress { enum Type { TYPE_UNSPECIFIED = 0; HOME = 1; WORK = 2; } optional string email = 1; repeated EmailType type = 2; } string full_name = 1; repeated EmailAddress email_addresses = 2; } In this example, in proto `field` could take one of the following values: * `full_name` for a violation in the `full_name` value * `email_addresses[0].email` for a violation in the `email` field of the first `email_addresses` message * `email_addresses[2].type[1]` for a violation in the second `type` value in the third `email_addresses` message. In JSON, the same values are represented as: * `fullName` for a violation in the `fullName` value * `emailAddresses[0].email` for a violation in the `email` field of the first `emailAddresses` message * `emailAddresses[2].type[1]` for a violation in the second `type` value in the third `emailAddresses` message. |
| `description` | No | `string` | A description of why the request element is bad. |
| `reason` | No | `string` | The reason of the field-level error. This is a constant value that identifies the proximate cause of the field-level error. It should uniquely identify the type of the FieldViolation within the scope of the ErrorInfo.domain. This should be at most 63 characters and match a regular expression of `[A-Z][A-Z0-9_]+[A-Z0-9]`, which represents UPPER_SNAKE_CASE. |
| `localizedMessage` | No | [`LocalizedMessage`](#schema-LocalizedMessage) | Provides a localized error message for field-level errors that is safe to return to the API consumer. |

```yaml
type: object
properties:
  field:
    type: string
    description: |-
      A path that leads to a field in the request body. The value will be a
      sequence of dot-separated identifiers that identify a protocol buffer
      field.

      Consider the following:

      message CreateContactRequest {
      message EmailAddress {
      enum Type {
      TYPE_UNSPECIFIED = 0;
      HOME = 1;
      WORK = 2;
      }

      optional string email = 1;
      repeated EmailType type = 2;
      }

      string full_name = 1;
      repeated EmailAddress email_addresses = 2;
      }

      In this example, in proto `field` could take one of the following values:

      * `full_name` for a violation in the `full_name` value
      * `email_addresses[0].email` for a violation in the `email` field of the
      first `email_addresses` message
      * `email_addresses[2].type[1]` for a violation in the second `type`
      value in the third `email_addresses` message.

      In JSON, the same values are represented as:

      * `fullName` for a violation in the `fullName` value
      * `emailAddresses[0].email` for a violation in the `email` field of the
      first `emailAddresses` message
      * `emailAddresses[2].type[1]` for a violation in the second `type`
      value in the third `emailAddresses` message.
  description:
    type: string
    description: A description of why the request element is bad.
  reason:
    type: string
    description: |-
      The reason of the field-level error. This is a constant value that
      identifies the proximate cause of the field-level error. It should
      uniquely identify the type of the FieldViolation within the scope of the
      ErrorInfo.domain. This should be at most 63
      characters and match a regular expression of `[A-Z][A-Z0-9_]+[A-Z0-9]`,
      which represents UPPER_SNAKE_CASE.
  localizedMessage:
    description: |-
      Provides a localized error message for field-level errors that is safe to
      return to the API consumer.
    $ref: "#/components/schemas/LocalizedMessage"
title: FieldViolation
description: A message type used to describe a single bad request field.
```

</details>

<details id="schema-BadRequestErrorDetail">
<summary>BadRequestErrorDetail</summary>

Describes violations in a client request. This error type focuses on the
syntactic aspects of the request.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `fieldViolations` | No | Array of [`BadRequestFieldViolation`](#schema-BadRequestFieldViolation) | Describes all violations in a client request. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  fieldViolations:
    type: array
    items:
      $ref: "#/components/schemas/BadRequestFieldViolation"
    description: Describes all violations in a client request.
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.BadRequest
title: BadRequestErrorDetail
description: |-
  Describes violations in a client request. This error type focuses on the
  syntactic aspects of the request.
required:
  - "@type"
```

</details>

<details id="schema-DebugInfo">
<summary>DebugInfo</summary>

Describes additional debugging info.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `stackEntries` | No | Array of `string` | The stack trace entries indicating where the error occurred. |
| `detail` | No | `string` | Additional debugging information provided by the server. |

```yaml
type: object
properties:
  stackEntries:
    type: array
    items:
      type: string
    description: The stack trace entries indicating where the error occurred.
  detail:
    type: string
    description: Additional debugging information provided by the server.
title: DebugInfo
description: Describes additional debugging info.
```

</details>

<details id="schema-DebugInfoErrorDetail">
<summary>DebugInfoErrorDetail</summary>

Describes additional debugging info.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `stackEntries` | No | Array of `string` | The stack trace entries indicating where the error occurred. |
| `detail` | No | `string` | Additional debugging information provided by the server. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  stackEntries:
    type: array
    items:
      type: string
    description: The stack trace entries indicating where the error occurred.
  detail:
    type: string
    description: Additional debugging information provided by the server.
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.DebugInfo
title: DebugInfoErrorDetail
description: Describes additional debugging info.
required:
  - "@type"
```

</details>

<details id="schema-ErrorInfo">
<summary>ErrorInfo</summary>

Describes the cause of the error with structured details.

Example of an error when contacting the "pubsub.googleapis.com" API when it
is not enabled:

{ "reason": "API_DISABLED"
"domain": "googleapis.com"
"metadata": {
"resource": "projects/123",
"service": "pubsub.googleapis.com"
}
}

This response indicates that the pubsub.googleapis.com API is not enabled.

Example of an error that is returned when attempting to create a Spanner
instance in a region that is out of stock:

{ "reason": "STOCKOUT"
"domain": "spanner.googleapis.com",
"metadata": {
"availableRegions": "us-central1,us-east2"
}
}

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `reason` | No | `string` | The reason of the error. This is a constant value that identifies the proximate cause of the error. Error reasons are unique within a particular domain of errors. This should be at most 63 characters and match a regular expression of `[A-Z][A-Z0-9_]+[A-Z0-9]`, which represents UPPER_SNAKE_CASE. |
| `domain` | No | `string` | The logical grouping to which the "reason" belongs. The error domain is typically the registered service name of the tool or product that generates the error. Example: "pubsub.googleapis.com". If the error is generated by some common infrastructure, the error domain must be a globally unique value that identifies the infrastructure. For Google API infrastructure, the error domain is "googleapis.com". |
| `metadata` | No | `object` | Additional structured details about this error. Keys must match a regular expression of `[a-z][a-zA-Z0-9-_]+` but should ideally be lowerCamelCase. Also, they must be limited to 64 characters in length. When identifying the current value of an exceeded limit, the units should be contained in the key, not the value.  For example, rather than `{"instanceLimit": "100/request"}`, should be returned as, `{"instanceLimitPerRequest": "100"}`, if the client exceeds the number of instances that can be created in a single (batch) request. |

```yaml
type: object
properties:
  reason:
    type: string
    description: |-
      The reason of the error. This is a constant value that identifies the
      proximate cause of the error. Error reasons are unique within a particular
      domain of errors. This should be at most 63 characters and match a
      regular expression of `[A-Z][A-Z0-9_]+[A-Z0-9]`, which represents
      UPPER_SNAKE_CASE.
  domain:
    type: string
    description: |-
      The logical grouping to which the "reason" belongs. The error domain
      is typically the registered service name of the tool or product that
      generates the error. Example: "pubsub.googleapis.com". If the error is
      generated by some common infrastructure, the error domain must be a
      globally unique value that identifies the infrastructure. For Google API
      infrastructure, the error domain is "googleapis.com".
  metadata:
    type: object
    additionalProperties:
      type: string
      title: value
    description: |-
      Additional structured details about this error.

      Keys must match a regular expression of `[a-z][a-zA-Z0-9-_]+` but should
      ideally be lowerCamelCase. Also, they must be limited to 64 characters in
      length. When identifying the current value of an exceeded limit, the units
      should be contained in the key, not the value.  For example, rather than
      `{"instanceLimit": "100/request"}`, should be returned as,
      `{"instanceLimitPerRequest": "100"}`, if the client exceeds the number of
      instances that can be created in a single (batch) request.
title: ErrorInfo
description: |-
  Describes the cause of the error with structured details.

  Example of an error when contacting the "pubsub.googleapis.com" API when it
  is not enabled:

  { "reason": "API_DISABLED"
  "domain": "googleapis.com"
  "metadata": {
  "resource": "projects/123",
  "service": "pubsub.googleapis.com"
  }
  }

  This response indicates that the pubsub.googleapis.com API is not enabled.

  Example of an error that is returned when attempting to create a Spanner
  instance in a region that is out of stock:

  { "reason": "STOCKOUT"
  "domain": "spanner.googleapis.com",
  "metadata": {
  "availableRegions": "us-central1,us-east2"
  }
  }
```

</details>

<details id="schema-ErrorInfoErrorDetail">
<summary>ErrorInfoErrorDetail</summary>

Describes the cause of the error with structured details.

Example of an error when contacting the "pubsub.googleapis.com" API when it
is not enabled:

{ "reason": "API_DISABLED"
"domain": "googleapis.com"
"metadata": {
"resource": "projects/123",
"service": "pubsub.googleapis.com"
}
}

This response indicates that the pubsub.googleapis.com API is not enabled.

Example of an error that is returned when attempting to create a Spanner
instance in a region that is out of stock:

{ "reason": "STOCKOUT"
"domain": "spanner.googleapis.com",
"metadata": {
"availableRegions": "us-central1,us-east2"
}
}

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `reason` | No | `string` | The reason of the error. This is a constant value that identifies the proximate cause of the error. Error reasons are unique within a particular domain of errors. This should be at most 63 characters and match a regular expression of `[A-Z][A-Z0-9_]+[A-Z0-9]`, which represents UPPER_SNAKE_CASE. |
| `domain` | No | `string` | The logical grouping to which the "reason" belongs. The error domain is typically the registered service name of the tool or product that generates the error. Example: "pubsub.googleapis.com". If the error is generated by some common infrastructure, the error domain must be a globally unique value that identifies the infrastructure. For Google API infrastructure, the error domain is "googleapis.com". |
| `metadata` | No | `object` | Additional structured details about this error. Keys must match a regular expression of `[a-z][a-zA-Z0-9-_]+` but should ideally be lowerCamelCase. Also, they must be limited to 64 characters in length. When identifying the current value of an exceeded limit, the units should be contained in the key, not the value.  For example, rather than `{"instanceLimit": "100/request"}`, should be returned as, `{"instanceLimitPerRequest": "100"}`, if the client exceeds the number of instances that can be created in a single (batch) request. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  reason:
    type: string
    description: |-
      The reason of the error. This is a constant value that identifies the
      proximate cause of the error. Error reasons are unique within a particular
      domain of errors. This should be at most 63 characters and match a
      regular expression of `[A-Z][A-Z0-9_]+[A-Z0-9]`, which represents
      UPPER_SNAKE_CASE.
  domain:
    type: string
    description: |-
      The logical grouping to which the "reason" belongs. The error domain
      is typically the registered service name of the tool or product that
      generates the error. Example: "pubsub.googleapis.com". If the error is
      generated by some common infrastructure, the error domain must be a
      globally unique value that identifies the infrastructure. For Google API
      infrastructure, the error domain is "googleapis.com".
  metadata:
    type: object
    additionalProperties:
      type: string
      title: value
    description: |-
      Additional structured details about this error.

      Keys must match a regular expression of `[a-z][a-zA-Z0-9-_]+` but should
      ideally be lowerCamelCase. Also, they must be limited to 64 characters in
      length. When identifying the current value of an exceeded limit, the units
      should be contained in the key, not the value.  For example, rather than
      `{"instanceLimit": "100/request"}`, should be returned as,
      `{"instanceLimitPerRequest": "100"}`, if the client exceeds the number of
      instances that can be created in a single (batch) request.
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.ErrorInfo
title: ErrorInfoErrorDetail
description: |-
  Describes the cause of the error with structured details.

  Example of an error when contacting the "pubsub.googleapis.com" API when it
  is not enabled:

  { "reason": "API_DISABLED"
  "domain": "googleapis.com"
  "metadata": {
  "resource": "projects/123",
  "service": "pubsub.googleapis.com"
  }
  }

  This response indicates that the pubsub.googleapis.com API is not enabled.

  Example of an error that is returned when attempting to create a Spanner
  instance in a region that is out of stock:

  { "reason": "STOCKOUT"
  "domain": "spanner.googleapis.com",
  "metadata": {
  "availableRegions": "us-central1,us-east2"
  }
  }
required:
  - "@type"
```

</details>

<details id="schema-Help">
<summary>Help</summary>

Provides links to documentation or for performing an out of band action.

For example, if a quota check failed with an error indicating the calling
project hasn't enabled the accessed service, this can contain a URL pointing
directly to the right place in the developer console to flip the bit.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `links` | No | Array of [`HelpLink`](#schema-HelpLink) | URL(s) pointing to additional information on handling the current error. |

```yaml
type: object
properties:
  links:
    type: array
    items:
      $ref: "#/components/schemas/HelpLink"
    description: URL(s) pointing to additional information on handling the current error.
title: Help
description: |-
  Provides links to documentation or for performing an out of band action.

  For example, if a quota check failed with an error indicating the calling
  project hasn't enabled the accessed service, this can contain a URL pointing
  directly to the right place in the developer console to flip the bit.
```

</details>

<details id="schema-HelpLink">
<summary>HelpLink</summary>

Describes a URL link.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `description` | No | `string` | Describes what the link offers. |
| `url` | No | `string` | The URL of the link. |

```yaml
type: object
properties:
  description:
    type: string
    description: Describes what the link offers.
  url:
    type: string
    description: The URL of the link.
title: Link
description: Describes a URL link.
```

</details>

<details id="schema-HelpErrorDetail">
<summary>HelpErrorDetail</summary>

Provides links to documentation or for performing an out of band action.

For example, if a quota check failed with an error indicating the calling
project hasn't enabled the accessed service, this can contain a URL pointing
directly to the right place in the developer console to flip the bit.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `links` | No | Array of [`HelpLink`](#schema-HelpLink) | URL(s) pointing to additional information on handling the current error. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  links:
    type: array
    items:
      $ref: "#/components/schemas/HelpLink"
    description: URL(s) pointing to additional information on handling the current error.
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.Help
title: HelpErrorDetail
description: |-
  Provides links to documentation or for performing an out of band action.

  For example, if a quota check failed with an error indicating the calling
  project hasn't enabled the accessed service, this can contain a URL pointing
  directly to the right place in the developer console to flip the bit.
required:
  - "@type"
```

</details>

<details id="schema-LocalizedMessage">
<summary>LocalizedMessage</summary>

Provides a localized error message that is safe to return to the user
which can be attached to an operation error.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `locale` | No | `string` | The locale used following the specification defined at https://www.rfc-editor.org/rfc/bcp/bcp47.txt. Examples are: "en-US", "fr-CH", "es-MX" |
| `message` | No | `string` | The localized error message in the above locale. |

```yaml
type: object
properties:
  locale:
    type: string
    description: |-
      The locale used following the specification defined at
      https://www.rfc-editor.org/rfc/bcp/bcp47.txt.
      Examples are: "en-US", "fr-CH", "es-MX"
  message:
    type: string
    description: The localized error message in the above locale.
title: LocalizedMessage
description: |-
  Provides a localized error message that is safe to return to the user
  which can be attached to an operation error.
```

</details>

<details id="schema-LocalizedMessageErrorDetail">
<summary>LocalizedMessageErrorDetail</summary>

Provides a localized error message that is safe to return to the user
which can be attached to an operation error.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `locale` | No | `string` | The locale used following the specification defined at https://www.rfc-editor.org/rfc/bcp/bcp47.txt. Examples are: "en-US", "fr-CH", "es-MX" |
| `message` | No | `string` | The localized error message in the above locale. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  locale:
    type: string
    description: |-
      The locale used following the specification defined at
      https://www.rfc-editor.org/rfc/bcp/bcp47.txt.
      Examples are: "en-US", "fr-CH", "es-MX"
  message:
    type: string
    description: The localized error message in the above locale.
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.LocalizedMessage
title: LocalizedMessageErrorDetail
description: |-
  Provides a localized error message that is safe to return to the user
  which can be attached to an operation error.
required:
  - "@type"
```

</details>

<details id="schema-PreconditionFailure">
<summary>PreconditionFailure</summary>

Describes what preconditions have failed.

For example, if an operation failed because it required the Terms of Service to be
acknowledged, it could list the terms of service violation in the
PreconditionFailure message.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `violations` | No | Array of [`PreconditionFailureViolation`](#schema-PreconditionFailureViolation) | Describes all precondition violations. |

```yaml
type: object
properties:
  violations:
    type: array
    items:
      $ref: "#/components/schemas/PreconditionFailureViolation"
    description: Describes all precondition violations.
title: PreconditionFailure
description: |-
  Describes what preconditions have failed.

  For example, if an operation failed because it required the Terms of Service to be
  acknowledged, it could list the terms of service violation in the
  PreconditionFailure message.
```

</details>

<details id="schema-PreconditionFailureViolation">
<summary>PreconditionFailureViolation</summary>

A message type used to describe a single precondition failure.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `type` | No | `string` | The type of PreconditionFailure. We recommend using a service-specific enum type to define the supported precondition violation subjects. For example, "TOS" for "Terms of Service violation". |
| `subject` | No | `string` | The subject, relative to the type, that failed. For example, "google.com/cloud" relative to the "TOS" type would indicate which terms of service is being referenced. |
| `description` | No | `string` | A description of how the precondition failed. Developers can use this description to understand how to fix the failure. For example: "Terms of service not accepted". |

```yaml
type: object
properties:
  type:
    type: string
    description: |-
      The type of PreconditionFailure. We recommend using a service-specific
      enum type to define the supported precondition violation subjects. For
      example, "TOS" for "Terms of Service violation".
  subject:
    type: string
    description: |-
      The subject, relative to the type, that failed.
      For example, "google.com/cloud" relative to the "TOS" type would indicate
      which terms of service is being referenced.
  description:
    type: string
    description: |-
      A description of how the precondition failed. Developers can use this
      description to understand how to fix the failure.

      For example: "Terms of service not accepted".
title: Violation
description: A message type used to describe a single precondition failure.
```

</details>

<details id="schema-PreconditionFailureErrorDetail">
<summary>PreconditionFailureErrorDetail</summary>

Describes what preconditions have failed.

For example, if an operation failed because it required the Terms of Service to be
acknowledged, it could list the terms of service violation in the
PreconditionFailure message.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `violations` | No | Array of [`PreconditionFailureViolation`](#schema-PreconditionFailureViolation) | Describes all precondition violations. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  violations:
    type: array
    items:
      $ref: "#/components/schemas/PreconditionFailureViolation"
    description: Describes all precondition violations.
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.PreconditionFailure
title: PreconditionFailureErrorDetail
description: |-
  Describes what preconditions have failed.

  For example, if an operation failed because it required the Terms of Service to be
  acknowledged, it could list the terms of service violation in the
  PreconditionFailure message.
required:
  - "@type"
```

</details>

<details id="schema-QuotaFailure">
<summary>QuotaFailure</summary>

Describes how a quota check failed.

For example if a daily limit was exceeded for the calling project,
a service could respond with a QuotaFailure detail containing the project
id and the description of the quota limit that was exceeded.  If the
calling project hasn't enabled the service in the developer console, then
a service could respond with the project id and set `service_disabled`
to true.

Also see RetryInfo and Help types for other details about handling a
quota failure.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `violations` | No | Array of [`QuotaFailureViolation`](#schema-QuotaFailureViolation) | Describes all quota violations. |

```yaml
type: object
properties:
  violations:
    type: array
    items:
      $ref: "#/components/schemas/QuotaFailureViolation"
    description: Describes all quota violations.
title: QuotaFailure
description: |-
  Describes how a quota check failed.

  For example if a daily limit was exceeded for the calling project,
  a service could respond with a QuotaFailure detail containing the project
  id and the description of the quota limit that was exceeded.  If the
  calling project hasn't enabled the service in the developer console, then
  a service could respond with the project id and set `service_disabled`
  to true.

  Also see RetryInfo and Help types for other details about handling a
  quota failure.
```

</details>

<details id="schema-QuotaFailureViolation">
<summary>QuotaFailureViolation</summary>

A message type used to describe a single quota violation.  For example, a
daily quota or a custom quota that was exceeded.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `subject` | No | `string` | The subject on which the quota check failed. For example, "clientip:<ip address of client>" or "project:<Google developer project id>". |
| `description` | No | `string` | A description of how the quota check failed. Clients can use this description to find more about the quota configuration in the service's public documentation, or find the relevant quota limit to adjust through developer console. For example: "Service disabled" or "Daily Limit for read operations exceeded". |
| `apiService` | No | `string` | The API Service from which the `QuotaFailure.Violation` orginates. In some cases, Quota issues originate from an API Service other than the one that was called. In other words, a dependency of the called API Service could be the cause of the `QuotaFailure`, and this field would have the dependency API service name. For example, if the called API is Kubernetes Engine API (container.googleapis.com), and a quota violation occurs in the Kubernetes Engine API itself, this field would be "container.googleapis.com". On the other hand, if the quota violation occurs when the Kubernetes Engine API creates VMs in the Compute Engine API (compute.googleapis.com), this field would be "compute.googleapis.com". |
| `quotaMetric` | No | `string` | The metric of the violated quota. A quota metric is a named counter to measure usage, such as API requests or CPUs. When an activity occurs in a service, such as Virtual Machine allocation, one or more quota metrics may be affected. For example, "compute.googleapis.com/cpus_per_vm_family", "storage.googleapis.com/internet_egress_bandwidth". |
| `quotaId` | No | `string` | The id of the violated quota. Also know as "limit name", this is the unique identifier of a quota in the context of an API service. For example, "CPUS-PER-VM-FAMILY-per-project-region". |
| `quotaDimensions` | No | `object` | The dimensions of the violated quota. Every non-global quota is enforced on a set of dimensions. While quota metric defines what to count, the dimensions specify for what aspects the counter should be increased. For example, the quota "CPUs per region per VM family" enforces a limit on the metric "compute.googleapis.com/cpus_per_vm_family" on dimensions "region" and "vm_family". And if the violation occurred in region "us-central1" and for VM family "n1", the quota_dimensions would be, { "region": "us-central1", "vm_family": "n1", } When a quota is enforced globally, the quota_dimensions would always be empty. |
| `quotaValue` | No | `integer or string` (int64) | The enforced quota value at the time of the `QuotaFailure`. For example, if the enforced quota value at the time of the `QuotaFailure` on the number of CPUs is "10", then the value of this field would reflect this quantity. |
| `futureQuotaValue` | No | `integer or string or null` (int64) | The new quota value being rolled out at the time of the violation. At the completion of the rollout, this value will be enforced in place of quota_value. If no rollout is in progress at the time of the violation, this field is not set. For example, if at the time of the violation a rollout is in progress changing the number of CPUs quota from 10 to 20, 20 would be the value of this field. |

```yaml
type: object
properties:
  subject:
    type: string
    description: |-
      The subject on which the quota check failed.
      For example, "clientip:<ip address of client>" or "project:<Google
      developer project id>".
  description:
    type: string
    description: |-
      A description of how the quota check failed. Clients can use this
      description to find more about the quota configuration in the service's
      public documentation, or find the relevant quota limit to adjust through
      developer console.

      For example: "Service disabled" or "Daily Limit for read operations
      exceeded".
  apiService:
    type: string
    description: |-
      The API Service from which the `QuotaFailure.Violation` orginates. In
      some cases, Quota issues originate from an API Service other than the one
      that was called. In other words, a dependency of the called API Service
      could be the cause of the `QuotaFailure`, and this field would have the
      dependency API service name.

      For example, if the called API is Kubernetes Engine API
      (container.googleapis.com), and a quota violation occurs in the
      Kubernetes Engine API itself, this field would be
      "container.googleapis.com". On the other hand, if the quota violation
      occurs when the Kubernetes Engine API creates VMs in the Compute Engine
      API (compute.googleapis.com), this field would be
      "compute.googleapis.com".
  quotaMetric:
    type: string
    description: |-
      The metric of the violated quota. A quota metric is a named counter to
      measure usage, such as API requests or CPUs. When an activity occurs in a
      service, such as Virtual Machine allocation, one or more quota metrics
      may be affected.

      For example, "compute.googleapis.com/cpus_per_vm_family",
      "storage.googleapis.com/internet_egress_bandwidth".
  quotaId:
    type: string
    description: |-
      The id of the violated quota. Also know as "limit name", this is the
      unique identifier of a quota in the context of an API service.

      For example, "CPUS-PER-VM-FAMILY-per-project-region".
  quotaDimensions:
    type: object
    additionalProperties:
      type: string
      title: value
    description: |-
      The dimensions of the violated quota. Every non-global quota is enforced
      on a set of dimensions. While quota metric defines what to count, the
      dimensions specify for what aspects the counter should be increased.

      For example, the quota "CPUs per region per VM family" enforces a limit
      on the metric "compute.googleapis.com/cpus_per_vm_family" on dimensions
      "region" and "vm_family". And if the violation occurred in region
      "us-central1" and for VM family "n1", the quota_dimensions would be,

      {
      "region": "us-central1",
      "vm_family": "n1",
      }

      When a quota is enforced globally, the quota_dimensions would always be
      empty.
  quotaValue:
    type:
      - integer
      - string
    format: int64
    description: |-
      The enforced quota value at the time of the `QuotaFailure`.

      For example, if the enforced quota value at the time of the
      `QuotaFailure` on the number of CPUs is "10", then the value of this
      field would reflect this quantity.
  futureQuotaValue:
    type:
      - integer
      - string
      - "null"
    format: int64
    description: |-
      The new quota value being rolled out at the time of the violation. At the
      completion of the rollout, this value will be enforced in place of
      quota_value. If no rollout is in progress at the time of the violation,
      this field is not set.

      For example, if at the time of the violation a rollout is in progress
      changing the number of CPUs quota from 10 to 20, 20 would be the value of
      this field.
title: Violation
description: |-
  A message type used to describe a single quota violation.  For example, a
  daily quota or a custom quota that was exceeded.
```

</details>

<details id="schema-QuotaFailureErrorDetail">
<summary>QuotaFailureErrorDetail</summary>

Describes how a quota check failed.

For example if a daily limit was exceeded for the calling project,
a service could respond with a QuotaFailure detail containing the project
id and the description of the quota limit that was exceeded.  If the
calling project hasn't enabled the service in the developer console, then
a service could respond with the project id and set `service_disabled`
to true.

Also see RetryInfo and Help types for other details about handling a
quota failure.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `violations` | No | Array of [`QuotaFailureViolation`](#schema-QuotaFailureViolation) | Describes all quota violations. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  violations:
    type: array
    items:
      $ref: "#/components/schemas/QuotaFailureViolation"
    description: Describes all quota violations.
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.QuotaFailure
title: QuotaFailureErrorDetail
description: |-
  Describes how a quota check failed.

  For example if a daily limit was exceeded for the calling project,
  a service could respond with a QuotaFailure detail containing the project
  id and the description of the quota limit that was exceeded.  If the
  calling project hasn't enabled the service in the developer console, then
  a service could respond with the project id and set `service_disabled`
  to true.

  Also see RetryInfo and Help types for other details about handling a
  quota failure.
required:
  - "@type"
```

</details>

<details id="schema-RequestInfo">
<summary>RequestInfo</summary>

Contains metadata about the request that clients can attach when filing a bug
or providing other forms of feedback.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `requestId` | No | `string` | An opaque string that should only be interpreted by the service generating it. For example, it can be used to identify requests in the service's logs. |
| `servingData` | No | `string` | Any data that was used to serve this request. For example, an encrypted stack trace that can be sent back to the service provider for debugging. |

```yaml
type: object
properties:
  requestId:
    type: string
    description: |-
      An opaque string that should only be interpreted by the service generating
      it. For example, it can be used to identify requests in the service's logs.
  servingData:
    type: string
    description: |-
      Any data that was used to serve this request. For example, an encrypted
      stack trace that can be sent back to the service provider for debugging.
title: RequestInfo
description: |-
  Contains metadata about the request that clients can attach when filing a bug
  or providing other forms of feedback.
```

</details>

<details id="schema-RequestInfoErrorDetail">
<summary>RequestInfoErrorDetail</summary>

Contains metadata about the request that clients can attach when filing a bug
or providing other forms of feedback.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `requestId` | No | `string` | An opaque string that should only be interpreted by the service generating it. For example, it can be used to identify requests in the service's logs. |
| `servingData` | No | `string` | Any data that was used to serve this request. For example, an encrypted stack trace that can be sent back to the service provider for debugging. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  requestId:
    type: string
    description: |-
      An opaque string that should only be interpreted by the service generating
      it. For example, it can be used to identify requests in the service's logs.
  servingData:
    type: string
    description: |-
      Any data that was used to serve this request. For example, an encrypted
      stack trace that can be sent back to the service provider for debugging.
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.RequestInfo
title: RequestInfoErrorDetail
description: |-
  Contains metadata about the request that clients can attach when filing a bug
  or providing other forms of feedback.
required:
  - "@type"
```

</details>

<details id="schema-ResourceInfo">
<summary>ResourceInfo</summary>

Describes the resource that is being accessed.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `resourceType` | No | `string` | A name for the type of resource being accessed, e.g. "sql table", "cloud storage bucket", "file", "Google calendar"; or the type URL of the resource: e.g. "type.googleapis.com/google.pubsub.v1.Topic". |
| `resourceName` | No | `string` | The name of the resource being accessed.  For example, a shared calendar name: "example.com_4fghdhgsrgh@group.calendar.google.com", if the current error is [CodePERMISSION_DENIED][CodePERMISSION_DENIED]. |
| `owner` | No | `string` | The owner of the resource (optional). For example, "user:<owner email>" or "project:<Google developer project id>". |
| `description` | No | `string` | Describes what error is encountered when accessing this resource. For example, updating a cloud project may require the `writer` permission on the developer console project. |

```yaml
type: object
properties:
  resourceType:
    type: string
    description: |-
      A name for the type of resource being accessed, e.g. "sql table",
      "cloud storage bucket", "file", "Google calendar"; or the type URL
      of the resource: e.g. "type.googleapis.com/google.pubsub.v1.Topic".
  resourceName:
    type: string
    description: |-
      The name of the resource being accessed.  For example, a shared calendar
      name: "example.com_4fghdhgsrgh@group.calendar.google.com", if the current
      error is
      [CodePERMISSION_DENIED][CodePERMISSION_DENIED].
  owner:
    type: string
    description: |-
      The owner of the resource (optional).
      For example, "user:<owner email>" or "project:<Google developer project
      id>".
  description:
    type: string
    description: |-
      Describes what error is encountered when accessing this resource.
      For example, updating a cloud project may require the `writer` permission
      on the developer console project.
title: ResourceInfo
description: Describes the resource that is being accessed.
```

</details>

<details id="schema-ResourceInfoErrorDetail">
<summary>ResourceInfoErrorDetail</summary>

Describes the resource that is being accessed.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `resourceType` | No | `string` | A name for the type of resource being accessed, e.g. "sql table", "cloud storage bucket", "file", "Google calendar"; or the type URL of the resource: e.g. "type.googleapis.com/google.pubsub.v1.Topic". |
| `resourceName` | No | `string` | The name of the resource being accessed.  For example, a shared calendar name: "example.com_4fghdhgsrgh@group.calendar.google.com", if the current error is [CodePERMISSION_DENIED][CodePERMISSION_DENIED]. |
| `owner` | No | `string` | The owner of the resource (optional). For example, "user:<owner email>" or "project:<Google developer project id>". |
| `description` | No | `string` | Describes what error is encountered when accessing this resource. For example, updating a cloud project may require the `writer` permission on the developer console project. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  resourceType:
    type: string
    description: |-
      A name for the type of resource being accessed, e.g. "sql table",
      "cloud storage bucket", "file", "Google calendar"; or the type URL
      of the resource: e.g. "type.googleapis.com/google.pubsub.v1.Topic".
  resourceName:
    type: string
    description: |-
      The name of the resource being accessed.  For example, a shared calendar
      name: "example.com_4fghdhgsrgh@group.calendar.google.com", if the current
      error is
      [CodePERMISSION_DENIED][CodePERMISSION_DENIED].
  owner:
    type: string
    description: |-
      The owner of the resource (optional).
      For example, "user:<owner email>" or "project:<Google developer project
      id>".
  description:
    type: string
    description: |-
      Describes what error is encountered when accessing this resource.
      For example, updating a cloud project may require the `writer` permission
      on the developer console project.
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.ResourceInfo
title: ResourceInfoErrorDetail
description: Describes the resource that is being accessed.
required:
  - "@type"
```

</details>

<details id="schema-RetryInfo">
<summary>RetryInfo</summary>

Describes when the clients can retry a failed request. Clients could ignore
the recommendation here or retry when this information is missing from error
responses.

It's always recommended that clients should use exponential backoff when
retrying.

Clients should wait until `retry_delay` amount of time has passed since
receiving the error response before retrying.  If retrying requests also
fail, clients should use an exponential backoff scheme to gradually increase
the delay between retries based on `retry_delay`, until either a maximum
number of retries have been reached or a maximum retry delay cap has been
reached.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `retryDelay` | No | [`Duration`](#schema-Duration) | Clients should wait at least this long between retrying the same request. |

```yaml
type: object
properties:
  retryDelay:
    description: Clients should wait at least this long between retrying the same request.
    $ref: "#/components/schemas/Duration"
title: RetryInfo
description: |-
  Describes when the clients can retry a failed request. Clients could ignore
  the recommendation here or retry when this information is missing from error
  responses.

  It's always recommended that clients should use exponential backoff when
  retrying.

  Clients should wait until `retry_delay` amount of time has passed since
  receiving the error response before retrying.  If retrying requests also
  fail, clients should use an exponential backoff scheme to gradually increase
  the delay between retries based on `retry_delay`, until either a maximum
  number of retries have been reached or a maximum retry delay cap has been
  reached.
```

</details>

<details id="schema-RetryInfoErrorDetail">
<summary>RetryInfoErrorDetail</summary>

Describes when the clients can retry a failed request. Clients could ignore
the recommendation here or retry when this information is missing from error
responses.

It's always recommended that clients should use exponential backoff when
retrying.

Clients should wait until `retry_delay` amount of time has passed since
receiving the error response before retrying.  If retrying requests also
fail, clients should use an exponential backoff scheme to gradually increase
the delay between retries based on `retry_delay`, until either a maximum
number of retries have been reached or a maximum retry delay cap has been
reached.

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `retryDelay` | No | [`Duration`](#schema-Duration) | Clients should wait at least this long between retrying the same request. |
| `@type` | Yes | `string` |  |

```yaml
type: object
properties:
  retryDelay:
    description: Clients should wait at least this long between retrying the same request.
    $ref: "#/components/schemas/Duration"
  "@type":
    type: string
    const: type.googleapis.com/google.rpc.RetryInfo
title: RetryInfoErrorDetail
description: |-
  Describes when the clients can retry a failed request. Clients could ignore
  the recommendation here or retry when this information is missing from error
  responses.

  It's always recommended that clients should use exponential backoff when
  retrying.

  Clients should wait until `retry_delay` amount of time has passed since
  receiving the error response before retrying.  If retrying requests also
  fail, clients should use an exponential backoff scheme to gradually increase
  the delay between retries based on `retry_delay`, until either a maximum
  number of retries have been reached or a maximum retry delay cap has been
  reached.
required:
  - "@type"
```

</details>

<details id="schema-UpdateSecretBody">
<summary>UpdateSecretBody</summary>

UpdateSecretRequest replaces secret material.
service_type must be empty for custom material (server-assigned)

| Property | Required | Type | Description |
| --- | --- | --- | --- |
| `serviceType` | No | `string` | service_type identifies the consuming service. |

```yaml
type: object
allOf:
  - properties:
      serviceType:
        type: string
        description: service_type identifies the consuming service.
  - oneOf:
      - type: object
        properties:
          custom:
            description: |-
              custom stores an opaque value with caller-authored injection metadata.
              UpdateSecretRequest.service_type is server-assigned for custom
              material; leave it empty.
            $ref: "#/components/schemas/CustomSecretMaterial"
        title: custom
        required:
          - custom
      - type: object
        properties:
          oauth:
            description: oauth stores OAuth refresh material.
            $ref: "#/components/schemas/OAuthRefreshMaterial"
        title: oauth
        required:
          - oauth
      - type: object
        properties:
          token:
            description: token stores opaque token material.
            $ref: "#/components/schemas/TokenSecretMaterial"
        title: token
        required:
          - token
title: UpdateSecretRequest
unevaluatedProperties: false
description: |-
  UpdateSecretRequest replaces secret material.
  service_type must be empty for custom material (server-assigned)
```

</details>
