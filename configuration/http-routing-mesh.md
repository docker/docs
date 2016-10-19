# Use the HTTP Routing Mesh

UCP provides the HTTP Routing Mesh, which adds to the networking capabilities of Docker Engine. Docker Engine provides load balancing and service discovery at the transport layer for TCP and UDP connections. UCP's HTTP Routing Mesh allows you to extend service discovery to have name-based virtual hosting for HTTP services.

See the [Docker Engine documentation on overlay networks](https://docs.docker.com/engine/swarm/networking/) for more information on what Docker Engine provides.

This feature is currently experimental.

## Enabling the HTTP Routing Mesh

To enable the HTTP Routing Mesh, go to the **UCP web UI**, navigate to the **Settings** page, and click the **Routing Mesh** tab.

<!-- todo: add screenshot -->

The default port for HTTP services is **80**. You may choose an alternate port on this screen.

Check the checkbox to enable the HTTP Routing Mesh. This will create a service called `ucp-hrm` and a network called `ucp-hrm`.

## Route to a service using the HTTP Routing Mesh

The HTTP Routing Mesh can route to a Docker service that runs a webserver. This service must meet three criteria:

* The service must be connected to the `ucp-hrm` network
* The service must publish one or more ports
* The service must have a `com.docker.ucp.mesh.http` label to specify the ports to route

The syntax for the `com.docker.ucp.mesh.http` label is a list of one or more values separated by commas. Each of these values is in the form of `internal_port=protocol://host`, where

* `internal_port` is the port the service is listening on (and may be omitted if there is only one port published)
* `protocol` is `http`
* `host` is the hostname that should be routed to this service

Examples:

A service based on the image `mywebserver` with a webserver running on port 8080 can be routed to `http://foo.example.com` can be created using the following:

```sh
docker service create -p 80:8080 --network ucp-hrm -l com.docker.ucp.mesh.http=http://foo.example.com mywebserver
```

Next, you will need to route the referenced domains to the HTTP Routing Mesh.

## Routing Domains to the HTTP Routing Mesh

The HTTP Routing Mesh uses the `Host` HTTP header to determine which service should receive a particular HTTP request. This is typically done using DNS and pointing one or more domains to one or more nodes in the UCP cluster. For more information, see the UCP reference architecture.

## Disabling the HTTP Routing Mesh

To disable the HTTP Routing Mesh, first ensure that all services that are using the HTTP Routing Mesh are disconnected from the **ucp-hrm** network.

Next, go to the **UCP web UI**, navigate to the **Settings** page, and click the **Routing Mesh** tab. Uncheck the checkbox to disable the HTTP Routing Mesh.

## Troubleshooting the HTTP Routing Mesh

Check the logs of the `ucp-controller` containers on your UCP controller nodes for logging from the HTTP Routing Mesh.
