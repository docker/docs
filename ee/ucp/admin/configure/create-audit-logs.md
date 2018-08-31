Audit logs can be considered a type of event that are more focused on external user/agent actions and security than understanding state or events of the system itself. They are a security-relevant chronological set of records documenting the sequence of activities that have affected system by individual users, administrators or other components of the system.

Audit Logs capture all HTTP actions (GET, PUT, POST, PATCH, DELETE) to all UCP API, Swarm API and Kubernetes API endpoints that are invoked (except for the ignored list) and sent to Docker Engine via stdout. zCreating audit logs is mainly CLI driven and is an UCP component that integrates with Swarm, K8s, and UCP APIs.

## Procedure

1. Download the UCP Client bundle.
2.  Retrieve JSON for current audit log configuration.
```
export DOCKER_CERT_PATH=~/ucp-bundle-dir/
curl --cert ${DOCKER_CERT_PATH}/cert.pem --key ${DOCKER_CERT_PATH}/key.pem --cacert ${DOCKER_CERT_PATH}/ca.pem -k -X GET https://ucp-domain/api/ucp/config/logging > auditlog.json
```
3. Modify the auditLevel field to metadata or request.
```
vi auditlog.json

 {"logLevel":"INFO","auditLevel":"metadata","supportDumpIncludeAuditLogs":false}
 ```
4. Send the JSON request for the auditlog config with the same API path but with the `PUT` method
```
curl --cert ${DOCKER_CERT_PATH}/cert.pem --key ${DOCKER_CERT_PATH}/key.pem --cacert ${DOCKER_CERT_PATH}/ca.pem -k -H "Content-Type: application/json" -X PUT --data $(cat auditlog.json) https://ucp-domain/api/ucp/config/logging
```

5. Create any workload or RBAC grants in Kube and generate a support dump to check the contents of ucp-controller.log file for audit log entries.
6.
6. Optionally, configure the Docker Engine driver to logstash and collect and query audit logs within ELK stack after deploying ELK: https://success.docker.com/article/elasticsearch-logstash-kibana-logging

### API Endpoints ignored

The following endpoints are ignored since they are API endpoints that are not considered security events and may cause a large amount of entries to be triggered.

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
