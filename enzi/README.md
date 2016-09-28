# eNZi

An Authentication and Authorization Service for Docker Commercial Products

## Goals

- Become a central location for user and organization Accounts
- Sync user accounts with an external LDAP directory service
- Manage organization members and teams
- Sync organization and team members with an external LDAP directory service
- Be an 'OpenID Connect' provider (OAuth 2.0/JSON Object Signing & Encryption)
- Provide user session management
- Provide resource authorization via labels
- Provide resource authorization delegation to Open ID Connect clients
- Provide a rich API for users and Open ID Connect clients to query:
    - identity
    - sessions
    - membership
    - teams
    - label permissions

## Installation (Develompent Environment)

### Create a TLS Keypair

For development, we'll be using the same keypair for all connections. While
this makes setup easy, it essentially makes the key pair a shared secret. In a
production environment, you should use a separate key pair for each server and
client component with a root CA key stored offline which is only used to issue
new certificates.

First, create a volume in Docker which will be used to store the TLS key pair
and use a container to create the TLS key pair and certificate and save them to
this volume:

```bash
docker run -i --rm -v enzi-tls:/tls alpine:3.4 sh << 'EOF'
    set -ex
    apk add --update openssl
    openssl genrsa -out /tls/key.pem 2048
    openssl req -new -x509 -key /tls/key.pem -out /tls/cert.pem -days 3650 -subj '/CN=*.enzi'
    cp /tls/cert.pem /tls/ca.pem
EOF
```

> **NOTE**: this certificate should work for all hosts on the `enzi` network.

### Create a Docker Network

Each node in the cluster will be able to discover and communicate with each
other over the network. This network is called `enzi`.

```bash
docker network create enzi
```

> **NOTE**: The certificate created in the previous step uses the CommonName
> `*.enzi` which will allow servers to authenticate to clients using any name
> on this network.

### Create a RethinkDB Cluster

The following command will create the first node in the RethinkDB cluster
(inline comments should be ignored by `bash`):

```bash
NUM=01
docker run -d \
    -v enzi-db-$NUM-data:/var/data `# create a volume for this node's data, mounted at /var/data` \
    -v enzi-tls:/tls               `# mount the tls key pair at /tls` \
    --net enzi                     `# add this container to the enzi network` \
    --name enzi-db-$NUM            `# the container will have the domain name "enzi-db-$NUM.enzi" on the network` \
    --net-alias rethinkdb          `# the domain name 'rethinkdb.enzi' will refer to the longest-running container with this alias` \
    jlhawn/rethinkdb:2.3.4         `# this is the container image to run. It has the entrypoint '/bin/rethinkdb'` \
        --bind all                 `# bind to the network interface` \
        --no-http-admin            `# the admin web console is insecure and useless` \
        --server-name enzi_db_$NUM \
        --canonical-address enzi-db-$NUM.enzi \
        --directory /var/data/rethinkdb \
        --driver-tls-key /tls/key.pem \
        --driver-tls-cert /tls/cert.pem \
        --driver-tls-ca /tls/ca.pem \
        --cluster-tls-key /tls/key.pem \
        --cluster-tls-cert /tls/cert.pem \
        --cluster-tls-ca /tls/ca.pem
```

> **NOTE**: You can add as many nodes to the RethinkDB cluster as you want! For
> each additional node you want to add, just increment the `NUM` environment
> variable in the first line of the above command and run the command again
> with the extra option `--join enzi-db-01.enzi`. A cluster of 3 RethinkDB
> nodes is required for high availability of the database, but remember to run
> `enzi-setup` (described in the next section) to reconfigure data replication
> after scaling the database cluster.

### Build and Run Setup Tool

Use the following command to build the required container image:

```bash
make image
```

This command creates a container to build the combined `enzi` program located
at `main.go`. It then runs `docker build` to create the `enzi` image.

Now, use the following command to run `enzi sync-db` which will initialize
the database tables. The data replication count is set to the current number of
nodes in the RethinkDB cluster:

```bash
docker run -it --rm \
    --net enzi \
    -v enzi-tls:/tls \
    docker/enzi:dev \
        --db-addr=rethinkdb.enzi \
        sync-db
```

> **NOTE**: You can run this tool whenever you want to reconfigure the
> database tables (e.g., add a new index, change replication count).

Then create the initial admin user using `enzi-setup create-admin`. You will be
prompted for the username and password.

```bash
docker run -it --rm \
    --net enzi \
    -v enzi-tls:/tls \
    docker/enzi:dev \
        --db-addr=rethinkdb.enzi \
        create-admin --interactive
```

### Run API Server

Use the following command to run an API Server in the background:

```bash
NUM=01
docker run -d \
    --name enzi-api-$NUM \
    --net enzi \
    --net-alias api \
    -v enzi-tls:/tls \
    -p 4443:4443 \
    docker/enzi:dev \
        --db-addr=rethinkdb.enzi \
        api \
            --enable-docs
