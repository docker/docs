---
title: Install Docker Trusted Registry
description: Learn how to install Docker Trusted Registry for production.
keywords: dtr, registry, install
---

>{% include enterprise_label_shortform.md %}

Docker Trusted Registry (DTR) is a containerized application that runs on a
swarm managed by the Universal Control Plane (UCP). It can be installed
on-premises or on a cloud infrastructure.

## Step 1. Validate the system requirements

Before installing DTR, make sure your
infrastructure meets the [system requirements](./system-requirements) that DTR
needs to run.

## Step 2. Install UCP

DTR requires Docker Universal Control Panel (UCP) to run. If UCP is not yet installed, refer to [install UCP for production](/ee/ucp/admin/install/).

>**Note**
>
> Prior to installing DTR:
> * When upgrading, upgrade UCP before DTR for each major version. For example,
> if you are upgrading four major versions, upgrade one major version at a
> time, first UCP, then DTR, and then repeat for the remaining three versions.
> * UCP upgraded to the most recent version before an
> initial install of DTR.
> * Docker Engine should be updated to the most recent version before
> installing or updating UCP.

DTR and UCP must not be installed on the same node, due to the potential for resource and port conflicts. Instead, install DTR on worker nodes that will be managed by UCP. Note also that DTR cannot be installed on a standalone Docker Engine.

![](../../images/install-dtr-1.svg)


## Step 3. Install DTR

1. Once UCP is installed, navigate to the UCP web interface as an admin. Expand
your profile on the left navigation pane, and select **Admin Settings > Docker Trusted Registry**.

    ![](../../images/install-dtr-2.png){: .with-border}

2. After you configure all the options, you should see a Docker CLI command that you can use to install DTR. Before you run the command, take note of the `--dtr-external-url` parameter:

      ```bash
      $ docker run -it --rm \
        {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} install \
        --dtr-external-url <dtr.example.com> \
        --ucp-node <ucp-node-name> \
        --ucp-username admin \
        --ucp-url <ucp-url>
      ```

    If you want to point this parameter to a load balancer that uses HTTP for
    health probes over port `80` or `443`, temporarily reconfigure the load
    balancer to use TCP over a known open port. Once DTR is installed, you can configure the load balancer however you need to.

3. Run the DTR install command on any node connected to the UCP cluster, and with the Docker Engine installed. DTR will not be installed on the node you run
the install command on. DTR will be installed on the UCP worker defined by the
`--ucp-node` flag.

    For example, you could SSH into a UCP node and run the DTR install command
    from there.  Running the installation command in interactive TTY or `-it` mode
    means you will be prompted for any required additional information.  [Learn more about installing DTR](/reference/dtr/2.7/cli/install/).

    Here are some useful options you can set during installation:
    * To install a specific version of DTR, replace `{{ page.dtr_version }}` with
    your desired version in the [installation command](#step-3-install-dtr)
    above. Find all DTR versions in the [DTR release notes](/ee/dtr/release-notes/)
    page.

    * DTR is deployed with self-signed certificates by default, so UCP might not be
    able to pull images from DTR. Use the `--dtr-external-url <dtr-domain>:<port>`
    optional flag during installation, or during a reconfiguration, so that UCP is
    automatically reconfigured to trust DTR.

    * With DTR 2.7, you can [enable browser authentication via client
    certificates](/ee/enable-client-certificate-authentication/) at install
    time. This bypasses the DTR login page and hides the logout button, thereby
    skipping the need for entering your username and password.

4. Verify that DTR is installed. Either:
    *  See `https://<ucp-fqdn>/manage/settings/dtr`, or;

    *  Navigate to **Admin Settings > Docker Trusted Registry** from the UCP web UI. Under the hood, UCP modifies `/etc/docker/certs.d` for each host and adds DTR's CA certificate. UCP can then pull images from DTR because the Docker Engine for each node in the UCP swarm has been configured to trust DTR.

5. Reconfigure your load balancer back to your desired protocol and port.


## Step 4. Check that DTR is running

1. In your browser, navigate to the UCP web interface.

2. Select **Shared Resources > Stacks** from the left navigation pane. You should see DTR listed as a stack.

3. To verify that DTR is accessible from the browser, enter your DTR IP address or FQDN on the address bar.
Since [HSTS (HTTP Strict-Transport-Security)
header](https://en.wikipedia.org/wiki/HTTP_Strict_Transport_Security) is included in all API responses,
make sure to specify the FQDN (Fully Qualified Domain Name) of your DTR prefixed with `https://`,
or your browser may refuse to load the web interface.

    ![](../../images/create-repository-1.png){: .with-border}


## Step 5. Configure DTR

After installing DTR, you should configure:

  * The certificates used for TLS communication. [Learn more](../configure/use-your-own-tls-certificates.md).
  * The storage backend to store the Docker images. [Learn more](../configure/external-storage/index.md).

### Web interface

  * To update your TLS certificates, access DTR from the browser and navigate to **System > General**.
  * To configure your storage backend, navigate to **System > Storage**. If you are upgrading and changing your existing storage backend, see [Switch storage backends](/ee/dtr/admin/configure/external-storage/storage-backend-migration/) for recommended steps.

### Command line interface

  To reconfigure DTR using the CLI, see the reference page for [the reconfigure command](/reference/dtr/2.7/cli/reconfigure/).

## Step 6. Test pushing and pulling

Now that you have a working installation of DTR, you should test that you can
push and pull images:

* [Configure your local Docker Engine](../../user/access-dtr/index.md)
* [Create a repository](../../user/manage-images/index.md)
* [Push and pull images](../../user/manage-images/pull-and-push-images.md)

## Step 7. Join replicas to the cluster

This step is optional.

To set up DTR for high availability,
you can add more replicas to your DTR cluster. Adding more replicas allows you
to load-balance requests across all replicas, and keep DTR working if a
replica fails.

For high-availability, you should set 3 or 5 DTR replicas. The replica nodes also need
to be managed by the same UCP.

To add replicas to a DTR cluster, use the [join](/reference/dtr/2.7/cli/join/) command:

1. Load your [UCP user bundle](/ee/ucp/user-access/cli/#use-client-certificates).

2.  Run the join command.

    ```bash
    docker run -it --rm \
      {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} join \
      --ucp-node <ucp-node-name> \
      --ucp-insecure-tls
    ```

    > --ucp-node
    >
    > The `<ucp-node-name>` following the `--ucp-node` flag is the target node to
    > install the DTR replica. This is NOT the UCP Manager URL.
    {: .important}

    When you join a replica to a DTR cluster, you need to specify the
    ID of a replica that is already part of the cluster. You can find an
    existing replica ID by going to the **Shared Resources > Stacks** page on UCP.
    
3. Check that all replicas are running.

    In your browser, navigate to UCP's
    web interface. Select **Shared Resources > Stacks**. All replicas should
    be displayed.

    ![](../../images/install-dtr-6.png){: .with-border}

## Where to go next

- [Install DTR offline](install-offline.md)
- [Upgrade DTR](../upgrade.md)
