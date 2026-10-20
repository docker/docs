---
title: Network and HTTP rule syntax
weight: 10
description: Syntax for sandbox network rules, including domain patterns, CIDR ranges, ports, and HTTP method and path rules.
keywords: docker sandboxes, network rules, rule syntax, wildcards, CIDR, HTTP rules, method and path
---

Network rules use `connect:tcp` for TCP and `connect:udp` for UDP. Resources are
hostnames, CIDR ranges, or ports. UDP requires
[experimental outbound UDP](../local/network-rules.md#allow-outbound-udp).
ICMP is blocked.

**Hostname patterns**

| Pattern               | Example           | Matches                                                      |
| --------------------- | ----------------- | ------------------------------------------------------------- |
| Exact hostname        | `example.com`     | `example.com` on any port, not subdomains                     |
| Single-level wildcard | `*.example.com`   | One subdomain level, any port: `api.example.com`              |
| Multi-level wildcard  | `**.example.com`  | Any depth, any port: `api.example.com`, `v2.api.example.com`  |
| Hostname with port    | `example.com:443` | `example.com` on port 443 only                                |

`example.com` and `*.example.com` don't cover each other. Specify both if you
need to match the root domain and its subdomains.

**CIDR ranges**

Both IPv4 and IPv6 notation are supported: `10.0.0.0/8`, `192.168.1.0/24`,
`2001:db8::/32`.

## HTTP method and path

A network rule matches a destination host on its own. An HTTP rule is a network
rule that also names an HTTP method and URL path, so a policy can allow reads
from an API without allowing writes to it.

An HTTP rule names one or more methods, a destination, and path patterns. What
each part accepts depends on where you configure the rule:

| Part        | Organization policy                        | Local policy                      |
| ----------- | ------------------------------------------ | --------------------------------- |
| Method      | One or more listed methods                 | `ANY`, or one or more methods     |
| Destination | A host or IP address                       | A host                            |
| Path        | One or more absolute path patterns         | One absolute path pattern         |

A destination can include a port, and a path pattern looks like `/api/**`.

A local rule doesn't accept an IP address or a CIDR range. Use a plain network
rule for those destinations. To cover a second path in a local policy, add a
second rule.

Every rule applies to at least one method. On the CLI, `--method ANY` covers
every HTTP method. In the composer, a rule with no methods selected covers
every method the composer lists. For the methods you can select individually, see
[Add a network rule](../organization/policies/network.md#add-a-network-rule) for an
organization policy and
[HTTP method and path rules](../local/network-rules.md#http-method-and-path-rules)
for a local one.

Path patterns follow the same wildcard rules as filesystem paths, where `*`
matches within one path segment and `**` matches any depth. A pattern without a
wildcard matches that path exactly, so `/repos` matches `/repos` and nothing
below it. Every pattern must start with `/` and can't contain a query string, a
fragment, or a `..` segment. A local rule's path must also be canonical, so it
can't contain percent-encoding, control characters, repeated or trailing
slashes, or a `.` segment. For the full list, see
[HTTP method and path rules](../local/network-rules.md#http-method-and-path-rules).

HTTP requests are evaluated against both layers. A network rule sets the
baseline for a host, and HTTP rules adjust individual methods and paths within
it:

| Rules that cover the host   | Result for an HTTP request                                              |
| --------------------------- | ----------------------------------------------------------------------- |
| Network allow only          | Allowed at any method and path                                          |
| HTTP allow only             | Allowed only where a rule matches the method and path. Anything else is denied |
| Network allow and HTTP deny | The denied methods and paths are blocked. The rest stay allowed          |
| Network deny                | Blocked. An HTTP allow can't reopen a denied host                        |

A network deny is therefore a floor that HTTP rules can't raise, while a
network allow is a ceiling that HTTP rules can carve into.

When a rule requires a method and path decision, the sandbox HTTP proxy
evaluates each request separately instead of deciding once per connection.

Those requests have to go through the proxy. A connection it can't inspect,
such as one it handles transparently, is blocked rather than evaluated.
Traffic that isn't HTTP, such as SSH, carries no method or path, so HTTP rules
never match it. Control those destinations with network rules.

For local and organization policy configuration, see
[Manage local network rules](../local/network-rules.md) and
[Network access policies](../organization/policies/network.md).
