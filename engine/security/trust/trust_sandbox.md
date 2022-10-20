---
description: Play in a trust sandbox
keywords: trust, security, root,  keys, repository, sandbox
title: Play in a content trust sandbox
redirect_from:
  - /security/trust/trust_sandbox/
---

This page explains how to set up and use a sandbox for experimenting with trust.
The sandbox allows you to configure and try trust operations locally without
impacting your production images.

Before working through this sandbox, you should have read through the
[trust overview](index.md).

### Prerequisites

These instructions assume you are running in Linux or macOS. You can run
this sandbox on a local machine or on a virtual machine. You need to
have privileges to run docker commands on your local machine or in the VM.

This sandbox requires you to install two Docker tools: Docker Engine >= 1.10.0
and Docker Compose >= 1.6.0. To install the Docker Engine, choose from the
[list of supported platforms](../../install/index.md). To install
Docker Compose, see the
[detailed instructions here](../../../compose/install/index.md).

## What is in the sandbox?

If you are just using trust out-of-the-box you only need your Docker Engine
client and access to the Docker Hub. The sandbox mimics a
production trust environment, and sets up these additional components.

| Service         | Description                                                                                                                   |
| --------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| server          | The service hosts the Notary Server, which manages the TUF database                                                           |
| signer          | The service hosts the Notary Signer, which manages the _timestamp_ and _snapshot_ private keys, and signing operations        |
| trustsandbox    | The service hosts a sandbox instance with Docker, used to communicate with the Notary Server, Signer and the registry server. |
| Registry server | A local registry service.                                                                                                     |

This means you run your own content trust (Notary) server and registry.
If you work exclusively with the Docker Hub, you would not need these components.
They are built into the Docker Hub for you. For the sandbox, however, you build
your own entire, mock production environment.

Within the `trustsandbox` container, you interact with your local registry rather
than the Docker Hub. This means your everyday image repositories are not used.
They are protected while you play.

When you play in the sandbox, you also create root and repository keys. The
sandbox is configured to store all the keys and files inside the `trustsandbox`
container. Since the keys you create in the sandbox are for play only,
destroying the container destroys them as well.

By using a docker-in-docker image for the `trustsandbox` container, you also
don't pollute your real Docker daemon cache with any images you push and pull.
The images are stored in an anonymous volume attached to this container,
and can be destroyed after you destroy the container.

## Build the sandbox

In this section, you use Docker Compose to specify how to set up and link together
the `trustsandbox` container, the Notary server, and the Registry server.

1.  Create a new `trustsandbox` directory and change into it.

    ```console
    $ mkdir trustsandbox
    $ cd trustsandbox
    ```

2.  Clone the official `notary` repository:

    ```console
    $ git clone https://github.com/notaryproject/notary.git
    ```

3.  Create a file called `docker-compose.yml` with your favorite editor and add the following to the new file:

    ```yaml
    services:
      notaryserver:
        image: notary:server-0.6.1-2
        networks:
          - sandbox
        depends_on:
          - notarysigner
        volumes:
          - ./notary:/notarydir
        command: |-
          sh -c "/notarydir/migrations/migrate.sh &&
          cd /notarydir/fixtures &&
          notary-server -config=/notarydir/fixtures/server-config-local.json"

      notarysigner:
        image: notary:signer-0.6.1-2
        networks:
          - sandbox
        volumes:
          - ./notary:/notarydir
        command: |-
          sh -c "/notarydir/migrations/migrate.sh &&
          cd /notarydir/fixtures &&
          notary-signer -config=/notarydir/fixtures/signer-config-local.json"

      trustsandbox:
        image: docker:20.10.20-dind
        networks:
          - sandbox
        volumes:
          - ./notary:/notary
        privileged: true
        entrypoint: ""
        command: |-
            sh -c '
                cp /notary/fixtures/root-ca.crt /usr/local/share/ca-certificates/root-ca.crt &&
                update-ca-certificates &&
                dockerd-entrypoint.sh --insecure-registry sandboxregistry:5000'

      sandboxregistry:
        image: registry:2.8.1
        networks:
          - sandbox

    networks:
      sandbox:
        external: false
    ```

