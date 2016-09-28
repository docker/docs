# Changelog

## [v0.3.0](https://github.com/docker/notary/releases/tag/v0.3.0) 5/11/2016
+ Root rotations
+ RethinkDB support as a storage backend for Server and Signer
+ A new TUF repo builder that merges server and client validation
+ Trust Pinning: configure known good key IDs and CAs to replace TOFU.
+ Add --input, --output, and --quiet flags to notary verify command
+ Remove local certificate store. It was redundant as all certs were also stored in the cached root.json
+ Cleanup of dead code in client side key storage logic
+ Update project to Go 1.6.1
+ Reorganize vendoring to meet Go 1.6+ standard. Still using Godeps to manage vendored packages
+ Add targets by hash, no longer necessary to have the original target data available
+ Active Key ID verification during signature verification
+ Switch all testing from assert to require, reduces noise in test runs
+ Use alpine based images for smaller downloads and faster setup times
+ Clean up out of data signatures when re-signing content
+ Set cache control headers on HTTP responses from Notary Server
+ Add sha512 support for targets
+ Add environment variable for delegation key passphrase
+ Reduce permissions requested by client from token server
+ Update formatting for delegation list output
+ Move SQLite dependency to tests only so it doesn't get built into official images
+ Fixed asking for password to list private repositories
+ Enable using notary client with username/password in a scripted fashion
+ Fix static compilation of client
+ Enforce TUF version to be >= 1, previously 0 was acceptable although unused
+ json.RawMessage should always be used as *json.RawMessage due to concepts of addressability in Go and effects on encoding

## [v0.2](https://github.com/docker/notary/releases/tag/v0.2.0) 2/24/2016
+ Add support for delegation roles in `notary` server and client
+ Add `notary CLI` commands for managing delegation roles: `notary delegation`
    + `add`, `list` and `remove` subcommands
+ Enhance `notary CLI` commands for adding targets to delegation roles
    + `notary add --roles` and `notary remove --roles` to manipulate targets for delegations
+ Support for rotating the snapshot key to one managed by the `notary` server
+ Add consistent download functionality to download metadata and content by checksum
+ Update `docker-compose` configuration to use official mariadb image
    + deprecate `notarymysql`
    + default to using a volume for `data` directory
    + use separate databases for `notary-server` and `notary-signer` with separate users
+ Add `notary CLI` command for changing private key passphrases: `notary key passwd`
+ Enhance `notary CLI` commands for importing and exporting keys
+ Change default `notary CLI` log level to fatal, introduce new verbose (error-level) and debug-level settings
+ Store roles as PEM headers in private keys, incompatible with previous notary v0.1 key format
    + No longer store keys as `<KEY_ID>_role.key`, instead store as `<KEY_ID>.key`; new private keys from new notary clients will crash old notary clients
+ Support logging as JSON format on server and signer
+ Support mutual TLS between notary client and notary server

## [v0.1](https://github.com/docker/notary/releases/tag/v0.1) 11/15/2015
+ Initial non-alpha `notary` version
+ Implement TUF (the update framework) with support for root, targets, snapshot, and timestamp roles
+ Add PKCS11 interface to store and sign with keys in HSMs (i.e. Yubikey)
