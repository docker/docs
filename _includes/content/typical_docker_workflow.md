### Typical Docker Platform Workflow

1. Get your code and its dependencies into Docker [containers](engine/getstarted/step_two.md):
   - [Write a Dockerfile](engine/getstarted/step_four.md) that specifies the execution
     environment and pulls in your code.
   - If your app depends on external applications (such as Redis, or
     MySQL), simply [find them on a registry such as Docker Hub](docker-hub/repos.md), and refer to them in
     [a Docker Compose file](compose/overview.md), along with a reference to your application, so they'll run
     simultaneously.
     - Software providers also distribute paid software via the [Docker Store](https://store.docker.com).
   - Build, then run your containers on a virtual host via [Docker Machine](machine/overview.md) as you develop.
2. Configure [networking](engine/tutorials/networkingcontainers.md) and
   [storage](engine/tutorials/dockervolumes.md) for your solution, if needed.
3. Upload builds to a registry ([ours](engine/tutorials/dockerrepos.md), [yours](/datacenter/dtr/2.0/index.md), or your cloud provider's), to collaborate with your team.
4. If you're gonna need to scale your solution across multiple hosts (VMs or physical machines), [plan
   for how you'll set up your Swarm cluster](engine/swarm/key-concepts.md) and [scale it to meet demand](engine/swarm/swarm-tutorial/index.md).
   - Note: Use [Universal Control Plane](/datacenter/ucp/1.1/overview.md) and you can manage your
     Swarm cluster using a friendly UI!
5. Finally, deploy to your preferred
   cloud provider (or, for redundancy, *multiple* cloud providers) with [Docker Cloud](/docker-cloud/index.md). Or, use [Docker Datacenter](https://www.docker.com/products/docker-datacenter), and deploy to your own on-premise hardware.
