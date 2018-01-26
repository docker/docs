---
description: Setup & Prerequisites
keywords: ibm cloud, ibm, iaas, tutorial
title: Docker EE for IBM Cloud setup & prerequisites
redirect_from:
---

## Docker Enterprise Edition (EE) for IBM Cloud

Docker EE for IBM Cloud is an unmanaged, native Docker environment within IBM Cloud that runs Docker Enterprise Edition software. Docker EE for IBM Cloud is available on **December 20th 2017 as a closed Beta**.

[Email IBM to request access to the closed beta](mailto:sealbou@us.ibm.com). In the welcome email you receive, you are given the Docker EE installation URL that you use for the beta.

## Prerequisites

To create a swarm cluster in IBM Cloud, you must have certain accounts, credentials, and environments set up.

### Accounts

If you do not have an IBM Cloud account, [register for a Pay As You Go IBM Cloud account](https://console.bluemix.net/registration/).

If you already have an IBM Cloud account, make sure that you can provision infrastructure resources. You might need to [upgrade or link your account](https://console.bluemix.net/docs/account/index.html#accounts).

For a full list of infrastructure permissions, see [What IBM Cloud infrastructure permissions do I need?](faqs.md). In general you need the ability
to provision the following types of resources:

  * File and block storage.
  * Load balancers.
  * SSH keys.
  * Subnet IPs.
  * Virtual server devices.
  * VLANs.

### Credentials

[Add your SSH key to IBM Cloud infrastructure](https://knowledgelayer.softlayer.com/procedure/add-ssh-key), label it, and note the label for use when [administering swarms](administering-swarms.md).

Log in to [IBM Cloud infrastructure](https://control.softlayer.com/), select your user profile, and under the **API Access Information** section retrieve your **API Username** and **Authentication Key**.

### Environment

If you have not already, [create an organization and space](https://console.bluemix.net/docs/admin/orgs_spaces.html#orgsspacesusers) in IBM Cloud. You must be the account owner or administrator to complete this step.

## Install the CLIs

To use Docker EE for IBM Cloud, you need the following CLIs:

* IBM Cloud CLI version.
* Docker for IBM Cloud plug-in.
* Optional: IBM Cloud Container Registry plug-in.

Steps:

1. Install the [IBM Cloud CLI](https://console.bluemix.net/docs/cli/reference/bluemix_cli/get_started.html#getting-started).

2. Log in to the IBM Cloud CLI. Enter your credentials when prompted. If you have a federated ID, use the `--sso` option.

   ```bash
   $ bx login [--sso]
   ```

3. Install the Docker EE for IBM Cloud plug-in. The prefix for running commands is `bx d4ic`.

   ```bash
   $ bx plugin install docker-for-ibm-cloud -r Bluemix
   ```

4. Optional: To manage a private IBM Cloud Container Registry, install the plug-in. The prefix for running commands is `bx cr`.

   ```bash
   $ bx plugin install container-registry -r Bluemix
   ```

5. Verify that the plug-ins have been installed properly:

   ```bash
  $ bx plugin list
   ```

## Set infrastructure environment variables

The Docker EE for IBM Cloud CLI plug-in simplifies your interaction with IBM Cloud infrastructure resources. As such, many `bx d4ic` commands require you to provide your infrastructure account user name and API key credentials.

Instead of including these in each command, you can set your environment variables.

Steps:

1. [Log in to IBM Cloud infrastructure user profile](https://control.bluemix.net/account/user/profile).

2. Under the **API Access Information** section, locate your **API Username** and **Authentication Key**.

3. Retrieve your Docker EE installation URL. For beta, you received this in your welcome email.

4. From the CLI, set the environment variables with your infrastructure credentials and your Docker EE installation URL:

   ```none
   export SOFTLAYER_USERNAME=user.name.1234567
   export SOFTLAYER_API_KEY=my_authentication_key
   export D4IC_DOCKER_EE_URL=my_docker-ee-url
   ```

5. Verify that your environment variables were set.

   ```bash
   $ env | grep SOFTLAYER && env | grep D4IC_DOCKER_EE_URL
   SOFTLAYER_API_KEY=my_authentication_key
   SOFTLAYER_USERNAME=user.name.1234567
   D4IC_DOCKER_EE_URL=my_docker-ee-url
   ```

## What's next?

* [Create a swarm](administering-swarms.md#create-swarms).
* [Access UCP](administering-swarms.md#access-ucp) and the [download client certificate bundle](administering-swarms.md#download-client-certificates).
* [Learn when to use UCP and the CLIs](administering-swarms.md#ucp-and-clis).
* [Configure DTR to use IBM Cloud Object Storage](dtr-ibm-cos.md).
