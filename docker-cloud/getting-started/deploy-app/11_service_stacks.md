---
description: Stackfiles for your service
keywords: Python, service, stack
redirect_from:
- /docker-cloud/getting-started/python/11_service_stacks/
title: Stackfiles for your service
---

## What are stack files?

A stack is a logical grouping of related services that are usually deployed
together and require each other to work as intended. If you are familiar with
*fig* or *Docker Compose* then you should feel right at home with **stacks**.
You can learn more about stacks [here](../../apps/stacks.md).

Stack files are YAML files, and you can learn more about the available syntax
[here](../../apps/stack-yaml-reference.md). You can also interact with stacks
using the [stack commands in our API](/apidocs/docker-cloud.md#stacks).

## Service definitions in the stack file

The services that you created in this tutorial form a stack with three services:
the load-balancer, the web application, and the redis cache.

Look at the file called `docker-cloud.yml` in your quickstart to see the stack
file that defines the three services (lb, web, redis) you created in the
previous steps, including all modifications and environment variables.

This is what the `docker-cloud.yml` file looks like. (If you are using the
quickstart-go version, you see `quickstart-go` instead of
`quickstart-python`.)

```yml
lb:
  image: dockercloud/haproxy
  autorestart: always
  links:
    - web
  ports:
    - "80:80"
  roles:
    - global
web:
  image: dockercloud/quickstart-python
  autorestart: always
  links:
    - redis
  environment:
    - NAME=Friendly Users
  deployment_strategy: high_availability
  target_num_containers: 4
redis:
  image: redis
  autorestart: always
  environment:
    - REDIS_PASS=password
    - REDIS_APPENDONLY=yes
    - REDIS_APPENDFSYNC=always
```

You can use this stack file to quickly deploy this cluster of three services to
another set of nodes. You can also edit the file to change the configuration.

## Run a stack

To create the services in a stack file you use the simple `stack up` command.

You can run this in the path containing your stackfile (docker-cloud.yml), like
so:

```bash
$ docker-cloud stack up
```

Or you can specify the YML file to use and its location:

```bash
$ docker-cloud up -f /usr/dockercloud/quickstart-python/docker-cloud.yml
```

## What's Next?

Next, we do some [data management with volumes](12_data_management_with_volumes.md).
