---
title: Fetch
description: Fetch content from remote URLs
---

This toolset allows your agent to fetch content from HTTP and HTTPS URLs. It
supports multiple URLs in a single call, respects robots.txt restrictions, and
can return content in different formats (text, markdown, or HTML).

## Usage

```yaml
toolsets:
  - type: fetch
    timeout: 30 # Optionally set a default timeout for HTTP requests
```

## Tools

| Name    | Description                                                      |
| ------- | ---------------------------------------------------------------- |
| `fetch` | Fetch content from one or more HTTP/HTTPS URLs with metadata     |

### fetch

Fetches content from one or more HTTP/HTTPS URLs and returns the response body
along with metadata such as status code, content type, and content length.

Args:

- `urls`: Array of URLs to fetch (required)
- `format`: The format to return the content in - `text`, `markdown`, or `html`
  (required)
- `timeout`: Request timeout in seconds, default is 30, maximum is 300
  (optional)

**Features:**

- Support for multiple URLs in a single call
- Returns response body and metadata (status code, content type, length)
- Automatic HTML to markdown or text conversion based on format parameter
- Respects robots.txt restrictions
- Configurable timeout per request

**Example:**

Single URL:
```json
{
  "urls": ["https://example.com"],
  "format": "markdown",
  "timeout": 60
}
```

Multiple URLs:
```json
{
  "urls": ["https://example.com", "https://another.com"],
  "format": "text"
}
```
