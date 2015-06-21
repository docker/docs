# Notary is still a work in progress and we invite contributions and reviews from the security community. It will need to go through a formal security review process before it should be used in production.

# Notary Server

Notary Server manages TUF data over an HTTP API compatible with the
[notary client](../notary/). The API is defined as part of the [distribution
specification](https://github.com/docker/distribution/blob/master/docs/spec/api.md)
as it is intended to be used side by side with a docker registry and
the URL paths are designed to be compatible with the existing registry
URL structure.

It may be configured to use either JWT or HTTP Basic Auth for authentication.
Currently it only supports MySQLfor store of the TUF data, we intend to
expand this to other storage options.

## Setup for Development

The notary repository comes with Dockerfiles and a docker-compose file
to faciliate development. Simply run the following commands to start
a notary server with a temporary MySQL database in containers:

```
$ docker-compose build
$ docker-compose up
```

If you are on Mac OSX with boot2docker or kitematic, you'll need to
update your hosts file such that the name `notary` is associated with
the IP address of your VM (for boot2docker, this can be determined
by running `boot2docker ip`, with kitematic, `echo $DOCKER_HOST` should
show the IP of the VM). If you are using the default Linux setup,
you need to add `127.0.0.1 notary` to your hosts file.

## Compiling Notary Server

From the root of this git repository, run `make binaries`. This will
compile the notary and notary-server applications and place them in
a `bin` directory at the root of the git repository (the `bin` directory
is ignored by the .gitignore file).

## Running Notary Server

The `notary-server` application has the following usage:

```
$ bin/notary-server --help
usage: bin/notary-serve
  -config="": Path to configuration file
  -debug=false: Enable the debugging server on localhost:8080
```

## Configuring Notary Server

The configuration file must be a json file with the following format:

```json
{
    "server": {
        "addr": ":4443",
        "tls_cert_file": "./fixtures/notary.pem",
        "tls_key_file": "./fixtures/notary.key"
    }
    "logging": {
        "level": 5
    }
}
```

The pem and key provided in fixtures are purely for local development and
testing. For production, you must create your own keypair and certificate,
either via the CA of your choice, or a self signed certificate.