4.  Save and close the file.

5.  Run the containers on your local system.

    ```console
    $ docker compose up -d
    ```

    The first time you run this, the docker-in-docker, Notary server, and registry
    images are downloaded from Docker Hub.

## Play in the sandbox

Now that everything is setup, you can go into your `trustsandbox` container and
start testing Docker content trust. From your host machine, obtain a shell
in the `trustsandbox` container.

  ```console
  $ docker compose exec trustsandbox sh
  /#
  ```

### Test some trust operations

Now, pull some images from within the `trustsandbox` container.

1.  Download a `docker` image to test with.

    ```console
    /# docker pull docker/trusttest
    docker pull docker/trusttest
    Using default tag: latest
    latest: Pulling from docker/trusttest

    b3dbab3810fc: Pull complete
    a9539b34a6ab: Pull complete
    Digest: sha256:d149ab53f8718e987c3a3024bb8aa0e2caadf6c0328f1d9d850b2a2a67f2819a
    Status: Downloaded newer image for docker/trusttest:latest
    ```

2.  Tag it to be pushed to our sandbox registry:

    ```console
    /# docker tag docker/trusttest sandboxregistry:5000/test/trusttest:latest
    ```

3.  Enable content trust.

    ```console
    /# export DOCKER_CONTENT_TRUST=1
    ```

4.  Identify the trust server.

    ```console
    /# export DOCKER_CONTENT_TRUST_SERVER=https://notaryserver:4443
    ```

    This step is only necessary because the sandbox is using its own server.
    Normally, if you are using the Docker Public Hub this step isn't necessary.

5.  Pull the test image.

    ```console
    /# docker pull sandboxregistry:5000/test/trusttest
    Using default tag: latest
    Error: remote trust data does not exist for sandboxregistry:5000/test/trusttest: notaryserver:4443 does not have trust data for sandboxregistry:5000/test/trusttest
    ```

    You see an error, because this content doesn't exist on the `notaryserver` yet.

6.  Push and sign the trusted image.

    ```console
    /# docker push sandboxregistry:5000/test/trusttest:latest
    The push refers to a repository [sandboxregistry:5000/test/trusttest]
    5f70bf18a086: Pushed
    c22f7bc058a9: Pushed
    latest: digest: sha256:7034d197b82fcb07299fda8b05c91d1601ce64f31bc102b1345d03a2953d210a size: 734
    Signing and pushing trust metadata
    You are about to create a new root signing key passphrase. This passphrase
    will be used to protect the most sensitive key in your signing system. Please
    choose a long, complex passphrase and be careful to keep the password and the
    key file itself secure and backed up. It is highly recommended that you use a
    password manager to generate the passphrase and keep it safe. There will be no
    way to recover this key. You can find the key in your config directory.
    Enter passphrase for new root key with ID 27ec255:
    Repeat passphrase for new root key with ID 27ec255:
    Enter passphrase for new repository key with ID 58233f9 (sandboxregistry:5000/test/trusttest):
    Repeat passphrase for new repository key with ID 58233f9 (sandboxregistry:5000/test/trusttest):
    Finished initializing "sandboxregistry:5000/test/trusttest"
    Successfully signed "sandboxregistry:5000/test/trusttest":latest
    ```

    Because you are pushing this repository for the first time, Docker creates
    new root and repository keys and asks you for passphrases with which to
    encrypt them. If you push again after this, it only asks you for repository
    passphrase so it can decrypt the key and sign again.

7.  Try pulling the image you just pushed:

    ```console
    /# docker pull sandboxregistry:5000/test/trusttest
    Using default tag: latest
    Pull (1 of 1): sandboxregistry:5000/test/trusttest:latest@sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926
    sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926: Pulling from test/trusttest
    Digest: sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926
    Status: Downloaded newer image for sandboxregistry:5000/test/trusttest@sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926
    Tagging sandboxregistry:5000/test/trusttest@sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926 as sandboxregistry:5000/test/trusttest:latest
    ```

### Test with malicious images

