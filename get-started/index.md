---
title: "Get Started, Part 1: Orientation and setup"
keywords: get started, setup, orientation, quickstart, intro, concepts, containers, docker desktop
description: Get oriented on some basics of Docker and install Docker Desktop.
redirect_from:
- /getstarted/
- /get-started/part1/
- /get-started/part6/
- /engine/getstarted/
- /learn/
- /engine/getstarted/step_one/
- /engine/getstarted/step_two/
- /engine/getstarted/step_three/
- /engine/getstarted/step_four/
- /engine/getstarted/step_five/
- /engine/getstarted/step_six/
- /engine/getstarted/last_page/
- /engine/getstarted-voting-app/
- /engine/getstarted-voting-app/node-setup/
- /engine/getstarted-voting-app/create-swarm/
- /engine/getstarted-voting-app/deploy-app/
- /engine/getstarted-voting-app/test-drive/
- /engine/getstarted-voting-app/customize-app/
- /engine/getstarted-voting-app/cleanup/
- /engine/userguide/intro/
- /mac/started/
- /windows/started/
- /linux/started/
- /getting-started/
- /mac/step_one/
- /windows/step_one/
- /linux/step_one/
- /engine/tutorials/dockerizing/
- /mac/step_two/
- /windows/step_two/
- /linux/step_two/
- /mac/step_three/
- /windows/step_three/
- /linux/step_three/
- /engine/tutorials/usingdocker/
- /mac/step_four/
- /windows/step_four/
- /linux/step_four/
- /engine/tutorials/dockerimages/
- /userguide/dockerimages/
- /engine/userguide/dockerimages/
- /mac/last_page/
- /windows/last_page/
- /linux/last_page/
- /mac/step_six/
- /windows/step_six/
- /linux/step_six/
- /engine/tutorials/dockerrepos/
- /userguide/dockerrepos/
- /engine/userguide/containers/dockerimages/
---

{% include_relative nav.html selected="1" %}

Welcome! We are excited that you want to learn Docker. The _Docker Get Started Tutorial_
teaches you how to:

1. Set up your Docker environment (on this page)
2. [Build an image and run it as one container](part2.md)
3. [Set up and use a Kubernetes environment on your development machine](part3.md)
4. [Set up and use a Swarm environment on your development machine](part4.md)
5. [Share your containerized applications on Docker Hub](part5.md)

## Docker concepts

Docker is a platform for developers and sysadmins to **build, share, and run**
applications with containers. The use of containers to deploy applications
is called _containerization_. Containers are not new, but their use for easily
deploying applications is.

Containerization is increasingly popular because containers are:

- Flexible: Even the most complex applications can be containerized.
- Lightweight: Containers leverage and share the host kernel,
  making them much more efficient in terms of system resources than virtual machines.
- Portable: You can build locally, deploy to the cloud, and run anywhere.
- Loosely coupled: Containers are highly self sufficient and encapsulated,
  allowing you to replace or upgrade one without disrupting others.
- Scalable: You can increase and automatically distribute container replicas across a datacenter.
- Secure: Containers apply aggressive constraints and isolations to processes without
  any configuration required on the part of the user.

![Containers are portable](images/laurel-docker-containers2019.png){:width="100%"}

### Images and containers

Fundamentally, a container is nothing but a running process,
with some added encapsulation features applied to it in order to keep it isolated from the host
and from other containers.
One of the most important aspects of container isolation is that each container interacts
with its own, private filesystem; this filesystem is provided by a Docker **image**.
An image includes everything needed to run an application -- the code or binary,
runtimes, dependencies, and any other filesystem objects required.

### Containers and virtual machines

A container runs _natively_ on Linux and shares the kernel of the host
machine with other containers. It runs a discrete process, taking no more memory
than any other executable, making it lightweight.

By contrast, a **virtual machine** (VM) runs a full-blown "guest" operating
system with _virtual_ access to host resources through a hypervisor. In general,
VMs incur a lot of overhead beyond what is being consumed by your application logic.

