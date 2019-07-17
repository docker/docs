---
title: Deploy an app from the UI
description: Learn how to deploy containerized applications on a cluster, with Docker Universal Control Plane.
keywords: ucp, deploy, application, stack, service, compose
---

With Docker Universal Control Plane you can deploy applications from the UI
using `docker-compose.yml` files. In this example, we're going to deploy an
application that allows users to vote on whether they prefer cats or dogs. 😺 🐶

## Deploy the voting application

The application we're going to deploy is composed of several services:

* `vote`: The web application that presents the voting interface via port 5000
* `result`: A web application that displays the voting results via port 5001
* `visualizer`: A web application that shows a map of the deployment of the
  various services across the available nodes via port 8080
* `redis`: Collects raw voting data and stores it in a key/value queue
* `db`: A PostgreSQL service which provides permanent storage on a host volume
* `worker`: A background service that transfers votes from the queue to permanent storage

In your browser, log in to the UCP web UI, and navigate to the
**Stacks** page. Click **Create Stack** to deploy a new application.

Give the application a name, like "VotingApp", and in the **Mode** field,
select **Services**.

Paste the following YAML into the **COMPOSE.YML** editor:

```none
version: "3"
services:

  redis:
    image: redis:alpine
    ports:
      - "6379"
    networks:
      - frontend
    deploy:
      replicas: 2
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure
  db:
    image: postgres:9.4
    volumes:
      - db-data:/var/lib/postgresql/data
    networks:
      - backend
    deploy:
      placement:
        constraints: [node.role == manager]
  vote:
    image: manomarks/examplevotingapp_vote
    ports:
      - 5000:80
    networks:
      - frontend
    depends_on:
      - redis
    deploy:
      replicas: 6
      update_config:
        parallelism: 2
      restart_policy:
        condition: on-failure
  result:
    image: manomarks/examplevotingapp_result
    ports:
      - 5001:80
    networks:
      - backend
    deploy:
      replicas: 2
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure

  worker:
    image: manomarks/examplevotingapp_worker
    networks:
      - frontend
      - backend
    deploy:
      mode: replicated
      replicas: 2
      labels: [APP=VOTING]
      restart_policy:
        condition: on-failure
        delay: 10s
        max_attempts: 3
        window: 120s
      placement:
        constraints: [node.role == worker]
  visualizer:
    image: manomarks/visualizer
    ports:
      - "8080:8080"
    stop_grace_period: 1m30s
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock"

networks:
  frontend:
  backend:

volumes:
  db-data:
```

![](../../images/deploy-app-ui-1.png){: .with-border}

> When "Services" is selected, you can define services in this YAML file that
have a `deploy:` key, which schedules the containers on certain nodes, defines
their restart behavior, configures the number of replicas, and so on. These
features are provided by the Compose V3 file format.
[Learn about Compose files](/compose/compose-file/).

Click **Create** to build and deploy the application. When you see the
**Created successfully** message, click **Done**.

In the left pane, click **Services** to see the details of the services that
you deployed across your nodes. Click the `VotingApp_vote` service and find
the **Published Endpoint** field in the details pane. Click the link to visit
the voting page, which is published on port `5000`.

![Screenshot of deployed service](../../images/deployed_visualizer.png){: .with-border}

Click **Cats** and **Dogs** a few times to register some votes, and notice
that each click is processed by a different container.

Go back to the **Services** page in the UCP web UI. Click the
`VotingApp_result` service and find the **Published Endpoint** field in
the details pane. It has the same URL as the other `VotingApp` services,
but it's published on port `5001`. Click the link to view the vote tally.

Back in the **Services** page, click the
`VotingApp_visualizer` service and find the **Published Endpoint** field in
the details pane. A link to your UCP instance's URL includes
the published port of the visualizer service, which is 8080 in this case.

Visiting this URL accesses the running instance of the `VotingApp_visualizer`
service in your browser, which shows a map of how this application was deployed:

![Screenshot of visualizer](../../images/deployed_visualizer_detail.png){: .with-border}

You can see some of the characteristics of the deployment specification
from the Compose file in play. For example, the manager node is running the
PostgreSQL container, as configured by setting `[node.role == manager]` as a
constraint in the `deploy` key for the `db` service.

## Limitations

There are some limitations when deploying docker-compose.yml applications from
the UI. You can't reference any external files, so the following Docker
Compose keywords are not supported:

* build
* dockerfile
* env_file

To overcome these limitations, you can
[deploy your apps from the CLI](deploy-app-cli.md).

Also, UCP doesn't store the compose file used to deploy the application. You can
use your version control system to persist that file.

## Where to go next

* [Deploy an app from the CLI](deploy-app-cli.md)
