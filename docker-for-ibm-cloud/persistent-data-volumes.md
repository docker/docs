---
description: Persistent data volumes
keywords: ibm persistent data volumes
title: Save persistent data in file storage volumes
---

Docker EE for IBM Cloud comes with the `d4ic-volume` plug-in preinstalled. With this plug-in, you can set up your cluster to save persistent data in your IBM Cloud infrastructure account's file storage volumes. Learn how to set up data volumes, create swarm services that use volumes, and clean up volumes.

## Set up file storage for persistent data volumes

With Docker EE for IBM Cloud, you can create new or use existing IBM Cloud infrastructure file storage to save persistent data in your cluster.

> Volumes
>
> Volumes are tied at the cluster level, not the IBM Cloud account level. Follow the steps in this document to set up storage volumes for your cluster. Do not use other methods such as UCP.

### Create file storage

Create an IBM Cloud infrastructure file storage volume from Docker EE for IBM Cloud. The volume might take a few minutes to provision. It has default settings of `Endurance` storage type, `2` IOPS, and `20` capacity (measured in GB).

If you want to change the default settings for storage type, IOPS, or capacity, [review file storage provisioning](https://console.bluemix.net/docs/infrastructure/FileStorage/index.html#provisioning) information. Note that for Docker EE for IBM Cloud, the minimum IOPS per GB is 2.

1. Connect to your cluster manager node:

   1. Get your cluster name by running `bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key`.
   2. Get your manager public IP by running `bx d4ic show --swarm-name my_swarm --sl-user user.name.1234567 --sl-api-key api_key`.
   3. Connect to the manager by running `ssh -A docker@managerIP -p 56422`.

2. Create the volume:

    ```bash
    $ docker volume create my_volume \
    --opt request=provision \
    --driver d4ic-volume
    ```

3. **Optional**: If you want to change the default settings of the volume, specify the options:

    ```bash
    $ docker volume create my_volume \
    --opt request=provision \
    --opt type=Performance \
    --opt iops=100 \
    --opt capacity=40 \
    --opt billingType=monthly \
    --driver d4ic-volume
    ```

    > File storage options
    >
    > If the options specified cannot be provisioned in IBM Cloud infrastructure file storage, the volume is not created. Change the values to ones within [the provisioning scope](https://console.bluemix.net/docs/infrastructure/FileStorage/index.html#provisioning) and try again.

4. Verify that the volume is created by inspecting it:

    ```bash
    $ docker volume ls

    DRIVER                  VOLUME NAME
    d4ic-volume: latest     my_volume

    $ docker volume inspect my_volume
    ```

    Example output:

    ```bash
    {% raw %}
    [
      {
          "Driver": "d4ic-volume:latest",
          "Labels": {},
          "Mountpoint": "my_file_storage_volume_mount_point",
          "Name": "my_volume",
          "Options": {
              "request": "provision"
          },
          "Scope": "global",
          "Status": {
              "Settings": {
                  "Capacity": 20,
                  "Datacenter": "wdc07",
                  "ID": 12345678,
                  "Iops": "",
                  "Mountpoint": "my_file_storage_volume_mount_point",
                  "Notes": "docker_volume_name:my_volume;docker_swarm_id:my_swarmID",
                  "Status": "PROVISION_COMPLETED",
                  "StorageType": "ENDURANCE_FILE_STORAGE"
                          }
                      }
        }
      ]
      {% endraw %}
      ```

> File storage provisioning
>
> The `docker volume create` request might take some time to provision the file storage volume in IBM Cloud infrastructure. If the request fails but a new file storage was created in your infrastructure account, a connection error might have disrupted provisioning. Follow the instructions for [using existing IBM Cloud file storage](#use-existing-ibm-cloud-file-storage) in your swarm volume.

Now [create a swarm service](#create-swarm-services-with-persistent-data) to use your persistent data volume.

### Use existing IBM Cloud file storage

Use an existing IBM Cloud infrastructure file storage volume with Docker EE for IBM Cloud. After you configure the volume to be used with Docker EE for IBM Cloud, you can [create swarm services](#create-swarm-services-with-persistent-data) that use the volume and [clean up](#clean-up-volumes-in-your-swarm) the volume.

1. Connect to your Docker EE for IBM Cloud cluster that you want to mount the volume to. Navigate to the directory where you [downloaded the UCP credentials](administering-swarms.md#download-client-certificates) and run the script. For example:

   ```bash
   $ cd filepath/to/certificate/repo && source env.sh
   ```

2. Retrieve the cluster ID:

   ```bash
   {% raw %}
   $ docker info --format={{.Swarm.Cluster.ID}}
   {% endraw %}
   ```

3. From your browser, log in to your [IBM Cloud infrastructure account](https://control.softlayer.com/) and access the file storage volume that you want to use.

4. Under notes, add the `docker_volume_name` field to the first line. Add a unique volume name and swarm ID.

    ```bash
    docker_volume_name:my_volume;docker_swarm_id:my_swarmID
    ```

    > Volume names
    >
    > Don't use the same Docker volume name for multiple file storage volumes within the same cluster!

5. **Optional**: If you have other notes in the file storage volume, add a semicolon after the Docker volume name. Make sure that the Docker volume name is on the first line. Example:

    ```bash
    docker_volume_name:my_volume;docker_swarm_id:my_swarmID;
    other_field:other_notes
    ```

Now [create a swarm service](#create-swarm-services-with-persistent-data) to use your persistent data volume.

## Create swarm services with persistent data

Before you begin creating services or running tasks for Docker EE for IBM Cloud swarms with persistent data, [set up file storage volumes](#set-up-file-storage-for-persistent-data-volumes). Volumes are shared across all instances of the service.

### Create a service with persistent data

You can create a service to schedule tasks across the worker nodes in your swarm.

Before you begin:

- Connect to the cluster manager node.

  1. Get your cluster name by running `bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key`.
  2. Get your manager public IP by running `bx d4ic show --swarm-name my_swarm --sl-user user.name.1234567 --sl-api-key api_key`.
  3. Connect to the manager by running `ssh -A docker@managerIP -p 56422`.

- [Set up file storage for persistent data](#set-up-file-storage-for-persistent-data-volumes).

Create a service that specifies the volume you want to use. The example creates _my_service_ that schedules a task to spawn swarm containers based on the Alpine image, creates 3 repliacs, mounts to _my_volume_, and sets the volume destination (dst) path within each container to the _/dst/directory_.

 ```bash
 $ docker service create --name my_service \
   --mount type=volume,source=my_volume,dst=/dst/directory,volume-driver=d4ic-volume \
   --replicas=3 \
   alpine ping 8.8.8.8
 ```

> Do not provision a volume when you create a Docker service
>
> When you use the `docker service create` command, do not specify the `volume-opt=request=provision` option. Instead, use the `docker volume create` [command](#create-file-storage) to provision new file storage volumes. You run this command only on one manager node for the swarm, and only once per shared file storage volume.

### Run a task with persistent data

You can run a task on a single Docker node that connects to your IBM Cloud infrastructure file storage volume. If you want to connect the volume to multiple containers in your Docker swarm, [create a service](#create-a-service-with-persistent-data) instead.

Before you begin:

- Connect to a node.
- [Set up file storage for persistent data](#set-up-file-storage-for-persistent-data-volumes).

Create a task that specifies the volume you want to use. The example creates a task that spawns an image based on the Busybox image, mounts it to _my_volume_, and creates the volume path within the container to the _/dst/directory_.

  ```bash
  $ docker run -it --volume my_volume:/dst/directory busybox sh
  ```

## Clean up volumes in your swarm

You can remove services with persistent data, delete volumes, or disconnect a IBM Cloud infrastructure file storage volume.

### Remove services with persistent data

You can remove the persistent data volume service from the Docker swarm. Your IBM Cloud infrastructure file storage volume still exists, and can be mounted to other swarms. If you want to use the service without persistent data, remove the service and create it again without mounting the volume.

Example command:

```bash
$ docker service rm my_service
```

### Delete volumes

You can permanently delete your IBM Cloud infrastructure file storage volume. Any data that is stored on the volume is lost when you delete it. Before deleting a volume, ensure that no service is using the volume for persistent data. You can check your services by using the `docker service inspect` [command](/engine/reference/commandline/service_inspect/).

Example command:

```bash
$ docker volume rm my_volume
```

### Disconnect volumes

You can disconnect a particular IBM Cloud infrastructure file storage volume.

1. Log in to your IBM Cloud infrastructure account and access the file storage volume that you want to disconnect.

2. Under notes, delete the entry `docker_volume_name:my_volume;docker_swarm_id:my_swarmID`.