```

> **NOTE**: You can add as many API Servers as you want! However, to access it
> from your web browser you must handle load balancing on your own, either via
> a proxy or DNS. Other parts of the system are able to access any API Server
> through the Docker network alias.

## Run Worker

Use the following command to run a Worker in the background:

```bash
NUM=01
docker run -d \
    --name enzi-worker-$NUM \
    --net enzi \
    --net-alias worker \
    -v enzi-worker-$NUM-data:/work \
    -v enzi-tls:/tls \
    docker/enzi:dev \
        --debug \
        --db-addr=rethinkdb.enzi \
        worker \
            --addr=enzi-worker-$NUM.enzi:4443
```

### View the Interactive Documentation

Congratulations, the system should now be up and running!

To view interactive documentation for the API, connect to an API Server from
your browser over HTTPS on port 4443 and visit `/v0/docs/`.

### Create the Example Service

eNZi comes with a demo OpenID Connect Client application which serves as an
example for how to integrate a service with eNZi as an authentication provider.
Source code for the Example Service can be found in the directory
`example-service` in this repository.

You should first build a container image for this service using the following
command:

```bash
make example-service
```

For this demo, it is required that you set `/etc/hosts` entries for the domains
`api.enzi` and `example-service.enzi` which both point to the IP
address of the Docker Host which is running the eNZi stack. If the IP of your
docker host is `192.168.56.103`, for example, you would add the following line
to your `/etc/hosts` file:

```
192.168.56.103 api.enzi example-service.enzi
```

Now, visit `https://api.enzi:4443/login` in your browser and login using
the admin username and password which you setup during installation earlier in
this guide.

Once logged in, you will be redirected to the interactive API documentation at
`/v0/docs/`. A session token has been set in your browser which gives you full
admin access when trying any of the API endpoints.

To install and run the Example Service, we have to register the Example Service
with eNZi. You can do this with the interactive API Documentation and the
following commands.

First, find the documentation for the `Create a Service` endpoint:

```
POST /v0/accounts/{accountNameOrID}/services
```

Click the header to expand the documentation and form for this API endpoint.
Fill out the `accountNameOrID` field using the username you used to login. The
body will be filled in with a JSON payload that we will prepare with the
following commands.

We must create a JSON escaped string of our `ca.pem` file which was generated
in an earlier step:

```bash
CA_BUNDLE=$(docker run -i --rm -v enzi-tls:/tls alpine:3.3 sh << 'EOF'
head -c -1 /tls/ca.pem | sed ':begin;$!N;s|\n|\\n|;tbegin'
EOF
)
```

Next, prepare a JSON payload which we will use with the API endpoint to
register a new service:


```bash
JSON_PAYLOAD='{
    "name": "Example Service",
    "description": "This is a simple example service.",
    "url": "https://example-service.enzi:4043",
    "privileged": true,
    "redirectURIs": [
        "https://example-service.enzi:4043/openid_callback"
    ],
    "jwksURIs": [
        "https://example-service.enzi:4043/openid_keys"
    ],
    "providerIdentities": [
        "api.enzi:4443"
    ],
    "caBundle": "'$CA_BUNDLE'"
}'
```

Read the descriptions for the fields on the right side of the body form if you
would like to know more about these options.

Now, print out the full payload:

```bash
echo "$JSON_PAYLOAD"
```

Copy and paste this JSON object into the `body` form field in your browser. Now
click the `Try it out!` button at the bottom of the form.

The form will now expand with details about the request and response. The
expected response code is 201. The response body will be a JSON document like
the one you submitted but will contain an `ownerID` field (the ID of your
admin user) and an `id` field which is the ID of the service. Copy this ID as
you will need it in a moment.

Finally, we can run the Example Service using the following command (if your
service ID was `acfcb3eb-1cc6-4728-be63-2433ee16101f`):

```bash
SERVICE_ID=acfcb3eb-1cc6-4728-be63-2433ee16101f
docker run -d \
    --name enzi-example-service \
    --net enzi \
    --net-alias example-service.enzi \
    -v enzi-tls:/tls \
    -p 4043:4043 \
    docker/enzi-example-service:dev \
        --provider-host api.enzi:4443 \
        --service-id $SERVICE_ID \
        --service-host example-service.enzi:4043
```

Now visit `https://example-service.enzi:4043` in your browser. You should see
a simple `Hello` message addressed to your username. The example service was
able to get your username via an OpenID connect authorization flow. It now has
a single-sign-on session with `https://api.enzi:4443`. If you visit
`https://api.enzi:4443/logout`, logout, and then visit the Example
Service page again, it will notice that you have logged out and will redirect
you to login. Go ahead and try it!
