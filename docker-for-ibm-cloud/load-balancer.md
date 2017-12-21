---
description: Load Balancer
keywords: IBM Cloud load balancer
title: Load balance Docker EE for IBM Cloud clusters
---

Docker Enterprise Edition (EE) for IBM Cloud deploys three load balancers to each cluster so that you can:

* [Access Docker EE Universal Control Plan (UCP) and cluster manager node](#manager-load-balancer).
* [Access Docker Trusted Registry (DTR)](#dtr-load-balancer).
* [Expose services created in the cluster](#service-load-balancer).

The load balancers are preconfigured for you. Do not change the configurations.

## Manager load balancer

The manager load balancer is preconfigured to connect your local Docker client, your cluster, `bx d4ic` commands, and Docker EE UCP.

Ports:

* UCP is listening on load balancer port 443.
* Agent traffic on the manager nodes is on port 56443.

Use the manager load balancer URL to access Docker EE [UCP](/datacenter/ucp/2.2/guides/).

1. Get your cluster's load balancer URL for UCP:

   ```bash
   $ bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key
   ```

2. In your browser, navigate to the URL and log in.

   **Tip**: Your user name is `admin` or the user name that your admin created for you. You got the password when you [created the cluster](administering-swarms.md#create-swarms) or when your admin created your credentials.


## DTR load balancer

The DTR load balancer is used to access DTR and run registry commands such as `docker push` or `docker pull`.

The DTR load balancer exposes the following ports:

* 443 for HTTPS
* 80 for HTTP

Use the load balancer to access [DTR](/datacenter/dtr/2.4/guides/).

1. Get the name of the cluster for which you want to access DTR:

   ```bash
   $ bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key
   ```

2. Get the DTR URL for the cluster:

   ```bash
   $ bx d4ic show --swarm-name my_swarm --sl-user user.name.1234567 --sl-api-key api_key
   ```

3. In your browser, navigate to the URL and log in. UCP and DTR share the same login.

   **Tip**: Your user name is `admin` or the user name that your admin created for you. You got the password when you [created the cluster](administering-swarms.md#create-swarms) or when your admin created your credentials.

## Service load balancer
When you create a service, any ports that are opened with `-p` are automatically exposed through the load balancer. For example:

```bash
$ docker service create --name nginx -p 80:80 nginx
```

This opens up port 80 on the load balancer, and directs any traffic
on that port to your service.

> Note: 10 ports on the service load balancer
>
> Each cluster's service load balancer can have 10 ports opened. If you create new services or update a service to publish it on a port but already used 10 ports, new ports are not added and the service cannot be accessed through the load balancer.
> If you need more than 10 ports, you can explore alternative solutions such as [UCP domain names](https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/use-domain-names-to-access-services/) or [Træfik](https://github.com/containous/traefik). You can also [create another cluster](administering-swarms.md#create-swarms).

### Accessing a service with the service load balancer

1. Connect to your Docker EE for IBM Cloud swarm. Navigate to the directory where you [downloaded the UCP credentials](administering-swarms.md#download-client-certificates) and run the script. For example:

   ```bash
   $ cd filepath/to/certificate/repo && source env.sh
   ```

2. Create the service specifying the port on which you want the service exposed. For example:

   ```bash
   $ docker service create --name go-demo \
     -e DB=go-demo-db \
     --network go-demo \
     --publish 8080:8080 \
     vfarcic/go-demo
   ```

3. Get the name of the cluster, and then use it to get the service load balancer URL:

   ```bash
   $ bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key
   $ bx d4ic show --swarm-name my_swarm --sl-user user.name.1234567 --sl-api-key api_key
   ```

4. To access a service that you have previously exposed on a port, use the `service-lb-url` that you retrieved. For example:

   ```bash
   $ curl https://service-lb-url:8080/demo/hello
   ```

### Services with SSL certificates
Use [IBM Cloud infrastructure SSL Certificates](https://knowledgelayer.softlayer.com/topic/ssl-certificates) to authenticate and encrypt online transactions that are transmitted through your cluster's load balancer.

When you create a certificate for your domain, specify the **Common Name**. When you create the Docker service, include the certificate common name to use the certificate for SSL termination for your service.

Before you begin:

1. Log in to [IBM Cloud infrastructure](https://control.softlayer.com/).

2. [Add or an import an SSL certificate](https://knowledgelayer.softlayer.com/topic/ssl-certificates) to use. In your infrastructure account, you can access the page from **Security** > **SSL** > **Certificates**.

3. Note the certificate **Common Name**.

**Start a service that uses SSL termination**: Start a service that listens on ports `80` and `443`. The service load balancer provides SSL termination on port `443` that uses your SSL certificate's common name, `com.ibm.d4ic.lb.cert=certificate-common-name`, when you create the service.

In the label, append `@HTTPS:port` to list the ports you want to expose.

For example:

```bash
$ docker service create --name name \
...
--label com.ibm.d4ic.lb.cert=certificate-common-name@HTTPS:443
...
```

To specify other or multiple ports, append them as follows:

* Links HTTPS to port 444: `--label com.ibm.d4ic.lb.cert=certificate-common-name@HTTPS:444`
* Links HTTPS to ports 444 and 8080: `--label com.ibm.d4ic.lb.cert=certificate-common-name@HTTPS:444,HTTPS:8080`

**Set a health check path**: By default, the service load balancer sets a health check path to `/`. If the service cannot respond with a 200 message to a `GET` request on the `/` path, then include a health monitor path label when you create the service. For example:

```bash
--label com.ibm.d4ic.healthcheck.path=/demo/hello@443
```

When the route is published, the health check is set to the path that you specify in the label. Choose a path that can respond with a 200 message to a `GET` request.

**Example command**: The following `docker service create` command expands on the example from the [previous section](#access-a-service-with-the-service-load-balancer) to create a demo service that is exposed on a different port than the default and includes a health check path. It is based on the `vfarcic/go-demo` image.

```bash
$ docker service create --name go-demo \
  -e DB=go-demo-db \
  --network go-demo \
  --publish 8080:8080 \
  --replicas 3 \
  --label com.ibm.d4ic.lb.cert=certificate-common-name@HTTPS:8080 \
  --label com.ibm.d4ic.healthcheck.path=/demo/hello@8080 \
  vfarcic/go-demo
```
