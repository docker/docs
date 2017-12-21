---
description: Scale your swarm stack
keywords: ibm cloud, ibm, iaas, tutorial
title: Modify your Docker EE for IBM Cloud swarm infrastructure
---

## Scale workers

Before you begin:

* [Create a cluster](administering-swarms.md#create-swarms).
* Retrieve your IBM Cloud infrastructure [API username and key](https://knowledgelayer.softlayer.com/procedure/retrieve-your-api-key).

Steps:

1. Connect to your Docker EE for IBM Cloud swarm. Navigate to the directory where you [downloaded the UCP credentials](administering-swarms.md#download-client-certificates) and run the script. For example:

    ```bash
    $ cd filepath/to/certificate/repo && source env.sh
    ```

2. Note the name of the cluster that you want to scale:

    ```bash
    $ bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key
    ```

3. Note the manager leader node:

    ```bash
    $ docker node ls
    ```

4. Get the public IP address of the leader node, replacing _my_swarm_ with the swarm you want to scale:

    ```bash
    $ bx d4ic show --swarm-name my_swarm --sl-user user.name.1234567 --sl-api-key api_key
    ```

5. Connect to the leader node using the _leaderIP_ you previously retrieved:

    ```bash
    $ ssh docker@leaderIP -p 56422
    ```

6. Use InfraKit to modify the number of swarm mode cluster resources. For example, the following commands set the target number of worker nodes in the cluster to 8. You can use the same commands to reduce the number of worker node instances.

    ```bash
    $ /var/ibm/d4ic/infrakit.sh local stack/vars change -c cluster_swarm_worker_size=8
    $ /var/ibm/d4ic/infrakit.sh local stack/groups commit-group file:////infrakit_files/defn-wkr-group.json
    ```