What happens when data is corrupted and you try to pull it when trust is
enabled? In this section, you go into the `sandboxregistry` and tamper with some
data. Then, you try and pull it.

1.  Leave the `trustsandbox` shell and container running.

2.  Open a new interactive terminal from your host, and obtain a shell into the
    `sandboxregistry` container.

    ```console
    $ docker compose exec sandboxregistry sh
    /#
    ```

3.  List the layers for the `test/trusttest` image you pushed:

    ```console
    /# ls -l /var/lib/registry/docker/registry/v2/repositories/test/trusttest/_layers/sha256
    total 12
    drwxr-xr-x    2 root     root          4096 Oct 20 15:25 4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1
    drwxr-xr-x    2 root     root          4096 Oct 20 15:25 cc7629d1331a7362b5e5126beb5bf15ca0bf67eb41eab994c719a45de53255cd
    drwxr-xr-x    2 root     root          4096 Oct 20 15:25 fbe94c69c9b964ffe67daa6e1654ab4f1db1dc14decc30f7c03ae43284cfa723
    ```

4.  Change into the registry storage for one of those layers (this is in a different directory):

    ```console
    /# cd /var/lib/registry/docker/registry/v2/blobs/sha256/cc/cc7629d1331a7362b5e5126beb5bf15ca0bf67eb41eab994c719a45de53255cd
    ```

5.  Add malicious data to one of the `trusttest` layers:

    ```console
    /# echo "Malicious data" > data
    ```

6.  Go back to your `trustsandbox` terminal.

7.  List the `trusttest` image.

    ```console
    /# docker image ls | grep trusttest
    REPOSITORY                            TAG                 IMAGE ID            CREATED             SIZE
    docker/trusttest                      latest    cc7629d1331a   7 years ago   5.03MB
    sandboxregistry:5000/test/trusttest   latest    cc7629d1331a   7 years ago   5.03MB
    ```

8.  Remove the `trusttest:latest` image from our local cache.

    ```console
    /# docker image rm -f cc7629d1331a
    Untagged: docker/trusttest:latest
    Untagged: sandboxregistry:5000/test/trusttest:latest
    Untagged: sandboxregistry:5000/test/trusttest@sha256:ebf59c538accdf160ef435f1a19938ab8c0d6bd96aef8d4ddd1b379edf15a926
    Deleted: sha256:cc7629d1331a7362b5e5126beb5bf15ca0bf67eb41eab994c719a45de53255cd
    Deleted: sha256:2a1f6535dc6816ffadcdbe20590045e6cbf048d63fd4cc753a684c9bc01abeea
    Deleted: sha256:c22f7bc058a9a8ffeb32989b5d3338787e73855bf224af7aa162823da015d44c
    ```

    Docker does not re-download images that it already has cached, but we want
    Docker to attempt to download the tampered image from the registry and reject
    it because it is invalid.

9.  Pull the image again. This downloads the image from the registry, because we don't have it cached.

    ```console
    /# docker pull sandboxregistry:5000/test/trusttest
    Using default tag: latest
    Pull (1 of 1): sandboxregistry:5000/test/trusttest:latest@sha256:7034d197b82fcb07299fda8b05c91d1601ce64f31bc102b1345d03a2953d210a
    sandboxregistry:5000/test/trusttest@sha256:7034d197b82fcb07299fda8b05c91d1601ce64f31bc102b1345d03a2953d210a: Pulling from test/trusttest
    fbe94c69c9b9: Pull complete 
    4f4fb700ef54: Pull complete 
    error pulling image configuration: download failed after attempts=6: unexpected EOF
    ```

    The pull did not complete because the trust system couldn't verify the
    image.

## More play in the sandbox

Now, you have a full Docker content trust sandbox on your local system,
feel free to play with it and see how it behaves. If you find any security
issues with Docker, feel free to send us an email at <security@docker.com>.

## Clean up your sandbox

When you are done, and want to clean up all the services you've started and any
anonymous volumes that have been created, just run the following command in the
directory where you've created your Docker Compose file:

  ```console
  $ docker compose down -v
  ```
