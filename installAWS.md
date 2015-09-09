
+++
title = "AWS Installation"
description = "Instructions for installing Docker Trusted Registry in Amazon Web Services"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, AWS, Amazon"]
[menu.main]
parent="smn_dhe"
weight=5
+++



# Installing Docker Trusted Registry in Amazon Web Services (AWS)


This Quick Start Guide gives you a hands-on look at how to install and use the Docker Trusted Registry  in an Amazon Web Services Virtual Private Cloud (AWS-VPC) environment. Specifically, this guide demonstrates the process of installing Docker Trusted Registry via an Amazon Machine Image (AMI), performing basic configuration, and then accessing images on the Docker Trusted Registry server from within your AWS VPC.

This guide walks you through the following steps:

1. Launch the Docker Trusted Registry EC2 Host in AWS
2. Configure the AWS components
3. Connect to the Docker Trusted Registry EC2 Host
4. Manage Docker Trusted Registry via the web administration interface
5. Complete a Docker image workflow (push and pull images)

This guide refers to two major components of a Docker Trusted Registry implementation in AWS:

1. The "Docker Trusted Registry EC2 Host". This is the Linux VM running in AWS that hosts the containers required to run Docker Trusted Registry Service.

2. The "Docker Trusted Registry Service". This is the private Docker Registry service that runs on the Docker Trusted Registry EC2 Host.

You should be able to complete this guide in about thirty minutes.

> **Note**: Amazon may occasionally change the appearance of the AWS web
> interface. This mean the AWS web interface may differ from this guide, but the
> overall process remains the same.

## Prerequisites

To complete this guide, you'll need:

* The Docker Hub user-name and password used to obtain the Docker Subscription licenses
* A Docker Trusted Registry license key. Either a purchased license or a trial license will work
* A commercially supported Docker Engine running within AWS
* An AWS account with the ability to launch EC2 instances
* The ability to modify Security Groups and Network ACLs in your AWS VPC
* Familiarity with how to manage resources in an AWS VPC.

## Launching the Docker Trusted Registry EC2 Host in AWS

First, retrieve a copy of the Docker Trusted Registry AMI from the AWS Marketplace. Do this by launching a new EC2 instance from your “EC2 Dashboard” by clicking the blue “Launch Instance” button.

Choose “AWS Marketplace” from the resulting screen, and type "Docker Trusted Registry" into the “Search AWS Marketplace Products” search box.

> **Note**: Currently, the Docker Trusted Registry AMI is only available for Ubuntu 14.04 LTS.

Select the Docker Trusted Registry AMI you wish to retrieve, and then select the instance-type based on your requirements. Then choose the option “Next: Configure Instance Details”.

At this point you must configure the Docker Trusted Registry EC2 Host according to the requirements of your particular environment. When doing so, consider the following:

* If you want your Docker Trusted Registry EC2 Host to be accessible from the internet, you will need to assign it an Elastic IP or a Public IP.
* You may also wish to Tag the Docker Trusted Registry instance with meaningful name.

The Docker Trusted Registry EC2 Host is managed over SSH, whereas the Docker Trusted Registry Service is managed over HTTPS. When launching the AMI for the first time, the wizard will prompt you to create a new “Security Group” with rules that allow SSH, HTTP, and HTTPS already created.

> **Note**: Make sure that you are launching your Docker Trusted Registry EC2 Host in the correct
> Region, VPC, and subnet.

Once you are satisfied with your Docker Trusted Registry EC2 Host's configuration details, click “Launch”.

You will now be prompted to associate the Docker Trusted Registry EC2 Host with a key pair. If you already have a key pair you would like to use, select it from the drop-down list of available key pairs and check the "Acknowledge" check-box. This will enable the “Launch Instances” button.

If you do not have an existing key pair, choose “Create a new key pair” from the first drop-down list, give the key pair a meaningful name, and click the “Download Key Pair” button. This will enable the “Launch Instances” button.

