+++
title = "Manually setting up a CA"
description = "Docker Universal Control Plane"
[menu.main]
parent="mn_ucp"
+++

# Manually setting up a CA

A few features of UCP require an external CA (cfssl or equivalent) to sign
certs. By default in developer mode, these services aren't created (maybe we'll
set them up someday, but there's some bootstrapping challenges...)  The
following instructions explain how to set this up manually so you can do
developer mode with CA support.


```bash
sudo mkdir -p /etc/docker/ssl/orca

sudo sh -c 'cat << EOF > /etc/docker/ssl/orca/config.json
{
    "roots": {
        "orca": {
            "key": "/etc/docker/ssl/orca/orca_ca_key.pem",
            "certificate": "/etc/docker/ssl/orca/orca_ca.pem"
        },
        "swarm": {
            "key": "/etc/docker/ssl/orca/swarm_ca_key.pem",
            "certificate": "/etc/docker/ssl/orca/swarm_ca.pem"
        }
    },
    "signing": {
        "default": {
            "expiry": "87600h"
        },
        "profiles": {
            "client": {
                    "usages": [
                            "signing",
                            "key encipherment",
                            "client auth"
                    ],
                    "expiry": "87600h"
            },
            "node": {
                    "usages": [
                            "signing",
                            "key encipherment",
                            "server auth",
                            "client auth"
                    ],
                    "expiry": "87600h"
            },
            "intermediate": {
                    "usages": [
                            "signing",
                            "key encipherment",
                            "cert sign",
                            "crl sign"
                    ],
                    "is_ca": true,
                    "expiry": "87600h"
            }
        }
    }
}
EOF'

sudo sh -c 'cat << EOF > /etc/docker/ssl/orca/orca_ca.json
{
    "key": {
        "algo": "rsa",
        "size": 4096
    },
    "CN": "Orca Root CA"
}
EOF'

sudo sh -c 'cat << EOF > /etc/docker/ssl/orca/swarm_ca.json
{
    "key": {
        "algo": "rsa",
        "size": 4096
    },
    "CN": "Swarm Root CA"
}
EOF'

docker run --rm -v /etc/docker/ssl/orca:/etc/docker/ssl/orca -w /etc/docker/ssl/orca dockerorca/orca-cfssl genkey -initca orca_ca.json | \
    docker run --rm -i -v /etc/docker/ssl/orca:/etc/docker/ssl/orca --entrypoint cfssljson -w /etc/docker/ssl/orca dockerorca/orca-cfssl -bare orca_ca

docker run --rm -v /etc/docker/ssl/orca:/etc/docker/ssl/orca -w /etc/docker/ssl/orca dockerorca/orca-cfssl genkey -initca swarm_ca.json | \
    docker run --rm -i -v /etc/docker/ssl/orca:/etc/docker/ssl/orca --entrypoint cfssljson -w /etc/docker/ssl/orca dockerorca/orca-cfssl -bare swarm_ca


# Just to keep the naming consistent...
sudo mv /etc/docker/ssl/orca/orca_ca-key.pem /etc/docker/ssl/orca/orca_ca_key.pem
sudo mv /etc/docker/ssl/orca/swarm_ca-key.pem /etc/docker/ssl/orca/swarm_ca_key.pem
```

Once you've generated the root cert (above) you can start the servers

```bash
docker run -d \
    -v /etc/docker/ssl/orca/orca_ca.pem:/etc/cfssl/ca.pem:ro \
    -v /etc/docker/ssl/orca/orca_ca_key.pem:/etc/cfssl/ca-key.pem:ro \
    -v /etc/docker/ssl/orca/config.json:/etc/cfssl/config.json:ro \
    --name orca-ca \
    dockerorca/orca-cfssl serve --address 0.0.0.0 -config config.json

docker run -d \
    -v /etc/docker/ssl/orca/swarm_ca.pem:/etc/cfssl/ca.pem:ro \
    -v /etc/docker/ssl/orca/swarm_ca_key.pem:/etc/cfssl/ca-key.pem:ro \
    -v /etc/docker/ssl/orca/config.json:/etc/cfssl/config.json:ro \
    --name orca-swarm-ca \
    dockerorca/orca-cfssl serve --address 0.0.0.0 -config config.json
```

