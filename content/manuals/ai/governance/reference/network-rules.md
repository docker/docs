---
title: Network and HTTP rule syntax
linkTitle: Network rules
weight: 10
description: Syntax for sandbox network rules, including hostname patterns, CIDR ranges, ports, and HTTP method and path rules.
keywords: docker sandboxes, network rules, rule syntax, wildcards, CIDR, HTTP rules, method and path
---

A network rule matches a destination. There are two kinds:

- A host rule names no method or path, so it matches all traffic to the
  destination. The Docker Home composer calls it an **All traffic** rule.
- An HTTP rule also names an HTTP method and URL path, with optional
  query-parameter restrictions, so a policy can allow selected API requests.

For how to write destinations, see [Destinations](#destinations).

## HTTP method and path

An HTTP rule names one or more methods, a destination, and path patterns. What
each part accepts depends on where you configure the rule:

| Part        | Organization policy                        | Local policy                      |
| ----------- | ------------------------------------------ | --------------------------------- |
| Method      | One or more listed methods                 | `ANY`, or one or more methods     |
| Destination | A host or IP address                       | A host                            |
| Path        | One or more absolute path patterns         | One absolute path pattern         |

A destination can include a port, and a path pattern looks like `/api/**`.

A local HTTP rule accepts only a host pattern, not an IP address or a CIDR
range. To match an IP address or a CIDR range, use a network rule without
`--method` or `--path`. A local HTTP rule also takes one path, so add a separate
rule for each path.

Every rule applies to at least one method. On the CLI, `--method ANY` covers
every HTTP method. In the composer, select at least one method. A rule with no
methods selected fails validation. For the methods you can select individually,
see
[What a network rule can match](../organization/policies/network.md#what-a-network-rule-can-match)
for an organization policy and
[Allow or block HTTP methods and paths](../local/network-rules.md#allow-or-block-http-methods-and-paths)
for a local one.

HTTP path matching ignores ASCII letter case: `/API/**` matches `/api/items`.
Non-ASCII characters remain case-sensitive. This applies to both allow and deny
rules.

Path patterns use the same wildcards as filesystem paths, where `*`
matches within one path segment and `**` matches any depth. A pattern without a
wildcard matches that path exactly, so `/repos` matches `/repos` and nothing
below it. Every pattern must start with `/` and can't contain a query string, a
fragment, or a `..` segment. A local rule's path must also be canonical, so it
can't contain percent-encoding, control characters, surrounding whitespace,
repeated or trailing slashes, or a `.` segment.

### Query parameters

An HTTP rule can also restrict query parameters independently of the path.

Use exact parameter names and values, or a single `*` at the start or end of a
value to match a suffix or prefix. A value of `*` matches any value, including
an empty one. Names cannot contain wildcards. Names and values ignore ASCII
letter case, so `Tag=PROD` matches `tag=prod`. Non-ASCII characters remain
case-sensitive.

Every specified parameter must match.
Parameter order doesn't matter, and requests may include additional
parameters. Repeating a parameter name in the rule requires separate matching
occurrences in the request. A rule without query restrictions matches
regardless of query parameters.

Write names and value patterns using URL query encoding: `+` or `%20` for a
space, `%2B` for a literal plus, and `%2A` for a literal asterisk. Keep query
parameters out of `--path`. For CLI examples, see
[Allow or block HTTP methods and paths](../local/network-rules.md#allow-or-block-http-methods-and-paths).
For the **Query parameters** field in an organization policy, see
[HTTP rules](../organization/policies/network.md#http-rules).

### Combine HTTP and host rules

HTTP requests are evaluated against both kinds of rule. A host rule sets the
baseline for a host, and HTTP rules narrow access by method, path, and query
parameters:

| Rules that cover the host | Result for an HTTP request                                                                              |
| ------------------------- | ------------------------------------------------------------------------------------------------------- |
| Host allow only           | Allowed at any method and path                                                                          |
| HTTP allow only           | Allowed only where a rule matches the method, path, and any query restrictions. Anything else is denied |
| Host allow and HTTP deny  | Requests matching the HTTP deny are blocked. The rest stay allowed                                      |
| Host deny                 | Blocked. An HTTP allow can't reopen a denied host                                                       |

A host deny is therefore a floor that HTTP rules can't raise, while a host
allow is a ceiling that HTTP rules can carve into.

When a rule requires an HTTP decision, the sandbox HTTP proxy
evaluates each request separately instead of deciding once per connection.

Those requests have to go through the proxy. A connection it can't inspect,
such as one it handles transparently, is blocked rather than evaluated.
Traffic that isn't HTTP, such as SSH, carries no method or path, so HTTP rules
never match it. Control those destinations with host rules.

## Destinations

Network rules use `connect:tcp` for TCP and `connect:udp` for UDP. Resources are
hostnames, CIDR ranges, or ports. UDP requires
[experimental outbound UDP](../local/network-rules.md#allow-outbound-udp).
ICMP is blocked.

### Hostname patterns

Hostname matching is case-insensitive: `EXAMPLE.com` matches `example.com`.

| Pattern               | Example           | Matches                                                      |
| --------------------- | ----------------- | ------------------------------------------------------------- |
| Exact hostname        | `example.com`     | `example.com` on any port, not subdomains                     |
| Single-level wildcard | `*.example.com`   | One subdomain level, any port: `api.example.com`              |
| Multi-level wildcard  | `**.example.com`  | Any depth, any port: `api.example.com`, `v2.api.example.com`  |
| Hostname with port    | `example.com:443` | `example.com` on port 443 only                                |

`example.com` and `*.example.com` don't cover each other. Specify both if you
need to match the root domain and its subdomains.

### CIDR ranges

Both IPv4 and IPv6 notation are supported: `10.0.0.0/8`, `192.168.1.0/24`,
`2001:db8::/32`.