When creating a new key pair, clicking the “Download Key Pair” button initiates a one-time operation that creates the key pair. So make sure you keep the downloaded key pair in a safe place as you will not be able to download it again.

Next, click the “Launch Instances” button.

Your Docker Trusted Registry EC2 Host will launch; you can view its status on the “Instances” page of your “EC2 Dashboard”. It may take a minute or two for your Docker Trusted Registry EC2 Host to reach the running state.

## Configuring AWS Components

Now that you have a Docker Trusted Registry EC2 Host up and running, you'll customize it to integrate with your infrastructure.

Start by configuring your AWS VPC to allow SSH and HTTP/HTTPS traffic to your Docker Trusted Registry EC2 Host.

### Allowing SSH and HTTP/HTTPS access to your Docker Trusted Registry instance

There are two places where you need to enable SSH and HTTP/HTTPS traffic:

1. All Security Groups associated with your Docker Trusted Registry EC2 Host
2. The Network ACL associated with the subnet in which your Docker Trusted Registry EC2 Host is running

#### Security Group configuration

> **Note**: If you configured the Security Group associated with your Docker Trusted Registry EC2
> Host to allow SSH and HTTP/HTTPS traffic when creating the instance, you can
> skip ahead to the next section and configure the Network ACL.

All Security Groups associated with your Docker Trusted Registry instance will need to allow SSH and HTTP/HTTPS traffic.
To ensure this, select your Docker Trusted Registry EC2 Host in your “EC2 dashboard” and click “view rules” from the “Description” tab as shown below. Three rules – allowing TCP ports 22, 80, and 443 – need to be present.

Any rule with a Source of "0.0.0.0/0" will allow any host from any network to connect over that protocol. This works but is not secure. For improved security, you should specify the IP address, or the network, that your management hosts are on.

#### Network ACL configuration

The Network ACL associated with the subnet where your Docker Trusted Registry EC2 Host is running needs to allow inbound SSH and HTTP/HTTPS traffic.

To ensure this, go to your “VPC Dashboard” and select the subnet that your Docker Trusted Registry EC2 Host is running in from the list of available subnets. Then select the “Network ACL” tab. Three rules (allowing TCP ports 22, 80, and 443) need to be present in the “Inbound” section. These rules must appear above the default “DENY” rule.

> **Note**: An ALLOW rule allowing “All Traffic” on “ALL” protocols, on “ALL”
> ports will allow the necessary SSH and HTTP/HTTPS traffic. However, it is more
> secure to create specific rules that only allow specific traffic types.

If you have not given your subnets meaningful names, you may need to obtain the “Subnet ID” in which your Docker Trusted Registry EC2 Host is running. You’ll find it on the “Instance” pane of the your “EC2 Dashboard”. From here you can select your Docker Trusted Registry EC2 Host and obtain its Subnet ID from the “Description” tab. Make a note of the Subnet ID and use it to locate the correct Subnet ID from the “VPC Dashboard”.

You must also make sure that appropriate outbound rules exist in the Network ACL. Commonly, outbound Network ACL rules allow all traffic. However, if your network security policy does not allow this, you will need to create rules that conform to your policy.

## Connecting to the Docker Trusted Registry EC2 Host

Now that you have configured Security Group and Network ACL rules, you can connect to the Docker Trusted Registry EC2 Host over SSH using the key pair associated with the instance and your “ec2-user” username. Beyond this, the Docker Trusted Registry AMI does not require any manual configuration in order to work for this quick start guide, so we won't be discussing further configuration of the Docker Trusted Registry EC2 Host.

When connecting to the Docker Trusted Registry EC2 Host, you will need its DNS name or IP address. This information can be obtained from the “Description” tab of your Docker Trusted Registry EC2 Host in your “EC2 Dashboard”. EC2 instances can have the following IP addresses:

* Private IP (accessible only from within your AWS VPC, as well as from networks connected to your VPC)
* Public IP (accessible from the internet, but will change when the Docker Trusted Registry EC2 Host is rebooted)
* Elastic IP (accessible from the internet and will not change when the Docker Trusted Registry EC2 Host is rebooted)

