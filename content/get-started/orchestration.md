---
title: Deployment and orchestration
keywords: orchestration, deploy, kubernetes, swarm,
description: Get oriented on some basics of Docker and install Docker Desktop.
---

Containerization provides an opportunity to move and scale applications to
clouds and data centers. Containers effectively guarantee that those applications run the
same way anywhere, allowing you to quickly and easily take advantage of all
these environments. Additionally, as you scale your applications up, you need some
tooling to help automate the maintenance of those applications, enable the
replacement of failed containers automatically, and manage the roll-out of
updates and reconfigurations of those containers during their lifecycle.

Tools to manage, scale, and maintain containerized applications are called
orchestrators. Two of the most popular orchestration tools are Kubernetes and
Docker Swarm. Docker Desktop provides development environments for both of these
orchestrators.

The advanced modules teach you how to:

1. [Set up and use a Kubernetes environment on your development machine](kube-deploy.md)
2. [Set up and use a Swarm environment on your development machine](swarm-deploy.md)

## Turn on Kubernetes

Docker Desktop sets up Kubernetes for you quickly and easily. Follow the setup and validation instructions appropriate for your operating system:

{{< tabs group="os" >}}
{{< tab name="Mac and Linux" >}}

### Mac

1.  From the Docker Dashboard, navigate to **Settings**, and select the **Kubernetes** tab.

2.  Select the checkbox labeled **Enable Kubernetes**, and select **Apply & Restart**. Docker Desktop automatically sets up Kubernetes for you. You'll know that Kubernetes has been successfully enabled when you see a green light beside 'Kubernetes _running_' in **Settings**.

3. To confirm that Kubernetes is up and running, create a text file called `pod.yaml` with the following content:

    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: demo
    spec:
      containers:
      - name: testpod
        image: alpine:latest
        command: ["ping", "8.8.8.8"]
    ```

    This describes a pod with a single container, isolating a simple ping to 8.8.8.8.

4. In a terminal, navigate to where you created `pod.yaml` and create your pod:

    ```console
    $ kubectl apply -f pod.yaml
    ```

5. Check that your pod is up and running:

    ```console
    $ kubectl get pods
    ```

    You should see something like:

    ```shell
    NAME      READY     STATUS    RESTARTS   AGE
    demo      1/1       Running   0          4s
    ```

6. Check that you get the logs you'd expect for a ping process:

    ```console
    $ kubectl logs demo
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

    ```console
    $ kubectl delete -f pod.yaml
    ```

{{< /tab >}}
{{< tab name="Windows" >}}

### Windows

1. From the Docker Dashboard, navigate to **Settings**, and select the **Kubernetes** tab.

2. Select the checkbox labeled **Enable Kubernetes**, and select **Apply & Restart**. Docker Desktop automatically sets up Kubernetes for you. You'll know that Kubernetes has been successfully enabled when you see a green light beside 'Kubernetes _running_' in the **Settings** menu.

3. To confirm that Kubernetes is up and running, create a text file called `pod.yaml` with the following content:

    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: demo
    spec:
      containers:
      - name: testpod
        image: alpine:latest
        command: ["ping", "8.8.8.8"]
    ```

    This describes a pod with a single container, isolating a simple ping to 8.8.8.8.

4. In PowerShell, navigate to where you created `pod.yaml` and create your pod:

    ```console
    $ kubectl apply -f pod.yaml
    ```

5. Check that your pod is up and running:

    ```console
    $ kubectl get pods
    ```

    You should see something like:

    ```shell
    NAME      READY     STATUS    RESTARTS   AGE
    demo      1/1       Running   0          4s
    ```

6. Check that you get the logs you'd expect for a ping process:

    ```console
    $ kubectl logs demo
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

    ```console
    $ kubectl delete -f pod.yaml
    ```

{{< /tab >}}
{{< /tabs >}}

## Enable Docker Swarm

Docker Desktop runs primarily on Docker Engine, which has everything you need to run a Swarm built in. Follow the setup and validation instructions appropriate for your operating system:

{{< tabs group="os" >}}
{{< tab name="Mac" >}}

### Mac

1. Open a terminal, and initialize Docker Swarm mode:

    ```console
    $ docker swarm init
    ```

    If all goes well, you should see a message similar to the following:

    ```shell
    Swarm initialized: current node (tjjggogqpnpj2phbfbz8jd5oq) is now a manager.

    To add a worker to this swarm, run the following command:

        docker swarm join --token SWMTKN-1-3e0hh0jd5t4yjg209f4g5qpowbsczfahv2dea9a1ay2l8787cf-2h4ly330d0j917ocvzw30j5x9 192.168.65.3:2377

    To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
    ```

2. Run a simple Docker service that uses an alpine-based filesystem, and isolates a ping to 8.8.8.8:

    ```console
    $ docker service create --name demo alpine:latest ping 8.8.8.8
    ```

3. Check that your service created one running container:

    ```console
    $ docker service ps demo
    ```

    You should see something like:

    ```shell
    ID                  NAME                IMAGE               NODE                DESIRED STATE       CURRENT STATE           ERROR               PORTS
    463j2s3y4b5o        demo.1              alpine:latest       docker-desktop      Running             Running 8 seconds ago
    ```

4. Check that you get the logs you'd expect for a ping process:

    ```console
    $ docker service logs demo
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

    ```console
    $ docker service rm demo
    ```

{{< /tab >}}
{{< tab name="Windows" >}}

### Windows

1. Open a PowerShell, and initialize Docker Swarm mode:

    ```console
    $ docker swarm init
    ```

    If all goes well, you should see a message similar to the following:

    ```shell
    Swarm initialized: current node (tjjggogqpnpj2phbfbz8jd5oq) is now a manager.

    To add a worker to this swarm, run the following command:

        docker swarm join --token SWMTKN-1-3e0hh0jd5t4yjg209f4g5qpowbsczfahv2dea9a1ay2l8787cf-2h4ly330d0j917ocvzw30j5x9 192.168.65.3:2377

    To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
    ```

2. Run a simple Docker service that uses an alpine-based filesystem, and isolates a ping to 8.8.8.8:

    ```console
    $ docker service create --name demo alpine:latest ping 8.8.8.8
    ```

3. Check that your service created one running container:

    ```console
    $ docker service ps demo
    ```

    You should see something like:

    ```shell
    ID                  NAME                IMAGE               NODE                DESIRED STATE       CURRENT STATE           ERROR               PORTS
    463j2s3y4b5o        demo.1              alpine:latest       docker-desktop      Running             Running 8 seconds ago
    ```

4. Check that you get the logs you'd expect for a ping process:

    ```console
    $ docker service logs demo
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

    ```console
    $ docker service rm demo
    ```

{{< /tab >}}
{{< /tabs >}}

## Conclusion

At this point, you've confirmed that you can run simple containerized workloads in Kubernetes and Swarm. The next step is to write a YAML file that describes how to run and manage these containers.

- [Deploy to Kubernetes](kube-deploy.md)
- [Deploy to Swarm](swarm-deploy.md)

## CLI references

Further documentation for all CLI commands used in this article are available here:

- [`kubectl apply`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#apply)
- [`kubectl get`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#get)
- [`kubectl logs`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#logs)
- [`kubectl delete`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#delete)
- [`docker swarm init`](/reference/cli/docker/swarm/init/)
- [`docker service *`](/reference/cli/docker/service/)