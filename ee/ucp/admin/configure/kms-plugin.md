---
title: KMS plugin support for UCP
description: Learn about the KMS plugin for UCP.
keywords: ucp, kms, kubernetes, plugin, configuration
---

>{% include enterprise_label_shortform.md %}

Docker Universal Control Plane (UCP) 3.2.5 adds support for a Key Management Service (KMS) plugin to allow access to third-party secrets management solutions, such as Vault. This plugin is used by UCP for access from Kubernetes clusters.

## Deployment

KMS must be deployed before a machine becomes a UCP manager or it may be considered unhealthy. UCP will not health check, clean up, or otherwise manage the KMS plugin.

## Configuration

KMS plugin configuration should be done through UCP. UCP will maintain ownership of the Kubernetes EncryptionConfig file, where the KMS plugin is configured for Kubernetes. UCP does not currently check this file’s contents after deployment. 

UCP adds new configuration options to the cluster configuration table. These options are not exposed through the web UI, but can be configured via the [API](https://docs.docker.com/ee/ucp/admin/configure/ucp-configuration-file/).

The following table shows the configuration options for the KMS plugin. These options are not required.

| Parameter        | Type | Description                              |
|------------------|------|------------------------------------------|
| `kms_enabled`    | bool | Determines if UCP should configure a KMS plugin. |
| `kms_name` | string  | Name of the KMS plugin resource (for example, “vault”). |
| `kms_endpoint`   | string  | Path of the KMS plugin socket. This path must refer to a UNIX socket on the host (for example, “/tmp/socketfile.sock”). UCP will bind mount this file to make it accessible to the API server. |
| `kms_cachesize`  | int  | Number of data encryption keys (DEKs) to be cached in the clear. |

## Where to go next
* [Using a KMS provider for data encryption](https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/)
* [Encrypting Secret Data at Rest](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/)
* [UCP API Documentation](https://docs.docker.com/reference/ucp/3.2/api/)
