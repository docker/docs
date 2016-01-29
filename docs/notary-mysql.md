<!--[metadata]>
+++
title = "Notary MySQL"
description = "Description of the Notary MySQL"
keywords = ["docker, notary, notary-mysql"]
[menu.main]
parent="mn_notary"
+++
<![end-metadata]-->

# Notary MySQL

The Notary MySQL is one of the backends for [Notary Server](notary-server.md) and
[Notary Signer](notary-signer.md).

### Recommendation
For security, especially in production deployments, one should create users
with restricted permissions and separate databases for the `server` and
`signer` since the `signer` only needs the `private_keys` table, and the
`server` only needs `timestamp_keys` and `tuf_files`.

We use such a setup in our compose file to provide people with more accurate
guidance in deploying their own instances.
