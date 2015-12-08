# Key/Value Store

In this release, UCP leverages the [etcd](https://github.com/coreos/etcd/) KV
store.

Under normal circumstances, you should not have to access the KV store
directly.  To mitigate unforeseen problems or change advanced settings,
you may be required by Docker support or your SE to change configuration
values or data in the store.

The following example demonstrates basic `curl` usage assuming you
have set up your environment with the downloaded ucp bundle.

The example below uses the [jq](https://stedolan.github.io/jq/) tool to
pretty print the resulting json.  This can be omitted for raw json output.


```bash
export KV_URL="https://$(echo $DOCKER_HOST | cut -f3 -d/ | cut -f1 -d:):12379"

curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${KV_URL}/v2/keys | jq "."
```


You can browse the UCP keys under `/v2/keys/ucp/` and swarm under
`/v2/keys/swarm` as well as modify by `POST`ing updated values to
workaround problems.  Further documentation for the etcd API is available
at https://github.com/coreos/etcd/blob/master/Documentation/api.md



### Learn about the certificates

The store is configured with mutual TLS to prevent unauthorized access.

All components in the system that require access to the KV store use
client certificates signed by the Swarm Root CA.  As admin account
certificates are also signed by this Swarm Root CA, administrators can
access the KV store using `curl` or other tools, provided the admin's
certificate is used as the client certificate.
