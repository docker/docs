---
title: Configure the Docker client to sign images
description:
---

After an administrator
[configures the UCP and DTR servers for content trust](admin_tasks.md) and
[delegates users to be able to sign images](admin_tasks.md#delegate-image-signing),
each of those users needs to configure their system for image signing. This
topic covers the required configuration steps, as well as how to sign and
push images.

Each user who will sign images needs to follow these steps.

## Import the user's signing key

Users who need to sign images should import the `key.pem` file from their UCP
client bundle into Docker. The user probably downloaded the client bundle when
they sent the administrator their `cert.pem` so the administrator could
[delegate signing to them](#delegate-image-signing).

Import the private key associated with the user certificate. You must specify
the trust directory `~/.docker/private`. If the `~/.docker/private` directory
does not yet exist, Notary will create it.

- **Linux or macOS**:

  ```bash
  $ notary -d ~/.docker/trust key import /path/to/key.pem
  ```

- **Windows**:

  ```powershell
  PS C:\> notary -d ~/.docker/trust key import /path/to/key.pem
  ```

You are prompted for a passphrase. Save it in a secure location such as a
password manager. You will need to provide the passphrase each time you sign an
image.

## Configure the Docker client

These steps may need to be performed for each Docker user who will push images
to the trusted repository, and also on each Docker client which should only be
allowed to pull and use trusted images.

### Linux or MacOS

1.  **Required**: Set the `DOCKER_CONTENT_TRUST` environment variable to `1`.
    You can do this temporarily or permanently.

    - To set the environment variable for the current command-line session, type
      the following into the terminal where you will be running `docker` commands:

      ```bash
      $ export DOCKER_CONTENT_TRUST=1
      ```

      This environment variable will be effective until you close the command
      prompt.

    - To set the environmment variable for just a single command, add it before
      the command:

      ```bash
      $ DOCKER_CONTENT_TRUST=1 docker pull...
      ```

    - To set the environment variable permanently, edit the `~/.profile` file
      and add the following line:

      ```bash
      export DOCKER_CONTENT_TRUST=1
      ```

2.  **If your DTR instance uses certificates not signed by a public certificate authority (CA)**:
    Configure the local Docker daemon and client to trust the DTR server's
    certificate. You need to do this step if you see an error like the following
    when you try to [sign and push an image](#sign-and-push-an-image).

    ```none
    x509: certificate signed by unknown authority
    ```

    This procedure is different if you are on a Linux or macOS client:

    - **Linux**:

      1.  Download the certificate add it to a subdirectory of the
          `/etc/docker/certs.d/` directory.

          ```bash
          $ sudo mkdir -p /etc/docker/certs.d/<dtr-domain-name>

          $ curl -k https://<host_or_ip_of_dtr_host>/ca -o <dtr-domain-name>.crt

          $ sudo mv <dtr-domain-name>.crt /etc/docker/certs.d/<dtr-domain-name>/ca.crt
          ```

      2.  Configure the Docker client to use certificates available to the
          Docker daemon by creating a symbolic link from `/etc/docker/certs.d/`
          to `~/.docker/tls/`

          ```bash
          $ ln -s /etc/docker/certs.d ~/.docker/tls
          ```

      3.  Restart Docker using one of the following commands:

          - `sudo systemctl restart docker`
          - `sudo service docker restart`

    - **macOS**:

      1.  Download the certificate and name the output file
          `<dtr-domain-name>.crt`.

          ```bash
          $ curl -k https://<host_or_ip_of_dtr_host>/ca -o <dtr-domain-name>.crt
          ```

      2.  Import the certificate into the macOS keychain. This example uses the
          command line, but you can use the **Keychain Access** application
          instead.

          ```bash
          $ sudo security add-trusted-cert -d \
              -r trustRoot \
              -k /Library/Keychains/System.keychain \
              <dtr-domain-name>.crt
          ```

      3.  Restart Docker for Mac. Click the Docker icon in the toolbar and click
          **restart**.

The Docker daemon and client now trust the DTR server. Continue to
[Sign and push an image](#sign-and-push-an-image).

### Windows

1.  Set the `DOCKER_CONTENT_TRUST` environment variable to `1`. You can do this
    temporarily or permanently.

    - To set the environment variable for the current PowerShell session, type the
      following into the PowerShell terminal where you will be running `docker` commands:

      ```powershell
      PS C:\> $env:DOCKER_CONTENT_TRUST = "1"
      ```

      This environment variable will be effective until you close the PowerShell
      session.

    - To set the environment variable permanently for the logged-in user, use the
      following command:

      ```powershell
      PS C:\> [Environment]::SetEnvironmentVariable("DOCKER_CONTENT_TRUST", "1", "User")
      ```

      The variable is set immediately.

      Whichever method you use, you can verify that the environment variable is set
      by typing `$Env:DOCKER_CONTENT_TRUST` at the command line.

2.  **If your DTR instance uses certificates not signed by a public certificate authority (CA)**:
    Configure the local Docker daemon and client to trust the DTR server's
    certificate. You need to do this step if you see an error like the following
    when you try to [sign and push an image](#sign-and-push-an-image).

    1.  Download the certificate by browsing to the URL
        `https://<dtr-domain-name>/ca`. The certificate is shown in the browser
        as a text file. Choose **File** / **Save As** and save the file as
        `<dtr-domain-name>.crt`.

    2.  Open Windows Explorer and go to the directory where you saved the file.
        Right-click `<dtr-domain-name>.crt` and choose **Install certificate**.

        - Select **Local machine** for the store location.

        - Select **Place all certificates in the following store**.

        - Click **Browse** and select **Trusted Root Certificate Authorities**.

        - Click **Finish**.

    3.  Restart Docker for Windows. Click the Docker icon in the Notifications
        area and click **Settings**. Click **Reset** and choose
        **Restart Docker**.

The Docker daemon and client now trust the DTR server. Continue to
[Sign and push an image](#sign-and-push-an-image).

## Sign and push an image

After [Configuring the signer's Notary and Docker clients](#onfigure-the-signers-otary-and-ocker-clients),
the user can sign and push images to Docker Trusted Registry. These steps are
the same on Linux, macOS, or Windows.

1.  Log into DTR.

    ```bash
    $ docker login <dtr_url>
    ```

    You are prompted for your DTR credentials.

2.  Tag the image with a tag in the format `<GUN>:imagename`. The following
    example tags the `ubuntu:16.04` image as `ubuntu` in your trusted
    repository. This will signal to the `docker push` command that the image tag
    contains a repository.

    ```bash
    $ docker tag ubuntu:16.04 dtr-example.com/engineering/testrepo:ubuntu
    ```

3.  Sign and push the tagged image, so that your deployments can use it. The
    following example signs and pushes the image created in the previous step.
    You are prompted for the delegation key passphrase.

    ```bash
    $ docker push dtr-example.com/engineering/testrepo:ubuntu

    The push refers to a repository [dtr-example.com/engineering/testrepo]
    5eb5bd4c5014: Pushed
    d195a7a18c70: Pushed
    af605e724c5a: Pushed
    59f161c3069d: Pushed
    4f03495a4d7d: Pushed
    ubuntu: digest: sha256:4c0b138bdaaefa6a1c290ba8d8a97a568f43c0f8f25c733af54d3999da12dfd4 size: 1357
    Signing and pushing trust metadata
    Enter passphrase for delegation key with ID ff97e18:
    Successfully signed "dtr-example.com/engineering/testrepo":ubuntu
    ```

4.  To test pulling the image, remove it locally, then pull it.

    ```bash
    $ docker image remove dtr-example.com/engineering/testrepo:ubuntu

    $ docker pull dtr-example.com/engineering/testrepo:ubuntu
    ```

5.  You can verify that the image exists in the repository using the DTR web UI.
    Go to the DTR web UI and click **Repositories**. Choose the repository and
    go to **Images**.

    ![Trusted image in DTR repository](/datacenter/images/signed_image_in_dtr.png)

The signed, trusted image is available in your trusted repository.

## Where to go next

* [Restrict services to worker nodes](restrict-services-to-worker-nodes.md)