Now you can generate a server cert (you might want to edit this to make the hostname/ip match your config

```bash
sudo sh -c 'cat << EOF > /etc/docker/ssl/orca/server.json
{
    "hosts": [
        "127.0.0.1"
    ],
    "key": {
        "algo": "rsa",
        "size": 4096
    },
    "CN": "My Server"
}
EOF'

docker run --rm --link=orca-swarm-ca:swarm_ca \
    -v /etc/docker/ssl/orca:/etc/docker/ssl/orca \
    -w /etc/docker/ssl/orca \
    --entrypoint=/bin/sh \
    dockerorca/orca-cfssl -c \
    'cfssl gencert -remote $SWARM_CA_PORT_8888_TCP_ADDR -profile=node server.json' | \
    docker run --rm -i -v /etc/docker/ssl/orca:/etc/docker/ssl/orca \
    -w /etc/docker/ssl/orca \
    --entrypoint=cfssljson dockerorca/orca-cfssl \
    -bare swarm_orca_server

docker run --rm --link=orca-ca:orca_ca \
    -v /etc/docker/ssl/orca:/etc/docker/ssl/orca \
    -w /etc/docker/ssl/orca \
    --entrypoint=/bin/sh \
    dockerorca/orca-cfssl -c \
    'cfssl gencert -remote $ORCA_CA_PORT_8888_TCP_ADDR -profile=node server.json' | \
    docker run --rm -i -v /etc/docker/ssl/orca:/etc/docker/ssl/orca \
    -w /etc/docker/ssl/orca \
    --entrypoint=cfssljson dockerorca/orca-cfssl \
    -bare orca_server

# Just to keep the naming consistent...
sudo mv /etc/docker/ssl/orca/swarm_orca_server-key.pem /etc/docker/ssl/orca/swarm_orca_server_key.pem
sudo mv /etc/docker/ssl/orca/orca_server-key.pem /etc/docker/ssl/orca/orca_server_key.pem
```


**Now you can run UCP**


(proxy)
```bash
docker run -d \
    -v /etc/docker/ssl/orca:/etc/docker/ssl/orca:ro \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e SSL_CA=/etc/docker/ssl/orca/swarm_ca.pem  \
    -e SSL_CERT=/etc/docker/ssl/orca/swarm_orca_server.pem  \
    -e SSL_KEY=/etc/docker/ssl/orca/swarm_orca_server_key.pem  \
    --name orca-proxy \
    dockerorca/orca-proxy
```

(rethinkdb)
```bash
docker run -d \
    --name orca-db \
    dockerorca/rethinkdb
```


Then orca itself

```bash
docker run --rm -it \
    -v /etc/docker/ssl/orca:/etc/docker/ssl/orca:ro \
    --link orca-ca:orca_ca \
    --link orca-swarm-ca:swarm_ca \
    --link orca-proxy:proxy \
    --link orca-db:rethinkdb \
    --name orca-controller \
    dockerorca/orca \
    --debug server \
    --docker tcp://proxy:2375 \
    --tls-ca-cert /etc/docker/ssl/orca/swarm_ca.pem \
    --tls-cert /etc/docker/ssl/orca/swarm_orca_server.pem \
    --tls-key /etc/docker/ssl/orca/swarm_orca_server_key.pem \
    --orca-tls-ca-cert /etc/docker/ssl/orca/orca_ca.pem \
    --orca-tls-cert /etc/docker/ssl/orca/orca_server.pem \
    --orca-tls-key /etc/docker/ssl/orca/orca_server_key.pem
```