![Container stack example](/images/Container%402x.png){:width="300px"} | ![Virtual machine stack example](/images/VM%402x.png){:width="300px"}

## Install Docker Desktop

The best way to get started developing containerized applications is with Docker Desktop, for OSX or Windows. Docker Desktop will allow you to easily set up Kubernetes or Swarm on your local development machine, so you can use all the features of the orchestrator you're developing applications for right away, no cluster required. Follow the installation instructions appropriate for your operating system:

 - [OSX](/docker-for-mac/install/){: target="_blank" class="_"}
 - [Windows](/docker-for-windows/install/){: target="_blank" class="_"}

## Enable Kubernetes

Docker Desktop will set up Kubernetes for you quickly and easily. Follow the setup and validation instructions appropriate for your operating system:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" href="#kubeosx">OSX</a></li>
  <li><a data-toggle="tab" href="#kubewin">Windows</a></li>
</ul>
<div class="tab-content">
  <div id="kubeosx" class="tab-pane fade in active">
{% capture local-content %}

#### OSX

1. After installing Docker Desktop, you should see a Docker icon in your menu bar. Click on it, and navigate **Preferences... -> Kubernetes**.

2. Check the checkbox labeled *Enable Kubernetes*, and click **Apply**. Docker Desktop will automatically set up Kubernetes for you. You'll know everything has completed successfully once you can click on the Docker icon in the menu bar, and see a green light beside 'Kubernetes is Running'.

3. In order to confirm that Kubernetes is up and running, create a text file called `pod.yaml` with the following content:

    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: demo
    spec:
      containers:
      - name: testpod
        image: alpine:3.5
        command: ["ping", "8.8.8.8"]
    ```

    This describes a pod with a single container, isolating a simple ping to 8.8.8.8.

4. In a terminal, navigate to where you created `pod.yaml` and create your pod:

    ```shell
    kubectl apply -f pod.yaml
    ```

5. Check that your pod is up and running:

    ```shell
    kubectl get pods
    ```

    You should see something like:

    ```shell
    NAME      READY     STATUS    RESTARTS   AGE
    demo      1/1       Running   0          4s
    ```

6. Check that you get the logs you'd expect for a ping process:

    ```shell
    kubectl logs demo
    ```

    You should see the output of a healthy ping process:

    ```shell
    PING 8.8.8.8 (8.8.8.8): 56 data bytes
    64 bytes from 8.8.8.8: seq=0 ttl=37 time=21.393 ms
    64 bytes from 8.8.8.8: seq=1 ttl=37 time=15.320 ms
    64 bytes from 8.8.8.8: seq=2 ttl=37 time=11.111 ms
    ...
    ```

7. Finally, tear down your test pod:

    ```shell
    kubectl delete -f pod.yaml
    ```

{% endcapture %}
{{ local-content | markdownify }}

</div>
<div id="kubewin" class="tab-pane fade" markdown="1">
{% capture localwin-content %}

#### Windows

1. After installing Docker Desktop, you should see a Docker icon in your system tray. Right-click on it, and navigate **Settings -> Kubernetes**.

2. Check the checkbox labeled *Enable Kubernetes*, and click **Apply**. Docker Desktop will automatically set up Kubernetes for you. Note this can take a significant amount of time (20 minutes). You'll know everything has completed successfully once you can right-click on the Docker icon in the menu bar, click **Settings**, and see a green light beside 'Kubernetes is running'.

3. In order to confirm that Kubernetes is up and running, create a text file called `pod.yaml` with the following content:

    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: demo
    spec:
      containers:
      - name: testpod
        image: alpine:3.5
        command: ["ping", "8.8.8.8"]
    ```

    This describes a pod with a single container, isolating a simple ping to 8.8.8.8.

4. In powershell, navigate to where you created `pod.yaml` and create your pod:

    ```shell
    kubectl apply -f pod.yaml
    ```

