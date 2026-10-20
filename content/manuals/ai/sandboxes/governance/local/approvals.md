---
title: Respond to approval requests
linkTitle: Respond to approvals
weight: 30
description: Approve or deny sandbox network requests that need your confirmation.
keywords: docker sandboxes, approval required, approve network access, sbx policy, local policy
---

When a destination needs approval, a sandbox can't reach it until you confirm
it. The request is blocked and the sandbox receives a message naming the
destination:

```plaintext
Approval required for api.example.com.

Review and respond with:
  sbx policy approval ls
```

If your organization
[configures a support message](../organization/policies/create-manage.md#configure-a-support-message), it
appears after the approval instructions.

The request that triggers the prompt doesn't wait for an answer. It's denied,
and approving the destination affects later requests. Agents that retry a
failed request pick up the new access on their next attempt. For others, run
the operation again.

List the destinations waiting for a response:

```console
$ sbx policy approval ls
APPROVAL                                              SANDBOX      TITLE                 DETAIL                                                                              OPTIONS
network:cWtN-4xUrNjxh4ouezcstgjrky6Rg57QfeRe3WEkyyY   my-sandbox   api.example.com:443   Protocol: TCP Resource type: domain approval required by policy "default network"   allow (Allow), dismiss (Dismiss)
```

Each entry names the destination, the sandbox that asked for it, and why it
needs approval. Under organization governance that reason names the policy, as
shown. Without it, the reason is that no matching allow rule covers the
destination. A request that an HTTP rule matches also shows the method and
path, such as `GET api.example.com:443/v1/data`.

To inspect a single entry, pass its ID to `sbx policy approval inspect`:

```console
$ sbx policy approval inspect network:cWtN-4xUrNjxh4ouezcstgjrky6Rg57QfeRe3WEkyyY
```

Respond by selecting one of the options the entry offers:

```console
$ sbx policy approval respond network:cWtN-4xUrNjxh4ouezcstgjrky6Rg57QfeRe3WEkyyY --option allow
```

Choosing `allow` grants access to that destination. Choosing `dismiss` leaves
it blocked, and the destination is requested again the next time the sandbox
tries to reach it. Repeated attempts collapse into a single entry, so a sandbox
retrying in a loop leaves one request to answer, not a queue of duplicates.

## What approving grants

Approving records a rule that allows the destination the sandbox actually
asked for, scoped to the sandbox that asked. Three things follow from that:

- The rule covers one destination, not the pattern the policy rule used. A
  policy that allows `*.example.com` with approval asks about
  `api.example.com` and `cdn.example.com` separately.
- A destination includes its port, so `api.example.com:443` and
  `api.example.com:8443` are approved separately.
- Another sandbox reaching the same destination asks again.

When an [HTTP rule](../reference/network-rules.md#http-method-and-path) in the policy matches
the request, approving covers the method and exact path the sandbox requested
rather than the whole destination. `GET /v1/data` and `POST /v1/data` on the
same host are approved separately. A request whose method or path can't be
recorded as a rule, such as a path with percent-encoding, is blocked without an
entry to respond to.

Approved destinations stay allowed until the rule is removed. List them with
`sbx policy ls --wide --created-via approval`, and remove one the same way as
any other local rule, with [`sbx policy rm network`](network-rules.md).

Approvals live in the local policy store, so [`sbx policy reset`](get-started.md#resetting)
removes all of them along with your other local rules. Each destination is
requested again the next time a sandbox reaches it.

## Troubleshooting

### A request is blocked with "Approval required"

The destination needs your confirmation. Either no allow or deny rule matches
it and your machine isn't under organization governance, or an organization
policy allows it but requires approval first.

Run `sbx policy approval ls` to see the pending request and respond to it with
`sbx policy approval respond`. Approving applies to later requests, not the one
that was blocked, so run the operation again afterward.
