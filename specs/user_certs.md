# Generating User certificates by Hand

This shouldn't be required once UCP can generate user certs, but might be useful in some instances...


All of these commands must be run against the engine where UCP is running
E.g. Locally, or via a machine driver that mounts your home directory.
If you can't mount home on the machine, then use docker-machine ssh
... to run these commands from within the host.


```bash
mkdir orca_user
cat << EOF > orca_user/client.json
{
    "hosts": [
        "127.0.0.1"
    ],
    "key": {
        "algo": "rsa",
        "size": 4096
    },
    "CN": "Orca User"
}
EOF

docker run --rm -it \
    --link=orca-ca:orca_ca \
    -v /etc/docker/ssl/orca:/etc/docker/ssl/orca \
    -v $(pwd)/orca_user:/orca_user \
    -w /orca_user \
    --entrypoint=/bin/sh \
    dockerorca/orca-cfssl -c \
    'cfssl gencert -remote $ORCA_CA_PORT_8888_TCP_ADDR -profile=client client.json | cfssljson -bare client '

cp /etc/docker/ssl/orca/orca_ca_chain.pem orca_user/ca.pem
mv orca_user/client.pem orca_user/cert.pem
mv orca_user/client-key.pem orca_user/key.pem
openssl x509 -pubkey -noout -in orca_user/cert.pem > orca_user/pub_key.pem
tar czvf orca_user.tgz orca_user
```

Now get the file off the machine.  If you're using docker machine, something along these lines might work

```bash
NAME=node0
docker-machine ssh ${NAME} "cat orca_user.tgz" | tar zxvf -
```

Now login to UCP, Click on "Security" and paste in the contents of the pub\_key.pem into the Public Key field

Then you should be able to run docker CLI commands with:

```bash
export DOCKER_TLS_VERIFY=1
export DOCKER_HOST=tcp://192.168.122.116:443
export DOCKER_CERT_PATH=$(pwd)/orca_user/
docker info
```
