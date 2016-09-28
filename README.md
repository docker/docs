dhe-deploy
==========

[![Circle CI](https://circleci.com/gh/docker/dhe-deploy.svg?style=shield&circle-token=ac92968ce5afea4db29bfd7410ab732ecfc6dbae)](https://circleci.com/gh/docker/dhe-deploy)
[![codecov](https://codecov.io/gh/docker/dhe-deploy/branch/master/graph/badge.svg?token=lIB2r9EPug)](https://codecov.io/gh/docker/dhe-deploy)


Docker Trusted Registry (DTR) consists of multiple containers in groups we call "replicas". All replicas are identical and contain one etcd, one rethinkdb, one registry/garant, one api server and one nginx. They are deployed using the `dtr` image, which servers as an installer and command line tool for configuring and scaling DTR.

Each released image that DTR uses will exist in the `docker` namespace. Example: `docker/dtr`.

Private development images will reside in the `dockerhubenterprise` namespace and will be suffixed with `-dev`. Example: `dockerhubenterprise/dtr-dev`. When doing a local build the images will be tagged using these names. After building you will use these images to run an instance of DTR.

Contributing
------------

To build and run your own instance of DTR, follow these steps:

1. Clone this repository: `git clone git@github.com:docker/dhe-deploy`
2. Make the images: `make clean && make` from the repository root
3. Install UCP: `docker run -it --rm --name ucp -v /var/run/docker.sock:/var/run/docker.sock docker/ucp:1.1.2 install --fresh-install --admin-username admin --admin-password password --controller-port 444 --host-address 172.17.0.1`
4. Install DTR: `docker run -it --rm dockerhubenterprise/dtr-dev install --ucp-url 172.17.0.1:444 --ucp-username admin --ucp-password password --ucp-insecure-tls --dtr-external-url 172.17.0.1`

To install on UCP, you need to follow their instructions for installing UCP. You also need to set up overlay networking. Then you can use the UCP install command here: https://docker.atlassian.net/wiki/display/DHE/DTR+hacks#DTRhacks-DTR2cli(bootstrapper)

#### Found bugs?
File a Github issue: https://github.com/docker/dhe-deploy/issues/new
* Include screenshots or gifs: http://recordit.co/
* (Optional) Add any related issue labels: https://help.github.com/articles/applying-labels-to-issues-and-pull-requests/

Releases
--------

https://docker.atlassian.net/wiki/display/DHE/DTR+hacks#DTRhacks-Release

Dependencies
------------

dhe-deploy has two types of dependencies: source code and docker images

Golang source dependencies are managed using [Godep](https://github.com/tools/godep) and stored in [/vendor](/vendor). You can see the selected revisions by viewing [Godeps/Godeps.json](Godeps/Godeps.json). If running `go` or `godep` commands locally, make sure that the environment variable `GO15VENDOREXPERIMENT` is set to 1 and that you have go1.5 installed. When saving or restoring go dependencies, use `godep save $(go list ./...| grep -v '/vendor\|bootstrap/migrate\|bootstrap/bootstrap')` or `godep restore $(go list ./... | grep -v '/vendor')` so that you do not install the dependencies of unused vendor packages.

Docker image dependencies are referenced in [constants.go](constants.go), specifically `RethinkdbImageName`, which references a specific tag of the official rethink repository.

**NOTE**: These tags **must** adhere to the following naming scheme and match the versions of their respective Godep-included code:
- For git tags:
  * `docker/dtr-<repository name>:<tag name>`. Example: `docker/dtr-distribution:v2.0.0-rc.3`
- For other git shas **(NOT TO BE INCLUDED IN CUSTOMER RELEASES)**:
  * `docker/dtr-<repository name>:<godep version string>`. Example: `docker/dtr-distribution:1.0.0-rc.1-8-gb9ef615`
