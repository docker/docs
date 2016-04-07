<!--[metadata]>
+++
title = "Installation"
description = "Trusted Registry Installation Overview"
keywords = ["docker, documentation, about, technology, install, enterprise, hub, CS engine, Docker Trusted Registry"]
[menu.main]
parent="workw_dtr"
identifier="workw_dtr_install"
weight=30
+++
<![end-metadata]-->

# Trusted Registry installation overview

Docker Trusted Registry is an enterprise-grade on-premises registry bundled with commercially supported Docker Engines (CS Engine). Use Docker Trusted Registry to manage your images, and the commercially supported Docker Engine to create, test, and share your application images. Together, these two Docker products can optimize your continuous integration (CI) and/or software deployment workflows.

Depending on your business requirements, there are two paths available for you to install Docker Trusted Registry (Trusted Registry). This document describes those options and prerequisites in order for you to make a decision that is best suited to your needs and provides the install directions for your selected path.

## Install options

You can install Trusted Registry on premises or through a cloud provider. Currently, Docker supports installation on any cloud provider.

## Get a license

Docker requires that you obtain a license to use the Trusted Registry. The installation path you choose (on premises or in the cloud) can affect the licensing methods available to you (bring your own license or cloud marketplace).

All installation paths  support a license which you buy outright from Docker and
apply during the installation process. If you would like, you can get a free
trial license that is good for 30 days. To get a free trial or buy a
license go to the [Subscription page](https://hub.docker.com/enterprise/)
on Docker Hub.

If you are installing on Microsoft Azure, you have the option of installing using the Virtual Hard Disk (VHD) in the Azure Marketplace. You should use a license you bought direct from Docker in this installation.

If you are installing on AWS, you have the option of installing using Amazon Machine Images (AMI). You can use a license you bought direct using the Docker's Bring Your Own License (BYOL) AMI. You can also choose to pay-as-you-go by installing with the AWS Business Day Support (BDS) AMI. Under the BDS model, your license is part of your Amazon Web Services (AWS) Business Support subscription.

## Plan your install

This section summarizes the process of installing Docker Trusted Registry.

**(Option 1) Install on physical infrastructure or a cloud provider**

  * Obtain a trial or paid license.
  * Install the commercially supported Docker Engine.
  * Install the Trusted Registry.

**(Option 2) Install using AWS AMI**

  * Decide if you are going to bring your own license or use a subscription.
  * Depending on your choice:
    * Obtain a trial or paid license and install the bring your own license (BYOL) AMI.
    * Install the pay-as-you-go business day subscription (BDS) AMI.

**After installing either option**

* Start the DTR admin console.
* If you have a license and not a subscription, install the license.
* Configure your DTR installation.
* Install additional CS engines on other systems.

Remember, your support is based on your type of license. Each license has a single Trusted Registry and one or more CS engines. Your support for CS Engine installations is limited to the number of engines identified by your license.

Docker Trusted Registry requires that you use the latest version of the commercially supported Docker Engine. This means that when you upgrade Trusted Registry, you must also upgrade to the latest CS Engine.
