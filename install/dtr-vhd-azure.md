+++
title = "Install on Microsoft Azure"
description = "Install Trusted Registry in Microsoft Azure (BYOL)"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, Azure, VHD, Microsoft"]
[menu.main]
parent="workw_dtr_install"
+++


# Install Trusted Registry on Microsoft Azure (BYOL)

This page explains how to install Docker Trusted Registry using a virtual hard
drive (VHD) in a Microsoft Azure environment. Azure is a cloud service which
means that you don't need to host the Trusted Registry your own hardware or
network. If you have not already done so, make sure you have first read the
[installation overview](index.md) for Trusted Registry.  

Before installing, you may want to <a href="http://www.docker.com/microsoft" target="_blank">read more information about running Docker with Microsoft</a>.


## Prerequisites

This installation requires that you "bring your own license" (BYOL).  This means
you need to have a [free trial license or buy a license](../license.md) from
Docker to run Trusted Registry on an Azure server. A license is linked to a
Docker Hub account. The account can be a personal account or an account
associated with your organization.  

Additionally, installing requires a Microsoft Azure account with the
ability to launch new instances. These installation instructions do not
require you to modify security groups or networks in Azure. However, if you are installing for production, authority to modify such settings is recommended.

You should be able to complete the installation in under thirty minutes.

> **Note**: Microsoft may occasionally change the appearance of the Azure web
> interface. So, the interface may differ from what you see here but the
> overall process remains the same.


## Launch the Trusted Registry VHD

1. Log into the <a href="https://portal.azure.com/#" target="_blank">Microsoft Azure portal</a>.

    ![Azure portal](../images/azure_portal.png)

2. Choose the + New option.

3. Choose the Marketplace option.

4. Search for `Docker Trusted Registry`.

    ![Azure filter](../images/azure_filter.png)

5. Double click Docker Trusted Registry.

    The system prompts you to review information about the registry.

6. Press Create.

    The system prompts you to enter basic configuration settings.

    ![Azure basics](../images/basic_configuration.png)

    For production, you should always choose to use an SSH public key. This
    example uses a trial version of Azure, so Password authentication is
    sufficient.

7. Press OK on the the default Size, Settings, and Summary pages.

    If you were going into production, the size and storage of an instance would
    depend on the load and configuration you were planning for. For this
    example, the defaults are sufficient.

8. When you reach the Buy page, press Purchase.

    The Docker Trusted Registry is a bring your own license (BYOL) purchase, so
    the cost of the purchase is 0.00 USD. That is because you should get the
    license through Docker. The use of the Azure instance is charged separately.

    After you press Purchase, Microsoft provisions your instance. Currently, the Azure VHD is an Ubuntu 14.04.3 LTS (GNU/Linux 3.16.0-49-generic x86_64) system.

9. After the provisioning completes, copy the IP address of your instance.

    ![Azure basics](../images/azure_ip.png)

10. In a terminal or through PuTTy, connect to your Trusted Registry instance.

    For example, to connect using SSH and a username/password, you'd do the following:

        $ ssh moxiegirl@40.117.88.185
          moxiegirl@40.117.88.185's password:
          Welcome to Ubuntu 14.04.3 LTS (GNU/Linux 3.16.0-49-generic x86_64)
           * Documentation:  https://help.ubuntu.com/
            System information as of Wed Nov 11 00:45:38 UTC 2015
            System load:  0.07               Processes:              287
            Usage of /:   12.1% of 28.80GB   Users logged in:        0
            Memory usage: 4%                 IP address for eth0:    10.1.0.4
            Swap usage:   0%                 IP address for docker0: 172.17.42.1
            Graph this data and manage this system at:
              https://landscape.canonical.com/
            Get cloud support with Ubuntu Advantage Cloud Guest:
              http://www.ubuntu.com/business/services/cloud

          Last login: Wed Nov 11 00:45:38 2015 from docker.static.monkeybrains.net

