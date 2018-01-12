---
description: Set up IBM Cloud Object Storage for Docker Trusted Registry
keywords: ibm, ibm cloud, registry, dtr, iaas, tutorial
title: Set up IBM Cloud Object Storage for Docker Trusted Registry
---

# Use DTR with Docker EE for IBM Cloud
To use [Docker Trusted Registry (DTR)](docs.docker.com/datacenter/dtr/2.4/guides/) with Docker EE for IBM Cloud, configure DTR to store images on [IBM Cloud Object Storage (COS)](https://ibm-public-cos.github.io/crs-docs/).

Your DTR files are stored in "buckets", and users have permissions to read, write, and delete files from those buckets. When you integrate DTR with IBM Cloud Object Storage, DTR sends all read and write operations to the COS bucket so that the images are persisted there.

## Configure DTR to use IBM Cloud Object Storage
To use IBM Cloud Object Storage, make sure that you have [an IBM Cloud Pay As You Go or Subscription account](https://console.bluemix.net/docs/pricing/index.html#accounts) that can provision infrastructure resources. You might need to [upgrade or link your account](https://console.bluemix.net/docs/pricing/index.html#accounts).

### Create an IBM Cloud Object Storage account instance

1. Log in to your [IBM Cloud infrastructure](https://control.softlayer.com/) account.
2. Click **Storage** > **Object Storage** > **Order Object Storage**.
3. Select the default storage type of **Cloud Object Storage - S3 API**, then click **Continue**.
4. Select the **Master Service Agreement** acknowledge check box, then click **Place Order**.
5. The new object storage account instance is provisioned shortly, and appears in the table. Click the **Short Description** field to identify what the instance is for, such as **my-d4ic-dtr-cos**.

### Create an IBM Cloud Object Storage bucket
Before you begin, [create an IBM Cloud Object Storage account instance](#create-an-ibm-cloud-object-storage-account-instance).

1. From your [IBM Cloud infrastructure](https://control.softlayer.com/) **Storage** > **Object Storage** page, select the storage account instance that you made previously.
2. From **Manage Buckets**, click the **Add Bucket** icon.
3. Configure your bucket.
   * For **Resiliency/Location**, select the region that you want to use, such as **Region - us south**. For higher resiliency, you can select a **Cross Region** option.
   * For **Storage Class**, keep the default of **Standard**.
   * For the **Bucket Name**, give your bucket a name, such as **dtr**. The name must be unique within your IBM Cloud infrastructure COS account; [learn more about naming requirements](https://ibm-public-cos.github.io/crs-docs/storing-and-retrieving-objects#using-buckets).
4. Click **Add**.

### Configure the S3 bucket access and permissions in DTR
Once you’ve created a bucket, you can configure DTR to use it. You use the information from your IBM Cloud Object Storage instance to configure DTR external storage.

Before you begin:

* Retrieve your [IBM Cloud infrastructure user name and API key](https://knowledgelayer.softlayer.com/procedure/retrieve-your-api-key).
* [Create a cluster](administering-swarms.md#create-swarms) and [set up UCP](administering-swarms.md#use-the-universal-control-plane).
* Get your UCP password by running `$ docker logs mycluster_ID`.

Steps:

1. From your [IBM Cloud infrastructure](https://control.softlayer.com/) **Storage** > **Object Storage** page, select the storage account instance that you made previously.
2. Select **Access & Permissions**. Keep this page handy, as you use its information to fill in the DTR external storage required fields later in these steps.
3. Retrieve your cluster's DTR URL. If you don't remember your cluster's name, you can use the `bx d4ic list` command:

   ```bash
   $ bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key
   $ bx d4ic show --swarm-name mycluster --sl-user user.name.1234567 --sl-api-key api_key
   ...
   Load Balancers
   ID                                     Name                Address                                          Type
   ...
   ID                                     mycluster-dtr       mycluster-dtr-1234567-wdc07.lb.bluemix.net       dtr
   ```

4. Log in to DTR with your credentials of `admin` and the UCP password that you previously retrieved, or the credentials that your admin assigned you.

   ```none
   https://mycluster-dtr-1234567-wdc07.lb.bluemix.net
   ```

5. From the DTR GUI, navigate to **Settings** > **Storage**.
6. Choose **Cloud** storage type, and then select the **Amazon S3** cloud storage provider.
7. Fill in the **S3 Settings** form with the information from your IBM Cloud Object Storage **dtr** bucket. The following table describes the required fields and what information to include:

   | Field | Description of what to include |
   | --- | --- |
   | Root directory | The path in the bucket where images are stored. You can leave this blank. |
   | AWS Region name | The AWS region does not affect IBM COS, so you can fill it in with any region, such as us-east-2. |
   | S3 bucket name  | The name of the IBM COS bucket that you previously made, **dtr**. |
   | AWS access key  | From **Access & Permissions** in your IBM COS bucket, expand **Access Keys** and include your **Access Key ID**. |
   | AWS secret key  | From **Access & Permissions** in your IBM COS bucket, expand **Access Keys** and include your **Secret Access Key**. |
   | Region endpoint | From **Access & Permissions** in your IBM COS bucket, expand **Region Endpoints** and include the **Public** endpoint for the region that you created your COS bucket in, such as **us-south**: `s3.us-south.objectstorage.softlayer.net`. |

8. From the DTR GUI, click **Save**.

## Configure your client
Docker EE for IBM Cloud uses a TLS certificate in its storage backend, so you must configure your Docker Engine to trust DTR.

Before you begin:

* Retrieve your [IBM Cloud infrastructure user name and API key](https://knowledgelayer.softlayer.com/procedure/retrieve-your-api-key).
* [Create a cluster](administering-swarms.md#create-swarms) and [set up UCP](administering-swarms.md#use-the-universal-control-plane).
* Get your UCP password by running `$ docker logs mycluster_ID`.
* [Configure DTR to use IBM Cloud Object Storage](#configure-dtr-to-use-ibm-cloud-object-storage).

1. List your swarm cluster's name and then retrieve its DTR URL:

   ```bash
   $ bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key
   $ bx d4ic show --swarm-name mycluster
   ...
   Load Balancers
   ID                                     Name                Address                                          Type
   ...
   ID                                     mycluster-dtr       mycluster-dtr-1234567-wdc07.lb.bluemix.net       dtr
   ```

2. Navigate to the DTR page for certificate authority by appending to the DTR URL: `/ca`. Log in with your credentials of `admin` and the UCP password that you previously retrieved, or the credentials that your admin assigned you.

   ```none
   https://mycluster-dtr-1234567-wdc07.lb.bluemix.net/ca
   ```

3. Follow the OS-specific instructions to [configure your host](/datacenter/dtr/2.4/guides/user/access-dtr/) to use the DTR CA.

## Verify that DTR is running
Verify that DTR is set up properly by [pulling and pushing an image](/datacenter/dtr/2.4/guides/user/manage-images/pull-and-push-images/).
