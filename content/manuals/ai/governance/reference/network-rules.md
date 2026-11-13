---
title: Network and HTTP rule syntax
linkTitle: Network rules
weight: 10
description: Syntax for sandbox network rules, including hostname patterns, CIDR ranges, ports, and HTTP method and path rules.
keywords: docker sandboxes, network rules, rule syntax, wildcards, CIDR, HTTP rules, method and path
---

A network rule matches a destination. There are two kinds:

- A host rule matches all traffic to a destination. In local policy, it's a rule
  without `--method`. In an organization policy, it's the **All traffic** rule
  type.
- An HTTP rule names an HTTP method and URL path, with optional query-parameter
  restrictions, so a policy can allow selected API requests. In local policy,
  it's a rule with `--method`. In an organization policy, it's the **HTTP** rule
  type.

To add rules, see [Manage local network rules](../local/network-rules.md) for
local policy or
[Create and manage organization policies](../organization/policies/create-manage.md)
for an organization policy. For every `sbx policy` flag, see the
[`sbx policy` CLI reference](/reference/cli/sbx/policy/).

## Host rules

A host rule matches all traffic to its destination. What each part accepts
depends on where you configure the rule:

| Part        | Organization policy                             | Local policy                                                                     |
| ----------- | ----------------------------------------------- | -------------------------------------------------------------------------------- |
| Destination | One or more hosts, IP addresses, or CIDR ranges | One or more hosts, IP addresses, or CIDR ranges, comma-separated                 |
| Protocol    | `TCP`, `UDP`, or `Both`                         | TCP for allow rules and TCP and UDP for deny rules, unless you pass `--protocol` |

### Destination

A destination is a hostname, an IP address, or a CIDR range. A hostname or IP
address can include a port, such as `example.com:443`. Write an IPv6 address in
brackets with a port, such as `[2001:db8::1]:443`, or as a CIDR range, such as
`2001:db8::1/128`.

#### Hostname patterns

Hostname matching is case-insensitive: `EXAMPLE.com` matches `example.com`.

| Pattern               | Example           | Matches                                                      |
| --------------------- | ----------------- | ------------------------------------------------------------- |
| Exact hostname        | `example.com`     | `example.com` on any port, not subdomains                     |
| Single-level wildcard | `*.example.com`   | One subdomain level, any port: `api.example.com`              |
| Multi-level wildcard  | `**.example.com`  | Any depth, any port: `api.example.com`, `v2.api.example.com`  |
| Hostname with port    | `example.com:443` | `example.com` on port 443 only                                |
| All hosts             | `**`              | Every hostname, any port                                      |

`example.com` and `*.example.com` don't cover each other. Specify both if you
need to match the root domain and its subdomains.

#### CIDR ranges

Both IPv4 and IPv6 notation are supported: `10.0.0.0/8`, `192.168.1.0/24`,
`2001:db8::/32`.

### Protocol

