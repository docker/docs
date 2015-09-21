# Install Docker Trusted Registry for AWS (Business Day Support)

This article walks you through the process of launching the *Docker Trusted Registry for AWS (Business Day Support)* AMI as an EC2 instance in the Amazon Web Services (AWS) cloud.

## Prerequisites

You can locate, install, and launch the AMI from the Amazon AWS Marketplace or from with the AWS EC2 Console by selecting the AMI from the "Launch Instance" dialog. Both the AWS Marketplace and the AWS EC2 Console require that you have an AWS account to launch the AMI.

If your account is supplied through your company, your company's administrator must have given you permissions to launch EC2 instances. If you receive a permissions error when following these instructions, contact your AWS administrator for help.

# Install procedure

These instructions show you how to locate, install, and launch Docker Trusted Registry (DTR) from Amazon AWS Marketplace. The AWS Marketplace allows you to do a "1-Click Launch" or "Manual Launch".

The *Manual Launch* allows you to launch using the AWS EC2 Console. It allows for fine control of EC2 instance settings such as:

- Instance type
- VPC settings
- Storage
- Instance tags
- Security Group settings

The *1-Click Launch* is quicker, provides default values for most settings, and dynamically updates the Cost Estimator. This install shows you how to do a *1-Click Launch*. The entire process should take about 15 minutes to complete.

## Locate the Docker Trusted Registry AMI

1. If you haven't already done so, open your browser to the <a href="https://aws.amazon.com/marketplace" targer="_blank">Amazon Marketplace</a>.

2. Search the Marketplace for "Docker Trusted Registry for AWS (Business Day Support)".

    > **Note:** Docker Trusted Registry may not be available in your particular Region. Over time, it will be made available in more Regions.

3. Select the "Docker Trusted Registry for AWS (Business Day Support)" AMI from the list of results.

    The Marketplace entry provides details on the product.

4. Press "Continue" to move to the launch step.

    If you are not logged into AWS, the system prompts you to.

5. Enter your AWS login credentials.

    When your login succeeds, the browser displays the "Launch on EC2" page.

6. Make sure that the "1-Click Launch" tab is selected.

## Deploy the 1-Click Launch

You can deploy Docker Trusted Registry to a private or public subnet. A private subnet provides added security but also prevents your Docker Trusted Registry instance from being directly addressable on the internet. If you choose to deploy to a private subnet, you may need to access your Docker Trusted Registry via a Bastion host or a management instance within your VPC.

These instructions launch an EC2 instance into a public subnet with a public IP so that gaining access to it in the "Connecting to the Docker Trusted Registry Administration web interface" section is simplified. 

The following steps walk you through the 1-Click Launch settings:

1. From the "Software Pricing" box, select a "Subscription Term" and an "Applicable Instance Type."

    These two options contribute to the overall cost of running your choice of EC2 instance. The combination of  these two fees make up the running costs of your EC2 instance, and are shown in the "Cost Estimator" box. Make sure you understand these costs before launching your instance.

> **REVIEW NOTE:** I need to put something in here about selecting the right instance type and mention what happens if it's too small and you need to up it to a bigger instance type - e.g. are their options to migrate images in your registry to a newer bigger DTR instance????????

2. Select the version you want to deploy from the list of available versions.

3. Select the Region you want to deploy to from the "Region" dropdown.

4. Select the VPC and Subnet you want to deploy to from the "VPC" and "Subnet" dropdowns.

5. From the Security Group box, select "Create new based on seller settings".

    ![](http://farm6.staticflickr.com/5719/21466010276_1bf996c189_b.jpg)

    This option has security implications. It allows incoming connections to the listed ports from any host or IP address. You should lock this down in line with your existing AWS security policies.

6. Select an existing or add a new key pair using the "Key Pair" box.

    If you choose to use an existing key pair, be sure to choose one that you have access to, as this cannot be changed after the instance is launched.

7. Review your choices and check the values in the Cost Estimator.

    Changing your selected Region and VPC settings can cause your selected EC2 instance type to reset to the default value of "m3.2xlarge".

8. If you are happy with your configuration and estimated charges, click "Launch with 1-Click".  

9. Go to the <a href="https://console.aws.amazon.com/ec2/v2/home">EC2 Dashboard</a> to view your instance.


## Connect to the Docker Trusted Registry Administration web interface

You administer your Docker Trusted Registry server via the Administration web
interface. You can configure your own custom DNS names for your EC2 instance
using CNAME records so forth. Or, you can use the default DNS names provided by
AWS. These instructions use the default DNS name provided by AWS.

The DTR Administration web interface is exposed on port 443 (HTTPS) of
the EC2 instance. To connect to the DTR Administration web
interface:

1. Log into the AWS Console.

2. Go to the EC2 Dashboard.

3. Choose the "Running Instances" option.

4. Select the Docker Trusted Registry EC2 instance.

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

    ![](http://farm6.staticflickr.com/5743/21497230382_e8523ac28c_b.jpg)

## Configure the Docker Trusted Registry Service

When you first login to the DTR Administration web interface you are prompted to configure the "Domain name" on the "General" tab of the "Settings" page. The Domain Name should be a fully qualified domain name that you have configured for your DTR service. Enter your desired domain name and click the "Save and restart" button at the bottom of the page.

After the DTR server restarts, return to the DTR Administration web interface. The browser displays another certificate related browser warning. Changing the Domain Name property of your DTR server generates a new self-signed certificate. Again, this is expected behavior and you can bypass the warning.

Log into the Docker Trusted Registry and change the default password for the "admin" account from the "Auth" tab on the "Settings" page in the DTR Administration web interface. 

Your Docker Trusted Registry server is now ready for use.

## Next Steps

For more information on using DTR, go to the
[User's Guide]({{< relref "userguide.md" >}}).
