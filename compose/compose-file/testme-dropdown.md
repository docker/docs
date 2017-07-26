---
description: Compose file reference
keywords: fig, composition, compose, docker
redirect_from:
- /compose/yml
- /compose/compose-file-v3.md
title: TEST ME
toc_max: 4
toc_min: 1
---



## Here is the test

Click the whale (![whale menu](/docker-for-mac/images/whale-x.png){: .inline}) to get Preferences and other options.

This is a test.{: .tryme}

<div class="tryme" id="mydiv">This is a paragraph in a div.</div>

Click to show/hide the example Compose file below.

<div class="panel-group" id="accordion" role="tablist" aria-multiselectable="true">

  <div class="panel panel-default">
    <div class="panel-heading" role="tab" id="headingThree">
      <h5 class="panel-title" id="collapsible-group-item-3"> <a class="" role="button" data-toggle="collapse" data-parent="#accordion" data-target="#collapseThree" aria-expanded="true" aria-controls="collapseThree"> Example Compose file version 3 <i class="fa fa-car"></i></a> </h5>
    </div>
    <div id="collapseThree" class="panel-collapse collapse" role="tabpanel" aria-labelledby="headingThree" aria-expanded="true">
      <div class="panel-body">
      <pre><code>
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
    image: dockersamples/examplevotingapp_vote:before
    ports:
      - 5000:80
    networks:
      - frontend
    depends_on:
      - redis
    deploy:
      replicas: 2
      update_config:
        parallelism: 2
      restart_policy:
        condition: on-failure
  result:
    image: dockersamples/examplevotingapp_result:before
    ports:
      - 5001:80
    networks:
      - backend
    depends_on:
      - db
    deploy:
      replicas: 1
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure

  worker:
    image: dockersamples/examplevotingapp_worker
    networks:
      - frontend
      - backend
    deploy:
      mode: replicated
      replicas: 1
      labels: [APP=VOTING]
      restart_policy:
        condition: on-failure
        delay: 10s
        max_attempts: 3
        window: 120s
      placement:
        constraints: [node.role == manager]

  visualizer:
    image: dockersamples/visualizer:stable
    ports:
      - "8080:8080"
    stop_grace_period: 1m30s
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock"
    deploy:
      placement:
        constraints: [node.role == manager]

networks:
    frontend:
    backend:

volumes:
    db-data:
    </code></pre>
      </div>
    </div>
  </div>
</div>

## TBD

Page continues here ..
