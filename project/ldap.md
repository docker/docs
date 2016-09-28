# Internal notes for AD/LDAP integrations

For all the magical credentials, check out
* https://docker.atlassian.net/wiki/pages/viewpage.action?pageId=18286175#LDAP/ADTesting-DTR1.3LDAPTesting

# Initial configuration

* Display the current auth configuration

```
export UCP_URL="https://$(echo $DOCKER_HOST | cut -f3 -d/ )"

curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${UCP_URL}/api/config/auth | jq "."
```

* Configure AD based on our test rig (replace all **XXX** secrets with the real stuff)

```
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    -XPOST -d \
'{"auth_type":2,"ldap_auth":{"recoveryAdminUsername":"dtradmin1","serverURL":"ldap://54.69.226.114","startTLS":true,"tlsSkipVerify":true,"readerDN":"ad\\demo","readerPassword":"XXX","userBaseDN":"dc=ad,dc=dckr,dc=org","userLoginAttrName":"sAMAccountName","userSearchFilter":"(memberOf=CN=dheusers,OU=Groups,DC=ad,DC=dckr,DC=org)","syncInterval":60}}' \
    ${UCP_URL}/api/config/auth
```


# Adding an AD group

This assumes you've set up your bundle/environment.

* Display the current teams

```
export UCP_URL="https://$(echo $DOCKER_HOST | cut -f3 -d/ )"

curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${UCP_URL}/api/teams | jq "."
```

* Add a team

```
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    -XPOST -d '{"name":"dtradmins","ldapdn":"CN=dtradmins,OU=Groups,DC=ad,DC=dckr,DC=org"}' \
    ${UCP_URL}/api/teams | jq "."
```

* Show details for a specific team (replace the ID from the list output above)

```
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${UCP_URL}/api/teams/2592e87dda67165d | jq "."
```

# Recovery

If you screwed up your LDAP configuration and need to get back to builtin auth, run the following commands on one of the controllers

```
docker exec -it ucp-kv etcdctl \
        --endpoint https://127.0.0.1:2379 \
        --ca-file /etc/docker/ssl/ca.pem \
        --cert-file /etc/docker/ssl/cert.pem \
        --key-file /etc/docker/ssl/key.pem \
        set /orca/v1/config/auth '{"auth_type":0}'
```

After doing this, you should be able to login using your "admin" account.
