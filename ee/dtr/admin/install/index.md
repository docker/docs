---
title: Install Docker Trusted Registry
description: Learn how to install Docker Trusted Registry for production.
keywords: dtr, registry, install
---

>{% include enterprise_label_shortform.md %}

Docker Trusted Registry (DTR) is a containerized application that runs on a
swarm managed by the Universal Control Plane (UCP). It can be installed
on-premises or on a cloud infrastructure.

DTR installation is comprised of seven distinct steps:

  1. [Validate System Requirements](validate-the-system-requirements)
  2. [Install UCP](install-ucp)
  3. [Install DTR](install-dtr)
  4. [Confirm that DTR is Running](confirm-that-dtr-is-running)
  5. [Configure DTR](configure-dtr)
  6. [Push/Pull Testing](push/pull-testing)
  7. [Join Replicas to the Cluster](join-replicas-top-the-cluster) (optional)

## 1. Validate the System Requirements

Before installing DTR, confirm that your infrastructure meets DTR's [system requirements](./system-requirements).

## 2. Install UCP

To run, DTR requires Docker Enterprise's Docker Universal Control Panel (UCP). If UCP is not yet installed, refer to [install UCP for production](/ee/ucp/admin/install/).

## 3. Install DTR

>**Note**
>
> Prior to installing DTR:
> * When upgrading, upgrade UCP before DTR for each major version. For example if you are upgrading four major versions, upgrade one major version at a time, first UCP, then DTR, and then repeat for the remaining three versions.
> * UCP upgraded to the most recent version before an initial install of DTR.
> * Docker Engine should be updated to the most recent version before installing or updating UCP.

DTR and UCP must not be installed on the same node, due to the potential for resource and port conflicts. Instead, install DTR on worker nodes that will be managed by UCP. Note also that DTR cannot be installed on a standalone Docker Engine.

![](../../images/install-dtr-1.svg)

1. As an admin, access the UCP web interface.

2. Expand your profile in the UI's left navigation pane and select **Admin Settings > Docker Trusted Registry**.

    ![](../../images/install-dtr-2.png){: .with-border}

3. Configure all DTR options. Once this is done, a Docker CLI command will display, which you can use to install DTR. Before running the command, though, take note of the `--dtr-external-url` parameter:

      ```bash
      $ docker run -it --rm \
        {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} install \
        --dtr-external-url <dtr.example.com> \
        --ucp-node <ucp-node-name> \
        --ucp-username admin \
        --ucp-url <ucp-url>
      ```

    To point the parameter to a load balancer that uses HTTP for health probes over port `80` or `443`, temporarily reconfigure the load balancer to use TCP over a known open port. Once DTR is installed, you can configure the load balancer as necessary.

4. Run the DTR install command on any node connected to the UCP cluster on which the Docker Engine is installed. Note that DTR will not be installed on the node on which the install command is run, but on the UCP worker defined by the `--ucp-node` flag. For example, it is possible to SSH into a UCP node and run the DTR install command from there. When running the installation command in interactive TTY or `-it` mode, you will be prompted for any required additional information. [Learn more about installing DTR](/reference/dtr/2.7/cli/install/).

    * To install a specific version of DTR, replace `{{ page.dtr_version }}` with the  desired version in the [installation command](#step-3-install-dtr). (All DTR versions can be found on the [DTR release notes](/ee/dtr/release-notes/) page.)

    * As DTR is deployed with self-signed certificates by default, UCP might not be able to pull images it. To automatically configure UCP to trust DTR, use the `--dtr-external-url <dtr-domain>:<port>` optional flag during installation, or during a reconfiguration.

    * With DTR 2.7, it is possible to [enable browser authentication via client
    certificates](/ee/enable-client-certificate-authentication/) at the time of install, and thus remove the need to enter user credentials.

5. Verify that DTR is installed, either by navigating to `https://<ucp-fqdn>/manage/settings/dtr`, or by navigating in the UCP web UI to **Admin Settings > Docker Trusted Registry**.

    Under the hood, UCP modifies `/etc/docker/certs.d` for each host and adds DTR's CA certificate. UCP can then pull images from DTR because the Docker Engine for each node in the UCP swarm has been configured to trust DTR.

6. Reconfigure the load balancer back to the desired protocol and port.


## 4. Confirm that DTR is Running

1. Using a browser, navigate to the UCP web interface.

2. Select **Shared Resources > Stacks** in the left navigation pane. DTR should be listed as a stack.

3. Verify that DTR is accessible from the browser by entering the DTR IP address or FQDN on the address bar. As an [HSTS (HTTP Strict-Transport-Security)
header](https://en.wikipedia.org/wiki/HTTP_Strict_Transport_Security) is included in all API responses, make sure to specify the FQDN (Fully Qualified Domain Name) of the DTR instance prefixed with `https://` or the web interface may fail to load.

    ![](../../images/create-repository-1.png){: .with-border}


## 5. Configure DTR

Once DTR has been successfully installed, configure the certificates used for TLS communication ([learn more](../configure/use-your-own-tls-certificates.md)) and the storage backend in which to store the Docker images ([learn more](../configure/external-storage/index.md)).

### Web Interface

  * To update the TLS certificates, access DTR from the browser and navigate to **System > General**.
  * To configure the storage backend, navigate to **System > Storage**. If you are upgrading and changing the existing storage backend, refer to [Switch storage backends](/ee/dtr/admin/configure/external-storage/storage-backend-migration/).

### Command Line Interface

  To reconfigure DTR using the CLI, refer to the [the reconfigure command](/reference/dtr/2.7/cli/reconfigure/) reference page.

## 6. Push and Pull Testing

Once the DTR instance is up and running, test whether images can be pushed and pulled.

* [Configure your local Docker Engine](../../user/access-dtr/index.md)
* [Create a repository](../../user/manage-images/index.md)
* [Push and pull images](../../user/manage-images/pull-and-push-images.md)

## 7. Join Replicas to the Cluster (optional)

To set up DTR for high availability, add more replicas to your DTR cluster. Doing this allows for the load-balancing of requests across all replicas, and it will keep DTR working in the event that a replica fails.

For high-availability, set 3 or 5 DTR replicas. in addition, t\he replica nodes must be managed by the same UCP.

To add replicas to a DTR cluster, use the [join](/reference/dtr/2.7/cli/join/) command:

1. Load the [UCP user bundle](/ee/ucp/user-access/cli/#use-client-certificates).

2.  Run the `join` command.

    ```bash
    docker run -it --rm \
      {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} join \
      --ucp-node <ucp-node-name> \
      --ucp-insecure-tls
    ```

    > Important
    >
    > The `<ucp-node-name>` following the `--ucp-node` flag is the target node on which to install the DTR replica. It is NOT the UCP Manager URL.
    {: .important}

    When joining a replica to a DTR cluster, it is necessary to specify the
    ID of a replica that is already part of the cluster. Locate an
    existing replica ID on the **Shared Resources > Stacks** UCP page.

3. Confirm that all replicas are running by navigating to UCP's web interface and selecting **Shared Resources > Stacks**. All replicas should display.

    ![](../../images/install-dtr-6.png){: .with-border}

## Where to go next

- [Install DTR offline](install-offline.md)
- [Upgrade DTR](../upgrade.md)