11. Check that the Trusted Registry containers are running on this host.

        $ sudo docker ps
        sudo docker ps
        CONTAINER ID        IMAGE                                          COMMAND                CREATED             STATUS              PORTS                                      NAMES
        361856c46c1d        docker/trusted-registry-nginx:1.3.3            "nginxWatcher"         7 weeks ago         Up 24 minutes       0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp   docker_trusted_registry_load_balancer    
        01d6c8204b8c        docker/trusted-registry-admin-server:1.3.3     "server"               7 weeks ago         Up 24 minutes       80/tcp                                     docker_trusted_registry_admin_server     
        5033f0a16a09        docker/trusted-registry-log-aggregator:1.3.3   "log-aggregator"       7 weeks ago         Up 24 minutes                                                  docker_trusted_registry_log_aggregator   
        63141333eab3        docker/trusted-registry-garant:1.3.3           "garant /config/gara   7 weeks ago         Up 24 minutes                                                  docker_trusted_registry_auth_server      
        47fb8f13038a        postgres:9.4.1                                 "/docker-entrypoint.   7 weeks ago         Up 24 minutes       5432/tcp                                   docker_trusted_registry_postgres        

12. Enter the `https://<host-ip>/`` your browser's address bar to display the Trusted Registry Administrator interface.

    Your browser warns you that this is an unsafe site, with a self-signed,
    untrusted certificate. At this point, this dialog is normal and expected;
    allow this connection temporarily.

# Set the Trusted Registry domain name

At this point, the Docker Trusted Registry Administrator site should warn that
the Domain Name is not set. While you can use the public IP address that the portal created, you may find it more convenient to create a fully qualified domain name (FQDN). Refer to the <a href="https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-create-fqdn-on-portal/" target="_blank">Microsoft Azure documentation</a>.

1. Select Settings from the global nav bar at the top of the page.

2. Set the Domain Name to the full host-name of your Docker Trusted Registry server.

3. Click the Save and Restart Docker Trusted Registry Server button to generate a new certificate.

    The certificate is used by both the Docker Trusted Registry Administrator
    web interface and the Docker Trusted Registry server.

3. After the server restarts, allow the connection to the untrusted Docker Trusted Registry web admin site.

    You see a warning notification that this instance of Docker Trusted Registry
    is unlicensed. You'll correct this in the next section.

## Apply your license

The Docker Trusted Registry services will not start until you apply your
license. To do that, you'll first download your license from the Docker Hub and
then upload it to your Docker Trusted Registry web admin server. Follow these
steps:

1. If needed, log back into the [Docker Hub](https://hub.docker.com)
   using the username you used when obtaining your license.

2. Under your name, go to Settings to display the Account Settings page.

3. Click the Licenses submenu to display the Licenses page.

    There is a list of available licenses.

4. Click the download button to obtain the license file you want.

5. Go to your Docker Trusted Registry instance in your browser.

6. Click Settings in the global nav bar.

7. Click License in the Settings nav bar.

8. Click the Choose File button and navigate to your license file.

9. Approve the selection to close the dialog and upload your file.

10. Click the Save and restart button.

    Docker Trusted Registry quits and then restarts with the applied the license.

11. Verify the acceptance of the license by confirming that the "Unlicensed
copy" warning is no longer present.

## Secure the Trusted Registry

Securing Docker Trusted Registry is **required**. You will not be able to push
or pull from Docker Trusted Registry until you secure it.

There are several options and methods for securing Docker Trusted Registry. For
more information, see the [configuration
documentation](../configure/configuration.md#security).

## Push and pull images

Now that you have Docker Trusted Registry configured with a "Domain Name" and
have your client Docker daemons configured with the required security settings,
you can test your setup by following the instructions for [Using Docker Trusted
Registry to Push and pull images](../userguide.md).

## Docker Trusted Registry web interface and registry authentication

By default, there is no authentication set on either the Docker Trusted Registry
web admin interface or the Docker Trusted Registry. You can restrict access
using an in-Docker Trusted Registry configured set of users (and passwords), or
you can configure Docker Trusted Registry to use LDAP-based authentication.

See [Docker Trusted Registry Authentication settings](../configure/configuration.md#authentication) for more details.

## See also

* [Upgrade information](upgrade.md) to upgrade either the Docker Trusted Registry or the commercially supported engine.
* [Install the CS Engine](install-csengine.md).
* To configure for your environment, see the [configuration instructions](../configure/configuration.md).
* To use Docker Trusted Registry, see the [User guide](../userguide.md).
* To make administrative changes, see the[Admin guide](../adminguide.md).
* To see previous changes, see [the release notes](../release-notes.md).
