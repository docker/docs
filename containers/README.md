
# Containers

These are the "accessory" containers with which Hub 2.0 is run.

## dnsmasq

dnsmasq is used to fake the `Origin` header in CORS requests. This is
necessary because the browser automatically sends `Origin: localhost`
(users can't modify it) and we need it to be in the `*.docker.com`
space, since staging is set up to handle single dot subdomains.

We've chosen `bagels.docker.com` as the development domain (something
that is unlikely to ever be deployed in production so that we won't
have to change the name in the future).

### prerequisites

```bash
cd $PROJECT
make dns
```

This runs `$PROJECT/containers/configure_system_dns.sh`, which will
add `bagels.docker.com` to your host system's `/etc/resolver/`. This
makes it so that `bagels.docker.com` will resolver to `boot2docker ip`.

### run

```bash
cd $PROJECT/containers/dnsmasq
docker build -t bagelteam/dnsmasq
docker run -itp 53:53/udp bagelteam/dnsmasq
```

## HAProxy

HAProxy is a load balancer used to terminate SSL.

Currently Out-of-Order.

```bash
docker run -itp 80:80 -p 443:433 bagelteam/haproxy
```

HAProxy will load balance `bagels.docker.com` across a single
container (hah), and more importantly, take care of SSL Offloading at
the load balancer. The image has it's own SSL certificates.
