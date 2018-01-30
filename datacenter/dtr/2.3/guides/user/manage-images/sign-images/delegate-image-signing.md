---
title: Delegate image signing
description: Learn how to grant permission for others to sign images in Docker Trusted Registry.
keywords: registry, sign, trust
---

Instead of signing all the images yourself, you can delegate that task
to other users.

A typical workflow looks like this:

1. A repository owner creates a repository in DTR, and initializes the trust
metadata for that repository
3. Team members download a UCP client bundle and share their public key
certificate with the repository owner
4. The repository owner delegates signing to the team members
5. Team members can sign images using the private keys in their UCP client
bundles

In this example, the IT ops team creates and initializes trust for the
`dev/nginx`. Then they allow users in the QA team to push and sign images in
that repository.

![teams](../../../images/delegate-image-signing-1.svg)

## Create a repository and initialize trust

A member of the IT ops team starts by configuring their
[Notary CLI client](../../access-dtr/configure-your-notary-client.md).

Then they create the `dev/nginx` repository,
[initialize the trust metadata](index.md) for that repository, and grant
write access to members of the QA team, so that they can push images to that
repository.

## Ask for the public key certificates

The member of the IT ops team then asks the QA team for their public key
certificate files that are part of their UCP client bundle.

If they don't have a UCP client bundle,
[they can download a new one](/datacenter/ucp/2.2/guides/user/access-ucp/cli-based-access.md).

## Delegate image signing

When delegating trust, you associate a public key certificate with a role name.
UCP requires that you delegate trust to two different roles:

* `targets/releases`
* `targets/<role>`, where `<role>` is the UCP team the user belongs to

In this example we delegate trust to `targets/releases` and `targets/qa`:

```none
# Delegate trust, and add that public key with the role targets/releases
notary delegation add --publish \
  dtr.example.org/dev/nginx targets/releases \
  --all-paths <user-1-cert.pem> <user-2-cert.pem>

# Delegate trust, and add that public key with the role targets/admin
notary delegation add --publish \
  dtr.example.org/dev/nginx targets/qa \
  --all-paths <user-1-cert.pem> <user-2-cert.pem>
```

Now members from the QA team just need to [configure their Notary CLI client
with UCP private keys](../../access-dtr/configure-your-notary-client.md)
before [pushing and signing images](index.md) into the `dev/nginx` repository.

## Where to go next

* [Manage trusted repositories](manage-trusted-repositories.md)
