---
description: Sharing data between 2 couchdb databases
keywords: docker, example, package installation, networking, couchdb,  data volumes
title: Dockerize a CouchDB service
redirect_from:
  - /engine/examples/couchdb_data_volumes/
---

> **Note**
>
> **If you don't like sudo** then see [*Giving non-root access*](../engine/install/linux-postinstall.md#manage-docker-as-a-non-root-user)

Here's an example of using data volumes to share the same data between
two CouchDB containers. This could be used for hot upgrades, testing
different versions of CouchDB on the same data, etc.

## Create first database

We're marking `/var/lib/couchdb` as a data volume.

```console
$ COUCH1=$(docker run -d -p 5984 -v /var/lib/couchdb shykes/couchdb:2013-05-03)
```

## Add data to the first database

We're assuming your Docker host is reachable at `localhost`. If not,
replace `localhost` with the public IP of your Docker host.

```console
$ HOST=localhost
$ URL="http://$HOST:$(docker port $COUCH1 5984 | grep -o '[1-9][0-9]*$')/_utils/"
$ echo "Navigate to $URL in your browser, and use the couch interface to add data"
```

## Create second database

This time, we're requesting shared access to `$COUCH1`'s volumes.

```console
$ COUCH2=$(docker run -d -p 5984 --volumes-from $COUCH1 shykes/couchdb:2013-05-03)
```

## Browse data on the second database

```console
$ HOST=localhost
$ URL="http://$HOST:$(docker port $COUCH2 5984 | grep -o '[1-9][0-9]*$')/_utils/"
$ echo "Navigate to $URL in your browser. You should see the same data as in the first database"'!'
```

Congratulations, you are now running two Couchdb containers, completely
isolated from each other *except* for their data.
