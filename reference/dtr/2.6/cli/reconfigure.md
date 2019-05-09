---
title: docker/dtr reconfigure
description: Change DTR configurations
keywords: dtr, cli, reconfigure
---

Change DTR configurations. 

## Usage

```bash
docker run -it --rm docker/dtr \
    reconfigure [command options]
```

## Description


This command changes DTR configuration settings. If you are using NFS as a storage volume, see [Use NFS](/ee/dtr/admin/configure/external-storage/nfs/) for details on changes to the reconfiguration process. 

DTR is restarted for the new configurations to take effect. To have no down
time, configure your DTR for high availability.


## Options

| Option                        | Environment Variable      | Description                                                                          |
|:------------------------------|:--------------------------|:-------------------------------------------------------------------------------------|
| `--debug` | $DEBUG | Enable debug mode for additional logs of this bootstrap container (the log level of downstream DTR containers can be set with `--log-level`). |
| `--dtr-ca` | $DTR_CA | Use a PEM-encoded TLS CA certificate for DTR. By default DTR generates a self-signed TLS certificate during deployment. You can use your own root CA public certificate with `--dtr-ca "$(cat ca.pem)"`. |
| `--dtr-cert` | $DTR_CERT | Use a PEM-encoded TLS certificate for DTR. By default DTR generates a self-signed TLS certificate during deployment. You can use your own public key certificate with `--dtr-cert "$(cat cert.pem)"`. If the certificate has been signed by an intermediate certificate authority, append its public key certificate at the end of the file to establish a chain of trust. |
| `--dtr-external-url` | $DTR_EXTERNAL_URL | URL of the host or load balancer clients use to reach DTR. When you use this flag, users are redirected to UCP for logging in. Once authenticated  they are redirected to the url you specify in this flag. If you don't use this flag, DTR  is deployed without single sign-on with UCP. Users and teams are shared but users login  separately into the two applications. You can enable and disable single sign-on in the DTR  settings. Format `https://host[:port]`, where port is the value you used  with `--replica-https-port`. Since [HSTS (HTTP Strict-Transport-Security) header](https://en.wikipedia.org/wiki/HTTP_Strict_Transport_Security) is included in all API responses, make sure to specify the FQDN (Fully Qualified Domain Name) of your DTR, or your browser may refuse to load the web interface. |
| `--dtr-key` | $DTR_KEY | Use a PEM-encoded TLS private key for DTR. By default DTR generates a self-signed TLS certificate during deployment. You can use your own TLS private key with `--dtr-key "$(cat key.pem)"`. |
| `--dtr-storage-volume` | $DTR_STORAGE_VOLUME | Customize the volume to store Docker images. By default DTR creates a volume to store the Docker images in the local  filesystem of the node where DTR is running, without high-availability. Use this flag to  specify a full path or volume name for DTR to store images. For high-availability, make  sure all DTR replicas can read and write data on this volume. If you're using NFS, use `--nfs-storage-url` instead. |
| `--enable-pprof` | $DTR_PPROF | Enables pprof profiling of the server. Use `--enable-pprof=false` to disable it. Once DTR is deployed with this flag, you can access the pprof endpoint for the api server  at `/debug/pprof`, and the registry endpoint at `/registry_debug_pprof/debug/pprof`. |
| `--existing-replica-id` | $DTR_REPLICA_ID | The ID of an existing DTR replica. To add, remove or modify DTR, you must connect to an existing healthy replica's database. |
| `--help-extended` | $DTR_EXTENDED_HELP | Display extended help text for a given command. |
| `--http-proxy` | $DTR_HTTP_PROXY | The HTTP proxy used for outgoing requests. |
| `--https-proxy` | $DTR_HTTPS_PROXY | The HTTPS proxy used for outgoing requests. |
| `--log-host` | $LOG_HOST | The syslog system to send logs to. The endpoint to send logs to. Use this flag if you set `--log-protocol` to `tcp` or `udp`. |
| `--log-level` | $LOG_LEVEL | Log level for all container logs when logging to syslog. Default: INFO. The supported log levels are `debug`, `info`, `warn`, `error`, or `fatal`. |
| `--log-protocol` | $LOG_PROTOCOL | The protocol for sending logs. Default is internal. By default, DTR internal components log information using the logger specified in the Docker daemon in the node where the DTR replica is deployed.   Use this option to send DTR logs to an external syslog system. The supported values are `tcp`, `udp`, and `internal`. Internal is the default option, stopping DTR from sending logs to an external system. Use this flag with `--log-host`. |
| `--nfs-storage-url` | $NFS_STORAGE_URL | When running DTR 2.5 (with experimental online garbage collection) and 2.6.0-2.6.3, there is an issue with [reconfiguring and restoring DTR with `--nfs-storage-url`](/ee/dtr/release-notes#version-26) which leads to erased tags. Make sure to [back up your DTR metadata](/ee/dtr/admin/disaster-recovery/create-a-backup/#back-up-dtr-metadata) before you proceed. To work around the issue, manually create a storage volume on each DTR node and reconfigure DTR with `--dtr-storage-volume` and your newly-created volume instead. See [Reconfigure Using a Local NFS Volume](https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#reconfigureusingalocalnfsvolume) for more details. To reconfigure DTR to stop using NFS, leave this option empty: `--nfs-storage-url ""`. See [USE NFS](/ee/dtr/admin/configure/external-storage/nfs/) for more details. [Upgrade to 2.6.4](/reference/dtr/2.6/cli/upgrade/) and follow [Best practice for data migration in 2.6.4](/ee/dtr/admin/configure/external-storage/storage-backend-migration/#best-practice-for-data-migration) when switching storage backends. |
| `--async-nfs` | $ASYNC_NFS | Use async NFS volume options on the replica specified in the `--existing-replica-id` option. The NFS configuration must be set with `--nfs-storage-url` explicitly to use this option. Using `--async-nfs` will bring down any containers on the replica that use the NFS volume, delete the NFS volume, bring it back up with the appropriate configuration, and restart any containers that were brought down.  |
| `--nfs-options` | $NFS_OPTIONS | Pass in NFS volume options verbatim for the replica specified in the `--existing-replica-id` option. The NFS configuration must be set with `--nfs-storage-url` explicitly to use this option. Specifying `--nfs-options` will pass in character-for-character the options specified in the argument when creating or recreating the NFS volume. For instance, to use NFS v4 with async, pass in "rw,nfsvers=4,async" as the argument.  |
| `--no-proxy` | $DTR_NO_PROXY | List of domains the proxy should not be used for. When using `--http-proxy` you can use this flag to specify a list  of domains that you don't want to route through the proxy. Format `acme.com[, acme.org]`. |
| `--replica-http-port` | $REPLICA_HTTP_PORT | The public HTTP port for the DTR replica. Default is `80`. This allows you to customize the HTTP port where users can reach DTR. Once users access  the HTTP port, they are redirected to use an HTTPS connection, using the port specified  with --replica-https-port. This port can also be used for unencrypted health checks. |
| `--replica-https-port` | $REPLICA_HTTPS_PORT | The public HTTPS port for the DTR replica. Default is `443`. This allows you to customize the HTTPS port where users can reach DTR. Each replica can  use a different port. |
| `--replica-rethinkdb-cache-mb` | $RETHINKDB_CACHE_MB | The maximum amount of space in MB for RethinkDB in-memory cache used by the given replica. Default is auto. Auto is `(available_memory - 1024) / 2`. This config allows changing the RethinkDB cache usage per replica. You need to run it once per replica to change each one. |
| `--storage-migrated` | $STORAGE_MIGRATED | A flag added in 2.6.4 which lets you indicate the migration status of your storage data. Specify this flag if you are migrating to a new storage backend and have already moved all contents from your old backend to your new one. If not specified, DTR will assume the new backend is empty during a backend storage switch, and consequently destroy your existing tags and related image metadata. |
| `--ucp-ca` | $UCP_CA | Use a PEM-encoded TLS CA certificate for UCP. Download the UCP TLS CA certificate from `https://<ucp-url>/ca`, and  use `--ucp-ca "$(cat ca.pem)"`. |
| `--ucp-insecure-tls` | $UCP_INSECURE_TLS | Disable TLS verification for UCP. The installation uses TLS but always trusts the TLS certificate used by UCP, which can lead to MITM (man-in-the-middle) attacks.  For production deployments, use `--ucp-ca "$(cat ca.pem)"` instead. |
| `--ucp-password` | $UCP_PASSWORD | The UCP administrator password. |
| `--ucp-url` | $UCP_URL | The UCP URL including domain and port. |
| `--ucp-username` | $UCP_USERNAME | The UCP administrator username. |

