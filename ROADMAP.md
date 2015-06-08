# Roadmap

The Trust project consists of a number of moving parts of which Vetinari is one. Vetinari is the front line metadata service
that clients interact with. It manages TUF metadata and interacts with a pluggable signing service to issue new TUF timestamp
files.

The Rufus repository is provided as our reference implementation of a signing service. It supports HSMs along with Ed25519
software signing.


## Trust Goals

- Provide a trust service that enables users to validate the integrity and provenance of content they download.
- Provide trust management such that a user can lock down the scope of their trust to individual signers if
  they so desire.
- Provide guarantees that servers used to store and distribute content have very limited scope and opportunity to tamper with the data.
- Provide a trust mechanism that can support insecure plain HTTP mirrors while still providing security to end consumers
