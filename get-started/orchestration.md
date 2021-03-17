---
title: "Orchestration"
keywords: orchestration, deploy, swarm,
description: Get oriented on some basics of Docker and install Docker Desktop.
---

The portability and reproducibility of a containerized process mean we have an opportunity to move and scale our containerized applications across clouds and datacenters. Containers effectively guarantee that those applications run the same way anywhere, allowing us to quickly and easily take advantage of all these environments. Furthermore, as we scale our applications up, we'll want some tooling to help automate the maintenance of those applications, able to replace failed containers automatically, and manage the rollout of updates and reconfigurations of those containers during their lifecycle.

Tools to manage, scale, and maintain containerized applications are called _orchestrators_. The orchestration module teaches you how to [Set up and use a Swarm environment on your development machine](swarm-deploy.md).

## Enable Docker Swarm

Docker Desktop runs primarily on Docker Engine, which has everything you need to run a Swarm built in. Follow the setup and validation instructions appropriate for your operating system:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" href="#swarmosx">Mac</a></li>
  <li><a data-toggle="tab" href="#swarmwin">Windows</a></li>
</ul>
<div class="tab-content">
  <div id="swarmosx" class="tab-pane fade in active">
{% capture local-content %}

### Mac

1.  Open a terminal, and initialize Docker Swarm mode:

    ```shell
    docker swarm init
    ```

    If all goes well, you should see a message similar to the following:

    ```shell
    Swarm initialized: current node (tjjggogqpnpj2phbfbz8jd5oq) is now a manager.

    To add a worker to this swarm, run the following command:

        docker swarm join --token SWMTKN-1-3e0hh0jd5t4yjg209f4g5qpowbsczfahv2dea9a1ay2l8787cf-2h4ly330d0j917ocvzw30j5x9 192.168.65.3:2377

    To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
    ```

2.  Run a simple Docker service that uses an alpine-based filesystem, and isolates a ping to 8.8.8.8:

    ```shell
    docker service create --name demo alpine:3.5 ping 8.8.8.8
    ```

3.  Check that your service created one running container:

    ```shell
    docker service ps demo
    ```

    You should see something like:

    ```shell
    ID                  NAME                IMAGE               NODE                DESIRED STATE       CURRENT STATE           ERROR               PORTS
    463j2s3y4b5o        demo.1              alpine:3.5          docker-desktop      Running             Running 8 seconds ago
    ```

4.  Check that you get the logs you'd expect for a ping process:

    ```shell
    docker service logs demo
    ```

    You should see the output of a healthy ping process:

    ```shell
    demo.1.463j2s3y4b5o@docker-desktop    | PING 8.8.8.8 (8.8.8.8): 56 data bytes
    demo.1.463j2s3y4b5o@docker-desktop    | 64 bytes from 8.8.8.8: seq=0 ttl=37 time=13.005 ms
    demo.1.463j2s3y4b5o@docker-desktop    | 64 bytes from 8.8.8.8: seq=1 ttl=37 time=13.847 ms
    demo.1.463j2s3y4b5o@docker-desktop    | 64 bytes from 8.8.8.8: seq=2 ttl=37 time=41.296 ms
    ...
    ```

5.  Finally, tear down your test service:

    ```shell
    docker service rm demo
    ```

{% endcapture %}
{{ local-content | markdownify }}

</div>
<div id="swarmwin" class="tab-pane fade" markdown="1">
{% capture localwin-content %}

### Windows

1.  Open a powershell, and initialize Docker Swarm mode:

    ```shell
    docker swarm init
    ```

    If all goes well, you should see a message similar to the following:

    ```shell
    Swarm initialized: current node (tjjggogqpnpj2phbfbz8jd5oq) is now a manager.

    To add a worker to this swarm, run the following command:

        docker swarm join --token SWMTKN-1-3e0hh0jd5t4yjg209f4g5qpowbsczfahv2dea9a1ay2l8787cf-2h4ly330d0j917ocvzw30j5x9 192.168.65.3:2377

    To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
    ```

2.  Run a simple Docker service that uses an alpine-based filesystem, and isolates a ping to 8.8.8.8:

    ```shell
    docker service create --name demo alpine:3.5 ping 8.8.8.8
    ```

3.  Check that your service created one running container:

    ```shell
    docker service ps demo
    ```

    You should see something like:

    ```shell
    ID                  NAME                IMAGE               NODE                DESIRED STATE       CURRENT STATE           ERROR               PORTS
    463j2s3y4b5o        demo.1              alpine:3.5          docker-desktop      Running             Running 8 seconds ago
    ```

4.  Check that you get the logs you'd expect for a ping process:

    ```shell
    docker service logs demo
    ```

    You should see the output of a healthy ping process:

    ```shell
    demo.1.463j2s3y4b5o@docker-desktop    | PING 8.8.8.8 (8.8.8.8): 56 data bytes
    demo.1.463j2s3y4b5o@docker-desktop    | 64 bytes from 8.8.8.8: seq=0 ttl=37 time=13.005 ms
    demo.1.463j2s3y4b5o@docker-desktop    | 64 bytes from 8.8.8.8: seq=1 ttl=37 time=13.847 ms
    demo.1.463j2s3y4b5o@docker-desktop    | 64 bytes from 8.8.8.8: seq=2 ttl=37 time=41.296 ms
    ...
    ```

5.  Finally, tear down your test service:

    ```shell
    docker service rm demo
    ```

{% endcapture %}
{{ localwin-content | markdownify }}
</div>
<hr>
</div>

## Conclusion

At this point, you've confirmed that you can run simple containerized workloads in Swarm. To learn how to write the stack file to help you run and manage containers on Swarm, see [Deploying to Swarm](swarm-deploy.md).

## CLI references

Further documentation for all CLI commands used in this article are available here:

- [`docker swarm init`](https://docs.docker.com/engine/reference/commandline/swarm_init/)
- [`docker service *`](https://docs.docker.com/engine/reference/commandline/service/)
