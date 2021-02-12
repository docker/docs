{% assign platform = include.platform %}

{% comment %}

Include a chunk of this file, using variables already set in the file
where you want to reuse the chunk.

Usage: {% include kubernetes-mac-win.md platform="mac" %}

{% endcomment %}

{% if platform == "mac" %}
  {% assign product = "Docker Desktop for Mac" %}

  {% capture min-version %}{{ product }} 18.06.0-ce-mac70 CE{% endcapture %}

  {% capture version-caveat %}
  **Kubernetes is only available in {{ min-version }} and higher.**
  {% endcapture %}

  {% capture local-kubectl-warning %}
> If you independently installed the Kubernetes CLI, `kubectl`, make sure that
> it is pointing to `docker-desktop` and not some other context such as
> `minikube` or a GKE cluster. Run: `kubectl config use-context docker-desktop`.
> If you experience conflicts with an existing `kubectl` installation, remove `/usr/local/bin/kubectl`.

  {% endcapture %}

  {% assign kubectl-path = "/usr/local/bin/kubectl" %}

{% elsif platform == "windows" %}
  {% assign product = "Docker Desktop for Windows" %}

  {% capture min-version %}{{ product }} 18.06.0-ce-win70 CE{% endcapture %}

  {% capture version-caveat %}
  **Kubernetes is only available in {{ min-version }} and higher.**
  {% endcapture %}

  {% capture local-kubectl-warning %}
  If you installed `kubectl` by another method, and experience conflicts, remove it.
  {% endcapture %}

  {% assign kubectl-path = "C:\>Program Files\Docker\Docker\Resources\bin\kubectl.exe" %}

{% endif %}

Docker Desktop includes a standalone Kubernetes server and client,
as well as Docker CLI integration. The Kubernetes server runs locally within
your Docker instance, is not configurable, and is a single-node cluster.

The Kubernetes server runs within a Docker container on your local system, and
is only for local testing. When Kubernetes support is enabled, you can deploy
your workloads, in parallel, on Kubernetes, Swarm, and as standalone containers.
Enabling or disabling the Kubernetes server does not affect your other
workloads.

See [{{ product }} > Getting started](/docker-for-{{ platform }}/#kubernetes) to
enable Kubernetes and begin testing the deployment of your workloads on
Kubernetes.

{{ kubectl-warning }}

## Use the kubectl command

The {{ platform }} Kubernetes integration provides the Kubernetes CLI command
at `{{ kubectl-path }}`. This location may not be in your shell's `PATH`
variable, so you may need to type the full path of the command or add it to
the `PATH`. For more information about `kubectl`, see the
[official `kubectl` documentation](https://kubernetes.io/docs/reference/kubectl/overview/).
You can test the command by listing the available nodes:

```bash
kubectl get nodes

NAME                 STATUS    ROLES     AGE       VERSION
docker-desktop       Ready     master    3h        v1.8.2
```
