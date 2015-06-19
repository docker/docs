# Notary

## Overview

Notary Server manages trust metadata as a complementary service to the registry.
It implements all endpoints under the `_trust` segment of the registry URLs.
Notary Server expects to manage TUF metadata and will do validation of one parent
level of content for any data uploaded to ensure repositories do not become
corrupted. This means either the keys in the root.json file will be used to
validate the uploaded role, or the keys in the immediate delegate parent will
be used.

Uploading a new root.json will be validated using the same token mechanism
present in the registry. A user having write permissions on a repository
will be sufficient to permit the uploading of a new root.json.

## Timestamping

TUF requires a timestamp file be regularly generated. To achieve any ease
of use, it is necessary that Notary Server is responsible for generating the
timestamp.json based on the snapshot.json created and uploaded by the
repository owner.

It is bad policy to place any signing keys in frontline servers. While
Notary Server is capable of supporting this behaviour we recommend using a
separate service and server with highly restricted permissions. Rufus
is provided as a reference implementation of a remote signer. An
implementation that satisfies the gRPC interface defined in Rufus will
satisfy Notary Server's requirements.

# Running

`# docker-compose build`
`# docker-compose up`
