---
title: Create UCP audit logs
description: Learn how to create audit logs of all activity in UCP
keywords: logs, ucp, swarm, kubernetes, audits
---

Audit logs are a chronological record of security-relevant activities by individual users, administrators or software components that have affected the system. They are focused on external user/agent actions and security rather than understanding state or events of the system itself.

Audit Logs capture all HTTP actions (GET, PUT, POST, PATCH, DELETE) to all UCP API, Swarm API and Kubernetes API endpoints that are invoked (except for the ignored list) and sent to Docker Engine via stdout. Creating audit logs is a UCP component that integrates with Swarm, K8s, and UCP APIs.

## Logging levels

To allow more control to administrators over the audit Logging, three audit logging levels are provided:

- None: audit logging is disabled
- Metadata: includes the following:
    - Method and API endpoint for the request
    - UCP user who made the request
    - Response Status (success or failure)
    - Timestamp of the call
    - Object ID of any created or updated resource (for create or update API calls). We do not include names of created or updated resources
    - License Key
    - Remote Address
- Request: includes all fields from the Metadata level as well as the request payload.

## Benefits

You can use audit logs to help with the following use cases:

- **Historical Troubleshooting** - Audit logs are helpful in determining a sequence of past events that explain why an issue occured.
- **Security Analysis and Auditing** - Security is one of the primary uses for audit logs. A full record of all user interactions with the container infrastructure gives your security team full visibility to questionable or attempted unauthorized accesses.
- **Chargeback** - You can use audit logs and information about the resources to generate chargeback information.
- **Alerting** - If there is a watch on an event stream or a notification created by the event, alerting features can be built on top of event tools that generate alerts for ops teams (PagerDuty, OpsGenie, Slack, or custom solutions).

## Procedure

1. Download the UCP Client bundle [Download client bundle from the command line](https://success.docker.com/article/download-client-bundle-from-the-cli).

2.  Retrieve JSON for current audit log configuration.

    ```
    export DOCKER_CERT_PATH=~/ucp-bundle-dir/
    curl --cert ${DOCKER_CERT_PATH}/cert.pem --key ${DOCKER_CERT_PATH}/key.pem --cacert ${DOCKER_CERT_PATH}/ca.pem -k -X GET https://ucp-domain/api/ucp/config/logging > auditlog.json
    ```

3. Open auditlog.json to modify the 'auditlevel' field to `metadata` or `request`.

    ```
    {
    	"logLevel": "INFO",
    	"auditLevel": "metadata",
    	"supportDumpIncludeAuditLogs": false
    }
    ```

4. Send the JSON request for the auditlog config with the same API path but with the `PUT` method.

    ```
    curl --cert ${DOCKER_CERT_PATH}/cert.pem --key ${DOCKER_CERT_PATH}/key.pem --cacert ${DOCKER_CERT_PATH}/ca.pem -k -H "Content-Type: application/json" -X PUT --data $(cat auditlog.json) https://ucp-domain/api/ucp/config/logging
    ```

5. Create any workload or RBAC grants in Kubernetes and generate a support dump to check the contents of ucp-controller.log file for audit log entries.

6. Optionally, configure the Docker Engine driver to logstash and collect and query audit logs within ELK stack after deploying ELK. (https://success.docker.com/article/elasticsearch-logstash-kibana-logging)

## Sample logs

Here is a sample audit log for a Kubernetes cluster.

```
{"audit"; {
      "metadata": {...},
      "level": "Metadata",
      "timestamp": "2018-08-07T22:10:35Z",
      "auditID": "7559d301-fa6b-4ad6-901c-b587fab75277",
      "stage": "RequestReceived",
      "requestURI": "/api/v1/namespaces/default/pods",
      "verb": "list",
      "user": {"username": "alice",...},
      "sourceIPs": ["127.0.0.1"],
      ...,
      "requestReceivedTimestamp": "2018-08-07T22:10:35.428850Z"}}
```
Here is a sample audit log for a Swarm cluster.

```
{"audit"; {
      "metadata": {...},
      "level": "Metadata",
      "timestamp": "2018-08-07T22:10:35Z",
      "auditID": "7559d301-94e7-4ad6-901c-b587fab31512",
      "stage": "RequestReceived",
      "requestURI": "/v1.30/configs/create",
      "verb": "post",
      "user": {"username": "alice",...},
      "sourceIPs": ["127.0.0.1"],
      ...,
      "requestReceivedTimestamp": "2018-08-07T22:10:35.428850Z"}}
```

## API endpoints ignored

The following API endpoints are ignored since they are not considered security events and may create a large amount of log entries.

- /_ping
- /ca
- /auth
- /trustedregistryca
- /kubeauth
- /metrics
- /info
- /version*
- /debug
- /openid_keys
- /apidocs
- /kubernetesdocs
- /manage
