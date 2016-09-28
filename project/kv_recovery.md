# Recovering from a loss of etcd cluster quorum

If you have multiple failures in your cluster and are unable to recover
enough controllers to restore quorum (1/2+1) you will have to perform
the following steps to rebuild and recover your cluser.


1. Stop **ALL** the ucp-kv containers that are still alive across all remaining controllers

    ```sh
docker stop ucp-kv
```

2. Select a node to perform the recovery on - This is ideally your first node you installed on, or a node where root CA material was replicated.
    * If you have not backed up your root certificate material, use the regen-certs command to regenerate new root CAs
        * **TODO** - make sure this flow works
    * Consider temporarily disabling multi-host networking - need to test this in 1.1 **TODO**
        * `mv /etc/docker/daemon.json /etc/docker/daemon.json.backup`
        * restart the daemon (OS specific)
        * Move the file back after completion of this workflow, and re-bounce the daemon
    * Connect to the engine on that node (bypassing swarm and UCP which are broken without a working KV)
3. Perform a backup of the KV store data in-place:  **The ucp-kv container MUST BE STOPPED FIRST** see step (1) above.

    ```sh
docker run --rm \
        -v ucp-kv:/data \
        --entrypoint ash \
        docker/ucp-etcd:1.0.4 \
        -c "mkdir /data/backup && etcdctl backup --data-dir /data --backup-dir /data/backup"
```

4. Temporarily start etcd using the backup to force creation of a new cluster with only one member:

    ```sh
docker run -d \
        --name ucp-kv-recover \
        -v ucp-kv:/data \
        docker/ucp-etcd:1.0.4 \
        --data-dir /data/backup \
        --force-new-cluster
```

5. Verify this temporary cluster of 1 is healthy

    ```sh
# This should report "healthy"
docker exec -it ucp-kv-recover etcdctl \
        cluster-health


# This should show "/docker" and "/orca"
docker exec -it ucp-kv-recover etcdctl \
        ls /
```

6. Stop the temporary cluster

    ```sh
docker stop ucp-kv-recover
docker rm ucp-kv-recover
```

7. Swap the backup in for the old data

    ```sh
docker run --rm \
        -v ucp-kv:/data \
        --entrypoint ash \
        docker/ucp-etcd:1.0.4 \
        -c "rm -rf /data/old ; mkdir /data/old && mv /data/member /data/old && mv /data/backup/member /data && rmdir /data/backup"
```

8. Restart the ucp-kv **on this node only** (Leave all the other ucp-kv containers on other nodes offline)

    ```sh
    docker start ucp-kv
```

9. Repair the peer URL
    * During the recovery, the peer URL will be reset to localhost and an incorrect port number, and must be changed back to the proper IP:port for this node before replicas can be (re)joined.

    ```sh
docker exec -it ucp-kv etcdctl \
        --endpoint https://127.0.0.1:2379 \
        --ca-file /etc/docker/ssl/ca.pem \
        --cert-file /etc/docker/ssl/cert.pem \
        --key-file /etc/docker/ssl/key.pem \
        member list
```

    ```sh
543711b6260001: name=orca-kv-192.168.122.6 peerURLs=http://localhost:7001,http://localhost:2380 clientURLs=https://192.168.122.6:12379
```
    * Take note the member ID, and IP address listed as the name and clientURLs settings.  Then run the `member update` command, which might look something like the following given the example output above. (replace the ID and IP address, but the port will always be 12380 for UCP)

    ```sh
docker exec -it ucp-kv etcdctl \
        --endpoint https://127.0.0.1:2379 \
        --ca-file /etc/docker/ssl/ca.pem \
        --cert-file /etc/docker/ssl/cert.pem \
        --key-file /etc/docker/ssl/key.pem \
        member update 543711b6260001 https://192.168.122.6:12380
```

    ```sh
Updated member with ID 543711b6260001 in cluster
```

10. Verify you can log in to UCP on this node. (you should use the IP address to bypass any load balancers you may be using)
    * Running a `docker logs -f ucp-kv` in a terminal may be helpful at this point to keep an eye on the health of the etcd cluster while (re)joining nodes
11. Re-join any replica nodes you want to be part of the cluster.
    * Run these commands on the other controller nodes in your cluster (not the controller you just recovered)
    * **YOU MUST UNINSTALL FIRST** to ensure no stale membership information leaks.  Failure to clean up the old installation/state on the replica will result in etcd getting stuck again due to loss of quorum.  If this happens, you can safely restart this procedure at step 1.

    ```sh
# Manually remove the ucp-kv container that you previouslly stopped
docker rm ucp-kv
docker run --rm -it --name ucp -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp:1.0.4 \
        uninstall -i

docker volume ls  | grep "ucp-"
# Remove any volumes
```
    * Re-join

    ```sh
docker run --rm -it --name ucp -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp:1.0.4 \
        join -i --replica [other args...]
```