5. Check that your pod is up and running:

    ```shell
    kubectl get pods
    ```

    You should see something like:

    ```shell
    NAME      READY     STATUS    RESTARTS   AGE
    demo      1/1       Running   0          4s
    ```

6. Check that you get the logs you'd expect for a ping process:

    ```shell
    kubectl logs demo
    ```

    You should see the output of a healthy ping process:

    ```shell
    PING 8.8.8.8 (8.8.8.8): 56 data bytes
    64 bytes from 8.8.8.8: seq=0 ttl=37 time=21.393 ms
    64 bytes from 8.8.8.8: seq=1 ttl=37 time=15.320 ms
    64 bytes from 8.8.8.8: seq=2 ttl=37 time=11.111 ms
    ...
    ```

7. Finally, tear down your test pod:

    ```shell
    kubectl delete -f pod.yaml
    ```

{% endcapture %}
{{ localwin-content | markdownify }}
</div>
<hr>
</div>

## Enable Docker Swarm

Docker Desktop runs primarily on Docker Engine, which has everything you need to run a Swarm built in. Follow the setup and validation instructions appropriate for your operating system:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" href="#swarmosx">OSX</a></li>
  <li><a data-toggle="tab" href="#swarmwin">Windows</a></li>
</ul>
<div class="tab-content">
  <div id="swarmosx" class="tab-pane fade in active">
{% capture local-content %}

#### OSX

1. Open a terminal, and initialize Docker Swarm mode:

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

2. Run a simple Docker service that uses an alpine-based filesystem, and isolates a ping to 8.8.8.8:

    ```shell
    docker service create --name demo alpine:3.5 ping 8.8.8.8
    ```

3. Check that your service created one running container:

    ```shell
    docker service ps demo
    ```

    You should see something like:

    ```shell
    ID                  NAME                IMAGE               NODE                DESIRED STATE       CURRENT STATE           ERROR               PORTS
    463j2s3y4b5o        demo.1              alpine:3.5          docker-desktop      Running             Running 8 seconds ago
    ```

4. Check that you get the logs you'd expect for a ping process:

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

5. Finally, tear down your test service:

    ```shell
    docker service rm demo
    ```

{% endcapture %}
{{ local-content | markdownify }}

</div>
<div id="swarmwin" class="tab-pane fade" markdown="1">
{% capture localwin-content %}

#### Windows

1. Open a powershell, and initialize Docker Swarm mode:

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

2. Run a simple Docker service that uses an alpine-based filesystem, and isolates a ping to 8.8.8.8:

    ```shell
    docker service create --name demo alpine:3.5 ping 8.8.8.8
    ```

3. Check that your service created one running container:

    ```shell
    docker service ps demo
    ```

    You should see something like:

    ```shell
    ID                  NAME                IMAGE               NODE                DESIRED STATE       CURRENT STATE           ERROR               PORTS
    463j2s3y4b5o        demo.1              alpine:3.5          docker-desktop      Running             Running 8 seconds ago
    ```

4. Check that you get the logs you'd expect for a ping process:

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

5. Finally, tear down your test service:

    ```shell
    docker service rm demo
    ```

{% endcapture %}
{{ localwin-content | markdownify }}
</div>
<hr>
</div>

## Conclusion

At this point, you've installed Docker Desktop on your development machine, and confirmed that you can run simple containerized workloads in Kuberentes and Swarm. In the next section, we'll start developing our first containerized application.

[On to Part 2 >>](part2.md){: class="button outline-btn" style="margin-bottom: 30px; margin-right: 100%"}

## CLI References

Further documentation for all CLI commands used in this article are available here:

- [`kubectl apply`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#apply)
- [`kubectl get`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#get)
- [`kubectl logs`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#logs)
- [`kubectl delete`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#delete)
- [`docker swarm init`](https://docs.docker.com/engine/reference/commandline/swarm_init/)
- [`docker service *`](https://docs.docker.com/engine/reference/commandline/service/)