Each rule applies to TCP, UDP, or both. In local policy, the `--protocol` flag
sets this. In an organization policy, the **Protocols** field sets it. UDP
requires [experimental outbound
UDP](../local/network-rules.md#allow-outbound-udp). ICMP is always blocked.

## HTTP rules

An HTTP rule needs a destination and at least one method. Paths and query
parameters are optional: without a path, the rule matches every path, and
without query parameters, it matches any query string. What each part accepts
depends on where you configure the rule:

| Part        | Organization policy                 | Local policy                            |
| ----------- | ----------------------------------- | --------------------------------------- |
| Destination | A host or IP address                | A host                                  |
| Method      | One or more listed methods          | `ANY`, or one or more methods           |
| Path        | Zero or more absolute path patterns | Zero or one absolute path pattern       |
| Query       | Zero or more `name=value` entries   | Zero or more `--query name=value` flags |

### Destination

A hostname destination uses the same [hostname patterns](#hostname-patterns) as
host rules.

A local HTTP rule accepts only a hostname pattern, not an IP address or a CIDR
range. To match an IP address or a CIDR range, use a network rule without
`--method` or `--path`.

### Method

Every rule applies to at least one method:

| Policy                             | Methods you can choose                                                                                                                            |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| Local policy                       | `GET`, `HEAD`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`, `CONNECT`, `TRACE`, or `ANY` for every method                                         |
| Organization policy in Docker Home | `GET`, `HEAD`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`, `Select all` for every listed method, or `Read-only` for `GET`, `HEAD`, and `OPTIONS` |

In local policy, method names are case-insensitive, and `ANY` can't be combined
with specific methods.

### Path

HTTP path matching ignores ASCII letter case: `/API/**` matches `/api/items`.
Non-ASCII characters remain case-sensitive. This applies to both allow and deny
rules.

A path pattern matches like this:

| Pattern            | Example     | Matches                                        |
| ------------------ | ----------- | ---------------------------------------------- |
| Exact path         | `/repos`    | `/repos` only, nothing below it                |
| Segment wildcard   | `/repos/*`  | `/repos/docs`, one path segment only           |
| Recursive wildcard | `/repos/**` | `/repos/docs`, `/repos/docs/issues`, any depth |

Every pattern must start with `/` and can't contain a query string, a fragment,
or a `..` segment.

A local HTTP rule takes one `--path` value, so add a separate rule for each
path. `--path` needs `--method`. Repeating `--path` keeps only the last value,
and a comma is read as part of the path. The value must also be canonical, so it
can't contain percent-encoding, control characters, surrounding whitespace,
repeated or trailing slashes, or a `.` segment.

### Query

An HTTP rule can also restrict query parameters independently of the path. In
local policy, `--query` needs `--method`.

Names and values ignore ASCII letter case, so `Tag=PROD` matches `tag=prod`.
Non-ASCII characters remain case-sensitive.

Use exact parameter names and values, or a single `*` at the start or end of a
value:

| Pattern     | Example         | Matches                                 |
| ----------- | --------------- | --------------------------------------- |
| Exact value | `team=platform` | `team=platform` only                    |
| Prefix      | `ref=release*`  | `ref=release-1.2`, `ref=release`        |
| Suffix      | `env=*-prod`    | `env=eu-prod`                           |
| Any value   | `tag=*`         | Any `tag` value, including an empty one |

A request matches a rule's query parameters like this:

| Behavior                        | Rule query           | Request query           | Matches |
| ------------------------------- | -------------------- | ----------------------- | ------- |
| Extra parameters allowed        | `team=platform`      | `?team=platform&page=2` | Yes     |
| Any order                       | `team=a`, `env=prod` | `?env=prod&team=a`      | Yes     |
| Every parameter required        | `team=a`, `env=prod` | `?team=a`               | No      |
| Repeated names match separately | `tag=a`, `tag=b`     | `?tag=a`                | No      |
| Rule without query parameters   | No entries           | `?team=platform&page=2` | Yes     |

Write names and value patterns using URL query encoding:

| Character        | Write as     |
| ---------------- | ------------ |
| Space            | `+` or `%20` |
| Literal plus     | `%2B`        |
| Literal asterisk | `%2A`        |

For CLI examples, see
[Let an agent read from an API without writing to it](../local/network-rules.md#let-an-agent-read-from-an-api-without-writing-to-it).
For the **Query parameters** field in an organization policy, see
[What a network rule can match](../organization/policies/network.md#http).

## How HTTP and host rules combine

HTTP requests are evaluated against both kinds of rule. A host rule sets the
baseline for a host, and HTTP rules narrow access by method, path, and query
parameters:

| Rules that cover the host | Result for an HTTP request                                                                              |
| ------------------------- | ------------------------------------------------------------------------------------------------------- |
| Host allow only           | Allowed at any method and path                                                                          |
| HTTP allow only           | Allowed only where a rule matches the method, path, and any query restrictions. Anything else is denied |
| Host allow and HTTP deny  | Requests matching the HTTP deny are blocked. The rest stay allowed                                      |
| Host deny                 | Blocked. An HTTP allow can't reopen a denied host                                                       |

HTTP rules apply only to HTTP requests that go through the sandbox's HTTP
proxy. When HTTP rules cover a host, a connection the proxy can't inspect, such
as one it handles transparently, is blocked. Traffic that isn't HTTP, such as
SSH, has no method or path for HTTP rules to match, so control it with host
rules.
