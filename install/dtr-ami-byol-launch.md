+++
title = "Install Docker Subscription for AWS (BYOL))"
description = "Install Docker Subscription for AWS (BYOL)"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, AWS, Amazon, AMI"]
[menu.main]
parent="smn_dhe_install"
weight=-1
+++

# Install Docker Subscription for AWS (BYOL)

This article walks you through the process of launching the *Docker Subscription
for AWS (BYOL)* AMI as an EC2 instance in the Amazon Web Services (AWS) cloud.
The Trusted Registry installation includes a single instance of the commercially
supported Docker Engine.

You can install additional instances of the commercially supported Docker Engine
either on premise or through AWS. For more inforamtion, see the [installation
overview](index.md) for Trusted Registry.

The remainder of this document refers to the running instance of this AMI as a “Trusted Registry”.

## Prerequisites

You need the following to complete this guide:

* An AWS account with permissions to launch EC2 Instances.
* A valid Docker subscription license.

Contact your AWS administrator if your AWS account is provided by your company and you do not have permissions to launch EC2 Instances.

If you do not have a valid Docker Subscription license the following options are available:

- Use the *Docker Trusted Registry for AWS (Business Day Support)* AMI and the *Docker Engine for AWS (Business Day Support)* AMI, both of which include the cost of a Docker Subscription.
- Register for a [Free 30 Day Trial](https://hub.docker.com/enterprise/trial/).
- [Contact Docker](https://www.docker.com/contact) to obtain a quote for a Docker Subscription.

# Install procedure

These instructions show you how to locate, install, and launch a Trusted Registry using the *Docker Subscription for AWS (BYOL)* AMI from Amazon's AWS Marketplace.

The AWS Marketplace allows you to do a "1-Click Launch" or "Manual Launch".

The *Manual Launch* allows you to launch using the AWS EC2 Console. It allows for fine control of EC2 instance settings such as:

- Instance type
- VPC settings
- Storage
- Instance tags
- Security Group settings

The *1-Click Launch* is quicker, provides default values for most settings, and dynamically updates the Cost Estimator. This install shows you how to do a *1-Click Launch*. The entire process should take about 15 minutes to complete.

## Locate the Docker Trusted Registry AMI

1. If you haven't already done so, open your browser to the <a href="https://aws.amazon.com/marketplace">Amazon AWS Marketplace</a>.

2. Search the Marketplace for "Docker Subscription for AWS (BYOL)".

3. Select the "Docker Subscription for AWS (BYOL)" AMI from the list of results.

    The Marketplace entry provides details on the product.

4. Press "Continue" to move to the launch step.

    If you are not logged into AWS, the system prompts you to.

5. Enter your AWS login credentials.

    When your login succeeds, the browser displays the "Launch on EC2" page.

6. Make sure that the "1-Click Launch" tab is selected.

## Deploy the 1-Click Launch

You can deploy a Trusted Registry instance to a private or public subnet. A private subnet provides added security but also prevents your Trusted Registry instance from being directly addressable on the internet. If you choose to deploy to a private subnet, you may need to access your Trusted Registry via a Bastion host or a management instance within your VPC.

These instructions launch a Trusted Registry on an EC2 instance in a public subnet with a public IP, so that gaining access to it in the "Connecting to the Docker Trusted Registry Administration web interface" section is simplified.

> **Note:** Deploying a Trusted Registry instance to an AWS Public Subnet will automatically assign it a Public IP and Public DNS. Do not forget that AWS Public IPs and Public DNS names change when an EC2 Instance is rebooted. If you want your Trusted Registry EC2 Instance to be directly accessible over the internet you should assign it an Elastic IP.

The following steps walk you through the 1-Click Launch settings:

1. Select the version you want to deploy from the list of available versions.

2. Select the Region you want to deploy to from the "Region" dropdown.

3. Select the EC2 Instance type

    Be sure to check the "Pricing Details" and "Cost Estimator" boxes when changing EC2 Instance types.

3. Select the VPC and Subnet you want to deploy to from the "VPC" and "Subnet" dropdowns.

4. From the Security Group box, select "Create new based on seller settings".

    ![](assets/aws-dtr-sg-rules.png)

    This option has security implications. It allows incoming connections to the listed ports from any host or IP address. You should lock this down in line with your existing AWS security policies..

5. Select an existing or add a new key pair using the "Key Pair" box.

    If you choose to use an existing key pair, be sure to choose one that you have access to, as this cannot be changed after the instance is launched.

6. Review your choices and check the values in the Cost Estimator.

    Changing your selected Region and VPC settings can cause your selected EC2 Instance type to reset to the default value of "m3.2xlarge".

7. If you are happy with your configuration and estimated charges, click "Launch with 1-Click".  

8. Go to the <a href="https://console.aws.amazon.com/ec2/v2/home">EC2 Dashboard</a> to view your instance.


## Connect to the Docker Trusted Registry Administration web interface

You administer your Trusted Registry server via the Administration web
interface (hereafter referred to as *DTR Administration web interface*).

You can configure your own custom DNS names for your EC2 instance
using CNAME records and so forth. Or, you can use the default DNS names provided by
AWS. These instructions use the default DNS name provided by AWS.

The DTR Administration web interface is exposed on port 443 (HTTPS) of
the EC2 instance. To connect to the DTR Administration web
interface:

1. Log into the AWS Console.

2. Go to the EC2 Dashboard.

3. Choose the "Running Instances" option.

4. Select the Trusted Registry EC2 instance.

5. Select the "Description" tab.

6. Locate the Public DNS or Public IP of the EC2 instance.

7. Copy the Public DNS or Public IP into your browser's address bar and press `return`.

    > **Note:** Connecting to the DTR Administration web
    interface may result in a certificate related browser warning. This is
    expected behavior and you can bypass the warning.

   The interface prompts you for the username and password.

8. Enter "admin" for the username.

9. For the password, use the EC2 Instance ID.

    You'll find the Instance ID on the "Description" tab on the EC2 Dashboard as shown in the image below:

    ![](assets/aws-instance-id.png)

## Configure the Docker Trusted Registry Service

When you first login to the DTR Administration web interface you are prompted to complete two configuration items:

1.  Configure the "Domain name" on the "General" tab of the "Settings" page.

    This should be a fully qualified domain name that you have configured for your Trusted Registry service.

    Enter your desired domain name and click the "Save and restart" button at the bottom of the page.

    After the Trusted Registry server restarts, return to the DTR Administration web interface. The browser displays another certificate related browser warning. Changing the Domain Name property of your Trusted Registry server generates a new self-signed certificate. Again, this is expected behavior and you can bypass the warning.

    Log back in to the DTR Administration web interface.

2. License your copy of Docker Trusted Registry from the "License" tab of the "Settings" page.

    Your Docker Trusted Registry license file is available from Docker Hub. To download it, login to Docker Hub and click your username in the top right corner. Choose "Settings" and select the "Licenses" tab. Click the download button beneath your license.

    ![](assets/dtr-license-download.png)

    From the Docker Trusted Registry Administration web interface, select "Settings" and then "License". Under the "Apply a new license" heading select "Choose File". Select your downloaded license file and click "Save and restart".

> **Note:** Restarting your Trusted Registry from the DTR Administration web interface, or as part of the above procedures, does not restart the EC2 instance. Therefore, the Public IP and Public DNS of the EC2 instance does not change.

Log into the DTR Administration web interface and change the default password for the "admin" account from the "Auth" tab on the "Settings" page.

Your Docker Trusted Registry server is now ready for use.

## Next Steps

For more information on using DTR, go to the
[User's Guide](https://docs.docker.com/docker-trusted-registry/userguide/).
