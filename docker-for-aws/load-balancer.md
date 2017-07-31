---
description: Load Balancer
keywords: aws load balancer elb
title: Configure the Docker for AWS load balancer
---

{% include d4a_buttons.md %}

## How does it work?

When you create a service, any ports that are exposed with `-p` are automatically exposed through the platform load balancer:

```bash
$ docker service create --name nginx -p 80:80 nginx
```

This opens up port 80 on the Elastic Load Balancer (ELB) and direct any traffic on that port, to your swarm service.

## How can I configure my load balancer to support SSL/TLS traffic?

Docker uses [Amazons' ACM service](https://aws.amazon.com/certificate-manager/), which provides free SSL/TLS certificates, and can be used with ELBs. You need to create a new certificate for your domain, and get the ARN for that certificate.

You add a label to your service to tell swarm that you want to use a given ACM cert for SSL connections to your service.

### Examples

Start a service and listen on the ELB with ports `80` and `443`. Port `443` is served using a SSL certificate from ACM, which is referenced by the ARN that is described in the service label `com.docker.aws.lb.arn`

```bash
$ docker service create \
  --name demo \
  --detach=true \
  --publish 80:80 \
  --publish 443:80 \
  --label com.docker.aws.lb.arn="arn:aws:acm:us-east-1:0123456789:certificate/c02117b6-2b5f-4507-8115-87726f4ab963" \
  yourname/your-image:latest
```

By default when you add an ACM ARN as a label, it listens on port `443`. If you want to change which port to listen too you append an `@` symbol and a list of ports you want to expose.

#### links SSL to port 443

```none
com.docker.aws.lb.arn="arn:..."
```

#### links SSL to port 444

```none
com.docker.aws.lb.arn="arn:...@444"
```

#### links SSL to ports 444 and 8080

```none
com.docker.aws.lb.arn="arn:...@444,8080"
```

### More full examples

Listen for HTTP on ports 80 and HTTPS on 444

```bash
$ docker service create \
  --name demo \
  --detach=true \
  --publish 80:80 \
  --publish 444:80 \
  --label com.docker.aws.lb.arn="arn:aws:acm:us-east-1:0123456789:certificate/c02117b6-2b5f-4507-8115-87726f4ab963@444" \
   yourname/your-image:latest
```

#### SSL listen on port 444 and 443

```bash
$ docker service create \
  --name demo \
  --detach=true \
  --publish 80:80 \
  --publish 444:80 \
  --label com.docker.aws.lb.arn="arn:aws:acm:us-east-1:0123456789:certificate/c02117b6-2b5f-4507-8115-87726f4ab963@443,444" \
   yourname/your-image:latest
```

#### SSL listen on port 8080

```bash
$ docker service create \
  --name demo \
  --detach=true \
  --publish 8080:80 \
  --label com.docker.aws.lb.arn="arn:aws:acm:us-east-1:0123456789:certificate/c02117b6-2b5f-4507-8115-87726f4ab963@8080" \
  yourname/your-image:latest
```

### Add a CNAME for your ELB

Once you have your ELB setup, with the correct listeners and certificates, you need to add a DNS CNAME that points to your ELB at your DNS provider.

### ELB SSL limitations

- There can only be one SSL certificate per ELB, so that means you can only have one label per swarm. If you add more than one `com.docker.aws.lb.arn` label per swarm the last one loaded overwrites the other one.
- If you remove the service that has the `com.docker.aws.lb.arn` label, it is removed from the ELB.
- If you edit the ELB config directly from the dashboard, the changes are removed after the next update.

## Can I manually change the ELB configuration?

No. If you make any manual changes to the ELB, they are removed the next time we update the ELB configuration based on any swarm changes. This is because the swarm service configuration is the source of record for service ports. If you add listeners to the ELB manually, they could conflict with what is in swarm, and cause issues.
