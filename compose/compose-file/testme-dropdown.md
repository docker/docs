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

<div class="panel panel-default">
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample1">
    Example Compose file version 3
    <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample1">
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
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample2"> Another Sample <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample2">
<p>
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.</p>
  </div>
</div>

## Here is another test

<div class="panel panel-default">
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample3"> First Cool Sample
    <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample3">
<p>
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.
</p>
    </div>
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample4"> Second Cool Sample <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapseSample4">
<p>
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.
</p>
  </div>
</div>

## TBD

Page continues here ..
