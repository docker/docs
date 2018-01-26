---
description: Set up IBM Cloud Object Storage for Docker Trusted Registry
keywords: ibm, ibm cloud, registry, dtr, iaas, tutorial
title: Set up IBM Cloud Object Storage for Docker Trusted Registry
---

# Use DTR with Docker EE for IBM Cloud

To use [Docker Trusted Registry (DTR)](/datacenter/dtr/2.4/guides/) with Docker EE for IBM Cloud, DTR stores images on [IBM Cloud Object Storage (COS)](https://ibm-public-cos.github.io/crs-docs/).

> COS for DTR enabled by default
>
> When you [create a cluster](administering-swarms.md#create-swarms), Docker EE for IBM Cloud orders an IBM Cloud Swift API Object Storage account and creates a container named `dtr-container`.
>
> If you used the `--disable-dtr-storage` parameter to prevent the `dtr-container` from being created, then follow the steps on this page to configure DTR to store images using IBM COS. You can order a new COS account, or use an existing one.
>
> Whether you use the default Swift or set up S3, the COS persists after you delete your cluster. To permanently remove the COS, see [Delete IBM Cloud Object Storage](#delete-ibm-cloud-object-storage).

With IBM Cloud Object Storage S3, your DTR files are stored in "buckets", and users have permissions to read, write, and delete files from those buckets. When you integrate DTR with IBM Cloud Object Storage, DTR sends all read and write operations to the COS bucket so that the images are persisted there.

## Configure DTR to use IBM Cloud Object Storage

To use IBM Cloud Object Storage, make sure that you have [an IBM Cloud Pay As You Go or Subscription account](https://console.bluemix.net/docs/pricing/index.html#accounts) that can provision infrastructure resources. You might need to [upgrade or link your account](https://console.bluemix.net/docs/pricing/index.html#accounts).

If you already have an IBM Cloud Object Storage account that you want to use for your swarm, you can skip the create instructions and configure your [S3](#configure-the-s3-bucket-access-and-permissions-in-dtr) or [Swift](#configure-the-regional-swift-api-container-access-in-dtr) Cloud Object Storage accounts for DTR.

## Create an IBM Cloud Object Storage account instance

1. Log in to your [IBM Cloud infrastructure](https://control.softlayer.com/) account.

2. Click **Storage** > **Object Storage** > **Order Object Storage**.

3. Select the storage type, **Cloud Object Storage - S3 API** or **Cloud Object Storage - Standard Regional Swift API**, then click **Continue**.

4. Select the **Master Service Agreement** acknowledge check box, then click **Place Order**.

5. The new object storage account instance is provisioned shortly, and appears in the table. Click the **Short Description** field to identify what the instance is for, such as **my-d4ic-dtr-cos**.

The next steps depend on which type of COS you created:

* [Configure S3 API](#configure-ibm-cloud-object-storage-s3-for-dtr)
* [Configure Regional Swift API](#configure-ibm-cloud-object-storage-regional-swift-api-for-dtr)

## Configure IBM Cloud Object Storage S3 for DTR

Create and configure your COS S3 bucket to use with DTR.

### Create an IBM Cloud Object Storage bucket

Before you begin, [create an IBM Cloud Object Storage account instance](#create-an-ibm-cloud-object-storage-account-instance).

1. From your [IBM Cloud infrastructure](https://control.softlayer.com/) **Storage** > **Object Storage** page, select the storage account instance that you made previously.

2. From **Manage Buckets**, click the **Add Bucket** icon.

3. Configure your bucket.
   * For **Resiliency/Location**, select the region that you want to use, such as **Region - us south**. For higher resiliency, you can select a **Cross Region** option.
   * For **Storage Class**, keep the default of **Standard**.
   * For the **Bucket Name**, give your bucket a name, such as **dtr**. The name must be unique within your IBM Cloud infrastructure COS account; [learn more about naming requirements](https://ibm-public-cos.github.io/crs-docs/storing-and-retrieving-objects#using-buckets).

4. Click **Add**.

Next, [configure the S3 bucket access and permissions in DTR](#configure-the-s3-bucket-access-and-permissions-in-dtr).

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


## Configure IBM Cloud Object Storage Regional Swift API for DTR

Create and configure your COS Regional Swift API container to use with DTR.

### Create an IBM Cloud Object Storage Regional Swift API container

Before you begin, [create an IBM Cloud Object Storage account instance](#create-an-ibm-cloud-object-storage-account-instance). You also can reuse a COS Regional Swift API account that you made with a previously deleted swarm. When you [create a new swarm](administering-swarms.md#create-swarms), include the `--disable-dtr-storage` parameter and then follow these steps to reuse the COS account.

1. From your [IBM Cloud infrastructure](https://control.softlayer.com/) **Storage** > **Object Storage** page, select the storage account instance that you made previously.

2. Select a data center location within the region that you created your swarm in. For example, by default swarms are created in `wdc07`, so select **Washington 1, wdc - wdc01**.

   | Docker for IBM Cloud location | Cloud Object Storage location |
   | --- | --- |
   | Dallas `dal12`, `dal13` | Dallas 5 `dal - dal05` |
   | Frankfurt `fra02` | Frankfurt 2 `fra - fra02` |
   | Hong Kong `hgk02` | Hong Kong 2 `hkg - hkg02` |
   | London `lon04` | London 2 `lon - lon02` |
   | Paris `par01` | Paris 1 `par - par01` |
   | Sydney `syd01`, `syd04` | Sydney 1 `syd - syd01` |
   | Toronto `tor01` | Toronto 1 `tor - tor01` |
   | Washington DC `wdc06`, `wdc07`| Washington 1 `wdc - wdc01` |

3. Click **Add Container**.

4. Name the container `dtr-container`.

5. Click the **View Credentials** link.

6. Copy the following information:

   * **Authentication Endpoint**: Public URL. For example, `https://wdc.objectstorage.softlayer.net/auth/v1.0/`.
   * **Username**: The user name for the COS account. For example, `IBMOS1234567-01:user.name.1234567`.
   * **API Key (Password)**: An alphanumeric password for the COS account.

Next, [configure the Regional Swift API container access and permissions in DTR](#configure-the-regional-swift-api-container-access-in-dtr).

### Configure the Regional Swift API container access in DTR

Once you’ve created a COS Swift API container, you can configure DTR to use it. You use the information from your IBM Cloud Object Storage instance to configure DTR external storage.

Before you begin:

* Retrieve your [IBM Cloud infrastructure user name and API key](https://knowledgelayer.softlayer.com/procedure/retrieve-your-api-key).
* [Create a cluster](administering-swarms.md#create-swarms) and [set up UCP](administering-swarms.md#use-the-universal-control-plane).
* Get your UCP password by running `$ docker logs mycluster_ID`.

Steps:

1. Retrieve the container name, public authentication endpoint, user name, and API key that you copied from the previous step.

2. Retrieve your cluster's DTR URL. If you don't remember your cluster's name, you can use the `bx d4ic list` command:

   ```bash
   $ bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key
   $ bx d4ic show --swarm-name mycluster --sl-user user.name.1234567 --sl-api-key api_key
   ...
   Load Balancers
   ID                                     Name                Address                                          Type
   ...
   ID                                     mycluster-dtr       mycluster-dtr-1234567-wdc07.lb.bluemix.net       dtr
   ```

3. Log in to DTR with your credentials of `admin` and the UCP password that you previously retrieved, or the credentials that your admin assigned you.

   ```none
   https://mycluster-dtr-1234567-wdc07.lb.bluemix.net
   ```

4. From the DTR GUI, navigate to **Settings** > **Storage**.

5. Choose **Cloud** storage type, and then select the **OpenStack Swift** cloud storage provider.

6. Fill in the **Swift Settings** form with the information from your IBM Cloud Object Storage **dtr-container** container:

   * **Authorization URL**: Fill in the public authentical endpoint that you previously retrieved. For example, `https://wdc.objectstorage.softlayer.net/auth/v1.0/`.
   * **Username**: Fill in the COS account user name that you previously retrieved. For example, `IBMOS1234567-01:user.name.1234567`.
   * **Password**: Fill in the COS account API Key (Password) that you previously retrieved.
   * **Container**: Fill in the container's name that you previously made, such as `dtr-container`.

7. From the DTR GUI, click **Save**.

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
   $ bx d4ic show --swarm-name mycluster --sl-user user.name.1234567 --sl-api-key api_key
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

## Delete IBM Cloud Object Storage

1. Log in to your [IBM Cloud infrastructure account](https://control.softlayer.com/).

2. Select **Storage** > **Object**.

3. For the object storage **Account Name** that you want to delete, click the cancel icon.

   > Find your object storage account name
   >
   > You can see the name of the swarm associated with the object storage account in the description field.
   > Alternatively, you can find the object storage account name for your swarm by running `bx d4ic show --swarm-name my_swarm --sl-user user.name.1234567 --sl-api-key api_key`.

4. Choose when you want to cancel the object storage account and click **Continue**.

5. Check the acknowledgement and click **Cancel Object Storage Account**.
