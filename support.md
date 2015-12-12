+++
title = "UCP Support Dumps"
description = "Docker Universal Control Plane"
[menu.ucp]
weight="2"
+++



# UCP Support Dumps

UCP supports generating support dumps across the entire swarm cluster,
leveraging the dsinfo container developed by Docker Support.

In this version of UCP, support dumps are only exposed via API, but this can be
fairly easily accessed with curl.  The following example shows how to download a
support bundle from your UCP server.

This example leverages curl, which most customers should have, as well
as a handy utility called [jq](https://stedolan.github.io/jq/).  If the
customer doesn't have jq or doesn't want to install it, they can manually
cut-and-paste the token output from the login command below.


```bash
# Replace with your UCP server IP or hostname
UCP=https://192.68.1.2
echo -n "Please enter your admin password"
read -s PASSWORD
TOKEN=$(curl --insecure -s -X POST -d "{\"username\":\"admin\",\"password\":\"${PASSWORD}\"}" "${UCP}/auth/login" | jq -r '.auth_token')

curl --insecure -s -H "X-Access-Token:admin:${TOKEN}" -X POST "${UCP}/api/support" > dump.zip
```

> **Hints**:

>* The UCP server doesn't like extra slashes at the beginning, so if you set your UCP variable with a trailing slash, then you'll get a 301 (redirect)

>* The token often has special characters in it, so if the user cuts and pastes, they may run into problems with the shell interpreting things like $

>* If you want to avoid the --insecure, you'll have to install the UCP server's cert locally.  See below...


## Trusting the UCP server

If you want to trust the UCP servers certificates on the local system, you can
use the following technique.  Note that the paths are somewhat system specific
and may vary from linux distro to distro.

Run the following when pointed at the machine the UCP server is running on
(either locally, or via DOCKER\_HOST).

```bash
sudo bash -c "docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        dockerorca/ucp \
        dump-certs > /usr/local/share/ca-certificates/orca.crt"
sudo update-ca-certificates
```
