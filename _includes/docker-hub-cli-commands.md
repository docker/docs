1.  Open a terminal and log into Docker Hub with the Docker CLI:

    ```
    $ docker login

    Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
    Username: gordon
    Password:
    WARNING! Your password will be stored unencrypted in /home/gwendolynne/.docker/config.json.
    Configure a credential helper to remove this warning. See
    https://docs.docker.com/engine/reference/commandline/login/#credentials-store
    ```

2. Search for the `busybox` image:

    ```
    $ docker search busybox

    NAME                        DESCRIPTION                            STARS   OFFICIAL   AUTOMATED
    busybox                     Busybox base image.                    1268    [OK]
    progrium/busybox                                                   66                  [OK]
    hypriot/rpi-busybox-httpd   Raspberry Pi compatible …   41
    radial/busyboxplus          Full-chain, Internet enabled, …   19                       [OK]
    ...
    ```

    > Private repos are not returned at the commandline. Go to the Docker Hub UI
    > to see your allowable repos.

3. Pull the official busybox image to your machine and list it (to ensure it was
   pulled):

    ```
    $ docker pull busybox

    Using default tag: latest
    latest: Pulling from library/busybox
    07a152489297: Pull complete
    Digest: sha256:141c253bc4c3fd0a201d32dc1f493bcf3fff003b6df416dea4f41046e0f37d47
    Status: Downloaded newer image for busybox:latest

    $ docker image ls

    REPOSITORY         TAG        IMAGE ID         CREATED         SIZE
    busybox            latest     8c811b4aec35     11 days ago     1.15MB

    ```

4. Tag the official image (to differentiate it), list it, and push it to your
   personal repo:

    ```
    $ docker tag busybox <DOCKER ID>/busybox:test-tag

    $ docker image ls

    REPOSITORY         TAG        IMAGE ID         CREATED         SIZE
    gordon/busybox     v1         8c811b4aec35     11 days ago     1.15MB
    busybox            latest     8c811b4aec35     11 days ago     1.15MB

    $ docker push <DOCKER ID>/busybox:test-tag
    ```

5. Log out from Docker Hub:

    ```
    $ docker logout
    ```

6. Log on to the [Docker Hub UI](https://hub.docker.com){: target="_blank" class="_"} and view the image you
   pushed.
