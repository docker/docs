---
description: Learn how to troubleshoot your DTR installation.
keywords: docker, registry, monitor, troubleshoot
title: Troubleshoot Docker Trusted Registry
---

High availability in DTR depends on having overlay networking working in UCP.
To manually test that overlay networking is working in UCP run the following
commands on two different UCP machines.

```
docker run -it --rm --net dtr-ol --name overlay-test1 --entrypoint sh docker/dtr
docker run -it --rm --net dtr-ol --name overlay-test2 --entrypoint ping docker/dtr -c 3 overlay-test1
```

You can create new overlay network for this test with `docker network create -d overlay network-name`.
You can also use any images that contain `sh` and `ping` for this test.

If the second command succeeds, overlay networking is working.

## DTR doesn't come up after a Docker restart

This is a known issue with Docker restart policies when DTR is running on the same
machine as a UCP controller. If this happens, you can simply restart the DTR replica
from the UCP UI under "Applications". The best workaround right now is to not run
DTR on the same node as a UCP controller.

## Etcd refuses to start after a Docker restart

If you see the following log message in etcd's logs after a DTR restart it means that
your DTR replicas are on machines that don't have their clocks synchronized. Etcd requires
synchronized clocks to function correctly.

```
2016-04-27 17:56:34.086748 W | rafthttp: the clock difference against peer aa4fdaf4c562342d is too high [8.484795885s > 1s]
```

## Accessing the RethinkDB Admin UI

 > Warning: This command will expose your database to the internet with no authentication. Use with caution.

Run this on the UCP node that has a DTR replica with the given replica id:

```bash
REPLICA_ID=""

docker run \
  --rm \
  --net dtr-ol \
  --name db-proxy \
  -v dtr-ca-$REPLICA_ID:/ca \
  -p 9999:8080 \
  rethinkdb:2.3 \
    rethinkdb \
      proxy \
      --bind all \
      --canonical-address db-proxy \
      --driver-tls-key /ca/rethink/key.pem \
      --driver-tls-cert /ca/rethink/cert.pem \
      --driver-tls-ca /ca/rethink/cert.pem \
      --cluster-tls-key /ca/rethink/key.pem \
      --cluster-tls-cert /ca/rethink/cert.pem \
      --cluster-tls-ca /ca/rethink/cert.pem \
      --join dtr-rethinkdb-$REPLICA_ID.dtr-ol
```

Options to make this more secure:

* Use `-p 127.0.0.1:9999:8080` to expose the admin UI only to localhost
* Use an SSH tunnel in combination with exposing the port only to localhost
* Use a firewall to limit which IPs are allowed to connect
* Use a second proxy with TLS and basic auth to provide secure access over the Internet

## Accessing etcd directly

You can execute etcd commands on a UCP node hosting a DTR replica using etcdctl
via the following docker command:

```
docker run --rm -v dtr-ca-$REPLICA_ID:/ca --net dtr-br -it --entrypoint /etcdctl docker/dtr-etcd:v2.2.4 --endpoint https://dtr-etcd-$REPLICA_ID.dtr-br:2379 --ca-file /ca/etcd/cert.pem --key-file /ca/etcd-client/key.pem --cert-file /ca/etcd-client/cert.pem
```