If you want to manage your Docker Trusted Registry instance from within your AWS VPC, choose the Private DNS or Private IP address.

If you want to manage your Docker Trusted Registry instance over the internet, choose its Public DNS, Elastic IP, or Public IP address.

## Managing the Docker Trusted Registry Service via the Administration web interface

You can now manage the Docker Trusted Registry Service via its Administration web interface over HTTPS. To connect, open a web browser and connect to the DNS name or IP address of your Docker Trusted Registry EC2 Host.

> **Note**: Connecting to the Docker Trusted Registry Service Administration web interface using the default, self-signed certificate will result in a browser warning. This is expected behavior, you can ignore the warning.

Be sure to connect using the correct DNS name or IP address. E.g., if connecting from within AWS, use the Private DNS or Private IP. If connecting from over the internet, use the Public DNS, Public IP, or Elastic IP.

> **Note**: By default, traffic to port 80 and 443 of your Docker Trusted Registry EC2 Host is
> automatically redirected to the Docker Trusted Registry Service Administration web
> interface.

You can perform most Docker Trusted Registry management tasks, including updating Docker Trusted Registry, from the Docker Trusted Registry Administration web interface. But first, two initial tasks must be completed:

1. Configure the Domain Name of your Docker Trusted Registry server
2. License your Docker Trusted Registry server

To configure the Domain Name, click “Settings” > “HTTP”, and enter the DNS name of your Docker Trusted Registry server in the text box titled “Domain Name”. In order to use the Docker Trusted Registry Service to push and pull Docker images from within AWS, you will want to use the AWS Private DNS name.

After configuring the Domain Name, restart Docker Trusted Registry by clicking the “Save and Restart Docker Trusted Registry Server” button.

> **Note**: Changing the Domain Name property of your Docker Trusted Registry server will generate a
> new self-signed certificate that is used by the Docker Trusted Registry Administration web
> interface and the Docker Trusted Registry server. Therefore, you will receive another certificate
> warning the first time you connect to the Docker Trusted Registry Administration web interface
> after changing its Domain Name. This is expected behavior, you can ignore the > warning.

To license your Docker Trusted Registry Service, click “Settings” > “License” and then click “Upload License”. Your license will normally be available for download from your Docker Hub account under “Settings” > “Enterprise Licenses”.

Once your license is uploaded, restart Docker Trusted Registry by clicking the “Save and Restart Docker Trusted Registry Server” button. This completes the basic configuration of Docker Trusted Registry. You can now start using it as an image Registry.

## Docker Image Workflow

This section will walk you through the process of pushing and pulling images to and from your Docker Trusted Registry server from another EC2 instance within your AWS VPC, from a peer VPC, or from a remote location connected via VPN. As such, this guide will use the Private DNS name of the Docker Trusted Registry EC2 Host when tagging and pushing the image.

To complete this section you will need two EC2 instances:

1. The Docker Trusted Registry EC2 Host you have already built and configured
2. A Docker client EC2 instance running commercially supported versions of [Docker Engine](https://www.docker.com/compatibility-maintenance) with at least one image stored locally.

The instructions in this section of the guide will assume the Docker client has a local Docker image called "jenkins", and that the Docker Trusted Registry Service has the following DNS name "ip-10-0-0-117.us-west-2.compute.internal". Your image name and DNS name for your Docker Trusted Registry Service will be different, so you will need to replace these values with the appropriate values for your environment.

> **Note**: Push and pull traffic to a Docker Trusted Registry Service is encrypted using
> SSL certificates. By default, Docker Trusted Registry installs with a self-signed certificate
> which you will need to either: (a) configure your Docker hosts to trust, or
> (b) configure your Docker hosts to ignore by using the `--insecure-registry`
> flag. Alternatively, you can generate and use your own SSL certificates.

### Pushing an image to Docker Trusted Registry Service

From the command line of the Docker client, run the following:

```
docker images
    REPOSITORY  TAG     IMAGE ID        CREATED     VIRTUAL SIZE
    jenkins     latest  4704aa632ce7    12 days ago 887.1 MB

```

> **Note**: Depending on your configuration, you may need to prefix your Docker commands with `sudo`.

You will now tag the local Jenkins image to associate it with a repo in your newly built Docker Trusted Registry server. To do this, type the following:
`docker tag jenkins ip-10-0-0-117.us-west-2.compute.internal/ci-infrastructure/jnkns-img`

This will tag a version of the local Jenkins image so that it can be stored in the "ip-10-0-0-117.us-2.compute.internal" registry in a repository called "ci-infrastructure" with the name "jnkns-img".

Run the `docker images` command again to verify the tag operation succeeded. If it did, you will see an additional tagged image associated with the repository used in the previous docker tag command.

```
docker images
REPOSITORY       TAG      IMAGE ID        CREATED         VIRTUAL SIZE
jenkins                               latest              4704aa632ce7   2 days ago      887.1 MB
ip-10-0-0-117.us-west-2.compute.internal/ci-infrastructure/jnkns-img   latest          4704aa632ce7    2 days ago      887.1 MB

```

Now that the image is tagged, it can be pushed to Docker Trusted Registry with the following command:

```

docker push ip-10-0-0-117.us-west-2.compute.internal/ci-infrastructure/jnkns-img
The push refers to a repository [ip-10-0-0-117.us-west-2.compute.internal/ci-infrastructure/jnkns-img] (len: 1)
4704aa632ce7: Image already exists
77f96086063d: Image successfully pushed
841f40a9f341: Image successfully pushed
8768f04b3a96: Image successfully pushed
fcd8dccdd336: Image successfully pushed
0087c04f8fb6: Image successfully pushed
5cb564bdbf98: Image successfully pushed
<output truncated>
Digest: sha256:1bf8c96ca484290178064e448ea69a55caa52f53ea7e279ff66f5c66625aff43

```

From the “System Health" page of the Docker Trusted Registry Administration web interface, you can view stats from your Docker Trusted Registry Service, including network throughput. The image below shows spikes in network throughput (related to the image_storage_1 image store) generated while the image was being pushed.

Your tagged image is now stored in the Docker Trusted Registry.

### Pulling an image from your Docker Trusted Registry Service

Now that your image is stored in your Docker Trusted Registry, you can pull that image from any supported Docker host that has access to the Registry.

From a Docker Host that has access to the Docker Trusted Registry server, run the following to pull the image locally:

```
docker pull  ip-10-0-0-117.us-west-2.compute.internal/ci-infrastructure/jnkns-img
latest: Pulling from ip-10-0-0-117.us-west-2.compute.internal/ci-infrastructure/jnkns-img
64e5325c0d9d: Extracting [=======>   ] 7.864 MB/51.36 MB
bf84c1d84a8f: Download complete
87de57de6955: Download complete
6a974bea7c0d: Download complete
06c293acac6e: Download complete
b8a058108e9e: Download complete
9aa09af53eee: Download complete
a0513c939a75: Download complete
f509350ab0be: Download complete
b0b7b9978dda: Download complete
6a0b67c37920: Downloading [===============>   ] 63.41 MB/199.1 MB
1f80eb0f8128: Download complete
1d1aa175e120: Download complete
<output truncated>
Digest: sha256:1bf8c96ca484290178064e448ea69a55caa52f53ea7e279ff66f5c66625aff43
Status: Downloaded newer image for ip-10-0-0-117.us-west-2.compute.internal/ci-infrastructure/jnkns-img:latest

```

Finally, run `docker images` again to verify that the image has been successfully pulled and stored locally:

```
docker images
REPOSITORY                                                           TAG        IMAGE ID        CREATED       VIRTUAL SIZE
ip-10-0-0-117.us-west-2.compute.internal/ci-infrastructure/jnkns-img         latest        4704aa632ce7         2 days ago          887.1 MB
```

