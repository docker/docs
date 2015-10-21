+++
title = "Trusted Registry Installation Overview"
description = "Trusted Registry Installation Overview"
keywords = ["docker, documentation, about, technology, install, enterprise, hub, CS engine, Docker Trusted Registry"]
[menu.main]
parent="smn_dhe"
identifier="smn_dhe_install"
+++

# Trusted Registry Installation Overview
Docker Trusted Registry is an enterprise grade on-premise registry bundled with commercially supported Docker Engines (CS Engine). This provides a base for you to build your Docker workflows.

Depending on your business requirements, there are two paths available for you to install Docker Trusted Registry (Trusted Registry). This document describes those options and prerequisites in order for you to make a decision that is best suited to your needs and provides the install directions for your selected path.

## Install Options

You can install Trusted Registry on premise or through any cloud provider. Installation through a cloud provider's marketplace is only supported on AWS and Microsoft Azure.

## Get a license

Docker requires that you obtain a license to use the Trusted Registry. The installation path you choose affects the licensing methods available to you.

All installation paths support a license which you buy outright from Docker and
apply during the installation process. If you would like, you can get a free
trial license that is good for 30 days. To get a free trial or buy a
license go to the [Subscription page](https://hub-beta.docker.com/enterprise/)
on Docker Hub.

If you are installing on AWS, you can use a license you bought direct using the Docker's Bring Your Own License (BYOL) AMI. You can also choose to pay by subscription installing with the AWS Business Day Support AMI. Under this model, your license is part of your Amazon Web Services (AWS) Business Support subscription.

Buying a license is the only licensing option available with Microsoft Azure.

## Plan your install

The following steps are the general order of how you would obtain and install Docker Trusted Registry. You need to:

1. Get a license (either paid or trial) or managed through the pay as you go option.
2. If installing on premise, install the commercially supported Engine for your production environment.
3. Install Docker Trusted Registry on top of the commercially supported Engine.
4. Apply the license (unless it is a part of your AWS subscription) and finish the install process.
5. Configure your environment.
6. Test your installation.
7. (Optional) Install the CS Engine on another machine or configure existing CS engines. Again, you will have to configure and test.

You will want to configure a new CS Engine or existing Docker Engines to the Trusted Registry. Remember, your support is based on your type of license which is limited to the Trusted Registry and CS engines.

Docker Trusted Registry requires that you use the latest version of the commercially supported Docker Engine. This means that when you upgrade, you will also be upgrading to the latest CS Engine. The instructions are the same, whether you plan to install for development or enterprise purposes. The difference between them is how extensive you choose to configure your environment.

## Next steps

Now that you have planned for your install, see the following documents:

* Get your [license]({{< relref "license.md" >}}).
* Download the [commercially supported Docker Engine]({{< relref "csengineinstall.md" >}}).
* Install by [bringing your own license (BYOL)]({{< relref "installAWS.md" >}}).
* Install through the [pay as you go option]({{< relref "ami-launch.md" >}}).
* Install [manually]({{< relref "install.md" >}}).

## See also

* To configure for your environment, see the
[configuration instructions]({{< relref "configuration.md" >}}).
* To use Docker Trusted Registry, see the [User guide]({{< relref "userguide.md" >}}).
* To make administrative changes, see the [Admin guide]({{< relref "adminguide.md" >}}).
* To see previous changes, see the [release notes]({{< relref "release-notes.md" >}}).
